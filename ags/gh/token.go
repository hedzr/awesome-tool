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
	"github.com/machinebox/graphql"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/ssh/terminal"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"syscall"
	"time"
)

const (
	ApiEntryV4 = "https://api.github.com/graphql"

	SepLine  = "_____"
	SepDigit = "___"

	ClientID               = "a1af758b9f273ef1c488"
	ClientSecret           = "38c223a7f008b21e018d231e0b0ce64711aa493b"
	asgoStatsTokenFilename = "$HOME/.asg.stats.token"
)

func GhApplyToken(req *graphql.Request) *graphql.Request {
	req.Header["Authorization"] = []string{"token " + GhRequestToken()}
	return req
}

func firstLine(s string) string {
	if i := strings.Index(s, "\n"); i >= 0 {
		return s[0:i]
	}
	return s
}

func GhRequestToken() (token string) {
	token = os.Getenv("ASGO_STATS_TOKEN")
	if len(token) == 0 {
		if b, err := ioutil.ReadFile(os.ExpandEnv(asgoStatsTokenFilename)); err == nil {
			token = firstLine(string(b))
			if len(token) == 40 {
				// logrus.Debugf("  ..using existed token : %v", token)
				return
			}
		}

		logrus.Debugf("  ..request and create authorization")

		// req and get an new authorization
		url := fmt.Sprintf("https://api.github.com/authorizations/clients/%v", ClientID)
		body := fmt.Sprintf(`{
  "client_secret": "%v",
  "scopes": [
    "public_repo"
  ],
  "note": "asgo-stats app",
  "fingerprint": "asgo.stats"
}`, ClientSecret)

		var ok bool
		var gr = make(map[string]interface{})
		gr = httpReadJson("PUT", url, body)
		logrus.Debugf(`token: %v (hashed: %v), updated at: %v`, gr["token"], gr["hashed_token"], gr["updated_at"])
		if token, ok = gr["token"].(string); ok {
			if len(token) == 0 {
				if token, ok = gr["hashed_token"].(string); ok {
					_ = ioutil.WriteFile(os.ExpandEnv(asgoStatsTokenFilename), []byte(token), 0600)
					return
				}

				url = gr["url"].(string)
				gr = httpReadJson("GET", url, "")
			}
			_ = ioutil.WriteFile(os.ExpandEnv(asgoStatsTokenFilename), []byte(token), 0600)
		}

	}
	return
}

func GhHttpClient() *http.Client {
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
		logrus.Fatal(err)
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
		// logrus.Debugf("  > %v:%v", username, password)
		req.SetBasicAuth(username, password)

		trans := &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: false}}
		client := &http.Client{Timeout: 300 * time.Second, Transport: trans}

		logrus.Debug("Headers requesting:")
		for k, v := range req.Header {
			logrus.Debugf("  - %v: %v", k, v)
		}

		if resp, err := client.Do(req); err != nil {
			logrus.Fatal(err)
		} else {
			defer resp.Body.Close()
			logrus.Debug("Headers returns:")
			for k, v := range resp.Header {
				logrus.Debugf("  - %v: %v", k, v)
			}

			var buf bytes.Buffer
			if _, err := io.Copy(&buf, resp.Body); err != nil {
				logrus.Fatal(errors.Wrap(err, "reading body"))
			}
			logrus.Debugf("<< %s", buf.String())
			if err := json.NewDecoder(&buf).Decode(&r); err != nil {
				logrus.Fatal(errors.Wrap(err, "decoding response"))
			}
		}
	}
	return
}
