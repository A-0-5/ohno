// Copyright © A.O.S, 2023.
// All Rights Reserved.
//
// author: A.O.S

// This code is derived from https://github.com/golang/tools/blob/master/cmd/stringer/stringer.go licensed under BSD Clause 3. Copyright notice of original
// work and the license terms are below

// Copyright 2014 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Copyright (c) 2009 The Go Authors. All rights reserved.

// Redistribution and use in source and binary forms, with or without
// modification, are permitted provided that the following conditions are
// met:

//    * Redistributions of source code must retain the above copyright
// notice, this list of conditions and the following disclaimer.
//    * Redistributions in binary form must reproduce the above
// copyright notice, this list of conditions and the following disclaimer
// in the documentation and/or other materials provided with the
// distribution.
//    * Neither the name of Google Inc. nor the names of its
// contributors may be used to endorse or promote products derived from
// this software without specific prior written permission.

// THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS
// "AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT
// LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR
// A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT
// OWNER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL,
// SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT
// LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE,
// DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY
// THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
// (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
// OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.

// ------------------

// ohnogen is a golang [stringer] based tool to automate the creation of helper
// methods that provide the name and description of enums(ints only). This tool
// also automatically ensures all your enum types will satisfy the [error] interface.
//
// # Installation
//
// You can install this tool by running the command below:
//
//	go install github.com/A-0-5/ohno/cmd/ohnogen@latest
//
// # Usage
//
//	Usage of ohnogen:
//		ohnogen [flags] -type T [directory]
//		ohnogen [flags] -type T files... # Must be a single package
//
// # Flags
//
//	Flags:
//	  -formatbase int
//	    	format in which the enum value needs to be printed in different use cases.
//	    	Valid options are 2(binary), 8(octal),10(decimal), 16(hex).
//	    	default -formatbase=10 (default 10)
//	  -ohno
//	    	generate the OhNo method for using with ohno package
//	  -output string
//	    	output file name; default srcdir/<type>_errors.go
//	  -tags string
//	    	comma-separated list of build tags to apply
//	  -trimprefix prefix
//	    	trim the prefix from the generated constant names
//	  -type string
//	    	comma-separated list of type names; must be set
//
// # Usage Example
//
// Given an enum  of Type T , it generates the following methods
//
//		func (t T) String() string
//		func (t T) Description() string
//	 	func (t T) Error() string
//	 	func (t T) Package() string
//	 	func (t T) Code() string
//
//		// This function gets generated only if -ohno flag is set
//		func (MyError) OhNo(message string, extra any, cause error,
//			sourceInfoType sourceinfo.SourceInfoType, timestamp time.Time,
//			timestampLayout string) (ohnoError error)
//
// The file is created in the same package and directory as the package that defines T.
// It has helpful defaults designed for use with go generate.
//
// ohnogen works best with constants that are consecutive values such as created using iota,
// but creates good code regardless.
//
// For example, given this snippet,
//
//	package somepkg
//
//	type MyError int
//
//	const (
//	  NotFound      MyError    = iota // Requested resource was not found
//	  Timeout                         // Operation timed out
//	  AlreadyExists                   // Resource already exists
//	  Internal                        // An internal error occurred
//	  Unknown       = Internal        // An unknown error occurred
//	)
//
// running this command
//
//	ohnogen -type=MyError -ohno
//
// in the same directory will create the file myerror_errors.go, in package somepkg,
// containing a definition of
//
//	func (MyError) String() string
//	func (MyError) Description() string
//	func (MyError) Error() string
//	func (MyError) Package() string
//	func (MyError) Code() string
//
//	// This function gets generated only if -ohno flag is set
//	func (MyError) OhNo(message string, extra any, cause error,
//		sourceInfoType sourceinfo.SourceInfoType, timestamp time.Time,
//		timestampLayout string) (ohnoError error)
//
// # String() method
//
// The String() method will translate the value of a MyError constant to the string representation
// of the respective constant name, so that the call
//
//	fmt.Print(somepkg.NotFound)
//
// will print the string
//
//	NotFound
//
// # Description() method
//
// The Description() method will translate the comment of a MyError constant to
// the string representation of the respective constant, so that the call
//
//	fmt.Print(somepkg.NotFound.Description())
//
// will print the string
//
//	Requested resource was not found
//
// # Error() method
//
// The Error() method is implemented to satisfy the error interface and provide
// a single error string containing the name and description of the respective
// constant, so that the call
//
//	fmt.Print(somepkg.NotFound.Error())
//
// will print the string
//
//	NotFound: Requested resource was not found
//
// You can also use the this as go error as this satisfies the error interface
// and wrap, join and compare like you would with any [error].
//
// # Package() method
//
// The Package() method gives the package name where this error constant is
// defined. The call
//
//	fmt.Print(somepkg.NotFound.Package())
//
// will print the string
//
//	somepkg
//
// # Code() method
//
// The Code() method returns the enum value as a string formatted to the base
// specified in formatbase argument. Assuming formatbase was 16 at the time of
// generation,  the call
//
//	fmt.Print(somepkg.NotFound.Code())
//
// will print the string
//
//	0x0
//
// # OhNo(...) method
//
// The OhNo(...) method constructs an [github.com/A-0-5/ohno/pkg/ohno.OhNoError] from the [github.com/A-0-5/ohno/pkg/ohno] package and
// returns it. This is useful if you want to provide additional context to your
// error like a message and custom fields. The call
//
//	 fmt.Print(somepkg.NotFound.OhNo("not_found message", "extra_msg", nil,
//			sourceinfo.ShortFileAndLineWithFunc, time.Now(), time.DateTime))
//
// will print the string
//
//	2023-09-30 20:51:43 main.go:25 (main.main): [0x0] somepkg.NotFound: Requested resource was not found, not_found message, extra_msg
//
// # More Info
//
// Typically this process would be run using go generate, like this:
//
//	//go:generate ohnogen -type=MyError -ohno
//
// If multiple constants have the same value, the lexically first matching name will
// be used (in the example, Unknown will print as "Internal").
//
// With no arguments, it processes the package in the current directory.
// Otherwise, the arguments must name a single directory holding a Go package
// or a set of Go source files that represent a single Go package.
//
// The -type flag accepts a comma-separated list of types so a single run can
// generate methods for multiple types. The default output file is t_errors.go,
// where t is the lower-cased name of the first type listed. It can be overridden
// with the -output flag.
//
// # The `-ohno` Flag
//
// This is a special flag which when set generates the OhNo method which allows
// you to add additional context to the error like source information,
// timestamp, custom message etc. refer the [ohno] package for more details or
// refer [examples] to see how to use them
//
// # Examples
//
// You can use this tool in multiple ways. Checkout the [examples] part of this
// module to understand how you can use this tool and the package.
//
// [stringer]: https://pkg.go.dev/golang.org/x/tools/cmd/stringer
// [error]: https://pkg.go.dev/builtin#error
// [examples]: https://pkg.go.dev/github.com/A-0-5/ohno/examples
package main

