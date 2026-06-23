# be — Polymarket AI 自动化预测与交易引擎

## 概述

基于 Polymarket 预测市场的 AI 自动化交易系统。通过多维度数据聚合（技术指标、市场情绪、新闻事件、链上数据），结合 **DeepSeek 多智能体架构** 生成校准概率预测，自动制定交易策略并执行。

核心流程：**市场扫描 → 数据聚合 → AI 预测 → 策略生成 → 交易执行 → 结算归还**

金库合约 PolyVault（Polygon 链上）管理整个资金流：策略从金库提取 USDC 下单，完成交易后归还资金到金库，形成闭环。

> **设计参考**: 本系统的 Python 原型版是 `polymarket_ai/backend/`，Go 版是其工程化重写。

### 六层架构总览

```
外源数据 (Polymarket Gamma / Binance / alternative.me / GNews / CryptoQuant / Polygon链上)
    │
    ▼
┌─ External/Adapter 层 ── HTTP/WebSocket/合约客户端 ───────────┐
│   binance.KlineAPI      polymarket.CLOB API                     │
│   polymarket.WebSocket  external.GNewsFetcher                   │
│   external.Fetcher      contract.VaultContractClient             │
└────────────────────────────────────────────────────────────────┘
    │
    ▼
┌─ Service 层 ── 核心业务逻辑 ───────────────────────────────────┐
│                                                                    │
│  MarketScanner  → 扫描 Polymarket，选出最佳 BTC 市场             │
│  DataAggregator → Binance 1h/1d K线 + 贪恐指数 + 新闻 + 链上    │
│  AIPredictor    → DeepSeek 校准概率预测（单模型）                │
│  StrategyGenerator → 三闸门过滤 + half-Kelly 仓位计算            │
│  TradeExecutor  → Vault 提款 → Polymarket CLOB 下单              │
│  PositionMonitor → 止盈 / 止损 / 预解析退出                      │
│  OrderBookMonitor → WebSocket 实时盘口 + 市场结算事件            │
└────────────────────────────────────────────────────────────────┘
    │
    ▼
┌─ Tasks 层 ── Asynq 异步任务 + 定时流水线 ────────────────────┐
│   daily:scan(08:00) → predict → strategy → execute             │
│   monitor:positions(5min) / settlement(1h) / vault(30min)      │
│   WebSocket 实时: market_resolved 可提前标记结算                │
└────────────────────────────────────────────────────────────────┘
    │
    ▼
┌─ Handler/API 层 ── Gin RESTful + Swagger ──────────────────────┐
│   /api/v1/markets  /predictions  /strategies  /trades  /tasks  │
│   /health  /metrics  /swagger/index.html                        │
└────────────────────────────────────────────────────────────────┘
    │
    ▼
┌─ DAO/Model 层 ── GORM + MySQL ────────────────────────────────┐
│   markets / predictions / strategies / trades / vault_snapshots  │
│   system_logs / users                                            │
└────────────────────────────────────────────────────────────────┘
    │
    ▼
┌─ Database / Cache 层 ── MySQL + Redis ─────────────────────────┐
│   MySQL: 业务持久化                                              │
│   Redis DB0: 缓存  Redis DB1: Asynq 任务队列                     │
└────────────────────────────────────────────────────────────────┘
```

---

## 技术栈

| 组件 | 技术 |
|------|------|
| 语言 | Go 1.22+ |
| Web 框架 | Gin + Sponge |
| 数据库 | MySQL 8.0+ (GORM) |
| 缓存 / 任务队列 | Redis 7+ (Asynq) |
| AI 模型 | DeepSeek (Eino ChatModel) |
| 数据源 | Binance API (行情)、alternative.me (恐惧贪婪指数)、GNews (新闻)、CryptoQuant (链上) |
| 链上交互 | Polymarket CLOB API / Gamma API / WebSocket、PolyVault 合约 (Polygon) |
| 合约交互 | Ethereum Go bindings (abigen 生成) |
| 监控 | Prometheus + Grafana (内置 Metrics) |
| 链路追踪 | Jaeger (可选) |
| 部署 | 裸机 / Docker / Kubernetes |

---

## 目录结构

