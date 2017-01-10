package fuego

import (
	"io/ioutil"

	"golang.org/x/net/context"
	"golang.org/x/oauth2/google"
)

const (
	firebaseDatabaseScope = "https://www.googleapis.com/auth/firebase.database"
	firebaseUserInfoScope = "https://www.googleapis.com/auth/userinfo.email"
)

// NewFStore builds a new store based on the 2 passed in params.
// Attempts to read the file from serviceAccountPath and parse it
// into a jwt.Config struct which provides an authorized http client.
// No validation is done to ensure a valid firebaseURL or a valid
// service account.
func NewFStore(firebaseURL string, serviceAccountPath string) (*FStore, error) {
	serviceAccountBytes, err := ioutil.ReadFile(serviceAccountPath)
	if err != nil {
		return nil, err
	}

	jwtConfig, err := google.JWTConfigFromJSON(serviceAccountBytes, firebaseDatabaseScope, firebaseUserInfoScope)
	if err != nil {
		return nil, err
	}

	ctx := context.TODO()
	client := jwtConfig.Client(ctx)

	fStore := &FStore{
		client:      client,
		FirebaseURL: firebaseURL,
	}

	return fStore, nil
}
