set version=v2.5.0
go generate -x
set GOOS=linux
go build -ldflags "-X main.Version=%version%"

set GOOS=windows
go build -ldflags "-X main.Version=%version%"