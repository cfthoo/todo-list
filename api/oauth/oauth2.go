package oauth2api

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/cfthoo/todo-app/pkg/controller"
	"github.com/cfthoo/todo-app/pkg/utils"
)

type UserInfo struct {
	Id    string `json:"id"`
	Email string `json:"email"`
}

type GithubUserInfo struct {
	Id    int    `json:"id"`
	Email string `json:"email"`
}

var UserId string

// LoginHandler handles and redirects to login page based on loginType ie google/fb/github
func LoginHandler(w http.ResponseWriter, r *http.Request) {

	// get loginType from url
	loginType := utils.StringExtractor(r.URL.String(), "/")

	// set configuration based on loginType
	configs, _, err := utils.SetOuathConfigByLoginType(loginType)
	if err != nil {
		fmt.Fprintln(w, err)
	}
	url := configs.AuthCodeURL("randomstate")

	http.Redirect(w, r, url, http.StatusSeeOther)

}

// CallbackHandler handles the callback action for loginType ie google/fb/github
func CallbackHandler(w http.ResponseWriter, r *http.Request) {

	// get loginType from url
	loginType := utils.StringExtractor(r.URL.String(), "/")

	//state
	state := r.URL.Query()["state"][0]
	if state != "randomstate" {
		fmt.Fprintln(w, "states dont match")
		return
	}

	//code
	code := r.URL.Query()["code"][0]

	// set configuration based on loginType
	configs, oauthUrl, err := utils.SetOuathConfigByLoginType(loginType)
	if err != nil {
		fmt.Fprintln(w, err)
	}

	// exchange code for token
	token, err := configs.Exchange(context.Background(), code)
	if err != nil {
		fmt.Fprintln(w, "Code-Token Exchange Failed")
	}

	// fetch user info with oauth api
	resp, err := utils.FetchUserData(loginType, oauthUrl, token.AccessToken)

	// Parse user data JSON Object
	defer resp.Body.Close()
	contents, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Fprintf(w, "failed read response: %s", err.Error())
		return
	}

	// get userId and to be used in other api ,ie listTodoList by userId
	// in real world scenario please use cookie to store the userId
	if loginType == "github" {
		var userInfo GithubUserInfo
		err = json.Unmarshal(contents, &userInfo)
		if err != nil {
			fmt.Fprintf(w, "failed to unmarshal JSON: %s", err.Error())
			return
		}
		UserId = strconv.Itoa(userInfo.Id)
	} else {
		var userInfo UserInfo
		err = json.Unmarshal(contents, &userInfo)
		if err != nil {
			fmt.Fprintf(w, "failed to unmarshal JSON: %s", err.Error())
			return
		}

		// it is recommeded to set user_id in cookies
		UserId = userInfo.Id

		// cookie := http.Cookie{
		// 	Name:     "UserId",
		// 	Value:    userInfo.Id,
		// 	Path:     "/",
		// 	MaxAge:   3600,
		// 	HttpOnly: true,
		// 	Secure:   true,
		// 	SameSite: http.SameSiteLaxMode,
		// }

		// http.SetCookie(w, &cookie)

	}

	// get jwtToken
	// in real world scenario please save the token somewhere ie DB ,cookie/session
	jwtToken, err := controller.CreateJWT([]byte(token.AccessToken))

	responseBody := "Authorized!\nPlease copy the token :" + jwtToken
	// send back response to browser
	fmt.Fprintln(w, responseBody)

}
