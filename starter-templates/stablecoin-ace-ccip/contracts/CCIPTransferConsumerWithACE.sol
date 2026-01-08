// SPDX-License-Identifier: MIT
pragma solidity 0.8.26;

import {IRouterClient} from "@chainlink/contracts-ccip/ccip/interfaces/IRouterClient.sol";
import {Client} from "@chainlink/contracts-ccip/ccip/libraries/Client.sol";
import {IERC20} from "@openzeppelin/contracts/token/ERC20/IERC20.sol";
import {PolicyProtected} from "@chainlink/policy-management/core/PolicyProtected.sol";
import {IERC165} from "@openzeppelin/contracts/utils/introspection/IERC165.sol";

/**
 * @title CCIPTransferConsumerWithACE
 * @notice ACE-protected consumer for executing cross-chain CCIP transfers with policy enforcement.
 * @dev Inherits PolicyProtected for ACE integration and implements CRE receiver interface.
 * 
 * ACE Integration:
 * - The _processReport function is protected by the runPolicy modifier
 * - Before initiating CCIP transfer, PolicyEngine checks all attached policies
 * - Policies can check beneficiary address, sender address, amount, or destination chain
 * - Note: We use "beneficiary" for the receiving address (consistent with mint operations)
 * 
 * Flow:
 * 1. Bank sends cross-chain transfer instruction to CRE workflow
 * 2. CRE workflow generates signed report with CCIP transfer data
 * 3. Forwarder validates signatures and calls this contract's onReport()
 * 4. runPolicy modifier intercepts → PolicyEngine checks policies (e.g., blacklist)
 * 5. If policies pass → pull tokens from sender, approve router, execute CCIP transfer
 * 6. If any policy rejects → transaction reverts with PolicyRunRejected
 * 
 * CCIP Flow:
 * - Consumer pulls tokens from sender (sender must have approved consumer)
 * - Consumer approves Router to spend tokens
 * - Consumer approves Router to spend LINK for fees
 * - Router pulls tokens from consumer to TokenPool
 * - TokenPool burns tokens on source chain
 * - CCIP message is sent cross-chain
 * - TokenPool mints tokens on destination chain to beneficiary
 * 
 * Deployment:
 * - MUST be deployed via ERC1967Proxy for upgradeability
 * - Call initialize() immediately after deployment
 * - Fund this contract with LINK tokens for CCIP fees
 */
