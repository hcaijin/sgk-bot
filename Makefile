SHELL:=/bin/bash
TARGET=sgkbot

all: win linux mac

win:
	GOOS=windows GOARCH=amd64 go build -o ./bin/${TARGET}.exe .
	GOOS=windows GOARCH=386 go build -o ./bin/${TARGET}-x86.exe .

linux:
	GOOS=linux GOARCH=amd64 go build -o ./bin/${TARGET}_${@} .

mac:
	GOOS=darwin GOARCH=amd64 go build -o ./bin/${TARGET}_${@} .

clean:
	rm -rf ./bin/${TARGET}*