import (
	"bytes"
	"flag"
	"fmt"
	"go/ast"
	"go/constant"
	"go/format"
	"go/token"
	"go/types"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"golang.org/x/tools/go/packages"
)

var (
	typeNames    = flag.String("type", "", "comma-separated list of type names; must be set")
	output       = flag.String("output", "", "output file name; default srcdir/<type>_errors.go")
	trimprefix   = flag.String("trimprefix", "", "trim the `prefix` from the generated constant names")
	ohnoFlag     = flag.Bool("ohno", false, "generate the OhNo method for using with ohno package")
	codeBaseFlag = flag.Int("formatbase", 10, "format in which the enum value needs to be printed in different use cases.\nValid options are 2(binary), 8(octal),10(decimal), 16(hex).\ndefault -formatbase=10")
	buildTags    = flag.String("tags", "", "comma-separated list of build tags to apply")
)

// Usage is a replacement usage function for the flags package.
func Usage() {
	fmt.Fprintf(os.Stderr, "Usage of ohnogen:\n")
	fmt.Fprintf(os.Stderr, "\tohnogen [flags] -type T [directory]\n")
	fmt.Fprintf(os.Stderr, "\tohnogen [flags] -type T files... # Must be a single package\n")
	fmt.Fprintf(os.Stderr, "For more information, see:\n")
	fmt.Fprintf(os.Stderr, "\thttps://pkg.go.dev/github.com/A-0-5/ohno/cmd/ohnogen\n")
	fmt.Fprintf(os.Stderr, "Flags:\n")
	flag.PrintDefaults()
}

