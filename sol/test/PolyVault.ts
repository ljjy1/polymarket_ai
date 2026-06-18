// 导入 Chai 断言库
import { expect } from "chai";
// Hardhat 运行时环境和升级插件
import hre from "hardhat";
import { upgrades } from "@openzeppelin/hardhat-upgrades";
// anyValue 用于事件参数匹配（忽略具体值）
import { anyValue } from "@nomicfoundation/hardhat-ethers-chai-matchers/withArgs";

// 通过 Hardhat network 创建 ethers.js 实例和升级 API
const connection = await hre.network.create();
const { ethers } = connection;
const upgradesApi = await upgrades(hre, connection);

// ======================== 测试常量定义 ========================

const INITIAL_USER_BALANCE = 1000_000000n; // 1000 USDC
const DEFAULT_WITHDRAWAL_DELAY = 3600n; // 1 hour
const DEFAULT_MAX_ALLOCATION = 5000n; // 50%
const DEFAULT_PERFORMANCE_FEE = 1000n; // 10%

function toUsdc(amount: bigint): bigint {
  return amount * 1_000_000n;
}

// ======================== 测试上下文接口 ========================

interface VaultContext {
  vault: any;
  usdc: any;
  vaultImpl: any;
  admin: any;
  strategist: any;
  guardian: any;
  feeRecipient: any;
  user1: any;
  user2: any;
}

// ======================== 基础部署函数 ========================

async function deployBase(): Promise<VaultContext> {
  const signers = await ethers.getSigners();
  const admin = signers[0];
  const strategist = signers[1];
  const guardian = signers[2];
  const feeRecipient = signers[3];
  const user1 = signers[4];
  const user2 = signers[5];

  // Deploy MockUSDC
  const MockUSDC = await ethers.getContractFactory("MockUSDC");
  const usdc = await MockUSDC.deploy();

  // Deploy PolyVault via UUPS proxy (handles implementation + proxy + initialize)
  const PolyVaultFactory = await ethers.getContractFactory("PolyVault");
  const vault = await upgradesApi.deployProxy(PolyVaultFactory, [
    await usdc.getAddress(),
    await admin.getAddress(),
    await strategist.getAddress(),
    await guardian.getAddress(),
    await feeRecipient.getAddress(),
    DEFAULT_WITHDRAWAL_DELAY,
    DEFAULT_MAX_ALLOCATION,
    DEFAULT_PERFORMANCE_FEE,
  ], { kind: "uups", initializer: "initialize" });

  // Get the implementation contract address for upgrade tests
  const vaultImplAddr = await upgradesApi.erc1967.getImplementationAddress(
    await vault.getAddress(),
  );
  const vaultImpl = await ethers.getContractAt("PolyVault", vaultImplAddr);

  // Mint USDC to users
  await usdc.mint(await user1.getAddress(), INITIAL_USER_BALANCE);
  await usdc.mint(await user2.getAddress(), INITIAL_USER_BALANCE);
  await usdc.mint(await strategist.getAddress(), 5000_000000n);

  return {
    vault, usdc, vaultImpl,
    admin, strategist, guardian, feeRecipient, user1, user2,
  };
}

// ============================================================
// 1. 初始化测试
// ============================================================

describe("PolyVault — 初始化", () => {
  let ctx: VaultContext;
  before(async () => { ctx = await deployBase(); });

  it("初始化名称、符号和资产代币正确", async () => {
    expect(await ctx.vault.name()).to.equal("PolyVault USDC");
    expect(await ctx.vault.symbol()).to.equal("pvUSDC");
    expect(await ctx.vault.asset()).to.equal(await ctx.usdc.getAddress());
  });

  it("初始化参数正确", async () => {
    expect(await ctx.vault.withdrawalDelay()).to.equal(DEFAULT_WITHDRAWAL_DELAY);
    expect(await ctx.vault.maxStrategyAllocation()).to.equal(DEFAULT_MAX_ALLOCATION);
    expect(await ctx.vault.performanceFee()).to.equal(DEFAULT_PERFORMANCE_FEE);
    expect(await ctx.vault.feeRecipient()).to.equal(await ctx.feeRecipient.getAddress());
    expect(await ctx.vault.minDeposit()).to.equal(toUsdc(1n));
    expect(await ctx.vault.maxDeposit()).to.equal(toUsdc(100_000n));
  });

  it("分配角色正确", async () => {
    expect(
      await ctx.vault.hasRole(await ctx.vault.DEFAULT_ADMIN_ROLE(), await ctx.admin.getAddress()),
    ).to.be.true;
    expect(
      await ctx.vault.hasRole(await ctx.vault.STRATEGIST_ROLE(), await ctx.strategist.getAddress()),
    ).to.be.true;
    expect(
      await ctx.vault.hasRole(await ctx.vault.GUARDIAN_ROLE(), await ctx.guardian.getAddress()),
    ).to.be.true;
  });

  it("重复初始化应回滚", async () => {
    await expect(
      ctx.vault.initialize(
        await ctx.usdc.getAddress(),
        await ctx.admin.getAddress(),
        await ctx.strategist.getAddress(),
        await ctx.guardian.getAddress(),
        await ctx.feeRecipient.getAddress(),
        DEFAULT_WITHDRAWAL_DELAY,
        DEFAULT_MAX_ALLOCATION,
        DEFAULT_PERFORMANCE_FEE,
      ),
    ).to.be.revertedWithCustomError(ctx.vault, "InvalidInitialization");
  });
});

// ============================================================
// 2. deposit
// ============================================================

