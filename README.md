# syn-scanner
Simple SYN/ACK scanner on go.
Pet project for GoLang skills promotion.

![icon.svg](assets/icon.png)

## Usage
``` shell
go run . scan syn <target> [-t count]
```

### Arguments
```
  <target> string
        target for scanning
```

### Flags
```
  -t int
        number of threads(streams) when scanning (default 64)
```
