package cidbot

import "fmt"

func GetUserFriends(userID uint64) ([]*User, error) {
	requestUrl := fmt.Sprintf("https://friends.roblox.com/v1/users/%d/friends", userID)
	response, err := GetRequest[ResponseData[*User]](requestUrl)
	if err != nil {
		return nil, err
	}
	return response.Data, nil
}
