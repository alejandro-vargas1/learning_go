# Generics in Golang

The generics are used to declare and used functions or types that are written to work any of a set of types provided by calling code.

The steps to show this are:

1. [[#Create the folder]]
2. [[#Add non-generic functions]]
3. [[#Add a generic function to handle multiple types]]
4. [[#Remove type arguments when calling the generic function]]
5. [[#Declare a type constraint]]


## Create the folder

Now it is necessary to create a folder named "generics"

```bash
mkdir generics
cd generics
```

Create the module with the next command

```bash 
go mod init example/generics
```


## Add non-generic functions

Will be add two functions that each add together the values to a map and return the total

We'll be using two different functions because we are working with two data types: int64 and float64

Create the main.go file, and at the top of the file specify the package

```go
package main
```

An standalone program will be in the main package, except for a library

Beneath the declaration of the package will be the two functions that will be using


```go
// SumInts adds together the values of m
func SumInts(m map[string]int64) int64 {  
    var s int64  
    for _, v := range m {  
       s += v  
    }  
    return s  
}  
  
// SumFloats adds together the values of m
func SumFloats(m map[string]float64) float64 {  
    var s float64  
    for _, v := range m {  
       s += v  
    }  
    return s  
}
```

This are the two functions to sum a map and the values in it, could be int or float depend the function. But it is necessary to create the main function to call this two functions

```go
func main() {  
    // Initialize a map for the integer values  
    ints := map[string]int64{  
       "first":  34,  
       "second": 12,  
    }  
  
    floats := map[string]float64{  
       "first":  35.98,  
       "second": 26.99,  
    }  
  
    fmt.Printf("Non-generic Sums: %v and %v\n",  
       SumInts(ints),  
       SumFloats(floats))  
}
```

With a generic function it is only needed one function instead of two

## Add a generic function to handle multiple types

With a generic function can be replaced the two functions above with only one.

It's needed to declare which types it supports, in addition to the function itself. The function should be called with the type arguments and the ordinary function arguments.


> Each type parameter has a type constraint that acts as a kind of meta-type for the type parameter. Each type constraint specifies the permissible type arguments that calling code can use for the respective type parameter.

\- [Golang documentation](https://go.dev/doc/tutorial/generics#:~:text=Each%20type%20parameter%20has%20a%20type%20constraint%20that%20acts%20as%20a%20kind%20of%20meta%2Dtype%20for%20the%20type%20parameter.%20Each%20type%20constraint%20specifies%20the%20permissible%20type%20arguments%20that%20calling%20code%20can%20use%20for%20the%20respective%20type%20parameter.)

The generic function will be:

```go
// SumIntsOrFloats sums the values of map m. It supports both int64 and float64
// as types for map values.
func SumIntsOrFloats[K comparable, V int64 | float64](m map[K]V) V{  
    var s V  
    for _, v := range m{  
       s += v  
    }  
    return s  
}
```

The next things to understand and consider about the function are:

- The function is declare with two type parameters (inside the square brackets []), K and V, and one argument that use the type parameters
- To the K parameter specify the type constraint "comparable". This allows that any type that can be used with the operand == or !=. ==This is necessary because Go requires that keys in maps are comparable==. 
- For V the type parameter it's a union of two types, int64 and float64, the "|" symbol put the constraint that can be any of those two types.
- When we specify the m argument of type "map[K]V", where K and V are already specified in the type parameters. K comparable it is necessary as the key.

The way to call this new generic function is:

```go
fmt.Printf("Generic Sums: %v and %v\n",   
    SumIntsOrFloats[string, int64](ints),   
    SumIntsOrFloats[string, float64](floats))
```


## Remove type arguments when calling the generic function

We are going to modify the original way to call the generic function, omitting this time to specify the data types of the map, because the compiler will infer the data types. This is not always possible, for example can not be made if the generic function does not have arguments, you will need to include the type arguments in the function call.

The function will look like this

```go
fmt.Printf("Generics Sums, type parameters inferd: %v and %v\n",   
    SumIntsOrFloats(ints),   
    SumIntsOrFloats(floats))
```

This time we are not going to call the function specifying the data types in the arguments.


## Declare a type constraint

In this part the constraints declared before will be move to it's own interface, so can be reused in multiples cases.

Can be declared a type constraint as an interface, the way to do it is

```go
type Number interface {  
    int64 | float64  
}
```

And this is going to be used in a new function which use this new type constraint, the function is

```go
// SunNumbers sums the values of map m. It supports both integers  
// and floats as map values
func SumNumbers[K comparable, V Number](m map[K]V) V {  
    var s V  
    for _, v := range m {  
       s += v  
    }  
    return s  
}
```

This function accept the type constraint in replace of the old function, but the way to call it is very similar. 