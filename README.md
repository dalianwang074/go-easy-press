# GO_EASY_PRESS 够易压

[![go](https://img.shields.io/badge/python-3.6-green.svg?style=plastic)](https://golang.google.cn/)
[![release](https://img.shields.io/badge/django-2.1-brightgreen.svg?style=plastic)](https://www.djangoproject.com/)

when you would be a http press test, what you like to use? jmeter ?
yes ,jemeter is wonderful , but it difficute to use and it make by java waste too much memory.

GOEASYPRESS while help you to resolve the problem.
first,it easy to use. in shell execute GOEASYPRESS ,it will be run. enter -help , it will show you the param. 

GoEasyPress is a tool for http/https press.

## Running GOEASYPRESS

./gep -m=POST -u="https://www.baidu.com" -c=2 -d=10
gep.exe -m=POST -u="https://www.baidu.com" -c=2 -d=10

## -help
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
  
  
