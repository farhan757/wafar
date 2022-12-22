package controllers

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"waapifar/dto"
	"waapifar/service"

	"github.com/gin-gonic/gin"
)

type MessageController struct{}

var messageService = new(service.MessageService)
var reqSend = new(dto.ReqSend)

func (ctrl MessageController) Send(c *gin.Context) {
	fmt.Printf("ClientIP: %s\n", c.ClientIP())

	// Call BindJSON to bind the received JSON to
	// newAlbum.
	if err := c.BindJSON(&reqSend); err != nil {
		c.JSON(http.StatusOK, gin.H{"message": "Invalid Request", "data": err})
	} else {
		messageService.SendMessage(reqSend)
		c.JSON(http.StatusOK, gin.H{"message": "Successfully"})
	}

}

func (ctrl MessageController) SendImage(c *gin.Context) {
	fmt.Printf("ClientIP: %s\n", c.ClientIP())

	c.Request.ParseMultipartForm(10 << 20)

	fl, handle, err := c.Request.FormFile("pict")

	if err != nil {
		fmt.Println("Error Retrieving the File")
		fmt.Println(err)
		return
	}

	defer fl.Close()
	fmt.Printf("Uploaded File: %+v\n", handle.Filename)
	fmt.Printf("File Size: %+v\n", handle.Size)
	fmt.Printf("MIME Header: %+v\n", handle.Header)

	fileBytes, err := ioutil.ReadAll(fl)
	if err != nil {
		fmt.Println(err)
	}

	var reqSendImg dto.ReqSendPict
	reqSendImg.No = c.Request.FormValue("no")
	reqSendImg.Pict = fileBytes
	reqSendImg.Text = c.Request.FormValue("text")

	if err := c.Bind(&reqSendImg); err != nil {
		c.JSON(http.StatusOK, gin.H{"message": "Invalid Request"})
	}

	messageService.SendGambar(&reqSendImg)
	c.JSON(http.StatusOK, gin.H{"message": "Successfully"})
}

func (ctrl MessageController) SendFile(c *gin.Context) {
	fmt.Printf("ClientIP: %s\n", c.ClientIP())

	c.Request.ParseMultipartForm(10 << 20)

	fl, handle, err := c.Request.FormFile("pdf")

	if err != nil {
		fmt.Println("Error Retrieving the File")
		fmt.Println(err)
		return
	}

	defer fl.Close()
	fmt.Printf("Uploaded File: %+v\n", handle.Filename)
	fmt.Printf("File Size: %+v\n", handle.Size)
	fmt.Printf("MIME Header: %+v\n", handle.Header)

	fileBytes, err := ioutil.ReadAll(fl)
	if err != nil {
		fmt.Println(err)
	}

	var reqSendPdf dto.ReqSendPdf
	reqSendPdf.No = c.Request.FormValue("no")
	reqSendPdf.Pdf = fileBytes
	reqSendPdf.Text = c.Request.FormValue("text")

	if err := c.Bind(&reqSendPdf); err != nil {
		c.JSON(http.StatusOK, gin.H{"message": "Invalid Request"})
	}

	messageService.SendFilePdf(&reqSendPdf)
	c.JSON(http.StatusOK, gin.H{"message": "Successfully"})
}
