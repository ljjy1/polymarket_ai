/**
 * Sepolia 测试网部署脚本
 *
 * 仅允许在 sepolia 网络上部署。USDC 地址通过环境变量传入，不自动部署。
 * 部署流程：
 *   1. 检查网络环境，非 sepolia 直接报错退出
 *   2. 检查 USDC_ADDRESS 环境变量
 *   3. 通过 UUPS 代理部署 PolyVault 并执行 initialize
 *   4. 将部署信息写入 deployments/<network.name>.json
 *
 * 使用:
 *   USDC_ADDRESS=0x... \
 *   npx hardhat run scripts/deploy-sepolia.ts --network sepolia
 *
 * 可选环境变量（默认使用部署账户）:
 *   STRATEGIST_ADDRESS=0x...   策略师地址
 *   GUARDIAN_ADDRESS=0x...     守护者地址
 *   FEE_RECIPIENT=0x...        业绩费接收地址
 */

import "dotenv/config";
import hre from "hardhat";
import { upgrades } from "@openzeppelin/hardhat-upgrades";
import * as fs from "fs";
import * as path from "path";
import { fileURLToPath } from "url";

const __filename = fileURLToPath(import.meta.url);
const __dirname = path.dirname(__filename);

// ======================== 常量 ========================

/** 允许部署的网络名称 */
const ALLOWED_NETWORK = "sepolia";

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
  if (networkName !== ALLOWED_NETWORK) {
    console.error(`\n  ❌ 禁止部署到 "${networkName}" 网络！`);
    console.error(`  ✅ 仅支持: ${ALLOWED_NETWORK}\n`);
    process.exitCode = 1;
    return;
  }
  console.log(`\n  🌐 网络: ${networkName}（继续部署...）\n`);

  // ── 3. 检查 USDC 地址 ────────────────────────────────
  const usdcAddress = process.env.USDC_ADDRESS;
  if (!usdcAddress) {
    console.error(`\n  ❌ 未设置 USDC_ADDRESS 环境变量！`);
    console.error(`  ✅ 请设置: export USDC_ADDRESS=0x...\n`);
    process.exitCode = 1;
    return;
  }
  if (!ethers.isAddress(usdcAddress)) {
    console.error(`\n  ❌ USDC_ADDRESS 格式无效: ${usdcAddress}\n`);
    process.exitCode = 1;
    return;
  }
  console.log(`  💰 USDC: ${usdcAddress}\n`);

  // ── 4. 获取签名者 ────────────────────────────────────
  const signers = await ethers.getSigners();
  const deployer = signers[0];
  console.log(`  🔑 部署账户: ${deployer.address}\n`);

  // ── 5. 读取角色地址（可覆盖，默认使用部署账户）───────
  const admin = deployer.address;
  const strategist = process.env.STRATEGIST_ADDRESS || deployer.address;
  const guardian = process.env.GUARDIAN_ADDRESS || deployer.address;
  const feeRecipient = process.env.FEE_RECIPIENT || deployer.address;

  console.log(`  👤 Admin:         ${admin}`);
  console.log(`  👤 Strategist:    ${strategist}`);
  console.log(`  👤 Guardian:      ${guardian}`);
  console.log(`  👤 Fee Recipient: ${feeRecipient}\n`);

  // ── 6. 部署 PolyVault（UUPS 代理） ──────────────────
  console.log("  📦 部署 PolyVault (UUPS 代理)...");
  const PolyVaultFactory = await ethers.getContractFactory("PolyVault");

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
    { initializer: "initialize", kind: "uups" },
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
      USDC: usdcAddress,
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
