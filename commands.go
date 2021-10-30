package main

import (
	"github.com/bwmarrin/discordgo"

	"fmt"
	"regexp"
	"strings"
	"time"
)

var (
	GlobalCommandArgSplitter *regexp.Regexp
)

const DiscordMessageMaxChars int = 2000

func InitCommands() {
	GlobalCommandArgSplitter = regexp.MustCompile(`(?i)(?:[^\s"]+\b|:|(")[^"]*("))+|[=!&|~+\-\*\/\%]`)
	if GlobalCommandArgSplitter == nil {
		Panik("commandArgSplitter failed to compile")
	}
}

func ProcCommands(Ses *discordgo.Session, Msg *discordgo.MessageCreate) {
	Args := GlobalCommandArgSplitter.FindAllString(Msg.Content, -1)
	if Args != nil && len(Args) >= 2 {
		CommandStr := Args[1]
		Args = Args[2:] // get array from 2 to n

		// Use this instead of interfaces
		// Clearer, more concise, easier to debug
		switch CommandStr {
		case "version":
			Ses.ChannelMessageSend(Msg.ChannelID, MsgVersion)
		case "help":
			Ses.ChannelMessageSend(Msg.ChannelID, MsgHelp)
		case "poll":
			CmdPoll(Ses, Msg, Args)
		case "roll":
			fallthrough
		case "calc":
			CmdRoll(Ses, Msg, Args)
		case "count-start":
			CmdCountStart(Ses, Msg, Args)
		case "count-stop":
			CmdCountStop(Ses, Msg, Args)
		default:
			ProcResponses(Ses, Msg)
		}

	} else {
		ProcResponses(Ses, Msg)
	}
}

func WrapCode(Str string) string {
	return "```" + Str + "```"
}

func CmdCountStart(Ses *discordgo.Session, Msg *discordgo.MessageCreate, Args []string) {
	if len(Args) == 0 {
		Reply := fmt.Sprintf(MsgHelpQuery, MsgCountStartHelp)
		Ses.ChannelMessageSend(Msg.ChannelID, Reply)
		return
	}

	// Check if entry exists in word count
	Entry, Success := GlobalWordCountDatabase[Msg.ChannelID]
	if Success {
		// The entry exists
		Reply := fmt.Sprintf(MsgCountStartAlready, Entry.Word, MsgCountStopHelp)
		Ses.ChannelMessageSend(Msg.ChannelID, Reply)
	} else {
		Word := strings.Join(Args, " ")
		Regex := regexp.MustCompile("(?i)" + Word)
		CurrentTime := time.Now().Unix()

		GlobalWordCountDatabase[Msg.ChannelID] = &Word_Count_Entry{
			Word:    Word,
			Count:   0,
			Regex:   Regex,
			Created: CurrentTime,
		}

		Reply := fmt.Sprintf(MsgCountStart, Word, MsgCountStopHelp)
		Ses.ChannelMessageSend(Msg.ChannelID, Reply)
	}
}

func CmdCountStop(Ses *discordgo.Session, Msg *discordgo.MessageCreate, Args []string) {
	if len(Args) > 0 {
		Reply := fmt.Sprintf(MsgHelpQuery, MsgCountStopHelp)
		Ses.ChannelMessageSend(Msg.ChannelID, Reply)
		return
	}

	Entry, Success := GlobalWordCountDatabase[Msg.ChannelID]
	if Success {
		delete(GlobalWordCountDatabase, Msg.ChannelID)
		Reply := fmt.Sprintf(MsgCountStop, Entry.Count, Entry.Word)
		Ses.ChannelMessageSend(Msg.ChannelID, Reply)
	} else {
		Reply := fmt.Sprintf(MsgCountStopAlready, MsgCountStopHelp)
		Ses.ChannelMessageSend(Msg.ChannelID, Reply)
	}

}

func CmdRoll(Ses *discordgo.Session, Msg *discordgo.MessageCreate, Args []string) {
	if len(Args) == 0 {
		Reply := fmt.Sprintf(MsgHelpQuery, MsgRollHelp)
		Ses.ChannelMessageSend(Msg.ChannelID, Reply)
		return
	}

	Infix := strings.Join(Args, "")

	fmt.Printf("Infix: %s\n", Infix)
	Postfix, InfixToPostfixErr := InfixToPostfix(Infix)
	if InfixToPostfixErr != nil {
		fmt.Println(InfixToPostfixErr)
		Reply := fmt.Sprintf(MsgHelpQuery, MsgRollHelp)
		Ses.ChannelMessageSend(Msg.ChannelID, Reply)
	}

	EvalResult, EvalPostfixErr := EvaluatePostfix(Postfix)
	if EvalPostfixErr != nil {
		fmt.Println(EvalPostfixErr)
		Reply := fmt.Sprintf(MsgHelpQuery, MsgRollHelp)
		Ses.ChannelMessageSend(Msg.ChannelID, Reply)
	}

	// Construct the Reply string with the format:
	// "<random messsage>
	// ```<dice info>```
	// <conclusion> + sum
	// If it's too long, omit <dice info>
	//
	// The code below could definitely be optimized for speed
	// but I am doing it #DLP #damnlazyprogramming
	//
	RandMsg := RandomMsg(MsgRollReplies)
	ConcludeStr := fmt.Sprintf(MsgRollConclude, EvalResult.Sum)
	InfoStr := "```" + EvalResult.DiceInfo + "```"

	TotalReplyLength := len(RandMsg) + len(InfoStr) + len(ConcludeStr)
	if TotalReplyLength > DiscordMessageMaxChars {
		Reply := RandMsg + ConcludeStr
		Ses.ChannelMessageSend(Msg.ChannelID, Reply)
	} else {
		Reply := RandMsg + InfoStr + ConcludeStr
		Ses.ChannelMessageSend(Msg.ChannelID, Reply)
	}
}

func CmdPoll(Ses *discordgo.Session, Msg *discordgo.MessageCreate, Args []string) {
	if len(Args) == 0 {
		Reply := fmt.Sprintf(MsgHelpQuery, MsgPollHelp)
		Ses.ChannelMessageSend(Msg.ChannelID, Reply)
		return
	}

	// combine all the args into a sentence
	Description := strings.Join(Args, " ")

	Embed := &discordgo.MessageEmbed{}
	Embed.Title = fmt.Sprintf(MsgPollTitle, Msg.Author.Username)
	Embed.Footer = &discordgo.MessageEmbedFooter{
		Text: MsgPollFooter,
	}
	Embed.Description = Description
	Embed.Color = 0xFFFFFF

	Reply, Err := Ses.ChannelMessageSendEmbed(Msg.ChannelID, Embed)
	if Err != nil {
		return
	}
	Ses.MessageReactionAdd(Reply.ChannelID, Reply.ID, "👍")
	Ses.MessageReactionAdd(Reply.ChannelID, Reply.ID, "👎")
}
