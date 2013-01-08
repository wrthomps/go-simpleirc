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
	Author	string
	Channel	string
	Content	string
	Host	string
	Timestamp	time.Time
	Type	string
}

func (msg *Message) PrettyPrint() {
	if msg == nil {
		return
	}
	log.Println("Author:", msg.Author)
	log.Println("Host:", msg.Host)
	log.Println("Type:", msg.Type)
	log.Println("Channel:", msg.Channel)
	log.Println("Content:", msg.Content)
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
	User:       "VoteBot",
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
	fmt.Println(fields)

	// Respond to ping/pong event to avoid ping timeout disconnects
	if fields[0] == "PING" {
		response := strings.Replace(line, "PING", "PONG", 1)
		fmt.Fprintf(bot.Connection, response)
		fmt.Println("Responded to PING event:", response)
		return nil
	}

	var msg Message

	// The author of the message is the first field up to a "!" character, except
	// for the leading colon
	msg.Author = (strings.Split(fields[0], "!"))[0][1:]
	// The author's host is everything else in the first field
	msg.Host = (strings.Split(fields[0], "!"))[1]
	// The type of message is the next whitespace-delimited field
	msg.Type = fields[1]
	// The channel is the next field
	msg.Channel = fields[2]
	// The content is everything else, except the leading colon for the last field
	msg.Content = string((strings.Join(fields[3:], " "))[1:])
	msg.Timestamp = time.Now()

	return &msg
}

// Checks a message to see if the bot should perform some command, and if so,
// performs it
func doCommand(bot *IRCBot, msg *Message) {
	// Right now the only command is !test
	if strings.Fields(msg.Content)[0] == "!test" {
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

