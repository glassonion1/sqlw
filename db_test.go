package sqlw_test

import (
	"database/sql"
	"errors"
	"log"
	"testing"

	"github.com/glassonion1/sqlw"
)

// nolint
func TestDBReadable(t *testing.T) {

	master, err := sql.Open("mysql", "root:password@tcp(:3306)/app")
	if err != nil {
		t.Error(err)
	}

	// Deliberately close for testing
	closedmaster, err := sql.Open("mysql", "root:password@tcp(:3306)/app")
	if err != nil {
		t.Error(err)
	}
	closedmaster.Close()

	replica1, err := sql.Open("mysql", "root:password@tcp(:3307)/app")
	if err != nil {
		log.Printf("error replica1: %v", err)
	}

	replica2, err := sql.Open("mysql", "root:password@tcp(:3308)/app")
	if err != nil {
		log.Printf("error replica2: %v", err)
	}

	// Deliberately close for testing
	closedreplica, err := sql.Open("mysql", "root:password@tcp(:3307)/app")
	if err != nil {
		log.Println(err)
	} else {
		closedreplica.Close()
	}

	tests := []struct {
		name    string
		in      *sqlw.DB
		wantErr bool
	}{
		{
			name:    "one master",
			in:      sqlw.NewDB(master),
			wantErr: false,
		},
		{
			name:    "one master and one replica",
			in:      sqlw.NewDB(master, replica1),
			wantErr: false,
		},
		{
			name:    "one master and two replicas",
			in:      sqlw.NewDB(master, replica1, replica2),
			wantErr: false,
		},

		{
			name:    "one closed master and two replicas",
			in:      sqlw.NewDB(closedmaster, replica1, replica2),
			wantErr: true,
		},
		{
			name:    "one closed master and one cloed replica and replica",
			in:      sqlw.NewDB(closedmaster, closedreplica, replica2),
			wantErr: true,
		},
		{
			name:    "one closed master and one closed replica",
			in:      sqlw.NewDB(closedmaster, closedreplica),
			wantErr: true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			err := tt.in.Readable()
			if tt.wantErr != (err != nil) {
				t.Errorf("wantErr: %v, err: %v", tt.wantErr, err)
			}
		})
	}
}

func TestDBWritable(t *testing.T) {

	master, err := sql.Open("mysql", "root:password@tcp(:3306)/app")
	if err != nil {
		t.Error(err)
	}

	// Deliberately close for testing
	closedmaster, err := sql.Open("mysql", "root:password@tcp(:3306)/app")
	if err != nil {
		t.Error(err)
	}
	closedmaster.Close()

	replica1, err := sql.Open("mysql", "root:password@tcp(:3307)/app")
	if err != nil {
		log.Println(err)
	}
	replica2, err := sql.Open("mysql", "root:password@tcp(:3308)/app")
	if err != nil {
		log.Println(err)
	}

	// Deliberately close for testing
	closedreplica, err := sql.Open("mysql", "root:password@tcp(:3307)/app")
	if err != nil {
		log.Println(err)
	} else {
		closedreplica.Close()
	}

	tests := []struct {
		name string
		in   *sqlw.DB
		err  error
	}{
		{
			name: "one master",
			in:   sqlw.NewDB(master),
			err:  nil,
		},
		{
			name: "one master",
			in:   sqlw.NewDB(closedmaster),
			err:  errors.New("sql: database is closed"),
		},
		{
			name: "one master and two replicas",
			in:   sqlw.NewDB(master, replica1, replica2),
			err:  nil,
		},
		{
			name: "one master and one replica",
			in:   sqlw.NewDB(master, closedreplica),
			err:  nil,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			err := tt.in.Writable()
			if err != nil && err.Error() != tt.err.Error() {
				t.Errorf("testing %s: should be error of %v but got: %v", tt.name, tt.err, err)
			}
		})
	}
}
