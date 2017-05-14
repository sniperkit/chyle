package senders

import (
	"github.com/antham/chyle/chyle/types"
)

// Sender define where the date must be sent
type Sender interface {
	Send(changelog *types.Changelog) error
}

// Send forward changelog to senders
func Send(senders *[]Sender, changelog *types.Changelog) error {
	for _, sender := range *senders {
		err := sender.Send(changelog)

		if err != nil {
			return err
		}
	}

	return nil
}

// CreateSenders build senders from a config
func CreateSenders(features Features, senders Config) *[]Sender {
	results := []Sender{}

	if features.GITHUBRELEASE {
		results = append(results, buildGithubReleaseSender(senders.GITHUBRELEASE))
	}

	if features.STDOUT {
		results = append(results, buildStdoutSender(senders.STDOUT))
	}

	return &results
}
