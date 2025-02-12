package roblox

import (
	"fmt"

	"github.com/RobloxUSArmyCID/CIDBot/requests"
)

type Badge struct{}

func getUserBadges(userID uint64) ([]*Badge, error) {
	requestUrl := fmt.Sprintf("https://badges.roblox.com/v1/users/%d/badges?limit=100", userID)
	response, err := requests.Get[requests.ResponseData[*Badge]](requestUrl)
	if err != nil {
		return nil, err
	}
	return response.Data, nil
}
