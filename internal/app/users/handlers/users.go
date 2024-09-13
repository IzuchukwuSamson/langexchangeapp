package handlers

import (
	"fmt"
	"net/http"

	"github.com/IzuchukwuSamson/lexi/utils"
)

func (u UserHandlers) Dashboard(rw http.ResponseWriter, r *http.Request) {
	fmt.Println("Dashboard route")
}

func (u UserHandlers) GetAllUsers(rw http.ResponseWriter, r *http.Request) {
	users, err := u.services.FetchAllUsers()
	if err != nil {
		u.log.Printf("error decoding json request: %v\n", err)
		utils.ReturnJSON(rw, utils.ErrMessage{Error: "error getting users"}, http.StatusInternalServerError)
		return
	}

	// Wrap the users in a data object
	response := map[string]interface{}{
		"data": users,
	}

	utils.ReturnJSON(rw, response, http.StatusOK)
}

// GetAllUsersById handles the HTTP request to fetch a user by ID
func (u UserHandlers) GetUserById(rw http.ResponseWriter, r *http.Request) {
	// Get the user ID from the query parameters
	userIdStr := r.URL.Query().Get("id")
	if userIdStr == "" {
		http.Error(rw, "No user ID provided", http.StatusBadRequest)
		return
	}

	// Fetch user by ID
	user, err := u.services.FetchUserByID(userIdStr)
	if err != nil {
		u.log.Printf("error fetching user: %v\n", err)
		http.Error(rw, "Error fetching user", http.StatusInternalServerError)
		return
	}

	// Prepare response data
	response := map[string]interface{}{
		"data": user,
	}

	// Return JSON response
	utils.ReturnJSON(rw, response, http.StatusOK)
}

func (u UserHandlers) GetCurrentUser(rw http.ResponseWriter, r *http.Request) {
	// Extract user ID from context
	userID, ok := r.Context().Value(utils.CtxKey("user")).(string)
	if !ok || userID == "" {
		u.log.Printf("user ID not found in context or is empty\n")
		utils.ReturnJSON(rw, utils.ErrMessage{Error: "user ID not found in context"}, http.StatusUnauthorized)
		return
	}

	// Fetch the current user based on userID
	user, err := u.services.FetchUserByID(userID)
	if err != nil {
		u.log.Printf("error fetching user: %v\n", err)
		utils.ReturnJSON(rw, utils.ErrMessage{Error: "error getting user"}, http.StatusInternalServerError)
		return
	}

	if user == nil {
		u.log.Printf("user with ID %s not found\n", userID)
		utils.ReturnJSON(rw, utils.ErrMessage{Error: "user not found"}, http.StatusNotFound)
		return
	}

	// Prepare response data
	response := map[string]interface{}{
		"data": user,
	}

	// Return JSON response
	utils.ReturnJSON(rw, response, http.StatusOK)
}

func (u UserHandlers) UserTestRoute(rw http.ResponseWriter, r *http.Request) {
	// Your route logic here
	rw.Write([]byte("This is a protected route"))
}

// GetAllUsersById handles the HTTP request to fetch a user by ID
// func (u UserHandlers) GetAllUsersById(rw http.ResponseWriter, r *http.Request) {
// 	// Parse the ID from the URL parameters
// 	idStr := r.URL.Query().Get("id")
// 	if idStr == "" {
// 		utils.ReturnJSON(rw, utils.ErrMessage{Error: "Missing user ID"}, http.StatusBadRequest)
// 		return
// 	}

// 	// Convert ID to integer
// 	id, err := strconv.Atoi(idStr)
// 	if err != nil {
// 		utils.ReturnJSON(rw, utils.ErrMessage{Error: "Invalid user ID"}, http.StatusBadRequest)
// 		return
// 	}

// 	// Fetch the user by ID
// 	user, err := u.services.FetchUserByID(id)
// 	if err != nil {
// 		u.log.Printf("error fetching user by ID: %v\n", err)
// 		if err.Error() == "user with ID "+strconv.Itoa(id)+" not found" {
// 			utils.ReturnJSON(rw, utils.ErrMessage{Error: "User not found"}, http.StatusNotFound)
// 		} else {
// 			utils.ReturnJSON(rw, utils.ErrMessage{Error: "Failed to fetch user"}, http.StatusInternalServerError)
// 		}
// 		return
// 	}

// 	// Wrap the user in a data object
// 	response := map[string]interface{}{
// 		"data": user,
// 	}

// 	utils.ReturnJSON(rw, response, http.StatusOK)
// }
