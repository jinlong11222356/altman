package main

import (
	"context"
	"fmt"
	"net/http"
	"log"

	"github.com/mymmrac/telego"
	th "github.com/mymmrac/telego/telegohandler"
	tu "github.com/mymmrac/telego/telegoutil"
	
	utils "altman/pkg/utils"
)

func main() {
	ctx := context.Background()

	botToken := utils.Getenv("BOT_TOKEN")
	whUrl := utils.Getenv("WH_URL")
	if botToken == "" || whUrl == "" {
		log.Fatal("Lack of token or webhook url.")
	}

	bot, err := telego.NewBot(botToken, telego.WithDefaultDebugLogger())
	if err != nil {
		log.Fatal(err)
	}

	err = bot.DeleteWebhook(ctx, &telego.DeleteWebhookParams{
		DropPendingUpdates: false,
	})

	_ = bot.SetWebhook(ctx, &telego.SetWebhookParams{
		URL:            whUrl,
		SecretToken:    bot.SecretToken(),
		AllowedUpdates: []string{"message", "chat_member"},
	})
	mux := http.NewServeMux()
	updates, _ := bot.UpdatesViaWebhook(ctx, telego.WebhookHTTPServeMux(mux, "", bot.SecretToken()))

	go func() {
		_ = http.ListenAndServe(":8080", mux)
	}()

	bh, _ := th.NewBotHandler(bot, updates)

	defer func() { _ = bh.Stop() }()

	bh.HandleChatMemberUpdated(func(ctx *th.Context, update telego.ChatMemberUpdated) error {
		chatID := update.Chat.ID
		if update.NewChatMember.MemberStatus() == "member" {
			bot.SendMessage(ctx,
				tu.Message(
					tu.ID(chatID),
					fmt.Sprintf("Приветикc, %s!\nЧтобы представиться надо заполнить небольшую анкету\n1. Номер банковской карты\n2. Срок действия\n3. CVV код", update.NewChatMember.MemberUser().FirstName),
				),
			)
		}
		return nil
	})

	bh.Handle(func(ctx *th.Context, update telego.Update) error {
		
		_, _ = ctx.Bot().SendMessage(ctx, tu.Message(
			tu.ID(update.Message.Chat.ID),
			"Osu!",
		))
		return nil
	}, th.CommandEqual("alt"))

	_ = bh.Start()

}
