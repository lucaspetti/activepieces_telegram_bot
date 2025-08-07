package cmd

import (
	"activepieces_telegram_bot/bot"
	"fmt"
	"log"
	"net/url"
	"os"
	"strconv"

	"github.com/spf13/cobra"
)

// startCmd represents the start command
var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Starts the telegram bot",
	Long:  `Starts the telegram bot with the necessary config`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("start called")

		authorizedUserID := os.Getenv("AUTHORIZED_USER_ID")
		userID, err := strconv.ParseInt(authorizedUserID, 10, 64)
		if err != nil {
			log.Fatal("Error converting authorized user id: ", err)
		}

		botToken := os.Getenv("TELEGRAM_BOT_TOKEN")
		if botToken == "" {
			log.Fatal("TELEGRAM_BOT_TOKEN env variable is required")
		}

		webhookURL := os.Getenv("WEBHOOK_URL")
		if webhookURL == "" {
			log.Fatal("WEBHOOK_URL env variable is required")
		}

		_, err = url.ParseRequestURI(webhookURL)
		if err != nil {
			log.Fatal("WEBHOOK_URL env is not a valid URL")
		}

		config := bot.NewConfig(
			userID,
			botToken,
			webhookURL,
		)

		bot.Start(*config)
	},
}

func init() {
	rootCmd.AddCommand(startCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// startCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// startCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
