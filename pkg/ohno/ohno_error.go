// Copyright © A.O.S, 2023.
// All Rights Reserved.
//
// author: A.O.S

package ohno

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/A-0-5/ohno/pkg/ohnoer"
	"github.com/A-0-5/ohno/pkg/sourceinfo"
)

const (
	separator      string = " "
	commaSeparator string = ", "
	newline        string = "\n"
	newlineTab     string = "\n-> "
)

// OhNoError is a structure which holds an error interface which satisfies the
// ohnoer.OhNoer interface.
type OhNoError struct {
	// An ohnoer.OhNoer interface error field
	ErrorCode error
	// A custom message for this instance of the error
	Message string
	// Any additional data to add context to this error
	Extra any
	// The error which led to this error being generated
	Cause error
	// File, Line & possibly Function name where this error was generated
	SourceInfo *sourceinfo.SourceInformation
	// Time at which this error occurred
	Timestamp time.Time
	// Layout in which the timestamp needs to be printed refer https://pkg.go.dev/time#pkg-constants
	TimestampLayout string
}

// This is the Error() method which satisfies the builtin [error] interface
// This prints the error in the format
//
//	timestamp file:line(function): [code]name: description, message, extra
//		cause(same representation as above with one indent)...
//
// [error]: https://pkg.go.dev/builtin#error
func (o *OhNoError) Error() string {

	var ob strings.Builder
	if !o.Timestamp.IsZero() {
		if o.TimestampLayout == "" {
			o.TimestampLayout = time.RFC3339Nano
		}
		ob.WriteString(o.Timestamp.Format(o.TimestampLayout))
		ob.WriteString(separator)
	}

	if o.SourceInfo != nil {
		ob.WriteString(o.SourceInfo.String())
		ob.WriteString(separator)
	}

	ob.WriteString(o.ErrorCode.Error())
	if o.Message != "" {
		ob.WriteString(commaSeparator)
		ob.WriteString(o.Message)
	}

	if o.Extra != nil {
		ob.WriteString(commaSeparator)
		fmt.Fprintf(&ob, "%+v", o.Extra)
	}

	if o.Cause != nil {
		ob.WriteString(newlineTab)
		ob.WriteString(o.Cause.Error())
	}

	return ob.String()

}

// This is a method implementation for usage with [errors.Is] in order to check
// if any errors in the chain match the current one
func (o *OhNoError) Is(target error) bool {
	if o.ErrorCode == nil {
		return false
	}

	return errors.Is(o.ErrorCode, target)
}

// This is a method implementation for usage with [errors.Unwrap] in order to
// get the wrapped/nested error
func (o *OhNoError) Unwrap() error {
	return o.Cause
}

// A simple yaml marshaler implementation for satisfying [yaml.Marshaler]
//
// [yaml.Marshaler]: https://pkg.go.dev/gopkg.in/yaml.v3#Marshaler
func (o *OhNoError) MarshalYAML() (interface{}, error) {
	if _, ok := o.ErrorCode.(ohnoer.OhNoer); !ok {
		return o.Error(), nil
	}

	marshalErr := o.marshalableError()
	return marshalErr, nil
}

// A simple json marshaler implementation for satisfying [encoding/json.Marshaler]
func (o *OhNoError) MarshalJSON() ([]byte, error) {
	if _, ok := o.ErrorCode.(ohnoer.OhNoer); !ok {
		return json.Marshal(o.Error())
	}

	marshalErr := o.marshalableError()
	return json.Marshal(marshalErr)
}

func (o *OhNoError) marshalableError() *ohNoMarshalError {
	marshalErr := new(ohNoMarshalError)

	ohnoer := o.ErrorCode.(ohnoer.OhNoer)
	marshalErr.Package = ohnoer.Package()
	marshalErr.Code = ohnoer.Code()
	marshalErr.Name = ohnoer.String()
	marshalErr.Description = ohnoer.Description()

	if !o.Timestamp.IsZero() {
		if o.TimestampLayout == "" {
			o.TimestampLayout = time.RFC3339Nano
		}

		marshalErr.TimeStamp = o.Timestamp.Format(o.TimestampLayout)
	}

	if o.SourceInfo != nil {
		marshalErr.SourceInfo = o.SourceInfo
	}

	if o.Cause != nil {
		marshalErr.CausedBy = o.Cause
	}

	return marshalErr
}

func (o *OhNoError) drillDownAndStackUp() error {
	errorList := []error{}
	errorList = o.removeCauseAndAppendRecursive(errorList)
	return &OhNoJoinError{
		Errors: errorList,
	}
}

func (o *OhNoError) removeCauseAndAppendRecursive(errs []error) []error {
	oNew := &OhNoError{
		ErrorCode:       o.ErrorCode,
		Message:         o.Message,
		Extra:           o.Extra,
		SourceInfo:      o.SourceInfo,
		Timestamp:       o.Timestamp,
		TimestampLayout: o.TimestampLayout,
	}

	errs = append(errs, oNew)
	if o.Cause == nil {
		return errs
	}

	if cause, ok := o.Cause.(*OhNoError); ok {
		errs = cause.removeCauseAndAppendRecursive(errs)
	} else if causes, ok := o.Cause.(*OhNoJoinError); ok {
		errs = append(errs, causes.Errors...)
	} else if causes, ok := o.Cause.(interface{ Unwrap() []error }); ok {
		errs = append(errs, causes.Unwrap()...)
	} else {
		errs = append(errs, o.Cause)
	}

	return errs
}

type ohNoMarshalError struct {
	Package     string                        `json:"package" yaml:"package"`
	Code        string                        `json:"code" yaml:"code"`
	Name        string                        `json:"name" yaml:"name"`
	Description string                        `json:"description" yaml:"description"`
	TimeStamp   string                        `json:"timestamp,omitempty" yaml:"timestamp,omitempty"`
	SourceInfo  *sourceinfo.SourceInformation `json:"source_information,omitempty" yaml:"source_information,omitempty"`
	CausedBy    error                         `json:"caused_by,omitempty" yaml:"caused_by,omitempty"`
}
