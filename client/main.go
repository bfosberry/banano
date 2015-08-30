package main

import (
	"log"

	"github.com/bfosberry/banano/nano"
)

func main() {
	log.Println("Starting..")
	repo := nano.NewRemoteRepository("localhost:8080")

	t1 := &nano.Thingey{
		ID:   "1",
		Data: "d1",
	}
	if err := repo.Create(t1); err != nil {
		log.Fatal(err)
	}
	log.Printf("Created %+v\n", t1.Data)

	if t2, err := repo.Get(t1.ID); err != nil {
		log.Fatal(err)
	} else {
		log.Printf("Got %+v\n", t2.Data)
	}

	if tList, err := repo.List(); err != nil {
		log.Fatal(err)
	} else {
		log.Printf("Listed %d items\n", len(tList))
	}

	if err := repo.Delete(t1); err != nil {
		log.Fatal(err)
	}
	log.Printf("Deleted %+v\n", t1.Data)

	if tList, err := repo.List(); err != nil {
		log.Fatal(err)
	} else {
		log.Printf("Listed %d items\n", len(tList))
	}
}
