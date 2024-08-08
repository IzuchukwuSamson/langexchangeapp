package services

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/IzuchukwuSamson/lexi/internal/app/users/models"
	"github.com/IzuchukwuSamson/lexi/utils"
)

type UserService struct {
	DB  *sql.DB
	ctx context.Context
	log *log.Logger
}

func NewUserService(db *sql.DB) UserServiceInterface {
	return &UserService{
		DB:  db,
		ctx: context.TODO(),
	}
}

// GetUserById implements UserInterface.
func (u *UserService) FetchUserById(id string) (*models.User, error) {
	var result models.User
	query := "SELECT id, username, email FROM users WHERE id = ?"
	err := u.DB.QueryRow(query, id).Scan(&result.ID, &result.Username, &result.Email)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// GetAllUsers implements UserInterface.
func (u *UserService) FetchAllUsers() ([]models.User, error) {
	query := `
        SELECT id, firstname, lastname, username, email, email_verified, password, phonenumber, role, last_active, created_at, updated_at 
        FROM users
    `
	// Execute the query
	rows, err := u.DB.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch users: %v", err)
	}
	defer rows.Close()

	// Prepare to store the results
	var users []models.User

	// Iterate over the result set
	for rows.Next() {
		var user models.User
		// Scan the result into the user struct
		err := rows.Scan(
			&user.ID,
			&user.FirstName,
			&user.LastName,
			&user.Username,
			&user.Email,
			&user.EmailVerified,
			&user.Password,
			&user.PhoneNumber,
			&user.Role,
			&user.LastActive,
			&user.CreatedAt,
			&user.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan user: %v", err)
		}
		users = append(users, user)
	}

	// Check for errors after iterating over the result set
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over results: %v", err)
	}

	return users, nil
}

func (u *UserService) CreateUser(user models.User) (*models.User, string, error) {
	// Transaction Begin: Start a new transaction with tx, err := u.DB.Begin().
	tx, err := u.DB.Begin()
	if err != nil {
		return nil, "", fmt.Errorf("failed to begin transaction: %v", err)
	}

	// Insert User: Insert the user data using tx.Exec() within the transaction context.
	query := `
        INSERT INTO users (firstname, lastname, username, email, email_verified, password, phonenumber, role, last_active, created_at, updated_at) 
        VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
    `
	result, err := tx.Exec(query, user.FirstName, user.LastName, user.Username, user.Email, user.EmailVerified, user.Password, user.PhoneNumber, user.Role, time.Now(), time.Now(), time.Now())
	if err != nil {
		// Transaction Rollback: Rollback the transaction with tx.Rollback() if any operation fails.
		tx.Rollback()
		return nil, "", fmt.Errorf("failed to insert user: %v", err)
	}

	// Retrieve User ID: Retrieve the last inserted user ID with result.LastInsertId().
	userID, err := result.LastInsertId()
	if err != nil {
		// Transaction Rollback: Rollback the transaction with tx.Rollback() if any operation fails.
		tx.Rollback()
		return nil, "", fmt.Errorf("failed to retrieve inserted user ID: %v", err)
	}

	user.ID = strconv.FormatInt(userID, 10)

	// Generate a verification code
	verificationCode := utils.GeneratePIN()

	// Insert Email Verification: Insert the email verification record within the same transaction.
	verificationQuery := `
        INSERT INTO email_verifications (user_id, email, code) 
        VALUES (?, ?, ?)
    `
	_, err = tx.Exec(verificationQuery, user.ID, user.Email, strconv.Itoa(verificationCode))
	if err != nil {
		// Transaction Rollback: Rollback the transaction with tx.Rollback() if any operation fails.
		tx.Rollback()
		return nil, "", fmt.Errorf("failed to insert email verification: %v", err)
	}

	// Transaction Commit: Commit the transaction with tx.Commit() if all operations succeed.
	if err := tx.Commit(); err != nil {
		// Transaction Rollback: Rollback the transaction with tx.Rollback() if any operation fails.
		tx.Rollback()
		return nil, "", fmt.Errorf("failed to commit transaction: %v", err)
	}

	return &user, strconv.Itoa(verificationCode), nil
}

func (u *UserService) VerifyEmail(verificationCode string) error {
	// Find the email verification record based on the verification code
	var emailVerificationID int
	var userID int
	query := "SELECT id, user_id FROM email_verifications WHERE code = ?"
	err := u.DB.QueryRow(query, verificationCode).Scan(&emailVerificationID, &userID)
	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("invalid or expired verification code for email")
		}
		u.log.Printf("error finding email verification: %v\n", err)
		return fmt.Errorf("internal server error")
	}

	// Update the user's IsActive status
	updateQuery := "UPDATE users SET email_verified = 1 WHERE id = ?"
	_, err = u.DB.Exec(updateQuery, userID)
	if err != nil {
		u.log.Printf("error updating user: %v\n", err)
		return fmt.Errorf("internal server error")
	}

	// Delete the email verification record
	deleteQuery := "DELETE FROM email_verifications WHERE id = ?"
	_, err = u.DB.Exec(deleteQuery, emailVerificationID)
	if err != nil {
		u.log.Printf("error deleting email verification record: %v\n", err)
		return fmt.Errorf("internal server error")
	}

	return nil
}

