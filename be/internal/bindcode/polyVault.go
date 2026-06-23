// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package bindcode

import (
	"errors"
	"math/big"
	"strings"

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/event"
)

// Reference imports to suppress errors if they are not otherwise used.
var (
	_ = errors.New
	_ = big.NewInt
	_ = strings.NewReader
	_ = ethereum.NotFound
	_ = bind.Bind
	_ = common.Big1
	_ = types.BloomLookup
	_ = event.NewSubscription
	_ = abi.ConvertType
)

// IPolyVaultWithdrawalRequest is an auto generated low-level Go binding around an user-defined struct.
type IPolyVaultWithdrawalRequest struct {
	Shares           *big.Int
	RequestTimestamp *big.Int
	Pending          bool
}

// PolyVaultMetaData contains all meta data concerning the PolyVault contract.
var PolyVaultMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[],\"name\":\"AccessControlBadConfirmation\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"},{\"internalType\":\"bytes32\",\"name\":\"neededRole\",\"type\":\"bytes32\"}],\"name\":\"AccessControlUnauthorizedAccount\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"target\",\"type\":\"address\"}],\"name\":\"AddressEmptyCode\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"maximum\",\"type\":\"uint256\"}],\"name\":\"DepositAboveMaximum\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"minimum\",\"type\":\"uint256\"}],\"name\":\"DepositBelowMinimum\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"DirectWithdrawDisabled\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"implementation\",\"type\":\"address\"}],\"name\":\"ERC1967InvalidImplementation\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"ERC1967NonPayable\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"allowance\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"needed\",\"type\":\"uint256\"}],\"name\":\"ERC20InsufficientAllowance\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"balance\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"needed\",\"type\":\"uint256\"}],\"name\":\"ERC20InsufficientBalance\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"approver\",\"type\":\"address\"}],\"name\":\"ERC20InvalidApprover\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"receiver\",\"type\":\"address\"}],\"name\":\"ERC20InvalidReceiver\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"}],\"name\":\"ERC20InvalidSender\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"}],\"name\":\"ERC20InvalidSpender\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"receiver\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"assets\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"max\",\"type\":\"uint256\"}],\"name\":\"ERC4626ExceededMaxDeposit\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"receiver\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"shares\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"max\",\"type\":\"uint256\"}],\"name\":\"ERC4626ExceededMaxMint\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"shares\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"max\",\"type\":\"uint256\"}],\"name\":\"ERC4626ExceededMaxRedeem\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"assets\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"max\",\"type\":\"uint256\"}],\"name\":\"ERC4626ExceededMaxWithdraw\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"EnforcedPause\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"ExpectedPause\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"FailedCall\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"requested\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"available\",\"type\":\"uint256\"}],\"name\":\"InsufficientShares\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"allocation\",\"type\":\"uint256\"}],\"name\":\"InvalidAllocation\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidInitialization\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"fee\",\"type\":\"uint256\"}],\"name\":\"InvalidPerformanceFee\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"delay\",\"type\":\"uint256\"}],\"name\":\"InvalidWithdrawalDelay\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"NoPendingWithdrawal\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"NotInitializing\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"ReentrancyGuardReentrantCall\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"}],\"name\":\"SafeERC20FailedOperation\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"requested\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"available\",\"type\":\"uint256\"}],\"name\":\"StrategyAllocationExceeded\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"UUPSUnauthorizedCallContext\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"slot\",\"type\":\"bytes32\"}],\"name\":\"UUPSUnsupportedProxiableUUID\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"WithdrawalAlreadyPending\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"currentTime\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"availableTime\",\"type\":\"uint256\"}],\"name\":\"WithdrawalDelayNotMet\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"ZeroAddress\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"ZeroAmount\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"Approval\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"assets\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"shares\",\"type\":\"uint256\"}],\"name\":\"Deposit\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"minDeposit\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"maxDeposit\",\"type\":\"uint256\"}],\"name\":\"DepositLimitsUpdated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"oldRecipient\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"newRecipient\",\"type\":\"address\"}],\"name\":\"FeeRecipientUpdated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"version\",\"type\":\"uint64\"}],\"name\":\"Initialized\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"oldAllocation\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"newAllocation\",\"type\":\"uint256\"}],\"name\":\"MaxStrategyAllocationUpdated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"Paused\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"oldFee\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"newFee\",\"type\":\"uint256\"}],\"name\":\"PerformanceFeeUpdated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"profit\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"fee\",\"type\":\"uint256\"}],\"name\":\"ProfitReported\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"previousAdminRole\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"newAdminRole\",\"type\":\"bytes32\"}],\"name\":\"RoleAdminChanged\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"}],\"name\":\"RoleGranted\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"}],\"name\":\"RoleRevoked\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"strategist\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"StrategyDeposit\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"strategist\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"StrategyWithdrawal\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"Transfer\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"Unpaused\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"implementation\",\"type\":\"address\"}],\"name\":\"Upgraded\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"receiver\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"assets\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"shares\",\"type\":\"uint256\"}],\"name\":\"Withdraw\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"user\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"shares\",\"type\":\"uint256\"}],\"name\":\"WithdrawalCancelled\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"oldDelay\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"newDelay\",\"type\":\"uint256\"}],\"name\":\"WithdrawalDelayUpdated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"user\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"shares\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"assets\",\"type\":\"uint256\"}],\"name\":\"WithdrawalExecuted\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"user\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"shares\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"timestamp\",\"type\":\"uint256\"}],\"name\":\"WithdrawalRequested\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"BASIS_POINTS\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"DEFAULT_ADMIN_ROLE\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"GUARDIAN_ROLE\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"MAX_PERFORMANCE_FEE\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"MAX_WITHDRAWAL_DELAY\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"MIN_WITHDRAWAL_DELAY\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"STRATEGIST_ROLE\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"UPGRADE_INTERFACE_VERSION\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"}],\"name\":\"allowance\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"approve\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"asset\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"availableBalance\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"balanceOf\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"cancelWithdraw\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"shares\",\"type\":\"uint256\"}],\"name\":\"convertToAssets\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"assets\",\"type\":\"uint256\"}],\"name\":\"convertToShares\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"decimals\",\"outputs\":[{\"internalType\":\"uint8\",\"name\":\"\",\"type\":\"uint8\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"assets\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"receiver\",\"type\":\"address\"}],\"name\":\"deposit\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"depositFromStrategy\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"executeWithdraw\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"feeRecipient\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"}],\"name\":\"getRoleAdmin\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"user\",\"type\":\"address\"}],\"name\":\"getWithdrawalRequest\",\"outputs\":[{\"components\":[{\"internalType\":\"uint256\",\"name\":\"shares\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"requestTimestamp\",\"type\":\"uint256\"},{\"internalType\":\"bool\",\"name\":\"pending\",\"type\":\"bool\"}],\"internalType\":\"structIPolyVault.WithdrawalRequest\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"grantRole\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"hasRole\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_usdc\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"_admin\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"_strategist\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"_guardian\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"_feeRecipient\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_withdrawalDelay\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"_maxAllocation\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"_performanceFee\",\"type\":\"uint256\"}],\"name\":\"initialize\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"maxDeposit\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"maxDeposit\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"maxMint\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"maxRedeem\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"maxStrategyAllocation\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"maxWithdraw\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"minDeposit\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"shares\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"receiver\",\"type\":\"address\"}],\"name\":\"mint\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"name\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"pause\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"paused\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"performanceFee\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"assets\",\"type\":\"uint256\"}],\"name\":\"previewDeposit\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"shares\",\"type\":\"uint256\"}],\"name\":\"previewMint\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"shares\",\"type\":\"uint256\"}],\"name\":\"previewRedeem\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"assets\",\"type\":\"uint256\"}],\"name\":\"previewWithdraw\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"proxiableUUID\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"redeem\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"callerConfirmation\",\"type\":\"address\"}],\"name\":\"renounceRole\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"shares\",\"type\":\"uint256\"}],\"name\":\"requestWithdraw\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"revokeRole\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_min\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"_max\",\"type\":\"uint256\"}],\"name\":\"setDepositLimits\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_recipient\",\"type\":\"address\"}],\"name\":\"setFeeRecipient\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_allocation\",\"type\":\"uint256\"}],\"name\":\"setMaxStrategyAllocation\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_fee\",\"type\":\"uint256\"}],\"name\":\"setPerformanceFee\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_delay\",\"type\":\"uint256\"}],\"name\":\"setWithdrawalDelay\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"strategyDebt\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes4\",\"name\":\"interfaceId\",\"type\":\"bytes4\"}],\"name\":\"supportsInterface\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"symbol\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"totalAssets\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"totalSupply\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"transfer\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"transferFrom\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"unpause\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newImplementation\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"}],\"name\":\"upgradeToAndCall\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"withdraw\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"withdrawToStrategy\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"withdrawalDelay\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
	Bin: "0x60a06040523073ffffffffffffffffffffffffffffffffffffffff1660809073ffffffffffffffffffffffffffffffffffffffff16815250348015610042575f5ffd5b5061005161005660201b60201c565b6101d1565b5f61006561015460201b60201c565b9050805f0160089054906101000a900460ff16156100af576040517ff92ee8a900000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b67ffffffffffffffff8016815f015f9054906101000a900467ffffffffffffffff1667ffffffffffffffff16146101515767ffffffffffffffff815f015f6101000a81548167ffffffffffffffff021916908367ffffffffffffffff1602179055507fc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d267ffffffffffffffff60405161014891906101b8565b60405180910390a15b50565b5f5f61016461016d60201b60201c565b90508091505090565b5f7ff0c57e16840df040f15088dc2f81fe391c3923bec73e23a9662efc9c229c6a005f1b905090565b5f67ffffffffffffffff82169050919050565b6101b281610196565b82525050565b5f6020820190506101cb5f8301846101a9565b92915050565b6080516156396101f75f395f818161323c01528181613291015261345001526156395ff3fe6080604052600436106103c2575f3560e01c80638456cb59116101f1578063b3d7f6b91161010c578063d2c13da51161009f578063e1f1c4a71161006e578063e1f1c4a714610eca578063e74b981b14610ef4578063ef8b30f714610f1c578063f8fd979514610f58576103c2565b8063d2c13da514610e02578063d547741f14610e2a578063d905777e14610e52578063dd62ed3e14610e8e576103c2565b8063c63d75b6116100db578063c63d75b614610d24578063c6e6f59214610d60578063ca1689d014610d9c578063ce96cb7714610dc6576103c2565b8063b3d7f6b914610c46578063b460af9414610c82578063ba08765214610cbe578063bdca916514610cfa576103c2565b80639aca479e11610184578063a7ab696111610153578063a7ab696114610b8c578063a9059cbb14610bb6578063ab2f0e5114610bf2578063ad3cb1cc14610c1c576103c2565b80639aca479e14610ae6578063a217fddf14610b0e578063a238f9df14610b38578063a378a32414610b62576103c2565b80638c661b5d116101c05780638c661b5d14610a0857806391d1485414610a4457806394bf804d14610a8057806395d89b4114610abc576103c2565b80638456cb591461098a57806384b76824146109a057806386f1ef0c146109b657806387788782146109de576103c2565b8063402d267d116102e15780635c975abb1161027457806370a082311161024357806370a08231146108d65780637403c6cd14610912578063745400c91461093a57806383c56b0214610962576103c2565b80635c975abb1461081e5780636083e59a146108485780636e553f651461087257806370897b23146108ae576103c2565b80634cdad506116102b05780634cdad506146107745780634eddea06146107b05780634f1ef286146107d857806352d1902d146107f4576103c2565b8063402d267d146106ba57806341b3d185146106f657806346904840146107205780634a5f2b5d1461074a576103c2565b8063248a9ca31161035957806336568abe1161032857806336568abe1461062857806338d52e0f146106505780633e9dc7621461067a5780633f4ba83a146106a4576103c2565b8063248a9ca31461057057806324ea54f4146105ac5780632f2ff15d146105d6578063313ce567146105fe576103c2565b8063095ea7b311610395578063095ea7b3146104925780630a28a477146104ce57806318160ddd1461050a57806323b872dd14610534576103c2565b806301e1d114146103c657806301ffc9a7146103f057806306fdde031461042c57806307a2d13a14610456575b5f5ffd5b3480156103d1575f5ffd5b506103da610f6e565b6040516103e791906146c2565b60405180910390f35b3480156103fb575f5ffd5b5061041660048036038101906104119190614741565b611000565b6040516104239190614786565b60405180910390f35b348015610437575f5ffd5b50610440611079565b60405161044d919061480f565b60405180910390f35b348015610461575f5ffd5b5061047c60048036038101906104779190614859565b611117565b60405161048991906146c2565b60405180910390f35b34801561049d575f5ffd5b506104b860048036038101906104b391906148de565b611129565b6040516104c59190614786565b60405180910390f35b3480156104d9575f5ffd5b506104f460048036038101906104ef9190614859565b61114b565b60405161050191906146c2565b60405180910390f35b348015610515575f5ffd5b5061051e61115e565b60405161052b91906146c2565b60405180910390f35b34801561053f575f5ffd5b5061055a6004803603810190610555919061491c565b611175565b6040516105679190614786565b60405180910390f35b34801561057b575f5ffd5b506105966004803603810190610591919061499f565b6111a3565b6040516105a391906149d9565b60405180910390f35b3480156105b7575f5ffd5b506105c06111cd565b6040516105cd91906149d9565b60405180910390f35b3480156105e1575f5ffd5b506105fc60048036038101906105f791906149f2565b6111f1565b005b348015610609575f5ffd5b50610612611213565b60405161061f9190614a4b565b60405180910390f35b348015610633575f5ffd5b5061064e600480360381019061064991906149f2565b611248565b005b34801561065b575f5ffd5b506106646112c3565b6040516106719190614a73565b60405180910390f35b348015610685575f5ffd5b5061068e6112f8565b60405161069b91906146c2565b60405180910390f35b3480156106af575f5ffd5b506106b86112fe565b005b3480156106c5575f5ffd5b506106e060048036038101906106db9190614a8c565b611333565b6040516106ed91906146c2565b60405180910390f35b348015610701575f5ffd5b5061070a61135c565b60405161071791906146c2565b60405180910390f35b34801561072b575f5ffd5b50610734611362565b6040516107419190614a73565b60405180910390f35b348015610755575f5ffd5b5061075e611387565b60405161076b91906146c2565b60405180910390f35b34801561077f575f5ffd5b5061079a60048036038101906107959190614859565b61138d565b6040516107a791906146c2565b60405180910390f35b3480156107bb575f5ffd5b506107d660048036038101906107d19190614ab7565b61139f565b005b6107f260048036038101906107ed9190614c21565b6113f7565b005b3480156107ff575f5ffd5b50610808611416565b60405161081591906149d9565b60405180910390f35b348015610829575f5ffd5b50610832611447565b60405161083f9190614786565b60405180910390f35b348015610853575f5ffd5b5061085c611469565b60405161086991906146c2565b60405180910390f35b34801561087d575f5ffd5b5061089860048036038101906108939190614c7b565b61146f565b6040516108a591906146c2565b60405180910390f35b3480156108b9575f5ffd5b506108d460048036038101906108cf9190614859565b611530565b005b3480156108e1575f5ffd5b506108fc60048036038101906108f79190614a8c565b6115ce565b60405161090991906146c2565b60405180910390f35b34801561091d575f5ffd5b5061093860048036038101906109339190614cb9565b611621565b005b348015610945575f5ffd5b50610960600480360381019061095b9190614859565b611aae565b005b34801561096d575f5ffd5b5061098860048036038101906109839190614859565b611cc7565b005b348015610995575f5ffd5b5061099e611f1a565b005b3480156109ab575f5ffd5b506109b4611f4f565b005b3480156109c1575f5ffd5b506109dc60048036038101906109d79190614859565b6120aa565b005b3480156109e9575f5ffd5b506109f2612148565b6040516109ff91906146c2565b60405180910390f35b348015610a13575f5ffd5b50610a2e6004803603810190610a299190614a8c565b61214e565b604051610a3b9190614dc8565b60405180910390f35b348015610a4f575f5ffd5b50610a6a6004803603810190610a6591906149f2565b6121d2565b604051610a779190614786565b60405180910390f35b348015610a8b575f5ffd5b50610aa66004803603810190610aa19190614c7b565b612243565b604051610ab391906146c2565b60405180910390f35b348015610ac7575f5ffd5b50610ad0612311565b604051610add919061480f565b60405180910390f35b348015610af1575f5ffd5b50610b0c6004803603810190610b079190614859565b6123af565b005b348015610b19575f5ffd5b50610b22612559565b604051610b2f91906149d9565b60405180910390f35b348015610b43575f5ffd5b50610b4c61255f565b604051610b5991906146c2565b60405180910390f35b348015610b6d575f5ffd5b50610b76612566565b604051610b8391906149d9565b60405180910390f35b348015610b97575f5ffd5b50610ba061258a565b604051610bad91906146c2565b60405180910390f35b348015610bc1575f5ffd5b50610bdc6004803603810190610bd791906148de565b61258f565b604051610be99190614786565b60405180910390f35b348015610bfd575f5ffd5b50610c066125b1565b604051610c1391906146c2565b60405180910390f35b348015610c27575f5ffd5b50610c30612636565b604051610c3d919061480f565b60405180910390f35b348015610c51575f5ffd5b50610c6c6004803603810190610c679190614859565b61266f565b604051610c7991906146c2565b60405180910390f35b348015610c8d575f5ffd5b50610ca86004803603810190610ca39190614de1565b612682565b604051610cb591906146c2565b60405180910390f35b348015610cc9575f5ffd5b50610ce46004803603810190610cdf9190614de1565b6126b5565b604051610cf191906146c2565b60405180910390f35b348015610d05575f5ffd5b50610d0e6126e8565b604051610d1b91906146c2565b60405180910390f35b348015610d2f575f5ffd5b50610d4a6004803603810190610d459190614a8c565b6126ee565b604051610d5791906146c2565b60405180910390f35b348015610d6b575f5ffd5b50610d866004803603810190610d819190614859565b612717565b604051610d9391906146c2565b60405180910390f35b348015610da7575f5ffd5b50610db0612729565b604051610dbd91906146c2565b60405180910390f35b348015610dd1575f5ffd5b50610dec6004803603810190610de79190614a8c565b61272f565b604051610df991906146c2565b60405180910390f35b348015610e0d575f5ffd5b50610e286004803603810190610e239190614859565b612738565b005b348015610e35575f5ffd5b50610e506004803603810190610e4b91906149f2565b6127e1565b005b348015610e5d575f5ffd5b50610e786004803603810190610e739190614a8c565b612803565b604051610e8591906146c2565b60405180910390f35b348015610e99575f5ffd5b50610eb46004803603810190610eaf9190614e31565b61280c565b604051610ec191906146c2565b60405180910390f35b348015610ed5575f5ffd5b50610ede61289c565b604051610eeb91906146c2565b60405180910390f35b348015610eff575f5ffd5b50610f1a6004803603810190610f159190614a8c565b6128a2565b005b348015610f27575f5ffd5b50610f426004803603810190610f3d9190614859565b6129b6565b604051610f4f91906146c2565b60405180910390f35b348015610f63575f5ffd5b50610f6c6129c8565b005b5f600454610f7a6112c3565b73ffffffffffffffffffffffffffffffffffffffff166370a08231306040518263ffffffff1660e01b8152600401610fb29190614a73565b602060405180830381865afa158015610fcd573d5f5f3e3d5ffd5b505050506040513d601f19601f82011682018060405250810190610ff19190614e83565b610ffb9190614edb565b905090565b5f7f7965db0b000000000000000000000000000000000000000000000000000000007bffffffffffffffffffffffffffffffffffffffffffffffffffffffff1916827bffffffffffffffffffffffffffffffffffffffffffffffffffffffff19161480611072575061107182612c97565b5b9050919050565b60605f611084612d00565b905080600301805461109590614f3b565b80601f01602080910402602001604051908101604052809291908181526020018280546110c190614f3b565b801561110c5780601f106110e35761010080835404028352916020019161110c565b820191905f5260205f20905b8154815290600101906020018083116110ef57829003601f168201915b505050505091505090565b5f611122825f612d27565b9050919050565b5f5f611133612d7f565b9050611140818585612d86565b600191505092915050565b5f611157826001612d98565b9050919050565b5f5f611168612d00565b9050806002015491505090565b5f5f61117f612d7f565b905061118c858285612df0565b611197858585612e83565b60019150509392505050565b5f5f6111ad612f73565b9050805f015f8481526020019081526020015f2060010154915050919050565b7f55435dd261a4b9b3364963f7738a7a662ad9c84396d64be3365284bb7f0a504181565b6111fa826111a3565b61120381612f9a565b61120d8383612fae565b50505050565b5f5f61121d6130a6565b90506112276130cd565b815f0160149054906101000a900460ff166112429190614f6b565b91505090565b611250612d7f565b73ffffffffffffffffffffffffffffffffffffffff168173ffffffffffffffffffffffffffffffffffffffff16146112b4576040517f6697b23200000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6112be82826130d4565b505050565b5f5f6112cd6130a6565b9050805f015f9054906101000a900473ffffffffffffffffffffffffffffffffffffffff1691505090565b60045481565b7f55435dd261a4b9b3364963f7738a7a662ad9c84396d64be3365284bb7f0a504161132881612f9a565b6113306131cc565b50565b5f7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff9050919050565b60025481565b60075f9054906101000a900473ffffffffffffffffffffffffffffffffffffffff1681565b610e1081565b5f611398825f612d27565b9050919050565b5f5f1b6113ab81612f9a565b82600281905550816003819055507fb2ad710f2954a5376267a683f9ece9ec46ee7dfb47075163379904ee941df8da83836040516113ea929190614f9f565b60405180910390a1505050565b6113ff61323a565b61140882613320565b6114128282613330565b5050565b5f61141f61344e565b7f360894a13ba1a3210667c828492db98dca3e2076cc3735a920a3ca505d382bbc5f1b905090565b5f5f6114516134d5565b9050805f015f9054906101000a900460ff1691505090565b60035481565b5f6114786134fc565b61148061353d565b6002548310156114cb57826002546040517fed28d0890000000000000000000000000000000000000000000000000000000081526004016114c2929190614f9f565b60405180910390fd5b60035483111561151657826003546040517fb074571a00000000000000000000000000000000000000000000000000000000815260040161150d929190614f9f565b60405180910390fd5b611520838361356a565b905061152a6135ea565b92915050565b5f5f1b61153c81612f9a565b6107d082111561158357816040517f29f99fe400000000000000000000000000000000000000000000000000000000815260040161157a91906146c2565b60405180910390fd5b5f6006549050826006819055507f607b1c943753982194530bf7133a5972ea2626e028005410efa54ab20035caf881846040516115c1929190614f9f565b60405180910390a1505050565b5f5f6115d8612d00565b9050805f015f8473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f2054915050919050565b5f61162a61360e565b90505f815f0160089054906101000a900460ff161590505f825f015f9054906101000a900467ffffffffffffffff1690505f5f8267ffffffffffffffff161480156116725750825b90505f60018367ffffffffffffffff161480156116a557505f3073ffffffffffffffffffffffffffffffffffffffff163b145b9050811580156116b3575080155b156116ea576040517ff92ee8a900000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6001855f015f6101000a81548167ffffffffffffffff021916908367ffffffffffffffff1602179055508315611737576001855f0160086101000a81548160ff0219169083151502179055505b5f73ffffffffffffffffffffffffffffffffffffffff168d73ffffffffffffffffffffffffffffffffffffffff16148061179c57505f73ffffffffffffffffffffffffffffffffffffffff168c73ffffffffffffffffffffffffffffffffffffffff16145b806117d257505f73ffffffffffffffffffffffffffffffffffffffff168973ffffffffffffffffffffffffffffffffffffffff16145b15611809576040517fd92e233d00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b610e1088108061181b575062093a8088115b1561185d57876040517f8ef4581900000000000000000000000000000000000000000000000000000000815260040161185491906146c2565b60405180910390fd5b6107d08611156118a457856040517f29f99fe400000000000000000000000000000000000000000000000000000000815260040161189b91906146c2565b60405180910390fd5b6127108711156118eb57866040517f6f6773380000000000000000000000000000000000000000000000000000000081526004016118e291906146c2565b60405180910390fd5b61195f6040518060400160405280600e81526020017f506f6c795661756c7420555344430000000000000000000000000000000000008152506040518060400160405280600681526020017f7076555344430000000000000000000000000000000000000000000000000000815250613621565b6119688d613637565b61197061364b565b611978613655565b6119845f5f1b8d612fae565b506119af7f17a8e30262c1f919c33056d877a3c22b95c2f5e4dac44683c1c2323cd79fbdb08c612fae565b506119da7f55435dd261a4b9b3364963f7738a7a662ad9c84396d64be3365284bb7f0a50418b612fae565b50875f8190555086600581905550856006819055508860075f6101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff160217905550620f424060028190555064174876e8006003819055508315611a9f575f855f0160086101000a81548160ff0219169083151502179055507fc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d26001604051611a96919061501b565b60405180910390a15b50505050505050505050505050565b611ab66134fc565b611abe61353d565b5f8103611af7576040517f1f2a200500000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60015f3373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f206002015f9054906101000a900460ff1615611b7b576040517fa2aed0f000000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b80611b85336115ce565b1015611bd25780611b95336115ce565b6040517fcb1d8bba000000000000000000000000000000000000000000000000000000008152600401611bc9929190614f9f565b60405180910390fd5b611bdd333083612e83565b60405180606001604052808281526020014281526020016001151581525060015f3373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205f820151815f0155602082015181600101556040820151816002015f6101000a81548160ff0219169083151502179055509050503373ffffffffffffffffffffffffffffffffffffffff167f24b91f4f47caf44230a57777a9be744924e82bf666f2d5702faf97df35e60f9f8242604051611cb4929190614f9f565b60405180910390a2611cc46135ea565b50565b7f17a8e30262c1f919c33056d877a3c22b95c2f5e4dac44683c1c2323cd79fbdb0611cf181612f9a565b611cf961353d565b5f8203611d32576040517f1f2a200500000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b611d66333084611d406112c3565b73ffffffffffffffffffffffffffffffffffffffff1661365f909392919063ffffffff16565b6004548210611ea7575f60045483611d7e9190615034565b90505f6004819055505f81118015611d9757505f600654115b8015611df057505f73ffffffffffffffffffffffffffffffffffffffff1660075f9054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1614155b15611ea1575f61271060065483611e079190615067565b611e1191906150d5565b9050611e6660075f9054906101000a900473ffffffffffffffffffffffffffffffffffffffff1682611e416112c3565b73ffffffffffffffffffffffffffffffffffffffff166136b49092919063ffffffff16565b7ffcfbfd1d7fecbea7809bda42bd54ffa877192d8f5170375720ba7197c80181bc8282604051611e97929190614f9f565b60405180910390a1505b50611ec0565b8160045f828254611eb89190615034565b925050819055505b3373ffffffffffffffffffffffffffffffffffffffff167fc6f6f91a48277d76f232cc08a9a30f6b05b3fd9b92c3180c25936e17a22a102583604051611f0691906146c2565b60405180910390a2611f166135ea565b5050565b7f55435dd261a4b9b3364963f7738a7a662ad9c84396d64be3365284bb7f0a5041611f4481612f9a565b611f4c613707565b50565b611f5761353d565b5f60015f3373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f209050806002015f9054906101000a900460ff16611fde576040517f9121b84f00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b5f815f0154905060015f3373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205f5f82015f9055600182015f9055600282015f6101000a81549060ff02191690555050612050303383612e83565b3373ffffffffffffffffffffffffffffffffffffffff167f2eed97477f07c07ec48f8f678f4e84f7c0de55bf33f51c3dc989b133530803198260405161209691906146c2565b60405180910390a250506120a86135ea565b565b5f5f1b6120b681612f9a565b6127108211156120fd57816040517f6f6773380000000000000000000000000000000000000000000000000000000081526004016120f491906146c2565b60405180910390fd5b5f6005549050826005819055507f5755997f472615e724fa79d9ae7a5a69a67307f0eef7b848b07f517abac1adaf818460405161213b929190614f9f565b60405180910390a1505050565b60065481565b61215661468a565b60015f8373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f206040518060600160405290815f820154815260200160018201548152602001600282015f9054906101000a900460ff1615151515815250509050919050565b5f5f6121dc612f73565b9050805f015f8581526020019081526020015f205f015f8473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205f9054906101000a900460ff1691505092915050565b5f61224c6134fc565b61225461353d565b5f61225e8461266f565b90506002548110156122ab57806002546040517fed28d0890000000000000000000000000000000000000000000000000000000081526004016122a2929190614f9f565b60405180910390fd5b6003548111156122f657806003546040517fb074571a0000000000000000000000000000000000000000000000000000000081526004016122ed929190614f9f565b60405180910390fd5b6123008484613776565b91505061230b6135ea565b92915050565b60605f61231c612d00565b905080600401805461232d90614f3b565b80601f016020809104026020016040519081016040528092919081815260200182805461235990614f3b565b80156123a45780601f1061237b576101008083540402835291602001916123a4565b820191905f5260205f20905b81548152906001019060200180831161238757829003601f168201915b505050505091505090565b7f17a8e30262c1f919c33056d877a3c22b95c2f5e4dac44683c1c2323cd79fbdb06123d981612f9a565b6123e161353d565b5f820361241a576040517f1f2a200500000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b5f612710600554612429610f6e565b6124339190615067565b61243d91906150d5565b9050808360045461244e9190614edb565b11156124b4575f6004548211612464575f612473565b600454826124729190615034565b5b905083816040517f5b0b300c0000000000000000000000000000000000000000000000000000000081526004016124ab929190614f9f565b60405180910390fd5b8260045f8282546124c59190614edb565b925050819055506124fe33846124d96112c3565b73ffffffffffffffffffffffffffffffffffffffff166136b49092919063ffffffff16565b3373ffffffffffffffffffffffffffffffffffffffff167fd5ad0f046bd35f48b421a3e575435de38cea1980177b1c6da935d2f26049f3fa8460405161254491906146c2565b60405180910390a2506125556135ea565b5050565b5f5f1b81565b62093a8081565b7f17a8e30262c1f919c33056d877a3c22b95c2f5e4dac44683c1c2323cd79fbdb081565b5f5481565b5f5f612599612d7f565b90506125a6818585612e83565b600191505092915050565b5f6125ba6112c3565b73ffffffffffffffffffffffffffffffffffffffff166370a08231306040518263ffffffff1660e01b81526004016125f29190614a73565b602060405180830381865afa15801561260d573d5f5f3e3d5ffd5b505050506040513d601f19601f820116820180604052508101906126319190614e83565b905090565b6040518060400160405280600581526020017f352e302e3000000000000000000000000000000000000000000000000000000081525081565b5f61267b826001612d27565b9050919050565b5f6040517ff8a6b6dc00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b5f6040517ff8a6b6dc00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6107d081565b5f7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff9050919050565b5f612722825f612d98565b9050919050565b60055481565b5f5f9050919050565b5f5f1b61274481612f9a565b610e10821080612756575062093a8082115b1561279857816040517f8ef4581900000000000000000000000000000000000000000000000000000000815260040161278f91906146c2565b60405180910390fd5b5f5f549050825f819055507f9c3f1b54b1487e018f1d0593ff5cf7fb625b2df6332c974a6cc56bb35887984181846040516127d4929190614f9f565b60405180910390a1505050565b6127ea826111a3565b6127f381612f9a565b6127fd83836130d4565b50505050565b5f5f9050919050565b5f5f612816612d00565b9050806001015f8573ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205f8473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205491505092915050565b61271081565b5f5f1b6128ae81612f9a565b5f73ffffffffffffffffffffffffffffffffffffffff168273ffffffffffffffffffffffffffffffffffffffff1603612913576040517fd92e233d00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b5f60075f9054906101000a900473ffffffffffffffffffffffffffffffffffffffff1690508260075f6101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff1602179055507faaebcf1bfa00580e41d966056b48521fa9f202645c86d4ddf28113e617c1b1d381846040516129a9929190615105565b60405180910390a1505050565b5f6129c1825f612d98565b9050919050565b6129d061353d565b5f60015f3373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f209050806002015f9054906101000a900460ff16612a57576040517f9121b84f00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b5f5f548260010154612a699190614edb565b905080421015612ab25742816040517fdaa936b0000000000000000000000000000000000000000000000000000000008152600401612aa9929190614f9f565b60405180910390fd5b5f825f015490505f612ac382611117565b905060015f3373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205f5f82015f9055600182015f9055600282015f6101000a81549060ff021916905550505f612b2e6112c3565b73ffffffffffffffffffffffffffffffffffffffff166370a08231306040518263ffffffff1660e01b8152600401612b669190614a73565b602060405180830381865afa158015612b81573d5f5f3e3d5ffd5b505050506040513d601f19601f82011682018060405250810190612ba59190614e83565b90505f5f828411612bbb57849150839050612bfa565b829050612bc783612717565b915084821115612bd5578491505b5f8286612be29190615034565b90505f811115612bf857612bf7303383612e83565b5b505b612c0430836137f6565b612c363382612c116112c3565b73ffffffffffffffffffffffffffffffffffffffff166136b49092919063ffffffff16565b3373ffffffffffffffffffffffffffffffffffffffff167f37ce46bc94895501203dc6abbf2b2e0d502e856e3cd90186faaba6dab7d316bb8383604051612c7e929190614f9f565b60405180910390a250505050505050612c956135ea565b565b5f7f01ffc9a7000000000000000000000000000000000000000000000000000000007bffffffffffffffffffffffffffffffffffffffffffffffffffffffff1916827bffffffffffffffffffffffffffffffffffffffffffffffffffffffff1916149050919050565b5f7f52c63247e1f47db19d5ce0460030c497f067ca4cebf71ba98eeadabe20bace00905090565b5f612d776001612d35610f6e565b612d3f9190614edb565b612d476130cd565b600a612d53919061525b565b612d5b61115e565b612d659190614edb565b8486613875909392919063ffffffff16565b905092915050565b5f33905090565b612d9383838360016138c2565b505050565b5f612de8612da46130cd565b600a612db0919061525b565b612db861115e565b612dc29190614edb565b6001612dcc610f6e565b612dd69190614edb565b8486613875909392919063ffffffff16565b905092915050565b5f612dfb848461280c565b90507fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff811015612e7d5781811015612e6e578281836040517ffb8f41b2000000000000000000000000000000000000000000000000000000008152600401612e65939291906152a5565b60405180910390fd5b612e7c84848484035f6138c2565b5b50505050565b5f73ffffffffffffffffffffffffffffffffffffffff168373ffffffffffffffffffffffffffffffffffffffff1603612ef3575f6040517f96c6fd1e000000000000000000000000000000000000000000000000000000008152600401612eea9190614a73565b60405180910390fd5b5f73ffffffffffffffffffffffffffffffffffffffff168273ffffffffffffffffffffffffffffffffffffffff1603612f63575f6040517fec442f05000000000000000000000000000000000000000000000000000000008152600401612f5a9190614a73565b60405180910390fd5b612f6e838383613a9f565b505050565b5f7f02dd7bc7dec4dceedda775e58dd541e08a116c6c53815c0bd028192f7b626800905090565b612fab81612fa6612d7f565b613cce565b50565b5f5f612fb8612f73565b9050612fc484846121d2565b61309b576001815f015f8681526020019081526020015f205f015f8573ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205f6101000a81548160ff021916908315150217905550613037612d7f565b73ffffffffffffffffffffffffffffffffffffffff168373ffffffffffffffffffffffffffffffffffffffff16857f2f8788117e7eff1d82e926ec794901d17c78024a50270940304540a733656f0d60405160405180910390a460019150506130a0565b5f9150505b92915050565b5f7f0773e532dfede91f04b12a73d3d2acd361424f41f76b4fb79f090161e36b4e00905090565b5f5f905090565b5f5f6130de612f73565b90506130ea84846121d2565b156131c1575f815f015f8681526020019081526020015f205f015f8573ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205f6101000a81548160ff02191690831515021790555061315d612d7f565b73ffffffffffffffffffffffffffffffffffffffff168373ffffffffffffffffffffffffffffffffffffffff16857ff6391f5c32d9c69d2a47ea670b442974b53935d1edc7fd64eb21e047a839171b60405160405180910390a460019150506131c6565b5f9150505b92915050565b6131d4613d1f565b5f6131dd6134d5565b90505f815f015f6101000a81548160ff0219169083151502179055507f5db9ee0a495bf2e6ff9c91a7834c1ba4fdd244a5e8aa4e537bd38aeae4b073aa613222612d7f565b60405161322f9190614a73565b60405180910390a150565b7f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff163073ffffffffffffffffffffffffffffffffffffffff1614806132e757507f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff166132ce613d5f565b73ffffffffffffffffffffffffffffffffffffffff1614155b1561331e576040517fe07c8dba00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b565b5f5f1b61332c81612f9a565b5050565b8173ffffffffffffffffffffffffffffffffffffffff166352d1902d6040518163ffffffff1660e01b8152600401602060405180830381865afa92505050801561339857506040513d601f19601f8201168201806040525081019061339591906152ee565b60015b6133d957816040517f4c9c8ce30000000000000000000000000000000000000000000000000000000081526004016133d09190614a73565b60405180910390fd5b7f360894a13ba1a3210667c828492db98dca3e2076cc3735a920a3ca505d382bbc5f1b811461343f57806040517faa1d49a400000000000000000000000000000000000000000000000000000000815260040161343691906149d9565b60405180910390fd5b6134498383613db2565b505050565b7f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff163073ffffffffffffffffffffffffffffffffffffffff16146134d3576040517fe07c8dba00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b565b5f7fcd5ed15c6e187e77e9aee88184c21f4f2182ab5827cb3b7e07fbedcd63f03300905090565b613504611447565b1561353b576040517fd93c066500000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b565b613545613e24565b613568600161355a613555613e65565b613e8e565b613e9790919063ffffffff16565b565b5f5f61357583611333565b9050808411156135c0578284826040517f79012fb20000000000000000000000000000000000000000000000000000000081526004016135b7939291906152a5565b60405180910390fd5b5f6135ca856129b6565b90506135df6135d7612d7f565b858784613e9e565b809250505092915050565b61360c5f6135fe6135f9613e65565b613e8e565b613e9790919063ffffffff16565b565b5f5f613618613f1f565b90508091505090565b613629613f48565b6136338282613f88565b5050565b61363f613f48565b61364881613fc4565b50565b613653613f48565b565b61365d613f48565b565b61366d848484846001614056565b6136ae57836040517f5274afe70000000000000000000000000000000000000000000000000000000081526004016136a59190614a73565b60405180910390fd5b50505050565b6136c183838360016140c7565b61370257826040517f5274afe70000000000000000000000000000000000000000000000000000000081526004016136f99190614a73565b60405180910390fd5b505050565b61370f6134fc565b5f6137186134d5565b90506001815f015f6101000a81548160ff0219169083151502179055507f62e78cea01bee320cd4e420270b5ea74000d11b0c9f74754ebdbfc544b05a25861375e612d7f565b60405161376b9190614a73565b60405180910390a150565b5f5f613781836126ee565b9050808411156137cc578284826040517f284ff6670000000000000000000000000000000000000000000000000000000081526004016137c3939291906152a5565b60405180910390fd5b5f6137d68561266f565b90506137eb6137e3612d7f565b858388613e9e565b809250505092915050565b5f73ffffffffffffffffffffffffffffffffffffffff168273ffffffffffffffffffffffffffffffffffffffff1603613866575f6040517f96c6fd1e00000000000000000000000000000000000000000000000000000000815260040161385d9190614a73565b60405180910390fd5b613871825f83613a9f565b5050565b5f6138a361388283614129565b801561389e57505f8480613899576138986150a8565b5b868809115b614156565b6138ae868686614161565b6138b89190614edb565b9050949350505050565b5f6138cb612d00565b90505f73ffffffffffffffffffffffffffffffffffffffff168573ffffffffffffffffffffffffffffffffffffffff160361393d575f6040517fe602df050000000000000000000000000000000000000000000000000000000081526004016139349190614a73565b60405180910390fd5b5f73ffffffffffffffffffffffffffffffffffffffff168473ffffffffffffffffffffffffffffffffffffffff16036139ad575f6040517f94280d620000000000000000000000000000000000000000000000000000000081526004016139a49190614a73565b60405180910390fd5b82816001015f8773ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205f8673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f20819055508115613a98578373ffffffffffffffffffffffffffffffffffffffff168573ffffffffffffffffffffffffffffffffffffffff167f8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b92585604051613a8f91906146c2565b60405180910390a35b5050505050565b5f613aa8612d00565b90505f73ffffffffffffffffffffffffffffffffffffffff168473ffffffffffffffffffffffffffffffffffffffff1603613afc5781816002015f828254613af09190614edb565b92505081905550613bce565b5f815f015f8673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f2054905082811015613b87578481846040517fe450d38c000000000000000000000000000000000000000000000000000000008152600401613b7e939291906152a5565b60405180910390fd5b828103825f015f8773ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f2081905550505b5f73ffffffffffffffffffffffffffffffffffffffff168373ffffffffffffffffffffffffffffffffffffffff1603613c175781816002015f8282540392505081905550613c63565b81815f015f8573ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205f82825401925050819055505b8273ffffffffffffffffffffffffffffffffffffffff168473ffffffffffffffffffffffffffffffffffffffff167fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef84604051613cc091906146c2565b60405180910390a350505050565b613cd882826121d2565b613d1b5780826040517fe2517d3f000000000000000000000000000000000000000000000000000000008152600401613d12929190615319565b60405180910390fd5b5050565b613d27611447565b613d5d576040517f8dfc202b00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b565b5f613d8b7f360894a13ba1a3210667c828492db98dca3e2076cc3735a920a3ca505d382bbc5f1b614240565b5f015f9054906101000a900473ffffffffffffffffffffffffffffffffffffffff16905090565b613dbb82614249565b8173ffffffffffffffffffffffffffffffffffffffff167fbc7cd75a20ee27fd9adebab32041f755214dbc6bffa90cc0225b39da2e5c2d3b60405160405180910390a25f81511115613e1757613e118282614312565b50613e20565b613e1f614403565b5b5050565b613e2c61443f565b15613e63576040517f3ee5aeb500000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b565b5f7f9b779b17422d0df92223018b32b4d1fa46e071723d6817e2486d003becc55f005f1b905090565b5f819050919050565b80825d5050565b613ea8848361445d565b613eb28382614474565b8273ffffffffffffffffffffffffffffffffffffffff168473ffffffffffffffffffffffffffffffffffffffff167fdcbc1c05240f31ff3ad067ef1ee35ce4997762752e3a095284754544f4c709d78484604051613f11929190614f9f565b60405180910390a350505050565b5f7ff0c57e16840df040f15088dc2f81fe391c3923bec73e23a9662efc9c229c6a005f1b905090565b613f506144f3565b613f86576040517fd7e6bcf800000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b565b613f90613f48565b5f613f99612d00565b905082816003019081613fac91906154d7565b5081816004019081613fbe91906154d7565b50505050565b613fcc613f48565b5f613fd56130a6565b90505f5f613fe284614511565b9150915081613ff2576012613ff4565b805b835f0160146101000a81548160ff021916908360ff16021790555083835f015f6101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff16021790555050505050565b5f5f6323b872dd60e01b9050604051815f525f1960601c87166004525f1960601c86166024528460445260205f60645f5f8c5af1925060015f511483166140b45783831516156140a8573d5f823e3d81fd5b5f883b113d1516831692505b806040525f606052505095945050505050565b5f5f63a9059cbb60e01b9050604051815f525f1960601c86166004528460245260205f60445f5f8b5af1925060015f5114831661411b57838315161561410f573d5f823e3d81fd5b5f873b113d1516831692505b806040525050949350505050565b5f60016002836003811115614141576141406155a6565b5b61414b91906155d3565b60ff16149050919050565b5f8115159050919050565b5f5f5f61416e86866145c5565b915091505f820361419357838181614189576141886150a8565b5b0492505050614239565b8184116141b2576141b16141ac5f8614601260116145e2565b6145fb565b5b5f8486880990508181118303925080820391505f855f038616905080860495508083049250600181825f0304019050808402831792505f600287600302189050808702600203810290508087026002038102905080870260020381029050808702600203810290508087026002038102905080870260020381029050808402955050505050505b9392505050565b5f819050919050565b5f8173ffffffffffffffffffffffffffffffffffffffff163b036142a457806040517f4c9c8ce300000000000000000000000000000000000000000000000000000000815260040161429b9190614a73565b60405180910390fd5b806142d07f360894a13ba1a3210667c828492db98dca3e2076cc3735a920a3ca505d382bbc5f1b614240565b5f015f6101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff16021790555050565b60605f61431f848461460c565b905080801561435557505f614332614620565b118061435457505f8473ffffffffffffffffffffffffffffffffffffffff163b115b5b1561436a57614362614627565b9150506143fd565b80156143ad57836040517f9996b3150000000000000000000000000000000000000000000000000000000081526004016143a49190614a73565b60405180910390fd5b5f6143b6614620565b11156143c9576143c4614644565b6143fb565b6040517fd6bda27500000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b505b92915050565b5f34111561443d576040517fb398979f00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b565b5f61445861445361444e613e65565b613e8e565b61464f565b905090565b6144706144686112c3565b83308461365f565b5050565b5f73ffffffffffffffffffffffffffffffffffffffff168273ffffffffffffffffffffffffffffffffffffffff16036144e4575f6040517fec442f050000000000000000000000000000000000000000000000000000000081526004016144db9190614a73565b60405180910390fd5b6144ef5f8383613a9f565b5050565b5f6144fc61360e565b5f0160089054906101000a900460ff16905090565b5f5f5f61451c614659565b90505f5f6145748660405160240160405160208183030381529060405263313ce56760e01b6020820180517bffffffffffffffffffffffffffffffffffffffffffffffffffffffff8381831617835250505050614662565b509150915061458283614683565b81801561459757506020614594614620565b10155b80156145a8575060ff8016815f1c11155b6145b3575f5f6145b9565b6001815f1c5b94509450505050915091565b5f5f5f198385098385029150818110828203039250509250929050565b5f6145ec84614156565b82841802821890509392505050565b634e487b715f52806020526024601cfd5b5f5f5f835160208501865af4905092915050565b5f3d905090565b606060405190503d81523d5f602083013e3d602001810160405290565b6040513d5f823e3d81fd5b5f815c9050919050565b5f604051905090565b5f5f5f60405f855160208701885afa92505f51915060205190509250925092565b8060405250565b60405180606001604052805f81526020015f81526020015f151581525090565b5f819050919050565b6146bc816146aa565b82525050565b5f6020820190506146d55f8301846146b3565b92915050565b5f604051905090565b5f5ffd5b5f5ffd5b5f7fffffffff0000000000000000000000000000000000000000000000000000000082169050919050565b614720816146ec565b811461472a575f5ffd5b50565b5f8135905061473b81614717565b92915050565b5f60208284031215614756576147556146e4565b5b5f6147638482850161472d565b91505092915050565b5f8115159050919050565b6147808161476c565b82525050565b5f6020820190506147995f830184614777565b92915050565b5f81519050919050565b5f82825260208201905092915050565b8281835e5f83830152505050565b5f601f19601f8301169050919050565b5f6147e18261479f565b6147eb81856147a9565b93506147fb8185602086016147b9565b614804816147c7565b840191505092915050565b5f6020820190508181035f83015261482781846147d7565b905092915050565b614838816146aa565b8114614842575f5ffd5b50565b5f813590506148538161482f565b92915050565b5f6020828403121561486e5761486d6146e4565b5b5f61487b84828501614845565b91505092915050565b5f73ffffffffffffffffffffffffffffffffffffffff82169050919050565b5f6148ad82614884565b9050919050565b6148bd816148a3565b81146148c7575f5ffd5b50565b5f813590506148d8816148b4565b92915050565b5f5f604083850312156148f4576148f36146e4565b5b5f614901858286016148ca565b925050602061491285828601614845565b9150509250929050565b5f5f5f60608486031215614933576149326146e4565b5b5f614940868287016148ca565b9350506020614951868287016148ca565b925050604061496286828701614845565b9150509250925092565b5f819050919050565b61497e8161496c565b8114614988575f5ffd5b50565b5f8135905061499981614975565b92915050565b5f602082840312156149b4576149b36146e4565b5b5f6149c18482850161498b565b91505092915050565b6149d38161496c565b82525050565b5f6020820190506149ec5f8301846149ca565b92915050565b5f5f60408385031215614a0857614a076146e4565b5b5f614a158582860161498b565b9250506020614a26858286016148ca565b9150509250929050565b5f60ff82169050919050565b614a4581614a30565b82525050565b5f602082019050614a5e5f830184614a3c565b92915050565b614a6d816148a3565b82525050565b5f602082019050614a865f830184614a64565b92915050565b5f60208284031215614aa157614aa06146e4565b5b5f614aae848285016148ca565b91505092915050565b5f5f60408385031215614acd57614acc6146e4565b5b5f614ada85828601614845565b9250506020614aeb85828601614845565b9150509250929050565b5f5ffd5b5f5ffd5b7f4e487b71000000000000000000000000000000000000000000000000000000005f52604160045260245ffd5b614b33826147c7565b810181811067ffffffffffffffff82111715614b5257614b51614afd565b5b80604052505050565b5f614b646146db565b9050614b708282614b2a565b919050565b5f67ffffffffffffffff821115614b8f57614b8e614afd565b5b614b98826147c7565b9050602081019050919050565b828183375f83830152505050565b5f614bc5614bc084614b75565b614b5b565b905082815260208101848484011115614be157614be0614af9565b5b614bec848285614ba5565b509392505050565b5f82601f830112614c0857614c07614af5565b5b8135614c18848260208601614bb3565b91505092915050565b5f5f60408385031215614c3757614c366146e4565b5b5f614c44858286016148ca565b925050602083013567ffffffffffffffff811115614c6557614c646146e8565b5b614c7185828601614bf4565b9150509250929050565b5f5f60408385031215614c9157614c906146e4565b5b5f614c9e85828601614845565b9250506020614caf858286016148ca565b9150509250929050565b5f5f5f5f5f5f5f5f610100898b031215614cd657614cd56146e4565b5b5f614ce38b828c016148ca565b9850506020614cf48b828c016148ca565b9750506040614d058b828c016148ca565b9650506060614d168b828c016148ca565b9550506080614d278b828c016148ca565b94505060a0614d388b828c01614845565b93505060c0614d498b828c01614845565b92505060e0614d5a8b828c01614845565b9150509295985092959890939650565b614d73816146aa565b82525050565b614d828161476c565b82525050565b606082015f820151614d9c5f850182614d6a565b506020820151614daf6020850182614d6a565b506040820151614dc26040850182614d79565b50505050565b5f606082019050614ddb5f830184614d88565b92915050565b5f5f5f60608486031215614df857614df76146e4565b5b5f614e0586828701614845565b9350506020614e16868287016148ca565b9250506040614e27868287016148ca565b9150509250925092565b5f5f60408385031215614e4757614e466146e4565b5b5f614e54858286016148ca565b9250506020614e65858286016148ca565b9150509250929050565b5f81519050614e7d8161482f565b92915050565b5f60208284031215614e9857614e976146e4565b5b5f614ea584828501614e6f565b91505092915050565b7f4e487b71000000000000000000000000000000000000000000000000000000005f52601160045260245ffd5b5f614ee5826146aa565b9150614ef0836146aa565b9250828201905080821115614f0857614f07614eae565b5b92915050565b7f4e487b71000000000000000000000000000000000000000000000000000000005f52602260045260245ffd5b5f6002820490506001821680614f5257607f821691505b602082108103614f6557614f64614f0e565b5b50919050565b5f614f7582614a30565b9150614f8083614a30565b9250828201905060ff811115614f9957614f98614eae565b5b92915050565b5f604082019050614fb25f8301856146b3565b614fbf60208301846146b3565b9392505050565b5f819050919050565b5f67ffffffffffffffff82169050919050565b5f819050919050565b5f615005615000614ffb84614fc6565b614fe2565b614fcf565b9050919050565b61501581614feb565b82525050565b5f60208201905061502e5f83018461500c565b92915050565b5f61503e826146aa565b9150615049836146aa565b925082820390508181111561506157615060614eae565b5b92915050565b5f615071826146aa565b915061507c836146aa565b925082820261508a816146aa565b915082820484148315176150a1576150a0614eae565b5b5092915050565b7f4e487b71000000000000000000000000000000000000000000000000000000005f52601260045260245ffd5b5f6150df826146aa565b91506150ea836146aa565b9250826150fa576150f96150a8565b5b828204905092915050565b5f6040820190506151185f830185614a64565b6151256020830184614a64565b9392505050565b5f8160011c9050919050565b5f5f8291508390505b60018511156151815780860481111561515d5761515c614eae565b5b600185161561516c5780820291505b808102905061517a8561512c565b9450615141565b94509492505050565b5f826151995760019050615254565b816151a6575f9050615254565b81600181146151bc57600281146151c6576151f5565b6001915050615254565b60ff8411156151d8576151d7614eae565b5b8360020a9150848211156151ef576151ee614eae565b5b50615254565b5060208310610133831016604e8410600b841016171561522a5782820a90508381111561522557615224614eae565b5b615254565b6152378484846001615138565b9250905081840481111561524e5761524d614eae565b5b81810290505b9392505050565b5f615265826146aa565b915061527083614a30565b925061529d7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff848461518a565b905092915050565b5f6060820190506152b85f830186614a64565b6152c560208301856146b3565b6152d260408301846146b3565b949350505050565b5f815190506152e881614975565b92915050565b5f60208284031215615303576153026146e4565b5b5f615310848285016152da565b91505092915050565b5f60408201905061532c5f830185614a64565b61533960208301846149ca565b9392505050565b5f819050815f5260205f209050919050565b5f6020601f8301049050919050565b5f82821b905092915050565b5f6008830261539c7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff82615361565b6153a68683615361565b95508019841693508086168417925050509392505050565b5f6153d86153d36153ce846146aa565b614fe2565b6146aa565b9050919050565b5f819050919050565b6153f1836153be565b6154056153fd826153df565b84845461536d565b825550505050565b5f5f905090565b61541c61540d565b6154278184846153e8565b505050565b5b8181101561544a5761543f5f82615414565b60018101905061542d565b5050565b601f82111561548f5761546081615340565b61546984615352565b81016020851015615478578190505b61548c61548485615352565b83018261542c565b50505b505050565b5f82821c905092915050565b5f6154af5f1984600802615494565b1980831691505092915050565b5f6154c783836154a0565b9150826002028217905092915050565b6154e08261479f565b67ffffffffffffffff8111156154f9576154f8614afd565b5b6155038254614f3b565b61550e82828561544e565b5f60209050601f83116001811461553f575f841561552d578287015190505b61553785826154bc565b86555061559e565b601f19841661554d86615340565b5f5b828110156155745784890151825560018201915060208501945060208101905061554f565b86831015615591578489015161558d601f8916826154a0565b8355505b6001600288020188555050505b505050505050565b7f4e487b71000000000000000000000000000000000000000000000000000000005f52602160045260245ffd5b5f6155dd82614a30565b91506155e883614a30565b9250826155f8576155f76150a8565b5b82820690509291505056fea26469706673582212201d7e1d84da95a711bfbbd2210c16b9e72c07bf39453906be319f7a205be5361264736f6c634300081c0033",
}

