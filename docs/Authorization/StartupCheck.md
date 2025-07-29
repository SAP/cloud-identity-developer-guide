#  Startup Check
The AMS client library modules need to load the authorization model, i.e. the AMS schema and base policies, before they perform authorization checks. Otherwise, runtime errors will be thrown when the application attempts to perform authorization checks.

::: danger Important
The application must ensure with a startup check that the AMS client library modules are ready before serving requests that require authorization checks.
:::

## Implementation
There are two strategies for the startup check:
1. **Synchronous Check**: The application waits synchronously for the AMS module during application startup to become ready before proceeding. This method will slightly delay the remaining setup, but typically not for long unless there is a problem with the AMS server.
1. **Asynchronous Check**: The application finishes its setup and exposes its endpoints. Unauthorized endpoints, such as a health check endpoint, can be served already while AMS is not ready. To prevent authorized endpoints from being called too early, it must include the readiness status of the AMS module in its health / readiness status reported to the cloud platform.

Some AMS modules support both strategies, while others only support one of them. In the former case, the application can choose the strategy that best fits its needs.

::: tip
In the asynchronous strategy, it should make sure to eventually throw an error when the AMS module does not become ready within a certain time frame.
:::

::: code-group

```js [CAP Node.js]
// server.js
const cds = require('@sap/cds');
const { amsCapPluginRuntime } = require("@sap/ams");

cds.on('served', async () => {
    // synchronous startup check by delaying remaining event handling up to 5s before throwing an error
    await amsCapPluginRuntime.ams.whenReady(5);
    console.log("AMS has become ready.");
});

// *or*: use amsCapPluginRuntime.ams.isReady() in /health endpoint
```

```js [Node.js]
// example for asynchronous startup check

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
        console.log('AMS has become ready.');
    } catch (e) {
        console.error('AMS did not become become ready in time:', e);
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
import com.sap.cloud.security.ams.dcl.client.pdp.PolicyDecisionPoint;
import static com.sap.cloud.security.ams.dcl.client.pdp.PolicyDecisionPoint.Parameters.STARTUP_HEALTH_CHECK_TIMEOUT;
import static com.sap.cloud.security.ams.dcl.client.pdp.PolicyDecisionPoint.Parameters.FAIL_ON_STARTUP_HEALTH_CHECK;

// Throws an error if the AMS module does not become ready within 5 seconds
PolicyDecisionPoint policyDecisionPoint = PolicyDecisionPoint.create(
    DEFAULT, 
    STARTUP_HEALTH_CHECK_TIMEOUT, 5L,
    FAIL_ON_STARTUP_HEALTH_CHECK, true
); 
```
:::