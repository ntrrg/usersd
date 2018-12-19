// Copyright 2018 Miguel Angel Rivera Notararigo. All rights reserved.
// This source code was released under the WTFPL.

package usersd

import (
	"encoding/json"
	"sync"
	"time"

	"github.com/blevesearch/bleve"
	"github.com/blevesearch/bleve/search/query"
	"github.com/dgraph-io/badger"
	"github.com/gofrs/uuid"
	"golang.org/x/crypto/bcrypt"
)

var (
	dbKeyPrefixUsers = "users-"

	// Password hashing strength.
	bcryptCost = 10
)

// User is an entity that may be authenticated and authorized.
type User struct {
	ID        string `json:"id"`
	Email     string `json:"email"`
	Mode      string `json:"mode"`
	Password  string `json:"password,omitempty"`
	Verified  bool   `json:"verified"`
	CreatedAt int64  `json:"createdAt"`
	LastLogin int64  `json:"lastLogin"`

	mu   sync.Mutex
	Data map[string]interface{} `json:"data,omitempty"`
}

// GetUser fetches a user with the given ID from the database.
func GetUser(tx *badger.Txn, id string) (*User, error) {
	item, err := tx.Get([]byte(dbKeyPrefixUsers + id))
	if err == badger.ErrKeyNotFound {
		return nil, ErrUserIDNotFound
	} else if err != nil {
		return nil, err
	}

	data, err := item.ValueCopy(nil)
	if err != nil {
		return nil, err
	}

	user := new(User)
	if err := json.Unmarshal(data, user); err != nil {
		return nil, err
	}

	return user, nil
}

// GetUsers fetches users that satisfies the given constraints.
func GetUsers(tx *badger.Txn, index bleve.Index, q string, sort ...string) ([]*User, error) {
	if q == "" && len(sort) == 0 {
		return getAllUsers(tx)
	}

	var (
		users []*User
		bq    query.Query
	)

	if q != "" {
		bq = bleve.NewQueryStringQuery(q)
	} else {
		bq = bleve.NewMatchAllQuery()
	}

	req := bleve.NewSearchRequest(bq)
	req.SortBy(sort)

	res, err := index.Search(req)
	if err != nil {
		return nil, err
	}

	for _, hit := range res.Hits {
		user, err := GetUser(tx, hit.ID)
		if err != nil {
			return nil, err
		}

		users = append(users, user)
	}

	return users, nil
}

// NewUser creates a user with the given arguments and writes it to the
// database. Receives a Badger transaction, a Bleve search index, an optional
// user ID, an email and a password; returns a User instance and an error if
// any.
func NewUser(tx *badger.Txn, index bleve.Index, rules []func(tx *badger.Txn, index bleve.Index, user *User, old *User) error, id, email, password string, data map[string]interface{}) (*User, error) {
	u := &User{ID: id, Email: email, Password: password, Data: data}
	if err := u.Write(tx, index, rules); err != nil {
		return nil, err
	}

	return u, nil
}

// Get gets the given value at the given key at Data field.
func (u *User) Get(key string) interface{} {
	v, ok := u.Data[key]

	if !ok {
		return nil
	}

	return v
}

// String implements fmt.Stringer.
func (u *User) String() string {
	return "(" + u.ID + ") " + u.Email
}

// Validate checks the user data and returns any errors.
func (u *User) Validate(tx *badger.Txn, index bleve.Index, rules []func(tx *badger.Txn, index bleve.Index, user *User, old *User) error) error {
	if rules == nil {
		rules = []func(tx *badger.Txn, index bleve.Index, user *User, old *User) error{
			UserIDValidator,
			UserModeValidator,
			UserEmailValidator,
			UserVerifiedValidator,
			UserPasswordValidator,
			UserCreatedAtValidator,
		}
	}

	old, err := GetUser(tx, u.ID)
	if err != nil && err != ErrUserIDNotFound {
		return err
	}

	errors := ValidationErrors{}

	for _, f := range rules {
		if err := f(tx, index, u, old); err != nil {
			errors = append(errors, err)
		}
	}

	if len(errors) > 0 {
		return errors
	}

	return nil
}

// Set sets the given value at the given key at Data field.
func (u *User) Set(key string, value interface{}) {
	u.mu.Lock()
	defer u.mu.Unlock()
	u.Data[key] = value
}

// Unset removes the given key from Data field.
func (u *User) Unset(key string) {
	u.mu.Lock()
	defer u.mu.Unlock()
	delete(u.Data, key)
}

