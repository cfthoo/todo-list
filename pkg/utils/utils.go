package utils

import (
	"errors"
	"net/http"
	"strings"

	"github.com/cfthoo/todo-app/pkg/config"
	"golang.org/x/oauth2"
)

// StringExtractor extracts input string based on searchString
func StringExtractor(input string, searchString string) string {

	outputString := ""

	index := strings.Index(input, searchString)

	if index != -1 {

		outputString = input[index+len(searchString):]
		index = strings.Index(outputString, searchString)

		if index != -1 {
			outputString = outputString[:index]
		}
	}
	return outputString
}

// SetOuathConfigByLoginType sets oauth configuration by loginType
func SetOuathConfigByLoginType(loginType string) (*oauth2.Config, string, error) {
	var configs *oauth2.Config
	oauthUrl := ""
	switch loginType {
	case "google":
		configs = config.SetupGoogleConfig()
		oauthUrl = config.GoogleOauthUrl
	case "fb":
		configs = config.SetupFbConfig()
		oauthUrl = config.FacebookOauthUrl
	case "github":
		configs = config.SetupGithubConfig()
		oauthUrl = config.GithubOauthUrl
	default:
		return nil, "", errors.New("Unknown login type")
	}

	return configs, oauthUrl, nil
}

// FetchUserData fetches user data based on loginType from different provider ie google/fb/github
func FetchUserData(loginType string, oauthUrl string, token string) (*http.Response, error) {

	if loginType != "github" {
		resp, err := http.Get(oauthUrl + token)
		if err != nil {
			return nil, errors.New("Failed to get User data")
		}
		return resp, nil

	} else {
		req, err := http.NewRequest("GET", oauthUrl, nil)
		if err != nil {

			return nil, errors.New("Failed to get User data")
		}

		req.Header.Set("Authorization", "token "+string(token))

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			return nil, errors.New("Failed to get User data")
		}

		return resp, nil
	}
}
