GOOS=windows GOARCH=amd64 go build -o build/cpq-windows-amd64.exe ./cmd/cpq
GOOS=windows GOARCH=386 go build -o build/cpq-windows-386.exe ./cmd/cpq
GOOS=linux GOARCH=amd64 go build -o build/cpq-linux-amd64 ./cmd/cpq
GOOS=linux GOARCH=386 go build -o build/cpq-linux-386 ./cmd/cpq
GOOS=darwin GOARCH=amd64 go build -o build/cpq-macos-amd64 ./cmd/cpq
GOOS=darwin GOARCH=386 go build -o build/cpq-macos-386 ./cmd/cpq
