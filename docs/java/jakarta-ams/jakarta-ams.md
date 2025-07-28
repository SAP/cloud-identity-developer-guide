# jakarta-ams

# Authorization Management Service Client Library for Jakarta EE Applications

The client library for the Authorization Management Service (AMS) enables Jakarta EE applications to make authorization decisions based on user policies.
Users are identified using JWT bearer tokens provided by either the SAP Business Technology Platform `xsuaa` or `identity` service.
To specify these policies, the application must use a Data Control Language (DCL). Administrators can customize and combine policies and assign them to users.
The library also supports controlling access to specific actions and resources, and allows for easy implementation of instance-based access control by specifying resource-specific attributes such as "cost center" or "country code".

## <a id="api-disclaimer"></a>Disclaimer on API Usage

This documentation provides information that might be useful in using Authorization Management Service. We will try to ensure that future versions of the APIs are backwards compatible to the immediately preceding version. This is also true for the API that is exposed with `com.sap.cloud.security.ams.dcl.client` package and its subpackages.\
Please check the [release notes](../releases.md) to stay tuned about new features and API modifications.

## Requirements

- Java 17
- Tomcat 10 or any other servlet that implements specification of the Jakarta EE platform

## Sample DCL

Every business application has to describe its security model in [dcl language](../../concepts/DeployDCL.md).<br>
To describe the names and types of attributes, a file named `schema.dcl` and must be located in the root folder, for example in `src/main/resources`.

```
SCHEMA {
    salesOrder: {
            type: number
    },
    CountryCode: String
}
```

Furthermore, the security model comprises rules and policies that can be used or enriched by the administrator to grant users permissions. They can be organized in packages.

```
POLICY readAll {
    GRANT read ON * WHERE CountryCode IS NOT RESTRICTED;
}

POLICY readAll_Europe {
    USE readAll RESTRICT CountryCode IN ('AT', 'BE', 'BG', ...);
}

POLICY adminAllSales {
    GRANT read, write, delete, activate ON salesOrders, salesOrderItems;
}

POLICY anyActionOnSales {
    GRANT * ON salesOrders, salesOrderItems;
}

POLICY readSalesOrders_Type {
    GRANT read ON salesOrders WHERE salesOrder.type BETWEEN 100 AND 500;
}

POLICY readAllOwnItems {
    GRANT read ON * WHERE author.createdBy = $user.user_uuid;
}
```

**Note:** As of version `0.9.0` the DCL compiler adds this `$user` default attributes, in case there is no `$user` present:

```TEXT
SCHEMA {
  ...
  $user: {
    user_uuid: String,
    email: String,
    groups: String[]
  },
  ...
}
```

