go get github.com/vugu/vugu
go install github.com/vugu/vugu/cmd/vugugen
COMMAND="$(go env GOPATH)/bin/vgrun -install-tools"
bash -c $COMMAND
echo "Got Tools"

COMMAND="$(go env GOPATH)/bin/vugugen -s"
bash -c $COMMAND
echo "Generated"
GOOS=js GOARCH=wasm go build -ldflags="-w -s" -o docs/main.wasm
echo "Built"
du -h docs/main.wasm

THING=$(cat docs/index.html | pcregrep -o1 '\?v=([0-9]+)')
THING=$((THING + 1))
sed -i '' "s/\?v=1/?v=${THING}/g" docs/index.html
echo "Incremented Version to ${THING}"
