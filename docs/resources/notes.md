# Resource: pipedrive_notes

## Example Usage

```hcl
resource "pipedrive_notes" "example" {
  content = "The content of the note in HTML format. Subject to sanitization on the back-end."
  deal_id = "123"
}
```

## Notes example when using a tracked deal

```hcl
data "pipedrive_organizations" "example"{
  id = "123"
}

resource "pipedrive_deals" "example" {
  title  = "Example Deals"
  org_id = data.pipedrive_organizations
}

resource "pipedrive_notes" "example" {
  content = "The content of the note in HTML format. Subject to sanitization on the back-end."
  deal_id = pipedrive_deals.example.id
}
```

## Argument Reference

The following arguments are supported:

* `content` - (Required) The content of the note in HTML format. Subject to sanitization on the back-end.
* `deal_id` - (Required) The ID of the deal the note will be attached to.

## Attributes Reference

In addition to all argument, the following attributes are exported:

* `id` - The ID of this resource
* `add_time` - The time the note was added
* `deal_attached` - The title of the deal attached to the note
* `update_time` - The time the note was updated

## Import

Pipedrive notes can be imported using the `id` eg,

`terraform import pipedrive_notes.example 123`
