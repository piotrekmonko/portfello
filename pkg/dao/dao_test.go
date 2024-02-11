package dao

import (
	"database/sql"
	"reflect"
	"testing"
)

func TestNilStr(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  sql.NullString
	}{
		{
			name:  "empty string",
			input: "",
			want: sql.NullString{
				String: "",
				Valid:  false,
			},
		},
		{
			name:  "valid string",
			input: "valid",
			want: sql.NullString{
				String: "valid",
				Valid:  true,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NilStr(tt.input); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NilStr() = %v, want %v", got, tt.want)
			}
		})
	}
}
