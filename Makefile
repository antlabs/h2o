build:
	go build ./cmd/h2o/h2o.go
	./h2o codemsg --code-msg --linecomment --type ErrNo ./testdata/err.go


# 忽略
guo.dev:
	make
	rm ~/go/bin/h2o
	cp h2o ~/go/bin

test.proto:
	./h2o pb -f ./testdata/usertoken.yaml &>t.proto
	protoc --go_out=. --go_opt=paths=source_relative \
    --go-grpc_out=. --go-grpc_opt=paths=source_relative \
    ./t.proto