```
├── cmd/be/                         # 程序入口
│   ├── main.go                     # 入口 (默认 localhost:8080)
│   └── initial/
│       ├── initApp.go              # 初始化配置/数据库/Redis/日志/链路追踪
│       ├── createService.go        # 组装所有服务 (HTTP/WS/Asynq)
│       └── close.go                # 资源释放
│
├── internal/
│   ├── binance/                    # Binance 行情 API 客户端
│   │   └── client.go              #   K线 / 24h Ticker / 技术指标计算
│   ├── external/                   # 外部 API 客户端
│   │   └── fetcher.go             #   GNews (新闻) + CryptoQuant (链上)
│   ├── polymarket/                 # Polymarket 交互
│   │   ├── client.go              #   CLOB REST API + Gamma API + 事件查询
│   │   └── websocket.go           #   WebSocket 盘口客户端 (book/price_change/best_bid_ask/market_resolved/...)
│   │
│   ├── config/be.go               # 配置结构体 (YAML → Go struct)
│   ├── database/                   # MySQL / Redis 连接管理
│   ├── bindcode/                   # 智能合约 Go 绑定 (abigen 生成)
│   │   └── polyVault.go           #   PolyVault 合约的读/写接口
│   ├── contract/                   # 合约客户端封装
│   │   └── client.go              #   VaultContractClient (只读 + 交易方法)
│   ├── model/                      # GORM 数据模型
│   │   ├── markets.go             #   预测市场 (event_slug, polymarket_token_id, threshold, prices)
│   │   ├── predictions.go         #   AI 预测结果 (概率/置信度/edge/raw_request/raw_response/data_snapshot)
│   │   ├── strategies.go          #   交易策略 (action/position_size/entry/take_profit/stop_loss/kelly_fraction)
│   │   ├── trades.go              #   交易记录 (order_id/filled_amount/price/pnl/close_reason)
│   │   ├── vault_snapshots.go     #   金库快照
│   │   ├── system_logs.go         #   系统日志
│   │   └── user.go                #   钱包用户
│   ├── dao/                        # 数据访问层 (GORM)
│   ├── cache/                      # Redis 缓存层
│   │
│   ├── service/                    # 核心业务逻辑
│   │   ├── market_scanner.go      #   扫描 Polymarket 市场 + 候选排序
│   │   ├── data_aggregator.go     #   多源数据聚合 (含1h K线技术指标)
│   │   ├── ai_predictor.go        #   AI 预测 (DeepSeek, 结构化输出)
│   │   ├── strategy_generator.go  #   Kelly公式策略生成 (三闸门)
│   │   ├── trade_executor.go      #   策略执行/CLOB 下单
│   │   ├── orderbook_monitor.go   #   WS 盘口实时监控 (app.IServer)
│   │   ├── position_monitor.go    #   持仓止盈止损监控
│   │   ├── vault_service.go       #   金库资产/余额查询
│   │   └── indicators.go          #   技术指标计算 (RSI/MACD/布林带/EMA/ATR)
│   │
│   ├── tasks/                      # Asynq 异步任务
│   │   ├── client.go             #   任务客户端 (Enqueue / EnqueueIn)
│   │   ├── server.go             #   任务处理服务器
│   │   ├── daily.go              #   日常流水线 (scan/predict/strategy/execute)
│   │   ├── monitor.go            #   监控任务 (positions/vault/health/settlement)
│   │   └── task_api.go           #   HTTP 手动触发的桥梁层
│   │
│   ├── handler/                    # HTTP 请求处理
│   ├── routers/                    # 路由注册
│   ├── types/                      # 请求/响应结构体
│   ├── ecode/                      # 错误码定义
│   └── server/                     # HTTP 服务配置 (中间件/超时/CORS)
│
├── configs/
│   └── be.yml                      # 主配置文件
├── deployments/                    # 部署
│   ├── binary/                    #   裸机部署
│   ├── docker-compose/            #   Docker Compose
│   └── kubernetes/                #   K8S 部署
├── scripts/
├── docs/                           # Swagger 文档 (自动生成)
├── Makefile
└── README.md
```

---

## 应用启动流程

