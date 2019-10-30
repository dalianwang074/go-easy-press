when you would be a http press test, what you like to use? jmeter ?
yes ,jemeter is wonderful , but it difficute to use and it make by java waste too much memory.

GOEASYPRESS while help you to resolve the problem.
first,it easy to use. in shell execute GOEASYPRESS ,it will be run. enter -help , it will show you the param. 

GoEasyPress is a tool for http/https press.

Usage:
	gep.exe <command>=[arguments]       //windows
	gep <command>=[arguments]           //linux/mac
	gep.exe -m=POST -u="https://www.baidu.com" -c=2 -d=10

The command are:
	-m		        method,POST/GET
	-c		        concurrent，并发数
	-d		        duration，持续时间（秒）
	-u		        url，请求路径,建议加双引号.-u="www.baidu.com"
	-head_file	    设置http请求header值：文件内容使用‘K=V’形式
			        -head_file=C:/Users/Desktop/head.txt
	-post_form_file	设置请求参数，提供两种方式：K=V键值对形式；JSON形式
			        -post_form_file=C:/Users/Desktop/post_form.txt
	-print_body	    是否显示返回数据，默认false.-print_body=true
	-is_send_once	是否只发一次数据，默认false.-is_send_once=true

thie is a simple run example:
gep.exe -m=post -c=10 -d=10 -u=http://localhost:8000/       //windows
gep -m=post -c=10 -d=10 -u=http://localhost:8000/           //linux/mac
and the result is ：
the avetime time is %!d(float64=26.908333333333335) ,total num is 120, success num is 120, fail num is 0 

second,it code by go language,use less memory and could be more concurrents.
is 4G memory cumputer , its easy to go 1000 concurrent.

GOEASYPRESS can be use in windows/linux/mac, whatever you want.

welcome to go GOEASYPRESS!!
