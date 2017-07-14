## func ChanOf(dir ChanDir, t Type) Type {}

到这里的时候，Make系列也只剩下一个Chan了，并且随着复杂度的提高，这里越来越多的只能讲使用上的技巧，而非原理。 

既然有MakeChan()，必然就有ChanOf(),先看一下ChanOf()的使用方式。  

首先，它要我们传入一个ChanDir的参数，这又是什么？  

在reflect包的type.go文件里面，有这么一个定义：  
`type ChanDir int`
看来ChanDir是一个整型，而且是最小的整形。  
接着，我又看到一个注释说明：  
```
ChanDir returns a channel type's direction.It panics if the type's Kind is not Chan.
```
channel 类型的流方向。  

```go
const (
	RecvDir ChanDir             = 1 << iota // <-chan
	SendDir                                 // chan<-
	BothDir = RecvDir | SendDir             // chan
)
```
在type.go中，有一套常量定义了ChanDir的值。  
分别是1，2，3，表示出，入，两者皆是。这个<-chan和chan<-是再熟悉不过了，是channel的惯用用法。  

其实在下面的另外一个方法里面，将上述的变量进行了转换。  
```go
func (d ChanDir) String() string {
	switch d {
	case SendDir:
		return "chan<-"
	case RecvDir:
		return "<-chan"
	case BothDir:
		return "chan"
	}
	return "ChanDir" + strconv.Itoa(int(d))
}
```
还记得这个String()么，是系统自动触发的转换，很明显对应的是三个不同的chan值。  

为了回顾一下channel的用法，这里写一个简单的channel方法来复习一下。  

```go
func Sum(val string,ch chan int){
	fmt.Println("desc:",val)
	ch <- 1
}
func main()  {
	chs := make(chan int)
	go Sum("dan.liu",chs)
	<-chs
}
result:
desc: dan.liu
```
这是最简单的一个并行的程序。向channel写入一个数所，最后又写出。这个ch<1和<-chs，恐怕正是我们所需要的。  
接下来，我们试着ChanOf()一下。  

```go
func main()  {
	type T string
	tt := reflect.TypeOf(T(""))
	chanValue := reflect.ChanOf(reflect.ChanDir(1),tt)
	fmt.Println(chanValue)
}

result:
<-chan main.T
```
似乎已经成功。
接下来我再用MakeChan（）试一下。
```go
func main()  {
	type T string
	tt := reflect.TypeOf(T(""))
	chanValue := reflect.ChanOf(reflect.ChanDir(1),tt)
	fmt.Println(chanValue)

	v := reflect.MakeChan(chanValue, 1)
	v.Send(reflect.ValueOf(T("hello")))
	sv1, _ := v.Recv()
	fmt.Println(sv1)

}
result:
panic: reflect.MakeChan: unidirectional channel type
```

报错了……看源代码。
```go
if typ.ChanDir() != BothDir {
	panic("reflect.MakeChan: unidirectional channel type")
}
```

奇怪，原来Make的时候强行要BothDir才行（即值3），于是我将上面的值改成3，`chanValue := reflect.ChanOf(reflect.ChanDir(1),tt)`  
再次运行  
```go
result:
chan main.T
hello
```

成功了。可以看到这样一个很简单的流程，创建一个channel，然后丢个数据过去（"hello"），又可以把数据recv回来。  
功能是实现了，但怎么看上去不实用呢？另外后面那个`v := reflect.MakeChan(chanValue, 1)`中的1又是什么东西？  
其实这个1看MakeChan()的定义
```go
func MakeChan(typ Type, buffer int) Value {}
```
是buffer的机制，buffer 1 表示能缓冲1个，超过了就会发生阻塞。做个实验：  
```go
v := reflect.MakeChan(chanValue, 1)

v.Send(reflect.ValueOf(T("hello one")))
v.Send(reflect.ValueOf(T("hello two")))
v.Send(reflect.ValueOf(T("hello three")))
sv, _ := v.Recv()
fmt.Println(sv)

result:
fatal error: all goroutines are asleep - deadlock!
```
把1改成3即正常。  

