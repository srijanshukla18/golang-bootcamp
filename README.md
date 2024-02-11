# golang-bootcamp

Raw dump of notes,


NEW TOPIC:
main.go too large not okay ofc
Packages in go

suppose helper functions,
helper.go, (in the same directory)
```golang
package main
import "strings"
func someFunc(){}
//same function, so its chill
//but now you can't run main.go directly, you also gotta run helper.go
//so you run
$ go run main.go helper.go
//OR, whole directory
$ go run .
```

you can have more packages
to, organize and logically group your code
each package gets its own folder

```golang
package helper
import (
"string"
)
func SomeFunc(){} // capitalized to EXPORT function
```
```golang
package main 
import (
"booking-app/helper" // read below about go.mod
)
helper.SomeFunc()
```

in your go.mod, 
first line was,
module booking-app

remember,
$ go run .
this includes packages folder 
so if the dir contains helper/whatever.go

NEW TOPIC:
3 levels of variable scope.
-package
-global
-local

local variable means,
can be inside a function or a code block like a for loop

package variable,
outside functions, in the beginning

gloabal variable,
same place as package variable, but capital
var MyVar = 0

NEW TOPIC:
maps
var userData = map[string]int // keys are string, vals are int
var mySlice = []string // slice of strings

//initialize empty map
var userData = make(map[string]int)

// intialize a list of maps
var bookings = make([]map[string]int, 0)
// this was list not slice, so we had to mention size
// we print it,
fmt.Printf("list of bookings is %v\n", bookings)

NEW TOPIC:
Struct
prefer struct when the attributes are certain. use maps when data can vary
User can be a struct. Predefined Structure.
Struct can have methods.
Structs can hold mixed data type
Struct are like classes. Kinda.
complex cases: struct containing a map or viceversa.

Ex:

// Struct define
type Userdata struct{
  firstName string
  lastName string
}

// Thats a new object
var userData = UserData {
  firstName: "Srijan",
  lastName: "Shukla,
}

var bookings = make([]UserData, 0)

NEW TOPIC:
look up fmt.printf vs fmt.Sprintf
time.Sleep(10*time.Second) // sleeps the main thread(goroutine)

new lightweight green thread
just add go keyword when calling the function
go myFunc()
thats it the abstraction takes care of making it run in the background.

Synchronizing the goroutines,
var wg = sync.WaitGroup{} // make the main thread wait to finish for a goroutine

like I mean, pseudocode
A()
B()
wg.Add(1) // we are saying there is one function that you need to wait for before exiting the main thread. Increases WaitGroup counter.
go C()
D()
E()
wg.Wait() // we are saying at the end of the main thread that hey you gotta be waiting for that one function remember? wait until WaitGroup counter is 0

Inside func C,
wg.Done() // decreases the WaitGroup counter by 1 

Channels are how goroutines talk to each other




NEW TOPIC:
func add(x int, y int) int {
	return x + y
}
// see how return type is mentioned?


func add(x, y int) int {
	return x + y
}
// x and y are of the same type

func swap(x, y string) (string, string) {
	return y, x
}
a, b := swap("hs", "asd")
// see how multiple are returned


func split(sum int) (x, y int) {
	x = sum * 4 / 9
	y = sum - x
	return
}
// see this? this is naked return. returns x and y
// only use in shorter functions

var c, python, java bool
// 3 booleans


// you can initialize things
var c, python, java = true, false, "no!"
var i, j int = 1, 2
c, python, java := true, false, "no!"


// TYPES:
bool
string
int  int8  int16  int32  int64
uint uint8 uint16 uint32 uint64 uintptr
byte // alias for uint8
rune // alias for int32
     // represents a Unicode code point
float32 float64
complex64 complex128


Variables declared without an explicit initial value are given their zero value.
The zero value is:
0 for numeric types,
false for the boolean type, and
"" (the empty string) for strings.

// type conversions,
i := 42
f := float64(i)
u := uint(f)




func (v *Vertex) Scale(f float64) {
	v.X = v.X * f
	v.Y = v.Y * f
}
//and then
v := Vertex{3, 4}
v.Scale(10)
// you see how you can do v.Scale, syntax for defining such a function is different


// v *Vertex args
// will want &v as input

//If a function is like this,
func ScaleFunc(v *Vertex, f float64) {}
//its called pointer receiver
//Then it means it can take a direct value or a pointer. it takes both.

var v Vertex
fmt.Println(v.Abs()) // OK
p := &v
fmt.Println(p.Abs()) // OK
//In this case, the method call p.Abs() is interpreted as (*p).Abs().

// Interfaces exist, 'type' implements a interface. Its implicit!
type I interface {
	M()
}
type T struct {
	S string
}
func (t T) M() {
	fmt.Println(t.S)
}