# Testing AMS Integration

This guide explains how to efficiently test your integration with the Authorization Management Service (**AMS**), ensuring fast feedback and reliable authorization logic.

## Why and When?

::: tip
Write tests early to save time and get your AMS setup running quickly. Well-designed tests help you catch issues in your authorization logic **before** deployment.
:::

To enhance productivity, tests should be executable locally — without requiring cloud resources or deploying the application. This enables rapid feedback cycles and lets you iterate on your authorization policies and application logic efficiently.

Local tests also make it easy to use a debugger or analyze debug logs without persisting sensitive information to understand unexpected authorization check behavior. This is invaluable for [troubleshooting](/Troubleshooting).

::: warning TEST BEFORE DEPLOYMENT
Nothing is more frustrating—both for you and the AMS support team—than spending time on a cloud deployment with limited debuggability, only to find out that your authorization logic is not working as expected. Local tests help you avoid this pitfall.
:::

::: tip CAP Hybrid Testing
In CAP projects, you can use [hybrid testing](https://cap.cloud.sap/docs/advanced/hybrid-testing) to run a local instance of the application with the AMS bundle from a productive landscape, allowing you to test with the same policies and assignments as in production.
:::

## Mock Authentication, Not Authorization

When testing the AMS integration of your application, always mock authentication, not authorization. This way, you test the same AMS code that runs in production. This means that local authorization tests are generally very reliable. Avoid mocking internal authorization logic, for example, library functions of AMS because it's complex and prone to API changes, making tests brittle and less reliable.

The authorization checks in the AMS libraries typically use the security context created by the BTP authentication libraries as input. To test effectively, it's best practice to mock the security context directly, either by mocking its properties or by constructing it from a test *JWT token*.

::: tip
In CAP, simply use the standard [mock users](https://cap.cloud.sap/docs/node.js/authentication#mock-users) for testing.
:::

The following SAP Identity Services token claims are relevant for authorization testing:
- `app_tid` - application-internal ID of the user's tenant (chosen by application in subscription callbacks)
- `scim_id` - SCIM ID of the user (used in conjunction with `app_tid` to determine assigned policies)
- `ias_apis` - contains the consumed APIs of SAP Cloud Identity Services in technical communication flows
- `azp`/ `sub` - used to distinguish between regular user and technical user tokens (technical user tokens have `azp = sub`)

[Node.js Sample](https://github.com/SAP-samples/ams-samples-node/blob/main/ams-express-shopping/auth/authenticate.js#L47-L59) / [Java Samples](https://github.com/SAP-samples/ams-samples-java/blob/main/jakarta-ams-sample/src/test/java/com/sap/cloud/security/samples/JavaServletIntegrationTest.java#L83-L97)

## Testing Without AMS Cloud Instance

For local tests without an AMS cloud instance, the following steps are required:

1. Assign policies to mocked users, so that the authorization checks in the tests use these policies.
1. Configure the AMS client library to use the mocked policy assignments and load a local version of the authorization model.

The two steps are described in the following sections.

### Assigning Policies to Mocked Users

In CAP applications, policies can be assigned to (both existing and custom) mocked users directly.\
In non-CAP applications, they are assigned to `app_tid` and `scim_id` pairs mocked during authentication:

::: code-group
```json [Node.js/Java] 
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

```json [Node.js (CAP)] 
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
                }
            }
        }
    }
} 
```

```yml [Java (CAP)]
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
:::

### Loading local DCL bundle

To test without an AMS cloud instance, the client library needs to use the local DCL files instead of the bundle from the AMS cloud instance.

#### Compiling DCL to DCN
Before running the tests, the local DCL files must be compiled to DCN files as input for the client library.

::: tip
In CAP Node.js projects, this is done automatically by `@sap/ams-dev` before `cds start/watch/test`. In other projects, you must manually set up a compilation step that runs before the tests.
:::

