# Logging

This section describes how to log results of authorization checks for both [debugging](#debug-logging) and [auditing](#audit-logging) purposes.

## Debug Logging
Debug Logs for the AMS modules can be enabled via environment variables to analyze unexpected authorization results of the application.
They can also be used to understand and verify the behavior of the AMS modules in early stages of development.

::: danger Warning
Debug logging should not be enabled in production environments as it may expose sensitive information and impact performance.

It is your responsibility to ensure that debug logs are not stored or transmitted in a way that could compromise security or violate data protection regulations.
:::

::: tip
In CAP projects, you can use [hybrid testing](https://cap.cloud.sap/docs/advanced/hybrid-testing) to run a local instance of the application with the AMS bundle from a productive landscape where it is safe to enable debug logging.
:::

### Local Debug Logging
To enable debug logging when starting the application locally, run with the following environment variables:

::: code-group
```bash [CAP Node.js]
DEBUG=ams cds watch
```

```bash [Node.js]
DEBUG=ams npm start
```

```yaml [CAP Java]
# application.yaml
logging:
    level:
        com.sap.cloud.security.ams: DEBUG
        com.sap.cloud.security.ams.dcl.capsupport: DEBUG

# or via environment variables:
LOGGING_LEVEL_COM_SAP_CLOUD_SECURITY_AMS=DEBUG
LOGGING_LEVEL_COM_SAP_CLOUD_SECURITY_AMS_DCL_CAPSUPPORT=DEBUG 
```

```yaml [Java]
# application.yaml
logging:
    level:
        com.sap.cloud.security.ams: DEBUG

# or via environment variables:
LOGGING_LEVEL_COM_SAP_CLOUD_SECURITY_AMS=DEBUG
```
:::

### Cloud Debug Logging
When the application is running in a cloud environment, you can enable debug logging by setting the corresponding environment variables in the application's cloud environment configuration.

::: tip SAP BTP Cloud Foundry
In the SAP BTP Cloud Foundry environment, you can manage [user-provided variables](https://help.sap.com/docs/btp/sap-business-technology-platform/manage-environment-variables#loio9984a29f721e4981ad6a0b0b0cb6b868__section_wgl_w3f_32c) and restart the application to enable debug logging.
:::

## Audit Logging
Audit logs are used to record results of authorization checks for compliance and security purposes.

To write audit logs, it is necessary to register to events of the AMS module. The event payload contains the result of the authorization check, including the principal, resource, action, and result.

::: code-group
```js [CAP Node.js]
const { amsCapPluginRuntime } = require("@sap/ams");

cds.on('bootstrap', () => {
    amsCapPluginRuntime.ams.on("authorizationCheck", event => {
        // build audit log payload from event payload and write to audit log 
    });
});
```

```js [Node.js]
ams.on("authorizationCheck", event => {
    // build audit log payload from event payload and write to audit log 
});
```

```java [CAP Java + Java]
import com.sap.cloud.security.ams.dcl.client.pdp.PolicyDecisionPoint;
import static com.sap.cloud.security.ams.factory.AmsPolicyDecisionPointFactory.DEFAULT;

PolicyDecisionPoint policyDecisionPoint = PolicyDecisionPoint.create(DEFAULT);
policyDecisionPoint.registerListener(new AmsAuditLogger());

public class AmsAuditLogger implements Consumer<PolicyEvaluationResult> {
    @Override
    public void accept(PolicyEvaluationResult result) {
        // build audit log payload from result and write to audit log 
    }
}
```
:::

[Node.js Details](/nodejs/sap_ams/sap_ams.md#Logging) / [Java Details](/java/jakarta-ams/jakarta-ams.md#Logging)


### Request correlation
When logging results, it is typically necessary to correlate the audit log entries with the original request, e.g. based on a `correlation_id`. Please refer to the details of the AMS module used in your application to learn how to this.