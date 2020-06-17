package main

import (
	"encoding/json"

	"github.com/algolia/algoliasearch-client-go/v3/algolia/opt"
	"github.com/algolia/algoliasearch-client-go/v3/algolia/search"
	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/structure"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func resourceRuleCreate(d *schema.ResourceData, m interface{}) error {
	client := *m.(*search.Client)
	index := client.InitIndex(d.Get("index").(string))

	conditionsList := d.Get("condition").([]interface{})
	conditionMap := conditionsList[0].(map[string]interface{})
	condition := search.RuleCondition{
		Anchoring:    search.RulePatternAnchoring(conditionMap["anchoring"].(string)),
		Pattern:      conditionMap["pattern"].(string),
		Alternatives: search.AlternativesEnabled(),
	}
	var consequence search.RuleConsequence
	consequenceJSON := []byte(d.Get("consequence").(string))
	json.Unmarshal(consequenceJSON, &consequence)

	rule := search.Rule{
		ObjectID:    uuid.New().String(),
		Condition:   condition,
		Consequence: consequence,
		Enabled:     opt.Enabled(d.Get("enabled").(bool)),
	}

	res, err := index.SaveRule(rule)
	if err != nil {
		return err
	}
	res.Wait()
	d.SetId(rule.ObjectID)
	return resourceRuleRead(d, m)
}

func flattenCondition(in search.RuleCondition) []interface{} {
	m := make(map[string]interface{})
	m["anchoring"] = in.Anchoring
	m["pattern"] = in.Pattern
	return []interface{}{m}
}

func flattenConsequence(in search.RuleConsequence) string {
	// Just using JSON for now
	consequenceJSON, _ := json.Marshal(in)
	return string(consequenceJSON)
}

func resourceRuleRead(d *schema.ResourceData, m interface{}) error {
	client := *m.(*search.Client)
	index := client.InitIndex(d.Get("index").(string))
	rule, err := index.GetRule(d.Id())
	if err != nil {
		d.SetId("")
		return nil
	}
	d.Set("enabled", rule.Enabled.Get())
	d.Set("condition", flattenCondition(rule.Condition))
	d.Set("consequence", flattenConsequence(rule.Consequence))

	return nil
}

func resourceRuleUpdate(d *schema.ResourceData, m interface{}) error {
	client := *m.(*search.Client)
	index := client.InitIndex(d.Get("index").(string))

	conditionsList := d.Get("condition").([]interface{})
	conditionMap := conditionsList[0].(map[string]interface{})
	condition := search.RuleCondition{
		Anchoring:    search.RulePatternAnchoring(conditionMap["anchoring"].(string)),
		Pattern:      conditionMap["pattern"].(string),
		Alternatives: search.AlternativesEnabled(),
	}
	var consequence search.RuleConsequence
	consequenceJSON := []byte(d.Get("consequence").(string))
	json.Unmarshal(consequenceJSON, &consequence)

	rule := search.Rule{
		ObjectID:    d.Id(),
		Condition:   condition,
		Consequence: consequence,
		Enabled:     opt.Enabled(d.Get("enabled").(bool)),
	}

	res, err := index.SaveRule(rule)
	if err != nil {
		return err
	}
	res.Wait()

	return resourceRuleRead(d, m)
}

func resourceRuleDelete(d *schema.ResourceData, m interface{}) error {
	client := *m.(*search.Client)
	index := client.InitIndex(d.Get("index").(string))

	id := d.Id()
	res, err := index.DeleteRule(id)
	if err != nil {
		return err
	}
	res.Wait()
	return nil
}

func resourceRule() *schema.Resource {
	return &schema.Resource{
		Create: resourceRuleCreate,
		Read:   resourceRuleRead,
		Update: resourceRuleUpdate,
		Delete: resourceRuleDelete,

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
					},
				},
			},
			"consequence": {
				Type:             schema.TypeString,
				Required:         true,
				ValidateFunc:     validation.StringIsJSON,
				DiffSuppressFunc: structure.SuppressJsonDiff,
			},
		},
	}
}
