# Fuzzing

With fuzzing can generate random data to test our code searching for bugs or vulnerabilities. To see how this works we are going to create a fuzz test with a function, debug and fix the issues.

The steps are the next

1. [[#Create a folder for your code]]
2. [[#Add code to test]]
3. [[#Add a unit test]]
4. [[#Add a fuzz test]]
5. [[#Fix the invalid string error]]
6. [[#Fix the double reverse error]]

## Create a folder for your code

First it is needed to create the folder where we want to locate our project, in this path, we are going to write the next commands

```bash
mkdir fuzz
cd fuzz
```

Now create the go module to keep tracking the dependencies

```bash
go mod init example/fuzz
```

For production could be different the module path

## Add code to test

The function will be to reverse a string. For this create a file named "main.go", and at the top file specify the package

```go
package main
```

And beneath declare the Reverse function

```go
func Reverse(s string) string{  
    b := []byte(s)  
    for i, j := 0, len(b)-1; i < len(b)/2; i, j = i+1, j-1{  
       b[i], b[j] = b[j], b[i]  
    }  
    return string(b)  
}
```

==This function converts to a slice of bytes because it is mutable, and the string not==.

Now declare the main function to call the reverse function

```go
func main() {  
    input := "The quick brown fox jumped over the lazy dog"  
    rev := Reverse(input)  
    doubleRev := Reverse(rev)  
    fmt.Printf("origin: %q\n", input)  
    fmt.Printf("Reversed: %q\n", rev)  
    fmt.Printf("Reversed again: %q\n", doubleRev)  
}
```

## Add a unit test

Add a unit test, to test the Reverse function

First create a new file named reverse_test.go

```go
package main  
  
import (  
    "testing"  
)  
  
func TestReverse(t *testing.T) {  
    testcases := []struct {  
       in, want string  
    }{  
       {"Hello, world", "dlrow ,olleH"},  
       {" ", " "},  
       {"!12345", "54321!"},  
    }  
    for _, tc := range testcases {  
       rev := Reverse(tc.in)  
       if rev != tc.want {  
          t.Errorf("Reverse: %q, want %q\n", rev, tc.want)  
       }  
    }  
}
```

With  this go file we can call our Reverse function, which has the first upper letter, making it accessible in all the main package.

To run this is with

```bash
go test
```

## Add a fuzz test

Now add a new function to work as fuzzy, to find new scenarios to the testing

```go
func FuzzReverse(f *testing.F){  
    testcases := []string{"Hello, world", " ", "!12345"}  
    for _, tc := range testcases{  
       f.Add(tc) // Use f.add to provide a seed corpus  
    }  
    f.Fuzz(func(t *testing.T, orig string) {  
       rev := Reverse(orig)  
       doubleRev := Reverse(rev)  
       if orig != doubleRev{  
          t.Errorf("Before: %q, after: %q\n", orig, doubleRev)  
       }  
       if utf8.ValidString(orig) && !utf8.ValidString(rev) {  
          t.Errorf("Reverse produced invalid string UTF-8 string %q\n", rev)  
       }  
    })  
}
```

There is something to keep in mind that the name of a unit test starts with TestXXX and a fuzz test starts with FuzzXXX, and takes as argument *testing.F instead of *testing.T. And it uses t.Fuzz instead of t.Run.

To run a fuzz test we modify a little the command to run the tests.

```bash
go test -fuzz=Fuzz
```

This test is going to return a failure, and to watch which is the problem we can see it in the ==testdata/fuzz/FuzzReverse==. And should appear the next content in the file

```bash
go test fuzz v1  
string("ч")
```

- The first line indicates the encoding version
- Each of the following lines indicates the type of the corpus entry

Run the test again without -fuzz flag, and now should be use the corpus that use saw before.


## Fix the invalid string error

Now we are going to follow a number of steps to analyze and fix the errors.

### Diagnose the error

It is necessary to debug the code to find which is the error. 

Based on the docs of the utf8.ValidString it says

> ValidString reports whether s consists entirely of valid UTF-8-encoded runes.

\- [utf8.ValidString documentation](https://pkg.go.dev/unicode/utf8#:~:text=ValidString%20reports%20whether%20s%20consists%20entirely%20of%20valid%20UTF%2D8%2Dencoded%20runes.)

With our current function the reverse is make byte-by-byte, but to solve this it is necessary to do it rune-by-rune

- **Byte**: A byte takes 8 bit storage, it is equivalent to uint8 data type.
- **Rune**: It takes 4 bytes or 32 bits of storage, it is equivalent to int32. Are used to represent the unicode characters, that are broader than ASCII. This are encoded in UTF-8 format.

To watch more about it is with

[[Runes vs bytes#^3cd654]]


#### Write the code

Replace the fuzz target within the FuzzReverse function with the next correction

```go
f.Fuzz(func(t *testing.T, orig string) {  
    rev := Reverse(orig)  
    doubleRev := Reverse(rev)  
    t.Logf("Number of runes: orig=%d, rev=%d, doubleRev=%d", utf8.RuneCountInString(orig), utf8.RuneCountInString(rev), utf8.RuneCountInString(doubleRev))  
    if orig != doubleRev {  
       t.Errorf("Before: %q, after: %q\n", orig, doubleRev)  
    }  
    if utf8.ValidString(orig) && !utf8.ValidString(rev) {  
       t.Errorf("Reverse produced invalid string UTF-8 string %q\n", rev)  
    }  
})
```

The only change was to add a new line to show a log with the error.

==The t.logf will print in the console if an error occurs, or if we execute the test with -v flag==

### Fix the error

To correct this error there is needed to work with runes instead of bytes, so we will change the Reverse() function

```go
func Reverse(s string) string {  
    r := []rune(s)  
    for i, j := 0, len(r)-1; i < len(r)/2; i, j = i+1, j-1 {  
       r[i], r[j] = r[j], r[i]  
    }  
    return string(r)  
}
```

The only change that was made is the part that we use an slice of runes instead of bytes.

With this solution we pass almost initially the bugs, but when we run again the test with fuzz there is an error. 


## Fix the double reverse error

In this section will be fixed the error in the double reverse

### Diagnose the error

This could be made with a debugger, but this time we'll going to use the log.

When we replace the slice of bytes with a slice of runes, there is a problem that ==the string accept bytes that are not valid to UTF-8==, but when we make it with runes, ==Go will encode to valid UTF-8 character==, making this change the original with the result.

Make some changes in the Reverse() function to see the difference between the bytes and the runes.

```go
func Reverse(s string) string {  
    fmt.Printf("input: %q\n", s)  
    r := []rune(s)  
    fmt.Printf("runes: %q\n", r)  
    for i, j := 0, len(r)-1; i < len(r)/2; i, j = i+1, j-1 {  
       r[i], r[j] = r[j], r[i]  
    }  
    return string(r)  
}
```

The changes let us know which are the difference between the string and the runes

And to check this we run our test to verify what is happening. For this we are going to select the testdata that show us the error

```
└── testdata
    └── fuzz
        └── FuzzReverse
            ├── 87c897a306ca8a64
            └── ba0160f680fd7fcb
```

And we run the fuzz test with the next command. The hash will be different in all the computers. This is to run exactly the test that we want to check

```bash
go test -run=FuzzReverse/ba0160f680fd7fcb
```

```bash
input: "\xd0"
runes: ['�']
input: "�"
runes: ['�']
--- FAIL: FuzzReverse (0.00s)
    --- FAIL: FuzzReverse/ba0160f680fd7fcb (0.00s)
        reverse_test.go:32: Number of runes: orig=1, rev=1, doubleRev=1
        reverse_test.go:34: Before: "\xd0", after: "�"
FAIL
exit status 1
FAIL    example/fuzz    0.001s
```

As we can see the input at the beginning and the second line are different. The first one is an invalid unicode

### Fix the error

In the Reverse() function we are going to constraint, to verify if the input is not a valid UTF-8

```go
func Reverse(s string) (string, error) {  
    if !utf8.ValidString(s){  
       return s, errors.New("Input is not valid UTF-8")  
    }  
    r := []rune(s)  
    for i, j := 0, len(r)-1; i < len(r)/2; i, j = i+1, j-1 {  
       r[i], r[j] = r[j], r[i]  
    }  
    return string(r), nil  
}
```

It verifies if the input is a valid string, but because now we return more than one value it is necessary to fix the main function at the moment to call it.

```go
func main() {  
    input := "The quick brown fox jumped over the lazy dog"  
    rev, revErr := Reverse(input)  
    doubleRev, doubleRevErr := Reverse(rev)  
    fmt.Printf("origin: %q\n", input)  
    fmt.Printf("Reversed: %q, err: %v\n", rev, revErr)  
    fmt.Printf("Reversed again: %q, err: %v\n", doubleRev, doubleRevErr)  
}
```

And at the tests need to be adapted to call the function with the two parameters and return when there is a problem.

```go
package main  
  
import (  
    "testing"  
    "unicode/utf8")  
  
func TestReverse(t *testing.T) {  
    testcases := []struct {  
       in, want string  
    }{  
       {"Hello, world", "dlrow ,olleH"},  
       {" ", " "},  
       {"!12345", "54321!"},  
    }  
    for _, tc := range testcases {  
       rev, err := Reverse(tc.in)  
       if err != nil {  
          return  
       }  
       if rev != tc.want {  
          t.Errorf("Reverse: %q, want %q\n", rev, tc.want)  
       }  
    }  
}  
  
func FuzzReverse(f *testing.F) {  
    testcases := []string{"Hello, world", " ", "!12345"}  
    for _, tc := range testcases {  
       f.Add(tc) // Use f.add to provide a seed corpus  
    }  
    f.Fuzz(func(t *testing.T, orig string) {  
       rev, err1 := Reverse(orig)  
       if err1 != nil {  
          return  
       }  
       doubleRev, err2 := Reverse(rev)  
       if err2 != nil {  
          return  
       }  
       t.Logf("Number of runes: orig=%d, rev=%d, doubleRev=%d", utf8.RuneCountInString(orig), utf8.RuneCountInString(rev), utf8.RuneCountInString(doubleRev))  
       if orig != doubleRev {  
          t.Errorf("Before: %q, after: %q\n", orig, doubleRev)  
       }  
       if utf8.ValidString(orig) && !utf8.ValidString(rev) {  
          t.Errorf("Reverse produced invalid string UTF-8 string %q\n", rev)  
       }  
    })  
}
```

And now run the test with the fuzz flag again, but with a time limit to avoid an infinite loop.

```bash
go test -fuzz=Fuzz -fuzztime 30s
```

With this command we can test with the fuzz. 

Now it does not return an error and pass the tests