describe("PolyVault — 存款", () => {
  let ctx: VaultContext;
  beforeEach(async () => { ctx = await deployBase(); });

  it("金额低于最小值应回滚", async () => {
    const amount = toUsdc(1n) - 1n;
    await expect(
      ctx.vault.connect(ctx.user1).deposit(amount, await ctx.user1.getAddress()),
    ).to.be.revertedWithCustomError(ctx.vault, "DepositBelowMinimum");
  });

  it("金额超过最大值应回滚", async () => {
    const amount = toUsdc(100_001n);
    await expect(
      ctx.vault.connect(ctx.user1).deposit(amount, await ctx.user1.getAddress()),
    ).to.be.revertedWithCustomError(ctx.vault, "DepositAboveMaximum");
  });

  it("暂停状态下应回滚", async () => {
    await ctx.vault.connect(ctx.guardian).pause();
    await expect(
      ctx.vault
        .connect(ctx.user1)
        .deposit(toUsdc(100n), await ctx.user1.getAddress()),
    ).to.be.revertedWithCustomError(ctx.vault, "EnforcedPause");
  });

  it("以 1:1 比例成功存款", async () => {
    const amount = toUsdc(100n);
    await ctx.usdc
      .connect(ctx.user1)
      .approve(await ctx.vault.getAddress(), amount);
    const shares = await ctx.vault
      .connect(ctx.user1)
      .deposit.staticCall(amount, await ctx.user1.getAddress());
    await ctx.vault.connect(ctx.user1).deposit(amount, await ctx.user1.getAddress());

    expect(shares).to.equal(amount);
    expect(await ctx.vault.balanceOf(await ctx.user1.getAddress())).to.equal(amount);
    expect(await ctx.usdc.balanceOf(await ctx.vault.getAddress())).to.equal(amount);
  });

  it("成功存入最小金额", async () => {
    const amount = toUsdc(1n);
    await ctx.usdc
      .connect(ctx.user1)
      .approve(await ctx.vault.getAddress(), amount);
    const shares = await ctx.vault
      .connect(ctx.user1)
      .deposit.staticCall(amount, await ctx.user1.getAddress());
    await ctx.vault.connect(ctx.user1).deposit(amount, await ctx.user1.getAddress());
    expect(shares).to.equal(amount);
  });

  it("成功存入最大金额", async () => {
    const amount = toUsdc(100_000n);
    await ctx.usdc.mint(await ctx.user1.getAddress(), amount);
    await ctx.usdc
      .connect(ctx.user1)
      .approve(await ctx.vault.getAddress(), amount);
    const shares = await ctx.vault
      .connect(ctx.user1)
      .deposit.staticCall(amount, await ctx.user1.getAddress());
    await ctx.vault.connect(ctx.user1).deposit(amount, await ctx.user1.getAddress());
    expect(shares).to.equal(amount);
  });

  it("多个用户同时存款", async () => {
    const amount1 = toUsdc(50n);
    const amount2 = toUsdc(75n);

    await ctx.usdc
      .connect(ctx.user1)
      .approve(await ctx.vault.getAddress(), amount1);
    await ctx.vault.connect(ctx.user1).deposit(amount1, await ctx.user1.getAddress());

    await ctx.usdc
      .connect(ctx.user2)
      .approve(await ctx.vault.getAddress(), amount2);
    await ctx.vault.connect(ctx.user2).deposit(amount2, await ctx.user2.getAddress());

    expect(await ctx.vault.balanceOf(await ctx.user1.getAddress())).to.equal(amount1);
    expect(await ctx.vault.balanceOf(await ctx.user2.getAddress())).to.equal(amount2);
  });
});

// ============================================================
// 3. mint
// ============================================================

describe("PolyVault — 铸造", () => {
  let ctx: VaultContext;
  beforeEach(async () => { ctx = await deployBase(); });

  it("所需资产低于最小值应回滚", async () => {
    await expect(
      ctx.vault.connect(ctx.user1).mint(1, await ctx.user1.getAddress()),
    ).to.be.revertedWithCustomError(ctx.vault, "DepositBelowMinimum");
  });

  it("暂停状态下应回滚", async () => {
    await ctx.vault.connect(ctx.guardian).pause();
    await expect(
      ctx.vault.connect(ctx.user1).mint(toUsdc(100n), await ctx.user1.getAddress()),
    ).to.be.revertedWithCustomError(ctx.vault, "EnforcedPause");
  });

  it("成功铸造精确份额", async () => {
    const shares = toUsdc(100n);
    const requiredAssets = await ctx.vault.previewMint(shares);

    await ctx.usdc.mint(await ctx.user1.getAddress(), requiredAssets);
    await ctx.usdc
      .connect(ctx.user1)
      .approve(await ctx.vault.getAddress(), requiredAssets);
    await ctx.vault.connect(ctx.user1).mint(shares, await ctx.user1.getAddress());

    expect(await ctx.vault.balanceOf(await ctx.user1.getAddress())).to.equal(shares);
  });
});

// ============================================================
// 4. direct withdraw disabled
// ============================================================

describe("PolyVault — 直接提款禁用", () => {
  let ctx: VaultContext;
  beforeEach(async () => {
    ctx = await deployBase();
    await ctx.usdc
      .connect(ctx.user1)
      .approve(await ctx.vault.getAddress(), toUsdc(100n));
    await ctx.vault.connect(ctx.user1).deposit(toUsdc(100n), await ctx.user1.getAddress());
  });

  it("直接提款应回滚", async () => {
    const addr = await ctx.user1.getAddress();
    await expect(
      ctx.vault.connect(ctx.user1).withdraw(toUsdc(100n), addr, addr),
    ).to.be.revertedWithCustomError(ctx.vault, "DirectWithdrawDisabled");
  });

  it("直接赎回应回滚", async () => {
    const addr = await ctx.user1.getAddress();
    await expect(
      ctx.vault.connect(ctx.user1).redeem(toUsdc(100n), addr, addr),
    ).to.be.revertedWithCustomError(ctx.vault, "DirectWithdrawDisabled");
  });

  it("maxWithdraw 返回 0", async () => {
    expect(await ctx.vault.maxWithdraw(await ctx.user1.getAddress())).to.equal(0);
  });

  it("maxRedeem 返回 0", async () => {
    expect(await ctx.vault.maxRedeem(await ctx.user1.getAddress())).to.equal(0);
  });
});

// ============================================================
// 5. requestWithdraw
// ============================================================

