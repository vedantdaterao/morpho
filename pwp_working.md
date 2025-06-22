# peer wire protocol

> client ------pwp------> peer

### Establish connection

| client     | wire | peer       | msg ID |
| ---------- | ---- | ---------- | ------ |
| handshake  | ->   | .          | -      |
| .          | <-   | handshake  | -      |
| .          | <-   | bitfield   | 5      |
| .          | <-   | interested | 2      |
| .          | <-   | keep-alive | -      |
| keep-alive | ->   | .          | -      |

### Send `interested`

| client     | wire | peer | msg ID |
| ---------- | ---- | ---- | ------ |
| interested | ->   | .    | 2      |

### If peer sends `unchoke`

| client  | wire | peer    | msg ID |
| ------- | ---- | ------- | ------ |
|         |
| .       | <-   | unchoke | 1      |
| request | ->   | .       | 6      |
| .       | <-   | piece   | 7      |
| .       | <-   | piece   | 7      |
| .       | <-   | piece   | 7      |
| .       | <-   | piece   | 7      |
| -       | -    | -       | -      |
| have    | ->   | .       | 4      |

---

## code

```go
peer = Connect(peer)
// fills data [ bitfield, interested (0/1), choke (0/1) ]

for {
    msg = peer.ReadMsg()
    // process msg
}

peer.WriteMsg()

buffer = peer.Listner() // return message
peer.conn.Close()
```
