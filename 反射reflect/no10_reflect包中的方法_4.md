学习完了reflect.Append()方法之后，其实我还有不少的问题。  
首先一个问题就是，reflect.Append()是为一个切片添加内容，那我为什么不用原生的Slice切片来添加呢？原生的应该更简单吧。  

```go
func main() {
	slice := []string{"Mon","Tues","Wed","Thur","Fri"}
	newSlice := append(slice,"sat","Sun");
	fmt.Println(newSlice)
}
```
同样的方法，完全实现了一样的功能。
