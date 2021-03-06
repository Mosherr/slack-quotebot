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
package quotebot

import (
	"gopkg.in/mgo.v2"
	"cloud.google.com/go/datastore"
	"golang.org/x/net/context"
	"log"
	"time"
	"math/rand"
)

func init() {
	var err error

	// [START mongo]
	// To use Mongo, uncomment the next lines and update the address string and
	// optionally, the credentials
	var cred = &mgo.Credential{
		Username:AuthUserName,
		Password:AuthPassword,
		Source:DB_NAME,
	}
	DB, err = newMongoDB(MongoDBHosts, cred)
	// [END mongo]

	// [START datastore]
	// To use Cloud Datastore, uncomment the following lines and update the
	// project ID.
	// More options can be set, see the google package docs for details:
	// http://godoc.org/golang.org/x/oauth2/google
	//
	//DB, err = configureDatastoreDB("slackquotebot")
	// [END datastore]

	if err != nil {
		log.Fatal(err)
	}

	rand.Seed(time.Now().UTC().UnixNano())

}

var (
	// Set the token variable. It is needed to verify that the
	// requests to the slash command come from Slack. It is provided for
	// you by Slack when you create the Slash command as a custom
	// integration. https://my.slack.com/services/new/slash-commands
	token string = "xxxxxx"

	// [START mongo]
	MongoDBHosts string = "xxxxxx"
	AuthUserName string = "xxxxxx"
	AuthPassword string = "xxxxxx"
	// [END mongo]

	DB QuoteDatabase

	// Force import of mgo library.
	_ mgo.Session
)

func configureDatastoreDB(projectID string) (QuoteDatabase, error) {
	ctx := context.Background()
	client, err := datastore.NewClient(ctx, projectID)
	if err != nil {
		return nil, err
	}
	return newDatastoreDB(client)
}
