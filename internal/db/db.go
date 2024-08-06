package db

import (
	"context"
	"database/sql"
	"fmt"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type DB struct {
	Mongo *mongo.Database
	Sql   *sql.DB
}

func NewDB(mongoDB *mongo.Database, sqlDB *sql.DB) *DB {
	return &DB{
		Mongo: mongoDB,
		Sql:   sqlDB,
	}
}

// Mongo initializes a mongo db connection.
func Mongo() *mongo.Database {
	uri := os.Getenv("MONGO_URI")

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
	return db
}

func SQL() (*sql.DB, error) {
	var db *sql.DB
	var err error

	dbDriver := os.Getenv("DB_DRIVER")
	switch dbDriver {
	case "sqlite":
		db, err = sql.Open("sqlite", os.Getenv("DB_SQLITE"))
	case "postgres":
		db, err = sql.Open("postgres", os.Getenv("PG_DSN"))
	case "mysql":
		dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?parseTime=True",
			os.Getenv("DB_USER"), os.Getenv("DB_PASS"),
			os.Getenv("DB_HOST"), os.Getenv("DB_NAME"))
		db, err = sql.Open("mysql", dsn)
	default:
		err = fmt.Errorf("unsupported db driver: %s", dbDriver)
	}

	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}
