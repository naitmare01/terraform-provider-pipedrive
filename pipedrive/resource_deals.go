package pipedrive

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceDeals() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"title": {
				Type:     schema.TypeString,
				Required: true,
			},
			"org_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"org_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "open",
				ValidateFunc: func(val any, key string) (warns []string, errs []error) {
					value := val.(string)
					expected := map[string]bool{
						"open":    true,
						"won":     true,
						"lost":    true,
						"deleted": true,
					}
					if !expected[value] {
						errs = append(errs, fmt.Errorf("%q is not a valid value. Please use open, won, lost, deleted", value))
					}
					return
				},
			},
		},
		CreateContext: resourceDealsCreate,
		ReadContext:   resourceDealsRead,
		UpdateContext: resourceDealsUpdate,
		DeleteContext: resourceDealsDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func resourceDealsCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics
	path := "deals"
	body := DealsBody(d)
	resp, _, _, err := m.(*Client).SendRequest("POST", path, body, 201)

	if err != nil {
		return diag.FromErr(err)
	}

	var result map[string]any
	err = json.Unmarshal([]byte(resp), &result)
	if err != nil {
		return diag.FromErr(err)
	}

	id := result["data"].(map[string]interface{})["id"]
	id_string := fmt.Sprintf("%v", id)
	org_name := result["data"].(map[string]interface{})["org_name"]

	d.Set("org_name", org_name)
	d.SetId(id_string)

	return diags
}

func resourceDealsRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics
	id := d.Id()
	path := "deals/" + id
	resp, _, _, err := m.(*Client).SendRequest("GET", path, nil, 200)

	if err != nil {
		return diag.FromErr(err)
	}

	var result map[string]any
	err = json.Unmarshal([]byte(resp), &result)
	if err != nil {
		return diag.FromErr(err)
	}

	org_name := result["data"].(map[string]interface{})["org_name"]
	title := result["data"].(map[string]interface{})["title"]
	status := result["data"].(map[string]interface{})["status"]

	d.Set("org_name", org_name)
	d.Set("title", title)
	d.Set("status", status)

	return diags
}

func resourceDealsUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	id := d.Id()
	path := "deals/" + id
	body := DealsBody(d)
	_, _, _, err := m.(*Client).SendRequest("PUT", path, body, 200)

	if err != nil {
		return diag.FromErr(err)
	}

	return resourceDealsRead(ctx, d, m)
}

func resourceDealsDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics
	id := d.Id()
	path := "deals/" + id
	resp, _, _, err := m.(*Client).SendRequest("DELETE", path, nil, 200)

	if err != nil {
		return diag.FromErr(err)
	}

	var result map[string]any
	err = json.Unmarshal([]byte(resp), &result)
	if err != nil {
		return diag.FromErr(err)
	}

	success := result["success"]

	if success == false {
		error_msg := result["error"]
		error_msg_verbose := result["error_info"]
		return diag.Errorf("%s %s", error_msg, error_msg_verbose)
	}

	// d.SetId("") is automatically called assuming delete returns no errors, but
	// it is added here for explicitness.
	d.SetId("")

	return diags
}
