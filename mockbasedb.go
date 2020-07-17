package testutils

import (
	"reflect"

	"github.com/go-pg/pg/v9"
)

type mockBaseDB struct {
	db *MockDB
}

func (db *mockBaseDB) RunInTransaction(fn func(tx *MockTx) error) error {
	tx := &MockTx{db: db.db}
	if err := fn(tx); err != nil {
		return err
	}

	return nil
}

func (db *mockBaseDB) Query(model interface{}, query interface{}, params ...interface{}) (pg.Result, error) {
	model = db.db.responses[0]
	db.db.responses = db.db.responses[1:]

	return (pg.Result)(nil), nil
}

func (db *mockBaseDB) QueryOne(model interface{}, query interface{}, params ...interface{}) (pg.Result, error) {
	model = db.db.responses[0]
	db.db.responses = db.db.responses[1:]

	return (pg.Result)(nil), nil
}

func (db *mockBaseDB) Select(model interface{}) error {
	model = db.db.responses[0]
	return nil
}

func (db *mockBaseDB) Insert(model ...Model) error {
	db.db.models = append(db.db.models, model...)
	return nil
}

func (db *mockBaseDB) Update(model Model) error {
	for i, r := range db.db.models {
		if reflect.TypeOf(r) == reflect.TypeOf(model) && r.GetID() == model.GetID() {
			db.db.models[i] = model
		}
	}
	return nil
}

func (db *mockBaseDB) Delete(model Model) error {
	for i, r := range db.db.models {
		if reflect.TypeOf(r) == reflect.TypeOf(model) && r.GetID() == model.GetID() {
			if len(db.db.models) > 1 {
				db.db.models = append(db.db.models[:i], db.db.models[i+1:]...)
			}
		}
	}
	return nil
}
