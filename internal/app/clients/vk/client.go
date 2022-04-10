package vk

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/SevereCloud/vksdk/v2/api"
	"github.com/SevereCloud/vksdk/v2/api/params"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"time"
)

type VKClient struct {
	 serviceClient *api.VK
	 groupClient *api.VK
}

func NewVKClient(serviceKey string, groupKey string) *VKClient {
	return &VKClient{
		serviceClient: api.NewVK(serviceKey),
		groupClient: api.NewVK(groupKey),
	}
}

func (vk *VKClient) CreatChat(title string) (int, error){
	chat := params.NewMessagesCreateChatBuilder()
	chat.Title(title)
	chatInfo, err := vk.groupClient.MessagesCreateChat(chat.Params)
	if err != nil {
		return 0, err
	}

	return chatInfo, nil
}

func (vk *VKClient) UploadChatPhoto(id int, fileHeader *multipart.FileHeader) error {
	chat := params.NewPhotosGetChatUploadServerBuilder()
	chat.ChatID(id)
	chatInfo, err := vk.groupClient.PhotosGetChatUploadServer(chat.Params)
	if err != nil {
		return err
	}
	client := &http.Client{
		Timeout: time.Second * 10,
	}
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	fw, err := writer.CreateFormField("file")
	if err != nil {
		return err
	}

	file, err := fileHeader.Open()
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = io.Copy(fw, file)
	if err != nil {
		return err
	}
	writer.Close()

	req, err := http.NewRequest("POST", chatInfo.UploadURL, bytes.NewReader(body.Bytes()))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())
	rsp, _ := client.Do(req)
	if rsp.StatusCode != http.StatusOK {
		log.Printf("Request failed with response code: %d", rsp.StatusCode)
	}
	defer rsp.Body.Close()

	bodyBytes, err := io.ReadAll(rsp.Body)
	if err != nil {
		return err
	}

	resp := struct {
		Response string
	}{}

	err = json.Unmarshal(bodyBytes, &resp)
	if err != nil {
		return err
	}

	chatUpload := params.NewMessagesSetChatPhotoBuilder()
	chatUpload.File(resp.Response)
	_, err = vk.groupClient.MessagesSetChatPhoto(chatUpload.Params)
	if err != nil {
		return err
	}

	return nil
}

func (vk *VKClient) GetChatLink(id int) (string, error){
	chat := params.NewMessagesGetInviteLinkBuilder()
	chat.PeerID(2000000000 + id)
	chat.Reset(true)
	chatInfo, err := vk.groupClient.MessagesGetInviteLink(chat.Params)
	if err != nil {
		return "", err
	}

	return chatInfo.Link, nil
}

func (vk *VKClient) CreatNotification() {
	not := params.NewNotificationsSendMessageBuilder()
	not.Message("Привет")
	not.UserIDs([]int{146506479})
	a, err := vk.serviceClient.NotificationsSendMessage(not.Params)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(a)
}