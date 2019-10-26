package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	"time"
)

//打印压测结果方法
func printResult(pressResult *PressResult, concurrent int) {
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

//读本地文件，主要是header和params文件
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

func printResultSchedule(pressResult *PressResult, concurrent int) {
	for {
		time.Sleep(time.Duration(1000) * time.Millisecond)
		printResult(pressResult, concurrent)
	}
}

func closeChen(pressResultChan chan PressResult) {
	wg.Wait()
	close(pressResultChan)
}

//处理参数方法
func checkParams(requestParams *RequestParams) string {
	for _, param := range os.Args {
		if strings.Contains(param, "-m") {
			requestParams.method = strings.Replace(param, "-m=", "", 1)
		} else if strings.Contains(param, "-u") {
			requestParams.url = strings.Replace(param, "-u=", "", 1)
		} else if strings.Contains(param, "-c") {
			p := strings.Replace(param, "-c=", "", 1)
			v, _ := strconv.ParseInt(p, 0, 64)
			requestParams.concurrent = int(v)
		} else if strings.Contains(param, "-d") {
			p := strings.Replace(param, "-d=", "", 1)
			v, _ := strconv.ParseInt(p, 0, 64)
			requestParams.duration = v
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
			return "help"
		} else if strings.Contains(param, "-head_file") {
			hf_path := strings.Replace(param, "-head_file=", "", 1)
			headerMap = make(map[string]string)
			requestParams.hf_value = readHeadFile(hf_path, "head_file")
			if strings.Contains(requestParams.hf_value, "read_error") {
				fmt.Println(requestParams.hf_value)
				return "read_error"
			}
		} else if strings.Contains(param, "-post_form_file") {
			hf_path := strings.Replace(param, "-post_form_file=", "", 1)
			postFormMap = make(map[string]string)
			requestParams.hf_value = readHeadFile(hf_path, "post_form_file")
			if strings.Contains(requestParams.hf_value, "read_error") {
				fmt.Println(requestParams.hf_value)
				return "read_error"
			}
		} else if strings.Contains(param, "-print_body") {
			requestParams.print_body = strings.Replace(param, "-print_body=", "", 1)
		} else if strings.Contains(param, "-is_send_once") {
			requestParams.is_send_once = strings.Replace(param, "-is_send_once=", "", 1)
		}

	}
	return ""
}
