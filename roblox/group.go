package roblox

import (
	"fmt"
	"slices"
	"strings"

	"github.com/RobloxUSArmyCID/CIDBot/requests"
	"github.com/RobloxUSArmyCID/CIDBot/roblox/usar"
)

type GroupAndRole struct {
	Group Group `json:"group"`
	Role  Role  `json:"role"`
}

type Group struct {
	ID          uint64      `json:"id"`
	Name        string      `json:"name"`
	MemberCount uint        `json:"memberCount"`
	Owner       *GroupOwner `json:"owner"`
}

type Role struct {
	ID   uint64 `json:"id"`
	Name string `json:"name"`
	Rank uint   `json:"rank"`
}

type GroupOwner struct {
	Name string `json:"username"`
	ID   uint64 `json:"userId"`
}

func (u *User) GetGroups() error {
	const (
		UsarGroupID = 3108077
		RankE1      = 5
	)

	groups, err := getUserGroups(u.ID)
	if err != nil {
		return err
	}

	isE1 := false
	isInUsar := false
	usarRank := "N/A"

	for _, group := range groups {
		if group.Group.ID == UsarGroupID {
			isInUsar = true
			usarRank = group.Role.Name
			isE1 = group.Role.Rank == RankE1
		}
	}

	u.mu.Lock()
	defer u.mu.Unlock()

	u.Groups = groups
	u.GetSuspiciousGroups()
	u.GetUsarUnits()
	u.IsE1 = isE1
	u.IsInUsar = isInUsar
	u.UsarRank = usarRank
	return nil
}

func getUserGroups(userID uint64) ([]GroupAndRole, error) {
	requestUrl := fmt.Sprintf("https://groups.roblox.com/v1/users/%d/groups/roles", userID)
	response, err := requests.Get[requests.ResponseData[GroupAndRole]](requestUrl)
	if err != nil {
		return nil, err
	}

	return response.Data, nil
}

func (u *User) GetSuspiciousGroups() {
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

	susGroups := []GroupAndRole{}
	for _, group := range u.Groups {
		groupNotAlreadyInList := !slices.Contains(susGroups, group)
		groupBelow30Members := group.Group.MemberCount <= 30
		groupBelow10KMembers := group.Group.MemberCount < 10_000
		_, groupIsAUsarGroup := usar.Groups[group.Group.ID]

		if groupBelow30Members && groupNotAlreadyInList && !groupIsAUsarGroup {
			susGroups = append(susGroups, group)
		}

		for _, keyword := range keywords {
			groupHasKeyword := strings.Contains(strings.ToLower(group.Group.Name), keyword)

			if groupHasKeyword && groupNotAlreadyInList && groupBelow10KMembers && !groupIsAUsarGroup {
				susGroups = append(susGroups, group)
			}
		}

	}

	u.mu.Lock()
	defer u.mu.Unlock()
	u.SuspiciousGroups = susGroups
}

func (u *User) GetUsarUnits() {
	units := []string{}
	for _, group := range u.Groups {
		if unit, ok := usar.Groups[group.Group.ID]; ok {
			units = append(units, unit)
		}
	}

	u.mu.Lock()
	defer u.mu.Unlock()
	u.UsarUnits = units
}
