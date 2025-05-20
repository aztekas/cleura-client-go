package cleura

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// Get Cloud Profile Data.
// This can be used get available kubernetes versions,machine types and images suitable for
// specification in shoot clusters/ worker groups.
func (c *Client) GetCloudProfile(gardenDomain string) (*CloudProfile, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/gardener/v1/%s/cloudprofile", c.HostURL, gardenDomain), nil)
	if err != nil {
		return nil, err
	}
	var profiles []CloudProfile
	profileData, err := c.doRequest(req, 200)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(profileData, &profiles)
	if err != nil {
		return nil, err
	}
	var cProfile CloudProfile
	if !(len(profiles) == 1) || !(profiles[0].Name == "cleuracloud") {
		return nil, fmt.Errorf("something is wrong with profile data content, got: %+v", profiles)
	}
	cProfile = profiles[0]

	return &cProfile, nil
}
