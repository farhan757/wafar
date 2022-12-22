package service

import (
	"context"
	"fmt"
	"waapifar/config"
	"waapifar/dto"

	"go.mau.fi/whatsmeow"
	waProto "go.mau.fi/whatsmeow/binary/proto"
	"google.golang.org/protobuf/proto"
)

// Message Service

type MessageService struct{}

func (s MessageService) SendMessage(req *dto.ReqSend) {

	//send
	recipient, ok := config.ParseJID(req.No)
	if !ok {
		return
	}

	msg := &waProto.Message{Conversation: proto.String(req.Text)}

	resp, err := config.GetClient().SendMessage(context.Background(), recipient, "", msg)
	if err != nil {
		fmt.Println("Error sending message: %v", err)
	} else {
		fmt.Println("Message sent (server timestamp: %s)", resp.Timestamp)
	}
}

func (s *MessageService) SendGambar(req *dto.ReqSendPict) {
	recipient, ok := config.ParseJID(req.No)
	if !ok {
		return
	}

	respi, erri := config.GetClient().Upload(context.Background(), req.Pict, whatsmeow.MediaImage)

	if erri != nil {
		return
	}

	msg := &waProto.ImageMessage{
		Caption:       proto.String(req.Text),
		Mimetype:      proto.String("image/png"),
		Url:           &respi.URL,
		DirectPath:    &respi.DirectPath,
		MediaKey:      respi.MediaKey,
		FileEncSha256: respi.FileEncSHA256,
		FileSha256:    respi.FileSHA256,
		FileLength:    &respi.FileLength,
	}

	resp, err := config.GetClient().SendMessage(context.Background(), recipient, "", &waProto.Message{
		ImageMessage: msg,
	})

	if err != nil {
		fmt.Println("Error sending message: %v", err)
	} else {
		fmt.Println("Message sent (server timestamp: %s)", resp.Timestamp)
	}
}

func (s *MessageService) SendFilePdf(req *dto.ReqSendPdf) {
	recipient, ok := config.ParseJID(req.No)
	if !ok {
		return
	}

	respi, erri := config.GetClient().Upload(context.Background(), req.Pdf, whatsmeow.MediaDocument)

	if erri != nil {
		return
	}

	msg := &waProto.DocumentMessage{
		Caption:       proto.String(req.Text),
		Url:           &respi.URL,
		DirectPath:    &respi.DirectPath,
		Mimetype:      proto.String("application/pdf"),
		MediaKey:      respi.MediaKey,
		FileEncSha256: respi.FileEncSHA256,
		FileSha256:    respi.FileSHA256,
		FileLength:    &respi.FileLength,
		FileName:      proto.String("Pdf Encrypt.pdf"),
	}

	resp, err := config.GetClient().SendMessage(context.Background(), recipient, "", &waProto.Message{
		DocumentMessage: msg,
	})

	if err != nil {
		fmt.Println("Error sending message: %v", err)
	} else {
		fmt.Println("Message sent (server timestamp: %s)", resp.Timestamp)
	}
}

func (s *MessageService) SendButton(req *dto.ReqSend) {
	//send
	recipient, ok := config.ParseJID(req.No)
	if !ok {
		return
	}

	this_message := &waProto.Message{
		TemplateMessage: &waProto.TemplateMessage{
			HydratedTemplate: &waProto.TemplateMessage_HydratedFourRowTemplate{
				TemplateId:          proto.String("id1"),
				HydratedContentText: proto.String(req.Text),
				HydratedButtons: []*waProto.HydratedTemplateButton{
					{
						Index: proto.Uint32(1),
						HydratedButton: &waProto.HydratedTemplateButton_QuickReplyButton{
							QuickReplyButton: &waProto.HydratedTemplateButton_HydratedQuickReplyButton{
								DisplayText: proto.String("Reply"),
								Id:          proto.String("id"),
							},
						},
					},
					{
						Index: proto.Uint32(2),
						HydratedButton: &waProto.HydratedTemplateButton_QuickReplyButton{
							QuickReplyButton: &waProto.HydratedTemplateButton_HydratedQuickReplyButton{
								DisplayText: proto.String("Menu"),
								Id:          proto.String("menu"),
							},
						},
					},
				},
			},
		},
	}

	resp, err := config.GetClient().SendMessage(context.Background(), recipient, "", this_message)
	if err != nil {
		fmt.Println("Error sending message: %v", err)
	} else {
		fmt.Println("Message sent (server timestamp: %s)", resp.Timestamp)
	}
}
