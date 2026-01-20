# Logging

This section describes how to log results of authorization checks for both [debugging](#debug-logging) and [auditing](#audit-logging) purposes.

## Debug Logging
Debug Logs for the Authorization Management Service (**AMS**) modules can be enabled via environment variables to analyze unexpected authorization results of the application.
They can also be used to understand and verify the behavior of the AMS modules in early stages of development.

::: danger Warning
Debug logging shouldn't be enabled in production environments. It might expose sensitive information and impact performance.

It's your responsibility to ensure that debug logs aren't stored or transmitted in a way that could compromise security or violate data protection regulations.
:::

::: tip
In CAP projects, you can use [hybrid testing](https://cap.cloud.sap/docs/advanced/hybrid-testing) to run a local instance of the application with the AMS bundle from a productive landscape where it is safe to enable debug logging.
:::

### Local Debug Logging
To enable debug logging when starting the application locally, enable it as follows:

::: code-group
```bash [CAP Node.js]
DEBUG=ams cds watch
```

```bash [Node.js]
DEBUG=ams npm start
```

```yaml [Spring / Spring Boot (CAP)]
# application.yaml
logging:
    level:
        com.sap.cloud.security.ams: DEBUG

# or via environment variables:
LOGGING_LEVEL_COM_SAP_CLOUD_SECURITY_AMS=DEBUG
```

``` [log4j.properties]
log4j.logger.com.sap.cloud.security.ams=DEBUG
```
:::

### Cloud Debug Logging
When the application is running in a cloud environment, you can enable debug logging by setting the corresponding environment variables in the application's cloud environment configuration.

::: tip SAP BTP Cloud Foundry
In the Cloud Foundry environment of SAP BTP, you can manage [user-provided variables](https://help.sap.com/docs/btp/sap-business-technology-platform/manage-environment-variables#loio9984a29f721e4981ad6a0b0b0cb6b868__section_wgl_w3f_32c) and restart the application to enable debug logging.
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

```java [Spring / Spring Boot (CAP)]
import com.sap.cloud.security.ams.api.AuthorizationCheckListener;
import com.sap.cloud.security.ams.events.*;

@Configuration
public class AmsEventConfiguration {

    @Autowired
    private AuthorizationManagementService ams;

    @PostConstruct
    public void registerAuthorizationCheckListener() {
        // Same handler for all events...
        ams.addAuthorizationCheckListener(
            AuthorizationCheckListener.fromConsumer(
                (AuthorizationCheckEvent event) -> AMS_EVENT_LOGGER.debug(event.toString()))
        );
        
        // ... or dedicated handlers for each event type
        ams.addAuthorizationCheckListener(new AuthorizationCheckListener() {
            @Override
            public void onPrivilegeCheck(PrivilegeCheckEvent event) {}     

            @Override
            public void onGetPotentialActions(GetPotentialActionsEvent event) {}

            @Override
            public void onGetPotentialResources(GetPotentialResourcesEvent event) {}
            
            @Override
            public void onGetPotentialPrivileges(GetPotentialPrivilegesEvent event) {}
        });
    }
}
```

```java [Java]
import com.sap.cloud.security.ams.api.*;
import com.sap.cloud.security.ams.events.*;
import jakarta.annotation.PostConstruct;
import org.springframework.context.annotation.Configuration;
import org.springframework.beans.factory.annotation.Autowired;

@Configuration
public class AmsEventConfiguration {

    @Autowired
    private AuthorizationManagementService ams;

    @PostConstruct
    public void registerAuthorizationCheckListener() {
        // Same handler for all events...
        ams.addAuthorizationCheckListener(
            AuthorizationCheckListener.fromConsumer(
                (AuthorizationCheckEvent event) -> AMS_EVENT_LOGGER.debug(event.toString()))
        );
        
        // ... or dedicated handlers for each event type
        ams.addAuthorizationCheckListener(new AuthorizationCheckListener() {
            @Override
            public void onPrivilegeCheck(PrivilegeCheckEvent event) {}

            @Override
            public void onGetPotentialActions(GetPotentialActionsEvent event) {}

            @Override
            public void onGetPotentialResources(GetPotentialResourcesEvent event) {}


            @Override
            public void onGetPotentialPrivileges(GetPotentialPrivilegesEvent event) {}
        });
    }
}
```
:::

[Node.js Details](/Libraries/nodejs/sap_ams/sap_ams.md#events-logging) / [Java Details](/Libraries/java/jakarta/jakarta-ams.md#logging)


### Request correlation
When logging results, it is typically necessary to correlate the audit log entries with the original request, e.g. based on a `correlation_id`.

In Node.js, this information is carried by the event payload. In Java, access depends on your application's web framework and is typically achieved by thread-local storage.