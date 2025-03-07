package main

import (
	"context"
	"fmt"
	"github.com/adlio/trello"
	"github.com/joho/godotenv"
	"golang.org/x/time/rate"
	"log"
	"os"
	"path/filepath"
	"sync"
	"time"
)

var (
	apiKey    string
	apiToken  string
	boardID   string
	listID    string
	client    *trello.Client
	limiter   *rate.Limiter
	fileCount int
	dir       string
)

func getListIdByName(listName string) (string, error) {
	if listName == "" {
		return "", fmt.Errorf("list name not provided")
	}

	board, err := client.GetBoard(boardID, trello.Defaults())
	if err != nil {
		return "", err
	}

	lists, err := board.GetLists(trello.Defaults())
	if err != nil {
		return "", err
	}

	for _, list := range lists {
		if list.Name == listName {
			return list.ID, nil
		}
	}

	return "", fmt.Errorf("list %s not found", listName)
}

func processTask(path string, info os.FileInfo, wg *sync.WaitGroup) error {
	defer wg.Done()
	if info.IsDir() {
		return nil // skip directories
	}
	if filepath.Ext(path) == ".txt" {
		fileCount++

		// Read file content
		content, err := os.ReadFile(path)
		if err != nil {
			return err
		}

		// Rate limit the requests
		err = limiter.Wait(context.Background())
		if err != nil {
			return err
		}

		// Create a new card on Trello
		card := trello.Card{
			Name:    info.Name(),     // Use file name as card title
			Desc:    string(content), // Use file content as card description
			IDBoard: boardID,
		}

		card.IDList = listID

		// Send the card to Trello
		err = client.CreateCard(&card)
		if err != nil {
			fmt.Println("Error creating card:", err)
			os.Remove(path)
			return nil
		}

		fmt.Printf("Created card for %s\n", info.Name())

		//remove file from directory
		os.Remove(path)
	}
	return nil
}

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	//load env variables
	apiKey = os.Getenv("API_KEY")
	apiToken = os.Getenv("API_TOKEN")
	boardID = os.Getenv("BOARD_ID")
	client = trello.NewClient(apiKey, apiToken)
	limiter = rate.NewLimiter(rate.Every(time.Second), 1) // adjust the rate limit as needed
	dir = "files"

	listID, err = getListIdByName(os.Getenv("LIST_NAME"))
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	start := time.Now()

	wg := sync.WaitGroup{}

	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		wg.Add(1)
		go processTask(path, info, &wg)

		return nil
	})

	if err != nil {
		log.Fatal(err)
	}

	wg.Wait()

	fmt.Println("Total files processed:", fileCount)
	fmt.Println("Execution time:", time.Since(start))
}
