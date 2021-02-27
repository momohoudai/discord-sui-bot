package main

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"math/rand"
	"strings"
	"time"
)

var (
	ResEat       Response_Eat
	ResIfArr     [4]*Response_If
	ResRandIfArr [4]*Response_Rand_If
)

func InitResponses() {
	ResIfArr[0] = &Response_If{
		Reply:         MsgBirthdayKaru,
		EitherContain: MsgBirthdayKaruKeywords,
		MustContain:   MsgBirthdayResponseKeywords,
	}
	ResIfArr[1] = &Response_If{
		Reply:         MsgBirthdayRudo,
		EitherContain: MsgBirthdayRudoKeywords,
		MustContain:   MsgBirthdayResponseKeywords,
	}
	ResIfArr[2] = &Response_If{
		Reply:         MsgBirthdayRako,
		EitherContain: MsgBirthdayRakoKeywords,
		MustContain:   MsgBirthdayResponseKeywords,
	}
	ResIfArr[3] = &Response_If{
		Reply:         MsgBirthdaySui,
		EitherContain: MsgBirthdaySuiKeywords,
		MustContain:   MsgBirthdayResponseKeywords,
	}

	ResRandIfArr[0] = &Response_Rand_If{
		Replies:  MsgHiResponses,
		Keywords: MsgHiKeywords,
	}
	ResRandIfArr[1] = &Response_Rand_If{
		Replies:  MsgThanksResponses,
		Keywords: MsgThanksKeywords,
	}
	ResRandIfArr[2] = &Response_Rand_If{
		Replies:  MsgByeResponses,
		Keywords: MsgByeKeywords,
	}
	ResRandIfArr[3] = &Response_Rand_If{
		Replies:  MsgSorryResponses,
		Keywords: MsgSorryKeywords,
	}

}

func ProcResponses(Ses *discordgo.Session, Msg *discordgo.MessageCreate) {
	if strings.Contains(Msg.Content, "sui") {
		// Don't do it dynamically; this should be way faster.
		if ResEat.Exec(Ses, Msg) {
			return
		}

		for _, element := range ResRandIfArr {
			if element.Exec(Ses, Msg) {
				return
			}
		}

		for _, element := range ResIfArr {
			if element.Exec(Ses, Msg) {
				return
			}
		}
	}

}

// Responses
type Response_Eat struct {
	EitherContain []string
}

func (r *Response_Eat) Exec(Ses *discordgo.Session, Msg *discordgo.MessageCreate) bool {
	Content := strings.ToLower(Msg.Content)
	rand.Seed(time.Now().Unix()) // initialize global pseudo random generator

	for _, Element := range r.EitherContain {
		if strings.Contains(Content, Element) {
			randNum := rand.Intn(len(MsgEatResponses))
			chosenResponse := MsgEatResponses[randNum]
			str := fmt.Sprintf("You should have **%s**", chosenResponse)
			Ses.ChannelMessageSend(Msg.ChannelID, str)
			return true
		}
	}
	return false
}

type Response_If struct {
	Reply         string
	EitherContain []string
	MustContain   []string
}

func (r *Response_If) Exec(Ses *discordgo.Session, Msg *discordgo.MessageCreate) bool {
	Content := strings.ToLower(Msg.Content)

	for _, E := range r.MustContain {
		if !strings.Contains(Content, E) {
			return false
		}
	}

	for _, E := range r.EitherContain {
		if strings.Contains(Content, E) {
			Ses.ChannelMessageSend(Msg.ChannelID, r.Reply)
			return true
		}
	}
	return false
}

type Response_Rand_If struct {
	Replies  []string
	Keywords []string
}

func (r *Response_Rand_If) Exec(Ses *discordgo.Session, Msg *discordgo.MessageCreate) bool {
	rand.Seed(time.Now().Unix()) // initialize global pseudo random generator
	content := strings.ToLower(Msg.Content)
	for _, element := range r.Keywords {
		if strings.Contains(content, element) {
			randNum := rand.Intn(len(r.Replies))
			reply := r.Replies[randNum]
			Ses.ChannelMessageSend(Msg.ChannelID, reply)
			return true
		}
	}
	return false
}
