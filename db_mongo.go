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
	"fmt"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"math/rand"
)

type mongoDB struct {
	conn *mgo.Session
	c    *mgo.Collection
}

// Ensure mongoDB conforms to the QuoteDatabase interface.
var _ QuoteDatabase = (*mongoDB)(nil)

const (
	DB_NAME       = "quotestore"
	DB_COLLECTION = "quotes"
)

// newMongoDB creates a new QuoteDatabase backed by a given Mongo server,
// authenticated with given credentials.
func newMongoDB(addr string, cred *mgo.Credential) (QuoteDatabase, error) {
	conn, err := mgo.Dial(addr)
	if err != nil {
		return nil, fmt.Errorf("mongo: could not dial: %v", err)
	}

	if cred != nil {
		if err := conn.Login(cred); err != nil {
			return nil, fmt.Errorf("mongo: Invalid credentials: %v", err)
		}
	}

	return &mongoDB{
		conn: conn,
		c:    conn.DB(DB_NAME).C(DB_COLLECTION),
	}, nil
}

// Close closes the database.
func (db *mongoDB) Close() {
	db.conn.Close()
}

// GetQuote retrieves a quote by its User.
func (db *mongoDB) GetQuote(usr string) (*Quote, error) {
	q := &Quote{}

	var err error
	var count int

	if usr == "random" {
		count, err = db.c.Count()
		if count > 0 {
			err = db.c.Find(bson.M{}).Skip(randomRecordNumber(0, count)).One(q)
		}
	} else {
		count, err = db.c.Find(bson.M{"user": usr}).Count()
		if count > 0 {
			err = db.c.Find(bson.M{"user": usr}).Skip(randomRecordNumber(0, count)).One(q)
		}
	}

	if err != nil {
		return nil, err
	} else {
		return q, nil
	}
}

// AddQuote saves a given quote.
func (db *mongoDB) AddQuote(q *Quote) (err error) {
	err = db.c.Insert(q)
	if err != nil {
		return fmt.Errorf("mongodb: could not add quote: %v", err)
	}
	return nil
}

// randomRecordNumber returns a positive number that fits in the bounds.
func randomRecordNumber(min int, max int) (int) {
	return rand.Intn(max - min) + min
}