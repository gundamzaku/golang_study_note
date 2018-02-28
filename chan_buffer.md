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
`
fatal error: all goroutines are asleep - deadlock!

goroutine 1 [chan send]:
`

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
正常
