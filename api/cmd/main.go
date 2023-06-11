package main

import (
	"log"
	"math/rand"
	"time"
)

func main() {
	rand.Seed(time.Now().UnixNano())
	turns := rand.Perm(4)
	log.Println(turns)
}
