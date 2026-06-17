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

	// async, poll the task til its done
	// TODO(verify-live): confirm which field is the poll id (request_id vs ref_id)
	// and which carrys the new snapshot id once we see a real 202 body
	// pollID := task.RefID
	pollID := task.RequestID
	if pollID == "" {
		pollID = task.RefID
	}
	done, err := client.WaitForTask(pollID, 5, 1800)
	if err != nil {
		return err
	}

	// figure out the new snap id. ref_id is the best guess, else scan the host
	snapID := done.RefID
	if snapID == "" {
		snapID, err = resolveLatestSnapshotForHost(client, req.SourceHostID)
		if err != nil {
			return err
		}
	}
	if snapID == "" {
		return fmt.Errorf("snapshot created but could not resolve its id from task %s", pollID)
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
	if err := client.DeleteDBSnapshot(d.Id()); err != nil {
		return err
	}
	d.SetId("")
	return nil
}

// fallback for when the task body doesnt hand back the new snap id. lists the
// hosts snap ids and grabs the last one. best effort, see the verify-live todo above
func resolveLatestSnapshotForHost(client *sdk.Credentials, hostID string) (string, error) {
	resp, err := client.GetHostDBSnapshotIDs(hostID)
	if err != nil {
		return "", err
	}
	if resp == nil || len(resp.DBSnapshotIDs) == 0 {
		return "", nil
	}
	return resp.DBSnapshotIDs[len(resp.DBSnapshotIDs)-1], nil
}
