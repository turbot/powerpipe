package apiclient

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/turbot/powerpipe/internal/service/api/dto"
)

type ApiClient struct {
	host       string
	port       string
	pathPrefix string
}

func NewApiClient(host string, port string, pathPrefix string) *ApiClient {
	return &ApiClient{
		host:       host,
		port:       port,
		pathPrefix: pathPrefix,
	}
}

var Client = NewApiClient("localhost", "9194", "/api/v0")

func (api *ApiClient) GetMod(ctx context.Context) (mod *dto.GetModResponse, err error) {
	mod = &dto.GetModResponse{}
	err = api.Get(ctx, "/mod", nil, &mod)
	return
}

func (api *ApiClient) InstallMod(ctx context.Context, modId string, install *dto.InstallModRequest) (response *dto.InstallModResponse, err error) {
	response = &dto.InstallModResponse{}
	err = api.Post(ctx, "/mod/:mod_id/dependency", map[string]string{"mod_id": modId}, install, response)
	return
}

func (api *ApiClient) Get(ctx context.Context, path string, params map[string]string, response interface{}) (err error) {
	url := api.getUrl(path, params)
	req, err := http.NewRequest("GET", url, nil)
	req = req.WithContext(ctx)
	if err != nil {
		return err
	}

	// Create an HTTP client and send the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&response)
	return err
}

func (api *ApiClient) Post(ctx context.Context, path string, params map[string]string, payload any, response any) (err error) {
	buffer := bytes.NewBuffer([]byte{})
	encoder := json.NewEncoder(buffer)
	encoder.Encode(payload)

	url := api.getUrl(path, params)
	req, err := http.NewRequest("POST", url, buffer)
	req = req.WithContext(ctx)
	if err != nil {
		return err
	}

	// Set the Content-Type header to indicate that you are sending JSON data
	req.Header.Set("Content-Type", "application/json")

	// Create an HTTP client and send the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&response)
	return err
}

func (api *ApiClient) getUrl(path string, params map[string]string) string {
	url := fmt.Sprintf("http://%s:%s%s%s", api.host, api.port, api.pathPrefix, path)
	for key, value := range params {
		toReplace := fmt.Sprintf("/:%s/", key)
		url = strings.ReplaceAll(url, toReplace, fmt.Sprintf("/%s/", value))
	}
	return url
}
