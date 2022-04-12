#/bin/bash
GOOS=windows GOARCH=386 go build -ldflags -H=windowsgui -o isagraf_repair_tool.exe src/main.go src/utils.go