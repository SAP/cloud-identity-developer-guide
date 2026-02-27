# ams-core

`ams-core` is the core module used to perform [authorization checks](https://help.sap.com/docs/identity-authentication/identity-authentication/configuring-authorization-policies?locale=en-US) in Java applications based on AMS policies.

## Installation

The module is available as a Maven dependency:

```xml
<dependency>
    <groupId>com.sap.cloud.security.ams</groupId>
    <artifactId>ams-core</artifactId>
</dependency>
```

## Usage

For detailed usage instructions including initialization, authorization checks, and error handling, see:

- [Authorization Bundle](/Authorization/AuthorizationBundle) - Client library initialization and startup checks
- [Authorization Checks](/Authorization/AuthorizationChecks) - Performing privilege checks and handling decisions

### Events/Logging

Consumer applications can listen to authorization check events of the `AuthorizationManagementService` instance to manually log authorization check results and/or create audit log events.

#### Authorization Check Events

Register an `AuthorizationCheckListener` to receive notifications about authorization checks:

```java
import com.sap.cloud.security.ams.api.AuthorizationCheckListener;
import com.sap.cloud.security.ams.api.DecisionResult;
import com.sap.cloud.security.ams.events.*;

// Option 1: Implement the AuthorizationCheckListener interface
ams.addAuthorizationCheckListener(new AuthorizationCheckListener() {
    @Override
    public void onPrivilegeCheck(PrivilegeCheckEvent event) {
        if (event.getResult() == DecisionResult.GRANTED) {
            logger.info("Privilege '{}' on '{}' was granted", 
                event.getAction(), event.getResource());
        }
    }

    @Override
    public void onGetPotentialActions(GetPotentialActionsEvent event) {
        logger.debug("Potential actions for '{}': {}", 
            event.getResource(), event.getPotentialActions());
    }

    @Override
    public void onGetPotentialResources(GetPotentialResourcesEvent event) {
        logger.debug("Potential resources: {}", event.getPotentialResources());
    }

    @Override
    public void onGetPotentialPrivileges(GetPotentialPrivilegesEvent event) {
        logger.debug("Potential privileges: {}", event.getPotentialPrivileges());
    }
});

// Option 2: Use the convenience method for simple logging scenarios
ams.addAuthorizationCheckListener(
    AuthorizationCheckListener.fromConsumer(event -> {
        logger.info("Authorization check: {}", event.getDescription());
    })
);
```

#### Event Types

The following event types are emitted during authorization operations:

| Event Type | Description | Key Properties |
|------------|-------------|----------------|
| `PrivilegeCheckEvent` | Emitted during a privilege check | `action`, `resource`, `result`, `input`, `defaultInput`, `dcn` |
| `GetPotentialActionsEvent` | Emitted when collecting potential actions for a resource | `resource`, `potentialActions` |
| `GetPotentialResourcesEvent` | Emitted when collecting potential resources | `potentialResources` |
| `GetPotentialPrivilegesEvent` | Emitted when collecting potential privileges | `potentialPrivileges` |

All events inherit from `AuthorizationCheckEvent` and include:
- `policies`: The fully-qualified policy names used for the check
- `limitingPolicies`: Policy names whose privileges were used as upper limit (if any)
- `description`: A human-readable description of the authorization check

#### PrivilegeCheckEvent Details

The `PrivilegeCheckEvent` provides comprehensive information about privilege checks:

```java
ams.addAuthorizationCheckListener(new AuthorizationCheckListener() {
    @Override
    public void onPrivilegeCheck(PrivilegeCheckEvent event) {
        String action = event.getAction();           // e.g., "read"
        String resource = event.getResource();       // e.g., "orders"
        DecisionResult result = event.getResult();   // GRANTED, DENIED, or CONDITIONAL
        String dcn = event.getDcn();                 // DCN condition in human-readable format
        
        // Input attributes used for the check
        Map<AttributeName, Object> input = event.getInput();
        Map<AttributeName, Object> defaultInput = event.getDefaultInput();
        
        // Policies involved
        Set<PolicyName> policies = event.getPolicies();
        Set<PolicyName> limitingPolicies = event.getLimitingPolicies();
    }
    
    // ... other methods
});
```

### Bundle Loading Error Handling

AMS uses a bundle loader internally to manage the policies and assignments bundle in the background, independently of incoming requests. The `AuthorizationManagementService` emits error events when bundle loading requests fail.

#### Error Event Types

There are two distinct error event types:

- **`BundleInitializationErrorEvent`**: Emitted when the initial bundle download fails and the instance is not yet ready for use.
- **`BundleRefreshErrorEvent`**: Emitted when a bundle refresh request fails. Since the library continuously polls, this doesn't necessarily mean the data is outdated, just that the polling attempt failed. The instance remains ready but if there have been recent policy or assignment changes, it cannot take them into account.

#### Handling Bundle Errors

Register an error listener to handle bundle loading errors:

```java
import com.sap.cloud.security.ams.error.AmsBackgroundException;
import com.sap.cloud.security.ams.error.BundleInitializationErrorEvent;
import com.sap.cloud.security.ams.error.BundleRefreshErrorEvent;

ams.addErrorListener(event -> {
    if (event instanceof BundleInitializationErrorEvent) {
        logger.error("AMS bundle initialization failed - service not ready: {}", 
            event.getException().getMessage());
        // Eventually the cloud platform will restart the application after a 
        // certain number of failed attempts to the health endpoint, so
        // typically no action besides logging is required here
    } else if (event instanceof BundleRefreshErrorEvent refreshError) {
        logger.warn("AMS bundle refresh failed (current bundle age: {} seconds): {}",
            refreshError.getSecondsSinceLastRefresh(),
            refreshError.getException().getMessage());
        // Consider taking action such as logging an error instead of a warning
        // when the bundle is stale for extended periods of time
    }
});
```

::: info Automatic Error Logging
If your application does not register any error listeners, bundle loading errors will be automatically logged:
- `BundleInitializationErrorEvent`: Logged at ERROR level
- `BundleRefreshErrorEvent`: Logged at WARN level with the seconds since last refresh
:::

::: tip Handling Initial Bundle Load Errors
Refer to the [Authorization Bundle](/Authorization/AuthorizationBundle) documentation for guidance on how to react when AMS fails to initialize the bundle. The error events emitted for this case are only intended to provide information about the failed requests.
:::