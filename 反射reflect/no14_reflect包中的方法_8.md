我的sb通过代码`sc := reflect.SliceOf(sb)`转成了sc  
接下来就要通过reflect.MakeSlice()把sc变成真正能用的slice类型。
mySlice := reflect.MakeSlice(sc,0,10)
在MakeSlice(sc,0,10)里面用
`s := sliceHeader{unsafe_NewArray(typ.Elem().(*rtype), cap), len, cap}`拼
然后再用  
`Value{typ.common(), unsafe.Pointer(&s), flagIndir | flag(Slice)}`
最后一个sclie就生成出来了。

写到这里，懂了吗？似乎不太懂。不懂？似乎又有点懂。总之，就这么创建了一个sclie出来。

关于slice的创建，我们知道的还有一种方式。
```go
func main()  {
	var my = make([]string,3)
	fmt.Print(my)
}
```
比reflect的方法要简单得多得多了。有什么区别呢？不过很可惜，这个make()也是不透明的。我们看不到具体的方法。  

好了，到这些我不得不再做一次对reflect.MakeSlice()和reflect.SliceOf()的回顾。  
有很多不理解的地方  
1、写入到集合，当然，这也可能是所有的slice对象都要被集合托管？  
2、typelinks()的机制，似乎与rsyn.pool有关系。  
3、newName()这个方法，我仍然还未有解读。  
4、fnv1的hash生成方式  

第一个问题，在重新看代码的时候，我发现cacheKey永远只会有一种类型一个key，因为他的生成规则是cacheKey{Slice, typ, nil, 0}，也就是说，只要是slice变量的话，不同的slice{变量类型}，他永远都是惟一的。这又是怎么回事呢？

做一个最简单的例子：
```go
//step.1
var my string = "hello"
sb := reflect.TypeOf(my)
fmt.Println(sb)
sc := reflect.SliceOf(sb)
fmt.Println(sc)
//step.2
var my2 string = "hello world"
sb2 := reflect.TypeOf(my2)
fmt.Println(sb2)
sc2 := reflect.SliceOf(sb2)
fmt.Println(sc2)
```
在第一次的时候，系统并未产生cache，而在第二次的时候，因为第一次的cache已经写入，故产生了cache，直接返回。这也验证了一点，cache只会有一个，然后不管你怎么生成多少slice（用reflect），都不会再去创建了，而是直接从cache中调取。

而当你上下两种分别赋属于不同的类型时  
```go
var my string = "hello"
sb := reflect.TypeOf(my)
fmt.Println(sb)
sc := reflect.SliceOf(sb)
fmt.Println(sc)
// 遍历map

var my2 bool = true
sb2 := reflect.TypeOf(my2)
fmt.Println(sb2)
sc2 := reflect.SliceOf(sb2)
fmt.Println(sc2)
```
因为rtype的不同，他也变成了两个cacheKey  

我在reflect包中写了一个方法：
```go
func GetLook(){

	for _, v := range lookupCache.m {
		println("test_")
		println(v)
	}
}

result:
test_
0x496200
test_
0x495880
```

确实产生了两个。  

那么现在的问题就可以聚焦于
```go
// Make a slice type.
var islice interface{} = ([]unsafe.Pointer)(nil)
prototype := *(**sliceType)(unsafe.Pointer(&islice))
slice := *prototype
slice.tflag = 0
slice.str = resolveReflectName(newName(s, "", "", false))
slice.hash = fnv1(typ.hash, '[')
slice.elem = typ
slice.ptrToThis = 0
```
产生一个slice type，这个type只做了一次，然后存入缓存，再也不会做了。

而正在我调试这段代码的时候，在上一节中困扰我的一个问题，突然看到了一点曙光。
```go
// Look in known types.
s := "[]" + typ.String()
for _, tt := range typesByString(s) {
	slice := (*sliceType)(unsafe.Pointer(tt))
	if slice.elem == typ {
		return cachePut(ckey, tt)
	}
}
```
就是这段之前我认为没有起到作用的代码，在我这次的时候的时候，突然被触发，并且直接返回了。  
再看一下注释：在已知的类型中查看。

可是为什么这次就触发了，而上次没有触发呢？

