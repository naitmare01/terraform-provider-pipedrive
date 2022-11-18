package pipedrive

import (
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
	apitoken string
	baseurl  string
}

// NewClient creates common settings
func NewClient(apitoken string) *Client {
	return &Client{
		apitoken: "?api_token=" + apitoken,
		baseurl:  "https://api.pipedrive.com/v1",
	}
}
