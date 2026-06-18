-- ============================================================
-- Polymarket AI 自动化预测金库系统 - 数据库表结构
-- 引擎: InnoDB, 字符集: utf8mb4, 排序规则: utf8mb4_unicode_ci
-- ============================================================

CREATE DATABASE IF NOT EXISTS polymarket_ai
    DEFAULT CHARACTER SET utf8mb4
    DEFAULT COLLATE utf8mb4_unicode_ci;

USE polymarket_ai;

-- -----------------------------------------------------------
-- 1. markets — 市场信息表
--    存储从 Polymarket Gamma API 扫描并选中的交易市场
-- -----------------------------------------------------------
CREATE TABLE markets (
                         id                      INT             NOT NULL AUTO_INCREMENT  COMMENT '主键ID',
                         created_at              DATETIME        NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
                         updated_at              DATETIME        NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
                         deleted_at              DATETIME        DEFAULT NULL             COMMENT '软删除时间（有值表示已删除）',

                         polymarket_condition_id VARCHAR(128)    NOT NULL COMMENT 'Polymarket 条件ID（唯一标识一个预测市场）',
                         polymarket_token_id     VARCHAR(128)    NOT NULL COMMENT 'Polymarket Yes Token ID（AI默认下注方向）',
                         event_slug              VARCHAR(255)    NOT NULL COMMENT '事件唯一标识符（Event Slug）',
                         question                VARCHAR(500)    NOT NULL COMMENT '预测市场问题标题（如"BTC年底>$100K?"）',
                         price_threshold         INT             NOT NULL COMMENT '价格阈值（用于筛选市场的价格门槛, 单位: 百分点）',
                         scan_date               DATE            NOT NULL COMMENT '扫描日期（幂等键，每日最多一条记录）',
                         target_date             DATE            NOT NULL COMMENT '市场预测目标日期（即 Polymarket 的到期日）',
                         current_yes_price       DECIMAL(38,18)  NOT NULL COMMENT '当前 Yes 代币价格（即市场概率）',
                         current_no_price        DECIMAL(38,18)  NOT NULL COMMENT '当前 No 代币价格',
                         selected_at             DATETIME        NOT NULL COMMENT '被策略选中的时间戳',
                         status                  VARCHAR(16)     NOT NULL DEFAULT 'active' COMMENT '市场状态: active-活跃, resolved-已结算, expired-已过期',
                         resolution              VARCHAR(8)      DEFAULT NULL COMMENT '结算结果: yes-是, no-否（仅在 resolved 后有值）',

                         PRIMARY KEY (id),
                         UNIQUE KEY uk_polymarket_condition_id (polymarket_condition_id),
                         UNIQUE KEY uk_scan_date (scan_date),
                         INDEX idx_target_date (target_date),
                         INDEX idx_scan_date (scan_date)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='Polymarket 市场信息表';


-- -----------------------------------------------------------
-- 2. predictions — AI 预测记录表
--    存储 AI 模型对某个市场生成的概率预测及分析详情
-- -----------------------------------------------------------
CREATE TABLE predictions (
                             id                      INT             NOT NULL AUTO_INCREMENT  COMMENT '主键ID',
                             created_at              DATETIME        NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
                             updated_at              DATETIME        NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
                             deleted_at              DATETIME        DEFAULT NULL             COMMENT '软删除时间（有值表示已删除）',

                             market_id               INT             NOT NULL COMMENT '关联市场ID（逻辑外键→markets.id, 无物理约束）',

                             predicted_probability   DECIMAL(38,18)  NOT NULL COMMENT 'AI预测的概率值（范围: 0.01 ~ 0.99）',
                             confidence              DECIMAL(38,18)  NOT NULL COMMENT 'AI对自己预测的置信度（范围: 0 ~ 1）',
                             direction               VARCHAR(16)     NOT NULL COMMENT '预测方向: bullish-看涨, bearish-看跌, neutral-中性',

                             key_factors             JSON            NOT NULL COMMENT '关键技术因素列表（JSON字符串数组）',
                             risk_factors            JSON            NOT NULL COMMENT '风险因素列表（JSON字符串数组）',
                             technical_analysis      TEXT            NOT NULL COMMENT '技术面分析详情',
                             sentiment_analysis      TEXT            NOT NULL COMMENT '市场情绪分析详情',
                             news_impact             TEXT            NOT NULL COMMENT '新闻影响分析详情',
                             onchain_analysis        TEXT            NOT NULL COMMENT '链上数据分析详情',
                             reasoning               TEXT            NOT NULL COMMENT 'AI推理过程的完整描述',

                             recommended_action      VARCHAR(16)     NOT NULL COMMENT '推荐操作: buy_yes-买入Yes, buy_no-买入No, skip-跳过',
                             market_probability      DECIMAL(38,18)  NOT NULL COMMENT '预测时刻的 Polymarket 市场概率（Yes价格）',
                             edge                    DECIMAL(38,18)  NOT NULL COMMENT 'AI概率与市场概率的差值 = predicted_probability - market_probability',

                             model_version           VARCHAR(32)     NOT NULL COMMENT 'AI模型版本号',
                             prompt_version          VARCHAR(16)     NOT NULL COMMENT '提示词模板版本号',
                             seed                    INT             NOT NULL COMMENT '模型推理使用的随机种子（用于复现）',

                             raw_request             JSON            NOT NULL COMMENT '发送给AI模型的原始请求体',
                             raw_response            JSON            NOT NULL COMMENT 'AI模型返回的原始响应体',
                             data_snapshot           JSON            NOT NULL COMMENT '预测时刻的各类数据快照',

                             tokens_used             INT             NOT NULL DEFAULT 0 COMMENT '本次预测消耗的Token数量',
                             latency_ms              INT             NOT NULL DEFAULT 0 COMMENT 'AI模型推理耗时（毫秒）',

                             PRIMARY KEY (id),
                             INDEX idx_market_id (market_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='AI预测记录表';


-- -----------------------------------------------------------
-- 3. strategies — 交易策略表
--    存储策略生成器根据AI预测产出的交易策略
-- -----------------------------------------------------------
CREATE TABLE strategies (
                            id                      INT             NOT NULL AUTO_INCREMENT  COMMENT '主键ID',
                            created_at              DATETIME        NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
                            updated_at              DATETIME        NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
                            deleted_at              DATETIME        DEFAULT NULL             COMMENT '软删除时间（有值表示已删除）',

                            prediction_id           INT             NOT NULL COMMENT '关联预测记录ID（逻辑外键→predictions.id, 无物理约束）',
                            market_id               INT             NOT NULL COMMENT '关联市场ID（逻辑外键→markets.id, 无物理约束）',

                            action                  VARCHAR(16)     NOT NULL COMMENT '策略动作: buy_yes-买入Yes, buy_no-买入No, skip-跳过',
                            side                    VARCHAR(8)      DEFAULT NULL COMMENT '交易方向: yes-Yes方, no-No方（skip时为NULL）',
                            position_size           DECIMAL(38,18)  NOT NULL COMMENT '仓位大小（USDC金额, skip时为0）',
                            entry_price             DECIMAL(38,18)  NOT NULL COMMENT '入场价格（skip时为0）',
                            take_profit             DECIMAL(38,18)  NOT NULL COMMENT '止盈价格（skip时为0）',
                            stop_loss               DECIMAL(38,18)  NOT NULL COMMENT '止损价格（skip时为0）',
                            kelly_fraction          DECIMAL(38,18)  NOT NULL COMMENT '凯利公式建议仓位比例（skip时为0）',
                            edge                    DECIMAL(38,18)  NOT NULL COMMENT '策略采用的价差（正负值，有符号）',
                            skip_reason             TEXT            NOT NULL COMMENT '跳过交易的详细原因（未跳过时为空字符串）',
                            status                  VARCHAR(16)     NOT NULL DEFAULT 'pending' COMMENT '策略状态: skipped-跳过, pending-待执行, executing-执行中, active-已开仓, closed-已平仓, failed-执行失败',
                            executed_at             DATETIME        DEFAULT NULL COMMENT '策略实际执行时间',

                            PRIMARY KEY (id),
                            INDEX idx_prediction_id (prediction_id),
                            INDEX idx_market_id (market_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='交易策略表';


-- -----------------------------------------------------------
-- 4. trades — 交易执行记录表
--    存储策略执行后实际在 Polymarket CLOB 上成交的订单
-- -----------------------------------------------------------
CREATE TABLE trades (
                        id                      INT             NOT NULL AUTO_INCREMENT  COMMENT '主键ID',
                        created_at              DATETIME        NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
                        updated_at              DATETIME        NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
                        deleted_at              DATETIME        DEFAULT NULL             COMMENT '软删除时间（有值表示已删除）',

                        strategy_id             INT             NOT NULL COMMENT '关联策略ID（逻辑外键→strategies.id, 无物理约束）',
                        market_id               INT             NOT NULL COMMENT '关联市场ID（逻辑外键→markets.id, 无物理约束）',

                        polymarket_order_id     VARCHAR(128)    NOT NULL COMMENT 'Polymarket CLOB 订单ID（全局唯一）',
                        side                    VARCHAR(8)      NOT NULL COMMENT '交易方向: yes-Yes方, no-No方',
                        action                  VARCHAR(8)      NOT NULL COMMENT '操作类型: buy-买入建仓, sell-卖出平仓',
                        amount                  DECIMAL(38,18)  NOT NULL COMMENT '交易金额（买入时为USDC支出, 卖出时为USDC收入）',
                        price                   DECIMAL(38,18)  NOT NULL COMMENT '成交单价',
                        shares                  DECIMAL(38,18)  NOT NULL COMMENT '成交份额数量',
                        status                  VARCHAR(16)     NOT NULL COMMENT '订单状态: pending-待成交, filled-已成交, partial-部分成交, cancelled-已取消, failed-失败',
                        fee                     DECIMAL(38,18)  NOT NULL DEFAULT 0 COMMENT '交易手续费（USDC）',
                        pnl                     DECIMAL(38,18)  DEFAULT NULL COMMENT '盈亏金额（平仓时有值, USDC）',
                        close_reason            VARCHAR(32)     DEFAULT NULL COMMENT '平仓原因: take_profit-止盈, stop_loss-止损, pre_resolution-到期前, manual-手动',
                        filled_at               DATETIME        DEFAULT NULL COMMENT '订单成交时间',
                        closed_at               DATETIME        DEFAULT NULL COMMENT '平仓时间',

                        PRIMARY KEY (id),
                        UNIQUE KEY uk_polymarket_order_id (polymarket_order_id),
                        INDEX idx_strategy_id (strategy_id),
                        INDEX idx_market_id (market_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='交易执行记录表';


-- -----------------------------------------------------------
-- 5. vault_snapshots — 金库快照表
--    定期记录 PolyVault 智能合约的状态快照
-- -----------------------------------------------------------
CREATE TABLE vault_snapshots (
                                 id                      INT             NOT NULL AUTO_INCREMENT  COMMENT '主键ID',
                                 created_at              DATETIME        NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
                                 updated_at              DATETIME        NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
                                 deleted_at              DATETIME        DEFAULT NULL             COMMENT '软删除时间（有值表示已删除）',

                                 total_assets            DECIMAL(38,18)  NOT NULL COMMENT '金库总资产（含链下策略债务, USDC）',
                                 share_price             DECIMAL(38,18)  NOT NULL COMMENT '当前份额价格（USDC/份额）',
                                 tvl                     DECIMAL(38,18)  NOT NULL COMMENT '锁定总价值（Total Value Locked, USDC）',
                                 depositor_count         INT             NOT NULL DEFAULT 0 COMMENT '存款人数量',
                                 deployed_amount         DECIMAL(38,18)  NOT NULL DEFAULT 0 COMMENT '已部署到链下策略的资金量（USDC）',
                                 snapshot_at             DATETIME        NOT NULL COMMENT '快照时间戳',

                                 PRIMARY KEY (id),
                                 INDEX idx_snapshot_at (snapshot_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='PolyVault金库快照表';


-- -----------------------------------------------------------
-- 6. system_logs — 系统日志表
--    存储系统运行日志，同时用作暂停/恢复标志位存储
-- -----------------------------------------------------------
CREATE TABLE system_logs (
                             id                      INT             NOT NULL AUTO_INCREMENT  COMMENT '主键ID',
                             created_at              DATETIME        NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
                             updated_at              DATETIME        NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
                             deleted_at              DATETIME        DEFAULT NULL             COMMENT '软删除时间（有值表示已删除）',

                             level                   VARCHAR(16)     NOT NULL COMMENT '日志级别: INFO, WARNING, ERROR, DEBUG',
                             source                  VARCHAR(64)     NOT NULL COMMENT '日志来源（模块名, 如 scanner, predictor, executor）',
                             message                 TEXT            NOT NULL COMMENT '日志消息内容',
                             context                 JSON            NOT NULL COMMENT '日志上下文信息（JSON格式, 包含额外结构化数据）',
                             trace_id                VARCHAR(64)     DEFAULT NULL COMMENT '链路追踪ID（用于关联同一请求链中的多条日志）',

                             PRIMARY KEY (id),
                             INDEX idx_source (source),
                             INDEX idx_trace_id (trace_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='系统日志表（同时用作暂停恢复标志位存储）';

-- -----------------------------------------------------------
-- 7. users — 用户表
--    存储通过 MetaMask 钱包登录的用户信息
-- -----------------------------------------------------------
CREATE TABLE users (
                       id                      INT             NOT NULL AUTO_INCREMENT  COMMENT '主键ID',
                       created_at              DATETIME        NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
                       updated_at              DATETIME        NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
                       deleted_at              DATETIME        DEFAULT NULL             COMMENT '软删除时间（有值表示已删除）',

                       wallet_address          VARCHAR(256)    NOT NULL COMMENT 'MetaMask 钱包地址',
                       nickname                VARCHAR(128)    DEFAULT NULL COMMENT '昵称',
                       avatar                  VARCHAR(512)    DEFAULT NULL COMMENT '头像 URL',
                       last_login_at           DATETIME        DEFAULT NULL COMMENT '最后登录时间',

                       PRIMARY KEY (id),
                       UNIQUE INDEX idx_wallet_address (wallet_address)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='用户表（MetaMask 钱包登录用户）';
