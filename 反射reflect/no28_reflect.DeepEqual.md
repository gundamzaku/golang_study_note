## func DeepEqual(x, y interface{}) bool {}

DeepEqual是我所学习的reflect包中的最后一个对外函数了，看本体代码很短
```go
func DeepEqual(x, y interface{}) bool {
	if x == nil || y == nil {
		return x == y
	}
	v1 := ValueOf(x)
	v2 := ValueOf(y)
	if v1.Type() != v2.Type() {
		return false
	}
	return deepValueEqual(v1, v2, make(map[visit]bool), 0)
}
```
然而它的注释却比代码长得多得多。我一扫而过，只记住了一个单词“deeply equal”，深度相等。所以我也很清楚地明白这个方法的定义。  
粗看代码，首先要两个传入的参数不为空，接着两个参数的Type()必须相等，这样还不够，符合条件了。还要走deepValueEqual（）方法进行验证。  
这个deepValueEqual(v1, v2, make(map[visit]bool), 0)，不仅仅是把v1,v2两个参数传了进去，还代入了个map变量。  

deepValueEqual()就不同了，这才是整个方法的核心部分。  

首先，他还会再验证一次Type()
```go
if !v1.IsValid() || !v2.IsValid() {
	return v1.IsValid() == v2.IsValid()
}
if v1.Type() != v2.Type() {
	return false
}
```
此还，还根据v1.Kind()，针对不同的类型（Array、Slice、Interface、Ptr、Struct、Map、Func），还进行了递归的操作。  
至于代码部分，大至了解一下就行了。  

最后，我们还是以实际的代码来了解这个函数的用法，而官方的例子是非常好的参考了。  
```go
type Loop *Loop
type Loopy interface{}

var loop1, loop2 Loop
var loopy1, loopy2 Loopy

type NotBasic Basic
type Basic struct {
	x int
	y float32
}
type DeepEqualTest struct {
	a, b interface{}
	eq   bool
}
var (
	fn1 func()             // nil.
	fn2 func()             // nil.
	fn3 = func() { fn1() } // Not nil.
)
type self struct{}
func main(){

	var deepEqualTests = []DeepEqualTest{
		// Equalities
		{nil, nil, true},
		{1, 1, true},
		{int32(1), int32(1), true},
		{0.5, 0.5, true},
		{float32(0.5), float32(0.5), true},
		{"hello", "hello", true},
		{make([]int, 10), make([]int, 10), true},
		{&[3]int{1, 2, 3}, &[3]int{1, 2, 3}, true},
		{Basic{1, 0.5}, Basic{1, 0.5}, true},
		{error(nil), error(nil), true},
		{map[int]string{1: "one", 2: "two"}, map[int]string{2: "two", 1: "one"}, true},
		{fn1, fn2, true},

		// Inequalities
		{1, 2, false},
		{int32(1), int32(2), false},
		{0.5, 0.6, false},
		{float32(0.5), float32(0.6), false},
		{"hello", "hey", false},
		{make([]int, 10), make([]int, 11), false},
		{&[3]int{1, 2, 3}, &[3]int{1, 2, 4}, false},
		{Basic{1, 0.5}, Basic{1, 0.6}, false},
		{Basic{1, 0}, Basic{2, 0}, false},
		{map[int]string{1: "one", 3: "two"}, map[int]string{2: "two", 1: "one"}, false},
		{map[int]string{1: "one", 2: "txo"}, map[int]string{2: "two", 1: "one"}, false},
		{map[int]string{1: "one"}, map[int]string{2: "two", 1: "one"}, false},
		{map[int]string{2: "two", 1: "one"}, map[int]string{1: "one"}, false},
		{nil, 1, false},
		{1, nil, false},
		{fn1, fn3, false},
		{fn3, fn3, false},
		{[][]int{{1}}, [][]int{{2}}, false},
		{math.NaN(), math.NaN(), false},
		{&[1]float64{math.NaN()}, &[1]float64{math.NaN()}, false},
		{&[1]float64{math.NaN()}, self{}, true},
		{[]float64{math.NaN()}, []float64{math.NaN()}, false},
		{[]float64{math.NaN()}, self{}, true},
		{map[float64]float64{math.NaN(): 1}, map[float64]float64{1: 2}, false},
		{map[float64]float64{math.NaN(): 1}, self{}, true},

		// Nil vs empty: not the same.
		{[]int{}, []int(nil), false},
		{[]int{}, []int{}, true},
		{[]int(nil), []int(nil), true},
		{map[int]int{}, map[int]int(nil), false},
		{map[int]int{}, map[int]int{}, true},
		{map[int]int(nil), map[int]int(nil), true},

		// Mismatched types
		{1, 1.0, false},
		{int32(1), int64(1), false},
		{0.5, "hello", false},
		{[]int{1, 2, 3}, [3]int{1, 2, 3}, false},
		{&[3]interface{}{1, 2, 4}, &[3]interface{}{1, 2, "s"}, false},
		{Basic{1, 0.5}, NotBasic{1, 0.5}, false},
		{map[uint]string{1: "one", 2: "two"}, map[int]string{2: "two", 1: "one"}, false},

		// Possible loops.
		{&loop1, &loop1, true},
		{&loop1, &loop2, true},
		{&loopy1, &loopy1, true},
		{&loopy1, &loopy2, true},
	}


	for _, test := range deepEqualTests {
		if test.b == (self{}) {
			test.b = test.a
		}
		fmt.Printf("A:%s deepEqual B:%v Result:%v \n",test.a,test.b,reflect.DeepEqual(test.a, test.b))
	}
}
```
结果一目了然。

<完>
