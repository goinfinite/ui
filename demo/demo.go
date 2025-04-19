package main

import (
	"context"
	"log"
	"os"
)

func main() {
	f, err := os.Create("index.html")
	if err != nil {
		log.Fatalf("CreateHtmlFileFailed: %v", err)
	}

	err = DemoIndex().Render(context.Background(), f)
	if err != nil {
		log.Fatalf("WriteHtmlFileFailed: %v", err)
	}
}
