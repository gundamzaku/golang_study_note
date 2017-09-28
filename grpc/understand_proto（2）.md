### 嵌套类型  

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
