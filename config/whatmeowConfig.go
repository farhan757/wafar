package config

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"waapifar/dto"

	_ "github.com/lib/pq"

	//_ "github.com/mattn/go-sqlite3"
	"github.com/mdp/qrterminal/v3"
	"go.mau.fi/whatsmeow"
	waProto "go.mau.fi/whatsmeow/binary/proto"
	"go.mau.fi/whatsmeow/store/sqlstore"
	"go.mau.fi/whatsmeow/types"
	"go.mau.fi/whatsmeow/types/events"
	waLog "go.mau.fi/whatsmeow/util/log"
	"google.golang.org/protobuf/proto"
)

var client *whatsmeow.Client
var req dto.ReqSend

type KeyBot struct {
	key_string  string
	desc_string string
}

//Init ...
func InitWA() {

	dbLog := waLog.Stdout("Database", "DEBUG", true)
	// Make sure you add appropriate DB connector imports, e.g. github.com/mattn/go-sqlite3 for SQLite
	//container, err := sqlstore.New("sqlite3", "file:wa.db?_foreign_keys=on", dbLog)
	container, err := sqlstore.New("postgres", "postgres://postgres:gns123@localhost:5432/wa_golang2?sslmode=disable", dbLog)
	if err != nil {
		panic(err)
	}
	// If you want multiple sessions, remember their JIDs and use .GetDevice(jid) or .GetAllDevices() instead.
	deviceStore, err := container.GetFirstDevice()
	//deviceStore, err := container.GetDevice()
	if err != nil {
		panic(err)
	}
	clientLog := waLog.Stdout("Client", "DEBUG", true)
	client = whatsmeow.NewClient(deviceStore, clientLog)
	client.AddEventHandler(eventHandler)

	if client.Store.ID == nil {
		// No ID stored, new login
		qrChan, _ := client.GetQRChannel(context.Background())
		err = client.Connect()
		if err != nil {
			panic(err)
		}
		for evt := range qrChan {
			if evt.Event == "code" {
				configqr := qrterminal.Config{
					Level:          qrterminal.L,
					Writer:         os.Stdout,
					HalfBlocks:     true,
					BlackChar:      qrterminal.BLACK_BLACK,
					WhiteBlackChar: qrterminal.WHITE_BLACK,
					WhiteChar:      qrterminal.WHITE_WHITE,
					BlackWhiteChar: qrterminal.BLACK_WHITE,
					QuietZone:      3,
				}
				// Render the QR code here
				// e.g. qrterminal.GenerateHalfBlock(evt.Code, qrterminal.L, os.Stdout)
				// or just manually `echo 2@... | qrencode -t ansiutf8` in a terminal
				fmt.Println("QR code:", evt.Code)
				qrterminal.GenerateWithConfig(evt.Code, configqr)
			} else {
				fmt.Println("Login event:", evt.Event)
			}
		}
	} else {
		// Already logged in, just connect
		err = client.Connect()
		if err != nil {
			panic(err)
		}
	}

	// Listen to Ctrl+C (you can also do something else that prevents the program from exiting)
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c

	client.Disconnect()
}

//GetDB ...
func GetClient() *whatsmeow.Client {
	return client
}

func eventHandler(evt interface{}) {

	switch v := evt.(type) {
	case *events.Message:
		//fmt.Println("Received a message!", v.Message.GetConversation())
		if !v.Info.IsFromMe && !v.Info.IsGroup && v.Info.MediaType == "" {
			fmt.Println("GetConversation : ", v.Message.GetConversation())
			fmt.Println("Sender : ", v.Info.Sender)
			fmt.Println("Sender Number : ", v.Info.Sender.User)
			fmt.Println("IsGroup : ", v.Info.IsGroup)
			fmt.Println("MessageSource : ", v.Info.MessageSource)
			fmt.Println("ID : ", v.Info.ID)
			fmt.Println("PushName : ", v.Info.PushName)
			fmt.Println("BroadcastListOwner : ", v.Info.BroadcastListOwner)
			fmt.Println("Category : ", v.Info.Category)
			fmt.Println("Chat : ", v.Info.Chat)
			fmt.Println("DeviceSentMeta : ", v.Info.DeviceSentMeta)
			fmt.Println("IsFromMe : ", v.Info.IsFromMe)
			fmt.Println("MediaType : ", v.Info.MediaType)
			fmt.Println("Multicast : ", v.Info.Multicast)
			fmt.Println("Info.Chat.Server : ", v.Info.Chat.Server)
			if v.Info.Chat.Server == "g.us" {
				groupInfo, err := client.GetGroupInfo(v.Info.Chat)
				fmt.Println("error GetGroupInfo : ", err)
				fmt.Println("Nama Group : ", groupInfo.GroupName.Name)
			}

			connStr := "postgres://postgres:gns123@localhost:5432/wa_golang2?sslmode=disable"
			db, err := sql.Open("postgres", connStr)

			if err != nil {
				fmt.Println(err)
			}
			res := db.QueryRow("select key_string, desc_string from whatsmeow_key_bot where key_string=$1;", strings.ToLower(v.Message.GetConversation()))

			var keybot KeyBot
			err = res.Scan(&keybot.key_string, &keybot.desc_string)

			if err != nil {
				fmt.Println(err)
			} else {
				req.No = v.Info.Sender.User
				req.Text = keybot.desc_string
				SendRepBot(&req)
			}
			//db.Close()
		}
	}
}

func SendRepBot(req *dto.ReqSend) {

	//send
	recipient, ok := ParseJID(req.No)
	if !ok {
		return
	}

	msg := &waProto.Message{Conversation: proto.String(req.Text)}

	resp, err := GetClient().SendMessage(context.Background(), recipient, "", msg)
	if err != nil {
		fmt.Println("Error sending message: %v", err)
	} else {
		fmt.Println("Message sent (server timestamp: %s)", resp.Timestamp)
	}
}

func ParseJID(arg string) (types.JID, bool) {
	if arg[0] == '+' {
		arg = arg[1:]
	}
	if !strings.ContainsRune(arg, '@') {
		return types.NewJID(arg, types.DefaultUserServer), true
	} else {
		recipient, err := types.ParseJID(arg)
		if err != nil {
			fmt.Println("Invalid JID %s: %v", arg, err)
			return recipient, false
		} else if recipient.User == "" {
			fmt.Println("Invalid JID %s: no server specified", arg)
			return recipient, false
		}
		return recipient, true
	}
}
