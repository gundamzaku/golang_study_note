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
	type T *uintptr
	tt := reflect.TypeOf(T(nil))
	chanValue := reflect.ChanOf(reflect.ChanDir(RecvDir),tt)
	fmt.Println(chanValue)
}
```
似乎已经成功。
