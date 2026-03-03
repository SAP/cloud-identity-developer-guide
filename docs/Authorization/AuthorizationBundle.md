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

##  Startup Check

While it is possible to synchronously block application startup until the AMS module becomes ready, we recommend including AMS in the application's **readiness probes**. This allows the application process to become healthy for the cloud platform but prevent traffic from being routed to the process until the AMS module is ready to serve authorization checks. 

If an application does not provide readiness probes, it can alternatively include the AMS readiness state in its **health endpoint**.

::: code-group
```js [Node.js (CAP)]
// server.js
const cds = require('@sap/cds');
const { amsCapPluginRuntime } = require("@sap/ams");

cds.on('served', async () => {
    // effectively a synchronous startup check that prevents the server
    // from serving requests for up to 30s before throwing an error
    await amsCapPluginRuntime.ams.whenReady(30);
    console.log("AMS has become ready.");
});

// *better*: use amsCapPluginRuntime.ams.isReady() in health or readiness endpoint
```

```js [Node.js]
// uses readiness status in health endpoint

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

```java [Java]
import java.util.concurrent.TimeUnit;
import java.util.concurrent.atomic.AtomicBoolean;

// Synchronous startup check:
// throws an error if the AMS module doesn't get ready within 30 seconds
ams.whenReady().get(30, TimeUnit.SECONDS);

// Asynchronous startup check: uses readiness status in health endpoint
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

```java [Spring Boot Readiness]
// The spring-boot-starter-ams-readiness module provides an auto-config for a
// AmsReadinessContributor bean. It autowires an AuthorizationManagementService
// instance and uses AvailabilityChangeEvent.publish
// to integrate its readiness state with Spring Boot's availability state.
// The AMS readiness starter is included transitively in all AMS Spring Boot starters.
```

```java [Spring Boot Health Actuator]
// The spring-boot-starter-ams-health and spring-boot-3-starter-ams-health modules
// provide an auto-config for a HealthIndicator bean that integrates with the
// Spring Boot Actuator health endpoint. It autowires an
// AuthorizationManagementService instance and includes its readiness state
// in the health status.
// The AMS health starter is NOT included transitively in any AMS Spring Boot starter.
```

:::