```
main.go
  └─ InitApp()
       ├─ 解析命令行参数 -c configs/be.yml
       ├─ 展开环境变量 ${VAR} → 实际值
       ├─ 初始化日志 (zap, 控制台/文件)
       ├─ 初始化链路追踪 (Jaeger, 可选)
       ├─ 初始化系统资源统计
       └─ 初始化 MySQL + Redis 连接
  └─ CreateServices()
       ├─ 创建 Gin HTTP 服务
       ├─ 创建 WebSocket 盘口监控 (polyClient + marketsDao)
       ├─ 创建 Asynq Scheduler (定时任务)
       ├─ 创建 DailyDeps + MonitorDeps
       │    ├─ 初始化 VaultContractClient
       │    │    └─ 有私钥时包含 Transactor (write 方法)
       │    ├─ 初始化 Polymarket 客户端
       │    └─ 注入到 TradeExecutor / PositionMonitor / DataAggregator
       ├─ 设置全局任务依赖 (供 HTTP API 使用)
       └─ 注册 Asynq 任务处理器
  └─ app.New(services, closes).Run()
       ├─ 启动 HTTP 服务
       ├─ 启动 WebSocket 盘口监控
       ├─ 启动 Asynq Scheduler (定时唤醒)
       └─ 启动 Asynq Server (处理队列)
```

---

## 核心业务详解

### 1. 市场扫描 — MarketScanner（定时 08:00 UTC）

扫描 Polymarket Gamma API，选出最适合 AI 交易的 BTC 预测市场。

**数据流：**
1. 调用 `GET /events?tag_slug=bitcoin&active=true` → 获取所有活跃 BTC 事件
2. 事件标题必须匹配 `/bitcoin\s+above/i` 正则
3. 目标的结算日期必须等于 **扫描日 + 2天**（T+2 规则）
4. 对每个事件，调用 `GET /events/{id}` 获取详情，拿到完整 `question` 字段
5. 从 `question` 中通过正则 `/above\s*\$?([0-9,.]+)/i` 提取价格阈值

**候选排序（6步过滤）：**

| 步骤 | 条件 | 作用 |
|------|------|------|
| 1. 目标日期 | `target_date == scan_date + 2 days` | 排除远期市场，T+2 小结算窗口更适合 AI |
| 2. 标题匹配 | `question` 含 "Bitcoin above $X" | 仅关注 BTC 门槛预测 |
| 3. 价格范围 | Yes 价格 ∈ **[0.35, 0.65]** | 排除确定性太高或太低的市场 |
| 4. 最小交易量 | 总交易量 ≥ **$10,000** | 排除流动性不足的市场 |
| 5. 排序 | `|yes_price - 0.5|` 升序 + 交易量降序 | 越接近 50% 越优先 |
| 6. 取最优 | 排序后第一个 | 每天只处理一个最具信噪比的市场 |

**幂等性：** 以 `scan_date` 为幂等键，同一天的扫描不会重复写入。

---

### 2. 数据聚合 — DataAggregator（预测前执行）

聚合 AI 预测所需的全维度上下文数据。

#### 外部数据源

| 数据源 | API 端点 | 获取数据 | 用途 |
|--------|---------|---------|------|
| **Binance** | `GET /api/v3/klines` | BTCUSDT 1小时K线(168根=7天) | 计算技术指标 |
| **Binance** | `GET /api/v3/klines` | BTCUSDT 日K线(30根) | 7日涨跌幅 |
| **Binance** | `GET /api/v3/ticker/24hr` | 24h价格变化、最高/低、成交量 | 基本面数据 |
| **alternative.me** | `GET /fng/?limit=7` | 今日恐惧贪婪指数 + 7日趋势 | 市场情绪 |
| **GNews** | `GET /search?q=bitcoin&max=10` | 最近BTC新闻标题 | 新闻影响分析 |
| **CryptoQuant** | `GET /v3/top_by_entity/...` | 交易所净流量、大额交易、活跃地址 | 链上分析 |

#### 技术指标（基于1小时K线）

| 指标 | 参数 | 说明 |
|------|------|------|
| RSI | 14周期 | 超买/超卖判断 |
| MACD | 12/26/9 | 趋势动能 + 信号线交叉 |
| 布林带 | 20周期, 2标准差 | 波动率 + 支撑/阻力 |
| EMA | 7 / 25 / 99 | 短/中/长期趋势 |
| ATR | 14周期 | 平均真实波幅 |

