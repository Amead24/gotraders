package utils

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/user"
	"strings"
)

type JsonError struct {
	ErrorMsg struct {
		Msg string `json:"message"`
	} `json:"error"`
}

func (je JsonError) Error() string {
	return fmt.Sprintln(je.ErrorMsg.Msg)
}

type Account struct {
	Credits        int    `json:"credits,omitempty"`
	JoinedAt       string `json:"joinedAt,omitempty"`
	ShipCount      int    `json:"shipCount,omitempty"`
	StructureCount int    `json:"structureCount,omitempty"`
	Username       string `json:"username,omitempty"`
}

// lesson 4 - boiler plate & memory management
// it's not clear if I should be rturning anything when modifying an object?
func Get(url string, params map[string]string, obj interface{}) error {
	token, err := getToken()
	if err != nil {
		return err
	}

	var queryParams []string
	for key, value := range params {
		queryParams = append(queryParams, fmt.Sprintf("?%s=%s", key, value))
	}

	queryParamString := strings.Join(queryParams, "")
	urlWithParams := fmt.Sprintf("%s?token=%s&%s", url, token, queryParamString)

	resp, _ := http.Get(urlWithParams)
	if resp.StatusCode >= 400 {
		bodyBytes, _ := ioutil.ReadAll(resp.Body)
		if os.Getenv("GOTRADERS_DEBUG") != "" {
			fmt.Printf("Non-200 response on %s\n", url)
			fmt.Println(string(bodyBytes))
		}

		var msg JsonError
		if err := json.Unmarshal(bodyBytes, &msg); err != nil {
			if terr, ok := err.(*json.UnmarshalTypeError); ok {
				log.Fatalf("Bad response; failed to unmarshal field %s \n", terr.Field)
			}
			return err
		}

		return msg
	}

	bodyBytes, _ := ioutil.ReadAll(resp.Body)
	if err := json.Unmarshal(bodyBytes, &obj); err != nil {
		if terr, ok := err.(*json.UnmarshalTypeError); ok {
			log.Fatalf("Bad response; failed to unmarshal field %s \n", terr.Field)
		}
		return err
	}

	return nil
}

func Post(url string, params map[string]string, obj interface{}) (bool, error) {
	token, err := getToken()
	if err != nil {
		return false, err
	}

	urlWithCreds := fmt.Sprintf("%s?token=%s", url, token)

	postBody, _ := json.Marshal(params)
	responseBody := bytes.NewBuffer(postBody)
	resp, _ := http.Post(urlWithCreds, "application/json", responseBody)
	if resp.StatusCode >= 400 {
		bodyBytes, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return false, err
		}

		if os.Getenv("GOTRADERS_DEBUG") != "" {
			fmt.Printf("Non-200 response on %s\n", url)
			fmt.Println(string(bodyBytes))
		}

		var msg JsonError
		if err := json.Unmarshal(bodyBytes, &msg); err != nil {
			if terr, ok := err.(*json.UnmarshalTypeError); ok {
				log.Fatalf("Bad response; failed to unmarshal field %s \n", terr.Field)
			}
			return false, err
		}

		return false, msg
	}

	// How am I supposed to know what errors to handle here? Or anywhere?!
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return false, err
	}

	if os.Getenv("GOTRADERS_DEBUG") != "" {
		fmt.Printf("~200 response on %s\n", url)
		fmt.Println(string(bodyBytes))
	}

	if err := json.Unmarshal(bodyBytes, &obj); err != nil {
		if terr, ok := err.(*json.UnmarshalTypeError); ok {
			log.Fatalf("Valid response; failed to unmarshal field %s \n", terr.Field)
		}
		return false, err
	}

	return true, nil
}

func GetServerHealth() (string, error) {
	resp, err := http.Get("http://api.spacetraders.io/game/status")

	if err != nil {
		return "", err
	}

	return resp.Status, nil
}

func SetUsernameAndToken(username string) (bool, error) {
	// would be niice to wrap all commands in tthe CLI
	// to not run unless this has already been called
	type User struct {
		Credits  int    `json:"credits,omitempty"`
		Loans    []int  `json:"loans,omitempty"`
		Ships    []int  `json:"ships,omitempty"`
		Username string `json:"username,omitempty"`
	}

	type Claim struct {
		Token string `json:"token,omitempty"`
		User  User   `json:"user,omitempty"`
	}

	// --- Given the token save to ~/.spacetravlers/credentials
	// These should probably be global variables
	usr, _ := user.Current()
	config_dir := fmt.Sprintf("%s/.spacetravlers", usr.HomeDir)
	if _, err := os.Stat(config_dir); os.IsNotExist(err) {
		err := os.Mkdir(config_dir, 0777)
		if err != nil {
			return false, err
		}
	}

	credrc := fmt.Sprintf("%s/credentials", config_dir)
	fmt.Printf("Creating credential file: %s\n", credrc)
	f, err := os.Create(credrc)
	if err != nil {
		return false, err
	}
	defer f.Close()

	var responseClaim Claim
	url := fmt.Sprintf("https://api.spacetraders.io/users/%s/claim", username)
	params := map[string]string{}

	ok, err := Post(url, params, &responseClaim)
	if !ok {
		return false, err
	}

	creds := fmt.Sprintf("username=%s\ntoken=%s\n", username, responseClaim.Token)
	w := bufio.NewWriter(f)
	w.WriteString(creds)
	w.Flush()

	return true, nil
}

func GetAccount() (Account, error) {
	type User struct {
		User Account `json:"user"`
	}

	var user User
	url := "https://api.spacetraders.io/my/account"
	params := map[string]string{}
	err := Get(url, params, &user)

	if err != nil {
		return Account{}, err
	}

	return user.User, nil
}

func getToken() (string, error) {
	usr, _ := user.Current()
	config_dir := fmt.Sprintf("%s/.spacetravlers", usr.HomeDir)
	credentials_file := fmt.Sprintf("%s/credentials", config_dir)
	file, err := os.Open(credentials_file)
	if err != nil {
		return "", err
	}
	defer file.Close()

	credentials_data := make(map[string]string)
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.Split(scanner.Text(), "=")
		credentials_data[line[0]] = line[1]
	}

	return credentials_data["token"], nil
}
