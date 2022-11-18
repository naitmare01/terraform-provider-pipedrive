# Data Source: pipedrive_organizations

## Example Usage

```hcl
data pipedrive_organizations "example"{
  id = "123"
}

output "example" {
  value = data.pipedrive_organizations.example
}
```

## Argument Reference

The following arguments are supported:

* **id** - (Required) The id of the organization.

## Attributes Reference

In addition to all argument, the following attributes are exported:

* **id** The ID of this resource.
* **title** The title of this resource
