# websocket_benchmark
Websocket Benchmarking


Build steps:
1. `git clone https://github.com/FarhadF/websocket_benchmark/`
2. `cd websocket_benchmark`
3. `go build -o websocket_benchmark main.go`

```
Usage:
  websocket_benchmark [flags]

Flags:
  -a, --address string   Websocket endpoint address (default "localhost:8080")
  -d, --duration int     Runtime Duration in seconds (default 60)
  -h, --help             help for websocket_benchmark
  -i, --interval int     Message sending Interval in seconds (default 1)
  -m, --message string   Message to send (default "{"message":"sample message"}")
  -p, --path string      Websocket endpoint relative path (default "/echo")
  -s, --sockets int      Number of Sockets to use (default 500)
  -t, --timeout int      Websocket handshake timeout in seconds (default 15)
  -v, --version          Prints version info
  ```
