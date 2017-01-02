package fuego

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
)

// FStore struct is used to interact with a firebase
// real time database
type FStore struct {
	client *http.Client

	FirebaseURL      string
	workingDirectory []string
}

// Cd (Change directory) emulates the cd command on a
// terminal. One major difference, the Cd command will never fail
// since firebase is a JSON store and not an actual directory
// structure.
func (fs *FStore) Cd(dir string) {
	components := strings.Split(dir, "/")
	for _, component := range components {
		if component == ".." {
			length := len(fs.workingDirectory)
			if length > 0 {
				fs.workingDirectory = fs.workingDirectory[:length-1]
			}
		} else {
			fs.workingDirectory = append(fs.workingDirectory, component)
		}
	}
}

// Wd (Working directory) returns the path for the "directory"
// the FStore is currently in.
func (fs *FStore) Wd() string {
	path := strings.Join(fs.workingDirectory, "/")
	return "/" + path
}

// Networking
// ----------------------------------------------------------------------------

func (fs *FStore) buildURL(path string) (string, error) {
	if path != "" {
		// santanize path here
		return fs.FirebaseURL + path + ".json", nil
	}

	return fs.FirebaseURL + ".json", nil
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

// FStore Get Operations
// ----------------------------------------------------------------------------

// Get performs a shallow get request for the given path
func (fs *FStore) Get(path string) (interface{}, error) {
	p, err := fs.buildURL(path)
	if err != nil {
		return nil, err
	}

	request, err := http.NewRequest("GET", p, nil)
	if err != nil {
		return nil, err
	}

	return fs.do(request)
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
