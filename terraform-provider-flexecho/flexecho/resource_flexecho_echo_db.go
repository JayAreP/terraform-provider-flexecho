package flexecho

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	sdk "github.com/silk-us/flexecho-go-sdk/flexecho"
)

// an echo-db clone. two create paths:
//   - from a snapshot (snapshot_id set)  -> POST /db_snapshots/{id}/echo_db
//   - standalone replicate (no snapshot) -> POST /echo_dbs
// delete is always DELETE /echo_dbs with {host_id, database_id}
func resourceFlexEchoEchoDB() *schema.Resource {
	return &schema.Resource{
		Create: resourceFlexEchoEchoDBCreate,
		Read:   resourceFlexEchoEchoDBRead,
		Delete: resourceFlexEchoEchoDBDelete,

		Schema: map[string]*schema.Schema{
			"snapshot_id": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "If set, clone FROM this snapshot. Otherwise a standalone replicate is performed.",
			},
			"source_host_id": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "Source host id (required for the standalone replicate path).",
			},
			"destination_host_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"destination_db_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"destination_db_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"target_state": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Default:     "online",
				Description: "recovery or online.",
			},
		},
	}
}

// started with only the from-snapshot path, no standalone replicate. added the
// branch later once i realized echo_db can clone without a snapshot too
// func resourceFlexEchoEchoDBCreate(d *schema.ResourceData, m interface{}) error {
// 	client := m.(*sdk.Credentials)
// 	dest := sdk.DestinationDB{
// 		HostID: d.Get("destination_host_id").(string),
// 		DBID:   d.Get("destination_db_id").(string),
// 		DBName: d.Get("destination_db_name").(string),
// 	}
// 	req := sdk.DeployRequestBody{
// 		Destinations: []sdk.DestinationDB{dest},
// 		TargetState:  d.Get("target_state").(string),
// 	}
// 	task, err := client.CreateEchoDBFromSnapshot(d.Get("snapshot_id").(string), req)
// 	if err != nil {
// 		return err
// 	}
// 	client.WaitForTask(task.RefID, 5, 1800)
// 	d.SetId(dest.DBID)
// 	return nil
// }

func resourceFlexEchoEchoDBCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(*sdk.Credentials)

	dest := sdk.DestinationDB{
		HostID: d.Get("destination_host_id").(string),
		DBID:   d.Get("destination_db_id").(string),
		DBName: d.Get("destination_db_name").(string),
	}

	var task *sdk.TaskStatusResponse
	var err error

	// snapshot_id set means clone from a snap, otherwise standalone replicate
	if snapID := d.Get("snapshot_id").(string); snapID != "" {
		req := sdk.DeployRequestBody{
			Destinations: []sdk.DestinationDB{dest},
			TargetState:  d.Get("target_state").(string),
		}
		task, err = client.CreateEchoDBFromSnapshot(snapID, req)
	} else {
		req := sdk.ReplicateRequest{
			SourceHostID: d.Get("source_host_id").(string),
			Destinations: []sdk.DestinationDB{dest},
			TargetState:  d.Get("target_state").(string),
		}
		task, err = client.CreateEchoDB(req)
	}
	if err != nil {
		return err
	}

	pollID := task.RequestID
	if pollID == "" {
		pollID = task.RefID
	}
	if _, err := client.WaitForTask(pollID, 5, 1800); err != nil {
		return err
	}

	// composite id of host + db so delete (needs both) can pull them back apart
	// d.SetId(dest.DBID)
	d.SetId(fmt.Sprintf("%s/%s", dest.HostID, dest.DBID))
	return resourceFlexEchoEchoDBRead(d, m)
}

func resourceFlexEchoEchoDBRead(d *schema.ResourceData, m interface{}) error {
	// no get-by-id for echo dbs in the spec, the create-time attrs are the source
	// of truth. nothing to refresh remotely without a topology scan
	// TODO(verify-live): maybe confirm existence via GET /topology databases
	return nil
}

func resourceFlexEchoEchoDBDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(*sdk.Credentials)
	req := sdk.DeleteRequest{
		HostID:     d.Get("destination_host_id").(string),
		DatabaseID: d.Get("destination_db_id").(string),
	}
	task, err := client.DeleteEchoDB(req)
	if err != nil {
		return err
	}
	pollID := task.RequestID
	if pollID == "" {
		pollID = task.RefID
	}
	if _, err := client.WaitForTask(pollID, 5, 1800); err != nil {
		return err
	}
	d.SetId("")
	return nil
}
