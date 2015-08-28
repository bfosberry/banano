package main

import (
	"log"

	"github.com/bfosberry/banano/nano"
)

func main() {
	log.Println("Starting..")
	repo := nano.NewLocalRepository()

	t1 := &nano.Thingey{
		ID:   "1",
		Data: "d1",
	}
	if err := repo.Create(t1); err != nil {
		log.Fatal(err)
	}

	if t2, err := repo.Get(t1.ID); err != nil {
		log.Fatal(err)
	} else {
		log.Println(t2.Data)
	}

	if tList, err := repo.List(); err != nil {
		log.Fatal(err)
	} else {
		log.Println(len(tList))
	}

	if err := repo.Delete(t1); err != nil {
		log.Fatal(err)
	}

	if tList, err := repo.List(); err != nil {
		log.Fatal(err)
	} else {
		log.Println(len(tList))
	}
}
