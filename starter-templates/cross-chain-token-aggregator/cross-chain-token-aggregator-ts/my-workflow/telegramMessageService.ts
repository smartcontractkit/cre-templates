import { cre, type Runtime, consensusIdenticalAggregation, HTTPSendRequester, ok } from "@chainlink/cre-sdk";
import { IConfig } from "./interfaces/IConfig";

// to get telegram token, make the user message on your telegram, then get chat id and store on chain
// if already present then send to that chat id

export const sendTelegramMessage = (runtime: Runtime<IConfig>, message: string) => {
    try {
        const chatId = runtime.getSecret({id: "TELEGRAM_CHAT_ID"}).result().value;
        const botToken = runtime.getSecret({id: "TELEGRAM_BOT_ACCESS_TOKEN"}).result().value;
        const httpClient = new cre.capabilities.HTTPClient();
    
        httpClient.sendRequest(
            runtime,
            _sendTelegramMessage,
            consensusIdenticalAggregation<boolean>()
        )(runtime.config, chatId, message, botToken).result();
    }
    catch (error) {
        runtime.log("Error sending notification on telegram, continuing workflow...");
    }
}

const _sendTelegramMessage = (sendRequester: HTTPSendRequester, config: IConfig, chatId: string, message: string, botToken: string): boolean => {
    let telegramUrl = config.telegramMessageApi;
    telegramUrl = telegramUrl.replace("{{TELEGRAM_BOT_ACCESS_TOKEN}}", botToken);
    telegramUrl = telegramUrl.replace("{{TELEGRAM_CHAT_ID}}", chatId);
    telegramUrl = telegramUrl.replace("{{MESSAGE}}", message);
    
    const messageReq = {
        url: telegramUrl,
        method: "GET" as const,
        cacheSettings: {
            store: true,
            maxAge: '60s',
        },
    };

    const result = sendRequester.sendRequest(messageReq).result();
    return ok(result);
}