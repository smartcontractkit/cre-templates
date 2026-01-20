// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import "@openzeppelin/contracts/access/AccessControl.sol";
import "@openzeppelin/contracts/token/ERC1155/extensions/ERC1155Burnable.sol";
import "@openzeppelin/contracts/utils/Strings.sol";
import {ReceiverTemplate} from "./interfaces/ReceiverTemplate.sol";

/**
 * @title TokenizedAssetPlatform
 * @dev A platform for real word assets(RWA).
 *      Supported assets includes invoice, loans, T-Bills, Carbon credit, etc., and the asset is control by issuer address
 *      ERC-1155 standard is used to manage multiple assets types, each assetId represents a specific asset issued
 *      life cycle of an asset includes:
 *      - Register: the admin registers a new asset for an issuer
 *      - Verify: issuer verify the asset and mark the asset active
 *      - Mint: issuer mints the asset for investors
 *      - Redeem/Burn: burn the tokens to redeem or settle
 *      - Transfer(optional): token holder transfer the token to recipients.
 *      - UpdateMetadata(optional): update the metadata of the asset
 *     
 *      Events will be emitted when key transactions happen for external retrieve and records
 *      Only key events will be used in CRE to update the offchain DB
 *      Asset issuers manage the asset with role ISSURE_ROLE
 *      The admin(there is only one admin for the platform) manage all assets with role DEFAULT_ADMIN_ROLE
 */
