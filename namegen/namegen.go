package namegen

import (
	"embed"
	"encoding/json"
	"log"
	"math/rand"
	"strings"
	"time"
)

//go:embed adjectives.json animals.json
var files embed.FS
var adjectives []string
var animals []string

func init() {
	adjectivesFile, err := files.ReadFile("adjectives.json")
	if err != nil {
		log.Println("failed to read adjectives.json")
		panic(err)
	}
	if err := json.Unmarshal(adjectivesFile, &adjectives); err != nil {
		log.Println("failed to unmarshal adjectives.json")
		panic(err)
	}

	animalsFile, err := files.ReadFile("animals.json")
	if err != nil {
		log.Println("failed to read animals.json")
		panic(err)
	}
	if err := json.Unmarshal(animalsFile, &animals); err != nil {
		log.Println("failed to unmarshal animals.json")
		panic(err)
	}

	rng = rand.New(rand.NewSource(time.Now().UnixNano()))
}

var rng *rand.Rand

func Generate() string {
	adjective := strings.Title(adjectives[rng.Intn(len(adjectives))])
	animal := strings.Title(animals[rng.Intn(len(animals))])
	return adjective + " " + animal
}
