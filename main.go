package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/shomali11/slacker"
	"github.com/spf13/viper"
)

func loadEnv() {
	viper.SetConfigFile(".env")
	viper.SetConfigType("env")
	viper.AutomaticEnv()
	err := viper.ReadInConfig()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

}

func printCommandEvents(analyticsChannel <-chan *slacker.CommandEvent) {
	for event := range analyticsChannel {
		fmt.Println("Command:", event.Command)
		fmt.Println("Timestamp:", event.Timestamp)
		fmt.Println("Parameters:", event.Parameters)
		fmt.Println("Event:", event.Event)
	}
}

func main() {
	loadEnv()
	bot := slacker.NewClient(viper.GetString("SLACK_BOT_TOKEN"), viper.GetString("SLACK_APP_TOKEN"))

	go printCommandEvents(bot.CommandEvents())

	bot.Command("my yob is <year>", &slacker.CommandDefinition{
		Description: "Returns your age",
		Example:     "my yob is 1980",
		Handler: func(botCtx slacker.BotContext, request slacker.Request, response slacker.ResponseWriter) {
			year := request.Param("year")
			if year == "" {
				response.Reply("Please specify your year of birth")
				return
			}
			yob, err := strconv.Atoi(year)
			if err != nil {
				response.Reply("Please specify a valid year")
				return
			}
			age := time.Now().Year() - yob
			response.Reply(fmt.Sprintf("Your age is %d", age))
		},
	})

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	err := bot.Listen(ctx)
	if err != nil {
		log.Fatal(err)
	}

}