describe("PolyVault — 请求提款", () => {
  let ctx: VaultContext;
  beforeEach(async () => {
    ctx = await deployBase();
    await ctx.usdc
      .connect(ctx.user1)
      .approve(await ctx.vault.getAddress(), toUsdc(100n));
    await ctx.vault.connect(ctx.user1).deposit(toUsdc(100n), await ctx.user1.getAddress());
  });

  it("数量为零应回滚", async () => {
    await expect(
      ctx.vault.connect(ctx.user1).requestWithdraw(0),
    ).to.be.revertedWithCustomError(ctx.vault, "ZeroAmount");
  });

  it("有待处理请求时再次请求应回滚", async () => {
    await ctx.vault.connect(ctx.user1).requestWithdraw(toUsdc(50n));
    await expect(
      ctx.vault.connect(ctx.user1).requestWithdraw(toUsdc(30n)),
    ).to.be.revertedWithCustomError(ctx.vault, "WithdrawalAlreadyPending");
  });

  it("余额不足应回滚", async () => {
    await expect(
      ctx.vault.connect(ctx.user1).requestWithdraw(toUsdc(200n)),
    ).to.be.revertedWithCustomError(ctx.vault, "InsufficientShares");
  });

  it("暂停状态下应回滚", async () => {
    await ctx.vault.connect(ctx.guardian).pause();
    await expect(
      ctx.vault.connect(ctx.user1).requestWithdraw(toUsdc(50n)),
    ).to.be.revertedWithCustomError(ctx.vault, "EnforcedPause");
  });

  it("锁定份额并记录请求", async () => {
    const userAddr = await ctx.user1.getAddress();
    const requestAmount = toUsdc(50n);

    const userBalanceBefore = await ctx.vault.balanceOf(userAddr);
    const vaultBalanceBefore = await ctx.vault.balanceOf(await ctx.vault.getAddress());

    await ctx.vault.connect(ctx.user1).requestWithdraw(requestAmount);

    const req = await ctx.vault.getWithdrawalRequest(userAddr);
    expect(req.shares).to.equal(requestAmount);
    expect(req.pending).to.be.true;
    expect(req.requestTimestamp).to.be.gt(0);

    expect(await ctx.vault.balanceOf(userAddr)).to.equal(userBalanceBefore - requestAmount);
    expect(await ctx.vault.balanceOf(await ctx.vault.getAddress())).to.equal(vaultBalanceBefore + requestAmount);
  });

  it("触发 WithdrawalRequested 事件", async () => {
    const userAddr = await ctx.user1.getAddress();

    await expect(
      ctx.vault.connect(ctx.user1).requestWithdraw(toUsdc(50n)),
    )
      .to.emit(ctx.vault, "WithdrawalRequested")
      .withArgs(userAddr, toUsdc(50n), anyValue);
  });
});

// ============================================================
// 6. cancelWithdraw
// ============================================================

describe("PolyVault — 取消提款", () => {
  let ctx: VaultContext;
  beforeEach(async () => {
    ctx = await deployBase();
    await ctx.usdc
      .connect(ctx.user1)
      .approve(await ctx.vault.getAddress(), toUsdc(100n));
    await ctx.vault.connect(ctx.user1).deposit(toUsdc(100n), await ctx.user1.getAddress());
    await ctx.vault.connect(ctx.user1).requestWithdraw(toUsdc(50n));
  });

  it("无待处理请求应回滚", async () => {
    await expect(
      ctx.vault.connect(ctx.user2).cancelWithdraw(),
    ).to.be.revertedWithCustomError(ctx.vault, "NoPendingWithdrawal");
  });

  it("返还份额并清除请求", async () => {
    const userAddr = await ctx.user1.getAddress();
    const balanceBefore = await ctx.vault.balanceOf(userAddr);

    await ctx.vault.connect(ctx.user1).cancelWithdraw();

    expect(await ctx.vault.balanceOf(userAddr)).to.equal(balanceBefore + toUsdc(50n));

    const req = await ctx.vault.getWithdrawalRequest(userAddr);
    expect(req.pending).to.be.false;
    expect(req.shares).to.equal(0);
  });

  it("触发 WithdrawalCancelled 事件", async () => {
    const userAddr = await ctx.user1.getAddress();

    await expect(ctx.vault.connect(ctx.user1).cancelWithdraw())
      .to.emit(ctx.vault, "WithdrawalCancelled")
      .withArgs(userAddr, toUsdc(50n));
  });
});

// ============================================================
// 7. executeWithdraw
// ============================================================

describe("PolyVault — 执行提款", () => {
  let ctx: VaultContext;
  beforeEach(async () => {
    ctx = await deployBase();
    await ctx.usdc
      .connect(ctx.user1)
      .approve(await ctx.vault.getAddress(), toUsdc(100n));
    await ctx.vault.connect(ctx.user1).deposit(toUsdc(100n), await ctx.user1.getAddress());
    await ctx.vault.connect(ctx.user1).requestWithdraw(toUsdc(50n));
  });

  it("无待处理请求应回滚", async () => {
    await ctx.vault.connect(ctx.user1).cancelWithdraw();
    await expect(
      ctx.vault.connect(ctx.user1).executeWithdraw(),
    ).to.be.revertedWithCustomError(ctx.vault, "NoPendingWithdrawal");
  });

  it("should revert when delay not met", async () => {
    await expect(
      ctx.vault.connect(ctx.user1).executeWithdraw(),
    ).to.be.revertedWithCustomError(ctx.vault, "WithdrawalDelayNotMet");
  });

  it("延迟后执行提款", async () => {
    const userAddr = await ctx.user1.getAddress();
    const usdcBefore = await ctx.usdc.balanceOf(userAddr);

    await ethers.provider.send("evm_increaseTime", [Number(DEFAULT_WITHDRAWAL_DELAY) + 1]);
    await ethers.provider.send("evm_mine");
    await ctx.vault.connect(ctx.user1).executeWithdraw();

    expect(await ctx.usdc.balanceOf(userAddr) - usdcBefore).to.equal(toUsdc(50n));

    const req = await ctx.vault.getWithdrawalRequest(userAddr);
    expect(req.pending).to.be.false;
    expect(await ctx.vault.balanceOf(userAddr)).to.equal(toUsdc(50n));
  });

  it("触发 WithdrawalExecuted 事件", async () => {
    const userAddr = await ctx.user1.getAddress();

    await ethers.provider.send("evm_increaseTime", [Number(DEFAULT_WITHDRAWAL_DELAY) + 1]);
    await ethers.provider.send("evm_mine");

    await expect(ctx.vault.connect(ctx.user1).executeWithdraw())
      .to.emit(ctx.vault, "WithdrawalExecuted")
      .withArgs(userAddr, toUsdc(50n), toUsdc(50n));
  });

  it("金库余额不足时执行部分提款", async () => {
    await ctx.vault
      .connect(ctx.admin)
      .setMaxStrategyAllocation(10000n);
    await ctx.vault.connect(ctx.strategist).withdrawToStrategy(toUsdc(80n));

    const userAddr = await ctx.user1.getAddress();
    const usdcBefore = await ctx.usdc.balanceOf(userAddr);

    await ethers.provider.send("evm_increaseTime", [Number(DEFAULT_WITHDRAWAL_DELAY) + 1]);
    await ethers.provider.send("evm_mine");
    await ctx.vault.connect(ctx.user1).executeWithdraw();

    expect(await ctx.usdc.balanceOf(userAddr) - usdcBefore).to.equal(toUsdc(20n));
    expect(await ctx.vault.balanceOf(userAddr)).to.be.gt(0);
  });

  it("在延迟临界时间点执行提款", async () => {
    const userAddr = await ctx.user1.getAddress();
    const usdcBefore = await ctx.usdc.balanceOf(userAddr);

    await ethers.provider.send("evm_increaseTime", [Number(DEFAULT_WITHDRAWAL_DELAY)]);
    await ethers.provider.send("evm_mine");
    await ctx.vault.connect(ctx.user1).executeWithdraw();

    expect(await ctx.usdc.balanceOf(userAddr) - usdcBefore).to.equal(toUsdc(50n));
  });
});

