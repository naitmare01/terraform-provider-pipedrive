# Data Source: pipedrive_notes

## Example Usage

```hcl
data "pipedrive_notes" "example" {
  id = 123
}

output "pipedrive_notes" {
  value = data.pipedrive_notes.example
}
```

## Argument Reference

The following arguments are supported:

* **id** - (Required) The ID of the note.

## Attributes Reference

In addition to all argument, the following attributes are exported:

* **id** The ID of the note.
* **add_time** The time the note was added.
* **content** The content of the note.
* **deal_attached** The title of the deal attached to the note
* **update_time** The time the note was updated.
