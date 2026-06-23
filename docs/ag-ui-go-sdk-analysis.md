# AG-UI Go SDK 代码分析文档

> 仓库地址: https://github.com/ag-ui-protocol/ag-ui/tree/main/sdks/community/go
> 模块路径: `github.com/ag-ui-protocol/ag-ui/sdks/community/go`

---

## 1. 整体架构概览

AG-UI Go SDK 实现了 AG-UI 协议的 Go 语言版本，提供了一套完整的 Agent-UI 交互基础设施，包括：

- **核心类型系统** — Message、Tool、Interrupt、RunAgentInput 等协议类型
- **事件驱动模型** — 完整的事件类型定义和生命周期管理
- **SSE 客户端** — 基于 Server-Sent Events 的流式通信客户端
- **编解码框架** — 支持多种编码格式（JSON/SSE）的编解码器体系
- **示例项目** — 包含完整的 Client（TUI）和 Server（Fiber HTTP）示例

**依赖项：**
- `github.com/google/uuid` — UUID 生成
- `github.com/sirupsen/logrus` — 结构化日志

---

## 2. 目录结构

```
sdks/community/go/
├── go.mod / go.sum
├── pkg/
│   ├── core/
│   │   ├── types/                    # 核心协议类型
│   │   │   ├── types.go              # Role, Message, Tool, RunAgentInput 等
│   │   │   └── message_helpers.go    # Message 辅助方法
│   │   └── events/                   # 事件系统
│   │       ├── events.go             # Event接口/BaseEvent/事件类型常量/序列验证
│   │       ├── message_events.go     # 文本消息事件 (start/content/end/chunk)
│   │       ├── run_events.go         # 运行/步骤事件 (started/finished/error/step)
│   │       ├── state_events.go       # 状态事件 (snapshot/delta) + JSON Patch
│   │       ├── activity_events.go    # 活动事件 (snapshot/delta)
│   │       ├── reasoning_events.go   # 推理事件 (reasoning message lifecycle)
│   │       ├── custom_events.go      # 自定义事件 + Raw事件
│   │       ├── decoder.go            # SSE -> Go 结构体 事件解码器
│   │       └── id_utils.go           # ID 生成器 (UUID / Timestamp)
│   ├── client/
│   │   └── sse/
│   │       └── client.go             # SSE 客户端 (HTTP 连接、流式读帧)
│   ├── encoding/                     # 编解码框架
│   │   ├── interface.go              # 编解码器接口 (Encoder/Decoder/StreamCodec)
│   │   ├── encoder/encoder.go        # 编码器实现
│   │   ├── sse/writer.go             # SSE 写入器
│   │   ├── json/                     # JSON 编解码实现
│   │   └── negotiation/              # 内容协商 (Accept/Content-Type)
│   └── errors/
│       ├── error_types.go            # 错误类型定义
│       └── error_utils.go            # 错误工具
└── example/
    ├── client/                       # TUI 客户端示例 (Bubble Tea)
    │   ├── cmd/main.go
    │   └── internal/{agent,event,message,ui}/
    └── server/                       # Fiber HTTP 服务端示例
        ├── cmd/main.go
        └── internal/{agentic,config,mcp,routes}/
```

---

## 3. 核心类型系统 (`pkg/core/types`)

### 3.1 Role（消息角色）

```go
type Role string
const (
    RoleDeveloper  Role = "developer"
    RoleSystem     Role = "system"
    RoleAssistant  Role = "assistant"
    RoleUser       Role = "user"
    RoleTool       Role = "tool"
    RoleActivity   Role = "activity"
    RoleReasoning  Role = "reasoning"
)
```

### 3.2 Message（消息结构）

```go
type Message struct {
    ID               string      `json:"id"`
    Role             Role        `json:"role"`
    Content          any         `json:"content,omitempty"`   // string | []InputContent | map
    Name             string      `json:"name,omitempty"`
    EncryptedContent string      `json:"encryptedContent,omitempty"`
    EncryptedValue   string      `json:"encryptedValue,omitempty"`
    ToolCalls        []ToolCall  `json:"toolCalls,omitempty"`
    ToolCallID       string      `json:"toolCallId,omitempty"`
    Error            string      `json:"error,omitempty"`
    ActivityType     string      `json:"activityType,omitempty"`
}
```

