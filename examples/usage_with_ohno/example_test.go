// Copyright © A.O.S, 2023.
// All Rights Reserved.
//
// author: A.O.S

package usage_with_ohno_test

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/A-0-5/ohno/examples/usage_with_ohno"
	"github.com/A-0-5/ohno/pkg/ohno"
	"github.com/A-0-5/ohno/pkg/sourceinfo"
	"gopkg.in/yaml.v3"
)

// This is a function which returns an error of type [github.com/A-0-5/ohno/pkg/ohno.OhNoError]
func Foo() error {
	// Here we can use the enum Fatal to generate an OhNoError with additional context
	return usage_with_ohno.Fatal.OhNo(
		"something really bad happened",             // This is the message we want to add
		struct{ BadStuff string }{"some data here"}, // Some additional context
		nil,                                 // Since this is the root there is no cause error
		sourceinfo.ShortFileAndLineWithFunc, // File, line and Function format
		time.Unix(0, 0).UTC(),               // This error occurred on Jan 1 1970 00:00
		time.DateTime,                       // We just need the format as date & time
	)
}

// In this example we are simply defining a function Foo() which can return an
// error which needs to printed to the console.
func ExampleMyFabulousOhNoError_usage() {
	// We just print what ever foo returns
	fmt.Println(Foo().Error())
	// Output:
	// 1970-01-01 00:00:00 example_test.go:23 (github.com/A-0-5/ohno/examples/usage_with_ohno_test.Foo): [0x6a]usage_with_ohno.Fatal: Help!!! Im dying!!!, something really bad happened, {BadStuff:some data here}
}

// This example demonstrates how an OhNoError can be checked directly against
// the enum. Here Foo() returns an OhNoError which is being compared with the
// code usage_with_ohno.Fatal which will match. You can also check one
// OhNoError with another directly but in that case the whole error should match.
func ExampleMyFabulousOhNoError_check() {
	err := Foo()
	if err != nil {
		if errors.Is(err, usage_with_ohno.Fatal) {
			fmt.Printf("oh no its fatal!!!\n\t%s", err.Error())
		}

		if errors.Is(err, err) {
			fmt.Println()
			fmt.Println("they match !!!")
		}
	}

	// Output:
	// oh no its fatal!!!
	// 	1970-01-01 00:00:00 example_test.go:23 (github.com/A-0-5/ohno/examples/usage_with_ohno_test.Foo): [0x6a]usage_with_ohno.Fatal: Help!!! Im dying!!!, something really bad happened, {BadStuff:some data here}
	// they match !!!
}

// You can even directly use switch case here to check against multiple errors
// after converting it to OhNoError and using the embedded ErrorCode field.
func ExampleMyFabulousOhNoError_switch() {
	err := Foo()
	prefix := "i don't know what this is"
	if err != nil {
		var ohNoErr *ohno.OhNoError
		if errors.As(err, &ohNoErr) {
			switch ohNoErr.ErrorCode {
			case usage_with_ohno.AlreadyExists:
				prefix = "well its there already"
			case usage_with_ohno.NotFound:
				prefix = "nah I didn't find it"
			case usage_with_ohno.Busy:
				prefix = "not now!!"
			case usage_with_ohno.Fatal:
				prefix = "oh no its fatal!!!"
			}
		}

		fmt.Printf("%s\n\t%s", prefix, err.Error())
	}

	// Output:
	// oh no its fatal!!!
	// 	1970-01-01 00:00:00 example_test.go:23 (github.com/A-0-5/ohno/examples/usage_with_ohno_test.Foo): [0x6a]usage_with_ohno.Fatal: Help!!! Im dying!!!, something really bad happened, {BadStuff:some data here}
}