func (u *UserService) FindOrCreateUser(userInfo map[string]interface{}) (*models.User, error) {
	var user models.User

	// Check if the user already exists
	query := "SELECT id, firstname, lastname, username, email, email_verified, password, phonenumber, role, last_active, created_at, updated_at FROM users WHERE email = ?"
	err := u.DB.QueryRow(query, userInfo["email"]).Scan(
		&user.ID,
		&user.FirstName,
		&user.LastName,
		&user.Username,
		&user.Email,
		&user.EmailVerified,
		&user.Password,
		&user.PhoneNumber,
		&user.Role,
		&user.LastActive,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err == nil {
		// User already exists
		return &user, nil
	}
	if err != nil && err != sql.ErrNoRows {
		// Some error occurred other than "no rows found"
		return nil, fmt.Errorf("failed to check user existence: %v", err)
	}

	// User does not exist, create a new one
	query = `
        INSERT INTO users (firstname, lastname, username, email, email_verified, password, phonenumber, role, last_active, created_at, updated_at) 
        VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
    `
	result, err := u.DB.Exec(query,
		userInfo["firstname"],
		userInfo["lastname"],
		userInfo["username"],
		userInfo["email"],
		userInfo["email_verified"],
		userInfo["password"],
		userInfo["phonenumber"],
		userInfo["role"],
		time.Now(),
		time.Now(),
		time.Now(),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %v", err)
	}

	// Get the inserted user ID
	userID, err := result.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve inserted user ID: %v", err)
	}

	user.ID = strconv.FormatInt(userID, 10)
	// user.ID = userID
	user.FirstName = userInfo["firstname"].(string)
	user.LastName = userInfo["lastname"].(string)
	user.Username = userInfo["username"].(string)
	user.Email = userInfo["email"].(string)
	user.EmailVerified = userInfo["email_verified"].(int)
	user.Password = userInfo["password"].(string)
	user.PhoneNumber = userInfo["phonenumber"].(string)
	user.Role = userInfo["role"].(string)
	user.LastActive = time.Now()
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()

	return &user, nil
}

func (u *UserService) UpdateUser(updatedUser *models.User) (*models.User, error) {
	query := `
        UPDATE users 
        SET firstname = ?, lastname = ?, username = ?, email = ?, password = ?, phonenumber = ?, role = ?, last_active = ?, updated_at = ?
        WHERE id = ?
    `
	_, err := u.DB.Exec(query,
		updatedUser.FirstName,
		updatedUser.LastName,
		updatedUser.Username,
		updatedUser.Email,
		updatedUser.Password,
		updatedUser.PhoneNumber,
		updatedUser.Role,
		time.Now(),
		time.Now(),
		updatedUser.ID,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to update user: %v", err)
	}

	// Fetch the updated user from the database to return
	var user models.User
	selectQuery := `
        SELECT id, firstname, lastname, username, email, email_verified, password, phonenumber, role, last_active, created_at, updated_at 
        FROM users 
        WHERE id = ?
    `
	err = u.DB.QueryRow(selectQuery, updatedUser.ID).Scan(
		&user.ID,
		&user.FirstName,
		&user.LastName,
		&user.Username,
		&user.Email,
		&user.EmailVerified,
		&user.Password,
		&user.PhoneNumber,
		&user.Role,
		&user.LastActive,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve updated user: %v", err)
	}

	return &user, nil
}

func (u *UserService) Update(user *models.User) error {
	query := `
        UPDATE users 
        SET email = ?, firstname = ?, lastname = ?, username = ?, phonenumber = ?, role = ?, last_active = ?, updated_at = ?
        WHERE id = ?
    `
	_, err := u.DB.Exec(query,
		user.Email,
		user.FirstName,
		user.LastName,
		user.Username,
		user.PhoneNumber,
		user.Role,
		time.Now(),
		time.Now(),
		user.ID,
	)
	return err
}

