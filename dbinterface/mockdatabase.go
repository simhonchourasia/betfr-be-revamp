package dbinterface

import (
	"database/sql"

	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

// MockDB is a mock implementation of DBInterface.
type MockDB struct {
	mock.Mock
}

// Implement the methods of DBInterface for MockDB.

func (m *MockDB) Exec(sql string, values ...interface{}) *gorm.DB {
	argsMock := m.Called(sql, values)
	return argsMock.Get(0).(*gorm.DB)
}

func (m *MockDB) Find(dest interface{}, conds ...interface{}) *gorm.DB {
	argsMock := m.Called(dest, conds)
	return argsMock.Get(0).(*gorm.DB)
}

func (m *MockDB) First(dest interface{}, conds ...interface{}) *gorm.DB {
	argsMock := m.Called(dest, conds)
	return argsMock.Get(0).(*gorm.DB)
}

func (m *MockDB) Create(value interface{}) *gorm.DB {
	argsMock := m.Called(value)
	return argsMock.Get(0).(*gorm.DB)
}

func (m *MockDB) Save(value interface{}) *gorm.DB {
	argsMock := m.Called(value)
	return argsMock.Get(0).(*gorm.DB)
}

func (m *MockDB) Delete(value interface{}, conds ...interface{}) *gorm.DB {
	argsMock := m.Called(value, conds)
	return argsMock.Get(0).(*gorm.DB)
}

func (m *MockDB) Where(query interface{}, args ...interface{}) *gorm.DB {
	argsMock := m.Called(query, args)
	return argsMock.Get(0).(*gorm.DB)
}

func (m *MockDB) Preload(query string, args ...interface{}) *gorm.DB {
	argsMock := m.Called(query, args)
	return argsMock.Get(0).(*gorm.DB)
}

func (m *MockDB) Begin(opts ...*sql.TxOptions) *gorm.DB {
	argsMock := m.Called(opts)
	return argsMock.Get(0).(*gorm.DB)
}

func (m *MockDB) Commit() *gorm.DB {
	argsMock := m.Called()
	return argsMock.Get(0).(*gorm.DB)
}

func (m *MockDB) Rollback() *gorm.DB {
	argsMock := m.Called()
	return argsMock.Get(0).(*gorm.DB)
}

func (m *MockDB) Raw(sql string, values ...interface{}) *gorm.DB {
	argsMock := m.Called(sql, values)
	return argsMock.Get(0).(*gorm.DB)
}

func (m *MockDB) Model(value interface{}) *gorm.DB {
	argsMock := m.Called(value)
	return argsMock.Get(0).(*gorm.DB)
}

// Ensure GormDB implements DBInterface.
var _ DBInterface = (*MockDB)(nil)
