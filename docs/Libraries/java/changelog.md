# Release Notes for AMS Client Library Java 

## Version 4

### 4.0.7

- Fix: Read Number constants from DCN as `Double` instead of `Long/Int` to avoid runtime errors when comparing with `Double` attribute input
- Fix `AmsCapAutoConfiguration`: Use @Order(-100) for `AmsUserInfoProvider` bean to make sure it runs late in the `UserInfoProvider` chain. For example, this fixes incompatibilities with DwcUserInfoProvider which must run before the `AmsUserInfoProvider` to extract user information from the token.

### 4.0.0 - 4.0.6*

Version 4 drastically changes the core API to streamline it with the Node.js library which received positive feedback since it introduced the same changes.

Instead of checking privileges on a `PolicyDecisionPoint` with an `Attributes` object, an `AuthorizationsProvider` prepares an `Authorizations` object for the same purpose. This separates *what* to check from *how* to check it. The necessary configuration for advanced authorization scenarios such as principal propagation or non-standard authorization strategies are configured once 
during application start. As a result, the authorization checks themselves remain straight-forward in version 4, with a focus on the application domain.

\* *first publicly available 4.x release version is 4.0.5*

### New features

- Spring Boot 3 **and** 4 support
- Official XSUAA legacy support via `HybridAuthorizationsProvider`
- Zero-Trust-Identity-Service (**ZTIS**) certificate support
::: tip ZTIS Auto-Configuration
There is out-of-the-box support for ZTIS service bindings via the Spring Boot starters.
:::
- Domain-Specific `Authorizations` by [wrapping](https://github.com/SAP-samples/ams-samples-java/blob/main/ams-javalin-shopping/src/main/java/com/sap/cloud/security/ams/samples/auth/AuthHandler.java#L68) `Authorizations` objects with [domain-specific methods](https://github.com/SAP-samples/ams-samples-java/blob/main/ams-javalin-shopping/src/main/java/com/sap/cloud/security/ams/samples/auth/ShoppingAuthorizations.java#L27-L46) for [better readability](https://github.com/SAP-samples/ams-samples-java/blob/main/ams-javalin-shopping/src/main/java/com/sap/cloud/security/ams/samples/service/OrdersService.java#L151-L153) and reusability of authorization checks across your application.
::: tip CdsAuthorizations
The CAP Spring Boot starter already wraps the standard `Authorizations` in a `CdsAuthorizations` adapter that provides CAP-specific methods for role checks.
:::
- Provided [CAP Spring beans](/Libraries/java/cap-ams.html#auto-configuration) for custom authorization checks
- Improved [Spring Security beans](/Libraries/java/spring-boot-ams#auto-configuration) for custom authorization checks
- New [event logging API](/Libraries/java/ams-core#events-logging) for logging authorization events
- Configuration options for [technical communication](/Authorization/TechnicalCommunication) scenarios via SAP Identity Service
- Customization of authorization strategy via `AuthorizationsProvider` interface, e.g. [custom user attribute injection](/Authorization/AuthorizationChecks#overriding-methods)
- JUnit 5+ extension for unit testing policy semantics without a full-blown integration test using [`ams-test`](/Libraries/java/ams-test).
- Detailed [**DEBUG**](/Troubleshooting) logging about construction of `Authorizations` from token
- **TRACE** logging of authorization bundle content and logic engine evaluations, showing how conditions are built and grounded with attribute input and how the predicates were evaluated
- Drastically reduced number of authorization checks in CAP requests, which improves debug log analysis
- New `Privilege`, `AttributeName` and `PolicyName` utility classes to define constants for the action/resource combinations of your application, as well as references to DCL attributes and policies, to avoid typos and increase readability.

::: tip
There is no more need to deal with `$app` and `$env` attribute prefixes as they are inferred automatically just like in DCL. There are both factory methods for dot notation (`of`) and array notation (`ofSegments`).
:::

### Removed Features

- Audit Logging library integration has been removed in favor of general event logging via the new API. Please refer to the [migration guide](/Libraries/java/v3/migration-v3-to-v4#event-logging) for details.

### Breaking Changes
CAP Java Applications should only need to do trivial changes when updating to version 4 unless they used `spring-ams`.

Please refer to the [migration guide](/Libraries/java/v3/migration-v3-to-v4) for details.

### Performance

Our performance tests have indicated that the performance impact of authorization checks with the AMS library was already negligible before. Although there are improvements in version 4, such as a reduction of redundant authorization checks in CAP applications, we did not measure a significant performance impact.

For example, for both library versions, the request latency for a localhost CAP OData endpoint with instance-based authorization filters was `<= 5ms` of which most of the time was likely spent on database and network handling instead of the AMS library.

## Version 3

### 3.8.0

- This release removes the dependencies from `com.sap.cloud.security.ams.dcl` artifacts. All required classes,
interfaces, etc., are now part of the `jakarta-ams` module using the same packages. So, everything should continue
to work without any changes. Please remove any direct dependencies on `com.sap.cloud.security.ams.dcl` artifacts.

### 3.7.0

- Maintenance release with updated dependencies and fixes for the Maven Central release process.

### 3.6.0

- The property `cds.security.mock.enabled` is now used to enable the mock users in the
  `cap-ams-support` module.
- A new property `ams.properties.bundleGatewayUpdater.maxFailedUpdates` is introduced to configure the maximum
  number of failed updates of the bundle gateway before it logs an error message. The default value is `0`.
