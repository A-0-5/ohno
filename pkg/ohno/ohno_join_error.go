// Copyright © A.O.S, 2023.
// All Rights Reserved.
//
// author: A.O.S

package ohno

// This is structural representation of multiple errors in the same level. It
// is implemented as an array of errors
type OhNoJoinError struct {
	Errors []error `json:"errors" yaml:"errors"`
}

// This is the Error() method which satisfies the builtin [error] interface
// This prints the errors in the format
//
//	timestamp file:line(function): [code]name: description, message, extra
//	timestamp file:line(function): [code]name: description, message, extra
//	...
//
// [error]: https://pkg.go.dev/builtin#error
func (oj *OhNoJoinError) Error() string {
	var b []byte
	for i, err := range oj.Errors {
		if i > 0 {
			b = append(b, '\n')
		}
		b = append(b, err.Error()...)
	}

	return string(b)
}

// This method is an implementation to satisfy [errors.Unwrap] usage. It
// returns the underlying error array
func (oj *OhNoJoinError) Unwrap() []error {
	return oj.Errors
}
