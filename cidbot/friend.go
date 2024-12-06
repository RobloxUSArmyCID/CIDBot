package cidbot

import (
	"fmt"

	"github.com/adrg/strutil"
	"github.com/adrg/strutil/metrics"
)

func GetUserFriends(userID uint64) ([]*User, error) {
	requestUrl := fmt.Sprintf("https://friends.roblox.com/v1/users/%d/friends", userID)
	response, err := GetRequest[ResponseData[*User]](requestUrl)
	if err != nil {
		return nil, err
	}
	return response.Data, nil
}

func GetSuspiciousFriends(user *User, friends []*User) (susFriends []*User) {
	jarowinkler := metrics.NewJaroWinkler()
	jarowinkler.CaseSensitive = false

	for _, friend := range friends {
		similarity := strutil.Similarity(friend.Name, user.Name, jarowinkler)
		if similarity >= 0.72 {
			susFriends = append(susFriends, friend)
		}
	}
	return susFriends
}
