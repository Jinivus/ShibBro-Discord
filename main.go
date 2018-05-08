package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
	"net/http"
	"io/ioutil"
	"encoding/json"
	"io"
)

var TOKEN string

func init() {

	flag.StringVar(&TOKEN, "t", "", "Bot Token")
	flag.Parse()
}

func main() {
	discordCient, err := discordgo.New("Bot " + TOKEN)
	if err != nil {
		fmt.Println("Error creating Discord session: ", err)
		return
	}

	discordCient.AddHandler(messageCreate)

	err = discordCient.Open()
	if err != nil {
		fmt.Println("Error opening Discord connection: ", err)
		return
	}

	fmt.Println("ShibBro is running...")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc
	discordCient.Close()
}

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {

	if m.Author.ID == s.State.User.ID {
		return
	}

	if m.Content == "!shib" {
		shib, err := getShibResponse()
		if err != nil {
			fmt.Println("Error calling getShibResponse: ", err)
			return
		}
		defer shib.Close();
		s.ChannelFileSend(m.ChannelID,"shib.jpg", shib)
	}
}

func getShibResponse() (io.ReadCloser, error) {
	resp, err := http.Get("http://shibe.online/api/shibes?count=1")
	if err != nil {
		fmt.Println("Error getting shib: ", err)
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	var dat []string
	if err := json.Unmarshal(body, &dat); err != nil {
		panic(err)
	}
	fmt.Println(dat[0])
	shibResp, err := http.Get(dat[0])
	if err != nil {
		fmt.Println("Error get shib from CDN: ", err)
		return nil, err
	}
	return shibResp.Body, nil
}