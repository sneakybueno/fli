package fuego

import (
	"fmt"
	"path"
	"strconv"
	"strings"
)

// FStore struct is used to interact with a firebase
// real time database
type FStore struct {
	fClient *FClient

	FirebaseURL      string
	workingDirectory []string
}

// NewFStore builds a new store based on the 2 passed in params.
// No validation is done to ensure a valid firebaseURL or a valid
// service account.
func NewFStore(firebaseURL string, serviceAccountPath string) (*FStore, error) {
	fClient, err := NewFClient(firebaseURL, serviceAccountPath)
	if err != nil {
		return nil, err
	}

	fStore := &FStore{
		fClient:     fClient,
		FirebaseURL: firebaseURL,
	}

	return fStore, nil
}

// FStore Directory Commands
// ----------------------------------------------------------------------------

// Cd (Change directory) emulates the cd command on a
// terminal. One major difference, the Cd command will never fail
// since firebase is a JSON store and not an actual directory
// structure.
func (fs *FStore) Cd(dir string) {
	// mimicing zsh behavior, if no arg is passed
	// return to root directory
	if dir == "" {
		fs.workingDirectory = []string{}
		return
	}

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
	return "~/" + fs.Wd() + " > "
}

// Wd (Working directory) returns the path for the "directory"
// the FStore is currently in.
func (fs *FStore) Wd() string {
	return path.Join(fs.workingDirectory...)
}

// WorkingDirectoryURL returns the firebase URL
// for the current working directory
func (fs *FStore) WorkingDirectoryURL() string {
	return fs.FirebaseURL + fs.Wd()
}

// Ls does a thing
func (fs *FStore) Ls() (string, error) {
	path := fs.Wd()
	data, err := fs.fClient.ShallowGet(path)
	if err != nil {
		return "", err
	}

	return firebaseDataToString(data)
}

// Private utilities
// ----------------------------------------------------------------------------

func firebaseDataToString(data interface{}) (string, error) {
	switch v := data.(type) {
	case int:
		return strconv.Itoa(v), nil
	case float64:
		return strconv.FormatFloat(v, 'f', -1, 64), nil
	case string:
		return v, nil
	case map[string]interface{}:
		keys := make([]string, 0, len(v))
		for k := range v {
			keys = append(keys, k)
		}
		return strings.Join(keys, "\t"), nil
	default:
		return "", fmt.Errorf("Error: Unsupported type %+v ", data)
	}
}
