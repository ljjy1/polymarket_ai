import "dotenv/config";
import hardhatToolboxMochaEthersPlugin from "@nomicfoundation/hardhat-toolbox-mocha-ethers";
import hardhatUpgrades from "@openzeppelin/hardhat-upgrades";
import hardhatIgnoreWarnings from "hardhat-ignore-warnings";
import {defineConfig } from "hardhat/config";

function env(name: string): string {
  const val = process.env[name];
  if (!val) {
    throw new Error(`${name} 未设置，请在 .env 文件中添加 ${name}=...`);
  }
  return val;
}

export default defineConfig({
  plugins: [hardhatToolboxMochaEthersPlugin, hardhatUpgrades, hardhatIgnoreWarnings],
  warnings: {
    "*": {
      "transient-storage": "off",
    },
  },
  solidity: {
    profiles: {
      default: {
        version: "0.8.28",
      },
      production: {
        version: "0.8.28",
        settings: {
          optimizer: {
            enabled: true,
            runs: 200,
          },
        },
      },
    },
  },
  networks: {
    hardhatMainnet: {
      type: "edr-simulated",
      chainType: "l1",
    },
    hardhatOp: {
      type: "edr-simulated",
      chainType: "op",
    },
    sepolia: {
      type: "http",
      chainType: "l1",
      url: env("SEPOLIA_RPC_URL"),
      accounts: [env("SEPOLIA_PRIVATE_KEY")],
    },
  },
});
