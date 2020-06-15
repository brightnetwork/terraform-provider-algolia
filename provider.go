package main

import (
	"github.com/algolia/algoliasearch-client-go/v3/algolia/search"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"application_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Algolia Application ID",
			},
			"api_key": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Algolia API Key",
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"algolia_synonym": resourceSynonym(),
		},
		ConfigureFunc: providerConfigure,
	}
}

func providerConfigure(data *schema.ResourceData) (interface{}, error) {
	applicationID := data.Get("application_id").(string)
	apiKey := data.Get("api_key").(string)
	client := search.NewClient(applicationID, apiKey)
	return client, nil
}
