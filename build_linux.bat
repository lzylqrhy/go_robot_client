set GOARCH=amd64
set GOOS=linux
go build -o bin/robots_linux ./
xcopy configs bin\configs /s /y
xcopy README.md bin\ /y