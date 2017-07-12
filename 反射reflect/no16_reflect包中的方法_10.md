接下来，我们再看一个make方法：`func MakeMap(typ Type) Value {}`吧，看上去还是比较简单的。从名字上看，必然是生成一个map的集合了。  
先写一个例子。  
```go
func main()  {
	var val string
	m:=reflect.MakeMap(reflect.TypeOf(val))
	fmt.Println(m)
}
result:
panic: reflect.MakeMap of non-map type
```
报错了，这是怎么回事？  
错误出在MakeMap代码内的验证上面：  
```go
if typ.Kind() != Map {
	panic("reflect.MakeMap of non-map type")
}
```
传入的不是map类型？我不是要生成一个map么，怎么一定要先定义成map类型了。先不管他了，既然人家的源代码有这种限制，只能照做，我重新改一下代码。  
```go
func main()  {
	var val map[string]string
	m:=reflect.MakeMap(reflect.TypeOf(val))
	fmt.Println(m)
}

result:
map[]
```
感觉把`m:=reflect.MakeMap(reflect.TypeOf(val))`这一段去掉也没关系啊。。  