func main() {
	log.SetFlags(0)
	log.SetPrefix("stringer: ")
	flag.Usage = Usage
	flag.Parse()
	if len(*typeNames) == 0 {
		flag.Usage()
		os.Exit(2)
	}

	codeBasePrefixString := ""
	switch *codeBaseFlag {
	case 2:
		codeBasePrefixString = "\"0b\" + "
	case 8:
		codeBasePrefixString = "\"0o\" + "
	case 10:
		codeBasePrefixString = ""
	case 16:
		codeBasePrefixString = "\"0x\" + "
	default:
		log.Fatalf("formatbase can only be one of 2,8,10,16 current value = %d", *codeBaseFlag)
	}

	types := strings.Split(*typeNames, ",")
	var tags []string
	if len(*buildTags) > 0 {
		tags = strings.Split(*buildTags, ",")
	}

	// We accept either one directory or a list of files. Which do we have?
	args := flag.Args()
	if len(args) == 0 {
		// Default: process whole package in current directory.
		args = []string{"."}
	}

	// Parse the package once.
	var dir string
	g := Generator{
		trimPrefix:     *trimprefix,
		lineComment:    true,
		ohnoEnable:     *ohnoFlag,
		codeBase:       *codeBaseFlag,
		codeBasePrefix: codeBasePrefixString,
	}

	// TODO(suzmue): accept other patterns for packages (directories, list of files, import paths, etc).
	if len(args) == 1 && isDirectory(args[0]) {
		dir = args[0]
	} else {
		if len(tags) != 0 {
			log.Fatal("-tags option applies only to directories, not when files are specified")
		}
		dir = filepath.Dir(args[0])
	}

	g.parsePackage(args, tags)

	// Print the header and package clause.
	g.Printf("// Code generated by \"ohnogen %s\"; DO NOT EDIT.\n", strings.Join(os.Args[1:], " "))
	g.Printf("\n")
	g.Printf("package %s", g.pkg.name)
	g.Printf("\n")
	if g.ohnoEnable {
		g.Printf("import (\n\t\"strconv\"\n\t\"time\"\n\t\"github.com/A-0-5/ohno/pkg/ohno\"\n\t\"github.com/A-0-5/ohno/pkg/sourceinfo\"\n)\n")
	} else {
		g.Printf("import \"strconv\"\n") // Used by all methods.
	}

	// Run generate for each type.
	for _, typeName := range types {
		g.generate(typeName)
	}

	// Format the output.
	src := g.format()

	// Write to file.
	outputName := *output
	if outputName == "" {
		baseName := fmt.Sprintf("%s_errors.go", types[0])
		outputName = filepath.Join(dir, strings.ToLower(baseName))
	}
	err := os.WriteFile(outputName, src, 0644)
	if err != nil {
		log.Fatalf("writing output: %s", err)
	}
}

// isDirectory reports whether the named file is a directory.
func isDirectory(name string) bool {
	info, err := os.Stat(name)
	if err != nil {
		log.Fatal(err)
	}
	return info.IsDir()
}

// Generator holds the state of the analysis. Primarily used to buffer
// the output for format.Source.
type Generator struct {
	buf bytes.Buffer // Accumulated output.
	pkg *Package     // Package we are scanning.

	trimPrefix     string
	lineComment    bool
	ohnoEnable     bool
	codeBase       int
	codeBasePrefix string

	logf func(format string, args ...interface{}) // test logging hook; nil when not testing
}

func (g *Generator) Printf(format string, args ...interface{}) {
	fmt.Fprintf(&g.buf, format, args...)
}

// File holds a single parsed file and associated data.
type File struct {
	pkg  *Package  // Package to which this file belongs.
	file *ast.File // Parsed AST.
	// These fields are reset for each type being generated.
	typeName string  // Name of the constant type.
	values   []Value // Accumulator for constant values of that type.

	trimPrefix  string
	lineComment bool
	ohnoEnable  bool
}

type Package struct {
	name  string
	defs  map[*ast.Ident]types.Object
	files []*File
}

// parsePackage analyzes the single package constructed from the patterns and tags.
// parsePackage exits if there is an error.
func (g *Generator) parsePackage(patterns []string, tags []string) {
	cfg := &packages.Config{
		Mode: packages.NeedName | packages.NeedTypes | packages.NeedTypesInfo | packages.NeedSyntax,
		// TODO: Need to think about constants in test files. Maybe write type_string_test.go
		// in a separate pass? For later.
		Tests:      false,
		BuildFlags: []string{fmt.Sprintf("-tags=%s", strings.Join(tags, " "))},
		Logf:       g.logf,
	}
	pkgs, err := packages.Load(cfg, patterns...)
	if err != nil {
		log.Fatal(err)
	}
	if len(pkgs) != 1 {
		log.Fatalf("error: %d packages matching %v", len(pkgs), strings.Join(patterns, " "))
	}
	g.addPackage(pkgs[0])
}

