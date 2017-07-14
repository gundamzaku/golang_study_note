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
```
