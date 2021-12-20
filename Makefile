build:
	go-bindata -pkg asset -o resource/asset/tpl.go resource/template/...
	go env -w GOOS="windows"
	go build -o bin/win-ctool.exe main.go
	go env -w GOOS="linux"
	go build -o bin/linux-ctool main.go
	go env -w GOOS="darwin"
	go build -o bin/mac-ctool main.go

push:
	git add .
	git commit -m "update"
	git push
	git tag -af v1.0.0 -m ""
	git push -f --tags