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
			in:   sqlw.ExportNewDB(master),
			err:  nil,
		},
		{
			name: "one master and one replica",
			in:   sqlw.ExportNewDB(master, replica1, replica2),
			err:  nil,
		},
		{
			name: "one master and two replicas",
			in:   sqlw.ExportNewDB(closedmaster, replica1, replica2),
			err:  nil,
		},
		{
			name: "one master and two replicas",
			in:   sqlw.ExportNewDB(closedmaster, closedreplica, replica2),
			err:  nil,
		},
		{
			name: "one master and one replica",
			in:   sqlw.ExportNewDB(closedmaster, closedreplica),
			err:  errors.New("connection refused from all databases"),
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			err := tt.in.Readable()
			if err != nil && err.Error() != tt.err.Error() {
				t.Errorf("testing %s: should be error of %v but got: %v", tt.name, tt.err, err)
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
			in:   sqlw.ExportNewDB(master),
			err:  nil,
		},
		{
			name: "one master",
			in:   sqlw.ExportNewDB(closedmaster),
			err:  errors.New("sql: database is closed"),
		},
		{
			name: "one master and two replicas",
			in:   sqlw.ExportNewDB(master, replica1, replica2),
			err:  nil,
		},
		{
			name: "one master and one replica",
			in:   sqlw.ExportNewDB(master, closedreplica),
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
