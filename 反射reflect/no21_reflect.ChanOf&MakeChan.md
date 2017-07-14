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

接下来问题来了，我拿并发不仅仅是做发送和接收的事情吧，比如……我想写个复杂一点的方法怎么办？可惜，找遍了网上都没有找到实现的方法，只好自己发挥想象力了。我原来的程序实现了发送和接受的两种行为。那么我把原生的方法稍微改一下能不能用呢？  
```go
func Sum(val string,ch chan int){
	fmt.Println("desc:",val)
	ch <- 1//这里是发送，是不是等同于前面的v.Send呢？
}
func main()  {
	chs := make(chan int)
	go Sum("dan.liu",chs)
	<-chs//这里是接受，是不是等同于前面的v.recv呢？
}
```

不如改写一下试试看吧。  
```go
type T string

func Sum(val string,chs reflect.Value){
	fmt.Println("desc:",val)
	chs.Send(reflect.ValueOf(T("hello one")))
}
func main()  {
	tt := reflect.TypeOf(T(""))
	chanValue := reflect.ChanOf(reflect.ChanDir(3),tt)
	fmt.Println(chanValue)

	chs := reflect.MakeChan(chanValue,1)
	go Sum("dan.liu",chs)
	chs.Recv()

}
result:
chan main.T
desc: dan.liu
```

成功了，如果我的逻辑成立的话，也就是说go xxx()这个方法，是需要我在外部自己实现的了。  

接着，我多开几个线程做一下测试：
```go

func main()  {

	tt := reflect.TypeOf(T(""))
	chanValue := reflect.ChanOf(reflect.ChanDir(3),tt)
	fmt.Println(chanValue)

	chs := reflect.MakeChan(chanValue,1024)
	var i int32
	for i = 0; i<=200; i++ {
		go Sum(i,chs)
	}
	for i = 0; i<=200; i++ {
		chs.Recv()
	}
}
```
同样没有什么问题。

其实，除了v.send和v.recv以外，在v的主体中，还存在两个方法，v.trysend和v.tryrecv，这两个到底和普通的send、recv有什么区别？  
```go
func (v Value) TrySend(x Value) bool {
	v.mustBe(Chan)
	v.mustBeExported()
	return v.send(x, true)
}
```
```go
func (v Value) Send(x Value) {
	v.mustBe(Chan)
	v.mustBeExported()
	v.send(x, false)
}
```
两个方法对比了一下，关键是在v.send()里面的true和false，原来这两个方法（首字母S大写）调的是内部的一个同名方法（首字母小写）。  
而小写的send()需要我们额外传入一个bool的参数，和是否阻塞有关系。  

追查下去，到
```go
func chansend(t *rtype, ch unsafe.Pointer, val unsafe.Pointer, nb bool) bool
```
这个方法在runtime包的chan.go中被实现，这里不谈源代码的全部，仅针对这一块看一下。  
```
 * If block is not nil,
 * then the protocol will not
 * sleep but return if it could
 * not complete.
 ```
 注释中写道，如果这个block（即nb）不为nil，那么这个协议将不会睡眠如果他还没有完全返回数据。  
 看代码吧，如果是在block为true的情况下，会执行  
 ```go
 gopark(nil, nil, "chan send (nil chan)", traceEvGoStop, 2)
 
 func gopark(unlockf func(*g, unsafe.Pointer) bool, lock unsafe.Pointer, reason string, traceEv byte, traceskip int) {}
 ```
 这么一段代码，看注释  
 ```go
// Puts the current goroutine into a waiting state and calls unlockf.
// If unlockf returns false, the goroutine is resumed.
// unlockf must not access this G's stack, as it may be moved between
// the call to gopark and the call to unlockf.
```

不行了，翻译不过来了。只是隐约看到放置一个当前的go协程放等候的状态，并且呼叫unlockf。如果unlockf这个方法返回false，这个go协程会恢复，unlockf必须不能访问G栈，它可以在调用gopark和调用unlockf之间移动？

