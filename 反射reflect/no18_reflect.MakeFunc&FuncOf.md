
## func FuncOf(in, out []Type, variadic bool) Type {}
传入in,out，均为[]Type类型的切片，variadic为布尔
返回Type变量

在学习reflect.MakeFunc()之前，根据以往的惯例，有Make必有Of，所以很理所当然地我们找到了reflect.FuncOf()方法。  
reflect.FuncOf()的源代码很长，在此处就不帖出来了。比较难理解的是，之前的xxOf()都是让传一个Type类型，现在这里是[]Type。

怎么理解呢，我们先想办法写一个例子看一下比较好。
嗯……首先，我是很痛苦的，因为我完全不理解[]Type是个什么类型。我在声明Slice的时候，用过[]String，[]Int，那么是不是声明成[]Type就可以了？  
当然不可以，因为系统里完全不存在[]Type这个变量声明。直接报错了。  
那怎么办？其实仔细想想，Type的T是大写，Type在reflect包里面，那么我直接声明reflect.Type不就行了？  
```go
func main()  {

	var s []reflect.Type
	s=append(s,reflect.TypeOf("h"),reflect.TypeOf(1))
	rs:= reflect.FuncOf(s,s,false)
	fmt.Println(rs)
}

result:
func(string, int) (string, int)
```
果然如此，而且这样一下，也理解了in和out的定义，in就是你传入的参数，out就是你return的参数。  

接下来就简单地过一遍源代码吧。  

和其它的xxOf一样，首先要验证参数合法性，然后上锁，缓存看一看，没有的话创建（Make）Func 类型，写进缓存，返回。  
当然了，细节上可没有这么简单。  

```go
// Make a func type.
var ifunc interface{} = (func())(nil)
prototype := *(**funcType)(unsafe.Pointer(&ifunc))
n := len(in) + len(out)
```
一开始先把ifunc这个地址先划分出来，然后计算机in和out的长度，以我上面的代码为例，n为4，表示有四个参数。  
```go
var ft *funcType
var args []*rtype
```
声明两个变量，一个是方法本身，一个是方法里面的参数。  

后面是一堆很长的switch，主要是针对参数的判断逻辑，一共有6种情况，4个以内的，8个以内的，16、32、64、128个以内的。
不同的位数，定义了不同的结构体。
```go
type funcTypeFixed4 struct {
	funcType
	args [4]*rtype
}
```
这方法里的4，也有8，16，32……等等，一共也有6个结构定。  
接下来，把`*ft = *prototype`将ft指到了**prototype上
