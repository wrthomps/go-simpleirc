package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"net/textproto"
	"strings"
	"time"
)

type IRCBot struct {
	Channel     string
	Connection  net.Conn
	Nick        string
	Pass        string
	Port        string
	Server      string
	User        string
}

type Message struct {
	Command    string
	Params     []string
	Prefix     string
	Timestamp  time.Time
	Trail      string
}

func (msg *Message) PrettyPrint() {
	if msg == nil {
		return
	}
	log.Println(msg.Timestamp.Format("[15:04]"))
	log.Println("Prefix:\t", msg.Prefix)
	log.Println("Command:\t", msg.Command)
	log.Println("Params:\t", msg.Params)
	log.Println("Trail:\t", msg.Trail)
}

// Create a new bot with the given server info
func createBot() *IRCBot {
	return &IRCBot {
	Server:     "irc.rizon.net",
	Port:       "6667",
	Nick:       "IdolBot",
	Channel:    "#ircbottestrelam",
	Pass:       "",
	Connection: nil,
	User:       "IdolBot",
	}
}

// Connects the bot to the server
func (bot *IRCBot) ServerConnect() (connection net.Conn, err error) {
	connection, err = net.Dial("tcp", bot.Server + ":" + bot.Port)
	if err != nil {
	  log.Fatal("Unable to connect to the specified server", err)
	}
	bot.Connection = connection
	log.Printf("Successfully connected to %s ($s)\n", bot.Server, bot.Connection.RemoteAddr())
	return bot.Connection, nil
}

// Parses a line from IRC into a Message struct
func parseMessage(bot *IRCBot, line string) *Message {
	fmt.Println("---")
	fmt.Println(line)
	fmt.Println("---")

	fields := strings.Fields(line)
	var msg Message

	// Used to isolate the command and parameters
	var prefixEnd int
	var trailBegin int

	// Collect the prefix if it exists. The prefix exists iff the first character of the
	// message is a colon, and the prefix is everything following the colon up to the first
	// space. Thus the prefix can only exist in the first field.
	prefix := fields[0]
	if prefix[0] == ':' {
		msg.Prefix = prefix[1:]
		prefixEnd = len(fields[0])
	} else {
		msg.Prefix = "" 
		prefixEnd = -1
	}

	// Collect the trail if it exists. It is everything after the occurrance of " :", note
	// the space.
	trailBegin = strings.Index(line, " :")
	
	// Only set the trail if there is a valid one. If " :" is the last part of the message
	// then there is no character afterward.
	if trailBegin >= 0 && trailBegin+2 < len(line) {
		msg.Trail = line[trailBegin+2:]
	} else {
		msg.Trail = ""
	}

	// Collect the command and parameters. They're everything between the prefix and trail.
	cmdAndParams := strings.Split(line[prefixEnd+1:trailBegin], " ")
	msg.Command = cmdAndParams[0]

	if len(cmdAndParams) > 1 {
		msg.Params = cmdAndParams[1:]
	} else {
		msg.Params = nil
	}

	return &msg
}

// Checks a message to see if the bot should perform some command, and if so,
// performs it
func doCommand(bot *IRCBot, msg *Message) {
	// Right now the only command is !test
	if strings.Fields(msg.Trail)[0] == "!test" {
		s := "PRIVMSG"
		s += " " + bot.Channel
		s += " :"
		s += "The quick brown fox jumps over the lazy dog."
		s += "\r\n"
		fmt.Println(s)
		fmt.Fprintf(bot.Connection, s)
	}
}

func main() {
	bot := createBot()
	connection, _ := bot.ServerConnect()
	fmt.Fprintf(connection, "USER %s 8 * :%s\n", bot.Nick, bot.Nick)
	fmt.Fprintf(connection, "NICK %s\n", bot.Nick)
	fmt.Fprintf(connection, "JOIN %s\n", bot.Channel)
	defer connection.Close()

	reader := bufio.NewReader(connection)
	respReq := textproto.NewReader(reader)

	// Wait for the initial overhead messages we don't need to respond to
	waitChannel := time.After(30 * time.Second)
	respond := false
	for {
		line, err := respReq.ReadLine()
		if err != nil {
	  		break
		}

		select {
		case _ = <-waitChannel:
			respond = true
		default:
		}

		if respond {
			msg := parseMessage(bot, line)
	  		msg.PrettyPrint()
			doCommand(bot, msg)
		} else {
			fmt.Println(line)
		}
	}
}

