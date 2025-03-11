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

	// O(n) with n-max = 100
	// MaxO(100) -- not worth of optimization
	susGroups := getSuspiciousGroups(groups)

	u.mu.Lock()
	defer u.mu.Unlock()

	u.Groups = groups
	u.IsE1 = isE1
	u.IsInUsar = isInUsar
	u.UsarRank = usarRank
	u.SuspiciousGroups = susGroups
	return nil
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
		if group.Group.MemberCount <= 30 && !slices.Contains(susGroups, group) {
			susGroups = append(susGroups, group)
		}

		for _, keyword := range keywords {
			if strings.Contains(strings.ToLower(group.Group.Name), keyword) && !slices.Contains(susGroups, group) {
				susGroups = append(susGroups, group)
			}
		}

	}

	return susGroups
}
