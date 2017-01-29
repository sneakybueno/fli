package fuego

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"

	"golang.org/x/oauth2/google"
)

const (
	firebaseDatabaseScope = "https://www.googleapis.com/auth/firebase.database"
	firebaseUserInfoScope = "https://www.googleapis.com/auth/userinfo.email"
)

// FClient is wrapper for http interactions with firebase's
// real time database. Init with a custom http client
// with proper auth tokens or config with a service account.
type FClient struct {
	client *http.Client

	FirebaseURL string
}

// NewFClient builds a firebase client based on the 2 passed in params.
// Attempts to read the file from serviceAccountPath and parse it
// into a jwt.Config struct which provides an authorized http client.
// No validation is done to ensure a valid firebaseURL or a valid
// service account.
func NewFClient(firebaseURL string, serviceAccountPath string) (*FClient, error) {
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

	fClient := &FClient{
		client:      client,
		FirebaseURL: firebaseURL,
	}

	return fClient, nil
}

// FStore Get Operations
// ----------------------------------------------------------------------------

// Get performs a http get request for the given path
func (fc *FClient) Get(path string, params map[string]string) (interface{}, error) {
	p, err := fc.buildURL(path)
	if err != nil {
		return nil, err
	}

	request, err := http.NewRequest("GET", p, nil)
	if err != nil {
		return nil, err
	}

	if len(params) > 0 {
		q := request.URL.Query()

		for key, value := range params {
			q.Add(key, value)
		}

		request.URL.RawQuery = q.Encode()
	}

	return fc.do(request)
}

// ShallowGet performs a http shallow get request for the given path
func (fc *FClient) ShallowGet(path string) (interface{}, error) {
	params := map[string]string{"shallow": "true"}
	return fc.Get(path, params)
}

// Networking
// ----------------------------------------------------------------------------

func (fc *FClient) buildURL(p string) (string, error) {
	if p != "" {
		// need to escape p properly here
		u := fc.FirebaseURL + p
		return u + ".json", nil
	}

	return fc.FirebaseURL + ".json", nil
}

func (fc *FClient) do(request *http.Request) (interface{}, error) {
	resp, err := fc.client.Do(request)
	if err != nil {
		return nil, err
	}

	decoder := json.NewDecoder(resp.Body)

	var b interface{}
	err = decoder.Decode(&b)
	if err != nil {
		return nil, err
	}

	return b, nil
}
