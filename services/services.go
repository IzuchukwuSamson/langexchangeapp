package services

import (
	"context"
	"errors"
	"fmt"
	"lexibuddy/models"
	"lexibuddy/utils"
	"log"
	"strconv"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type UserService struct {
	userCollection          *mongo.Collection
	emailVerification       *mongo.Collection
	passwordResetCollection *mongo.Collection
	DB                      *mongo.Database
	ctx                     context.Context
	tokenBlacklist          map[string]bool
	log                     *log.Logger
}

func NewUserService(db *mongo.Database, l *log.Logger) UserInterface {
	return &UserService{
		DB:                      db,
		userCollection:          db.Collection(utils.UsersTable),
		emailVerification:       db.Collection(utils.EmailVerificationTable),
		passwordResetCollection: db.Collection(utils.PasswordReset),
		ctx:                     context.TODO(),
		log:                     l,
	}
}

// GetUserById implements UserInterface.
func (u *UserService) FetchUserById(id string) (*models.User, error) {
	var result models.User
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	err = u.userCollection.FindOne(u.ctx, bson.M{"_id": objectId}).Decode(&result)
	return &result, err
}

// GetAllUsers implements UserInterface.
func (u *UserService) FetchAllUsers() ([]models.User, error) {
	var users []models.User
	res, err := u.userCollection.Find(u.ctx, bson.M{})
	if err != nil {
		return nil, err
	}

	return users, res.All(u.ctx, &users)
}

func (u *UserService) CreateUser(user models.User) (*models.User, string, error) {
	res, err := u.userCollection.InsertOne(u.ctx, user)
	if err != nil {
		return nil, "", err
	}

	// Assign the generated ID to the user object
	user.ID = res.InsertedID.(primitive.ObjectID)

	// Generate a verification PIN
	verificationCode := utils.GeneratePIN()

	// Create EmailVerification object
	emailVerification := models.EmailVerification{
		UserID: user.ID,
		Email:  user.Email,
		Code:   strconv.Itoa(verificationCode),
	}

	// Insert email verification into emailVerification collection
	_, err = u.emailVerification.InsertOne(u.ctx, emailVerification)
	if err != nil {
		// Rollback user creation if email verification fails
		_, rollbackErr := u.userCollection.DeleteOne(u.ctx, bson.M{"_id": user.ID})
		if rollbackErr != nil {
			u.log.Printf("Error rolling back user creation: %v", rollbackErr)
		}
		return nil, "", fmt.Errorf("failed to verify email: %v", err)
	}

	return &user, emailVerification.Code, nil
}

func (u *UserService) VerifyEmail(verificationCode string) error {
	// Find the email verification record based on verification code and email
	var emailVerification models.EmailVerification
	err := u.emailVerification.FindOne(u.ctx, bson.M{"code": verificationCode}).Decode(&emailVerification)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return fmt.Errorf("invalid or expired verification code for email")
		}
		u.log.Printf("error finding email verification: %v\n", err)
		return fmt.Errorf("internal server error")
	}

	// Update the user's IsActive status
	filter := bson.M{"_id": emailVerification.UserID}
	update := bson.M{"$set": bson.M{"is_active": 1}}
	_, err = u.userCollection.UpdateOne(u.ctx, filter, update)
	if err != nil {
		u.log.Printf("error updating user: %v\n", err)
		return fmt.Errorf("internal server error")
	}

	// Delete the email verification record
	_, err = u.emailVerification.DeleteOne(u.ctx, bson.M{"_id": emailVerification.ID})
	if err != nil {
		u.log.Printf("error deleting email verification record: %v\n", err)
		return fmt.Errorf("internal server error")
	}

	return nil
}

func (u *UserService) FindOrCreateUser(userInfo map[string]interface{}) (*models.User, error) {
	email, ok := userInfo["email"].(string)
	if !ok {
		return nil, errors.New("email not found in user info")
	}

	// Find user by email
	user, err := u.GetUserByEmail(email)
	if err != nil {
		return nil, err
	}

	if user != nil {
		// User exists, return the user
		return user, nil
	}

	// User does not exist, create a new user
	newUser := models.User{
		ID:    primitive.NewObjectID(),
		Email: email,
		// CreatedAt: time.Now(),
	}

	// 	Populate additional fields if available
	// if name, ok := userInfo["name"].(string); ok {
	// 	newUser.Name = name
	// }
	// if picture, ok := userInfo["picture"].(string); ok {
	// 	newUser.Picture = picture
	// }

	createdUser, _, err := u.CreateUser(newUser)
	if err != nil {
		return nil, err
	}

	return createdUser, nil
}

