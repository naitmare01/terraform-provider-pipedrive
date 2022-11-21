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
		},
		DataSourcesMap: map[string]*schema.Resource{
			"pipedrive_deals":         dataSourceDeals(),
			"pipedrive_organizations": dataSourceOrganizations(),
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

	title := d.Get("title").(string)
	status := d.Get("status").(string)
	org_id := d.Get("org_id").(string)

	if title != "" {
		result += `"title": "` + title + `",`
	}
	if status != "" {
		result += `"status": "` + status + `",`
	}
	if org_id != "" {
		result += `"org_id": "` + org_id + `",`
	}

	result = strings.TrimSuffix(result, ",")
	result = resultstart + result + resultend
	returnresult := strings.NewReader(result)

	return returnresult
}
