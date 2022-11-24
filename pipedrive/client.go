package pipedrive

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// NewClient creates common settings
func NewClient(apitoken string) *Client {
	return &Client{
		apitoken:   "?api_token=" + apitoken,
		baseurl:    "https://api.pipedrive.com/v1/",
		httpClient: &http.Client{},
	}
}

func DealsBody(d *schema.ResourceData) DealsBodyPost {
	body := DealsBodyPost{
		Title: d.Get("title").(string),
		OrgID: d.Get("org_id").(string),
	}

	status, status_true := d.GetOk("status")

	if status_true {
		body.Status = status.(string)
	}

	return body
}

func NotesBody(d *schema.ResourceData) NotesBodyPost {
	body := NotesBodyPost{
		Content: d.Get("content").(string),
		DealID:  d.Get("deal_id").(string),
	}

	return body
}

func (c *Client) SendRequest(method string, path string, payload interface{}, statusCode int) (value string, respheaders string, respCode int, err error) {
	apiurl := c.baseurl + path + c.apitoken
	client := &http.Client{Timeout: 10 * time.Second}

	b := new(bytes.Buffer)
	err = json.NewEncoder(b).Encode(payload)
	if err != nil {
		return "", "", 0, err
	}

	req, err := http.NewRequest(method, apiurl, b)
	if err != nil {
		return "", "", 0, err
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return "", "", resp.StatusCode, err
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", "", resp.StatusCode, err
	}
	resp.Body.Close()

	strbody := string(body)

	respHeaders := resp.Header
	headers, err := json.Marshal(respHeaders)
	if err != nil {
		return "", "", resp.StatusCode, err
	}

	if statusCode != 0 {
		if resp.StatusCode != statusCode {
			return "", "", 0, fmt.Errorf("[ERROR] unexpected status code got: %v expected: %v \n %v", resp.StatusCode, statusCode, strbody)
		}
	}

	return strbody, string(headers), resp.StatusCode, nil
}
