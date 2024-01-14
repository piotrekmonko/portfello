package dao

import (
	"context"
	"database/sql"
	"github.com/piotrekmonko/portfello/pkg/logz"
	"reflect"
	"testing"
)

func TestDAO_BeginTx_Rollback(t *testing.T) {
	type fields struct {
		log     logz.Logger
		db      DBTX
		txDepth int
		Queries *Queries
	}
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name              string
		fields            fields
		args              args
		want              *DAO
		want1             func() error
		wantErr           bool
		wantRollbackerErr error
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			q := &DAO{
				log:     tt.fields.log,
				db:      tt.fields.db,
				txDepth: tt.fields.txDepth,
				Queries: tt.fields.Queries,
			}
			tTx, rollbacker, err := q.BeginTx(tt.args.ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("BeginTx() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(tTx, tt.want) {
				t.Errorf("BeginTx() got = %v, want %v", tTx, tt.want)
			}
			rollbackerErr := rollbacker()
			if !reflect.DeepEqual(rollbackerErr, tt.wantRollbackerErr) {
				t.Errorf("BeginTx() got1 = %v, want %v", rollbackerErr, tt.wantRollbackerErr)
			}
		})
	}
}

func TestDAO_BeginTx_Commit(t *testing.T) {
	type fields struct {
		log     logz.Logger
		db      DBTX
		txDepth int
		Queries *Queries
	}
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			q := &DAO{
				log:     tt.fields.log,
				db:      tt.fields.db,
				txDepth: tt.fields.txDepth,
				Queries: tt.fields.Queries,
			}
			if err := q.Commit(tt.args.ctx); (err != nil) != tt.wantErr {
				t.Errorf("Commit() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestNewDAO(t *testing.T) {
	type args struct {
		ctx context.Context
		log logz.Logger
		dsn string
	}
	tests := []struct {
		name    string
		args    args
		want    *sql.DB
		want1   *DAO
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1, err := NewDAO(tt.args.ctx, tt.args.log, tt.args.dsn)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewDAO() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewDAO() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("NewDAO() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

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