// ============================================================
// 8. withdrawToStrategy
// ============================================================

describe("PolyVault — 策略提款", () => {
  let ctx: VaultContext;
  beforeEach(async () => {
    ctx = await deployBase();
    await ctx.usdc
      .connect(ctx.user1)
      .approve(await ctx.vault.getAddress(), toUsdc(100n));
    await ctx.vault.connect(ctx.user1).deposit(toUsdc(100n), await ctx.user1.getAddress());
  });

  it("数量为零应回滚", async () => {
    await expect(
      ctx.vault.connect(ctx.strategist).withdrawToStrategy(0),
    ).to.be.revertedWithCustomError(ctx.vault, "ZeroAmount");
  });

  it("非策略师应回滚", async () => {
    await expect(
      ctx.vault.connect(ctx.user1).withdrawToStrategy(toUsdc(50n)),
    ).to.be.revertedWithCustomError(ctx.vault, "AccessControlUnauthorizedAccount");
  });

  it("超出分配额度应回滚", async () => {
    await ctx.vault.connect(ctx.strategist).withdrawToStrategy(toUsdc(40n));
    await expect(
      ctx.vault.connect(ctx.strategist).withdrawToStrategy(toUsdc(20n)),
    ).to.be.revertedWithCustomError(ctx.vault, "StrategyAllocationExceeded");
  });

  it("转账 USDC 并记录债务", async () => {
    const strategistAddr = await ctx.strategist.getAddress();
    const balanceBefore = await ctx.usdc.balanceOf(strategistAddr);

    await ctx.vault.connect(ctx.strategist).withdrawToStrategy(toUsdc(50n));

    expect(await ctx.usdc.balanceOf(strategistAddr) - balanceBefore).to.equal(toUsdc(50n));
    expect(await ctx.vault.strategyDebt()).to.equal(toUsdc(50n));
    expect(await ctx.vault.availableBalance()).to.equal(toUsdc(50n));
  });

  it("触发 StrategyWithdrawal 事件", async () => {
    const strategistAddr = await ctx.strategist.getAddress();

    await expect(
      ctx.vault.connect(ctx.strategist).withdrawToStrategy(toUsdc(30n)),
    )
      .to.emit(ctx.vault, "StrategyWithdrawal")
      .withArgs(strategistAddr, toUsdc(30n));
  });

  it("允许提取到最大分配额度", async () => {
    await ctx.vault.connect(ctx.strategist).withdrawToStrategy(toUsdc(50n));
    expect(await ctx.vault.strategyDebt()).to.equal(toUsdc(50n));
  });
});

// ============================================================
// 9. depositFromStrategy
// ============================================================

describe("PolyVault — 策略存款", () => {
  let ctx: VaultContext;
  beforeEach(async () => {
    ctx = await deployBase();
    await ctx.usdc
      .connect(ctx.user1)
      .approve(await ctx.vault.getAddress(), toUsdc(100n));
    await ctx.vault.connect(ctx.user1).deposit(toUsdc(100n), await ctx.user1.getAddress());
    await ctx.vault.connect(ctx.strategist).withdrawToStrategy(toUsdc(50n));
  });

  it("数量为零应回滚", async () => {
    await expect(
      ctx.vault.connect(ctx.strategist).depositFromStrategy(0),
    ).to.be.revertedWithCustomError(ctx.vault, "ZeroAmount");
  });

  it("非策略师应回滚", async () => {
    await expect(
      ctx.vault.connect(ctx.user1).depositFromStrategy(toUsdc(50n)),
    ).to.be.revertedWithCustomError(ctx.vault, "AccessControlUnauthorizedAccount");
  });

  it("未授权应回滚", async () => {
    await expect(
      ctx.vault.connect(ctx.strategist).depositFromStrategy(toUsdc(50n)),
    ).to.be.revertedWithCustomError(ctx.usdc, "ERC20InsufficientAllowance");
  });

  it("无利润时返还本金", async () => {
    const strategistAddr = await ctx.strategist.getAddress();
    const feeBefore = await ctx.usdc.balanceOf(await ctx.feeRecipient.getAddress());

    await ctx.usdc
      .connect(ctx.strategist)
      .approve(await ctx.vault.getAddress(), toUsdc(50n));
    await ctx.vault.connect(ctx.strategist).depositFromStrategy(toUsdc(50n));

    expect(await ctx.vault.strategyDebt()).to.equal(0);
    expect(await ctx.vault.availableBalance()).to.equal(toUsdc(100n));
    expect(await ctx.usdc.balanceOf(await ctx.feeRecipient.getAddress())).to.equal(feeBefore);
  });

  it("分配利润费用", async () => {
    const profit = toUsdc(10n);
    const returnAmount = toUsdc(50n) + profit;

    await ctx.usdc.mint(await ctx.strategist.getAddress(), profit);
    await ctx.usdc
      .connect(ctx.strategist)
      .approve(await ctx.vault.getAddress(), returnAmount);
    await ctx.vault.connect(ctx.strategist).depositFromStrategy(returnAmount);

    expect(await ctx.vault.strategyDebt()).to.equal(0);
    expect(await ctx.usdc.balanceOf(await ctx.feeRecipient.getAddress())).to.equal(toUsdc(1n));
  });

  it("有利润时触发 ProfitReported 事件", async () => {
    const profit = toUsdc(10n);
    const returnAmount = toUsdc(50n) + profit;

    await ctx.usdc.mint(await ctx.strategist.getAddress(), profit);
    await ctx.usdc
      .connect(ctx.strategist)
      .approve(await ctx.vault.getAddress(), returnAmount);

    await expect(
      ctx.vault.connect(ctx.strategist).depositFromStrategy(returnAmount),
    )
      .to.emit(ctx.vault, "ProfitReported")
      .withArgs(profit, toUsdc(1n));
  });

  it("触发 StrategyDeposit 事件", async () => {
    const strategistAddr = await ctx.strategist.getAddress();

    await ctx.usdc
      .connect(ctx.strategist)
      .approve(await ctx.vault.getAddress(), toUsdc(50n));

    await expect(
      ctx.vault.connect(ctx.strategist).depositFromStrategy(toUsdc(50n)),
    )
      .to.emit(ctx.vault, "StrategyDeposit")
      .withArgs(strategistAddr, toUsdc(50n));
  });

  it("亏损时减少债务", async () => {
    await ctx.usdc
      .connect(ctx.strategist)
      .approve(await ctx.vault.getAddress(), toUsdc(30n));
    await ctx.vault.connect(ctx.strategist).depositFromStrategy(toUsdc(30n));

    expect(await ctx.vault.strategyDebt()).to.equal(toUsdc(20n));
    expect(await ctx.usdc.balanceOf(await ctx.feeRecipient.getAddress())).to.equal(0);
  });

  it("业绩费为 0 时不收费", async () => {
    await ctx.vault
      .connect(ctx.admin)
      .setPerformanceFee(0);

    const profit = toUsdc(20n);
    await ctx.usdc.mint(await ctx.strategist.getAddress(), profit);
    await ctx.usdc
      .connect(ctx.strategist)
      .approve(await ctx.vault.getAddress(), toUsdc(70n));
    await ctx.vault.connect(ctx.strategist).depositFromStrategy(toUsdc(70n));

    expect(await ctx.usdc.balanceOf(await ctx.feeRecipient.getAddress())).to.equal(0);
  });

  it("多个策略周期", async () => {
    await ctx.usdc
      .connect(ctx.strategist)
      .approve(await ctx.vault.getAddress(), toUsdc(50n));
    await ctx.vault.connect(ctx.strategist).depositFromStrategy(toUsdc(50n));
    expect(await ctx.vault.strategyDebt()).to.equal(0);

    await ctx.vault.connect(ctx.strategist).withdrawToStrategy(toUsdc(30n));
    await ctx.usdc
      .connect(ctx.strategist)
      .approve(await ctx.vault.getAddress(), toUsdc(30n));
    await ctx.vault.connect(ctx.strategist).depositFromStrategy(toUsdc(30n));
    expect(await ctx.vault.strategyDebt()).to.equal(0);

    await ctx.vault.connect(ctx.strategist).withdrawToStrategy(toUsdc(20n));
    await ctx.usdc
      .connect(ctx.strategist)
      .approve(await ctx.vault.getAddress(), toUsdc(20n));
    await ctx.vault.connect(ctx.strategist).depositFromStrategy(toUsdc(20n));
    expect(await ctx.vault.strategyDebt()).to.equal(0);
  });
});

