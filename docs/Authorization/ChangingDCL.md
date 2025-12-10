# Changing DCL

Once your application is delivered and in use (especially for multi-tenant applications), any changes to the Authorization model must be done with special care.

## Allowed Changes

The following changes are currently allowed:

- Adding new policies
- Adding new schema attributes
- Changing actions and resources\* of policies
- Changing conditions of existing policies\*
    - **Exception**: `IS [NOT] RESTRICTED` conditions may not be removed as that will break admin policies which restrict these attributes.

\* In applications with multiple microservices, it is important that all microservices can handle the refactoring, see [Migration Strategy](#migration-strategy).

## Forbidden Changes

The following changes will lead to errors in admin policies that re-use the affected base policies.

- Deleting policies
- Renaming policies (equivalent to deleting and adding new policy)
- Moving policies to a different package (equivalent to renaming)
- Moving policies to a different file in the same package
- Renaming DCL files (equivalent to moving policies to a different file)
- Deleting `IS (NOT) RESTRICTED` conditions from a policy.
- Adding or removing the  ```DEFAULT``` keyword of policies.
- Adding or removing the  ```INTERNAL``` keyword of policies.

## Consequences of Forbidden Changes

Admin policies with errors, e.g. an invalid reference to a base policy that has been deleted, are removed from the authorization bundle.
Users that had these policies assigned will effectively lose these authorizations after the policy deployment is mirrored to applications.
Policies with errors will be shown to the Administrator in the administration console of SAP Cloud Identity Services inside the authorization policies of the application.
These policies can then be corrected or removed.

Example for a forbidden change that would **break** admin policies:

**Policy before upgrade**
```sql
POLICY ReadSalesOrders {
    GRANT read ON SalesOrders WHERE Country IS NOT RESTRICTED;
}
```

**Policy after upgrade**
```sql
POLICY ReadSalesOrders {
    GRANT read ON SalesOrders WHERE CustomerCountry IS NOT RESTRICTED;
}
```

**Broken admin policy built on top of the base policy**
```sql
POLICY salesOrderRestricted {
    USE shopping.ReadSalesOrders
        RESTRICT Country = 'DE'; -- ⚡️⚡️⚡️
}
```

::: warning
Forbidden changes are not technically prevented. Please take utmost care when changing the delivered
policies. If you need to do any changes, you must inform your customers and allow a migration period.
:::

## Migration Strategy
Internal changes, such as changing actions, resources, or conditions, can be achieved in a three-step process if there are multiple microservices involved:

1. Add the new *actions*/*resources*/*conditions* as an additional grant to the affected base policy.
2. Migrate the microservices individually to check for the new *actions*/*resources*/*conditions*.
3. Remove the grants with the old *actions*/*resources*/*conditions* from the affected base policy.