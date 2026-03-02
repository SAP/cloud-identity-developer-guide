# spring-boot-ams

This `spring-boot-ams` module integrates AMS into Spring applications for standard `Spring Security` based
enforcement
of authorization.

## Installation

Use the Spring Boot starter module:

```xml

<dependency>
    <groupId>com.sap.cloud.security.ams</groupId>
    <artifactId>spring-boot-starter-ams</artifactId>
</dependency>
```

::: tip CAP applications
For CAP applications, use `spring-boot-starter-cap-ams` instead.
:::

## Public API

The following classes are part of the stable public API and can be freely used by consumer applications.

- `com.sap.cloud.security.ams.spring.AmsRouteSecurity`
- `com.sap.cloud.security.ams.spring.annotations.CheckPrivilege`
- `com.sap.cloud.security.ams.spring.annotations.PrecheckPrivilege`
- `com.sap.cloud.security.ams.spring.annotations.AmsAttribute`

::: warning Semantic Versioning Notice
Classes and packages **not listed above** are internal implementation details. They may change, be renamed, or be removed in minor or patch releases without notice. Do not depend on internal classes in production code.
:::

## Auto-Configuration

The starter automatically configures:

- `AuthorizationManagementService` bean from SAP Identity Service binding
- Request-scoped `Authorizations` proxy bean for the current request
- `AmsRouteSecurity` bean for route-level authorization
- `AmsMethodSecurity` bean for method-level authorization
- Readiness state integration (via `spring-boot-starter-ams-readiness`)

## Declarative Authorization

The starter provides two approaches for declarative authorization checks in addition
to [programmatic checks](/Authorization/AuthorizationChecks):

### Route-Level Security

Use `AmsRouteSecurity` to secure HTTP endpoints in your Spring Security configuration:

```java

@Bean
public SecurityFilterChain filterChain(HttpSecurity http, AmsRouteSecurity ams) throws Exception {
    http.authorizeHttpRequests(authz -> authz
            // Require unconditional privilege grant
            .requestMatchers(DELETE, "/products/**").access(ams.checkPrivilege("delete", "products"))
            // Allow conditional access (service layer must enforce conditions)
            .requestMatchers(GET, "/products/**").access(ams.precheckPrivilege("read", "products"))
            .anyRequest().authenticated()
    );
    return http.build();
}
```

| Method                                | Description                                                         |
|---------------------------------------|---------------------------------------------------------------------|
| `checkPrivilege(action, resource)`    | Requires unconditional grant. Rejects conditional or denied access. |
| `precheckPrivilege(action, resource)` | Allows conditional access. Only rejects definitely denied access.   |

::: warning
When using `precheckPrivilege`, your service layer **must** perform additional contextual authorization checks to
enforce any conditions on the access.
:::

### Method-Level Security

Enable method security in your configuration:

```java

@Configuration
@EnableWebSecurity
@EnableMethodSecurity
public class SecurityConfiguration {
    // ...
}
```

Use `@CheckPrivilege` and `@PrecheckPrivilege` annotations on service methods:

```java

@Service
public class ProductsService {

    @CheckPrivilege(action = "read", resource = "products")
    public List<Product> getProducts() {
        return database.getProducts();
    }

    @CheckPrivilege(action = "delete", resource = "products")
    public void deleteProduct(int productId) {
        database.deleteProduct(productId);
    }
}
```

Use `@AmsAttribute` to pass method parameters as input for condition evaluation:

```java

@Service
public class OrdersService {

    @CheckPrivilege(action = "create", resource = "orders")
    public Order createOrder(
            Product product,
            int quantity,
            @AmsAttribute(name = "order.total") double totalAmount,
            @AmsAttribute(name = "product.category") String productCategory) {
        // Method executes only if privilege is granted for the given attribute values
        // ...
    }
}
```

| Annotation           | Description                                                               |
|----------------------|---------------------------------------------------------------------------|
| `@CheckPrivilege`    | Requires unconditional grant. Use for actions that cannot be filtered.    |
| `@PrecheckPrivilege` | Allows conditional access. Use when service layer can enforce conditions. |
| `@AmsAttribute`      | Maps method parameters to AMS schema attributes for condition evaluation. |

::: danger AOP Proxy Limitation
Method security annotations **do not apply** when methods within the same class call each other via `this`. This is a
general Spring AOP limitation, not AMS-specific. Internal method calls bypass the proxy that enforces the security
annotations.

```java

@Service
public class OrdersService {

    @CheckPrivilege(action = "delete", resource = "orders")
    public void deleteOrder(int orderId) {
        // Security check is enforced when called from outside
    }

    public void processOrders() {
        // ⚠️ This internal call bypasses security!
        this.deleteOrder(123);
    }
}
```

To ensure security checks are applied, either:

- Call secured methods from outside the class (via the Spring proxy)
- Inject the service into itself and call via the injected reference
- Use programmatic authorization checks with the `Authorizations` bean
  :::

## Configuration Properties

Configure the starter in `application.yml`:

```yaml
sap:
  ams:
    bundle-loader:
      polling-interval: 20000      # Bundle update polling interval in ms (default: 20000)
      initial-retry-delay: 1000    # Initial retry delay after failure in ms (default: 1000)
      max-retry-delay: 20000       # Maximum retry delay in ms (default: 20000)
      retry-delay-factor: 2        # Exponential backoff factor (default: 2)

    method-security:
      enabled: true      # Enable method-level security (default: true)

    actuator:
      health:
        enabled: true    # Enable health indicator (default: true)
      readiness:
        enabled: true    # Enable readiness state contributor (default: true)
```