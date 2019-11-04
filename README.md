# GO_EASY_PRESS 够易压

[![go](https://camo.githubusercontent.com/a6e1294d35e4a5bf8f9af4b16b1db2f2444f0549/68747470733a2f2f676f7265706f7274636172642e636f6d2f62616467652f6769746875622e636f6d2f646576656c6f7065722d6c6561726e696e672f6c6561726e696e672d676f6c616e67)](https://golang.google.cn/)



- when you would be a http press test, what you like to use? jmeter ?
yes ,jemeter is wonderful , but it difficute to use and it make by java waste too much memory.

- go_easy_press while help you to resolve the problem.
- first,it easy to use. we can use it in shell , when you enter -help , it will show you the useful.
- second,it code by go language,use less memory and could be more concurrents.is 2G memory cumputer , its easy to go 1000 concurrent.

GO_EASY_PRESS is a tool for http/https press.

## Running go_easy_press
- ./gep -help
- ./gep -m=POST -u="https://www.baidu.com" -c=2 -d=10     //linux/mac
- gep.exe -m=POST -u="https://www.baidu.com" -c=2 -d=10   //windows


