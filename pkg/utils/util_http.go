// Copyright 2023 QINIU. All rights reserved
// @Description:
// @Version: 1.0.0
// @Date: 2023/07/12 17:43
// @Author: liangfengyuan@qiniu.com

package utils

import (
	"errors"
	"fmt"
	"github.com/fengyuan-liang/GoKit/future"
	"github.com/fengyuan-liang/jet-web-fasthttp/pkg/xlog"
	"io"
	"net/http"
	"net/url"
	"os"
	"runtime/debug"
	"strings"
	"time"
)

const (
	default_timeout = time.Second * 10
	post            = "POST"
	get             = "GET"
)

var (
	ErrUnauthorized = errors.New("unauthorized")
)

func PostJson[T any](url string, requestBody any) (*T, error) {
	return PostJsonByReqFunc[T](url, requestBody, nil)
}

// PostJsonFuture 返回一个future
func PostJsonFuture[T any](url string, requestBody any) *future.Future[*T] {
	return future.FutureFunc[*T](PostJson[T], url, requestBody)
}

// PostJsonByReqFunc 默认只处理json格式的请求并将响应包装为obj
func PostJsonByReqFunc[T any](url string, requestBody any, f func(req *http.Request)) (*T, error) {
	return baseReq[T](func() *http.Request {
		var (
			req *http.Request
			err error
		)
		// 创建一个post请求
		if req, err = http.NewRequest(post, url, strings.NewReader(ObjToJsonStr(requestBody))); err != nil {
			fmt.Println("请求创建失败,", err)
			return nil
		}
		req.Header.Add("Content-Type", "application/json")
		if f != nil {
			f(req)
		}
		return req
	})
}

func PostJsonByReqFuncFuture[T any](url string, requestBody any, f func(req *http.Request)) *future.Future[*T] {
	return future.FutureFunc[*T](PostJsonByReqFunc[T], url, requestBody, f)
}

func Get[T any](url string) (*T, error) {
	return GetByReqFunc[T](url, nil)
}

func GetByReqFunc[T any](url string, f func(req *http.Request)) (*T, error) {
	return baseReq[T](func() *http.Request {
		var (
			req *http.Request
			err error
		)
		// 创建一个get请求
		if req, err = http.NewRequest(get, url, nil); err != nil {
			fmt.Println("请求创建失败,", err)
			return nil
		}
		if f != nil {
			f(req)
		}
		return req
	})
}

func init() {
	os.Setenv("HTTP_PROXY", "")
	os.Setenv("HTTPS_PROXY", "")
}

var (
	client      = &http.Client{Timeout: time.Second * 20}
	httpUtilLog = xlog.NewWith("http_util_log")
)

func baseReq[T any](f func() *http.Request) (*T, error) {
	defer func() {
		if e := recover(); e != nil {
			httpUtilLog.Errorf("baseReq panic:%v", e)
			fmt.Println(string(debug.Stack()))
		}
	}()
	var (
		req *http.Request
		err error
	)
	req = f()
	// 发送请求并获取响应
	resp, err := client.Do(req)
	if resp == nil {
		xlog.Errorf("err:%v", err)
		PrintReq(req)
		return nil, errors.New("resp is empty")
	}
	if err != nil || resp.StatusCode != 200 {
		respBody := printInfo(resp, err, req)
		if resp.StatusCode == http.StatusUnauthorized {
			return nil, ErrUnauthorized
		}
		return nil, errors.New(respBody)
	}
	defer func(Body io.ReadCloser) {
		e := Body.Close()
		if e != nil {
			fmt.Println("流关闭失败:", err)
		}
	}(resp.Body)

	data, e := io.ReadAll(resp.Body)

	if e != nil {
		fmt.Println("流读取失败:", e)
		return nil, err
	}
	t := new(T)
	err = ByteToObj(data, t)
	if err != nil {
		fmt.Printf("响应转换失败:%v\n", err)
		fmt.Printf("请求响应:\n%v\n", string(data))
		return nil, err
	}
	return t, nil
}

func printInfo(resp *http.Response, err error, req *http.Request) string {
	fmt.Printf("code:%v,请求发送失败:%v\n", resp.StatusCode, err)
	fmt.Println("================ 打印请求 ==================")
	// 打印请求
	PrintReq(req)
	fmt.Println("=============== 打印响应 ==================")
	// 打印响应
	respBody := printResp(resp)
	fmt.Println("========================================")
	return respBody
}

func printResp(resp *http.Response) string {
	// 打印响应的状态码和内容
	fmt.Println("状态码:", resp.StatusCode)
	data, _ := io.ReadAll(resp.Body)
	strData := string(data)
	resp.Body = io.NopCloser(strings.NewReader(strData))
	fmt.Println("响应内容:", strData)
	return strData
}

func PrintReq(req *http.Request) {
	// 打印url
	fmt.Printf("url:[%s]\n", req.URL.String())
	// 打印请求行
	fmt.Println(req.Method, req.URL.Path, req.Proto)
	// 打印请求头
	for key, values := range req.Header {
		for _, value := range values {
			fmt.Println(key+":", value)
		}
	}
	// 打印空行
	fmt.Println()
	// 打印请求体（如果有）
	if req.Method != http.MethodGet {
		bodyBytes, _ := io.ReadAll(req.Body)
		req.Body = io.NopCloser(strings.NewReader(string(bodyBytes)))
		fmt.Println(string(bodyBytes))
	}
}

func getHostFromURL(urlString string) (string, error) {
	parsedURL, err := url.Parse(urlString)
	if err != nil {
		return "", err
	}

	host := parsedURL.Hostname()
	return host, nil
}
