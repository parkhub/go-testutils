package testutils

import (
	"reflect"
	"testing"
)

func TestMockTx(t *testing.T) {
	emptyModel := &TestModel{}

	t.Run("Implements Tx interface", func(t *testing.T) {
		db := MockDB{}
		var tx interface{}
		tx, _ = db.Begin()
		if _, ok := tx.(Tx); !ok {
			t.Fatal("MockTx does not of interface type Tx")
		}
	})

	t.Run("Query", func(t *testing.T) {
		tm := &TestModel{ID: 1, Name: "Test Model"}
		db := &MockDB{responses: []interface{}{*tm}}
		response := &TestModel{}

		err := db.RunInTransaction(func(tx Tx) error {
			_, err := tx.Query(response, "SELECT whatever FROM fake_table")
			return err
		})
		if err != nil {
			t.Fatal(err)
		}
		if response.Equals(emptyModel) {
			t.Fatal("Returned TestModel is empty")
		}
		if !response.Equals(tm) {
			t.Fatal("response struct doesn't match queued response")
		}
	})

	t.Run("QueryOne", func(t *testing.T) {
		tm := &TestModel{ID: 1, Name: "Test Model"}
		db := &MockDB{responses: []interface{}{*tm}}
		response := &TestModel{}

		err := db.RunInTransaction(func(tx Tx) error {
			_, err := tx.Query(response, "SELECT whatever FROM fake_table")
			return err
		})
		if err != nil {
			t.Fatal(err)
		}
		if response.Equals(emptyModel) {
			t.Fatal("Returned TestModel is empty")
		}
		if !response.Equals(tm) {
			t.Fatal("response struct doesn't match queued response")
		}
	})

	t.Run("Select", func(t *testing.T) {
		tm := &TestModel{ID: 1, Name: "Test Model"}
		db := &MockDB{responses: []interface{}{*tm}}
		response := &TestModel{}

		err := db.RunInTransaction(func(tx Tx) error {
			err := db.Select(response)
			return err
		})

		if err != nil {
			t.Fatal(err)
		}
		if response.Equals(emptyModel) {
			t.Fatal("Returned TestModel is empty")
		}
		if !response.Equals(tm) {
			t.Fatal("response struct doesn't match queued response")
		}
	})

	t.Run("Insert", func(t *testing.T) {
		tm := &TestModel{ID: 1, Name: "Test Model"}
		db := &MockDB{}

		err := db.RunInTransaction(func(tx Tx) error {
			err := tx.Insert(tm)
			return err
		})
		if err != nil {
			t.Fatal(err)
		}

		n := len(db.models)
		if n != 1 {
			t.Fatal("expected 1 model in queue; found ", n)
		}
		eT := reflect.TypeOf(tm)
		rT := reflect.TypeOf(db.models[0])
		if rT != eT {
			t.Fatal("expected model in queue of type ", eT.String(), "; found ", rT.String())
		}
		iM := db.models[0]
		if !iM.Equals(tm) {
			t.Fatalf("MockDB model (%v) does not match expected (%v)", iM, *tm)
		}
	})

	t.Run("Update", func(t *testing.T) {
		tm := &TestModel{ID: 1, Name: "Test Model"}
		db := MockDB{models: []Model{tm}}

		um := &TestModel{ID: 1, Name:"Updated Model"}
		err := db.RunInTransaction(func (tx Tx) error {
			err := tx.Update(um)
			return err
		})
		if err != nil {
			t.Fatal(err)
		}
		fm := db.models[0]
		if !fm.Equals(um) {
			t.Fatalf("MockDB model (%v) does not match expected (%v)", fm, um)
		}
	})

	t.Run("Delete", func(t *testing.T) {
		tm := &TestModel{ID: 1, Name: "Test Model"}
		db := MockDB{models: []Model{tm}}

		err := db.RunInTransaction(func (tx Tx) error {
			err := tx.Delete(tm)
			return err
		})
		if err != nil {
			t.Fatal(err)
		}
		if len(db.models) != 0 {
			t.Fatal("MockDB.models should be empty")
		}
	})
}