看不懂，算了。还是回过去想想这个block到底是怎么个用法吧。  

一开始我仅仅是做了一下简单的替换，把Send换成TrySend，把recv换成了TryRecv，运行了一下，没有任何问题，看样子是要换一些极端的情况。  

比如我把chs.Recv()去掉，除了数据显示不出来，不管用哪种Send数据都没有什么异常。
接着我又把chs.Sendf去掉，用chs.Recv()的时候，数据是有的，但也报错了。  
```
desc: dan.liu
fatal error: all goroutines are asleep - deadlock!
```
所有的协程都睡着了！死锁！

于是我换成chs.TryRecv()，没报错，不过数据也没显示。看这样子是当中进行了一些异常的判断。  

我再试一下早先那个报错的代码  
```go
	tt := reflect.TypeOf(T(""))
	chanValue := reflect.ChanOf(reflect.ChanDir(3),tt)
	v := reflect.MakeChan(chanValue, 1)

	v.TrySend(reflect.ValueOf(T("hello one")))
	v.TrySend(reflect.ValueOf(T("hello two")))
	v.TrySend(reflect.ValueOf(T("hello three")))
	sv, _ := v.Recv()
	fmt.Println(sv)
```
事实上我改成了TrySend()以后，报错就没有了。看这样子，它的容错性非常高。  

MakeChan的用法就看到这里了，channel是一块内容蛮多的模块，在线程中能讲的不过是九牛二毛罢了

现在，我需要做的是把MakeChan和Select两个方法整合起来。承接之前的代码，我进行了一点调整。  

```go
type T string

func Sum(val string,chs reflect.Value){
	fmt.Println("desc:",val)
	chs.Send(reflect.ValueOf(T("hello one")))
}
func main()  {

	tt := reflect.TypeOf(T(""))
	chanValue := reflect.ChanOf(reflect.ChanDir(3),tt)
	fmt.Println(chanValue)

	chs := reflect.MakeChan(chanValue,1)
	//go Sum("dan.liu",chs)
	//chs.Recv()
	var cases []reflect.SelectCase

	cases = append(cases,reflect.SelectCase{
		Dir:  reflect.SelectRecv,
		Chan: reflect.ValueOf(chs),
		Send: reflect.ValueOf(2),
	})
	fmt.Println(cases)

	c,r,ok := reflect.Select(cases)
	fmt.Println(c)
	fmt.Println(r)
	fmt.Println(ok)

}
result:
panic: reflect.Select: RecvDir case has Send value
```
结果并没有达到我想象中的效果，非常沮丧，又不知道问题出在哪里。折腾了老半天，还是老老实实地去看源代码。  
在value.go的
```go
func Select(cases []SelectCase) (chosen int, recv Value, recvOK bool) {}
```
中定位到了问题所在。  
```go
case SelectRecv:
	if c.Send.IsValid() {
		panic("reflect.Select: RecvDir case has Send value")
	}
	ch := c.Chan
	if !ch.IsValid() {
		break
	}
	~~ch.mustBe(Chan)~~
```
ch是我们传进来的`Chan: reflect.ValueOf(chs),`仔细一看，傻比了，chs在`reflect.MakeChan(chanValue,1)`的时候已经被转成Value类型了，现在我又传了一次。

我把这段变回成`Chan: chs`  

可是又报了一个新的错：
```
panic:reflect.Select: RecvDir case has Send value
```
真是好事多磨。从字面的意思上，似乎是说我有一个Send value。  
```go
cases = append(cases,reflect.SelectCase{
	Dir:  reflect.SelectRecv,
	Chan: chs,
	Send: reflect.ValueOf(2),
})
```
那么把下面的~~Send: reflect.ValueOf(2)~~去掉，因为我的Dir是reflect.SelectRecv，并不需要Send啊。  
去掉以后错误解除，但是……
`fatal error: all goroutines are asleep - deadlock!`  
见鬼，怎么产生死锁了！  

