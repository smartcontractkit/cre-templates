import { DynamoDBDocumentClient } from "@aws-sdk/lib-dynamodb";
import { DynamoDBClient } from "@aws-sdk/client-dynamodb";
import { GetCommand, PutCommand } from "@aws-sdk/lib-dynamodb";
import * as crypto from "crypto";
import https from "https";
import { URL } from "url";

/**
 * Constants for HTTP status codes used throughout the handler.
 */
const STATUS_OK = 200;
const STATUS_BAD_REQUEST = 400;
const STATUS_NOT_FOUND = 404;
const STATUS_SERVER_ERROR = 500;

/**
 * Name of the DynamoDB table to store asset states.
 */
const TABLE_NAME = "AssetState"; // Update the value if you use a table with different name

/**
 * Required fields for each action to validate incoming parameters.
 */
const REQUIRED_FIELDS = {
  read: ["assetId"],
  AssetRegistered: ["assetId", "issuer", "initialSupply", "assetName"],
  AssetVerified: ["assetId", "isValid"],
  TokensMinted: ["assetId", "amount"],
  TokensRedeemed: ["assetId", "amount"],
  sendNotification: ["assetId", "apiUrl"],
};

/**
 * Helper function to build a standardized Lambda response object.
 * @param {number} statusCode - The HTTP status code.
 * @param {object} body - The response body object.
 * @returns {object} Lambda response with statusCode and JSON-stringified body.
 */
const buildResponse = (statusCode, body) => ({
  statusCode,
  body: JSON.stringify(body),
});

/**
 * Validates that all required fields for the given action are present in params.
 * @param {string} action - The action being performed.
 * @param {object} params - The parsed request parameters.
 * @throws {Error} If validation fails.
 */
const validateParams = (action, params) => {
  const required = REQUIRED_FIELDS[action];
  if (!required || !required.every((field) => params[field] != null)) {
    throw new Error("Missing required parameters");
  }
};

/**
 * Retrieves an item from DynamoDB by assetId.
 * @param {object} client - The DynamoDB document client.
 * @param {string|number} assetId - The ID of the asset.
 * @returns {object} The item from DynamoDB, or empty object if not found.
 */
const getItem = async (client, assetId) => {
  const command = new GetCommand({
    TableName: TABLE_NAME,
    Key: { AssetId: assetId },
  });
  const { Item } = await client.send(command);
  return Item || {};
};

/**
 * Puts an item into DynamoDB.
 * @param {object} client - The DynamoDB document client.
 * @param {object} item - The item to store.
 */
const putItem = async (client, item) => {
  const command = new PutCommand({
    TableName: TABLE_NAME,
    Item: item,
  });
  await client.send(command);
};


/**
 * Handlers for each supported action. Each returns a promise resolving to the response data.
 */
