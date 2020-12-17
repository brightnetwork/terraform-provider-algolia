---
page_title: "algolia_rule Resource - terraform-provider-algolia"
subcategory: ""
description: |-
  
---

# Resource `algolia_rule`





## Schema

### Required

- **condition** (Block List, Min: 1, Max: 1) (see [below for nested schema](#nestedblock--condition))
- **consequence** (String, Required)
- **index** (String, Required) Algolia Index

### Optional

- **enabled** (Boolean, Optional)
- **id** (String, Optional) The ID of this resource.

<a id="nestedblock--condition"></a>
### Nested Schema for `condition`

Required:

- **anchoring** (String, Required)
- **pattern** (String, Required)

Optional:

- **alternatives** (Boolean, Optional)
- **context** (String, Optional)


