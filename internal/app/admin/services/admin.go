package services

import (
	"context"
	"database/sql"
	"errors"
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

func (a *AdminService) GetAdminByEmail(ad string) (*adminModel.Admin, error) {
	var admin adminModel.Admin

	// Prepare the query to select the admin by email
	query := "SELECT id, name, email, password FROM admins WHERE email = ?"

	// Execute the query
	err := a.DB.QueryRow(query, admin.Email).Scan(&admin.ID, &admin.Name, &admin.Email, &admin.Password)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("admin not found")
		}
		return nil, err
	}

	return &admin, nil

}

// UpdateAdminPassword updates the password of an existing admin in the database.
func (s *AdminService) UpdateAdminPassword(admin *adminModel.Admin) error {
	// Construct the SQL query for updating the password
	query := `UPDATE admins SET password = ? WHERE email = ?`

	// Execute the query with the hashed password and the admin's email
	_, err := s.DB.Exec(query, admin.Password, admin.Email)
	if err != nil {
		return fmt.Errorf("error updating admin password: %w", err)
	}

	return nil
}

// func (a *AdminService) CreateAdmin(admin *adminModel.Admin) error {
// 	// Prepare the SQL statement for inserting a new admin
// 	query := "INSERT INTO admins (id, name, email, password) VALUES (?, ?, ?, ?)"

// 	// Execute the insertion
// 	_, err := a.DB.Exec(query, admin.ID, admin.Name, admin.Email, admin.Password)
// 	if err != nil {
// 		return err
// 	}

// 	return nil
// }

// FetchAllUsers implements AdminServiceInterface.
func (u *AdminService) FetchAllUsers() ([]userModel.User, error) {
	u.log.Println("Fetch All Users")
	return nil, nil
}
