package flexecho

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	sdk "github.com/silk-us/flexecho-go-sdk/flexecho"
)

func resourceFlexEchoDBSnapshot() *schema.Resource {
	return &schema.Resource{
		Create: resourceFlexEchoDBSnapshotCreate,
		Read:   resourceFlexEchoDBSnapshotRead,
		Delete: resourceFlexEchoDBSnapshotDelete,
		Importer: &schema.ResourceImporter{StateContext: schema.ImportStatePassthroughContext},

		Schema: map[string]*schema.Schema{
			"database_ids": {
				Type:        schema.TypeList,
				Required:    true,
				ForceNew:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "Database ids to capture in the snapshot.",
			},
			"source_host_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Host id the databases live on.",
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
			// computed
			"host_name": {Type: schema.TypeString, Computed: true},
			"timestamp": {Type: schema.TypeInt, Computed: true},
		},
	}
}

// first version fired the capture and set the id straight off task.RefID without
// waiting. tf came back and the snap wasnt queryable yet, the read blew up.
// now it polls the task then resolves the id. keeping the old one around
// func resourceFlexEchoDBSnapshotCreate(d *schema.ResourceData, m interface{}) error {
// 	client := m.(*sdk.Credentials)
// 	ids := []string{}
// 	for _, v := range d.Get("database_ids").([]interface{}) {
// 		ids = append(ids, v.(string))
// 	}
// 	req := sdk.CaptureRequest{
// 		DatabaseIDs:  ids,
// 		SourceHostID: d.Get("source_host_id").(string),
// 		NamePrefix:   d.Get("name_prefix").(string),
// 	}
// 	task, err := client.CreateDBSnapshot(req)
// 	if err != nil {
// 		return err
// 	}
// 	d.SetId(task.RefID)
// 	return resourceFlexEchoDBSnapshotRead(d, m)
// }

func resourceFlexEchoDBSnapshotCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(*sdk.Credentials)

	// flatten the list into plain strings
	ids := []string{}
	for _, v := range d.Get("database_ids").([]interface{}) {
		ids = append(ids, v.(string))
	}

	req := sdk.CaptureRequest{
		DatabaseIDs:      ids,
		SourceHostID:     d.Get("source_host_id").(string),
		NamePrefix:       d.Get("name_prefix").(string),
		ConsistencyLevel: d.Get("consistency_level").(string),
		UseVSS:           d.Get("use_vss").(bool),
	}

	task, err := client.CreateDBSnapshot(req)
	if err != nil {
		return err
	}

	// the create 202 only gives a ref_id; request_id is NOT the queryable task id
	// (the real one is the task's taskid). so wait on the task by matching ref_id
	if _, err := client.WaitForTaskByRef(task.RefID, 5, 1800); err != nil {
		return err
	}

	// ref_id is the task ref, NOT the snapshot id, so dont trust it for the id.
	// resolve the new snap off the snapshot list instead
	snapID, err := resolveLatestSnapshotForHost(client, req.SourceHostID)
	if err != nil {
		return err
	}
	if snapID == "" {
		return fmt.Errorf("snapshot task (ref %s) finished but no snapshot turned up for host %s", task.RefID, req.SourceHostID)
	}
	d.SetId(snapID)
	return resourceFlexEchoDBSnapshotRead(d, m)
}

func resourceFlexEchoDBSnapshotRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*sdk.Credentials)
	snap, err := client.GetDBSnapshot(d.Id())
	if err != nil {
		return err
	}
	if snap == nil {
		d.SetId("")
		return nil
	}
	d.Set("source_host_id", snap.HostID)
	d.Set("host_name", snap.HostName)
	d.Set("consistency_level", snap.ConsistencyLevel)
	d.Set("timestamp", snap.Timestamp)
	return nil
}

func resourceFlexEchoDBSnapshotDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(*sdk.Credentials)
	task, err := client.DeleteDBSnapshot(d.Id())
	if err != nil {
		return err
	}

	// delete is async like create. wait on the task (matched by ref_id) so a failed
	// delete (eg the snap still has dependent echo dbs) surfaces, instead of terraform
	// dropping it from state while the array still has it. if delete ever comes back
	// synchronous (no ref_id) we just skip the wait
	if task != nil && task.RefID != "" {
		if _, err := client.WaitForTaskByRef(task.RefID, 5, 1800); err != nil {
			return err
		}
	}

	d.SetId("")
	return nil
}

// task body doesnt hand back the new snap id, so resolve it off the global snap
// list -- same list Read uses, so the readback is guaranteed to find it. newest
// timestamp for the host wins
func resolveLatestSnapshotForHost(client *sdk.Credentials, hostID string) (string, error) {
	// resp, err := client.GetHostDBSnapshotIDs(hostID)
	// return resp.DBSnapshotIDs[len(resp.DBSnapshotIDs)-1], nil // per-host ids didnt line up with the global list
	snaps, err := client.GetDBSnapshots()
	if err != nil {
		return "", err
	}

	newestID := ""
	newestTS := -1
	for _, s := range snaps {
		// skip other hosts, but only if the list actually populates host_id
		if hostID != "" && s.HostID != "" && s.HostID != hostID {
			continue
		}
		if s.Timestamp >= newestTS {
			newestTS = s.Timestamp
			newestID = s.ID
		}
	}
	return newestID, nil
}
