package roblox

import (
	"fmt"

	"github.com/RobloxUSArmyCID/CIDBot/requests"
)

type Thumbnail struct {
	ImageUrl string `json:"imageUrl"`
}

func GetUserThumbnail(userID uint64) (*string, error) {
	requestUrl := fmt.Sprintf("https://thumbnails.roblox.com/v1/users/avatar-headshot?userIds=%d&size=150x150&format=Webp&isCircular=false", userID)
	response, err := requests.Get[requests.ResponseData[*Thumbnail]](requestUrl)
	if err != nil {
		return nil, err
	}
	return &response.Data[0].ImageUrl, nil
}
