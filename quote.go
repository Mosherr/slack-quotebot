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
	"gopkg.in/mgo.v2/bson"
	"time"
)

// Quote holds metadata about a quote.
type Quote struct {
	Id        bson.ObjectId `json:"id" bson:"_id,omitempty"`
	User      string `json:"user"`
	Text      string `json:"text"`
	AddedBy   string `json:"added_by"`
	DateAdded time.Time `json:"date_added"`
	InsertId  int  `json:"insert_id"`
}

// QuoteDatabase provides thread-safe access to a database of quotes.
type QuoteDatabase interface {
	// GetQuote retrieves a quote by its User.
	GetQuote(usr string) (*Quote, error)

	// AddQuote saves a given quote, assigning it a new ID.
	AddQuote(q *Quote) (err error)

	Close()
}
