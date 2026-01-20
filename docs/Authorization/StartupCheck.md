#  Startup Check
The client library modules for the Authorization Management Service (**AMS**) must load the authorization model - the schema and policies - before they perform authorization checks. Otherwise, runtime errors are thrown when the application attempts to perform authorization checks.

::: danger Important
The application must ensure with a startup check that the AMS client library modules are ready before serving requests that require authorization checks.
:::

## Implementation
There are two strategies for the startup check:
1. **Synchronous Check**: The application waits synchronously for the AMS module during application startup until it becomes ready before proceeding. This method slightly delays the remaining setup, but typically not for long unless there is a problem with the AMS server.
1. **Asynchronous Check**: The application finishes its setup and exposes its endpoints. Unauthorized endpoints, such as a health check endpoint, can be served already while AMS is not ready. To prevent authorized endpoints from being called too early, it must include the readiness status of the AMS module in its health / readiness status reported to the cloud platform.

Some AMS modules support both strategies, while others only support one of them. In the former case, the application can choose the strategy that best fits its needs.

::: tip
In the asynchronous strategy, make sure that an error is thrown eventually when the AMS module isn't ready within a certain time frame.
:::

::: code-group
```js [Node.js (CAP)]
// server.js
const cds = require('@sap/cds');
const { amsCapPluginRuntime } = require("@sap/ams");

cds.on('served', async () => {
    // synchronous startup check that prolongs the event handling up to 5s before throwing an error
    await amsCapPluginRuntime.ams.whenReady(5);
    console.log("AMS has become ready.");
});

// *or*: use amsCapPluginRuntime.ams.isReady() in /health endpoint
```

```js [Node.js]
// Asynchronous startup check: uses readiness status in health endpoint

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

```java [Spring Boot/Spring Boot (CAP)]
// The spring-boot-starter-ams and spring-boot-starter-cap-ams modules
// have auto-config for a HealthIndicator that integrates with the
// Spring Boot Actuator health endpoint.

@Bean
@ConditionalOnMissingBean(name = "amsHealthIndicator")
public HealthIndicator amsHealthIndicator(AuthorizationManagementService ams) {
    LOG.debug("Creating AMS health indicator");
    return () -> {
        if (ams.isReady()) {
            Long secondsSinceRefresh = ams.getBundleAge();
            return Health.up()
                    .withDetail("bundleAge", secondsSinceRefresh != null ? secondsSinceRefresh + "s" : "?")
                    .build();
        } else {
            return Health.down()
                    .withDetail("reason", "Initial authorization bundle not yet received")
                    .build();
        }
    };
}
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

:::