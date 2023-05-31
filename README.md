# h2o
[![Go](https://github.com/antlabs/h2o/workflows/Go/badge.svg)](https://github.com/antlabs/h2o/actions)
[![codecov](https://codecov.io/gh/antlabs/h2o/branch/master/graph/badge.svg)](https://codecov.io/gh/antlabs/h2o)

脚手架工具，统一的dsl，方便生成一些代码, 静态MVC。持有model，其余的就是交给h2o生成代码。最多写写logic
目前可以直接生成的代码有:
1. http client：实现一个第三方sdk，so easy
2. http server：实现服务再不用纠结结构体的定义
3. curl 命令： 方便自测和api交流
4. protobuf定义: 从json或者yaml生成

## Install
```bash
go install github.com/antlabs/h2o/cmd/h2o@latest
```
## 命令行
* codemsg: 用于生成错误码相关代码
* http: 生成http client 和http server(代码)
* curl: 生成curl命令
* jsonstruct: 字符串转golang struct定义
* yamlstruct: 字符串转golang struct定义
* pb: dsl生成protobuf定义
```
Usage:
    ./h2o [Options] <Subcommand>

Options:
    -h,--help     print the help information

Subcommand:
    codemsg       Generate code in codemsg format from constants
    curl          gen curl command
    http          gen http code
    jsonstruct    Generate structure from json
    pb            gen protobuf code
    transport     transport
    yamlstruct    Generate structure from yaml
```
## json 子命令
```bash
# 从json文件中生成结构体, -n 选项表示分拆
h2o jsonstruct -f ./test.yaml -n
# 从stdin生成结构体, 按ctrl+d 结束
h2o jsonstruct -f -
```
## yaml 子命令
```bash
# 从json文件中生成结构体， -n 选项表示分拆
h2o yamlstruct -f ./test.yaml -n
h2o yamlstruct -f -
```

## codemsg 子命令文档
[codemsg](./codemsg.md)

## http 子命令文档
[http](./http.md)
