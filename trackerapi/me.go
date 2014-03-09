package trackerapi

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/ecin/clirescue/cmdutil"
)

const (
	ME_URL = "me/"
	ME_FILENAME = "me"
)

var (
	Stdout       *os.File   = os.Stdout
)

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

func Me() {
	if (currentUser.APIToken == "") {
		getCredentials()
	}

	body, err := requestBody()
	if err != nil {
		fmt.Println(err)
		return
	}
	parse(body)
}

func requestBody() ([]byte, error) {

	req, err := http.NewRequest("GET", BASE_URL + ME_URL, nil)
	if (err != nil) {
		return nil, err
	}

	if currentUser.APIToken != "" {
		req.Header.Add("X-TrackerToken", currentUser.APIToken)
	} else {
		req.SetBasicAuth(currentUser.Username, currentUser.Password)
	}
	return getResponseBody(req)
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

func getCredentials() (string, string) {
	fmt.Fprint(Stdout, "Username: ")
	var username = cmdutil.ReadLine()
	cmdutil.Silence()
	fmt.Fprint(Stdout, "Password: ")

	var password = cmdutil.ReadLine()
	currentUser.Login(username, password)
	cmdutil.Unsilence()
}