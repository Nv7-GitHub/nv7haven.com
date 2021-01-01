go generate
echo "Generated"
GOOS=js GOARCH=wasm go build -ldflags="-w -s" -o docs/main.wasm
echo "Built"
du -h docs/main.wasm

THING=$(cat docs/index.html | pcregrep -o1 '\?v=(\d+)')
NUM=$((THING + 1))
sed -i '' "s/?v=${THING}/?v=${NUM}/g" docs/index.html
echo "Incremented Version to ${NUM}"