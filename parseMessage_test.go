package main

import (
	"testing"
)

func TestParseMessageBasic(t *testing.T) {
	testMessage := ":username@localhost PRIVMSG #test :Hello, world!"
	bot := createBot()
	msg := parseMessage(bot, testMessage)

	if msg.Prefix != "username@localhost" || msg.Command != "PRIVMSG" || msg.Params[0] != "#test" || len(msg.Params) != 1 || msg.Trail != "Hello, world!" {
		t.Fail()
	}
	
	respondToRequest(bot, msg)
}

func TestParseMessagePing(t *testing.T) {
	testMessage := "PING :irc.notreal.com"
	bot := createBot()
	msg := parseMessage(bot, testMessage)

	if msg.Prefix != "" || msg.Command != "PING" || msg.Params != nil || msg.Trail != "irc.notreal.com" {
		t.Fail()
	}

	respondToRequest(bot, msg)

	testMessage = ":server.internet.com PING :irc.notreal.com"
	msg = parseMessage(bot, testMessage)

	if msg.Prefix != "server.internet.com" || msg.Command != "PING" || msg.Params != nil || msg.Trail != "irc.notreal.com" {
		t.Fail()
	}

	respondToRequest(bot, msg)
}

func TestParseMessageVersion(t *testing.T) {
	testMessage := ":server.internet.com VERSION"
	bot := createBot()
	msg := parseMessage(bot, testMessage)

	respondToRequest(bot, msg)
}

func TestNotRecognized(t *testing.T) {
	testMessage := ":server.internet.com NYI params params :Hello, world!"
	bot := createBot()
	msg := parseMessage(bot, testMessage)

	respondToRequest(bot, msg)
}