// PolyVaultABI is the input ABI used to generate the binding from.
// Deprecated: Use PolyVaultMetaData.ABI instead.
var PolyVaultABI = PolyVaultMetaData.ABI

// PolyVaultBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use PolyVaultMetaData.Bin instead.
var PolyVaultBin = PolyVaultMetaData.Bin

// DeployPolyVault deploys a new Ethereum contract, binding an instance of PolyVault to it.
func DeployPolyVault(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *PolyVault, error) {
	parsed, err := PolyVaultMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(PolyVaultBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &PolyVault{PolyVaultCaller: PolyVaultCaller{contract: contract}, PolyVaultTransactor: PolyVaultTransactor{contract: contract}, PolyVaultFilterer: PolyVaultFilterer{contract: contract}}, nil
}

// PolyVault is an auto generated Go binding around an Ethereum contract.
type PolyVault struct {
	PolyVaultCaller     // Read-only binding to the contract
	PolyVaultTransactor // Write-only binding to the contract
	PolyVaultFilterer   // Log filterer for contract events
}

// PolyVaultCaller is an auto generated read-only Go binding around an Ethereum contract.
type PolyVaultCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// PolyVaultTransactor is an auto generated write-only Go binding around an Ethereum contract.
type PolyVaultTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// PolyVaultFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type PolyVaultFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// PolyVaultSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type PolyVaultSession struct {
	Contract     *PolyVault        // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// PolyVaultCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type PolyVaultCallerSession struct {
	Contract *PolyVaultCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts    // Call options to use throughout this session
}

// PolyVaultTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type PolyVaultTransactorSession struct {
	Contract     *PolyVaultTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts    // Transaction auth options to use throughout this session
}

// PolyVaultRaw is an auto generated low-level Go binding around an Ethereum contract.
type PolyVaultRaw struct {
	Contract *PolyVault // Generic contract binding to access the raw methods on
}

// PolyVaultCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type PolyVaultCallerRaw struct {
	Contract *PolyVaultCaller // Generic read-only contract binding to access the raw methods on
}

// PolyVaultTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type PolyVaultTransactorRaw struct {
	Contract *PolyVaultTransactor // Generic write-only contract binding to access the raw methods on
}

// NewPolyVault creates a new instance of PolyVault, bound to a specific deployed contract.
func NewPolyVault(address common.Address, backend bind.ContractBackend) (*PolyVault, error) {
	contract, err := bindPolyVault(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &PolyVault{PolyVaultCaller: PolyVaultCaller{contract: contract}, PolyVaultTransactor: PolyVaultTransactor{contract: contract}, PolyVaultFilterer: PolyVaultFilterer{contract: contract}}, nil
}

// NewPolyVaultCaller creates a new read-only instance of PolyVault, bound to a specific deployed contract.
func NewPolyVaultCaller(address common.Address, caller bind.ContractCaller) (*PolyVaultCaller, error) {
	contract, err := bindPolyVault(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &PolyVaultCaller{contract: contract}, nil
}

// NewPolyVaultTransactor creates a new write-only instance of PolyVault, bound to a specific deployed contract.
func NewPolyVaultTransactor(address common.Address, transactor bind.ContractTransactor) (*PolyVaultTransactor, error) {
	contract, err := bindPolyVault(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &PolyVaultTransactor{contract: contract}, nil
}

// NewPolyVaultFilterer creates a new log filterer instance of PolyVault, bound to a specific deployed contract.
func NewPolyVaultFilterer(address common.Address, filterer bind.ContractFilterer) (*PolyVaultFilterer, error) {
	contract, err := bindPolyVault(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &PolyVaultFilterer{contract: contract}, nil
}

// bindPolyVault binds a generic wrapper to an already deployed contract.
func bindPolyVault(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := PolyVaultMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_PolyVault *PolyVaultRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _PolyVault.Contract.PolyVaultCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_PolyVault *PolyVaultRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _PolyVault.Contract.PolyVaultTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_PolyVault *PolyVaultRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _PolyVault.Contract.PolyVaultTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_PolyVault *PolyVaultCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _PolyVault.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_PolyVault *PolyVaultTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _PolyVault.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_PolyVault *PolyVaultTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _PolyVault.Contract.contract.Transact(opts, method, params...)
}

// BASISPOINTS is a free data retrieval call binding the contract method 0xe1f1c4a7.
//
// Solidity: function BASIS_POINTS() view returns(uint256)
func (_PolyVault *PolyVaultCaller) BASISPOINTS(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _PolyVault.contract.Call(opts, &out, "BASIS_POINTS")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// BASISPOINTS is a free data retrieval call binding the contract method 0xe1f1c4a7.
//
// Solidity: function BASIS_POINTS() view returns(uint256)
func (_PolyVault *PolyVaultSession) BASISPOINTS() (*big.Int, error) {
	return _PolyVault.Contract.BASISPOINTS(&_PolyVault.CallOpts)
}

// BASISPOINTS is a free data retrieval call binding the contract method 0xe1f1c4a7.
//
// Solidity: function BASIS_POINTS() view returns(uint256)
func (_PolyVault *PolyVaultCallerSession) BASISPOINTS() (*big.Int, error) {
	return _PolyVault.Contract.BASISPOINTS(&_PolyVault.CallOpts)
}

// DEFAULTADMINROLE is a free data retrieval call binding the contract method 0xa217fddf.
//
// Solidity: function DEFAULT_ADMIN_ROLE() view returns(bytes32)
func (_PolyVault *PolyVaultCaller) DEFAULTADMINROLE(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _PolyVault.contract.Call(opts, &out, "DEFAULT_ADMIN_ROLE")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// DEFAULTADMINROLE is a free data retrieval call binding the contract method 0xa217fddf.
//
// Solidity: function DEFAULT_ADMIN_ROLE() view returns(bytes32)
func (_PolyVault *PolyVaultSession) DEFAULTADMINROLE() ([32]byte, error) {
	return _PolyVault.Contract.DEFAULTADMINROLE(&_PolyVault.CallOpts)
}

// DEFAULTADMINROLE is a free data retrieval call binding the contract method 0xa217fddf.
//
// Solidity: function DEFAULT_ADMIN_ROLE() view returns(bytes32)
func (_PolyVault *PolyVaultCallerSession) DEFAULTADMINROLE() ([32]byte, error) {
	return _PolyVault.Contract.DEFAULTADMINROLE(&_PolyVault.CallOpts)
}

// GUARDIANROLE is a free data retrieval call binding the contract method 0x24ea54f4.
//
// Solidity: function GUARDIAN_ROLE() view returns(bytes32)
func (_PolyVault *PolyVaultCaller) GUARDIANROLE(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _PolyVault.contract.Call(opts, &out, "GUARDIAN_ROLE")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// GUARDIANROLE is a free data retrieval call binding the contract method 0x24ea54f4.
//
// Solidity: function GUARDIAN_ROLE() view returns(bytes32)
func (_PolyVault *PolyVaultSession) GUARDIANROLE() ([32]byte, error) {
	return _PolyVault.Contract.GUARDIANROLE(&_PolyVault.CallOpts)
}

// GUARDIANROLE is a free data retrieval call binding the contract method 0x24ea54f4.
//
// Solidity: function GUARDIAN_ROLE() view returns(bytes32)
func (_PolyVault *PolyVaultCallerSession) GUARDIANROLE() ([32]byte, error) {
	return _PolyVault.Contract.GUARDIANROLE(&_PolyVault.CallOpts)
}

// MAXPERFORMANCEFEE is a free data retrieval call binding the contract method 0xbdca9165.
//
// Solidity: function MAX_PERFORMANCE_FEE() view returns(uint256)
func (_PolyVault *PolyVaultCaller) MAXPERFORMANCEFEE(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _PolyVault.contract.Call(opts, &out, "MAX_PERFORMANCE_FEE")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// MAXPERFORMANCEFEE is a free data retrieval call binding the contract method 0xbdca9165.
//
// Solidity: function MAX_PERFORMANCE_FEE() view returns(uint256)
func (_PolyVault *PolyVaultSession) MAXPERFORMANCEFEE() (*big.Int, error) {
	return _PolyVault.Contract.MAXPERFORMANCEFEE(&_PolyVault.CallOpts)
}

// MAXPERFORMANCEFEE is a free data retrieval call binding the contract method 0xbdca9165.
//
// Solidity: function MAX_PERFORMANCE_FEE() view returns(uint256)
func (_PolyVault *PolyVaultCallerSession) MAXPERFORMANCEFEE() (*big.Int, error) {
	return _PolyVault.Contract.MAXPERFORMANCEFEE(&_PolyVault.CallOpts)
}

// MAXWITHDRAWALDELAY is a free data retrieval call binding the contract method 0xa238f9df.
//
// Solidity: function MAX_WITHDRAWAL_DELAY() view returns(uint256)
func (_PolyVault *PolyVaultCaller) MAXWITHDRAWALDELAY(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _PolyVault.contract.Call(opts, &out, "MAX_WITHDRAWAL_DELAY")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// MAXWITHDRAWALDELAY is a free data retrieval call binding the contract method 0xa238f9df.
//
// Solidity: function MAX_WITHDRAWAL_DELAY() view returns(uint256)
func (_PolyVault *PolyVaultSession) MAXWITHDRAWALDELAY() (*big.Int, error) {
	return _PolyVault.Contract.MAXWITHDRAWALDELAY(&_PolyVault.CallOpts)
}

// MAXWITHDRAWALDELAY is a free data retrieval call binding the contract method 0xa238f9df.
//
// Solidity: function MAX_WITHDRAWAL_DELAY() view returns(uint256)
func (_PolyVault *PolyVaultCallerSession) MAXWITHDRAWALDELAY() (*big.Int, error) {
	return _PolyVault.Contract.MAXWITHDRAWALDELAY(&_PolyVault.CallOpts)
}

// MINWITHDRAWALDELAY is a free data retrieval call binding the contract method 0x4a5f2b5d.
//
// Solidity: function MIN_WITHDRAWAL_DELAY() view returns(uint256)
func (_PolyVault *PolyVaultCaller) MINWITHDRAWALDELAY(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _PolyVault.contract.Call(opts, &out, "MIN_WITHDRAWAL_DELAY")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// MINWITHDRAWALDELAY is a free data retrieval call binding the contract method 0x4a5f2b5d.
//
// Solidity: function MIN_WITHDRAWAL_DELAY() view returns(uint256)
func (_PolyVault *PolyVaultSession) MINWITHDRAWALDELAY() (*big.Int, error) {
	return _PolyVault.Contract.MINWITHDRAWALDELAY(&_PolyVault.CallOpts)
}

// MINWITHDRAWALDELAY is a free data retrieval call binding the contract method 0x4a5f2b5d.
//
// Solidity: function MIN_WITHDRAWAL_DELAY() view returns(uint256)
func (_PolyVault *PolyVaultCallerSession) MINWITHDRAWALDELAY() (*big.Int, error) {
	return _PolyVault.Contract.MINWITHDRAWALDELAY(&_PolyVault.CallOpts)
}

// STRATEGISTROLE is a free data retrieval call binding the contract method 0xa378a324.
//
// Solidity: function STRATEGIST_ROLE() view returns(bytes32)
func (_PolyVault *PolyVaultCaller) STRATEGISTROLE(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _PolyVault.contract.Call(opts, &out, "STRATEGIST_ROLE")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// STRATEGISTROLE is a free data retrieval call binding the contract method 0xa378a324.
//
// Solidity: function STRATEGIST_ROLE() view returns(bytes32)
func (_PolyVault *PolyVaultSession) STRATEGISTROLE() ([32]byte, error) {
	return _PolyVault.Contract.STRATEGISTROLE(&_PolyVault.CallOpts)
}

// STRATEGISTROLE is a free data retrieval call binding the contract method 0xa378a324.
//
// Solidity: function STRATEGIST_ROLE() view returns(bytes32)
func (_PolyVault *PolyVaultCallerSession) STRATEGISTROLE() ([32]byte, error) {
	return _PolyVault.Contract.STRATEGISTROLE(&_PolyVault.CallOpts)
}

// UPGRADEINTERFACEVERSION is a free data retrieval call binding the contract method 0xad3cb1cc.
//
// Solidity: function UPGRADE_INTERFACE_VERSION() view returns(string)
func (_PolyVault *PolyVaultCaller) UPGRADEINTERFACEVERSION(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _PolyVault.contract.Call(opts, &out, "UPGRADE_INTERFACE_VERSION")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// UPGRADEINTERFACEVERSION is a free data retrieval call binding the contract method 0xad3cb1cc.
//
// Solidity: function UPGRADE_INTERFACE_VERSION() view returns(string)
func (_PolyVault *PolyVaultSession) UPGRADEINTERFACEVERSION() (string, error) {
	return _PolyVault.Contract.UPGRADEINTERFACEVERSION(&_PolyVault.CallOpts)
}

// UPGRADEINTERFACEVERSION is a free data retrieval call binding the contract method 0xad3cb1cc.
//
// Solidity: function UPGRADE_INTERFACE_VERSION() view returns(string)
func (_PolyVault *PolyVaultCallerSession) UPGRADEINTERFACEVERSION() (string, error) {
	return _PolyVault.Contract.UPGRADEINTERFACEVERSION(&_PolyVault.CallOpts)
}

// Allowance is a free data retrieval call binding the contract method 0xdd62ed3e.
//
// Solidity: function allowance(address owner, address spender) view returns(uint256)
func (_PolyVault *PolyVaultCaller) Allowance(opts *bind.CallOpts, owner common.Address, spender common.Address) (*big.Int, error) {
	var out []interface{}
	err := _PolyVault.contract.Call(opts, &out, "allowance", owner, spender)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Allowance is a free data retrieval call binding the contract method 0xdd62ed3e.
//
// Solidity: function allowance(address owner, address spender) view returns(uint256)
func (_PolyVault *PolyVaultSession) Allowance(owner common.Address, spender common.Address) (*big.Int, error) {
	return _PolyVault.Contract.Allowance(&_PolyVault.CallOpts, owner, spender)
}

// Allowance is a free data retrieval call binding the contract method 0xdd62ed3e.
//
// Solidity: function allowance(address owner, address spender) view returns(uint256)
func (_PolyVault *PolyVaultCallerSession) Allowance(owner common.Address, spender common.Address) (*big.Int, error) {
	return _PolyVault.Contract.Allowance(&_PolyVault.CallOpts, owner, spender)
}

// Asset is a free data retrieval call binding the contract method 0x38d52e0f.
//
// Solidity: function asset() view returns(address)
func (_PolyVault *PolyVaultCaller) Asset(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _PolyVault.contract.Call(opts, &out, "asset")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Asset is a free data retrieval call binding the contract method 0x38d52e0f.
//
// Solidity: function asset() view returns(address)
func (_PolyVault *PolyVaultSession) Asset() (common.Address, error) {
	return _PolyVault.Contract.Asset(&_PolyVault.CallOpts)
}

// Asset is a free data retrieval call binding the contract method 0x38d52e0f.
//
// Solidity: function asset() view returns(address)
func (_PolyVault *PolyVaultCallerSession) Asset() (common.Address, error) {
	return _PolyVault.Contract.Asset(&_PolyVault.CallOpts)
}

// AvailableBalance is a free data retrieval call binding the contract method 0xab2f0e51.
//
// Solidity: function availableBalance() view returns(uint256)
func (_PolyVault *PolyVaultCaller) AvailableBalance(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _PolyVault.contract.Call(opts, &out, "availableBalance")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// AvailableBalance is a free data retrieval call binding the contract method 0xab2f0e51.
//
// Solidity: function availableBalance() view returns(uint256)
func (_PolyVault *PolyVaultSession) AvailableBalance() (*big.Int, error) {
	return _PolyVault.Contract.AvailableBalance(&_PolyVault.CallOpts)
}

// AvailableBalance is a free data retrieval call binding the contract method 0xab2f0e51.
//
// Solidity: function availableBalance() view returns(uint256)
func (_PolyVault *PolyVaultCallerSession) AvailableBalance() (*big.Int, error) {
	return _PolyVault.Contract.AvailableBalance(&_PolyVault.CallOpts)
}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address account) view returns(uint256)
func (_PolyVault *PolyVaultCaller) BalanceOf(opts *bind.CallOpts, account common.Address) (*big.Int, error) {
	var out []interface{}
	err := _PolyVault.contract.Call(opts, &out, "balanceOf", account)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address account) view returns(uint256)
func (_PolyVault *PolyVaultSession) BalanceOf(account common.Address) (*big.Int, error) {
	return _PolyVault.Contract.BalanceOf(&_PolyVault.CallOpts, account)
}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address account) view returns(uint256)
func (_PolyVault *PolyVaultCallerSession) BalanceOf(account common.Address) (*big.Int, error) {
	return _PolyVault.Contract.BalanceOf(&_PolyVault.CallOpts, account)
}

// ConvertToAssets is a free data retrieval call binding the contract method 0x07a2d13a.
//
// Solidity: function convertToAssets(uint256 shares) view returns(uint256)
func (_PolyVault *PolyVaultCaller) ConvertToAssets(opts *bind.CallOpts, shares *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _PolyVault.contract.Call(opts, &out, "convertToAssets", shares)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// ConvertToAssets is a free data retrieval call binding the contract method 0x07a2d13a.
//
// Solidity: function convertToAssets(uint256 shares) view returns(uint256)
func (_PolyVault *PolyVaultSession) ConvertToAssets(shares *big.Int) (*big.Int, error) {
	return _PolyVault.Contract.ConvertToAssets(&_PolyVault.CallOpts, shares)
}

// ConvertToAssets is a free data retrieval call binding the contract method 0x07a2d13a.
//
// Solidity: function convertToAssets(uint256 shares) view returns(uint256)
func (_PolyVault *PolyVaultCallerSession) ConvertToAssets(shares *big.Int) (*big.Int, error) {
	return _PolyVault.Contract.ConvertToAssets(&_PolyVault.CallOpts, shares)
}

// ConvertToShares is a free data retrieval call binding the contract method 0xc6e6f592.
//
// Solidity: function convertToShares(uint256 assets) view returns(uint256)
func (_PolyVault *PolyVaultCaller) ConvertToShares(opts *bind.CallOpts, assets *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _PolyVault.contract.Call(opts, &out, "convertToShares", assets)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// ConvertToShares is a free data retrieval call binding the contract method 0xc6e6f592.
//
// Solidity: function convertToShares(uint256 assets) view returns(uint256)
func (_PolyVault *PolyVaultSession) ConvertToShares(assets *big.Int) (*big.Int, error) {
	return _PolyVault.Contract.ConvertToShares(&_PolyVault.CallOpts, assets)
}

// ConvertToShares is a free data retrieval call binding the contract method 0xc6e6f592.
//
// Solidity: function convertToShares(uint256 assets) view returns(uint256)
func (_PolyVault *PolyVaultCallerSession) ConvertToShares(assets *big.Int) (*big.Int, error) {
	return _PolyVault.Contract.ConvertToShares(&_PolyVault.CallOpts, assets)
}

// Decimals is a free data retrieval call binding the contract method 0x313ce567.
//
// Solidity: function decimals() view returns(uint8)
func (_PolyVault *PolyVaultCaller) Decimals(opts *bind.CallOpts) (uint8, error) {
	var out []interface{}
	err := _PolyVault.contract.Call(opts, &out, "decimals")

	if err != nil {
		return *new(uint8), err
	}

	out0 := *abi.ConvertType(out[0], new(uint8)).(*uint8)

	return out0, err

}

// Decimals is a free data retrieval call binding the contract method 0x313ce567.
//
// Solidity: function decimals() view returns(uint8)
func (_PolyVault *PolyVaultSession) Decimals() (uint8, error) {
	return _PolyVault.Contract.Decimals(&_PolyVault.CallOpts)
}

// Decimals is a free data retrieval call binding the contract method 0x313ce567.
//
// Solidity: function decimals() view returns(uint8)
func (_PolyVault *PolyVaultCallerSession) Decimals() (uint8, error) {
	return _PolyVault.Contract.Decimals(&_PolyVault.CallOpts)
}

// FeeRecipient is a free data retrieval call binding the contract method 0x46904840.
//
// Solidity: function feeRecipient() view returns(address)
func (_PolyVault *PolyVaultCaller) FeeRecipient(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _PolyVault.contract.Call(opts, &out, "feeRecipient")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// FeeRecipient is a free data retrieval call binding the contract method 0x46904840.
//
// Solidity: function feeRecipient() view returns(address)
func (_PolyVault *PolyVaultSession) FeeRecipient() (common.Address, error) {
	return _PolyVault.Contract.FeeRecipient(&_PolyVault.CallOpts)
}

// FeeRecipient is a free data retrieval call binding the contract method 0x46904840.
//
// Solidity: function feeRecipient() view returns(address)
func (_PolyVault *PolyVaultCallerSession) FeeRecipient() (common.Address, error) {
	return _PolyVault.Contract.FeeRecipient(&_PolyVault.CallOpts)
}

// GetRoleAdmin is a free data retrieval call binding the contract method 0x248a9ca3.
//
// Solidity: function getRoleAdmin(bytes32 role) view returns(bytes32)
func (_PolyVault *PolyVaultCaller) GetRoleAdmin(opts *bind.CallOpts, role [32]byte) ([32]byte, error) {
	var out []interface{}
	err := _PolyVault.contract.Call(opts, &out, "getRoleAdmin", role)

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// GetRoleAdmin is a free data retrieval call binding the contract method 0x248a9ca3.
//
// Solidity: function getRoleAdmin(bytes32 role) view returns(bytes32)
func (_PolyVault *PolyVaultSession) GetRoleAdmin(role [32]byte) ([32]byte, error) {
	return _PolyVault.Contract.GetRoleAdmin(&_PolyVault.CallOpts, role)
}

// GetRoleAdmin is a free data retrieval call binding the contract method 0x248a9ca3.
//
// Solidity: function getRoleAdmin(bytes32 role) view returns(bytes32)
func (_PolyVault *PolyVaultCallerSession) GetRoleAdmin(role [32]byte) ([32]byte, error) {
	return _PolyVault.Contract.GetRoleAdmin(&_PolyVault.CallOpts, role)
}

// GetWithdrawalRequest is a free data retrieval call binding the contract method 0x8c661b5d.
//
// Solidity: function getWithdrawalRequest(address user) view returns((uint256,uint256,bool))
func (_PolyVault *PolyVaultCaller) GetWithdrawalRequest(opts *bind.CallOpts, user common.Address) (IPolyVaultWithdrawalRequest, error) {
	var out []interface{}
	err := _PolyVault.contract.Call(opts, &out, "getWithdrawalRequest", user)

	if err != nil {
		return *new(IPolyVaultWithdrawalRequest), err
	}

	out0 := *abi.ConvertType(out[0], new(IPolyVaultWithdrawalRequest)).(*IPolyVaultWithdrawalRequest)

	return out0, err

}

// GetWithdrawalRequest is a free data retrieval call binding the contract method 0x8c661b5d.
//
// Solidity: function getWithdrawalRequest(address user) view returns((uint256,uint256,bool))
func (_PolyVault *PolyVaultSession) GetWithdrawalRequest(user common.Address) (IPolyVaultWithdrawalRequest, error) {
	return _PolyVault.Contract.GetWithdrawalRequest(&_PolyVault.CallOpts, user)
}

// GetWithdrawalRequest is a free data retrieval call binding the contract method 0x8c661b5d.
//
// Solidity: function getWithdrawalRequest(address user) view returns((uint256,uint256,bool))
func (_PolyVault *PolyVaultCallerSession) GetWithdrawalRequest(user common.Address) (IPolyVaultWithdrawalRequest, error) {
	return _PolyVault.Contract.GetWithdrawalRequest(&_PolyVault.CallOpts, user)
}

// HasRole is a free data retrieval call binding the contract method 0x91d14854.
//
// Solidity: function hasRole(bytes32 role, address account) view returns(bool)
func (_PolyVault *PolyVaultCaller) HasRole(opts *bind.CallOpts, role [32]byte, account common.Address) (bool, error) {
	var out []interface{}
	err := _PolyVault.contract.Call(opts, &out, "hasRole", role, account)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// HasRole is a free data retrieval call binding the contract method 0x91d14854.
//
// Solidity: function hasRole(bytes32 role, address account) view returns(bool)
func (_PolyVault *PolyVaultSession) HasRole(role [32]byte, account common.Address) (bool, error) {
	return _PolyVault.Contract.HasRole(&_PolyVault.CallOpts, role, account)
}

// HasRole is a free data retrieval call binding the contract method 0x91d14854.
//
// Solidity: function hasRole(bytes32 role, address account) view returns(bool)
func (_PolyVault *PolyVaultCallerSession) HasRole(role [32]byte, account common.Address) (bool, error) {
	return _PolyVault.Contract.HasRole(&_PolyVault.CallOpts, role, account)
}

// MaxDeposit is a free data retrieval call binding the contract method 0x402d267d.
//
// Solidity: function maxDeposit(address ) view returns(uint256)
func (_PolyVault *PolyVaultCaller) MaxDeposit(opts *bind.CallOpts, arg0 common.Address) (*big.Int, error) {
	var out []interface{}
	err := _PolyVault.contract.Call(opts, &out, "maxDeposit", arg0)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// MaxDeposit is a free data retrieval call binding the contract method 0x402d267d.
//
// Solidity: function maxDeposit(address ) view returns(uint256)
func (_PolyVault *PolyVaultSession) MaxDeposit(arg0 common.Address) (*big.Int, error) {
	return _PolyVault.Contract.MaxDeposit(&_PolyVault.CallOpts, arg0)
}

// MaxDeposit is a free data retrieval call binding the contract method 0x402d267d.
//
// Solidity: function maxDeposit(address ) view returns(uint256)
func (_PolyVault *PolyVaultCallerSession) MaxDeposit(arg0 common.Address) (*big.Int, error) {
	return _PolyVault.Contract.MaxDeposit(&_PolyVault.CallOpts, arg0)
}

// MaxDeposit0 is a free data retrieval call binding the contract method 0x6083e59a.
//
// Solidity: function maxDeposit() view returns(uint256)
func (_PolyVault *PolyVaultCaller) MaxDeposit0(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _PolyVault.contract.Call(opts, &out, "maxDeposit0")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// MaxDeposit0 is a free data retrieval call binding the contract method 0x6083e59a.
//
// Solidity: function maxDeposit() view returns(uint256)
func (_PolyVault *PolyVaultSession) MaxDeposit0() (*big.Int, error) {
	return _PolyVault.Contract.MaxDeposit0(&_PolyVault.CallOpts)
}

// MaxDeposit0 is a free data retrieval call binding the contract method 0x6083e59a.
//
// Solidity: function maxDeposit() view returns(uint256)
func (_PolyVault *PolyVaultCallerSession) MaxDeposit0() (*big.Int, error) {
	return _PolyVault.Contract.MaxDeposit0(&_PolyVault.CallOpts)
}

// MaxMint is a free data retrieval call binding the contract method 0xc63d75b6.
//
// Solidity: function maxMint(address ) view returns(uint256)
func (_PolyVault *PolyVaultCaller) MaxMint(opts *bind.CallOpts, arg0 common.Address) (*big.Int, error) {
	var out []interface{}
	err := _PolyVault.contract.Call(opts, &out, "maxMint", arg0)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// MaxMint is a free data retrieval call binding the contract method 0xc63d75b6.
//
// Solidity: function maxMint(address ) view returns(uint256)
func (_PolyVault *PolyVaultSession) MaxMint(arg0 common.Address) (*big.Int, error) {
	return _PolyVault.Contract.MaxMint(&_PolyVault.CallOpts, arg0)
}

// MaxMint is a free data retrieval call binding the contract method 0xc63d75b6.
//
// Solidity: function maxMint(address ) view returns(uint256)
func (_PolyVault *PolyVaultCallerSession) MaxMint(arg0 common.Address) (*big.Int, error) {
	return _PolyVault.Contract.MaxMint(&_PolyVault.CallOpts, arg0)
}

// MaxRedeem is a free data retrieval call binding the contract method 0xd905777e.
//
// Solidity: function maxRedeem(address ) pure returns(uint256)
func (_PolyVault *PolyVaultCaller) MaxRedeem(opts *bind.CallOpts, arg0 common.Address) (*big.Int, error) {
	var out []interface{}
	err := _PolyVault.contract.Call(opts, &out, "maxRedeem", arg0)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// MaxRedeem is a free data retrieval call binding the contract method 0xd905777e.
//
// Solidity: function maxRedeem(address ) pure returns(uint256)
func (_PolyVault *PolyVaultSession) MaxRedeem(arg0 common.Address) (*big.Int, error) {
	return _PolyVault.Contract.MaxRedeem(&_PolyVault.CallOpts, arg0)
}

// MaxRedeem is a free data retrieval call binding the contract method 0xd905777e.
//
// Solidity: function maxRedeem(address ) pure returns(uint256)
func (_PolyVault *PolyVaultCallerSession) MaxRedeem(arg0 common.Address) (*big.Int, error) {
	return _PolyVault.Contract.MaxRedeem(&_PolyVault.CallOpts, arg0)
}

// MaxStrategyAllocation is a free data retrieval call binding the contract method 0xca1689d0.
//
// Solidity: function maxStrategyAllocation() view returns(uint256)
func (_PolyVault *PolyVaultCaller) MaxStrategyAllocation(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _PolyVault.contract.Call(opts, &out, "maxStrategyAllocation")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// MaxStrategyAllocation is a free data retrieval call binding the contract method 0xca1689d0.
//
// Solidity: function maxStrategyAllocation() view returns(uint256)
func (_PolyVault *PolyVaultSession) MaxStrategyAllocation() (*big.Int, error) {
	return _PolyVault.Contract.MaxStrategyAllocation(&_PolyVault.CallOpts)
}

// MaxStrategyAllocation is a free data retrieval call binding the contract method 0xca1689d0.
//
// Solidity: function maxStrategyAllocation() view returns(uint256)
func (_PolyVault *PolyVaultCallerSession) MaxStrategyAllocation() (*big.Int, error) {
	return _PolyVault.Contract.MaxStrategyAllocation(&_PolyVault.CallOpts)
}

// MaxWithdraw is a free data retrieval call binding the contract method 0xce96cb77.
//
// Solidity: function maxWithdraw(address ) pure returns(uint256)
func (_PolyVault *PolyVaultCaller) MaxWithdraw(opts *bind.CallOpts, arg0 common.Address) (*big.Int, error) {
	var out []interface{}
	err := _PolyVault.contract.Call(opts, &out, "maxWithdraw", arg0)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// MaxWithdraw is a free data retrieval call binding the contract method 0xce96cb77.
//
// Solidity: function maxWithdraw(address ) pure returns(uint256)
func (_PolyVault *PolyVaultSession) MaxWithdraw(arg0 common.Address) (*big.Int, error) {
	return _PolyVault.Contract.MaxWithdraw(&_PolyVault.CallOpts, arg0)
}

// MaxWithdraw is a free data retrieval call binding the contract method 0xce96cb77.
//
// Solidity: function maxWithdraw(address ) pure returns(uint256)
func (_PolyVault *PolyVaultCallerSession) MaxWithdraw(arg0 common.Address) (*big.Int, error) {
	return _PolyVault.Contract.MaxWithdraw(&_PolyVault.CallOpts, arg0)
}

// MinDeposit is a free data retrieval call binding the contract method 0x41b3d185.
//
// Solidity: function minDeposit() view returns(uint256)
func (_PolyVault *PolyVaultCaller) MinDeposit(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _PolyVault.contract.Call(opts, &out, "minDeposit")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// MinDeposit is a free data retrieval call binding the contract method 0x41b3d185.
//
// Solidity: function minDeposit() view returns(uint256)
func (_PolyVault *PolyVaultSession) MinDeposit() (*big.Int, error) {
	return _PolyVault.Contract.MinDeposit(&_PolyVault.CallOpts)
}

// MinDeposit is a free data retrieval call binding the contract method 0x41b3d185.
//
// Solidity: function minDeposit() view returns(uint256)
func (_PolyVault *PolyVaultCallerSession) MinDeposit() (*big.Int, error) {
	return _PolyVault.Contract.MinDeposit(&_PolyVault.CallOpts)
}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() view returns(string)
func (_PolyVault *PolyVaultCaller) Name(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _PolyVault.contract.Call(opts, &out, "name")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() view returns(string)
func (_PolyVault *PolyVaultSession) Name() (string, error) {
	return _PolyVault.Contract.Name(&_PolyVault.CallOpts)
}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() view returns(string)
func (_PolyVault *PolyVaultCallerSession) Name() (string, error) {
	return _PolyVault.Contract.Name(&_PolyVault.CallOpts)
}

// Paused is a free data retrieval call binding the contract method 0x5c975abb.
//
// Solidity: function paused() view returns(bool)
func (_PolyVault *PolyVaultCaller) Paused(opts *bind.CallOpts) (bool, error) {
	var out []interface{}
	err := _PolyVault.contract.Call(opts, &out, "paused")

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// Paused is a free data retrieval call binding the contract method 0x5c975abb.
//
// Solidity: function paused() view returns(bool)
func (_PolyVault *PolyVaultSession) Paused() (bool, error) {
	return _PolyVault.Contract.Paused(&_PolyVault.CallOpts)
}

// Paused is a free data retrieval call binding the contract method 0x5c975abb.
//
// Solidity: function paused() view returns(bool)
func (_PolyVault *PolyVaultCallerSession) Paused() (bool, error) {
	return _PolyVault.Contract.Paused(&_PolyVault.CallOpts)
}

// PerformanceFee is a free data retrieval call binding the contract method 0x87788782.
//
// Solidity: function performanceFee() view returns(uint256)
func (_PolyVault *PolyVaultCaller) PerformanceFee(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _PolyVault.contract.Call(opts, &out, "performanceFee")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// PerformanceFee is a free data retrieval call binding the contract method 0x87788782.
//
// Solidity: function performanceFee() view returns(uint256)
func (_PolyVault *PolyVaultSession) PerformanceFee() (*big.Int, error) {
	return _PolyVault.Contract.PerformanceFee(&_PolyVault.CallOpts)
}

// PerformanceFee is a free data retrieval call binding the contract method 0x87788782.
//
// Solidity: function performanceFee() view returns(uint256)
func (_PolyVault *PolyVaultCallerSession) PerformanceFee() (*big.Int, error) {
	return _PolyVault.Contract.PerformanceFee(&_PolyVault.CallOpts)
}

// PreviewDeposit is a free data retrieval call binding the contract method 0xef8b30f7.
//
// Solidity: function previewDeposit(uint256 assets) view returns(uint256)
func (_PolyVault *PolyVaultCaller) PreviewDeposit(opts *bind.CallOpts, assets *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _PolyVault.contract.Call(opts, &out, "previewDeposit", assets)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// PreviewDeposit is a free data retrieval call binding the contract method 0xef8b30f7.
//
// Solidity: function previewDeposit(uint256 assets) view returns(uint256)
func (_PolyVault *PolyVaultSession) PreviewDeposit(assets *big.Int) (*big.Int, error) {
	return _PolyVault.Contract.PreviewDeposit(&_PolyVault.CallOpts, assets)
}

// PreviewDeposit is a free data retrieval call binding the contract method 0xef8b30f7.
//
// Solidity: function previewDeposit(uint256 assets) view returns(uint256)
func (_PolyVault *PolyVaultCallerSession) PreviewDeposit(assets *big.Int) (*big.Int, error) {
	return _PolyVault.Contract.PreviewDeposit(&_PolyVault.CallOpts, assets)
}

// PreviewMint is a free data retrieval call binding the contract method 0xb3d7f6b9.
//
// Solidity: function previewMint(uint256 shares) view returns(uint256)
func (_PolyVault *PolyVaultCaller) PreviewMint(opts *bind.CallOpts, shares *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _PolyVault.contract.Call(opts, &out, "previewMint", shares)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// PreviewMint is a free data retrieval call binding the contract method 0xb3d7f6b9.
//
// Solidity: function previewMint(uint256 shares) view returns(uint256)
func (_PolyVault *PolyVaultSession) PreviewMint(shares *big.Int) (*big.Int, error) {
	return _PolyVault.Contract.PreviewMint(&_PolyVault.CallOpts, shares)
}

// PreviewMint is a free data retrieval call binding the contract method 0xb3d7f6b9.
//
// Solidity: function previewMint(uint256 shares) view returns(uint256)
func (_PolyVault *PolyVaultCallerSession) PreviewMint(shares *big.Int) (*big.Int, error) {
	return _PolyVault.Contract.PreviewMint(&_PolyVault.CallOpts, shares)
}

// PreviewRedeem is a free data retrieval call binding the contract method 0x4cdad506.
//
// Solidity: function previewRedeem(uint256 shares) view returns(uint256)
func (_PolyVault *PolyVaultCaller) PreviewRedeem(opts *bind.CallOpts, shares *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _PolyVault.contract.Call(opts, &out, "previewRedeem", shares)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// PreviewRedeem is a free data retrieval call binding the contract method 0x4cdad506.
//
// Solidity: function previewRedeem(uint256 shares) view returns(uint256)
func (_PolyVault *PolyVaultSession) PreviewRedeem(shares *big.Int) (*big.Int, error) {
	return _PolyVault.Contract.PreviewRedeem(&_PolyVault.CallOpts, shares)
}

// PreviewRedeem is a free data retrieval call binding the contract method 0x4cdad506.
//
// Solidity: function previewRedeem(uint256 shares) view returns(uint256)
func (_PolyVault *PolyVaultCallerSession) PreviewRedeem(shares *big.Int) (*big.Int, error) {
	return _PolyVault.Contract.PreviewRedeem(&_PolyVault.CallOpts, shares)
}

// PreviewWithdraw is a free data retrieval call binding the contract method 0x0a28a477.
//
// Solidity: function previewWithdraw(uint256 assets) view returns(uint256)
func (_PolyVault *PolyVaultCaller) PreviewWithdraw(opts *bind.CallOpts, assets *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _PolyVault.contract.Call(opts, &out, "previewWithdraw", assets)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// PreviewWithdraw is a free data retrieval call binding the contract method 0x0a28a477.
//
// Solidity: function previewWithdraw(uint256 assets) view returns(uint256)
func (_PolyVault *PolyVaultSession) PreviewWithdraw(assets *big.Int) (*big.Int, error) {
	return _PolyVault.Contract.PreviewWithdraw(&_PolyVault.CallOpts, assets)
}

// PreviewWithdraw is a free data retrieval call binding the contract method 0x0a28a477.
//
// Solidity: function previewWithdraw(uint256 assets) view returns(uint256)
func (_PolyVault *PolyVaultCallerSession) PreviewWithdraw(assets *big.Int) (*big.Int, error) {
	return _PolyVault.Contract.PreviewWithdraw(&_PolyVault.CallOpts, assets)
}

// ProxiableUUID is a free data retrieval call binding the contract method 0x52d1902d.
//
// Solidity: function proxiableUUID() view returns(bytes32)
func (_PolyVault *PolyVaultCaller) ProxiableUUID(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _PolyVault.contract.Call(opts, &out, "proxiableUUID")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// ProxiableUUID is a free data retrieval call binding the contract method 0x52d1902d.
//
// Solidity: function proxiableUUID() view returns(bytes32)
func (_PolyVault *PolyVaultSession) ProxiableUUID() ([32]byte, error) {
	return _PolyVault.Contract.ProxiableUUID(&_PolyVault.CallOpts)
}

// ProxiableUUID is a free data retrieval call binding the contract method 0x52d1902d.
//
// Solidity: function proxiableUUID() view returns(bytes32)
func (_PolyVault *PolyVaultCallerSession) ProxiableUUID() ([32]byte, error) {
	return _PolyVault.Contract.ProxiableUUID(&_PolyVault.CallOpts)
}

// Redeem is a free data retrieval call binding the contract method 0xba087652.
//
// Solidity: function redeem(uint256 , address , address ) pure returns(uint256)
func (_PolyVault *PolyVaultCaller) Redeem(opts *bind.CallOpts, arg0 *big.Int, arg1 common.Address, arg2 common.Address) (*big.Int, error) {
	var out []interface{}
	err := _PolyVault.contract.Call(opts, &out, "redeem", arg0, arg1, arg2)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Redeem is a free data retrieval call binding the contract method 0xba087652.
//
// Solidity: function redeem(uint256 , address , address ) pure returns(uint256)
func (_PolyVault *PolyVaultSession) Redeem(arg0 *big.Int, arg1 common.Address, arg2 common.Address) (*big.Int, error) {
	return _PolyVault.Contract.Redeem(&_PolyVault.CallOpts, arg0, arg1, arg2)
}

// Redeem is a free data retrieval call binding the contract method 0xba087652.
//
// Solidity: function redeem(uint256 , address , address ) pure returns(uint256)
func (_PolyVault *PolyVaultCallerSession) Redeem(arg0 *big.Int, arg1 common.Address, arg2 common.Address) (*big.Int, error) {
	return _PolyVault.Contract.Redeem(&_PolyVault.CallOpts, arg0, arg1, arg2)
}

// StrategyDebt is a free data retrieval call binding the contract method 0x3e9dc762.
//
// Solidity: function strategyDebt() view returns(uint256)
func (_PolyVault *PolyVaultCaller) StrategyDebt(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _PolyVault.contract.Call(opts, &out, "strategyDebt")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// StrategyDebt is a free data retrieval call binding the contract method 0x3e9dc762.
//
// Solidity: function strategyDebt() view returns(uint256)
func (_PolyVault *PolyVaultSession) StrategyDebt() (*big.Int, error) {
	return _PolyVault.Contract.StrategyDebt(&_PolyVault.CallOpts)
}

// StrategyDebt is a free data retrieval call binding the contract method 0x3e9dc762.
//
// Solidity: function strategyDebt() view returns(uint256)
func (_PolyVault *PolyVaultCallerSession) StrategyDebt() (*big.Int, error) {
	return _PolyVault.Contract.StrategyDebt(&_PolyVault.CallOpts)
}

// SupportsInterface is a free data retrieval call binding the contract method 0x01ffc9a7.
//
// Solidity: function supportsInterface(bytes4 interfaceId) view returns(bool)
func (_PolyVault *PolyVaultCaller) SupportsInterface(opts *bind.CallOpts, interfaceId [4]byte) (bool, error) {
	var out []interface{}
	err := _PolyVault.contract.Call(opts, &out, "supportsInterface", interfaceId)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// SupportsInterface is a free data retrieval call binding the contract method 0x01ffc9a7.
//
// Solidity: function supportsInterface(bytes4 interfaceId) view returns(bool)
func (_PolyVault *PolyVaultSession) SupportsInterface(interfaceId [4]byte) (bool, error) {
	return _PolyVault.Contract.SupportsInterface(&_PolyVault.CallOpts, interfaceId)
}

// SupportsInterface is a free data retrieval call binding the contract method 0x01ffc9a7.
//
// Solidity: function supportsInterface(bytes4 interfaceId) view returns(bool)
func (_PolyVault *PolyVaultCallerSession) SupportsInterface(interfaceId [4]byte) (bool, error) {
	return _PolyVault.Contract.SupportsInterface(&_PolyVault.CallOpts, interfaceId)
}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() view returns(string)
func (_PolyVault *PolyVaultCaller) Symbol(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _PolyVault.contract.Call(opts, &out, "symbol")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() view returns(string)
func (_PolyVault *PolyVaultSession) Symbol() (string, error) {
	return _PolyVault.Contract.Symbol(&_PolyVault.CallOpts)
}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() view returns(string)
func (_PolyVault *PolyVaultCallerSession) Symbol() (string, error) {
	return _PolyVault.Contract.Symbol(&_PolyVault.CallOpts)
}

// TotalAssets is a free data retrieval call binding the contract method 0x01e1d114.
//
// Solidity: function totalAssets() view returns(uint256)
func (_PolyVault *PolyVaultCaller) TotalAssets(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _PolyVault.contract.Call(opts, &out, "totalAssets")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// TotalAssets is a free data retrieval call binding the contract method 0x01e1d114.
//
// Solidity: function totalAssets() view returns(uint256)
func (_PolyVault *PolyVaultSession) TotalAssets() (*big.Int, error) {
	return _PolyVault.Contract.TotalAssets(&_PolyVault.CallOpts)
}

// TotalAssets is a free data retrieval call binding the contract method 0x01e1d114.
//
// Solidity: function totalAssets() view returns(uint256)
func (_PolyVault *PolyVaultCallerSession) TotalAssets() (*big.Int, error) {
	return _PolyVault.Contract.TotalAssets(&_PolyVault.CallOpts)
}

// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
//
// Solidity: function totalSupply() view returns(uint256)
func (_PolyVault *PolyVaultCaller) TotalSupply(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _PolyVault.contract.Call(opts, &out, "totalSupply")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
//
// Solidity: function totalSupply() view returns(uint256)
func (_PolyVault *PolyVaultSession) TotalSupply() (*big.Int, error) {
	return _PolyVault.Contract.TotalSupply(&_PolyVault.CallOpts)
}

// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
//
// Solidity: function totalSupply() view returns(uint256)
func (_PolyVault *PolyVaultCallerSession) TotalSupply() (*big.Int, error) {
	return _PolyVault.Contract.TotalSupply(&_PolyVault.CallOpts)
}

// Withdraw is a free data retrieval call binding the contract method 0xb460af94.
//
// Solidity: function withdraw(uint256 , address , address ) pure returns(uint256)
func (_PolyVault *PolyVaultCaller) Withdraw(opts *bind.CallOpts, arg0 *big.Int, arg1 common.Address, arg2 common.Address) (*big.Int, error) {
	var out []interface{}
	err := _PolyVault.contract.Call(opts, &out, "withdraw", arg0, arg1, arg2)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Withdraw is a free data retrieval call binding the contract method 0xb460af94.
//
// Solidity: function withdraw(uint256 , address , address ) pure returns(uint256)
func (_PolyVault *PolyVaultSession) Withdraw(arg0 *big.Int, arg1 common.Address, arg2 common.Address) (*big.Int, error) {
	return _PolyVault.Contract.Withdraw(&_PolyVault.CallOpts, arg0, arg1, arg2)
}

// Withdraw is a free data retrieval call binding the contract method 0xb460af94.
//
// Solidity: function withdraw(uint256 , address , address ) pure returns(uint256)
func (_PolyVault *PolyVaultCallerSession) Withdraw(arg0 *big.Int, arg1 common.Address, arg2 common.Address) (*big.Int, error) {
	return _PolyVault.Contract.Withdraw(&_PolyVault.CallOpts, arg0, arg1, arg2)
}

// WithdrawalDelay is a free data retrieval call binding the contract method 0xa7ab6961.
//
// Solidity: function withdrawalDelay() view returns(uint256)
func (_PolyVault *PolyVaultCaller) WithdrawalDelay(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _PolyVault.contract.Call(opts, &out, "withdrawalDelay")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// WithdrawalDelay is a free data retrieval call binding the contract method 0xa7ab6961.
//
// Solidity: function withdrawalDelay() view returns(uint256)
func (_PolyVault *PolyVaultSession) WithdrawalDelay() (*big.Int, error) {
	return _PolyVault.Contract.WithdrawalDelay(&_PolyVault.CallOpts)
}

// WithdrawalDelay is a free data retrieval call binding the contract method 0xa7ab6961.
//
// Solidity: function withdrawalDelay() view returns(uint256)
func (_PolyVault *PolyVaultCallerSession) WithdrawalDelay() (*big.Int, error) {
	return _PolyVault.Contract.WithdrawalDelay(&_PolyVault.CallOpts)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(address spender, uint256 value) returns(bool)
func (_PolyVault *PolyVaultTransactor) Approve(opts *bind.TransactOpts, spender common.Address, value *big.Int) (*types.Transaction, error) {
	return _PolyVault.contract.Transact(opts, "approve", spender, value)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(address spender, uint256 value) returns(bool)
func (_PolyVault *PolyVaultSession) Approve(spender common.Address, value *big.Int) (*types.Transaction, error) {
	return _PolyVault.Contract.Approve(&_PolyVault.TransactOpts, spender, value)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(address spender, uint256 value) returns(bool)
func (_PolyVault *PolyVaultTransactorSession) Approve(spender common.Address, value *big.Int) (*types.Transaction, error) {
	return _PolyVault.Contract.Approve(&_PolyVault.TransactOpts, spender, value)
}

// CancelWithdraw is a paid mutator transaction binding the contract method 0x84b76824.
//
// Solidity: function cancelWithdraw() returns()
func (_PolyVault *PolyVaultTransactor) CancelWithdraw(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _PolyVault.contract.Transact(opts, "cancelWithdraw")
}

// CancelWithdraw is a paid mutator transaction binding the contract method 0x84b76824.
//
// Solidity: function cancelWithdraw() returns()
func (_PolyVault *PolyVaultSession) CancelWithdraw() (*types.Transaction, error) {
	return _PolyVault.Contract.CancelWithdraw(&_PolyVault.TransactOpts)
}

// CancelWithdraw is a paid mutator transaction binding the contract method 0x84b76824.
//
// Solidity: function cancelWithdraw() returns()
func (_PolyVault *PolyVaultTransactorSession) CancelWithdraw() (*types.Transaction, error) {
	return _PolyVault.Contract.CancelWithdraw(&_PolyVault.TransactOpts)
}

// Deposit is a paid mutator transaction binding the contract method 0x6e553f65.
//
// Solidity: function deposit(uint256 assets, address receiver) returns(uint256)
func (_PolyVault *PolyVaultTransactor) Deposit(opts *bind.TransactOpts, assets *big.Int, receiver common.Address) (*types.Transaction, error) {
	return _PolyVault.contract.Transact(opts, "deposit", assets, receiver)
}

// Deposit is a paid mutator transaction binding the contract method 0x6e553f65.
//
// Solidity: function deposit(uint256 assets, address receiver) returns(uint256)
func (_PolyVault *PolyVaultSession) Deposit(assets *big.Int, receiver common.Address) (*types.Transaction, error) {
	return _PolyVault.Contract.Deposit(&_PolyVault.TransactOpts, assets, receiver)
}

// Deposit is a paid mutator transaction binding the contract method 0x6e553f65.
//
// Solidity: function deposit(uint256 assets, address receiver) returns(uint256)
func (_PolyVault *PolyVaultTransactorSession) Deposit(assets *big.Int, receiver common.Address) (*types.Transaction, error) {
	return _PolyVault.Contract.Deposit(&_PolyVault.TransactOpts, assets, receiver)
}

// DepositFromStrategy is a paid mutator transaction binding the contract method 0x83c56b02.
//
// Solidity: function depositFromStrategy(uint256 amount) returns()
func (_PolyVault *PolyVaultTransactor) DepositFromStrategy(opts *bind.TransactOpts, amount *big.Int) (*types.Transaction, error) {
	return _PolyVault.contract.Transact(opts, "depositFromStrategy", amount)
}

// DepositFromStrategy is a paid mutator transaction binding the contract method 0x83c56b02.
//
// Solidity: function depositFromStrategy(uint256 amount) returns()
func (_PolyVault *PolyVaultSession) DepositFromStrategy(amount *big.Int) (*types.Transaction, error) {
	return _PolyVault.Contract.DepositFromStrategy(&_PolyVault.TransactOpts, amount)
}

// DepositFromStrategy is a paid mutator transaction binding the contract method 0x83c56b02.
//
// Solidity: function depositFromStrategy(uint256 amount) returns()
func (_PolyVault *PolyVaultTransactorSession) DepositFromStrategy(amount *big.Int) (*types.Transaction, error) {
	return _PolyVault.Contract.DepositFromStrategy(&_PolyVault.TransactOpts, amount)
}

// ExecuteWithdraw is a paid mutator transaction binding the contract method 0xf8fd9795.
//
// Solidity: function executeWithdraw() returns()
func (_PolyVault *PolyVaultTransactor) ExecuteWithdraw(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _PolyVault.contract.Transact(opts, "executeWithdraw")
}

// ExecuteWithdraw is a paid mutator transaction binding the contract method 0xf8fd9795.
//
// Solidity: function executeWithdraw() returns()
func (_PolyVault *PolyVaultSession) ExecuteWithdraw() (*types.Transaction, error) {
	return _PolyVault.Contract.ExecuteWithdraw(&_PolyVault.TransactOpts)
}

// ExecuteWithdraw is a paid mutator transaction binding the contract method 0xf8fd9795.
//
// Solidity: function executeWithdraw() returns()
func (_PolyVault *PolyVaultTransactorSession) ExecuteWithdraw() (*types.Transaction, error) {
	return _PolyVault.Contract.ExecuteWithdraw(&_PolyVault.TransactOpts)
}

// GrantRole is a paid mutator transaction binding the contract method 0x2f2ff15d.
//
// Solidity: function grantRole(bytes32 role, address account) returns()
func (_PolyVault *PolyVaultTransactor) GrantRole(opts *bind.TransactOpts, role [32]byte, account common.Address) (*types.Transaction, error) {
	return _PolyVault.contract.Transact(opts, "grantRole", role, account)
}

// GrantRole is a paid mutator transaction binding the contract method 0x2f2ff15d.
//
// Solidity: function grantRole(bytes32 role, address account) returns()
func (_PolyVault *PolyVaultSession) GrantRole(role [32]byte, account common.Address) (*types.Transaction, error) {
	return _PolyVault.Contract.GrantRole(&_PolyVault.TransactOpts, role, account)
}

// GrantRole is a paid mutator transaction binding the contract method 0x2f2ff15d.
//
// Solidity: function grantRole(bytes32 role, address account) returns()
func (_PolyVault *PolyVaultTransactorSession) GrantRole(role [32]byte, account common.Address) (*types.Transaction, error) {
	return _PolyVault.Contract.GrantRole(&_PolyVault.TransactOpts, role, account)
}

// Initialize is a paid mutator transaction binding the contract method 0x7403c6cd.
//
// Solidity: function initialize(address _usdc, address _admin, address _strategist, address _guardian, address _feeRecipient, uint256 _withdrawalDelay, uint256 _maxAllocation, uint256 _performanceFee) returns()
func (_PolyVault *PolyVaultTransactor) Initialize(opts *bind.TransactOpts, _usdc common.Address, _admin common.Address, _strategist common.Address, _guardian common.Address, _feeRecipient common.Address, _withdrawalDelay *big.Int, _maxAllocation *big.Int, _performanceFee *big.Int) (*types.Transaction, error) {
	return _PolyVault.contract.Transact(opts, "initialize", _usdc, _admin, _strategist, _guardian, _feeRecipient, _withdrawalDelay, _maxAllocation, _performanceFee)
}

// Initialize is a paid mutator transaction binding the contract method 0x7403c6cd.
//
// Solidity: function initialize(address _usdc, address _admin, address _strategist, address _guardian, address _feeRecipient, uint256 _withdrawalDelay, uint256 _maxAllocation, uint256 _performanceFee) returns()
func (_PolyVault *PolyVaultSession) Initialize(_usdc common.Address, _admin common.Address, _strategist common.Address, _guardian common.Address, _feeRecipient common.Address, _withdrawalDelay *big.Int, _maxAllocation *big.Int, _performanceFee *big.Int) (*types.Transaction, error) {
	return _PolyVault.Contract.Initialize(&_PolyVault.TransactOpts, _usdc, _admin, _strategist, _guardian, _feeRecipient, _withdrawalDelay, _maxAllocation, _performanceFee)
}

// Initialize is a paid mutator transaction binding the contract method 0x7403c6cd.
//
// Solidity: function initialize(address _usdc, address _admin, address _strategist, address _guardian, address _feeRecipient, uint256 _withdrawalDelay, uint256 _maxAllocation, uint256 _performanceFee) returns()
func (_PolyVault *PolyVaultTransactorSession) Initialize(_usdc common.Address, _admin common.Address, _strategist common.Address, _guardian common.Address, _feeRecipient common.Address, _withdrawalDelay *big.Int, _maxAllocation *big.Int, _performanceFee *big.Int) (*types.Transaction, error) {
	return _PolyVault.Contract.Initialize(&_PolyVault.TransactOpts, _usdc, _admin, _strategist, _guardian, _feeRecipient, _withdrawalDelay, _maxAllocation, _performanceFee)
}

// Mint is a paid mutator transaction binding the contract method 0x94bf804d.
//
// Solidity: function mint(uint256 shares, address receiver) returns(uint256)
func (_PolyVault *PolyVaultTransactor) Mint(opts *bind.TransactOpts, shares *big.Int, receiver common.Address) (*types.Transaction, error) {
	return _PolyVault.contract.Transact(opts, "mint", shares, receiver)
}

// Mint is a paid mutator transaction binding the contract method 0x94bf804d.
//
// Solidity: function mint(uint256 shares, address receiver) returns(uint256)
func (_PolyVault *PolyVaultSession) Mint(shares *big.Int, receiver common.Address) (*types.Transaction, error) {
	return _PolyVault.Contract.Mint(&_PolyVault.TransactOpts, shares, receiver)
}

// Mint is a paid mutator transaction binding the contract method 0x94bf804d.
//
// Solidity: function mint(uint256 shares, address receiver) returns(uint256)
func (_PolyVault *PolyVaultTransactorSession) Mint(shares *big.Int, receiver common.Address) (*types.Transaction, error) {
	return _PolyVault.Contract.Mint(&_PolyVault.TransactOpts, shares, receiver)
}

// Pause is a paid mutator transaction binding the contract method 0x8456cb59.
//
// Solidity: function pause() returns()
func (_PolyVault *PolyVaultTransactor) Pause(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _PolyVault.contract.Transact(opts, "pause")
}

// Pause is a paid mutator transaction binding the contract method 0x8456cb59.
//
// Solidity: function pause() returns()
func (_PolyVault *PolyVaultSession) Pause() (*types.Transaction, error) {
	return _PolyVault.Contract.Pause(&_PolyVault.TransactOpts)
}

// Pause is a paid mutator transaction binding the contract method 0x8456cb59.
//
// Solidity: function pause() returns()
func (_PolyVault *PolyVaultTransactorSession) Pause() (*types.Transaction, error) {
	return _PolyVault.Contract.Pause(&_PolyVault.TransactOpts)
}

// RenounceRole is a paid mutator transaction binding the contract method 0x36568abe.
//
// Solidity: function renounceRole(bytes32 role, address callerConfirmation) returns()
func (_PolyVault *PolyVaultTransactor) RenounceRole(opts *bind.TransactOpts, role [32]byte, callerConfirmation common.Address) (*types.Transaction, error) {
	return _PolyVault.contract.Transact(opts, "renounceRole", role, callerConfirmation)
}

// RenounceRole is a paid mutator transaction binding the contract method 0x36568abe.
//
// Solidity: function renounceRole(bytes32 role, address callerConfirmation) returns()
func (_PolyVault *PolyVaultSession) RenounceRole(role [32]byte, callerConfirmation common.Address) (*types.Transaction, error) {
	return _PolyVault.Contract.RenounceRole(&_PolyVault.TransactOpts, role, callerConfirmation)
}

// RenounceRole is a paid mutator transaction binding the contract method 0x36568abe.
//
// Solidity: function renounceRole(bytes32 role, address callerConfirmation) returns()
func (_PolyVault *PolyVaultTransactorSession) RenounceRole(role [32]byte, callerConfirmation common.Address) (*types.Transaction, error) {
	return _PolyVault.Contract.RenounceRole(&_PolyVault.TransactOpts, role, callerConfirmation)
}

// RequestWithdraw is a paid mutator transaction binding the contract method 0x745400c9.
//
// Solidity: function requestWithdraw(uint256 shares) returns()
func (_PolyVault *PolyVaultTransactor) RequestWithdraw(opts *bind.TransactOpts, shares *big.Int) (*types.Transaction, error) {
	return _PolyVault.contract.Transact(opts, "requestWithdraw", shares)
}

// RequestWithdraw is a paid mutator transaction binding the contract method 0x745400c9.
//
// Solidity: function requestWithdraw(uint256 shares) returns()
func (_PolyVault *PolyVaultSession) RequestWithdraw(shares *big.Int) (*types.Transaction, error) {
	return _PolyVault.Contract.RequestWithdraw(&_PolyVault.TransactOpts, shares)
}

// RequestWithdraw is a paid mutator transaction binding the contract method 0x745400c9.
//
// Solidity: function requestWithdraw(uint256 shares) returns()
func (_PolyVault *PolyVaultTransactorSession) RequestWithdraw(shares *big.Int) (*types.Transaction, error) {
	return _PolyVault.Contract.RequestWithdraw(&_PolyVault.TransactOpts, shares)
}

// RevokeRole is a paid mutator transaction binding the contract method 0xd547741f.
//
// Solidity: function revokeRole(bytes32 role, address account) returns()
func (_PolyVault *PolyVaultTransactor) RevokeRole(opts *bind.TransactOpts, role [32]byte, account common.Address) (*types.Transaction, error) {
	return _PolyVault.contract.Transact(opts, "revokeRole", role, account)
}

// RevokeRole is a paid mutator transaction binding the contract method 0xd547741f.
//
// Solidity: function revokeRole(bytes32 role, address account) returns()
func (_PolyVault *PolyVaultSession) RevokeRole(role [32]byte, account common.Address) (*types.Transaction, error) {
	return _PolyVault.Contract.RevokeRole(&_PolyVault.TransactOpts, role, account)
}

// RevokeRole is a paid mutator transaction binding the contract method 0xd547741f.
//
// Solidity: function revokeRole(bytes32 role, address account) returns()
func (_PolyVault *PolyVaultTransactorSession) RevokeRole(role [32]byte, account common.Address) (*types.Transaction, error) {
	return _PolyVault.Contract.RevokeRole(&_PolyVault.TransactOpts, role, account)
}

// SetDepositLimits is a paid mutator transaction binding the contract method 0x4eddea06.
//
// Solidity: function setDepositLimits(uint256 _min, uint256 _max) returns()
func (_PolyVault *PolyVaultTransactor) SetDepositLimits(opts *bind.TransactOpts, _min *big.Int, _max *big.Int) (*types.Transaction, error) {
	return _PolyVault.contract.Transact(opts, "setDepositLimits", _min, _max)
}

// SetDepositLimits is a paid mutator transaction binding the contract method 0x4eddea06.
//
// Solidity: function setDepositLimits(uint256 _min, uint256 _max) returns()
func (_PolyVault *PolyVaultSession) SetDepositLimits(_min *big.Int, _max *big.Int) (*types.Transaction, error) {
	return _PolyVault.Contract.SetDepositLimits(&_PolyVault.TransactOpts, _min, _max)
}

// SetDepositLimits is a paid mutator transaction binding the contract method 0x4eddea06.
//
// Solidity: function setDepositLimits(uint256 _min, uint256 _max) returns()
func (_PolyVault *PolyVaultTransactorSession) SetDepositLimits(_min *big.Int, _max *big.Int) (*types.Transaction, error) {
	return _PolyVault.Contract.SetDepositLimits(&_PolyVault.TransactOpts, _min, _max)
}

// SetFeeRecipient is a paid mutator transaction binding the contract method 0xe74b981b.
//
// Solidity: function setFeeRecipient(address _recipient) returns()
func (_PolyVault *PolyVaultTransactor) SetFeeRecipient(opts *bind.TransactOpts, _recipient common.Address) (*types.Transaction, error) {
	return _PolyVault.contract.Transact(opts, "setFeeRecipient", _recipient)
}

// SetFeeRecipient is a paid mutator transaction binding the contract method 0xe74b981b.
//
// Solidity: function setFeeRecipient(address _recipient) returns()
func (_PolyVault *PolyVaultSession) SetFeeRecipient(_recipient common.Address) (*types.Transaction, error) {
	return _PolyVault.Contract.SetFeeRecipient(&_PolyVault.TransactOpts, _recipient)
}

// SetFeeRecipient is a paid mutator transaction binding the contract method 0xe74b981b.
//
// Solidity: function setFeeRecipient(address _recipient) returns()
func (_PolyVault *PolyVaultTransactorSession) SetFeeRecipient(_recipient common.Address) (*types.Transaction, error) {
	return _PolyVault.Contract.SetFeeRecipient(&_PolyVault.TransactOpts, _recipient)
}

// SetMaxStrategyAllocation is a paid mutator transaction binding the contract method 0x86f1ef0c.
//
// Solidity: function setMaxStrategyAllocation(uint256 _allocation) returns()
func (_PolyVault *PolyVaultTransactor) SetMaxStrategyAllocation(opts *bind.TransactOpts, _allocation *big.Int) (*types.Transaction, error) {
	return _PolyVault.contract.Transact(opts, "setMaxStrategyAllocation", _allocation)
}

// SetMaxStrategyAllocation is a paid mutator transaction binding the contract method 0x86f1ef0c.
//
// Solidity: function setMaxStrategyAllocation(uint256 _allocation) returns()
func (_PolyVault *PolyVaultSession) SetMaxStrategyAllocation(_allocation *big.Int) (*types.Transaction, error) {
	return _PolyVault.Contract.SetMaxStrategyAllocation(&_PolyVault.TransactOpts, _allocation)
}

// SetMaxStrategyAllocation is a paid mutator transaction binding the contract method 0x86f1ef0c.
//
// Solidity: function setMaxStrategyAllocation(uint256 _allocation) returns()
func (_PolyVault *PolyVaultTransactorSession) SetMaxStrategyAllocation(_allocation *big.Int) (*types.Transaction, error) {
	return _PolyVault.Contract.SetMaxStrategyAllocation(&_PolyVault.TransactOpts, _allocation)
}

// SetPerformanceFee is a paid mutator transaction binding the contract method 0x70897b23.
//
// Solidity: function setPerformanceFee(uint256 _fee) returns()
func (_PolyVault *PolyVaultTransactor) SetPerformanceFee(opts *bind.TransactOpts, _fee *big.Int) (*types.Transaction, error) {
	return _PolyVault.contract.Transact(opts, "setPerformanceFee", _fee)
}

// SetPerformanceFee is a paid mutator transaction binding the contract method 0x70897b23.
//
// Solidity: function setPerformanceFee(uint256 _fee) returns()
func (_PolyVault *PolyVaultSession) SetPerformanceFee(_fee *big.Int) (*types.Transaction, error) {
	return _PolyVault.Contract.SetPerformanceFee(&_PolyVault.TransactOpts, _fee)
}

// SetPerformanceFee is a paid mutator transaction binding the contract method 0x70897b23.
//
// Solidity: function setPerformanceFee(uint256 _fee) returns()
func (_PolyVault *PolyVaultTransactorSession) SetPerformanceFee(_fee *big.Int) (*types.Transaction, error) {
	return _PolyVault.Contract.SetPerformanceFee(&_PolyVault.TransactOpts, _fee)
}

// SetWithdrawalDelay is a paid mutator transaction binding the contract method 0xd2c13da5.
//
// Solidity: function setWithdrawalDelay(uint256 _delay) returns()
func (_PolyVault *PolyVaultTransactor) SetWithdrawalDelay(opts *bind.TransactOpts, _delay *big.Int) (*types.Transaction, error) {
	return _PolyVault.contract.Transact(opts, "setWithdrawalDelay", _delay)
}

// SetWithdrawalDelay is a paid mutator transaction binding the contract method 0xd2c13da5.
//
// Solidity: function setWithdrawalDelay(uint256 _delay) returns()
func (_PolyVault *PolyVaultSession) SetWithdrawalDelay(_delay *big.Int) (*types.Transaction, error) {
	return _PolyVault.Contract.SetWithdrawalDelay(&_PolyVault.TransactOpts, _delay)
}

// SetWithdrawalDelay is a paid mutator transaction binding the contract method 0xd2c13da5.
//
// Solidity: function setWithdrawalDelay(uint256 _delay) returns()
func (_PolyVault *PolyVaultTransactorSession) SetWithdrawalDelay(_delay *big.Int) (*types.Transaction, error) {
	return _PolyVault.Contract.SetWithdrawalDelay(&_PolyVault.TransactOpts, _delay)
}

// Transfer is a paid mutator transaction binding the contract method 0xa9059cbb.
//
// Solidity: function transfer(address to, uint256 value) returns(bool)
func (_PolyVault *PolyVaultTransactor) Transfer(opts *bind.TransactOpts, to common.Address, value *big.Int) (*types.Transaction, error) {
	return _PolyVault.contract.Transact(opts, "transfer", to, value)
}

// Transfer is a paid mutator transaction binding the contract method 0xa9059cbb.
//
// Solidity: function transfer(address to, uint256 value) returns(bool)
func (_PolyVault *PolyVaultSession) Transfer(to common.Address, value *big.Int) (*types.Transaction, error) {
	return _PolyVault.Contract.Transfer(&_PolyVault.TransactOpts, to, value)
}

// Transfer is a paid mutator transaction binding the contract method 0xa9059cbb.
//
// Solidity: function transfer(address to, uint256 value) returns(bool)
func (_PolyVault *PolyVaultTransactorSession) Transfer(to common.Address, value *big.Int) (*types.Transaction, error) {
	return _PolyVault.Contract.Transfer(&_PolyVault.TransactOpts, to, value)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(address from, address to, uint256 value) returns(bool)
func (_PolyVault *PolyVaultTransactor) TransferFrom(opts *bind.TransactOpts, from common.Address, to common.Address, value *big.Int) (*types.Transaction, error) {
	return _PolyVault.contract.Transact(opts, "transferFrom", from, to, value)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(address from, address to, uint256 value) returns(bool)
func (_PolyVault *PolyVaultSession) TransferFrom(from common.Address, to common.Address, value *big.Int) (*types.Transaction, error) {
	return _PolyVault.Contract.TransferFrom(&_PolyVault.TransactOpts, from, to, value)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(address from, address to, uint256 value) returns(bool)
func (_PolyVault *PolyVaultTransactorSession) TransferFrom(from common.Address, to common.Address, value *big.Int) (*types.Transaction, error) {
	return _PolyVault.Contract.TransferFrom(&_PolyVault.TransactOpts, from, to, value)
}

// Unpause is a paid mutator transaction binding the contract method 0x3f4ba83a.
//
// Solidity: function unpause() returns()
func (_PolyVault *PolyVaultTransactor) Unpause(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _PolyVault.contract.Transact(opts, "unpause")
}

// Unpause is a paid mutator transaction binding the contract method 0x3f4ba83a.
//
// Solidity: function unpause() returns()
func (_PolyVault *PolyVaultSession) Unpause() (*types.Transaction, error) {
	return _PolyVault.Contract.Unpause(&_PolyVault.TransactOpts)
}

// Unpause is a paid mutator transaction binding the contract method 0x3f4ba83a.
//
// Solidity: function unpause() returns()
func (_PolyVault *PolyVaultTransactorSession) Unpause() (*types.Transaction, error) {
	return _PolyVault.Contract.Unpause(&_PolyVault.TransactOpts)
}

// UpgradeToAndCall is a paid mutator transaction binding the contract method 0x4f1ef286.
//
// Solidity: function upgradeToAndCall(address newImplementation, bytes data) payable returns()
func (_PolyVault *PolyVaultTransactor) UpgradeToAndCall(opts *bind.TransactOpts, newImplementation common.Address, data []byte) (*types.Transaction, error) {
	return _PolyVault.contract.Transact(opts, "upgradeToAndCall", newImplementation, data)
}

// UpgradeToAndCall is a paid mutator transaction binding the contract method 0x4f1ef286.
//
// Solidity: function upgradeToAndCall(address newImplementation, bytes data) payable returns()
func (_PolyVault *PolyVaultSession) UpgradeToAndCall(newImplementation common.Address, data []byte) (*types.Transaction, error) {
	return _PolyVault.Contract.UpgradeToAndCall(&_PolyVault.TransactOpts, newImplementation, data)
}

// UpgradeToAndCall is a paid mutator transaction binding the contract method 0x4f1ef286.
//
// Solidity: function upgradeToAndCall(address newImplementation, bytes data) payable returns()
func (_PolyVault *PolyVaultTransactorSession) UpgradeToAndCall(newImplementation common.Address, data []byte) (*types.Transaction, error) {
	return _PolyVault.Contract.UpgradeToAndCall(&_PolyVault.TransactOpts, newImplementation, data)
}

// WithdrawToStrategy is a paid mutator transaction binding the contract method 0x9aca479e.
//
// Solidity: function withdrawToStrategy(uint256 amount) returns()
func (_PolyVault *PolyVaultTransactor) WithdrawToStrategy(opts *bind.TransactOpts, amount *big.Int) (*types.Transaction, error) {
	return _PolyVault.contract.Transact(opts, "withdrawToStrategy", amount)
}

// WithdrawToStrategy is a paid mutator transaction binding the contract method 0x9aca479e.
//
// Solidity: function withdrawToStrategy(uint256 amount) returns()
func (_PolyVault *PolyVaultSession) WithdrawToStrategy(amount *big.Int) (*types.Transaction, error) {
	return _PolyVault.Contract.WithdrawToStrategy(&_PolyVault.TransactOpts, amount)
}

// WithdrawToStrategy is a paid mutator transaction binding the contract method 0x9aca479e.
//
// Solidity: function withdrawToStrategy(uint256 amount) returns()
func (_PolyVault *PolyVaultTransactorSession) WithdrawToStrategy(amount *big.Int) (*types.Transaction, error) {
	return _PolyVault.Contract.WithdrawToStrategy(&_PolyVault.TransactOpts, amount)
}

// PolyVaultApprovalIterator is returned from FilterApproval and is used to iterate over the raw logs and unpacked data for Approval events raised by the PolyVault contract.
type PolyVaultApprovalIterator struct {
	Event *PolyVaultApproval // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *PolyVaultApprovalIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(PolyVaultApproval)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(PolyVaultApproval)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *PolyVaultApprovalIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *PolyVaultApprovalIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// PolyVaultApproval represents a Approval event raised by the PolyVault contract.
type PolyVaultApproval struct {
	Owner   common.Address
	Spender common.Address
	Value   *big.Int
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterApproval is a free log retrieval operation binding the contract event 0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925.
//
// Solidity: event Approval(address indexed owner, address indexed spender, uint256 value)
func (_PolyVault *PolyVaultFilterer) FilterApproval(opts *bind.FilterOpts, owner []common.Address, spender []common.Address) (*PolyVaultApprovalIterator, error) {

	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}
	var spenderRule []interface{}
	for _, spenderItem := range spender {
		spenderRule = append(spenderRule, spenderItem)
	}

	logs, sub, err := _PolyVault.contract.FilterLogs(opts, "Approval", ownerRule, spenderRule)
	if err != nil {
		return nil, err
	}
	return &PolyVaultApprovalIterator{contract: _PolyVault.contract, event: "Approval", logs: logs, sub: sub}, nil
}

// WatchApproval is a free log subscription operation binding the contract event 0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925.
//
// Solidity: event Approval(address indexed owner, address indexed spender, uint256 value)
func (_PolyVault *PolyVaultFilterer) WatchApproval(opts *bind.WatchOpts, sink chan<- *PolyVaultApproval, owner []common.Address, spender []common.Address) (event.Subscription, error) {

	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}
	var spenderRule []interface{}
	for _, spenderItem := range spender {
		spenderRule = append(spenderRule, spenderItem)
	}

	logs, sub, err := _PolyVault.contract.WatchLogs(opts, "Approval", ownerRule, spenderRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(PolyVaultApproval)
				if err := _PolyVault.contract.UnpackLog(event, "Approval", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseApproval is a log parse operation binding the contract event 0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925.
//
// Solidity: event Approval(address indexed owner, address indexed spender, uint256 value)
func (_PolyVault *PolyVaultFilterer) ParseApproval(log types.Log) (*PolyVaultApproval, error) {
	event := new(PolyVaultApproval)
	if err := _PolyVault.contract.UnpackLog(event, "Approval", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// PolyVaultDepositIterator is returned from FilterDeposit and is used to iterate over the raw logs and unpacked data for Deposit events raised by the PolyVault contract.
type PolyVaultDepositIterator struct {
	Event *PolyVaultDeposit // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *PolyVaultDepositIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(PolyVaultDeposit)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(PolyVaultDeposit)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *PolyVaultDepositIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *PolyVaultDepositIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// PolyVaultDeposit represents a Deposit event raised by the PolyVault contract.
type PolyVaultDeposit struct {
	Sender common.Address
	Owner  common.Address
	Assets *big.Int
	Shares *big.Int
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterDeposit is a free log retrieval operation binding the contract event 0xdcbc1c05240f31ff3ad067ef1ee35ce4997762752e3a095284754544f4c709d7.
//
// Solidity: event Deposit(address indexed sender, address indexed owner, uint256 assets, uint256 shares)
func (_PolyVault *PolyVaultFilterer) FilterDeposit(opts *bind.FilterOpts, sender []common.Address, owner []common.Address) (*PolyVaultDepositIterator, error) {

	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}
	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}

	logs, sub, err := _PolyVault.contract.FilterLogs(opts, "Deposit", senderRule, ownerRule)
	if err != nil {
		return nil, err
	}
	return &PolyVaultDepositIterator{contract: _PolyVault.contract, event: "Deposit", logs: logs, sub: sub}, nil
}

// WatchDeposit is a free log subscription operation binding the contract event 0xdcbc1c05240f31ff3ad067ef1ee35ce4997762752e3a095284754544f4c709d7.
//
// Solidity: event Deposit(address indexed sender, address indexed owner, uint256 assets, uint256 shares)
func (_PolyVault *PolyVaultFilterer) WatchDeposit(opts *bind.WatchOpts, sink chan<- *PolyVaultDeposit, sender []common.Address, owner []common.Address) (event.Subscription, error) {

	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}
	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}

	logs, sub, err := _PolyVault.contract.WatchLogs(opts, "Deposit", senderRule, ownerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(PolyVaultDeposit)
				if err := _PolyVault.contract.UnpackLog(event, "Deposit", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseDeposit is a log parse operation binding the contract event 0xdcbc1c05240f31ff3ad067ef1ee35ce4997762752e3a095284754544f4c709d7.
//
// Solidity: event Deposit(address indexed sender, address indexed owner, uint256 assets, uint256 shares)
func (_PolyVault *PolyVaultFilterer) ParseDeposit(log types.Log) (*PolyVaultDeposit, error) {
	event := new(PolyVaultDeposit)
	if err := _PolyVault.contract.UnpackLog(event, "Deposit", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// PolyVaultDepositLimitsUpdatedIterator is returned from FilterDepositLimitsUpdated and is used to iterate over the raw logs and unpacked data for DepositLimitsUpdated events raised by the PolyVault contract.
type PolyVaultDepositLimitsUpdatedIterator struct {
	Event *PolyVaultDepositLimitsUpdated // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *PolyVaultDepositLimitsUpdatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(PolyVaultDepositLimitsUpdated)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(PolyVaultDepositLimitsUpdated)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *PolyVaultDepositLimitsUpdatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *PolyVaultDepositLimitsUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// PolyVaultDepositLimitsUpdated represents a DepositLimitsUpdated event raised by the PolyVault contract.
type PolyVaultDepositLimitsUpdated struct {
	MinDeposit *big.Int
	MaxDeposit *big.Int
	Raw        types.Log // Blockchain specific contextual infos
}

// FilterDepositLimitsUpdated is a free log retrieval operation binding the contract event 0xb2ad710f2954a5376267a683f9ece9ec46ee7dfb47075163379904ee941df8da.
//
// Solidity: event DepositLimitsUpdated(uint256 minDeposit, uint256 maxDeposit)
func (_PolyVault *PolyVaultFilterer) FilterDepositLimitsUpdated(opts *bind.FilterOpts) (*PolyVaultDepositLimitsUpdatedIterator, error) {

	logs, sub, err := _PolyVault.contract.FilterLogs(opts, "DepositLimitsUpdated")
	if err != nil {
		return nil, err
	}
	return &PolyVaultDepositLimitsUpdatedIterator{contract: _PolyVault.contract, event: "DepositLimitsUpdated", logs: logs, sub: sub}, nil
}

// WatchDepositLimitsUpdated is a free log subscription operation binding the contract event 0xb2ad710f2954a5376267a683f9ece9ec46ee7dfb47075163379904ee941df8da.
//
// Solidity: event DepositLimitsUpdated(uint256 minDeposit, uint256 maxDeposit)
func (_PolyVault *PolyVaultFilterer) WatchDepositLimitsUpdated(opts *bind.WatchOpts, sink chan<- *PolyVaultDepositLimitsUpdated) (event.Subscription, error) {

	logs, sub, err := _PolyVault.contract.WatchLogs(opts, "DepositLimitsUpdated")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(PolyVaultDepositLimitsUpdated)
				if err := _PolyVault.contract.UnpackLog(event, "DepositLimitsUpdated", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseDepositLimitsUpdated is a log parse operation binding the contract event 0xb2ad710f2954a5376267a683f9ece9ec46ee7dfb47075163379904ee941df8da.
//
// Solidity: event DepositLimitsUpdated(uint256 minDeposit, uint256 maxDeposit)
func (_PolyVault *PolyVaultFilterer) ParseDepositLimitsUpdated(log types.Log) (*PolyVaultDepositLimitsUpdated, error) {
	event := new(PolyVaultDepositLimitsUpdated)
	if err := _PolyVault.contract.UnpackLog(event, "DepositLimitsUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// PolyVaultFeeRecipientUpdatedIterator is returned from FilterFeeRecipientUpdated and is used to iterate over the raw logs and unpacked data for FeeRecipientUpdated events raised by the PolyVault contract.
type PolyVaultFeeRecipientUpdatedIterator struct {
	Event *PolyVaultFeeRecipientUpdated // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *PolyVaultFeeRecipientUpdatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(PolyVaultFeeRecipientUpdated)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(PolyVaultFeeRecipientUpdated)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *PolyVaultFeeRecipientUpdatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *PolyVaultFeeRecipientUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// PolyVaultFeeRecipientUpdated represents a FeeRecipientUpdated event raised by the PolyVault contract.
type PolyVaultFeeRecipientUpdated struct {
	OldRecipient common.Address
	NewRecipient common.Address
	Raw          types.Log // Blockchain specific contextual infos
}

// FilterFeeRecipientUpdated is a free log retrieval operation binding the contract event 0xaaebcf1bfa00580e41d966056b48521fa9f202645c86d4ddf28113e617c1b1d3.
//
// Solidity: event FeeRecipientUpdated(address oldRecipient, address newRecipient)
func (_PolyVault *PolyVaultFilterer) FilterFeeRecipientUpdated(opts *bind.FilterOpts) (*PolyVaultFeeRecipientUpdatedIterator, error) {

	logs, sub, err := _PolyVault.contract.FilterLogs(opts, "FeeRecipientUpdated")
	if err != nil {
		return nil, err
	}
	return &PolyVaultFeeRecipientUpdatedIterator{contract: _PolyVault.contract, event: "FeeRecipientUpdated", logs: logs, sub: sub}, nil
}

// WatchFeeRecipientUpdated is a free log subscription operation binding the contract event 0xaaebcf1bfa00580e41d966056b48521fa9f202645c86d4ddf28113e617c1b1d3.
//
// Solidity: event FeeRecipientUpdated(address oldRecipient, address newRecipient)
func (_PolyVault *PolyVaultFilterer) WatchFeeRecipientUpdated(opts *bind.WatchOpts, sink chan<- *PolyVaultFeeRecipientUpdated) (event.Subscription, error) {

	logs, sub, err := _PolyVault.contract.WatchLogs(opts, "FeeRecipientUpdated")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(PolyVaultFeeRecipientUpdated)
				if err := _PolyVault.contract.UnpackLog(event, "FeeRecipientUpdated", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseFeeRecipientUpdated is a log parse operation binding the contract event 0xaaebcf1bfa00580e41d966056b48521fa9f202645c86d4ddf28113e617c1b1d3.
//
// Solidity: event FeeRecipientUpdated(address oldRecipient, address newRecipient)
func (_PolyVault *PolyVaultFilterer) ParseFeeRecipientUpdated(log types.Log) (*PolyVaultFeeRecipientUpdated, error) {
	event := new(PolyVaultFeeRecipientUpdated)
	if err := _PolyVault.contract.UnpackLog(event, "FeeRecipientUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// PolyVaultInitializedIterator is returned from FilterInitialized and is used to iterate over the raw logs and unpacked data for Initialized events raised by the PolyVault contract.
type PolyVaultInitializedIterator struct {
	Event *PolyVaultInitialized // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *PolyVaultInitializedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(PolyVaultInitialized)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(PolyVaultInitialized)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *PolyVaultInitializedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *PolyVaultInitializedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// PolyVaultInitialized represents a Initialized event raised by the PolyVault contract.
type PolyVaultInitialized struct {
	Version uint64
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterInitialized is a free log retrieval operation binding the contract event 0xc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d2.
//
// Solidity: event Initialized(uint64 version)
func (_PolyVault *PolyVaultFilterer) FilterInitialized(opts *bind.FilterOpts) (*PolyVaultInitializedIterator, error) {

	logs, sub, err := _PolyVault.contract.FilterLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return &PolyVaultInitializedIterator{contract: _PolyVault.contract, event: "Initialized", logs: logs, sub: sub}, nil
}

// WatchInitialized is a free log subscription operation binding the contract event 0xc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d2.
//
// Solidity: event Initialized(uint64 version)
func (_PolyVault *PolyVaultFilterer) WatchInitialized(opts *bind.WatchOpts, sink chan<- *PolyVaultInitialized) (event.Subscription, error) {

	logs, sub, err := _PolyVault.contract.WatchLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(PolyVaultInitialized)
				if err := _PolyVault.contract.UnpackLog(event, "Initialized", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseInitialized is a log parse operation binding the contract event 0xc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d2.
//
// Solidity: event Initialized(uint64 version)
func (_PolyVault *PolyVaultFilterer) ParseInitialized(log types.Log) (*PolyVaultInitialized, error) {
	event := new(PolyVaultInitialized)
	if err := _PolyVault.contract.UnpackLog(event, "Initialized", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// PolyVaultMaxStrategyAllocationUpdatedIterator is returned from FilterMaxStrategyAllocationUpdated and is used to iterate over the raw logs and unpacked data for MaxStrategyAllocationUpdated events raised by the PolyVault contract.
type PolyVaultMaxStrategyAllocationUpdatedIterator struct {
	Event *PolyVaultMaxStrategyAllocationUpdated // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *PolyVaultMaxStrategyAllocationUpdatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(PolyVaultMaxStrategyAllocationUpdated)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(PolyVaultMaxStrategyAllocationUpdated)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *PolyVaultMaxStrategyAllocationUpdatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *PolyVaultMaxStrategyAllocationUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// PolyVaultMaxStrategyAllocationUpdated represents a MaxStrategyAllocationUpdated event raised by the PolyVault contract.
type PolyVaultMaxStrategyAllocationUpdated struct {
	OldAllocation *big.Int
	NewAllocation *big.Int
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterMaxStrategyAllocationUpdated is a free log retrieval operation binding the contract event 0x5755997f472615e724fa79d9ae7a5a69a67307f0eef7b848b07f517abac1adaf.
//
// Solidity: event MaxStrategyAllocationUpdated(uint256 oldAllocation, uint256 newAllocation)
func (_PolyVault *PolyVaultFilterer) FilterMaxStrategyAllocationUpdated(opts *bind.FilterOpts) (*PolyVaultMaxStrategyAllocationUpdatedIterator, error) {

	logs, sub, err := _PolyVault.contract.FilterLogs(opts, "MaxStrategyAllocationUpdated")
	if err != nil {
		return nil, err
	}
	return &PolyVaultMaxStrategyAllocationUpdatedIterator{contract: _PolyVault.contract, event: "MaxStrategyAllocationUpdated", logs: logs, sub: sub}, nil
}

// WatchMaxStrategyAllocationUpdated is a free log subscription operation binding the contract event 0x5755997f472615e724fa79d9ae7a5a69a67307f0eef7b848b07f517abac1adaf.
//
// Solidity: event MaxStrategyAllocationUpdated(uint256 oldAllocation, uint256 newAllocation)
func (_PolyVault *PolyVaultFilterer) WatchMaxStrategyAllocationUpdated(opts *bind.WatchOpts, sink chan<- *PolyVaultMaxStrategyAllocationUpdated) (event.Subscription, error) {

	logs, sub, err := _PolyVault.contract.WatchLogs(opts, "MaxStrategyAllocationUpdated")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(PolyVaultMaxStrategyAllocationUpdated)
				if err := _PolyVault.contract.UnpackLog(event, "MaxStrategyAllocationUpdated", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseMaxStrategyAllocationUpdated is a log parse operation binding the contract event 0x5755997f472615e724fa79d9ae7a5a69a67307f0eef7b848b07f517abac1adaf.
//
// Solidity: event MaxStrategyAllocationUpdated(uint256 oldAllocation, uint256 newAllocation)
func (_PolyVault *PolyVaultFilterer) ParseMaxStrategyAllocationUpdated(log types.Log) (*PolyVaultMaxStrategyAllocationUpdated, error) {
	event := new(PolyVaultMaxStrategyAllocationUpdated)
	if err := _PolyVault.contract.UnpackLog(event, "MaxStrategyAllocationUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// PolyVaultPausedIterator is returned from FilterPaused and is used to iterate over the raw logs and unpacked data for Paused events raised by the PolyVault contract.
type PolyVaultPausedIterator struct {
	Event *PolyVaultPaused // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *PolyVaultPausedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(PolyVaultPaused)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(PolyVaultPaused)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *PolyVaultPausedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *PolyVaultPausedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// PolyVaultPaused represents a Paused event raised by the PolyVault contract.
type PolyVaultPaused struct {
	Account common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterPaused is a free log retrieval operation binding the contract event 0x62e78cea01bee320cd4e420270b5ea74000d11b0c9f74754ebdbfc544b05a258.
//
// Solidity: event Paused(address account)
func (_PolyVault *PolyVaultFilterer) FilterPaused(opts *bind.FilterOpts) (*PolyVaultPausedIterator, error) {

	logs, sub, err := _PolyVault.contract.FilterLogs(opts, "Paused")
	if err != nil {
		return nil, err
	}
	return &PolyVaultPausedIterator{contract: _PolyVault.contract, event: "Paused", logs: logs, sub: sub}, nil
}

// WatchPaused is a free log subscription operation binding the contract event 0x62e78cea01bee320cd4e420270b5ea74000d11b0c9f74754ebdbfc544b05a258.
//
// Solidity: event Paused(address account)
func (_PolyVault *PolyVaultFilterer) WatchPaused(opts *bind.WatchOpts, sink chan<- *PolyVaultPaused) (event.Subscription, error) {

	logs, sub, err := _PolyVault.contract.WatchLogs(opts, "Paused")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(PolyVaultPaused)
				if err := _PolyVault.contract.UnpackLog(event, "Paused", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParsePaused is a log parse operation binding the contract event 0x62e78cea01bee320cd4e420270b5ea74000d11b0c9f74754ebdbfc544b05a258.
//
// Solidity: event Paused(address account)
func (_PolyVault *PolyVaultFilterer) ParsePaused(log types.Log) (*PolyVaultPaused, error) {
	event := new(PolyVaultPaused)
	if err := _PolyVault.contract.UnpackLog(event, "Paused", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// PolyVaultPerformanceFeeUpdatedIterator is returned from FilterPerformanceFeeUpdated and is used to iterate over the raw logs and unpacked data for PerformanceFeeUpdated events raised by the PolyVault contract.
type PolyVaultPerformanceFeeUpdatedIterator struct {
	Event *PolyVaultPerformanceFeeUpdated // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *PolyVaultPerformanceFeeUpdatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(PolyVaultPerformanceFeeUpdated)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(PolyVaultPerformanceFeeUpdated)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *PolyVaultPerformanceFeeUpdatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *PolyVaultPerformanceFeeUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// PolyVaultPerformanceFeeUpdated represents a PerformanceFeeUpdated event raised by the PolyVault contract.
type PolyVaultPerformanceFeeUpdated struct {
	OldFee *big.Int
	NewFee *big.Int
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterPerformanceFeeUpdated is a free log retrieval operation binding the contract event 0x607b1c943753982194530bf7133a5972ea2626e028005410efa54ab20035caf8.
//
// Solidity: event PerformanceFeeUpdated(uint256 oldFee, uint256 newFee)
func (_PolyVault *PolyVaultFilterer) FilterPerformanceFeeUpdated(opts *bind.FilterOpts) (*PolyVaultPerformanceFeeUpdatedIterator, error) {

	logs, sub, err := _PolyVault.contract.FilterLogs(opts, "PerformanceFeeUpdated")
	if err != nil {
		return nil, err
	}
	return &PolyVaultPerformanceFeeUpdatedIterator{contract: _PolyVault.contract, event: "PerformanceFeeUpdated", logs: logs, sub: sub}, nil
}

// WatchPerformanceFeeUpdated is a free log subscription operation binding the contract event 0x607b1c943753982194530bf7133a5972ea2626e028005410efa54ab20035caf8.
//
// Solidity: event PerformanceFeeUpdated(uint256 oldFee, uint256 newFee)
func (_PolyVault *PolyVaultFilterer) WatchPerformanceFeeUpdated(opts *bind.WatchOpts, sink chan<- *PolyVaultPerformanceFeeUpdated) (event.Subscription, error) {

	logs, sub, err := _PolyVault.contract.WatchLogs(opts, "PerformanceFeeUpdated")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(PolyVaultPerformanceFeeUpdated)
				if err := _PolyVault.contract.UnpackLog(event, "PerformanceFeeUpdated", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParsePerformanceFeeUpdated is a log parse operation binding the contract event 0x607b1c943753982194530bf7133a5972ea2626e028005410efa54ab20035caf8.
//
// Solidity: event PerformanceFeeUpdated(uint256 oldFee, uint256 newFee)
func (_PolyVault *PolyVaultFilterer) ParsePerformanceFeeUpdated(log types.Log) (*PolyVaultPerformanceFeeUpdated, error) {
	event := new(PolyVaultPerformanceFeeUpdated)
	if err := _PolyVault.contract.UnpackLog(event, "PerformanceFeeUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// PolyVaultProfitReportedIterator is returned from FilterProfitReported and is used to iterate over the raw logs and unpacked data for ProfitReported events raised by the PolyVault contract.
type PolyVaultProfitReportedIterator struct {
	Event *PolyVaultProfitReported // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *PolyVaultProfitReportedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(PolyVaultProfitReported)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(PolyVaultProfitReported)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *PolyVaultProfitReportedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *PolyVaultProfitReportedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// PolyVaultProfitReported represents a ProfitReported event raised by the PolyVault contract.
type PolyVaultProfitReported struct {
	Profit *big.Int
	Fee    *big.Int
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterProfitReported is a free log retrieval operation binding the contract event 0xfcfbfd1d7fecbea7809bda42bd54ffa877192d8f5170375720ba7197c80181bc.
//
// Solidity: event ProfitReported(uint256 profit, uint256 fee)
func (_PolyVault *PolyVaultFilterer) FilterProfitReported(opts *bind.FilterOpts) (*PolyVaultProfitReportedIterator, error) {

	logs, sub, err := _PolyVault.contract.FilterLogs(opts, "ProfitReported")
	if err != nil {
		return nil, err
	}
	return &PolyVaultProfitReportedIterator{contract: _PolyVault.contract, event: "ProfitReported", logs: logs, sub: sub}, nil
}

// WatchProfitReported is a free log subscription operation binding the contract event 0xfcfbfd1d7fecbea7809bda42bd54ffa877192d8f5170375720ba7197c80181bc.
//
// Solidity: event ProfitReported(uint256 profit, uint256 fee)
func (_PolyVault *PolyVaultFilterer) WatchProfitReported(opts *bind.WatchOpts, sink chan<- *PolyVaultProfitReported) (event.Subscription, error) {

	logs, sub, err := _PolyVault.contract.WatchLogs(opts, "ProfitReported")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(PolyVaultProfitReported)
				if err := _PolyVault.contract.UnpackLog(event, "ProfitReported", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseProfitReported is a log parse operation binding the contract event 0xfcfbfd1d7fecbea7809bda42bd54ffa877192d8f5170375720ba7197c80181bc.
//
// Solidity: event ProfitReported(uint256 profit, uint256 fee)
func (_PolyVault *PolyVaultFilterer) ParseProfitReported(log types.Log) (*PolyVaultProfitReported, error) {
	event := new(PolyVaultProfitReported)
	if err := _PolyVault.contract.UnpackLog(event, "ProfitReported", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// PolyVaultRoleAdminChangedIterator is returned from FilterRoleAdminChanged and is used to iterate over the raw logs and unpacked data for RoleAdminChanged events raised by the PolyVault contract.
type PolyVaultRoleAdminChangedIterator struct {
	Event *PolyVaultRoleAdminChanged // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *PolyVaultRoleAdminChangedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(PolyVaultRoleAdminChanged)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(PolyVaultRoleAdminChanged)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *PolyVaultRoleAdminChangedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *PolyVaultRoleAdminChangedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// PolyVaultRoleAdminChanged represents a RoleAdminChanged event raised by the PolyVault contract.
type PolyVaultRoleAdminChanged struct {
	Role              [32]byte
	PreviousAdminRole [32]byte
	NewAdminRole      [32]byte
	Raw               types.Log // Blockchain specific contextual infos
}

// FilterRoleAdminChanged is a free log retrieval operation binding the contract event 0xbd79b86ffe0ab8e8776151514217cd7cacd52c909f66475c3af44e129f0b00ff.
//
// Solidity: event RoleAdminChanged(bytes32 indexed role, bytes32 indexed previousAdminRole, bytes32 indexed newAdminRole)
func (_PolyVault *PolyVaultFilterer) FilterRoleAdminChanged(opts *bind.FilterOpts, role [][32]byte, previousAdminRole [][32]byte, newAdminRole [][32]byte) (*PolyVaultRoleAdminChangedIterator, error) {

	var roleRule []interface{}
	for _, roleItem := range role {
		roleRule = append(roleRule, roleItem)
	}
	var previousAdminRoleRule []interface{}
	for _, previousAdminRoleItem := range previousAdminRole {
		previousAdminRoleRule = append(previousAdminRoleRule, previousAdminRoleItem)
	}
	var newAdminRoleRule []interface{}
	for _, newAdminRoleItem := range newAdminRole {
		newAdminRoleRule = append(newAdminRoleRule, newAdminRoleItem)
	}

	logs, sub, err := _PolyVault.contract.FilterLogs(opts, "RoleAdminChanged", roleRule, previousAdminRoleRule, newAdminRoleRule)
	if err != nil {
		return nil, err
	}
	return &PolyVaultRoleAdminChangedIterator{contract: _PolyVault.contract, event: "RoleAdminChanged", logs: logs, sub: sub}, nil
}

// WatchRoleAdminChanged is a free log subscription operation binding the contract event 0xbd79b86ffe0ab8e8776151514217cd7cacd52c909f66475c3af44e129f0b00ff.
//
// Solidity: event RoleAdminChanged(bytes32 indexed role, bytes32 indexed previousAdminRole, bytes32 indexed newAdminRole)
func (_PolyVault *PolyVaultFilterer) WatchRoleAdminChanged(opts *bind.WatchOpts, sink chan<- *PolyVaultRoleAdminChanged, role [][32]byte, previousAdminRole [][32]byte, newAdminRole [][32]byte) (event.Subscription, error) {

	var roleRule []interface{}
	for _, roleItem := range role {
		roleRule = append(roleRule, roleItem)
	}
	var previousAdminRoleRule []interface{}
	for _, previousAdminRoleItem := range previousAdminRole {
		previousAdminRoleRule = append(previousAdminRoleRule, previousAdminRoleItem)
	}
	var newAdminRoleRule []interface{}
	for _, newAdminRoleItem := range newAdminRole {
		newAdminRoleRule = append(newAdminRoleRule, newAdminRoleItem)
	}

	logs, sub, err := _PolyVault.contract.WatchLogs(opts, "RoleAdminChanged", roleRule, previousAdminRoleRule, newAdminRoleRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(PolyVaultRoleAdminChanged)
				if err := _PolyVault.contract.UnpackLog(event, "RoleAdminChanged", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseRoleAdminChanged is a log parse operation binding the contract event 0xbd79b86ffe0ab8e8776151514217cd7cacd52c909f66475c3af44e129f0b00ff.
//
// Solidity: event RoleAdminChanged(bytes32 indexed role, bytes32 indexed previousAdminRole, bytes32 indexed newAdminRole)
func (_PolyVault *PolyVaultFilterer) ParseRoleAdminChanged(log types.Log) (*PolyVaultRoleAdminChanged, error) {
	event := new(PolyVaultRoleAdminChanged)
	if err := _PolyVault.contract.UnpackLog(event, "RoleAdminChanged", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// PolyVaultRoleGrantedIterator is returned from FilterRoleGranted and is used to iterate over the raw logs and unpacked data for RoleGranted events raised by the PolyVault contract.
type PolyVaultRoleGrantedIterator struct {
	Event *PolyVaultRoleGranted // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *PolyVaultRoleGrantedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(PolyVaultRoleGranted)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(PolyVaultRoleGranted)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *PolyVaultRoleGrantedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *PolyVaultRoleGrantedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// PolyVaultRoleGranted represents a RoleGranted event raised by the PolyVault contract.
type PolyVaultRoleGranted struct {
	Role    [32]byte
	Account common.Address
	Sender  common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterRoleGranted is a free log retrieval operation binding the contract event 0x2f8788117e7eff1d82e926ec794901d17c78024a50270940304540a733656f0d.
//
// Solidity: event RoleGranted(bytes32 indexed role, address indexed account, address indexed sender)
func (_PolyVault *PolyVaultFilterer) FilterRoleGranted(opts *bind.FilterOpts, role [][32]byte, account []common.Address, sender []common.Address) (*PolyVaultRoleGrantedIterator, error) {

	var roleRule []interface{}
	for _, roleItem := range role {
		roleRule = append(roleRule, roleItem)
	}
	var accountRule []interface{}
	for _, accountItem := range account {
		accountRule = append(accountRule, accountItem)
	}
	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}

	logs, sub, err := _PolyVault.contract.FilterLogs(opts, "RoleGranted", roleRule, accountRule, senderRule)
	if err != nil {
		return nil, err
	}
	return &PolyVaultRoleGrantedIterator{contract: _PolyVault.contract, event: "RoleGranted", logs: logs, sub: sub}, nil
}

// WatchRoleGranted is a free log subscription operation binding the contract event 0x2f8788117e7eff1d82e926ec794901d17c78024a50270940304540a733656f0d.
//
// Solidity: event RoleGranted(bytes32 indexed role, address indexed account, address indexed sender)
func (_PolyVault *PolyVaultFilterer) WatchRoleGranted(opts *bind.WatchOpts, sink chan<- *PolyVaultRoleGranted, role [][32]byte, account []common.Address, sender []common.Address) (event.Subscription, error) {

	var roleRule []interface{}
	for _, roleItem := range role {
		roleRule = append(roleRule, roleItem)
	}
	var accountRule []interface{}
	for _, accountItem := range account {
		accountRule = append(accountRule, accountItem)
	}
	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}

	logs, sub, err := _PolyVault.contract.WatchLogs(opts, "RoleGranted", roleRule, accountRule, senderRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(PolyVaultRoleGranted)
				if err := _PolyVault.contract.UnpackLog(event, "RoleGranted", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseRoleGranted is a log parse operation binding the contract event 0x2f8788117e7eff1d82e926ec794901d17c78024a50270940304540a733656f0d.
//
// Solidity: event RoleGranted(bytes32 indexed role, address indexed account, address indexed sender)
func (_PolyVault *PolyVaultFilterer) ParseRoleGranted(log types.Log) (*PolyVaultRoleGranted, error) {
	event := new(PolyVaultRoleGranted)
	if err := _PolyVault.contract.UnpackLog(event, "RoleGranted", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// PolyVaultRoleRevokedIterator is returned from FilterRoleRevoked and is used to iterate over the raw logs and unpacked data for RoleRevoked events raised by the PolyVault contract.
type PolyVaultRoleRevokedIterator struct {
	Event *PolyVaultRoleRevoked // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *PolyVaultRoleRevokedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(PolyVaultRoleRevoked)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(PolyVaultRoleRevoked)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *PolyVaultRoleRevokedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *PolyVaultRoleRevokedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// PolyVaultRoleRevoked represents a RoleRevoked event raised by the PolyVault contract.
type PolyVaultRoleRevoked struct {
	Role    [32]byte
	Account common.Address
	Sender  common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterRoleRevoked is a free log retrieval operation binding the contract event 0xf6391f5c32d9c69d2a47ea670b442974b53935d1edc7fd64eb21e047a839171b.
//
// Solidity: event RoleRevoked(bytes32 indexed role, address indexed account, address indexed sender)
func (_PolyVault *PolyVaultFilterer) FilterRoleRevoked(opts *bind.FilterOpts, role [][32]byte, account []common.Address, sender []common.Address) (*PolyVaultRoleRevokedIterator, error) {

	var roleRule []interface{}
	for _, roleItem := range role {
		roleRule = append(roleRule, roleItem)
	}
	var accountRule []interface{}
	for _, accountItem := range account {
		accountRule = append(accountRule, accountItem)
	}
	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}

	logs, sub, err := _PolyVault.contract.FilterLogs(opts, "RoleRevoked", roleRule, accountRule, senderRule)
	if err != nil {
		return nil, err
	}
	return &PolyVaultRoleRevokedIterator{contract: _PolyVault.contract, event: "RoleRevoked", logs: logs, sub: sub}, nil
}

// WatchRoleRevoked is a free log subscription operation binding the contract event 0xf6391f5c32d9c69d2a47ea670b442974b53935d1edc7fd64eb21e047a839171b.
//
// Solidity: event RoleRevoked(bytes32 indexed role, address indexed account, address indexed sender)
func (_PolyVault *PolyVaultFilterer) WatchRoleRevoked(opts *bind.WatchOpts, sink chan<- *PolyVaultRoleRevoked, role [][32]byte, account []common.Address, sender []common.Address) (event.Subscription, error) {

	var roleRule []interface{}
	for _, roleItem := range role {
		roleRule = append(roleRule, roleItem)
	}
	var accountRule []interface{}
	for _, accountItem := range account {
		accountRule = append(accountRule, accountItem)
	}
	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}

	logs, sub, err := _PolyVault.contract.WatchLogs(opts, "RoleRevoked", roleRule, accountRule, senderRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(PolyVaultRoleRevoked)
				if err := _PolyVault.contract.UnpackLog(event, "RoleRevoked", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseRoleRevoked is a log parse operation binding the contract event 0xf6391f5c32d9c69d2a47ea670b442974b53935d1edc7fd64eb21e047a839171b.
//
// Solidity: event RoleRevoked(bytes32 indexed role, address indexed account, address indexed sender)
func (_PolyVault *PolyVaultFilterer) ParseRoleRevoked(log types.Log) (*PolyVaultRoleRevoked, error) {
	event := new(PolyVaultRoleRevoked)
	if err := _PolyVault.contract.UnpackLog(event, "RoleRevoked", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// PolyVaultStrategyDepositIterator is returned from FilterStrategyDeposit and is used to iterate over the raw logs and unpacked data for StrategyDeposit events raised by the PolyVault contract.
type PolyVaultStrategyDepositIterator struct {
	Event *PolyVaultStrategyDeposit // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *PolyVaultStrategyDepositIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(PolyVaultStrategyDeposit)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(PolyVaultStrategyDeposit)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *PolyVaultStrategyDepositIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *PolyVaultStrategyDepositIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// PolyVaultStrategyDeposit represents a StrategyDeposit event raised by the PolyVault contract.
type PolyVaultStrategyDeposit struct {
	Strategist common.Address
	Amount     *big.Int
	Raw        types.Log // Blockchain specific contextual infos
}

// FilterStrategyDeposit is a free log retrieval operation binding the contract event 0xc6f6f91a48277d76f232cc08a9a30f6b05b3fd9b92c3180c25936e17a22a1025.
//
// Solidity: event StrategyDeposit(address indexed strategist, uint256 amount)
func (_PolyVault *PolyVaultFilterer) FilterStrategyDeposit(opts *bind.FilterOpts, strategist []common.Address) (*PolyVaultStrategyDepositIterator, error) {

	var strategistRule []interface{}
	for _, strategistItem := range strategist {
		strategistRule = append(strategistRule, strategistItem)
	}

	logs, sub, err := _PolyVault.contract.FilterLogs(opts, "StrategyDeposit", strategistRule)
	if err != nil {
		return nil, err
	}
	return &PolyVaultStrategyDepositIterator{contract: _PolyVault.contract, event: "StrategyDeposit", logs: logs, sub: sub}, nil
}

// WatchStrategyDeposit is a free log subscription operation binding the contract event 0xc6f6f91a48277d76f232cc08a9a30f6b05b3fd9b92c3180c25936e17a22a1025.
//
// Solidity: event StrategyDeposit(address indexed strategist, uint256 amount)
func (_PolyVault *PolyVaultFilterer) WatchStrategyDeposit(opts *bind.WatchOpts, sink chan<- *PolyVaultStrategyDeposit, strategist []common.Address) (event.Subscription, error) {

	var strategistRule []interface{}
	for _, strategistItem := range strategist {
		strategistRule = append(strategistRule, strategistItem)
	}

	logs, sub, err := _PolyVault.contract.WatchLogs(opts, "StrategyDeposit", strategistRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(PolyVaultStrategyDeposit)
				if err := _PolyVault.contract.UnpackLog(event, "StrategyDeposit", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseStrategyDeposit is a log parse operation binding the contract event 0xc6f6f91a48277d76f232cc08a9a30f6b05b3fd9b92c3180c25936e17a22a1025.
//
// Solidity: event StrategyDeposit(address indexed strategist, uint256 amount)
func (_PolyVault *PolyVaultFilterer) ParseStrategyDeposit(log types.Log) (*PolyVaultStrategyDeposit, error) {
	event := new(PolyVaultStrategyDeposit)
	if err := _PolyVault.contract.UnpackLog(event, "StrategyDeposit", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// PolyVaultStrategyWithdrawalIterator is returned from FilterStrategyWithdrawal and is used to iterate over the raw logs and unpacked data for StrategyWithdrawal events raised by the PolyVault contract.
type PolyVaultStrategyWithdrawalIterator struct {
	Event *PolyVaultStrategyWithdrawal // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *PolyVaultStrategyWithdrawalIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(PolyVaultStrategyWithdrawal)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(PolyVaultStrategyWithdrawal)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *PolyVaultStrategyWithdrawalIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *PolyVaultStrategyWithdrawalIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// PolyVaultStrategyWithdrawal represents a StrategyWithdrawal event raised by the PolyVault contract.
type PolyVaultStrategyWithdrawal struct {
	Strategist common.Address
	Amount     *big.Int
	Raw        types.Log // Blockchain specific contextual infos
}

// FilterStrategyWithdrawal is a free log retrieval operation binding the contract event 0xd5ad0f046bd35f48b421a3e575435de38cea1980177b1c6da935d2f26049f3fa.
//
// Solidity: event StrategyWithdrawal(address indexed strategist, uint256 amount)
func (_PolyVault *PolyVaultFilterer) FilterStrategyWithdrawal(opts *bind.FilterOpts, strategist []common.Address) (*PolyVaultStrategyWithdrawalIterator, error) {

	var strategistRule []interface{}
	for _, strategistItem := range strategist {
		strategistRule = append(strategistRule, strategistItem)
	}

	logs, sub, err := _PolyVault.contract.FilterLogs(opts, "StrategyWithdrawal", strategistRule)
	if err != nil {
		return nil, err
	}
	return &PolyVaultStrategyWithdrawalIterator{contract: _PolyVault.contract, event: "StrategyWithdrawal", logs: logs, sub: sub}, nil
}

// WatchStrategyWithdrawal is a free log subscription operation binding the contract event 0xd5ad0f046bd35f48b421a3e575435de38cea1980177b1c6da935d2f26049f3fa.
//
// Solidity: event StrategyWithdrawal(address indexed strategist, uint256 amount)
func (_PolyVault *PolyVaultFilterer) WatchStrategyWithdrawal(opts *bind.WatchOpts, sink chan<- *PolyVaultStrategyWithdrawal, strategist []common.Address) (event.Subscription, error) {

	var strategistRule []interface{}
	for _, strategistItem := range strategist {
		strategistRule = append(strategistRule, strategistItem)
	}

	logs, sub, err := _PolyVault.contract.WatchLogs(opts, "StrategyWithdrawal", strategistRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(PolyVaultStrategyWithdrawal)
				if err := _PolyVault.contract.UnpackLog(event, "StrategyWithdrawal", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseStrategyWithdrawal is a log parse operation binding the contract event 0xd5ad0f046bd35f48b421a3e575435de38cea1980177b1c6da935d2f26049f3fa.
//
// Solidity: event StrategyWithdrawal(address indexed strategist, uint256 amount)
func (_PolyVault *PolyVaultFilterer) ParseStrategyWithdrawal(log types.Log) (*PolyVaultStrategyWithdrawal, error) {
	event := new(PolyVaultStrategyWithdrawal)
	if err := _PolyVault.contract.UnpackLog(event, "StrategyWithdrawal", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// PolyVaultTransferIterator is returned from FilterTransfer and is used to iterate over the raw logs and unpacked data for Transfer events raised by the PolyVault contract.
type PolyVaultTransferIterator struct {
	Event *PolyVaultTransfer // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *PolyVaultTransferIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(PolyVaultTransfer)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(PolyVaultTransfer)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *PolyVaultTransferIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *PolyVaultTransferIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// PolyVaultTransfer represents a Transfer event raised by the PolyVault contract.
type PolyVaultTransfer struct {
	From  common.Address
	To    common.Address
	Value *big.Int
	Raw   types.Log // Blockchain specific contextual infos
}

// FilterTransfer is a free log retrieval operation binding the contract event 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef.
//
// Solidity: event Transfer(address indexed from, address indexed to, uint256 value)
func (_PolyVault *PolyVaultFilterer) FilterTransfer(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*PolyVaultTransferIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _PolyVault.contract.FilterLogs(opts, "Transfer", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &PolyVaultTransferIterator{contract: _PolyVault.contract, event: "Transfer", logs: logs, sub: sub}, nil
}

// WatchTransfer is a free log subscription operation binding the contract event 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef.
//
// Solidity: event Transfer(address indexed from, address indexed to, uint256 value)
func (_PolyVault *PolyVaultFilterer) WatchTransfer(opts *bind.WatchOpts, sink chan<- *PolyVaultTransfer, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _PolyVault.contract.WatchLogs(opts, "Transfer", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(PolyVaultTransfer)
				if err := _PolyVault.contract.UnpackLog(event, "Transfer", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseTransfer is a log parse operation binding the contract event 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef.
//
// Solidity: event Transfer(address indexed from, address indexed to, uint256 value)
func (_PolyVault *PolyVaultFilterer) ParseTransfer(log types.Log) (*PolyVaultTransfer, error) {
	event := new(PolyVaultTransfer)
	if err := _PolyVault.contract.UnpackLog(event, "Transfer", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// PolyVaultUnpausedIterator is returned from FilterUnpaused and is used to iterate over the raw logs and unpacked data for Unpaused events raised by the PolyVault contract.
type PolyVaultUnpausedIterator struct {
	Event *PolyVaultUnpaused // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *PolyVaultUnpausedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(PolyVaultUnpaused)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(PolyVaultUnpaused)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *PolyVaultUnpausedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *PolyVaultUnpausedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// PolyVaultUnpaused represents a Unpaused event raised by the PolyVault contract.
type PolyVaultUnpaused struct {
	Account common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterUnpaused is a free log retrieval operation binding the contract event 0x5db9ee0a495bf2e6ff9c91a7834c1ba4fdd244a5e8aa4e537bd38aeae4b073aa.
//
// Solidity: event Unpaused(address account)
func (_PolyVault *PolyVaultFilterer) FilterUnpaused(opts *bind.FilterOpts) (*PolyVaultUnpausedIterator, error) {

	logs, sub, err := _PolyVault.contract.FilterLogs(opts, "Unpaused")
	if err != nil {
		return nil, err
	}
	return &PolyVaultUnpausedIterator{contract: _PolyVault.contract, event: "Unpaused", logs: logs, sub: sub}, nil
}

// WatchUnpaused is a free log subscription operation binding the contract event 0x5db9ee0a495bf2e6ff9c91a7834c1ba4fdd244a5e8aa4e537bd38aeae4b073aa.
//
// Solidity: event Unpaused(address account)
func (_PolyVault *PolyVaultFilterer) WatchUnpaused(opts *bind.WatchOpts, sink chan<- *PolyVaultUnpaused) (event.Subscription, error) {

	logs, sub, err := _PolyVault.contract.WatchLogs(opts, "Unpaused")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(PolyVaultUnpaused)
				if err := _PolyVault.contract.UnpackLog(event, "Unpaused", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseUnpaused is a log parse operation binding the contract event 0x5db9ee0a495bf2e6ff9c91a7834c1ba4fdd244a5e8aa4e537bd38aeae4b073aa.
//
// Solidity: event Unpaused(address account)
func (_PolyVault *PolyVaultFilterer) ParseUnpaused(log types.Log) (*PolyVaultUnpaused, error) {
	event := new(PolyVaultUnpaused)
	if err := _PolyVault.contract.UnpackLog(event, "Unpaused", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// PolyVaultUpgradedIterator is returned from FilterUpgraded and is used to iterate over the raw logs and unpacked data for Upgraded events raised by the PolyVault contract.
type PolyVaultUpgradedIterator struct {
	Event *PolyVaultUpgraded // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *PolyVaultUpgradedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(PolyVaultUpgraded)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(PolyVaultUpgraded)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *PolyVaultUpgradedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *PolyVaultUpgradedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// PolyVaultUpgraded represents a Upgraded event raised by the PolyVault contract.
type PolyVaultUpgraded struct {
	Implementation common.Address
	Raw            types.Log // Blockchain specific contextual infos
}

// FilterUpgraded is a free log retrieval operation binding the contract event 0xbc7cd75a20ee27fd9adebab32041f755214dbc6bffa90cc0225b39da2e5c2d3b.
//
// Solidity: event Upgraded(address indexed implementation)
func (_PolyVault *PolyVaultFilterer) FilterUpgraded(opts *bind.FilterOpts, implementation []common.Address) (*PolyVaultUpgradedIterator, error) {

	var implementationRule []interface{}
	for _, implementationItem := range implementation {
		implementationRule = append(implementationRule, implementationItem)
	}

	logs, sub, err := _PolyVault.contract.FilterLogs(opts, "Upgraded", implementationRule)
	if err != nil {
		return nil, err
	}
	return &PolyVaultUpgradedIterator{contract: _PolyVault.contract, event: "Upgraded", logs: logs, sub: sub}, nil
}

// WatchUpgraded is a free log subscription operation binding the contract event 0xbc7cd75a20ee27fd9adebab32041f755214dbc6bffa90cc0225b39da2e5c2d3b.
//
// Solidity: event Upgraded(address indexed implementation)
func (_PolyVault *PolyVaultFilterer) WatchUpgraded(opts *bind.WatchOpts, sink chan<- *PolyVaultUpgraded, implementation []common.Address) (event.Subscription, error) {

	var implementationRule []interface{}
	for _, implementationItem := range implementation {
		implementationRule = append(implementationRule, implementationItem)
	}

	logs, sub, err := _PolyVault.contract.WatchLogs(opts, "Upgraded", implementationRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(PolyVaultUpgraded)
				if err := _PolyVault.contract.UnpackLog(event, "Upgraded", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseUpgraded is a log parse operation binding the contract event 0xbc7cd75a20ee27fd9adebab32041f755214dbc6bffa90cc0225b39da2e5c2d3b.
//
// Solidity: event Upgraded(address indexed implementation)
func (_PolyVault *PolyVaultFilterer) ParseUpgraded(log types.Log) (*PolyVaultUpgraded, error) {
	event := new(PolyVaultUpgraded)
	if err := _PolyVault.contract.UnpackLog(event, "Upgraded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// PolyVaultWithdrawIterator is returned from FilterWithdraw and is used to iterate over the raw logs and unpacked data for Withdraw events raised by the PolyVault contract.
type PolyVaultWithdrawIterator struct {
	Event *PolyVaultWithdraw // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *PolyVaultWithdrawIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(PolyVaultWithdraw)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(PolyVaultWithdraw)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *PolyVaultWithdrawIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *PolyVaultWithdrawIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// PolyVaultWithdraw represents a Withdraw event raised by the PolyVault contract.
type PolyVaultWithdraw struct {
	Sender   common.Address
	Receiver common.Address
	Owner    common.Address
	Assets   *big.Int
	Shares   *big.Int
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterWithdraw is a free log retrieval operation binding the contract event 0xfbde797d201c681b91056529119e0b02407c7bb96a4a2c75c01fc9667232c8db.
//
// Solidity: event Withdraw(address indexed sender, address indexed receiver, address indexed owner, uint256 assets, uint256 shares)
func (_PolyVault *PolyVaultFilterer) FilterWithdraw(opts *bind.FilterOpts, sender []common.Address, receiver []common.Address, owner []common.Address) (*PolyVaultWithdrawIterator, error) {

	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}
	var receiverRule []interface{}
	for _, receiverItem := range receiver {
		receiverRule = append(receiverRule, receiverItem)
	}
	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}

	logs, sub, err := _PolyVault.contract.FilterLogs(opts, "Withdraw", senderRule, receiverRule, ownerRule)
	if err != nil {
		return nil, err
	}
	return &PolyVaultWithdrawIterator{contract: _PolyVault.contract, event: "Withdraw", logs: logs, sub: sub}, nil
}

// WatchWithdraw is a free log subscription operation binding the contract event 0xfbde797d201c681b91056529119e0b02407c7bb96a4a2c75c01fc9667232c8db.
//
// Solidity: event Withdraw(address indexed sender, address indexed receiver, address indexed owner, uint256 assets, uint256 shares)
func (_PolyVault *PolyVaultFilterer) WatchWithdraw(opts *bind.WatchOpts, sink chan<- *PolyVaultWithdraw, sender []common.Address, receiver []common.Address, owner []common.Address) (event.Subscription, error) {

	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}
	var receiverRule []interface{}
	for _, receiverItem := range receiver {
		receiverRule = append(receiverRule, receiverItem)
	}
	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}

	logs, sub, err := _PolyVault.contract.WatchLogs(opts, "Withdraw", senderRule, receiverRule, ownerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(PolyVaultWithdraw)
				if err := _PolyVault.contract.UnpackLog(event, "Withdraw", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseWithdraw is a log parse operation binding the contract event 0xfbde797d201c681b91056529119e0b02407c7bb96a4a2c75c01fc9667232c8db.
//
// Solidity: event Withdraw(address indexed sender, address indexed receiver, address indexed owner, uint256 assets, uint256 shares)
func (_PolyVault *PolyVaultFilterer) ParseWithdraw(log types.Log) (*PolyVaultWithdraw, error) {
	event := new(PolyVaultWithdraw)
	if err := _PolyVault.contract.UnpackLog(event, "Withdraw", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// PolyVaultWithdrawalCancelledIterator is returned from FilterWithdrawalCancelled and is used to iterate over the raw logs and unpacked data for WithdrawalCancelled events raised by the PolyVault contract.
type PolyVaultWithdrawalCancelledIterator struct {
	Event *PolyVaultWithdrawalCancelled // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *PolyVaultWithdrawalCancelledIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(PolyVaultWithdrawalCancelled)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(PolyVaultWithdrawalCancelled)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *PolyVaultWithdrawalCancelledIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *PolyVaultWithdrawalCancelledIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// PolyVaultWithdrawalCancelled represents a WithdrawalCancelled event raised by the PolyVault contract.
type PolyVaultWithdrawalCancelled struct {
	User   common.Address
	Shares *big.Int
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterWithdrawalCancelled is a free log retrieval operation binding the contract event 0x2eed97477f07c07ec48f8f678f4e84f7c0de55bf33f51c3dc989b13353080319.
//
// Solidity: event WithdrawalCancelled(address indexed user, uint256 shares)
func (_PolyVault *PolyVaultFilterer) FilterWithdrawalCancelled(opts *bind.FilterOpts, user []common.Address) (*PolyVaultWithdrawalCancelledIterator, error) {

	var userRule []interface{}
	for _, userItem := range user {
		userRule = append(userRule, userItem)
	}

	logs, sub, err := _PolyVault.contract.FilterLogs(opts, "WithdrawalCancelled", userRule)
	if err != nil {
		return nil, err
	}
	return &PolyVaultWithdrawalCancelledIterator{contract: _PolyVault.contract, event: "WithdrawalCancelled", logs: logs, sub: sub}, nil
}

// WatchWithdrawalCancelled is a free log subscription operation binding the contract event 0x2eed97477f07c07ec48f8f678f4e84f7c0de55bf33f51c3dc989b13353080319.
//
// Solidity: event WithdrawalCancelled(address indexed user, uint256 shares)
func (_PolyVault *PolyVaultFilterer) WatchWithdrawalCancelled(opts *bind.WatchOpts, sink chan<- *PolyVaultWithdrawalCancelled, user []common.Address) (event.Subscription, error) {

	var userRule []interface{}
	for _, userItem := range user {
		userRule = append(userRule, userItem)
	}

	logs, sub, err := _PolyVault.contract.WatchLogs(opts, "WithdrawalCancelled", userRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(PolyVaultWithdrawalCancelled)
				if err := _PolyVault.contract.UnpackLog(event, "WithdrawalCancelled", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseWithdrawalCancelled is a log parse operation binding the contract event 0x2eed97477f07c07ec48f8f678f4e84f7c0de55bf33f51c3dc989b13353080319.
//
// Solidity: event WithdrawalCancelled(address indexed user, uint256 shares)
func (_PolyVault *PolyVaultFilterer) ParseWithdrawalCancelled(log types.Log) (*PolyVaultWithdrawalCancelled, error) {
	event := new(PolyVaultWithdrawalCancelled)
	if err := _PolyVault.contract.UnpackLog(event, "WithdrawalCancelled", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// PolyVaultWithdrawalDelayUpdatedIterator is returned from FilterWithdrawalDelayUpdated and is used to iterate over the raw logs and unpacked data for WithdrawalDelayUpdated events raised by the PolyVault contract.
type PolyVaultWithdrawalDelayUpdatedIterator struct {
	Event *PolyVaultWithdrawalDelayUpdated // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *PolyVaultWithdrawalDelayUpdatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(PolyVaultWithdrawalDelayUpdated)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(PolyVaultWithdrawalDelayUpdated)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *PolyVaultWithdrawalDelayUpdatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *PolyVaultWithdrawalDelayUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// PolyVaultWithdrawalDelayUpdated represents a WithdrawalDelayUpdated event raised by the PolyVault contract.
type PolyVaultWithdrawalDelayUpdated struct {
	OldDelay *big.Int
	NewDelay *big.Int
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterWithdrawalDelayUpdated is a free log retrieval operation binding the contract event 0x9c3f1b54b1487e018f1d0593ff5cf7fb625b2df6332c974a6cc56bb358879841.
//
// Solidity: event WithdrawalDelayUpdated(uint256 oldDelay, uint256 newDelay)
func (_PolyVault *PolyVaultFilterer) FilterWithdrawalDelayUpdated(opts *bind.FilterOpts) (*PolyVaultWithdrawalDelayUpdatedIterator, error) {

	logs, sub, err := _PolyVault.contract.FilterLogs(opts, "WithdrawalDelayUpdated")
	if err != nil {
		return nil, err
	}
	return &PolyVaultWithdrawalDelayUpdatedIterator{contract: _PolyVault.contract, event: "WithdrawalDelayUpdated", logs: logs, sub: sub}, nil
}

// WatchWithdrawalDelayUpdated is a free log subscription operation binding the contract event 0x9c3f1b54b1487e018f1d0593ff5cf7fb625b2df6332c974a6cc56bb358879841.
//
// Solidity: event WithdrawalDelayUpdated(uint256 oldDelay, uint256 newDelay)
func (_PolyVault *PolyVaultFilterer) WatchWithdrawalDelayUpdated(opts *bind.WatchOpts, sink chan<- *PolyVaultWithdrawalDelayUpdated) (event.Subscription, error) {

	logs, sub, err := _PolyVault.contract.WatchLogs(opts, "WithdrawalDelayUpdated")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(PolyVaultWithdrawalDelayUpdated)
				if err := _PolyVault.contract.UnpackLog(event, "WithdrawalDelayUpdated", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseWithdrawalDelayUpdated is a log parse operation binding the contract event 0x9c3f1b54b1487e018f1d0593ff5cf7fb625b2df6332c974a6cc56bb358879841.
//
// Solidity: event WithdrawalDelayUpdated(uint256 oldDelay, uint256 newDelay)
func (_PolyVault *PolyVaultFilterer) ParseWithdrawalDelayUpdated(log types.Log) (*PolyVaultWithdrawalDelayUpdated, error) {
	event := new(PolyVaultWithdrawalDelayUpdated)
	if err := _PolyVault.contract.UnpackLog(event, "WithdrawalDelayUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// PolyVaultWithdrawalExecutedIterator is returned from FilterWithdrawalExecuted and is used to iterate over the raw logs and unpacked data for WithdrawalExecuted events raised by the PolyVault contract.
type PolyVaultWithdrawalExecutedIterator struct {
	Event *PolyVaultWithdrawalExecuted // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *PolyVaultWithdrawalExecutedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(PolyVaultWithdrawalExecuted)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(PolyVaultWithdrawalExecuted)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *PolyVaultWithdrawalExecutedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *PolyVaultWithdrawalExecutedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// PolyVaultWithdrawalExecuted represents a WithdrawalExecuted event raised by the PolyVault contract.
type PolyVaultWithdrawalExecuted struct {
	User   common.Address
	Shares *big.Int
	Assets *big.Int
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterWithdrawalExecuted is a free log retrieval operation binding the contract event 0x37ce46bc94895501203dc6abbf2b2e0d502e856e3cd90186faaba6dab7d316bb.
//
// Solidity: event WithdrawalExecuted(address indexed user, uint256 shares, uint256 assets)
func (_PolyVault *PolyVaultFilterer) FilterWithdrawalExecuted(opts *bind.FilterOpts, user []common.Address) (*PolyVaultWithdrawalExecutedIterator, error) {

	var userRule []interface{}
	for _, userItem := range user {
		userRule = append(userRule, userItem)
	}

	logs, sub, err := _PolyVault.contract.FilterLogs(opts, "WithdrawalExecuted", userRule)
	if err != nil {
		return nil, err
	}
	return &PolyVaultWithdrawalExecutedIterator{contract: _PolyVault.contract, event: "WithdrawalExecuted", logs: logs, sub: sub}, nil
}

// WatchWithdrawalExecuted is a free log subscription operation binding the contract event 0x37ce46bc94895501203dc6abbf2b2e0d502e856e3cd90186faaba6dab7d316bb.
//
// Solidity: event WithdrawalExecuted(address indexed user, uint256 shares, uint256 assets)
func (_PolyVault *PolyVaultFilterer) WatchWithdrawalExecuted(opts *bind.WatchOpts, sink chan<- *PolyVaultWithdrawalExecuted, user []common.Address) (event.Subscription, error) {

	var userRule []interface{}
	for _, userItem := range user {
		userRule = append(userRule, userItem)
	}

	logs, sub, err := _PolyVault.contract.WatchLogs(opts, "WithdrawalExecuted", userRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(PolyVaultWithdrawalExecuted)
				if err := _PolyVault.contract.UnpackLog(event, "WithdrawalExecuted", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseWithdrawalExecuted is a log parse operation binding the contract event 0x37ce46bc94895501203dc6abbf2b2e0d502e856e3cd90186faaba6dab7d316bb.
//
// Solidity: event WithdrawalExecuted(address indexed user, uint256 shares, uint256 assets)
func (_PolyVault *PolyVaultFilterer) ParseWithdrawalExecuted(log types.Log) (*PolyVaultWithdrawalExecuted, error) {
	event := new(PolyVaultWithdrawalExecuted)
	if err := _PolyVault.contract.UnpackLog(event, "WithdrawalExecuted", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// PolyVaultWithdrawalRequestedIterator is returned from FilterWithdrawalRequested and is used to iterate over the raw logs and unpacked data for WithdrawalRequested events raised by the PolyVault contract.
type PolyVaultWithdrawalRequestedIterator struct {
	Event *PolyVaultWithdrawalRequested // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *PolyVaultWithdrawalRequestedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(PolyVaultWithdrawalRequested)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(PolyVaultWithdrawalRequested)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *PolyVaultWithdrawalRequestedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *PolyVaultWithdrawalRequestedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// PolyVaultWithdrawalRequested represents a WithdrawalRequested event raised by the PolyVault contract.
type PolyVaultWithdrawalRequested struct {
	User      common.Address
	Shares    *big.Int
	Timestamp *big.Int
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterWithdrawalRequested is a free log retrieval operation binding the contract event 0x24b91f4f47caf44230a57777a9be744924e82bf666f2d5702faf97df35e60f9f.
//
// Solidity: event WithdrawalRequested(address indexed user, uint256 shares, uint256 timestamp)
func (_PolyVault *PolyVaultFilterer) FilterWithdrawalRequested(opts *bind.FilterOpts, user []common.Address) (*PolyVaultWithdrawalRequestedIterator, error) {

	var userRule []interface{}
	for _, userItem := range user {
		userRule = append(userRule, userItem)
	}

	logs, sub, err := _PolyVault.contract.FilterLogs(opts, "WithdrawalRequested", userRule)
	if err != nil {
		return nil, err
	}
	return &PolyVaultWithdrawalRequestedIterator{contract: _PolyVault.contract, event: "WithdrawalRequested", logs: logs, sub: sub}, nil
}

// WatchWithdrawalRequested is a free log subscription operation binding the contract event 0x24b91f4f47caf44230a57777a9be744924e82bf666f2d5702faf97df35e60f9f.
//
// Solidity: event WithdrawalRequested(address indexed user, uint256 shares, uint256 timestamp)
func (_PolyVault *PolyVaultFilterer) WatchWithdrawalRequested(opts *bind.WatchOpts, sink chan<- *PolyVaultWithdrawalRequested, user []common.Address) (event.Subscription, error) {

	var userRule []interface{}
	for _, userItem := range user {
		userRule = append(userRule, userItem)
	}

	logs, sub, err := _PolyVault.contract.WatchLogs(opts, "WithdrawalRequested", userRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(PolyVaultWithdrawalRequested)
				if err := _PolyVault.contract.UnpackLog(event, "WithdrawalRequested", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseWithdrawalRequested is a log parse operation binding the contract event 0x24b91f4f47caf44230a57777a9be744924e82bf666f2d5702faf97df35e60f9f.
//
// Solidity: event WithdrawalRequested(address indexed user, uint256 shares, uint256 timestamp)
func (_PolyVault *PolyVaultFilterer) ParseWithdrawalRequested(log types.Log) (*PolyVaultWithdrawalRequested, error) {
	event := new(PolyVaultWithdrawalRequested)
	if err := _PolyVault.contract.UnpackLog(event, "WithdrawalRequested", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
