package main

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func resourceRule() *schema.Resource {
	return &schema.Resource{
		// Create: resourceRuleCreate,
		// Read:   resourceRuleRead,
		// Update: resourceRuleUpdate,
		// Delete: resourceRuleDelete,

		Schema: map[string]*schema.Schema{
			"index": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Algolia Index",
			},
			"enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			"condition": {
				Type:     schema.TypeList,
				Required: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"pattern": {
							Type:     schema.TypeString,
							Required: true,
						},
						"anchoring": {
							Type:     schema.TypeString,
							Required: true,
						},
						"alternatives": {
							Type:     schema.TypeBool,
							Required: true,
						},
					},
				},
			},
			"consequence": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringIsJSON,
			},
		},
	}
}