::: code-group
```json [Node.js]
// package.json
"scripts": {
        "jest": "NODE_ENV=test npx jest",
        "pretest": "npx compile-dcl -d auth/dcl -o test/dcn", // [!code ++]
},
"devDependencies": {
        "@sap/ams-dev": "^2", // [!code ++]
}
```

```xml [Java (CAP)]
<!-- srv/pom.xml -->
<build>
    <plugins>
        <!-- STEPS TO BUILD CDS MODEL AND GENERATE POJOs -->
        <plugin>
            <groupId>com.sap.cds</groupId>
            <artifactId>cds-maven-plugin</artifactId>
            <executions>
                <!-- ... -->

                <!-- DCL -> DCN compilation before tests --> // [!code ++:14]
                <execution>
                    <id>compile-dcl</id>
                    <goals>
                        <goal>npx</goal>
                    </goals>
                    <phase>generate-test-resources</phase>
                    <configuration>
                        <arguments>--package=@sap/ams-dev compile-dcl
                            -d ${project.basedir}/src/main/resources/ams
                            -o ${project.basedir}/target/generated-test-sources/dcn
                        </arguments>
                    </configuration>
                </execution>
            </executions>
        </plugin>
    </plugins>
</build>
```

```xml [Java]
<build>
        <plugins>
            <!-- ... -->

            <!-- Plugin for DCL -> DCN compilation before tests --> // [!code ++:31]
            <plugin>
                <groupId>com.github.eirslett</groupId>
                <artifactId>frontend-maven-plugin</artifactId>
                <version>1.14.1</version>
                <executions>
                    <execution>
                        <id>install node and npm</id>
                        <goals>
                            <goal>install-node-and-npm</goal>
                        </goals>
                        <phase>generate-test-resources</phase>
                        <configuration>
                            <nodeVersion>v24.11.0</nodeVersion>
                        </configuration>
                    </execution>
                    <execution>
                        <id>compile-dcl</id>
                        <goals>
                            <goal>npx</goal>
                        </goals>
                        <phase>generate-test-resources</phase>
                        <configuration>
                            <arguments>--package=@sap/ams-dev compile-dcl
                                -d ${project.basedir}/src/main/resources/dcl
                                -o ${project.basedir}/target/generated-test-sources/dcn
                            </arguments>
                        </configuration>
                    </execution>
                </executions>
            </plugin>
        </plugins>
    </build>
```
:::

#### Loading DCN
To load the compiled DCN files, the AMS client library needs to be configured to do so before tests.

::: tip
In CAP Node.js projects, this is done automatically by the AMS modules if `requires.auth.kind = mocked`.

In CAP Java projects, this is done automatically if `spring-boot-starter-cap-ams-test` is on the classpath and `cds.security.mock.enabled` is set to `true`.
:::

::: code-group
```js [Node.js]
// test-specific application setup

const ams = AuthorizationManagementService.fromLocalDcn("./test/dcn", {
    assignments: "./test/mockPolicyAssignments.json"
});
```

```java [Java]
// test-specific application setup

AuthorizationManagementService ams = null;
try {
    LocalAuthorizationManagementServiceConfig amsTestConfig = 
        new LocalAuthorizationManagementServiceConfig()
            .withPolicyAssignmentsPath(
                Path.of("src", "test", "resources", "mockPolicyAssignments.json"));
    ams = AuthorizationManagementServiceFactory
            .fromLocalDcn(
                Path.of("target", "generated-test-sources", "ams", "dcn"), amsTestConfig);

    ams.whenReady().get(3, TimeUnit.SECONDS);
} catch (TimeoutException e) {
    throw new RuntimeException("Local AMS test client did not become ready within timeout", e);
} catch (Exception e) {
    throw new RuntimeException("Failed to create AMS client", e);
}
```
:::

##  Test policies

The DCL package called `local` has a special semantic: it's intended for DCL files with policies that are only relevant for testing, but not for production. Policies in the `local` package are ignored during base policy upload even if they are present in the archive.

This allows you to test policies that are restrictions of base policies without shipping them to customers. Typically, such policies would be created by an administrator at runtime in the `SCI admin cockpit`.