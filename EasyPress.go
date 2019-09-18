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
	total_time float64
	total_num int
	success_num int
	fail_num int
}

/**
参数：
-m=POST  method POST/GET
-c=10    Concurrent 并发数
-d=10    duration   持续时间 秒
-u=http://localhost:8000/    url

两种实现方式：
加锁：数据绝对没有并发问题，代码实现简单清晰。
管道：实现繁琐，会有少量并发问题
 */

func main()  {

	var method string
	var url string
	var concurrent int
	var duration int64

	for _,param := range os.Args {
		if strings.Contains(param,"-m") {
			method = strings.Replace(param,"-m=","",1)
		}else if strings.Contains(param,"-u") {
			url = strings.Replace(param,"-u=","",1)
		}else if strings.Contains(param,"-c") {
			p  := strings.Replace(param,"-c=","",1)
			v, _ := strconv.ParseInt(p, 0, 64)
			concurrent = int(v)
		}else if strings.Contains(param,"-d") {
			p := strings.Replace(param,"-d=","",1)
			v, _ := strconv.ParseInt(p, 0, 64)
			duration = v
		}else if strings.Contains(param,"-h") {
			fmt.Println("-m=POST  method POST/GET")
			fmt.Println("-c=10    Concurrent 并发数")
			fmt.Println("-d=10    duration   持续时间 秒")
			fmt.Println("-u=http://localhost:8000/")
			return
		}
	}

	var pressResult * press_result = &press_result{total_num:0,fail_num:0}

	fmt.Println(url)

	//lock
	for i:=0;i<concurrent;i++ {
		go httpRequestLock(url,pressResult,duration,method)
	}

	time.Sleep(1000 * (time.Duration(duration) + 5) * time.Millisecond)

	var aveTimeStr string = fmt.Sprintf("%.2f", (pressResult.total_time / float64(pressResult.success_num)), 64)
	fmt.Printf("the avetime time is %s ,total num is %d, success num is %d, fail num is %d \n",aveTimeStr,pressResult.total_num,pressResult.success_num,pressResult.fail_num)
}


var lock sync.Mutex

/**
发生http请求方法。For循环发送，发送完睡眠900MS.
 */
func httpRequestLock(url string,pressResult * press_result,durationTime int64,method string) {

	enterTime := time.Now().UnixNano() / 1e6
	endTime := time.Now().UnixNano() / 1e6

	client := &http.Client{}

	for (endTime - enterTime) < (durationTime * 1000) {

		startTime := time.Now().UnixNano() / 1e6

		//resp, err :=   http.Get("http://tsp-uat-vhlgw.faw.cn:60080/fawtsp/faw-test-demo/test/demo?name=g3")
		req, err := http.NewRequest(method, url, strings.NewReader("n=n"))
		//req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		resp, err := client.Do(req)
		if resp != nil {
			defer resp.Body.Close()
		}
		if err != nil {
			fmt.Println(err)
		}


		endTime = time.Now().UnixNano() / 1e6
		responseTime := endTime - startTime;
		fmt.Printf("response time is %d \n" , responseTime)

		statusCode := resp.StatusCode
		_, err = ioutil.ReadAll(resp.Body)
		if err != nil {
			// handle error
		}

		lock.Lock()
		if statusCode == 200 {
			pressResult.success_num ++
			pressResult.total_time = pressResult.total_time + float64(responseTime)
		}else {
			pressResult.fail_num ++
		}
		pressResult.total_num ++
		lock.Unlock()

		time.Sleep(900 * time.Millisecond)
	}

}



