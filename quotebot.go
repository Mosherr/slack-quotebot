// Copyright 2016 Google Inc. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Package quotebot demonstrates how to create an App Engine application as a
// Slack slash command.
package quotebot

import (
"encoding/json"
"html/template"
"math/rand"
"net/http"

"google.golang.org/appengine"
"google.golang.org/appengine/log"
	"strings"
)

type slashResponse struct {
	ResponseType string `json:"response_type"`
	Text         string `json:"text"`
}

var indexTmpl = template.Must(template.ParseFiles("index.html"))

func init() {
	http.HandleFunc("/", handleAction)
}

// If the token parameter doesn't match the secret token provided by Slack, we
// reject the request.  To let the app know what this token is, when you create
// the custom integration, populate the token variable in config.go.
func handleAction(w http.ResponseWriter, r *http.Request) {
	//if token != "" && r.PostFormValue("token") != token {
	//	http.Error(w, "Invalid Slack token.", http.StatusBadRequest)
	//	return
	//}

	w.Header().Set("content-type", "application/json")

	input := r.PostFormValue("text")
	//parts := strings.Fields(input)

	parts := strings.Split(input, "-")

	var resp *slashResponse

	if len(parts) == 1 {
		// random quote
		resp = handleGetQuote()
	} else {
		cmd := parts[1]
		switch cmd{
			case "-addquote":
				str := strings.Replace(cmd, "-addquote", "", -1)
				resp = &slashResponse{
					ResponseType: "in_channel",
					Text:         str,
				}
			case "-getquote":
				resp = &slashResponse{
					ResponseType: "in_channel",
					Text:         "getquote",
				}
			default:
				resp = &slashResponse{
					ResponseType: "in_channel",
					Text:         input,
				}
		}

	}

	if err := json.NewEncoder(w).Encode(resp); err != nil {
		c := appengine.NewContext(r)
		log.Errorf(c, "Error encoding JSON: %s", err)
		http.Error(w, "Error encoding JSON.", http.StatusInternalServerError)
		return
	}

	// split the input into commands
	// commands prefixed with '-'

	// -get random
	// -get @user
	// -add @user txt
	//switch input {
	
	//}
}

// Takes an option
// random or no option will return a random quote
// If username is passed try to find a quote by them, if none exists return error
func handleGetQuote() (*slashResponse) {
	resp := &slashResponse{
		ResponseType: "in_channel",
		Text:         quotes[rand.Intn(len(quotes))],
	}

	return resp
}

// takes two options
// username to save quote for
// text of the quote
// if both options are not passed error
func handleAddQuote(w http.ResponseWriter, r *http.Request) {

}
