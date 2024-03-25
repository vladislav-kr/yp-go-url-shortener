package logger

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

)

func TestMustLogger(t *testing.T) {

	cases := []struct {
		name  string
		level string
	}{
		{
			name:  fmt.Sprintf("level %s", DebugLevel),
			level: DebugLevel,
		},
		{
			name:  fmt.Sprintf("level %s", InfoLevel),
			level: InfoLevel,
		},
		{
			name:  fmt.Sprintf("level %s", WarnLevel),
			level: WarnLevel,
		},
		{
			name:  fmt.Sprintf("level %s", ErrorLevel),
			level: ErrorLevel,
		},
		{
			name:  fmt.Sprintf("level %s", DPanicLevel),
			level: DPanicLevel,
		},
		{
			name:  fmt.Sprintf("level %s", PanicLevel),
			level: PanicLevel,
		},
		{
			name:  fmt.Sprintf("level %s", FatalLevel),
			level: FatalLevel,
		},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			log := MustLogger(tc.level)

			assert.NotEmpty(t, log)
		})
	}

}
