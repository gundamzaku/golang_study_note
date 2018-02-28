如何理解有缓存的通道和无缓存的通道  

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
正常，从这个示例中可以看出来，第二段代码的c有两个缓存的通道，可以放两个东西，所以没有报错。
而第一段代码的c是没有缓存通道的，不能放东西。所以一放就报错了。  

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
这里是没有缓存的时候，感觉代码有点随机，我先第1次扔下鱼饵，这个消息还没有及时显示出来，鱼已经钓上来了。这就和打雷一样，先看到闪电，才听到雷声。  
第2次就正常了，第3次又和第1次一样。  
按网上所说，无缓存的时候是一种阻塞的模式，那么如果抛开这个打雷闪电的顺序问题的话。可以看出顺序：  
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
这一次加了1个缓存（即1个通道），这次可以看到，前面2次钓鱼因为缓存的关系，顺序正常了。
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
这一次加了2个缓存（即2个通道），这次可以看到，我先扔1个鱼饵到池塘里面。然后还没等鱼钓上来，另2个鱼饵也已经扔下去了。  
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
