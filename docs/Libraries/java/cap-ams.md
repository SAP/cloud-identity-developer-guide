# cap-ams

This `cap-ams` module integrates AMS with CAP Java applications.

## Installation

Use the Spring Boot starter module:

```xml
<dependency>
    <groupId>com.sap.cloud.security.ams</groupId>
    <artifactId>spring-boot-starter-cap-ams</artifactId>
</dependency>
```

::: tip
For non-Spring-Boot CAP Java applications, please open a support ticket to discuss integration options. The `cap-ams` module does not require Spring Boot, but it is tough to provide a flexible starter with replaceable, configurable bootstrapping without Spring's dependency injection and configuration features.
:::

## Auto-Configuration

The starter automatically configures:

- `AuthorizationManagementService` bean from SAP Identity Service binding
- CAP authorization integration for automatic role computation and AMS filter annotations
- Request-scoped `CdsAuthorizations` proxy for manual authorization checks of the current user
- Readiness state integration (via `spring-boot-starter-ams-readiness`)

## CAP Authorization Integration

See the [CAP Integration](/CAP/Basics) documentation for the relationship between cds annotations and enforcement with authorization policies.

## Configuration Properties

Configure the starter in `application.yaml` using CDS properties:

```yaml
cds:
  security:
    authorization:
      ams:
        edge-service:
          url: http://localhost:8080   # Edge service URL (optional)
        bundle-loader:
          polling-interval: 20000      # Bundle update polling interval in ms (default: 20000)
          initial-retry-delay: 1000    # Initial retry delay after failure in ms (default: 1000)
          max-retry-delay: 20000       # Maximum retry delay in ms (default: 20000)
          retry-delay-factor: 2        # Exponential backoff factor (default: 2)
        features:
          generateExists: true         # Generate EXISTS predicates for filter attributes behind 1:N associations (default: true)
```
