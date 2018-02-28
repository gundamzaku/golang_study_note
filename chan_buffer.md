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
````fish := make(chan int)
钓起 1 条鱼
第1次扔下鱼饵
第2次扔下鱼饵
钓起 2 条鱼
钓起 3 条鱼
第3次扔下鱼饵
鱼饵用完了
鱼钓完了
```
