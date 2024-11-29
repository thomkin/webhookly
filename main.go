package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"

	"github.com/go-playground/webhooks/v6/github"
	"gopkg.in/yaml.v2"
)

type Handler struct {
	Ref  string
	Path string
}

type Handlers map[string]Handler

func main() {
	var config struct {
		Secret   string   `yaml:"secret"`
		Cert     string   `yaml:"cert"`
		Key      string   `yaml:"key"`
		Handlers Handlers `yaml:"handlers"`
	}
	data, err := os.ReadFile("config.yaml")
	if err != nil {
		log.Fatal(err)
	}
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		log.Fatal(err)
	}

	hook, _ := github.New(github.Options.Secret(config.Secret))
	http.HandleFunc("/github", func(w http.ResponseWriter, r *http.Request) {
		payload, err := hook.Parse(r, github.PushEvent, github.PullRequestEvent)
		fmt.Printf("payload: %v\n", payload)
		if err != nil {
			if err == github.ErrEventNotFound {
				println("Event not found")
				// ok event wasn't one of the ones asked to be parsed
			}
		}
		switch payload.(type) {

		case github.PushPayload:
			push := payload.(github.PushPayload)
			ref := push.Ref
			handler, ok := config.Handlers[ref]
			if !ok {
				println("Handler not found for ref: " + ref)
				return
			}

			if _, err := os.Stat(handler.Path); os.IsNotExist(err) {
				log.Fatalf("Handler %s does not exist", handler.Path)
			}

			cmd := exec.Command(handler.Path)
			output, err := cmd.CombinedOutput()
			if err != nil {
				log.Fatal(err)
			}
			fmt.Printf("%s", output)

		case github.PullRequestPayload:
			pullRequest := payload.(github.PullRequestPayload)
			// Do whatever you want from here...
			fmt.Printf("PullRequest: %+v", pullRequest)
		}
	})

	http.ListenAndServeTLS(":80", config.Cert, config.Key, nil)

}
