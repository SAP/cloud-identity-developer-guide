# Spring Boot Starter for Authorization Management Service (AMS)

This Spring Boot Starter provides auto-configuration for integrating Spring Boot applications with the SAP Authorization Management Service (AMS).

:::warning
For CAP applications, use 
:::

## Overview

This module contains the auto-configuration classes that enable seamless integration of AMS with Spring Boot applications. It automatically configures:

- `AuthorizationManagementService` bean
- `AuthProvider` implementation
- Request-scoped `Authorizations` proxy
- AMS route security components
- JWT authentication converter for Spring Security
- Health indicator for Spring Boot Actuator (optional)
- Method security integration

## Maven Dependency

```xml
<dependency>
    <groupId>com.sap.cloud.security.ams.client</groupId>
    <artifactId>spring-boot-starter-ams</artifactId>
    <version>${sap.cloud.security.ams.client.version}</version>
</dependency>
```

## Configuration

The starter can be configured via application properties:

```yaml
sap:
  spring:
    ams:
      enabled: true                    # Enable/disable AMS autoconfiguration (default: true)
      start-immediately: true          # Auto-start AMS service on application startup (default: true)
      startup-check-enabled: true      # Perform readiness check during startup (default: true)
      startup-timeout-seconds: 30      # Timeout for startup readiness check (default: 30)
```

## Disabling Auto-configuration

To disable the auto-configuration:

```java
@SpringBootApplication(exclude = {
    AmsAutoConfiguration.class,
    AmsMethodSecurityAutoConfiguration.class,
    AmsHealthIndicatorAutoConfiguration.class
})
public class Application {
    // ...
}
```

Or via properties:

```yaml
sap:
  spring:
    ams:
      enabled: false
```

## What's Included

This starter includes:

1. **Core AMS Spring Integration** (`spring-ams`) - The base library with authorization components
2. **Auto-configuration Classes**:
   - `AmsAutoConfiguration` - Main configuration for AMS beans
   - `AmsMethodSecurityAutoConfiguration` - Method-level security integration
   - `AmsHealthIndicatorAutoConfiguration` - Health indicator for Actuator
   - `AmsZtisConfiguration` - Optional ZTIS integration

## Requirements

- Java 17+
- Spring Boot 3.5.0+
- Spring Security 6.5.0+

## For More Information

For detailed usage instructions and examples, see the [spring-ams module documentation](../spring/README.md).