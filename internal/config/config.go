package config

import (
	"fmt"
	"github.com/alexedwards/scs/v2"
	"golang.org/x/oauth2"
	"html/template"
	"log"
	"os"
)

// AppConfig holds the application config
type AppConfig struct {
	UseCache      bool
	TemplateCache map[string]*template.Template
	InfoLog       log.Logger
	Prod          bool
	Session       *scs.SessionManager
}

// AuthDiscordConfig returns oauth2 config for Discord
func AuthDiscordConfig() *oauth2.Config {

	discordEndpoint := oauth2.Endpoint{
		AuthURL:  "https://discord.com/api/oauth2/authorize",
		TokenURL: "https://discord.com/api/oauth2/token",
	}

	scopes := []string{"identify"}

	protocol := "http://"
	if os.Getenv("SECURE") == "true" {
		protocol = "https://"
	}

	config := &oauth2.Config{
		ClientID:     os.Getenv("DISCORD_CLIENT_ID"),
		ClientSecret: os.Getenv("DISCORD_SECRET"),
		Endpoint:     discordEndpoint,
		RedirectURL:  fmt.Sprint(protocol + os.Getenv("DOMAIN") + "/user/auth/callback"),
		Scopes:       scopes,
	}

	return config
}