// addPackage adds a type checked Package and its syntax files to the generator.
func (g *Generator) addPackage(pkg *packages.Package) {
	g.pkg = &Package{
		name:  pkg.Name,
		defs:  pkg.TypesInfo.Defs,
		files: make([]*File, len(pkg.Syntax)),
	}

	for i, file := range pkg.Syntax {
		g.pkg.files[i] = &File{
			file:        file,
			pkg:         g.pkg,
			trimPrefix:  g.trimPrefix,
			lineComment: g.lineComment,
			ohnoEnable:  g.ohnoEnable,
		}
	}
}

// generate produces the String method for the named type.
func (g *Generator) generate(typeName string) {
	values := make([]Value, 0, 100)
	for _, file := range g.pkg.files {
		// Set the state for this run of the walker.
		file.typeName = typeName
		file.values = nil
		if file.file != nil {
			ast.Inspect(file.file, file.genDecl)
			values = append(values, file.values...)
		}
	}

	if len(values) == 0 {
		log.Fatalf("no values defined for type %s", typeName)
	}

	signed := values[0].signed
	// Generate code that will fail if the constants change value.
	g.Printf("func _() {\n")
	g.Printf("\t// An \"invalid array index\" compiler error signifies that the constant values have changed.\n")
	g.Printf("\t// Re-run the stringer command to generate them again.\n")
	g.Printf("\tvar x [1]struct{}\n")
	for _, v := range values {
		g.Printf("\t_ = x[%s - %s]\n", v.originalName, v.str)
	}
	g.Printf("}\n")
	runs := splitIntoRuns(values)
	// The decision of which pattern to use depends on the number of
	// runs in the numbers. If there's only one, it's easy. For more than
	// one, there's a tradeoff between complexity and size of the data
	// and code vs. the simplicity of a map. A map takes more space,
	// but so does the code. The decision here (crossover at 10) is
	// arbitrary, but considers that for large numbers of runs the cost
	// of the linear scan in the switch might become important, and
	// rather than use yet another algorithm such as binary search,
	// we punt and use a map. In any case, the likelihood of a map
	// being necessary for any realistic example other than bitmasks
	// is very low. And bitmasks probably deserve their own analysis,
	// to be done some other day.
	switch {
	case len(runs) == 1:
		g.buildOneRun(runs, typeName)
	case len(runs) <= 10:
		g.buildMultipleRuns(runs, typeName)
	default:
		g.buildMap(runs, typeName)
	}
	g.Printf("\n")
	g.Printf(errFunc, typeName)
	g.Printf("\n")
	g.Printf(pkgFunc, typeName, g.pkg.name)
	g.Printf("\n")

	formatInt := "FormatUint"
	typeCast := "uint64"
	if signed {
		formatInt = "FormatInt"
		typeCast = "int64"
	}

	g.Printf(codeFunc, typeName, formatInt, typeCast, g.codeBase, g.codeBasePrefix)
	if *ohnoFlag {
		g.Printf("\n")
		g.Printf(ohNoFunc, typeName, g.pkg.name)
	}
}

// splitIntoRuns breaks the values into runs of contiguous sequences.
// For example, given 1,2,3,5,6,7 it returns {1,2,3},{5,6,7}.
// The input slice is known to be non-empty.
func splitIntoRuns(values []Value) [][]Value {
	// We use stable sort so the lexically first name is chosen for equal elements.
	sort.Stable(byValue(values))
	// Remove duplicates. Stable sort has put the one we want to print first,
	// so use that one. The String method won't care about which named constant
	// was the argument, so the first name for the given value is the only one to keep.
	// We need to do this because identical values would cause the switch or map
	// to fail to compile.
	j := 1
	for i := 1; i < len(values); i++ {
		if values[i].value != values[i-1].value {
			values[j] = values[i]
			j++
		}
	}
	values = values[:j]
	runs := make([][]Value, 0, 10)
	for len(values) > 0 {
		// One contiguous sequence per outer loop.
		i := 1
		for i < len(values) && values[i].value == values[i-1].value+1 {
			i++
		}
		runs = append(runs, values[:i])
		values = values[i:]
	}
	return runs
}

