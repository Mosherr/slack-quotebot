package quotebot

import (
	"fmt"
	"cloud.google.com/go/datastore"
	"golang.org/x/net/context"
)

// datastoreDB persists quotes to Cloud Datastore.
// https://cloud.google.com/datastore/docs/concepts/overview
type datastoreDB struct {
	client *datastore.Client
}

// Ensure datastoreDB conforms to the QuoteDatabase interface.
var _ QuoteDatabase = &datastoreDB{}

// newDatastoreDB creates a new QuoteDatabase backed by Cloud Datastore.
// See the datastore and google packages for details on creating a suitable Client:
// https://godoc.org/cloud.google.com/go/datastore
func newDatastoreDB(client *datastore.Client) (QuoteDatabase, error) {
	ctx := context.Background()
	// Verify that we can communicate and authenticate with the datastore service.
	t, err := client.NewTransaction(ctx)
	if err != nil {
		return nil, fmt.Errorf("datastoredb: could not connect: %v", err)
	}
	if err := t.Rollback(); err != nil {
		return nil, fmt.Errorf("datastoredb: could not connect: %v", err)
	}
	return &datastoreDB{
		client: client,
	}, nil
}

// Close closes the database.
func (db *datastoreDB) Close() {
	// No op.
}

// GetQuote retrieves a quote by its User.
func (db *datastoreDB) GetQuote(usr string) (*Quote, error) {
	q := &Quote{}
	ctx := context.Background()

	var err error
	var count int
	var quotes []*Quote

	if usr == "random" {
		query := datastore.NewQuery("Quote")
		count, err = db.client.Count(ctx, query)
		if count > 0 {
			query = datastore.NewQuery("Quote").
				Filter("Id =", randomRecordNumber(0, count)).
				Limit(1)
			_, err = db.client.GetAll(ctx, query, &quotes)
		}
	} else {
		query := datastore.NewQuery("Quote").
			Filter("User =", usr)
		count, err = db.client.Count(ctx, query)
		if count > 0 {
			query := datastore.NewQuery("Quote").
				Filter("User =", usr).
				Filter("Id =", randomRecordNumber(0, count)).
				Limit(1)
			_, err = db.client.GetAll(ctx, query, &quotes)
		}
	}

	if err != nil {
		return nil, err
	} else {
		q = quotes[0]
		return q, nil
	}
}

// AddQuote saves a given quote.
func (db *datastoreDB) AddQuote(q *Quote) (err error) {
	ctx := context.Background()
	k := datastore.IncompleteKey("Quote", nil)
	k, err = db.client.Put(ctx, k, q)
	if err != nil {
		return fmt.Errorf("datastoredb: could not put Quote: %v", err)
	}
	return nil
}