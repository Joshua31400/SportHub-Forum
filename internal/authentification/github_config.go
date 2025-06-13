package authentification

import (
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/github"
	"os"
)

var GitHubConfig *oauth2.Config

// InitGitHubConfig initializes GitHub OAuth2 configuration
func InitGitHubConfig() {
	clientID := os.Getenv("GITHUB_CLIENT_ID")
	clientSecret := os.Getenv("GITHUB_CLIENT_SECRET")
	redirectURL := os.Getenv("GITHUB_REDIRECT_URL")

	GitHubConfig = &oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		RedirectURL:  redirectURL,
		Scopes: []string{
			"user:email",
			"read:user",
		},
		Endpoint: github.Endpoint,
	}
}