const handlers = {
  /**
   * Reads the asset state from DynamoDB.
   */
  read: async (client, { assetId }) => ({ data: await getItem(client, assetId) }),

  /**
   * Registers a new asset in DynamoDB.
   */
  AssetRegistered: async (client, { assetId, issuer, initialSupply, assetName }) => {
    const item = {
      AssetId: assetId,
      AssetName: assetName,
      Issuer: issuer,
      Supply: initialSupply,
      Uid: crypto.randomUUID(),
    };
    await putItem(client, item);
    return { message: "Asset registered successfully" };
  },

  /**
   * Verifies an asset by updating its Verified status.
   */
  AssetVerified: async (client, { assetId, isValid }) => {
    const current = await getItem(client, assetId);
    const updated = { ...current, Verified: isValid };
    await putItem(client, updated);
    return { message: "Asset verified successfully", isValid };
  },

  /**
   * Mints new tokens by incrementing the TokenMinted count.
   */
  TokensMinted: async (client, { assetId, amount }) => {
    const current = await getItem(client, assetId);
    const currentAmount = BigInt(current.TokenMinted || 0n);
    const updated = {
      ...current,
      TokenMinted: (currentAmount + BigInt(amount)).toString(),
    };
    await putItem(client, updated);
    return { message: "New Token minted successfully", amount };
  },

  /**
   * Redeems tokens by incrementing the TokenRedeemed count.
   */
  TokensRedeemed: async (client, { assetId, amount }) => {
    const current = await getItem(client, assetId);
    const currentAmount = BigInt(current.TokenRedeemed || 0n);
    const updated = {
      ...current,
      TokenRedeemed: (currentAmount + BigInt(amount)).toString(),
    };
    await putItem(client, updated);
    return { message: "Token redeemed successfully", amount };
  },

  /**
   * Sends a POST notification to the provided API URL using asset data.
   * Note: 
   * Requires CRE workflow deployed on mainnet for full functionality.
   * In the demo, the action will not be used. 
   * the purpose of the snippet codes is to show how to send a POST request to CRE.
   */
  sendNotification: async (client, { assetId, apiUrl }) => {
    const item = await getItem(client, assetId);
    if (!item?.Uid) {
      throw new Error("Asset UID not found");
    }

    const postData = JSON.stringify({
      assetId: Number(assetId),
      uid: item.Uid,
    });

    const parsedUrl = new URL(apiUrl);
    const options = {
      hostname: parsedUrl.hostname,
      port: parsedUrl.port || 443,
      path: parsedUrl.pathname + parsedUrl.search,
      method: "POST",
      headers: {
        "Content-Type": "application/json",
        "Content-Length": Buffer.byteLength(postData),
      },
    };

    const response = await new Promise((resolve, reject) => {
      const req = https.request(options, (res) => {
        let data = "";
        res.on("data", (chunk) => (data += chunk));
        res.on("end", () => resolve({ statusCode: res.statusCode, body: data }));
      });
      req.on("error", reject);
      req.write(postData);
      req.end();
    });

    if (response.statusCode >= 200 && response.statusCode < 300) {
      return {
        message: "POST request sent successfully",
        assetId: Number(assetId),
        uid: item.Uid,
        apiResponse: response.body,
      };
    } else {
      throw new Error(`POST request failed: ${response.body}`);
    }
  },
};

/**
 * Main Lambda handler function.
 * Processes incoming events, validates, executes actions, and returns responses.
 * @param {object} event - The Lambda event object.
 * @returns {Promise<object>} The Lambda response.
 */
export const handler = async (event) => {
  const yourAwsRegion = ""; // input your AWS region here 

  if (!yourAwsRegion) {
    return buildResponse(STATUS_BAD_REQUEST, { error: "region is null, please define the region." });
  }
  
  // Initialize DynamoDB client with the specified region.
  const client = new DynamoDBDocumentClient(new DynamoDBClient({ region: yourAwsRegion }));

  // Parse request body as JSON.
  let params;
  try {
    params = JSON.parse(event.body || "{}");
  } catch {
    return buildResponse(STATUS_BAD_REQUEST, { error: "Invalid JSON in request body" });
  }

  const { action } = params;

  try {
    // Validate action exists and is supported.
    if (!action || !handlers[action]) {
      return buildResponse(STATUS_BAD_REQUEST, { error: "Invalid action" });
    }

    // Validate required parameters for the action.
    validateParams(action, params);
    
    // Execute the specific action handler.
    const result = await handlers[action](client, params);
    return buildResponse(STATUS_OK, result);
  } catch (error) {
    console.error("Error:", error);
    
    // Handle specific errors with appropriate status codes.
    if (error.message === "Asset UID not found") {
      return buildResponse(STATUS_NOT_FOUND, { error: error.message });
    }
    if (error.message.startsWith("POST request failed")) {
      return buildResponse(STATUS_BAD_REQUEST, { error: error.message });
    }
    return buildResponse(STATUS_SERVER_ERROR, { error: "Internal server error", details: error.message });
  }
};