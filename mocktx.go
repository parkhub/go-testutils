package testutils

import (
	"encoding/json"
	"fmt"
	"reflect"

	"github.com/go-pg/pg/v9"
)

// MockTx implements the Tx interface to mock a pg.Tx instance
type MockTx struct {
	db   *MockDB
	open bool
	models  []Model
}

func (tx *MockTx) RunInTransaction(fn func(tx *MockTx) error) error {
	if err := fn(tx); err != nil {
		return err
	}
	if err := tx.Commit(); err != nil {
		return err
	}
	if err := tx.Close(); err != nil {
		return err
	}

	return nil
}

// Query is an alias for DB.Query
func (tx *MockTx) Query(model interface{}, query interface{}, params ...interface{}) (pg.Result, error) {
	return tx.db.Query(model, query, params...)
}

// QueryOne is an alias for DB.QueryOne
func (tx *MockTx) QueryOne(model interface{}, query interface{}, params ...interface{}) (pg.Result, error) {
	return tx.db.QueryOne(model, query, params...)
}

// Select is an alias for DB.Select
func (tx *MockTx) Select(model interface{}) error {
	return tx.db.Select(model)
}

// Insert is an alias for DB.Insert
func (tx *MockTx) Insert(model ...interface{}) error {
	// return tx.db.Insert(model...)
	tms := make([]Model, len(model))
	for i, m := range model {
	tms[i] = m.(Model)
	}
	tx.models = append(tx.models, tms...)
	return nil
}

// Update is an alias for DB.Update
func (tx *MockTx) Update(model interface{}) error {
	// return tx.db.Update(model)
	for i, r := range tx.models {
		m, ok := r.(Model)
		if !ok {
			return fmt.Errorf("model is not a Model: %v", r)
		}
		if reflect.TypeOf(m) == reflect.TypeOf(model) && m.GetID() == model.(Model).GetID() {
			reflect.ValueOf(tx.models[i]).Elem().Set(reflect.ValueOf(model).Elem())
			return nil
		}
	}
	return fmt.Errorf("%s model with ID %s not found to update",
		reflect.TypeOf(model).String(),
		model.(Model).GetID())
}

// Delete is an alias for DB.Delete
func (tx *MockTx) Delete(model interface{}) error {
	// return tx.db.Delete(model)
	for i, r := range tx.models {
		m, ok := r.(Model)
		if !ok {
			return fmt.Errorf("model is not a Model: %v", r)
		}
		if reflect.TypeOf(m) == reflect.TypeOf(model) && m.GetID() == model.(Model).GetID() {
			tx.models = append(tx.models[:i], tx.models[i+1:]...)
			return nil
		}
	}
	return fmt.Errorf("%s model with ID %s not found to delete",
		reflect.TypeOf(model).String(),
		model.(Model).GetID())
}

// Commit commits the transaction.
func (tx *MockTx) Commit() error {
	tx.db.models = tx.models
	tx.db.models = append(tx.models)
	tx.models = nil
	return nil
}

// Rollback aborts the transaction.
func (tx *MockTx) Rollback() error {
	tx.models = nil
	return nil
}

// Close calls Rollback if the tx has not already been committed or rolled back.
func (tx *MockTx) Close() error {
	defer func() { tx.open = false }()
	if tx.open {
		return tx.Rollback()
	}
	return nil
}

func (tx *MockTx) Model(model ...interface{}) Query {
	numModels := len(model)
	queryModels := make([]Model, numModels, numModels)
	for i, m := range model {
		queryModels[i] = m.(Model)
	}
	return &MockQuery{
		db: tx.db,
		queryModels: queryModels,
	}
}

// MarshalModels returns a pretty string of JSON for logging out the contents of
// the MockTx models
func (tx *MockTx) MarshalModels() (string, error) {
	bytes, err := json.MarshalIndent(tx.models, "", "  ")
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}
