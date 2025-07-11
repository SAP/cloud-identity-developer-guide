# @sap/ams

## CAP integration

AMS can be used for authorization in CAP applications to provide both role and instance-based authorization management at runtime.
The integration is based on the standard cds annotations for authorization via roles and optional ams-specific annotations for instance-based authorization filters. 

For production, AMS is meant to be used with SAP Cloud Identity Services as authentication solution but mocked authentication can be used to test authorization without the need for SAP Cloud Identity Services tokens. This is useful when the application is started locally or to execute automated tests.

When deployed, the application's authorization policies are managed in your application using the administration console of SAP Cloud Identity Services. During development, policies can be edited in the IDE and assigned to mocked users via the `cds env` configuration of non-production profiles.

The plugin runtime has the following expectations on the project environment. If your projects differs from this, for example due to a custom auth middleware, you can customize the plugin via [cds env configuration](#configuration) and the [plugin runtime configuration](#plugin-runtime).

| **Default**                                     | **Value**                                                                       | **Customize**
|------------------------------------------       | ------------------------------------------------------------------------------- | -----------------
| SAP Cloud Identity Services credentials location| `cds.env.requires.auth.credentials`                                             | [Provide credentials manually](#custom-sap-cloud-identity-services-credential-location)
| amsPluginRuntime.authProvider.xssecAuthProvider | Defaults to `IdentityServiceAuthProvider`                                       | [Replace default XssecAuthProvider](#custom-xssecauthprovider)
| @sap/xssec SecurityContext location             | `IdentityServiceSecurityContext` expected under `cds.context.http.req.authInfo` | [Replace default CdsAuthProvider](#custom-cdsauthprovider)
| amsPluginRuntime.authProvider                   | Defaults to `CdsXssecAuthProvider`                                              | [Replace default CdsAuthProvider](#custom-cdsauthprovider)

### cds add ams
The `cds add ams` command configures the application for AMS.

It installs both the AMS runtime plugin for cds (@sap/ams) and the AMS development plugin for cds (@sap/ams-dev):

```shell
npm i @sap/ams
npm i --save-dev @sap/ams-dev
```

Additionally, it configures the application's deployment artefacts (`mta`, `helm`, `cf-manifest`) for AMS, for example by adding configuration for the [ams policy deployer application](#base-policy-upload).

### Features
#### Role-based authorization
AMS can be used to assign roles from the cds authorization model to users via authorization policies. For example, the following policy would grant the `Admin` role when assigned to a user:

```SQL
POLICY Admin {
  ASSIGN ROLE Admin;
}
```

The AMS plugin implements a middleware that computes the roles of SAP Cloud Identity Services users before each request by overriding the [`user.is`](https://cap.cloud.sap/docs/node.js/authentication#user-is) function.

#### Instance-based authorization
Policies that assign roles can be extended with attribute filters for instance-based authorization. This allows administrators to create custom policies at runtime for fine-grained control. This is most useful to give customer administrators in multi-tenant applications fine-grained control over their user's rights.

Via `@ams.attributes` annotations, the AMS attributes are mapped to elements (or association paths) in the cds model. Whenever requests access cds resources with those elements, the result is filtered based on the attribute conditions computed by AMS. The conditions are based on the `WHERE` condition behind role assignments in policies.

`ams.attributes` annotations are supported on *aspects*, *entities* and *actions/functions bound to a single entity* as those are the [cds resources that support *where* conditions](https://cap.cloud.sap/docs/guides/security/authorization#supported-combinations-with-cds-resources).

Example annotations
```js
aspect media {
    price: Integer;
    genre: Association to Genres;
}

annotate media with @ams.attributes: {
    price: (price),
    genre: (genre.name)
};

@restrict: [{ grant:['READ'], to: ['Reader'] }]
entity Books : media {
  // ...
}
```

Example schema.dcl
```sql
SCHEMA {
  genre: String,
  price: Number
}
```

Example basePolicies.dcl
```sql
POLICY "Reader" {
    ASSIGN ROLE "Reader" WHERE genre IS NOT RESTRICTED AND price IS NOT RESTRICTED;
}
```

Example admin policy (created at runtime via the administration console of SAP Cloud Identity Services)
```sql
POLICY JuniorReader {
    USE "Reader" RESTRICT genre IN ('Fantasy', 'Fairy Tale'), price < 20;
}
```

#### Validation
The AMS plugin [@sap/ams](https://www.npmjs.com/package/@sap/ams) adds a [custom build task](https://cap.cloud.sap/docs/guides/deployment/custom-builds#custom-build-plugins) for *ams*.

It validates `@ams.attributes` annotations for syntactic correctness and type coherence during `cds build`, and whenever a model is loaded if the application was started via `cds serve`, `cds watch` or `cds.test`. This gives early feedback about the correctness of the annotations during development:

- validates that `@ams.attributes` annotations map AMS attributes syntactically correct to cds elements via expression syntax.
- if a manually written/adjusted `schema.dcl` is used, validates that all AMS attributes mapped via `@ams.attributes` annotations exist and have a type that fits each cds element to which they are mapped.
- if a generated `schema.dcl` is used, validates that the inferred type of each AMS attribute is coherent across all `@ams.attributes` mappings in which it is mapped to a cds element.

#### Base Policy Generation
Unless disabled, the AMS build task generates base policies for roles that occur in the `@requires` and `@restrict` annotations of the cds model.

Example annotation
```js
@restrict: [{ grant:['READ'], to: ['Reader', 'Inquisitor'] }]
entity Books as projection on my.Books { *,
```

Example `basePolicies.dcl`
```sql
POLICY "Reader" {
  ASSIGN ROLE "Reader";
}

POLICY "Inquisitor" {
  ASSIGN ROLE "Inquisitor";
}
```

It also generates a `schema.dcl` that defines AMS attributes with inferred types based on `@ams.attributes` annotation for [instance-based authorization](#instance-based-authorization).

:information_source: Policy generation also occurs during `cds serve`, `cds watch` and `cds.test` to react to changes of cds annotations.

DCL Files that have been modified manually will not be overridden during generation to allow manual changes of the schema and base policies. To force re-generation of a generated DCL file, delete it prior to the next DCL generation.

#### Base Policy Upload
During `cds build`, a policy deployer application will be generated in `<cds.build.target>/policies`:

[*Node.js default*] `gen/policies`\
[*Java default*] `srv/src/gen/policies`

It requires a certificate-based binding to the Identity service and must be deployed together with the application to upload the base policies to the AMS.
For `helm` chart deployments, it can be built via `containerize` as a Node.js image and deployed with the `content-deployment` helm template.

#### Testing policies
The DCL package called `local` has a special semantic. It is meant for DCL files with policies that are only relevant for testing but not for production. Its policies are ignored during the base policy upload, even if they are contained during the upload.

For example, you can create fictitious admin policies inside this package to test whether extensions of base policies work as expected.

`@sap/ams-dev` automatically compiles DCL files to the `DCN` format which is required for local policy evaluations. This happens when the application is started via `cds start`, `cds watch` or via `cds.test`, so that the application should be able to do authorization checks via AMS even during development without deploying the policies first to the SAP Cloud Identity Services.

#### Mocked user testing
For testing and development purposes, policies can be assigned to mocked users via the `cds env` configuration of non-production profiles:

```json
{
    "requires": {
        "auth": {
            "[development]": {
                "kind": "mocked",
                "users": {
                    "carol": {
                        "policies": [
                            "cap.Reader"
                        ]
                    },
                    "dave": {
                        "policies": [
                            "cap.Admin",
                            "cap.Reader"
                        ]
```

It is important to assign policies via their fully-qualified name which includes the DCL package (`cap` in this example).

Of course, you can still assign roles via the `roles` array directly to mocked users.
Assigning policies instead of roles is mostly useful for testing instance-based authorization via AMS as the attribute filters only apply to roles assigned via AMS policies.

#### Hybrid testing
If [autoDeployDcl](#configuration) is enabled when bound to an `ias` instance for authentication, for example during [Hybrid testing](https://cap.cloud.sap/docs/advanced/hybrid-testing), the AMS plugin uploads the base policies to the AMS server instead of compiling them to `DCN`. From there, they will be downloaded into the DCN engine shortly after that via polling and subsequently used for authorization checks.

**Be very careful though with `autoDeployDcl` and do not enable it when bound against a productive system or it will override the deployed base policies with the current development state!**

An application bound to an `ias` instance for authentication will always download its policy bundle from the corresponding AMS instance. This means, hybrid testing can be used to run an application locally with the policies from an AMS instance (including admin policies created at runtime) without overriding them. Downloading policies in hybrid mode does not require `autoDeployDcl` to be enabled.

### Configuration
The AMS plugins are configured inside the `requires.auth.ams` property of the [cds env](https://cap.cloud.sap/docs/node.js/cds-env#project-settings).\
It supports the following properties with the following [`default`]:

- **generateDcl** *true/false* [`true`]: unless set to `false`, generates `basePolicies.dcl` and `schema.dcl` from the cds model (see [Base Policy Generaiton](#base-policy-generation))
- **dclRoot** *string* [`ams/dcl` / `srv/src/main/resources/ams` (Java)]: the root DCL folder (containing the `schema.dcl`) which is used for generating DCL, compiling DCL to DCN during development, uploading DCL etc.
- **dclGenerationPackage** *string* [`cap`]: name of the DCL package to which basePolicies.dcl is generated
- **dcnRoot** *string* [`gen/dcn` / `srv/src/gen/ams/dcn` (Java)]:  folder for DCL to DCN compilation results during development (see [Testing Policies](#testing-policies))
- **policyDeployerRoot** *string* [`gen/policies` / `srv/src/gen/policies` (Java)]:  folder of the ams policy deployer application created during `cds build` (see [Base Policy Upload](#base-policy-upload))
- **authPushDcl** *true/false* [`false`]:  if enabled, uploads the base policies to the AMS server (see [Hybrid testing](#hybrid-testing)

All AMS properties also work lowercased (for example `generatedcl`) and this casing has priority of the camelCase (for example `generateDcl`) version of properties. This means, all [cds env sources](https://cap.cloud.sap/docs/node.js/cds-env#sources-for-cds-env) including case-insensitive ones are supported such as setting properties via environment variables (`CDS_REQUIRES_AUTH_AMS_GENERATEDCL`) which gets mapped to lowercased versions of the property. 

#### Plugin Runtime
It is possible to replace the following defaults in the runtime of the plugin to configure it for non-standard project environments.

##### Custom SAP Cloud Identity Services credential location

If the SAP Cloud Identity Services credentials are not available under the default location (`cds.env.requires.auth.credentials`), you need to manually provide them:

server.js
```js
const { amsCapPluginRuntime } = require("@sap/ams");

amsCapPluginRuntime.credentials = { ... } // manually provide the SAP Cloud Identity Services credentials from service binding
```

##### Custom XssecAuthProvider

It is possible to override the `XssecAuthProvider` implementation used by the default `CdsAuthProvider` internally to a different implementation.
For example, the following snippet shows how it can be replaced in projects that authenticate both via SAP Cloud Identity Services and XSUAA.

server.js
```js
const { amsCapPluginRuntime, HybridAuthProvider } = require("@sap/ams");

const mapScope = (scope, securityContext) => scope; // your custom scope to policy mapper
amsCapPluginRuntime.authProvider.xssecAuthProvider = new HybridAuthProvider(amsCapPluginRuntime.ams, mapScope) // authorization for both SAP Cloud Identity Services and XSUAA tokens
```

##### Custom CdsAuthProvider

server.js
```js
const { amsCapPluginRuntime } = require("@sap/ams");

amsCapPluginRuntime.authProvider = new MyCustomCdsAuthProvider(); // your custom CdsAuthProvider implementation if you do not authorize based on @sap/xssec SecurityContexts
```

#### Technical communication
By default, the plugin runtime uses an [IdentityServiceAuthProvider](#identityserviceauthprovider) which supports technical communication via SAP Cloud Identity Services out-of-the-box.
You can access it as follows in the default plugin runtime to configure which policies to use for technical communication scenarios:

```js
const { amsCapPluginRuntime, TECHNICAL_USER_FLOW, PRINCIPAL_PROPAGATION_FLOW } = require("@sap/ams");
const { mapTechnicalUserApi, mapPrincipalPropagationApi } = require('../ams/apis.js'); // your custom API mappers

const identityServiceAuthProvider = amsCapPluginRuntime.authProvider.xssecAuthProvider;
identityServiceAuthProvider
  .withApiMapper(mapTechnicalUserApi, TECHNICAL_USER_FLOW)
  .withApiMapper(mapPrincipalPropagationApi, PRINCIPAL_PROPAGATION_FLOW);
```

### Logging
The AMS CAP plugins log to namespace `ams` in CAP. To see [debug logs](https://cap.cloud.sap/docs/node.js/cds-log#debug-env-variable) during development, turn it on for this namespace, for example via

```shell
DEBUG=ams cds watch
```

## deploy-dcl script

The script deploys a DCL bundle (including schema.dcl, DCL root package and all sub-packages) to the Identity service instance from the environment (see `deploy-dcl --help`):

```
Usage: deploy-dcl -d [DCL_ROOT_DIR] -c [CREDENTIALS_FILE] -n [DEPLOYER_APP_NAME]

Options:
      --help         Show help                                         [boolean]
      --version      Show version number                               [boolean]
  -d, --dcl          [optional] path to the directory that contains the DCL root
                     package. If a path is provided via environment variable
                     AMS_DCL_ROOT, it overrides this option.
                                                       [string] [default: "dcl"]
  -c, --credentials  [optional] path to a JSON file containing the credentials
                     object of an Identity service binding. If omitted, will try
                     to find and use an Identity service binding from the
                     process environment.                               [string]
  -n, --name         [optional] a descriptive name of this deployer application
                     to trace back the currently deployed DCL bundle on the AMS
                     server to its source when DCL is deployed from more than
                     one source. If a name is provided directly via environment
                     variable AMS_APP_NAME or indirectly as application_name via
                     VCAP_APPLICATION on Cloud Foundry or the pod name on K8s,
                     it overrides this option.
                                   [string] [default: "@sap/ams:deploy-dcl"]

Examples:
  deploy-dcl                                Pushes the DCL content in ./dcl
                                            (including schema.dcl, DCL root
                                            package and all subpackages) to the
                                            identity service instance from the
                                            environment.
  deploy-dcl -d src/dcl -c config/ias.json  Pushes the DCL content from
  -n bookshop-dcl-deployer                  ./src/dcl using the SAP Cloud Identity Services
                                            credentials in ./config/ias.json.
                                            The deployer app name for this
                                            upload will be set to
                                            "bookshop-dcl-deployer" to be able
                                            to trace back the upload source to
                                            this deployer.
```