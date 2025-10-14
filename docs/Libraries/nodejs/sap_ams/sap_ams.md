# @sap/ams

`@sap/ams` is the Node.Js runtime library used to perform [authorization checks](https://help.sap.com/docs/identity-authentication/identity-authentication/configuring-authorization-policies?locale=en-US) in applications based on AMS policies.

Additionally, it is a cds plugin that provides [AMS CAP dev-time features](/CAP/cds-Plugin) for both Node.js and Java CAP applications.

The module [@sap/ams-dev](https://www.npmjs.com/package/@sap/ams-dev) provides the corresponding tooling support for local [testing](/Authorization/Testing) until a Node.js DCL compiler is available.

## Installation
The module is available in the public [npmjs](https://www.npmjs.com/) repository:

```sh
npm install @sap/ams
```

::: danger Important
Always keep the version of this dependency up-to-date since it's a **crucial part of your application's security**:

```bash
npm update @sap/ams # or: npm update
```
This is especially important when you deploy your application with a `package-lock.json` that locks the version that gets installed to a fixed version, until the next `npm update`.
:::

:::tip
To check which version of the module is installed, run:

```bash
npm list @sap/ams
```
This prints a dependency tree that shows which versions of the module are installed, including transitive dependencies.
:::

## Version 3
Version 3 drastically changes the core API. Instead of checking privileges on a `PolicyDecisionPoint` with an `Attributes` object, an `AuthProvider` prepares an `Authorizations` object for the same purpose. This separates *what* to check from *how* to check it. The necessary configuration for advanced authorization scenarios, such as principal propagation or non-standard authorization strategies, are configured once during application start. As a result, the authorization checks themselves remain straightforward in version 3, with a focus on the application domain.

New features:

- Out-of-the-box support for technical communication scenarios using SAP Cloud Identity Services
- Flexible configuration and extensibility for non-standard authorization strategies, for example when authenticating both via XSUAA and SAP Cloud Identity Services tokens
- Exports Typescript Types for a better development experience
- Improved events that allow correlating authorization checks with requests for logging and auditing
- Support for SAP Cloud Identity Services credentials with certificates changing at runtime, for example when using ZTIS or mounted Kyma service bindings

#### Version 3 Migration
For Non-CAP Node.js applications, please open a [Support](/Support) ticket if you need help with the migration of a productive application.

CAP Node.js applications should **not** need to make changes when updating to version 3 

## Usage

The following snippets show you how to use the core API of the library. For more details on the concepts, see [Authorization Checks](/Authorization/AuthorizationChecks).

::: tip CAP applications
For CAP Node.js applications, most of this section is irrelevant because the library setup and authorization checks are fully automated.
:::

### Set-up
The following code snippets show the setup of the `@sap/ams` library in a Node.js application. We recommend to carve out the Authorization Management Service (**AMS**) setup into a separate module, for example, to `ams.js`, and import the AMS-related objects in your `server.js`.

```js
const { AuthorizationManagementService, IdentityServiceAuthProvider, AmsError } = require("@sap/ams");

let ams;
if (process.env.NODE_ENV === 'test') {
  // For local tests, initialize AMS with a locally compiled policy bundle and mocked policy assignments
  ams = AuthorizationManagementService.fromLocalDcn("./test/dcn", { // your compile target directory of the @sap/ams-dev compile-dcl script
    assignments: "./test/mockPolicyAssignments.json" // a file path of your choice or an in-memory structure
  });
} else {
  // For production, initialize AMS with the cloud policy and assignment bundle from the AMS server
  const identityService = ... // your @sap/xssec 4 IdentityService instance used for authentication
  ams = AuthorizationManagementService.fromIdentityService(identityService);
}

// Use the standard IdentityServiceAuthProvider or if necessary, write your own AuthProvider implementation instead
const authProvider = new IdentityServiceAuthProvider(ams);

// Recommendation: Register a middleware behind the authentication middleware that allows authorization checks via the req object
const amsMw = authProvider.getMiddleware(); // IdentityServiceAuthProvider provides a pre-configured instance of AmsMiddleware
app.use(/^\/(?!health).*/i, authenticate, amsMw.authorize());
```

::: warning Important
Implement a [startup check](/Authorization/StartupCheck) to ensure that the `AuthorizationManagementService` instance is ready for authorization checks before serving authorized endpoints.
:::

### Authorization checks
Authorization checks are performed on `Authorization` objects created by implementations of the `AuthProvider` interface.

`Authorizations` can be created manually for a specific (sub)set of policies:
```js
const authorizations = ams.getAuthorizations({
  policies: ["dclPackage.dclSubpackage.policyName1", "dclPackage.dclSubpackage.policyName2"],
  includeDefaultPolicies: false
});
```

Typically, however, authorizations are derived from the security context that is created from a valid token during authentication. For example, `IdentityServiceAuthProvider` creates `Authorizations` from an `IdentityServiceSecurityContext`.

```js
const securityContext = await createSecurityContext(identityService, { req }); // @sap/xssec 4 authentication
const authorizations = await authProvider.getAuthorizations(securityContext);
```

If the `AmsMiddleware#authorize()` middleware is registered, it streamlines this process for all requests by creating an `Authorizations` instance and placing it on the `req` object under the `AMS_AUTHORIZATIONS` symbol exported by `@sap/ams`.

```js
const { AMS_AUTHORIZATIONS } = require("@sap/ams");

function createOrder(req, res) {
  const authorizations = req[AMS_AUTHORIZATIONS];
}
```

A typical authorization check to verify if a user is allowed to create orders in the `EU` region looks like this:

```js
const input = { "region" : "EU" };
const decision = authorizations.checkPrivilege('create', 'orders', input);
``` 

### Handling Decisions
An authorization check returns a `Decision` object. A decision can be in one of three states:

- `granted`: the checked privilege is unconditionally granted
- `denied`: the checked privilege is unconditionally denied
- `conditional`: the checked privilege might be granted depending on a condition whose attributes have not been grounded to values in the attribute input.

It provides the following methods to distinguish between these cases:

```js
let decision;

// context-free privilege-check
decision = authorizations.checkPrivilege('delete', 'orders')
if(!decision.isGranted()) {
    return res.sendStatus(403);
}  

// contextual privilege check for a single entity
decision = authorizations.checkPrivilege('create', 'orders', { "$app.product.category" : "accessory" });
if(!decision.isGranted()) {
  return res.sendStatus(403);
}

// contextual privilege check for many entities
decision = authorizations.checkPrivilege('read', 'orders');
if(decision.isDenied()) {
  return res.sendStatus(403);
} else if (decision.isGranted()) {
  // definitive GRANT without outstanding WHERE condition
  return res.json(db.readAllOrders());
} else {
  // instance-based GRANT with outstanding WHERE condition
  const filter = decision.visitDCN(convertCall, convertValue); // convert condition to a filter for the DB layer
  return res.json(db.readOrders(filter));
}
```

### AmsMiddleware
Besides [`authorize()`](#authorization-checks), the `AmsMiddleware` provides additional handlers to define [declarative](/Authorization/AuthorizationChecks#declarative-authorization-checks) privilege (pre-)checks on the endpoint layer:

```js
// returns 403 when no definitive (= no outstanding WHERE condition) GRANT delete ON orders is assigned to user
app.delete('/orders/:id', amsMw.checkPrivilege('delete', 'orders'), deleteOrder);

// returns 403 when no GRANT read ON orders is assigned to user. A potential WHERE condition for the GRANT is acceptable. It must be evaluated in the service handler.
app.get('/orders', amsMw.precheckPrivilege('read', 'orders'), getOrders);

// returns 403 when no definitive GRANT read ON products or no (definitive or conditional) GRANT create ON orders is assigned to user
app.post('/orders', amsMw.checkPrivilege('read', 'products'), amsMw.precheckPrivilege('create', 'orders'), createOrder);
```

### Events/Logging
[Debug Logs](/Authorization/Logging#debug-logging) of `@sap/ams` can be enabled by including `ams` in the `DEBUG` environment variable:

```sh
DEBUG=ams node server.js
DEBUG=xssec,ams node server.js # also log @sap/xssec 4 debug logs...
```

Consumer applications can listen to events of the `AuthorizationManagementService` instance to manually log authorization check results and/or create audit log events:

```js
ams.on("authorizationCheck", event => {
  if(event.type === "checkPrivilege") {
    if (event.decision.isGranted()) {
      console.log(`Privilege '${event.action} ${event.resource}' for ${event.context.token.scimId} was granted based on input`, event.input);
    }
  } else if (event.type === "getPotentialPrivileges") {
    ...
  }
});
```

### Error Handling

If necessary, you can identify errors from @sap/ams via instanceof, for example, in the express error handler:

```js
app.use((err, req, res, next) => {
    if(err instanceof AmsError) {
      // Error from @sap/ams library
    } else {
      // other Error
    }

    return res.sendStatus(500);
});
```

#### Bundle Loading Errors

AMS uses a bundle loader internally to manage the policies and assignments bundle in the background, independently of incoming requests. Instances of `AuthorizationManagementService` emit two distinct error event types when bundle loading requests fail:

- **`bundleInitializationError`**: Emitted when the initial bundle download fails and the instance is not yet ready for use.
- **`bundleRefreshError`**: Emitted when a bundle refresh request fails (includes time since last successful refresh). Since the library continuously polls, this doesn't necessarily mean the data is outdated, just that the polling attempt failed. The instance remains ready but if there have been recent policy or assignment changes, it cannot take them into account.

You can distinguish between these event types using the `type` property and handle them according to your requirements:

```js
ams.on("error", event => {
  if (event.type === "bundleInitializationError") {
    console.error("AMS bundle initialization failed - service not ready:", event.error);
    // Eventually the separate startup check calling the whenReady function will reject,
    // so typically no action besides logging is required here
  } else if (event.type === "bundleRefreshError") {
    console.warn(`AMS bundle refresh failed (current bundle age: ${event.secondsSinceLastRefresh} seconds):`,event.error);
    // Consider taking action such as logging an error instead of a warning when the bundle is stale for
    // extended periods of time
  }
});
```

::: info Automatic Error Logging
If your application does not subscribe to the "error" event, these bundle events will be automatically logged to the console. This is due to Node.js's special handling of events with the name "error".
:::

::: tip Handling Initial Bundle Load Errors
Refer to the [Startup Check](/Authorization/StartupCheck) documentation for guidance on how to react when AMS fails to initialize the bundle. The error events emitted for this case are only intended to provide information about the failed requests.
:::

### Testing
See the central [Testing](/Authorization/Testing) documentation for details.

## Configuration

### Memory Consumption
The memory that `@sap/ams` needs depends on the number of tenants, users and policy assignments in the application. 
To approximately calculate the memory usage you can use the following formula: 
````
Memory(MB) = 6.54 + (AssignmentCount × 0.000117)
````
Although the memory usage depends on tenants, users and policy assignemnts, we found out that the driving factor behind large bundle sizes is primarily the number of policy assignments which naturally increases with a larger number of tenants and users. Our experiments found the above formula is a simple and practical way to estimate bundle sizes.

In the following table you can find some example sizes: 
|  Users  | Tenants | Assignments | Total Memory | Memory Growth | Memory/User | Memory/Assignment |
|:-------:|:-------:|:-----------:|:------------:|:-------------:|:-----------:|:-----------------:|
| 0       | 0       | 0           | 6.54MB       | 0MB           | -           | -                 |
| 10      | 1       | 19          | 6.67MB       | 0.13MB        | 13.0KB      | 6.8KB             |
| 100     | 1       | 167         | 6.72MB       | 0.18MB        | 1.8KB       | 1.1KB             |
| 1,000   | 10      | 1,901       | 6.97MB       | 0.43MB        | 0.43KB      | 0.23KB            |
| 10,000  | 100     | 19,164      | 9.06MB       | 2.52MB        | 0.25KB      | 0.13KB            |
| 50,000  | 100     | 95,867      | 18.28MB      | 11.74MB       | 0.24KB      | 0.12KB            |
| 100,000 | 100     | 191,446     | 29.12MB      | 22.58MB       | 0.23KB      | 0.12KB            |

The analysis revealed that memory usage per user and per policy assignment scales linearly. However, the distribution of users across tenants significantly impacts memory consumption:
- 1 tenant: 7.4KB per user
- 10 tenants: 0.43KB per user
- 100 tenants: 0.24KB per user

With an increasing number of tenants, the memory usage per user gets more efficient, however the overall memory consumption will increase as more tenants are added. 

## API

### AuthorizationManagementService

#### Construction
- **`fromIdentityService(identityService, config?): AuthorizationManagementService`**  
  Creates an instance using the DCN and policy assignments fetched with SAP Cloud Identity Services credentials.  
  - `identityService` (object): SAP Cloud Identity Services object with **certificate-based** credentials.  

- **`fromLocalDcn(dcnRoot, config?): AuthorizationManagementService`**  
  Creates an instance using locally compiled DCL files for testing.  
  - `dcnRoot` (string): Root directory of the DCN bundle.  
  - `config` (object, optional):  
    - `watch` (boolean, default: `false`): Watch for file changes.  
    - `assignments` (string | PolicyAssignments, optional): Path to JSON file or `PolicyAssignments` object.  
    - `debounceDelay` (number, default: `1000`): Debounce delay in ms for changes of local DCL files.  
    - `start` (boolean, default: `true`): Control whether to immediately start downloading the AMS bundle.

If an instance has been constructed with `config.start=false`, the loading of the AMS bundle must be started manually. This is useful when ZTIS is used, and the credentials do not yet contain a certificate when the instance is created:

```js
const ams = AuthorizationManagementService.fromIdentityService(identityService, { start: false });
// fill credentials with certificate asynchronously from ZTIS
getCertificateFromZTIS().then((cert, key) => {
  identityService.setCertificateAndKey(cert, key);
  ams.start();
});
```

#### Readiness Checks
- **`whenReady(timeoutSeconds = 0): Promise<void>`**  
  Returns a Promise that resolves once the instance is ready for authorization checks. If it hasn't received policies and assignments after the specified timeout interval, the Promise is rejected.  
  - `timeoutSeconds` (number): Maximum waiting time in seconds.  

- **`isReady(): boolean`**  
  Synchronously checks if the instance is ready for authorization checks.  

---

### Authorizations

An abstract representation of authorizations determined by the strategy of the [AuthProvider](#authproviderc) from which it was constructed.

#### Methods

- **`constructor(ams: AuthorizationManagementService, policySet: PolicySet, context: any): Authorizations`**  

- **`checkPrivilege(action: string, resource: string, input?: AttributeInput): Decision`**  
  Checks if the action is allowed on the resource.  
  - `input` (AttributeInput, optional): A flat input object that grounds attribute names to values, for example, `{ "product.category" : "accessory" }`.
  Attributes that are not grounded in the input are interpreted as *unknowns* and may result in a conditional decision.

- **`getPotentialResources(): Set<string>`**  
  Collects all resources for which at least one action is potentially granted, ignoring conditions.  

- **`getPotentialActions(resource: string): Set<string>`**  
  Collects all actions that are potentially granted for a given resource, ignoring conditions.  

- **`getPotentialPrivileges(): Array<{action: string, resource: string}>`**  
  Collects all action/resource combinations that are potentially granted, ignoring conditions.  

- **`withDefaultInput(input: AttributeInput): Authorizations`**  
  Sets default input used for all authorization checks.  
  - `input` (AttributeInput, optional): A flat input object that grounds fully-qualified attribute names to values, for example, `{ "$env.$user.origin" : "EU" }`

- **`limitedTo(other: Authorizations): Authorizations`**  
  Limits the authorizations of this instance to the authorizations of another instance. Subsequent authorization checks on this instance use the logical intersection of its authorizations and those of the other Authorization instances.

---

### Decision

Represents the result of an authorization check. A decision can have one of three states: *granted*, *denied*, or *conditional*.

#### Methods

- **`isGranted(): boolean`**  
  Returns true if the authorization check resulted in a definitive GRANT with no outstanding conditions.

- **`isDenied(): boolean`**  
  Returns true if the authorization check resulted in a definitive DENY with no outstanding conditions.

- **`<T,V>visit(visitCall: CallVisitor, visitValue: ValueVisitor) : T`**
  This method can be used to visit the condition tree bottom-up. The visitor calls `visitValue` whenever it encounters a value (attribute reference or literal) or `visitCall` when it encounters a function call in the condition, for example, a call to the "EQ" function to compare two arguments for equality.
    - `visitCall` ((call : string, args : (DcnReference|DcnValue|V)[]) => T): A function that visits the given call and its arguments, for example, to transform `("EQ", args)` => `"args[0] = args[1]"`. The call names are the constants from `DclConstants.operators`.
    - `transformValue` ((value : DcnReference|DcnValue) => DcnReference|DcnValue|V): A function that visits the given attribute reference or literal, for example, to translate AMS references to database field names.
    - {ref:string} `DcnReference`
    - {number|string|boolean|number[]|string[]|boolean[]} `DcnValue`

- **`filterUnknown(unknowns: string[]): Decision`**  
  Returns a new `Decision` instance that is the result of keeping only the fully-qualified attributes as *unknown*, evaluating the remaining attributes as *unset*.

- **`apply(flatInput: { [attributeName: string]: DcnValue }): Decision`**  
  **EXPERIMENTAL**: If you plan to use this method in production, please open a ServiceNow ticket on component `BC-CP-CF-SEC-LIB`.
  Uses the data provided in `flatInput` to create a new Decision. Simple example on a decision `d` that represents `a = 3` and `b = 4`:
  ```js
  d.apply({ a: 3 }) // returns a Decision that represents b=4
  d.apply({ a: 3, b: 4 }) // returns a Decision that represents true
  d.apply({ a: 1 }) // returns a Decision that represents false
  d.apply({ a:3 }).apply({ b: 4 }) // returns a Decision that represents true
  ```
   
### Events

Instances of `AuthorizationManagementService` emit the following events to which consumers can subscribe using the `on(eventName: string, function(event: AmsEvent) : void)` method.

- **`authorizationCheck`**: Emitted during various authorization operations by the methods of [Authorizations](#authorizations). The event object contains the following properties based on the type of operation:
  - **`type`**: The type of the event. Possible values are:
    - `"checkPrivilege"`: Emitted during a privilege check. Additional payload:
      - **`action`**: The action being checked.
      - **`resource`**: The resource being checked.
      - **`input`**: The input used for the authorization check.
      - **`decision`**: The decision of the authorization check.
    - `"getPotentialActions"`: Emitted when collecting potential actions for a resource. Additional payload:
      - **`resource`**: The resource for which actions are being collected.
      - **`potentialActions`**: The set of actions that are potentially granted for the given resource.
    - `"getPotentialResources"`: Emitted when collecting potential resources. Additional payload:
      - **`potentialResources`**: The set of resources for which at least one action is potentially granted.
    - `"getPotentialPrivileges"`: Emitted when collecting potential privileges. Additional payload:
      - **`potentialPrivileges`**: The list of potentially granted privileges, each containing:
        - **`action`**: The action.
        - **`resource`**: The resource.
  - **`authorizations`**: The `Authorizations` instance that triggered the event.
  - **`context`**: The context of the event from which the `Authorizations` instance was created.

- **`error`**: Emitted when an error occurs in a background operation. The event object contains the following properties:
  - **`type`**: The type of the event. Possible values are:
    - `"bundleInitializationError"`: Emitted when the initial bundle download fails and the instance is not yet ready.
    - `"bundleRefreshError"`: Emitted when the bundle loader fails to refresh the current policies and assignments bundle, for example, due to a failed request to the AMS server. Since the library continuously polls, this doesn't necessarily mean the data is outdated. Additional payload:
      - **`secondsSinceLastRefresh`**: Time in seconds since the last successful bundle refresh.
  - **`error`**: The `AmsError` instance that describes the refresh error.



## CAP Runtime Integration

This section provides details for the customization of the `@sap/ams` **runtime** functionality for authorization checks in CAP Node.js applications.\
The `@sap/ams` **cds build plugin** functionality is documented on the [cds plugin](/CAP/cds-Plugin#Configuration) page.

### Runtime Features

1. `@sap/ams` implements a middleware that injects additional roles to [cds.context.user](https://cap.cloud.sap/docs/node.js/authentication#cds-user) based on policies before each request by overriding the [`user.is`](https://cap.cloud.sap/docs/node.js/authentication#user-is) function.
2. Additionally, it integrates into [instance-based authorization](https://cap.cloud.sap/docs/guides/security/authorization#instance-based-auth) of the cds framework by dynamically adjusting the `where` condition for the [privileges](https://cap.cloud.sap/docs/guides/security/authorization#restrict-annotation) used in a request if a user's role assignments are restricted based on DCL conditions.

### Runtime Customization
It is possible to replace the following defaults in the runtime of the plugin to configure it for non-standard project environments.

#### Custom SAP Cloud Identity Services credential location

If the SAP Cloud Identity Services credentials aren't available under the default location (`cds.env.requires.auth.credentials`), you must provide provide them manually:

::: code-group

```js [server.js]
const { amsCapPluginRuntime } = require("@sap/ams");

amsCapPluginRuntime.credentials = { ... } // manually provide the SAP Cloud Identity Services credentials from service binding
```
:::

#### Custom XssecAuthProvider

It's possible to override the `XssecAuthProvider` implementation used inside the default `CdsAuthProvider` with a different implementation.

For example, the following snippet shows how it can be replaced in projects where the authentication middleware provides either an `IdentityServiceSecurityContext` or an `XsuaaSecurityContext` and both shall be authorized based on policies.

::: code-group
```js [server.js]
const { amsCapPluginRuntime, CdsXssecAuthProvider, HybridAuthProvider } = require("@sap/ams");

const mapScope = (scope, securityContext) => scope; // your custom scope to policy mapper
amsCapPluginRuntime.authProvider.xssecAuthProvider = new HybridAuthProvider(amsCapPluginRuntime.ams, mapScope) // authorization for both SAP Cloud Identity Services and XSUAA tokens
```

```ts [server.ts]
const { amsCapPluginRuntime, HybridAuthProvider } = require("@sap/ams");

const mapScope = (scope, securityContext) => scope; // your custom scope to policy mapper
(amsCapPluginRuntime.authProvider as CdsXssecAuthProvider).xssecAuthProvider = new HybridAuthProvider(amsCapPluginRuntime.ams, mapScope) // authorization for both SAP Cloud Identity Services and XSUAA tokens
```
:::

##### Custom CdsAuthProvider

If your CAP project doesn't use the `@sap/xssec` library for authentication, you can provide a custom implementation of the `CdsAuthProvider` interface that derives the user's `Authorizations` with a custom strategy from the `cds.context` object:

::: code-group
```js [server.js]
const { amsCapPluginRuntime } = require("@sap/ams");

amsCapPluginRuntime.authProvider = new MyCustomCdsAuthProvider();
```
:::

### Debug Logging
`@sap/ams` logs to namespace `ams` in CAP. 

To see CDS debug logs, [turn them on](https://cap.cloud.sap/docs/node.js/cds-log#debug-env-variable) for this namespace, for example, using

```shell
DEBUG=ams cds watch
```
