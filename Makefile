build:
	go build ./cmd/h2o/h2o.go

# 忽略
guo.dev:
	make
	rm ~/go/bin/h2o
	cp h2o ~/go/bin
