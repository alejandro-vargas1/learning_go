package main

import (
	"example.com/greetings"
	"fmt"
	"log"
)

func main() {
	// set the properties of the predefined Logger, including
	// the log entry prefix and a flag to disable printing
	// the time, source file and line number
	log.SetPrefix("greetings: ")
	log.SetFlags(0)

	// A slice of names
	names := []string{"Laura", "Alejandro", "Soledad"}

	// Request a greeting message for the names
	messages, err := greetings.Hellos(names)
	// If an error was returned, print it to the console and
	// exit the program
	if err != nil {
		log.Fatal(err)
	}

	// If no error was returned, print the returned map of
	// messages to the console
	fmt.Println(messages)
}

/*
Because it is local should be go mod edit -replace example.com/greetings=../greetings
This .. is because is like the UNIX shell to navigate throught directories.

After go mod tidy to search for needed packages

1. go mod edit -replace example.com/greetings=../greetings
2. go mod tidy


The command go list it's very usefull to get information about the package
and possible directories for the installation of the binary

The utilities are very usefull to know about the packages

The command to know where will be installed is "go list -f {{.Target}}"
*/
