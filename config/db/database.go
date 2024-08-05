package db

import (
	"context"
	"fmt"
	"os"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type DB struct {
	Mongo *mongo.Database
	Sql   *gorm.DB
}

func NewDB(mongoDB *mongo.Database, sqlDb *gorm.DB) *DB {
	return &DB{
		Mongo: mongoDB,
		Sql:   sqlDb,
	}
}

// Mongo initializes a mongo db connection.
func Mongo() (*mongo.Database, error) {
	// uri := os.Getenv("MONGO_URI")
	uri := "mongodb://127.0.0.1:27017"

	bsonOpts := &options.BSONOptions{
		UseJSONStructTags: false,
	}
	clientOpts := options.Client().
		ApplyURI(uri).
		SetBSONOptions(bsonOpts)
	client, err := mongo.Connect(context.TODO(), clientOpts)
	if err != nil {
		panic(err)
	}
	if err := client.Ping(context.TODO(), nil); err != nil {
		panic(err)
	}
	db := client.Database(os.Getenv("DATABASE_NAME"))
	return db, nil
}

func SQL() (*gorm.DB, error) {
	var db *gorm.DB
	var err error

	var opts = &gorm.Config{}
	if os.Getenv("DEBUG") == "true" {
		opts = &gorm.Config{
			Logger: logger.Default.LogMode(logger.Info),
		}
	}
	switch os.Getenv("DB_DRIVER") {
	case "postgres":
		db, err = gorm.Open(postgres.Open(os.Getenv("PG_DSN")), opts)
	case "mysql":
		dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?parseTime=True",
			os.Getenv("DB_USER"), os.Getenv("DB_PASS"),
			os.Getenv("DB_HOST"), os.Getenv("DB_NAME"))
		db, err = gorm.Open(mysql.Open(dsn), opts)
	}
	// db.AutoMigrate(&data.User{}, &data.Profile{}, &data.Ideal{}, &data.Social{}, &data.Admin{}, &data.Image{})
	return db, err
}