// format returns the gofmt-ed contents of the Generator's buffer.
func (g *Generator) format() []byte {
	src, err := format.Source(g.buf.Bytes())
	if err != nil {
		// Should never happen, but can arise when developing this code.
		// The user can compile the output to see the error.
		log.Printf("warning: internal error: invalid Go generated: %s", err)
		log.Printf("warning: compile the package to analyze the error")
		return g.buf.Bytes()
	}
	return src
}

// Value represents a declared constant.
type Value struct {
	originalName string // The name of the constant.
	name         string // The name with trimmed prefix.
	// The value is stored as a bit pattern alone. The boolean tells us
	// whether to interpret it as an int64 or a uint64; the only place
	// this matters is when sorting.
	// Much of the time the str field is all we need; it is printed
	// by Value.String.
	value       uint64 // Will be converted to int64 when needed.
	signed      bool   // Whether the constant is a signed type.
	str         string // The string representation given by the "go/constant" package.
	description string
}

func (v *Value) String() string {
	return v.str
}

// byValue lets us sort the constants into increasing order.
// We take care in the Less method to sort in signed or unsigned order,
// as appropriate.
type byValue []Value

func (b byValue) Len() int      { return len(b) }
func (b byValue) Swap(i, j int) { b[i], b[j] = b[j], b[i] }
func (b byValue) Less(i, j int) bool {
	if b[i].signed {
		return int64(b[i].value) < int64(b[j].value)
	}
	return b[i].value < b[j].value
}

// genDecl processes one declaration clause.
func (f *File) genDecl(node ast.Node) bool {
	decl, ok := node.(*ast.GenDecl)
	if !ok || decl.Tok != token.CONST {
		// We only care about const declarations.
		return true
	}
	// The name of the type of the constants we are declaring.
	// Can change if this is a multi-element declaration.
	typ := ""
	// Loop over the elements of the declaration. Each element is a ValueSpec:
	// a list of names possibly followed by a type, possibly followed by values.
	// If the type and value are both missing, we carry down the type (and value,
	// but the "go/types" package takes care of that).
	for _, spec := range decl.Specs {
		vspec := spec.(*ast.ValueSpec) // Guaranteed to succeed as this is CONST.
		if vspec.Type == nil && len(vspec.Values) > 0 {
			// "X = 1". With no type but a value. If the constant is untyped,
			// skip this vspec and reset the remembered type.
			typ = ""

			// If this is a simple type conversion, remember the type.
			// We don't mind if this is actually a call; a qualified call won't
			// be matched (that will be SelectorExpr, not Ident), and only unusual
			// situations will result in a function call that appears to be
			// a type conversion.
			ce, ok := vspec.Values[0].(*ast.CallExpr)
			if !ok {
				continue
			}
			id, ok := ce.Fun.(*ast.Ident)
			if !ok {
				continue
			}
			typ = id.Name
		}
		if vspec.Type != nil {
			// "X T". We have a type. Remember it.
			ident, ok := vspec.Type.(*ast.Ident)
			if !ok {
				continue
			}
			typ = ident.Name
		}
		if typ != f.typeName {
			// This is not the type we're looking for.
			continue
		}
		// We now have a list of names (from one line of source code) all being
		// declared with the desired type.
		// Grab their names and actual values and store them in f.values.
		for _, name := range vspec.Names {
			if name.Name == "_" {
				continue
			}
			// This dance lets the type checker find the values for us. It's a
			// bit tricky: look up the object declared by the name, find its
			// types.Const, and extract its value.
			obj, ok := f.pkg.defs[name]
			if !ok {
				log.Fatalf("no value for constant %s", name)
			}
			info := obj.Type().Underlying().(*types.Basic).Info()
			if info&types.IsInteger == 0 {
				log.Fatalf("can't handle non-integer constant type %s", typ)
			}
			value := obj.(*types.Const).Val() // Guaranteed to succeed as this is CONST.
			if value.Kind() != constant.Int {
				log.Fatalf("can't happen: constant is not an integer %s", name)
			}
			i64, isInt := constant.Int64Val(value)
			u64, isUint := constant.Uint64Val(value)
			if !isInt && !isUint {
				log.Fatalf("internal error: value of %s is not an integer: %s", name, value.String())
			}
			if !isInt {
				u64 = uint64(i64)
			}
			v := Value{
				originalName: name.Name,
				value:        u64,
				signed:       info&types.IsUnsigned == 0,
				str:          value.String(),
			}
			if c := vspec.Comment; f.lineComment && c != nil && len(c.List) == 1 {
				v.description = strings.TrimSpace(c.Text())
			}

			v.name = strings.TrimPrefix(v.originalName, f.trimPrefix)

			f.values = append(f.values, v)
		}
	}
	return false
}

