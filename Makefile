build:

codemsg:
	go build ./cmd/h2o/h2o.go
	- rm ./testdata/errno_string.go
	- rm ./testdata/errno_codemsg.go
	- rm ./testdata/errno_grpc_status.go.go
	./h2o codemsg --code-msg --linecomment --string --string-method String2 --grpc --type ErrNo ./testdata/err.go

protoc:
	go build ./cmd/h2o/h2o.go
	rm -rf ./usertoken
	./h2o pb -f ./testdata/usertoken.yaml
	cat ./usertoken/usertoken.proto

# 忽略
update:
	make
	rm ~/go/bin/h2o
	cp h2o ~/go/bin

test.proto:
	./h2o pb -f ./testdata/usertoken.yaml &>t.proto
	protoc --go_out=. --go_opt=paths=source_relative \
    --go-grpc_out=. --go-grpc_opt=paths=source_relative \
    ./t.proto

