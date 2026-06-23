/**
 * 本地部署脚本
 *
 * 支持两种部署模式：
 *   模式 A — EDR 内嵌节点（hardhatMainnet / hardhatOp），Hardhat 自动启动临时节点
 *   模式 B — 外部已启动节点（localhost），需先手动启动节点再运行此脚本
 *
 * 部署流程：
 *   1. 初始化运行时，检查网络环境，非本地网络直接报错退出
 *   2. 部署 MockUSDC（本地测试用 USDC）
 *   3. 通过 UUPS 代理部署 PolyVault 并执行 initialize
 *   4. 将部署信息写入 deployments/<network.name>.json
 *
 * 使用:
 *   # 模式 A：EDR 内嵌节点
 *   npx hardhat run scripts/deploy-local.ts --network hardhatMainnet
 *
 *   # 模式 B：连接已运行的本地节点
 *   # 终端 1：npx hardhat node（或 anvil，anvil --code-size-limit 30000）
 *   # 终端 2：
 *   npx hardhat run scripts/deploy-local.ts --network localhost
 *   # 或自定义 RPC 地址：
 *   LOCAL_RPC_URL=http://127.0.0.1:8545 npx hardhat run scripts/deploy-local.ts --network localhost
 *
 *   支持的 network: hardhatMainnet, hardhatOp, localhost
 */

import hre from "hardhat";
import { upgrades } from "@openzeppelin/hardhat-upgrades";
import * as fs from "fs";
import * as path from "path";
import { fileURLToPath } from "url";

const __filename = fileURLToPath(import.meta.url);
const __dirname = path.dirname(__filename);

// ======================== 常量 ========================

/** 允许部署的本地网络名称 */
const LOCAL_NETWORKS = new Set(["hardhat", "hardhatMainnet", "hardhatOp", "localhost"]);

/** 部署信息输出目录 */
const DEPLOYMENTS_DIR = path.resolve(__dirname, "..", "deployments");

// ======================== 参数 ========================

/** PolyVault initialize 参数（根据实际需要修改这里） */
const VAULT_PARAMS = {
  withdrawalDelay: 3600n,        // 1 小时（范围: 1h ~ 7d）
  maxAllocation: 5000n,          // 50%（基点，10000 = 100%）
  performanceFee: 1000n,         // 10%（基点，最大 2000 = 20%）
};

// ======================== 主逻辑 ========================

async function main() {
  // ── 1. 初始化 Hardhat v3 运行时 ──────────────────────
  const connection = await hre.network.create();
  const { ethers } = connection;
  const upgradesApi = await upgrades(hre, connection);

  // ── 2. 环境检查 ──────────────────────────────────────
  const networkName = connection.networkName;
  if (!LOCAL_NETWORKS.has(networkName)) {
    console.error(`\n  ❌ 禁止部署到非本地网络 "${networkName}"！`);
    console.error(`  ✅ 仅支持: ${[...LOCAL_NETWORKS].join(", ")}\n`);
    process.exitCode = 1;
    return;
  }
  console.log(`\n  🌐 网络: ${networkName}（本地环境，继续部署...）\n`);

  // ── 3. 连接可用性检查（模式 B：已启动的节点必须可达） ──
  try {
    const net = await ethers.provider.getNetwork();
    const blockNumber = await ethers.provider.getBlockNumber();
    console.log(`     📡 RPC 可达, 链 ID: ${net.chainId}, 最新区块: ${blockNumber}\n`);
  } catch (err: any) {
    console.error(`\n  ❌ 无法连接到 ${networkName} 节点！`);
    console.error(`     请先在终端 1 中启动本地节点:`);
    console.error(`       npx hardhat node`);
    console.error(`     然后在终端 2 中运行部署脚本:\n`);
    process.exitCode = 1;
    return;
  }

  // ── 4. 获取签名者 ────────────────────────────────────
  const signers = await ethers.getSigners();
  const [deployer] = signers;
  console.log(`  🔑 部署账户: ${deployer.address}\n`);

  // ── 5. 部署 MockUSDC ─────────────────────────────────
  console.log("  📦 部署 MockUSDC...");
  const MockUSDC = await ethers.getContractFactory("MockUSDC");
  const usdc = await MockUSDC.deploy();
  await usdc.waitForDeployment();
  const usdcAddress = await usdc.getAddress();
  console.log(`     ✅ MockUSDC: ${usdcAddress}`);

  // ── 6. 部署 PolyVault（UUPS 代理） ──────────────────
  console.log("\n  📦 部署 PolyVault (UUPS 代理)...");
  const PolyVaultFactory = await ethers.getContractFactory("PolyVault");

  const admin = deployer.address;
  const strategist = signers[1].address;
  const guardian = signers[2].address;
  const feeRecipient = signers[3].address;

  const vault = await upgradesApi.deployProxy(
    PolyVaultFactory,
    [
      usdcAddress,
      admin,
      strategist,
      guardian,
      feeRecipient,
      VAULT_PARAMS.withdrawalDelay,
      VAULT_PARAMS.maxAllocation,
      VAULT_PARAMS.performanceFee,
    ],
    { initializer: "initialize" },
  );
  await vault.waitForDeployment();

  const proxyAddress = await vault.getAddress();

  // 获取实现合约地址（通过 upgrades 插件 erc1967 命名空间）
  const implAddress = await upgradesApi.erc1967.getImplementationAddress(proxyAddress);

  console.log(`     ✅ PolyVault Proxy:  ${proxyAddress}`);
  console.log(`     ✅ PolyVault Impl:   ${implAddress}`);

  // ── 7. 保存部署信息到 JSON ──────────────────────────
  if (!fs.existsSync(DEPLOYMENTS_DIR)) {
    fs.mkdirSync(DEPLOYMENTS_DIR, { recursive: true });
  }

  const deployRecord = {
    network: networkName,
    timestamp: new Date().toISOString(),
    chainId: (await ethers.provider.getNetwork()).chainId.toString(),
    deployer: deployer.address,
    contracts: {
      MockUSDC: usdcAddress,
      PolyVault: {
        proxy: proxyAddress,
        implementation: implAddress,
      },
    },
    params: {
      admin,
      strategist,
      guardian,
      feeRecipient,
      ...VAULT_PARAMS,
    },
  };

  const outPath = path.join(DEPLOYMENTS_DIR, `${networkName}.json`);
  fs.writeFileSync(
    outPath,
    JSON.stringify(
      deployRecord,
      (_key, value) => (typeof value === "bigint" ? value.toString() : value),
      2,
    ),
  );
  console.log(`\n  📄 部署记录已保存 → ${outPath}\n`);
}

main().catch((err) => {
  console.error(err);
  process.exitCode = 1;
});
