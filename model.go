package ginCore

import (
	"fmt"
	"gorm.io/gorm"
	"time"
)

type BaseModel struct {
}

func (m *BaseModel) BeforeSave(tx *gorm.DB) (err error) {
	fmt.Println("before save")
	return
}

func (m *BaseModel) AfterSave(tx *gorm.DB) (err error) {
	fmt.Println("after save")
	return
}

func (m *BaseModel) BeforeCreate(tx *gorm.DB) (err error) {
	now := time.Now()
	tx.Statement.SetColumn("CreatedAt", now)
	tx.Statement.SetColumn("UpdatedAt", now)
	tx.Statement.SetColumn("Flag", 1)
	tx.Statement.SetColumn("State", 1)
	tx.Statement.SetColumn("Status", 1)
	return
}

func (m *BaseModel) AfterCreate(tx *gorm.DB) (err error) {

	return
}

func (m *BaseModel) BeforeUpdate(tx *gorm.DB) (err error) {
	now := time.Now()
	tx.Statement.SetColumn("UpdatedAt", now)
	return
}

func (m *BaseModel) AfterUpdate(tx *gorm.DB) (err error) {

	return
}
