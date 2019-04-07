build:
	mkdir -p ./bin
	GOOS=linux  go build  -o ./bin/awswatcher cmd/main.go
	cd ./bin;zip awswatcher.zip awswatcher

