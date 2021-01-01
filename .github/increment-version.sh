go get github.com/vugu/vugu
go install github.com/vugu/vugu/cmd/vugugen
echo "Got Tools"

COMMAND="$(go env GOPATH)/bin/vugugen -s"
bash -c $COMMAND
bash -c $COMMAND
bash -c $COMMAND
bash -c $COMMAND
echo "Generated"
GOOS=js GOARCH=wasm go build -ldflags="-w -s" -o docs/main.wasm
echo "Built"
du -h docs/main.wasm

THING=$(cat docs/index.html | pcregrep -o1 '\?v=([0-9]+)')
THING=$((THING + 1))
sed -i '' "s/\?v=1/?v=${THING}/g" docs/index.html
echo "Incremented Version to ${THING}"
