package pyaml

import (
	stdjson "encoding/json"
	"fmt"
	"net/http"
	"sort"
	"strings"

	"github.com/antlabs/h2o/model"
	"github.com/antlabs/tostruct/header"
	"github.com/antlabs/tostruct/json"
	"github.com/antlabs/tostruct/name"
	"github.com/antlabs/tostruct/option"
	"github.com/antlabs/tostruct/protobuf"
)

func GetBody(h model.Muilt, isProtobuf bool) (reqBody Body, defReqBody []model.KeyVal[string, string], respBody Body, err error) {

	newReqType := h.Req.NewType
	if isProtobuf {
		newReqType = h.Req.NewProtobufType
	}

	reqBody, defReqBody, err = getBody(h.Req.Name,
		h.Req.Body,
		newReqType,
		h.Req.Encode,
		h.Req.UseDefault.Body, isProtobuf)
	if err != nil {
		fmt.Printf("get request body:%s\n", err)
		return
	}

	newRespType := h.Resp.NewType
	if isProtobuf {
		newRespType = h.Resp.NewProtobufType
	}

	respBody, _, err = getBody(h.Resp.Name, h.Resp.Body, newRespType, model.Encode{}, nil, isProtobuf)
	if err != nil {
		fmt.Printf("get response body:%s \n", err)
		all, _ := stdjson.Marshal(h.Resp.Body)
		fmt.Println(string(all), err, h.Resp.Body == nil)
		return
	}
	return
}

func GetHeader(h model.Muilt, opt ...option.OptionFunc) (reqHeader Header, defReqHeader []model.KeyVal[string, string],
	respHeader Header, defRespHeader []model.KeyVal[string, string], err error) {

	reqHeader, defReqHeader, err = getHeader(h.Req.Name+"Header", h.Req.Header, h.Req.UseDefault.Header, opt...)
	if err != nil {
		fmt.Printf("get request header:%s, %v\n", err, h.Req.Header)
		return
	}

	respHeader, _, err = getHeader(h.Resp.Name+"Header", h.Resp.Header, nil, opt...)
	if err != nil {
		fmt.Printf("get request body:%s\n", err)
		return
	}
	return
}

func getBody(bodyName string, bodyData any, newType map[string]string, encode model.Encode, bodyDefKey []string, isProtobuf bool) (
	body Body,
	rvDefaultBody []model.KeyVal[string, string],
	err error) {

	body.Name = bodyName + "Body"

	tagName := "json"
	if encode.Body == model.WWWForm {
		tagName = "form"
	}

	getVal := make(map[string]any)

	for _, v := range bodyDefKey {
		getVal[v] = ""
	}

	var data []byte
	switch v := bodyData.(type) {
	case map[string]any:
		if isProtobuf {
			data, err = protobuf.Marshal(v,
				option.WithStructName(body.Name),
				option.WithSpecifyType(newType))
		} else {
			data, err = json.Marshal(v,
				option.WithStructName(body.Name),
				option.WithTagName(tagName),
				option.WithSpecifyType(newType),
				option.WithGetRawValue(getVal))
		}
	case []any:
		if isProtobuf {
			data, err = protobuf.Marshal(v,
				option.WithStructName(body.Name),
				option.WithSpecifyType(newType),
			)
		} else {
			data, err = json.Marshal(v,
				option.WithStructName(body.Name),
				option.WithTagName(tagName),
				option.WithSpecifyType(newType),
				option.WithGetRawValue(getVal))
		}
	default:
		body.Name = ""
	}

	if len(getVal) > 0 {
		rvDefaultBody = make([]model.KeyVal[string, string], 0, len(getVal))
		for k, v := range getVal {
			fieldName, _ := name.GetFieldAndTagName(k)
			rvDefaultBody = append(rvDefaultBody,
				(&model.KeyVal[string, string]{Key: fieldName, Val: fmt.Sprint(v), RawVal: v}).SetIs())

		}
		sort.Slice(rvDefaultBody, func(i, j int) bool {
			return rvDefaultBody[i].Key < rvDefaultBody[j].Key
		})
	}

	body.StructType = string(data)
	return
}

func getHeader(headerName string, headerSlice []string, defaultHeader []string, opt ...option.OptionFunc) (
	htmpl Header,
	rvDefaultHeader []model.KeyVal[string, string],
	err error) {

	// http header
	if len(headerSlice) == 0 {
		return
	}

	hmap := sliceToHTTPHeader(headerSlice)
	htmpl.Name = headerName

	getVal := make(map[string]any)

	for _, v := range defaultHeader {
		_, ok := hmap[v]
		if !ok {
			continue
		}

		getVal[v] = ""
	}

	opt = append(opt, option.WithStructName(htmpl.Name),
		option.WithTagName("header"), option.WithTagNameFromKey(), option.WithGetRawValue(getVal))
	var data []byte
	data, err = header.Marshal(hmap, opt...)
	if err != nil {
		return
	}

	if len(getVal) > 0 {
		rvDefaultHeader = make([]model.KeyVal[string, string], 0, len(getVal))
		for k, v := range getVal {
			fieldName, _ := name.GetFieldAndTagName(k)
			rvDefaultHeader = append(rvDefaultHeader, (&model.KeyVal[string, string]{Key: fieldName, Val: fmt.Sprint(v), RawVal: v}).SetIs())
		}
		sort.Slice(rvDefaultHeader, func(i, j int) bool {
			return rvDefaultHeader[i].Key < rvDefaultHeader[j].Key
		})
	}

	htmpl.StructType = string(data)
	return
}

func sliceToHTTPHeader(headerSlice []string) http.Header {

	hmap := make(http.Header)
	for _, v := range headerSlice {
		pos := strings.Index(v, ":")
		if pos == -1 {
			continue
		}

		val := v[pos+1:]
		if len(val) == 0 {
			continue
		}
		if val[0] == ' ' {
			val = val[1:]
		}
		hmap.Set(v[:pos], val)
	}
	return hmap
}
