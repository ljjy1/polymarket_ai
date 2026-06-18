// SPDX-License-Identifier: MIT
pragma solidity ^0.8.28;

import {IERC4626} from "@openzeppelin/contracts/interfaces/IERC4626.sol";

/// @title IPolyVault - PolyVault USDC 金库接口
/// @notice 扩展 ERC4626，增加了延迟提款、策略管理和基于角色的访问控制
interface IPolyVault is IERC4626 {
    // ========== 结构体 ==========

    /// @notice 提款请求结构体
    /// @dev 记录用户的待处理提款请求信息
    struct WithdrawalRequest {
        uint256 shares;              // 请求提款的份额数量
        uint256 requestTimestamp;    // 请求发起的时间戳
        bool pending;               // 是否处于待处理状态
    }

    // ========== 事件 ==========

    /// @notice 当用户发起提款请求时触发
    /// @param user 发起提款请求的用户地址
    /// @param shares 请求提款的份额数量
    /// @param timestamp 请求发起的时间戳
    event WithdrawalRequested(address indexed user, uint256 shares, uint256 timestamp);

    /// @notice 当提款执行成功时触发
    /// @param user 执行提款的用户地址
    /// @param shares 实际赎回的份额数量
    /// @param assets 实际获得的资产数量（USDC）
    event WithdrawalExecuted(address indexed user, uint256 shares, uint256 assets);

    /// @notice 当用户取消待处理的提款请求时触发
    /// @param user 取消提款请求的用户地址
    /// @param shares 被释放的份额数量
    event WithdrawalCancelled(address indexed user, uint256 shares);

    /// @notice 当策略管理者从金库提取资金时触发
    /// @param strategist 执行提款操作的管理者地址
    /// @param amount 提取的金额
    event StrategyWithdrawal(address indexed strategist, uint256 amount);

    /// @notice 当策略管理者将资金存入金库时触发
    /// @param strategist 执行存款操作的管理者地址
    /// @param amount 存入的金额
    event StrategyDeposit(address indexed strategist, uint256 amount);

    /// @notice 当金库报告利润时触发
    /// @param profit 产生的利润金额
    /// @param fee 从中收取的费用金额
    event ProfitReported(uint256 profit, uint256 fee);

    /// @notice 当绩效费用率更新时触发
    /// @param oldFee 旧的费率（基点）
    /// @param newFee 新的费率（基点）
    event PerformanceFeeUpdated(uint256 oldFee, uint256 newFee);

    /// @notice 当提款延迟时间更新时触发
    /// @param oldDelay 旧的延迟时间（秒）
    /// @param newDelay 新的延迟时间（秒）
    event WithdrawalDelayUpdated(uint256 oldDelay, uint256 newDelay);

    /// @notice 当最大策略分配比例更新时触发
    /// @param oldAllocation 旧的分配比例（基点）
    /// @param newAllocation 新的分配比例（基点）
    event MaxStrategyAllocationUpdated(uint256 oldAllocation, uint256 newAllocation);

    /// @notice 当存款限额更新时触发
    /// @param minDeposit 新的最小存款金额
    /// @param maxDeposit 新的最大存款金额
    event DepositLimitsUpdated(uint256 minDeposit, uint256 maxDeposit);

    /// @notice 当费用接收地址更新时触发
    /// @param oldRecipient 旧的费用接收地址
    /// @param newRecipient 新的费用接收地址
    event FeeRecipientUpdated(address oldRecipient, address newRecipient);

    // ========== 错误定义 ==========

    /// @notice 存款金额低于最小限额时抛出
    error DepositBelowMinimum(uint256 amount, uint256 minimum);

    /// @notice 存款金额超过最大限额时抛出
    error DepositAboveMaximum(uint256 amount, uint256 maximum);

    /// @notice 用户没有待处理的提款请求时抛出
    error NoPendingWithdrawal();

    /// @notice 用户已有待处理的提款请求时抛出
    error WithdrawalAlreadyPending();

    /// @notice 提款延迟时间未到时抛出
    /// @param currentTime 当前时间
    /// @param availableTime 可提款时间
    error WithdrawalDelayNotMet(uint256 currentTime, uint256 availableTime);

    /// @notice 份额不足时抛出
    error InsufficientShares(uint256 requested, uint256 available);

