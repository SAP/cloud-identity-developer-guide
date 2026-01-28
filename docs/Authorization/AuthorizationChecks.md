# Authorization Checks

In this section, we cover the basic concepts of authorization checks with the Authorization Management Service (**AMS**).

::: tip
In CAP applications, it's typically not necessary to implement authorization checks programmatically. Instead, authorization requirements are [declared](#declarative-authorization-checks) via [annotations](https://cap.cloud.sap/docs/guides/security/authorization#requires). The AMS modules perform the resulting authorization checks dynamically for the application.

Since CAP has role-based authorization, authorization policies and authorization checks in CAP follow a [*role-based*](/CAP/Basics#role-policies) paradigm instead of the standard *action*/*resource* paradigm documented below.
:::

## Actions and Resources

Authorization policies grant the right for one (or multiple) *actions* on one (or multiple) *resources*. For example:

```dcl
POLICY ReadProducts {
    GRANT read ON products;
}
```

Therefore, a typical authorization check answers the question whether a user is allowed to perform a specific action on a specific resource, for example, whether a user is allowed to read products.

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
Decision decision = authorizations.checkPrivilege("read", "products");
if(decision.isGranted()) {
    // user is allowed to read products
} else {
    // user is not allowed to read products
}
```

[Node.js Details](/Libraries/nodejs/sap_ams/sap_ams.md#authorization-checks) / [Java Details](/Libraries/java/jakarta/jakarta-ams.md#authorization-checks)

:::

## Authorizations

Authorization checks are performed with an `Authorizations` object. It represents the set of authorization policies applicable for the current request and -in more complex scenarios- how to apply them.

Typically, for each request, a new `Authorizations` object is created after authentication which can be used for one or many authorization checks during the request processing.

::: tip AmsCapAuthorizations
In CAP Applications, the `AmsCapAuthorizations` interface is used instead. It extends the standard `Authorizations` interface with CAP-specific methods for role-based authorization checks which delegate internally to an `Authorizations` object.
:::

### AuthProvider

To access the `Authorizations` object for the current request, an `AuthProvider` is used. It determines which policies apply and provides default values for authorization attributes such as `$user.email`. Typically, the `Authorizations` are built from the thread-local security context after authentication.

For the standard SAP BTP security service offerings, use the built-in `AuthProvider` implementations.
They implement the recommended authorization strategies, including advanced scenarios like inbound request flows from external applications.

- **`IdentityServiceAuthProvider`** (recommended default): Derives authorizations from SAP Identity Service security context
- **`HybridAuthProvider`**: For applications that have migrated from XSUAA to AMS - extends `IdentityServiceAuthProvider`, allows the mapping of scopes from XSUAA security contexts to base policies in AMS

::: code-group

```js [Node.js]
const { IdentityServiceAuthProvider, HybridAuthProvider } = require('@sap/ams');

const authProvider = new IdentityServiceAuthProvider(ams);

// or for migrated XSUAA applications:
const scopeToPolicyMapper = (scope) => {
    const scopeToPoliciesMap = {
        'ProductReader': ['shopping.ReadProducts'],
        'Customer': ['shopping.ReadProducts', 'shopping.CreateOrders']
    };
    return scopeToPoliciesMap[scope] || [];
};
const authProvider = new HybridAuthProvider(ams, scopeToPolicyMapper);
```

```java [Java]
import com.sap.cloud.security.ams.core.*;

AuthProvider authProvider = new IdentityServiceAuthProvider(ams);

// or for migrated XSUAA applications:
Map<String, Set<String>> scopeToPoliciesMap = Map.of(
    "ProductReader", Set.of("shopping.ReadProducts"),
    "Customer", Set.of("shopping.ReadProducts", "shopping.CreateOrders")
);

HybridAuthProvider authProvider = new HybridAuthProvider(
    ams,
    ScopeMapper.ofMapMultiple(scopeToPoliciesMap)
);
```

:::

Alternatively, you can implement a custom `AuthProvider`, e.g. to derive applicable policies from additional sources.

<details>
<summary> Custom AuthProviders</summary>

**Customizing standard AuthProviders:**

You can extend the built-in `AuthProvider` implementations to customize the behavior, e.g., to apply more than those policies which are assigned to users in the user directory:

```java
public class CustomIdentityServiceAuthProvider extends IdentityServiceAuthProvider {
    public CustomIdentityServiceAuthProvider(AuthorizationManagementService ams) {
        super(ams);
    }
    
    @Override
    protected Authorizations getUserAuthorizations(SapIdToken token) {
        Authorizations authorizations = super.getUserAuthorizations(token);

        Set<String> policies = new HashSet<>(authorizations.getPolicies());
        policies.addAll(/* add policies from other sources for this user */);
        authorizations.setPolicies(policies);

        return authorizations;
    }
}

AuthProvider authProvider = new CustomIdentityServiceAuthProvider(ams);
```

**Custom AuthProvider implementation:**

You can also implement a custom `AuthProvider` with your own logic for determining which policies apply:

```java
public class CustomAuthProvider implements AuthProvider {
    private final AuthorizationManagementService ams;
    
    public CustomAuthProvider(AuthorizationManagementService ams) {
        this.ams = ams;
    }
    
    @Override
    public Authorizations getAuthorizations() {
        // Custom logic to determine which policies apply
        Set<String> policies = determinePoliciesFromContext();
        return ams.getAuthorizations(policies);
    }
}

AuthProvider authProvider = new CustomAuthProvider(ams);
```

</details>

## Conditional Policies

Grants of authorization policies can be made conditional on dynamic data.
After declaring relevant attributes in a schema, policies can reference those in where-conditions.
This is usually used to filter the entities of a resource on which the action is allowed.
However, conditions may also be based on other data in the context of the authorization check, e.g. a specific user property.

**Example** A policy can grant the right to read products only from a specific category:

```dcl
SCHEMA {
    category: String;
}

POLICY ReadEquipment {
    GRANT read ON products WHERE category = 'Equipment';
}
```

### Checking Authorizations with fixed attributes

When performing an authorization check, the values of relevant attributes may already be known.
In this case, those attribute values can be provided as part of the authorization check.

**Example** The check should be performed for a particular product category:

::: code-group
```js [Node.js]
const decision = authorizations.checkPrivilege(
    'read', 'products', { category: 'Equipment' });
if (decision.isGranted()) {
    // user is allowed to read products in the 'Equipment' category
} else {
    // user is not allowed to read products in the 'Equipment' category
}
```

```java [Java]
Decision decision = authorizations.checkPrivilege(
    "read", "products", Map.of("category", "Equipment"));
if (decision.isGranted()) {
    // user is allowed to read products in the 'Equipment' category
} else {
    // user is not allowed to read products in the 'Equipment' category
}
```

[Node.js Details](/Libraries/nodejs/sap_ams/sap_ams.md#handling-decisions) / [Java Details](/Libraries/java/jakarta/jakarta-ams.md#has-user-read-access-to-salesorders-resource-with-countrycode-de-and-salesorder-type-4711)
:::

### Checking Authorizations with variable attributes

When the values of attributes relevant for the authorization check are part of the entities to be protected, the authorization check can still be performed.
In this case, the decision resulting from the authorization check becomes conditional.
Depending on the data storage used by the application, the actual attribute values can be applied later when querying the data.
The application has two options:

1. Loop over all entities and check access for each instance individually.
2. Delegate the filtering process, e.g., by using a database query based on the conditional decision.

##### Looping
The first option is easier to implement and is fine when only a few instances are involved:

**Example**

::: code-group
```js [Node.js]
const catalog = [
    { name: 'Notebook', category: 'Equipment' },
    { name: 'Printer', region: 'Equipment' },
    { name: 'Toner', region: 'OfficeSupplies' }
];

const decision = authorizations.checkPrivilege('read', 'products');
const accessibleProducts = 
    catalog
    .filter(product => {
        return decision.apply({ 
            '$app.category' : product.category
        }).isGranted();
    });
```

```java [Java]
List<Map<String, Product>> catalog = List.of(
    Product.create("Notebook", "Equipment"),
    Product.create("Printer", "Equipment"),
	Product.create("Toner", "OfficeSupplies")
);

List<Product> accessibleProducts = 
    catalog.stream()
    .filter(product -> authorizations.checkPrivilege("read", "products",
            Map.of("$app.category", product.getCategory())
    ).isGranted())
    .collect(Collectors.toList());
```
:::

However, this strategy can lead to performance issues for larger collections, for which thousands of values must be checked individually.

##### Filtering
The second option is to filter the entities before they enter the application.
This is more efficient because it reduces the number of instances in the application memory to those instances the user is allowed to access.
However, this strategy is non-trivial to implement because it requires traversing the condition tree and translating it into a query language expression.

::: tip CAP Projects
In CAP projects, this translation is implemented out-of-the-box by the AMS plugins which translate filter conditions imposed by authorization policies to *CQL* expressions.
:::

For non-CAP projects, we aim to provide extractors for standard query languages. We recommend contacting us for assistance with the existing API or discuss a feature request for missing extractors for your query format.

As of today, there is an extractor for SQL queries available in the Java AMS library:

::: code-group
```java [Java]
// extractor can be built once per handler
SqlExtractor sqlExtractor = new SqlExtractor(Map.of(
    AttributeName.of("category"), "CategoryName")
);

Decision decision = authorizations.checkPrivilege("read", "products");
SqlExtractor.SqlResult sqlCondition = decision.visit(sqlExtractor);

String sqlQuery = String.format("SELECT * FROM Products WHERE %s;",
    sqlCondition.getSqlTemplate());
List<Product> accessibleProducts =
    db.query(sqlQuery, sqlCondition.getParameters(), Product.class);
```

```js [Node.js]
// Equivalent to Java snippet coming soon
```
:::

## Declarative Authorization Checks
Instead of manually implementing authorization checks scattered over the code base, it improves maintainability to declare the required privileges for different parts of the application.
This can, for example, be done centrally or by using code annotations.

```dcl
POLICY ReadProducts {
    GRANT read ON products;
}

POLICY OrderOfficeSupplies {
    GRANT create ON orders WHERE category = 'OfficeSupplies';
}
```

::: code-group
```js [Node.js (express)]
const app = express();
app.use(/^\/(?!health).*/i, authenticate, amsMw.authorize());

app.get('/products', amsMw.checkPrivilege('read', 'products'), getOrders);
app.post('/orders', amsMw.precheckPrivilege('create', 'orders'), createOrder);
```

```java [Spring (Route Security)]
import com.sap.cloud.security.ams.spring.authorization.AmsRouteSecurity;

@Configuration
@EnableWebSecurity
public class SecurityConfiguration {

    @Bean
    public SecurityFilterChain filterChain(HttpSecurity http,
            AmsRouteSecurity via) {
        http.authorizeHttpRequests(authz -> authz
                .requestMatchers(GET, "/products/**")
                    .access(via.checkPrivilege("read", "products"))
                .requestMatchers(POST, "/orders/**")
                    .access(via.precheckPrivilege("create", "orders")));
    
        return http.build();
    }
}
```

```java [Spring (Method Security)]
import com.sap.cloud.security.ams.spring.authorization.annotations.AmsAttribute;
import com.sap.cloud.security.ams.spring.authorization.annotations.CheckPrivilege;
import com.sap.cloud.security.ams.spring.authorization.annotations.PrecheckPrivilege;

/**
 * Performs an order creation, secured with instance-based authorization.
 *
 * @param product the product
 * @param quantity the quantity
 * @param productCategory the product category (used for authorization)
 * @return the created order
 */
@CheckPrivilege(action="create", resource="orders")
public Order createOrder(
    Product product,
    int quantity,
    @AmsAttribute(name="product.category") String productCategory) {
        if(!Objects.equals(product.getCategory(), productCategory)) {
            throw new IllegalArgumentException(
                "Authorization attribute for product category does not match the product");
        }
        
        // ... create order implementation
}
```

```cds [CAP]
// use standard cds @requires or @restrict annotations

service ProductService {
    @(restrict: [ { grant: 'READ', to: 'ReadProducts' } ])
    entity Products as projection on my.db.Products;
}

service OrderService {
    @(restrict: [ { 
        grant: ['READ', 'WRITE'],
        to: 'CreateOrders',
        // dynamically extended at runtime with product category = 'OfficeSupplies' filter
        where: 'createdBy = $user.email'
    } ])
    entity Orders as projection on my.db.Orders;
}
```

[Node.js Details](/Libraries/nodejs/sap_ams/sap_ams.md#amsmiddleware) / [Spring Route Security Details](/Libraries/java/spring/spring-ams.md#route-security) / [Spring Method Security Details](/Libraries/java/spring/spring-ams.md#method-security) / [CAP Details](/CAP/Basics)
:::

### Advantages
Declarative authorization checks have several advantages:
  - concise syntax
  - provides central overview of required privileges for different parts of the application
  - allows changing required privileges without touching the implementation of service handlers
  - prevents accidental information leaks, for example by returning 404 instead of 403 while fetching database entities to gather information for an authorization check in the service handler

### Limitations
::: warning Conditional Policies with Instance-Based Access
Not all declaration methods support checks for *action*/*resource* pairs with instance-based access conditions. In this case, they can only be used for pre-checks but the service handler must perform an additional check for the condition. This is because the condition check requires additional attribute input, typically involving information from the database.
:::