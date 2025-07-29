# Authorization Checks

In this section, we will cover the basic concepts of authorization checks with the Authorization Management Service (AMS).

::: tip
In CAP applications, it is typically not necessary to implement authorization checks programmatically. Instead, authorization requirements are [declared](#declarative-authorization-checks) via [annotations](https://cap.cloud.sap/docs/guides/security/authorization#requires). The resulting authorization checks are performed dynamically for the application by the AMS modules.

As CAP has role-based authorization, AMS policies and authorization checks in CAP follow a [*role-based*](/CAP/Basics#role-policies) paradigm instead of the standard *action*/*resource* paradigm documented below.
:::

## Actions and Resources

Authorization policies grant the right for a one (or multiple) *actions* on one (or multiple) *resources*. For example:

```dcl
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

[Node.js Details](/Libraries/nodejs/sap_ams/sap_ams.md) / [Java Details](/Libraries/java/jakarta-ams/jakarta-ams.md)

:::

## Conditional Policies

Grants of authorization policies can be restricted by conditions to filter the entities of a resource on which the action is allowed. For example, a policy can grant the right to read products only if the product is in stock:

```dcl
POLICY ReadProducts {
    GRANT read ON products WHERE stock > 0;
}
```

### Conditional checks for single entity

When checking the authorization for a single entity, the condition is evaluated against the attributes of that entity. For example, if the product has `stock = 10`, the check will return `true`:

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

[Node.js Details](/Libraries/nodejs/sap_ams/sap_ams.md) / [Java Details](/Libraries/java/jakarta-ams/jakarta-ams.md)
:::

### Conditional checks for multiple entities

When checking the authorization for multiple entities, the application has two options:

1. It can fetch all entities from the database, loop over them and check access for each entity individually.
2. It can filter the entities in the database query based on the condition.

#### Looping
The first option is simpler to implement and may be sufficient for resources with few entities, such as a list of landscapes:

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

```dcl [DCL]
POLICY AccessEUCanaryLandscapes {
    GRANT access ON landscape WHERE name LIKE '%canary' AND region = 'EU';
}
```
:::

However, this strategy can lead to performance issues for larger entity sets, for which thousands of values would need to be checked individually.

#### Filtering
The second option is to filter the entities already in the database query. This is more efficient, as it reduces the number of entities in application memory to those that the user is allowed to access. However, this strategy is non-trivial to implement as it requires traversing the condition tree and translating it into a query language condition.

In CAP projects, it is implemented out-of-the-box by the libraries to dynamically translate filter conditions imposed by authorization policies to *CQL* conditions. For non-CAP projects, we recommend to contact us for assistance with the existing API or discuss a feature request for a standard transformer to the required query format.

## Declarative Authorization Checks
Instead of manually implementing authorization checks in the code, it is sometimes more elegant to impose them automatically with declarations for required privileges.
For example, in CAP applications, the standard `@restrict` and `@requires` annotations are used to make checks for roles with AMS.
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

```dcl [DCL]
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

[Node.js Details](/Libraries/nodejs/sap_ams/sap_ams.md#amsmiddleware) / [Java Spring Details](/Libraries/java/spring-ams/spring-ams.md)

:::

### Advantages
Declarative authorization checks have several advantages:
  - concise syntax
  - provides central overview of required privileges for different parts of the application
  - allows changing required privileges without touching the implementation of service handlers
  - prevents accidental information leaks, for example by returning 404 instead of 403 while fetching database entities for the actual authorization check in the service handler

### Disadvantages
However, this approach is not enough for *action*/*resource* pairs for which conditional access may be granted. The best we can do in this case, is to do a pre-check for the action and resource, and then let the service handler perform an additional check for the condition. This is because the condition check requires additional attribute input, typically involving the database.