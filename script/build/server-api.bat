set GOOS=windows
set GOARCH=amd64
go build -o ../../builds/server-api/server-api_win64.exe ../../src/server-api

set GOOS=linux
set GOARCH=amd64
go build -o ../../builds/server-api/server-api_linux64 ../../src/server-api