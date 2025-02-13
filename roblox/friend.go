package roblox

import (
	"fmt"

	"github.com/RobloxUSArmyCID/CIDBot/requests"
	"github.com/adrg/strutil"
	"github.com/adrg/strutil/metrics"
)

func getUserFriends(userID uint64) ([]*User, error) {
	requestUrl := fmt.Sprintf("https://friends.roblox.com/v1/users/%d/friends", userID)
	response, err := requests.Get[requests.ResponseData[*User]](requestUrl)
	if err != nil {
		return nil, err
	}
	return response.Data, nil
}

func getSuspiciousFriends(username string, friends []*User) []*User {
	jarowinkler := metrics.NewJaroWinkler()
	jarowinkler.CaseSensitive = false

	susFriends := []*User{}

	for _, friend := range friends {
		similarity := strutil.Similarity(friend.Name, username, jarowinkler)
		if similarity >= 0.72 {
			susFriends = append(susFriends, friend)
		}
	}

	return susFriends
}
