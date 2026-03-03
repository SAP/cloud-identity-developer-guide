# Migration Guide: v3 → v4

This guide helps you migrate your Java application from AMS Client Library version 3.x to version 4.x. Version 4
introduces significant changes to the module structure, dependency coordinates, and API design.

::: warning Upgrade Recommended
Version 3.x is deprecated. All new projects should use version 4.x, and existing projects should migrate as soon as
possible.
:::

## Overview of Changes

Version 4 brings the following major changes:

| Area            | v3                                  | v4                                        |
|-----------------|-------------------------------------|-------------------------------------------|
| Maven Group ID  | `com.sap.cloud.security.ams.client` | `com.sap.cloud.security.ams`              |
| BOM Artifact    | `parent`                            | `ams-bom`                                 |
| DCL Compilation | Maven plugins                       | `@sap/ams-dev` npm package                |
| DCL Output      | `target/dcl_opa`                    | `target/generated-test-resources/ams/dcn` |
| Spring Security | Expression handlers                 | `AmsRouteSecurity`                        |
| CAP Integration | `cap-ams-support`                   | `spring-boot-starter-cap-ams`             |

## Maven Dependencies

### BOM Configuration

**v3 Configuration:**

```xml

<dependencyManagement>
    <dependencies>
        <dependency>
            <groupId>com.sap.cloud.security.ams.client</groupId>
            <artifactId>parent</artifactId>
            <version>${sap.cloud.security.ams.version}</version>
            <type>pom</type>
            <scope>import</scope>
        </dependency>
    </dependencies>
</dependencyManagement>
```

**v4 Configuration:**

```xml

<dependencyManagement>
    <dependencies>
        <dependency>
            <groupId>com.sap.cloud.security.ams</groupId>
            <artifactId>ams-bom</artifactId>
            <version>${sap.cloud.security.ams.version}</version>
            <type>pom</type>
            <scope>import</scope>
        </dependency>
    </dependencies>
</dependencyManagement>
```

### Artifact Mapping

The following table shows the mapping from v3 artifacts to v4 artifacts:

| v3 Artifact                              | v4 Artifact                        | Description                           |
|------------------------------------------|------------------------------------|---------------------------------------|
| `jakarta-ams`                            | `ams-core`                         | Core authorization library            |
| `java-ams-test`                          | `ams-test`                         | Testing utilities for plain Java      |
| `spring-ams`                             | `spring-boot-starter-ams`          | Spring Boot integration               |
| `spring-boot-starter-ams-resourceserver` | `spring-boot-starter-ams`          | Spring Security resource server       |
| `spring-boot-starter-ams-test`           | `spring-boot-starter-ams-test`     | Spring Boot test utilities            |
| `cap-ams-support`                        | `spring-boot-starter-cap-ams`      | CAP Java integration                  |
| N/A                                      | `spring-boot-starter-cap-ams-test` | CAP Java test utilities               |
| N/A                                      | `spring-boot-3-starter-ams-health` | Spring Boot Actuator health indicator |

### Example: Spring Boot Application

**v3:**

```xml

<dependency>
    <groupId>com.sap.cloud.security.ams.client</groupId>
    <artifactId>spring-boot-starter-ams-resourceserver</artifactId>
    <version>${sap.cloud.security.ams.version}</version>
</dependency>
<dependency>
<groupId>com.sap.cloud.security.ams.client</groupId>
<artifactId>spring-boot-starter-ams-test</artifactId>
<version>${sap.cloud.security.ams.version}</version>
<scope>test</scope>
</dependency>
```

**v4:**

```xml

<dependency>
    <groupId>com.sap.cloud.security.ams</groupId>
    <artifactId>spring-boot-starter-ams</artifactId>
</dependency>
<dependency>
<groupId>com.sap.cloud.security.ams</groupId>
<artifactId>spring-boot-starter-ams-test</artifactId>
<scope>test</scope>
</dependency>
        <!-- Optional: Health indicator -->
<dependency>
<groupId>com.sap.cloud.security.ams</groupId>
<artifactId>spring-boot-3-starter-ams-health</artifactId>
</dependency>
```

