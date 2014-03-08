package trackerapi

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"io/ioutil"

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
	setCredentials()
	parse(makeRequest())
}

func makeRequest() []byte {

	client := &http.Client{}
	req, err := http.NewRequest("GET", BASE_URL + ME_URL, nil)
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
	if (currentUser.APIToken != "") {
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