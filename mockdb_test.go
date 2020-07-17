package testutils

import (
	"reflect"
	"testing"
)

type TestModel struct {
	ID   int    `pg:"id" json:"id"`
	Name string `pg:"name" json:"name"`
}

func (tm *TestModel) GetID() string {
	return string((*tm).ID)
}

func (tm *TestModel) Equals(i interface{}) bool {
	b, ok := i.(*TestModel)
	if !ok {
		return false
	}
	return tm.ID == b.ID && tm.Name == b.Name
}

func TestMockDB(t *testing.T) {
	t.Run("Implements DB interface", func(t *testing.T) {
		var db interface{}
		db = NewMockDB()
		if _, ok := db.(DB); !ok {
			t.Fatal("MockDB does not of interface type DB")
		}
	})

	t.Run("QueueResponses", func(t *testing.T) {
		db := &MockDB{}
		tm := &TestModel{ID: 1, Name: "Test Model"}
		db.QueueResponses(tm)
		n := len(db.responses)
		if n != 1 {
			t.Fatal("expected 1 response in queue; found ", n)
		}
		eT := reflect.TypeOf(tm)
		rT := reflect.TypeOf(db.responses[0])
		if rT != eT {
			t.Fatal("expected response in queue of type ", eT.String(), "; found ", rT.String())
		}
		if db.responses[0] != tm {
			t.Fatalf("response in queue (%v) does not match expected (%v)", db.responses[0], tm)
		}
	})

	t.Run("QueueModels", func(t *testing.T) {
		db := &MockDB{}
		tm := &TestModel{ID: 1, Name: "Test Model"}
		db.QueueModels(tm)
		n := len(db.models)
		if n != 1 {
			t.Fatal("expected 1 model in queue; found ", n)
		}
		eT := reflect.TypeOf(tm)
		rT := reflect.TypeOf(db.models[0])
		if rT != eT {
			t.Fatal("expected model in queue of type ", eT.String(), "; found ", rT.String())
		}
		if db.models[0].GetID() != tm.GetID() {
			t.Fatalf("model in queue (%v) does not match expected (%v)", db.models[0], tm)
		}
	})

	t.Run("Query", func(t *testing.T) {
		tm := &TestModel{ID: 1, Name: "Test Model"}
		emptyModel := &TestModel{}
		db := &MockDB{responses: []interface{}{*tm}}

		response := &TestModel{}
		_, err := db.Query(response, "SELECT whatever FROM fake_table")
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
		emptyModel := &TestModel{}
		db := &MockDB{responses: []interface{}{*tm}}

		response := &TestModel{}
		_, err := db.Query(response, "SELECT whatever FROM fake_table")
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
		emptyModel := &TestModel{}
		db := &MockDB{responses: []interface{}{*tm}}

		response := &TestModel{}
		err := db.Select(response)
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

		if err := db.Insert(tm); err != nil {
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
			t.Fatalf("MockDB model (%v) does not match expected (%v)", iM, tm)
		}
	})

	t.Run("Update", func(t *testing.T) {
		tm := &TestModel{ID: 1, Name: "Test Model"}
		db := &MockDB{models: []Model{tm}}

		tm2 := &TestModel{ID: 1, Name: "Updated Model"}
		if err := db.Update(tm2); err != nil {
			t.Fatal(err)
		}
		um := db.models[0]
		if !um.Equals(tm2) {
			t.Fatalf("MockDB model (%v) does not match expected (%v)", *um.(*TestModel), *tm)
		}
	})

	t.Run("Delete", func(t *testing.T) {
		tm := &TestModel{ID: 1, Name: "Test Model"}
		db := &MockDB{models: []Model{tm}}

		if err := db.Delete(tm); err != nil {
			t.Fatal(err)
		}
		if len(db.models) != 0 {
			t.Fatal("MockDB.models should be empty")
		}
	})

	t.Run("Find", func(t *testing.T) {
		tm := &TestModel{ID: 1, Name: "Test Model"}
		db := &MockDB{models: []Model{tm}}

		fm, err := db.Find(tm)
		if err != nil {
			t.Fatal(err)
		}
		if !fm.Equals(tm) {
			t.Fatal("Did not find test model")
		}
	})
}