The client libraries then will map the token claims of the same name to the respective attributes for the Policy Decision Point request.
Consequently, the attribute like `$user.user_uuid` is prefilled and can be overwritten with a custom `AttributesProcessor` as described [here](./jakarta-ams.md#customize-attributes).

## Configuration

### Maven Dependencies

```xml
<dependency>
    <groupId>com.sap.cloud.security.ams.client</groupId>
    <artifactId>jakarta-ams</artifactId>
    <version>${sap.cloud.security.ams.client.version}</version>
</dependency>
```

## Usage

### Setup PolicyDecisionPoint

```java
PolicyDecisionPoint policyDecisionPoint = PolicyDecisionPoint.create(DEFAULT); // use DclFactoryBase.DEFAULT
```

### Initial application startup check

Initial startup check will check if the bundle containing the defined policies are accessible to the application at the application initialization step.
This functionality is provided by `PolicyDecisionPoint` and upon application startup it will perform periodical health checks until `OK` state is reached or until the `startupHealthCheckTimeout` has elapsed, then it will throw a `PolicyEngineCommunicationException`.

Startup check behaviour can be configured by the following 2 parameters:

- `STARTUP_HEALTH_CHECK_TIMEOUT` - max wait time for startup check in seconds. 0 or negative values disable the startup check. :bulb: If left unset uses default value of 30 seconds
- `FAIL_ON_STARTUP_HEALTH_CHECK` - boolean whether the application should start, if policies are not available i.e. `PolicyEngineCommunicationException` won't be thrown. :bulb: If left unset uses default value of `true`

like so:

```java
//Sets the max health check wait time to 10 seconds and disables error throwing if OK state isn't reached
PolicyDecisionPoint policyDecisionPoint = PolicyDecisionPoint.create(DEFAULT, 
        STARTUP_HEALTH_CHECK_TIMEOUT, 10L, // use PolicyDecisionPoint.Parameters.STARTUP_HEALTH_CHECK_TIMEOUT and PolicyDecisionPoint.Parameters.FAIL_ON_STARTUP_HEALTH_CHECK to configure PolicyDecisionPoint with startup check parameters
        FAIL_ON_STARTUP_HEALTH_CHECK, false
        );
```

### Overview `PolicyDecisionPoint` methods

| method                                        | description                                                                                                                     |
| :-------------------------------------------- | :------------------------------------------------------------------------------------------------------------------------------ |
| `boolean getHealthStatus().getHealthState()`  | Returns policy decision point readiness status. Use this within your health endpoint.                                           |
| `boolean getHealthStatus().getBundleStatus()` | Returns detailed information about the configured bundles. Use this for advanced health checks for your application.            |
| `boolean allow(Attributes attributes)`        | Returns `true` if the user is authorized for `action`, `resource` and `attributes`. <br> Throws an `PolicyEvaluationException`. |

Further details and further methods are documented in JavaDoc.

### Exceptions

The `PolicyDecisionPoint` raises the following unchecked `RuntimeException`s:

- `IllegalArgumentException` or `NullpointerException`:\
  if the policy decision point is parameterized wrongly.
- `PolicyEvaluationException`:\
  represents all kind of issues that occurs during policy evaluation and interaction with the policy engine, e.g. misconfiguration, unknown dcl package or lack of information provided to perform the policy evaluation.
- `PolicyEngineCommunicationException` can be thrown at the application startup indicating that the `PolicyDecisionPoint` bundles are not ready to be used for authorization evaluation. This exception will be thrown only if the `PolicyDecisionPoint` startup check is enabled (`STARTUP_HEALTH_CHECK_TIMEOUT` < 0) and if the `FAIL_ON_STARTUP_HEALTH_CHECK` == `true`

### Is Policy Engine in healthy state?

```java
if (policyDecisionPoint.getHealthStatus().getHealthState() == HealthState.OK) {
       // The HealthState is OK
}
```

### Allow

#### Pre-fill user-attributes from OIDC token

In case you use `com.sap.cloud.security:java-api` or `com.sap.cloud.security:java-security` for token validation you can simply create a `Principal` instance within the same thread. It derives the principal information from the OIDC token, that is stored in `SecurityContext`:

```java
Principal principal = Principal.create();
```

> Alternatively you can also build the `Principal` using the `PrincipalBuilder`.\
> E.g. `PrincipalBuilder.create("the-zone-id", "the-user-id").build();`

#### Has user from OIDC token `read` access to any resources?

```java
Attributes attributes = principal.getAttributes()
                .setAction("read");

boolean isUserAuthorized = policyDecisionPoint.allow(attributes);
```

#### Has user `read` access to `salesOrders` resource?

```java
Attributes attributes = principal.getAttributes()
                .setAction("read")
                .setResource("salesOrders");

boolean isUserAuthorized = policyDecisionPoint.allow(attributes);
```

#### Has user `read` access to `salesOrders` resource with `CountryCode` = "DE" and `salesOrder.type` = 4711?

```java
Attributes attributes = principal.getAttributes()
                .setAction("write")
                .setResource("salesOrders")
                .app().value("CountryCode", "IT")
                      .structure("salesOrder").value("type", 4711)
                .attributes();

boolean isUserAuthorized = policyDecisionPoint.allow(attributes);
```

> In case the Attribute is of type `number` its relevant that you pass the value as Integer or Double but NOT as String object.

### Customize Attributes

It is possible to modify the `Attributes` that are sent as request to the `PolicyDecisionPoint`, by implementing the `AttributesProcessor`
interface whose `processAttributes` method is called _each_ time when the `Principal` object is built to fill its' field `attributes`.
The attributes can be later accessed by calling `Principal.getAttributes()` method.

For example, this can be used to explicitly specify the list of policies that are used during evaluations for the given `Principal`.

Implementations of `AttributesProcessor` can be registered using Java's `ServiceLoader` mechanism as follows:

- Create an SPI configuration file with name `com.sap.cloud.security.ams.api.AttributesProcessor` in the `src/main/resources/META-INF/services` directory.
- Enter the fully qualified name of your `AttributesProcessor` implementation class, e.g. `com.mypackage.MyCustomAttributesProcessor`.
- The implementation could look like this:

```java
public class MyCustomAttributesProcessor implements AttributesProcessor {

    @Override
    public void processAttributes(Principal principal) {
        principal.getAttributes().app().value("customAttribute", "theAttributeValue").attributes();
        
        // simplified - without NullPointer check
        HashMap<String, Object> userMap = new HashMap<>();
    userMap.put("groups", SecurityContext.getToken().getClaimAsStringList("groups"));
    principal.getAttributes().env().value("$user", userMap).attributes();
    }
    
    // ----Optional------ 
    // @Override
    // public int getPriority() {
    //   return AttributesProcessor.HIGHEST_PRECEDENCE;
    // }
}
```

Note that more than one evaluation may be performed for the same Principle, so the implementations of `AttributesProcessor` should be [idempotent](https://en.wikipedia.org/wiki/Idempotence) to prevent issues when called more than once for the same Principle.

If multiple `AttributesProcessor` implementations are loaded, the optional `getPriority` method can be implemented to take control over the order in which the implementations are called.

> :heavy_exclamation_mark: [Consider also limitation of API liability](#api-disclaimer).
>
> :warning: It is not possible to use `AttributesProcessor` with classes managed by Spring dependency injection mechanism(e.g. Autowired beans),
> as the lifecycle of the ServiceLoader and Spring DI mechanisms are different.

## Test utilities

You can find the test utilities documented in the `java-ams-test` module.

## Logging

Additionally `jakarta-ams` client library provides these two predefined loggers:

1. `PolicyEvaluationSlf4jLogger` logs policy evaluation results including their requests using the
   SLF4J logger. <br>This logger is by default registered as listener to the
   default `PolicyDecisionPoint implementations such as "default" and "test-server".
2. `PolicyEvaluationV2AuditLogger` writes audit logs for all policy evaluation results
   using SAP internal `com.sap.xs.auditlog:audit-java-client-api` client library (version 2).
   <br>The purpose of that data access audit message is to complement application's audit logs with
   the policy evaluation result documenting "WHETHER" and "WHY" the policy engine has granted or denied access.

#### Application Logging

This library uses [slf4j](http://www.slf4j.org/) for logging. It only ships the [slf4j-api module](https://mvnrepository.com/artifact/org.slf4j/slf4j-api) and no actual logger implementation.
For the logging to work `slf4j` needs to find a valid logger implementation at runtime.
If your application is deployed via SAP Java buildpack then you will have one available and logging should just work.

Otherwise, you need to add a logger implementation dependency to your application.
See the slf4j [documentation](http://www.slf4j.org/codes.html#StaticLoggerBinder)
for more information and a [list](http://www.slf4j.org/manual.html#swapping) of available logger options.

In detail, `PolicyEvaluationSlf4jLogger` writes

- `ERROR` messages only in exceptional cases, e.g. unavailability of the policy engine
- `WARN` messages in case of denied policy evaluations
- `INFO` messages for all other non-technical policy evaluation results.

> :bulb: It's recommended to accept messages of severity `INFO` for this logger `com.sap.cloud.security.ams.logging.*`.

#### Audit Logging

In order to enable that audit logging the application needs to configure and register this `PolicyEvaluationV2AuditLogger` as listener to the `PolicyDecisionPoint` implementation as done in the [jakarta-security-ams sample](https://github.com/SAP-samples/ams-samples-java/blob/main/jakarta-ams-sample/src/main/java/com/sap/cloud/security/samples/filter/PolicyDecisionAuditLogFilter.java).<br>
In order to correlate this audit log message with logs written for the same request context `PolicyEvaluationV2AuditLogger` also fills `sap-passport` if provided with the Mapped Diagnostic Context (MDC) context. For that applications have to

- leverage a slf4j implementation that supports MDC like [logback](http://logback.qos.ch/manual/mdc.html)
- provide dependencies to [Audit Log Service Java Client](https://github.wdf.sap.corp/xs-audit-log/audit-java-client)
- fetch `sap_passport e.g. from http header
- enrich the MDC as following:<br>

```java
MDC.put("sap-passport", ((HttpServletRequest) request).getHeader("sap_passport"));
```

Finally, your application needs to be bound to an audit log service instance:<br>

```
cf create-service auditlog standard <the-service-instance>
cf bind-service <the-application> <the-service-instance>
```

## Common Pitfalls

#### The `PolicyDecisionPoint.allowFilterClause()` returns a `FilterClause` with 'false' as Condition:

```
Unknowns:
  ["$app"."user"]
Ignores:
  []
Condition:
  false
```

... instead of the expected `Call` object, that can be mixed into SQL statement:

```
Unknowns:
  ["$app"."user"]
Ignores:
  []
Condition:
  or(
    eq("$app"."user"."createdBy", the-user-id)
    eq("$app"."user"."updatedBy", the-user-id)
  )
```

#### The `PolicyDecisionPoint.allowFilterClause()` returns a `FilterClause` with 'true' as Condition instead of the expected `Call` object.

This is not an error, in that case the access is not limited by a condition.

## Samples

- [jakarta-security sample](https://github.com/SAP-samples/ams-samples-java/tree/main/jakarta-ams-sample)
