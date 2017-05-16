//Copyright 2017 Mosherr
//
//Permission is hereby granted, free of charge,
// to any person obtaining a copy of this software and associated documentation files (the "Software"),
// to deal in the Software without restriction,
// including without limitation the rights to use, copy, modify, merge, publish,
// distribute, sublicense, and/or sell copies of the Software,
// and to permit persons to whom the Software is furnished to do so,
// subject to the following conditions:
//
//The above copyright notice and this permission notice shall be included in all copies
// or substantial portions of the Software.
//
//THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND,
// EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT.
// IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM,
// DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE,
// ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS
// IN THE SOFTWARE.
package main

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"
	"unicode"
	"os"
	"github.com/Sirupsen/logrus"
)

var (
	logger         = logrus.WithField("cmd", "quotebot")
)

type slashResponse struct {
	ResponseType string `json:"response_type"`
	Text         string `json:"text"`
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		logger.WithField("PORT", port).Fatal("$PORT must be set")
	}

	http.HandleFunc("/", handleAction)
	logger.Println(http.ListenAndServe(":"+port, nil))
}

// If the token parameter doesn't match the secret token provided by Slack, we
// reject the request.  To let the app know what this token is, when you create
// the custom integration, populate the token variable in config.go.
func handleAction(w http.ResponseWriter, r *http.Request) {
	if r.PostFormValue("token") != token {
		http.Error(w, "Invalid Slack token.", http.StatusBadRequest)
		return
	}

	w.Header().Set("content-type", "application/json")

	input := r.PostFormValue("text")

	parts := parseSlackInput(input)
	var str string
	var usr string

	var resp *slashResponse

	if len(parts) == 1 {
		// random quote
		resp = handleGetQuote("random")
	} else {
		cmd := parts[0]
		switch cmd{
			case "-a":
				if len(parts) != 3 {
					// error not enough options
					resp = &slashResponse{
						ResponseType: "in_channel",
						Text:         "Incorrect options to add quote.",
					}
				} else {
					usr = parts[1]
					str = parts[2]
					addedBy := r.PostFormValue("user_name")
					resp = handleAddQuote(usr, str, addedBy)
				}
			case "-g":
				usr = parts[1]
				resp = handleGetQuote(usr)
			default:
				resp = handleGetQuote("random")
		}

	}

	if err := json.NewEncoder(w).Encode(resp); err != nil {
		http.Error(w, "Error encoding JSON.", http.StatusInternalServerError)
		return
	}
}

// Takes a user to return a quote from
// "random" will return a random quote from any user
// If username is passed try to find a quote by them, if none exists return error
func handleGetQuote(usr string) (*slashResponse) {

	var result string

	quote, err := DB.GetQuote(usr)

	if err != nil {
		// failed to get result from db
		result = "No quote found for " + usr
	}

	if quote.User == "" {
		// If no match from the db
		result = "No quote found for " + usr
	}

	if len(result) == 0 {
		result = quote.Text + " -" + quote.User
	}

	resp := &slashResponse{
		ResponseType: "in_channel",
		Text:         result,
	}

	return resp
}

// username to save quote for
// text of the quote
// user who added the quote
func handleAddQuote(usr string, quoteText string, addedBy string) (*slashResponse) {

	var result string

	insert := &Quote{
		User: usr,
		Text: quoteText,
		AddedBy: addedBy,
		DateAdded: time.Now(),
	}

	err := DB.AddQuote(insert)

	if err != nil {
		// error inserting
		result = "Failed to save quote"
	} else {
		result = quoteText + " added for user " + usr
	}

	resp := &slashResponse{
		ResponseType: "in_channel",
		Text:         result,
	}

	return resp
}

// Parse the slack input
func parseSlackInput(text string)(parsedInput []string) {
	lastQuote := rune(0)
	f := func(c rune) bool {
		switch {
		case c == lastQuote:
			lastQuote = rune(0)
			return false
		case lastQuote != rune(0):
			return false
		case unicode.In(c, unicode.Quotation_Mark):
			lastQuote = c
			return false
		default:
			return unicode.IsSpace(c)

		}
	}

	m := strings.FieldsFunc(text, f)

	return m
}
