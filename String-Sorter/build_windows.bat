go version
go env
go get -u all
go-winres make
go build -o nss.exe -ldflags="-s -w"