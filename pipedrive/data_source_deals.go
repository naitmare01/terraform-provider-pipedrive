package pipedrive

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
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
		},
	}
}

func dataSourceDealsRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := &http.Client{Timeout: 10 * time.Second}
	id := d.Get("id").(string)
	//apikey := d.Get("apikey").(string)
	apikey := m.(*Client).apitoken
	apiurl := fmt.Sprintf("https://api.pipedrive.com/v1/deals/%s?api_token=%s", id, apikey)

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
	d.Set("title", title)

	// always run
	d.SetId(strconv.FormatInt(time.Now().Unix(), 10))

	return diags
}
