package account

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"os/user"
	"strings"
)

type Credentials struct {
	Username string
	Token    string
}

type Account struct {
	Credits        int    `json:"credits,omitempty"`
	JoinedAt       string `json:"joinedAt,omitempty"`
	ShipCount      int    `json:"shipCount,omitempty"`
	StructureCount int    `json:"structureCount,omitempty"`
	Username       string `json:"username,omitempty"`
}

func (a Account) String() string {
	return fmt.Sprintf("Username:%s\nCredits:%d\nShips:%d\nStructures:%d", a.Username, a.Credits, a.ShipCount, a.StructureCount)
}

type User struct {
	User Account `json:"user"`
}

func SetUsernameAndToken(username string) (string, error) {
	// would be niice to wrap all commands in tthe CLI
	// to not run unless this has already been called
	type User struct {
		Credits  int    `json:"credits,omitempty"`
		Loans    []int  `json:"loans,omitempty"`
		Ships    []int  `json:"ships,omitempty"`
		Username string `json:"username"`
	}

	type Claim struct {
		Token string `json:"token"`
		User  User   `json:"user"`
	}

	url := fmt.Sprintf("https://api.spacetraders.io/users/%s/claim", username)

	b := new(bytes.Buffer)
	json.NewEncoder(b)
	resp, _ := http.Post(url, "application/json:charset=utf-8", b)
	defer resp.Body.Close()

	// create a new struct
	var responseClaim Claim

	// given the response, format it into the new struct
	err := json.NewDecoder(resp.Body).Decode(&responseClaim)
	if err != nil {
		return "", err
	}

	// --- Given the token save to ~/.spacetravlers/credentials
	// These should probably be global variables
	usr, _ := user.Current()
	config_dir := fmt.Sprintf("%s/.spacetravlers", usr.HomeDir)
	if _, err := os.Stat(config_dir); os.IsNotExist(err) {
		err := os.Mkdir(config_dir, 0777)
		if err != nil {
			return "", err
		}
	}

	f, err := os.Create(fmt.Sprintf("%s/credentials", config_dir))
	if err != nil {
		return "", err
	}
	defer f.Close()

	creds := fmt.Sprintf("username=%s\ntoken=%s", username, responseClaim.Token)
	w := bufio.NewWriter(f)
	w.WriteString(creds)
	w.Flush()

	return resp.Status, nil
}

func GetUsernameAndToken() (Credentials, error) {
	// probably should be a private function
	usr, _ := user.Current()
	config_dir := fmt.Sprintf("%s/.spacetravlers", usr.HomeDir)
	credentials_file := fmt.Sprintf("%s/credentials", config_dir)
	file, err := os.Open(credentials_file)
	if err != nil {
		return Credentials{}, err
	}
	defer file.Close()

	credentials_data := make(map[string]string)
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.Split(scanner.Text(), "=")
		credentials_data[line[0]] = line[1]
	}

	return Credentials{Username: credentials_data["username"], Token: credentials_data["token"]}, nil
}

func GetAccount(token string, username string) (Account, error) {
	var creds Credentials

	if token == "" || username == "" {
		// try to load credentials based on ~/.spacetravlers/credentials file
		creds, _ = GetUsernameAndToken()
	} else {
		// should we validate these are correct?
		creds = Credentials{
			Username: username,
			Token:    token,
		}
	}

	url := fmt.Sprintf("https://api.spacetraders.io/my/account?token=%s", url.QueryEscape(creds.Token))
	resp, err := http.Get(url)
	if err != nil {
		return Account{}, err
	}

	// given the response, format it into the new struct
	var user User
	err = json.NewDecoder(resp.Body).Decode(&user)
	if err != nil {
		return Account{}, err
	}

	return user.User, nil
}
