package roblox_test

import (
	"slices"
	"testing"

	"github.com/RobloxUSArmyCID/CIDBot/roblox"
)

func TestGroupWithKeywordAndManyMembers(t *testing.T) {
	user := &roblox.User{
		Groups: []*roblox.GroupAndRole{
			{
				Group: roblox.Group{ // pass
					Name:        "The Killing Force", // keyword
					MemberCount: 10_000_000,
				},
			},
		},
	}

	want := []*roblox.GroupAndRole{}

	user.GetSuspiciousGroups()

	if !slices.Equal(user.SuspiciousGroups, want) {
		t.Error(`Group "The Killing Force" with 10M members was tagged as suspicious.`)
	}

}

func TestGroupWithKeywordAndLittleMembers(t *testing.T) {
	user := &roblox.User{
		Groups: []*roblox.GroupAndRole{
			{
				Group: roblox.Group{ // fail
					Name:        "The Killing Force", // keyword
					MemberCount: 100,
				},
			},
		},
	}

	want := []*roblox.GroupAndRole{
		{
			Group: roblox.Group{ // fail
				Name:        "The Killing Force", // keyword
				MemberCount: 100,
			},
		},
	}

	user.GetSuspiciousGroups()

	if !slices.Equal(user.SuspiciousGroups, want) {
		t.Error(`Group "The Killing Force" with 100 members was not tagged as suspicious.`)
	}

}

func TestUSARGroupWithKeyword(t *testing.T) {
	user := &roblox.User{
		Groups: []*roblox.GroupAndRole{
			{
				Group: roblox.Group{ // pass
					Name:        "Forces Command", // keyword
					ID:          3198375,
					MemberCount: 100,
				},
			},
		},
	}

	want := []*roblox.GroupAndRole{}

	user.GetSuspiciousGroups()

	if !slices.Equal(user.SuspiciousGroups, want) {
		t.Error(`USAR Group "Forces Command" with 100 members was tagged as suspicious.`)
	}
}
