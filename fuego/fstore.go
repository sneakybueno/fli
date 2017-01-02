package fuego

import (
	"encoding/json"
	"net/http"
	"path"
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

// Prompt retuns a string to be displayed
// as a prompt to the user
func (fs *FStore) Prompt() string {
	return "~/" + fs.Wd()
}

// Wd (Working directory) returns the path for the "directory"
// the FStore is currently in.
func (fs *FStore) Wd() string {
	return path.Join(fs.workingDirectory...)
}

// WorkingDirectoryURL returns the firebase URL
// for the current working directory
func (fs *FStore) WorkingDirectoryURL() string {
	return path.Join(fs.FirebaseURL, fs.Wd())
}

func (fs *FStore) Ls() (string, error) {
	path := fs.Wd()
	_, err := fs.ShallowGet(path)
	if err != nil {
		return "", err
	}
	// format v

	return "", nil
}

// Networking
// ----------------------------------------------------------------------------

func (fs *FStore) buildURL(p string) (string, error) {
	if p != "" {
		u := path.Join(fs.FirebaseURL, p)
		return u + ".json", nil
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

	return fs.do(request)
}