// ============================================================
// 10. totalAssets
// ============================================================

describe("PolyVault — 总资产", () => {
  let ctx: VaultContext;
  beforeEach(async () => { ctx = await deployBase(); });

  it("无策略债务时等于金库余额", async () => {
    await ctx.usdc
      .connect(ctx.user1)
      .approve(await ctx.vault.getAddress(), toUsdc(100n));
    await ctx.vault.connect(ctx.user1).deposit(toUsdc(100n), await ctx.user1.getAddress());

    expect(await ctx.vault.totalAssets()).to.equal(toUsdc(100n));
  });

  it("包含策略债务", async () => {
    await ctx.usdc
      .connect(ctx.user1)
      .approve(await ctx.vault.getAddress(), toUsdc(100n));
    await ctx.vault.connect(ctx.user1).deposit(toUsdc(100n), await ctx.user1.getAddress());
    await ctx.vault.connect(ctx.strategist).withdrawToStrategy(toUsdc(30n));

    expect(await ctx.vault.totalAssets()).to.equal(toUsdc(100n));
  });

  it("策略盈利后总资产增加", async () => {
    await ctx.usdc
      .connect(ctx.user1)
      .approve(await ctx.vault.getAddress(), toUsdc(100n));
    await ctx.vault.connect(ctx.user1).deposit(toUsdc(100n), await ctx.user1.getAddress());
    await ctx.vault.connect(ctx.strategist).withdrawToStrategy(toUsdc(50n));

    await ctx.usdc.mint(await ctx.strategist.getAddress(), toUsdc(10n));
    await ctx.usdc
      .connect(ctx.strategist)
      .approve(await ctx.vault.getAddress(), toUsdc(60n));
    await ctx.vault.connect(ctx.strategist).depositFromStrategy(toUsdc(60n));

    expect(await ctx.vault.totalAssets()).to.equal(toUsdc(109n));
  });
});

// ============================================================
// 11. admin functions
// ============================================================