    /// @notice 无效的提款延迟时间时抛出
    error InvalidWithdrawalDelay(uint256 delay);

    /// @notice 无效的绩效费率时抛出
    error InvalidPerformanceFee(uint256 fee);

    /// @notice 无效的分配比例时抛出
    error InvalidAllocation(uint256 allocation);

    /// @notice 策略分配额度超限时抛出
    error StrategyAllocationExceeded(uint256 requested, uint256 available);

    /// @notice 零地址错误时抛出
    error ZeroAddress();

    /// @notice 零金额错误时抛出
    error ZeroAmount();

    /// @notice 直接提款被禁用时抛出
    error DirectWithdrawDisabled();

    // ========== 延迟提款功能 ==========

    /// @notice 发起延迟提款请求，在金库中锁定份额
    /// @dev 用户调用此函数发起提款请求，份额将被锁定直到延迟期结束
    /// @param shares 要提款的份额数量
    function requestWithdraw(uint256 shares) external;

    /// @notice 取消待处理的提款请求，收回被锁定的份额
    /// @dev 用户可以在延迟期内取消提款请求，释放被锁定的份额
    function cancelWithdraw() external;

    /// @notice 在延迟期过后执行提款
    /// @dev 用户必须在延迟期结束后才能执行提款，获得对应的资产
    function executeWithdraw() external;

    // ========== 策略管理功能 ==========

    /// @notice 从金库提取USDC用于Polymarket交易（仅策略管理者可调用）
    /// @dev 策略管理者调用此函数将资金从金库提取到策略合约进行交易
    /// @param amount 要提取的金额
    function withdrawToStrategy(uint256 amount) external;

    /// @notice 交易后将USDC返还给金库，自动分配利润费用
    /// @dev 策略管理者调用此函数将资金和利润返还给金库，系统会自动计算并扣除绩效费用
    /// @param amount 要存入的金额
    function depositFromStrategy(uint256 amount) external;

    // ========== 管理员功能 ==========

    /// @notice 更新提款延迟时间
    /// @dev 只有管理员可以调用，影响所有用户的提款等待时间
    /// @param delay 新的延迟时间（秒）
    function setWithdrawalDelay(uint256 delay) external;

    /// @notice 设置最大策略分配比例（以基点为单位）
    /// @dev 只有管理员可以调用，限制策略可使用的最大资金比例
    /// @param allocation 新的分配比例（基点，如 5000 = 50%）
    function setMaxStrategyAllocation(uint256 allocation) external;

    /// @notice 设置绩效费率（以基点为单位，最高20%）
    /// @dev 只有管理员可以调用，影响策略利润的收取比例
    /// @param fee 新的费率（基点，如 200 = 2%）
    function setPerformanceFee(uint256 fee) external;

    /// @notice 设置接收绩效费用的地址
    /// @dev 只有管理员可以调用
    /// @param recipient 新的费用接收地址
    function setFeeRecipient(address recipient) external;

    /// @notice 设置最小和最大存款限额
    /// @dev 只有管理员可以调用，控制用户单次存款的金额范围
    /// @param min 最小存款金额
    /// @param max 最大存款金额
    function setDepositLimits(uint256 min, uint256 max) external;

    // ========== 视图函数 ==========

    /// @notice 获取金库中可用的USDC余额（不包括策略中的债务）
    /// @dev 返回当前金库持有的USDC数量，不包括已部署到策略的资金
    /// @return 可用的USDC余额
    function availableBalance() external view returns (uint256);

    /// @notice 获取用户的待处理提款请求
    /// @param user 要查询的用户地址
    /// @return 用户的提款请求结构体
    function getWithdrawalRequest(address user) external view returns (WithdrawalRequest memory);

    /// @notice 当前已部署到策略中的USDC金额
    /// @return 策略债务（已部署的资金总额）
    function strategyDebt() external view returns (uint256);

    /// @notice 当前的提款延迟时间
    /// @return 延迟时间（秒）
    function withdrawalDelay() external view returns (uint256);

    /// @notice 最大策略分配比例（基点）
    /// @return 分配比例（基点）
    function maxStrategyAllocation() external view returns (uint256);

    /// @notice 绩效费率（基点）
    /// @return 费率（基点）
    function performanceFee() external view returns (uint256);

    /// @notice 接收绩效费用的地址
    /// @return 费用接收地址
    function feeRecipient() external view returns (address);
}