// Helpers

// usize returns the number of bits of the smallest unsigned integer
// type that will hold n. Used to create the smallest possible slice of
// integers to use as indexes into the concatenated strings.
func usize(n int) int {
	switch {
	case n < 1<<8:
		return 8
	case n < 1<<16:
		return 16
	default:
		// 2^32 is enough constants for anyone.
		return 32
	}
}

// declareIndexAndNameVars declares the index slices and concatenated names
// strings representing the runs of values.
func (g *Generator) declareIndexAndNameVars(runs [][]Value, typeName string) {
	var indexes, names []string
	var dIndexes, dNames []string
	for i, run := range runs {
		index, name := g.createIndexAndNameDecl(run, typeName, fmt.Sprintf("_%d", i))
		dIndex, dName := g.createIndexAndNameDescDecl(run, typeName, fmt.Sprintf("_%d", i))
		if len(run) != 1 {
			indexes = append(indexes, index)
			dIndexes = append(dIndexes, dIndex)
		}
		names = append(names, name)
		dNames = append(dNames, dName)
	}
	g.Printf("const (\n")
	for _, name := range names {
		g.Printf("\t%s\n", name)
	}
	for _, dName := range dNames {
		g.Printf("\t%s\n", dName)
	}
	g.Printf(")\n\n")

	if len(indexes) > 0 {
		g.Printf("var (")
		for _, index := range indexes {
			g.Printf("\t%s\n", index)
		}

		for _, dIndex := range dIndexes {
			g.Printf("\t%s\n", dIndex)
		}
		g.Printf(")\n\n")
	}
}

// declareIndexAndNameVar is the single-run version of declareIndexAndNameVars
func (g *Generator) declareIndexAndNameVar(run []Value, typeName string) {
	index, name := g.createIndexAndNameDecl(run, typeName, "")
	dIdx, dName := g.createIndexAndNameDescDecl(run, typeName, "")
	g.Printf("const (\n\t%s\n", name)
	g.Printf("\t%s\n", dName)
	g.Printf(")\n\n")
	g.Printf("var (\n\t%s\n", index)
	g.Printf("\t%s\n", dIdx)
	g.Printf(")\n\n")
}

// createIndexAndNameDecl returns the pair of declarations for the run. The caller will add "const" and "var".
func (g *Generator) createIndexAndNameDecl(run []Value, typeName string, suffix string) (string, string) {
	b := new(bytes.Buffer)
	indexes := make([]int, len(run))
	for i := range run {
		b.WriteString(run[i].name)
		indexes[i] = b.Len()
	}
	nameConst := fmt.Sprintf("_%s_name%s = %q", typeName, suffix, b.String())
	nameLen := b.Len()
	b.Reset()
	fmt.Fprintf(b, "_%s_index%s = [...]uint%d{0, ", typeName, suffix, usize(nameLen))
	for i, v := range indexes {
		if i > 0 {
			fmt.Fprintf(b, ", ")
		}
		fmt.Fprintf(b, "%d", v)
	}
	fmt.Fprintf(b, "}")
	return b.String(), nameConst
}

func (g *Generator) createIndexAndNameDescDecl(run []Value, typeName string, suffix string) (string, string) {
	b := new(bytes.Buffer)
	indexes := make([]int, len(run))
	for i := range run {
		b.WriteString(run[i].description)
		indexes[i] = b.Len()
	}
	nameConst := fmt.Sprintf("_%s_desc_name%s = %q", typeName, suffix, b.String())
	nameLen := b.Len()
	b.Reset()
	fmt.Fprintf(b, "_%s_desc_index%s = [...]uint%d{0, ", typeName, suffix, usize(nameLen))
	for i, v := range indexes {
		if i > 0 {
			fmt.Fprintf(b, ", ")
		}
		fmt.Fprintf(b, "%d", v)
	}
	fmt.Fprintf(b, "}")
	return b.String(), nameConst
}

// declareNameVars declares the concatenated names string representing all the values in the runs.
func (g *Generator) declareNameVars(runs [][]Value, typeName string, suffix string) {
	g.Printf("const (\n\t_%s_name%s = \"", typeName, suffix)
	for _, run := range runs {
		for i := range run {
			g.Printf("%s", run[i].name)
		}
	}
	g.Printf("\"\n")
	g.Printf("\t_%s_desc_name%s = \"", typeName, suffix)
	for _, run := range runs {
		for i := range run {
			g.Printf("%s", run[i].description)
		}
	}
	g.Printf("\"\n)\n\n")
}

