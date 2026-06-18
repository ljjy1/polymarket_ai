import { execSync } from "node:child_process";
import { existsSync, mkdirSync, readdirSync, readFileSync, rmSync, writeFileSync } from "node:fs";
import { join } from "node:path";
import { fileURLToPath } from "node:url";

const __dirname = join(fileURLToPath(import.meta.url), "..");
const rootDir = join(__dirname, "..");
const artifactsDir = join(rootDir, "artifacts", "contracts");
const outputDir = join(__dirname, "files");
const pkgName = process.argv[2] || "bindcode";

// 1. 重新编译合约
console.log("重新编译合约...");
execSync("npx hardhat compile", { cwd: rootDir, stdio: "inherit" });

// 2. 清理并创建输出目录
if (existsSync(outputDir)) {
  rmSync(outputDir, { recursive: true });
}
mkdirSync(outputDir, { recursive: true });

// 3. 遍历 artifacts
const contractFiles = [];
function walk(dir) {
  for (const entry of readdirSync(dir, { withFileTypes: true })) {
    const fullPath = join(dir, entry.name);
    if (entry.isDirectory()) {
      if (entry.name.endsWith(".t.sol")) continue;
      if (entry.name === "libraries" || entry.name === "interfaces") continue;
      walk(fullPath);
    } else if (entry.isFile() && entry.name.endsWith(".json")) {
      if (entry.name.endsWith(".dbg.json")) continue;
      contractFiles.push(fullPath);
    }
  }
}
walk(artifactsDir);

let count = 0;
for (const file of contractFiles) {
  const data = JSON.parse(readFileSync(file, "utf-8"));
  const { contractName, abi, bytecode } = data;

  // 跳过接口/抽象合约（bytecode 为空或仅为 "0x"）
  if (!bytecode || bytecode === "0x") {
    console.log(`  跳过 ${contractName}（无 bytecode）`);
    continue;
  }

  // 创建合约目录
  const contractDir = join(outputDir, contractName);
  mkdirSync(contractDir, { recursive: true });

  // 提取 .abi、.bin 和 .json（包含 abi + bytecode）
  const abiFile = join(contractDir, `${contractName}.abi`);
  const binFile = join(contractDir, `${contractName}.bin`);
  const jsonFile = join(contractDir, `${contractName}.json`);
  writeFileSync(abiFile, JSON.stringify(abi, null, 2), "utf-8");
  writeFileSync(binFile, bytecode, "utf-8");
  writeFileSync(jsonFile, JSON.stringify({ abi, bytecode }, null, 2), "utf-8");

  // 生成 Go 绑定代码（首字母小写作为输出文件名）
  const firstLower = contractName.charAt(0).toLowerCase() + contractName.slice(1);
  const goFile = join(contractDir, `${firstLower}.go`);

  try {
    execSync(
      `abigen --abi "${abiFile}" --bin "${binFile}" --pkg ${pkgName} --type "${contractName}" --out "${goFile}"`,
      { stdio: "inherit" },
    );
  } catch {
    console.log(`  ⚠️  abigen 未安装或执行失败，跳过 Go 代码生成（${contractName}）`);
  }

  count++;
  console.log(`  生成 ${contractName}`);
}

console.log(`\n完成！处理了 ${count} 个合约，文件在 ${outputDir} 下`);
