### 如何理解有缓存的通道和无缓存的通道  

Unbuffered channels combine communication—the exchange of a value—with synchronization—guaranteeing that two calculations (goroutines) are in a known state.

不加buffer的channels结合了通信和以同步来进行值的换取，保证两个计算(goroutines)处于已知状态。（我乱翻的囧）  

A buffered channel can be used like a semaphore, for instance to limit throughput. In this example, incoming requests are passed to handle, which sends a value into the channel, processes the request, and then receives a value from the channel to ready the “semaphore” for the next consumer. The capacity of the channel buffer limits the number of simultaneous calls to process.  

带buffered的channel能像信号量一样使用，比如用来限制吞吐量。在一个例子中，到来的请求被处理，向channel中发送了一个值，处理请求，然后从channel中接受这个值准备“信号量”给到下一个消费者。channel bufferr 的能力就是限制被同时由队列调用的数量。  

`Semaphore又称信号量，是操作系统中的一个概念，Semaphore（信号量）是用来控制同时访问特定资源的线程数量，它通过协调各个线程，以保证合理的使用公共资源。`  

好吧，看了半天也没看明白。  
再抄一段书上的话吧《Action In GO》：  

无缓冲的通道（unbuffered channel）是指在接收前没有能力保存任何值的通道。这种类型的通道要求发送goroutine和接收goroutine同时准备好，才能完成发送和接收操作。如果两个goroutine没有同时准备好，通道会导致先执行发送或接收操作的goroutine阻塞等待。这种对通道进行发送和接收的交互行为本身就是同步的。其中任意一个操作都无法离开另一个操作单独存在。  

有缓冲的通道（buffered channel）是一种在被接收前能存储一个或者多个值的通道。这种类型的通道并不强制要求goroutine之间必须同时完成发送和接收。通道会阻塞发送和接收动作的条件也会不同。只有在通道中没有要接收的值时，接收动作才会阻塞。只有在通道没有可用缓冲区容纳被发送的值时，发送动作才会阻塞。这导致有缓冲的通道和无缓冲的通道之间的一个很大的不同：无缓冲的通道保证进行发送和接收的goroutine会在同一时间进行数据交换；有缓冲的通道没有这种保证。  

还是看不懂也没关系，那就举例子好了。  


```go
func main() {

	//创建一个通道
	c := make(chan int)
	c<-100
	c<-99
	fmt.Println(<-c)
	fmt.Println(<-c)
}
```

fatal error: all goroutines are asleep - deadlock!  
goroutine 1 [chan send]:


```go
func main() {

	//创建一个通道
	c := make(chan int,2)
	c<-100
	c<-99
	fmt.Println(<-c)
	fmt.Println(<-c)
}
```
正常，从这个示例中可以看出来，第二段代码的c有两个缓冲的通道，可以放两个东西，所以没有报错。  
而第一段代码的c是没有缓冲通道的，不能放东西。所以一放就报错了。  

chan的主要用处是在go run上面，所以还是要拿一个go run的例子看一下。  

```go
package main

import (
	"fmt"
)

//开始钓鱼
func fishing(fish chan int,done chan bool)  {
	for i := 1; ; i++ {
		j, more := <-fish
		if more {
			fmt.Println("钓起",j,"条鱼")
		} else {
			fmt.Println("鱼钓完了")
			done <- true
			return
		}
	}
}

func main() {

	fish := make(chan int)
	done := make(chan bool)

	go fishing(fish,done)	//开始钓鱼

	fish<-1
	fmt.Println("第1次扔下鱼饵")

	fish<-2
	fmt.Println("第2次扔下鱼饵")

	fish<-3
	fmt.Println("第3次扔下鱼饵")


	close(fish)
	fmt.Println("鱼饵用完了")

	<-done
}
```

其中的`fish := make(chan int)`代码分别改成：  
```go 
fish := make(chan int)
fish := make(chan int，1)
fish := make(chan int，2)
fish := make(chan int，3)
```
可以看到不同的结果：
```
fish := make(chan int)  
钓起 1 条鱼
第1次扔下鱼饵
第2次扔下鱼饵
钓起 2 条鱼
钓起 3 条鱼
第3次扔下鱼饵
鱼饵用完了
鱼钓完了
```
这里是没有缓冲的时候，感觉代码有点随机，我先第1次扔下鱼饵，这个消息还没有及时显示出来，鱼已经钓上来了。这就和打雷一样，先看到闪电，才听到雷声。  
第2次就正常了，第3次又和第1次一样。  
按网上所说，无缓冲的时候是一种阻塞的模式，那么如果抛开这个打雷闪电的顺序问题的话。可以看出顺序：  
第1次（下饵->上钓）（开始准备第2个鱼饵）  
第2次（下饵->上钓）（开始准备第3个鱼饵）  
第3次（下饵->上钓）  

```
fish := make(chan int,1)  
第1次扔下鱼饵
第2次扔下鱼饵
钓起 1 条鱼
钓起 2 条鱼
钓起 3 条鱼
第3次扔下鱼饵
鱼饵用完了
鱼钓完了
```
这一次加了1个缓冲（即1个通道），这次可以看到，前面2次钓鱼因为缓冲的关系，顺序正常了。  
可以理解为，我先扔了1个鱼饵到水里，然后可以不用等鱼上钩，再扔了第2个下去。  
  
顺序就是：  
第1次（第1个下饵）（第2个下饵）  
第2次（第1条鱼上钩）（第3个下饵）  
第3次（第2条鱼上钩）  
第4次（第3条鱼上钩）  

```
fish := make(chan int,2)  

第1次扔下鱼饵
第2次扔下鱼饵
第3次扔下鱼饵
鱼饵用完了
钓起 1 条鱼
钓起 2 条鱼
钓起 3 条鱼
鱼钓完了
```
这一次加了2个缓冲（即2个通道），这次可以看到，我先扔1个鱼饵到池塘里面。然后还没等鱼钓上来，另2个鱼饵也已经扔下去了。  
顺序就是：  
第1次（第1个下饵）（第2个下饵）（第3个下饵）  
第2次（第1条鱼上钩）  
第3次（第2条鱼上钩）  
第4次（第3条鱼上钩）  
```
fish := make(chan int,3)  

第1次扔下鱼饵
第2次扔下鱼饵
第3次扔下鱼饵
鱼饵用完了
钓起 1 条鱼
钓起 2 条鱼
钓起 3 条鱼
鱼钓完了
```
这次和上次一样  
