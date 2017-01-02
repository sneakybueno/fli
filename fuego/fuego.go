package fuego

import (
	"encoding/json"
	"log"
	"net/http"
	"time"
)

const (
	fStoreDefaultTimeout = 60
)

// FStore struct is used to interact with a firebase
// real time database
type FStore struct {
	client *http.Client

	FirebaseURL string
}

func NewFStore(firebaseURL string) *FStore {
	client := &http.Client{
		Timeout: time.Second * fStoreDefaultTimeout,
	}

	return &FStore{
		client:      client,
		FirebaseURL: firebaseURL,
	}
}

func (fs *FStore) buildURL(path string) (string, error) {
	if path != "" {
		// santanize path here
		return fs.FirebaseURL + path + ".json", nil
	}

	return fs.FirebaseURL + ".json", nil
}

// SGet is a helper to perform a shallow get at the current
// store's path
func (fs *FStore) SGet() {

}

// ShallowGet performs a shallow get request for the given path
func (fs *FStore) ShallowGet(path string) (interface{}, error) {
	p, err := fs.buildURL(path)
	if err != nil {
		return nil, err
	}

	request, err := http.NewRequest("GET", p, nil)
	if err != nil {
		return nil, err
	}

	q := request.URL.Query()
	q.Add("shallow", "true")
	request.URL.RawQuery = q.Encode()

	log.Println(request.URL.String())
	return fs.do(request)
}

func (fs *FStore) do(request *http.Request) (interface{}, error) {
	resp, err := fs.client.Do(request)
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
