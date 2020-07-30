package testutils

import (
	"errors"
	"reflect"

	"github.com/go-pg/pg/v9"
)

// MockTx implements mocks a go-pg query
type MockQuery struct {
	db *MockDB
	tx *MockTx
	queryModels []Model
}

func (q *MockQuery) Clone() Query {
	return &MockQuery{
		db: q.db,
		tx: q.tx,
		queryModels: append(q.queryModels),
	}
}

func (q *MockQuery) DB(db DB) Query {
	q.db = db.(*MockDB)
	return q
}

// Deleted adds `WHERE deleted_at IS NOT NULL` clause for soft deleted models.
func (q *MockQuery) Deleted() Query {
	return q
}

// AllWithDeleted changes query to return all rows including soft deleted ones.
func (q *MockQuery) AllWithDeleted() Query {
	return q
}

// With adds subq as common table expression with the given name.
func (q *MockQuery) With(name string, subq Query) Query {
	return q
}

func (q *MockQuery) WithInsert(name string, subq Query) Query {
	return q
}

func (q *MockQuery) WithUpdate(name string, subq Query) Query {
	return q
}

func (q *MockQuery) WithDelete(name string, subq Query) Query {
	return q
}

// WrapWith creates new Query and adds to it current query as
// common table expression with the given name.
func (q *MockQuery) WrapWith(name string) Query {
	return q
}

func (q *MockQuery) Table(tables ...string) Query {
	return q
}

func (q *MockQuery) TableExpr(expr string, params ...interface{}) Query {
	return q
}

func (q *MockQuery) Distinct() Query {
	return q
}

func (q *MockQuery) DistinctOn() Query {
	return q
}

func (q *MockQuery) Column(columns ...string) Query {
	return q
}

func (q *MockQuery) ColumnExpr(expr string, params ...interface{}) Query {
	return q
}

func (q *MockQuery) ExcludeColumn(columns ...string) Query {
	return q
}

func (q *MockQuery) Relation(name string, apply ...func(Query) (Query, error)) Query {
	return q
}

func (q *MockQuery) Set(set string, params ...interface{}) Query {
	panic("implement me")
	// for _, qm := range q.queryModels {
	// 	qmValue := reflect.Indirect(reflect.ValueOf(qm))
	// 	structType := qmValue.Type()
	// 	// Loop through each field
	// 	for i := 0; i < qmValue.NumField(); i++ {
	// 		// If field is exported...
	// 		structField := structType.Field(i)
	// 		if structField.PkgPath == "" {
	// 			qmField := qmValue.Field(i)
	// 			if structField.Tag.Get("pg") == set {
	// 				qmField.Set(reflect.ValueOf(params))
	// 				return q
	// 			}
	// 		}
	// 	}
	// }
	// return q
}

// Value overwrites model value for the column in INSERT and UPDATE queries.
func (q *MockQuery) Value(column string, value string, params ...interface{}) Query {
	for _, qm := range q.queryModels {
		qmValue := reflect.Indirect(reflect.ValueOf(qm))
		structType := qmValue.Type()
		// Loop through each field
		for i := 0; i < qmValue.NumField(); i++ {
			// If field is exported...
			structField := structType.Field(i)
			if structField.PkgPath == "" {
				qmField := qmValue.Field(i)
				if structField.Tag.Get("pg") == column {
					qmField.Set(reflect.ValueOf(value))
					return q
				}
			}
		}
	}
	return q
}

func (q *MockQuery) Where(condition string, params ...interface{}) Query {
	return q
}

func (q *MockQuery) WhereOr(condition string, params ...interface{}) Query {
	return q
}

func (q *MockQuery) WhereGroup(fn func(Query) (Query, error)) Query {
	return q
}

func (q *MockQuery) WhereNotGroup(fn func(Query) (Query, error)) Query {
	return q
}

func (q *MockQuery) WhereOrGroup(fn func(Query) (Query, error)) Query {
	return q
}

func (q *MockQuery) WhereOrNotGroup(fn func(query Query) (Query, error)) Query {
	return q
}

// WhereIn is a shortcut for Where and pg.In:
func (q *MockQuery) WhereIn(where string, slice interface{}) Query {
	return q
}

