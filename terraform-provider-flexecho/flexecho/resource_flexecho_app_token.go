package flexecho

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	sdk "github.com/silk-us/flexecho-go-sdk/flexecho"
)

// manages a flex app token (core /api/v2/flex_app_tokens). create + delete, REad walks the list
func resourceFlexEchoAppToken() *schema.Resource {
	return &schema.Resource{
		Create: resourceFlexEchoAppTokenCreate,
		Read:   resourceFlexEchoAppTokenRead,
		Delete: resourceFlexEchoAppTokenDelete,

		Schema: map[string]*schema.Schema{
			"app_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			// computed
			"ttl":       {Type: schema.TypeInt, Computed: true},
			"valid":     {Type: schema.TypeBool, Computed: true},
			"expire_ts": {Type: schema.TypeInt, Computed: true},
		},
	}
}

func resourceFlexEchoAppTokenCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(*sdk.Credentials)
	tok, err := client.CreateAppToken(d.Get("app_name").(string), d.Get("description").(string))
	if err != nil {
		return err
	}
	d.SetId(tok.ID)
	// d.Set("description", tok.Description)
	d.Set("ttl", tok.TTL)
	d.Set("valid", tok.Valid)
	d.Set("expire_ts", tok.ExpireTS)
	return nil
}

// wanted a single GetAppToken(id) here but theres no get-by-id, so list + match
// func resourceFlexEchoAppTokenRead(d *schema.ResourceData, m interface{}) error {
// 	client := m.(*sdk.Credentials)
// 	tok, err := client.GetAppToken(d.Id())
// 	if err != nil {
// 		return err
// 	}
// 	if tok == nil {
// 		d.SetId("")
// 		return nil
// 	}
// 	d.Set("ttl", tok.TTL)
// 	d.Set("valid", tok.Valid)
// 	return nil
// }

func resourceFlexEchoAppTokenRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*sdk.Credentials)
	tokens, err := client.GetAppTokens()
	if err != nil {
		return err
	}
	// no get-by-id so walk the list and match on our stored id
	for _, t := range tokens {
		if t.ID == d.Id() {
			d.Set("description", t.Description)
			d.Set("ttl", t.TTL)
			d.Set("valid", t.Valid)
			d.Set("expire_ts", t.ExpireTS)
			return nil
		}
	}

	d.SetId("") // not in the list anymore, its gone
	return nil
}

func resourceFlexEchoAppTokenDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(*sdk.Credentials)
	if err := client.DeleteAppToken(d.Id()); err != nil {
		return err
	}
	d.SetId("")
	return nil
}