describe("PolyVault — 管理员功能", () => {
  let ctx: VaultContext;
  beforeEach(async () => { ctx = await deployBase(); });

  it("更新提款延迟时间", async () => {
    const newDelay = 7200n;

    await expect(
      ctx.vault.connect(ctx.admin).setWithdrawalDelay(newDelay),
    )
      .to.emit(ctx.vault, "WithdrawalDelayUpdated")
      .withArgs(DEFAULT_WITHDRAWAL_DELAY, newDelay);

    expect(await ctx.vault.withdrawalDelay()).to.equal(newDelay);
  });

  it("延迟时间过短应回滚", async () => {
    await expect(
      ctx.vault.connect(ctx.admin).setWithdrawalDelay(1800n),
    ).to.be.revertedWithCustomError(ctx.vault, "InvalidWithdrawalDelay");
  });

  it("延迟时间过长应回滚", async () => {
    await expect(
      ctx.vault.connect(ctx.admin).setWithdrawalDelay(8n * 86400n),
    ).to.be.revertedWithCustomError(ctx.vault, "InvalidWithdrawalDelay");
  });

  it("非管理员设置延迟应回滚", async () => {
    await expect(
      ctx.vault.connect(ctx.user1).setWithdrawalDelay(7200n),
    ).to.be.revertedWithCustomError(ctx.vault, "AccessControlUnauthorizedAccount");
  });

  it("更新最大策略分配比例", async () => {
    const newAllocation = 7500n;

    await expect(
      ctx.vault.connect(ctx.admin).setMaxStrategyAllocation(newAllocation),
    )
      .to.emit(ctx.vault, "MaxStrategyAllocationUpdated")
      .withArgs(DEFAULT_MAX_ALLOCATION, newAllocation);

    expect(await ctx.vault.maxStrategyAllocation()).to.equal(newAllocation);
  });

  it("分配比例超过 100% 应回滚", async () => {
    await expect(
      ctx.vault.connect(ctx.admin).setMaxStrategyAllocation(10001n),
    ).to.be.revertedWithCustomError(ctx.vault, "InvalidAllocation");
  });

  it("更新业绩费", async () => {
    const newFee = 500n;

    await expect(ctx.vault.connect(ctx.admin).setPerformanceFee(newFee))
      .to.emit(ctx.vault, "PerformanceFeeUpdated")
      .withArgs(DEFAULT_PERFORMANCE_FEE, newFee);

    expect(await ctx.vault.performanceFee()).to.equal(newFee);
  });

  it("业绩费超过最大值应回滚", async () => {
    await expect(
      ctx.vault.connect(ctx.admin).setPerformanceFee(2001n),
    ).to.be.revertedWithCustomError(ctx.vault, "InvalidPerformanceFee");
  });

  it("更新费用接收地址", async () => {
    const newRecipient = await ctx.user2.getAddress();

    await expect(ctx.vault.connect(ctx.admin).setFeeRecipient(newRecipient))
      .to.emit(ctx.vault, "FeeRecipientUpdated")
      .withArgs(await ctx.feeRecipient.getAddress(), newRecipient);

    expect(await ctx.vault.feeRecipient()).to.equal(newRecipient);
  });

  it("费用接收地址为零应回滚", async () => {
    await expect(
      ctx.vault.connect(ctx.admin).setFeeRecipient(ethers.ZeroAddress),
    ).to.be.revertedWithCustomError(ctx.vault, "ZeroAddress");
  });

  it("更新存款限额", async () => {
    const newMin = toUsdc(10n);
    const newMax = toUsdc(50_000n);

    await expect(ctx.vault.connect(ctx.admin).setDepositLimits(newMin, newMax))
      .to.emit(ctx.vault, "DepositLimitsUpdated")
      .withArgs(newMin, newMax);

    expect(await ctx.vault.minDeposit()).to.equal(newMin);
    expect(await ctx.vault.maxDeposit()).to.equal(newMax);
  });

  it("更新限额后成功存款", async () => {
    await ctx.vault.connect(ctx.admin).setDepositLimits(toUsdc(1n), toUsdc(200_000n));

    await ctx.usdc.mint(await ctx.user1.getAddress(), toUsdc(200_000n));
    await ctx.usdc
      .connect(ctx.user1)
      .approve(await ctx.vault.getAddress(), toUsdc(200_000n));
    const shares = await ctx.vault
      .connect(ctx.user1)
      .deposit.staticCall(toUsdc(200_000n), await ctx.user1.getAddress());
    await ctx.vault
      .connect(ctx.user1)
      .deposit(toUsdc(200_000n), await ctx.user1.getAddress());

    expect(shares).to.equal(toUsdc(200_000n));
  });
});

// ============================================================
// 12. pause/unpause
// ============================================================

describe("PolyVault — 暂停/恢复", () => {
  let ctx: VaultContext;
  beforeEach(async () => { ctx = await deployBase(); });

  it("守护者暂停合约", async () => {
    await ctx.vault.connect(ctx.guardian).pause();

    await ctx.usdc
      .connect(ctx.user1)
      .approve(await ctx.vault.getAddress(), toUsdc(100n));
    await expect(
      ctx.vault.connect(ctx.user1).deposit(toUsdc(100n), await ctx.user1.getAddress()),
    ).to.be.revertedWithCustomError(ctx.vault, "EnforcedPause");
  });

  it("非守护者暂停应回滚", async () => {
    await expect(ctx.vault.connect(ctx.user1).pause()).to.be.revertedWithCustomError(ctx.vault, "AccessControlUnauthorizedAccount");
  });

  it("守护者恢复并允许存款", async () => {
    await ctx.vault.connect(ctx.guardian).pause();
    await ctx.vault.connect(ctx.guardian).unpause();

    await ctx.usdc
      .connect(ctx.user1)
      .approve(await ctx.vault.getAddress(), toUsdc(100n));
    const shares = await ctx.vault
      .connect(ctx.user1)
      .deposit.staticCall(toUsdc(100n), await ctx.user1.getAddress());
    await ctx.vault.connect(ctx.user1).deposit(toUsdc(100n), await ctx.user1.getAddress());

    expect(shares).to.equal(toUsdc(100n));
  });

  it("非守护者恢复应回滚", async () => {
    await ctx.vault.connect(ctx.guardian).pause();
    await expect(ctx.vault.connect(ctx.user1).unpause()).to.be.revertedWithCustomError(ctx.vault, "AccessControlUnauthorizedAccount");
  });
});

// ============================================================
// 13. view functions
// ============================================================

describe("PolyVault — 视图函数", () => {
  let ctx: VaultContext;
  before(async () => { ctx = await deployBase(); });

  it("默认返回空提款请求", async () => {
    const req = await ctx.vault.getWithdrawalRequest(await ctx.user1.getAddress());
    expect(req.pending).to.be.false;
    expect(req.shares).to.equal(0);
  });

  it("请求后返回正确的提款请求", async () => {
    // 重新部署以获取干净状态（这个测试会修改状态）
    const localCtx = await deployBase();
    await localCtx.usdc
      .connect(localCtx.user1)
      .approve(await localCtx.vault.getAddress(), toUsdc(100n));
    await localCtx.vault.connect(localCtx.user1).deposit(toUsdc(100n), await localCtx.user1.getAddress());
    await localCtx.vault.connect(localCtx.user1).requestWithdraw(toUsdc(50n));

    const req = await localCtx.vault.getWithdrawalRequest(await localCtx.user1.getAddress());
    expect(req.shares).to.equal(toUsdc(50n));
    expect(req.pending).to.be.true;
    expect(req.requestTimestamp).to.be.gt(0);
  });

  it("初始策略债务为零", async () => {
    expect(await ctx.vault.strategyDebt()).to.equal(0);
  });
});

// ============================================================
// 14. UUPS upgrade
// ============================================================

