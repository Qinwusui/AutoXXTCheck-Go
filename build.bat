@echo off
@REM Linux amd64
go env -w GOOS=linux
go env -w GOARCH=amd64
go build -o AutoCheck-Linux-amd64 ./main
@REM Linux arm64
go env -w GOARCH=arm64
go build -o AutoCheck-Linux-arm64 ./main
@REM Windows arm64
go env -w GOOS=windows
go build -o AutoCheck-Windows-arm64.exe ./main
@REM Windows amd64
go env -w GOARCH=amd64
go build -o AutoCheck-Windows-amd64.exe ./main