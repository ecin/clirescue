package trackerapi

import (
	"os"
	"io/ioutil"

	"github.com/ecin/clirescue/cmdutil"
	"github.com/ecin/clirescue/user"
)

const (
	TOKEN_FILENAME = "user_token"
	BASE_URL = "https://www.pivotaltracker.com/services/v5/"
)

var (
	cache_dir = os.Getenv("HOME") + "/.clirescue/"
	currentUser  *user.User = user.NewUser(readToken())
)

func saveToken (token string) {
	os.Mkdir(cache_dir, os.ModeDir | os.ModePerm)
	ioutil.WriteFile(cache_dir + TOKEN_FILENAME, []byte(currentUser.APIToken), os.ModePerm)
}

func readToken () string {
	tokenPath := cache_dir + TOKEN_FILENAME
	if (!cmdutil.FileExists(tokenPath)) {
		return ""
	}

	bytes, error := ioutil.ReadFile(tokenPath)
	if error != nil {
		return ""
	}
	return string(bytes)
}