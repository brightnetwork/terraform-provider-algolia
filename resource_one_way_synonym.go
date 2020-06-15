package main

import (
	"github.com/algolia/algoliasearch-client-go/v3/algolia/search"
	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceOneWaySynonymCreate(d *schema.ResourceData, m interface{}) error {
	client := *m.(*search.Client)
	index := client.InitIndex(d.Get("index").(string))

	synonym := search.NewOneWaySynonym(
		uuid.New().String(),
		d.Get("input").(string),
		castStringList(d.Get("synonyms").([]interface{}))...,
	)

	res, err := index.SaveSynonym(synonym)
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
	res, err := index.SaveSynonym(synonym)
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
	res, err := index.DeleteSynonym(id)
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

		Schema: map[string]*schema.Schema{
			"index": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Algolia Index",
			},
			"input": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "Search Term",
			},
			"synonyms": &schema.Schema{
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
