#!/bin/sh

# ARM
env GOOS=linux GOARCH=arm GOARM=5 go build -o client_linux_arm5 github.com/messyidea/sufttun/client
env GOOS=linux GOARCH=arm GOARM=6 go build -o client_linux_arm6 github.com/messyidea/sufttun/client
env GOOS=linux GOARCH=arm GOARM=7 go build -o client_linux_arm7 github.com/messyidea/sufttun/client
env GOOS=linux GOARCH=arm GOARM=5 go build -o server_linux_arm5 github.com/messyidea/sufttun/server
env GOOS=linux GOARCH=arm GOARM=6 go build -o server_linux_arm6 github.com/messyidea/sufttun/server
env GOOS=linux GOARCH=arm GOARM=7 go build -o server_linux_arm7 github.com/messyidea/sufttun/server
tar -zcf kcptun-linux-arm567.tar.gz client_linux_arm* server_linux_arm*
md5 kcptun-linux-arm567.tar.gz

# AMD64
env GOOS=linux GOARCH=amd64 go build -o client_linux_amd64 github.com/messyidea/sufttun/client
env GOOS=linux GOARCH=amd64 go build -o server_linux_amd64 github.com/messyidea/sufttun/server
tar -zcf kcptun-linux-amd64.tar.gz client_linux_amd64 server_linux_amd64
md5 kcptun-linux-amd64.tar.gz
env GOOS=darwin GOARCH=amd64 go build -o client_darwin_amd64 github.com/messyidea/sufttun/client
env GOOS=darwin GOARCH=amd64 go build -o server_darwin_amd64 github.com/messyidea/sufttun/server
tar -zcf kcptun-darwin-amd64.tar.gz client_darwin_amd64 server_darwin_amd64
md5 kcptun-darwin-amd64.tar.gz
env GOOS=windows GOARCH=amd64 go build -o client_windows_amd64.exe github.com/messyidea/sufttun/client
env GOOS=windows GOARCH=amd64 go build -o server_windows_amd64.exe github.com/messyidea/sufttun/server
tar -zcf kcptun-windows-amd64.tar.gz client_windows_amd64.exe server_windows_amd64.exe
md5 kcptun-windows-amd64.tar.gz

# 386
env GOOS=linux GOARCH=386 go build -o client_linux_386 github.com/messyidea/sufttun/client
env GOOS=linux GOARCH=386 go build -o server_linux_386 github.com/messyidea/sufttun/server
tar -zcf kcptun-linux-386.tar.gz client_linux_386 server_linux_386
md5 kcptun-linux-386.tar.gz
env GOOS=darwin GOARCH=386 go build -o client_darwin_386 github.com/messyidea/sufttun/client
env GOOS=darwin GOARCH=386 go build -o server_darwin_386 github.com/messyidea/sufttun/server
tar -zcf kcptun-darwin-386.tar.gz client_darwin_386 server_darwin_386
md5 kcptun-darwin-386.tar.gz
env GOOS=windows GOARCH=386 go build -o client_windows_386.exe github.com/messyidea/sufttun/client
env GOOS=windows GOARCH=386 go build -o server_windows_386.exe github.com/messyidea/sufttun/server
tar -zcf kcptun-windows-386.tar.gz client_windows_386.exe server_windows_386.exe
md5 kcptun-windows-386.tar.gz
