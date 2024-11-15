## 分散式聊天匹配伺服器

這是一個使用 Go 語言開發的分散式聊天匹配伺服器，專注於探索分散式系統的設計與實作，同時打造一個具備高效能、易於擴展的架構。

---

### 系統架構概述

- **WebSocket 模組**：處理用戶端連線及即時訊息交換。
- **Pub/Sub 模組**：支援 Redis，透過介面實現抽象，方便切換到其他消息中介。
- **業務邏輯模組**：專注於匹配邏輯與使用者行為處理，與通訊層分離以提升開發效率。
- **可擴展性**：利用 Go 的介面與模組化設計，實現靈活的架構。

### 訊息流程圖

```
Client --> Node --> Broker --> Node --> Another Client
```

每個使用者發送訊息時，都會由節點先做處理再判斷發給 `Local Client` 或是 `Broker`，之後的 Node 會直接傳給 `Client`。

### 架構

```
├── cmd              <- 執行檔
├── consumer         <- 負責處理業務邏輯
│   ├── handler      <- 業務邏輯處理
│   └── request_msg  <- 定義訊息
├── controller       <- http控制器(使用者會從者裡建立ws連線)
└── websocket        <- 負責處理通訊連線
    └── pubsub       <- 背後的分散式處理
```

#### Websocket

```
└── websocket
    ├── client.go        <- 使用者的實例(當使用者建立連線後，會建立一個實例控制使用者)
    ├── hub.go           <- 類似中央控制，掌管每個使用者實例，並提供廣播或搜尋使用者功能、也負責Pub\Sub訊息處理
    ├── message.go       <- 定義消息
    └── pubsub           <- 透過Pub\Sub進行分散式的連線
        ├── interface.go
        └── redis.go
```

##### PubSubInterface
```
	Publish([]byte) error <- 負責將事件發布至其他Node
	Subscribe() []byte    <- 訂閱事件
```

Hub會訂閱(Subscribe) `PubSub` 取得其他 Node 的訊息，並判斷訊息內容是否在這台機器裡面，如果是則傳給`Client`，不是的話則省略。

#### 業務邏輯核心

```
├── consumer
│   ├── handler
│   │   ├── join_room_handler.go   <- 處理使用者進入房間、配對
│   │   ├── leave_room_handler.go  <- 處理使用者離開房間
│   │   └── on_message_handler.go  <- 處理使用者發送訊息
│   ├── message_consumer.go        <- 訊息會從這邊進入再分配至不同的handler
```

將業務邏輯分開，開發者可以專注在業務邏輯之中，當要傳送訊息給特定使用者時，只要呼叫 `SendMsgToClient` ，這個 Func 背後會自動判斷該傳給 `Broker` 或是 `Local User`。

### 下一步計劃
- 添加連線池與訊息緩存功能，進一步提升效能。
- 支援更多 Pub/Sub 後端（如 Kafka、NATS）。
- 模擬多節點場景，驗證系統在分散式環境中的穩定性。