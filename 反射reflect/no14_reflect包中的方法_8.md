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
