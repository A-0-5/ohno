// Copyright © A.O.S, 2023.
// All Rights Reserved.
//
// author: A.O.S

package usage_without_ohno_test

import (
	"errors"
	"fmt"

	"github.com/A-0-5/ohno/examples/usage_without_ohno"
)

// In this example we are simply defining a function Foo() which can return an
// error which needs to printed to the console.
func ExampleMyFabulousError_usage() {
	// We just print what ever foo returns
	fmt.Println(Foo().Error())
	// Output:
	// [0x6a]usage_without_ohno.Fatal: Help!!! Im dying!!!
}

// With the error generated using ohnogen you can use the errors
// package method to check if an error matches a specific error
func ExampleMyFabulousError_check() {
	err := Foo()
	if err != nil {
		if errors.Is(err, usage_without_ohno.Fatal) {
			fmt.Printf("oh no its fatal!!!\n\t%s", err.Error())
		}
	}

	// Output:
	// oh no its fatal!!!
	// 	[0x6a]usage_without_ohno.Fatal: Help!!! Im dying!!!
}

// You can even directly use switch case here to check against multiple errors
func ExampleMyFabulousError_switch() {
	err := Foo()
	prefix := ""
	if err != nil {
		switch err {
		case usage_without_ohno.AlreadyExists:
			prefix = "well its there already"
		case usage_without_ohno.NotFound:
			prefix = "nah I didn't find it"
		case usage_without_ohno.Busy:
			prefix = "not now!!"
		case usage_without_ohno.Fatal:
			prefix = "oh no its fatal!!!"
		default:
			prefix = "i don't know what this is"
		}

		fmt.Printf("%s\n\t%s", prefix, err.Error())
	}

	// Output:
	// oh no its fatal!!!
	// 	[0x6a]usage_without_ohno.Fatal: Help!!! Im dying!!!
}

// You can wrap these errors like any other error using fmt.Errorf
func ExampleMyFabulousError_wrapUnwrap() {
	err := fmt.Errorf("Ive wrapped\n\t%w\nin this error", usage_without_ohno.Busy)
	fmt.Println(err.Error())
	fmt.Println()
	busyError := errors.Unwrap(err)
	fmt.Println("after unwrapping it we get")
	fmt.Println(busyError.Error())
	// Output:
	// Ive wrapped
	// 	[0x68]usage_without_ohno.Busy: I'm busy rn, can we do this later?
	// in this error
	//
	// after unwrapping it we get
	// [0x68]usage_without_ohno.Busy: I'm busy rn, can we do this later?
}

// You can also join multiple errors using the errors.Join method
func ExampleMyFabulousError_join() {
	// Here we join 4 errors
	err := errors.Join(
		usage_without_ohno.AlreadyExists,
		usage_without_ohno.Busy,
		usage_without_ohno.Internal,
		usage_without_ohno.NotFound,
	)

	// When we print the output will be as follows
	fmt.Println(err.Error())

	// Output:
	// [0x65]usage_without_ohno.AlreadyExists: I have this already!
	// [0x68]usage_without_ohno.Busy: I'm busy rn, can we do this later?
	// [0x66]usage_without_ohno.Internal: Its not you, its me :(
	// [0x64]usage_without_ohno.NotFound: I didn't find what you were looking for!
}

// This is a function which returns an error of type [MyFabulousError]
func Foo() error {
	return usage_without_ohno.Fatal
}
