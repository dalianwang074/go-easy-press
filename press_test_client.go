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
-m=post  method post/get
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
			fmt.Println("-m=post  method post/get")
			fmt.Println("-c=10    Concurrent 并发数")
			fmt.Println("-d=10    duration   持续时间 秒")
			fmt.Println("-u=http://localhost:8000/")
			return
		}
	}

	var pressResult * press_result = &press_result{total_num:0,fail_num:0}

	//lock
	for i:=0;i<concurrent;i++ {
		go httpRequestLock(url,pressResult,duration,method)
	}

	//chan  需要 sleep之后 close(resultChan)
	//resultChan := make(chan press_result,200)
	//for i:=0;i<threadNum;i++ {
	//	go httpPostChan(url,durationTime,resultChan)
	//}
	//go ave_calculate(resultChan,pressResult)


	time.Sleep(1000 * (time.Duration(duration) + 5) * time.Millisecond)
	//close(resultChan)

	fmt.Printf("the avetime time is %d ,total num is %d, success num is %d, fail num is %d \n",
		pressResult.total_time / float64(pressResult.success_num),pressResult.total_num,pressResult.success_num,pressResult.fail_num)
}


var lock sync.Mutex
func httpRequestLock(url string,pressResult * press_result,durationTime int64,method string) {

	enterTime := time.Now().UnixNano() / 1e6
	endTime := time.Now().UnixNano() / 1e6

	client := &http.Client{}

	for (endTime - enterTime) < (durationTime * 1000) {

		startTime := time.Now().UnixNano() / 1e6

		//resp, err := http.Post(url,"json",strings.NewReader("n=n"))
		//resp, err := http.Get(url)

		req, err := http.NewRequest(method, url, strings.NewReader("n=n"))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

		resp, err := client.Do(req)
		defer resp.Body.Close()
		if err != nil {
			fmt.Println(err)
		}

		endTime = time.Now().UnixNano() / 1e6
		v := endTime - startTime;
		fmt.Printf("response time is %d \n" , v)

		lock.Lock()
		if resp.StatusCode == 200 {
			pressResult.success_num ++
			pressResult.total_time = pressResult.total_time + float64(v)
		}else {
			pressResult.fail_num ++
		}
		pressResult.total_num ++
		lock.Unlock()

		_, err = ioutil.ReadAll(resp.Body)
		if err != nil {
			// handle error
		}

		time.Sleep(900 * time.Millisecond)
	}

}




func ave_calculate(resultChan chan press_result,pressResult * press_result){
	for rc := range resultChan {
		pressResult.total_time = pressResult.total_time + rc.total_time
		pressResult.success_num = pressResult.success_num + rc.success_num
		pressResult.total_num = pressResult.total_num + rc.total_num
		pressResult.fail_num = pressResult.fail_num + rc.fail_num
		fmt.Printf("the success num is: %d ,and the total_time is %d \n",rc.success_num,rc.total_time)
	}
}

func httpPostChan(url string,durationTime int64,resultChan chan press_result) {

	var pr press_result = press_result{total_num:0,fail_num:0}

	enterTime := time.Now().UnixNano() / 1e6
	endTime := time.Now().UnixNano() / 1e6

	for (endTime - enterTime) < (durationTime * 1000) {

		startTime := time.Now().UnixNano() / 1e6
		//resp, err := http.Post(url,"json",strings.NewReader("name=g1"))
		resp, err := http.Get(url)
		if err != nil {
			fmt.Println(err)
		}
		defer resp.Body.Close()

		endTime = time.Now().UnixNano() / 1e6
		v := endTime - startTime;
		fmt.Printf("receive time is %d \n" , v)

		if resp.StatusCode == 200 {
			pr.success_num++
			pr.total_time = pr.total_time + float64(v)
		}else {
			pr.fail_num ++
		}
		pr.total_num ++

		_, err = ioutil.ReadAll(resp.Body)
		if err != nil {
			// handle error
		}

		resultChan <- pr

		time.Sleep(900 * time.Millisecond)
	}

}
