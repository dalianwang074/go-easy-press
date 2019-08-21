when you would be a http press test, what you like to use? jmeter ?
yes ,jemeter is wonderful , but it difficute to use and it make by java waste too much memory.

GOEASYPRESS while help you to resolve the problem.
first,it easy to use. in shell execute GOEASYPRESS ,it will be run. enter -h , it will show you the param. 

-m=post  method post/get
-c=10    Concurrent 并发数
-d=10    duration   持续时间 秒
-u=http://localhost:8000/

thie is a simple run example:
press_test_client.exe -m=post -c=10 -d=10 -u=http://localhost:8000/
and the result is ：
the avetime time is %!d(float64=26.908333333333335) ,total num is 120, success num is 120, fail num is 0 

second,it code by go language,use less memory and could be more concurrents.
is 4G memory cumputer , its easy to go 1000 concurrent.

GOEASYPRESS can be use in windows/linux/mac, whatever you want.

welcome to go GOEASYPRESS!!
