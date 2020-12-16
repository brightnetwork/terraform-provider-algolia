# terraform-provider-algolia

A work-in-progress terraform provider for managing Algolia search indexes. Currently only supports synonyms (regular & one-way) and query rules.

## Syntax

```tf
provider "algolia" {
  api_key        = "<SNIP>"
  application_id = "<SNIP>"
}

resource "algolia_regular_synonym" "law" {
  index    = "CommonIndex_production"
  synonyms = ["law", "commercial law", "legal"]
}

resource "algolia_one_way_synonym" "accounting" {
  index    = "CommonIndex_production"
  input    = "accounting"
  synonyms = ["Accounting, Tax & Audit"]
}

resource "algolia_rule" "test" {
  index = "dan_test"

  condition {
    pattern   = "success story"
    anchoring = "contains"
  }

  consequence = <<EOF
{
    "params": {
        "query": {
            "edits": [
                {
                    "type": "remove",
                    "delete": "success story"
                }
            ]
        },
        "filters": "_objectModel:web.models.SuccessStory"
    },
    "filterPromotes": true
}
EOF
}
```

## Import

Resources can be imported using the following commands:

```sh
$ terraform import algolia_regular_synonym.foo index_name:syn-123456789-0
$ terraform import algolia_one_way_synonym.bar index_name:syn-123456789-1
$ terraform import algolia_rule.baz index_name:qr-123456789-0
```
