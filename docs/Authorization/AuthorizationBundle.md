# Authorization Bundle

The *Authorization Bundle* is a container for [authorization policies](#authorization-policies) that sets the context for authorization checks in an application. Bundles are compiled centrally by the Authorization Management Service (**AMS**) from where client applications download the bundle of their respective [AMS service instance](/Authorization/GettingStarted#provisioning-of-ams-instances) to perform authorization checks.

The bundle contains both the authorization policies defined as part of the application's source code as well as the policies created by administrators at runtime. For this reason, clients regularly poll for changes to keep the local copy up to date when administrators make changes (see [Client Library Initialization](#client-library-initialization)).

## Authorization Policies

Authorization policies grant the right to perform actions on protected resources in an application. They can be assigned to users to control access to various parts of the application.

Developers can define a set of base policies that can be assigned directly or used as building blocks by the application administrators to create additional so-called admin policies at runtime.

### DCL

Authorization policies are defined in a domain-specific language called Data Control Language (**DCL**) that supports conditions that can be used to grant fine-grained access to resources.

#### Example
Here is an example of authorization policies defined in DCL:

```dcl
SCHEMA {
    category: String;
}

POLICY ReadProducts {
    GRANT read ON products WHERE category IS NOT RESTRICTED;
}

POLICY ReadOfficeSupplies {
    USE ReadProducts RESTRICT category = 'OfficeSupplies';
}
```

#### Specification
The complete specification for DCL can be found in the [SAP Help Portal](https://help.sap.com/docs/cloud-identity-services/cloud-identity-services/data-control-language-dcl).

#### Deployment
Please refer to the [Deploying DCL](/Authorization/DeployDCL) page for instructions on how to deploy DCL policies to an AMS service instance.

## Client Library Initialization

To initialize the AMS client libraries, an instance of the `AuthorizationManagementService` class must be created. In production, applications create an instance from **certificate-based** credentials for mTLS authentication with the AMS service to download the authorization bundle. These credentials are typically provided in the form of a SAP BTP service binding for the SAP Cloud Identity Services (**SCI**) offering.

::: code-group

```js [Node.js]
const { AuthorizationManagementService } = require("@sap/ams");

// pass your @sap/xssec IdentityService instance used for authentication
// which was created from certificate-based SCI credentials or
// alternatively, a fixed { credentials ... } object directly
const ams = AuthorizationManagementService
    .fromIdentityService(identityService);
```

```java [Java]
import com.sap.cloud.security.ams.api.AuthorizationManagementService;
// from com.sap.cloud.environment.servicebinding:java-sap-vcap-services
import com.sap.cloud.environment.servicebinding.api.DefaultServiceBindingAccessor;


ServiceBinding identityServiceBinding = DefaultServiceBindingAccessor
    .getInstance().getServiceBindings().stream()
    .filter(binding -> "identity".equals(binding.getServiceName().orElse(null)))
    .findFirst()
    .orElse(null);
                
 AuthorizationManagementService ams = AuthorizationManagementService
    .fromIdentityServiceBinding(identityServiceBinding);
```

:::

::: danger Important
After creating the `AuthorizationManagementService` instance, the application must ensure with a [startup check](#startup-check) that the instance is ready before accepting traffic that requires authorization checks.
:::

::: tip
The AMS client libraries integrate into different web frameworks, such as [CAP](https://cap.cloud.sap/docs/) or [Spring Security](https://spring.io/projects/spring-security). The respective [Spring Boot starters](/Authorization/GettingStarted#java) and [Node.js CAP plugin](/Authorization/GettingStarted#node-js) automatically create the `AuthorizationManagementService` instance from the SCI service binding in the application's environment, so manual initialization is not required in these cases.
:::

### Certificate Configuration

For SAP BTP service bindings with `"credential-type": "X509_PROVIDED"` or `"credential-type": "X509_ATTESTED"`, the certificate and key required for mTLS authentication with AMS is not included in the service binding and must be provided by the application before the library instantiation.

::: tip X509_GENERATED
SAP BTP service bindings with `"credential-type": "X509_GENERATED"` already contain the client certificate and key. No certificate configuration is needed in this case.
:::

::: code-group

```js [Node.js]
// Update the identityService object passed to
// fromIdentityService with the certificate information.
// cert and key must be PEM-encoded strings
identityService.setCertificateAndKey(cert, key);

// then create the AMS instance as usual
const ams = AuthorizationManagementService
    .fromIdentityService(identityService);
```

```js [Node.js (CAP)]
// server.js
const cds = require('@sap/cds');
const { amsCapPluginRuntime } = require("@sap/ams");

// Extend the AMS CAP plugin runtime credentials with the provided certificate/key.
cds.on('served', async () => {
    // assuming cert and key are PEM-encoded strings
    amsCapPluginRuntime.credentials.certificate = cert;
    amsCapPluginRuntime.credentials.key = key; 
});

// If certificate/key change during runtime, update the properties inside the credentials object again.
someCredentialsWatcher.on('change', (newCert, newKey) => {
    // assuming newCert and newKey are PEM-encoded strings
    amsCapPluginRuntime.credentials.certificate = newCert;
    amsCapPluginRuntime.credentials.key = newKey;
});
```

```java [Java]
import com.sap.cloud.security.ams.api.AuthorizationManagementService;
import com.sap.cloud.security.ams.config.CloudAuthorizationManagementServiceConfig;
import java.security.KeyStore;

// The KeyStore must contain exactly one private key entry with no password (empty password).
KeyStore keyStore = // load KeyStore containing client certificate and private key

CloudAuthorizationManagementServiceConfig config = new CloudAuthorizationManagementServiceConfig()
    .withKeyStore(keyStore);

AuthorizationManagementService ams = AuthorizationManagementService
    .fromIdentityServiceBinding(identityServiceBinding, config);
```

```java [Spring Boot]
import org.springframework.context.annotation.Bean;
import org.springframework.beans.factory.annotation.Qualifier;
import java.security.KeyStore;

@Bean
@Qualifier("amsKeyStore")
public KeyStore amsKeyStore() {
    // The KeyStore must contain exactly one private key entry with no password (empty password).
    KeyStore keyStore = // load KeyStore containing client certificate and private key
    return keyStore;
}
```

:::

##  Startup Check

The application must ensure that the `AuthorizationManagementService` instance has loaded its initial authorization bundle before it performs authorization checks. Otherwise, exceptions are thrown because no authorization data is available yet.

There are two general approaches to achieve this:

- **Synchronous startup check**: block application startup until the instance is ready, failing after a timeout. This is the most robust option because it prevents *any* premature usage of the instance during startup, including *internal* usage by frameworks such as CAP — not just *external* traffic.
- **Readiness endpoint integration**: expose the instance's readiness state through a health or readiness endpoint so that the platform does not route traffic to the process until the instance is ready. This only delays authorization checks triggered by *external* traffic.

::: tip Spring Boot starters do this for you
The [Spring Boot starters](/Authorization/GettingStarted#java) perform a synchronous startup check out of the box, so no manual startup check is required there. Its timeout can be adjusted and it can be disabled via configuration (see the `Spring Boot` examples below).

Other setups — including the [Node.js CAP plugin](/Authorization/GettingStarted#node-js) — must implement the startup check themselves, as shown in the examples below.
:::

The following examples show, per framework, how the startup check is performed or configured, and how to expose the readiness state in a custom health endpoint. For Spring Boot the check runs automatically and the tabs only show how to adjust the timeout or opt out. Disabling the check is discouraged: the application then becomes responsible for ensuring readiness before performing authorization checks (e.g. via the programmatic APIs shown here).

::: code-group
```js [Node.js (CAP)]
// server.js
const cds = require('@sap/cds');
const { amsCapPluginRuntime } = require("@sap/ams");

cds.on('served', async () => {
    // CAP awaits all 'served' handlers to finish before it serves traffic.
    // So even though whenReady is awaited asynchronously here, it effectively
    // blocks application startup until AMS is ready — or aborts it if the handler
    // throws after the 30s timeout.
    await amsCapPluginRuntime.ams.whenReady(30);
    console.log("AMS has become ready.");
});
```

```js [Node.js]
// demonstrates readiness status via health endpoint

let isReady = false;
const healthCheck = (req, res) => {
    if (isReady) {
        res.json({ status: 'UP' });
    } else {
        res.status(503).json({ status: 'DOWN', message: 'Service is not ready' });
    }
};

const amsStartupCheck = async () => {
    try {
        await ams.whenReady(AMS_STARTUP_TIMEOUT);
        isReady = true;
        console.log("AMS is ready now.");
    } catch (e) {
        console.error("AMS didn't get ready in time:", e);
        process.exit(1);
    }
};

app.get('/health', healthCheck);
const server = app.listen(PORT, () => {
    console.log(`Server is listening on port ${PORT}`);
});

amsStartupCheck();
```

```yaml [Spring Boot (CAP)]
# The synchronous startup check runs automatically.
# Adjust its timeout or opt out via configuration:
cds:
  security:
    authorization:
      ams:
        startup-check:
          enabled: true   # perform synchronous startup check (default: true)
          timeout: 30s    # fail startup if not ready within this duration (default: 30s)
```

```yaml [Spring Boot]
# The synchronous startup check runs automatically.
# Adjust its timeout or opt out via configuration:
sap:
  ams:
    startup-check:
      enabled: true   # perform synchronous startup check (default: true)
      timeout: 30s    # fail startup if not ready within this duration (default: 30s)
```

```java [Java]
import java.time.Duration;
import java.util.concurrent.TimeUnit;
import java.util.concurrent.atomic.AtomicBoolean;

// Synchronous startup check:
// throws an error if the AMS module doesn't get ready within 30 seconds
ams.awaitReady(Duration.ofSeconds(30));

// Asynchronous startup check: demonstrates readiness status in health endpoint
private static final AtomicBoolean isReady = new AtomicBoolean(false);

app.get("/health", ctx -> {
    if (isReady.get()) {
        ctx.json(HealthStatus.up());
    } else {
        ctx.status(503).json(HealthStatus.down("Service is not ready"));
    }
});

// Wait up to 30s for AMS to become ready
ams.whenReady().orTimeout(30, TimeUnit.SECONDS).thenRun(() -> {
    isReady.set(true);
    LOG.info("AMS is ready, application is now ready to serve requests");
}).exceptionally(ex -> {
    LOG.error("AMS failed to become ready within the timeout", ex);
    System.exit(1);
    return null;
});
```

:::

::: tip Expose AMS readiness in the Spring Boot Actuator health endpoint
Beyond the startup check, applications can expose the AMS readiness state through the **Spring Boot Actuator health endpoint** (for example, to observe the time since the last bundle refresh at runtime) by including one of the optional health starters:

- `spring-boot-starter-ams-health` for Spring Boot 4
- `spring-boot-3-starter-ams-health` for Spring Boot 3

These provide an auto-configured `HealthIndicator` bean that autowires the `AuthorizationManagementService` and includes its readiness state in the health status. They are *not* included transitively by any AMS starter.
:::