package fuego

import (
	"net/http"
	"time"
)

const (
	fStoreDefaultTimeout = 60
)

func NewFStore(firebaseURL string) *FStore {
	client := &http.Client{
		Timeout: time.Second * fStoreDefaultTimeout,
	}

	return &FStore{
		client:      client,
		FirebaseURL: firebaseURL,
	}
}
