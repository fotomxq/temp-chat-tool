set GOOS=windows
set GOARCH=amd64
go build -o ../../builds/server-api/server-api_win64.exe ../../server-api

set GOOS=linux
set GOARCH=amd64
go build -o ../../builds/server-api/server-api_linux64 ../../server-api