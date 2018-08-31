# go-portscanner

Check if a port is open using check-host.net

## API

```go
// Check if TCP port on a public ip address is open
CheckTCP(ip string, port uint16) (bool, error)
```
