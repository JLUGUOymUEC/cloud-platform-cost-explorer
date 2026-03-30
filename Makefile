.PHONY: proto
# 检查是否规范 并生成代码
proto:
	buf lint 
	buf generate
#构建部分
.PHONY: build
build: proto
	go build -o bin/server ./server
	go build -o bin/client ./client

.PHONY: run-server
run-server: build
	./bin/server/main.go

.PHONY: run-client
run-client: build
	./bin/client/main.go

.PHONY: clean
clean:
	rm -rf bin/
	rm -rf gen/
#go mod tidy 清除不需要的依赖项，go mod vendor 将依赖项复制到vendor目录中，buf mod update 更新buf模块
.PHONY: deps
deps:
	go mod tidy 
	go mod vendor
	buf mod update

.PHONY: all
all: deps proto build #执行所有任务