package SlackApi

import (
	"encoding/json"
	"errors"
	"flag"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/slack-go/slack"
)

var (
	signingSecret string
)

func Slash(w http.ResponseWriter, r *http.Request) {
	var signingSecret string
	flag.StringVar(&signingSecret, "secret", os.Getenv("SLACK_SIGINING_SECRET"), "Your Slack app's signing secret")
	flag.Parse()

	verifier, err := slack.NewSecretsVerifier(r.Header, signingSecret)
	if err != nil {
		log.Fatal(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	r.Body = ioutil.NopCloser(io.TeeReader(r.Body, &verifier))
	s, err := slack.SlashCommandParse(r)
	if err != nil {
		log.Fatal(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err = verifier.Ensure(); err != nil {
		log.Fatal(err)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	switch s.Command {
	case "/hello":
		params := &slack.Msg{Text: s.Text}
		b, err := json.Marshal(params)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(b)
	default:
		log.Fatal(errors.New("No command specified"))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
