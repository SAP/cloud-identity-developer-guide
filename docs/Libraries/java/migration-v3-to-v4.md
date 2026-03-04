# Migration Guide: v3 → v4

This guide helps you migrate your Java application from AMS Client Library version 3.x to version 4.x.

::: warning Recommendation
We recommend all applications to upgrade to version 4 to benefit from simpler customization, stream-lined
authorization strategies and security updates. Version 3 will not receive new features and will only get
critical bug fixes until its end of maintenance.
:::

## Overview of Changes

Version 4 introduces **significant** changes to the API, lifecycle objects and dependency modules.
However, breaking changes for **CAP applications** and **Spring Boot applications** are limited compared to plain Java integrations due to the annotation-based authorization checks.

### Lifecycle Object Changes
The core API changed from `PolicyDecisionPoint`, `Attributes` and `AttributesProcessor` to the following interfaces:

- `Authorizations`: Main API used for performing authorization checks\
(1 instance *per request context*)
- `AuthorizationManagementService`: Library instantiation, configuration, event logging\
(usually 1 instance *per application*)
- `AuthorizationsProvider`: Gives access to `Authorizations` for current principal, implementing official authorization strategies for different token flows with customization options\
(usually 1 instance *per application*)


## Dependency Migration

- Remove the previous AMS maven dependencies for group id `com.sap.cloud.security.ams.client`.\*
- Add the [recommended dependencies](/Authorization/GettingStarted#dependency-setup) for group id `com.sap.cloud.security.ams`.

\* You can keep the old `dcl-compiler-plugin` for now. However, there will be an improved `dcl-compiler-plugin` available very soon.




## Runtime Code Migration

### Initialization and Configuration

Replace `PolicyDecisionPoint` initialization with `AuthorizationManagementService` and `AuthorizationsProvider` as documented [here](/Authorization/AuthorizationBundle#client-library-initialization).

::: tip
If you are using a Spring Boot starter, the `AuthorizationManagementService` and `AuthorizationsProvider` are auto-configured and available for injection. You can customize them by overriding the beans or configuring them.
:::

### Startup Checks

- Remove manual startup checks and configuration properties.
- Implement startup check as described [here](/Authorization/AuthorizationBundle#startup-check)

::: tip
If you use a Spring Boot starter, it automatically integrates into Spring Boot readiness. Alternatively, you can use the Spring Boot 3 or 4 starter for health actuator integration.
:::

### AttributesProcessor Removal

- Remove any implementations of the `AttributesProcessor` interface and the meta data configuration for the service loader.

Typical use cases for `AttributesProcessor` such as [technical communication](/Authorization/TechnicalCommunication), [XSUAA scope mapping](/Authorization/AuthorizationChecks#hybridauthorizationsprovider) or [custom user attribute injection](/Authorization/AuthorizationChecks#overriding-methods) can now be implemented much simpler via `AuthorizationsProvider` configuration.

### PolicyDecisionPoint checks

- Rewrite `PolicyDecisionPoint#allow` checks with `Authorizations#checkPrivilege` as documented [here](/Authorization/AuthorizationChecks#performing-authorization-checks).

**Example**:

::: code-group

```java [v3]
import com.sap.cloud.security.ams.api.Principal;
import static com.sap.cloud.security.ams.dcl.client.pdp.Attributes.Names.APP;
import static com.sap.cloud.security.ams.dcl.client.pdp.Attributes.Names.ENV;

Principal principal = Principal.create();

// definite allow
Attributes attributes =
    principal
        .getAttributes()
        .setAction("read")
        .setResource("salesOrders");

boolean allowed = policyDecisionPoint.allow(attributes);

// allow that ignores any conditions
Attributes attributes =
    principal
        .getAttributes()
        .setAction("read")
        .setResource("salesOrders")
        .setIgnores(List.of(APP, ENV));

boolean allowed = policyDecisionPoint.allow(attributes);
```

```java [v4]
import static com.sap.cloud.security.ams.api.Principal.fromSecurityContext;

Authorizations authorizations = authProvider
    .getAuthorizations(fromSecurityContext());

// definite allow
boolean allowed = authorizations
    .checkPrivilege("read", "salesOrders")
    .isGranted();

// allow that ignores any conditions
boolean allowed = !authorizations
    .checkPrivilege("read", "salesOrders")
    .isDenied();
```

:::

### Spring Route Security

- Replace `SecurityExpressionHandler` with `AmsRouteSecurity` (CAP: `AmsCdsRouteSecurity`) bean in `SecurityFilterChain`.
- Update route authorization checks based on following mapping:

| v3 Route Check Syntax                                        | v4 Route Check Syntax                          |
|--------------------------------------------------------------|------------------------------------------------|
| `hasBaseAuthority("action", "resource")`                     | `precheckPrivilege("action", "resource")`      |
| `forAction("action")`                                        | `checkPrivilege("action", "*")`                |
| `forResource("resource")`                                    | `checkPrivilege("*", "resource")`              |
| `forResourceAction("resource", "action")`                    | `checkPrivilege("action", "resource")`         |
| `forResourceAction("resource", "action", attributes...)`     | use method security instead                    |

**Example**:

::: code-group
```java [v3]
@Bean
public SecurityFilterChain filterChain(
        HttpSecurity http,
        SecurityExpressionHandler<RequestAuthorizationContext> amsHttpExpressionHandler) {

    WebExpressionAuthorizationManager readOrders =
            new WebExpressionAuthorizationManager("hasBaseAuthority('read', 'orders')");
    readOrders.setExpressionHandler(amsHttpExpressionHandler);

    http.authorizeHttpRequests(authz -> authz
            .requestMatchers(GET, "/orders/**").access(readOrders));
    return http.build();
}
```

```java [v4]
@Bean
public SecurityFilterChain filterChain(HttpSecurity http, AmsRouteSecurity via) {

    http.authorizeHttpRequests(authz -> authz
            .requestMatchers(GET, "/orders/**")
                .access(via.precheckPrivilege("read", "orders")));
    return http.build();
}
```
:::

### Spring Method Security

- Replace `@PreAuthorize` annotations with v3 AMS expressions by the new AMS annotations.
- For methods with attributes, use `@AmsAttribute` on parameters to pass them to the authorization check.

| v3 Method Security Syntax                                           | v4 Method Security Syntax                                          |
|---------------------------------------------------------------------|--------------------------------------------------------------------|
| `@PreAuthorize("forAction('read')")`                                | `@CheckPrivilege(action = "read", resource = "*")`                 |
| `@PreAuthorize("forResource('products')")`                          | `@CheckPrivilege(action = "*", resource = "products")`             |
| `@PreAuthorize("forResourceAction('products', 'read')")`            | `@CheckPrivilege(action = "read", resource = "products")`          |

::: code-group
```java [v3]
@PreAuthorize("forResourceAction('products', 'read')")
public List<Product> getProducts() { ... }

@PreAuthorize("forResourceAction('products', 'read', 'product.category:string=' + #category)")
public List<Product> getProductsByCategory(@PathVariable String category) { ... }
```

```java [v4]
@CheckPrivilege(action = "read", resource = "products")
public List<Product> getProducts() { ... }

@CheckPrivilege(action = "read", resource = "products")
public List<Product> getProductsByCategory(@AmsAttribute(name = "product.category") String category) { ... }
```
:::




## Test Setup Migration

### DCL Output Directory

Replace the DCL output directory with the new default output directory for AMS DCN test resources in DCL compiler maven plugin.

| v3 Output         | v4 Output                                  |
|-------------------|--------------------------------------------|
| `target/dcl_opa/` | `target/generated-test-resources/ams/dcn/` |

### CAP Java Configuration

- Remove test sources property from `application.yaml`:

```yaml
cds:
  security:
    authorization:
      ams:
        test-sources: "" # empty uses default srv/target/dcl_opa
```

:::tip
In v4, the existence of `spring-boot-starter-ams-cap-test` on the classpath determines whether AMS will try to load local DCN. For this reason, make sure to keep it test-scoped.
:::

### Spring Security Tests

The `MockOidcTokenRequestPostProcessor.userWithPolicies` from `jakarta-ams-test` has been removed because the real AMS production code can now be tested.
It requires the definition of a [policy assignments](/Authorization/Testing#assigning-policies-to-mocked-users) map from which AMS determines the used policies based on the `app_tid` and `scim_id` claims of the token, and for advanced token flows: other claims as needed.





## Leverage new Features

### Unit Testing Policies

There is a simple new method of unit testing policy semantics without a full-blown integration test using [`ams-test`](/Libraries/java/ams-test).

### Domain-specific Authorization Checks

You can implement an `AuthorizationsAdapter<T>` to [wrap](https://github.com/SAP-samples/ams-samples-java/blob/new_lib_v4/ams-javalin-shopping/src/main/java/com/sap/cloud/security/ams/samples/auth/AuthHandler.java#L68) `Authorizations` objects with [domain-specific methods](https://github.com/SAP-samples/ams-samples-java/blob/new_lib_v4/ams-javalin-shopping/src/main/java/com/sap/cloud/security/ams/samples/auth/ShoppingAuthorizations.java#L27-L46) for [better readability](https://github.com/SAP-samples/ams-samples-java/blob/new_lib_v4/ams-javalin-shopping/src/main/java/com/sap/cloud/security/ams/samples/service/OrdersService.java#L151-L153) and reusability of authorization checks across your application.

::: tip CdsAuthorizations
The CAP Spring Boot starter already wraps the standard `Authorizations` in a `CdsAuthorizations` adapter that provides CAP-specific methods for role checks.
:::

### Simple Name Constants

Use the new `Privilege`, `AttributeName` and `PolicyName` utility classes to define constants for the action/resource combinations of your application, as well as references to DCL attributes and policies, to avoid typos and increase readability.

::: tip
There is no more need to deal with `$app` and `$env` attribute prefixes as they are inferred automatically just like in DCL. There are both factory methods for dot notation (`of`) and array notation (`ofSegments`).
:::

### Event Logging

Use the new [event logging API](/Libraries/java/ams-core#events-logging) to log authorization events as per your needs.

### Better DEBUG Logging

When things don't work as expected, you can benefit from [improved debug logging](/Troubleshooting.html#generic-privilege-check) of the AMS Client Library. By setting the log level to DEBUG for the `com.sap.cloud.security.ams` package, you can get detailed insights into the decisions that lead to the construction of the `Authorizations` object and the results of authorization checks, including which policies were evaluated and which attributes were considered.

### New TRACE Logging

To understand the detailed steps taken by the internal logic engine for evaluation of policy conditions, you can set the log level to TRACE for the `com.sap.cloud.security.ams` package. This will provide a step-by-step trace of the evaluation process, showing how conditions are built and grounded with attribute input and how the predicates were evaluated.

Additionally, it provides insights into the content of the authorization bundle.