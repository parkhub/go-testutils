package testutils

import (
	"encoding/json"
	"fmt"
	"reflect"

	"github.com/go-pg/pg/v9"
)

// MockDB implements the DB interface to mock a pg.DB instance
type MockDB struct {
	responses []interface{}
	models    []Model
}

// NewMockDB creates a new mock database client for unit tests
func NewMockDB() *MockDB {
	db := MockDB{}
	return &db
}

// Begin starts a transaction. Most callers should use RunInTransaction instead.
func (db *MockDB) Begin() (*MockTx, error) {
	tx := &MockTx{db: db, open: true, models: append(db.models)}
	return tx, nil
}

// QueueResponses allows a test to add an ordered list of mock responses to the
// database for Query, QueryOne, and Select calls
func (db *MockDB) QueueResponses(response ...interface{}) {
	db.responses = append(db.responses, response...)
}

// QueueResponses allows a test to add a list of mock data models to the
// database for Update and Delete calls
func (db *MockDB) QueueModels(model ...Model) {
	db.models = append(db.models, model...)
}

// RunInTransaction runs a function in a transaction. If function
// returns an error transaction is rollbacked, otherwise transaction
// is committed.
func (db *MockDB) RunInTransaction(fn func(tx Tx) error) error {
	tx := &MockTx{db: db, open: true, models: append(db.models)}
	if err := fn(tx); err != nil {
		return err
	}
	if err := tx.Commit(); err != nil {
		_ = tx.Rollback()
		return err
	}
	if err := tx.Close(); err != nil {
		return err
	}

	return nil
}

// Query executes a query that returns rows, typically a SELECT.
// The params are for any placeholders in the query.
func (db *MockDB) Query(model, query interface{}, params ...interface{}) (pg.Result, error) {
	if len(db.responses) == 0 {
		return nil, nil
	}
	reflect.ValueOf(model).Elem().Set(reflect.ValueOf(db.responses[0]))
	db.responses = db.responses[1:]

	return (pg.Result)(nil), nil
}

// QueryOne acts like Query, but query must return only one row. It
// returns ErrNoRows error when query returns zero rows or
// ErrMultiRows when query returns multiple rows.
func (db *MockDB) QueryOne(model, query interface{}, params ...interface{}) (pg.Result, error) {
	if len(db.responses) == 0 {
		return nil, pg.ErrNoRows
	}
	reflect.ValueOf(model).Elem().Set(reflect.ValueOf(db.responses[0]))
	db.responses = db.responses[1:]

	return (pg.Result)(nil), nil
}

// Select finds a model in the models slice
func (db *MockDB) Select(model interface{}) error {
	if len(db.responses) == 0 {
		return nil
	}
	reflect.ValueOf(model).Elem().Set(reflect.ValueOf(db.responses[0]))
	return nil
}

// Insert appends a model to the models slice
func (db *MockDB) Insert(model ...interface{}) error {
	tms := make([]Model, len(model))
	for i, m := range model {
		tms[i] = m.(Model)
	}
	db.models = append(db.models, tms...)
	return nil
}

// Update finds a model in the models slice based on its GetID() and updates it,
// or returns an error if it is not found
func (db *MockDB) Update(model interface{}) error {
	for i, r := range db.models {
		m, ok := r.(Model)
		if !ok {
			return fmt.Errorf("model is not a Model: %v", r)
		}
		if reflect.TypeOf(m) == reflect.TypeOf(model) && m.GetID() == model.(Model).GetID() {
			reflect.ValueOf(db.models[i]).Elem().Set(reflect.ValueOf(model).Elem())
			return nil
		}
	}
	return fmt.Errorf("%s model with ID %s not found to update",
		reflect.TypeOf(model).String(),
		model.(Model).GetID())
}

// Delete finds a model in the DB and removes it, or returns an error if it is
// not found
func (db *MockDB) Delete(model interface{}) error {
	for i, r := range db.models {
		m, ok := r.(Model)
		if !ok {
			return fmt.Errorf("model is not a Model: %v", r)
		}
		if reflect.TypeOf(m) == reflect.TypeOf(model) && m.GetID() == model.(Model).GetID() {
			db.models = append(db.models[:i], db.models[i+1:]...)
			return nil
		}
	}
	return fmt.Errorf("%s model with ID %s not found to delete",
		reflect.TypeOf(model).String(),
		model.(Model).GetID())
}

// Find searches through the MockDB models and returns a model of matching type
// and ID if it exists, or nil if not.
func (db *MockDB) Find(model Model) (Model, error) {
	for _, r := range db.models {
		if reflect.TypeOf(r) == reflect.TypeOf(model) && r.GetID() == model.(Model).GetID() {
			return r, nil
		}
	}
	return nil, nil
}

// MarshalModels returns a pretty string of JSON for logging out the contents of
// the MockDB models
func (db *MockDB) MarshalModels() (string, error) {
	bytes, err := json.MarshalIndent(db.models, "", "  ")
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

// MarshalModels returns a pretty string of JSON for logging out the contents of
// the MockDB responses
func (db *MockDB) MarshalResponses() (string, error) {
	bytes, err := json.MarshalIndent(db.responses, "", "  ")
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}