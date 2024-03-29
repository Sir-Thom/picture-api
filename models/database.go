package models

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
)

func Database() (*gorm.DB, error) {

	dsn := "host=192.168.1.67 user=server password=$458TT*#17 dbname=hentai port=5432 "
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return db, nil

}
