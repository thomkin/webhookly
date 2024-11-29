package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"

	"github.com/go-playground/webhooks/v6/github"
	"gopkg.in/yaml.v2"
)

type Handler struct {
	Path string `yaml:"path"`
}

type Handlers map[string]Handler

func main() {
	var config struct {
		Secret   string   `yaml:"secret"`
		Cert     string   `yaml:"cert"`
		Key      string   `yaml:"key"`
		Port     string   `yaml:"port"`
		Handlers Handlers `yaml:"handlers"`
	}

	configFile := flag.String("c", "config.yaml", "path to config file")
	flag.Parse()

	data, err := os.ReadFile(*configFile)
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

		if err != nil {
			if err == github.ErrEventNotFound {
				return
			}
		}

		switch payload.(type) {

		case github.PushPayload:
			push := payload.(github.PushPayload)
			ref := push.Ref
			handlerExecution(ref, config.Handlers)

			// case github.PullRequestPayload:
			// 	pullRequest := payload.(github.PullRequestPayload)
			// 	ref := pullRequest.
			// 	handlerExecution(ref, config.Handlers)

		}
	})

	http.ListenAndServeTLS(":"+config.Port, config.Cert, config.Key, nil)
}

func handlerExecution(ref string, handlers Handlers) {
	handler, ok := handlers[ref]
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
		fmt.Println(err)
	}

	//TODO: we can make this log a little nicer at some point
	fmt.Printf("%s", output)

}
