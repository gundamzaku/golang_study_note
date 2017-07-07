学习完了reflect.Append()方法之后，其实我还有不少的问题。  
首先一个问题就是，reflect.Append()是为一个切片添加内容，那我为什么不用原生的Slice切片来添加呢？原生的应该更简单吧。  

```go
func main() {
	slice := []string{"Mon","Tues","Wed","Thur","Fri"}
	newSlice := append(slice,"sat","Sun");
	fmt.Println(newSlice)
}
```
同样的方法名，完全实现了一样的功能。那它和`reflect.Append()`有什么区别呢？  
append()方法在builtin包中的builtin.go文件中被声明，
`func append(slice []Type, elems ...Type) []Type`

可惜我找遍了整个src目录，都没有找到append()的实现方法，看网上有提到`总的意思是我在源码找到的 builtin.go 只是一个文档，并不是源代码，这些函数是经过特殊处理的，我们不能直接去使用他们。`，按这个说法的话，这个方法对我们而言，是不透明的了。这大概也是反射和原生的一个区别。
不过看注释的时候，注意到了有趣的一点说明：  
```go
// As a special case, it is legal to append a string to a byte slice, like this:
作为一个特别的安全，将字符串追加到byte类型的切片中是合法的。
// slice = append([]byte("hello "), "world"...)
```
很有意思的小技艺，能记的话就记一下。  
