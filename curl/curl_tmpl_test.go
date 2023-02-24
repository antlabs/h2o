package curl

import (
	"os"
	"testing"
)

func TestCurlTmpl(t *testing.T) {

	all := `{"appkey":"XXXX#XXXX","callId":"XXXX#XXXX","chat_type":"roster","eventType":"chat","from":"XXXX#XXXX","host":"XXXX","msg_id":"9XXXX2","payload":{"operation":"add","reason":"","roster_ver":"","status":{"error_code":"ok"}},"security":"XXXX","timestamp":1642648175092,"to":"tst01"}`
	tmpl := curlTmpl{Method: "POST", Header: []string{"A: Avalue", "B: Bvalue"}, URL: "www.qq.com", Data: all}
	tmpl.Gen(os.Stdout)
}
