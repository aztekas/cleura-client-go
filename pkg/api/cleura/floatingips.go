package cleura

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func (c *Client) GetFloatingIPs(region string, projectId string) ([]FloatingIp, error) {
	// https://rest.cleura.cloud/networking/v1/floatingips/<region>/<project_id>

	req, err := http.NewRequest("GET", fmt.Sprintf("%s/networking/v1/floatingips/%s/%s", c.HostURL, region, projectId), nil)
	if err != nil {
		return nil, err
	}
	var floatingIps []FloatingIp
	floatingIpData, err := c.doRequest(req, 200)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(floatingIpData, &floatingIps)
	if err != nil {
		return nil, err
	}

	return floatingIps, nil
}