describe("PolyVault — UUPS 升级", () => {
  let ctx: VaultContext;
  beforeEach(async () => { ctx = await deployBase(); });

  it("非管理员升级应回滚", async () => {
    await expect(
      ctx.vault
        .connect(ctx.user1)
        .upgradeToAndCall(
          "0x" + "f".repeat(40),
          "0x",
        ),
    ).to.be.revertedWithCustomError(ctx.vault, "AccessControlUnauthorizedAccount");
  });

  it("管理员成功升级", async () => {
    const PolyVaultFactory = await ethers.getContractFactory("PolyVault", ctx.admin);
    await upgradesApi.upgradeProxy(
      await ctx.vault.getAddress(),
      PolyVaultFactory,
      { kind: "uups" },
    );

    const newImplAddr = await upgradesApi.erc1967.getImplementationAddress(
      await ctx.vault.getAddress(),
    );
    expect(newImplAddr).to.not.equal(ethers.ZeroAddress);

    // Verify ERC1967 storage slot matches
    const implSlot = "0x360894a13ba1a3210667c828492db98dca3e2076cc3735a920a3ca505d382bbc";
    const currentImplStorage = await ethers.provider.send("eth_getStorageAt", [
      await ctx.vault.getAddress(), implSlot, "latest",
    ]);
    const actualImpl = "0x" + currentImplStorage.slice(-40);
    expect(actualImpl.toLowerCase()).to.equal(newImplAddr.toLowerCase());
  });

  it("升级后状态和功能保持正常", async () => {
    // NOTE: 此测试使用同一个 PolyVault 作为新实现，仅用于验证升级流程
    // 正常场景下，应部署一个继承自 PolyVault 的升级合约（如 PolyVaultV2），
    // 重写或新增函数，例如：
    //
    //   contract PolyVaultV2 is PolyVault {
    //       function newFeature() external onlyRole(DEFAULT_ADMIN_ROLE) returns (uint256) {
    //           return 42;
    //       }
    //   }
    //
    // 先将用户1存入一些 USDC，确保升级后状态保持
    const depositAmount = toUsdc(200n);
    await ctx.usdc.connect(ctx.user1).approve(await ctx.vault.getAddress(), depositAmount);
    await ctx.vault.connect(ctx.user1).deposit(depositAmount, await ctx.user1.getAddress());
    const balanceBefore = await ctx.vault.balanceOf(await ctx.user1.getAddress());

    // 部署并升级到新实现（此处仅为测试目的复用 PolyVault）
    const PolyVaultFactory = await ethers.getContractFactory("PolyVault", ctx.admin);
    await upgradesApi.upgradeProxy(
      await ctx.vault.getAddress(),
      PolyVaultFactory,
      { kind: "uups" },
    );

    // 验证升级后 storage slot 已更新
    const newImplAddr = await upgradesApi.erc1967.getImplementationAddress(
      await ctx.vault.getAddress(),
    );
    const implSlot = "0x360894a13ba1a3210667c828492db98dca3e2076cc3735a920a3ca505d382bbc";
    const currentImplStorage = await ethers.provider.send("eth_getStorageAt", [
      await ctx.vault.getAddress(), implSlot, "latest",
    ]);
    expect(("0x" + currentImplStorage.slice(-40)).toLowerCase()).to.equal(newImplAddr.toLowerCase());

    // 验证升级后用户余额保持不变（storage 被保留）
    const balanceAfter = await ctx.vault.balanceOf(await ctx.user1.getAddress());
    expect(balanceAfter).to.equal(balanceBefore);

    // 验证升级后功能仍然正常
    const total = await ctx.vault.totalAssets();
    expect(total).to.equal(depositAmount);

    // 验证提款仍然正常
    const withdrawalShares = toUsdc(50n);
    const requestTx = await ctx.vault.connect(ctx.user1).requestWithdraw(withdrawalShares);
    await requestTx.wait();

    // 快进时间
    const delay = await ctx.vault.withdrawalDelay();
    const increment = delay + 3600n;
    await ethers.provider.send("evm_increaseTime", [Number(increment)]);
    await ethers.provider.send("evm_mine", []);

    // 执行提款
    await ctx.vault.connect(ctx.user1).executeWithdraw();
    const balanceFinal = await ctx.vault.balanceOf(await ctx.user1.getAddress());
    expect(balanceFinal).to.be.lessThan(balanceBefore);
  });
});

// ============================================================
// 15. role access
// ============================================================

describe("PolyVault — 角色权限", () => {
  let ctx: VaultContext;
  before(async () => { ctx = await deployBase(); });

  it("策略师不能暂停合约", async () => {
    await expect(ctx.vault.connect(ctx.strategist).pause()).to.be.revertedWithCustomError(ctx.vault, "AccessControlUnauthorizedAccount");
  });

  it("守护者不能调用策略提款", async () => {
    await expect(
      ctx.vault.connect(ctx.guardian).withdrawToStrategy(toUsdc(50n)),
    ).to.be.revertedWithCustomError(ctx.vault, "AccessControlUnauthorizedAccount");
  });

  it("普通用户不能调用管理员功能", async () => {
    await expect(
      ctx.vault.connect(ctx.user1).setWithdrawalDelay(7200n),
    ).to.be.revertedWithCustomError(ctx.vault, "AccessControlUnauthorizedAccount");
    await expect(
      ctx.vault.connect(ctx.user1).setPerformanceFee(500n),
    ).to.be.revertedWithCustomError(ctx.vault, "AccessControlUnauthorizedAccount");
    await expect(
      ctx.vault.connect(ctx.user1).setFeeRecipient(await ctx.user2.getAddress()),
    ).to.be.revertedWithCustomError(ctx.vault, "AccessControlUnauthorizedAccount");
    await expect(
      ctx.vault.connect(ctx.user1).setMaxStrategyAllocation(7500n),
    ).to.be.revertedWithCustomError(ctx.vault, "AccessControlUnauthorizedAccount");
    await expect(
      ctx.vault.connect(ctx.user1).setDepositLimits(toUsdc(1n), toUsdc(100n)),
    ).to.be.revertedWithCustomError(ctx.vault, "AccessControlUnauthorizedAccount");
  });
});

// ============================================================
// 16. edge cases
// ============================================================

describe("PolyVault — 边界情况", () => {
  let ctx: VaultContext;
  beforeEach(async () => { ctx = await deployBase(); });

  it("支持 请求→取消→重新请求 流程", async () => {
    await ctx.usdc
      .connect(ctx.user1)
      .approve(await ctx.vault.getAddress(), toUsdc(100n));
    await ctx.vault.connect(ctx.user1).deposit(toUsdc(100n), await ctx.user1.getAddress());

    await ctx.vault.connect(ctx.user1).requestWithdraw(toUsdc(30n));
    await ctx.vault.connect(ctx.user1).cancelWithdraw();
    await ctx.vault.connect(ctx.user1).requestWithdraw(toUsdc(40n));

    const req = await ctx.vault.getWithdrawalRequest(await ctx.user1.getAddress());
    expect(req.shares).to.equal(toUsdc(40n));
  });

  it("多个用户同时请求提款", async () => {
    await ctx.usdc
      .connect(ctx.user1)
      .approve(await ctx.vault.getAddress(), toUsdc(100n));
    await ctx.vault.connect(ctx.user1).deposit(toUsdc(100n), await ctx.user1.getAddress());

    const user2Addr = await ctx.user2.getAddress();
    await ctx.usdc.mint(user2Addr, toUsdc(200n));
    await ctx.usdc
      .connect(ctx.user2)
      .approve(await ctx.vault.getAddress(), toUsdc(200n));
    await ctx.vault.connect(ctx.user2).deposit(toUsdc(200n), user2Addr);

    await ctx.vault.connect(ctx.user1).requestWithdraw(toUsdc(30n));
    await ctx.vault.connect(ctx.user2).requestWithdraw(toUsdc(50n));

    const req1 = await ctx.vault.getWithdrawalRequest(await ctx.user1.getAddress());
    expect(req1.shares).to.equal(toUsdc(30n));
    expect(req1.pending).to.be.true;

    const req2 = await ctx.vault.getWithdrawalRequest(user2Addr);
    expect(req2.shares).to.equal(toUsdc(50n));
    expect(req2.pending).to.be.true;
  });

  it("请求全部余额", async () => {
    await ctx.usdc
      .connect(ctx.user1)
      .approve(await ctx.vault.getAddress(), toUsdc(50n));
    await ctx.vault.connect(ctx.user1).deposit(toUsdc(50n), await ctx.user1.getAddress());

    await ctx.vault.connect(ctx.user1).requestWithdraw(toUsdc(50n));

    expect(await ctx.vault.balanceOf(await ctx.user1.getAddress())).to.equal(0);
    expect(
      (await ctx.vault.getWithdrawalRequest(await ctx.user1.getAddress())).pending,
    ).to.be.true;
  });
});

