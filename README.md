# terraform-provider-algolia

A custom terraform provider for managing our Algolia search indexes. Currently only supports synonyms (regular & one-way).

The hope is that we'll open-source this when it's more mature and fully-featured.

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
```

## Import

Resources can be imported using the following commands:

```sh
$ terraform import algolia_regular_synonym.foo index_name:syn-123456789-0
$ terraform import algolia_one_way_synonym.bar index_name:syn-123456789-1
```

## Local Development

Create a test terraform config in the root of the project (make sure to specify a test index).

The provider can be built and reloaded using:

```sh
go build -o terraform-provider-algolia && terraform init
```

## Building for Terraform Cloud

Until Terraform 0.13 is released[1], it's tricky to use custom providers with Terraform Cloud. For now, we have to (cross-)compile the provider for linux/amd64 using `make build` and commit the resulting binary along with the code changes. Our `terraform-algolia` repo then includes this project as a git submodule.

## Known Issues

* Terraform import is currently broken when using Terraform Cloud [2].


## References

[1] https://www.hashicorp.com/blog/announcing-providers-in-the-new-terraform-registry/

[2] https://discuss.hashicorp.com/t/how-do-i-pass-sensitive-variable-values-when-running-a-command-locally/7266/2