package api

import (
	"net/http"
	"net/url"
)

const (
	EmailServerHost = "http://127.0.0.1:8090"
)

func SendEmail(sub, body, to string) {
	http.PostForm(EmailServerHost,
		url.Values{
			"sub":  {sub},
			"body": {body},
			"to":   {to},
		})
}
