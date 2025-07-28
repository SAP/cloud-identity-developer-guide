# Authorization Management Service Client Library for Spring Boot Applications
The client library for the Authorization Management Service (AMS) enables Spring Boot 3 applications to make authorization decisions based on user policies.
Users are identified using JWT bearer tokens provided by either the SAP Business Technology Platform `xsuaa` or `identity` service.
To specify these policies, the application must use a Data Control Language (DCL). Administrators can customize and combine policies and assign them to users.
The library also supports controlling access to specific actions and resources, and allows for easy implementation of instance-based access control by specifying resource-specific attributes such as "cost center" or "country code".

### Supported Environments
- Cloud Foundry
- Kubernetes/Kyma

## <a id="api-disclaimer"></a>Disclaimer on API Usage
This documentation provides information that might be useful in using Authorization Management Service. We will try to ensure that future versions of the APIs are backwards compatible to the immediately preceding version. 
This is also true for the API that is exposed with `com.sap.cloud.security.ams.dcl.client` package and its subpackages.  
Please check the [release notes](../releases.md) or the release notes to stay tuned about new features and API modifications.

## Requirements
- Java 17
- Spring Boot 3.4.0-SNAPSHOT or later
- Spring Security 6.0.0 or later

## Setup

### Maven Dependencies
When using [Spring Security OAuth 2.0 Resource Server](https://docs.spring.io/spring-security/reference/servlet/oauth2/resource-server/index.html) 
for authentication in your Spring Boot application, you can leverage the provided Spring Boot Starter:

```xml
<dependency>
    <groupId>com.sap.cloud.security.ams.client</groupId>
    <artifactId>spring-boot-starter-ams-resourceserver</artifactId>
    <version>${sap.cloud.security.ams.client.version}</version>
</dependency>
```
> It's possible to disable all auto configurations using `com.sap.cloud.security.ams.auto=false`.

### Base DCL

Every business application has to describe its security model in [dcl language](../../concepts/DeployDCL.md).<br>
To describe the names and types of attributes, a file named `schema.dcl` and must be located in the root folder, for example in ``src/main/resources``.

```
SCHEMA {
    salesOrder: {
            type: number
    },
    CountryCode: String
}
```

Furthermore, the security model comprises rules and policies that can be used or 
enriched by the administrator to grant users permissions. They can be organized in packages.
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
```

## Usage
This section explains how to configure the health check and how to enforce authorization decisions by integrating the Authorization Management Service with Spring Security.

### Authorization decision point readiness evaluation

#### Initial application startup check

Initial startup check will check if the bundle containing the defined policies are accessible to the application at the application initialization step.
This functionality is provided by `PolicyDecisionPoint` and upon application startup will 
perform a periodical health checks until `OK` state is reached or 
until health check times out with the max default timeout of 30s, then it will throw an `PolicyEngineCommunicationException` error.
Startup check behaviour can be configured by the following 2 parameters:
- `com.sap.dcl.client.startupHealthCheckTimeout` - max wait time for startup check in seconds
- `com.sap.dcl.client.failOnStartupCheck` - boolean whether the application should start, if policies are not available i.e. `PolicyEngineCommunicationException` won't be thrown

These parameters can be configured as:
- system properties: `System.setProperty("com.sap.dcl.client.startupHealthCheckTimeout", "0")`
- spring properties: `com.sap.dcl.client.startupHealthCheckTimeout:0` :bulb: spring properties will only work with autoconfiguration enabled
- maven command line argument: `-Dcom.sap.dcl.client.startupHealthCheckTimeout=0`
- environment variable

and they take the following precedence:
1. spring properties
2. system properties
3. maven command line argument
4. environment variable

#### Configure endpoint with health check
For web applications, Cloud Foundry recommends setting the health check type to `http` instead of a simple port check (default). 
That means you can configure an endpoint that returns http status code HTTP `200` if application is healthy. 
This check should include an availability check of the policy decision runtime:

````java
import com.sap.cloud.security.ams.dcl.client.pdp.HealthState;

@SpringBootApplication
@RestController
public class DemoApplication {

    @Autowired
    PolicyDecisionPoint policyDecisionPoint;

    @GetMapping(value = "/health")
    @ResponseStatus(OK)
    public String healthCheck() {
        if (policyDecisionPoint.getHealthStatus.getHealthState == HealthState.OK){
            return "{\"status\":\"UP\"}";
        }
        throw new HttpServerErrorException(SERVICE_UNAVAILABLE, "Policy engine is not reachable.");
    }
}
````
Find further documentation on [Cloud Foundry Documentation: Using App Health Checks](https://docs.cloudfoundry.org/devguide/deploy-apps/healthchecks.html).

#### Advanced health status check 
`PolicyDecisionPoint` exposes a `getHealthStatus()` method that provides fine granular policy
decision point health status information together with information of each policy bundle status. 

```java
import com.sap.cloud.security.ams.dcl.client.pdp.HealthState;

public class DemoApplication {
    @Autowired
    PolicyDecisionPoint policyDecisionPoint;

    public String healthCheck() {
        PolicyDecisionPointHealthStatus pdpStatus = policyDecisionPoint.getHealthStatus();
        if (pdpStatus.getHealthState() == HealthState.UPDATE_FAILED) { //Checks whether bundle update was successful
            Map<String, BundleStatus> bundles = pdpStatus.getBundleStatus();
            bundles.entrySet().forEach(b -> 
            {
                if (b.getValue().hasBundleError()) {
                    //Process the failing bundle update ...
                    LOG.error("{} bundle does not contain latest updates, last successful bundle activation was at {}", b.getKey, b.getValue().getLastSuccessfulActivation());
                }
            });
        }
    }
}
```

### Access control on request level

This section explains how to configure authentication and authorization checks for REST APIs 
exposed by your Servlet based Spring application. You are going to configure the ``SecurityFilterChain`` for ``HttpServletRequest``.

By default, Spring Security’s authorization will require all requests to be authenticated. 
We can configure Spring Security to have different rules by adding more rules in order of precedence.

In addition to the [common built-in Spring Security authorization rules](https://docs.spring.io/spring-security/reference/servlet/authorization/authorize-http-requests.html#authorize-requests) like ``hasAuthority(String authority)`` this library provides another one, namely: ``hasBaseAuthority(String action, String resource)``. 
With ``hasBaseAuthority`` security expression you can easily specify "start-stop-conditions". 
For instance, you can validate if the user's principal has the necessary authorization to 
execute a specific ``action`` on a certain ``resource``. It's also allowed to use a ``'*'`` to 
be less restrictive. In case the authorization check fails, an exception of type ``AccessDeniedException`` is thrown.

> :information_source: If a user only has limited access to a resource, for example, the user can only _'read'_ _'sales orders'_ within their own country, this restriction is not accounted for with ``hasBaseAuthority('read', 'salesOrders')``.

> :information_source: ``.access("hasBaseAuthority('read', '*')")`` is different from the ``.hasAuthority("read")`` policy. While both allow _'read'_ access, the ``.hasAuthority("read")`` policy may currently contain or could in the future have attribute-level constraints.

#### Example Security Configuration
````java
@Configuration
@PropertySource(factory = IdentityServicesPropertySourceFactory.class, ignoreResourceNotFound = true, value = { "" })
public class SecurityConfiguration {

    @Autowired
    SecurityExpressionHandler<FilterInvocation> amsWebExpressionHandler;

    @Autowired
  SecurityExpressionHandler<RequestAuthorizationContext> amsHttpExpressionHandler;

  @Bean
  public SecurityFilterChain filterChain(HttpSecurity http) throws Exception {
    WebExpressionAuthorizationManager hasBaseAuthority = new WebExpressionAuthorizationManager(
        "hasBaseAuthority('read', 'salesOrders')");
    hasBaseAuthority.setExpressionHandler(amsHttpExpressionHandler);
    http.authorizeHttpRequests(authz -> {
          authz.requestMatchers(GET, "/health", "/", "/uiurl")
              .permitAll();
          authz.requestMatchers(GET, "/salesOrders/**")
              .access(hasBaseAuthority);
          authz.requestMatchers(GET, "/authorized")
              .hasAuthority("view").anyRequest().authenticated();
        })
        .oauth2ResourceServer(oauth2 ->
            oauth2.jwt(jwt -> jwt.jwtAuthenticationConverter(amsAuthenticationConverter)));
    return http.build();
  }

}
````

See [Spring Authorization Documentation](https://docs.spring.io/spring-security/reference/servlet/authorization/index.html) for in-depth understanding and comprehensive details about its usage and implementation.

### Access control on method level

`PolicyDecisionPointSecurityExpression` extends the [common built-in Spring Security Expressions](https://docs.spring.io/spring-security/reference/servlet/authorization/authorize-http-requests.html#authorize-requests). 
This can be used to control access on method level using [Method Security](https://docs.spring.io/spring-security/reference/servlet/authorization/method-security.html).

> :heavy_exclamation_mark: In Spring Boot versions `>= 6.1`, the Java compiler needs to be configured with the [-parameters flag](https://github.com/spring-projects/spring-framework/wiki/Spring-Framework-6.1-Release-Notes#parameter-name-retention) to use the expressions below that refer to method parameters.

#### Overview of AMS Spring Security Expressions

| method                                                                            | Description                                                                                                                                                                                                                                                                                                                                                | Example                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                     |
|-----------------------------------------------------------------------------------|------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|-------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| `boolean hasAuthority(String action)`                                             | returns `true` if the authenticated user has the permission to perform the given ``action``. This check can only be applied for very trivial policies, e.g. `POLICY myPolicy { GRANT action ON *; }`.                                                                                                                                                      |                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                             |
| `boolean hasBaseAuthority(String action, String resource)`                        | returns `true` if the authenticated user has in principal the permission to perform a given ``action`` on a given ``resource``. It's also allowed to use a ``'*'`` to be less restrictive.                                                                                                                                                                 |                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                             | 
| `boolean forAction(String action, String... attributes)`                          | returns `true` if the user is authorized to perform a dedicated `action`. The `resource` in the DCL rules needs to be declared as `*` for this case. If the DCL rules depend on attributes that are not automatically filled (either by default or an `AttributesProvider`) then their values need to be provided as `attributes` arguments.<sup>`*`</sup> | Has user `read` access to any resources?<br/>`@PreAuthorize("forAction('read')")`                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                           | 
| `boolean forResource(String resource, String... attributes)`                      | returns `true` if the user is authorized to access the given `resource`. The `action` in the DCL rules needs to be declared as `*` for this case. If the DCL rules depend on attributes that are not automatically filled (either by default or an `AttributesProvider`) then their values need to be provided as `attributes` arguments.<sup>`*`</sup>    | Has user any access to `salesOrder` resource?<br/>`@PreAuthorize("forResource('salesOrders')")`                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                             |
| `boolean forResourceAction(String resource, String action, String... attributes)` | returns `true` if the user is authorized to perform a dedicated `action` on the given `resource`. If the DCL rules depend on attributes that are not automatically filled (either by default or an `AttributesProvider`) then their values need to be provided as `attributes` arguments. <sup>`*`</sup>                                                   | Has user `read` access to `salesOrders` resource?<br/>`@PreAuthorize("forResourceAction('salesOrders', 'read')")` <br/><br/>Has user `read` access to `salesOrders` resource with `CountryCode` = "DE" and `salesOrder.type` = 247?<br/><br/>Consider a RESTful application that looks up a salesOrder by country and type from the URL path in the format, e.g. ``/readByCountryAndType/{countryCode}/{type}``. <br/>You are able to refer to path variables within a URL:<br/><pre>@PreAuthorize("forResourceAction('salesOrders', 'read', 'CountryCode:string=' + #countryCode, 'salesOrder.type:number=' + #type)")<br/>@GetMapping(value = "/readByCountryAndType/{countryCode}/{type}")<br/>public String readSelectedSalesOrder(@PathVariable String countryCode, @PathVariable String type){}</pre> |

> :bulb: <sup>`*`</sup>
> the `attributes` [varargs](https://docs.oracle.com/javase/1.5.0/docs/guide/language/varargs.html) shall have the
> format: `<attribute name>:<string|number| boolean>=<attribute value>`. Currently, these data types are
> supported: `string`, `number` and `boolean`.

### Authorization hybrid mode - Xsuaa and AMS authorization checks

The hybrid authorization mode provides a way to use Xsuaa tokens and their scopes, as well as IAS tokens with AMS policies.

There are 2 ways to enforce authorization in the hybrid mode:

1. **Xsuaa scope check**: this option performs as before for Xsuaa issued tokens and authorizations are defined with a
   help
   of `hasAuthority("Role")`
   method provided
   by [Spring Security](https://docs.spring.io/spring-security/reference/servlet/authorization/authorize-http-requests.html#authorizing-endpoints).
   This option is enabled by default within
   an *autoconfiguration* `ResourceServerAuthConverterAutoConfiguration.java` ,
   if the service configuration has an Xsuaa binding.
2. **AMS Policy Authorization**: this option disables the Xsuaa scope check, requiring the Spring
   property `sap.security.hybrid.xsuaa.scope-check.enabled `to be set to false. Authorization is carried out
   through `PolicyDecisionPoint` and enforced
   by [AMS Spring Security Expressions](#overview-of-ams-spring-security-expressions). However, for this mode to work,
   Xsuaa Scopes need to be converted into AMS policies. This can be done with the help of [AttributesProccessor](../jakarta-ams/jakarta-ams.md#customize-attributes), but
   please note this conversion process is not provided by this library.

### Additional Information

#### How to use ``PolicyDecisionPoint`` API

````java
@Autowired
private PolicyDecisionPoint policyDecisionPoint;

    ...
    Principal principal=Principal.create();
    Attributes attributes=principal.getAttributes().setAction("read");
    boolean isReadAllowed=policyDecisionPoint.allow(attributes);
````
See [here](../jakarta-ams/jakarta-ams#overview-policydecisionpoint-methods) for ``PolicyDecisionPoint`` API documentation.

> :heavy_exclamation_mark: [Consider also limitation of api liability](#api-disclaimer).

#### Pre-fill user-attributes from OIDC token
You can simply create a ``Principal`` instance within the same thread. Use ``Principal.create()`` to derive the principal information from the OIDC token, that is stored in ``SecurityContextHolder``.

> Alternatively, you can also build the ``Principal`` using the ``PrincipalBuilder``.   
> E.g. `PrincipalBuilder.create("the-zone-id", "the-user-id").build();` 


## Testing
You can find the test utilities `spring-ams-test-starter` module.

## Audit Logging
Please check out this [documentation](../jakarta-ams/jakarta-ams.md#audit-logging).

## Troubleshooting

For troubleshooting purposes check [here](../../Troubleshooting.md).

### Set DEBUG log level

First, configure the Debug log level for Spring Framework Web and all Security related libs. 
This can be done as part of your `application.yml` or `application.properties` file.

```yaml
logging.level:
  com.sap.cloud.security: DEBUG       # set SAP-class loggers to DEBUG; set to ERROR for production setup
  org.springframework: ERROR          # set to DEBUG to see all beans loaded and auto-config conditions met
  org.springframework.security: DEBUG # set to ERROR for production setup
  org.springframework.web: DEBUG      # set to ERROR for production setup
```

Then, in case you like to see what different filters are applied to particular request then set debug flag to true in `@EnableWebSecurity` annotation:
```java
@Configuration
@EnableWebSecurity(debug = true) // TODO "debug" may include sensitive information. Do not use in a production system!
public class SecurityConfiguration {
   ...
}
```

:bulb: Remember to restage your application for the changes to take effect.

## Sample
- [spring-security sample using Spring Security](https://github.com/SAP-samples/ams-samples-java/tree/main/spring-security-ams)