contract TokenizedAssetPlatform is ERC1155, AccessControl, ERC1155Burnable, ReceiverTemplate {
    using Strings for uint256;

    bytes32 public constant ISSUER_ROLE = keccak256("ISSUER_ROLE");
    bytes32 public constant ADMIN_ROLE = keccak256("ADMIN_ROLE");

    struct Asset {
        string name;          // Asset name(eg. Invoice-01, T-Bill-01, Loan-01, Carbon-Credit-01, etc.)
        string symbol;        // Asset symbol(eg. IVC1, TB1, LN1, CC1, etc)
        string assetType;     // Asset type(eg. receipt, loan, t-bill, carbon_credit)
        address issuer;       // Asset issuer address
        uint256 totalSupply;  // Asset total supply
        bool active;          // Mark if asset is verified or not by issuer
        string uid;           // uid in db, uid is only updated from updateMetadata called by CRE (optional)
    }

    struct UpdateMetadataCallData {
        uint256 assetId;
        string newUri;
    }

    // Map asset id to Asset, asset id is a auto-increment number
    mapping(uint256 => Asset) public assets;
    // Map asset id to issuer (the mapping is used to verify an address's permission)
    mapping(uint256 => address) public assetIssuers;
    // auto-increment asset ID
    uint256 private _nextAssetId = 1;

    // key events for asset life cycle
    event AssetRegistered(
        uint256 indexed assetId,
        address indexed issuer,
        uint256 indexed initialSupply,
        string name,
        string symbol,
        string assetType
    );

    event AssetVerified(
        uint256 indexed assetId, 
        bool indexed isValid, 
        string verificationDetails
    );
    
    event TokensMinted(
        uint256 indexed assetId,
        uint256 indexed amount,
        address indexed to,
        string reason
    );

    
    event TokensRedeemed(
        uint256 indexed assetId,
        uint256 indexed amount,
        address indexed account,
        string settlementDetails
    );

    // optional asset events for
    event TokensTransferred(uint256 indexed assetId, address indexed from, address indexed to, uint256 amount);
    event AssetPaused(uint256 indexed assetId, address indexed issuer);
    event AssetUnpaused(uint256 indexed assetId, address indexed issuer);
    event AssetUidUpdated(uint256 indexed assetId, string newUri);

    /**
     * @dev constructor, platform admin is set.
     */
    constructor(address forwardAddr) ERC1155("") ReceiverTemplate(forwardAddr){
        _grantRole(DEFAULT_ADMIN_ROLE, msg.sender);
        _grantRole(ADMIN_ROLE, msg.sender);
    }

    /**
   * @notice This internal function contains the core business logic.
   */
  function _processReport(bytes calldata report) internal override {
    // Decode the report bytes into assetId and UID
    (uint64 assetId, string memory _uid) = abi.decode(
            report,
            (uint64, string)
        );

    require(assetId != 0, "invalid asset Id");
    
    updateSystemOfRecordMetadata(assetId, _uid);
  }

    /**
     * @dev register new asset. The func can only be called by ADMIN_ROLE, ISSUER_ROLE is granted to a issuer.
     * @param name asset name
     * @param symbol asset symbol
     * @param assetType asset type
     * @param initialSupply initial total supply for the asset
     * @param assetUid asset db uid
     * @param issuer issuer address
     */
    function registerAsset(
        string memory name,
        string memory symbol,
        string memory assetType,
        uint256 initialSupply,
        string memory assetUid,
        address issuer
    ) public onlyRole(ADMIN_ROLE) {
        require(issuer != address(0), "Invalid issuer address");
        uint256 assetId = _nextAssetId++;
        
        assets[assetId] = Asset({
            name: name,
            symbol: symbol,
            assetType: assetType,
            issuer: issuer,
            totalSupply: initialSupply,
            active: false,
            uid: assetUid
        });
        
        assetIssuers[assetId] = issuer;
        
        // grant ISSUER_ROLE to a issuer who manages the asset
        _grantRole(ISSUER_ROLE, issuer);
        
        emit AssetRegistered(assetId, issuer, initialSupply, name, symbol, assetType);
    }

    /**
     * @dev the func can only be called by the issuer or amin to verify the validity of asset.
     * @param assetId asset ID
     * @param isValid the asset is valid or not
     * @param verificationDetails any detail for the asset verification (used in the event)
     */
    function verifyAsset(uint256 assetId, bool isValid, string memory verificationDetails) public {
        require(_isIssuerOrAdmin(assetId, msg.sender), "Not authorized to verify");
        assets[assetId].active = isValid;
        emit AssetVerified(assetId, isValid, verificationDetails);
    }

    /**
     * @dev mint token for an asset. The func can only be called by issuer and admin when the asset is active
     * @param to recipient address
     * @param assetId asset id
     * @param amount amount to mint
     * @param reason reason to mint (used in event)
     */
    function mint(address to, uint256 assetId, uint256 amount, string memory reason)
        public
    {
        require(_isIssuerOrAdmin(assetId, msg.sender), "Not authorized to mint");
        require(assets[assetId].active, "Asset is not active");
        require(to != address(0), "Invalid recipient");
        
        assets[assetId].totalSupply += amount;
        _mint(to, assetId, amount, "");
        
        emit TokensMinted(assetId, amount, to, reason);
    }

    /**
     * @dev redeem token, the func can only be called by asset token holders.
     * @param assetId asset ID
     * @param amount amount of token to redeem
     * @param settlementDetails settlement details(off-chain records)
     */
    function redeem(uint256 assetId, uint256 amount, string memory settlementDetails)
        public
    {
        require(assets[assetId].active, "Asset is not active");
        // holders can redeem tokens.
        if (!_isIssuerOrAdmin(assetId, msg.sender)) {
            require(balanceOf(msg.sender, assetId) >= amount, "Insufficient balance");
        }
        
        assets[assetId].totalSupply -= amount;
        _burn(msg.sender, assetId, amount);
        
        emit TokensRedeemed(assetId, amount, msg.sender, settlementDetails);
    }

    /**
     * @dev override burn to avoid assets burned without calling redeem
     */
    function burn(address /*account*/, uint256 /*id*/, uint256 /*value*/) public virtual override {
        revert("Use redeem function for settlement and totalSupply tracking");
    }

    /**
     * @dev override burnBatch to avoid assets burned without calling redeem
     */
    function burnBatch(address /*account*/, uint256[] memory /*ids*/, uint256[] memory /*values*/) public virtual override {
        revert("Use redeem function for settlement and totalSupply tracking");
    }

    /**
     * @dev update the asset metadata. the func can only be called the asset issuer
     * @param assetId asset Id
     * @param _uid generated by offchain service
     */
    function updateSystemOfRecordMetadata(uint256 assetId, string memory _uid) public {
        assets[assetId].uid = _uid;
        emit AssetUidUpdated(assetId, _uid);
    }

    // ERC1155 overrides to promise only active asset can be transferred
    function _update(
        address from,
        address to,
        uint256[] memory assetIds,
        uint256[] memory amounts
    ) internal virtual override(ERC1155) {
        // Check if the asset is verified
        for (uint256 i = 0; i < assetIds.length; i++) {
            require(assets[assetIds[i]].active, "Asset is not active");
        }
        super._update(from, to, assetIds, amounts);
        
        // emit the TokenTransferred when it is not burn or mint
        if (from != address(0) && to != address(0)) {
            for (uint256 i = 0; i < assetIds.length; i++) {
                if (amounts[i] > 0) {
                    emit TokensTransferred(assetIds[i], from, to, amounts[i]);
                }
            }
        }
    }

    function uid(uint256 assetId) public view returns (string memory) {
        return assets[assetId].uid;
    }

    // internal function: check if account is the asset issuer
    function _isIssuer(uint256 assetId, address account) internal view returns (bool) {
        return assetIssuers[assetId] == account;
    }

    // internal function: check if the account is asset issuer or admin
    function _isIssuerOrAdmin(uint256 assetId, address account) internal view returns (bool) {
        return _isIssuer(assetId, account) || hasRole(ADMIN_ROLE, account);
    }

    // view function to retrieve the total supply of the asset
    function totalSupply(uint256 id) public view returns (uint256) {
        return assets[id].totalSupply;
    }

    function supportsInterface(bytes4 interfaceId) public view override(ERC1155, AccessControl, ReceiverTemplate)
        returns (bool)
    {
        return super.supportsInterface(interfaceId);
    }
}