### Example: CAP Java Application

**v3:**

```xml

<dependency>
    <groupId>com.sap.cloud.security.ams.client</groupId>
    <artifactId>jakarta-ams</artifactId>
    <version>${sap.cloud.security.ams.version}</version>
</dependency>
<dependency>
<groupId>com.sap.cloud.security.ams.client</groupId>
<artifactId>cap-ams-support</artifactId>
<version>${sap.cloud.security.ams.version}</version>
</dependency>
```

**v4:**

```xml

<dependency>
    <groupId>com.sap.cloud.security.ams</groupId>
    <artifactId>spring-boot-starter-cap-ams</artifactId>
</dependency>
<dependency>
<groupId>com.sap.cloud.security.ams</groupId>
<artifactId>spring-boot-starter-cap-ams-test</artifactId>
<scope>test</scope>
</dependency>
```

## DCL Compilation

### Migration from Maven Plugins to npm

Version 3 used Maven plugins (`dcl-compiler-plugin` and `dcl-surefire-plugin`) for DCL compilation and testing. Version
4 uses the `@sap/ams-dev` npm package instead.

**v3 Maven Plugin Configuration:**

```xml

<plugin>
    <groupId>com.sap.cloud.security.ams.client</groupId>
    <artifactId>dcl-compiler-plugin</artifactId>
    <version>${sap.cloud.security.ams.version}</version>
    <executions>
        <execution>
            <goals>
                <goal>compile</goal>
            </goals>
            <configuration>
                <verbose>true</verbose>
                <failOn>deprecation</failOn>
                <sourceDirectory>${project.basedir}/dcldeployer/dcl</sourceDirectory>
                <dcn>true</dcn>
                <dcnParameter>pretty</dcnParameter>
                <compileTestToDcn>true</compileTestToDcn>
            </configuration>
        </execution>
    </executions>
</plugin>
<plugin>
<groupId>com.sap.cloud.security.ams.client</groupId>
<artifactId>dcl-surefire-plugin</artifactId>
<version>${sap.cloud.security.ams.version}</version>
<executions>
    <execution>
        <goals>
            <goal>test</goal>
        </goals>
        <configuration>
            <verbose>true</verbose>
            <prefix>${project.groupId}__${project.artifactId}__</prefix>
        </configuration>
    </execution>
</executions>
</plugin>
```

**v4 npm-based Configuration:**

For Spring Boot applications, use `frontend-maven-plugin`:

```xml

<plugin>
    <groupId>com.github.eirslett</groupId>
    <artifactId>frontend-maven-plugin</artifactId>
    <version>1.14.1</version>
    <executions>
        <execution>
            <id>install node and npm</id>
            <goals>
                <goal>install-node-and-npm</goal>
            </goals>
            <phase>generate-test-resources</phase>
            <configuration>
                <nodeVersion>v24.11.0</nodeVersion>
            </configuration>
        </execution>
        <execution>
            <id>compile-dcl</id>
            <goals>
                <goal>npx</goal>
            </goals>
            <phase>generate-test-resources</phase>
            <configuration>
                <arguments>--package=@sap/ams-dev compile-dcl
                    -d ${project.basedir}/src/main/resources/dcl
                    -o ${project.build.directory}/generated-test-resources/ams/dcn
                </arguments>
            </configuration>
        </execution>
    </executions>
</plugin>
```

For CAP Java applications, use `cds-maven-plugin`:

```xml

<plugin>
    <groupId>com.sap.cds</groupId>
    <artifactId>cds-maven-plugin</artifactId>
    <executions>
        <!-- ... other executions ... -->
        <execution>
            <id>compile-dcl</id>
            <goals>
                <goal>npx</goal>
            </goals>
            <phase>generate-test-resources</phase>
            <configuration>
                <arguments>--package=@sap/ams-dev compile-dcl
                    -d ${project.basedir}/src/main/resources/ams
                    -o ${project.build.directory}/generated-test-resources/ams/dcn
                </arguments>
            </configuration>
        </execution>
    </executions>
</plugin>
```

