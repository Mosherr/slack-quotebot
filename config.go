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
	"log"
)

func init() {
	var err error
	var cred = &mgo.Credential{
		Username:AuthUserName,
		Password:AuthPassword,
		Source:DB_NAME,
	}

	DB, err = newMongoDB(MongoDBHosts, cred)

	if err != nil {
		log.Fatal(err)
	}
}

var (
	// Set the token variable. It is needed to verify that the
	// requests to the slash command come from Slack. It is provided for
	// you by Slack when you create the Slash command as a custom
	// integration. https://my.slack.com/services/new/slash-commands
	token string = "xxxxxx"

	MongoDBHosts string = "xxxxxx"
	AuthUserName string = "xxxxxx"
	AuthPassword string = "xxxxxx"

	DB QuoteDatabase

	// Force import of mgo library.
	_ mgo.Session
)