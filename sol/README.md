# PolyVault — USDC 策略金库

PolyVault 是一个基于 **ERC4626** 标准的可升级 USDC 金库合约，支持延迟提款、策略管理和基于角色的访问控制。合约使用 UUPS 代理模式部署。

## 核心功能

- **存款/铸造** — 存入 USDC 获得生息份额代币（pvUSDC），支持最低/最高存款限额
- **延迟提款** — 两步提款流程：`requestWithdraw()` 锁定份额 → `executeWithdraw()` 在延迟期后提取 USDC；直接 `withdraw()`/`redeem()` 被禁用
- **策略管理** — 策略师可从金库提取 USDC 用于链上交易（如 Polymarket），返还时自动计算并分配业绩费
- **权限控制** — 基于 `AccessControl`：管理员（`DEFAULT_ADMIN_ROLE`）、策略师（`STRATEGIST_ROLE`）、守护者（`GUARDIAN_ROLE`）
- **暂停机制** — 守护者可紧急暂停存款和提款请求
- **UUPS 升级** — 通过 UUPS 代理模式支持合约升级，仅管理员可授权

## 项目结构

```
sol/
├── contracts/
│   ├── interfaces/
│   │   └── IPolyVault.sol          # 金库接口定义（事件/错误/函数签名）
│   ├── mocks/
│   │   └── MockUSDC.sol            # 测试用模拟 USDC（6 位小数）
│   ├── test/
│   │   └── ERC1967ProxyHelper.sol  # Foundry 测试用代理部署辅助合约
│   ├── PolyVault.sol               # 核心金库合约（ERC4626 + 可升级）
│   └── PolyVault.t.sol             # Foundry Solidity 测试（~2300 行）
├── scripts/
│   ├── deploy-local.ts             # 本地网络部署脚本
│   └── export-abi.ts               # ABI 导出脚本
├── test/
│   └── PolyVault.ts                # Hardhat + Mocha TypeScript 测试
├── exports/abi/                    # ABI 导出目录（由 export-abi 生成）
├── deployments/                    # 部署记录目录（由 deploy-local 生成）
├── artifacts/                      # 编译产物（自动生成）
├── cache/                          # 编译缓存（自动生成）
├── hardhat.config.ts               # Hardhat v3 配置
├── tsconfig.json                   # TypeScript 配置
├── package.json                    # 项目依赖
└── .gitignore                      # Git 忽略规则
```

## 前置条件

- Node.js >= 20
- [Foundry](https://book.getfoundry.sh/getting-started/installation)（`forge`、`cast`）

```bash
# 安装依赖
npm install
```

## 编译

```bash
npx hardhat compile
```

## 测试

### 运行全部测试（Foundry + Mocha）

```bash
npx hardhat test
```

### 仅运行 Foundry 测试（Solidity）

```bash
npx hardhat test solidity
```

### 仅运行 Mocha 测试（TypeScript）

```bash
npx hardhat test mocha
```

## 测试覆盖率

Hardhat 3 内置 Solidity 代码覆盖率支持，只需在测试命令后添加 `--coverage` 参数：

```bash
# 运行全部测试并生成覆盖率报告（终端 + HTML）
npx hardhat test --coverage

# 仅对 Solidity 测试生成覆盖率
npx hardhat test solidity --coverage

# 仅对 Mocha 测试生成覆盖率
npx hardhat test mocha --coverage
```

覆盖率报告输出到 `coverage/html/` 目录，浏览器打开即可查看。

## 本地部署

脚本仅允许在本地网络（`hardhatMainnet` / `hardhatOp` / `localhost`）上执行，非本地网络会直接拦截报错。

```bash
# 部署到本地主网模拟环境
npx hardhat run scripts/deploy-local.ts --network hardhatMainnet

# 部署到本地 Optimism 模拟环境
npx hardhat run scripts/deploy-local.ts --network hardhatOp
```

部署完成后，合约地址和参数信息会保存到 `deployments/<network>.json`。

### 部署参数（在 `scripts/deploy-local.ts` 中修改）

| 参数 | 默认值 | 说明 |
|---|---|---|
| `withdrawalDelay` | `3600` (1h) | 提款延迟时间（范围: 1h ~ 7d） |
| `maxAllocation` | `5000` (50%) | 最大策略分配比例（基点，10000 = 100%） |
| `performanceFee` | `1000` (10%) | 业绩费比例（基点，最大 20%） |

## 测试网部署

脚本仅允许在 `sepolia` 网络上执行。USDC 地址通过环境变量传入，不会自动部署 MockUSDC。

### 前置条件

在项目根目录创建 `.env` 文件：

```bash
SEPOLIA_RPC_URL=https://sepolia.infura.io/v3/你的API_KEY
SEPOLIA_PRIVATE_KEY=0x你的私钥
```

### 部署命令

```bash
USDC_ADDRESS=0x1c7D4B196Cb0C7B01d743Fbc6116a902379C7238 \
npx hardhat run scripts/deploy-sepolia.ts --network sepolia
```

可选通过环境变量指定角色地址（默认使用部署账户）：

```bash
USDC_ADDRESS=0x1c7D4B196Cb0C7B01d743Fbc6116a902379C7238 \
STRATEGIST_ADDRESS=0x... \
GUARDIAN_ADDRESS=0x... \
FEE_RECIPIENT=0x... \
npx hardhat run scripts/deploy-sepolia.ts --network sepolia
```

部署完成后，合约地址和参数信息会保存到 `deployments/sepolia.json`。

## 提取ABI和BIN 后续用于后端连接 (需要合约已经编译了)
```shell
#在终端 sol目录下运行
jq '.abi' artifacts/contracts/PolyVault.sol/PolyVault.json > PolyVault.abi  
 
jq -r '.bytecode' artifacts/contracts/PolyVault.sol/PolyVault.json > PolyVault.bin
```

```shell
#脚本生成 go合约绑定文件和对应abi和bin json 需要提前安装好jq和abigen
#在终端 sol目录下运行
# 默认 --pkg bindcode
node abi-bin-go-code/extract.mjs
# 自定义包名mycustompkg
node abi-bin-go-code/extract.mjs mycustompkg
```

