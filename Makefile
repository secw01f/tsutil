build:
	go build -o bin/ .

compile:
	GOOS=linux GOARCH=amd64 go build -o bin/linux-amd64/ .
	tar -czf tsutil-linux-amd64.tar.gz -C bin/linux-amd64 .
	GOOS=linux GOARCH=arm64 go build -o bin/linux-arm64/ .
	tar -czf tsutil-linux-arm64.tar.gz -C bin/linux-arm64 .
	GOOS=darwin GOARCH=amd64 go build -o bin/macos-amd64/ .
	tar -czf tsutil-macos-amd64.tar.gz -C bin/macos-amd64 .
	GOOS=darwin GOARCH=arm64 go build -o bin/macos-arm64/ .
	tar -czf tsutil-macos-arm64.tar.gz -C bin/macos-arm64 .