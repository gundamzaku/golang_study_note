整理一下现在的思绪，有必要再回顾一下之前的一些知识点。  

先看上一篇中无法执行的代码。

```go
var a = [5]int {1, 2, 3, 4, 5}
var b = [5]int {6, 7, 8, 9, 0}
c := reflect.Copy(reflect.ValueOf(a),reflect.ValueOf(b))
```

问题就出在reflect.ValueOf()上面，