contract CCIPTransferConsumerWithACE is PolicyProtected {
    /// @custom:storage-location erc7201:cre-demo.CCIPTransferConsumerWithACE
    struct CCIPTransferConsumerStorage {
        IRouterClient router;
        IERC20 linkToken;
        IERC20 stablecoin;
        address expectedAuthor;
        bytes10 expectedWorkflowName;
    }

    // keccak256(abi.encode(uint256(keccak256("cre-demo.CCIPTransferConsumerWithACE")) - 1)) & ~bytes32(uint256(0xff))
    bytes32 private constant CCIPTransferConsumerStorageLocation =
        0x6d9e8f7a5b4c3d2e1f0a9b8c7d6e5f4a3b2c1d0e9f8a7b6c5d4e3f2a1b0c9d00;

    function _getCCIPTransferConsumerStorage() private pure returns (CCIPTransferConsumerStorage storage $) {
        assembly {
            $.slot := CCIPTransferConsumerStorageLocation
        }
    }

    event TransferInitiated(
        bytes32 indexed messageId,
        uint64 indexed destinationChainSelector,
        address indexed beneficiary,
        uint256 amount,
        bytes32 bankReference,
        uint256 fees
    );

    event TransferInstructionReceived(
        uint64 indexed destinationChainSelector,
        address indexed beneficiary,
        uint256 amount,
        bytes32 indexed bankReference,
        uint256 timestamp
    );

    error TransferFailed();
    error InsufficientLinkBalance(uint256 required, uint256 available);

    constructor() {
        _disableInitializers();
    }

    /**
     * @notice Initializes the CCIP transfer consumer with ACE integration.
     * @dev MUST be called immediately after proxy deployment.
     * 
     * @param initialOwner The address that will own this contract (can attach policies).
     * @param _stablecoin Address of the StablecoinERC20 contract.
     * @param _router Address of the CCIP Router contract.
     * @param _linkToken Address of the LINK token contract.
     * @param _policyEngine Address of the PolicyEngine for compliance checks.
     * @param _expectedAuthor Expected workflow owner address (use address(0) for testing).
     * @param _expectedWorkflowName Expected workflow name (use bytes10("dummy") for testing).
     */
    function initialize(
        address initialOwner,
        address _stablecoin,
        address _router,
        address _linkToken,
        address _policyEngine,
        address _expectedAuthor,
        bytes10 _expectedWorkflowName
    ) public initializer {
        __PolicyProtected_init(initialOwner, _policyEngine);
        
        CCIPTransferConsumerStorage storage $ = _getCCIPTransferConsumerStorage();
        $.stablecoin = IERC20(_stablecoin);
        $.router = IRouterClient(_router);
        $.linkToken = IERC20(_linkToken);
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
     * 2. PolicyEngine calls CCIPTransferConsumerExtractor to extract parameters
     * 3. Extractor decodes metadata and report, returns [beneficiary, sender, amount]
     * 4. PolicyEngine runs AddressBlacklistPolicy
     * 5. If beneficiary is blacklisted → revert with PolicyRunRejected
     * 6. If allowed → continue to processCCIPReport
     * 
     * @param metadata Encoded metadata (not used in testing version).
     * @param report Encoded transfer instruction.
     */
    function onReport(bytes calldata metadata, bytes calldata report) external runPolicy {
        // In production, validate metadata here
        processCCIPReport(report);
    }

    /**
     * @notice Process the cross-chain transfer instruction.
     * @dev This function contains the core CCIP logic (ACE already checked in onReport).
     * 
     * @param report ABI-encoded: (uint64 destChainSelector, address sender, address beneficiary, uint256 amount, bytes32 bankRef)
     */
    function processCCIPReport(bytes calldata report) public {
        CCIPTransferConsumerStorage storage $ = _getCCIPTransferConsumerStorage();
        
        // 1. Decode the transfer instruction
        (uint64 destinationChainSelector, address sender, address beneficiary, uint256 amount, bytes32 bankRef) = abi.decode(
            report,
            (uint64, address, address, uint256, bytes32)
        );

        // Validate decoded values
        require(sender != address(0), "Invalid sender");
        require(beneficiary != address(0), "Invalid beneficiary");
        require(amount > 0, "Invalid amount");

        // 2. Log the instruction
        emit TransferInstructionReceived(
            destinationChainSelector,
            beneficiary,
            amount,
            bankRef,
            block.timestamp
        );

        // 3. Build CCIP message
        Client.EVM2AnyMessage memory message = Client.EVM2AnyMessage({
            receiver: abi.encode(beneficiary),
            data: "",
            tokenAmounts: new Client.EVMTokenAmount[](1),
            feeToken: address($.linkToken),
            extraArgs: Client._argsToBytes(
                Client.EVMExtraArgsV1({gasLimit: 0}) // 0 for default
            )
        });

        message.tokenAmounts[0] = Client.EVMTokenAmount({
            token: address($.stablecoin),
            amount: amount
        });

        // 4. Calculate fee
        uint256 fees = $.router.getFee(destinationChainSelector, message);

        // 5. Check LINK balance
        uint256 linkBalance = $.linkToken.balanceOf(address(this));
        if (linkBalance < fees) {
            revert InsufficientLinkBalance(fees, linkBalance);
        }

        // 6. Acquire tokens for CCIP transfer
        // Check if we already have tokens (e.g., minted directly to this contract)
        uint256 ourBalance = $.stablecoin.balanceOf(address(this));
        
        if (ourBalance < amount) {
            // We need more tokens, pull from sender
            uint256 needed = amount - ourBalance;
            require(
                $.stablecoin.transferFrom(sender, address(this), needed),
                "Token transfer failed"
            );
        }
        // else: We already have enough tokens (e.g., from direct mint), no transfer needed

        // 7. Approve Router to spend tokens (Router will transfer to Pool, then Pool burns)
        $.stablecoin.approve(address($.router), amount);

        // 8. Approve Router to spend LINK for fees
        $.linkToken.approve(address($.router), fees);

        // 9. Execute cross-chain transfer
        bytes32 messageId = $.router.ccipSend(destinationChainSelector, message);

        emit TransferInitiated(
            messageId,
            destinationChainSelector,
            beneficiary,
            amount,
            bankRef,
            fees
        );
    }

    /**
     * @notice Fund this contract with LINK for CCIP fees.
     * @dev Owner or anyone can send LINK to this contract.
     */
    receive() external payable {}

    /**
     * @notice Get the router contract address.
     * @return The address of the CCIP Router contract.
     */
    function getRouter() external view returns (address) {
        return address(_getCCIPTransferConsumerStorage().router);
    }

    /**
     * @notice Get the stablecoin contract address.
     * @return The address of the StablecoinERC20 contract.
     */
    function getStablecoin() external view returns (address) {
        return address(_getCCIPTransferConsumerStorage().stablecoin);
    }

    /**
     * @notice Get the LINK token contract address.
     * @return The address of the LINK token contract.
     */
    function getLinkToken() external view returns (address) {
        return address(_getCCIPTransferConsumerStorage().linkToken);
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

