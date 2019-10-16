package main

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

/**
参数：
-m=POST  method POST/GET
-c=10    Concurrent 并发数
-d=10    duration   持续时间 秒
-u=http://localhost:8000/    url
-help
-head_file=C:/Users/Administrator/Desktop/head   设置header值：文件内容使用‘K=V’形式
-post_form_file=C:/Users/Administrator/Desktop/head   设置请求参数，提供两种方式：K=V键值对形式；JSON形式；
-print_body=true   是否显示返回数据
-is_send_once=true   是否只发一次数据

两种实现方式：

加锁：数据绝对没有并发问题，代码实现简单清晰。
管道：实现繁琐，会有少量并发问题
*/

//压测结果
type press_result struct {
	total_time  float64
	total_num   int
	success_num int
	fail_num    int
}

//请求参数
type request_params struct {
	print_body   string //是否打印body
	is_send_once string //是否只发送一次
	json_param   []byte //json参数
}

var wg sync.WaitGroup
var headerMap map[string]string   //用来存放header
var postFormMap map[string]string //用来存放postform
var requestParams *request_params //用来存放请求参数

func main() {

	var method string
	var url string
	var concurrent int
	var duration int64 = 1
	var hf_value string //header文件内容

	requestParams = &request_params{print_body: "fable", is_send_once: "fable"}

	for _, param := range os.Args {
		if strings.Contains(param, "-m") {
			method = strings.Replace(param, "-m=", "", 1)
		} else if strings.Contains(param, "-u") {
			url = strings.Replace(param, "-u=", "", 1)
		} else if strings.Contains(param, "-c") {
			p := strings.Replace(param, "-c=", "", 1)
			v, _ := strconv.ParseInt(p, 0, 64)
			concurrent = int(v)
		} else if strings.Contains(param, "-d") {
			p := strings.Replace(param, "-d=", "", 1)
			v, _ := strconv.ParseInt(p, 0, 64)
			duration = v
		} else if strings.Contains(param, "-help") {
			fmt.Println("---下面是基本请求参数配置项---")
			fmt.Println("-m=POST //method:POST/GET")
			fmt.Println("-c=10 //Concurrent:并发数")
			fmt.Println("-d=10 //duration:持续时间（秒）")
			fmt.Println("-u=http://localhost:8000/ //请求路径,建议加双引号-u=\"www.baidu.com\"。")
			fmt.Println("---下面是复杂请求参数配置项---")
			fmt.Println("-head_file=C:/Users/Administrator/Desktop/head.txt //设置http请求header值：文件内容使用‘K=V’形式")
			fmt.Println("-post_form_file=C:/Users/Administrator/Desktop/post_form.txt //设置请求参数，提供两种方式：K=V键值对形式；JSON形式；")
			fmt.Println("-print_body=true //是否显示返回数据")
			fmt.Println("-is_send_once=true //是否只发一次数据")
			return
		} else if strings.Contains(param, "-head_file") {
			hf_path := strings.Replace(param, "-head_file=", "", 1)
			headerMap = make(map[string]string)
			hf_value = readHeadFile(hf_path, "head_file")
			if strings.Contains(hf_value, "read_error") {
				fmt.Println(hf_value)
				return
			}
		} else if strings.Contains(param, "-post_form_file") {
			hf_path := strings.Replace(param, "-post_form_file=", "", 1)
			postFormMap = make(map[string]string)
			hf_value = readHeadFile(hf_path, "post_form_file")
			if strings.Contains(hf_value, "read_error") {
				fmt.Println(hf_value)
				return
			}
		} else if strings.Contains(param, "-print_body") {
			requestParams.print_body = strings.Replace(param, "-print_body=", "", 1)
		} else if strings.Contains(param, "-is_send_once") {
			requestParams.is_send_once = strings.Replace(param, "-is_send_once=", "", 1)
		}

	}

	var pressResult *press_result = &press_result{total_num: 0, fail_num: 0}

	fmt.Println(url)
	wg.Add(concurrent)

	for i := 0; i < concurrent; i++ {
		go httpRequestLock(url, pressResult, duration, method, hf_value)
	}

	var num int64 = 0
	for num < duration {
		time.Sleep(time.Duration(1000) * time.Millisecond)
		go printResult(pressResult, concurrent)
		num++
	}

	wg.Wait()
	printResult(pressResult, concurrent)
}

func printResult(pressResult *press_result, concurrent int) {
	aveTimeFloat := pressResult.total_time / float64(pressResult.total_num)
	var aveTimeStr string = fmt.Sprintf("%.2f", aveTimeFloat, 64)
	qps := (1000 / aveTimeFloat) * float64(concurrent)
	var qpsStr string = fmt.Sprintf("%.2f", qps, 64)
	var failPercent float64 = 0
	var failPercentStr string = "0.0"
	if pressResult.total_num > 0 {
		failPercent = float64(pressResult.fail_num) / float64(pressResult.total_num) * 100
		failPercentStr = fmt.Sprintf("%.2f", failPercent, 64)
	}
	fmt.Printf("avetime time is %s ,the qps is %s ,total num is %d, success num is %d, fail num is %d, fail percent is %s \n",
		aveTimeStr, qpsStr, pressResult.total_num, pressResult.success_num, pressResult.fail_num, failPercentStr)
}

var lock sync.Mutex

/**
发生http请求方法。For循环发送
*/
func httpRequestLock(requestPath string, pressResult *press_result, durationTime int64, method string, hf_value string) {

	enterTime := time.Now().UnixNano() / 1e6
	endTime := time.Now().UnixNano() / 1e6

	for (endTime - enterTime) < (durationTime * 1000) {

		startTime := time.Now().UnixNano() / 1e6

		var resp *http.Response
		resp = requestGP(requestPath, method)

		endTime = time.Now().UnixNano() / 1e6
		responseTime := endTime - startTime

		statusCode := resp.StatusCode
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			fmt.Println(err)
		}
		if requestParams.print_body == "true" {
			fmt.Printf("responseTime is %d \n", responseTime)
			fmt.Println(string(body))
		}

		lock.Lock()
		if statusCode == 200 {
			pressResult.success_num++
		} else {
			pressResult.fail_num++
		}
		pressResult.total_time += float64(responseTime)
		pressResult.total_num++
		lock.Unlock()

		//只发一次请求
		if requestParams.is_send_once == "true" {
			break
		}

	}
	wg.Done()
}

//封装负载request请求
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
	}
	client := &http.Client{Transport: tr}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
	}
	return resp
}

func readHeadFile(filePath string, fileType string) (value string) {
	if fileObj, err := os.Open(filePath); err == nil {
		defer fileObj.Close()
		if contents, err := ioutil.ReadAll(fileObj); err == nil {
			hf_value := string(contents)
			if strings.Index(hf_value, "{") == 0 && strings.Contains(hf_value, "}") {
				requestParams.json_param = contents
			}

			lines := strings.Split(hf_value, "\n")
			for _, line := range lines {
				if line == "" {
					continue
				}

				kvs := strings.Split(line, "=")
				value := ""
				if len(kvs) > 1 {
					value = kvs[1]
					value = strings.Replace(value, "\r", "", 1)
					value = strings.Replace(value, "\n", "", 1)
				}
				if fileType == "head_file" {
					headerMap[kvs[0]] = value
				} else if fileType == "post_form_file" {
					postFormMap[kvs[0]] = value
				}

			}
			return "read_success"
		} else {
			return "read_error:读取文件异常"
		}
	} else {
		return "read_error:找不到文件：" + filePath
	}
}
