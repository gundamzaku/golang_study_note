
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
接下来，把`*ft = *prototype`将ft指到了*prototype上
```go
for _, in := range in {
	t := in.(*rtype)
	args = append(args, t)
	hash = fnv1(hash, byte(t.hash>>24), byte(t.hash>>16), byte(t.hash>>8), byte(t.hash))
}
```
遍历in，加到args里面去。  

```go
for _, out := range out {
	t := out.(*rtype)
	args = append(args, t)
	hash = fnv1(hash, byte(t.hash>>24), byte(t.hash>>16), byte(t.hash>>8), byte(t.hash))
}
```
遍历out，加到args里面去。  

到这里，我突然发现我漏了一个参数。  
```go
if variadic {
	hash = fnv1(hash, 'v')
}
```
真是见鬼，在FuncOf(in, out []Type, variadic bool)里面，variadic是作为第三个参数传入的，可是他是什么意思呢？  
看一下注释  
```
The variadic argument controls whether the function is variadic. 
可变参数控制这个方法是不是可变。。。
```
废话，可我还是没弄懂你是干嘛的啊。。  
不如我还是先用代码做测试吧，最早的代码，我把
```diff
func main()  {

	var s []reflect.Type
	s=append(s,reflect.TypeOf("h"),reflect.TypeOf(1))
	+rs:= reflect.FuncOf(s,s, false)
	-rs:= reflect.FuncOf(s,s, true)
	fmt.Println(rs)
}

result:
panic: reflect.FuncOf: last arg of variadic func must be slice
```
直接报错，原因就出在源码的第一行验证上面  
```go
if variadic && (len(in) == 0 || in[len(in)-1].Kind() != Slice) {
	panic("reflect.FuncOf: last arg of variadic func must be slice")
}
```
`in[len(in)-1].Kind() != Slice`,现在打印出来的Kind()是2，Int类型，看样子是我传入的In有问题了。  
动手改造一下。  
```go
func main()  {
	var s []reflect.Type
	var e []reflect.Type
	var vr []string
	s=append(s,reflect.TypeOf("h"),reflect.TypeOf(vr))
	e=append(e,reflect.TypeOf("h"))
	rs:= reflect.FuncOf(s,e,true)
	fmt.Println(rs)
}
result:
func(string, ...string) string
```
完美诠释了variadic的用法！  

接着再回来funcOf上面，下一段  
```go
ft.tflag = 0
ft.hash = hash
ft.inCount = uint16(len(in))
ft.outCount = uint16(len(out))
if variadic {
	ft.outCount |= 1 << 15
}
```

这里就是函数内部的常规定义了，注意的是，variadic再次出现，并且这次明显是用来控制outCount的，如果是variadic为true的话，outCount的大小是……32769？似乎有点夸张。并且目前我也不知道他的具体用意。  

再接下来，就是我们的老朋友funcLookupCache了，先上锁，现去查缓存。   
```go
funcLookupCache.RLock()
for _, t := range funcLookupCache.m[hash] {
	if haveIdenticalUnderlyingType(&ft.rtype, t, true) {
		funcLookupCache.RUnlock()
		return t
	}
}
funcLookupCache.RUnlock()
```
很显然，第一次没有。

没有缓存，再跑一次  
```go
// Not in cache, lock and retry.
funcLookupCache.Lock()
defer funcLookupCache.Unlock()
if funcLookupCache.m == nil {//肯定是空的
	funcLookupCache.m = make(map[uint32][]*rtype)
}
for _, t := range funcLookupCache.m[hash] {
	if haveIdenticalUnderlyingType(&ft.rtype, t, true) {
		return t
	}
}
```
```go
str := funcStr(ft)
```
接成字符串，之前的xxOf都是直接手写的，这里估计太复杂了，就做成了一个函数，最后呈现在的结果以我的代码为主的话就是：func(string, ...string) string

```go
for _, tt := range typesByString(str) {
	if haveIdenticalUnderlyingType(&ft.rtype, tt, true) {
		funcLookupCache.m[hash] = append(funcLookupCache.m[hash], tt)
		return tt
	}
}
```
可惜，typesByString(str)也拆不出相对应的类型出来，直接就是nil跳过了。  

```go
// Populate the remaining fields of ft and store in cache.
ft.str = resolveReflectName(newName(str, "", "", false))
ft.ptrToThis = 0
funcLookupCache.m[hash] = append(funcLookupCache.m[hash], &ft.rtype)

return &ft.rtype
```

最后，串成一个完整的函数类型的格式，并且写入缓存，返回。  
看不懂没关系，知道大概怎么个流程就行了。
