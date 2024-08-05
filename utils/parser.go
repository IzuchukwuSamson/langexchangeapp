package utils

import (
	"encoding/json"
	"io"
	"net/http"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// ReturnJSON writes the JSON encoding of data in the response writer
func ReturnJSON(rw http.ResponseWriter, data interface{}, code int) {
	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(code)
	err := json.NewEncoder(rw).Encode(data)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}
}

// FromJSON decodes a JSON request body and stores it in dest
func FromJSON(r io.Reader, dest interface{}) error {
	dec := json.NewDecoder(r)
	return dec.Decode(&dest)
}

// GetContextAsString gets the value of a context stored in the request
func GetContextAsString(r *http.Request, key string) string {
	return r.Context().Value(CtxKey(key)).(string)
}

// ConvertStringToObjectID converts a string to a primitive.ObjectID type
func ConvertStringToObjectID(rw http.ResponseWriter, r *http.Request, s string) (primitive.ObjectID, error) {
	objectId, err := primitive.ObjectIDFromHex(s)
	if err != nil {
		ReturnJSON(rw, ErrMessage{Error: "invalid id"}, http.StatusUnprocessableEntity)
		return primitive.ObjectID{}, err
	}
	return objectId, nil
}
