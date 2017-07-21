首先要安装Graphviz  
windows的话，去  
http://www.graphviz.org/Download_windows.php  
下载，并安装  

记得把bin目录追加到环境变量里面，另外他需要firefox启动，也把firefox.exe加到环境变量里面。  

接下来就可以在代码层面部署了。  

先写一段代码：

```go
```
编译  
go build  

生成文件    
reflect --cpuprofile=xxx.prof  

进入pprof中
go tool pprof reflect.exe reflect.prof

pprof>>web  
可以生成视图并查看  

其它的命令：
list func(函数名），可以看到具体函数的执行顺序list。  
weblist func(函数名），可以在浏览器上直接看，很直观。  

