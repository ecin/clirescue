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

func Me() {
	
	req, err := buildRequest()
	body, err := getResponseBody(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	parse(body)
}

func buildRequest() (*http.Request, error) {
	req, err := http.NewRequest("GET", BASE_URL + ME_URL, nil)
	if (err != nil) {
		return req, err
	}

	if user_token != "" {
		req.Header.Add("X-TrackerToken", user_token)
	} else {
		username, password := getCredentials()
		req.SetBasicAuth(username, password)
	}

	return req, nil
}

func parse(body []byte) {
	type MeResponse struct {
		APIToken string `json:"api_token"`
	}

	var meResp = new(MeResponse)
	err := json.Unmarshal(body, &meResp)
	if err != nil {
		fmt.Println("error:", err)
	}

	user_token = meResp.APIToken
	saveToken()
}

func getCredentials() (string, string) {
	fmt.Fprint(Stdout, "Username: ")
	var username = cmdutil.ReadLine()
	cmdutil.Silence()
	fmt.Fprint(Stdout, "Password: ")

	var password = cmdutil.ReadLine()
	cmdutil.Unsilence()
	return username, password
}