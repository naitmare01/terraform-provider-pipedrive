package pipedrive

import (
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// Provider -
func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"apikey": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
				DefaultFunc: schema.EnvDefaultFunc("PIPEDRIVE_PASSWORD", nil),
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"pipedrive_deals": resourceDeals(),
			"pipedrive_notes": resourceNotes(),
		},
		DataSourcesMap: map[string]*schema.Resource{
			"pipedrive_deals":         dataSourceDeals(),
			"pipedrive_organizations": dataSourceOrganizations(),
			"pipedrive_notes":         dataSourceNotes(),
		},
		ConfigureFunc: providerConfigure,
	}
}

func providerConfigure(d *schema.ResourceData) (interface{}, error) {
	apitoken := d.Get("apikey").(string)
	return NewClient(apitoken), nil
}

type Client struct {
	apitoken string "omitempty"
	baseurl  string "omitempty"
}

// NewClient creates common settings
func NewClient(apitoken string) *Client {
	return &Client{
		apitoken: "?api_token=" + apitoken,
		baseurl:  "https://api.pipedrive.com/v1",
	}
}

func DealsBody(d *schema.ResourceData) *strings.Reader {
	result := ""
	resultstart := `{`
	resultend := `}`

	title, title_true := d.GetOk("title")
	status, status_true := d.GetOk("status")
	org_id, org_id_true := d.GetOk("org_id")
	content, content_true := d.GetOk("content")
	deal_id, deal_id_true := d.GetOk("deal_id")

	if title_true {
		result += `"title": "` + title.(string) + `",`
	}
	if status_true {
		result += `"status": "` + status.(string) + `",`
	}
	if org_id_true {
		result += `"org_id": "` + org_id.(string) + `",`
	}
	if content_true {
		result += `"content": "` + content.(string) + `",`
	}
	if deal_id_true {
		result += `"deal_id": "` + deal_id.(string) + `",`
	}

	result = strings.TrimSuffix(result, ",")
	result = resultstart + result + resultend
	returnresult := strings.NewReader(result)

	return returnresult
}
