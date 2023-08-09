package dbo

import (
	"echoapptpl/model"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

type DBO struct {
	*gorm.DB
}

func New(dsn string) *DBO {
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
			NoLowerCase:   true,
		},
		DisableForeignKeyConstraintWhenMigrating: true,
		Logger:                                   logger.Default.LogMode(logger.Info),
	})

	if err != nil {
		panic(err)
	}
	return &DBO{db}
}

func (dbo *DBO) Sync() {
	if err := dbo.DB.AutoMigrate(&model.User{}); err != nil {
		panic(err)
	}
}
