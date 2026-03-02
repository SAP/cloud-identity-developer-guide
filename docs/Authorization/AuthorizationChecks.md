# Authorization Checks

In this section, we cover the basic concepts of authorization checks with the Authorization Management Service (**AMS
**).

::: tip
In CAP applications, it's typically not necessary to implement authorization checks programmatically. Instead,
authorization requirements are [declared](#declarative-authorization-checks)
via [annotations](https://cap.cloud.sap/docs/guides/security/authorization#requires). The AMS modules perform the
resulting authorization checks dynamically for the application.

Since CAP has role-based authorization, authorization policies and authorization checks in CAP follow a [
*role-based*](/CAP/Basics#role-policies) paradigm instead of the standard *action*/*resource* paradigm documented below.
:::

## Actions and Resources

Authorization policies grant the right for one (or multiple) *actions* on one (or multiple) *resources*. For example:

```dcl
POLICY ReadProducts {
    GRANT read ON products;
}
```

Therefore, a typical authorization check answers the question whether a user is allowed to perform a specific action on
a specific resource, for example, whether a user is allowed to read products.

::: code-group

```js [Node.js]
const decision = authorizations.checkPrivilege('read', 'products');
if (decision.isGranted()) {
    // user is allowed to read products
} else {
    // user is not allowed to read products
}
```

```java [Java]
Decision decision = authorizations.checkPrivilege("read", "products");
if(decision.

isGranted()){
        // user is allowed to read products
        }else{
        // user is not allowed to read products
        }
```

:::

Instead of checking a single action on a single resource, we can also query AMS for
a [list of action/resource privileges](#querying-potential-privileges) that are granted to the user.

## Authorizations

Authorization checks are performed with an `Authorizations` object. It represents the set of authorization policies
applicable for the current request and -in more complex authentication scenarios- how to combine authorizations from
different layers (user, technical client, etc.).

The `Authorizations` are built **once per request** after authentication and then used for all authorization checks
performed in this request context.

::: tip CdsAuthorizations
In CAP Applications, the `CdsAuthorizations` interface is used instead. It decorates the standard `Authorizations`
interface with additional CAP-specific methods for role-based authorization checks which delegate internally to an
`Authorizations` object.
:::

## AuthorizationsProvider

To create and access the `Authorizations` object for the current request, an `AuthorizationsProvider` is used. It
determines which policies apply and provides default values for authorization attributes of the principal such as
`$user.email`.

::: warning Standard implementations
For the standard SAP BTP security service offerings, we highly recommend using the built-in `AuthorizationsProvider`
implementations.
They implement the officially recommended authorization strategies correctly, including more complex scenarios like
inbound request flows from external applications.
When internals change, these implementations will be patched in a backward-compatible way to ensure applications share a
streamlined, well-tested implementation.
:::

### SciAuthorizationsProvider

The `SciAuthorizationsProvider` (Node.js: `IdentityServiceAuthProvider`) is the recommended default for applications
using SAP Cloud Identity Services for authentication. It derives authorizations from SAP Identity Service token
principals.

::: code-group

```js [Node.js]
const {IdentityServiceAuthProvider} = require('@sap/ams');

const authProvider = new IdentityServiceAuthProvider(ams);
```

```java [Java]
import com.sap.cloud.security.ams.core.SciAuthorizationsProvider;

SciAuthorizationsProvider<Authorizations> authProvider
        = SciAuthorizationsProvider.create(ams);
```

:::

The `SciAuthorizationsProvider` combines authorizations from two sources:

- **User Authorizations**: Policies assigned to the authenticated user in the SAP Cloud Identity Services directory.
- **Client Authorizations**: Policies derived from technical communication scenarios (e.g., consumed App-to-App APIs,
  BTP service plans).

By default*, the authorizations of these two layers are combined as follows:

| User Authorizations | Client Authorizations | Result                                                                     |
|---------------------|-----------------------|----------------------------------------------------------------------------|
| present             | null                  | User authorizations are used directly (e.g. named user token)              |
| null                | present               | Client authorizations are used directly (e.g. technical user token)        |
| present             | present               | Logical intersection of both is granted (e.g. principal propagation token) |
| null                | null                  | Fully denied, empty authorizations (unexpected scenario)                   |

\* *In the future, it might be possible to explicitly decide
for [principal propagation tokens](/Authorization/TechnicalCommunication) how the authorizations
should be enforced with a configuration property of the App-to-App dependency. In that case, this default logic
would be overridden based on this information in the token.*

**Example**:

Consider an authorization check with [conditional policies](#conditional-policies).

| User Authorizations      | Client Authorizations    | Effective Condition                                      |
|--------------------------|--------------------------|----------------------------------------------------------|
| `category = 'Equipment'` | null                     | `category = 'Equipment'`                                 |
| null                     | `category = 'Equipment'` | `category = 'Equipment'`                                 |
| `category = 'Equipment'` | `price < 100`            | `category = 'Equipment' AND price < 100`                 |
| null                     | null                     | Fully denied, empty authorizations (unexpected scenario) 

#### Customization

`SciAuthorizationsProvider` supports customization through configuration methods and method overriding.

The current configuration methods are for [Technical Communication](/Authorization/TechnicalCommunication). They are
described on that page in detail.

##### Overriding Methods

You can override the `getUserAuthorizations`,
`getClientAuthorizations` and the methods for building default input for authorization checks if necessary.

**Example: Customizing Default Input for Authorization Checks**

In this example, we override `getDefaultInput` to include a custom user attribute (`$user.division`) from a token
claim that is not included by default:

::: code-group

```js [Node.js]
class CustomAuthProvider extends IdentityServiceAuthProvider {
    /**
     * @param {import("@sap/xssec").IdentityServiceSecurityContext} securityContext
     */
    getInput(securityContext) {
        const defaultInput = super.getInput(securityContext);

        const division = securityContext.token.payload.division;
        if (division) {
            defaultInput["$user.division"] = division;
        }

        return defaultInput;
    }
}
```

```java [Java]
import com.sap.cloud.security.ams.core.SciAuthorizationsProvider;
import com.sap.cloud.security.ams.api.*;
import com.sap.cloud.security.ams.api.expression.AttributeName;

import java.util.Map;

public class CustomAuthorizationsProvider extends SciAuthorizationsProvider<Authorizations> {

    private static final AttributeName $USER_DIVISION = AttributeName.of("$user.division");

    public CustomAuthorizationsProvider(AuthorizationManagementService ams) {
        super(ams, AuthorizationsAdapter::identity);
    }

    @Override
    protected Map<AttributeName, Object> getDefaultInput(Principal principal) {
        Map<AttributeName, Object> defaultInput = super.getDefaultInput(principal);

        principal.getClaimAsString("division")
                .ifPresent(division -> defaultInput.put($USER_DIVISION, division));

        return defaultInput;
    }
}
```

:::

### HybridAuthorizationsProvider

The `HybridAuthorizationsProvider` (Node.js: `HybridAuthProvider`) is recommended for applications that have migrated
from XSUAA to AMS. It extends `SciAuthorizationsProvider` with additional support for XSUAA tokens by mapping XSUAA
scopes to AMS base policies.

**Scope to Policy Mapping**

When an XSUAA token is received, the provider extracts scopes from the token and maps them to AMS policies using a
configured `ScopeMapper`. Relevant scopes are typically prefixed with the application's `xsappname`
(e.g., `na-foobar!t4711`) in the token.

**Example Mapping** (xsappname: `na-foobar!t4711`)

| Scope in Token                  | Mapped Policy                                          |
|---------------------------------|--------------------------------------------------------|
| `na-foobar!t4711.ProductReader` | `shopping.ReadProducts`                                |
| `na-foobar!t4711.ProductAdmin`  | `shopping.ReadProducts`, `shopping.WriteProducts`      |
| `openid`                        | *(no policy - generic scope without xsappname prefix)* |
| `na-foobar!t4711.UnknownScope`  | *(no policy - not in mapping)*                         |

::: code-group

```js [Node.js]
const {HybridAuthProvider} = require('@sap/ams');

const scopeToPolicyMapper = (scope) => {
    const scopeToPoliciesMap = {
        'na-foobar!t4711.ProductReader': ['shopping.ReadProducts'],
        'na-foobar!t4711.ProductAdmin': ['shopping.ReadProducts', 'shopping.WriteProducts'],
    };
    return scopeToPoliciesMap[scope] || [];
};

const authProvider = new HybridAuthProvider(ams, scopeToPolicyMapper);
```

```java [Java]
import com.sap.cloud.security.ams.core.HybridAuthorizationsProvider;
import com.sap.cloud.security.ams.api.ScopeMapper;
import com.sap.cloud.security.ams.api.PolicyName;

PolicyName READ_PRODUCTS = PolicyName.of("shopping.ReadProducts");
PolicyName WRITE_PRODUCTS = PolicyName.of("shopping.WriteProducts");

Map<String, Set<PolicyName>> scopeToPoliciesMap = Map.of(
        "ProductReader", Set.of(READ_PRODUCTS),
        "ProductAdmin", Set.of(READ_PRODUCTS, WRITE_PRODUCTS)
);

HybridAuthorizationsProvider<?> authProvider = HybridAuthorizationsProvider
        .create(ams, ScopeMapper.ofMapMultiple(scopeToPoliciesMap))
        .withXsAppName("na-foobar!t4711");
```

:::

### Custom Implementation

If necessary, you can also implement a custom `AuthorizationsProvider` with your own logic for determining which
policies apply:

::: code-group

```js [Node.js (Typescript)]
import {Authorizations, Types, XssecAuthProvider} from "@sap/ams";
import {SecurityContext, Service, Token, XsuaaSecurityContext} from "@sap/xssec";

class XsuaaAuthProvider
    extends XssecAuthProvider<SecurityContext<Service, Token

>>
implements
XssecAuthProvider < XsuaaSecurityContext > {

    getAuthorizations(securityContext: XsuaaSecurityContext): Promise<Authorizations> {
        throw new Error("Method not implemented.");
    }
    getInput(securityContext: XsuaaSecurityContext): Types.AttributeInput {
        throw new Error("Method not implemented.");
    }
    supportsSecurityContext(securityContext: XsuaaSecurityContext): void {
        throw new Error("Method not implemented.");
    }
}
```

```java [Java]
import com.sap.cloud.security.ams.api.*;
import com.sap.cloud.security.ams.api.PolicyName;

public class CustomAuthorizationsProvider implements AuthorizationsProvider<Authorizations> {
    private final AuthorizationManagementService ams;

    public CustomAuthorizationsProvider(AuthorizationManagementService ams) {
        this.ams = ams;
    }

    @Override
    public Authorizations getAuthorizations(Principal principal) {
        // Custom logic to determine which policies apply
        Set<PolicyName> policies = determinePoliciesFromContext(principal);
        return ams.getAuthorizations(policies);
    }

    private Set<PolicyName> determinePoliciesFromContext(Principal principal) {
        // Your custom policy resolution logic here
        return Set.of(PolicyName.of("myapp.DefaultPolicy"));
    }
}

AuthorizationsProvider<?> authProvider = new CustomAuthorizationsProvider(ams);
```

:::

## Conditional Policies

Grants of authorization policies can be made conditional on dynamic data.
After declaring relevant attributes in a schema, policies can reference those in where-conditions.
This is usually used to filter the entities of a resource on which the action is allowed.
However, conditions may also be based on other data in the context of the authorization check, e.g. a specific user
property.

**Example** A policy can grant the right to read products only from a specific category:

```dcl
SCHEMA {
    category: String;
}

POLICY ReadEquipment {
    GRANT read ON products WHERE category = 'Equipment';
}
```

### Fixed Attribute values

When performing an authorization check, the values of relevant attributes may already be known.
In this case, those attribute values can be provided as part of the authorization check.

**Example** The check should be performed for a particular product category:

::: code-group

```js [Node.js]
const decision = authorizations.checkPrivilege(
    'read', 'products', {category: 'Equipment'});
if (decision.isGranted()) {
    // user is allowed to read products in the 'Equipment' category
} else {
    // user is not allowed to read products in the 'Equipment' category
}
```

```java [Java]
Decision decision = authorizations.checkPrivilege(
        "read", "products", Map.of("category", "Equipment"));
if(decision.

isGranted()){
        // user is allowed to read products in the 'Equipment' category
        }else{
        // user is not allowed to read products in the 'Equipment' category
        }
```

[Node.js Details](/Libraries/nodejs/sap_ams/sap_ams.md#handling-decisions) / [Java Details](/Libraries/java/jakarta/jakarta-ams.md#decision)
:::

### Dynamic Attribute values

When the values of the relevant attributes cannot be provided as part of the authorization check, the authorization
check can still be performed.
In this case, the decision resulting from the authorization check is typically conditional - unless a policy explicitly
grants unrestricted or fully restricted access based on these attributes.

The application has two options to handle conditional decisions:

1. Loop over each entity instance individually, apply the entity's attribute values to the decision and check whether
   the access is granted.
2. Delegate the filtering process to the data retrieval, e.g., to a database query based on the conditional decision.

##### Looping

The first option is easier to implement and is fine when only a few instances are involved:

**Example**

::: code-group

```js [Node.js]
const catalog = [
    {name: 'Notebook', category: 'Equipment'},
    {name: 'Printer', region: 'Equipment'},
    {name: 'Toner', region: 'OfficeSupplies'}
];

const decision = authorizations.checkPrivilege('read', 'products');
const accessibleProducts =
    catalog
        .filter(product => {
            return decision.apply({
                '$app.category': product.category
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
                .filter(product -> decision.apply(
                        Map.of(AttributeName.of("category"), product.getCategory())
                ).isGranted())
                .collect(Collectors.toList());
```

:::

However, this strategy can lead to performance issues for larger collections, for which thousands of values must be
checked individually.

##### Filtering

The second option is to filter the entities before they enter the application.
This is more efficient because it reduces the number of instances in the application memory to those instances that the
user is allowed to access.
However, this strategy is non-trivial to implement because it requires traversing the condition tree and translating it
into a query language expression.

::: tip CAP Projects
In CAP projects, this translation is implemented out-of-the-box by the AMS plugins which translate filter conditions
imposed by authorization policies to *CQL*/*CXN* expressions.
:::

For non-CAP projects, we aim to provide extractors for standard query languages. We recommend contacting us for
assistance with the existing API or discuss a feature request for missing extractors for your query format.

As of today, there is a basic extractor for SQL queries available in the Java AMS library:

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

We can add configuration options and features on request.

## Declarative Authorization Checks

Instead of manually implementing authorization checks scattered over the code base, it improves maintainability to
declare the required privileges for different parts of the application.
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
@CheckPrivilege(action = "create", resource = "orders")
public Order createOrder(
        Product product,
        int quantity,
        @AmsAttribute(name = "product.category") String productCategory) {
    if (!Objects.equals(product.getCategory(), productCategory)) {
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
- prevents accidental information leaks, for example by returning 404 instead of 403 while fetching database entities to
  gather information for an authorization check in the service handler

### Limitations

::: warning Conditional Policies with Instance-Based Access
Not all declaration methods support checks for *action*/*resource* pairs with instance-based access conditions. In this
case, they can only be used for pre-checks but the service handler must perform an additional check for the condition.
This is because the condition check requires additional attribute input, typically involving information from the
database.
:::

## Querying Potential Privileges

In addition to checking specific privileges, applications can query which actions, resources, or privileges are granted
to the user. These methods are useful for **pre-checks**, such as determining which UI elements to display to a user
before they attempt an action.

::: warning Conditions are Ignored
These methods ignore any conditions on grants during evaluation. The returned actions, resources, or privileges may
still depend on conditions, and an additional `checkPrivilege` call **must** be performed before actually allowing the
action on the resource.
:::

### getPotentialResources

Collects all resources for which at least one action is potentially granted:

::: code-group

```js [Node.js]
const potentialResources = authorizations.getPotentialResources();
// Returns: Set<string>, e.g., Set { 'products', 'orders', 'customers' }

for (const resource of potentialResources) {
    console.log(`User may have access to: ${resource}`);
}
```

```java [Java]
Set<String> potentialResources = authorizations.getPotentialResources();
// Returns: Set<String>, e.g., ["products", "orders", "customers"]

for(
String resource :potentialResources){
        System.out.

println("User may have access to: "+resource);
}
```

:::

### getPotentialActions

Collects all actions that are potentially granted for a given resource:

::: code-group

```js [Node.js]
const potentialActions = authorizations.getPotentialActions('products');
// Returns: Set<string>, e.g., Set { 'read', 'create', 'update' }

if (potentialActions.has('delete')) {
    // Show delete button in UI
}
```

```java [Java]
Set<String> potentialActions = authorizations.getPotentialActions("products");
// Returns: Set<String>, e.g., ["read", "create", "update"]

if(potentialActions.

contains("delete")){
        // Show delete button in UI
        }
```

:::

### getPotentialPrivileges

Collects all action/resource combinations that are potentially granted:

::: code-group

```js [Node.js]
const potentialPrivileges = authorizations.getPotentialPrivileges();
// Returns: Array<{action: string, resource: string}>
// e.g., [{ action: 'read', resource: 'products' }, { action: 'create', resource: 'orders' }]

for (const privilege of potentialPrivileges) {
    console.log(`User may: ${privilege.action} on ${privilege.resource}`);
}
```

```java [Java]
Set<Privilege> potentialPrivileges = authorizations.getPotentialPrivileges();
// Returns: Set<Privilege>
// e.g., [Privilege("read", "products"), Privilege("create", "orders")]

for(
Privilege privilege :potentialPrivileges){
        System.out.

println("User may: "+privilege.action() +" on "+privilege.

resource());
        }
```

:::

### Use Cases

These methods are particularly useful for:

- **UI Rendering**: Determine which menu items, buttons, or sections to display based on the user's potential
  authorizations.
- **Feature Toggles**: Enable or disable features in the UI based on whether the user might have access.
- **Navigation Guards**: Pre-filter accessible routes or views.

::: tip Remember to Verify
The results of these methods should only be used for UI hints and pre-checks. Always perform an actual `checkPrivilege`
call when the user attempts to execute an action to ensure proper authorization enforcement.
:::