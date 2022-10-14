@echo off
@REM go env -w GOOS=linux
@REM go env -w GOARCH=amd64
@REM go build -o AutoCheck-linux-amd64 ./main

@REM go env -w GOARCH=arm64
@REM go build -o AutoCheck-linux-arm64 ./main
rm .\AutoCheck-*
go env -w GOOS=windows
go env -w GOARCH=amd64
go build -o AutoCheck-win-amd64 ./main

@REM go env -w GOARCH=arm64
@REM go build -o AutoCheck-win-arm64 ./main
