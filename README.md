# h2o
脚手架工具，统一的dsl，方便生成一些代码

## Install
```bash
go install github.com/antlabs/h2o/tree/master/cmd/h2o@latest
```

## json 子命令
```bash
# 从json文件中生成结构体, -n 选项表示分拆
h2o json -f ./test.yaml -n
# 从stdin生成结构体
h2o json -f -
```
## yaml 子命令
```bash
# 从json文件中生成结构体， -n 选项表示分拆
h2o yaml -f ./test.yaml -n
h2o yaml -f -
```
