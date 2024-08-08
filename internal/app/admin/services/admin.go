package services

import (
	"context"
	"database/sql"
	"log"

	"github.com/IzuchukwuSamson/lexi/internal/app/users/models"
)

type AdminService struct {
	DB  *sql.DB
	log *log.Logger
	ctx context.Context
}

// AdminEmailExists implements AdminServiceInterface.

func NewAdminService(db *sql.DB) AdminServiceInterface {
	return &AdminService{
		DB:  db,
		ctx: context.TODO(),
	}
}

func (*AdminService) AdminEmailExists(email string) (bool error) {
	panic("unimplemented")
	// check
}

// NewAdminEmail implements AdminServiceInterface.
func (*AdminService) NewAdminEmail() {
	panic("unimplemented")
}

// FetchAllUsers implements AdminServiceInterface.
func (u *AdminService) FetchAllUsers() ([]models.User, error) {
	u.log.Println("Fetch All Users")
	return nil, nil
}
