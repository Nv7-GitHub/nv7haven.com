cd "$(go env GOPATH)/src/github.com/Nv7-Github/chess"
GOOS=js GOARCH=wasm go build -o "$(go env GOPATH)/src/github.com/Nv7-Github/nv7haven.com/docs/chess.wasm" -ldflags="-s -w"
cd "$(go env GOPATH)/src/github.com/Nv7-Github/nv7haven.com"