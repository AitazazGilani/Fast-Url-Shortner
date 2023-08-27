package model

import(
	"fmt"
    "gorm.io/gorm"                    // GORM ORM
    "gorm.io/driver/postgres"         // PostgreSQL driver for GORM
)


type shortURL struct {
	ID			uint64	`json:"id" gorm:"primaryKey"`
	Redirect	string	`json:"redirect" gorm:"not null"`
	NewURL		string	`json:"newurl" gorm:"unique;not null"`
	Clicked		uint64	`json:"clicked"`
	Random		bool	`json:"random"`
}

var db *gorm.DB

func Setup() {
	fmt.Println("setting up database")
	// Define the PostgreSQL connection string
	dsn := "user=your_database_user password=your_database_password dbname=your_database_name host=localhost port=5432 sslmode=disable TimeZone=UTC"
	
	var err error
	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	err = db.AutoMigrate(&shortURL{})
	if err != nil {
		fmt.Println(err)
	}

}