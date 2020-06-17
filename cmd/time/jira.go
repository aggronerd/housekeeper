// Copyright 2019 Gregory Doran <greg@gregorydoran.co.uk>.
// All rights reserved.

package time

import (
	"github.com/andygrunwald/go-jira"
	"github.com/spf13/viper"
	"log"
	"net/http"
)

type transport interface {
	Client() *http.Client
}

func getTransport() transport {
	authMethod := viper.GetString("jira.auth")
	switch authMethod {
	case "basic":
		username := viper.GetString("jira.username")
		password := viper.GetString("jira.password")
		transport := jira.BasicAuthTransport{
			Username: username,
			Password: password,
		}
		return &transport
	default:
		log.Fatal("Jira is not configured properly")
	}
	return nil
}

func getJiraClient() *jira.Client {
	transport := getTransport()
	client, err := jira.NewClient(transport.Client(), viper.GetString("jira.url"))

	if err != nil {
		log.Fatalf("Failed to connect to Jira, please check your configuration: %s", err)
	}

	return client
}
