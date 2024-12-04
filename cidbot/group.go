package cidbot

import (
	"fmt"
)

type Group struct {
	Group internalGroup `json:"group"`
	Role  role          `json:"role"`
}

type internalGroup struct {
	ID          uint64 `json:"id"`
	Name        string `json:"name"`
	MemberCount uint   `json:"memberCount"`
}

type role struct {
	ID   uint64 `json:"id"`
	Name string `json:"name"`
	Rank uint   `json:"rank"`
}

func GetUserGroups(userID uint64) ([]*Group, error) {
	requestUrl := fmt.Sprintf("https://users.roblox.com/v1/users/%d", userID)
	response, err := GetRequest[ResponseData[*Group]](requestUrl)
	if err != nil {
		return nil, err
	}
	return response.Data, nil
}
