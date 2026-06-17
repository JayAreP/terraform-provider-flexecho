package flexecho

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	sdk "github.com/silk-us/flexecho-go-sdk/flexecho"
)

func resourceFlexEchoHost() *schema.Resource {
	return &schema.Resource{
		Create: resourceFlexEchoHostCreate,
		Read:   resourceFlexEchoHostRead,
		Delete: resourceFlexEchoHostDelete,
		// no real update beyond a re-PUT, so just recreate on change
		Importer: &schema.ResourceImporter{StateContext: schema.ImportStatePassthroughContext},

		Schema: map[string]*schema.Schema{
			"host_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The host id. Must be 3-32 chars, start with a letter.",
			},
			"db_vendor": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Default:     "mssql",
				Description: "Database vendor: mssql or oracledb.",
			},
			"sdp_id": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "Optional SDP id to associate with the host.",
			},
			// computed read-backs
			"host_name":         {Type: schema.TypeString, Computed: true},
			"is_connected":      {Type: schema.TypeBool, Computed: true},
			"agent_version":     {Type: schema.TypeString, Computed: true},
			"db_engine_version": {Type: schema.TypeString, Computed: true},
			"token": {
				Type:        schema.TypeString,
				Computed:    true,
				Sensitive:   true,
				Description: "Agent token returned at host creation.",
			},
		},
	}
}

// older create didnt stash the agent token, only found out it's create-only when
// a later read came back blank for it. set it before the read now
// func resourceFlexEchoHostCreate(d *schema.ResourceData, m interface{}) error {
// 	client := m.(*sdk.Credentials)
// 	hostID := d.Get("host_id").(string)
// 	req := sdk.CreateHostRequest{DBVendor: d.Get("db_vendor").(string)}
// 	if _, err := client.CreateHost(hostID, req); err != nil {
// 		return err
// 	}
// 	d.SetId(hostID)
// 	return resourceFlexEchoHostRead(d, m)
// }

func resourceFlexEchoHostCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(*sdk.Credentials)
	hostID := d.Get("host_id").(string)

	req := sdk.CreateHostRequest{
		DBVendor: d.Get("db_vendor").(string),
		SDPID:    d.Get("sdp_id").(string),
	}
	resp, err := client.CreateHost(hostID, req)
	if err != nil {
		return err
	}

	d.SetId(hostID)
	// token only comes back on the create, STash it now or its gone
	if resp != nil && resp.Token != "" {
		d.Set("token", resp.Token)
	}
	return resourceFlexEchoHostRead(d, m)
}

func resourceFlexEchoHostRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*sdk.Credentials)
	h, err := client.GetHost(d.Id())
	if err != nil {
		return err
	}
	// gone, drop it from state
	// if h == nil {
	if h == nil || h.HostID == "" {
		d.SetId("")
		return nil
	}
	d.Set("host_id", h.HostID)
	d.Set("db_vendor", h.DBVendor)
	d.Set("host_name", h.HostName)
	d.Set("is_connected", h.IsConnected)
	d.Set("agent_version", h.AgentVersion)
	d.Set("db_engine_version", h.DBEngineVersion)
	return nil
}

func resourceFlexEchoHostDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(*sdk.Credentials)
	if err := client.DeleteHost(d.Id()); err != nil {
		return err
	}
	d.SetId("")
	return nil
}
