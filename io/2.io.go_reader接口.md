在io.go的代码段中，有一段都是用来定义接口的。  
```go
1、type Reader interface {}  
2、type Writer interface {}  
3、type Closer interface {}  
4、type Seeker interface {}  
5、type ReadWriter interface {}  
6、type ReadCloser interface {}  
7、type WriteCloser interface {}  
8、type ReadWriteCloser interface {}  
9、type ReadSeeker interface {}  
10、type WriteSeeker interface {}  
11、type ReadWriteSeeker interface {}  
12、type ReaderFrom interface {}  
13、type WriterTo interface {}  
14、type ReaderAt interface {}  
15、type WriterAt interface {}  
16、type ByteReader interface {}  
17、type ByteScanner interface {}  
18、type ByteWriter interface {}  
19、type RuneReader interface {}  
20、type RuneScanner interface {}  
21、type stringWriter interface {}  //私有
```
一共21个……这么多。目前来说，我并不知道他们的实际用处。  

接下来，又定义了几组对外的公有函数：   
```go
1、func WriteString(w Writer, s string) (n int, err error) {}  
2、func ReadAtLeast(r Reader, buf []byte, min int) (n int, err error) {}  
3、func ReadFull(r Reader, buf []byte) (n int, err error) {}  
4、func CopyN(dst Writer, src Reader, n int64) (written int64, err error) {}  
5、func Copy(dst Writer, src Reader) (written int64, err error) {}  
6、func CopyBuffer(dst Writer, src Reader, buf []byte) (written int64, err error) {}  
7、func copyBuffer(dst Writer, src Reader, buf []byte) (written int64, err error) {}  //私有
8、func LimitReader(r Reader, n int64) Reader { return &LimitedReader{r, n} }  
9、func TeeReader(r Reader, w Writer) Reader {}
```

最后，定义了几组对象：
```go
type LimitedReader struct {
	R Reader // underlying reader
	N int64  // max bytes remaining
}
func (l *LimitedReader) Read(p []byte) (n int, err error) {}
func NewSectionReader(r ReaderAt, off int64, n int64) *SectionReader {}  
```

```go
type SectionReader struct {
	r     ReaderAt
	base  int64
	off   int64
	limit int64
}
func (s *SectionReader) Read(p []byte) (n int, err error) {}
func (s *SectionReader) Seek(offset int64, whence int) (int64, error) {}
func (s *SectionReader) ReadAt(p []byte, off int64) (n int, err error) {}
func (s *SectionReader) Size() int64 { return s.limit - s.base }
```

```go
type teeReader struct { //私有
	r Reader
	w Writer
}
func (t *teeReader) Read(p []byte) (n int, err error) {
```

以上就是io.go的全貌了。那么怎么用起来呢？先从第一个方法看起吧。

```go
type Reader interface {
	Read(p []byte) (n int, err error)
}
```
第一个是接口，其实在io的操作里面，Reader和Write是必不可少的东西，可是如此一来，这个接口是在哪里实现的呢？  
在网上，我找到一个例子。  

```go
func main()  {
	data, err := ReadFrom(strings.NewReader("from string"), 12)
	fmt.Println(data)
	fmt.Println(err)
}

func ReadFrom(reader io.Reader, num int) ([]byte, error) {
	p := make([]byte, num)
	n, err := reader.Read(p)
	if n > 0 {
		return p[:n], nil
	}
	return p, err
}
result:
[102 114 111 109 32 115 116 114 105 110 103]
<nil>
```
可以看到，传入的是一个io.Reader类型，而这个类型经由strings.NewReader("from string")转化。  
于是我又转到string包中，发现了对Reader接口的实行。  
```go
func (r *Reader) Read(b []byte) (n int, err error) {
	if r.i >= int64(len(r.s)) {
		return 0, io.EOF
	}
	r.prevRune = -1
	n = copy(b, r.s[r.i:])
	r.i += int64(n)
	return
}
```
它是继承自
```go
// A Reader implements the io.Reader, io.ReaderAt, io.Seeker, io.WriterTo,
// io.ByteScanner, and io.RuneScanner interfaces by reading
// from a string.
type Reader struct {
	s        string
	i        int64 // current reading index
	prevRune int   // index of previous rune; or < 0
}
```
从注释上就可以看出，这个结构体又实现了io.Reader, io.ReaderAt, io.Seeker, io.WriterTo，io.ByteScanner, and io.RuneScanner这么多的接口。。  

先撸一遍代码的顺序  
step.1:`strings.NewReader("from string")`  
step.2:转入string包的reader.go  
step.3:返回`func NewReader(s string) *Reader { return &Reader{s, 0, -1} }`把Reader结构体（其实是对象）返回。  
step.4:`p := make([]byte, num)`初始化一个比特变量  
step.5:`n, err := reader.Read(p)`调用reader对象的Read方法，返回比特数据。  

总得来说，Reader这个接口是交给其它的方法去自行实现的接口，他本身不作任何数据处理。除了上面的读取字符串以外。还有：  
os.Stdin读取输入的流  
位于os包的file.go文件中 
```go
func (f *File) Read(b []byte) (n int, err error) {
	if err := f.checkValid("read"); err != nil {
		return 0, err
	}
	n, e := f.read(b)
	if n == 0 && len(b) > 0 && e == nil {
		return 0, io.EOF
	}
	if e != nil {
		err = &PathError{"read", f.name, e}
	}
	return n, err
}
```

还有位于读取文件：os.Open（同上）  

其实我最想了解的是网络io一块的内容，在net包中，不过暂时先放一下了，等下次读到net的时候再深入一下。  
