package pipedrive

import (
	"context"
	"encoding/json"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceNotes() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceNotesRead,
		Schema: map[string]*schema.Schema{ //TODO add option get all notes
			"id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"content": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"add_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"update_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"deal_attached": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceNotesRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics
	id := d.Get("id").(string)
	path := "notes/" + id
	resp, _, _, err := m.(*Client).SendRequest("GET", path, nil, 200)

	if err != nil {
		return diag.FromErr(err)
	}

	var result map[string]any
	err = json.Unmarshal([]byte(resp), &result)
	if err != nil {
		return diag.FromErr(err)
	}

	content := result["data"].(map[string]interface{})["content"]
	add_time := result["data"].(map[string]interface{})["add_time"]
	update_time := result["data"].(map[string]interface{})["update_time"]
	deal_attached := result["data"].(map[string]interface{})["deal"].(map[string]interface{})["title"]

	d.Set("content", content)
	d.Set("add_time", add_time)
	d.Set("update_time", update_time)
	d.Set("deal_attached", deal_attached)

	// always run
	d.SetId(id)

	return diags
}
