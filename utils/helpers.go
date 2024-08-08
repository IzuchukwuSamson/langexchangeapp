package utils

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"reflect"
)

// generic helper functions

func IsZeroValue(field reflect.Value) bool {
	zero := reflect.Zero(field.Type()).Interface()

	switch field.Kind() {
	case reflect.Slice, reflect.Array, reflect.Chan, reflect.Map:
		return field.Len() == 0
	default:
		return reflect.DeepEqual(zero, field.Interface())
	}

}

// GenerateRandomNumber generates a random string of 4 digits
func GenerateRandomNumber() (string, error) {
	n, err := rand.Int(rand.Reader, big.NewInt(10000))
	if err != nil {

		return "", err
	}
	randomNum := fmt.Sprintf("%04d", n)
	return randomNum, nil
}
