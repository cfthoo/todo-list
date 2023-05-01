package config

import (
	"os"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/facebook"
	"golang.org/x/oauth2/github"
	"golang.org/x/oauth2/google"
)

const GoogleOauthUrl = "https://www.googleapis.com/oauth2/v2/userinfo?access_token="
const FacebookOauthUrl = "https://graph.facebook.com/v13.0/me?fields=id,name,email,picture&access_token&access_token="
const GithubOauthUrl = "https://api.github.com/user"

func SetupGoogleConfig() *oauth2.Config {
	conf := &oauth2.Config{
		ClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
		ClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
		RedirectURL:  "http://localhost:8080/google/callback",
		Scopes: []string{
			"https://www.googleapis.com/auth/userinfo.email",
			"https://www.googleapis.com/auth/userinfo.profile",
		},
		Endpoint: google.Endpoint,
	}
	return conf
}

func SetupFbConfig() *oauth2.Config {
	conf := &oauth2.Config{
		ClientID:     os.Getenv("FB_CLIENT_ID"),
		ClientSecret: os.Getenv("FB_CLIENT_SECRET"),
		RedirectURL:  "http://localhost:8080/fb/callback",
		Scopes: []string{
			"email",
			"public_profile",
		},
		Endpoint: facebook.Endpoint,
	}
	return conf
}

func SetupGithubConfig() *oauth2.Config {
	conf := &oauth2.Config{
		ClientID:     os.Getenv("GITHUB_CLIENT_ID"),
		ClientSecret: os.Getenv("GITHUB_CLIENT_SECRET"),
		RedirectURL:  "http://localhost:8080/github/callback",
		Scopes: []string{
			"repo",
			"user",
		},
		Endpoint: github.Endpoint,
	}
	return conf
}
