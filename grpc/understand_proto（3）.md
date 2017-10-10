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

### 包和名称的解析  

在protocol buffer语言中的类型名称解析像C++ 一样：首先，从最内部查找，然后至次内部`就是由内向外`，以此类推。每一个包都会被看成是其父类包的“内部（inner）”。一个“.”（例如，.foo.bar.Baz）意味着从外层范围开始向内。（有点搞，不过知道包的概念的话基本上也就这套路）

protocol buffer编译器会通过分析导入的.proto文件来解析所有的类型名称。各语言产生的代码都知道如何去访问该语言的各类型。即使在范围规则上有所不同。  

## 定义服务  

如果你想在RPC（Remote Procedure Call 远程过程调用）系统中使用你的消息类型，你能在.proto文件中定义RPC服务的接口，并且protocol buffer编译器将产生服务接口代码，同时在你所选择的语言中进行存根`这个解释起来有点拗口，摘一段网上的说明：存根类是一个类，它实现了一个接口，但是实现后的每个方法都是空的。 `。嗯……例如，如果你想定义一个RPC服务的方法使你能接受你的查询请求（SearchRequest）并返回查询响应，你能在你的.proto文件中定义它，如下：

```go
service SearchService {
  rpc Search (SearchRequest) returns (SearchResponse);
}
```
最易懂的使用protocol buffers的RPC系统是gRPC：由Google开发的开源性的，与平台无关`就是没有程序语言限制`的RPC系统。gRPC用protocol buffers的话工作的很有效率，能让你用特别的protocol buffer编译器插件直接从你的.proto文件产生对应的RPC代码。`这段就是王婆卖瓜`  

后面几段全是自己在吹逼的话，不翻了。要注意的是有一个第三方的项目，基于protocol buffers的一些扩展和现实，有兴趣可以看一下。  

https://github.com/google/protobuf/blob/master/docs/third_party.md  

## JSON映射  

Proto3支持规范化的JSON编码，使其在系统之中能更简单的分享数据。这些编码在下面的表中将一个一个地进行说明：  

<table border=1>
<tr><th>proto3</th><th>JSON</th><th width="25%">JSON example</th><th>Notes</th></tr>
<tr><td>message</td><td>object</td><td><code>{"fBar": v,
 "g": null,
 …}</code>
</td><td>Generates JSON objects. Message field names are mapped to lowerCamelCase and become JSON object keys. <code>null</code> is accepted and treated as the default value of the corresponding field type.</td></tr>
<tr><td>enum</td><td>string</td><td><code>"FOO_BAR"</code></td><td>The name of the enum value as specified in proto is used.</td></tr>
<tr><td>map&lt;K,V&gt;</td><td>object</td><td><code>{"k": v, …}</code></td><td>All keys are converted to strings.</td></tr>
<tr><td>repeated V</td><td>array</td><td><code>[v, …]</code></td><td><code>null</code> is accepted as the empty list [].</td></tr>
<tr><td>bool</td><td>true, false</td><td><code>true, false</code></td><td></td></tr>
<tr><td>string</td><td>string</td><td><code>"Hello World!"</code></td><td></td></tr>
<tr><td>bytes</td><td>base64 string</td><td><code>"YWJjMTIzIT8kKiYoKSctPUB+"</code></td><td>JSON value will be the data encoded as a string using standard base64 encoding with paddings. Either standard or URL-safe base64 encoding with/without paddings are accepted.</td></tr>
<tr><td>int32, fixed32, uint32</td><td>number</td><td><code>1, -10, 0</code></td><td>JSON value will be a decimal number. Either numbers or strings are accepted.</td></tr>
<tr><td>int64, fixed64, uint64</td><td>string</td><td><code>"1", "-10"</code></td><td>JSON value will be a decimal string. Either numbers or strings are accepted.</td></tr>
<tr><td>float, double</td><td>number</td><td><code>1.1, -10.0, 0, "NaN", "Infinity"</code></td><td>JSON value will be a number or one of the special string values "NaN", "Infinity", and "-Infinity". Either numbers or strings are accepted. Exponent notation is also accepted. </td></tr>
<tr><td>Any</td><td><code>object</code></td><td><code>{"@type": "url", "f": v, … }</code></td><td>If the Any contains a value that has a special JSON mapping, it will be converted as follows: <code>{"@type": xxx, "value": yyy}</code>. Otherwise, the value will be converted into a JSON object, and the <code>"@type"</code> field will be inserted to indicate the actual data type.</td></tr>
<tr><td>Timestamp</td><td>string</td><td><code>"1972-01-01T10:00:20.021Z"</code></td><td>Uses RFC 3339, where generated output will always be Z-normalized and uses 0, 3, 6 or 9 fractional digits.</td></tr>
<tr><td>Duration</td><td>string</td><td><code>"1.000340012s", "1s"</code></td><td>Generated output always contains 0, 3, 6, or 9 fractional digits, depending on required precision. Accepted are any fractional digits (also none) as long as they fit into nano-seconds precision.</td></tr>
<tr><td>Struct</td><td><code>object</code></td><td><code>{ … }</code></td><td>Any JSON object. See <code>struct.proto</code>.</td></tr><!--TODO: add link once we've figured out where we're putting doc for provided proto types-->
<tr><td>Wrapper types</td><td>various types</td><td><code>2, "2", "foo", true, "true", null, 0, …</code></td><td>Wrappers use the same representation in JSON as the wrapped primitive type, except that <code>null</code> is allowed and preserved during data conversion and transfer.</td></tr>
<tr><td>FieldMask</td><td>string</td><td><code>"f.fooBar,h"</code></td><td>See <code>fieldmask.proto</code>.</td></tr><!--TODO: add link once we've figured out where we're putting doc for provided proto types-->
<tr><td>ListValue</td><td>array</td><td><code>[foo, bar, …]</code></td></td><td></td></tr>
<tr><td>Value</td><td>value</td><td></td><td>Any JSON value</td></tr>
<tr><td>NullValue</td><td>null</td><td></td><td>JSON null</td></tr>
</table>