**容错：** 新闻获取失败降级为空列表继续执行，不阻断流水线。

**聚合产物：** `MarketDataBundle` — 包含 BTC 现货价、目标阈值、距结算时间、24h 高/低/量、所有技术指标、恐惧贪婪指数+7日趋势、新闻标题、链上数据。

---

### 3. AI 预测 — AIPredictor

使用 DeepSeek ChatModel 生成概率预测。

**提示词设计（参考 Python 参考版）：**

**System Prompt** 核心规则（6条）：
1. 输出**校准概率** (0-1)，而非置信度加权赌注
2. 仅使用提供的数据，不编造指标
3. `predicted_probability` 与 `confidence` 是两个独立维度
4. 短结算窗口（48小时内）的确定性更高
5. 仅输出符合 JSON Schema 的有效 JSON，禁止 markdown
6. **警告**：下游系统仅在 `abs(edge) >= 0.25` 且 `confidence >= 0.6` 时才会交易

**User Prompt** 包含：
- 市场问题 + 目标日期 + 距结算时间
- Polymarket 当前 Yes/No 价格 + 交易量
- BTC 现货价 / 距阈值百分比 / 涨跌幅 / 高/低价
- 技术指标文本描述
- 恐惧贪婪指数 + 7日趋势
- 新闻标题列表
- 链上数据（交易所净流量、大额交易、活跃地址变化）

**结构化输出 Schema：**
```json
{
  "predicted_probability": 0.42,   // 校准概率 [0.01, 0.99]
  "confidence": 0.65,              // 置信度 [0, 1]
  "direction": "bearish_below",    // 方向判断
  "key_factors": ["..."],           // 关键因素列表
  "risk_factors": ["..."],         // 风险因素列表
  "technical_analysis": "...",     // 技术分析文本
  "sentiment_analysis": "...",     // 情绪分析文本
  "news_impact": "...",           // 新闻影响
  "onchain_analysis": "...",      // 链上分析
  "reasoning": "...",             // 推理过程
  "recommended_action": "skip"    // buy_yes / buy_no / skip
}
```

**防御性解析：** 即使模型返回 markdown 表格/散文格式，解析器会依次尝试：清除 markdown 包裹 → 直接解析 JSON → 从任意文本中提取 `{...}` JSON 对象。失败后重试一次。

---

### 4. 策略生成 — StrategyGenerator

**三闸门硬性过滤（与 Python 参考版一致）：**

```
                    ┌──────────────┐
       Prediction   │              │
       ────────────→│   Gate 1    │ ← abs(edge) >= 0.25
                    │  Edge 检查   │    (AI概率 - 市场价)
                    └──────┬───────┘
                           │ 通过
                           ▼
                    ┌──────────────┐
                    │   Gate 2    │ ← confidence >= 0.6
                    │  置信度检查  │
                    └──────┬───────┘
                           │ 通过
                           ▼
                    ┌──────────────┐
                    │   Gate 3    │ ← recommended_action
                    │  行动一致性  │    与 edge 符号一致
                    └──────┬───────┘
                           │ 通过
                           ▼
                    Kelly 仓位计算
```

**仓位计算（通过三个闸门后）：**

| 参数 | Python参考值 | Go 项目值 | 说明 |
|------|-------------|-----------|------|
| `kellyMultiplier` | 0.5 | 0.5 | half-Kelly，降低风险 |
| `maxPositionPct` | 0.10 | 0.10 | 最大仓位 = 金库余额的10% |
| `takeProfitFactor` | 0.7 | 0.7 | 止盈 = 入场价 + `|edge| × 0.7` |
| `stopLossFactor` | 0.5 | 0.5 | 止损 = 入场价 − `|edge| × 0.5` |

策略最终状态：`pending` → `executing` → `active` / `skipped` / `failed` → `closed`。

---

### 5. 交易执行 — TradeExecutor

1. 从 **金库合约** 提取策略仓位对应的 USDC
2. 通过 **Polymarket CLOB API** 以限价单买入对应 side 的份额
3. 至少部分成交即视为成功，记录 `trades` 表
4. 策略状态标记为 `active`

---

### 6. 仓位监控 — PositionMonitor（每5秒循环）

对每个活跃策略的持仓做三件事：

