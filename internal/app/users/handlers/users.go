package handlers

import (
	"fmt"
	"net/http"
)

func (u UserHandlers) Dashboard(rw http.ResponseWriter, r *http.Request) {
	fmt.Println("Dashboard route")
}
