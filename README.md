Test Utilities Package
======================

The testutils package provides a mock database to replace the go-pg/v9 database
client in tests and includes mock transactions.

Replacing *pg.DB
----------------

In order to use the `testutils.MockDB` type in tests, the functions expecting
the database client as a parameter should use the `testutils.DB` interface
instead of `*pg.DB`. The application should pass a `testutils.DBWrapper`
reference to these functions, which wraps `pg.DB` to implement the
`testutils.DB` interface. Tests for these methods will pass a `testutils.MockDB`
reference. The `testutils.DB` and `testutils.Tx` interfaces implement common
methods used from the `pg.DB` and `pg.Tx` types.

To use a type with the mock database types, the type must implement the
`testutils.Model` interface. The required functions may be defined in the test
file(s) for the type, and need not appear in the application code.

Differences Between `pg.DB` and the `testutils.DB` interface
------------------------------------------------------------

`(*pg.DB) RunInTransaction (func (*pg.Tx) error) error` method calls will need
to be updated to pass a function with signature `(func (testutils.Tx) error)`.
Use of the `DB` interface is otherwise the same.

Interfaces Provided
-------------------

### Model

The `Model` interface must be implemented by types inserted, updated, and
deleted from the mock database. It allows the mock database to find a "matching"
record without a SQL parser (and becoming something too similar to a real
database, which exceeds the intent of this package).

#### `GetID() string`

GetID should return a `string` that represents a unique identifier for this type.
In most cases, it may be the primary key cast to a `string`.

#### `Equals(interface{}) bool`

Equals returns whether or not the receiver and passed variable match for all
important fields. Some fields cannot be known (for example, created and
modified timestamps) and should not be compared for equality in this function.
It mainly serves two uses: to check if the receiver matches an empty struct and
to see if the receiver matches all important fields in a struct of the same type
that represents the expected value. It should return false if the type of the
passed value and the receiver do not match.

### BaseDB

`go-pg`'s DB type includes `pg.BaseDB` by composition, so it is included here,
but no matching mock type is implemented. The `testutils.DB` interface includes
these and the `testutils.MockDB` type implements these methods directly.

#### `RunInTransaction(fn func(Tx) error) error`

RunInTransaction runs a function in a transaction. If function returns an error
transaction is rollbacked, otherwise transaction is committed.

#### `Query(model interface{}, query interface{}, params ...interface{}) (pg.Result, error)`
Query executes a query that returns rows, typically a SELECT. The params are for
any placeholders in the query.

#### `QueryOne(model interface{}, query interface{}, params ...interface{}) (pg.Result, error)`
QueryOne acts like Query, but query must return only one row. It returns
ErrNoRows error when query returns zero rows or ErrMultiRows when query
returns multiple rows.

### DB

The `testutils.DB` interface includes the `testutils.BaseDB` interface.

### Tx
#### `Query(model interface{}, query interface{}, params ...interface{}) (pg.Result, error)`

Query is an alias for DB.Query

#### `QueryOne(model interface{}, query interface{}, params ...interface{}) (pg.Result, error)`

QueryOne is an alias for DB.QueryOne

#### `Select(model interface{}) error`

Select is an alias for DB.Select

#### `Insert(model ...interface{}) error`

Insert is an alias for DB.Insert

#### `Update(model interface{}) error`

Update is an alias for DB.Update

#### `Delete(model interface{}) error`

Delete is an alias for DB.Delete

Additional Functions
--------------------

#### `func NewMockDB() *MockDB`

NewMockDB creates a new mock database client of unit tests.

#### `func (db *MockDB) QueueResponses(response ...interface{})`

QueueResponses inserts data of any type into the mock database. Data will be
returned, one value at a time, to the Query, QueryOne, and Select functions in
the order it was inserted without any attempt to parse the query or match
conditions or IDs. Data inserted into QueueResponses does not need to implement
the `testutils.Model` interface.

#### `func (db *MockDB) QueueModels(model ...Model)`

QueueModels inserts structs into the mock database without using MockDB.Insert
to set up a test including Update or Delete calls. Values passed to QueueModels
must implement the `testutils.Model` interface.

#### `func (db *MockDB) Find(model Model) (Model, error)`

Find returns the value in the mock database models that matches the type and ID
of the provided model if it exists and nil if it doesn't.

#### `func (db *MockDB) MarshalModels() (string, error)`

MarshalModels returns an indented string of JSON for logging out the contents of
the MockDB models. It isn't useful for automated tests, but can give the
developer a way to see what data has been stored.

#### `func (db *MockDB) MarshalResponses() (string, error)`

MarshalResponses returns an indented string of JSON for logging out the contents
of the MockDB models. It isn't useful for automated tests, but can give the
developer a way to see what data has queued for response.

#### `func Diff(a interface{}, b interface{}) (map[string][]interface{}, error)`

Diff compares two values of the same type. If they are different types, an
error is returned. If they are the same type, it returns a map of names of
field struct fields that did not match as keys and a slice containing the
unequal field values as values