// WhereInMulti is a shortcut for Where and pg.InMulti:
func (q *MockQuery) WhereInMulti(where string, values ...interface{}) Query {
	return q
}

func (q *MockQuery) WherePK() Query {
	return q
}

func (q *MockQuery) WhereStruct(strct interface{}) Query {
	return q
}

func (q *MockQuery) Join(join string, params ...interface{}) Query {
	return q
}

func (q *MockQuery) JoinOn(condition string, params ...interface{}) Query {
	return q
}

func (q *MockQuery) JoinOnOr(condition string, params ...interface{}) Query {
	return q
}

func (q *MockQuery) Group(columns ...string) Query {
	return q
}

func (q *MockQuery) GroupExpr(group string, params ...interface{}) Query {
	return q
}

func (q *MockQuery) Having(having string, params ...interface{}) Query {
	return q
}

func (q *MockQuery) Union(other Query) Query {
	return q
}

func (q *MockQuery) UnionAll(other Query) Query {
	return q
}

func (q *MockQuery) Intersect(other Query) Query {
	return q
}

func (q *MockQuery) IntersectAll(other Query) Query {
	return q
}

func (q *MockQuery) Except(other Query) Query {
	return q
}

func (q *MockQuery) ExceptAll(other Query) Query {
	return q
}

func (q *MockQuery) Order(orders ...string) Query {
	return q
}

// Order adds sort order to the Query.
func (q *MockQuery) OrderExpr(order string, params ...interface{}) Query {
	return q
}

func (q *MockQuery) Limit(n int) Query {
	return q
}

func (q *MockQuery) Offset(n int) Query {
	return q
}

func (q *MockQuery) OnConflict(s string, params ...interface{}) Query {
	return q
}

func (q *MockQuery) Returning(s string, params ...interface{}) Query {
	return q
}

func (q *MockQuery) For(s string, params ...interface{}) Query {
	return q
}

// Apply calls the fn passing the Query as an argument.
func (q *MockQuery) Apply(fn func(Query) (Query, error)) Query {
	return q
}

// Count returns number of rows matching the query using count aggregate function.
func (q *MockQuery) Count() (int, error) {
	if len(q.db.responses) == 0 {
		return 0, errors.New("no results in response queue")
	}
	r0 := q.db.responses[0]
	q.db.responses = q.db.responses[1:]

	// If the "count" value itself is inserted as a response
	if i, ok := q.db.responses[0].(int); ok {
		return i, nil
	}

	// If the size of the response should be counted
	rVal := reflect.ValueOf(r0)
	if rVal.Kind() == reflect.Slice {
		return rVal.Len(), nil
	}
	return 0, errors.New("next response in queue is not a slice")
}

func (q *MockQuery) First() error {
	// What is this supposed to do?
	panic("implement me")
}

func (q *MockQuery) Last() error {
	// What is this supposed to do?
	panic("implement me")
}

// Select selects the model
func (q *MockQuery) Select(values ...interface{}) error {
	return q.db.Select(q.queryModels[0])
}

// SelectAndCount runs Select and Count in two goroutines,
// waits for them to finish and returns the result. If query limit is -1
// it does not select any data and only counts the results.
func (q *MockQuery) SelectAndCount(values ...interface{}) (count int, firstErr error) {
	if len(q.db.responses) == 0 {
		return 0, errors.New("no results in response queue")
	}
	r0 := q.db.responses[0]

	// If the size of the response should be counted
	rVal := reflect.ValueOf(r0)
	if rVal.Kind() == reflect.Slice {
		err := q.db.Select(q.queryModels[0])
		return rVal.Len(), err
	}
	return 0, errors.New("next response in queue is not a slice")
}

// SelectAndCountEstimate runs Select and CountEstimate in two goroutines,
// waits for them to finish and returns the result. If query limit is -1
// it does not select any data and only counts the results.
func (q *MockQuery) SelectAndCountEstimate(threshold int, values ...interface{}) (count int, firstErr error) {
	return q.SelectAndCount(values...)
}

// ForEach calls the function for each row returned by the query
// without loading all rows into the memory.
//
// Function can accept a struct, a pointer to a struct, an orm.Model,
// or values for the columns in a row. Function must return an error.
// func (q *MockQuery) ForEach(fn interface{}) error {
// 	panic("implement me")
// }

