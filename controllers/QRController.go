package controllers

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/mdp/qrterminal/v3"
	"go.mau.fi/whatsmeow"
)

type QRController struct{}

var client *whatsmeow.Client
var qrcode string

func (ctrl QRController) GetQR(c *gin.Context) {
	fmt.Printf("ClientIP: %s\n", c.ClientIP())

	qrChan, _ := client.GetQRChannel(context.Background())
	err := client.Connect()
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
		qrcode = evt.Code
	}
	c.JSON(http.StatusOK, gin.H{"message": "Successfully", "QRCODE": qrcode})

}
