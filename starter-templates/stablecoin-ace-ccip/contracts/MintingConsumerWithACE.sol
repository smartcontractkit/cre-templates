// SPDX-License-Identifier: MIT
pragma solidity 0.8.26;

import {StablecoinERC20} from "./StablecoinERC20.sol";
import {PolicyProtected} from "@chainlink/policy-management/core/PolicyProtected.sol";
import {IERC165} from "@openzeppelin/contracts/utils/introspection/IERC165.sol";

/**
 * @title MintingConsumerWithACE
 * @notice ACE-protected consumer contract for minting/redeeming stablecoins with policy enforcement.
 * @dev Inherits PolicyProtected for ACE integration and implements CRE receiver interface.
 * 
 * ACE Integration:
 * - The _processReport function is protected by the runPolicy modifier
 * - Before minting/burning, the PolicyEngine checks all attached policies
 * - Policies can check beneficiary address, amount, or any extracted parameter
 * 
 * Flow:
 * 1. Bank sends SWIFT MT103 message to CRE workflow (HTTP trigger)
 * 2. CRE workflow performs PoR check (validates sufficient reserves)
 * 3. CRE workflow generates signed report with mint/redeem instruction
 * 4. Forwarder validates signatures and calls this contract's onReport()
 * 5. runPolicy modifier intercepts → PolicyEngine checks policies (e.g., blacklist)
 * 6. If policies pass → mint or burn stablecoins
 * 7. If any policy rejects → transaction reverts with PolicyRunRejected
 * 
 * Deployment:
 * - MUST be deployed via ERC1967Proxy for upgradeability
 * - Call initialize() immediately after deployment
 */
contract MintingConsumerWithACE is PolicyProtected {
    /// @custom:storage-location erc7201:cre-demo.MintingConsumerWithACE
    struct MintingConsumerStorage {
        StablecoinERC20 stablecoin;
        address expectedAuthor;
        bytes10 expectedWorkflowName;
    }

    // keccak256(abi.encode(uint256(keccak256("cre-demo.MintingConsumerWithACE")) - 1)) & ~bytes32(uint256(0xff))
    bytes32 private constant MintingConsumerStorageLocation =
        0x3f8e2a1b5c9d7e6f4a3b8c2d1e9f7a6b5c4d3e2f1a0b9c8d7e6f5a4b3c2d1e00;

    function _getMintingConsumerStorage() private pure returns (MintingConsumerStorage storage $) {
        assembly {
            $.slot := MintingConsumerStorageLocation
        }
    }

    // Instruction types
    uint8 constant INSTRUCTION_MINT = 1;
    uint8 constant INSTRUCTION_REDEEM = 2;

    event MintInstructionReceived(
        address indexed recipient,
        uint256 amount,
        bytes32 indexed bankReference,
        uint256 timestamp
    );

    event RedeemInstructionReceived(
        address indexed account,
        uint256 amount,
        bytes32 indexed bankReference,
        uint256 timestamp
    );

    event MintExecuted(
        address indexed recipient,
        uint256 amount,
        bytes32 indexed bankReference
    );

    event RedeemExecuted(
        address indexed account,
        uint256 amount,
        bytes32 indexed bankReference
    );

    error MintFailed();
    error RedeemFailed();
    error InvalidInstructionType();

    constructor() {
        _disableInitializers();
    }

    /**
     * @notice Initializes the minting consumer with ACE integration.
     * @dev MUST be called immediately after proxy deployment.
     * 
     * @param initialOwner The address that will own this contract (can attach policies).
     * @param _stablecoin Address of the StablecoinERC20 contract.
     * @param _policyEngine Address of the PolicyEngine for compliance checks.
     * @param _expectedAuthor Expected workflow owner address (use address(0) for testing).
     * @param _expectedWorkflowName Expected workflow name (use bytes10("dummy") for testing).
     */
    function initialize(
        address initialOwner,
        address _stablecoin,
        address _policyEngine,
        address _expectedAuthor,
        bytes10 _expectedWorkflowName
    ) public initializer {
        __PolicyProtected_init(initialOwner, _policyEngine);
        
        MintingConsumerStorage storage $ = _getMintingConsumerStorage();
        $.stablecoin = StablecoinERC20(_stablecoin);
        $.expectedAuthor = _expectedAuthor;
        $.expectedWorkflowName = _expectedWorkflowName;
    }

    /**
     * @notice Receive report from CRE Forwarder with ACE policy enforcement.
     * @dev This is the entry point called by the CRE Forwarder after signature validation.
     *      The runPolicy modifier triggers ACE to check all attached policies.
     * 
     * Policy Check Flow:
     * 1. runPolicy modifier intercepts this call
     * 2. PolicyEngine calls OnReportMintingExtractor to extract parameters
     * 3. Extractor decodes metadata and report, returns [beneficiary, amount]
     * 4. PolicyEngine runs AddressBlacklistPolicy
     * 5. If beneficiary is blacklisted → revert with PolicyRunRejected
     * 6. If allowed → continue to processMintReport
     * 
     * @param metadata Encoded metadata (not used in testing version).
     * @param report Encoded instruction: (uint8 type, address account, uint256 amount, bytes32 bankRef).
     */
    function onReport(bytes calldata metadata, bytes calldata report) external runPolicy {
        // In production, validate metadata here (workflow name, author, etc.)
        // For testing/demo purposes, we skip validation
        
        processMintReport(report);
    }

    /**
     * @notice Process the mint or redeem instruction.
     * @dev This function contains the core mint/burn logic (ACE already checked in onReport).
     * 
     * @param report ABI-encoded: (uint8 instructionType, address account, uint256 amount, bytes32 bankRef)
     */
    function processMintReport(bytes calldata report) public {
        MintingConsumerStorage storage $ = _getMintingConsumerStorage();
        
        // Decode the report
        (uint8 instructionType, address account, uint256 amount, bytes32 bankRef) = abi.decode(
            report,
            (uint8, address, uint256, bytes32)
        );

        if (instructionType == INSTRUCTION_MINT) {
            // Log the mint instruction
            emit MintInstructionReceived(account, amount, bankRef, block.timestamp);

            // Execute the mint
            try $.stablecoin.mint(account, amount) {
                emit MintExecuted(account, amount, bankRef);
            } catch {
                revert MintFailed();
            }
        } else if (instructionType == INSTRUCTION_REDEEM) {
            // Log the redeem instruction
            emit RedeemInstructionReceived(account, amount, bankRef, block.timestamp);

            // Execute the burn
            try $.stablecoin.burnFrom(account, amount) {
                emit RedeemExecuted(account, amount, bankRef);
            } catch {
                revert RedeemFailed();
            }
        } else {
            revert InvalidInstructionType();
        }
    }

    /**
     * @notice Get the stablecoin contract address.
     * @return The address of the StablecoinERC20 contract.
     */
    function getStablecoin() external view returns (address) {
        return address(_getMintingConsumerStorage().stablecoin);
    }

    /**
     * @notice ERC165 interface support.
     * @dev Overrides PolicyProtected's supportsInterface to include CRE receiver interface.
     * @param interfaceId The interface identifier to check.
     * @return True if the interface is supported.
     */
    function supportsInterface(bytes4 interfaceId) 
        public 
        view 
        virtual 
        override 
        returns (bool) 
    {
        return interfaceId == this.onReport.selector || super.supportsInterface(interfaceId);
    }
}

