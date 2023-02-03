package pyaml

type Query struct {
	Name       string
	StructType string
}

type Body struct {
	Name       string
	StructType string
}

type Header struct {
	Name       string
	StructType string
}

type Req struct {
	Name   string //请求的结构体名
	Query  Query  //Query string 名
	Body   Body   //body名
	Header Header //header名
}

type Resp struct {
	Name   string
	Body   Body   //响应body结构体名
	Header Header //响应header结构体名
}

type ReqResp struct {
	Req  Req  //请求
	Resp Resp //响应
}

type TypeTmpl struct {
	PackageName string
	ReqResp     []ReqResp
}
