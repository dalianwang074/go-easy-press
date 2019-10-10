package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

type press_result struct {
	total_time  float64
	total_num   int
	success_num int
	fail_num    int
}

/**
参数：
-m=POST  method POST/GET
-c=10    Concurrent 并发数
-d=10    duration   持续时间 秒
-u=http://localhost:8000/    url
-help
-head_file=C:/Users/Administrator/Desktop/head   设置header值：文件内容使用‘K=V’形式
-post_form_file=C:/Users/Administrator/Desktop/head   post_form：文件内容使用‘K=V’形式

两种实现方式：
加锁：数据绝对没有并发问题，代码实现简单清晰。
管道：实现繁琐，会有少量并发问题
*/

var wg sync.WaitGroup
var headerMap map[string]string   //用来存放header
var postFormMap map[string]string //用来存放postform

func main() {

	var method string
	var url string
	var concurrent int
	var duration int64
	var hf_value string //header文件内容

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
			fmt.Println("-m=POST	method:POST/GET")
			fmt.Println("-c=10	Concurrent:并发数")
			fmt.Println("-d=10	duration:持续时间（秒）")
			fmt.Println("-u=http://localhost:8000/	请求路径")
			fmt.Println("---下面是复杂请求参数配置项---")
			fmt.Println("-head_file=C:/Users/Administrator/Desktop/head.txt		设置http请求header值：文件内容使用‘K=V’形式")
			fmt.Println("-post_form_file=C:/Users/Administrator/Desktop/post_form.txt	设置post_form表单值：文件内容使用‘K=V’形式")
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
	aveTimeFloat := pressResult.total_time / float64(pressResult.success_num)
	var aveTimeStr string = fmt.Sprintf("%.2f", aveTimeFloat, 64)
	qps := (1000 / aveTimeFloat) * float64(concurrent)
	var qpsStr string = fmt.Sprintf("%.2f", qps, 64)
	fmt.Printf("the avetime time is %s ,the qps is %s ,total num is %d, success num is %d, fail num is %d \n",
		aveTimeStr, qpsStr, pressResult.total_num, pressResult.success_num, pressResult.fail_num)
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
		ioutil.ReadAll(resp.Body)

		lock.Lock()
		if statusCode == 200 {
			pressResult.success_num++
			pressResult.total_time = pressResult.total_time + float64(responseTime)
		} else {
			pressResult.fail_num++
		}
		pressResult.total_num++
		lock.Unlock()

	}
	wg.Done()
}

//封装负载request请求
func requestGP(requestPath string, method string) *http.Response {

	var bodystr string

	//将header文件内容设置到请求头中
	if len(postFormMap) > 0 {
		var r http.Request
		r.ParseForm()
		for kv := range postFormMap {
			r.Form.Add(kv, postFormMap[kv])
		}
		bodystr = strings.TrimSpace(r.Form.Encode())
	}

	req, _ := http.NewRequest(method, requestPath, strings.NewReader(bodystr))

	//将header文件内容设置到请求头中
	for kv := range headerMap {
		req.Header.Set(kv, headerMap[kv])
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if resp != nil {
		defer resp.Body.Close()
	}
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
			lines := strings.Split(hf_value, "\n")
			for _, line := range lines {
				kvs := strings.Split(line, "=")
				value := strings.Replace(kvs[1], "\r", "", 1)
				value = strings.Replace(value, "\n", "", 1)
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