我对这段的理解就是，Go在系统中没有找到slice这个属性的空间定义的时候，会产生一个slice的分配（如果简单理解可以想象成是一个模板），然后这个模板就被缓存起来了，每次就会直接调用，而不需要再生成。一种变量类型（string,int)对应一个。  

range typesByString(s)是去系统中对应系统中已经存在的类型，如果存在的话，就直接返回给你（或者可以简单的认为，系统中默认已经生成好的slice），不需要你额外地再去创建，如果真是这样设定的话。。我想，那我自定义一种类型，系统岂不是必然会跳过这一步？也就是我前面碰到的不触发的情况？

于是我决定自定义一种type  
```go
type student struct {
	name string
	age int
}

func main()  {
	var my student
	sb := reflect.TypeOf(my)
	fmt.Println(sb)
	sc := reflect.SliceOf(sb)
	fmt.Println(sc)
}

……
value.go中
println(slice.tflag)
println(slice.str)
println(slice.hash)
println(slice.elem)

reulst:
0
-1
1383674263
0x49fc00
```

果然如此，虽然typesByString(s)这个方法的机制我不是特别了解，但是功能已经比较清晰。

最后，在这一节中，我再看一下newName()这个方法吧。  
newName作为resolveReflectName（）方法的参数传值，
`slice.str = resolveReflectName(newName(s, "", "", false))`
而resolveReflectName要求传入的参数必须是字节数据。
```go
type name struct {
	bytes *byte
}
```

因为newName()里面的值都是固定的，也方便我将这段代码抽出来好好测一下。  

```go
type name struct {
	bytes *byte
}
func newName(n, tag, pkgPath string, exported bool) name {

	if len(n) > 1<<16-1 {	//1<<16-1是65535，控制长度的
		panic("reflect.nameFrom: name too long: " + n)
	}
	if len(tag) > 1<<16-1 {
		panic("reflect.nameFrom: tag too long: " + tag)
	}

	var bits byte	//通过成一个字节类型

	l := 1 + 2 + len(n) //n长度是8，l给了11位长度
	//exported占1位
	if exported { //这里exported是false，跳过
		bits |= 1 << 0
	}
	//tag的话bits占2位，l还要加2位并加上tag的长度
	if len(tag) > 0 {//没有tag，跳过
		l += 2 + len(tag)
		bits |= 1 << 1
	}
	//pkgPath要占4位
	if pkgPath != "" {//为空，跳过
		bits |= 1 << 2
	}
	//根据l的字节长度，开始产生一个切片。
	b := make([]byte, l)
	//bits没生成出来，所以是0
	b[0] = bits
	b[1] = uint8(len(n) >> 8)//0，我不知道这段的意义
	b[2] = uint8(len(n))	//8
	copy(b[3:], n)//把n 复制到b[3:]后面
	if len(tag) > 0 {
		tb := b[3+len(n):]
		tb[0] = uint8(len(tag) >> 8)
		tb[1] = uint8(len(tag))
		copy(tb[2:], tag)
	}

	if pkgPath != "" {
		panic("reflect: creating a name with a package path is not supported")
	}
	return name{bytes: &b[0]}
}

func main() {
	rs:=newName("[]string","", "", false)
	fmt.Println(rs.bytes)
}
```

最后总结，就是生成一个变量的一段特定字节格式，对于name的描述，最好的解释还是只能看文档注释。  
```
// name is an encoded type name with optional extra data.
name是一种编码类型名称和一些非强制的扩展数据
//
// The first byte is a bit field containing:
第一个byte是bit字段包含
//
//	1<<0 the name is exported
	1这个名字已经被输出
//	1<<1 tag data follows the name
	2标签数据跟着名字
//	1<<2 pkgPath nameOff follows the name and tag
	4pkgPath nameOff跟着名字和标签后面
	
	好吧……完全不明白啊
//
// The next two bytes are the data length:

//
//	 l := uint16(data[1])<<8 | uint16(data[2])
//
// Bytes [3:3+l] are the string data.
//
// If tag data follows then bytes 3+l and 3+l+1 are the tag length,
// with the data following.
//
// If the import path follows, then 4 bytes at the end of
// the data form a nameOff. The import path is only set for concrete
// methods that are defined in a different package than their type.
//
// If a name starts with "*", then the exported bit represents
// whether the pointed to type is exported.
```
