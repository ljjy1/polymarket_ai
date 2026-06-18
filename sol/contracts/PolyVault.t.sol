// SPDX-License-Identifier: MIT
// 许可证标识：MIT 开源协议
pragma solidity ^0.8.28;

// 导入 Foundry 测试框架的基础测试合约
import {Test} from "forge-std/src/Test.sol";
// 导入待测试的 PolyVault 金库合约
import {PolyVault} from "./PolyVault.sol";
// 导入 MockUSDC 模拟代币合约
import {MockUSDC} from "./mocks/MockUSDC.sol";
// 导入 OpenZeppelin 的 ERC1967 代理合约（UUPS 代理）
import {ERC1967Proxy} from "@openzeppelin/contracts/proxy/ERC1967/ERC1967Proxy.sol";
import {IPolyVault} from "./interfaces/IPolyVault.sol";

/**
 * @title PolyVaultTest
 * @notice PolyVault 金库合约的完整测试套件
 * @dev 使用 forge-std 测试框架，通过 ERC1967Proxy 代理部署合约进行测试
 */
contract PolyVaultTest is Test {
    // ========== 测试合约状态变量 ==========

    /// @notice PolyVault 金库代理实例（所有测试交互的目标）
    PolyVault public vault;
    /// @notice PolyVault 逻辑合约实例（用于验证代理存储分离）
    PolyVault public vaultImpl;
    /// @notice MockUSDC 模拟代币实例
    MockUSDC public usdc;

    // ---------- 测试账户地址 ----------

    /// @notice 管理员地址，拥有 DEFAULT_ADMIN_ROLE
    address public admin = address(0x1);
    /// @notice 策略师地址，拥有 STRATEGIST_ROLE
    address public strategist = address(0x2);
    /// @notice 守护者地址，拥有 GUARDIAN_ROLE
    address public guardian = address(0x3);
    /// @notice 业绩费接收地址
    address public feeRecipient = address(0x4);
    /// @notice 普通用户1
    address public user1 = address(0x10);
    /// @notice 普通用户2
    address public user2 = address(0x11);

    // ---------- 测试常量 ----------

    /// @notice 每个用户的初始 USDC 余额：1000 USDC
    uint256 public constant INITIAL_USER_BALANCE = 1000e6;
    /// @notice 默认提款延迟时间：1小时
    uint256 public constant DEFAULT_WITHDRAWAL_DELAY = 1 hours;
    /// @notice 默认最大策略分配比例：50%（5000基点）
    uint256 public constant DEFAULT_MAX_ALLOCATION = 5000;
    /// @notice 默认业绩费率：10%（1000基点）
    uint256 public constant DEFAULT_PERFORMANCE_FEE = 1000;

    // ========== 事件定义（用于 vm.expectEmit 事件断言验证） ==========

    /// @notice 提款请求事件
    event WithdrawalRequested(address indexed user, uint256 shares, uint256 timestamp);
    /// @notice 提款执行事件
    event WithdrawalExecuted(address indexed user, uint256 shares, uint256 assets);
    /// @notice 提款取消事件
    event WithdrawalCancelled(address indexed user, uint256 shares);
    /// @notice 策略提款事件（金库 → 策略）
    event StrategyWithdrawal(address indexed strategist, uint256 amount);
    /// @notice 策略存款事件（策略 → 金库）
    event StrategyDeposit(address indexed strategist, uint256 amount);
    /// @notice 利润报告事件
    event ProfitReported(uint256 profit, uint256 fee);
    /// @notice 业绩费更新事件
    event PerformanceFeeUpdated(uint256 oldFee, uint256 newFee);
    /// @notice 提款延迟更新事件
    event WithdrawalDelayUpdated(uint256 oldDelay, uint256 newDelay);
    /// @notice 最大策略分配更新事件
    event MaxStrategyAllocationUpdated(uint256 oldAllocation, uint256 newAllocation);
    /// @notice 存款限制更新事件
    event DepositLimitsUpdated(uint256 minDeposit, uint256 maxDeposit);
    /// @notice 费用接收地址更新事件
    event FeeRecipientUpdated(address oldRecipient, address newRecipient);

    // ========== 辅助函数 ==========

    /**
     * @notice 部署 PolyVault 逻辑合约 + ERC1967Proxy，并在构造函数中执行初始化
     * @dev 将 initialize 调用编码为 data 传入代理构造函数，类似实际部署流程
     * @param usdc_ USDC 代币地址
     * @param admin_ 管理员地址
     * @param strategist_ 策略师地址
     * @param guardian_ 守护者地址
     * @param feeRecipient_ 费用接收地址
     * @param withdrawalDelay_ 提款延迟时间
     * @param maxAllocation_ 最大策略分配比例
     * @param performanceFee_ 业绩费率
     * @return 返回代理地址，转换为 PolyVault 接口
     */
    function _deployVaultProxy(
        address usdc_,
        address admin_,
        address strategist_,
        address guardian_,
        address feeRecipient_,
        uint256 withdrawalDelay_,
        uint256 maxAllocation_,
        uint256 performanceFee_
    ) internal returns (PolyVault) {
        // 1. 部署逻辑合约（UUPS 升级的目标实现）
        PolyVault impl = new PolyVault();

        // 2. 部署 ERC1967Proxy，将 initialize 编码为初始化 data
        ERC1967Proxy proxy = new ERC1967Proxy(
            address(impl),
            // 将 initialize 函数调用编码为 abi 数据，代理构造函数会 delegatecall 执行它
            abi.encodeWithSelector(
                PolyVault.initialize.selector,
                usdc_,
                admin_,
                strategist_,
                guardian_,
                feeRecipient_,
                withdrawalDelay_,
                maxAllocation_,
                performanceFee_
            )
        );

        // 3. 将代理地址强制转换为 PolyVault 类型返回
        return PolyVault(address(proxy));
    }

    // ========== 测试前置部署（setUp） ==========

    /**
     * @notice 每个测试方法执行前的准备工作
     * @dev 部署 MockUSDC、逻辑合约、ERC1967Proxy，完成初始化并铸造测试用 USDC
     */
    function setUp() public {
        // ---- 部署 MockUSDC ----
        // 创建一个 6 位小数的模拟 USDC 代币
        usdc = new MockUSDC();

        // ---- 部署 PolyVault 逻辑合约 ----
        // 部署逻辑合约实现，保存引用用于后续 UUPS 升级相关验证
        vaultImpl = new PolyVault();

        // ---- 部署 ERC1967Proxy 并执行初始化 ----
        // ERC1967Proxy 的构造函数接受两个参数：
        //   1. implementation：逻辑合约地址
        //   2. _data：初始化调用数据（delegatecall 执行 initialize）
        ERC1967Proxy proxy = new ERC1967Proxy(
            address(vaultImpl),
            // 将 initialize 函数选择器 + 参数编码为 calldata
            abi.encodeWithSelector(
                PolyVault.initialize.selector,
                address(usdc),
                admin,
                strategist,
                guardian,
                feeRecipient,
                DEFAULT_WITHDRAWAL_DELAY,
                DEFAULT_MAX_ALLOCATION,
                DEFAULT_PERFORMANCE_FEE
            )
        );

        // ---- 将代理地址赋值给 vault 变量 ----
        // 通过类型转换将 proxy 地址作为 PolyVault 接口使用
        // 所有对 vault 的调用都会通过 delegatecall 转发到逻辑合约
        vault = PolyVault(address(proxy));

        // 验证代理存储槽已正确设置（ERC1967 标准槽）
        bytes32 implSlot = 0x360894a13ba1a3210667c828492db98dca3e2076cc3735a920a3ca505d382bbc;
        address storedImpl = address(uint160(uint256(vm.load(address(vault), implSlot))));
        assertEq(storedImpl, address(vaultImpl));

        // ---- 给测试用户铸造 USDC ----
        // user1 和 user2 各获得 1000 USDC
        usdc.mint(user1, INITIAL_USER_BALANCE);
        usdc.mint(user2, INITIAL_USER_BALANCE);
        // 策略师获得 5000 USDC（用于策略返还测试）
        usdc.mint(strategist, 5000e6);
    }

    // ==================== 初始化测试 ====================

    /**
     * @notice 测试：传入零地址作为 USDC 地址时应该 revert
     * @dev 验证 initialize 中 address(0) 检查是否生效
     */
    function test_RevertInitialize_InvalidUsdcAddress() public {
        // 部署逻辑合约（复用）
        PolyVault impl = new PolyVault();
        // 预期抛出 ZeroAddress 错误
        vm.expectRevert(IPolyVault.ZeroAddress.selector);
        // 使用零地址部署代理，构造函数中的 initialize 应该 revert
        new ERC1967Proxy(
            address(impl),
            abi.encodeWithSelector(
                PolyVault.initialize.selector,
                address(0),    // USDC 地址为零地址，非法
                admin,
                strategist,
                guardian,
                feeRecipient,
                DEFAULT_WITHDRAWAL_DELAY,
                DEFAULT_MAX_ALLOCATION,
                DEFAULT_PERFORMANCE_FEE
            )
        );
    }

    /**
     * @notice 测试：提款延迟时间小于最小值（1小时）时应该 revert
     * @dev 验证 MIN_WITHDRAWAL_DELAY 边界检查
     */
    function test_RevertInitialize_InvalidWithdrawalDelay_TooShort() public {
        // 部署逻辑合约（复用）
        PolyVault impl = new PolyVault();
        // 预期抛出 InvalidWithdrawalDelay 错误，参数为 30 minutes
        vm.expectRevert(
            abi.encodeWithSelector(IPolyVault.InvalidWithdrawalDelay.selector, 30 minutes)
        );
        // 30分钟 < 1小时的最小限制，应该 revert
        new ERC1967Proxy(
            address(impl),
            abi.encodeWithSelector(
                PolyVault.initialize.selector,
                address(usdc),
                admin,
                strategist,
                guardian,
                feeRecipient,
                30 minutes,    // 小于 1 小时的强制最小值
                DEFAULT_MAX_ALLOCATION,
                DEFAULT_PERFORMANCE_FEE
            )
        );
    }

    /**
     * @notice 测试：提款延迟时间大于最大值（7天）时应该 revert
     * @dev 验证 MAX_WITHDRAWAL_DELAY 边界检查
     */
    function test_RevertInitialize_InvalidWithdrawalDelay_TooLong() public {
        // 部署逻辑合约（复用）
        PolyVault impl = new PolyVault();
        // 预期抛出 InvalidWithdrawalDelay 错误，参数为 8 days
        vm.expectRevert(
            abi.encodeWithSelector(IPolyVault.InvalidWithdrawalDelay.selector, 8 days)
        );
        // 8天 > 7天的最大限制，应该 revert
        new ERC1967Proxy(
            address(impl),
            abi.encodeWithSelector(
                PolyVault.initialize.selector,
                address(usdc),
                admin,
                strategist,
                guardian,
                feeRecipient,
                8 days,    // 大于 7 天的强制最大值
                DEFAULT_MAX_ALLOCATION,
                DEFAULT_PERFORMANCE_FEE
            )
        );
    }

    /**
     * @notice 测试：业绩费率超过最大值（20%）时应该 revert
     * @dev 验证 MAX_PERFORMANCE_FEE = 2000（20%）的边界检查
     */
    function test_RevertInitialize_PerformanceFeeExceedsMax() public {
        // 部署逻辑合约（复用）
        PolyVault impl = new PolyVault();
        // 预期抛出 InvalidPerformanceFee 错误，参数为 2001
        vm.expectRevert(
            abi.encodeWithSelector(IPolyVault.InvalidPerformanceFee.selector, 2001)
        );
        // 2001基点 > 2000基点的最大限制，应该 revert
        new ERC1967Proxy(
            address(impl),
            abi.encodeWithSelector(
                PolyVault.initialize.selector,
                address(usdc),
                admin,
                strategist,
                guardian,
                feeRecipient,
                DEFAULT_WITHDRAWAL_DELAY,
                DEFAULT_MAX_ALLOCATION,
                2001    // 超过最大业绩费 20% 的限制
            )
        );
    }

    /**
     * @notice 测试：策略分配比例超过 100%（10000基点）时应该 revert
     * @dev 验证 BASIS_POINTS 边界检查
     */
    function test_RevertInitialize_AllocationExceedsBasisPoints() public {
        // 部署逻辑合约（复用）
        PolyVault impl = new PolyVault();
        // 预期抛出 InvalidAllocation 错误，参数为 10001
        vm.expectRevert(
            abi.encodeWithSelector(IPolyVault.InvalidAllocation.selector, 10001)
        );
        // 10001基点 > 10000基点的最大限制（超过100%），应该 revert
        new ERC1967Proxy(
            address(impl),
            abi.encodeWithSelector(
                PolyVault.initialize.selector,
                address(usdc),
                admin,
                strategist,
                guardian,
                feeRecipient,
                DEFAULT_WITHDRAWAL_DELAY,
                10001,    // 超过 100% 的限制
                DEFAULT_PERFORMANCE_FEE
            )
        );
    }

    /**
     * @notice 测试：金库初始化成功，所有参数设置正确
     * @dev 验证初始化后 name、symbol、asset 等参数与输入一致
     */
    function test_Initialize_Success() public view {
        // 验证代币名称
        assertEq(vault.name(), "PolyVault USDC");
        // 验证代币符号
        assertEq(vault.symbol(), "pvUSDC");
        // 验证底层资产地址为 usdc
        assertEq(address(vault.asset()), address(usdc));
        // 验证提款延迟时间
        assertEq(vault.withdrawalDelay(), DEFAULT_WITHDRAWAL_DELAY);
        // 验证最大策略分配比例
        assertEq(vault.maxStrategyAllocation(), DEFAULT_MAX_ALLOCATION);
        // 验证业绩费率
        assertEq(vault.performanceFee(), DEFAULT_PERFORMANCE_FEE);
        // 验证费用接收地址
        assertEq(vault.feeRecipient(), feeRecipient);
        // 验证最小存款金额（默认 1 USDC）
        assertEq(vault.minDeposit(), 1e6);
        // 验证最大存款金额（默认 100,000 USDC）
        assertEq(vault.maxDeposit(), 100_000e6);
    }

    /**
     * @notice 测试：初始化后的角色分配正确
     * @dev 验证 admin、strategist、guardian 分别拥有对应的角色
     */
    function test_Initialize_RolesCorrect() public view {
        // 验证管理员拥有 DEFAULT_ADMIN_ROLE
        assertTrue(vault.hasRole(vault.DEFAULT_ADMIN_ROLE(), admin));
        // 验证策略师拥有 STRATEGIST_ROLE
        assertTrue(vault.hasRole(vault.STRATEGIST_ROLE(), strategist));
        // 验证守护者拥有 GUARDIAN_ROLE
        assertTrue(vault.hasRole(vault.GUARDIAN_ROLE(), guardian));
    }

    /**
     * @notice 测试：不同代理实例的存储是相互隔离的
     * @dev 验证每个代理有独立的存储空间，互不影响
     */
    function test_Proxy_StorageIsSeparateFromImplementation() public {
        // 给一个随机地址铸造 USDC（供第二个代理使用）
        usdc.mint(address(0x999), 1000e6);

        // 部署第二个独立的代理实例
        PolyVault vault2 = _deployVaultProxy(
            address(usdc),
            admin,
            strategist,
            guardian,
            feeRecipient,
            DEFAULT_WITHDRAWAL_DELAY,
            DEFAULT_MAX_ALLOCATION,
            DEFAULT_PERFORMANCE_FEE
        );

        // 验证 vault1 和 vault2 是不同的合约地址
        assertFalse(address(vault) == address(vault2));

        // vault2 刚部署没有任何存款，totalAssets 为 0
        assertEq(vault2.totalAssets(), 0);

        // vault1 在 setUp 中也未存款，所以也是 0
        assertEq(vault.totalAssets(), 0);
    }

    /**
     * @notice 测试：不能对已初始化的代理再次调用 initialize
     * @dev 验证 initializer modifier 防止重复初始化
     */
    function test_RevertInitialize_AlreadyInitialized() public {
        // 预期 revert（Initializable 合约的重复初始化保护）
        vm.expectRevert();
        // 再次调用 initialize 应该失败
        vault.initialize(
            address(usdc),
            admin,
            strategist,
            guardian,
            feeRecipient,
            DEFAULT_WITHDRAWAL_DELAY,
            DEFAULT_MAX_ALLOCATION,
            DEFAULT_PERFORMANCE_FEE
        );
    }

    // ==================== 存款测试 ====================

    /**
     * @notice 测试：存款金额低于最小限额时应该 revert
     * @dev 验证 minDeposit 边界检查（1 USDC）
     */
    function test_RevertDeposit_BelowMinimum() public {
        // 存款金额为 0.5 USDC，低于 1 USDC 的最小限额
        uint256 amount = 0.5e6;

        // 预期抛出 DepositBelowMinimum 错误，包含金额和最小值
        vm.expectRevert(
            abi.encodeWithSelector(IPolyVault.DepositBelowMinimum.selector, amount, 1e6)
        );
        // 执行存款应 revert
        vm.prank(user1);
        vault.deposit(amount, user1);
    }

    /**
     * @notice 测试：存款金额超过最大限额时应该 revert
     * @dev 验证 maxDeposit 边界检查（100,000 USDC）
     */
    function test_RevertDeposit_AboveMaximum() public {
        // 存款金额为 200,000 USDC，超过 100,000 USDC 的最大限额
        uint256 amount = 200_000e6;

        // 预期抛出 DepositAboveMaximum 错误
        vm.expectRevert(
            abi.encodeWithSelector(IPolyVault.DepositAboveMaximum.selector, amount, 100_000e6)
        );
        // 执行存款应 revert
        vm.prank(user1);
        vault.deposit(amount, user1);
    }

    /**
     * @notice 测试：合约暂停时存款应该 revert
     * @dev 验证 whenNotPaused modifier 在暂停状态下阻止存款
     */
    function test_RevertDeposit_WhenPaused() public {
        // guardian 调用 pause 暂停合约
        vm.prank(guardian);
        vault.pause();

        // 预期 revert（Pausable 的 whenNotPaused 校验）
        vm.expectRevert();
        // 暂停状态下存款应 revert
        vm.prank(user1);
        vault.deposit(100e6, user1);
    }

    /**
     * @notice 测试：成功存款后检查状态变化
     * @dev 验证金库余额、用户份额和初始 1:1 汇率
     */
    function test_Deposit_Success() public {
        // 存款金额 100 USDC
        uint256 amount = 100e6;

        // 以 user1 身份执行操作
        vm.startPrank(user1);
        // 授权金库从 user1 转账 USDC
        usdc.approve(address(vault), amount);
        // 执行存款，获得对应份额
        uint256 shares = vault.deposit(amount, user1);
        vm.stopPrank();

        // 验证金库的 USDC 余额增加
        assertEq(usdc.balanceOf(address(vault)), amount);
        // 验证用户获得相应份额
        assertEq(vault.balanceOf(user1), shares);
        // 验证初始汇率为 1:1（1 USDC = 1 pvUSDC）
        assertEq(shares, amount);
    }

    /**
     * @notice 测试：存入最小金额（1 USDC）应该成功
     * @dev 验证最小存款边界值可以通过
     */
    function test_Deposit_MinimumAmount() public {
        // 最小存款金额 1 USDC
        uint256 amount = 1e6;

        // 以 user1 身份执行操作
        vm.startPrank(user1);
        // 授权金库使用 USDC
        usdc.approve(address(vault), amount);
        // 存款 1 USDC，获得份额
        uint256 shares = vault.deposit(amount, user1);
        vm.stopPrank();

        // 验证获得 1:1 的份额
        assertEq(shares, amount);
    }

    /**
     * @notice 测试：存入最大金额（100,000 USDC）应该成功
     * @dev 验证最大存款边界值可以通过
     */
    function test_Deposit_MaximumAmount() public {
        // 最大存款金额 100,000 USDC
        uint256 amount = 100_000e6;
        // 额外给 user1 铸造足够多的 USDC（原有 1000 不够）
        usdc.mint(user1, amount);

        // 以 user1 身份执行操作
        vm.startPrank(user1);
        // 授权金库使用 USDC
        usdc.approve(address(vault), amount);
        // 存款 100,000 USDC
        uint256 shares = vault.deposit(amount, user1);
        vm.stopPrank();

        // 验证获得 1:1 的份额
        assertEq(shares, amount);
    }

    /**
     * @notice 测试：多个用户分别存款，验证各自余额正确
     * @dev 验证多用户场景下份额和总资产计算正确
     */
    function test_Deposit_MultipleUsers() public {
        // user1 存款 50 USDC
        uint256 amount1 = 50e6;
        // user2 存款 75 USDC
        uint256 amount2 = 75e6;

        // user1 存款操作
        vm.startPrank(user1);
        usdc.approve(address(vault), amount1);
        vault.deposit(amount1, user1);
        vm.stopPrank();

        // user2 存款操作
        vm.startPrank(user2);
        usdc.approve(address(vault), amount2);
        vault.deposit(amount2, user2);
        vm.stopPrank();

        // 验证 user1 的份额
        assertEq(vault.balanceOf(user1), amount1);
        // 验证 user2 的份额
        assertEq(vault.balanceOf(user2), amount2);
        // 验证金库总 USDC 余额
        assertEq(usdc.balanceOf(address(vault)), amount1 + amount2);
    }

    // ==================== 铸造份额测试 ====================

    /**
     * @notice 测试：铸造份额时所需资产低于最小限额应该 revert
     * @dev 验证 mint 路径的 minDeposit 检查
     */
    function test_RevertMint_BelowMinimum() public {
        // 预期抛出 DepositBelowMinimum 错误
        vm.expectRevert(
            abi.encodeWithSelector(IPolyVault.DepositBelowMinimum.selector, 1, 1e6)
        );
        // 铸造 1 份额，需要约 1 wei USDC，远小于最小限额
        vault.mint(1, user1);
    }

    /**
     * @notice 测试：铸造份额时所需资产超过最大限额应该 revert
     * @dev 验证 mint 路径的 maxDeposit 检查
     */
    function test_RevertMint_AboveMaximum() public {
        // 铸造 200,000 份额，所需资产超过最大限额
        vm.expectRevert(
            abi.encodeWithSelector(IPolyVault.DepositAboveMaximum.selector, 200_000e6, 100_000e6)
        );
        // previewMint(200000e6) 会返回 200,000 USDC > 100,000 最大限额
        vault.mint(200_000e6, user1);
    }

    /**
     * @notice 测试：暂停时铸造份额应该 revert
     * @dev 验证暂停状态对 mint 的限制
     */
    function test_RevertMint_WhenPaused() public {
        // guardian 暂停合约
        vm.prank(guardian);
        vault.pause();

        // 预期 revert
        vm.expectRevert();
        // 铸造 100 份额应失败
        vault.mint(100e6, user1);
    }

    /**
     * @notice 测试：成功铸造指定份额
     * @dev 验证 mint 基础流程
     */
    function test_Mint_Success() public {
        // 目标份额 100 pvUSDC
        uint256 shares = 100e6;

        // 以 user1 身份操作
        vm.startPrank(user1);
        // 授权金库使用 USDC（100e6 USDC ≈ 100e6 份额）
        usdc.approve(address(vault), shares);
        // 铸造 100 份额
        vault.mint(shares, user1);
        vm.stopPrank();

        // 验证用户拥有 100 份额
        assertEq(vault.balanceOf(user1), shares);
    }

    /**
     * @notice 测试：使用 previewMint 精确计算所需资产后铸造指定份额
     * @dev 验证 mint 和 previewMint 的一致性
     */
    function test_Mint_PreciseShares() public {
        // 目标份额 50 pvUSDC
        uint256 targetShares = 50e6;

        // 预览铸造 50 份额所需的 USDC 数量
        uint256 requiredAssets = vault.previewMint(targetShares);
        // 给 user1 铸造恰好所需的 USDC
        usdc.mint(user1, requiredAssets);

        // 以 user1 身份操作
        vm.startPrank(user1);
        // 授权金库使用所需数量的 USDC
        usdc.approve(address(vault), requiredAssets);
        // 铸造 50 份额
        vault.mint(targetShares, user1);
        vm.stopPrank();

        // 验证用户拥有 50 份额
        assertEq(vault.balanceOf(user1), targetShares);
    }

    // ==================== 直接提款禁用测试 ====================

    /**
     * @notice 测试：直接调用 withdraw 应该 revert
     * @dev 验证 withdraw 被禁用，必须使用 requestWithdraw + executeWithdraw
     */
    function test_RevertWithdraw_Disabled() public {
        // 先让 user1 存款 100 USDC
        vm.startPrank(user1);
        usdc.approve(address(vault), 100e6);
        vault.deposit(100e6, user1);

        // 预期抛出 DirectWithdrawDisabled 错误
        vm.expectRevert(IPolyVault.DirectWithdrawDisabled.selector);
        // 直接提款应该被拒绝
        vault.withdraw(100e6, user1, user1);
        vm.stopPrank();
    }

    /**
     * @notice 测试：直接调用 redeem 应该 revert
     * @dev 验证 redeem 被禁用，必须使用两步提款流程
     */
    function test_RevertRedeem_Disabled() public {
        // 先让 user1 存款 100 USDC
        vm.startPrank(user1);
        usdc.approve(address(vault), 100e6);
        vault.deposit(100e6, user1);

        // 预期抛出 DirectWithdrawDisabled 错误
        vm.expectRevert(IPolyVault.DirectWithdrawDisabled.selector);
        // 直接赎回应该被拒绝
        vault.redeem(100e6, user1, user1);
        vm.stopPrank();
    }

    /**
     * @notice 测试：maxWithdraw 返回 0
     * @dev 由于直接提款被禁用，maxWithdraw 应返回 0
     */
    function test_MaxWithdraw_ReturnsZero() public view {
        // 验证不可直接提款
        assertEq(vault.maxWithdraw(user1), 0);
    }

    /**
     * @notice 测试：maxRedeem 返回 0
     * @dev 由于直接赎回被禁用，maxRedeem 应返回 0
     */
    function test_MaxRedeem_ReturnsZero() public view {
        // 验证不可直接赎回
        assertEq(vault.maxRedeem(user1), 0);
    }

    // ==================== 请求提款测试 ====================

    /**
     * @notice 测试：请求提款数量为 0 时应该 revert
     * @dev 验证零金额检查
     */
    function test_RevertRequestWithdraw_ZeroAmount() public {
        // 预期抛出 ZeroAmount 错误
        vm.expectRevert(IPolyVault.ZeroAmount.selector);
        // 请求提款 0 份额应 revert
        vault.requestWithdraw(0);
    }

    /**
     * @notice 测试：已有待处理提款请求时再次请求应该 revert
     * @dev 验证 WithdrawalAlreadyPending 错误
     */
    function test_RevertRequestWithdraw_AlreadyPending() public {
        // 先让 user1 存款 100 USDC
        vm.startPrank(user1);
        usdc.approve(address(vault), 100e6);
        vault.deposit(100e6, user1);

        // 第一次请求提款 50 份额，成功
        vault.requestWithdraw(50e6);

        // 预期抛出 WithdrawalAlreadyPending 错误
        vm.expectRevert(IPolyVault.WithdrawalAlreadyPending.selector);
        // 再次请求提款应该 revert
        vault.requestWithdraw(30e6);
        vm.stopPrank();
    }

    /**
     * @notice 测试：请求提款数量超过余额时应该 revert
     * @dev 验证 InsufficientShares 错误
     */
    function test_RevertRequestWithdraw_InsufficientShares() public {
        // 用户存款 50 USDC，只有 50 份额
        vm.startPrank(user1);
        usdc.approve(address(vault), 50e6);
        vault.deposit(50e6, user1);
        vm.stopPrank();

        // 预期抛出 InsufficientShares 错误，请求 100，可用 50
        vm.expectRevert(
            abi.encodeWithSelector(IPolyVault.InsufficientShares.selector, 100e6, 50e6)
        );
        // 请求提款 100 份额（超过拥有的 50 份额）
        vm.prank(user1);
        vault.requestWithdraw(100e6);
    }

    /**
     * @notice 测试：暂停时请求提款应该 revert
     * @dev 验证暂停状态对 requestWithdraw 的限制
     */
    function test_RevertRequestWithdraw_WhenPaused() public {
        // 用户先存款 100 USDC
        vm.startPrank(user1);
        usdc.approve(address(vault), 100e6);
        vault.deposit(100e6, user1);
        vm.stopPrank();

        // guardian 暂停合约
        vm.prank(guardian);
        vault.pause();

        // user1 尝试请求提款
        vm.prank(user1);
        // 预期 revert
        vm.expectRevert();
        vault.requestWithdraw(50e6);
    }

    /**
     * @notice 测试：请求提款成功后份额锁定到金库合约
     * @dev 验证提款请求结构体字段和份额转移
     */
    function test_RequestWithdraw_Success() public {
        // 用户存款 100 USDC
        vm.startPrank(user1);
        usdc.approve(address(vault), 100e6);
        vault.deposit(100e6, user1);
        vm.stopPrank();
        // 记录请求前的用户份额和金库份额
        uint256 userSharesBefore = vault.balanceOf(user1);
        uint256 vaultSharesBefore = vault.balanceOf(address(vault));

        // user1 请求提款 50 份额
        vm.prank(user1);
        vault.requestWithdraw(50e6);

        // 获取用户的提款请求信息
        IPolyVault.WithdrawalRequest memory req = vault.getWithdrawalRequest(user1);

        // 验证请求份额为 50
        assertEq(req.shares, 50e6);
        // 验证请求处于待处理状态
        assertTrue(req.pending);
        // 验证请求时间戳已设置（大于 0）
        assertGt(req.requestTimestamp, 0);
        // 验证用户份额减少 50（已锁定到金库）
        assertEq(vault.balanceOf(user1), userSharesBefore - 50e6);
        // 验证金库持有的份额增加 50
        assertEq(vault.balanceOf(address(vault)), vaultSharesBefore + 50e6);
    }

    /**
     * @notice 测试：请求提款成功时发出 WithdrawalRequested 事件
     * @dev 验证事件数据的正确性
     */
    function test_RequestWithdraw_EmitsEvent() public {
        // 用户存款 100 USDC
        vm.startPrank(user1);
        usdc.approve(address(vault), 100e6);
        vault.deposit(100e6, user1);
        vm.stopPrank();

        // 以 user1 身份操作
        vm.prank(user1);
        // 预期发出 WithdrawalRequested 事件，包含 user1、50 份额、当前时间戳
        vm.expectEmit(true, true, true, true);
        emit WithdrawalRequested(user1, 50e6, block.timestamp);
        // 请求提款 50 份额
        vault.requestWithdraw(50e6);
    }

    /**
     * @notice 测试：余额为 0 的用户请求提款时应该 revert
     * @dev 验证空用户提款的边界情况
     */
    function test_RevertRequestWithdraw_ZeroBalanceUser() public {
        // 一个从未存款的地址
        address emptyUser = address(0x99);

        // 预期抛出 InsufficientShares 错误
        vm.expectRevert(
            abi.encodeWithSelector(IPolyVault.InsufficientShares.selector, 1e6, 0)
        );
        // 空用户尝试提款应 revert
        vm.prank(emptyUser);
        vault.requestWithdraw(1e6);
    }

    // ==================== 取消提款测试 ====================

    /**
     * @notice 测试：没有待处理请求时取消提款应该 revert
     * @dev 验证 NoPendingWithdrawal 错误
     */
    function test_RevertCancelWithdraw_NoPendingRequest() public {
        // 预期抛出 NoPendingWithdrawal 错误
        vm.expectRevert(IPolyVault.NoPendingWithdrawal.selector);
        // 没有待处理请求就取消应 revert
        vault.cancelWithdraw();
    }

    /**
     * @notice 测试：成功取消提款，份额退还用户
     * @dev 验证取消后份额返还和请求记录清除
     */
    function test_CancelWithdraw_Success() public {
        // 用户存款 100 USDC
        vm.startPrank(user1);
        usdc.approve(address(vault), 100e6);
        vault.deposit(100e6, user1);
        vm.stopPrank();
        // user1 请求提款 50 份额
        vm.prank(user1);
        vault.requestWithdraw(50e6);

        // 记录取消前的用户余额
        uint256 userSharesBefore = vault.balanceOf(user1);

        // user1 取消提款
        vm.prank(user1);
        vault.cancelWithdraw();

        // 验证用户份额增加了 50（50 + 之前剩余的 50 = 100，回到了初始状态）
        assertEq(vault.balanceOf(user1), userSharesBefore + 50e6);

        // 验证提款请求的待处理状态已清除
        IPolyVault.WithdrawalRequest memory req = vault.getWithdrawalRequest(user1);
        assertFalse(req.pending);
        // 验证请求份额重置为 0
        assertEq(req.shares, 0);
    }

    /**
     * @notice 测试：取消提款成功时发出 WithdrawalCancelled 事件
     * @dev 验证事件数据正确
     */
    function test_CancelWithdraw_EmitsEvent() public {
        // 用户存款 100 USDC -> 请求提款 50 份额
        vm.startPrank(user1);
        usdc.approve(address(vault), 100e6);
        vault.deposit(100e6, user1);
        vm.stopPrank();
        vm.prank(user1);
        vault.requestWithdraw(50e6);

        // user1 取消提款，预期发出 WithdrawalCancelled 事件
        vm.prank(user1);
        vm.expectEmit(true, true, true, true);
        emit WithdrawalCancelled(user1, 50e6);
        vault.cancelWithdraw();
    }

    // ==================== 执行提款测试 ====================

    /**
     * @notice 测试：没有待处理请求时执行提款应该 revert
     * @dev 验证 NoPendingWithdrawal 错误
     */
    function test_RevertExecuteWithdraw_NoPendingRequest() public {
        // 预期抛出 NoPendingWithdrawal 错误
        vm.expectRevert(IPolyVault.NoPendingWithdrawal.selector);
        // 没有待处理请求就执行提款应 revert
        vault.executeWithdraw();
    }

    /**
     * @notice 测试：延迟时间未到时执行提款应该 revert
     * @dev 验证 WithdrawalDelayNotMet 时间检查
     */
    function test_RevertExecuteWithdraw_DelayNotMet() public {
        // 用户存款 100 USDC 并请求提款 50 份额
        vm.startPrank(user1);
        usdc.approve(address(vault), 100e6);
        vault.deposit(100e6, user1);
        vm.stopPrank();
        vm.prank(user1);
        vault.requestWithdraw(50e6);

        // user1 尝试立即执行提款（延迟时间未到）
        vm.prank(user1);
        // 预期抛出 WithdrawalDelayNotMet 错误
        vm.expectRevert(
            abi.encodeWithSelector(
                IPolyVault.WithdrawalDelayNotMet.selector,
                block.timestamp,
                block.timestamp + DEFAULT_WITHDRAWAL_DELAY
            )
        );
        vault.executeWithdraw();
    }

    /**
     * @notice 测试：延迟期满后成功执行提款
     * @dev 验证完整提款流程：请求 → 等待 → 执行
     */
    function test_ExecuteWithdraw_Success() public {
        // 用户存款 100 USDC 并请求提款 50 份额
        vm.startPrank(user1);
        usdc.approve(address(vault), 100e6);
        vault.deposit(100e6, user1);
        vm.stopPrank();
        vm.prank(user1);
        vault.requestWithdraw(50e6);

        // 使用 warp 跳过延迟时间（1小时 + 1秒）
        vm.warp(block.timestamp + DEFAULT_WITHDRAWAL_DELAY + 1);

        // 记录执行前的用户 USDC 余额
        uint256 userUsdcBefore = usdc.balanceOf(user1);

        // user1 执行提款
        vm.prank(user1);
        vault.executeWithdraw();

        // 验证用户 USDC 增加了 50 USDC
        uint256 userUsdcAfter = usdc.balanceOf(user1);
        assertEq(userUsdcAfter - userUsdcBefore, 50e6);

        // 验证提款请求状态已清除
        IPolyVault.WithdrawalRequest memory req = vault.getWithdrawalRequest(user1);
        assertFalse(req.pending);

        // 验证用户的剩余份额（50 份额已销毁，50 份额留存）
        assertEq(vault.balanceOf(user1), 50e6);
    }

    /**
     * @notice 测试：执行提款成功时发出 WithdrawalExecuted 事件
     * @dev 验证事件数据正确
     */
    function test_ExecuteWithdraw_EmitsEvent() public {
        // 用户存款 100 USDC 并请求提款 50 份额
        vm.startPrank(user1);
        usdc.approve(address(vault), 100e6);
        vault.deposit(100e6, user1);
        vm.stopPrank();
        vm.prank(user1);
        vault.requestWithdraw(50e6);

        // 跳过延迟时间
        vm.warp(block.timestamp + DEFAULT_WITHDRAWAL_DELAY + 1);

        // user1 执行提款，预期发出 WithdrawalExecuted 事件
        vm.prank(user1);
        vm.expectEmit(true, true, true, true);
        emit WithdrawalExecuted(user1, 50e6, 50e6);
        vault.executeWithdraw();
    }

    /**
     * @notice 测试：金库余额不足时执行部分提款
     * @dev 验证资金被策略占用时，只能提取金库现有余额
     */
    function test_ExecuteWithdraw_PartialWhenInsufficientFunds() public {
        // 用户存款 100 USDC 并请求提款全部 100 份额
        vm.startPrank(user1);
        usdc.approve(address(vault), 100e6);
        vault.deposit(100e6, user1);
        vm.stopPrank();
        vm.prank(user1);
        vault.requestWithdraw(100e6);

        // 策略师从金库提取 80 USDC，金库只剩 20 USDC
        // 先提高最大策略分配比例到 100%
        vm.prank(admin);
        vault.setMaxStrategyAllocation(10000);
        vm.prank(strategist);
        vault.withdrawToStrategy(80e6);

        // 跳过延迟时间
        vm.warp(block.timestamp + DEFAULT_WITHDRAWAL_DELAY + 1);

        // 记录执行前的用户 USDC 余额
        uint256 userUsdcBefore = usdc.balanceOf(user1);

        // user1 执行提款
        vm.prank(user1);
        vault.executeWithdraw();

        // 由于金库只有 20 USDC，用户只能提取 20 USDC
        assertEq(usdc.balanceOf(user1) - userUsdcBefore, 20e6);

        // 部分份额销毁，剩余份额应退还给用户
        assertGt(vault.balanceOf(user1), 0);
    }

    /**
     * @notice 测试：在延迟时间精确到达时（不提前不延后）可以执行提款
     * @dev 验证时间边界条件：block.timestamp == availableTime 时允许执行
     */
    function test_ExecuteWithdraw_AtExactDelaySecond() public {
        // 用户存款 100 USDC 并请求提款 50 份额
        vm.startPrank(user1);
        usdc.approve(address(vault), 100e6);
        vault.deposit(100e6, user1);
        vm.stopPrank();
        vm.prank(user1);
        vault.requestWithdraw(50e6);

        // 将时间正好推进到 request.timestamp + withdrawalDelay
        vm.warp(block.timestamp + DEFAULT_WITHDRAWAL_DELAY);

        // user1 执行提款（应该成功，因为时间刚好满足条件）
        vm.prank(user1);
        vault.executeWithdraw();

        // 验证用户收到 50 USDC（初始 1000 - 存款 100 + 提款 50 = 950）
        assertEq(usdc.balanceOf(user1), INITIAL_USER_BALANCE - 100e6 + 50e6);
    }

    // ==================== 策略提款测试 ====================

    /**
     * @notice 测试：策略提款金额为 0 时应该 revert
     * @dev 验证零金额检查
     */
    function test_RevertWithdrawToStrategy_ZeroAmount() public {
        // 预期抛出 ZeroAmount 错误
        vm.expectRevert(IPolyVault.ZeroAmount.selector);
        // 策略师提款 0 金额应 revert
        vm.prank(strategist);
        vault.withdrawToStrategy(0);
    }

    /**
     * @notice 测试：非策略师调用 withdrawToStrategy 应该 revert
     * @dev 验证 STRATEGIST_ROLE 权限控制
     */
    function test_RevertWithdrawToStrategy_NotStrategist() public {
        // 非策略师调用预期 revert（AccessControl 校验）
        vm.expectRevert();
        vault.withdrawToStrategy(50e6);
    }

    /**
     * @notice 测试：策略提款超过最大分配比例时应该 revert
     * @dev 验证 maxStrategyAllocation 限制
     */
    function test_RevertWithdrawToStrategy_AllocationExceeded() public {
        // 用户先存款 100 USDC，总资产为 100 USDC
        vm.startPrank(user1);
        usdc.approve(address(vault), 100e6);
        vault.deposit(100e6, user1);
        vm.stopPrank();

        // 策略师提取 40 USDC（占总资产 40%，在 50% 限制内）
        vm.prank(strategist);
        vault.withdrawToStrategy(40e6);

        // 策略师尝试再提取 20 USDC（累计 60%，超过 50% 限制）
        vm.prank(strategist);
        vm.expectRevert();
        vault.withdrawToStrategy(20e6);
    }

    /**
     * @notice 测试：策略师成功从金库提取 USDC
     * @dev 验证资金转移和债务记录
     */
    function test_WithdrawToStrategy_Success() public {
        // 用户存款 100 USDC
        vm.startPrank(user1);
        usdc.approve(address(vault), 100e6);
        vault.deposit(100e6, user1);
        vm.stopPrank();

        // 策略师提取 50 USDC
        vm.prank(strategist);
        vault.withdrawToStrategy(50e6);

        // 验证策略师收到 50 USDC（初始 5000 + 提取 50 = 5050）
        assertEq(usdc.balanceOf(strategist), 5000e6 + 50e6);
        // 验证策略债务增加 50
        assertEq(vault.strategyDebt(), 50e6);
        // 验证金库可用余额减少到 50 USDC
        assertEq(vault.availableBalance(), 50e6);
    }

    /**
     * @notice 测试：策略提款成功时发出 StrategyWithdrawal 事件
     * @dev 验证事件数据正确
     */
    function test_WithdrawToStrategy_EmitsEvent() public {
        // 用户存款 100 USDC
        vm.startPrank(user1);
        usdc.approve(address(vault), 100e6);
        vault.deposit(100e6, user1);
        vm.stopPrank();

        // 策略师提取 30 USDC，预期发出 StrategyWithdrawal 事件
        vm.prank(strategist);
        vm.expectEmit(true, true, true, true);
        emit StrategyWithdrawal(strategist, 30e6);
        vault.withdrawToStrategy(30e6);
    }

    /**
     * @notice 测试：策略师可以提款到最大分配限制（50%）
     * @dev 验证最大分配边界可以通过
     */
    function test_WithdrawToStrategy_FullAllocation() public {
        // 用户存款 100 USDC
        vm.startPrank(user1);
        usdc.approve(address(vault), 100e6);
        vault.deposit(100e6, user1);
        vm.stopPrank();

        // 策略师提取 50 USDC（恰好为最大分配比例 50%）
        vm.prank(strategist);
        vault.withdrawToStrategy(50e6);

        // 验证债务为 50 USDC
        assertEq(vault.strategyDebt(), 50e6);
    }

    // ==================== 策略存款测试 ====================

    /**
     * @notice 测试：策略存款金额为 0 时应该 revert
     * @dev 验证零金额检查
     */
    function test_RevertDepositFromStrategy_ZeroAmount() public {
        // 预期抛出 ZeroAmount 错误
        vm.prank(strategist);
        vm.expectRevert(IPolyVault.ZeroAmount.selector);
        // 存入 0 金额应 revert
        vault.depositFromStrategy(0);
    }

    /**
     * @notice 测试：非策略师调用 depositFromStrategy 应该 revert
     * @dev 验证 STRATEGIST_ROLE 权限控制
     */
    function test_RevertDepositFromStrategy_NotStrategist() public {
        // 非策略师调用预期 revert
        vm.expectRevert();
        vault.depositFromStrategy(100);
    }

    /**
     * @notice 测试：未授权 USDC 时策略存款应该 revert
     * @dev 验证 SafeERC20 的 safeTransferFrom 检查
     */
    function test_RevertDepositFromStrategy_WithoutApproval() public {
        // 用户存款 100 USDC
        vm.startPrank(user1);
        usdc.approve(address(vault), 100e6);
        vault.deposit(100e6, user1);
        vm.stopPrank();

        // 策略师提取 50 USDC
        vm.prank(strategist);
        vault.withdrawToStrategy(50e6);

        // 策略师没有授权金库使用 USDC，直接存款应 revert
        vm.prank(strategist);
        vm.expectRevert();
        vault.depositFromStrategy(50e6);
    }

    /**
     * @notice 测试：策略师返还本金（无利润）
     * @dev 验证本金返还后债务清零，无费用产生
     */
    function test_DepositFromStrategy_ReturnPrincipal() public {
        // 用户存款 100 USDC
        vm.startPrank(user1);
        usdc.approve(address(vault), 100e6);
        vault.deposit(100e6, user1);
        vm.stopPrank();

        // 策略师提取 50 USDC
        vm.prank(strategist);
        vault.withdrawToStrategy(50e6);

        // 记录费用接收地址的初始余额
        uint256 feeRecipientBalanceBefore = usdc.balanceOf(feeRecipient);

        // 策略师批准并返还 50 USDC（正好是本金，无利润）
        vm.prank(strategist);
        usdc.approve(address(vault), 50e6);
        vm.prank(strategist);
        vault.depositFromStrategy(50e6);

        // 验证债务清零
        assertEq(vault.strategyDebt(), 0);
        // 验证金库可用余额恢复（50 + 50 = 100 USDC）
        assertEq(vault.availableBalance(), 100e6);
        // 验证费用接收地址余额不变（无利润，不收费）
        assertEq(usdc.balanceOf(feeRecipient), feeRecipientBalanceBefore);
    }

    /**
     * @notice 测试：策略师返还资金并产生利润，业绩费被收取
     * @dev 验证利润分配和业绩费计算正确
     */
    function test_DepositFromStrategy_WithProfit() public {
        // 用户存款 100 USDC
        vm.startPrank(user1);
        usdc.approve(address(vault), 100e6);
        vault.deposit(100e6, user1);
        vm.stopPrank();

        // 策略师提取 50 USDC
        vm.prank(strategist);
        vault.withdrawToStrategy(50e6);

        // 策略师赚取了 10 USDC 利润，额外铸造 10 USDC
        usdc.mint(strategist, 10e6);

        // 策略师批准并返还 60 USDC（50 本金 + 10 利润）
        vm.prank(strategist);
        usdc.approve(address(vault), 60e6);
        vm.prank(strategist);
        vault.depositFromStrategy(60e6);

        // 验证债务清零
        assertEq(vault.strategyDebt(), 0);

        // 验证业绩费 = 10 USDC * 10% = 1 USDC
        assertEq(usdc.balanceOf(feeRecipient), 1e6);
    }

    /**
     * @notice 测试：策略存款产生利润时发出 ProfitReported 事件
     * @dev 验证 ProfitReported 事件的 profit 和 fee 参数
     */
    function test_DepositFromStrategy_EmitsProfitReported() public {
        // 用户存款 100 USDC
        vm.startPrank(user1);
        usdc.approve(address(vault), 100e6);
        vault.deposit(100e6, user1);
        vm.stopPrank();

        // 策略师提取 50 USDC
        vm.prank(strategist);
        vault.withdrawToStrategy(50e6);

        // 给策略师铸造 10 USDC 利润
        usdc.mint(strategist, 10e6);

        // 策略师批准返还 60 USDC
        vm.prank(strategist);
        usdc.approve(address(vault), 60e6);

        // 预期发出 ProfitReported 事件，利润 10e6，费用 1e6
        vm.prank(strategist);
        vm.expectEmit(true, true, true, true);
        emit ProfitReported(10e6, 1e6);
        vault.depositFromStrategy(60e6);
    }

    /**
     * @notice 测试：策略存款成功时发出 StrategyDeposit 事件
     * @dev 验证 StrategyDeposit 事件数据正确
     */
    function test_DepositFromStrategy_EmitsStrategyDeposit() public {
        // 用户存款 100 USDC
        vm.startPrank(user1);
        usdc.approve(address(vault), 100e6);
        vault.deposit(100e6, user1);
        vm.stopPrank();

        // 策略师提取 50 USDC
        vm.prank(strategist);
        vault.withdrawToStrategy(50e6);

        // 策略师批准并返还 50 USDC
        vm.prank(strategist);
        usdc.approve(address(vault), 50e6);

        // 预期发出 StrategyDeposit 事件
        vm.prank(strategist);
        vm.expectEmit(true, true, true, true);
        emit StrategyDeposit(strategist, 50e6);
        vault.depositFromStrategy(50e6);
    }

    /**
     * @notice 测试：策略亏损时部分返还，债务减少但不产生业绩费
     * @dev 验证亏损场景下债务部分扣除且不收费
     */
    function test_DepositFromStrategy_WithLoss() public {
        // 用户存款 100 USDC
        vm.startPrank(user1);
        usdc.approve(address(vault), 100e6);
        vault.deposit(100e6, user1);
        vm.stopPrank();

        // 策略师提取 50 USDC
        vm.prank(strategist);
        vault.withdrawToStrategy(50e6);

        // 策略亏损，只返还 30 USDC（不是 50）
        vm.prank(strategist);
        usdc.approve(address(vault), 30e6);
        vm.prank(strategist);
        vault.depositFromStrategy(30e6);

        // 验证债务减少到 20 USDC（50 - 30 = 20）
        assertEq(vault.strategyDebt(), 20e6);
        // 验证没有业绩费产生
        assertEq(usdc.balanceOf(feeRecipient), 0);
    }

    /**
     * @notice 测试：业绩费为 0 时，即使有利润也不收取费用
     * @dev 验证 performanceFee = 0 的场景
     */
    function test_DepositFromStrategy_ZeroFeeCollectsNoFee() public {
        // 用户存款 100 USDC
        vm.startPrank(user1);
        usdc.approve(address(vault), 100e6);
        vault.deposit(100e6, user1);
        vm.stopPrank();

        // 管理员将业绩费设为 0
        vm.prank(admin);
        vault.setPerformanceFee(0);

        // 策略师提取 50 USDC
        vm.prank(strategist);
        vault.withdrawToStrategy(50e6);

        // 给策略师铸造 20 USDC 利润
        usdc.mint(strategist, 20e6);

        // 策略师返还 70 USDC（50 本金 + 20 利润）
        vm.prank(strategist);
        usdc.approve(address(vault), 70e6);
        vm.prank(strategist);
        vault.depositFromStrategy(70e6);

        // 业绩费为 0，费用接收地址余额不应变化
        assertEq(usdc.balanceOf(feeRecipient), 0);
    }

    /**
     * @notice 测试：策略师多次提取和返还操作
     * @dev 验证多次策略操作后债务状态正确
     */
    function test_DepositFromStrategy_DoubleReturn() public {
        // 用户存款 100 USDC
        vm.startPrank(user1);
        usdc.approve(address(vault), 100e6);
        vault.deposit(100e6, user1);
        vm.stopPrank();

        // 第一次：策略师提取 30 USDC 并返还 30 USDC
        vm.prank(strategist);
        vault.withdrawToStrategy(30e6);

        vm.prank(strategist);
        usdc.approve(address(vault), 30e6);
        vm.prank(strategist);
        vault.depositFromStrategy(30e6);

        // 验证债务清零
        assertEq(vault.strategyDebt(), 0);

        // 第二次：策略师提取 20 USDC 并返还 20 USDC
        vm.prank(strategist);
        vault.withdrawToStrategy(20e6);

        vm.prank(strategist);
        usdc.approve(address(vault), 20e6);
        vm.prank(strategist);
        vault.depositFromStrategy(20e6);

        // 验证债务清零
        assertEq(vault.strategyDebt(), 0);
    }

    // ==================== 总资产测试 ====================

    /**
     * @notice 测试：只有金库余额时 totalAssets 等于金库余额
     * @dev 验证无策略债务时的总资产计算
     */
    function test_TotalAssets_OnlyVaultBalance() public {
        // 用户存款 100 USDC
        vm.startPrank(user1);
        usdc.approve(address(vault), 100e6);
        vault.deposit(100e6, user1);
        vm.stopPrank();
        // 总资产 = 金库余额 100 USDC
        assertEq(vault.totalAssets(), 100e6);
    }

    /**
     * @notice 测试：totalAssets 包含策略债务
     * @dev 验证总资产 = 金库余额 + 策略债务
     */
    function test_TotalAssets_IncludesStrategyDebt() public {
        // 用户存款 100 USDC
        vm.startPrank(user1);
        usdc.approve(address(vault), 100e6);
        vault.deposit(100e6, user1);
        vm.stopPrank();

        // 策略师提取 30 USDC
        vm.prank(strategist);
        vault.withdrawToStrategy(30e6);

        // 总资产 = 金库余额(70e6) + 策略债务(30e6) = 100e6
        assertEq(vault.totalAssets(), 100e6);
    }

    /**
     * @notice 测试：策略盈利后 totalAssets 增加
     * @dev 验证利润扣除业绩费后总资产正确增长
     */
    function test_TotalAssets_AfterProfit() public {
        // 用户存款 100 USDC
        vm.startPrank(user1);
        usdc.approve(address(vault), 100e6);
        vault.deposit(100e6, user1);
        vm.stopPrank();

        // 策略师提取 50 USDC
        vm.prank(strategist);
        vault.withdrawToStrategy(50e6);

        // 策略师赚取 10 USDC 利润
        usdc.mint(strategist, 10e6);

        // 策略师返还 60 USDC
        vm.prank(strategist);
        usdc.approve(address(vault), 60e6);
        vm.prank(strategist);
        vault.depositFromStrategy(60e6);

        // 总资产 = 100(原始) + 10(利润) - 1(业绩费10%) = 109 USDC
        assertEq(vault.totalAssets(), 109e6);
    }

    /**
     * @notice 测试：策略亏损后 totalAssets 保持不变（亏损由债务体现）
     * @dev 验证亏损场景总资产不变，债务持续记录
     */
    function test_TotalAssets_AfterLoss() public {
        // 用户存款 100 USDC
        vm.startPrank(user1);
        usdc.approve(address(vault), 100e6);
        vault.deposit(100e6, user1);
        vm.stopPrank();

        // 策略师提取 50 USDC
        vm.prank(strategist);
        vault.withdrawToStrategy(50e6);

        // 策略师亏损，只返还 30 USDC
        vm.prank(strategist);
        usdc.approve(address(vault), 30e6);
        vm.prank(strategist);
        vault.depositFromStrategy(30e6);

        // 总资产 = 金库余额(50 + 30 = 80) + 剩余债务(20) = 100 USDC
        // 总资产保持不变，因为亏损部分仍记录为策略债务
        assertEq(vault.totalAssets(), 100e6);
    }

    // ==================== 管理员功能测试 ====================

    /**
     * @notice 测试：管理员成功更新提款延迟时间
     * @dev 验证 setWithdrawalDelay 功能
     */
    function test_SetWithdrawalDelay_Success() public {
        // 新的提款延迟时间 2 小时
        uint256 newDelay = 2 hours;

        // 管理员调用，预期发出 WithdrawalDelayUpdated 事件
        vm.prank(admin);
        vm.expectEmit(true, true, true, true);
        emit WithdrawalDelayUpdated(DEFAULT_WITHDRAWAL_DELAY, newDelay);
        vault.setWithdrawalDelay(newDelay);

        // 验证延迟时间已更新
        assertEq(vault.withdrawalDelay(), newDelay);
    }

    /**
     * @notice 测试：设置过短的提款延迟时间应该 revert
     * @dev 验证最小延迟时间边界
     */
    function test_RevertSetWithdrawalDelay_TooShort() public {
        // 管理员设置 30 分钟（小于 1 小时的最小值）
        vm.prank(admin);
        vm.expectRevert(
            abi.encodeWithSelector(IPolyVault.InvalidWithdrawalDelay.selector, 30 minutes)
        );
        vault.setWithdrawalDelay(30 minutes);
    }

    /**
     * @notice 测试：设置过长的提款延迟时间应该 revert
     * @dev 验证最大延迟时间边界
     */
    function test_RevertSetWithdrawalDelay_TooLong() public {
        // 管理员设置 8 天（大于 7 天的最大值）
        vm.prank(admin);
        vm.expectRevert(
            abi.encodeWithSelector(IPolyVault.InvalidWithdrawalDelay.selector, 8 days)
        );
        vault.setWithdrawalDelay(8 days);
    }

    /**
     * @notice 测试：非管理员设置提款延迟时间应该 revert
     * @dev 验证 DEFAULT_ADMIN_ROLE 权限控制
     */
    function test_RevertSetWithdrawalDelay_NotAdmin() public {
        // 非管理员调用预期 revert
        vm.expectRevert();
        vm.prank(user1);
        vault.setWithdrawalDelay(2 hours);
    }

    /**
     * @notice 测试：管理员成功更新最大策略分配比例
     * @dev 验证 setMaxStrategyAllocation 功能
     */
    function test_SetMaxStrategyAllocation_Success() public {
        // 新的分配比例 75%
        uint256 newAllocation = 7500;

        // 管理员调用，预期发出 MaxStrategyAllocationUpdated 事件
        vm.prank(admin);
        vm.expectEmit(true, true, true, true);
        emit MaxStrategyAllocationUpdated(DEFAULT_MAX_ALLOCATION, newAllocation);
        vault.setMaxStrategyAllocation(newAllocation);

        // 验证分配比例已更新
        assertEq(vault.maxStrategyAllocation(), newAllocation);
    }

    /**
     * @notice 测试：设置超过 100% 的策略分配比例应该 revert
     * @dev 验证 BASIS_POINTS 边界
     */
    function test_RevertSetMaxStrategyAllocation_ExceedsBasisPoints() public {
        // 管理员设置 10001 基点（超过 100%）
        vm.prank(admin);
        vm.expectRevert(
            abi.encodeWithSelector(IPolyVault.InvalidAllocation.selector, 10001)
        );
        vault.setMaxStrategyAllocation(10001);
    }

    /**
     * @notice 测试：管理员成功更新业绩费率
     * @dev 验证 setPerformanceFee 功能
     */
    function test_SetPerformanceFee_Success() public {
        // 新的业绩费率 5%
        uint256 newFee = 500;

        // 管理员调用，预期发出 PerformanceFeeUpdated 事件
        vm.prank(admin);
        vm.expectEmit(true, true, true, true);
        emit PerformanceFeeUpdated(DEFAULT_PERFORMANCE_FEE, newFee);
        vault.setPerformanceFee(newFee);

        // 验证业绩费已更新
        assertEq(vault.performanceFee(), newFee);
    }

    /**
     * @notice 测试：设置超过最大业绩费率应该 revert
     * @dev 验证 MAX_PERFORMANCE_FEE = 2000 边界
     */
    function test_RevertSetPerformanceFee_ExceedsMax() public {
        // 管理员设置 2001 基点（超过 20% 的最大值）
        vm.prank(admin);
        vm.expectRevert(
            abi.encodeWithSelector(IPolyVault.InvalidPerformanceFee.selector, 2001)
        );
        vault.setPerformanceFee(2001);
    }

    /**
     * @notice 测试：管理员成功更新费用接收地址
     * @dev 验证 setFeeRecipient 功能
     */
    function test_SetFeeRecipient_Success() public {
        // 新的费用接收地址
        address newRecipient = address(0x99);

        // 管理员调用，预期发出 FeeRecipientUpdated 事件
        vm.prank(admin);
        vm.expectEmit(true, true, true, true);
        emit FeeRecipientUpdated(feeRecipient, newRecipient);
        vault.setFeeRecipient(newRecipient);

        // 验证费用接收地址已更新
        assertEq(vault.feeRecipient(), newRecipient);
    }

    /**
     * @notice 测试：设置零地址为费用接收地址应该 revert
     * @dev 验证 ZeroAddress 检查
     */
    function test_RevertSetFeeRecipient_ZeroAddress() public {
        // 管理员设置零地址
        vm.prank(admin);
        vm.expectRevert(IPolyVault.ZeroAddress.selector);
        vault.setFeeRecipient(address(0));
    }

    /**
     * @notice 测试：管理员成功更新存款限额
     * @dev 验证 setDepositLimits 功能
     */
    function test_SetDepositLimits_Success() public {
        // 新的最小存款 10 USDC，最大 50,000 USDC
        uint256 newMin = 10e6;
        uint256 newMax = 50_000e6;

        // 管理员调用，预期发出 DepositLimitsUpdated 事件
        vm.prank(admin);
        vm.expectEmit(true, true, true, true);
        emit DepositLimitsUpdated(newMin, newMax);
        vault.setDepositLimits(newMin, newMax);

        // 验证存款限额已更新
        assertEq(vault.minDeposit(), newMin);
        assertEq(vault.maxDeposit(), newMax);
    }

    /**
     * @notice 测试：更新存款限额后，可以按新限额存款
     * @dev 验证新限额生效
     */
    function test_DepositWithUpdatedLimits() public {
        // 管理员将限额设为 0.5 USDC 到 200,000 USDC
        vm.prank(admin);
        vault.setDepositLimits(0.5e6, 200_000e6);

        // 存入 0.5 USDC（之前低于旧的最小限额 1 USDC）
        uint256 amount = 0.5e6;

        vm.startPrank(user1);
        usdc.approve(address(vault), amount);
        uint256 shares = vault.deposit(amount, user1);
        vm.stopPrank();

        // 验证存款成功
        assertEq(shares, amount);
    }

    // ==================== 暂停 / 恢复测试 ====================

    /**
     * @notice 测试：守护者可以成功暂停合约
     * @dev 验证 GUARDIAN_ROLE 的暂停能力
     */
    function test_Pause_ByGuardian() public {
        // guardian 调用 pause
        vm.prank(guardian);
        vault.pause();

        // 暂停后存款应 revert
        vm.startPrank(user1);
        usdc.approve(address(vault), 100e6);
        vm.expectRevert();
        vault.deposit(100e6, user1);
        vm.stopPrank();
    }

    /**
     * @notice 测试：非守护者调用 pause 应该 revert
     * @dev 验证暂停权限控制
     */
    function test_RevertPause_NotGuardian() public {
        // 非守护者调用预期 revert
        vm.expectRevert();
        vault.pause();
    }

    /**
     * @notice 测试：守护者可以成功恢复合约
     * @dev 验证 unpause 后功能恢复正常
     */
    function test_Unpause_ByGuardian() public {
        // guardian 暂停合约
        vm.prank(guardian);
        vault.pause();

        // guardian 恢复合约
        vm.prank(guardian);
        vault.unpause();

        // 恢复后应该可以正常存款
        vm.startPrank(user1);
        usdc.approve(address(vault), 100e6);
        uint256 shares = vault.deposit(100e6, user1);
        vm.stopPrank();

        // 验证存款成功
        assertEq(shares, 100e6);
    }

    /**
     * @notice 测试：非守护者调用 unpause 应该 revert
     * @dev 验证恢复权限控制
     */
    function test_RevertUnpause_NotGuardian() public {
        // guardian 先暂停
        vm.prank(guardian);
        vault.pause();

        // 非守护者恢复应 revert
        vm.expectRevert();
        vault.unpause();
    }

    // ==================== 可用余额测试 ====================

    /**
     * @notice 测试：存款后 availableBalance 等于金库余额
     * @dev 验证无策略债务时的可用余额
     */
    function test_AvailableBalance_AfterDeposit() public {
        // 用户存款 100 USDC
        vm.startPrank(user1);
        usdc.approve(address(vault), 100e6);
        vault.deposit(100e6, user1);
        vm.stopPrank();
        // 可用余额 = 金库 USDC 余额 = 100 USDC
        assertEq(vault.availableBalance(), 100e6);
    }

    /**
     * @notice 测试：availableBalance 排除策略债务部分
     * @dev 验证已部署到策略的资金不计入可用余额
     */
    function test_AvailableBalance_ExcludesStrategyDebt() public {
        // 用户存款 100 USDC
        vm.startPrank(user1);
        usdc.approve(address(vault), 100e6);
        vault.deposit(100e6, user1);
        vm.stopPrank();

        // 策略师提取 30 USDC
        vm.prank(strategist);
        vault.withdrawToStrategy(30e6);

        // 可用余额 = 金库余额 = 70 USDC（30 USDC 在策略中，不计入）
        assertEq(vault.availableBalance(), 70e6);
    }

    // ==================== 视图函数测试 ====================

    /**
     * @notice 测试：未请求提款的用户，查询结果为默认空结构体
     * @dev 验证 getWithdrawalRequest 默认返回值
     */
    function test_GetWithdrawalRequest_DefaultIsEmpty() public view {
        // 获取 user1 的提款请求（从未请求过）
        IPolyVault.WithdrawalRequest memory req = vault.getWithdrawalRequest(user1);
        // 验证不是待处理状态
        assertFalse(req.pending);
        // 验证份额为 0
        assertEq(req.shares, 0);
    }

    /**
     * @notice 测试：请求提款后 getWithdrawalRequest 返回正确的请求信息
     * @dev 验证请求数据的持久化
     */
    function test_GetWithdrawalRequest_AfterRequest() public {
        // 用户存款 100 USDC 并请求提款 50 份额
        vm.startPrank(user1);
        usdc.approve(address(vault), 100e6);
        vault.deposit(100e6, user1);
        vm.stopPrank();
        vm.prank(user1);
        vault.requestWithdraw(50e6);

        // 获取提款请求信息
        IPolyVault.WithdrawalRequest memory req = vault.getWithdrawalRequest(user1);
        // 验证份额为 50
        assertEq(req.shares, 50e6);
        // 验证处于待处理状态
        assertTrue(req.pending);
        // 验证时间戳已记录
        assertGt(req.requestTimestamp, 0);
    }

    /**
     * @notice 测试：初始策略债务为 0
     * @dev 验证默认值
     */
    function test_StrategyDebt_DefaultZero() public view {
        // 初始时策略债务为 0
        assertEq(vault.strategyDebt(), 0);
    }

    // ==================== UUPS 升级测试 ====================

    /**
     * @notice 测试：非管理员调用 upgradeToAndCall 应该 revert
     * @dev 验证升级权限控制
     */
    function test_RevertUpgrade_NotAdmin() public {
        // 非管理员尝试升级到随机地址
        vm.prank(user1);
        vm.expectRevert();
        vault.upgradeToAndCall(address(0xdead), "");
    }

    /**
     * @notice 测试：直接调用逻辑合约的 upgradeToAndCall 应该 revert
     * @dev 验证 UUPS 的代理上下文保护（逻辑合约不自毁）
     */
    function test_RevertUpgrade_FromImplementationDirectly() public {
        // 直接对逻辑合约（非代理）调用升级
        vm.prank(admin);
        vm.expectRevert();
        vaultImpl.upgradeToAndCall(address(0xdead), "");
    }

    /**
     * @notice 测试：管理员成功升级到新的逻辑合约
     * @dev 验证完整升级流程，检查 ERC1967 存储槽
     */
    function test_Upgrade_Success() public {
        // 部署一个新的逻辑合约
        PolyVault newImpl = new PolyVault();

        // 管理员执行升级，将实现地址更新为 newImpl
        vm.prank(admin);
        vault.upgradeToAndCall(address(newImpl), "");

        // 读取 ERC1967 实现地址存储槽（标准位置）
        bytes32 implSlot = 0x360894a13ba1a3210667c828492db98dca3e2076cc3735a920a3ca505d382bbc;
        address currentImpl = address(
            uint160(uint256(vm.load(address(vault), implSlot)))
        );
        // 验证存储槽中的实现地址已更新
        assertEq(currentImpl, address(newImpl));
    }

    /**
     * @notice 测试：代理和数据合约的存储是隔离的
     * @dev 验证逻辑合约的 storage 不被代理操作影响
     */
    function test_Proxy_ImplementationStorageUnchanged() public {
        // 通过代理执行存款操作
        vm.startPrank(user1);
        usdc.approve(address(vault), 100e6);
        vault.deposit(100e6, user1);
        vm.stopPrank();

        // 代理的 storage 有正确的值
        assertEq(vault.totalAssets(), 100e6);
        assertEq(vault.balanceOf(user1), 100e6);

        // 逻辑合约的 storage 不受代理操作影响（通过 vm.load 验证 ERC20 余额存储槽）
        // 逻辑合约未被初始化，因此 balanceOf 为 0
        assertEq(vaultImpl.balanceOf(user1), 0);
    }

    // ==================== 角色访问测试 ====================

    /**
     * @notice 测试：策略师不能暂停合约
     * @dev 验证 pause 需要 GUARDIAN_ROLE
     */
    function test_RevertStrategist_Pause() public {
        // 策略师尝试暂停
        vm.prank(strategist);
        vm.expectRevert();
        vault.pause();
    }

    /**
     * @notice 测试：守护者不能调用策略提款
     * @dev 验证 withdrawToStrategy 需要 STRATEGIST_ROLE
     */
    function test_RevertGuardian_WithdrawToStrategy() public {
        // 守护者尝试策略提款
        vm.prank(guardian);
        vm.expectRevert();
        vault.withdrawToStrategy(50e6);
    }

    /**
     * @notice 测试：普通用户不能调用任何管理员函数
     * @dev 批量验证所有管理员函数的权限控制
     */
    function test_RevertUser_AdminFunctions() public {
        // 以 user1 身份测试所有管理员函数
        vm.startPrank(user1);

        // 设置提款延迟应 revert
        vm.expectRevert();
        vault.setWithdrawalDelay(2 hours);

        // 设置业绩费应 revert
        vm.expectRevert();
        vault.setPerformanceFee(500);

        // 设置费用接收地址应 revert
        vm.expectRevert();
        vault.setFeeRecipient(address(0x99));

        // 设置最大策略分配应 revert
        vm.expectRevert();
        vault.setMaxStrategyAllocation(7500);

        // 设置存款限额应 revert
        vm.expectRevert();
        vault.setDepositLimits(1e6, 100e6);

        vm.stopPrank();
    }

    // ==================== 边界情况测试 ====================

    /**
     * @notice 测试：请求提款 → 取消 → 再次请求提款的完整流程
     * @dev 验证多次请求取消操作后状态正确
     */
    function test_RequestThenCancelThenRequestAgain() public {
        // 用户存款 100 USDC
        vm.startPrank(user1);
        usdc.approve(address(vault), 100e6);
        vault.deposit(100e6, user1);
        vm.stopPrank();

        // 第一次请求提款 30 份额
        vm.prank(user1);
        vault.requestWithdraw(30e6);

        // 取消提款
        vm.prank(user1);
        vault.cancelWithdraw();

        // 第二次请求提款 40 份额
        vm.prank(user1);
        vault.requestWithdraw(40e6);

        // 验证最终请求的份额为 40
        IPolyVault.WithdrawalRequest memory req = vault.getWithdrawalRequest(user1);
        assertEq(req.shares, 40e6);
    }

    /**
     * @notice 测试：多个用户同时有待处理的提款请求
     * @dev 验证多用户提款请求状态的独立性
     */
    function test_MultipleUserWithdrawRequests() public {
        // user1 存款 100 USDC，user2 存款 200 USDC
        vm.startPrank(user1);
        usdc.approve(address(vault), 100e6);
        vault.deposit(100e6, user1);
        vm.stopPrank();
        vm.startPrank(user2);
        usdc.approve(address(vault), 200e6);
        vault.deposit(200e6, user2);
        vm.stopPrank();

        // user1 请求提款 30 份额
        vm.prank(user1);
        vault.requestWithdraw(30e6);

        // user2 请求提款 50 份额
        vm.prank(user2);
        vault.requestWithdraw(50e6);

        // 验证 user1 的请求
        IPolyVault.WithdrawalRequest memory req1 = vault.getWithdrawalRequest(user1);
        assertEq(req1.shares, 30e6);
        assertTrue(req1.pending);

        // 验证 user2 的请求
        IPolyVault.WithdrawalRequest memory req2 = vault.getWithdrawalRequest(user2);
        assertEq(req2.shares, 50e6);
        assertTrue(req2.pending);
    }

    /**
     * @notice 测试：策略部分提款后，新用户存款不影响现有债务
     * @dev 验证债务和可用余额的独立性
     */
    function test_DepositAfterPartialStrategyWithdrawal() public {
        // user1 存款 100 USDC
        vm.startPrank(user1);
        usdc.approve(address(vault), 100e6);
        vault.deposit(100e6, user1);
        vm.stopPrank();

        // 策略师提取 30 USDC
        vm.prank(strategist);
        vault.withdrawToStrategy(30e6);

        // user2 存入 50 USDC
        vm.startPrank(user2);
        usdc.approve(address(vault), 50e6);
        vault.deposit(50e6, user2);
        vm.stopPrank();

        // 总资产 = 100 + 50 = 150 USDC
        assertEq(vault.totalAssets(), 150e6);
        // 策略债务保持 30 USDC
        assertEq(vault.strategyDebt(), 30e6);
        // 可用余额 = 70(剩余) + 50(新存) = 120 USDC
        assertEq(vault.availableBalance(), 120e6);
    }

    /**
     * @notice 测试：用户可以请求提款其全部余额
     * @dev 验证全额请求的边界情况
     */
    function test_RequestWithdraw_ExactBalance() public {
        // user1 存款 50 USDC
        vm.startPrank(user1);
        usdc.approve(address(vault), 50e6);
        vault.deposit(50e6, user1);
        vm.stopPrank();

        // 用户请求提款全部 50 份额
        vm.prank(user1);
        vault.requestWithdraw(50e6);

        // 验证用户余额为 0
        assertEq(vault.balanceOf(user1), 0);
        // 验证请求处于待处理状态
        assertTrue(vault.getWithdrawalRequest(user1).pending);
    }

    // ==================== 集成测试 ====================

    /**
     * @notice 集成测试：完整的存款 → 请求提款 → 取消 → 再次请求 → 执行提款周期
     * @dev 模拟用户的完整交互流程
     */
    function test_Integration_DepositWithdrawCycle() public {
        // 1. 用户存入 100 USDC
        vm.startPrank(user1);
        usdc.approve(address(vault), 100e6);
        vault.deposit(100e6, user1);
        vm.stopPrank();
        assertEq(vault.balanceOf(user1), 100e6);

        // 2. 用户请求提款 50 份额
        vm.prank(user1);
        vault.requestWithdraw(50e6);
        assertEq(vault.balanceOf(user1), 50e6);

        // 3. 用户取消提款，份额全部退回
        vm.prank(user1);
        vault.cancelWithdraw();
        assertEq(vault.balanceOf(user1), 100e6);

        // 4. 用户再次请求提款 50 份额
        vm.prank(user1);
        vault.requestWithdraw(50e6);

        // 5. 等待延迟期结束
        vm.warp(block.timestamp + DEFAULT_WITHDRAWAL_DELAY + 1);

        // 6. 用户执行提款
        vm.prank(user1);
        vault.executeWithdraw();

        // 验证用户获得 50 USDC（初始 1000 - 存款 100 + 提款 50 = 950）
        assertEq(usdc.balanceOf(user1), INITIAL_USER_BALANCE - 100e6 + 50e6);
        // 验证用户剩余 50 份额
        assertEq(vault.balanceOf(user1), 50e6);
    }

    /**
     * @notice 集成测试：完整的策略盈利周期
     * @dev 模拟用户存款 → 策略提款 → 盈利返还 → 验证费用的完整流程
     */
    function test_Integration_StrategyProfitCycle() public {
        // 1. 用户存入 1000 USDC
        vm.startPrank(user1);
        usdc.approve(address(vault), 1000e6);
        vault.deposit(1000e6, user1);
        vm.stopPrank();

        // 2. 策略师提取 500 USDC
        vm.prank(strategist);
        vault.withdrawToStrategy(500e6);
        assertEq(vault.strategyDebt(), 500e6);

        // 3. 策略盈利，返还 600 USDC（500 本金 + 100 利润）
        usdc.mint(strategist, 100e6);
        vm.prank(strategist);
        usdc.approve(address(vault), 600e6);
        vm.prank(strategist);
        vault.depositFromStrategy(600e6);

        // 4. 验证债务清零
        assertEq(vault.strategyDebt(), 0);
        // 验证业绩费 = 100 * 10% = 10 USDC
        assertEq(usdc.balanceOf(feeRecipient), 10e6);
        // 验证总资产 = 1000 + 100 - 10 = 1090 USDC
        assertEq(vault.totalAssets(), 1000e6 + 100e6 - 10e6);
    }

    /**
     * @notice 集成测试：完整的策略亏损周期
     * @dev 模拟用户存款 → 策略提款 → 亏损返还 → 验证无费用的完整流程
     */
    function test_Integration_StrategyLossCycle() public {
        // 1. 用户存入 1000 USDC
        vm.startPrank(user1);
        usdc.approve(address(vault), 1000e6);
        vault.deposit(1000e6, user1);
        vm.stopPrank();

        // 2. 策略师提取 500 USDC
        vm.prank(strategist);
        vault.withdrawToStrategy(500e6);

        // 3. 策略亏损，只返还 300 USDC
        vm.prank(strategist);
        usdc.approve(address(vault), 300e6);
        vm.prank(strategist);
        vault.depositFromStrategy(300e6);

        // 4. 验证剩余债务 200 USDC
        assertEq(vault.strategyDebt(), 200e6);
        // 验证无业绩费产生
        assertEq(usdc.balanceOf(feeRecipient), 0);
    }

    /**
     * @notice 集成测试：多轮策略操作（一盈一亏）
     * @dev 验证多次策略操作后状态一致性
     */
    function test_Integration_MultipleStrategyCycles() public {
        // 用户存入 1000 USDC
        vm.startPrank(user1);
        usdc.approve(address(vault), 1000e6);
        vault.deposit(1000e6, user1);
        vm.stopPrank();

        // ---- 第一轮：盈利 ----
        vm.prank(strategist);
        vault.withdrawToStrategy(300e6);

        // 铸造 50 USDC 利润并返还
        usdc.mint(strategist, 50e6);
        vm.prank(strategist);
        usdc.approve(address(vault), 350e6);
        vm.prank(strategist);
        vault.depositFromStrategy(350e6);

        // 验证第一轮业绩费 = 50 * 10% = 5 USDC
        uint256 feeAfterFirst = usdc.balanceOf(feeRecipient);
        assertEq(feeAfterFirst, 5e6);

        // ---- 第二轮：亏损 ----
        vm.prank(strategist);
        vault.withdrawToStrategy(200e6);

        // 只返还 100 USDC（亏损 100）
        vm.prank(strategist);
        usdc.approve(address(vault), 100e6);
        vm.prank(strategist);
        vault.depositFromStrategy(100e6);

        // 验证剩余债务 100 USDC
        assertEq(vault.strategyDebt(), 100e6);
        // 验证费用没有增加（亏损不收费）
        assertEq(usdc.balanceOf(feeRecipient), feeAfterFirst);
    }

    // ==================== 重复初始化测试 ====================

    /**
     * @notice 测试：已初始化的合约不能再次初始化
     * @dev 验证 Initializable 的防重复初始化保护
     */
    function test_RevertReinitialize() public {
        // 预期 revert
        vm.expectRevert();
        // 尝试再次初始化 vault（已在 setUp 中完成初始化）
        vault.initialize(
            address(usdc),
            admin,
            strategist,
            guardian,
            feeRecipient,
            DEFAULT_WITHDRAWAL_DELAY,
            DEFAULT_MAX_ALLOCATION,
            DEFAULT_PERFORMANCE_FEE
        );
    }

    // ==================== ERC1967 存储槽验证 ====================

    /**
     * @notice 测试：验证 ERC1967 实现地址存储槽的正确性
     * @dev 直接读取代理的 ERC1967 存储槽，确认实现地址正确设置
     */
    function test_Proxy_ERC1967StorageSlot() public view {
        // ERC1967 标准定义的实现地址存储槽位置
        // keccak256("eip1967.proxy.implementation") - 1
        bytes32 implSlot = 0x360894a13ba1a3210667c828492db98dca3e2076cc3735a920a3ca505d382bbc;
        // 使用 vm.load 直接读取代理合约的存储槽
        address storedImpl = address(uint160(uint256(vm.load(address(vault), implSlot))));
        // 验证存储的实现地址等于我们部署的逻辑合约
        assertEq(storedImpl, address(vaultImpl));
    }
}