```
PositionMonitor.CheckStrategies()
  │
  ├─ 1. 预解析退出: 距结算 ≤ 30分钟 → 强制平仓
  │    (避免结算前低流动性震荡)
  │
  ├─ 2. 止盈检查: 当前价格 ≥ take_profit → 卖出
  │
  └─ 3. 止损检查: 当前价格 ≤ stop_loss → 卖出
```

**价格获取：** 通过 Polymarket CLOB API 实时获取订单簿买卖中间价，不再依赖 DB 中缓存的静态价格。

---

### 7. 实时盘口监控 — OrderBookMonitor（WebSocket，持续运行）

通过 Polymarket WebSocket 实时订阅活跃市场的盘口数据。

#### 订阅的事件

| 事件类型 | 触发条件 | 业务用途 |
|---------|---------|---------|
| **book** | 盘口快照 | 初始价格缓存 |
| **price_change** | 盘口价位变动 | 实时更新内存价格 |
| **best_bid_ask** | 最优买卖价变动（需 custom_feature） | 直接用 server 计算的精准 bid/ask 更新 |
| **last_trade_price** | 最新成交 | 追踪成交价偏离盘口中间价，>15% 触发预警 |
| **new_market** | 新市场创建（需 custom_feature） | 检查 tag 含 bitcoin/btc 的潜在新市场 |
| **market_resolved** | 市场结算（需 custom_feature） | 更新 DB 状态为 resolved（带并发控制） |
| **tick_size_change** | 最小报价单位变动 | 日志记录 |

#### 并发控制

- **`market_resolved`**：使用 `sync.Map` 记录已处理的 `event_slug`，WS 重连重复接收也不重复执行
- **内存价格更新**：通过 `m.mu` 读写锁保护
- **定时任务保留**：cron `HandleSettlement`（每小时）仍正常运行。WS 的 `market_resolved` 仅更新 DB 状态，PnL 结算由 cron 统一处理

#### 数据库价格同步

`OrderBookMonitor` 会定期将内存中的最新价格批量回写 `markets` 表的 `current_yes_price` / `current_no_price` 字段。

---

## 业务流程全景

### 每日流水线（任务链驱动）

只有 `daily:scan` 有定时 cron，其余任务通过 Asynq 任务链自动触发：

```
08:00 UTC cron → market_scan
      └─ MarketScanner.Scan()
         ├─ Polymarket Gamma API 获取活跃市场
         ├─ 6步过滤 + 排序 → 选出最佳候选
         ├─ 写入 markets 表 (status=active)
         └─ 入队 predict 任务 (10分钟延迟)

08:10+ → ai_predict (延迟10分钟，等待数据就绪)
      └─ DataAggregator.Aggregate(market)
      │    ├─ Binance: 1h K线 + 技术指标
      │    ├─ GNews: 新闻头条
      │    ├─ CryptoQuant: 链上数据
      │    └─ alternative.me: 恐惧贪婪 + 7日趋势
      ├─ AIPredictor.Predict(bundle)
      │    ├─ System Prompt 6条规则
      │    ├─ User Prompt 全维度数据
      │    └─ Structured Output 解析 (含防御性回退)
      ├─ 写入 predictions 表
      └─ 入队 strategy 任务 (立即执行)

→ strategy_generator
      ├─ 读取金库链上余额
      ├─ 三闸门过滤 (edge ≥ 0.25 / confidence ≥ 0.6 / 行动一致性)
      ├─ half-Kelly 计算仓位
      ├─ 写入 strategies 表
      └─ 非 skip → 入队 execute 任务

→ execute
      ├─ Vault.withdrawToStrategy() → 提取资金
      ├─ Polymarket CLOB 限价单执行
      └─ 记录 trades 表 / 策略状态 active
```

### 持续监控定时任务

| 频率 | 任务 | 说明 |
|------|------|------|
| 实时 | WebSocket 盘口 | 内存价格 + 事件处理 (book, best_bid_ask, market_resolved, ...) |
| 每5秒 | 持仓监控 (循环) | 止盈 / 止损 / 预解析退出 |
| 每1小时 | 结算检查 (cron) | 遍历活跃策略 → 检查 resolution → 计算 PnL → 归还资金 |
| 每30分钟 | 金库快照 (cron) | 链上只读方法 → 写入 vault_snapshots 表 |
| 每10分钟 | 健康检查 (cron) | 数据库连通性 Ping |