// Write writes the user data to the database.
func (u *User) Write(tx *badger.Txn, index bleve.Index, rules []func(tx *badger.Txn, index bleve.Index, user *User, old *User) error) error {
	if err := u.Validate(tx, index, rules); err != nil {
		return err
	}

	v, err := json.Marshal(u)
	if err != nil {
		return err
	}

	if err := tx.Set([]byte(dbKeyPrefixUsers+u.ID), v); err != nil {
		return err
	}

	if err := index.Index(u.ID, u); err != nil {
		return err
	}

	return nil
}

// UserIDValidator validates the user ID.
func UserIDValidator(tx *badger.Txn, index bleve.Index, user *User, old *User) error {
	if user.ID == "" {
		x, err := uuid.NewV4()
		if err != nil {
			return err
		}

		user.ID = x.String()
	}

	return nil
}

// UserModeValidator validates from where the user were created. If none is
// given, the user will get 'local' as value.
func UserModeValidator(tx *badger.Txn, index bleve.Index, user *User, old *User) error {
	if user.Mode == "" {
		user.Mode = "local"
	}

	return nil
}

// UserEmailValidator validates the user email.
func UserEmailValidator(tx *badger.Txn, index bleve.Index, user *User, old *User) error {
	if user.Email == "" {
		return ErrUserEmailEmpty
	}

	q := `+email:"` + user.Email + `"`

	if old != nil {
		q = `-id:"` + user.ID + `" ` + q
	}

	users, err := GetUsers(tx, index, q)
	if err == nil && len(users) > 0 {
		return ErrUserEmailExists
	}

	return nil
}

// UserVerifiedValidator validates that the user email is verified.
func UserVerifiedValidator(tx *badger.Txn, index bleve.Index, user *User, old *User) error {
	if old == nil || user.Email != old.Email {
		user.Verified = false
	}

	return nil
}

// UserPasswordValidator validates the user password.
func UserPasswordValidator(tx *badger.Txn, index bleve.Index, user *User, old *User) error {
	if user.Password == "" {
		return ErrUserPasswordEmpty
	}

	password := user.Password

	if _, err := bcrypt.Cost([]byte(password)); err == nil {
		return nil
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcryptCost)
	if err != nil {
		return ErrUserPasswordHash.Format(err)
	}

	user.Password = string(hash)
	return nil
}

// UserCreatedAtValidator validates the users creation date.
func UserCreatedAtValidator(tx *badger.Txn, index bleve.Index, user *User, old *User) error {
	if old == nil {
		user.CreatedAt = time.Now().Unix()
	} else {
		user.CreatedAt = old.CreatedAt
	}

	return nil
}

func getAllUsers(tx *badger.Txn) ([]*User, error) {
	var users []*User

	it := tx.NewIterator(badger.DefaultIteratorOptions)
	defer it.Close()

	var (
		v []byte

		prefix = []byte(dbKeyPrefixUsers)
	)

	for it.Seek(prefix); it.ValidForPrefix(prefix); it.Next() {
		var err error
		item := it.Item()
		v, err = item.ValueCopy(v)
		if err != nil {
			return nil, err
		}

		user := new(User)
		if err := json.Unmarshal(v, user); err != nil {
			return nil, err
		}

		users = append(users, user)
	}

	return users, nil
}

// // CreateUserJSON creates a user with the given JSON, populates some required
// // data if missing and write it to the database.
// func CreateUserJSON(data []byte) (*User, error) {
// 	u, err := NewUserJSON(data)
//
// 	if err != nil {
// 		return nil, err
// 	}
//
// 	return NewUser(u.ID, u.Data)
// }
//
// // Delete removes the user from the database.
// func (u *User) Delete() error {
// 	txn := db.NewTransaction(true)
// 	defer txn.Discard()
//
// 	if err := txn.Delete([]byte(u.ID)); err != nil {
// 		msg := lError + " Can't remove the user (%v) from the database -> %v"
// 		l.Printf(msg, u.ID, err)
// 		return err
// 	}
//
// 	if err := index.Delete(u.ID); err != nil {
// 		msg := lError + " Can't remove the user (%v) from the search index -> %v"
// 		l.Printf(msg, u.ID, err)
// 		return err
// 	}
//
// 	if err := txn.Commit(); err != nil {
// 		msg := lError + " Can't commit changes to the database -> %v"
// 		l.Printf(msg, err)
// 		return err
// 	}
//
// 	if debug {
// 		l.Printf(lDebug+" User (%v) data removed", u.ID)
// 	}
//
// 	return nil
// }
//
// // Update updates the user data in the database.
// func (u *User) Update() error {
// 	if err := u.Validate(UserValidationRules); err != nil {
// 		return err
// 	}
//
// 	if _, err := GetUser(u.ID); err != nil {
// 		if err == ErrUserNotFound {
// 			l.Printf(lError+" Can't update, the user doesn't exists (%v)", u.ID)
// 		}
//
// 		return err
// 	}
//
// 	return u.write()
// }
