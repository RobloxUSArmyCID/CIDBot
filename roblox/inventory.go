package roblox

import (
	"fmt"

	"github.com/RobloxUSArmyCID/CIDBot/requests"
)

type CanViewInventory struct {
	CanView bool `json:"canView"`
}

func (u *User) GetInventoryVisibility() error {
	requestUrl := fmt.Sprintf("https://inventory.roblox.com/v1/users/%d/can-view-inventory", u.ID)
	response, err := requests.Get[CanViewInventory](requestUrl)
	if err != nil {
		return err
	}
	u.mu.Lock()
	defer u.mu.Unlock()
	u.CanViewInventory = response.CanView
	return nil
}
