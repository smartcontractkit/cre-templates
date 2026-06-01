import { IBridge } from "./IBridge";
import { AcrossBridge } from "./integrations/across";
import { ChainlinkCCIPBridge } from "./integrations/chainlinkCCIP";

export class BridgeFactory {
  static getBridge(config: string): IBridge {
    switch (config.toLowerCase()) {
      case "across":
        return new AcrossBridge();
      case "chainlink_ccip":
        return new ChainlinkCCIPBridge();
      default:
        throw new Error("Invalid bridge type");
    }
  }
}