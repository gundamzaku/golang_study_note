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

第一个问题，在重新看代码的时候，我发现cacheKey永远只会有一个，因为他的生成规则是cacheKey{Slice, typ, nil, 0}，也就是说，只要是slice变量的话，他永远都是惟一的。这又是怎么回事呢？

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
