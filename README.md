## Import `.txt` files to trello board list

### Prerequisites
1. Get `api key` and `token` from [power-up](https://trello.com/power-ups/admin) page.
2. Get `board id` from trello board url. For example, if the url is `https://trello.com/b/abc123/board-name`, then `abc123` is the board id.
3. Create board `list` in trello board.

### Usage
1. Create `.env` file by copying `.env.example` file and fill the required fields
2. Put `.txt` files in `files` directory. Name of the file will be the name of the card, content of the file will be the description of the card
3. Run `go mod download`
4. Run `go run main.go`