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

	"github.com/bmf-san/akashigo/client"
	"github.com/bmf-san/akashigo/stamp"
	"github.com/bmf-san/akashigo/types"
	"github.com/slack-go/slack"
)

var (
	signingSecret string
	apiclient     *client.Client
)

func init() {
	// TODO: apiTokenはストレージから取得するようにするので、これは後で破棄
	apiToken := os.Getenv("AKASHI_API_TOKEN")
	companyID := os.Getenv("AKASHI_COMPANY_ID")
	apiclient = client.New(apiToken, companyID)
}

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
<<<<<<< HEAD
	case "/hello":
		params := &slack.Msg{Text: s.Text}
		b, err := json.Marshal(params)
=======
	case "/akashi":
		params := &slack.Msg{Text: s.Text}
		p, err := json.Marshal(params)
		if err != nil {
			log.Fatal(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		// TODO: paramsのバリデーションをする

		model := &stamp.Stamp{
			Client: apiclient,
		}
		s, err := model.Stamp(stamp.StampParams{
			Token:    apiclient.APIToken,
			Type:     types.StampNumber(string(p)),
			Timezone: "+09:00",
		})
>>>>>>> 61d7416 ([update] minumum implement)
		if err != nil {
			log.Fatal(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		b, err := json.Marshal(s)
		if err != nil {
			log.Fatal(err)
			w.WriteHeader(http.StatusInternalServerError)
		}

		// TODO: 期限切れの場合はトークン再発行APIをコールして、再度リクエスト
		//

		w.Header().Set("Content-Type", "application/json")
		w.Write(b)
	default:
		log.Fatal(errors.New("Invalid command"))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
