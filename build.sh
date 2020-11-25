go generate
echo "Generated"
GOOS=js GOARCH=wasm go build -ldflags="-w -s" -o docs/main.wasm
echo "Built"
du -h docs/main.wasm