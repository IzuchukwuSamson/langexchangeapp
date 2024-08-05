package users

import (
	"fmt"
	"io"
	"lexibuddy/config/db"
	"lexibuddy/services"
	"log"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

// func Test_Abs(t *testing.T) {
// 	t.Log("similar to fmt.Println() and concurrently safe")
// 	t.Fail()    // will show a test case that has failed in the results
// 	t.FailNow() // t.Fail + safely exit without continuing
// 	t.Error()   // t.Log() + t.Fail()
// 	t.Fatal()   // t.Log() + t.FailNow()
// }

func Test_Signup(t *testing.T) {
	// fake server
	mongodb, err := db.Mongo()
	if err != nil {
		log.Fatal(err)
	}

	logger := log.New(os.Stdout, fmt.Sprintf("%s:", os.Getenv("APP_NAME")), log.LstdFlags)

	// initialze db
	dbConn := db.NewDB(mongodb, nil)
	// Initialize the service
	userService := services.NewUserService(dbConn.Mongo, logger)

	// Initialize the handlers
	uh := NewUserHandlers(log.Default(), userService)

	// request / call the handler
	rw := httptest.NewRecorder()

	data := `
		{
			"firstname": "Izu",
			"lastname": "samson",
			"email": "sam12@test.com",
			"password": "pass1234",
		}

	`

	req := httptest.NewRequest("POST", "/signup", strings.NewReader(data))

	uh.Signup(rw, req)
	res := rw.Result()

	c, err := io.ReadAll(res.Body)
	if err != nil {
		t.Error(err)
	}

	fmt.Println(c)

	// send a request to the fake server

	// evaluate the response

}
