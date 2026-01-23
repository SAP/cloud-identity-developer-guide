# jakarta-ams

`jakarta-ams` is the core module used to perform [authorization checks](https://help.sap.com/docs/identity-authentication/identity-authentication/configuring-authorization-policies?locale=en-US) in applications based on AMS policies.

::: tip CAP applications
CAP applications use this library transitively via `cap-ams`. The CAP module performs authorization checks automatically under the hood, so you typically don't need to use jakarta-ams directly.
:::

::: tip Spring Boot applications
Spring Boot applications use this library transitively via `spring-ams`. The Spring module brings support for declarative method- and route-level authorization checks, auto-configuration for the AMS module, and beans for simplified usage.
:::

## Installation

The module is available as a Maven dependency. Add the following to your `pom.xml`:

```xml
<dependency>
    <groupId>com.sap.cloud.security.ams</groupId>
    <artifactId>jakarta-ams</artifactId>
    <version>${sap.cloud.security.ams.version}</version>
</dependency>
```

::: danger Important
Always keep the version of this dependency up-to-date since it's a **crucial part of your application's security**:

```bash
mvn versions:use-latest-versions -Dincludes=com.sap.cloud.security.ams:jakarta-ams
```
This is especially important when you deploy your application with a locked dependency version.
:::

:::tip
To check which version of the module is installed, run:

```bash
mvn dependency:tree -Dincludes=com.sap.cloud.security.ams:jakarta-ams
```
This prints a dependency tree that shows which versions of the module are installed, including transitive dependencies.
:::

## Usage

The following snippets show you how to use the core API of the library. For more details on the authorization concepts, see [Authorization Checks](/Authorization/AuthorizationChecks).

### Setup

#### Initializing AuthorizationManagementService

The `AuthorizationManagementService` is the central entry point for performing authorization checks. It manages the policy bundle and provides `Authorizations` objects for checking privileges.

```java
import com.sap.cloud.security.ams.v4.api.*;
import com.sap.cloud.security.ams.v4.core.*;
import com.sap.cloud.security.config.*;
import com.sap.cloud.environment.servicebinding.api.*;

AuthorizationManagementService ams;

// For production: Initialize from SAP Identity BTP Service binding
ServiceBinding identityBinding = DefaultServiceBindingAccessor.getInstance()
    .getServiceBindings().stream()
    .filter(b -> "identity".equals(b.getServiceName().orElse("").toLowerCase()))
    .findFirst()
    .orElseThrow(() -> new IllegalStateException("No identity service binding found"));

// Alternative for local testing: Initialize from pre-compiled DCN directory    
AuthorizationManagementService ams = AuthorizationManagementServiceFactory.fromLocalDcn(
    Path.of("./test/dcn"),  // your compile target directory
    new LocalAuthorizationManagementServiceConfig()
        .withPolicyAssignmentsPath(Path.of("./test/mockPolicyAssignments.json"))
);
```

::: warning Important
Implement a [startup check](/Authorization/StartupCheck) to ensure that the `AuthorizationManagementService` instance is ready for authorization checks before serving authorized endpoints.
:::

## Configuration

## Error Handling

The `AuthorizationManagementService` emits error events for bundle loading failures. You can register listeners to handle these errors:

```java
ams.addErrorListener(event -> {
    if (event instanceof BundleInitializationErrorEvent) {
        // Initial bundle download failed - service not ready
        logger.error("AMS bundle initialization failed", event.getError());
    } else if (event instanceof BundleRefreshErrorEvent refreshError) {
        // Bundle refresh failed - service remains ready with stale data
        logger.warn("AMS bundle refresh failed (age: {} seconds)", 
            refreshError.getSecondsSinceLastRefresh(), event.getError());
    }
});
```

::: tip Handling Initial Bundle Load Errors
Refer to the [Startup Check](/Authorization/StartupCheck) documentation for guidance on how to react when AMS fails to initialize the bundle.
:::

## API Reference

### AuthorizationManagementServiceFactory

Static factory methods for creating `AuthorizationManagementService` instances:

| Method | Description |
|--------|-------------|
| `fromIasServiceConfiguration(config)` | Creates instance from SAP Identity Service OAuth2 configuration (production) |
| `fromIasServiceConfiguration(config, amsConfig)` | Same as above with custom AMS configuration |
| `fromIdentityServiceBinding(binding)` | Creates instance from SAP Identity Service binding obtained via BTP Service Binding Library (production) |
| `fromIdentityServiceBinding(binding, config)` | Same as above with custom configuration |
| `fromLocalDcn(path)` | Creates instance from local DCN directory (testing) |
| `fromLocalDcn(path, config)` | Same as above with custom configuration |

### AuthorizationManagementService

Central service for managing authorization state and providing authorizations.

#### Lifecycle Methods

| Method | Description |
|--------|-------------|
| `void start()` | Start the service and begin loading the policy bundle |
| `void stop()` | Stop the service and release all resources |
| `boolean isReady()` | Check if the service is ready for authorization checks |
| `CompletableFuture<Void> whenReady()` | Returns a future that completes when ready |

#### Authorization Methods

| Method | Description |
|--------|-------------|
| `Authorizations getAuthorizations(Set<String> policies)` | Get an `Authorizations` instance for the specified policies |

#### Error Handling

| Method | Description |
|--------|-------------|
| `void addErrorListener(Consumer<AmsBackgroundException> listener)` | Register an error event listener |
| `void removeErrorListener(Consumer<AmsBackgroundException> listener)` | Remove an error event listener |

### Authorizations

Interface for performing authorization checks based on a set of policies.

#### Privilege Check Methods

| Method | Description |
|--------|-------------|
| `Decision checkPrivilege(String action, String resource)` | Check if privilege is granted with no additional input |
| `Decision checkPrivilege(String action, String resource, Map<String, Object> input)` | Check if privilege is granted with attribute input |
| `Decision checkPrivilege(String action, String resource, Map<String, Object> input, Set<String> unknowns)` | Check with input and filter condition to contain only specified unknowns |
| `Decision checkPrivilege(Privilege privilege)` | Check using a Privilege object |
| `Decision checkPrivilege(Privilege privilege, Map<String, Object> input)` | Check using a Privilege object with input |
| `Decision checkPrivilege(Privilege privilege, Map<String, Object> input, Set<String> unknowns)` | Check using a Privilege object with input and unknowns |

#### Query Methods

| Method | Description |
|--------|-------------|
| `Set<String> getPotentialActions(String resource)` | Get all actions potentially granted for a resource (ignoring conditions) |
| `Set<String> getPotentialResources()` | Get all resources with at least one potentially granted action |
| `Set<Privilege> getPotentialPrivileges()` | Get all potentially granted action/resource combinations |

#### Configuration Methods

| Method | Description |
|--------|-------------|
| `Set<String> getPolicies()` | Get the set of policies used for authorization checks |
| `void setPolicies(Set<String> policies)` | Set the policies to use for authorization checks |
| `Map<String, Object> getDefaultInput()` | Get default input used for all checks |
| `void setDefaultInput(Map<String, Object> input)` | Set default input used for all checks |
| `Authorizations getLimit()` | Get the authorizations that limit this instance |
| `void setLimit(Authorizations limit)` | Limit this instance's authorizations to another instance's authorizations |

### Decision

Result of an authorization check, representing granted, denied, or conditional access.

| Method | Description |
|--------|-------------|
| `boolean isGranted()` | Returns true if unconditionally granted |
| `boolean isDenied()` | Returns true if unconditionally denied |
| `boolean isConditional()` | Returns true if access depends on a condition |
| `<T> T visit(DcnConditionVisitor<T> visitor)` | Visit the condition tree to extract/transform it (e.g., to SQL) |

### DcnConditionVisitor

Visitor interface for transforming condition trees into other representations (e.g., SQL WHERE clauses, predicates).

```java
public interface DcnConditionVisitor<T> {
    T visitCall(String operator, List<T> arguments);
    T visitAttributeName(AttributeName attributeName);
    T visitValue(Object value);
}
```

Built-in implementations:
- `SqlExtractor` - Converts conditions to SQL WHERE clauses

### Supporting Types

#### Privilege

Represents an action/resource pair:

```java
public class Privilege {
    private final String action;
    private final String resource;
    
    public Privilege(String action, String resource) { ... }
    public String getAction() { ... }
    public String getResource() { ... }
}
```

#### AttributeName

Represents a fully-qualified attribute name in the AMS schema:

```java
public class AttributeName {
    public String getFullyQualifiedName() { ... }
    // e.g., "$app.product.category"
}
```