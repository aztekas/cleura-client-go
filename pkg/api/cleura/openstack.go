package cleura

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func (c *Client) ListDomains() (*[]OpenstackDomain, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/accesscontrol/v1/openstack/domains", c.HostURL), nil)
	//https://rest.cleura.cloud/accesscontrol/v1/openstack/domains
	if err != nil {
		return nil, err
	}

	body, err := c.doRequest(req, 200)
	if err != nil {
		return nil, err
	}
	domains := []OpenstackDomain{}
	err = json.Unmarshal(body, &domains)
	if err != nil {
		return nil, err
	}
	//shoots = append(shoots, shoot)
	return &domains, nil
}

func (c *Client) ListProjects(domain_id string) (*[]OpenstackProject, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/accesscontrol/v1/openstack/%s/projects", c.HostURL, domain_id), nil)
	//https://rest.cleura.cloud/accesscontrol/v1/openstack/:domainId/projects
	if err != nil {
		return nil, err
	}

	body, err := c.doRequest(req, 200)
	if err != nil {
		return nil, err
	}
	projects := []OpenstackProject{}
	err = json.Unmarshal(body, &projects)
	if err != nil {
		return nil, err
	}
	//shoots = append(shoots, shoot)
	return &projects, nil
}
