// Copyright © A.O.S, 2023.
// All Rights Reserved.
//
// author: A.O.S

// package sourceinfo is a helper around runtime Caller which simplifies
// fetching source information using frames.
package sourceinfo // import "github.com/A-0-5/ohno/pkg/sourceinfo"

import (
	"path"
	"runtime"
	"strconv"
)

// This enum indicates the source information format
type SourceInfoType int

const (
	// Source information won't be retrieved
	NoSourceInfo SourceInfoType = iota
	// Full file name and line number
	FullFileAndLine
	// Full file name and line number with function name
	FullFileAndLineWithFunc
	// Short file name and line number
	ShortFileAndLine
	// Short file name and line number with function name
	ShortFileAndLineWithFunc
)

const (
	DefaultCallDepth int = 1
)

// This structure contains the information about the source code.
type SourceInformation struct {
	File     string `json:"file" yaml:"file"`
	Function string `json:"function,omitempty" yaml:"function,omitempty"`
	Line     int    `json:"line" yaml:"line"`
}

// This function prints the source information in the format
//
//	file:line(function):
func (s *SourceInformation) String() string {
	funcName := ""
	if s.Function != "" {
		funcName = " (" + s.Function + ")"
	}

	return s.File + ":" + strconv.Itoa(s.Line) + funcName + ":"
}

// This method gets the source information based on the type passed. Passing
// [sourceinfo.NoSourceInfo] will cause this function to return a nil pointer.
// If any error is encountered or if this function could not retrieve the
// source information then the returned pointer will be nil.
func GetSourceInformation(callDepth int, sourceInfoType SourceInfoType) (sourceInfo *SourceInformation) {
	if sourceInfoType == NoSourceInfo {
		return
	}

	if callDepth == 0 {
		callDepth = DefaultCallDepth
	}

	rpc := make([]uintptr, 1)
	n := runtime.Callers(callDepth+1, rpc[:])
	if n < 1 {
		return
	}
	frame, _ := runtime.CallersFrames(rpc).Next()
	if frame.PC == 0 {
		return
	}

	sourceInfo = &SourceInformation{
		File: frame.File,
		Line: frame.Line,
	}

	if sourceInfoType == ShortFileAndLine ||
		sourceInfoType == ShortFileAndLineWithFunc {
		sourceInfo.File = path.Base(frame.File)
	}

	if sourceInfoType == ShortFileAndLineWithFunc ||
		sourceInfoType == FullFileAndLineWithFunc {
		sourceInfo.Function = frame.Function
	}

	return
}
