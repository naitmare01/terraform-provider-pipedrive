package pipedrive

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

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
	client := &http.Client{Timeout: 10 * time.Second}
	title := d.Get("title").(string)
	org_id := d.Get("org_id").(string)
	apikey := m.(*Client).apitoken
	baseurl := m.(*Client).baseurl
	apiurl := fmt.Sprintf("%s/deals%s", baseurl, apikey)
	payload := strings.NewReader(`{
		"title": "` + title + `",` +
		`"org_id": "` + org_id + `"}`)

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
	org_name := result["data"].(map[string]interface{})["org_name"]

	d.Set("org_name", org_name)
	// always run
	d.SetId(id_string)

	return diags
}

func resourceDealsRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics
	return diags
}

func resourceDealsUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics
	return diags
}

func resourceDealsDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics
	return diags
}