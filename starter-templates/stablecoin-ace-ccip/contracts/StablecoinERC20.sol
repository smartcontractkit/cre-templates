// SPDX-License-Identifier: MIT
pragma solidity ^0.8.24;

import {ERC20} from "@openzeppelin/contracts/token/ERC20/ERC20.sol";
import {ERC20Burnable} from "@openzeppelin/contracts/token/ERC20/extensions/ERC20Burnable.sol";
import {Ownable} from "@openzeppelin/contracts/access/Ownable.sol";
import {IERC165} from "@openzeppelin/contracts/utils/introspection/IERC165.sol";
import {IERC20} from "@openzeppelin/contracts/token/ERC20/IERC20.sol";

// CCIP TokenPool compatibility interface
interface IBurnMintERC20 is IERC20 {
    function mint(address account, uint256 amount) external;
    function burnFrom(address account, uint256 amount) external;
}

/**
 * @title StablecoinERC20
 * @notice Simplified stablecoin for CRE bank mint demo with CCIP compatibility
 * @dev Role-based minting for authorized minters (MintingConsumer, TokenPools)
 */
contract StablecoinERC20 is ERC20, ERC20Burnable, Ownable, IERC165 {
    mapping(address => bool) private s_minters;
    mapping(address => bool) private s_burners;
    
    event MinterAdded(address indexed minter);
    event MinterRemoved(address indexed minter);
    event BurnerAdded(address indexed burner);
    event BurnerRemoved(address indexed burner);
    
    error OnlyMinter();
    error OnlyBurner();
    
    constructor(
        string memory name,
        string memory symbol
    ) ERC20(name, symbol) Ownable(msg.sender) {}
    
    modifier onlyMinter() {
        if (!s_minters[msg.sender]) revert OnlyMinter();
        _;
    }
    
    modifier onlyBurner() {
        if (!s_burners[msg.sender]) revert OnlyBurner();
        _;
    }
    
    function grantMintRole(address minter) external onlyOwner {
        s_minters[minter] = true;
        emit MinterAdded(minter);
    }
    
    function revokeMintRole(address minter) external onlyOwner {
        s_minters[minter] = false;
        emit MinterRemoved(minter);
    }
    
    function grantBurnRole(address burner) external onlyOwner {
        s_burners[burner] = true;
        emit BurnerAdded(burner);
    }
    
    function revokeBurnRole(address burner) external onlyOwner {
        s_burners[burner] = false;
        emit BurnerRemoved(burner);
    }
    
    function isMinter(address account) external view returns (bool) {
        return s_minters[account];
    }
    
    function isBurner(address account) external view returns (bool) {
        return s_burners[account];
    }
    
    /**
     * @notice Grant both mint and burn roles (used by CCIP TokenPools)
     * @dev Convenience function for CCIP integration
     */
    function grantMintAndBurnRoles(address minterBurner) external onlyOwner {
        s_minters[minterBurner] = true;
        emit MinterAdded(minterBurner);
        s_burners[minterBurner] = true;
        emit BurnerAdded(minterBurner);
    }
    
    /**
     * @notice Mint tokens (called by authorized minters like MintingConsumer)
     * @dev Only addresses with minter role can call this
     */
    function mint(address account, uint256 amount) external onlyMinter {
        _mint(account, amount);
    }
    
    /**
     * @notice Burn tokens from account (called by authorized burners like MintingConsumer)
     * @dev Only addresses with burner role can call this
     * @dev Overrides ERC20Burnable to add role check for redemption flow
     */
    function burnFrom(address account, uint256 amount) public override onlyBurner {
        _burn(account, amount);
    }
    
    /**
     * @notice IERC165 support for CCIP TokenPool compatibility
     * @dev Returns true for IBurnMintERC20, IERC20, and IERC165 interfaces
     */
    function supportsInterface(bytes4 interfaceId) public pure override returns (bool) {
        return
            interfaceId == type(IERC20).interfaceId ||
            interfaceId == type(IBurnMintERC20).interfaceId ||
            interfaceId == type(IERC165).interfaceId;
    }
}