---

## WebSocket 事件详解

### SDK 支持的事件类型

多亏了 SDK 的 `HandleUnknown` 兜底 + `custom_feature_enabled=true`，以下事件均被接收：

| 事件 | SDK 回调 | 数据结构 |
|------|---------|---------|
| `book` | `OnBook` | `{asset_id, bids[{price,size}], asks[{price,size}]}` |
| `price_change` | `OnPriceChange` | `{asset_id, changes[{price,size,side,best_bid,best_ask}]}` |
| `last_trade_price` | `OnLastTradePrice` | `{asset_id, price, side, size}` |
| `tick_size_change` | `OnTickSizeChange` | `{asset_id, market, old_tick_size, new_tick_size}` |
| `best_bid_ask` | `HandleUnknown` | `{asset_id, market, best_bid, best_ask, spread}` |
| `new_market` | `HandleUnknown` | `{id, question, market, slug, assets_ids, outcomes, tags}` |
| `market_resolved` | `HandleUnknown` | `{id, question, market, slug, winning_asset_id, winning_outcome}` |

### 自动重连

断连后 3 秒自动重试，无限重试，日志记录每次断连原因。

---

## 链上金库 (PolyVault)

系统通过 PolyVault 合约管理 USDC 资金，形成"金库 → 策略 → 交易 → 结算 → 金库"闭环。

### 合约交互

| 方法 | 类型 | 用途 | 调用时机 |
|------|------|------|---------|
| `availableBalance()` | view | 金库可用 USDC | 策略生成 / 金库快照 |
| `totalAssets()` | view | 总资产 (余额+债务) | 金库快照 |
| `totalSupply()` | view | 总份额发行量 | 金库快照 |
| `strategyDebt()` | view | 已部署策略资金 | 金库快照 |
| `convertToAssets(shares)` | view | 份额→资产转换 | 金库快照 |
| `withdrawToStrategy(amount)` | 写 (STRATEGIST_ROLE) | 提取到策略 | 交易执行 |
| `depositFromStrategy(amount)` | 写 (STRATEGIST_ROLE) | 归还到金库 | 结算检查 |

### 客户端封装

- **只读方法** (`VaultReader`)：始终可用，不需要私钥
- **交易方法** (`VaultTransactor`)：仅在配置 `strategistPrivateKey` 后可用，否则返回 `ErrTransactorNotInitialized`

---

## API 接口

所有 API 统一前缀 `/api/v1`。CRUD 接口需要 JWT 认证。

### 公共接口（无需认证）

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | `/health` | 健康检查 |
| GET | `/ping` | 连通测试 |
| GET | `/codes` | 所有错误码 |
| POST | `/api/v1/auth/login` | MetaMask 登录 (获取 nonce) |
| POST | `/api/v1/auth/signature` | 钱包签名 → JWT |
| GET | `/swagger/index.html` | Swagger 文档 |

### 业务 CRUD

| 资源 | 方法 | 路径 |
|------|------|------|
| 市场 (markets) | GET/POST/PUT/DELETE | `/api/v1/markets[/:id]` |
| 预测 (predictions) | GET/POST/PUT/DELETE | `/api/v1/predictions[/:id]` |
| 策略 (strategies) | GET/POST/PUT/DELETE | `/api/v1/strategies[/:id]` |
| 交易 (trades) | GET/POST/PUT/DELETE | `/api/v1/trades[/:id]` |
| 金库快照 | GET/DELETE | `/api/v1/vault-snapshots[/:id]` |
| 系统日志 | GET/DELETE | `/api/v1/system-logs[/:id]` |
| 仪表盘 | GET | `/api/v1/stats/dashboard` |

### 任务手动触发

所有定时任务均可通过 HTTP 同步执行，方便调试：

**日常任务：**

| 方法 | 路径 | 请求体 | 说明 |
|------|------|--------|------|
| POST | `/api/v1/tasks/scan` | 无 | 触发市场扫描 |
| POST | `/api/v1/tasks/aggregate` | `{"marketId": 1}` | 数据聚合 |
| POST | `/api/v1/tasks/predict` | `{"marketId": 1}` | AI 预测 |
| POST | `/api/v1/tasks/strategy` | `{"predictionId":1,"marketId":1}` | 策略生成 |
| POST | `/api/v1/tasks/execute` | `{"marketId":1,"strategyId":1}` | 交易执行 |

