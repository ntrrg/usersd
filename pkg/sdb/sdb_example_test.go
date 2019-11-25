// Copyright 2019 Miguel Angel Rivera Notararigo. All rights reserved.
// This source code was released under the MIT license.

package sdb_test

import (
	"errors"
	"fmt"
	"io/ioutil"

	"github.com/blevesearch/bleve"
	"github.com/blevesearch/bleve/analysis/analyzer/keyword"
	"nt.web.ve/go/usersd/pkg/sdb"
)

type Car struct {
	ID    int
	Brand string
	Model string
}

type Person struct {
	ID, Name string
	Email    string
	Alive    bool
	Numbers  []int
	Car      Car
	Family   []Person

	Doctype string // Search index document type.
}

var p = Person{
	ID:    "ntrrg",
	Name:  "Miguel Angel Rivera Notararigo",
	Email: "ntrrg@example.com",
	Alive: true,

	Numbers: []int{11, 2, 0},

	Car: Car{
		ID:    1,
		Brand: "Toyota",
		Model: "Corolla Araya",
	},

	Family: []Person{
		{ID: "alirio", Name: "Alirio Rivera"},
		{ID: "assdro", Name: "Alessandro Notararigo"},
	},

	Doctype: "people",
}

func Example() {
	dir, err := ioutil.TempDir("", "sdb")
	if err != nil {
		panic(err)
	}

	opts := sdb.DefaultOptions(dir)

	// Advanced document mapping

	keywordField := bleve.NewTextFieldMapping()
	keywordField.Analyzer = keyword.Name

	peopleMapping := bleve.NewDocumentMapping()
	peopleMapping.AddFieldMappingsAt("Doctype", keywordField)
	// Without this, 'rrg' would match with 'ntrrg', 'atrrg', etc...
	peopleMapping.AddFieldMappingsAt("ID", keywordField)
	// Without this, 'example.com' would match any person with an email from this
	// domain.
	peopleMapping.AddFieldMappingsAt("Email", keywordField)
	// Without this, boolean fields couldn't be compared with 'true' of 'false'.
	peopleMapping.AddFieldMappingsAt("Alive", keywordField)

	opts.Bleve.DocMappings["people"] = peopleMapping

	db, err := sdb.OpenWith(opts)
	if err != nil {
		panic(err)
	}

	// If no advanced options are needed, all the previous lines could be
	// replaced by:
	//
	//   db, err := sdb.Open("")
	//   if err != nil {
	//     panic(err)
	//   }

	defer db.Close()

	fmt.Printf("Initial -> %s: %s\n", p.ID, p.Name)

	writeData(db)
	getData(db)
	deleteData(db)

	// Output:
	// Initial -> ntrrg: Miguel Angel Rivera Notararigo
	// Get -> ntrrg: Miguel Angel Rivera Notararigo
	// Find -> (Name:miguel): ["ntrrg"]
	// Find -> (Email:example.com): []
	// Delete -> ntrrg: Not found
}

func writeData(db *sdb.DB) {
	tx := db.NewTx(sdb.RW)
	defer tx.Discard()

	if err := tx.Set([]byte(p.ID), p); err != nil {
		panic(err)
	}

	if err := tx.Commit(); err != nil {
		panic(err)
	}
}

func getData(db *sdb.DB) {
	tx := db.NewTx(sdb.RW)
	defer tx.Discard()

	p2 := Person{}
	if err := tx.Get([]byte(p.ID), &p2); err != nil {
		panic(err)
	}

	fmt.Printf("Get -> %s: %s\n", p2.ID, p2.Name)

	q := "Name:miguel"

	keys, err := tx.Find(q) // Any document with "miguel" in its name.
	if err != nil {
		panic(err)
	}

	fmt.Printf("Find -> (%s): %q\n", q, keys)

	q = "Email:example.com"

	keys, err = tx.Find(q) // Any document with "example.com" as email.
	if err != nil {
		panic(err)
	}

	fmt.Printf("Find -> (%s): %q\n", q, keys)
}

func deleteData(db *sdb.DB) {
	tx := db.NewTx(sdb.RW)
	defer tx.Discard()

	if err := tx.Delete([]byte(p.ID)); err != nil {
		panic(err)
	}

	p3 := Person{}
	if err := tx.Get([]byte(p.ID), &p3); errors.Is(err, sdb.ErrKeyNotFound) {
		p3.ID = p.ID
		p3.Name = "Not found"
	} else if err != nil {
		panic(err)
	}

	if err := tx.Commit(); err != nil {
		panic(err)
	}

	fmt.Printf("Delete -> %s: %s\n", p3.ID, p3.Name)
}
