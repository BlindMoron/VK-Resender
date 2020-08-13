package main

import (
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/Vorkytaka/easyvk-go/easyvk"
	vkapi "github.com/himidori/golang-vk-api"
	"github.com/urShadow/go-vk-api"
)

//Accounts information
type Accounts struct {
	phone    []string
	token    []string
	password []string
	vkid     []int64
}

func main() {
	//Accounts info init
	api := vk.New("ru")
	var accounts Accounts
	accounts.phone = []string{"phone1", "phone2"}
	accounts.token = []string{"token1", "token2"} //managers tokens, get here https://vkhost.github.io
	accounts.password = []string{"Pass1", "Pass2"}
	accounts.vkid = []int64{1, 2}
	var groupAPIKey = "123"
	//Api init
	err := api.Init(groupAPIKey) //group token
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("Start OK.")
	}
	go autoLike(accounts)
	api.OnNewMessage(func(msg *vk.LPMessage) {
		go messageMonitoring(api, msg, accounts)
	})
	api.RunLongPoll()
}

//Resend message from group to manager
func messageMonitoring(api *vk.VK, msg *vk.LPMessage, accounts Accounts) {
	if msg.Flags&vk.FlagMessageOutBox == 0 {
		if msg.Text == "/end" { //Program close from vk
			os.Exit(2)
		}
		for i := 0; i < len(accounts.vkid); i++ {
			if msg.FromID != accounts.vkid[i] { //On new message from client resend it to managers
				fmt.Println("Message from: " + strconv.FormatInt(msg.FromID, 10))
				api.Messages.Send(vk.RequestParams{
					"peer_id":          strconv.FormatInt(accounts.vkid[i], 10),
					"message":          strconv.FormatInt(msg.FromID, 10),
					"forward_messages": strconv.FormatInt(msg.ID, 10),
				})
			} else { //Resende message from manager to client
				h := strings.Split(msg.Text, " ")
				runes := []rune(msg.Text)
				fmt.Println("Message to: " + h[0])
				api.Messages.Send(vk.RequestParams{
					"peer_id": h[0],
					"message": string(runes[len(h[0]):]),
				})
			}
		}
	}
}

func autoLike(accounts Accounts) {
	for {
		for i := 0; i < len(accounts.phone); i++ {
			vk := easyvk.WithToken(accounts.token[i])
			fmt.Println("Auth vk ", vk.AccessToken)
			client, err := vkapi.NewVKClient(vkapi.DeviceIPhone, accounts.phone[i], accounts.password[i])
			wall, err := client.WallGet(-139912528, 2, nil) //get wall posts
			if err != nil {
				fmt.Println("Cant get wall posts: ", err)
			}
			for j := 0; j < len(wall.Posts); j++ {
				_, err = vk.Likes.Add(easyvk.PostLikeType, wall.Posts[j].OwnerID, uint(wall.Posts[j].ID)) //Add like to last post
				if err != nil {
					fmt.Println("Cant place like: ", err)
				} else {
					fmt.Print("Like placed by: ", accounts.phone[i])
					fmt.Print(" at: ", time.Now())
					fmt.Println(" placed on: ", wall.Posts[j].ID)
				}
			}
			n := time.Duration(rand.Int63n(15))
			fmt.Println(n)
			time.Sleep(time.Minute * n) //wait random time to place another like from manager
		}
		time.Sleep(time.Hour * 6) //wait for 6 hours
	}
}
