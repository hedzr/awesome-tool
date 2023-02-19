/*
 * Copyright © 2019 Hedzr Yeh.
 */

package gh

import (
	"bufio"
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"syscall"
	"time"

	"github.com/hedzr/cmdr"
	"github.com/machinebox/graphql"
	"golang.org/x/crypto/ssh/terminal"
	"gopkg.in/hedzr/errors.v3"
)

const (
	ApiEntryV4 = "https://api.github.com/graphql"

	SepLine  = "_____"
	SepDigit = "___"

	ClientID               = "a1af758b9f273ef1c488"
	ClientSecret           = "38c223a7f008b21e018d231e0b0ce64711aa493b"
	asgoStatsTokenFilename = "$HOME/.asg.stats.token"
)

func ApplyToken(req *graphql.Request) *graphql.Request {
	req.Header["Authorization"] = []string{"token " + RequestToken()}
	return req
}

func firstLine(s string) string {
	if i := strings.Index(s, "\n"); i >= 0 {
		return s[0:i]
	}
	return s
}

func RequestToken() (token string) {
	token = os.Getenv("ASGO_STATS_TOKEN")
	if len(token) == 0 {
		if b, err := ioutil.ReadFile(os.ExpandEnv(asgoStatsTokenFilename)); err == nil {
			token = firstLine(string(b))
			if len(token) == 40 {
				// cmdr.Logger.Debugf("  ..using existed token : %v", token)
				return
			}
		}

		cmdr.Logger.Debugf("  ..request and create authorization")

		fingerprint, fTimes := "awesome-tool.stats", 1

	RetryGetToken:
		// req and get an new authorization
		url := fmt.Sprintf("https://api.github.com/authorizations/clients/%v", ClientID)
		body := fmt.Sprintf(`{
  "client_secret": "%v",
  "scopes": [
    "public_repo"
  ],
  "note": "asgo-stats app",
  "fingerprint": "%v"
}`, ClientSecret, fingerprint)

		var ok bool
		var gr map[string]interface{}
		gr = httpReadJson("PUT", url, body)
		cmdr.Logger.Debugf(`token: %v (hashed: %v), updated at: %v`, gr["token"], gr["hashed_token"], gr["updated_at"])
		if token, ok = gr["token"].(string); ok {
			if len(token) == 0 {
				if token, ok = gr["hashed_token"].(string); ok {
					// _ = ioutil.WriteFile(os.ExpandEnv(asgoStatsTokenFilename), []byte(token), 0600)
					cmdr.Logger.Warnf("The token for fingerprint '%v' cannot be re-fetched, you MUST have to request a new one.", fingerprint)
					fingerprint = fmt.Sprintf("awesome-tool.stats.%d", fTimes)
					fTimes++
					goto RetryGetToken
				}

				url = gr["url"].(string)
				gr = httpReadJson("GET", url, "")
			}
			_ = ioutil.WriteFile(os.ExpandEnv(asgoStatsTokenFilename), []byte(token), 0600)
			cmdr.Logger.Infof(`token: %v (hashed: %v), updated at: %v. fingerprint is %v`, gr["token"], gr["hashed_token"], gr["updated_at"], fingerprint)
		}

	}
	return
}

func GithubHttpClient() *http.Client {
	h := &http.Client{}
	return h
}

func WithOAuth2Token() graphql.ClientOption {
	return func(c *graphql.Client) {
	}
}

func readUsername(tip string) string {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print(tip)
	text, _ := reader.ReadString('\n')
	return firstLine(text)
}

func readPassword(tip string) string {
	var bytePassword []byte
	var err error
	fmt.Print(tip)
	if bytePassword, err = terminal.ReadPassword(int(syscall.Stdin)); err != nil {
		fmt.Println() // it's necessary to add a new line after user's input
		return ""
	}
	fmt.Println() // it's necessary to add a new line after user's input
	return firstLine(string(bytePassword))
}

var (
	username, password string
)

func httpReadJson(method, url, body string) (r map[string]interface{}) {
	r = make(map[string]interface{})
	var requestBody bytes.Buffer
	requestBody.WriteString(body)
	if req, err := http.NewRequest(method, url, &requestBody); err != nil {
		cmdr.Logger.Fatalf("error: %v", err)
	} else {
		// req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Accept", "application/json; charset=utf-8")
		// Accept-Encoding: gzip

		if len(username) == 0 {
			username = readUsername("Enter your GitHub username: ")
		}
		if len(password) == 0 {
			password = readPassword("Enter your GitHub password: ")
		}
		// cmdr.Logger.Debugf("  > %v:%v", username, password)
		req.SetBasicAuth(username, password)

		trans := &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: false}}
		client := &http.Client{Timeout: 300 * time.Second, Transport: trans}

		cmdr.Logger.Debugf("Headers requesting:")
		for k, v := range req.Header {
			cmdr.Logger.Debugf("  - %v: %v", k, v)
		}

		if resp, err := client.Do(req); err != nil {
			cmdr.Logger.Fatalf("error: %v", err)
		} else {
			defer resp.Body.Close()
			cmdr.Logger.Debugf("Headers returns:")
			for k, v := range resp.Header {
				cmdr.Logger.Debugf("  - %v: %v", k, v)
			}

			var buf bytes.Buffer
			if _, err := io.Copy(&buf, resp.Body); err != nil {
				cmdr.Logger.Fatalf("error: %v", errors.New("reading body").WithErrors(err))
			}
			cmdr.Logger.Debugf("<< %s", buf.String())
			if err := json.NewDecoder(&buf).Decode(&r); err != nil {
				cmdr.Logger.Fatalf("error: %v", errors.New("decoding response").WithErrors(err))
			}
		}
	}
	return
}
