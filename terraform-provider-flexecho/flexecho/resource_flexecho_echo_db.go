package flexecho

import (
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	sdk "github.com/silk-us/flexecho-go-sdk/flexecho"
)

// old meta - snapshot_id + a single destination_host/db/name, modeled on the wrong
// request. the real /echo_dbs body is a ReplicateRequest (source + database_ids + a
// destinations array). kept for reference:
// func resourceFlexEchoEchoDB() *schema.Resource {
// 	return &schema.Resource{
// 		Create: resourceFlexEchoEchoDBCreate,
// 		Read:   resourceFlexEchoEchoDBRead,
// 		Delete: resourceFlexEchoEchoDBDelete,
// 		Schema: map[string]*schema.Schema{
// 			"snapshot_id": {
// 				Type:        schema.TypeString,
// 				Optional:    true,
// 				ForceNew:    true,
// 				Description: "If set, clone FROM this snapshot. Otherwise a standalone replicate is performed.",
// 			},
// 			"source_host_id": {
// 				Type:        schema.TypeString,
// 				Optional:    true,
// 				ForceNew:    true,
// 				Description: "Source host id (required for the standalone replicate path).",
// 			},
// 			"destination_host_id": {Type: schema.TypeString, Required: true, ForceNew: true},
// 			"destination_db_id":   {Type: schema.TypeString, Required: true, ForceNew: true},
// 			"destination_db_name": {Type: schema.TypeString, Required: true, ForceNew: true},
// 			"target_state": {
// 				Type:        schema.TypeString,
// 				Optional:    true,
// 				ForceNew:    true,
// 				Default:     "online",
// 				Description: "recovery or online.",
// 			},
// 		},
// 	}
// }

// an echo-db clone (standalone replicate, ReplicateRequest). fans a source host's dbs
// out to one or more destination hosts. delete tears down each destination clone
// (DELETE /echo_dbs is per host+db)
func resourceFlexEchoEchoDB() *schema.Resource {
	return &schema.Resource{
		Create: resourceFlexEchoEchoDBCreate,
		Read:   resourceFlexEchoEchoDBRead,
		Delete: resourceFlexEchoEchoDBDelete,

		Schema: map[string]*schema.Schema{
			"source_host_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Source host id the databases live on.",
			},
			"database_ids": {
				Type:        schema.TypeList,
				Required:    true,
				ForceNew:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "Database ids to replicate from the source host.",
			},
			"destination": {
				Type:        schema.TypeList,
				Required:    true,
				ForceNew:    true,
				MinItems:    1,
				Description: "One or more clone destinations. Repeat the block per destination.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"host_id": {Type: schema.TypeString, Required: true},
						"db_id":   {Type: schema.TypeString, Required: true},
						"db_name": {Type: schema.TypeString, Required: true},
					},
				},
			},
			"name_prefix": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Default:     "snap",
				Description: "Snapshot name prefix (4-20 chars, ^[a-z][a-z0-9_-]+$).",
			},
			"consistency_level": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Default:     "application",
				Description: "crash or application.",
			},
			"use_vss": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: true,
				Default:  false,
			},
			"target_state": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Default:     "online",
				Description: "recovery or online.",
			},
			"backup_session_timeout_sec": {
				Type:     schema.TypeInt,
				Optional: true,
				ForceNew: true,
			},
			"restore_session_timeout_sec": {
				Type:     schema.TypeInt,
				Optional: true,
				ForceNew: true,
			},
		},
	}
}

// older single-destination / from-snapshot create. replaced by the ReplicateRequest
// build below. kept for reference:
// func resourceFlexEchoEchoDBCreate(d *schema.ResourceData, m interface{}) error {
// 	client := m.(*sdk.Credentials)
// 	dest := sdk.DestinationDB{
// 		HostID: d.Get("destination_host_id").(string),
// 		DBID:   d.Get("destination_db_id").(string),
// 		DBName: d.Get("destination_db_name").(string),
// 	}
// 	var task *sdk.TaskStatusResponse
// 	var err error
// 	if snapID := d.Get("snapshot_id").(string); snapID != "" {
// 		req := sdk.DeployRequestBody{Destinations: []sdk.DestinationDB{dest}, TargetState: d.Get("target_state").(string)}
// 		task, err = client.CreateEchoDBFromSnapshot(snapID, req)
// 	} else {
// 		req := sdk.ReplicateRequest{SourceHostID: d.Get("source_host_id").(string), Destinations: []sdk.DestinationDB{dest}, TargetState: d.Get("target_state").(string)}
// 		task, err = client.CreateEchoDB(req)
// 	}
// 	if err != nil {
// 		return err
// 	}
// 	if _, err := client.WaitForTaskByRef(task.RefID, 5, 1800); err != nil {
// 		return err
// 	}
// 	d.SetId(fmt.Sprintf("%s/%s", dest.HostID, dest.DBID))
// 	return resourceFlexEchoEchoDBRead(d, m)
// }

func resourceFlexEchoEchoDBCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(*sdk.Credentials)

	// database_ids -> []string
	dbIDs := []string{}
	for _, v := range d.Get("database_ids").([]interface{}) {
		dbIDs = append(dbIDs, v.(string))
	}

	// destination blocks -> []DestinationDB
	dests := []sdk.DestinationDB{}
	for _, raw := range d.Get("destination").([]interface{}) {
		dm := raw.(map[string]interface{})
		dests = append(dests, sdk.DestinationDB{
			HostID: dm["host_id"].(string),
			DBID:   dm["db_id"].(string),
			DBName: dm["db_name"].(string),
		})
	}

	req := sdk.ReplicateRequest{
		SourceHostID:     d.Get("source_host_id").(string),
		DatabaseIDs:      dbIDs,
		Destinations:     dests,
		NamePrefix:       d.Get("name_prefix").(string),
		ConsistencyLevel: d.Get("consistency_level").(string),
		UseVSS:           d.Get("use_vss").(bool),
		TargetState:      d.Get("target_state").(string),
	}
	// timeouts are optional -- only send when set (>0)
	if v := d.Get("backup_session_timeout_sec").(int); v > 0 {
		req.BackupSessionTimeoutSec = &v
	}
	if v := d.Get("restore_session_timeout_sec").(int); v > 0 {
		req.RestoreSessionTimeoutSec = &v
	}

	task, err := client.CreateEchoDB(req)
	if err != nil {
		return err
	}
	if _, err := client.WaitForTaskByRef(task.RefID, 5, 1800); err != nil {
		return err
	}

	// no single backend id for a replicate -- key the resource off source + the dests
	keys := []string{req.SourceHostID}
	for _, dst := range dests {
		keys = append(keys, dst.HostID+":"+dst.DBID)
	}
	d.SetId(strings.Join(keys, "/"))
	return resourceFlexEchoEchoDBRead(d, m)
}

func resourceFlexEchoEchoDBRead(d *schema.ResourceData, m interface{}) error {
	// no get-by-id for echo dbs in the spec, the create-time attrs are the source
	// of truth. nothing to refresh remotely without a topology scan
	// TODO(verify-live): maybe confirm existence via GET /topology databases
	return nil
}

// older single-db delete. now iterates the destinations:
// func resourceFlexEchoEchoDBDelete(d *schema.ResourceData, m interface{}) error {
// 	client := m.(*sdk.Credentials)
// 	req := sdk.DeleteRequest{
// 		HostID:     d.Get("destination_host_id").(string),
// 		DatabaseID: d.Get("destination_db_id").(string),
// 	}
// 	task, err := client.DeleteEchoDB(req)
// 	if err != nil {
// 		return err
// 	}
// 	if _, err := client.WaitForTaskByRef(task.RefID, 5, 1800); err != nil {
// 		return err
// 	}
// 	d.SetId("")
// 	return nil
// }

func resourceFlexEchoEchoDBDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(*sdk.Credentials)

	// echo_dbs delete is per host+db, so tear down each destination clone
	// TODO(verify-live): confirm DeleteRequest.database_id wants the dest db_id (not db_name)
	for _, raw := range d.Get("destination").([]interface{}) {
		dm := raw.(map[string]interface{})
		req := sdk.DeleteRequest{
			HostID:     dm["host_id"].(string),
			DatabaseID: dm["db_id"].(string),
		}
		task, err := client.DeleteEchoDB(req)
		if err != nil {
			return err
		}
		if _, err := client.WaitForTaskByRef(task.RefID, 5, 1800); err != nil {
			return err
		}
	}

	d.SetId("")
	return nil
}
