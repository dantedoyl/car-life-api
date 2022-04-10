package vk

import (
	"fmt"
	"github.com/SevereCloud/vksdk/v2/api"
	"github.com/SevereCloud/vksdk/v2/api/params"
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