### DCL Source Directory

The default location for DCL files has changed:

| Application Type       | v3 Location                   | v4 Location                   |
|------------------------|-------------------------------|-------------------------------|
| Plain Java/Spring Boot | `dcldeployer/dcl/`            | `src/main/resources/dcl/`     |
| CAP Java               | `srv/src/main/resources/ams/` | `srv/src/main/resources/ams/` |

### DCL Output Directory

| v3 Output         | v4 Output                                  |
|-------------------|--------------------------------------------|
| `target/dcl_opa/` | `target/generated-test-resources/ams/dcn/` |

## DCL Syntax Updates

### Schema Definition

The schema keyword and type names are now in uppercase:

**v3:**

```dcl
schema {
    salesOrder: {
        type: number
    },
    CountryCode: string
}
```

**v4:**

```dcl
SCHEMA {
    SalesOrder: {
        Type: Number
    }
    CountryCode: String
}
```

::: tip Compatibility
While lowercase syntax may still work for backward compatibility, it is recommended to update to the uppercase format
for consistency with the current specification.
:::

### Policy Syntax

Policy syntax remains largely unchanged. Both `POLICY` and `GRANT`/`ASSIGN` keywords were already uppercase in v3.

## Spring Security Configuration

### Route-Level Security

The way you configure route-level security with AMS has changed significantly.

**v3 Configuration:**

```java
import org.springframework.security.access.expression.SecurityExpressionHandler;
import org.springframework.security.web.access.expression.WebExpressionAuthorizationManager;
import org.springframework.security.web.access.intercept.RequestAuthorizationContext;

@Configuration
public class SecurityConfiguration {

    @Autowired
    Converter<Jwt, AbstractAuthenticationToken> amsAuthenticationConverter;

    @Bean
    public SecurityFilterChain filterChain(
            HttpSecurity http,
            SecurityExpressionHandler<RequestAuthorizationContext> amsHttpExpressionHandler)
            throws Exception {

        WebExpressionAuthorizationManager hasBaseAuthority =
                new WebExpressionAuthorizationManager("hasBaseAuthority('read', 'salesOrders')");
        hasBaseAuthority.setExpressionHandler(amsHttpExpressionHandler);

        http.authorizeHttpRequests(authz -> {
                    authz.requestMatchers(GET, "/health").permitAll();
                    authz.requestMatchers(GET, "/salesOrders/**").access(hasBaseAuthority);
                    authz.anyRequest().authenticated();
                })
                .oauth2ResourceServer(oauth2 -> oauth2.jwt(
                        jwt -> jwt.jwtAuthenticationConverter(amsAuthenticationConverter)));

        return http.build();
    }
}
```

**v4 Configuration:**

```java
import com.sap.cloud.security.ams.spring.AmsRouteSecurity;

import static com.sap.cloud.security.ams.api.Privilege.of;

@Configuration
@EnableWebSecurity
@EnableMethodSecurity
public class SecurityConfiguration {

    // Define privileges as constants
    private static final Privilege READ_SALES_ORDERS = of("read", "salesOrders");

    @Bean
    public SecurityFilterChain filterChain(HttpSecurity http, AmsRouteSecurity via)
            throws Exception {

        http.authorizeHttpRequests(authz -> {
                    authz.requestMatchers(GET, "/health").permitAll();
                    authz.requestMatchers(GET, "/salesOrders/**")
                            .access(via.checkPrivilege(READ_SALES_ORDERS));
                    authz.anyRequest().authenticated();
                })
                .oauth2ResourceServer(oauth2 -> oauth2.jwt(Customizer.withDefaults()));

        return http.build();
    }
}
```

### Key Changes