// ============================================================
// 17. integration
// ============================================================

describe("PolyVault — 集成测试", () => {
  let ctx: VaultContext;
  beforeEach(async () => { ctx = await deployBase(); });

  it("完整 存款→请求→取消→重新请求→提款 流程", async () => {
    const userAddr = await ctx.user1.getAddress();

    await ctx.usdc
      .connect(ctx.user1)
      .approve(await ctx.vault.getAddress(), toUsdc(100n));
    await ctx.vault.connect(ctx.user1).deposit(toUsdc(100n), userAddr);
    expect(await ctx.vault.balanceOf(userAddr)).to.equal(toUsdc(100n));

    await ctx.vault.connect(ctx.user1).requestWithdraw(toUsdc(50n));
    expect(await ctx.vault.balanceOf(userAddr)).to.equal(toUsdc(50n));

    await ctx.vault.connect(ctx.user1).cancelWithdraw();
    expect(await ctx.vault.balanceOf(userAddr)).to.equal(toUsdc(100n));

    await ctx.vault.connect(ctx.user1).requestWithdraw(toUsdc(50n));

    await ethers.provider.send("evm_increaseTime", [Number(DEFAULT_WITHDRAWAL_DELAY) + 1]);
    await ethers.provider.send("evm_mine");

    const usdcBefore = await ctx.usdc.balanceOf(userAddr);
    await ctx.vault.connect(ctx.user1).executeWithdraw();

    expect(await ctx.usdc.balanceOf(userAddr) - usdcBefore).to.equal(toUsdc(50n));
    expect(await ctx.vault.balanceOf(userAddr)).to.equal(toUsdc(50n));
  });

  it("完整策略盈利周期", async () => {
    await ctx.usdc.mint(await ctx.user1.getAddress(), toUsdc(900n));
    await ctx.usdc
      .connect(ctx.user1)
      .approve(await ctx.vault.getAddress(), toUsdc(1000n));
    await ctx.vault.connect(ctx.user1).deposit(toUsdc(1000n), await ctx.user1.getAddress());

    await ctx.vault.connect(ctx.strategist).withdrawToStrategy(toUsdc(500n));
    expect(await ctx.vault.strategyDebt()).to.equal(toUsdc(500n));

    await ctx.usdc.mint(await ctx.strategist.getAddress(), toUsdc(100n));
    await ctx.usdc
      .connect(ctx.strategist)
      .approve(await ctx.vault.getAddress(), toUsdc(600n));
    await ctx.vault.connect(ctx.strategist).depositFromStrategy(toUsdc(600n));

    expect(await ctx.vault.strategyDebt()).to.equal(0);
    expect(await ctx.usdc.balanceOf(await ctx.feeRecipient.getAddress())).to.equal(toUsdc(10n));
    expect(await ctx.vault.totalAssets()).to.equal(toUsdc(1090n));
  });

  it("完整策略亏损周期", async () => {
    await ctx.usdc.mint(await ctx.user1.getAddress(), toUsdc(900n));
    await ctx.usdc
      .connect(ctx.user1)
      .approve(await ctx.vault.getAddress(), toUsdc(1000n));
    await ctx.vault.connect(ctx.user1).deposit(toUsdc(1000n), await ctx.user1.getAddress());

    await ctx.vault.connect(ctx.strategist).withdrawToStrategy(toUsdc(500n));

    await ctx.usdc
      .connect(ctx.strategist)
      .approve(await ctx.vault.getAddress(), toUsdc(300n));
    await ctx.vault.connect(ctx.strategist).depositFromStrategy(toUsdc(300n));

    expect(await ctx.vault.strategyDebt()).to.equal(toUsdc(200n));
    expect(await ctx.usdc.balanceOf(await ctx.feeRecipient.getAddress())).to.equal(0);
  });

  it("should handle multiple strategy cycles (profit + loss)", async () => {
    await ctx.usdc.mint(await ctx.user1.getAddress(), toUsdc(900n));
    await ctx.usdc
      .connect(ctx.user1)
      .approve(await ctx.vault.getAddress(), toUsdc(1000n));
    await ctx.vault.connect(ctx.user1).deposit(toUsdc(1000n), await ctx.user1.getAddress());

    await ctx.vault.connect(ctx.strategist).withdrawToStrategy(toUsdc(300n));
    await ctx.usdc.mint(await ctx.strategist.getAddress(), toUsdc(50n));
    await ctx.usdc
      .connect(ctx.strategist)
      .approve(await ctx.vault.getAddress(), toUsdc(350n));
    await ctx.vault.connect(ctx.strategist).depositFromStrategy(toUsdc(350n));

    expect(await ctx.usdc.balanceOf(await ctx.feeRecipient.getAddress())).to.equal(toUsdc(5n));

    await ctx.vault.connect(ctx.strategist).withdrawToStrategy(toUsdc(200n));
    await ctx.usdc
      .connect(ctx.strategist)
      .approve(await ctx.vault.getAddress(), toUsdc(100n));
    await ctx.vault.connect(ctx.strategist).depositFromStrategy(toUsdc(100n));

    expect(await ctx.vault.strategyDebt()).to.equal(toUsdc(100n));
    expect(await ctx.usdc.balanceOf(await ctx.feeRecipient.getAddress())).to.equal(toUsdc(5n));
  });
});
