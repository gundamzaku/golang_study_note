## func ArrayOf(count int, elem Type) Type {}

从功能上来讲，ArrayOf()的注释已经讲得很清楚了。  
```
// ArrayOf returns the array type with the given count and element type.
// For example, if t represents int, ArrayOf(5, t) represents [5]int.
```
也就是说你传他一个数组的长度，和数组的类型，他就返回给你这么个数组变量。  
```go
func main()  {
	var t string
	rs:= reflect.ArrayOf(5,reflect.TypeOf(t))
	fmt.Println(rs.String())
}
result:
[5]string
```
不过功能虽然简单，实现的代码可是超长的。  
同样的，他要上锁和读取缓存。这听上去似乎和之前的MakeXX函数类似，其实这也算是一个MakeArray()的函数吧。  
```go
// Look in cache.
ckey := cacheKey{Array, typ, nil, uintptr(count)}
if array := cacheGet(ckey); array != nil {
	return array
}
很容易可以看到这个cache的机制，同样的，他会以数组类型ID+要生成的数组类型（比如string）+生成的数组大小来创建一个key。
```

确认缓存中不存在以后，他便会以Array的规则来Make一个Array出来，这段代码就不去理解了。  

值得注意的是这里有一个常量的定义：`const maxPtrmaskBytes = 2048`，这代表了数组的最大长度是2048，那超过会怎么样？  