| Aspect              | v3                                            | v4                                               |
|---------------------|-----------------------------------------------|--------------------------------------------------|
| Expression Handler  | `amsHttpExpressionHandler`                    | `AmsRouteSecurity`                               |
| Authorization Check | `hasBaseAuthority('action', 'resource')`      | `via.checkPrivilege(Privilege)`                  |
| JWT Converter       | Manual `amsAuthenticationConverter` injection | Auto-configured with `Customizer.withDefaults()` |
| Pre-check Support   | N/A                                           | `via.precheckPrivilege(Privilege)`               |

### Privilege Constants

It's recommended to define privilege constants in a dedicated class:

```java
import com.sap.cloud.security.ams.api.Privilege;

public final class Privileges {
    public static final Privilege READ_PRODUCTS = Privilege.of("read", "products");
    public static final Privilege CREATE_ORDERS = Privilege.of("create", "orders");
    public static final Privilege DELETE_ORDERS = Privilege.of("delete", "orders");
    public static final Privilege READ_ORDERS = Privilege.of("read", "orders");

    private Privileges() {
    }
}
```

## Core API Changes (Plain Java / Javalin)

### Authorization Client Initialization

**v3:**

```java
import com.sap.cloud.security.ams.dcl.client.pdp.PolicyDecisionPoint;

import static com.sap.cloud.security.ams.factory.AmsPolicyDecisionPointFactory.DEFAULT;

PolicyDecisionPoint policyDecisionPoint = PolicyDecisionPoint.create(DEFAULT);
```

**v4:**

```java
import com.sap.cloud.security.ams.api.AuthorizationManagementService;
import com.sap.cloud.environment.servicebinding.api.ServiceBinding;
import com.sap.cloud.environment.servicebinding.api.DefaultServiceBindingAccessor;

ServiceBinding identityBinding = DefaultServiceBindingAccessor.getInstance()
        .getServiceBindings().stream()
        .filter(binding -> "identity".equals(binding.getServiceName().orElse(null)))
        .findFirst()
        .orElseThrow(() -> new IllegalStateException("No identity binding found"));

AuthorizationManagementService ams =
        AuthorizationManagementService.fromIdentityServiceBinding(identityBinding);
```

### Performing Authorization Checks

**v3:**

```java
import com.sap.cloud.security.ams.api.Principal;
import com.sap.cloud.security.ams.dcl.client.pdp.Attributes;

Attributes attributes = Principal.create()
        .getAttributes()
        .setAction("read")
        .setResource("salesOrders");

if(!policyDecisionPoint.

allow(attributes)){
        // Access denied
        }
```

**v4:**

```java
import com.sap.cloud.security.ams.api.Authorizations;
import com.sap.cloud.security.ams.api.Privilege;
import com.sap.cloud.security.ams.api.Principal;
import com.sap.cloud.security.ams.core.SciAuthorizationsProvider;

// Create an authorizations provider
SciAuthorizationsProvider<Authorizations> authProvider =
        SciAuthorizationsProvider.create(ams, Authorizations::of);

        // Get authorizations for current principal
        Authorizations authorizations = authProvider.getAuthorizations(
                Principal.fromSecurityContext());

        // Check privilege
        Privilege readSalesOrders = Privilege.of("read", "salesOrders");
if(authorizations.

        checkPrivilege(readSalesOrders).

        isDenied()){
        // Access denied
        }
```

### Custom Authorizations Class

In v4, you can create custom `Authorizations` implementations:

```java
import com.sap.cloud.security.ams.api.Authorizations;
import com.sap.cloud.security.ams.api.AuthorizationsData;

public class ShoppingAuthorizations implements Authorizations {
    private final AuthorizationsData data;

    private ShoppingAuthorizations(AuthorizationsData data) {
        this.data = data;
    }

    public static ShoppingAuthorizations of(AuthorizationsData data) {
        return new ShoppingAuthorizations(data);
    }

    @Override
    public AuthorizationsData getData() {
        return data;
    }

    // Add domain-specific methods
    public boolean canReadProducts() {
        return !checkPrivilege(Privilege.of("read", "products")).isDenied();
    }
}
```

