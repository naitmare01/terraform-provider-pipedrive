package pipedrive

import (
	"context"
	"encoding/json"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceOrganizations() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceOrganizationsRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"add_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceOrganizationsRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics
	id := d.Get("id").(string)
	path := "organizations/" + id
	resp, _, _, err := m.(*Client).SendRequest("GET", path, nil, 200)

	if err != nil {
		return diag.FromErr(err)
	}

	var result map[string]any
	err = json.Unmarshal([]byte(resp), &result)
	if err != nil {
		return diag.FromErr(err)
	}

	name := result["data"].(map[string]interface{})["name"]
	add_time := result["data"].(map[string]interface{})["add_time"]

	d.Set("name", name)
	d.Set("add_time", add_time)

	d.SetId(id)

	return diags
}
