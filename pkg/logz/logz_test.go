package logz

import (
	"context"
	"fmt"
	"github.com/piotrekmonko/portfello/pkg/conf"
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
)

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
	logger := NewLogger(&conf.Logging{
		Level:  "debug",
		Format: "dev",
	})
	assert.NotNil(t, logger)
}
