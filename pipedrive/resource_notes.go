package pipedrive

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceNotes() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"content": {
				Type:     schema.TypeString,
				Required: true,
			},
			"deal_id": {
				Type:     schema.TypeString,
				Required: true,
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
		CreateContext: resourceNotesCreate,
		ReadContext:   resourceNotesRead,
		UpdateContext: resourceNotesUpdate,
		DeleteContext: resourceNotesDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func resourceNotesCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics
	path := "notes"
	body := NotesBody(d)
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
	add_time := result["data"].(map[string]interface{})["add_time"]
	update_time := result["data"].(map[string]interface{})["update_time"]
	deal_attached := result["data"].(map[string]interface{})["deal"].(map[string]interface{})["title"]

	d.Set("add_time", add_time)
	d.Set("update_time", update_time)
	d.Set("deal_attached", deal_attached)
	d.SetId(id_string)

	return diags
}

func resourceNotesRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics
	id := d.Id()
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
	deal_id := result["data"].(map[string]interface{})["deal_id"]
	add_time := result["data"].(map[string]interface{})["add_time"]
	update_time := result["data"].(map[string]interface{})["update_time"]
	deal_attached := result["data"].(map[string]interface{})["deal"].(map[string]interface{})["title"]

	d.Set("content", content)
	d.Set("deal_id", deal_id)
	d.Set("add_time", add_time)
	d.Set("update_time", update_time)
	d.Set("deal_attached", deal_attached)

	return diags
}

func resourceNotesUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	id := d.Id()
	path := "notes/" + id
	body := NotesBody(d)
	_, _, _, err := m.(*Client).SendRequest("PUT", path, body, 200)

	if err != nil {
		return diag.FromErr(err)
	}

	return resourceNotesRead(ctx, d, m)
}

func resourceNotesDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics
	id := d.Id()
	path := "notes/" + id
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
