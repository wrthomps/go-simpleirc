package main

import (
	"log"
)

var CODE_LIST = map[string] func(bot *IRCBot, msg *Message) *Message {
	"PING" : pingPong,
	"PRIVMSG" : parseCommand,
	"VERSION" : version,
}

// Send a pre-constructed message
func sendMessage(bot *IRCBot, msg *Message) {
	rawMsg := ""

	// Convert the prefix back to its raw form if it exists
	if msg.Prefix != "" {
		rawMsg += ":"
		rawMsg += msg.Prefix
		rawMsg += " "
	}

	// Add the command onto the message string
	rawMsg += msg.Command
	rawMsg += " "

	for i := 0; i < len(msg.Params); i++ {
		rawMsg += msg.Params[i]
		rawMsg += " "
	}

	if msg.Trail != "" {
		rawMsg += ":"
		rawMsg += msg.Trail
	}

	log.Println(rawMsg)
	// fmt.Fprintf(bot.Connection, rawMsg)
}

// Returns a pong response
// Example:
// PONG :irc.example.com
func pingPong(bot *IRCBot, msg *Message) *Message {
	var resp Message

	resp.Command = "PONG"
	resp.Trail = msg.Trail

	// Send the message now that we've created it
	defer sendMessage(bot, &resp)
	return &resp
}

// Echos the struct representation of the message
// TODO: Have it look through a list of implemented commands and execute
//	one that may be requested
func parseCommand(bot *IRCBot, msg *Message) *Message {
	msg.PrettyPrint()
	return msg
}

func version(bot *IRCBot, msg *Message) *Message {
	var resp Message

	resp.Command = "VERSION"
	resp.Params = make([]string, 1)
	resp.Params[0] = "go-simpleIRC 0.0.2 (x86-64)"

	defer sendMessage(bot, &resp)
	return &resp
}