关键点：
- `Content` 字段为 `any` 类型，根据 Role 不同可以是 `string`（通用）、`[]InputContent`（User 多模态）、`map`（Activity）
- 提供了 `ContentString()`、`ContentInputContents()`、`ContentActivity()` 辅助方法
- 支持 `snake_case` 兼容反序列化（通过自定义 `UnmarshalJSON`）
- 包含加密内容字段用于状态连续性

### 3.3 多模态输入 (`InputContent`)

```go
type InputContent struct {
    Type     string               `json:"type"`     // text/binary/image/audio/video/document
    Text     string               `json:"text,omitempty"`
    Source   *InputContentSource  `json:"source,omitempty"` // data/url
    Data     string               `json:"data,omitempty"`   // base64
    // ... 其他字段
}
```

支持类型：`text`, `binary`, `image`, `audio`, `video`, `document`

### 3.4 Tool & ToolCall

```go
type Tool struct {
    Name        string `json:"name"`
    Description string `json:"description"`
    Parameters  any    `json:"parameters"` // JSON Schema
}

type ToolCall struct {
    ID       string       `json:"id"`
    Type     string       `json:"type"` // "function"
    Function FunctionCall `json:"function"`
}
```

### 3.5 Interrupt（中断/暂停机制）

```go
type Interrupt struct {
    ID             string         `json:"id"`
    Reason         string         `json:"reason"`
    Message        string         `json:"message,omitempty"`
    ResponseSchema map[string]any `json:"responseSchema,omitempty"`
    ExpiresAt      string         `json:"expiresAt,omitempty"`
    // ...
}
```

支持两种退出状态：
- `RunFinishedOutcomeTypeSuccess` — 正常运行结束
- `RunFinishedOutcomeTypeInterrupt` — 因中断暂停，等待用户输入后恢复

### 3.6 RunAgentInput（Agent 运行输入）

```go
type RunAgentInput struct {
    ThreadID       string         `json:"threadId"`
    RunID          string         `json:"runId"`
    ParentRunID    *string        `json:"parentRunId,omitempty"`
    State          any            `json:"state"`
    Messages       []Message      `json:"messages"`
    Tools          []Tool         `json:"tools"`
    Context        []Context      `json:"context"`
    ForwardedProps any            `json:"forwardedProps"`
    Resume         []ResumeEntry  `json:"resume,omitempty"` // 恢复暂停的 run
}
```

---

## 4. 事件系统 (`pkg/core/events`)

### 4.1 事件接口

```go
type Event interface {
    Type() EventType
    Timestamp() *int64
    SetTimestamp(timestamp int64)
    ThreadID() string
    RunID() string
    Validate() error
    ToJSON() ([]byte, error)
    GetBaseEvent() *BaseEvent
}
```

所有事件共享的基类：

```go
type BaseEvent struct {
    EventType   EventType `json:"type"`
    TimestampMs *int64    `json:"timestamp,omitempty"`
    RawEvent    any       `json:"rawEvent,omitempty"`
}
```

### 4.2 完整事件类型清单

| 分类 | 事件类型 | 说明 |
|------|---------|------|
| **Run** | `RUN_STARTED` | 运行开始，携带 threadId + runId |
| | `RUN_FINISHED` | 运行结束，可携带 result / outcome (success/interrupt) |
| | `RUN_ERROR` | 运行出错，携带 code + message |
| **Step** | `STEP_STARTED` | 步骤开始 |
| | `STEP_FINISHED` | 步骤结束 |
| **Text Message** | `TEXT_MESSAGE_START` | 文本消息开始 (messageId + role) |
| | `TEXT_MESSAGE_CONTENT` | 文本消息内容块 (messageId + delta) |
| | `TEXT_MESSAGE_END` | 文本消息结束 |
| | `TEXT_MESSAGE_CHUNK` | 文本消息便捷块 (可携带 messageId/role/delta) |
| **Tool Call** | `TOOL_CALL_START` | 工具调用开始 |
| | `TOOL_CALL_ARGS` | 工具调用参数 |
| | `TOOL_CALL_END` | 工具调用结束 |
| | `TOOL_CALL_CHUNK` | 工具调用便捷块 |
| | `TOOL_CALL_RESULT` | 工具调用结果 |
| **State** | `STATE_SNAPSHOT` | 状态完整快照 |
| | `STATE_DELTA` | 状态增量更新 (JSON Patch RFC 6902) |
| **Messages** | `MESSAGES_SNAPSHOT` | 消息列表快照 |
| **Activity** | `ACTIVITY_SNAPSHOT` | 活动消息快照 |
| | `ACTIVITY_DELTA` | 活动增量更新 |
| **Reasoning** | `REASONING_START` / `REASONING_END` | 推理阶段开始/结束 |
| | `REASONING_MESSAGE_START/CONTENT/END/CHUNK` | 推理消息生命周期 |
| | `REASONING_ENCRYPTED_VALUE` | 加密推理值 |
| **Thinking（废弃）** | `THINKING_START/END/...` | 旧版思考事件，已废弃 |
| **通用** | `RAW` | 原始透传事件 |
| | `CUSTOM` | 自定义事件 |

