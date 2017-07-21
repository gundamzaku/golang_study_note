官方教程：  
https://blog.golang.org/profiling-go-programs

首先要安装Graphviz  
windows的话，去  
http://www.graphviz.org/Download_windows.php  
下载，并安装  

记得把bin目录追加到环境变量里面，另外他需要firefox启动，也把firefox.exe加到环境变量里面。  

接下来就可以在代码层面部署了。  

先写一段代码：

```go
var cpuprofile = flag.String("cpuprofile", "", "write cpu profile to file")

func main() {
	flag.Parse()
	fmt.Println(*cpuprofile)
	if *cpuprofile != "" {
		f, err := os.Create(*cpuprofile)
		if err != nil {
			log.Fatal(err)
		}
		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
		var i int
		var val []int
		for i=0;i<10000 ;i++  {
			fmt.Println(i)
			val = append(val,i)
		}
	}
}
```
编译  
go build  

生成文件    
reflect --cpuprofile=xxx.prof  

进入pprof中
go tool pprof reflect.exe reflect.prof

pprof>>web  
可以生成视图并查看  
web mapaccess1//暂时不解  

其它的命令：
list func(函数名），可以看到具体函数的执行顺序list。

weblist func(函数名），可以在浏览器上直接看，很直观。 

top5，可以看到负载最高的数据  
```
(pprof) top5
360ms of 360ms total (  100%)
Showing top 5 nodes out of 16 (cum >= 350ms)
      flat  flat%   sum%        cum   cum%
     330ms 91.67% 91.67%      330ms 91.67%  runtime.cgocall
      10ms  2.78% 94.44%       10ms  2.78%  fmt.newPrinter
      10ms  2.78% 97.22%       10ms  2.78%  runtime.mallocgc
      10ms  2.78%   100%       10ms  2.78%  sync.(*Mutex).Unlock
         0     0%   100%      350ms 97.22%  fmt.Fprintln
```
top5 -cum，以递减的形式。

