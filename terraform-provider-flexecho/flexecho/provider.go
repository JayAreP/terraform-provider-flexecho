package flexecho

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// the flexecho provider. creds up top, then the resources / datasources hang off it
func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"server": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("SILK_FLEX_SERVER", nil),
				Description: "IP address or hostname of the Silk Flex management console.",
			},
			"token": {
				Type:        schema.TypeString,
				Required:    true,
				Sensitive:   true,
				DefaultFunc: schema.EnvDefaultFunc("SILK_FLEX_TOKEN", nil),
				Description: "Bearer token used to authenticate against the Flex API.",
			},
		},

		ResourcesMap: map[string]*schema.Resource{
			"flexecho_host":        resourceFlexEchoHost(),
			"flexecho_db_snapshot": resourceFlexEchoDBSnapshot(),
			"flexecho_echo_db":     resourceFlexEchoEchoDB(),
			"flexecho_app_token":   resourceFlexEchoAppToken(),
		},

		DataSourcesMap: map[string]*schema.Resource{
			"flexecho_host":          dataSourceFlexEchoHost(),
			"flexecho_hosts":         dataSourceFlexEchoHosts(),
			"flexecho_db_snapshots":  dataSourceFlexEchoDBSnapshots(),
			"flexecho_topology":      dataSourceFlexEchoTopology(),
		},

		ConfigureFunc: providerConfigure,
	}
}

// early version new'd the sdk client right here. moved it behind Config.Client()
// so the provider doesnt import the sdk directly
// func providerConfigure(d *schema.ResourceData) (interface{}, error) {
// 	server := d.Get("server").(string)
// 	token := d.Get("token").(string)
// 	return sdk.Connect(server, token), nil
// }

// pull server + token off the provider block and hand back a built sdk client
func providerConfigure(d *schema.ResourceData) (interface{}, error) {
	config := Config{
		Server: d.Get("server").(string),
		Token:  d.Get("token").(string),
	}
	return config.Client()
}