### 4.3 事件序列验证

`ValidateSequence(events []Event)` 函数提供事件序列验证，确保：
- Run 不能重复启动
- 不能在未启动时结束
- Step/Message/ToolCall 的 start-end 配对正确
- 不能在已结束的 run 下重启

### 4.4 事件解码器

`EventDecoder` 提供 SSE 事件名 -> Go 结构体的映射，支持所有已知事件类型的自动解码。

### 4.5 ID 生成器

提供两种实现：
- `DefaultIDGenerator` — 基于 UUID v4，格式 `{type}-{uuid}`（如 `run-xxx`, `msg-xxx`）
- `TimestampIDGenerator` — 基于时间戳+UUID 前缀，格式 `{type}-{timestamp}-{shortUUID}`

全局便捷函数：`GenerateRunID()`, `GenerateMessageID()`, `GenerateToolCallID()`, `GenerateThreadID()`, `GenerateStepID()`

---

## 5. SSE 客户端 (`pkg/client/sse`)

### 5.1 Client 结构

```go
type Config struct {
    Endpoint       string        // AG-UI 服务端地址
    APIKey         string        // API Key
    AuthHeader     string        // 自定义认证头（默认 Authorization）
    AuthScheme     string        // 认证方案（默认 Bearer）
    ConnectTimeout time.Duration // 连接超时（默认 30s）
    ReadTimeout    time.Duration // 读取超时（默认 5min）
    BufferSize     int           // 帧缓冲区大小（默认 100）
    Logger         *logrus.Logger
}
```

### 5.2 核心用法

```go
// 创建客户端
client := sse.NewClient(sse.Config{
    Endpoint: "http://localhost:8080/agentic",
    APIKey:   "your-key",
})

// 发起流式请求
frames, errors, err := client.Stream(sse.StreamOptions{
    Context: ctx,
    Payload: types.RunAgentInput{
        ThreadID: "thread-xxx",
        RunID:    "run-xxx",
        Messages: messages,
        // ...
    },
})

// 读取 SSE 帧
for {
    select {
    case frame := <-frames:
        // frame.Data 是原始字节，需通过 EventDecoder 解码
    case err := <-errors:
        // 处理错误
    case <-ctx.Done():
        return
    }
}
```

### 5.3 连接特性

- 使用 HTTP POST 发送 `application/json` 请求体
- 设置 `Accept: text/event-stream` 请求头
- 自动处理 `Authorization` 头（Bearer 方案）
- 支持自定义认证头和方案
- 异步读取流数据，通过 channel 输出帧
- 支持 context 取消和读超时

---

## 6. 编解码框架 (`pkg/encoding`)

### 6.1 接口层次

采用**接口隔离原则**，定义了细粒度的接口：

```
Encoder ────────── 编码单/多事件
Decoder ────────── 解码单/多事件

StreamEncoder ──── 流式编码（chan Event -> io.Writer）
StreamDecoder ──── 流式解码（io.Reader -> chan Event）

Codec     ─────── Encoder + Decoder + ContentTypeProvider
StreamCodec ───── Codec + 流式操作
FullStreamCodec ─ Codec + StreamCodec + 会话管理

CodecFactory ───── 创建 Codec 的工厂
StreamCodecFactory  创建 StreamCodec 的工厂
ContentNegotiator   内容协商（Accept 头解析）
```

