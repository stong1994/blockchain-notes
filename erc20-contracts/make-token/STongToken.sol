pragma solidity ^0.5.2;

// import "./ERC20.sol";
import "./ERC20Detail.sol";
import "./ERC20Burnable.sol";

/**
 * @title SimpleToken
 * @dev Very simple ERC20 Token example, where all tokens are pre-assigned to the creator.
 * Note they can later distribute these tokens as they wish using `transfer` and other
 * `ERC20` functions.
 */
contract STongToken is ERC20Burnable, ERC20Detailed {
    uint8 public constant DECIMALS = 18;
    uint256 public constant INITIAL_SUPPLY = 10000 * (10 ** uint256(DECIMALS));

    /**
     * @dev Constructor that gives msg.sender all of existing tokens.
     */
    constructor () public payable ERC20Detailed("STongToken", "STT", DECIMALS) {
        _mint(msg.sender, INITIAL_SUPPLY);
    }
}
