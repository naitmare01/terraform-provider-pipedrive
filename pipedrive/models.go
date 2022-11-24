package pipedrive

import (
	"net/http"
)

type DealsBodyPost struct {
	Title  string `json:"title,omitempty"`
	OrgID  string `json:"org_id,omitempty"`
	Status string `json:"status,omitempty"`
}

type Client struct {
	apitoken   string "omitempty"
	baseurl    string "omitempty"
	httpClient *http.Client
}

type NotesBodyPost struct {
	Content string `json:"content,omitempty"`
	DealID  string `json:"deal_id,omitempty"`
}
