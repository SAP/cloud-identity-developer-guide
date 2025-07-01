# Authorization Checks

In this section, we will cover the basic concepts of authorization checks with the Authorization Management Service (AMS).

## Actions and Resources

Authorization policies grant the right for a one (or multiple) *actions* on one (or multiple) *resources*. For example:

```SQL
POLICY ReadProducts {
    GRANT read ON products;
}
```

Therefore, a typical authorization check answers the question whether a user is allowed to perform a specific action on a specific resource. For example, whether a user is allowed to read products.

::: code-group
```js [Node.js]
const decision = authorizations.checkPrivilege('read', 'products');
if(decision.isGranted()) {
    // user is allowed to read products
} else {
    // user is not allowed to read products
}
```

```java [Java]
Attributes attributes = principal.getAttributes()
                .setAction("read")
                .setResource("products");

if(policyDecisionPoint.allow(attributes)) {
    // user is allowed to read products
} else {
    // user is not allowed to read products
}
```

[Node.js Details](/nodejs//sap_ams/sap_ams.md) / [Java Details](/java/jakarta-ams/jakarta-ams.md)

:::

## Conditional Policies

Grants of authorization policies can be restricted by conditions to filter the resources on which the action is allowed. For example, a policy can grant the right to read products only if the product is in stock:

```SQL
POLICY ReadProducts {
    GRANT read ON products IF stock > 0;
}
```

### Authorization checks for single resource

When checking the authorization for a single resource, the condition is evaluated against the resource attributes of that resource. For example, if the product has `stock = 10`, the check will return `true`:

::: code-group
```js [Node.js]
const decision = authorizations.checkPrivilege('read', 'products', { stock: 10 });
if(decision.isGranted()) {
    // user is allowed to read product with stock 10
} else {
    // user is not allowed to read product with stock 10
}
```

```java [Java]
Attributes attributes = principal.getAttributes()
                .setAction("read")
                .setResource("products")
                .app().value("stock", 10)
                .attributes();

if(policyDecisionPoint.allow(attributes)) {
    // user is allowed to read product with stock 10
} else {
    // user is not allowed to read product with stock 10
}
```

[Node.js Details](/nodejs//sap_ams/sap_ams.md) / [Java Details](/java/jakarta-ams/jakarta-ams.md)
:::

### Authorization checks for multiple resources

When checking the authorization for multiple resources, the application has two options:

1. It can fetch all resources from the database, loop over them and check access for each resource individually.
2. It can filter the resources in the database query based on the condition.

#### Looping
The first option is simpler to implement and may be sufficient for resources with a small domain, such as a list of landscapes:

::: code-group
```js [Node.js]
const decision = authorizations.checkPrivilege('access', 'landscape');
const landscapes = [
    { name: 'eu10', region: 'EU' },
    { name: 'eu10-canary', region: 'EU' },
    { name: 'us5', region: 'US' }
];

const accessibleLandscapes = 
    landscapes
    .filter(landscape => {
        return decision.apply({ 
            '$app.landscape.name' : landscape.name,
            '$app.landscape.region' : landscape.region
        }).isGranted();
    });
```

```sql [DCL]
POLICY AccessEUCanaryLandscapes {
    GRANT access ON landscape WHERE name LIKE '%canary' AND region = 'EU';
}
```
:::

However, this strategy can lead to performance issues for larger resource sets, such as REST entities, for which thousands of values would need to be checked individually.

#### Filtering
The second option is to filter the resources already in the database query. This is more efficient, as it reduces the number of resources fetched from the database to only those that the user is allowed to access. However, this strategy is non-trivial to implement as it requires traversing the condition tree and translating it into a query language condition.

In CAP projects, it is implemented out-of-the-box by the libraries to dynamically translate filter conditions imposed by authorization policies to *CQL* conditions. For non-CAP projects, we recommend to contact us for assistance with the existing API or discuss a feature request for a standard transformer to the required query format.

## Declarative Authorization Checks
Instead of manually implementing authorization checks in the code, it is sometimes more elegant to impose them automatically with declarations for required privileges.
As CAP applications are role-based, the standard `@restrict` and `@requires` annotations are used to make checks for roles with AMS.
In non-CAP applications, there are other ways to impose authorization checks by defining required privileges (i.e. *action*/*resource* pairs) on service endpoints:

::: code-group
```js [Node.js]
const app = express();
app.use(/^\/(?!health).*/i, authenticate, amsMw.authorize());

app.get('/orders', amsMw.precheckPrivilege('read', 'orders'), getOrders);
app.post('/orders', amsMw.checkPrivilege('read', 'products'), amsMw.precheckPrivilege('create', 'orders'), createOrder);
```

```java [Java]
// Example for Spring request matcher coming soon
```

```sql [DCL]
POLICY ReadProducts {
    GRANT read ON products;
}

POLICY ReadOrders {
    GRANT read ON orders WHERE order.createdBy = $user.email;
}

POLICY OrderAccessory {
    GRANT create ON orders WHERE product.category IN 'accesory';
}
```

[Node.js Details](/nodejs//sap_ams/sap_ams.md) / [Java Details](/java/jakarta-ams/jakarta-ams.md)

:::

### Advantages
The advantage of declarative instead of explicit authorization checks, is that they can typically be defined centrally. This gives a central overview of the required privileges per service endpoint and allows changing required privileges without touching the implementation.

### Disadvantages
However, this approach typically does not work well for *action*/*resource* pairs for which conditional access may be granted. The best we can do in this case, is to do a pre-check for the action and resource, and then let the application code handle the condition check. This is because the condition check requires access to the resource attributes.