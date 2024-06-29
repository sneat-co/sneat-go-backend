package gfuncs

import (
	"cloud.google.com/go/firestore"
	"context"
	"fmt"
	"log"
	"time"
)

type dbClient interface {
	RunTransaction(ctx context.Context, f func(ctx context.Context, transaction *firestore.Transaction) error, opts ...firestore.TransactionOption) error
	Collection(path string) *firestore.CollectionRef
}

var db dbClient

var newDbClient = func() (dbClient, error) {
	return firestore.NewClient(context.Background(), "sneat-team")
}

var users *firestore.CollectionRef

func initDb() {
	var err error
	if db, err = newDbClient(); err != nil {
		logus.Fatalf("Failed to init firestore client: %v", err)
	}
	users = db.Collection("users")
}

// AuthEvent is the payload of a Firestore Auth event.
type AuthEvent struct {
	UID           string `json:"uid"`
	DisplayName   string `json:"displayName,omitempty"`
	Email         string `json:"email,omitempty"`
	EmailVerified bool   `json:"emailVerified,omitempty"`
	Metadata      struct {
		CreatedAt time.Time `json:"createdAt,omitempty"`
	} `json:"metadata"`
}

// User record
type User struct {
	Created       time.Time `json:"created" firestore:"created"`
	Title         string    `json:"title,omitempty" firestore:"title,omitempty"`
	Email         string    `json:"email,omitempty" firestore:"email,omitempty"`
	EmailVerified bool      `json:"emailVerified,omitempty" firestore:"emailVerified,omitempty"`
}

// AuthUserCreated is triggered by Firestore Auth event
func AuthUserCreated(ctx context.Context, e AuthEvent) error {
	if e.UID == "" {
		return nil
	}
	if db == nil {
		initDb()
	}
	if err := db.RunTransaction(ctx, func(ctx context.Context, transaction *firestore.Transaction) (err error) {
		userDocRef := users.Doc(e.UID)
		return transaction.Create(
			userDocRef,
			User{
				Email:         e.Email,
				EmailVerified: e.EmailVerified,
				Title:         e.DisplayName,
				Created:       e.Metadata.CreatedAt,
			})
	}); err != nil {
		return fmt.Errorf("failed to create user record: %w", err)
	}
	return nil
}
