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
