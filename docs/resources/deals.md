# Resource: pipedrive_deals

## Example Usage

```hcl
resource "pipedrive_deals" "example" {
  title  = "Example Deals"
  org_id = "123"
}
```

## Deals example when using a tracked organization

```hcl
data "pipedrive_organizations" "example"{
  id = "123"
}

resource "pipedrive_deals" "example" {
  title  = "Example Deals"
  org_id = data.pipedrive_organizations
}
```

## Argument Reference

The following arguments are supported:

* `title` - (Required) The title of the deals that will be created in pipedrive.
* `org_id` - (Required) The organization ID of the organization where the deal will be created.
* `status` - (Optional) `open`, `won`, `lost`, `deleted`. If omitted, status will be set to open.

## Attributes Reference

In addition to all argument, the following attributes are exported:

* `id` - The ID of this resource

## Import

Pipedrive deals can be imported using the `id` eg,

`terraform import pipedrive_deals.example 123`
