package models

import "github.com/jinzhu/gorm"

// NewDB Creates a new gorm database connection
func NewDB(dialect string, datasource string) (*gorm.DB, error) {
	db, err := gorm.Open(dialect, datasource)
	if err != nil {
		return nil, err
	}
	return db, err
}

// InitTables creates some testing tables for development
func InitTables(db *gorm.DB) {
	if db.HasTable(&Channel{}) {
		db.DropTable(&Channel{})
	}

	db.CreateTable(&Channel{})
	db.Create(&Channel{Username: "g33kidd", StreamKey: "sk_g33kidd_13245"})
}
