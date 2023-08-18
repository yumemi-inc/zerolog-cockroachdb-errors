package cockroachdberrors

import (
	"bytes"
	"regexp"
	"testing"

	"github.com/cockroachdb/errors"
	"github.com/rs/zerolog"
)

func TestLogStack(t *testing.T) {
	zerolog.ErrorStackMarshaler = MarshalStack

	out := &bytes.Buffer{}
	log := zerolog.New(out)

	err := errors.Wrap(errors.New("error message"), "from error")
	log.Log().Stack().Err(err).Msg("")

	got := out.String()
	want := `\{"stack":\[\{"stacktrace":\[\{"source":"stacktrace_test.go","line":"18","func":"TestLogStack"\},.*\]\},\{"details":\["from error"\]\},\{"stacktrace":\[\{"source":"stacktrace_test.go","line":"18","func":"TestLogStack"\},.*\]\},\{"details":\["error message"\]\}],"error":"from error: error message"\}`
	if ok, _ := regexp.MatchString(want, got); !ok {
		t.Errorf("invalid log output:\ngot:  %v\nwant: %v", got, want)
	}
}

func TestLogStackFromContext(t *testing.T) {
	zerolog.ErrorStackMarshaler = MarshalStack

	out := &bytes.Buffer{}
	log := zerolog.New(out).With().Stack().Logger() // calling Stack() on log context instead of event

	err := errors.Wrap(errors.New("error message"), "from error")
	log.Log().Err(err).Msg("") // not explicitly calling Stack()

	got := out.String()
	want := `\{"stack":\[\{"stacktrace":\[\{"source":"stacktrace_test.go","line":"34","func":"TestLogStackFromContext"\},.*\]\},\{"details":\["from error"\]\},\{"stacktrace":\[\{"source":"stacktrace_test.go","line":"34","func":"TestLogStackFromContext"\},.*\]\},\{"details":\["error message"\]\}],"error":"from error: error message"\}`
	if ok, _ := regexp.MatchString(want, got); !ok {
		t.Errorf("invalid log output:\ngot:  %v\nwant: %v", got, want)
	}
}

func BenchmarkLogStack(b *testing.B) {
	zerolog.ErrorStackMarshaler = MarshalStack
	out := &bytes.Buffer{}
	log := zerolog.New(out)
	err := errors.Wrap(errors.New("error message"), "from error")
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		log.Log().Stack().Err(err).Msg("")
		out.Reset()
	}
}
