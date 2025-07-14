# Troubleshooting
This section describes a systematic approach to troubleshooting issues in applications using AMS modules, with a focus on unexpected authorization results.

## Analyzing Issues
To analyze unexpected authorization behavior of your application, please follow these steps:

1. Rule out authentication issues:
   - Ensure that the user is authenticated and that is is an authorization issue, not an authentication issue.
1. Setup a reproducible test case in a safe environment:
   - Make sure you can reproduce the issue in a safe environment, preferably in a local unit test setup that isolates the problem. If the error occurs in production, consider mocking a security context from the token's claims in a local test setup.
1. Enable [Debug Logs](/concepts/Logging):
   - Enable debug logs in the safe environment to get detailed information about the authorization process. This will help you identify the root cause of the issue.
1. Reproduce the request that causes the issue in the safe environment.
1. Understand the authorization check:
    - Which **action** is checked on which **resource**?
    - Which **policies** have been applied for the check?
    - Does the authorization check depend on schema **attributes** and if so, which condition do you expect based on the assigned policies?
    - Which **attribute inputs** have been used for the check and which attributes have not been grounded to a value (= *unknown* attributes)?
    - What is the **DCN condition** returned by the DCN engine for the check?
    - Has additional logic for **technical communication** been applied before the check, e.g. an upper authorization limit based on consumed API or an overwrite of the used policies based on consumed API?
1. Find Root Cause:
    - Based on the debug logs and your understanding of the authorization check, identify the root cause of the issue, e.g. *Incorrect mocked policy assignment*.\
    For a list of common root causes, see the section below.

## Solving Issues
Try to solve the issue based on your analysis and double-checking the documentation for guidance on specific features or configurations that may be relevant to your issue.

If you cannot solve the problem on your own, there are two options:

1. You believe you encountered a bug in the AMS modules (-> **Support ticket**)
1. You believe your application setup is not correct and you need help fixing it (-> **Consulting**)

For both of these options, we have documented the recommended steps to take in the [Support](/Support) guide.

## Common Root Causes

### Unexpected 403
If your application is returning a 403 Forbidden error unexpectedly, it may be due to one of the following reasons:

- Wrong or outdated AMS modules installed for the project setup.\
**Solution**: Refer to the [Getting Started](/concepts/GettingStarted.md) guide to ensure you have the recommended (combination of) AMS modules installed for your project type and **update to the latest versions**.

- Incorrect mocked policy assignment.\
**Solution**: Double-check the mocked policy assignments for [syntactical correctness](/concepts/Testing#assigning-policies-to-mocked-users). Make sure that the policy names are **fully-qualified** names, i.e. they are prefixed with the DCL package in which the policy is defined, e.g. `cap.Reader`.

- Incorrect API permission group -> Policy mapping.\
**Solution**: Ensure that the API permission group are [correctly mapped](/concepts/TechnicalCommunication#mapping-implementation) to **fully-qualified policy names**, i.e. the resulting policy names are prefixed with the DCL package in which the policy is defined, e.g. `cap.Reader`.

- Cached policy assignments.\
**Solution**: After changing policy assignments, give it some time before retrying the request. The policy assignments are cached and regularly refreshed. Typically, it takes about 15s for a change to take effect.

### Unexpected access granted
If your application is granting access to resources unexpectedly, consider the following:

- Multiple policies assigned.\
**Solution**: When removing policies to manually check that access is no longer granted without the policy, make sure there are no other policies assigned to the user that grant the necessary privilege.

- Cached policy assignments.\
**Solution**: After changing policy assignments, give it some time before retrying the request. The policy assignments are cached and regularly refreshed. Typically, it takes about 15s for a change to take effect.

### No authorization data loaded/received errors
If your application is experiencing errors related to no authorization data being loaded or received, consider the following:

- Incorrect test setup.\
**Solution**: Ensure that all steps of the [Testing](/concepts/Testing.md) guide are followed correctly, including the setup of the DCL compiler and loading of the DCN bundle.

- Missing startup check.\
**Solution**: Ensure that you have correctly implemented a [startup check](/concepts/Setup#startup-check) in your application that is successful before your application starts processing requests.