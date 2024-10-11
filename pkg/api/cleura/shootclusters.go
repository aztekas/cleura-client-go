package cleura

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

func (c *Client) GetShootCluster(clusterName string, clusterRegion string, clusterProject string) (*ShootClusterResponse, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/gardener/v1/public/shoot/%s/%s/%s", c.HostURL, clusterRegion, clusterProject, clusterName), nil)
	//https://rest.cleura.cloud/gardener/v1/:gardenDomain/shoot/:region/:project/:shootName
	if err != nil {
		return nil, err
	}

	body, err := c.doRequest(req, 200)
	if err != nil {
		return nil, err
	}
	shoot := ShootClusterResponse{}
	err = json.Unmarshal(body, &shoot)
	if err != nil {
		return nil, err
	}
	return &shoot, nil
}

func (c *Client) ListShootClusters(clusterRegion string, clusterProject string) ([]ShootClusterResponse, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/gardener/v1/public/shoot/%s/%s", c.HostURL, clusterRegion, clusterProject), nil)
	//https://rest.cleura.cloud/gardener/v1/:gardenDomain/shoot/:region/:project
	if err != nil {
		return nil, err
	}

	body, err := c.doRequest(req, 200)
	if err != nil {
		return nil, err
	}
	shoots := []ShootClusterResponse{}
	err = json.Unmarshal(body, &shoots)
	if err != nil {
		return nil, err
	}
	return shoots, nil
}

func (c *Client) CreateShootCluster(clusterRegion string, clusterProject string, shootClusterRequest ShootClusterRequest) (*ShootClusterCreateResponse, error) {
	//https://rest.cleura.cloud/gardener/v1/:gardenDomain/shoot/:region/:project
	crJsonByte, err := json.Marshal(shootClusterRequest)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/gardener/v1/public/shoot/%s/%s", c.HostURL, clusterRegion, clusterProject), strings.NewReader(string(crJsonByte)))
	if err != nil {
		return nil, err
	}
	body, err := c.doRequest(req, 200)
	if err != nil {
		return nil, err
	}
	var createdShootCluster ShootClusterCreateResponse
	err = json.Unmarshal(body, &createdShootCluster)
	if err != nil {
		return nil, err
	}

	return &createdShootCluster, nil
}

func (c *Client) DeleteShootCluster(clusterName string, clusterRegion string, clusterProject string) (string, error) {
	//https://rest.cleura.cloud/gardener/v1/:gardenDomain/shoot/:region/:project/:shoot
	req, err := http.NewRequest("DELETE", fmt.Sprintf("%s/gardener/v1/public/shoot/%s/%s/%s", c.HostURL, clusterRegion, clusterProject, clusterName), nil)
	if err != nil {
		return "", err
	}
	body, err := c.doRequest(req, 202)
	if err != nil {
		return "", err
	}
	return string(body), nil
}

func (c *Client) UpdateShootCluster(clusterRegion string, clusterProject string, clusterName string, shootClusterUpdateRequest ShootClusterRequest) (*ShootClusterResponse, error) {
	crJsonByte, err := json.Marshal(shootClusterUpdateRequest)
	if err != nil {
		return nil, err
	}
	//https://rest.cleura.cloud/gardener/v1/:gardenDomain/shoot/:region/:project/:shoot
	req, err := http.NewRequest("PUT", fmt.Sprintf("%s/gardener/v1/public/shoot/%s/%s/%s", c.HostURL, clusterRegion, clusterProject, clusterName), strings.NewReader(string(crJsonByte)))
	if err != nil {
		return nil, err
	}
	body, err := c.doRequest(req, 202)
	if err != nil {
		return nil, err
	}
	var createdShootCluster ShootClusterResponse
	err = json.Unmarshal(body, &createdShootCluster)
	if err != nil {
		return nil, err
	}

	return &createdShootCluster, nil
}

func (c *Client) AddWorkerGroup(clusterName string, clusterRegion string, clusterProject string, workerGroupRequest WorkerGroupRequest) (*ShootClusterResponse, error) {
	//https://rest.cleura.cloud/gardener/v1/:gardenDomain/shoot/:region/:project/:shoot/worker
	wgrJsonByte, err := json.Marshal(workerGroupRequest)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/gardener/v1/public/shoot/%s/%s/%s/worker", c.HostURL, clusterRegion, clusterProject, clusterName), strings.NewReader(string(wgrJsonByte)))
	if err != nil {
		return nil, err
	}
	body, err := c.doRequest(req, 202)
	if err != nil {
		return nil, err
	}
	var updatedShootCluster ShootClusterResponse
	err = json.Unmarshal(body, &updatedShootCluster)
	if err != nil {
		return nil, err
	}

	return &updatedShootCluster, nil
}

