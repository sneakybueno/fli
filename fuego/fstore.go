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

// Prompt retuns a string to be displayed
// as a prompt to the user
func (fs *FStore) Prompt() string {
	return "~/" + fs.BuildWorkingDirectoryPath(".") + " > "
}

// Wd (Working directory) returns the path for the "directory"
// the FStore is currently in.
func (fs *FStore) Wd() string {
	return fs.BuildWorkingDirectoryPath(".")
}

// FirebaseURLFromWorkingDirectory builds the firebase URL
// relative to the working directory. See BuildWorkingDirectoryPath
// for more info on how the path is built.
// Pass "" or "." to return the URL of the current working directory
func (fs *FStore) FirebaseURLFromWorkingDirectory(path string) string {
	return fs.FirebaseURL + fs.BuildWorkingDirectoryPath(path)
}

// BuildWorkingDirectoryPath builds the relative path
// based on the current working directory.
// Example: if the cwd = "/users" and path = "1234",
// it will return users/1234
// Pass "" or "." to return working directory path
func (fs *FStore) BuildWorkingDirectoryPath(p string) string {
	if p == "" || p == "." {
		return path.Join(fs.workingDirectory...)
	}

	wd := fs.workingDirectory[:]

	components := strings.Split(p, "/")
	for _, component := range components {
		if component == ".." {
			length := len(wd)
			if length > 0 {
				wd = wd[:length-1]
			}
		} else {
			wd = append(wd, component)
		}
	}

	return path.Join(wd...)
}

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

// Ls does a thing
func (fs *FStore) Ls(p string) (string, error) {
	path := fs.BuildWorkingDirectoryPath(p)
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
