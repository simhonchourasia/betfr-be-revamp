package dbinterface

import (
	"database/sql"

	"gorm.io/gorm"
)

// DBInterface is an interface for gorm.DB to enable mocking.
type DBInterface interface {
	Exec(sql string, values ...interface{}) *gorm.DB
	Find(dest interface{}, conds ...interface{}) *gorm.DB
	First(dest interface{}, conds ...interface{}) *gorm.DB
	Create(value interface{}) *gorm.DB
	Save(value interface{}) *gorm.DB
	Delete(value interface{}, conds ...interface{}) *gorm.DB
	Where(query interface{}, args ...interface{}) *gorm.DB
	Preload(query string, args ...interface{}) *gorm.DB
	Begin(...*sql.TxOptions) *gorm.DB
	Commit() *gorm.DB
	Rollback() *gorm.DB
	Raw(sql string, values ...interface{}) *gorm.DB
	Model(value interface{}) *gorm.DB
	// Scopes(funcs ...func(*gorm.DB) *gorm.DB) *gorm.DB
}

// GormDB is a wrapper around gorm.DB to implement DBInterface.
type GormDB struct {
	*gorm.DB
}

// New creates a new instance of GormDB.
func New(db *gorm.DB) *GormDB {
	return &GormDB{DB: db}
}

// Ensure GormDB implements DBInterface.
var _ DBInterface = (*GormDB)(nil)