**监控任务：**

| 方法 | 路径 | 说明 |
|------|------|------|
| POST | `/api/v1/tasks/positions` | 持仓监控 |
| POST | `/api/v1/tasks/vault-snapshot` | 金库快照 |
| POST | `/api/v1/tasks/health-check` | 健康检查 |
| POST | `/api/v1/tasks/settlement` | 结算检查 |

---

## 数据库模型

| 表 | 说明 | 关键字段 |
|----|------|---------|
| `markets` | 预测市场 | `polymarket_token_id`(唯一), `event_slug`, `question`, `price_threshold`, `current_yes/no_price`, `target_date`, `status`(active/resolved/expired), `resolution` |
| `predictions` | AI 预测 | `market_id`, `predicted_probability`, `confidence`, `direction`, `edge`, `raw_request`(JSON), `raw_response`(JSON), `data_snapshot`(JSON), `prompt_version` |
| `strategies` | 交易策略 | `prediction_id`, `market_id`, `action`(buy_yes/buy_no/skip), `position_size`, `entry_price`, `take_profit`, `stop_loss`, `kelly_fraction`, `edge`, `status`(pending/executing/active/skipped/closed/failed) |
| `trades` | 交易记录 | `strategy_id`, `market_id`, `polymarket_order_id`, `side`(yes/no), `amount`, `price`, `shares`, `status`(pending/filled/partial/cancelled/failed), `pnl`, `close_reason` |
| `vault_snapshots` | 金库快照 | `total_assets`, `share_price`, `tvl`, `depositor_count`, `deployed_amount` |
| `system_logs` | 系统日志 | `level`, `module`, `message`, `stack_trace` |
| `users` | 钱包用户 | `wallet_address`, `nonce`, `jwt_token` |

---

## 配置

主配置文件：`configs/be.yml`

```yaml
app:
  name: "be"
  env: "dev"                        # dev/prod/test
  enableMetrics: true               # Prometheus 指标
  enableTrace: false                # Jaeger 链路追踪
  cacheType: "redis"                # memory/redis

deepseek:
  apiKey: "${DEEPSEEK_API_KEY}"
  model: "deepseek-chat"
  temperature: 0.2
  useAgentPredictor: false          # true=多智能体, false=单模型

polymarket:
  clobApiUrl: "https://clob.polymarket.com"
  privateKey: "${POLYMARKET_PRIVATE_KEY}"
  gammaApiUrl: "https://gamma-api.polymarket.com"
  wsMarketUrl: "wss://ws-subscriptions-clob.polymarket.com/ws/market"

vault:
  rpcUrl: "https://polygon-rpc.com"
  contractAddress: ""               # PolyVault 合约地址
  strategistPrivateKey: "${STRATEGIST_PRIVATE_KEY}"

binance:
  baseUrl: "https://api.binance.com"
  klineLimit1h: 168                 # 1小时K线取7天
  klineLimit1d: 30                  # 日K线取30天

fearGreedIndex:
  url: "https://api.alternative.me/fng"
  limit: 7                          # 7日趋势

strategy:
  minEdge: 0.25                     # 最小期望收益
  minConfidence: 0.6                # 最小置信度
  maxPositionPct: 0.10              # 单笔最大仓位
  kellyMultiplier: 0.5              # Kelly 系数 (半凯利)

database:
  driver: mysql
  mysql:
    dsn: "user:pass@(127.0.0.1:3306)/polymarket_ai?parseTime=true&loc=Local&charset=utf8mb4"

redis:
  dsn: "default:pass@127.0.0.1:6379/0"
```

---

## 环境变量

| 变量 | 必填 | 说明 |
|------|------|------|
| `DEEPSEEK_API_KEY` | ✅ | DeepSeek API Key |
| `POLYMARKET_PRIVATE_KEY` | ✅ | Polymarket 签名钱包私钥 |
| `POLYMARKET_API_KEY` / `API_SECRET` / `PASSPHRASE` | ❌ | API 密钥认证 |
| `STRATEGIST_PRIVATE_KEY` | ❌* | PolyVault 策略角色私钥 |
| `GNEWS_API_KEY` | ❌ | GNews 新闻 API Key |
| `CRYPTOQUANT_API_KEY` | ❌ | CryptoQuant 链上数据 Key |

