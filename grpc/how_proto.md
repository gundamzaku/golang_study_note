在这里我们先对proto有一个比较大致的了解。  

```go
syntax = "proto3";
```

在文档中有解释：  
The first line of the file specifies that you're using proto3 syntax: if you don't do this the protocol buffer compiler will assume you are using proto2

就当是初始化吧，如果不指定这个，默认会帮你指定成:proto2  

```go
option java_multiple_files = true;
option java_package = "io.grpc.examples.helloworld";
option java_outer_classname = "HelloWorldProto";
```
后面这三行，怎么看都是java的东西，我现在编译的是go环境，先跳过吧。我现在没有用java的proto插件，估计现在也不认。  

```go
package helloworld;
```
顾名思议，就是指定生成的文件的包名。  

```go
// The greeting service definition.
service Greeter {
  // Sends a greeting
  rpc SayHello (HelloRequest) returns (HelloReply) {}
}
```

这一段是定义一个服务。  
为了验证一下，我把Greeter改成Greeter2017  
可以看一下一共生成了几个文件：  

```go
type Greeter2017Client interface {
	// Sends a greeting
	SayHello(ctx context.Context, in *HelloRequest, opts ...grpc.CallOption) (*HelloReply, error)
}

type greeter2017Client struct {
	cc *grpc.ClientConn
}

func NewGreeter2017Client(cc *grpc.ClientConn) Greeter2017Client {
	return &greeter2017Client{cc}
}
func (c *greeter2017Client) SayHello(ctx context.Context, in *HelloRequest, opts ...grpc.CallOption) (*HelloReply, error) {
	out := new(HelloReply)
	err := grpc.Invoke(ctx, "/helloworld.Greeter2017/SayHello", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Server API for Greeter2017 service

type Greeter2017Server interface {
	// Sends a greeting
	SayHello(context.Context, *HelloRequest) (*HelloReply, error)
}

func RegisterGreeter2017Server(s *grpc.Server, srv Greeter2017Server) {
	s.RegisterService(&_Greeter2017_serviceDesc, srv)
}

func _Greeter2017_SayHello_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(HelloRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(Greeter2017Server).SayHello(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/helloworld.Greeter2017/SayHello",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(Greeter2017Server).SayHello(ctx, req.(*HelloRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _Greeter2017_serviceDesc = grpc.ServiceDesc{
	ServiceName: "helloworld.Greeter2017",
	HandlerType: (*Greeter2017Server)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "SayHello",
			Handler:    _Greeter2017_SayHello_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "helloworld.proto",
}
```

哇……这么多。都干嘛的？我不知道。  

在server端，有一段代码是调用这个的：`pb.RegisterGreeter2017Server(s, &server{})`  
注册一个服务，我们就当这个是默认要加的吧。  

内部的  
```
rpc SayHello (HelloRequest) returns (HelloReply) {}
```

表示定义了一个方法，传入HelloRequest，返回HelloReply  
而这两个，则在下面定义：  
```go
// The request message containing the user's name.
message HelloRequest {
  string name = 1;
}

// The response message containing the greetings
message HelloReply {
  string message = 1;
}
```

在Server端我们可以看到这个SayHello的方法被实现：  
```go
// SayHello implements helloworld.GreeterServer
func (s *server) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloReply, error) {
	return &pb.HelloReply{Message: "Hello " + in.Name}, nil
}
```
传入`in *pb.HelloRequest`,传出`*pb.HelloReply`  

在client端可以看到调用的代码：  
```go
r, err := c.SayHello(context.Background(), &pb.HelloRequest{Name: name})
```

`context.Background()`是默认的不用管它。 
name是一个字符串，如前面的定义`string name`  
我改一下代码：
```go
name = "dan"
r, err := c.SayHello(context.Background(), &pb.HelloRequest{Name: name})
```

输出就是：`2017/09/27 22:29:20 Greeting: Hello dan`  

原来如此，那现在意味着我也可以依样画葫芦来写一个了。  

比如我自己在service里面多定义一个sayBye的方法  

```go
service Greeter2017 {
  // Sends a greeting
  rpc SayHello (HelloRequest) returns (HelloReply) {}
  
  rpc SayBye (HelloRequest) returns (HelloReply) {}
}
```
编译一下。
```
protoc --go_out=plugins=grpc:.  helloworld.proto
```
会发现helloworld.pb.go里面多了一个方法：  
```go
func (c *greeter2017Client) SayBye(ctx context.Context, in *HelloRequest, opts ...grpc.CallOption) (*HelloReply, error){}
```

很顺利，接着在Server端里增加这个方法的实现：  

```go
// SayBye implements helloworld.GreeterServer
func (s *server) SayBye(ctx context.Context, in *pb.HelloRequest) (*pb.HelloReply, error) {
	return &pb.HelloReply{Message: "ByeBye " + in.Name}, nil
}
```

在client里面改一下  
```go
name = "dan"
r, err := c.SayBye(context.Background(), &pb.HelloRequest{Name: name})
```

执行后输出： 2017/09/27 22:39:52 Greeting: ByeBye dan  

好了，成功地自定义了一个方法，就是这么简单。
