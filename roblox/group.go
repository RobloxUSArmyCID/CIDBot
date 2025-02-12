package roblox

import (
	"fmt"
	"slices"
	"strings"

	"github.com/RobloxUSArmyCID/CIDBot/requests"
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

func getUserGroups(userID uint64) ([]*Group, error) {
	requestUrl := fmt.Sprintf("https://groups.roblox.com/v2/users/%d/groups/roles", userID)
	response, err := requests.Get[requests.ResponseData[*Group]](requestUrl)
	if err != nil {
		return nil, err
	}
	return response.Data, nil
}

func (g *Group) IsSuspicious() bool {
	return slices.Contains(getSuspiciousGroups([]*Group{g}), g)
}

func getSuspiciousGroups(groups []*Group) []*Group {
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
	}

	susGroups := []*Group{}
	for _, group := range groups {
		if group.Group.MemberCount <= 30 {
			susGroups = append(susGroups, group)
		}
		if slices.Contains(keywords, strings.ToLower(group.Group.Name)) {
			susGroups = append(susGroups, group)
		}
	}

	return susGroups
}
