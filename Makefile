
linux:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ./bin/game_server_linux main.go

win:
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o ./bin/game_server_win.exe main.go

mac:
	CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -o ./bin/game_server_mac main.go

vendor:
	export GOPROXY=https://goproxy.io,direct && go mod tidy && go mod vendor