// buildOneRun generates the variables and String method for a single run of contiguous values.
func (g *Generator) buildOneRun(runs [][]Value, typeName string) {
	values := runs[0]
	g.Printf("\n")
	g.declareIndexAndNameVar(values, typeName)
	// The generated code is simple enough to write as a Printf format.
	lessThanZero := ""
	if values[0].signed {
		lessThanZero = "i < 0 || "
	}
	if values[0].value == 0 { // Signed or unsigned, 0 is still 0.
		g.Printf(stringOneRun, typeName, usize(len(values)), lessThanZero)
		g.Printf("\n")
		g.Printf(stringOneRunDesc, typeName, usize(len(values)), lessThanZero)
	} else {
		g.Printf(stringOneRunWithOffset, typeName, values[0].String(), usize(len(values)), lessThanZero)
		g.Printf("\n")
		g.Printf(stringOneRunWithOffsetDesc, typeName, values[0].String(), usize(len(values)), lessThanZero)
	}
}

// Arguments to format are:
//
//	[1]: type name
//	[2]: size of index element (8 for uint8 etc.)
//	[3]: less than zero check (for signed types)
const stringOneRun = `// Returns the error name as string
func (i %[1]s) String() string {
	if %[3]si >= %[1]s(len(_%[1]s_index)-1) {
		return "%[1]s(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _%[1]s_name[_%[1]s_index[i]:_%[1]s_index[i+1]]
}
`

const stringOneRunDesc = `// Returns the description string
func (i %[1]s) Description() string {
	if %[3]si >= %[1]s(len(_%[1]s_desc_index)-1) {
		return "%[1]s(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _%[1]s_desc_name[_%[1]s_desc_index[i]:_%[1]s_desc_index[i+1]]
}
`

// Arguments to format are:
//	[1]: type name
//	[2]: lowest defined value for type, as a string
//	[3]: size of index element (8 for uint8 etc.)
//	[4]: less than zero check (for signed types)
/*
 */
const stringOneRunWithOffset = `// Returns the error name as string
func (i %[1]s) String() string {
	i -= %[2]s
	if %[4]si >= %[1]s(len(_%[1]s_index)-1) {
		return "%[1]s(" + strconv.FormatInt(int64(i + %[2]s), 10) + ")"
	}
	return _%[1]s_name[_%[1]s_index[i] : _%[1]s_index[i+1]]
}
`

const stringOneRunWithOffsetDesc = `// Returns the description string
func (i %[1]s) Description() string {
	i -= %[2]s
	if %[4]si >= %[1]s(len(_%[1]s_desc_index)-1) {
		return "%[1]s(" + strconv.FormatInt(int64(i + %[2]s), 10) + ")"
	}
	return _%[1]s_desc_name[_%[1]s_desc_index[i] : _%[1]s_desc_index[i+1]]
}
`

