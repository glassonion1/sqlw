package sqlw_test

import (
	"testing"

	"github.com/glassonion1/sqlw"
)

func TestSQLQuery(t *testing.T) {
	tests := []struct {
		name string
		in   sqlw.SQLQuery
		err  error
	}{
		{
			name: "select statemant(lowercase)",
			in:   sqlw.SQLQuery("select * from hoge"),
			err:  nil,
		},
		{
			name: "select statement(uppercase)",
			in:   sqlw.SQLQuery("SELECT * FROM hoge"),
			err:  nil,
		},
		{
			name: "select statement",
			in:   sqlw.SQLQuery("SeLECt * FROM hoge"),
			err:  nil,
		},
		{
			name: "not select",
			in:   sqlw.SQLQuery("SELTC * FROM hoge"),
			err:  sqlw.ErrNotSQLQuery,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			err := tt.in.Validate()
			if err != nil && err.Error() != tt.err.Error() {
				t.Errorf("testing %s: should be error of %v but got: %v", tt.name, tt.err, err)
			}
		})
	}
}

func TestSQLMutation(t *testing.T) {
	tests := []struct {
		name string
		in   sqlw.SQLMutation
		err  error
	}{
		{
			name: "insert statement(lowercase)",
			in:   sqlw.SQLMutation("insert into hoge values()"),
			err:  nil,
		},
		{
			name: "insert statement(uppercase)",
			in:   sqlw.SQLMutation("INSERT into hoge values()"),
			err:  nil,
		},
		{
			name: "insert statement",
			in:   sqlw.SQLMutation("inSErt into hoge values()"),
			err:  nil,
		},
		{
			name: "update statement(lowercase)",
			in:   sqlw.SQLMutation("update hoge set"),
			err:  nil,
		},
		{
			name: "update statement(uppercase)",
			in:   sqlw.SQLMutation("UPDATE hoge SET"),
			err:  nil,
		},
		{
			name: "update statement",
			in:   sqlw.SQLMutation("UpdaTe hoge set"),
			err:  nil,
		},
		{
			name: "delete statement(lowercase)",
			in:   sqlw.SQLMutation("delete from hoge"),
			err:  nil,
		},
		{
			name: "delete statement(uppercase)",
			in:   sqlw.SQLMutation("DELETE FROM hoge"),
			err:  nil,
		},
		{
			name: "delete statement",
			in:   sqlw.SQLMutation("Delete from hoge"),
			err:  nil,
		},
		{
			name: "not insert update delete",
			in:   sqlw.SQLMutation("select into hoge values()"),
			err:  sqlw.ErrNotSQLMutation,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			err := tt.in.Validate()
			if err != nil && err.Error() != tt.err.Error() {
				t.Errorf("testing %s: should be error of %v but got: %v", tt.name, tt.err, err)
			}
		})
	}
}