// By default the OhNo method allows you wrap an error inside an OhNoError when
// creating one. This example shows you how to do it. Typically if you get any
// error from your dependencies and you want to wrap them into the error you
// wish to emit then you would use this feature.
func ExampleMyFabulousOhNoError_wrapUnwrap() {
	// Let the error that Foo() returns be called fooErr
	fooErr := Foo()

	// Now we can wrap this inside another OhNoError. Let us use Unknown as a wrapper error for this.
	myError := usage_with_ohno.Unknown.OhNo(
		"Ive got no clue as to what happened", // message string
		12345,                                 // extra any - this can be anything , even nil if you don't want it to be there
		fooErr,                                // cause error - fooErr caused this error
		sourceinfo.NoSourceInfo,               //  sourceInfoType sourceinfo.SourceInfoType - I dont need the source information for this
		time.Unix(300, 0).UTC(),               // timestamp time.Time - 1 Jan 1970 00:05
		time.DateTime,                         // timestampLayout string - Date Time format
	)

	// Let us print this error
	fmt.Println(myError.Error())

	// Now lets unwrap this and see whats inside (we already know :P)

	causeOfMyError := errors.Unwrap(myError)

	// We already know that fooErr was wrapped inside myError, lets confirm
	if errors.Is(causeOfMyError, fooErr) {
		fmt.Println("causeOfMyError is fooErr")
	}

	// We can also see that its a Fatal error type
	if errors.Is(causeOfMyError, usage_with_ohno.Fatal) {
		fmt.Println("causeOfMyError is usage_with_ohno.Fatal")
	}

	// Output:
	// 1970-01-01 00:05:00 [0x67]usage_with_ohno.Unknown: I don't know what happened, Ive got no clue as to what happened, 12345
	// -> 1970-01-01 00:00:00 example_test.go:23 (github.com/A-0-5/ohno/examples/usage_with_ohno_test.Foo): [0x6a]usage_with_ohno.Fatal: Help!!! Im dying!!!, something really bad happened, {BadStuff:some data here}
	// causeOfMyError is fooErr
	// causeOfMyError is usage_with_ohno.Fatal
}

// You can use the default errors.Join or use ohno.Join method to join multiple
// errors. The advantage of using ohno.Join is that it is friendly with json
// and yaml marshaling
func ExampleMyFabulousOhNoError_join() {
	// Lets take whatever error foo returns
	fooErr := Foo()

	// Now lets create our own error
	myErr := usage_with_ohno.Internal.OhNo(
		"something went wrong",  // message string
		"some data",             // extra any
		nil,                     // cause error
		sourceinfo.NoSourceInfo, // sourceInfoType sourceinfo.SourceInfoType
		time.Unix(300, 0).UTC(), // timestamp time.Time
		time.DateTime,           // timestampLayout string
	)

	// When printing the errors will be printed in the same order you've sent
	// them to ohno.Join
	joinErr := ohno.Join(myErr, fooErr)

	fmt.Println(joinErr.Error())

	// Output:
	// 1970-01-01 00:05:00 [0x66]usage_with_ohno.Internal: Its not you, its me :(, something went wrong, some data
	// 1970-01-01 00:00:00 example_test.go:23 (github.com/A-0-5/ohno/examples/usage_with_ohno_test.Foo): [0x6a]usage_with_ohno.Fatal: Help!!! Im dying!!!, something really bad happened, {BadStuff:some data here}
}

// As a part of the structured representation goals of this package the
// OhNoError can be marshalled in json and yaml formats.
func ExampleMyFabulousOhNoError_marshal() {
	fooErr := Foo()

	fooJson, err := json.MarshalIndent(fooErr, "", "  ")
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	fooYaml, err := yaml.Marshal(fooErr)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	fmt.Println("json:")
	fmt.Println(string(fooJson))
	fmt.Println()
	fmt.Println("yaml:")
	fmt.Println(string(fooYaml))

	// Output:
	// json:
	// {
	//   "package": "usage_with_ohno",
	//   "code": "0x6a",
	//   "name": "Fatal",
	//   "description": "Help!!! Im dying!!!",
	//   "timestamp": "1970-01-01 00:00:00",
	//   "source_information": {
	//     "file": "example_test.go",
	//     "line": 23,
	//     "function": "github.com/A-0-5/ohno/examples/usage_with_ohno_test.Foo"
	//   }
	// }
	//
	// yaml:
	// package: usage_with_ohno
	// code: "0x6a"
	// name: Fatal
	// description: Help!!! Im dying!!!
	// timestamp: "1970-01-01 00:00:00"
	// source_information:
	//     file: example_test.go
	//     line: 23
	//     function: github.com/A-0-5/ohno/examples/usage_with_ohno_test.Foo
}

