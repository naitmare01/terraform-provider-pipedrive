package pipedrive

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceOrganizations() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"add_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
		CreateContext: resourceOrganizationsCreate,
		ReadContext:   resourceOrganizationsRead,
		UpdateContext: resourceOrganizationsUpdate,
		DeleteContext: resourceOrganizationsDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func resourceOrganizationsCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	path := "organizations"
	body := OrganizationsBody(d)
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
	name := result["data"].(map[string]interface{})["name"]
	add_time := result["data"].(map[string]interface{})["add_time"]

	d.Set("name", name)
	d.Set("add_time", add_time)
	d.SetId(id_string)

	return diags
}

func resourceOrganizationsRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	id := d.Id()
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

	return diags
}

func resourceOrganizationsUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	id := d.Id()
	path := "organizations/" + id
	body := NotesBody(d)
	_, _, _, err := m.(*Client).SendRequest("PUT", path, body, 200)

	if err != nil {
		return diag.FromErr(err)
	}

	return resourceOrganizationsRead(ctx, d, m)
}

func resourceOrganizationsDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	id := d.Id()
	path := "organizations/" + id
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
