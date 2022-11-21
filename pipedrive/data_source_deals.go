package pipedrive

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceDeals() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceDealsRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"title": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"org_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceDealsRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := &http.Client{Timeout: 10 * time.Second}
	id := d.Get("id").(string)
	apikey := m.(*Client).apitoken
	baseurl := m.(*Client).baseurl
	apiurl := fmt.Sprintf("%s/deals/%s%s", baseurl, id, apikey)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	req, err := http.NewRequest("GET", apiurl, nil)
	if err != nil {
		return diag.FromErr(err)
	}

	r, err := client.Do(req)
	if err != nil {
		return diag.FromErr(err)
	}
	defer r.Body.Close()

	var result map[string]any
	err = json.NewDecoder(r.Body).Decode(&result)
	if err != nil {
		return diag.FromErr(err)
	}
	title := result["data"].(map[string]interface{})["title"]
	org_name := result["data"].(map[string]interface{})["org_name"]
	status := result["data"].(map[string]interface{})["status"]

	d.Set("title", title)
	d.Set("org_name", org_name)
	d.Set("status", status)

	// always run
	d.SetId(id)

	return diags
}
