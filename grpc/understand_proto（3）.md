## Maps  
如果你想创建一个关联map作为你数据定义的一部分，protocol buffers提供了一种便捷的语法：

```go
map<key_type, value_type> map_field = N;
```

`key_type`可以是任意的整形或字符类型（除了浮点指针类型和bytes的任意可变类型）。注意枚举并不是一个有效的key_type。value_type的话，除了其它的map以外的任意类型都是可以的。  

举个例子，如果你想创建一个项目的map，每一个项目message均与字符串键值进行关联。你能如下进行定义：  

```go
map<string, Project> projects = 3;
```

* map字段不能重复  
* Wire格式排序和map数值的map迭代排序是未定义的，因此你不能依赖你的map组进行特定的排序。
* 当.proto产生文本格式时，maps已经按key进行了排序。数字键会按数字进行排序。
* 当从wire中解析或当合并时，如果发现有重复的map keys的话，那么用其中的最后一个key。当解析文本格式中的map时，如果有重复的keys的话，解析会失败。

生成map的API目前已经被所有支持proto3的语言可用了。你能在相应的语言的<a href="https://developers.google.com/protocol-buffers/docs/reference/overview">相关文档</a>中找到关于这个map API的内容。（见鬼去吧……）

### 向后兼容  

map语法在wire上等同于下面的内容，因此只要是实现了Protocol Buffer`( 简称 Protobuf) 是 Google 公司内部的混合语言数据标准`的，就算不支持maps，仍然可以处理你的数据：  

```go
message MapFieldEntry {
  key_type key = 1;
  value_type value = 2;
}

repeated MapFieldEntry map_field = N;
```

好吧。。现在还用不上，等用到的时候再回头看一下。  

## 包  

你能在.proto文件中增加一个可选的包分类符来防止在协议消息类型`protocol message types`中产生同名的冲突。

```go
package foo.bar;
message Open { ... }
```
有点像命名空间吧。

接着当定义你的消息类型中的字段时你可以用包分类符：

```go
message Foo {
  ...
  foo.bar.Open open = 1;
  ...
}
```

包修饰符影响生成的代码的方式取决于你所选择的语言：  

* 在C++中，产生的classes会被包装在c++的命名空间内。例如，Open将在foo::bar空间中。`其实C++我不懂`  
* 在Java中，包会被当作Java包来使用，除非你在.proto文件中明确提供一个可选的java_package。  
* 在Python中，包的指令是被忽略的，因为根据其所在文件系统中的位置Python会自行组织模块(modules）。  
* 在Go中，包被作为Go的包名而使用，除非你在.proto文件中明确提供一个可选的go_package。  
* 在Ruby中……
* 在JavaNano中……
* 在C#中……

这三个不翻了，我没学过@_@