> *`STRATEGIST_PRIVATE_KEY` 仅在部署金库合约后必填，否则交易方法不可用。

---

## 启动

### 前置条件

- Go 1.22+
- MySQL 8.0+
- Redis 7+
- DeepSeek API Key
- Polygon RPC 节点（可选，金库功能需要）

### 1. 初始化数据库

```sql
CREATE DATABASE IF NOT EXISTS polymarket_ai
  CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci;
```

表结构需要手动创建（应用不会自动迁移），参考 `internal/model/` 下的模型定义。

### 2. 配置并运行

```bash
# 复制配置
cp configs/be.yml configs/be.yml

# 设置环境变量
export DEEPSEEK_API_KEY="sk-xxx"
export POLYMARKET_PRIVATE_KEY="0x..."
export STRATEGIST_PRIVATE_KEY="0x..."

# 编译运行
make run

# 测试页面: http://localhost:8080/test
# Swagger:  http://localhost:8080/swagger/index.html
```

### 3. 生成 API 文档

```bash
make docs
```

---

## Makefile 命令

| 命令 | 说明 |
|------|------|
| `make run` | 编译并启动 |
| `make build` | 构建 Linux amd64 二进制 |
| `make docs` | 生成 Swagger 文档 |
| `make test` | 运行单元测试 |
| `make cover` | 测试覆盖率 |
| `make ci-lint` | 代码规范检查 |
| `make run-nohup` | 后台运行 |
| `make run-docker` | Docker 部署 |
| `make deploy-k8s` | K8S 部署 |

---

## 监控与可观测性

| 功能 | 配置 | 端点 |
|------|------|------|
| Prometheus 指标 | `app.enableMetrics: true` | `GET /metrics` |
| Jaeger 链路追踪 | `app.enableTrace: true` | 需部署 Jaeger |
| 性能分析 | `app.enableHTTPProfile: true` | `GET /debug/pprof/` |
| 系统资源统计 | `app.enableStat: true` | 定期打印 CPU/内存 |
| 限流保护 | `app.enableLimit: true` | 自适应限流 |
| 熔断保护 | `app.enableCircuitBreaker: true` | 错误率熔断 |

---

## 容错设计

| 场景 | 处理方式 |
|------|---------|
| 新闻 API 失败 | 降级为空列表，不阻断流程 |
| AI 模型返回非 JSON | 自动提取 `{...}` 对象 + 重试一次 |
| WebSocket 断连 | 3 秒自动重连，无限重试 |
| WebSocket 重启重复收 `market_resolved` | `sync.Map` 去重 |
| 金库余额查询失败 | fallback 到 $100,000 默认值 |
| 数据库查询失败 | 返回错误，任务重试（Asynq 25次指数退避） |
| 链上交易失败 (`withdrawToStrategy`) | 策略标记 `failed`，不继续下单 |

---

## 常见问题

### Q: 启动报错 "panic: init config error"

检查 `configs/be.yml` 路径和 YAML 格式。

### Q: 如何只测试 AI 预测而不执行真实交易？

设置 `strategy.minConfidence: 1.5`（大于 1.0 使所有策略跳过），或直接调用 `/api/v1/tasks/predict`。

### Q: WebSocket 断连怎么办？

内置 3 秒自动重连，断连日志记录原因，无需手动干预。

### Q: 需要哪些外部 API？

| 服务 | 用途 | 必填 | 免费额度 |
|------|------|------|---------|
| DeepSeek | AI 预测 | ✅ | 注册赠送 |
| Binance | BTC 行情 | ✅ | 公开 API 免费 |
| Polymarket | 交易 + 市场数据 | ✅ | 只需 Gas |
| Polygon RPC | 链上金库 | ❌ | 公开 RPC 免费 |
| alternative.me | 恐惧贪婪指数 | ❌ | 公开 API 免费 |
| GNews | BTC 新闻 | ❌ | 免费版有速率限制 |
| CryptoQuant | 链上分析 | ❌ | 免费版数据有限 |
