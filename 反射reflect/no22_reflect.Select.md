## func Select(cases []SelectCase) (chosen int, recv Value, recvOK bool) {}

这个是上一课的补充，同属于channel里面的内容。所以特地跟着后面一节里来看（其实是我漏了）  

在Go的Chan里面，本身就存在select的概念，相对还是有点复杂的，所以要先重新学习一下。  

看网上很多例子都写得非常复杂和杂乱，所以我自己写了一个最简单的例子：  

```go
func Sum(val string,ch chan int){
	fmt.Println("desc:",val)
	ch <- 1
}
func main()  {
	chs := make(chan int)
	go Sum("dan.liu",chs)
	select {
		case msg1 := <- chs:
			fmt.Println(msg1)
	}
}

result:
desc: dan.liu
1
```
说白了，select就是switch那样，来判断我make()出来的channel的返回。可能一个看不太出来，那我弄两个channel。  

```go
func Sum(val string,ch chan int){
	fmt.Println("desc:",val)
	ch <- 1
}

func Sub(val string,ch chan int){
	fmt.Println("kill:",val)
	ch <- 100
}

func main()  {
	chs := make(chan int)
	chb := make(chan int)

	go Sum("dan.liu",chs)
	go Sub("master.dan",chb)

	select {
		case msg1 := <- chs:
			fmt.Println("get sum ",msg1)
		case msg2 := <- chb:
			fmt.Println("get sub",msg2)
	}
}
result1:
desc: dan.liu
get sum  1
kill: master.dan

result2:
kill: master.dan
get sub 100
```
当我弄成这种方式，以后，程序一直不很稳定。返回的数据也不稳定。  
这其实是我完全没有理清楚select的机制，实际上既然丢了两个channel上去，select也要丢两个。把代码改一改  
```go
for i := 0; i < 2; i++ {
	select {
		case msg1 := <- chs:
			fmt.Println("get sum ",msg1)
		case msg2 := <- chb:
			fmt.Println("get sub",msg2)
	}
}
```
当然，select还有写入的机制，同样看一段代码：  
```go
func Sum(val string,ch chan int){
	fmt.Println("desc:",val)
	fmt.Println(<-ch)
	ch<-1
}

func main()  {
	chs := make(chan int)

	go Sum("dan.liu",chs)
	select {
	case chs<-99:
		fmt.Println("success to write")
	}
	<-chs
}
result:
desc: dan.liu
99
success to write
```

基本上，原生的select我们是了解了，现在回到反射上面，不再多作文章。  
我们现在要用反射来产生一个select，怎么弄呢？  
看代码：
```go
func Select(cases []SelectCase) (chosen int, recv Value, recvOK bool) {}
```
他要我传入一个SelectCase切片的变量，这是什么鬼东西？  
在Value.go中我们找到了这个的定义。  
```go
type SelectCase struct {
	Dir  SelectDir // direction of case
	Chan Value     // channel to use (for send or receive)
	Send Value     // value to send (for send)
}
```
原来是一个结构体，也就是每一个case都必须包括这三个变量。  