## CAP Java Configuration

### Local Development

**v3** required explicit configuration of test sources in `application.yaml`:

```yaml
cds:
  security:
    authorization:
      ams:
        test-sources: "" # empty uses default srv/target/dcl_opa
```

**v4** auto-detects local test mode when mock users are configured with policies. No explicit `test-sources`
configuration is required:

```yaml
cds:
  security:
    mock:
      users:
        admin:
          password: admin
          roles:
            - admin
        stock-manager:
          policies:
            - cap.StockManager
        stock-manager-fiction:
          policies:
            - local.StockManagerFiction
```

::: info Automatic Detection
In v4, when mock users with `policies` are maintained for a profile, the policy assignment via mock users is active by
default.
:::

## Testing

### Spring Security Tests

**v3:**

```java
import com.sap.cloud.security.test.extension.SecurityTestExtension;

import static com.sap.cloud.security.ams.spring.test.resourceserver.MockOidcTokenRequestPostProcessor.userWithPolicies;
import static com.sap.cloud.security.ams.spring.test.resourceserver.MockOidcTokenRequestPostProcessor.userWithoutPolicies;
import static com.sap.cloud.security.config.Service.IAS;

@SpringBootTest
@AutoConfigureMockMvc
@ActiveProfiles("test")
class MyControllerTest {

    @RegisterExtension
    static SecurityTestExtension extension =
            SecurityTestExtension.forService(IAS).setPort(MOCK_SERVER_PORT);

    @Autowired
    private MockMvc mockMvc;

    @Test
    void testReadWithPermission() throws Exception {
        mockMvc.perform(get("/read")
                        .with(userWithPolicies(extension.getContext(), "common.readAll")))
                .andExpect(status().isOk());
    }

    @Test
    void testReadWithoutPermission() throws Exception {
        mockMvc.perform(get("/read")
                        .with(userWithoutPolicies(extension.getContext())))
                .andExpect(status().isForbidden());
    }
}
```

**v4:**

```java
import org.springframework.boot.test.context.SpringBootTest;
import org.springframework.test.context.ActiveProfiles;

@SpringBootTest
@ActiveProfiles("test")
@Import(TestSecurityConfiguration.class)
class MyControllerTest {

    @Autowired
    private MockMvc mockMvc;

    @Test
    void testReadWithPermission() throws Exception {
        String userJwt = loadJwtFromFile("User_with_read_policy.json");
        mockMvc.perform(get("/read")
                        .header("Authorization", "Bearer " + userJwt))
                .andExpect(status().isOk());
    }

    @Test
    void testReadWithoutPermission() throws Exception {
        String userJwt = loadJwtFromFile("User_without_policies.json");
        mockMvc.perform(get("/read")
                        .header("Authorization", "Bearer " + userJwt))
                .andExpect(status().isForbidden());
    }
}
```

**Test Security Configuration (v4):**

```java

@TestConfiguration
public class TestSecurityConfiguration {

    private final Base64JwtDecoder base64JwtDecoder = Base64JwtDecoder.getInstance();

    @Bean
    @Primary
    public JwtDecoder jwtDecoder() {
        return token -> {
            DecodedJwt decodedJwt = base64JwtDecoder.decode(token);
            SapIdToken sapIdToken = new SapIdToken(decodedJwt);
            SecurityContext.setToken(sapIdToken);

            Map<String, Object> headers = sapIdToken.getHeaders();
            Map<String, Object> claims = sapIdToken.getClaims();
            Instant issuedAt = claims.containsKey("iat")
                    ? Instant.ofEpochSecond(((Number) claims.get("iat")).longValue())
                    : Instant.now();
            Instant expiresAt = Optional.ofNullable(sapIdToken.getExpiration())
                    .orElse(Instant.now().plusSeconds(3600));

            return new Jwt(token, issuedAt, expiresAt, headers, claims);
        };
    }
}
```

