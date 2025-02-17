package roblox

import (
	"fmt"
	"slices"
	"strings"

	"github.com/RobloxUSArmyCID/CIDBot/requests"
)

type GroupAndRole struct {
	Group Group `json:"group"`
	Role  Role  `json:"role"`
}

type Group struct {
	ID          uint64 `json:"id"`
	Name        string `json:"name"`
	MemberCount uint   `json:"memberCount"`
	Owner       *User  `json:"owner"`
}

type Role struct {
	ID   uint64 `json:"id"`
	Name string `json:"name"`
	Rank uint   `json:"rank"`
}

func getUserGroups(userID uint64) ([]*GroupAndRole, error) {
	requestUrl := fmt.Sprintf("https://groups.roblox.com/v1/users/%d/groups/roles", userID)
	response, err := requests.Get[requests.ResponseData[*GroupAndRole]](requestUrl)
	if err != nil {
		return nil, err
	}
	return response.Data, nil
}

func (g *GroupAndRole) IsSuspicious() bool {
	return slices.Contains(getSuspiciousGroups([]*GroupAndRole{g}), g)
}

func getSuspiciousGroups(groups []*GroupAndRole) []*GroupAndRole {
	keywords := []string{
		"syndicate",
		"group",
		"fam",
		"family",
		"legacy",
		"bloodline",
		"divine sister",
		"pmc",
		"task force",
		"royalty",
		"force",
		"company",
		"intelligence",
		"mi5",
		"mi6",
		"mic",
		"intel",
	}

	susGroups := []*GroupAndRole{}
	for _, group := range groups {
		if group.Group.MemberCount <= 30 ||
			slices.Contains(keywords, strings.ToLower(group.Group.Name)) {
			susGroups = append(susGroups, group)
		}
	}

	return susGroups
}