func (u *UserService) GetUserByEmail(email string) (*models.User, error) {
	var user models.User
	var deletedAt sql.NullTime // Use sql.NullTime to handle possible NULL values

	query := `SELECT id, firstname, lastname, email, email_verified, password, username, phonenumber, last_active, created_at, updated_at, deleted_at FROM users WHERE email = ?`

	// Execute the query
	row := u.DB.QueryRow(query, email)

	// Scan the result into the user struct
	err := row.Scan(
		&user.ID,
		&user.FirstName,
		&user.LastName,
		&user.Email,
		&user.EmailVerified,
		&user.Password,
		&user.Username,
		&user.PhoneNumber,
		&user.LastActive,
		&user.CreatedAt,
		&user.UpdatedAt,
		&deletedAt,
	)

	// Handle errors
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user with email %s not found", email)
		}
		log.Printf("error querying user: %v\n", err)
		return nil, fmt.Errorf("internal server error")
	}

	return &user, nil
}

// UpdateUser implements UserInterface.
func (u *UserService) EditUser(id string, updatedUser models.User) (*models.User, error) {
	// Prepare the SQL query to update user information
	query := `
        UPDATE users 
        SET firstname = ?, lastname = ?, email = ?, password = ?, phonenumber = ?
        WHERE id = ?
    `

	// Execute the query
	_, err := u.DB.Exec(query,
		updatedUser.FirstName,
		updatedUser.LastName,
		updatedUser.Email,
		updatedUser.Password,
		updatedUser.PhoneNumber,
		id,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to update user: %v", err)
	}

	// Return the updated user
	return &updatedUser, nil
}

// DeleteUser implements UserInterface.
func (u *UserService) RemoveUser(id string) (*models.User, error) {
	// Convert ID string to integer if needed (MySQL usually uses integer IDs)
	// Assuming `id` is an integer in the database
	userID, err := strconv.Atoi(id)
	if err != nil {
		return nil, fmt.Errorf("invalid user ID: %v", err)
	}

	// Find the user to be removed
	var deletedUser models.User
	query := "SELECT id, firstname, lastname, email, email_verified, password, phonenumber, role, last_active, created_at, updated_at FROM users WHERE id = ?"
	err = u.DB.QueryRow(query, userID).Scan(
		&deletedUser.ID,
		&deletedUser.FirstName,
		&deletedUser.LastName,
		&deletedUser.Email,
		&deletedUser.EmailVerified,
		&deletedUser.Password,
		&deletedUser.PhoneNumber,
		&deletedUser.Role,
		&deletedUser.LastActive,
		&deletedUser.CreatedAt,
		&deletedUser.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("failed to retrieve user: %v", err)
	}

	// Delete the user
	deleteQuery := "DELETE FROM users WHERE id = ?"
	_, err = u.DB.Exec(deleteQuery, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to delete user: %v", err)
	}

	return &deletedUser, nil
}

func (u *UserService) GeneratePasswordResetToken(user models.User) (string, error) {
	resetToken := utils.GeneratePIN()
	expirationTime := time.Now().Add(1 * time.Hour)

	// Insert the password reset token into the password_resets table
	query := `
		INSERT INTO password_resets (user_id, code, expires_at) 
		VALUES (?, ?, ?)
	`
	_, err := u.DB.Exec(query, user.ID, strconv.Itoa(resetToken), expirationTime)
	if err != nil {
		u.log.Printf("error inserting password reset token: %v\n", err)
		return "", fmt.Errorf("failed to insert password reset token: %v", err)
	}

	return strconv.Itoa(resetToken), nil
}

func (u *UserService) GetPasswordResetByCode(code string) (*models.PasswordReset, error) {
	var passwordReset models.PasswordReset

	// Query to find the password reset token by code
	query := `
		SELECT id, user_id, code, expires_at 
		FROM password_resets 
		WHERE code = ?
	`

	row := u.DB.QueryRow(query, code)

	// Scan the result into the PasswordReset model
	err := row.Scan(
		&passwordReset.ID,
		&passwordReset.UserID,
		&passwordReset.Code,
		&passwordReset.ExpiresAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			// Handle case where no row was found
			return nil, fmt.Errorf("no password reset found for code: %s", code)
		}
		// Handle other SQL errors
		return nil, fmt.Errorf("failed to retrieve password reset: %v", err)
	}

	return &passwordReset, nil
}

func (u *UserService) UpdateUserPassword(user *models.User) error {
	// Define the update query
	query := `
		UPDATE users 
		SET password = ? 
		WHERE id = ?
	`

	// Execute the update query
	_, err := u.DB.Exec(query, user.Password, user.ID)
	if err != nil {
		return fmt.Errorf("failed to update user password: %v", err)
	}

	return nil
}

func (u *UserService) DeletePasswordReset(id int64) error {
	// Define the delete query
	query := `
		DELETE FROM password_resets 
		WHERE id = ?
	`

	// Execute the delete query
	_, err := u.DB.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete password reset record: %v", err)
	}

	return nil
}

func (u *UserService) InvalidateToken(token string) error {
	return nil
}

// Optionally add a method to check if a token is invalidated
func (u *UserService) IsTokenInvalidated(token string) bool {
	return false
}
