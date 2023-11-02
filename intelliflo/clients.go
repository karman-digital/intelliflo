package intelliflo

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/hashicorp/go-retryablehttp"
	intelliflomodels "github.com/karman-digital/intelliflo/models"
)

func (c *credentials) CreateNewClient(clientObj intelliflomodels.Client) (intelliflomodels.Client, error) {
	responseClient := intelliflomodels.Client{}
	clientBody, err := json.Marshal(clientObj)
	if err != nil {
		return responseClient, fmt.Errorf("error converting to body: %v", err)
	}
	clientReader := bytes.NewReader(clientBody)
	req, err := retryablehttp.NewRequest("POST", "https://api.gb.intelliflo.net/v2/clients", clientReader)
	if err != nil {
		return responseClient, fmt.Errorf("error creating request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header["x-api-key"] = []string{c.apiKey.String()}
	req.Header["authorization"] = []string{fmt.Sprintf("Bearer %s", c.accessToken)}
	resp, err := c.client.Do(req)
	if err != nil {
		return responseClient, fmt.Errorf("error making post request: %v", err)
	}
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return responseClient, fmt.Errorf("error reading body: %v", err)
	}
	if resp.StatusCode != http.StatusCreated {
		return responseClient, fmt.Errorf("error when creating client: %s", string(respBody))
	}
	defer resp.Body.Close()
	err = json.Unmarshal(respBody, &responseClient)
	if err != nil {
		return responseClient, fmt.Errorf("error unmarshalling response: %v", err)
	}
	return responseClient, nil
}

func (c *credentials) GetClients(options ...intelliflomodels.GetOptions) (intelliflomodels.Clients, error) {
	var clients intelliflomodels.Clients
	req, err := retryablehttp.NewRequest("GET", "https://api.gb.intelliflo.net/v2/clients", nil)
	if err != nil {
		return clients, fmt.Errorf("error creating request: %v", err)
	}
	if len(options) > 0 {
		q := req.URL.Query()
		if options[0].Skip != 0 {
			q.Add("skip", fmt.Sprintf("%d", options[0].Skip))
		}
		if options[0].Top != 0 {
			q.Add("top", fmt.Sprintf("%d", options[0].Top))
		} else {
			q.Add("top", "500")
		}
		if options[0].Filter != "" {
			q.Add("filter", options[0].Filter)
		}
		req.URL.RawQuery = q.Encode()
	} else {
		q := req.URL.Query()
		q.Add("top", "500")
		req.URL.RawQuery = q.Encode()
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header["x-api-key"] = []string{c.apiKey.String()}
	req.Header["authorization"] = []string{fmt.Sprintf("Bearer %s", c.accessToken)}
	resp, err := c.client.Do(req)
	if err != nil {
		return clients, fmt.Errorf("error making post request: %v", err)
	}
	defer resp.Body.Close()
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return clients, fmt.Errorf("error reading body: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		return clients, fmt.Errorf("error returned by endpoint, status code: %d, body: %s", resp.StatusCode, respBody)
	}
	err = json.Unmarshal(respBody, &clients)
	if err != nil {
		return clients, fmt.Errorf("error parsing body: %v", err)
	}
	return clients, nil
}
