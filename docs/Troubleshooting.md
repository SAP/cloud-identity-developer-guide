# Troubleshooting Guide

## Authorization Problems

This guide provides a systematic approach to troubleshooting authorization problems, such as unexpected authorization results in your application. Following these steps will help you identify the root cause of the issue efficiently.

### Step 1: Preparation

Before you start debugging, make sure you can reproduce the issue in an environment where one of two things are possible:
- Enabling **[Debug Logs](/concepts/Logging#debug-logging)** without violating data protection regulations by logging sensitive information.
- Debugging the application with a debugger attached to the process.

::: tip
There are very few issues that occur only in production environments and require analysis in a deployed application.

Most issues can be reproduced and fixed **much faster** by writing a test for the failing authorization logic, even cross-application integration via technical communication.
Refer to the **[Testing Guide](/concepts/Testing)** for guidance.
:::

### Step 2: Analyze Authorization Check

Reproduce the issue and check which part of the authorization check does not meet your expectations:

- Which **action** is checked on which **resource**?
- Which **principal**/**policies** was used?
- Which **attribute input** was provided?
- What was the resulting **DCN condition** from the authorization engine?
- Has any customization logic been applied that influenced the check, such as technical communication logic?

### Step 3: Check for Common Root Causes

The following are common root causes for authorization problems:

### Unexpected `403 Forbidden`

- **Outdated Dependencies**: Ensure you are using the correct combination of the AMS modules for your project setup and in the latest versions as recommended in the **[Getting Started Guide](/concepts/GettingStarted.md)**.
- **Caching Delays**: Policy assignment changes usually take up to 15 seconds to propagate due to caching (but can sometimes take longer). Wait a moment before retrying.

::: warning Missing policy assignments in Token
The policies that are assigned to a user are **not** contained in his SCI user tokens. This is normal.

A token refresh is **NOT** necessary after making changes to policy assignments. The AMS module retrieves the user policies separately from the token directly from the AMS server.
:::

- **Incorrect Policy Names**: Policies must be **fully-qualified**, including the DCL package name (e.g., `cap.Reader`). This applies to:
   -   Mocked policy assignments in tests.
   -   Policy mappings for technical communication.

#### Unexpected `Access Granted`

- A user might have multiple policies assigned that grant the same privilege. Ensure you've unassigned all of them before testing that a specific policy does not grant access.

#### Authorization bundle issues

- **No authorization data loaded error**: The application must wait for the **[AMS Startup Check](/concepts/StartupCheck)** before making authorization checks.

::: warning
The startup check is also necessary before unit tests that perform authorization checks.
:::

- **Incorrect Test Setup**: Make sure to follow the **[Testing Guide](/concepts/Testing.md)** carefully by including all steps, such as DCN compilation and DCN loading.

### Step 4: Solving the Issue

If you've identified a misconfiguration or a setup issue in your application, try to resolve the issue based on your findings and the documentation.

::: info Consulting
If you need help with fixing your application setup, this is a **consulting** request.
:::

For both consulting and issues inside the AMS modules or services, such as bugs or outages, please refer to our **[Support Guide](/Support)** for how to get help.