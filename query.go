package testutils

import (
	"github.com/go-pg/pg/v9"
)

type Query interface {
	// Clones the current query
	Clone() Query

	// Context(c context.Context) Query

	DB(db DB) Query

	Model(model ...interface{}) Query

	// TableModel() orm.TableModel

	// Deleted adds `WHERE deleted_at IS NOT NULL` clause for soft deleted models.
	Deleted() Query

	AllWithDeleted() Query

	With(name string, subq Query) Query

	WithInsert(name string, subq Query) Query

	WithUpdate(name string, subq Query) Query

	WithDelete(name string, subq Query) Query

	// WrapWith creates new Query and adds to it current query as
	// common table expression with the given name.
	WrapWith(name string) Query

	Table(tables ...string) Query

	TableExpr(expr string, params ...interface{}) Query

	Distinct() Query

	DistinctOn() Query

	Column(columns ...string) Query

	ColumnExpr(expr string, params ...interface{}) Query

	ExcludeColumn(columns ...string) Query

	Relation(name string, apply ...func(Query) (Query, error)) Query

	Set(set string, params ...interface{}) Query

	Value(column string, value string, params ...interface{}) Query

	Where(condition string, params ...interface{}) Query

	WhereOr(condition string, params ...interface{}) Query

	WhereGroup(fn func(Query) (Query, error)) Query

	WhereNotGroup(fn func(Query) (Query, error)) Query

	WhereOrGroup(fn func(Query) (Query, error)) Query

	WhereOrNotGroup(fn func(query Query) (Query, error)) Query

	WhereIn(where string, slice interface{}) Query

	WhereInMulti(where string, values ...interface{}) Query

	WherePK() Query

	WhereStruct(strct interface{}) Query

	Join(join string, params ...interface{}) Query

	JoinOn(condition string, params ...interface{}) Query

	JoinOnOr(condition string, params ...interface{}) Query

	Group(columns ...string) Query

	GroupExpr(group string, params ...interface{}) Query

	Having(having string, params ...interface{}) Query

	Union(other Query) Query

	UnionAll(other Query) Query

	Intersect(other Query) Query

	IntersectAll(other Query) Query

	Except(other Query) Query

	ExceptAll(other Query) Query

	Order(orders ...string) Query

	OrderExpr(order string, params ...interface{}) Query

	Limit(n int) Query

	Offset(n int) Query

	OnConflict(s string, params ...interface{}) Query

	Returning(s string, params ...interface{}) Query

	For(s string, params ...interface{}) Query

	Apply(fn func(Query) (Query, error)) Query

	Count() (int, error)

	First() error

	Last() error

	Select(values ...interface{}) error

	SelectAndCount(values ...interface{}) (count int, firstErr error)

	SelectAndCountEstimate(threshold int, values ...interface{}) (count int, firstErr error)

	// ForEach(fn interface{}) error

	Insert(values ...interface{}) (pg.Result, error)

	SelectOrInsert(values ...interface{}) (inserted bool, _ error)

	Update(scan ...interface{}) (pg.Result, error)

	UpdateNotZero(scan ...interface{}) (pg.Result, error)

	Delete(values ...interface{}) (pg.Result, error)

	ForceDelete(values ...interface{}) (pg.Result, error)

	// CreateTable(opt *orm.CreateTableOptions) error

	// DropTable(opt *orm.DropTableOptions) error

	// Exec(query interface{}, params ...interface{}) (pg.Result, error)

	// ExecOne(query interface{}, params ...interface{}) (pg.Result, error)

	Query(model, query interface{}, params ...interface{}) (pg.Result, error)

	QueryOne(model, query interface{}, params ...interface{}) (pg.Result, error)

	// CopyFrom(r io.Reader, query interface{}, params ...interface{}) (pg.Result, error)

	// CopyTo(w io.Writer, query interface{}, params ...interface{}) (pg.Result, error)

	// AppendQuery(fmter orm.QueryFormatter, b []byte) ([]byte, error)

	Exists() (bool, error)
}
