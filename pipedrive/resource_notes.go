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
	client := &http.Client{Timeout: 10 * time.Second}
	apikey := m.(*Client).apitoken
	baseurl := m.(*Client).baseurl
	apiurl := fmt.Sprintf("%s/notes%s", baseurl, apikey)
	payload := DealsBody(d)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	req, err := http.NewRequest("POST", apiurl, payload)
	if err != nil {
		return diag.FromErr(err)
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")

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
	client := &http.Client{Timeout: 10 * time.Second}
	id := d.Id()
	apikey := m.(*Client).apitoken
	baseurl := m.(*Client).baseurl
	apiurl := fmt.Sprintf("%s/notes/%s%s", baseurl, id, apikey)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	req, err := http.NewRequest("GET", apiurl, nil)
	if err != nil {
		return diag.FromErr(err)
	}

	r, err := client.Do(req)
	if r.StatusCode == 404 && err != nil {
		d.SetId("")
		diag.FromErr(err)
	}
	if err != nil {
		return diag.FromErr(err)
	}
	defer r.Body.Close()

	var result map[string]any
	err = json.NewDecoder(r.Body).Decode(&result)
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
	client := &http.Client{Timeout: 10 * time.Second}
	id := d.Id()
	apikey := m.(*Client).apitoken
	baseurl := m.(*Client).baseurl
	apiurl := fmt.Sprintf("%s/notes/%s%s", baseurl, id, apikey)
	payload := DealsBody(d)

	req, err := http.NewRequest("PUT", apiurl, payload)

	if err != nil {
		return diag.FromErr(err)
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")

	r, err := client.Do(req)
	if err != nil {
		return diag.FromErr(err)
	}
	defer r.Body.Close()

	return resourceNotesRead(ctx, d, m)
}

func resourceNotesDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := &http.Client{Timeout: 10 * time.Second}
	id := d.Id()
	apikey := m.(*Client).apitoken
	baseurl := m.(*Client).baseurl
	apiurl := fmt.Sprintf("%s/notes/%s%s", baseurl, id, apikey)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics
	req, err := http.NewRequest("DELETE", apiurl, nil)

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
