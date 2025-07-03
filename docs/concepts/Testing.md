# Testing AMS Integration

This guide explains how to efficiently test your AMS (Authorization Management Service) integration, ensuring fast feedback and reliable authorization logic.

## Why and When?

::: tip
Write tests early to save time and get your AMS setup running quickly. Well-designed tests help you catch issues in your authorization logic before deployment.
:::

To maximize productivity, tests should be executable locally—without requiring cloud resources or deploying the application. This enables rapid feedback cycles and lets you iterate on your authorization policies and application logic efficiently.

Local tests also make it easy to use a debugger or analyze debug logs to understand unexpected authorization check behavior. This is invaluable for [troubleshooting](/Troubleshooting).

::: warning TEST BEFORE DEPLOYMENT
Nothing is more frustrating—both for you and the AMS support team—than spending time on a cloud deployment with limited debuggability, only to find out that your authorization logic is not working as expected. Local tests help you avoid this pitfall.
:::

## Mock Authentication, Not Authorization

When testing your application's AMS integration, always mock authentication—not authorization. This way you test the same AMS code that runs in production which means, local authorization tests are generally very reliable. Avoid mocking internal authorization logic, as it is complex and may change, making tests brittle and less reliable.

The authorization checks in the AMS libraries typically use the security context created by the BTP authentication libraries as input. To test effectively, it is best practice to mock the security context directly—either by mocking its properties or by constructing it from a test *JWT token*.

::: info
In CAP, simply use the standard [mock users](https://cap.cloud.sap/docs/node.js/authentication#mock-users) for testing.
:::

The following SAP Identity Service token claims are relevant for authorization testing:
- `app_tid` - application-internal ID of the user's tenant (chosen by application in subscription callbacks)
- `scim_id` - SCIM ID of the user (used in conjunction with `app_tid` to determine assigned policies)
- `ias_apis` - contains the consumed SCI APIs in technical communication flows
- `azp`/ `sub` - used to distinguish between regular user and technical user tokens (technical user tokens have `azp = sub`)

[Node.js Sample](https://github.com/SAP-samples/ams-samples-node/blob/main/ams-express-shopping/auth/authenticate.js#L47-L59) / [Java Samples](https://github.com/SAP/cloud-security-services-integration-library/tree/main/java-security-test#samples)

## Testing Without AMS Cloud Instance

For local tests without an AMS cloud instance, the following steps are required:

1. Assign policies to mocked users, so that the authorization checks in the tests use these policies.
1. Configure the AMS client library to use the mocked policy assignments and load a local version of the authorization model.

The two steps are described in the following sections.

### Assigning Policies to Mocked Users

It depends on the application type how policies are assigned to mocked users. In CAP applications, they are assigned to mock users. In non-CAP applications, it depends on the language of the AMS client library.

::: code-group
```json [CAP Node.js] 
// cds.env source
{
    "requires": {
        "auth": {
            "[development]": {
                "kind": "mocked",
                "users": {
                    "alice": {
                        "policies": [ // [!code ++:4]
                            "shopping.CreateOrders",
                            "shopping.DeleteOrders"
                        ],
                    },
                    "bob": {
                        "policies": [ // [!code ++:3]
                            "local.OrderAccessory"
                        ]
                    } 
```

```yml [CAP Java]
# application.yml
cds:
  security:
    mock.users:
      alice:
        policies: // [!code ++:3]
          - shopping.CreateOrders
          - shopping.DeleteOrders
      bob:
        policies: // [!code ++:2]
          - local.OrderAccessory
```

```json [Node.js] 
// mockPolicyAssignments.json
{
    "defaultTenant": {
        "alice": [
            "shopping.CreateOrders",
            "shopping.DeleteOrders"
        ],
        "bob": [
            "local.OrderAccessory"
        ],
        "carol": []
    }
}
```

```java [Java]
// OrderControllerTest.java
@Test
  void requestWithCreateOrders_ok(SecurityTestContext context) throws IOException {
    String jwt =
        context
            .getPreconfiguredJwtGenerator()
            .withClaimValues("test_policies", "shopping.CreateOrders") // [!code ++]
            .createToken()
            .getTokenValue();

    HttpGet request = createGetRequest(jwt);
```
:::

### Loading local DCL bundle

To test with local DCL files, the AMS client library needs to be configured to load a local DCL bundle instead of the one from the AMS cloud instance. This is done by specifying the path to the local bundle. It depends on the language of the AMS client library how this is done.

#### Compiling DCL to DCN
Before running the tests, the local DCL files need to be compiled to DCN files as input for the AMS client library:

::: tip
In CAP Node.js projects, this is done automatically by `@sap/ams-dev` before `cds start/watch/test`. In other projects, you need to compile the DCL files manually.
:::

::: code-group
```json [Node.js]
// package.json
"scripts": {
        "jest": "NODE_ENV=test npx jest",
        "pretest": "npx compile-dcl -d auth/dcl -o test/dcn", // [!code focus]
},
"devDependencies": {
        "@sap/ams-dev": "^2", // [!code focus]
}
```

``` [Java]
https://github.wdf.sap.corp/CPSecurity/cloud-authorization-client-library-java/blob/3391eb3c4bd8ef9bf5c8a361d35379918571990e/docs/maven-plugins.md#-maven-dcl-compiler
```
:::

#### Loading DCN
To load the compiled DCN files, the AMS client library needs to be configured accordingly:

::: tip
In CAP Node.js projects, this is done automatically by the `@sap/ams` runtime when `requires.auth.kind = mocked`.
:::

::: code-group
```json [Node.js]
// package.json
"scripts": {
        "jest": "NODE_ENV=test npx jest",
        "pretest": "npx compile-dcl -d auth/dcl -o test/dcn", // [!code focus]
},
"devDependencies": {
        "@sap/ams-dev": "^2", // [!code focus]
}
```

``` [Java]
https://github.wdf.sap.corp/CPSecurity/cloud-authorization-client-library-java/blob/3391eb3c4bd8ef9bf5c8a361d35379918571990e/docs/maven-plugins.md#-maven-dcl-compiler
```
:::

##  Test policies

The DCL package called `local` has a special semantic: it is intended for DCL files with policies only relevant for testing, not for production. Policies in the `local` package are ignored during base policy upload, even if present in the archive.

This allows you to test policies that are restrictions of base policies without shipping them to customers. Typically, such policies would be created by an administrator at runtime in the SCI cockpit.