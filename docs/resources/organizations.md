# Resource: pipedrive_organizations

## Example Usage

```hcl
resource "pipedrive_organizations" "example" {
  name = "Example organizations"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the organization.

## Attributes Reference

In addition to all argument, the following attributes are exported:

* `id` - The ID of this resource
* `add_time` - The time the note was added
* `add_time` - The creation date & time of the organization in UTC

## Import

Pipedrive notes can be imported using the `id` eg,

`terraform import pipedrive_organizations.example 123`
