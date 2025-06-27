# 🗨️ Realtime Chat Application

This is a lightweight, in-memory real-time chat system implemented in Go, supporting:

- Real-time messaging using WebSocket
- Offline message buffering
- Guaranteed message ordering
- Message acknowledgment
- REST API for chat history
- CLI client simulation


## Functional Requirements Coverage

| Feature                        | Description                                                                                       | Implementation Highlights                              |
|-------------------------------|---------------------------------------------------------------------------------------------------|--------------------------------------------------------|
| Real-Time Messaging           | User A can send a message to User B via WebSocket and vice versa                                  | WebSocket server with per-user connections             |
| Offline Message Buffering     | If a recipient is offline, messages are buffered in-memory                                        | `sync.Map` with per-user queues                        |
| Buffered Delivery             | On reconnection, buffered messages are delivered **in order**                                     | Messages dequeued and sent upon re-login               |
| Message Acknowledgments       | Every sent message receives an acknowledgment from the server                                     | `{"status":"ok"}` response on message receipt          |
| Chat History via REST API     | `/messages?user1=A&user2=B` returns chat history                                                  | HTTP GET handler accessing in-memory message store     |
| Guaranteed Ordering           | Messages are delivered in the order they were received                                            | Global message buffer preserves order                  |
| Thread Safety                 | Concurrent reads/writes are safe                                                                  | Protected via `sync.Mutex` and `sync.Map`             |
| Logging                       | Key events (connect, disconnect, buffer, delivery) are logged                                     | Console logs via `log.Printf`                         |

---

## Getting Started

### Build Project Binaries

```bash
# Build both client and server binaries
make build-local

# OR build individually
make build-app       # Server binary
make build-client    # Client binary
```
This will generate:
- dist/app
- dist/client

#### Run Server:
```bash
./dist/app
```
You should see logs such as:
```
2025/06/27 18:21:48 Chat service started on :8080
```

#### Run Clients
You can open two terminals to simulate User A and B.

Terminal 1 - User A:
```bash
./dist/client -user=A -recipient=B
```

Terminal 1 - User B:
```bash
./dist/client -user=B -recipient=A
```

#### Sending Messages

client expects message:
```
<message>
```

Example:
```bash
# From A
Hello B!
```

Output:
```
📨 Sent message to B: Hello, B!
✅ Acknowledgment: ok
```

#### REST API: Chat History

Fetch chat history using:
```bash
curl "http://localhost:8080/messages?user1=A&user2=B"
```

Output:
```json
[
  {"from":"A","to":"B","content":"Hello, B!","timestamp":...},
  {"from":"B","to":"A","content":"Hey A!","timestamp":...}
]
```

## Demo

https://github.com/user-attachments/assets/a8bfcbb1-94d8-4497-a0a0-1c134b4fbe83



