package provider

import (
	"fmt"
	"strings"

	"github.com/algolia/algoliasearch-client-go/v3/algolia/opt"
	"github.com/algolia/algoliasearch-client-go/v3/algolia/search"
	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceOneWaySynonymCreate(d *schema.ResourceData, m interface{}) error {
	client := *m.(*search.Client)
	index := client.InitIndex(d.Get("index").(string))

	synonym := search.NewOneWaySynonym(
		uuid.New().String(),
		d.Get("input").(string),
		castStringList(d.Get("synonyms").([]interface{}))...,
	)

	res, err := index.SaveSynonym(synonym, opt.ForwardToReplicas(d.Get("forward_to_replicas").(bool)))
	if err != nil {
		return err
	}
	res.Wait()
	d.SetId(synonym.ObjectID())
	return resourceOneWaySynonymRead(d, m)
}

func resourceOneWaySynonymRead(d *schema.ResourceData, m interface{}) error {
	client := *m.(*search.Client)
	index := client.InitIndex(d.Get("index").(string))

	synonym, err := index.GetSynonym(d.Id())
	if err != nil {
		d.SetId("")
		return nil
	}
	d.Set("input", synonym.(search.OneWaySynonym).Input)
	d.Set("synonyms", synonym.(search.OneWaySynonym).Synonyms)
	return nil
}

func resourceOneWaySynonymUpdate(d *schema.ResourceData, m interface{}) error {
	client := *m.(*search.Client)
	index := client.InitIndex(d.Get("index").(string))

	synonym := search.NewOneWaySynonym(
		d.Id(),
		d.Get("input").(string),
		castStringList(d.Get("synonyms").([]interface{}))...,
	)
	res, err := index.SaveSynonym(synonym, opt.ForwardToReplicas(d.Get("forward_to_replicas").(bool)))
	if err != nil {
		return err
	}
	res.Wait()
	return resourceOneWaySynonymRead(d, m)
}

func resourceOneWaySynonymDelete(d *schema.ResourceData, m interface{}) error {
	client := *m.(*search.Client)
	index := client.InitIndex(d.Get("index").(string))

	id := d.Id()
	res, err := index.DeleteSynonym(id, opt.ForwardToReplicas(d.Get("forward_to_replicas").(bool)))
	if err != nil {
		return err
	}
	res.Wait()
	return nil
}

func resourceOneWaySynonym() *schema.Resource {
	return &schema.Resource{
		Create: resourceOneWaySynonymCreate,
		Read:   resourceOneWaySynonymRead,
		Update: resourceOneWaySynonymUpdate,
		Delete: resourceOneWaySynonymDelete,
		Importer: &schema.ResourceImporter{
			State: func(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
				// d.Id() here is the last argument passed to the `terraform import RESOURCE_TYPE.RESOURCE_NAME RESOURCE_ID` command
				parts := strings.SplitN(d.Id(), ":", 2)

				if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
					return nil, fmt.Errorf("unexpected format of ID (%s), expected index:id", d.Id())
				}

				d.Set("index", parts[0])
				d.SetId(parts[1])

				return []*schema.ResourceData{d}, nil
			},
		},

		Schema: map[string]*schema.Schema{
			"index": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Algolia Index",
			},
			"input": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Search Term",
			},
			"forward_to_replicas": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			"synonyms": {
				Type:        schema.TypeList,
				Required:    true,
				Description: "List of synonyms",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}
