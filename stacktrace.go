package cockroachdberrors

import (
	"slices"

	"github.com/cockroachdb/errors"
	"github.com/cockroachdb/errors/errbase"
)

type state struct {
	b []byte
}

// Write implement fmt.Formatter interface.
func (s *state) Write(b []byte) (n int, err error) {
	s.b = b
	return len(b), nil
}

// Width implement fmt.Formatter interface.
func (s *state) Width() (wid int, ok bool) {
	return 0, false
}

// Precision implement fmt.Formatter interface.
func (s *state) Precision() (prec int, ok bool) {
	return 0, false
}

// Flag implement fmt.Formatter interface.
func (s *state) Flag(c int) bool {
	return false
}

type Frame struct {
	Source   string `json:"source"`
	Line     string `json:"line"`
	Function string `json:"func"`
}

type Stack struct {
	Details    []string `json:"details,omitempty"`
	Stacktrace []Frame  `json:"stacktrace,omitempty"`
}

func frameField(f errbase.StackFrame, c rune) string {
	s := &state{}
	f.Format(s, c)

	return string(s.b)
}

func marshalStack(err error) []Stack {
	var stacks []Stack
	if inner := errors.UnwrapOnce(err); inner != nil {
		stacks = marshalStack(inner)
	} else {
		stacks = make([]Stack, 0)
	}

	var stacktrace []Frame
	var details []string

	if e, ok := err.(errbase.StackTraceProvider); ok {
		st := e.StackTrace()
		stacktrace = make([]Frame, 0, len(st))

		for _, f := range st {
			stacktrace = append(stacktrace, Frame{
				Source:   frameField(f, 's'),
				Line:     frameField(f, 'd'),
				Function: frameField(f, 'n'),
			})
		}
	} else {
		details = errors.GetSafeDetails(err).SafeDetails
	}

	stacks = append(stacks, Stack{
		Details:    details,
		Stacktrace: stacktrace,
	})

	return stacks
}

// MarshalStack implements cockroachdb/errors stack trace marshaling.
//
// zerolog.ErrorStackMarshaler = MarshalStack
func MarshalStack(err error) any {
	stacks := marshalStack(err)
	slices.Reverse(stacks)

	return stacks
}