func (u *UserService) UpdateUser(updatedUser *models.User) (*models.User, error) {
	// Find the existing user by ID
	existingUser, err := u.FetchUserById(updatedUser.ID.Hex())
	if err != nil {
		return nil, err
	}
	if existingUser == nil {
		return nil, errors.New("user not found")
	}

	// Update the user fields
	existingUser.Email = updatedUser.Email
	// existingUser.Name = updatedUser.Name
	// existingUser.Picture = updatedUser.Picture
	// existingUser.UpdatedAt = time.Now()

	// Save the updated user to the database
	// Save the updated user to the database
	err = u.Update(existingUser)
	if err != nil {
		return nil, err
	}

	return existingUser, nil
}

func (u *UserService) Update(user *models.User) error {
	_, err := u.userCollection.UpdateOne(
		u.ctx,
		bson.M{"_id": user.ID},
		bson.M{"$set": bson.M{
			"email": user.Email,
			// "name":       user.Name,
			// "picture":    user.Picture,
			// "updated_at": user.UpdatedAt,
		}},
		options.Update().SetUpsert(true),
	)
	return err
}

func (u *UserService) GetUserByEmail(email string) (*models.User, error) {
	var result models.User
	err := u.userCollection.FindOne(u.ctx, bson.M{"email": email}).Decode(&result)
	return &result, err
}

// UpdateUser implements UserInterface.
func (u *UserService) EditUser(id string, updatedUser models.User) (*models.User, error) {
	// Convert ID string to primitive.ObjectID
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	// Define update fields
	update := bson.M{
		"$set": bson.M{
			"firstname":   updatedUser.FirstName,
			"lastname":    updatedUser.LastName,
			"email":       updatedUser.Email,
			"password":    updatedUser.Password,
			"phonenumber": updatedUser.PhoneNumber,
		},
	}

	// Perform the update operation
	_, err = u.userCollection.UpdateOne(u.ctx, bson.M{"_id": objectID}, update)
	if err != nil {
		return nil, err
	}

	// Return the updated user
	return &updatedUser, nil

}

// DeleteUser implements UserInterface.
func (u *UserService) RemoveUser(id string) (*models.User, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	// Find the user to be removed
	var deletedUser models.User
	err = u.userCollection.FindOneAndDelete(u.ctx, bson.M{"_id": objectID}).Decode(&deletedUser)
	if err != nil {
		return nil, err
	}

	return &deletedUser, nil
}

func (u *UserService) GeneratePasswordResetToken(user models.User) (string, error) {
	resetToken := utils.GeneratePIN()

	expirationTime := time.Now().Add(1 * time.Hour)

	passwordReset := models.PasswordReset{
		UserID:    user.ID,
		Code:      strconv.Itoa(resetToken),
		ExpiresAt: expirationTime,
	}

	_, err := u.passwordResetCollection.InsertOne(u.ctx, passwordReset)
	if err != nil {
		u.log.Printf("error inserting password reset token: %v\n", err)
		return "", err
	}

	return strconv.Itoa(resetToken), nil
}

func (u *UserService) GetPasswordResetByCode(code string) (*models.PasswordReset, error) {
	var passwordReset models.PasswordReset
	err := u.passwordResetCollection.FindOne(u.ctx, bson.M{"code": code}).Decode(&passwordReset)
	if err != nil {
		return nil, err
	}
	return &passwordReset, nil
}

func (u *UserService) UpdateUserPassword(user *models.User) error {
	// Update the user's password
	_, err := u.userCollection.UpdateOne(
		u.ctx,
		bson.M{"_id": user.ID},
		bson.M{"$set": bson.M{"password": user.Password}},
	)
	if err != nil {
		return err
	}

	return err
}

func (u *UserService) DeletePasswordReset(id primitive.ObjectID) error {
	_, err := u.passwordResetCollection.DeleteOne(u.ctx, bson.M{"_id": id})
	if err != nil {
		return fmt.Errorf("failed to delete password reset record: %v", err)
	}
	return nil
}

func (u *UserService) InvalidateToken(token string) error {
	// Add the token to the blacklist
	u.tokenBlacklist[token] = true
	return nil
}

// Optionally add a method to check if a token is invalidated
func (u *UserService) IsTokenInvalidated(token string) bool {
	return u.tokenBlacklist[token]
}