### JWT Test Files (v4)

Create JSON files in `src/test/resources/jwt/` containing user claims and policies:

```json
{
  "sub": "alice",
  "iss": "https://test.accounts.ondemand.com",
  "aud": "test-client",
  "iat": 1704067200,
  "exp": 1893456000,
  "ams_policies": [
    "shopping.CreateOrders",
    "shopping.DeleteOrders"
  ]
}
```

## Package Migrations

Update your import statements according to this mapping:

| v3 Package                                              | v4 Package                          |
|---------------------------------------------------------|-------------------------------------|
| `com.sap.cloud.security.ams.dcl.client.pdp`             | `com.sap.cloud.security.ams.api`    |
| `com.sap.cloud.security.ams.factory`                    | `com.sap.cloud.security.ams.api`    |
| `com.sap.cloud.security.ams.spring.test.resourceserver` | `com.sap.cloud.security.ams.spring` |

## Health Indicator (New in v4)

v4 introduces a Spring Boot Actuator health indicator for AMS:

```xml

<dependency>
    <groupId>com.sap.cloud.security.ams</groupId>
    <artifactId>spring-boot-3-starter-ams-health</artifactId>
</dependency>
```

Configure in `application.yaml`:

```yaml
management:
  endpoint:
    health:
      show-components: always
      show-details: always
  endpoints:
    web:
      exposure:
        include: health
  health:
    defaults.enabled: false
    ping.enabled: true
    db.enabled: true
    ams.enabled: true
```

## Migration Checklist

Use this checklist to track your migration progress:

- [ ] Update Maven BOM from `parent` to `ams-bom`
- [ ] Update Maven Group ID from `com.sap.cloud.security.ams.client` to `com.sap.cloud.security.ams`
- [ ] Update artifact names according to the mapping table
- [ ] Remove `dcl-compiler-plugin` and `dcl-surefire-plugin`
- [ ] Add npm-based DCL compilation (`@sap/ams-dev`)
- [ ] Move DCL files to new location (if applicable)
- [ ] Update DCL output directory in configuration
- [ ] Update schema syntax to uppercase (optional but recommended)
- [ ] Update Spring Security configuration to use `AmsRouteSecurity`
- [ ] Update authorization checks to use `Privilege` class
- [ ] Update test configurations
- [ ] Remove explicit `test-sources` configuration (CAP)
- [ ] Add health indicator dependency (optional)
- [ ] Update import statements

## Troubleshooting

### Common Issues

**1. ClassNotFoundException for v3 Classes**

If you see errors like `ClassNotFoundException: com.sap.cloud.security.ams.dcl.client.pdp.PolicyDecisionPoint`, ensure
you've updated all dependencies to v4 artifacts.

**2. DCL Compilation Errors**

If DCL compilation fails with the npm-based approach:

- Ensure Node.js is installed (v18+ recommended)
- Check that the DCL source directory path is correct
- Verify the output directory exists or can be created

**3. Spring Security Authorization Failures**

If authorization checks fail after migration:

- Verify `AmsRouteSecurity` is injected correctly
- Check that privilege actions and resources match your DCL policies
- Ensure JWT decoder is configured with `Customizer.withDefaults()`

**4. CAP Mock Users Not Working**

If mock users with policies aren't being recognized:

- Ensure you're using the correct profile
- Verify policy names include the correct package prefix (e.g., `cap.StockManager`)
- Check that DCL files are compiled and available in the target directory

## Further Resources

- [AMS Core Module](/Libraries/java/ams-core)
- [Spring Boot AMS Integration](/Libraries/java/spring-boot-ams)
- [CAP AMS Integration](/Libraries/java/cap-ams)
- [Authorization Checks](/Authorization/AuthorizationChecks)
- [Java Sample Applications](https://github.com/SAP-samples/ams-samples-java)