// There may be cases where you may get an error with deep nesting and when
// marshalling it into either text, json or yaml you may want to improve
// readability by flattening out the structure. To do that you can take any
// error of OhNoError type and then convert them OhNoJoinError type
func ExampleMyFabulousOhNoError_convertToJoinError() {
	// Lets begin with fooErr which we now know has no nested errors
	fooErr := Foo()

	// Lets wrap it with one other error
	barErr := usage_with_ohno.Internal.OhNo(
		"this error wraps fooErr",
		"level-1 nesting",
		fooErr, // This is where we wrap
		sourceinfo.NoSourceInfo,
		time.Unix(300, 0).UTC(),
		time.DateTime,
	)

	// Lets wrap barErr with one more error
	bazErr := usage_with_ohno.Unknown.OhNo(
		"this error wraps barErr",
		"level-2 nesting",
		barErr,
		sourceinfo.NoSourceInfo,
		time.Unix(600, 0).UTC(),
		time.DateTime,
	)

	flattenedErr := ohno.ConvertToJoinError(bazErr)

	bazJson, err := json.MarshalIndent(bazErr, "", "  ")
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	flattenedJson, err := json.MarshalIndent(flattenedErr, "", "  ")
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	fmt.Println("error is wrapped so all nested errors will be indented")
	fmt.Println(bazErr.Error())
	fmt.Println()
	fmt.Println("all nested errors are flattened, no indentation")
	fmt.Println(flattenedErr.Error())
	fmt.Println()
	fmt.Println("----")
	fmt.Println()
	fmt.Print("json representation of nested error")
	fmt.Println(string(bazJson))
	fmt.Println()
	fmt.Print("json representation of flattened error")
	fmt.Println(string(flattenedJson))

	// Output:
	// error is wrapped so all nested errors will be indented
	// 1970-01-01 00:10:00 [0x67]usage_with_ohno.Unknown: I don't know what happened, this error wraps barErr, level-2 nesting
	// -> 1970-01-01 00:05:00 [0x66]usage_with_ohno.Internal: Its not you, its me :(, this error wraps fooErr, level-1 nesting
	// -> 1970-01-01 00:00:00 example_test.go:23 (github.com/A-0-5/ohno/examples/usage_with_ohno_test.Foo): [0x6a]usage_with_ohno.Fatal: Help!!! Im dying!!!, something really bad happened, {BadStuff:some data here}
	//
	// all nested errors are flattened, no indentation
	// 1970-01-01 00:10:00 [0x67]usage_with_ohno.Unknown: I don't know what happened, this error wraps barErr, level-2 nesting
	// 1970-01-01 00:05:00 [0x66]usage_with_ohno.Internal: Its not you, its me :(, this error wraps fooErr, level-1 nesting
	// 1970-01-01 00:00:00 example_test.go:23 (github.com/A-0-5/ohno/examples/usage_with_ohno_test.Foo): [0x6a]usage_with_ohno.Fatal: Help!!! Im dying!!!, something really bad happened, {BadStuff:some data here}
	//
	// ----
	//
	// json representation of nested error{
	//   "package": "usage_with_ohno",
	//   "code": "0x67",
	//   "name": "Unknown",
	//   "description": "I don't know what happened",
	//   "timestamp": "1970-01-01 00:10:00",
	//   "caused_by": {
	//     "package": "usage_with_ohno",
	//     "code": "0x66",
	//     "name": "Internal",
	//     "description": "Its not you, its me :(",
	//     "timestamp": "1970-01-01 00:05:00",
	//     "caused_by": {
	//       "package": "usage_with_ohno",
	//       "code": "0x6a",
	//       "name": "Fatal",
	//       "description": "Help!!! Im dying!!!",
	//       "timestamp": "1970-01-01 00:00:00",
	//       "source_information": {
	//         "file": "example_test.go",
	//         "line": 23,
	//         "function": "github.com/A-0-5/ohno/examples/usage_with_ohno_test.Foo"
	//       }
	//     }
	//   }
	// }
	//
	// json representation of flattened error{
	//   "errors": [
	//     {
	//       "package": "usage_with_ohno",
	//       "code": "0x67",
	//       "name": "Unknown",
	//       "description": "I don't know what happened",
	//       "timestamp": "1970-01-01 00:10:00"
	//     },
	//     {
	//       "package": "usage_with_ohno",
	//       "code": "0x66",
	//       "name": "Internal",
	//       "description": "Its not you, its me :(",
	//       "timestamp": "1970-01-01 00:05:00"
	//     },
	//     {
	//       "package": "usage_with_ohno",
	//       "code": "0x6a",
	//       "name": "Fatal",
	//       "description": "Help!!! Im dying!!!",
	//       "timestamp": "1970-01-01 00:00:00",
	//       "source_information": {
	//         "file": "example_test.go",
	//         "line": 23,
	//         "function": "github.com/A-0-5/ohno/examples/usage_with_ohno_test.Foo"
	//       }
	//     }
	//   ]
	// }
}
