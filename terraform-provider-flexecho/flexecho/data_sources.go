package flexecho

import (
	"encoding/json"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	sdk "github.com/silk-us/flexecho-go-sdk/flexecho"
)

// single host
func dataSourceFlexEchoHost() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceFlexEchoHostRead,
		Schema: map[string]*schema.Schema{
			"host_id":           {Type: schema.TypeString, Required: true},
			"db_vendor":         {Type: schema.TypeString, Computed: true},
			"host_name":         {Type: schema.TypeString, Computed: true},
			"is_connected":      {Type: schema.TypeBool, Computed: true},
			"agent_version":     {Type: schema.TypeString, Computed: true},
			"db_engine_version": {Type: schema.TypeString, Computed: true},
		},
	}
}

func dataSourceFlexEchoHostRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*sdk.Credentials)
	h, err := client.GetHost(d.Get("host_id").(string))
	if err != nil {
		return err
	}

	d.SetId(h.HostID)
	d.Set("db_vendor", h.DBVendor)
	d.Set("host_name", h.HostName)
	d.Set("is_connected", h.IsConnected)
	d.Set("agent_version", h.AgentVersion)
	d.Set("db_engine_version", h.DBEngineVersion)
	return nil
}

// all hosts
func dataSourceFlexEchoHosts() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceFlexEchoHostsRead,
		Schema: map[string]*schema.Schema{
			"hosts": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"host_id":      {Type: schema.TypeString, Computed: true},
						"db_vendor":    {Type: schema.TypeString, Computed: true},
						"host_name":    {Type: schema.TypeString, Computed: true},
						"is_connected": {Type: schema.TypeBool, Computed: true},
					},
				},
			},
		},
	}
}

// old flatten just appended the structs and let tf choke on the go field names.
// switched to explicit maps keyed by the wire names
// func dataSourceFlexEchoHostsRead(d *schema.ResourceData, m interface{}) error {
// 	client := m.(*sdk.Credentials)
// 	hosts, err := client.GetHosts()
// 	if err != nil {
// 		return err
// 	}
// 	out := []interface{}{}
// 	for _, h := range hosts {
// 		out = append(out, h)
// 	}
// 	d.SetId("flexecho_hosts")
// 	return d.Set("hosts", out)
// }

func dataSourceFlexEchoHostsRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*sdk.Credentials)
	hosts, err := client.GetHosts()
	if err != nil {
		return err
	}
	// flatten each host down to the handful of feilds we expose
	out := make([]map[string]interface{}, 0, len(hosts))
	for _, h := range hosts {
		out = append(out, map[string]interface{}{
			"host_id":      h.HostID,
			"db_vendor":    h.DBVendor,
			"host_name":    h.HostName,
			"is_connected": h.IsConnected,
		})
	}

	d.SetId("flexecho_hosts")
	return d.Set("hosts", out)
}

// db snapshots
func dataSourceFlexEchoDBSnapshots() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceFlexEchoDBSnapshotsRead,
		Schema: map[string]*schema.Schema{
			"host_id": {Type: schema.TypeString, Optional: true},
			"snapshots": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id":                {Type: schema.TypeString, Computed: true},
						"host_id":           {Type: schema.TypeString, Computed: true},
						"host_name":         {Type: schema.TypeString, Computed: true},
						"consistency_level": {Type: schema.TypeString, Computed: true},
						"timestamp":         {Type: schema.TypeInt, Computed: true},
					},
				},
			},
		},
	}
}

func dataSourceFlexEchoDBSnapshotsRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*sdk.Credentials)
	snaps, err := client.GetDBSnapshots()
	if err != nil {
		return err
	}
	// optional host_id narrows the list. empty = return em all
	filterHost := d.Get("host_id").(string)
	out := make([]map[string]interface{}, 0, len(snaps))
	for _, s := range snaps {
		// if filterHost == s.HostID {
		if filterHost != "" && s.HostID != filterHost {
			continue
		}
		out = append(out, map[string]interface{}{
			"id":                s.ID,
			"host_id":           s.HostID,
			"host_name":         s.HostName,
			"consistency_level": s.ConsistencyLevel,
			"timestamp":         s.Timestamp,
		})
	}
	d.SetId("flexecho_db_snapshots")
	return d.Set("snapshots", out)
}

// topology, handed back as a json string since the nested shape is big and loose
func dataSourceFlexEchoTopology() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceFlexEchoTopologyRead,
		Schema: map[string]*schema.Schema{
			"json": {Type: schema.TypeString, Computed: true},
		},
	}
}

func dataSourceFlexEchoTopologyRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*sdk.Credentials)
	topo, err := client.GetTopology()
	if err != nil {
		return err
	}
	// b, err := json.MarshalIndent(topo, "", "  ")
	b, err := json.Marshal(topo)
	if err != nil {
		return err
	}
	d.SetId("flexecho_topology")
	return d.Set("json", string(b))
}
