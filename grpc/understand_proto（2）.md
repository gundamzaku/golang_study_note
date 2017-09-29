## 嵌套类型  

一看例子就明白  

```go
message SearchResponse {
  message Result {
    string url = 1;
    string title = 2;
    repeated string snippets = 3;
  }
  repeated Result results = 1;
}
```
message里面可以定义message  

引用的话就是：  

```go
message SomeOtherMessage {
  SearchResponse.Result result = 1;
}
```  

这样造出来的是什么样的怪物代码？  

```go
type SomeOtherMessage struct {
	Result *SearchResponse_Result `protobuf:"bytes,1,opt,name=result" json:"result,omitempty"`
}
```
一个结构体，里面一个指针字段。  

指向了  
```go
type SearchResponse_Result struct {
	Url      string   `protobuf:"bytes,1,opt,name=url" json:"url,omitempty"`
	Title    string   `protobuf:"bytes,2,opt,name=title" json:"title,omitempty"`
	Snippets []string `protobuf:"bytes,3,rep,name=snippets" json:"snippets,omitempty"`
}
```
SearchResponse也还在，也就一个数组指针定段。
```go
type SearchResponse struct {
	Results []*SearchResponse_Result `protobuf:"bytes,1,rep,name=results" json:"results,omitempty"`
}
```

还可以尝试嵌套，不过其原理是一样的。  
```go
message Outer {                  // Level 0
  message MiddleAA {  // Level 1
    message Inner {   // Level 2
      int64 ival = 1;
      bool  booly = 2;
    }
  }
  message MiddleBB {  // Level 1
    message Inner {   // Level 2
      int32 ival = 1;
      bool  booly = 2;
    }
  }
}
``` 

## 更新消息类型  
如果有一个message不再满足你所有的需求的时候，比如你要加字段了，那么记住下面几点：

1、不要改变已经存在的字段的数字标签  

2、假如你要加新字段，message被序列化成代码时，你的旧message格式仍然可以被你新生成的代码所解析。任何一个messages在你使用老的message格式进行新的生成，被序列化为代码的时候。你应该记住，新代码的元素的的默认值和旧代码生成的message可以完全交互。类似的，你新代码产生的消息可以被你旧代码所解析：当解析的时候，旧的二进制会简单地忽略新的字段。详细的可以参考未知字段部分。  

3、字段可被移除，只要标签号码没有在你更新后的消息类型中使用。你可以用一个重新命名的字段来替代，也许加一个后缀"OBSOLETE_"，或者保留标签，这样防止将来这个号码被意外地使用。  
int32, uint32, int64, uint64, and bool都是相容的，这意味着你能互相将他们进行改变，而不用breaking forwards（这个实在是不知道怎么翻了）或向后兼容。如果从wire（还是不知道怎么翻译，似乎是C++中的一种类型）中被解析的数字不能与对应的类型匹配，你将会如同C++中语法一样把数字转化成对应的那个类型（例：假如一个64位数字被作为一个int32位来读取，他将被转化为32位）  

4、sint32 和 sint64相互兼容但是和其它的数字类型不兼容（这是什么数字类型……）  

5、只要比特是有效的UTF-8格式，那他和字符也是相互兼容的。  

6、内嵌的消息体和包含消息体编码版本的比特类型相互兼容。  

7、fixed32和sfixed32相互兼容，fixed64和sfixed64也相互兼容。  

8、枚举与terms of wire(这似乎是C++中的术语）格式中的int32, uint32, int64和 uint64（注意如果他们不匹配其值将会被转化）。可是，要认识到当消息被反序 列化的时候，客户端代码处理它们是非常困难的：例如，未识别的proto3枚举类型将被保存在消息体中，但是当消息体被反序例化怎么表示它，仍然要依赖于相应的语言。Int字段将总是只保存它的值。  

翻完了，什么乱七八糟的，老实说我这个翻的人都没看懂。不管他了，直接跳过吧。  

## 未知的字段  
一串英文说明，似乎也不涉及到proto的语法，直接跳过，英文好的可以认真读一下，我是没有时间了。  

## 任意  

任意是一种类型：`google.protobuf.Any`  
```go
import "google/protobuf/any.proto";

message ErrorStatus {
  string message = 1;
  repeated google.protobuf.Any details = 2;
}
```
需要import一个google的库才能用。  
看说明是针对JAVA的pack() 和 unpack()，C++的PackFrom() 和 UnpackTo()的。记得这点就好，如果以后用到再说。  

## oneOf之一  

这又是一个大的模块
说明：如果你的消息体有很多字段并且在同一时间最多只有一个字段会被设置，你能强制这个行为，并且通过使用oneof特性来节省内存。  

看上去有点深奥。

### 使用oneOf  
```go
message SampleMessage {
  oneof test_oneof {
    string name = 4;
    SubMessage sub_message = 9;
  }
}
```
这个官方的例子不好，少了一段：  
```go
message SubMessage {
  string username = 1;
}
```

* 在oneof中，oneof定义中可以再套oneof，不过repeated 关键字是不能被使用的。  
在生成的代码中，oneof字段也有getters 和setters，如同常规的字段一样。 
你能得到一个特别的方法去检查在oneof中哪一个值被设定。其它的，自个儿去看你所选择的语言的对应的API参考吧。  

### oneOf的特性 

* 用了一个oneOf以后，会自动清除掉其它的oneOf成员。也就是你设置了多个oneOf，只有最后一个是有值有效的。
```go
SampleMessage message;
message.set_name("name");
CHECK(message.has_name());
message.mutable_sub_message();   // Will clear name field.
CHECK(!message.has_name());
```
这是官方给出的例子，只是一个通用的说明而已，如果你想确实的运行，这个是不能用的。以我现在学的Go为例，我需要到Go的专门的API页面去看一下：  
https://developers.google.com/protocol-buffers/docs/reference/go-generated#singular-scalar-proto3  

可以看到一个非常棒的，也确实可以用的例子：  
```go
package account;
message Profile {
  oneof avatar {
    string image_url = 1;
    bytes image_data = 2;
  }
}
```
编译，生成GO文件。  

我自己写了一个客户端的代码：  
```go
p := new(account.Profile)
p.Avatar = &account.Profile_ImageUrl{"http://example.com/image.png"}
fmt.Println(p.GetAvatar())
p.Avatar = &account.Profile_ImageData{nil}
fmt.Println(p.GetAvatar())
```
这个就算是Set了吧。意思大概就是一次只能设置一个字段？  

* 如果解析在wire中的同一个oneof中的多个成员，只有最后一个成员是被解析的消息体所使用的。（好吧，看到有wire就表示我不是很理解了）  
* oneof不能repeated。  
* 反射API在oneof字段下是可以工作的。  
* 如果你使用C++……算了我不翻了，反正我目前是不会用到C++的。  
* 还是在C++中……我还是不翻。  

### 向后兼容的问题 

在增加或移除oneof字段的时候一定要小心。如果检查下来oneof的值返回None/NOT_SET，这会意味着oneof将不能被设置或它已经在字段中被其它版本的oneof所设置。因此你将无法知道在wire中的未知字段是否是oneof的成员。  

### Tag重用的问题  

* <b>从oneof中移入或移出字段：</b>在message被序列化和解析的时候你将失去一些你的信息（一些字段将被清除）。
* <b>删除或增加字段：</b>这在message被序列化和解析的时候可能清除你当前设置oneof字段。
* <b>分割或合并字段：</b>这和常规字段的处理方式一样。

看（翻译）到这里，我也快吐了。下次再继续吧。  


