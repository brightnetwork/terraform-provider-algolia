package main

import (
	"github.com/algolia/algoliasearch-client-go/v3/algolia/search"
	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceRegularSynonymCreate(d *schema.ResourceData, m interface{}) error {
	client := *m.(*search.Client)
	index := client.InitIndex(d.Get("index").(string))

	id := uuid.New().String()
	synonym := search.NewRegularSynonym(id, castStringList(d.Get("synonyms").([]interface{}))...)
	res, err := index.SaveSynonym(synonym)
	if err != nil {
		return err
	}
	res.Wait()
	d.SetId(synonym.ObjectID())
	return resourceRegularSynonymRead(d, m)
}

func resourceRegularSynonymRead(d *schema.ResourceData, m interface{}) error {
	client := *m.(*search.Client)
	index := client.InitIndex(d.Get("index").(string))

	synonym, err := index.GetSynonym(d.Id())
	if err != nil {
		d.SetId("")
		return nil
	}
	d.Set("synonyms", synonym.(search.RegularSynonym).Synonyms)
	return nil
}

func resourceRegularSynonymUpdate(d *schema.ResourceData, m interface{}) error {
	client := *m.(*search.Client)
	index := client.InitIndex(d.Get("index").(string))

	synonym := search.NewRegularSynonym(d.Id(), castStringList(d.Get("synonyms").([]interface{}))...)
	res, err := index.SaveSynonym(synonym)
	if err != nil {
		return err
	}
	res.Wait()
	return resourceRegularSynonymRead(d, m)
}

func resourceRegularSynonymDelete(d *schema.ResourceData, m interface{}) error {
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

func resourceRegularSynonym() *schema.Resource {
	return &schema.Resource{
		Create: resourceRegularSynonymCreate,
		Read:   resourceRegularSynonymRead,
		Update: resourceRegularSynonymUpdate,
		Delete: resourceRegularSynonymDelete,

		Schema: map[string]*schema.Schema{
			"index": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Algolia Index",
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