// Insert inserts the model
func (q *MockQuery) Insert(values ...interface{}) (pg.Result, error) {
	inserts := make([]interface{}, len(q.queryModels), len(q.queryModels))
	for i, m := range q.queryModels {
		inserts[i] = m
	}
	return (pg.Result)(nil), q.db.Insert(inserts...)
}

// SelectOrInsert selects the model inserting one if it does not exist.
// It returns true when model was inserted.
func (q *MockQuery) SelectOrInsert(values ...interface{}) (inserted bool, _ error) {
	inserted = false
	for _, m := range q.queryModels {
		f, err := q.db.Find(m)
		if err != nil {
			return false, err
		}
		if f == nil {
			err = q.db.Insert(m)
			if err != nil {
				return false, err
			}
			inserted = true
		} else {
			_ = q.db.Select(m)
		}
	}
	return inserted, nil
}

// Update updates the model.
func (q *MockQuery) Update(scan ...interface{}) (pg.Result, error) {
	for _, m := range q.queryModels {
		err := q.db.Update(m)
		if err != nil {
			return (pg.Result)(nil), err
		}
	}
	return (pg.Result)(nil), nil
}

// Update updates the model omitting fields with zero values such as:
//   - empty string,
//   - 0,
//   - zero time,
//   - empty map or slice,
//   - byte array with all zeroes,
//   - nil ptr,
//   - types with method `IsZero() == true`.
func (q *MockQuery) UpdateNotZero(scan ...interface{}) (pg.Result, error) {
	pgResult := (pg.Result)(nil)
	for _, qm := range q.queryModels {
		for _, m := range q.db.models {
			// If the model in the db is a match
			if reflect.TypeOf(m) == reflect.TypeOf(qm) && m.GetID() == qm.GetID() {
				qmValue := reflect.Indirect(reflect.ValueOf(qm))
				mValue := reflect.Indirect(reflect.ValueOf(m))
				structType := qmValue.Type()
				// Loop through each field
				for i := 0; i < qmValue.NumField(); i++ {
					// If field is exported...
					if structType.Field(i).PkgPath == "" {
						qmField := qmValue.Field(i)
						// set if query model's field is non-zero
						if !qmField.IsZero() {
							mValue.Field(i).Set(qmField)
						}
					}
				}
				// Assume only one model will match
				break
			}
		}
	}
	return pgResult, nil
}

// Delete deletes the model. When model has deleted_at column the row
// is soft deleted instead.
func (q *MockQuery) Delete(values ...interface{}) (pg.Result, error) {
	pgResult := (pg.Result)(nil)
	for _, m := range q.queryModels {
		if err := q.db.Delete(m); err != nil {
			return pgResult, err
		}
	}
	return pgResult, nil
}

// Delete forces delete of the model with deleted_at column.
func (q *MockQuery) ForceDelete(values ...interface{}) (pg.Result, error) {
	pgResult := (pg.Result)(nil)
	for _, m := range q.queryModels {
		if err := q.db.Delete(m); err != nil {
			return pgResult, err
		}
	}
	return pgResult, nil
}

// Exec is an alias for DB.Exec.
// func (q *MockQuery) Exec(query interface{}, params ...interface{}) (pg.Result, error) {
	// return q.db.Exec(query, params...)
	// panic("implement me")
// }

// ExecOne is an alias for DB.ExecOne.
// func (q *MockQuery) ExecOne(query interface{}, params ...interface{}) (pg.Result, error) {
// 	panic("implement me")
// }

// Query is an alias for DB.Query.
func (q *MockQuery) Query(model, query interface{}, params ...interface{}) (pg.Result, error) {
	return q.db.Query(q.queryModels[0], query, params)
}

// QueryOne is an alias for DB.QueryOne.
func (q *MockQuery) QueryOne(model, query interface{}, params ...interface{}) (pg.Result, error) {
	return q.db.QueryOne(q.queryModels[0], query, params)
}

// Exists returns true or false depending if there are any rows matching the query.
func (q *MockQuery) Exists() (bool, error) {
	model, err := q.db.Find(q.queryModels[0])
	if err != nil {
		return false, err
	}
	if model == nil {
		return false, nil
	}
	return true, nil
}
