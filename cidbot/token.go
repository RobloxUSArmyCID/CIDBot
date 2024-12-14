package cidbot

import (
	"errors"
	"flag"
	"os"
)

var token = flag.String("token", "", "The bot's authentication token")
var tokenPath = flag.String("token-path", "", "The path to a file containing the bot's authentication token")

func ParseToken() (*string, error) {
	flag.Parse()

	if *token == "" && *tokenPath == "" {
		return nil, errors.New("token and token_path not provided (pick one)")
	}

	if *token != "" && *tokenPath != "" {
		return nil, errors.New("both token and token_path were provided")
	}

	if *token != "" {
		return token, nil
	}

	if *tokenPath != "" {
		file, err := os.ReadFile(*tokenPath)

		if err != nil {
			return nil, err
		}

		token := string(file)

		return &token, nil
	}

	panic("unreachable")
}
