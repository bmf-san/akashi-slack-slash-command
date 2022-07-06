package SlackApi

import (
	"context"
	"encoding/base64"
	"errors"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/bmf-san/akashigo/client"
	"github.com/bmf-san/akashigo/stamp"
	"github.com/bmf-san/akashigo/types"
	"github.com/slack-go/slack"
	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
)

var (
	spreadSheetID string
	sheetClient   *sheets.Service
)

func init() {
	spreadSheetID = os.Getenv("SPREAD_SHEET_ID")
	serviceAccount := os.Getenv("SERVICE_ACCOUNT")

	dec, err := base64.StdEncoding.DecodeString(serviceAccount)
	if err != nil {
		log.Fatal(err)
	}
	cred := option.WithCredentialsJSON(dec)

	srv, err := sheets.NewService(context.TODO(), cred)
	if err != nil {
		log.Fatal(err)
		panic(err)
	}

	sheetClient = srv
}

func getAkashiClient(userName string, sheetClient *sheets.Service) (*client.Client, error) {
	resp, err := sheetClient.Spreadsheets.Values.Get(spreadSheetID, "シート1!A:D").Do()
	if err != nil {
		log.Fatal(err)
	}

	if len(resp.Values) == 0 {
		return nil, errors.New("no data")
	}

	var companyID string
	var apiToken string
	for _, row := range resp.Values {
		if row[1].(string) == userName {
			companyID = row[2].(string)
			apiToken = row[3].(string)
		}
	}

	return client.New(apiToken, companyID), nil
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
	case "/akashi":
		apiClient, err := getAkashiClient(s.UserName, sheetClient)
		if err != nil {
			log.Fatal(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		model := &stamp.Stamp{
			Client: apiClient,
		}
		stamp, err := model.Stamp(stamp.StampParams{
			Token:    apiClient.APIToken,
			Type:     types.StampNumber(s.Text),
			Timezone: "+09:00",
		})
		if err != nil {
			log.Fatal(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		if stamp.Success {
			typeStr := types.StampText(stamp.Response.Type)
			w.Write([]byte(fmt.Sprintf("%s %s", typeStr, stamp.Response.StampedAt)))
		} else {
			var msg string
			for _, v := range stamp.Errors {
				msg += v.Message
				if v.Code == "ERR300003" {
					msg += "スプレッドシートのakashi_api_tokenを更新してください。"
				}
			}
			w.Write([]byte(msg))
		}
	default:
		log.Fatal(errors.New("invalid command"))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
