package main

import (
	"context"
	"fmt"
	"net/http"
	"log"
	"os"
    "strconv"

	"github.com/mymmrac/telego"
	th "github.com/mymmrac/telego/telegohandler"
	tu "github.com/mymmrac/telego/telegoutil"
	
    "github.com/jackc/pgx/v5/pgxpool"

	utils "altman/pkg/utils"

	"altman/internal/db"
)

func main() {
	ctx := context.Background()

	botToken := utils.Getenv("BOT_TOKEN")
	whUrl := utils.Getenv("WH_URL")
	dbUrl := utils.Getenv("DB_URL")
    initPrompt := utils.Getenv("INIT_PROMPT")


	groupID, _ := strconv.ParseInt(utils.Getenv("GROUP_ID"),10,64)
	moderGroupID, err := strconv.ParseInt(utils.Getenv("MODER_GROUP_ID"),10,64)

	pool, err := pgxpool.New(ctx, dbUrl)
    if err != nil {
        fmt.Fprintf(os.Stderr, "Unable to create connection pool: %v\n", err)
        os.Exit(1)
    }
    defer pool.Close()

    q := db.New(pool)

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
    
	//  NEW MEMBER
	// bh.HandleChatMemberUpdated(func(ctx *th.Context, update telego.ChatMemberUpdated) error {
	// 	chatID := update.Chat.ID
	// 	if update.NewChatMember.MemberStatus() == "member" {
	// 		bot.SendMessage(ctx,
	// 			tu.Message(
	// 				tu.ID(chatID),
	// 				fmt.Sprintf("", update.NewChatMember.MemberUser().FirstName),
	// 			),
	// 		)
	// 	}
	// 	return nil
	// })

	// PING COMMAND
	bh.Handle(func(ctx *th.Context, update telego.Update) error {
		
		_, _ = ctx.Bot().SendMessage(ctx, tu.Message(
			tu.ID(update.Message.Chat.ID),
			"Osu!",
		))
		return nil

	}, th.CommandEqual("alt"))

	//  START COMMAND
	 bh.Handle(func(ctx *th.Context, update telego.Update) error {		
	 	_, _ = ctx.Bot().SendMessage(ctx, tu.Message(
	 		tu.ID(update.Message.Chat.ID),
	 		initPrompt,
	 	))
	 	return nil

	 }, th.CommandEqual("start"))

	//  SUGGESTION
	bh.HandleMessage(func(ctx *th.Context, message telego.Message) error {		
	    userChatID := message.Chat.ID
		
		if userChatID != moderGroupID && userChatID != groupID{ // request
				userMessageID := message.MessageID
				botMessage, _ := ctx.Bot().SendMessage(ctx, tu.Message(
				tu.ID(moderGroupID),
				message.Text,
			))
			err = q.CreateSuggestion(ctx, db.CreateSuggestionParams{
				UserChatID    : userChatID,
				UserMessageID : int64(userMessageID),
				BotMessageID  : int64(botMessage.MessageID),
			})
		}
		if err != nil {
			log.Println(err)
		}

		if userChatID == moderGroupID && message.ReplyToMessage != nil { // response
			requestInfo, err := q.GetSender(ctx, int64(message.ReplyToMessage.MessageID))
			if err != nil {
				log.Println("Get sender db request failed")
				return nil
			}
			_, _ = ctx.Bot().SendMessage(ctx, tu.Message(tu.ID(requestInfo.UserChatID), message.Text).WithReplyParameters(&telego.ReplyParameters{ MessageID: int(requestInfo.UserMessageID)}))
		}
            
		return nil
	})

	_ = bh.Start()

}