### 6.2 SSE Writer (`pkg/encoding/sse`)

SSE 写入器将 Go 事件转为标准 SSE 格式：

```go
sseWriter := sse.NewSSEWriter()

// 写入事件（自动添加 "data: " 前缀）
sseWriter.WriteEvent(ctx, writer, event)

// 写入原始字节
sseWriter.WriteBytes(ctx, writer, data)

// 写入注释
sseWriter.WriteComment(ctx, writer, "some comment")
```

---

## 7. 示例项目

### 7.1 Server 端 (`example/server`)

技术栈：**Fiber v3（HTTP）+ langchaingo（LLM）+ MCP**

```
入口: cmd/main.go
├── config/       — 配置管理（环境变量 + 命令行参数）
├── routes/       — 路由定义
│   └── POST /agentic → SSE 流式响应
├── agentic/      — Agent 逻辑
│   ├── agentic.go    — LLM 调用（Claude 3 Haiku + langchaingo Agent）
│   ├── handler.go    — 回调处理器（将 LLM 输出转为事件流）
│   └── data/reminder.md — 提示词模板
└── mcp/          — MCP Server 适配器
```

**关键流程：**

1. 客户端 POST JSON 到 `/agentic`
2. Server 解析请求，发送 `RUN_STARTED` 事件
3. 调用 `agentic.ProcessInput()` 处理用户输入
4. 内部调用 Claude LLM，通过回调处理器将输出转为事件流
5. 发送 `RUN_FINISHED` 事件

### 7.2 Client 端 (`example/client`)

技术栈：**Bubble Tea（TUI）+ SSE Client**

```
入口: cmd/main.go
├── agent/chat.go  — 使用 SSE Client 连接服务端
├── event/         — 事件处理（解码器）
├── message/       — 消息处理
└── ui/            — Bubble Tea TUI 界面
```

**关键流程：**

1. 启动 Bubble Tea TUI
2. 用户输入通过 channel 传递给 agent 协程
3. agent 使用 `sse.Client` 连接服务端，流式读取响应
4. 解码 SSE 帧为事件，更新 TUI 界面

---

## 8. JSON Patch 支持

SDK 内置 JSON Patch (RFC 6902) 支持：

```go
type JSONPatchOperation struct {
    Op    string `json:"op"`    // add/remove/replace/move/copy/test
    Path  string `json:"path"`  // JSON Pointer
    Value any    `json:"value,omitempty"`
    From  string `json:"from,omitempty"`
}
```

用于 `STATE_DELTA` 和 `ACTIVITY_DELTA` 事件的增量更新。

---

## 9. 使用流程图

```
┌──────────┐     POST /agentic      ┌──────────┐
│  Client  │ ──────────────────────> │  Server  │
│ (SSE)    │                         │ (Fiber)  │
│          │ <── SSE Stream ────── │          │
└──────────┘   事件流                └──────────┘
                                         │
                                    ┌────▼─────┐
                                    │  LLM     │
                                    │ (Claude) │
                                    └──────────┘

事件流序列：
  RUN_STARTED → TEXT_MESSAGE_START → TEXT_MESSAGE_CONTENT* → TEXT_MESSAGE_END → RUN_FINISHED

可选：
  - TOOL_CALL_START → TOOL_CALL_ARGS → TOOL_CALL_END → TOOL_CALL_RESULT
  - STATE_SNAPSHOT / STATE_DELTA
  - REASONING_START → REASONING_MESSAGE_* → REASONING_END
```

---

## 10. 关键设计特点

1. **接口隔离原则** — 编码框架定义了大量细粒度接口，通过组合实现复杂功能
2. **snake_case 兼容** — 所有 JSON 反序列化支持 camelCase 和 snake_case
3. **构建器模式** — 事件使用 Functional Options 模式（如 `WithRole()`, `WithAutoRunID()`）
4. **事件序列验证** — 内置 `ValidateSequence` 确保事件生命周期正确
5. **可扩展性** — 支持 `CUSTOM` 和 `RAW` 事件类型，业务可自定义扩展
6. **对象池** — `encoding/pool.go` 中使用了对象池优化性能
7. **中断/恢复机制** — 通过 `Interrupt` 和 `ResumeEntry` 支持人机交互暂停恢复
