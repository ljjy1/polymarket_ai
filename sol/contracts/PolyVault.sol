// SPDX-License-Identifier: MIT
pragma solidity ^0.8.28;

import "./interfaces/IPolyVault.sol";
import {AccessControlUpgradeable} from "@openzeppelin/contracts-upgradeable/access/AccessControlUpgradeable.sol";
import {ERC4626Upgradeable} from "@openzeppelin/contracts-upgradeable/token/ERC20/extensions/ERC4626Upgradeable.sol";
import {IERC20} from "@openzeppelin/contracts/token/ERC20/IERC20.sol";
import {IERC4626} from "@openzeppelin/contracts/interfaces/IERC4626.sol";
import {PausableUpgradeable} from "@openzeppelin/contracts-upgradeable/utils/PausableUpgradeable.sol";
import {ReentrancyGuardTransient} from "@openzeppelin/contracts/utils/ReentrancyGuardTransient.sol";
import {SafeERC20} from "@openzeppelin/contracts/token/ERC20/utils/SafeERC20.sol";
import {UUPSUpgradeable} from "@openzeppelin/contracts-upgradeable/proxy/utils/UUPSUpgradeable.sol";


contract PolyVault is
    ERC4626Upgradeable,           // 继承ERC4626金库标准
    AccessControlUpgradeable,     // 继承基于角色的访问控制
    PausableUpgradeable,          // 继承暂停功能
    ReentrancyGuardTransient,     // 继承重入保护
    UUPSUpgradeable,            // 继承UUPS升级模式
    IPolyVault
{
    //使用安全转账
    using SafeERC20 for IERC20;

    // 策略师角色，负责执行交易策略
    bytes32 public constant STRATEGIST_ROLE = keccak256("STRATEGIST_ROLE");
    // 守护者角色，负责紧急暂停
    bytes32 public constant GUARDIAN_ROLE = keccak256("GUARDIAN_ROLE");

    uint256 public constant MIN_WITHDRAWAL_DELAY = 1 hours;      // 最小提款延迟时间：1小时
    uint256 public constant MAX_WITHDRAWAL_DELAY = 7 days;       // 最大提款延迟时间：7天
    uint256 public constant MAX_PERFORMANCE_FEE = 2000;          // 最大业绩费：20%（2000基点）
    uint256 public constant BASIS_POINTS = 10_000;               // 基点分母：10000（100%）

    // ========== STATE ==========

    uint256 public withdrawalDelay;                              // 提款延迟时间（秒）
    mapping(address => IPolyVault.WithdrawalRequest) private _withdrawalRequests;  // 用户地址 => 提款请求映射

    uint256 public minDeposit;    // 最小存款金额
    uint256 public maxDeposit;    // 最大存款金额

    uint256 public strategyDebt;           // 策略债务：已部署到策略的USDC数量
    uint256 public maxStrategyAllocation;  // 最大策略分配比例（基点）

    uint256 public performanceFee;   // 业绩费比例（基点）
    address public feeRecipient;     // 费用接收地址


    //下面这个注释 这个注释是 Upgrades 插件识别的标准注解，我们的 constructor 里只有 _disableInitializers() ——这正是官方推荐的可升级合约标准写法，所以添加注解后就不需要测试文件做任何让步了。
    /// @custom:oz-upgrades-unsafe-allow constructor
    constructor(){
        //设置值允许初始化一次
        _disableInitializers();
    }

    /// @notice Initialize the vault with all required parameters
    /// @dev 初始化金库的所有必要参数
    /// @param _usdc USDC代币地址
    /// @param _admin 管理员地址，拥有DEFAULT_ADMIN_ROLE
    /// @param _strategist 策略师地址，拥有STRATEGIST_ROLE
    /// @param _guardian 守护者地址，拥有GUARDIAN_ROLE
    /// @param _feeRecipient 业绩费接收地址
    /// @param _withdrawalDelay 提款延迟时间（秒），范围：1小时到7天
    /// @param _maxAllocation 最大策略分配比例（基点），范围：0到10000
    /// @param _performanceFee 业绩费比例（基点），最大2000（20%）
    function initialize(
        address _usdc,
        address _admin,
        address _strategist,
        address _guardian,
        address _feeRecipient,
        uint256 _withdrawalDelay,
        uint256 _maxAllocation,
        uint256 _performanceFee
    ) external initializer {
        //验证关键地址不为零
        if(_usdc == address(0) || _admin == address(0) || _feeRecipient == address (0))
            revert ZeroAddress();
        // 验证提款延迟时间在允许范围内
        if (_withdrawalDelay < MIN_WITHDRAWAL_DELAY || _withdrawalDelay > MAX_WITHDRAWAL_DELAY)
            revert InvalidWithdrawalDelay(_withdrawalDelay);
        // 验证业绩费不超过最大值
        if (_performanceFee > MAX_PERFORMANCE_FEE)
            revert InvalidPerformanceFee(_performanceFee);
        // 验证分配比例不超过100%
        if (_maxAllocation > BASIS_POINTS)
            revert InvalidAllocation(_maxAllocation);

        // 初始化ERC20代币名称和符号
        __ERC20_init("PolyVault USDC", "pvUSDC");
        // 初始化ERC4626，设置基础资产为USDC
        __ERC4626_init(IERC20(_usdc));
        // 初始化访问控制
        __AccessControl_init();
        // 初始化暂停功能
        __Pausable_init();

        // 授予管理员角色（DEFAULT_ADMIN_ROLE是admin角色的管理员）
        _grantRole(DEFAULT_ADMIN_ROLE, _admin);
        // 授予策略师角色
        _grantRole(STRATEGIST_ROLE, _strategist);
        // 授予守护者角色
        _grantRole(GUARDIAN_ROLE, _guardian);

        // 设置状态变量
        withdrawalDelay = _withdrawalDelay;          // 提款延迟时间
        maxStrategyAllocation = _maxAllocation;      // 最大策略分配比例
        performanceFee = _performanceFee;            // 业绩费比例
        feeRecipient = _feeRecipient;                // 费用接收地址

        // 设置默认存款限额
        minDeposit = 1e6;          // 最小存款：1 USDC（USDC有6位小数）
        maxDeposit = 100_000e6;    // 最大存款：100,000 USDC
    }

    /// 存入assets数量的USDC 将assets数量的USDC存入当前合约 然后计算对应份额转入receiver
    /// @notice Deposit USDC into the vault with limit checks
    /// @dev 存入USDC到金库，包含金额限制检查
    /// @param assets 要存入的USDC数量（以最小单位计，1 USDC = 1e6）
    /// @param receiver 接收份额代币的地址
    /// @return 铸造的份额数量
    function deposit(
        uint256 assets,
        address receiver
    ) public override(ERC4626Upgradeable, IERC4626) whenNotPaused nonReentrant returns (uint256) {
        // 检查存款金额不低于最小限额
        if (assets < minDeposit) revert DepositBelowMinimum(assets, minDeposit);
        // 检查存款金额不超过最大限额
        if (assets > maxDeposit) revert DepositAboveMaximum(assets, maxDeposit);
        // 调用父类的deposit方法执行实际存款
        return super.deposit(assets, receiver);
    }

    /// 表明需要获取shares数量的份额 然后计算这个份额需要的USDC数量 如果receiver账户USDC数量足够 将计算后的USDC数量从receiver转入本合约 然后将shares数量的份额转入receiver
    /// @notice Mint exact shares with limit checks on the required assets
    /// @dev 铸造精确数量的份额，对所需资产金额进行检查
    /// @param shares 要铸造的份额数量
    /// @param receiver 接收份额的地址
    /// @return 所需的基础资产数量
    function mint(
        uint256 shares,
        address receiver
    ) public override(ERC4626Upgradeable, IERC4626) whenNotPaused nonReentrant returns (uint256) {
        // 预览铸造这些份额所需的资产数量
        uint256 assets = previewMint(shares);
        // 检查所需资产不低于最小限额
        if (assets < minDeposit) revert DepositBelowMinimum(assets, minDeposit);
        // 检查所需资产不超过最大限额
        if (assets > maxDeposit) revert DepositAboveMaximum(assets, maxDeposit);
        // 调用父类的mint方法执行实际铸造
        return super.mint(shares, receiver);
    }

    /// @notice Direct withdraw is disabled, use requestWithdraw + executeWithdraw
    /// @dev 直接提款已禁用，必须使用requestWithdraw和executeWithdraw两步流程
    /// @param 未使用
    /// @param 未使用
    /// @param 未使用
    /// @return 始终回滚
    function withdraw(uint256, address, address) public pure override(ERC4626Upgradeable, IERC4626) returns (uint256) {
        revert DirectWithdrawDisabled();  // 回滚并提示直接提款已禁用
    }

    /// @notice Direct redeem is disabled, use requestWithdraw + executeWithdraw
    /// @dev 直接赎回已禁用，必须使用两步流程
    /// @param 未使用
    /// @param 未使用
    /// @param 未使用
    /// @return 始终回滚
    function redeem(uint256, address, address) public pure override(ERC4626Upgradeable, IERC4626) returns (uint256) {
        revert DirectWithdrawDisabled();  // 回滚并提示直接赎回已禁用
    }

    /// @notice Returns 0 since direct withdraw is disabled
    /// @dev 由于直接提款被禁用，总是返回0
    /// @param 未使用的地址参数
    /// @return 0，表示不能直接提款
    function maxWithdraw(address) public pure override(ERC4626Upgradeable, IERC4626) returns (uint256) {
        return 0;  // 返回0表示不允许直接提款
    }

    /// @notice Returns 0 since direct redeem is disabled
    /// @dev 由于直接赎回被禁用，总是返回0
    /// @param 未使用的地址参数
    /// @return 0，表示不能直接赎回
    function maxRedeem(address) public pure override(ERC4626Upgradeable, IERC4626) returns (uint256) {
        return 0;  // 返回0表示不允许直接赎回
    }

    /// @notice Request a delayed withdrawal by locking shares in the vault
    /// @dev 通过锁定金库中的份额来请求延迟提款
    /// @param shares 要提取的份额数量
    function requestWithdraw(uint256 shares) external whenNotPaused nonReentrant {
        // 验证份额数量不为零
        if (shares == 0) revert ZeroAmount();
        // 验证没有待处理的提款请求
        if (_withdrawalRequests[msg.sender].pending) revert WithdrawalAlreadyPending();
        // 验证用户余额足够
        if (balanceOf(msg.sender) < shares)
            revert InsufficientShares(shares, balanceOf(msg.sender));

        // 将用户的份额转移到合约地址（锁定）
        _transfer(msg.sender, address(this), shares);

        // 创建提款请求记录
        _withdrawalRequests[msg.sender] = IPolyVault.WithdrawalRequest({
            shares: shares,                    // 锁定份额数量
            requestTimestamp: block.timestamp, // 当前时间戳
            pending: true                      // 标记为待处理
        });

        // 触发提款请求事件
        emit WithdrawalRequested(msg.sender, shares, block.timestamp);
    }

    /// @notice Cancel a pending withdrawal and reclaim locked shares
    /// @dev 取消待处理的提款请求并取回锁定的份额
    function cancelWithdraw() external nonReentrant {
        // 获取用户的提款请求（storage指针，可修改状态）
        IPolyVault.WithdrawalRequest storage req = _withdrawalRequests[msg.sender];
        // 验证存在待处理的提款请求
        if (!req.pending) revert NoPendingWithdrawal();

        // 保存要取回的份额数量
        uint256 shares = req.shares;
        // 删除提款请求（重置结构体）
        delete _withdrawalRequests[msg.sender];

        // 将锁定的份额从合约转回给用户
        _transfer(address(this), msg.sender, shares);

        // 触发提款取消事件
        emit WithdrawalCancelled(msg.sender, shares);
    }

    /// @notice Execute a withdrawal after the delay period has passed
    /// @dev 在延迟期过后执行提款操作
    function executeWithdraw() external nonReentrant {
        // 获取用户的提款请求
        IPolyVault.WithdrawalRequest storage req = _withdrawalRequests[msg.sender];
        // 验证存在待处理的提款请求
        if (!req.pending) revert NoPendingWithdrawal();

        // 计算可以执行提款的最早时间
        uint256 availableTime = req.requestTimestamp + withdrawalDelay;
        // 验证当前时间已达到可提款时间
        if (block.timestamp < availableTime)
            revert WithdrawalDelayNotMet(block.timestamp, availableTime);

        // 保存请求的份额数量
        uint256 requestedShares = req.shares;
        // 将份额转换为对应的资产数量
        uint256 assets = convertToAssets(requestedShares);
        // 删除提款请求
        delete _withdrawalRequests[msg.sender];

        // 获取合约当前的USDC余额
        uint256 vaultBalance = IERC20(asset()).balanceOf(address(this));
        uint256 sharesToBurn;     // 要销毁的份额数量
        uint256 assetsToTransfer; // 要转账的资产数量

        // 如果金库余额足够支付全部提款
        if (assets <= vaultBalance) {
            sharesToBurn = requestedShares;     // 销毁所有请求的份额
            assetsToTransfer = assets;          // 转账全部资产
        } else {
            // 金库余额不足（资金已部署到策略中），执行部分提款
            assetsToTransfer = vaultBalance;    // 只能转账金库中现有的余额
            sharesToBurn = convertToShares(vaultBalance);  // 按比例计算要销毁的份额
            // 确保不销毁超过请求的份额
            if (sharesToBurn > requestedShares) {
                sharesToBurn = requestedShares;
            }
            // 计算多余的份额（未被兑换的部分）
            uint256 excessShares = requestedShares - sharesToBurn;
            // 如果有多余份额，退还给用户
            if (excessShares > 0) {
                _transfer(address(this), msg.sender, excessShares);
            }
        }

        // 销毁合约持有的份额
        _burn(address(this), sharesToBurn);
        // 安全转账USDC给用户
        IERC20(asset()).safeTransfer(msg.sender, assetsToTransfer);

        // 触发提款执行事件
        emit WithdrawalExecuted(msg.sender, sharesToBurn, assetsToTransfer);
    }



    /// @notice Withdraw USDC from vault to trade on Polymarket
    /// @dev 从金库提取USDC用于Polymarket交易，仅策略师可调用
    /// @param amount 要提取的USDC数量
    function withdrawToStrategy(
        uint256 amount
    ) external onlyRole(STRATEGIST_ROLE) nonReentrant {
        // 验证提取金额不为零
        if (amount == 0) revert ZeroAmount();

        // 计算基于总资产和最大分配比例的最大允许金额
        uint256 maxAllowed = (totalAssets() * maxStrategyAllocation) / BASIS_POINTS;
        // 检查提取后策略债务是否超过最大允许值
        if (strategyDebt + amount > maxAllowed) {
            // 计算当前可用的策略额度
            uint256 available = maxAllowed > strategyDebt ? maxAllowed - strategyDebt : 0;
            revert StrategyAllocationExceeded(amount, available);
        }

        // 增加策略债务
        strategyDebt += amount;
        // 安全转账USDC给策略师
        IERC20(asset()).safeTransfer(msg.sender, amount);

        // 触发策略提款事件
        emit StrategyWithdrawal(msg.sender, amount);
    }


    /// @notice Return USDC to vault after trading, automatically distributes profit fee
    /// @dev 交易后将USDC返回到金库，自动分配利润费用
    /// @param amount 返还的USDC数量
    function depositFromStrategy(
        uint256 amount
    ) external onlyRole(STRATEGIST_ROLE) nonReentrant {
        // 验证返还金额不为零
        if (amount == 0) revert ZeroAmount();

        // 将USDC从策略师转移到金库
        IERC20(asset()).safeTransferFrom(msg.sender, address(this), amount);

        // 如果返还金额大于等于策略债务（有利润）
        if (amount >= strategyDebt) {
            // 计算利润：返还金额 - 债务
            uint256 profit = amount - strategyDebt;
            // 重置策略债务为0
            strategyDebt = 0;

            // 如果有利润且设置了业绩费且有费用接收地址
            if (profit > 0 && performanceFee > 0 && feeRecipient != address(0)) {
                // 计算业绩费
                uint256 fee = (profit * performanceFee) / BASIS_POINTS;
                // 转账业绩费给接收地址
                IERC20(asset()).safeTransfer(feeRecipient, fee);
                // 触发利润报告事件
                emit ProfitReported(profit, fee);
            }
        } else {
            // 如果返还金额小于策略债务（有亏损），减少债务
            strategyDebt -= amount;
        }

        // 触发策略存款事件
        emit StrategyDeposit(msg.sender, amount);
    }


    /// @notice Total assets includes USDC in vault + USDC deployed to strategy
    /// @dev 总资产包括金库中的USDC加上部署到策略的USDC
    /// @return 总资产数量
    function totalAssets() public view override(ERC4626Upgradeable, IERC4626) returns (uint256) {
        // 金库余额 + 策略债务（已部署但未归还的资金）
        return IERC20(asset()).balanceOf(address(this)) + strategyDebt;
    }

    /// @notice Update the withdrawal delay period
    /// @dev 更新提款延迟时间，仅管理员可调用
    /// @param _delay 新的提款延迟时间（秒）
    function setWithdrawalDelay(uint256 _delay) external onlyRole(DEFAULT_ADMIN_ROLE) {
        // 验证延迟时间在允许范围内
        if (_delay < MIN_WITHDRAWAL_DELAY || _delay > MAX_WITHDRAWAL_DELAY)
            revert InvalidWithdrawalDelay(_delay);
        // 保存旧值用于事件
        uint256 oldDelay = withdrawalDelay;
        // 更新延迟时间
        withdrawalDelay = _delay;
        // 触发更新事件
        emit WithdrawalDelayUpdated(oldDelay, _delay);
    }


    /// @notice Set maximum strategy allocation in basis points
    /// @dev 设置最大策略分配比例（基点），仅管理员可调用
    /// @param _allocation 新的分配比例（基点），10000 = 100%
    function setMaxStrategyAllocation(
        uint256 _allocation
    ) external onlyRole(DEFAULT_ADMIN_ROLE) {
        // 验证分配比例不超过100%
        if (_allocation > BASIS_POINTS) revert InvalidAllocation(_allocation);
        // 保存旧值
        uint256 oldAllocation = maxStrategyAllocation;
        // 更新分配比例
        maxStrategyAllocation = _allocation;
        // 触发更新事件
        emit MaxStrategyAllocationUpdated(oldAllocation, _allocation);
    }

    /// @notice Set performance fee in basis points (max 20%)
    /// @dev 设置业绩费比例（基点），最大20%，仅管理员可调用
    /// @param _fee 新的业绩费比例（基点）
    function setPerformanceFee(uint256 _fee) external onlyRole(DEFAULT_ADMIN_ROLE) {
        // 验证费用不超过最大值
        if (_fee > MAX_PERFORMANCE_FEE) revert InvalidPerformanceFee(_fee);
        // 保存旧值
        uint256 oldFee = performanceFee;
        // 更新业绩费
        performanceFee = _fee;
        // 触发更新事件
        emit PerformanceFeeUpdated(oldFee, _fee);
    }

    /// @notice Set the address that receives performance fees
    /// @dev 设置业绩费接收地址，仅管理员可调用
    /// @param _recipient 新的费用接收地址
    function setFeeRecipient(address _recipient) external onlyRole(DEFAULT_ADMIN_ROLE) {
        // 验证地址不为零地址
        if (_recipient == address(0)) revert ZeroAddress();
        // 保存旧值
        address oldRecipient = feeRecipient;
        // 更新接收地址
        feeRecipient = _recipient;
        // 触发更新事件
        emit FeeRecipientUpdated(oldRecipient, _recipient);
    }


    /// @notice Set min and max deposit limits
    /// @dev 设置存款限额，仅管理员可调用
    /// @param _min 最小存款金额
    /// @param _max 最大存款金额
    function setDepositLimits(
        uint256 _min,
        uint256 _max
    ) external onlyRole(DEFAULT_ADMIN_ROLE) {
        // 更新存款限额
        minDeposit = _min;
        maxDeposit = _max;
        // 触发更新事件
        emit DepositLimitsUpdated(_min, _max);
    }

    /// @notice Pause all deposits and withdrawal requests
    /// @dev 暂停所有存款和提款请求，仅守护者可调用
    function pause() external onlyRole(GUARDIAN_ROLE) {
        _pause();  // 调用Pausable的_pause方法
    }

    /// @notice Unpause the contract
    /// @dev 取消暂停合约，仅守护者可调用
    function unpause() external onlyRole(GUARDIAN_ROLE) {
        _unpause();  // 调用Pausable的_unpause方法
    }


    /// @notice Authorize UUPS upgrade (admin only)
    /// @dev 授权UUPS升级，仅管理员可调用
    /// @param 新的实现合约地址（未使用）
    function _authorizeUpgrade(address) internal override onlyRole(DEFAULT_ADMIN_ROLE) {
        // 仅检查调用者是否有管理员权限，不执行其他操作
    }

    /// @notice Get the USDC balance available in the vault (excludes strategy debt)
    /// @dev 获取金库中可用的USDC余额（不包括策略债务）
    /// @return 金库中的USDC余额
    function availableBalance() public view returns (uint256) {
        // 返回金库合约中的USDC余额（不包括部署到策略的资金）
        return IERC20(asset()).balanceOf(address(this));
    }


    /// @notice Get a user's pending withdrawal request
    /// @dev 获取用户的待处理提款请求
    /// @param user 用户地址
    /// @return WithdrawalRequest结构体，包含份额、请求时间和待处理状态
    function getWithdrawalRequest(
        address user
    ) external view returns (IPolyVault.WithdrawalRequest memory) {
        // 返回用户的提款请求信息
        return _withdrawalRequests[user];
    }

}
