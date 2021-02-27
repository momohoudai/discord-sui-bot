package main

import (
	"github.com/bwmarrin/discordgo"
	"github.com/robfig/cron/v3"

	"fmt"
	"io/ioutil"
	"math/rand"
	"regexp"
	"strings"
	"time"
)

var (
	GlobalBotId       string
	GlobalArgSplitter *regexp.Regexp
)

func MessageHandler(Ses *discordgo.Session, Msg *discordgo.MessageCreate) {
	User := Msg.Author
	if User.ID == GlobalBotId || User.Bot {
		return
	}

	// For message parsing
	UpdateWordCount(Msg)

	// handle prefix
	Content := strings.ToLower(Msg.Content)
	if strings.HasPrefix(Content, "sui") {
		go func() {
			defer Kalm(Ses, Msg, "ProcCommands")
			ProcCommands(Ses, Msg)
		}()
	} else {
		go func() {
			defer Kalm(Ses, Msg, "ProcResponses")
			ProcResponses(Ses, Msg)
		}()
	}
}

func ReadyHandler(Ses *discordgo.Session, Ready *discordgo.Ready) {
	Err := Ses.UpdateListeningStatus("'sui help'")
	if Err != nil {
		fmt.Println("Error attempting to set my status")
	}
	servers := Ses.State.Guilds
	fmt.Printf("SuiBot has started on %d servers\n", len(servers))
}

func Panik(Format string, a ...interface{}) {
	panic(fmt.Sprintf(Format, a...))
}

func Kalm(Ses *discordgo.Session, Msg *discordgo.MessageCreate, Name string) {
	if R := recover(); R != nil {
		fmt.Printf("[%s] Recovered: %v\n", Name, R)
		Ses.ChannelMessageSend(Msg.ChannelID, MsgGenericFail)
	}
}

func main() {
	rand.Seed(time.Now().UnixNano())
	Token, ReadFileErr := ioutil.ReadFile("TOKEN")
	if ReadFileErr != nil {
		Panik("Cannot read or find TOKEN file\n")
	}

	InitWordCount()
	InitCommands()
	InitResponses()

	// Scheduler
	Cron := cron.New()
	Cron.AddFunc("@every 1h", SchedCleanWordCountDatabase)

	Cron.Start()
	defer Cron.Stop()

	// Configure Discord
	Discord, Err := discordgo.New("Bot " + string(Token))
	if Err != nil {
		Panik("Cannot initialize discord: %s\n", Err.Error())
	}
	User, Err := Discord.User("@me")
	if Err != nil {
		Panik("Error retrieving account: %s\n", Err.Error())
	}
	GlobalBotId = User.ID
	Discord.AddHandler(MessageHandler)
	Discord.AddHandler(ReadyHandler)

	Err = Discord.Open()
	if Err != nil {
		Panik("Error retrieving account: %s\n", Err.Error())
	}
	defer Discord.Close()

	<-make(chan struct{})

}
