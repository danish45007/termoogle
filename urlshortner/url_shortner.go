package urlshortner

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/raahii/kutt-go"
)

func getEnvVariable(key string) string {

	// load .env file
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	return os.Getenv(key)
}

func UrlShortner(target string) string {
	cli := kutt.NewClient(getEnvVariable("API_KEY_KUTT"))
	URL, err := cli.Submit(
		target,
	)
	if err != nil {
		log.Fatal(err)
	}
	return URL.ShortURL
}
