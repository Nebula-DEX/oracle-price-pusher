# oracle-price-pusher

Binary pushes oracles from the [twelvedata](https://twelvedata.com) api to the vega protocol chain.

# Usage

```shell
# build the binary
go build -o oracle-price-pusher ./

# run the process
./oracle-price-pusher --config config.toml
```