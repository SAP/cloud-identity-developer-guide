# spring-ams

## Spring Boot Starter

::: tip Spring Boot Auto-Configuration
Add a dependency to the [spring-boot-starter-ams](./spring-boot-starter-ams.md) module instead of `spring-ams` for auto-configuration with Spring Boot applications.
:::

## Spring Route-Level Security

ToDo

## Spring AOP Method Security with Type-Safe Annotations

#### Examples

**`@CheckPrivilege`** - Type-safe annotation for privilege checks:
```java
@CheckPrivilege(action = "delete", resource = "orders")
public void deleteOrder(@AmsAttribute(name = "id") String orderId) {
    // Only executes if user has delete:orders privilege for this order
}
```

**`@PrecheckPrivilege`** - Type-safe annotation for pre-checks:
```java
@PrecheckPrivilege(action = "read", resource = "orders")
public List<Order> getOrders() {
    // Allows conditional access, filter in service layer
}
```

**`@AmsAttribute`** - Marks parameters to pass as instance attributes:
```java
@CheckPrivilege(action = "update", resource = "orders")
public void updateOrder(
        @AmsAttribute(name = "id") String orderId,
        @AmsAttribute(name = "status") String status) {
}
```

### Implementation Details

The annotations use Spring's **meta-annotation pattern**:

```java
@PreAuthorize("@methodSecurity.checkPrivilege(#root)")
public @interface CheckPrivilege {
    String action();
    String resource();
}
```

When Spring Security processes the method:
1. `@PreAuthorize` calls `methodSecurity.checkPrivilege(#root)`
2. The `#root` SpEL variable contains the `MethodSecurityExpressionOperations`
3. `MethodSecurity` extracts the `MethodInvocation` from the root
4. It reads the `@CheckPrivilege` annotation metadata
5. It extracts `@AmsAttribute` parameters
6. It calls `Authorizations.checkPrivilege()` with the extracted data

### Setup

#### Enable Method Security

```java
@Configuration
@EnableWebSecurity
@EnableMethodSecurity  // Required
public class SecurityConfiguration {
    // ... your config
}
```

#### MethodSecurity bean

The `MethodSecurity` bean of `spring-ams` is automatically registered when you use the [spring-boot-starter-ams](./spring-boot-starter-ams.md) module. It is mandatory for the annotations to work.