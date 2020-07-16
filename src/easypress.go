package main

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
	"sync"
	"time"
	"net"
)

/**
参数：
-m=POST  method POST/GET
-c=10    Concurrent 并发数
-d=10    duration   持续时间 秒
-u=http://localhost:8000/    url
-help
-head_file=C:/Users/Administrator/Desktop/head   设置header值：文件内容使用‘K=V’形式
-form_file=C:/Users/Administrator/Desktop/head   设置请求参数，提供两种方式：K=V键值对形式；JSON形式；
-print_body=true   是否显示返回数据
-send_once=true   是否只发一次数据

*/

//压测结果
type PressResult struct {
	total_time  float64
	total_num   int
	success_num int
	fail_num    int
}

//请求参数
type RequestParams struct {
	print_body   string //是否打印body
	send_once string //是否只发送一次
	json_param   []byte //json参数
	method       string
	url          string
	concurrent   int
	duration     int64
	hf_value     string //header文件内容
}

var wg sync.WaitGroup
var headerMap map[string]string   //用来存放header
var postFormMap map[string]string //用来存放postform
var requestParams *RequestParams  //用来存放请求参数

func main() {

	requestParams = &RequestParams{print_body: "false", send_once: "false", duration: 1,concurrent:1}
	paramResult := checkParams(requestParams)
	if len(paramResult) > 0 {
		return
	}

	fmt.Println(requestParams.url)
	wg.Add(requestParams.concurrent)

	pressResultChan := make(chan PressResult, 1024)

	for i := 0; i < requestParams.concurrent; i++ {
		go httpRequest(pressResultChan)
	}

	pressResultValue := &PressResult{}
	go closeChen(pressResultChan)

	for pressResult := range pressResultChan {
		pressResultValue.total_time += pressResult.total_time
		pressResultValue.total_num += pressResult.total_num
		pressResultValue.success_num += pressResult.success_num
		pressResultValue.fail_num += pressResult.fail_num
		if pressResultValue.total_num%100 == 0 {
			printResult(pressResultValue, requestParams.concurrent)
		}
	}

	printResult(pressResultValue, requestParams.concurrent)
}

/**
发送http请求方法。For循环发送
*/
func httpRequest(pressResultChan chan PressResult) {

	enterTime := time.Now().UnixNano() / 1e6
	endTime := time.Now().UnixNano() / 1e6

	for (endTime - enterTime) < (requestParams.duration * 1000) {

		startTime := time.Now().UnixNano() / 1e6

		var resp *http.Response
		resp = requestGP(requestParams.url, requestParams.method)

		endTime = time.Now().UnixNano() / 1e6
		responseTime := endTime - startTime

		var pressResult PressResult = PressResult{total_num: 0, fail_num: 0}

		if resp == nil {
			pressResult.fail_num++
			pressResult.total_time += float64(responseTime)
			pressResult.total_num++
			pressResultChan <- pressResult
			break
		}

		statusCode := resp.StatusCode
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			fmt.Println(err)
		}
		if requestParams.print_body == "true" {
			fmt.Printf("responseTime is %d \n", responseTime)
			fmt.Println(string(body))
		}



		if statusCode == 200 {
			pressResult.success_num++
		} else {
			pressResult.fail_num++
		}
		pressResult.total_time += float64(responseTime)
		pressResult.total_num++

		pressResultChan <- pressResult

		//只发一次请求
		if requestParams.send_once == "true" {
			break
		}
	}
	wg.Done()
}

//封装request请求
func requestGP(requestPath string, method string) *http.Response {

	var bodystr string
	var bodyReader io.Reader

	//将参数文件内容设置到请求头中
	if len(postFormMap) > 0 {
		var r http.Request
		r.ParseForm()
		for kv := range postFormMap {
			r.Form.Add(kv, postFormMap[kv])
		}
		bodystr = strings.TrimSpace(r.Form.Encode())
		bodyReader = strings.NewReader(bodystr)
	}
	if len(requestParams.json_param) > 0 {
		bodyReader = bytes.NewBuffer(requestParams.json_param)
	}

	req, _ := http.NewRequest(method, requestPath, bodyReader)

	//将header文件内容设置到请求头中
	for kv := range headerMap {
		req.Header.Set(kv, headerMap[kv])
	}

	//跳过证书认证
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		DialContext:(&net.Dialer{
			Timeout:60 * time.Second,
			KeepAlive:60 * time.Second,
		}).DialContext,
	}
	client := &http.Client{Transport: tr}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
	}
	return resp
}
