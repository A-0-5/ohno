// Copyright © A.O.S, 2023.
// All Rights Reserved.
//
// author: A.O.S

// This example demonstrates how to use the errors generated by [ohnogen] tool
// when generated without passing the -ohno flag. When you generate errors
// without the -ohno flag you will be able to use the enums as errors without
// being able to add additional context like a message, extra data, timestamp,
// source location etc. Lets begin by defining a custom type [MyFabulousError]
//
// [ohnogen]: https://pkg.go.dev/github.com/A-0-5/ohno/cmd/ohnogen
package usage_without_ohno

//go:generate go run ../../cmd/ohnogen/main.go -type=MyFabulousError -formatbase=16 -output=example_errors.go

// We first define a custom type like the one below
type MyFabulousError int

// Now here we define all the possible error values this package can return and
// run the command. (here formatbase=16 will make all the codes print in hex
// representation)
//
//	ohnogen -type=MyFabulousError -formatbase=16 -output=example_errors.go
//
// As this example purely concentrates on using the enums directly without
// depending on the ohno package the -ohno flag is omitted
const (
	NotFound      MyFabulousError = 100 + iota // I didn't find what you were looking for!
	AlreadyExists                              // I have this already!
	Internal                                   // Its not you, its me :(
	Unknown                                    // I don't know what happened
	Busy                                       // I'm busy rn, can we do this later?
	Unauthorised                               // You ain't got the creds to do this
	Fatal                                      // Help!!! Im dying!!!
)
