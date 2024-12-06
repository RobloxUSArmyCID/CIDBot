package cidbot

import "fmt"

type Badge struct {}

func GetUserBadges(userID uint64) ([]*Badge, error) {
	requestUrl := fmt.Sprintf("https://badges.roblox.com/v1/users/%d/badges?limit=100", userID)
	response, err := GetRequest[ResponseData[*Badge]](requestUrl)
	if err != nil {
		return nil, err
	}
	return response.Data, nil
}