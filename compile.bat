go generate
set GOOS=linux
go build

set GOOS=windows
go build