func (c *Client) UpdateWorkerGroup(clusterName string, clusterRegion string, clusterProject string, workerName string, workerGroupRequest WorkerGroupRequest) (*ShootClusterResponse, error) {
	// https://rest.cleura.cloud/gardener/v1/:gardenDomain/shoot/:region/:project/:shoot/worker/:workerName
	wgrJsonByte, err := json.Marshal(workerGroupRequest)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest("PUT", fmt.Sprintf("%s/gardener/v1/public/shoot/%s/%s/%s/worker/%s", c.HostURL, clusterRegion, clusterProject, clusterName, workerName), strings.NewReader(string(wgrJsonByte)))
	if err != nil {
		return nil, err
	}
	body, err := c.doRequest(req, 202)
	if err != nil {
		return nil, err
	}
	var updatedShootCluster ShootClusterResponse
	err = json.Unmarshal(body, &updatedShootCluster)
	if err != nil {
		return nil, err
	}

	return &updatedShootCluster, nil
}

func (c *Client) DeleteWorkerGroup(clusterName string, clusterRegion string, clusterProject string, workerName string) (*ShootClusterResponse, error) {
	//https://rest.cleura.cloud/gardener/v1/:gardenDomain/shoot/:region/:project/:shoot/worker/:worker

	req, err := http.NewRequest("DELETE", fmt.Sprintf("%s/gardener/v1/public/shoot/%s/%s/%s/worker/%s", c.HostURL, clusterRegion, clusterProject, clusterName, workerName), nil)
	if err != nil {
		return nil, err
	}
	body, err := c.doRequest(req, 202)
	if err != nil {
		return nil, err
	}
	var updatedShootCluster ShootClusterResponse
	err = json.Unmarshal(body, &updatedShootCluster)
	if err != nil {
		return nil, err
	}
	return &updatedShootCluster, nil
}

func (c *Client) GenerateKubeConfig(clusterRegion string, clusterProject string, clusterName string, durationSeconds int64) ([]byte, error) {
	//https://rest.cleura.cloud/gardener/v1/public/shoot/kna1/b5d2bf2c162444f4918aaa4cb534a612/myshoot/adminkubeconfig

	type Config struct {
		ExpirationSeconds int64 `json:"expirationSeconds"`
	}
	type Request struct {
		Config Config `json:"config"`
	}
	kubeConfigRequest := Request{
		Config: Config{
			ExpirationSeconds: durationSeconds,
		},
	}
	requestJsonByte, err := json.Marshal(kubeConfigRequest)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/gardener/v1/public/shoot/%s/%s/%s/adminkubeconfig", c.HostURL, clusterRegion, clusterProject, clusterName), strings.NewReader(string(requestJsonByte)))
	if err != nil {
		return nil, err
	}
	body, err := c.doRequest(req, 200)
	if err != nil {
		return nil, err
	}
	return body, nil
}

// Hibernate.
func (c *Client) HibernateCluster(clusterRegion string, clusterProject string, clusterName string) error {
	// https://rest.cleura.cloud/gardener/v1/:gardenDomain/shoot/:region/:project/:shoot/hibernate
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/gardener/v1/public/shoot/%s/%s/%s/hibernate", c.HostURL, clusterRegion, clusterProject, clusterName), nil)
	if err != nil {
		return err
	}
	_, err = c.doRequest(req, 202)
	if err != nil {
		return err
	}
	return nil
}

// Wake up call.
func (c *Client) WakeUpCluster(clusterRegion string, clusterProject string, clusterName string) error {
	// https://rest.cleura.cloud/gardener/v1/:gardenDomain/shoot/:region/:project/:shoot/wakeup
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/gardener/v1/public/shoot/%s/%s/%s/wakeup", c.HostURL, clusterRegion, clusterProject, clusterName), nil)
	if err != nil {
		return err
	}
	_, err = c.doRequest(req, 202)
	if err != nil {
		return err
	}
	return nil
}
