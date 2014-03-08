package trackerapi

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	u "os/user"

	"github.com/ecin/clirescue/cmdutil"
	"github.com/ecin/clirescue/user"
)

const (
	TOKEN_FILENAME = "user_token"
	TOKEN_DIR = "/.clirescue/"
)

var (
	URL          string     = "https://www.pivotaltracker.com/services/v5/me"
	FileLocation string     = homeDir() + "/.tracker"
	currentUser  *user.User = user.New()
	Stdout       *os.File   = os.Stdout
)

func Me() {
	setCredentials()
	parse(makeRequest())
	ioutil.WriteFile(FileLocation, []byte(currentUser.APIToken), 0644)
}

func makeRequest() []byte {

	client := &http.Client{}
	req, err := http.NewRequest("GET", URL, nil)
	if currentUser.APIToken != "" {
		req.Header.Add("X-TrackerToken", currentUser.APIToken)
	} else {
		req.SetBasicAuth(currentUser.Username, currentUser.Password)
	}
	resp, err := client.Do(req)
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Print(err)
	}
	fmt.Printf("\n****\nAPI response: \n%s\n", string(body))
	return body
}

func parse(body []byte) {
	var meResp = new(MeResponse)
	err := json.Unmarshal(body, &meResp)
	if err != nil {
		fmt.Println("error:", err)
	}

	currentUser.APIToken = meResp.APIToken

	saveToken(currentUser.APIToken)
}

func setCredentials() {
	userToken := readToken()
	if (userToken != "") {
		currentUser.APIToken = userToken
		return;
	}

	fmt.Fprint(Stdout, "Username: ")
	var username = cmdutil.ReadLine()
	cmdutil.Silence()
	fmt.Fprint(Stdout, "Password: ")

	var password = cmdutil.ReadLine()
	currentUser.Login(username, password)
	cmdutil.Unsilence()
}

func homeDir() string {
	usr, _ := u.Current()
	return usr.HomeDir
}

type MeResponse struct {
	APIToken string `json:"api_token"`
	Username string `json:"username"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Initials string `json:"initials"`
	Timezone struct {
		Kind      string `json:"kind"`
		Offset    string `json:"offset"`
		OlsonName string `json:"olson_name"`
	} `json:"time_zone"`
}

func saveToken (token string) {
	tokenDir := os.Getenv("HOME") + TOKEN_DIR
	os.Mkdir(tokenDir, os.ModeDir | os.ModePerm)
	ioutil.WriteFile(tokenDir + TOKEN_FILENAME, []byte(currentUser.APIToken), os.ModePerm)
}

func readToken () string {
	tokenPath := os.Getenv("HOME") + TOKEN_DIR + TOKEN_FILENAME
	if (!fileExists(tokenPath)) {
		return ""
	}

	bytes, error := ioutil.ReadFile(tokenPath)
	if error != nil {
		return ""
	}
	return string(bytes)

}

func fileExists(path string) bool {
  if _, err := os.Stat(path); os.IsNotExist(err) {
    return false
  }
  return true
}