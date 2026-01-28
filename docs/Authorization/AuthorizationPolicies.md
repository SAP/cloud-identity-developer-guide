# Authorization Policies

Authorization policies grant the right to perform actions on protected resources in an application. They can be assigned to users to control access to various parts of the application.

Developers can define a set of base policies that can be assigned directly or used as building blocks by the application administrators to create additional so-called admin policies at runtime.

## DCL

Authorization policies are defined in a domain-specific language called Data Control Language (**DCL**) that supports conditions that can be used to grant fine-grained access to resources.

### Example
Here is an example of authorization policies defined in DCL:

```dcl
SCHEMA {
    category: String;
}

POLICY ReadProducts {
    GRANT read ON products WHERE category IS NOT RESTRICTED;
}

POLICY ReadOfficeSupplies {
    USE ReadProducts RESTRICT category = 'OfficeSupplies';
}
```

### Specification
The complete specification for DCL can be found in the [SAP Help Portal](https://help.sap.com/docs/cloud-identity-services/cloud-identity-services/data-control-language-dcl).