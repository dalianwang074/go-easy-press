package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
)

//打印压测结果方法
func printResult(pressResult *PressResult, concurrent int) {
	aveTimeFloat := pressResult.total_time / float64(pressResult.total_num)
	tps := (1000 / aveTimeFloat) * float64(concurrent)
	var failPercent float64 = 0
	if pressResult.total_num > 0 {
		failPercent = float64(pressResult.fail_num) / float64(pressResult.total_num) * 100
	}
	fmt.Printf("avetime time is %dms ,the tps is %d ,total num is %d, fail num is %d, fail percent is %d%% \n",
		int(aveTimeFloat), int(tps), pressResult.total_num, pressResult.fail_num, int(failPercent))
}

//读本地文件，主要是header和body文件
func readFile(filePath string, fileType string) (value string) {
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
				} else if fileType == "body_file" {
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
			fmt.Println("GoEasyPress is a tool for http/https press.")
			fmt.Println()
			fmt.Println("Usage:")
			fmt.Println("\tgep.exe <command>=[arguments]\t //windows")
			fmt.Println("\tgep <command>=[arguments]\t //linux/mac")
			fmt.Println("\tgep.exe -m=POST -u=\"https://www.baidu.com\" -c=2 -d=10")
			fmt.Println()
			fmt.Println("The Command :")
			fmt.Println("\t-m\tmethod,POST/GET")
			fmt.Println("\t-c\tconcurrent，并发数")
			fmt.Println("\t-d\tduration，持续时间（秒）")
			fmt.Println("\t-u\turl，请求路径,建议加双引号.-u=\"www.baidu.com\"")
			fmt.Println("\t-head_file\t设置header参数，-head_file=C:/Users/head.txt")
			fmt.Println("\t-body_file\t设置BODY参数，文本提供两种方式：K=V键值对形式；JSON形式")
			fmt.Println("\t-print_body\t是否显示返回数据，默认false.-print_body=true")
			fmt.Println("\t-send_once\t是否只发一次数据，默认false.-send_once=true")
			fmt.Println()
			fmt.Println("The Explain :")
			fmt.Println("\t请求读取文本内容默认使用‘K=V’形式存储")
			return "help"
		} else if strings.Contains(param, "-head_file") {
			hf_path := strings.Replace(param, "-head_file=", "", 1)
			headerMap = make(map[string]string)
			requestParams.hf_value = readFile(hf_path, "head_file")
			if strings.Contains(requestParams.hf_value, "read_error") {
				fmt.Println(requestParams.hf_value)
				return "read_error"
			}
		} else if strings.Contains(param, "-body_file") {
			hf_path := strings.Replace(param, "-body_file=", "", 1)
			postFormMap = make(map[string]string)
			requestParams.hf_value = readFile(hf_path, "body_file")
			if strings.Contains(requestParams.hf_value, "read_error") {
				fmt.Println(requestParams.hf_value)
				return "read_error"
			}
		} else if strings.Contains(param, "-print_body") {
			requestParams.print_body = strings.Replace(param, "-print_body=", "", 1)
		} else if strings.Contains(param, "-send_once") {
			requestParams.send_once = strings.Replace(param, "-send_once=", "", 1)
		}

	}
	return ""
}
