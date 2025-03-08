package auth

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/gotd/td/telegram"
	"github.com/gotd/td/telegram/auth"
	"github.com/gotd/td/tg"
)

func StartAuth(ctx context.Context) (*telegram.Client, *tg.User, error) {
	apiID, err := strconv.Atoi(os.Getenv("TELEGRAM_API_ID"))
	if err != nil {
		log.Fatalf("Invalid TELEGRAM_API_ID: %v", err)
	}

	apiHash := os.Getenv("TELEGRAM_API_HASH")
	if apiHash == "" {
		log.Fatal("TELEGRAM_API_HASH not set")
	}

	// Create a new Telegram client
	client := telegram.NewClient(apiID, apiHash, telegram.Options{})

	err = client.Run(ctx, func(ctx context.Context) error {
		log.Println("Connecting to Telegram...")
		auth.
			flow := auth.NewFlow(
			auth.Terminal(auth.CodeAuthenticator(func(ctx context.Context, sentCode *tg.AuthSentCode) (string, error) {
				log.Printf("Enter the code sent to your Telegram account:")
				var code string
				_, err := fmt.Scan(&code)
				return code, err
			})),
			auth.SendCodeOptions{},
		)

		// Perform authentication
		if err := client.Auth().IfNecessary(ctx, flow); err != nil {
			return err
		}

		self, err := client.Self(ctx)
		if err != nil {
			return err
		}

		log.Printf("Logged in as %s", self.Username)
		return nil
	})

	if err != nil {
		return nil, nil, err
	}

	user, err := client.Self(ctx)
	if err != nil {
		return nil, nil, err
	}

	return client, user, nil
}
