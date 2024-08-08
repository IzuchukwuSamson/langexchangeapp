package services

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"strconv"
	"time"

	adminModel "github.com/IzuchukwuSamson/lexi/internal/app/admin/models"
	userModel "github.com/IzuchukwuSamson/lexi/internal/app/users/models"
	"github.com/IzuchukwuSamson/lexi/utils"
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

func (a *AdminService) AdminEmailExists(email string) (bool, error) {
	// Create a context with a 5-second timeout
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Define the query to check if the email exists in the admins table
	query := "SELECT EXISTS(SELECT 1 FROM admins WHERE email = ?)"

	var exists bool
	// Use QueryRowContext to execute the query within the context
	err := a.DB.QueryRowContext(ctx, query, email).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("error checking if email exists: %v", err)
	}

	return exists, nil
}

// NewAdminEmail implements AdminServiceInterface.
func (a *AdminService) GetUserRoleByID(adminID int) (string, error) {
	// Create a context with a 5-second timeout
	ctx, cancel := context.WithTimeout(a.ctx, 5*time.Second)
	defer cancel()

	// Define the query to get the user's role by ID
	query := "SELECT role FROM admins WHERE id = ?"

	var role string
	// Use QueryRowContext to execute the query within the context
	err := a.DB.QueryRowContext(ctx, query, adminID).Scan(&role)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", fmt.Errorf("no user found with ID %d", adminID)
		}
		return "", fmt.Errorf("error getting user role: %v", err)
	}

	return role, nil
}

func (a *AdminService) NewAdminEmail(ad adminModel.Admin) (*adminModel.Admin, error) {
	if ad.Role == "" {
		ad.Role = utils.AdminRole
	}
	// Create a context with a 5-second timeout
	ctx, cancel := context.WithTimeout(a.ctx, 5*time.Second)
	defer cancel()

	query := "INSERT INTO admins (email, role) VALUES (?, ?)"

	result, err := a.DB.ExecContext(ctx, query, ad.Email, ad.Role)
	if err != nil {
		return nil, fmt.Errorf("error inserting new email: %v", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("error getting last insert ID: %v", err)
	}

	ad.ID = strconv.FormatInt(id, 10)

	return &ad, nil
}

// FetchAllUsers implements AdminServiceInterface.
func (u *AdminService) FetchAllUsers() ([]userModel.User, error) {
	u.log.Println("Fetch All Users")
	return nil, nil
}
