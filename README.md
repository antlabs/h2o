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
## codemsg 子命令
只要实现code，自动生成String()类型和CodeMsg类型。可以在错误码代码节约时间  
1.给每个ErrNo类型生成String()方法, 内容就是注释里面的。  
2.给每个ErrNo类型生成CodeMsg{}结构
```go
// 这段代码是我们要写的。
package demo

type ErrNo int32 // 这里的类型可以自定义，--type后面类型

const (
	ENo ErrNo = 1003 // 号码出错

	ENotFound ErrNo = 1004 // 找不到
)

```

生成的CodeMsg代码如下
```go
// Code generated by "h2o codemsg --code-msg --linecomment --type ErrNo ./err.go"; DO NOT EDIT."

package demo

import (
	"encoding/json"
	"strings"
)

type CodeMsg struct {
	Code    "json:\"Code\""
	Message string "json:\"Message\""
}

func (x *CodeMsg) Error() string {
	all, _ := json.Marshal(x)
	var b strings.Builder
	b.Write(all)
	return b.String()
}

func NewCodeMsg(code) error {
	return &CodeMsg{
		Code:    code,
		Message: code.String(),
	}
}

var (
	ErrCodeMsgENo error = NewCodeMsg(ENo) //号码出错

	ErrCodeMsgENotFound error = NewCodeMsg(ENotFound) //找不到

)

```
```bash
h2o codemsg --code-msg --linecomment --type ErrNo ./err.go
```