// buildMultipleRuns generates the variables and String method for multiple runs of contiguous values.
// For this pattern, a single Printf format won't do.
func (g *Generator) buildMultipleRuns(runs [][]Value, typeName string) {
	g.Printf("\n")
	g.declareIndexAndNameVars(runs, typeName)
	g.Printf("// Returns the error name as string\nfunc (i %s) String() string {\n", typeName)
	g.Printf("\tswitch {\n")
	for i, values := range runs {
		if len(values) == 1 {
			g.Printf("\tcase i == %s:\n", &values[0])
			g.Printf("\t\treturn _%s_name_%d\n", typeName, i)
			continue
		}
		if values[0].value == 0 && !values[0].signed {
			// For an unsigned lower bound of 0, "0 <= i" would be redundant.
			g.Printf("\tcase i <= %s:\n", &values[len(values)-1])
		} else {
			g.Printf("\tcase %s <= i && i <= %s:\n", &values[0], &values[len(values)-1])
		}
		if values[0].value != 0 {
			g.Printf("\t\ti -= %s\n", &values[0])
		}
		g.Printf("\t\treturn _%s_name_%d[_%s_index_%d[i]:_%s_index_%d[i+1]]\n",
			typeName, i, typeName, i, typeName, i)
	}
	g.Printf("\tdefault:\n")
	g.Printf("\t\treturn \"%s(\" + strconv.FormatInt(int64(i), 10) + \")\"\n", typeName)
	g.Printf("\t}\n")
	g.Printf("}\n\n")

	g.Printf("// Returns the description string\nfunc (i %s) Description() string {\n", typeName)
	g.Printf("\tswitch {\n")
	for i, values := range runs {
		if len(values) == 1 {
			g.Printf("\tcase i == %s:\n", &values[0])
			g.Printf("\t\treturn _%s_desc_name_%d\n", typeName, i)
			continue
		}
		if values[0].value == 0 && !values[0].signed {
			// For an unsigned lower bound of 0, "0 <= i" would be redundant.
			g.Printf("\tcase i <= %s:\n", &values[len(values)-1])
		} else {
			g.Printf("\tcase %s <= i && i <= %s:\n", &values[0], &values[len(values)-1])
		}
		if values[0].value != 0 {
			g.Printf("\t\ti -= %s\n", &values[0])
		}
		g.Printf("\t\treturn _%s_desc_name_%d[_%s_desc_index_%d[i]:_%s_desc_index_%d[i+1]]\n",
			typeName, i, typeName, i, typeName, i)
	}
	g.Printf("\tdefault:\n")
	g.Printf("\t\treturn \"%s(\" + strconv.FormatInt(int64(i), 10) + \")\"\n", typeName)
	g.Printf("\t}\n")
	g.Printf("}\n")
}

// buildMap handles the case where the space is so sparse a map is a reasonable fallback.
// It's a rare situation but has simple code.
func (g *Generator) buildMap(runs [][]Value, typeName string) {
	g.Printf("\n")
	g.declareNameVars(runs, typeName, "")
	g.Printf("\nvar (\n\t_%s_map = map[%s]string{\n", typeName, typeName)
	n := 0
	for _, values := range runs {
		for _, value := range values {
			g.Printf("\t\t%s: _%s_name[%d:%d],\n", &value, typeName, n, n+len(value.name))
			n += len(value.name)
		}
	}
	g.Printf("\t}\n\n")
	g.Printf("\t_%s_desc_map = map[%s]string{\n", typeName, typeName)
	n1 := 0
	for _, values := range runs {
		for _, value := range values {
			g.Printf("\t\t%s: _%s_desc_name[%d:%d],\n", &value, typeName, n1, n1+len(value.description))
			n1 += len(value.description)
		}
	}
	g.Printf("\t}\n")
	g.Printf(")\n\n")
	g.Printf(stringMap, typeName)
	g.Printf("\n")
	g.Printf(stringDescMap, typeName)
}

// Argument to format is the type name.
const stringMap = `// Returns the error name as string
func (i %[1]s) String() string {
	if str, ok := _%[1]s_map[i]; ok {
		return str
	}
	return "%[1]s(" + strconv.FormatInt(int64(i), 10) + ")"
}
`

// Argument to format is the type name.
const stringDescMap = `// Returns the description string
func (i %[1]s) Description() string {
	if str, ok := _%[1]s_desc_map[i]; ok {
		return str
	}
	return "%[1]s(" + strconv.FormatInt(int64(i), 10) + ")"
}
`

const errFunc = `// Returns the error's string representation
// [CODE]PACKAGE_NAME.ERROR_NAME: DESCRIPTION
func (i %[1]s) Error() string {
	return "[" + i.Code() + "]" + i.Package() + "." + i.String() + ": " + i.Description()
}
`

const pkgFunc = `// Returns the package name
func (i %[1]s) Package() string {
	return "%[2]s"
}
`

const codeFunc = `// Returns the integer code string as per the format base provided
func (i %[1]s) Code() string {
	return %[5]sstrconv.%[2]s(%[3]s(i), %[4]d)
}
`

const ohNoFunc = `// Generate a new error of [ohno.OhNoError] type with the data provided
// timestamp is optional, empty [timestampLayout] will assume default timestamp 
// of RFC3339Nano,  if you do not want source information to be captured pass 
// [sourceinfo.NoSourceInfo] for the sourceInfoType parameter.
//
// [timestampLayout]: https://pkg.go.dev/time#pkg-constants
func (i %[1]s) OhNo(message string, extra any, cause error, sourceInfoType sourceinfo.SourceInfoType, timestamp time.Time, timestampLayout string) (ohnoError error) {
	return ohno.New(i, message, extra, cause, sourceInfoType, sourceinfo.DefaultCallDepth+1, timestamp, timestampLayout)
}
`
