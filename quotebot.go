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
	"math/rand"
	"net/http"
	"google.golang.org/appengine"
	"google.golang.org/appengine/log"
	"strings"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"time"
	"math"
)

type slashResponse struct {
	ResponseType string `json:"response_type"`
	Text         string `json:"text"`
}

type Quote struct {
	Id        bson.ObjectId `json:"id" bson:"_id,omitempty"`
	User      string `json:"user"`
	Text      string `json:"text"`
	AddedBy   string `json:"aded_by"`
	DateAdded time.Time `json:"date_added"`
}

const (
	DB_NAME       = "quotestore"
	DB_COLLECTION = "quotes"
)

func init() {
	http.HandleFunc("/", handleAction)
}

func ensureIndex(s *mgo.Session) {
	session := s.Copy()
	defer session.Close()

	c := session.DB(DB_NAME).C(DB_COLLECTION)

	index := mgo.Index{
		Key:        []string{"user"},
		Unique:     true,
		DropDups:   true,
		Background: true,
		Sparse:     true,
	}
	err := c.EnsureIndex(index)
	if err != nil {
		panic(err)
	}
}

// If the token parameter doesn't match the secret token provided by Slack, we
// reject the request.  To let the app know what this token is, when you create
// the custom integration, populate the token variable in config.go.
func handleAction(w http.ResponseWriter, r *http.Request) {
	if token != "" && r.PostFormValue("token") != token {
		http.Error(w, "Invalid Slack token.", http.StatusBadRequest)
		return
	}

	// We need this object to establish a session to our MongoDB.
	mongoDBDialInfo := &mgo.DialInfo{
		Addrs:    []string{MongoDBHosts},
		Timeout:  60 * time.Second,
		Database: AuthDatabase,
		Username: AuthUserName,
		Password: AuthPassword,
	}

	session, err := mgo.DialWithInfo(mongoDBDialInfo)

	if err != nil {
		panic(err)
	}
	defer session.Close()

	session.SetMode(mgo.Monotonic, true)
	ensureIndex(session)

	w.Header().Set("content-type", "application/json")

	input := r.PostFormValue("text")

	parts := strings.Split(input, "-")

	var resp *slashResponse

	if len(parts) == 1 {
		// random quote
		resp = handleGetQuote("random", session)
	} else {
		cmdParts := strings.Split(input, " ")
		cmd := cmdParts[0]
		switch cmd{
			case "-a":
				if len(parts) != 3 {
					// error not enough options
					resp = &slashResponse{
						ResponseType: "in_channel",
						Text:         "Incorrect options to add quote.",
					}
				} else {
					usr := strings.Replace(cmdParts[1], "-", "", -1)
					str := strings.Replace(input, "-a", "", -1)
					str = strings.Replace(input, cmdParts[1], "", -1)
					addedBy := r.PostFormValue("user_name")
					resp = handleAddQuote(usr, str, addedBy, session)
				}
			case "-g":
				str := strings.Replace(input, "-g", "", -1)
				resp = handleGetQuote(str, session)
			default:
				resp = handleGetQuote(cmd, session)
		}

	}

	if err := json.NewEncoder(w).Encode(resp); err != nil {
		c := appengine.NewContext(r)
		log.Errorf(c, "Error encoding JSON: %s", err)
		http.Error(w, "Error encoding JSON.", http.StatusInternalServerError)
		return
	}
}

// Takes a user to return a quote from
// "random" will return a random quote from any user
// If username is passed try to find a quote by them, if none exists return error
func handleGetQuote(usr string, s *mgo.Session) (*slashResponse) {

	var result string

	session := s.Copy()
	defer session.Close()

	c := session.DB(DB_NAME).C(DB_COLLECTION)

	var quote Quote

	var count, err = c.Count()

	if err != nil {
		// failed to get result from db
		result = "No quote found for" + usr
	} else {
		var randomInt = int(math.Floor( float64(rand.Int() * count )))

		if usr == "random" {
			err = c.Find(bson.M{}).Skip(randomInt).One(&quote)
		} else {
			err = c.Find(bson.M{"user": usr}).Skip(randomInt).One(&quote)
		}

		if err != nil {
			// failed to get result from db
			result = "No quote found for" + usr
		}

		if quote.User == "" {
			// If no match from the db
			result = "No quote found for" + usr
		}
	}

	resp := &slashResponse{
		ResponseType: "in_channel",
		Text:         result,
	}

	return resp
}

// takes two options
// username to save quote for
// text of the quote
// if both options are not passed error
func handleAddQuote(usr string, quoteText string, addedBy string, s *mgo.Session) (*slashResponse) {
	session := s.Copy()
	defer session.Close()

	var result string

	c := session.DB(DB_NAME).C(DB_COLLECTION)

	insert := &Quote{User: usr, Text: quoteText, AddedBy: addedBy, DateAdded: time.Now()}

	err := c.Insert(insert)
	if err != nil {
		// error inserting
		result = "Failed to save quote"
	} else {
		result = "'" + quoteText +"'" + "added for user" + usr
	}

	resp := &slashResponse{
		ResponseType: "in_channel",
		Text:         result,
	}

	return resp
}
