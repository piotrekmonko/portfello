package logz

import (
	"context"
	"fmt"
	"github.com/piotrekmonko/portfello/pkg/config"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap/zaptest"
	"reflect"
	"strings"
	"testing"
)

type TestLogger struct {
	testing.TB
	*Log
	Messages []string
}

func NewTestLogger(tb testing.TB) *TestLogger {
	tl := &TestLogger{
		TB:       tb,
		Messages: make([]string, 0),
	}
	tl.Log = &Log{SugaredLogger: zaptest.NewLogger(tl).Sugar()}
	return tl
}

func (t *TestLogger) Logf(format string, args ...interface{}) {
	m := fmt.Sprintf(format, args...)
	m = m[strings.IndexByte(m, '\t')+1:] // strip the timestamp and its following tab
	t.Messages = append(t.Messages, m)
	t.TB.Log(m)
}

func (t *TestLogger) AssertMessages(msgs ...string) {
	assert.Equal(t.TB, msgs, t.Messages, "logged messages did not match")
}

func TestFromCtx(t *testing.T) {
	tests := []struct {
		name  string
		input context.Context
		want  []any
	}{
		{
			name:  "nil ctx",
			input: nil,
			want:  []any{},
		},
		{
			name:  "empty ctx",
			input: context.Background(),
			want:  []any{},
		},
		{
			name:  "full ctx",
			input: WithCtx(context.Background(), "a", "b"),
			want:  []any{"a", "b"},
		},
		{
			name:  "doubly appended ctx",
			input: WithCtx(WithCtx(context.Background(), "a", "b"), "c", "d"),
			want:  []any{"a", "b", "c", "d"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := FromCtx(tt.input); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("FromCtx() = %+v, want %+v", got, tt.want)
			}
		})
	}
}

func TestLog_Levels(t *testing.T) {
	ctx := WithCtx(context.Background(), "a", "b")
	log := NewTestLogger(t)
	log.With("with", "fields").Debugw(ctx, "msg debug", "c", 1)
	log.Named("test-name").Infow(ctx, "msg info", "d", 2)
	log.Warnw(ctx, "msg warn", "e", 3)
	_ = log.Errorw(ctx, fmt.Errorf("an error"), "msg error", "f", 4)

	log.AssertMessages(
		`DEBUG	msg debug	{"with": "fields", "a": "b", "c": 1}`,
		`INFO	test-name	msg info	{"a": "b", "d": 2}`,
		`WARN	msg warn	{"a": "b", "e": 3}`,
		`ERROR	msg error	{"error": "an error", "a": "b", "f": 4}`,
	)
}

func TestNewLogger(t *testing.T) {
	logger := NewLogger(&config.Logging{
		Level:  "debug",
		Format: "dev",
	})
	assert.NotNil(t, logger)
}
