package main

import (
	"math/rand"
	"strings"
)

const MsgVersion = "I'm SuiBot v6.0.0, written in Golang!"

const MsgBirthdayKaru = "Ah...erm...K-Karu's birthday is on **April 12th**...I remember okay!? Are you testing me, mouu >//<"
const MsgBirthdayRudo = "Oh, Rudo's birthday is on **November 30th**. You should've looked at his face when I brought back a cake for him ^^;"
const MsgBirthdayRako = "Fumu, Rako's birthday should be on **October 31st**. It's the only occasion Rudo remembers and he always asks me for advice around that period ww"
const MsgBirthdaySui = "My birthday is on **June 24th**. Erm, I don't really watch soccer but my friends told me that it's the same birthday as Nakamura Shunsuke...?"

// Can't be const...maybe one day golang will allow it
var MsgBirthdayKaruKeywords = []string{"karu", "karu's"}
var MsgBirthdayRudoKeywords = []string{"rudo's", "rudo"}
var MsgBirthdayRakoKeywords = []string{"rako's", "rako"}
var MsgBirthdaySuiKeywords = []string{"your", "sui's"}

var MsgBirthdayResponseKeywords = []string{"birthday"}
var msgEatResponseKeywords = []string{"eat", "food", "hungry"}

// Hi
var MsgHiResponses = []string{
	"Oh hey!",
	"Hi! How's your day?",
	"Hello!",
}
var MsgHiKeywords = []string{"hello", "hi", "yo"}
var MsgSorryResponses = []string{"I'm sorry T_T"}
var MsgSorryKeywords = []string{"wrong", "no", "bad"}

var MsgThanksResponses = []string{
	"You're welcome!",
	"Ah, it's not a bother ^^",
	"Only here to help!　（｀・ω・´）",
}
var MsgThanksKeywords = []string{"thank you", "thanks", "thx"}

var MsgByeKeywords = []string{"bye", "goodbye", "cya", "gnite", "good night"}
var MsgByeResponses = []string{
	"Oh, see you!",
	"Come again next time ^^",
	"Oh, you are going? Bye!",
	"Feel free to drop by again.",
}

var MsgEatResponses = []string{
	"something with curry",
	"something soupy",
	"something with rice",
	"something with bread",
	"something with noodles",

	// specific
	"udon",
	"soba",
	"sushi",
	"ramen",
	"pasta",
	"pizza",
	"burger",
	"wrap",
	"sandwich",

	// meat
	"something with beef",
	"something with chicken",
	"something with pork",
	"something with fish",
	"something with meat",
	"something with vegetables",
	"something with tofu",

	//cultural
	"Indian",
	"Western",
	"Japanese",
	"Korean",
	"Chinese",
	"Italian",
	"Mexican",
	"Turkish",
	"Local",
}

const MsgPollTitle = "Poll Created By: %s"
const MsgPollFooter = "React to vote!"

const MsgHelpQuery = "Did you do it correctly? ```%s```"
const MsgPollHelp = "poll: Creates a poll!\n\t> sui poll <description_of_the_poll>\n\t> sui poll \"Is this moe?\"\n\t> sui poll \"Do you love me?\""
const MsgRollHelp = "roll: Rolls a die\n\t> sui roll 2d6"
const MsgCountStartHelp = "count-start: Starts counting a word or phrase from people in the server\n\t> sui count-start <string_you_want_to_count>\n\t> sui count-start moe\n\t> sui count-start \"so moe?!\""
const MsgCountStopHelp = "count-stop: Stops counting for count-start\n\t> sui count-stop"

var MsgHelp = strings.Join([]string{"```", MsgRollHelp, MsgPollHelp, MsgCountStartHelp, MsgCountStopHelp, "```"}, "\n")

const MsgGenericFail = "Sorry, something went wrong...contact Momo? (´・ω・`)"

var MsgRollReplies = []string{
	"K-Korokoroko~! 👀\n",
	"*Takes out :abacus: :face_with_monocle:* \n",
	"Let me open my caculator app ^^;\n",
	"Wait, I hadn't had my tea but er...I'll try! ><\n",
	"Er, that's too hard. Let me ask Karu!\n",
	"*frowns, scratches head* fumu fumu...\n",
}

const MsgRollConclude = "I got it! It's: **%d**!"
const MsgCountStartAlready = "I'm already counting '%s' for this channel! Stop counting with ```%s```"
const MsgCountStart = "Counting '%s' this channel! Remember to stop with ```%s```"
const MsgCountStop = "I counted a total of **%d** '%s'! Have a nice day! :relaxed:"
const MsgCountStopAlready = "You didn't ask me to start counting! Use the command below! :sweat_smile: ```%s```"

func RandomMsg(MsgArr []string) string {
	RandIndex := rand.Intn(len(MsgArr))
	return MsgArr[RandIndex]
}
