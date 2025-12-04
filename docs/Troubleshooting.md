# Troubleshooting Guide

## Authorization Problems

This guide provides a systematic approach to troubleshooting authorization problems, such as unexpected authorization results in your application. By following these steps, the guide helps you to efficiently identify the root cause of the issue.

### Step 1: Preparation

Before you start debugging, make sure you can reproduce the issue in an environment where one of two things are possible:
- Enabling **[Debug Logs](/Authorization/Logging#debug-logging)** without violating data protection regulations by logging sensitive information.
- Debugging the application with a debugger attached to the process.

::: tip
There are very few issues that occur only in production environments and require analysis in a deployed application.

Most issues can be reproduced and fixed **much faster** by writing a test for the failing authorization logic, even cross-application integration via technical communication.
Refer to the **[Testing Guide](/Authorization/Testing)** for guidance.
:::

### Step 2: Analyze Authorization Check

Reproduce the issue and check the debug logs to see which part of the authorization check doesn't meet your expectations.

#### CAP role check

In CAP applications, the following output is expected when the AMS plugin determines cds roles from authorization policies:

::: code-group
```bash [Node.js]
[ams] - Determined potential actions for resource '$SCOPES': Reader {
      resource: '$SCOPES',
      potentialActions: Set(1) { 'Reader' },
      policies: [ 'local.JuniorReader' ],
      limitingPolicies: [ 'internal.ReadCatalog' ],
      correlation_id: '057c0278-1d2f-4797-8762-435e948db42d'
    }
```

```text [Java]
An example for debug log output in Java is coming soon.
```
:::

- Which **principal**/**policies** was used? -> `local.JuniorReader`
- Which **roles** have been granted? -> `Reader`
- Has any customization logic been applied that influenced the check, such as technical communication logic? -> yes, the `internal.ReadCatalog` policy is used as an upper limit for the roles due to consumption of an App-to-App provided API that was mapped to this policy.

::: info
If the CAP application uses `@ams.attributes` annotations for instance-based authorization, additional [generic privilege checks](#generic-privilege-check) are logged per role to determine the filter conditions.
:::

#### Generic Privilege check

::: code-group
```bash [Node.js]
[ams] - Privilege check for 'Reader' on '$SCOPES' was conditional. {
      action: 'Reader',
      resource: '$SCOPES',
      input: {},
      defaultInput: { '$user.scim_id': '7c1d5056-b334-44ab-b2e1-14afba1f3b8b' },
      result: 'conditional',
      dcn: "($app.genre IN ['Fantasy', 'Fairy Tale', 'Mystery'] AND $app.stock < 30)",
      policies: [ 'local.JuniorReader' ],
      limitingPolicies: [ 'internal.ReadCatalog' ],
      correlation_id: '057c0278-1d2f-4797-8762-435e948db42d'
    }
```

```text [Java]
An example for debug log output in Java is coming soon.
As of today, the logs show the content of the Attributes object.
The limiting policies can be inferred from the scopeFilter property.
```

:::

- Which **action** is checked on which **resource**? -> `Reader` on `$SCOPES`
::: tip CAP applications
The action in CAP application is the cds role name for which a check is performed.
The resource is always the special resource [`$SCOPES`](/CAP/Basics#role-policies) and *not* the entity name.
:::
- Which **principal**/**policies** was used? -> `local.JuniorReader`
- Which **attribute input** was provided? -> no check-spefific input, only default input for `$user.scim_id`
- What was the resulting **DCN condition** from the authorization engine? -> `($app.genre IN ['Fantasy', 'Fairy Tale', 'Mystery'] AND $app.stock < 30)`
- Has any customization logic been applied that influenced the check, such as technical communication logic? -> yes, the `internal.ReadCatalog` policy is used as an upper limit for the permissions due to consumption of an App-to-App provided API that was mapped to this policy.

### Step 3: Check for Common Root Causes

In case your logs contain unexpected information, check for explanations in the common root causes below before opening a new ticket.

### Step 4: Solving the Issue

If you've identified a misconfiguration or a setup issue in your application, try to resolve the issue based on your findings and the documentation.

::: info Consulting
If you need help with fixing your application setup, this is a **consulting** request.
:::

For both consulting and issues inside the AMS modules or services, such as bugs or outages, please refer to our **[Support Guide](/Support)** for how to get help.



## Common Root Causes

### Unexpected `403 Forbidden`

#### Outdated Dependencies

Ensure that you are using the correct combination of the Authorization Management Service (**AMS**) modules for your project setup and in the latest versions as recommended in the **[Getting Started Guide](/Authorization/GettingStarted.md)**.

#### Wrong policies used

You should see debug output like this that shows how the list of policies for an authorization check was determined:

::: code-group
```bash [Node.js]
[ams] - User with scim_id  7c1d5056-b334-44ab-b2e1-14afba1f3b8b  in tenant  93127375-beab-43c7-b5a8-c1bc2265997d  has the following assigned policies:  [ 'local.JuniorReader' ]
```

```text [Java]
An example for debug log output in Java is coming soon.
```
:::

If the policies listed here do not match your expectations, double-check the assignment of the policies in the administration console.

Check that the application successfully received a bundle update from the AMS server after a short waiting period.
Policy assignment changes usually take up to 15 seconds to propagate due to caching (but can sometimes take longer). Wait a moment before retrying.

::: warning Missing policy assignments in Token
As of today, the policies that are assigned to a user are **not** contained in his SCI user tokens. This is normal.

A token refresh is **NOT** necessary after making changes to policy assignments. The AMS module retrieves the user policies separately from the token directly from the AMS server.
:::

If the policies listed here are different from the ones used for the authorization check (see Step 2), check for any customization logic that modifies the policies used for the check, such as overridden `AuthProvider` logic or an `AttributeProcessor` implementation in Java.

#### Unexpected limiting policies

##### Typo in API names

Make sure there are no typos in API names in the mapping functions.

::: warning
If you want to enable the special [`principal-propagation`](/Authorization/TechnicalCommunication#authorization-via-api-permission-groups) API, make sure to name it correctly. The name is case-sensitive and must be exactly `principal-propagation` and *not* `principle-propagation`.
:::

##### Incorrect Policy Names 
Policy names must be **fully-qualified**, including the DCL package name (for example, `cap.Reader`). This applies to:
   -   Mocked policy assignments in tests.
   -   Policy mappings for technical communication.

### Unexpected `2xx` Response

If you receive an unexpected successful response (HTTP 2xx) when you expect a `403 Forbidden`, check the following common root causes:

#### Left-over policy assignments

- A user might have multiple policies assigned that grant the same privilege. Ensure you've unassigned all of them before testing that a specific policy doesn't grant access.

#### Non-AMS authorization logic

Check if there is any non-AMS authorization logic in your application that could grant access independently from AMS policies.

For example, the standard CAP authentication handlers use the list of ias_apis to automatically grant cds roles with the same name in case of technical user tokens.

### Unexpected HTTP response in CAP applications

The AMS CAP plugins are not responsible for the HTTP response codes or messages returned by the application. Instead, they integrate into the CAP framework itself, e.g. by providing user roles and filter conditions, but it is the framework itself which decides whether requests are forbidden and with which message. In particular, users may receive roles from other sources than the AMS plugin, which can therefore lead to access being granted even though the roles from the AMS plugin would not grant access.

### Authorization bundle issues

- **No authorization data loaded error**: The application must wait for the **[AMS Startup Check](/Authorization/StartupCheck)** before making authorization checks.

::: warning
The startup check is also necessary before unit tests that perform authorization checks.
:::

- **Incorrect Test Setup**: Make sure to follow the **[Testing Guide](/Authorization/Testing.md)** carefully by including all steps, such as DCN compilation and DCN loading.