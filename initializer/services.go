package initializer

import (
	"database/sql"

	adminServices "github.com/IzuchukwuSamson/lexi/internal/app/admin/services"
	userServices "github.com/IzuchukwuSamson/lexi/internal/app/users/services"
	"github.com/IzuchukwuSamson/lexi/internal/db"
	"go.mongodb.org/mongo-driver/mongo"
)

type Store struct {
	User  userServices.UserServiceInterface
	Admin adminServices.AdminServiceInterface
	// Redis *redis.Client
}

func newStoreDB(sqldb *sql.DB, mongodb *mongo.Database) *Store {
	store := Store{
		User:  userServices.NewUserService(sqldb),
		Admin: adminServices.NewAdminService(sqldb),
	}
	return &store
}

func Services(db *db.DB) *Store {
	switch {
	case db.Mongo != nil:
		store := newStoreDB(nil, db.Mongo)
		// store.Redis = redis
		return store
	case db.Sql != nil:
		store := newStoreDB(db.Sql, nil)
		// store.Redis = redis
		return store
	default:
		panic("no database was set")
	}
}
