## http
本文档介绍http子命令的用法，主要涉及两个flag，client和server
* client: 主要场景主要用于快速实现一些第三方sdk。
* server: 主要场景主要快速实现一些callback服务

## 一个简单的配置文件
该配置文件使用yaml描述http元数据。请求体里面支持curl命令(目前不支持form表单，有时间再支持一波)
保存为apptoken.yaml 文件名
```yaml
---
package: "message" #生成go的包名
protobuf:
  package: "message.v1" #生成protobuf时使用，可以先忽略这个配置
  go_package: "./pb/v1/message" # 生成protobuf时，指定生成的go package的名字，可以先忽略这个配置
init:
  handler: New #函数名, 名字是New认识是构造函数
  rvStruct:
    name: Message #返回的结构体名, 如果生成客户端代码，这个是必须的
    field:
      Host: ''   #生成http 客户端结构体，自己配置的字段
      OrgName: '' #生成http 客户端结构体，自己配置的字段
      AppName: ''  #生成http 客户端结构体，自己配置的字段
      Username: '' #生成http 客户端结构体，自己配置的字段

multi:
- handler: Message.SendTxtMessageAsync  #包名.方法名, 这样命名
  req:
    # 必填字段，例 method: POST
    method: 
    curl: >-
      curl -X POST -i 'http://{{.Host}}/{{.OrgName}}/{{.AppName}}/messages/users' -H 'Content-Type: application/json' -H 'Accept: application/json' -H 'Authorization: Bearer <YourAppToken>' -d '{"from": "user1","to": ["user2"],"type": "txt","body": {"msg": "testmessages"}}'

    # 必填字段
    # url: https://host/port
    url: 
    template:
      url: true
    # 直接使用配置文件里面的值
    useDefault:
      header:
      #- Content-Type
      #- Accept
      body:
        #- .grant_type
    header:  #可配置的http header，目前选配
    - 'Content-Type: application/json'
    - 'Accept: application/json'
    encode: 
      body:
    body:

  resp:
    newType:
    # 响应体，直接复制一个json就行
    body: >-
      {
      "path": "/messages/users",
      "uri": "https://XXXX/XXXX/XXXX/messages/users",
      "timestamp": 1657254052191,
      "organization": "XXXX",
      "application": "e82bcc5f-XXXX-XXXX-a7c1-92de917ea2b0",
      "action": "post",
      "data": {
        "user2": "1029457500870543736"
      },
      "duration": 0,
      "applicationName": "XXXX"
      }

```

##  生成客户端代码
h2o http -f ./apptoken.yaml --client

## 生成服务端代码
h2o http -f ./apptoken.yaml --server
