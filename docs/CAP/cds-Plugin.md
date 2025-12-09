# The cds Plugin `@sap/ams`

The Authorization Management Service (**AMS**) Nodejs module `@sap/ams` functions as a plugin for the [cds](https://cap.cloud.sap/docs/tools/cds-cli) CLI tool.
It adds a [custom build task](https://cap.cloud.sap/docs/guides/deployment/custom-builds#custom-build-plugins) for *ams* that automatically runs during `cds build` to provide the following features.

::: tip
In Node.js CAP projects, the task occurs also before/during `cds serve`, `cds watch` and `cds.test`.
:::

## Features

### DCL Generation
Unless disabled, the AMS build task generates DCL files from the cds model.


::: tip
DCL files that have been modified manually aren't overridden during generation. This allows manual changes of the schema and base policies. To force a repeated generation of a generated DCL file, delete it prior to the next DCL generation.
:::

##### Base Policy Generation
The *ams* build task generates base policies for roles that occur in the `@requires` and `@restrict` annotations of the cds model:

::: code-group
```cds [BookService.cds]
@restrict: [{ grant:['READ'], to: ['Reader', 'Inquisitor'] }]
entity Books as projection on my.Books { *,
```

```dcl [basePolicies.dcl]
POLICY "Reader" {
  ASSIGN ROLE "Reader";
}

POLICY "Inquisitor" {
  ASSIGN ROLE "Inquisitor";
}
```
:::

##### Schema Generation
It also generates a `schema.dcl` that defines AMS attributes with inferred types based on `@ams.attributes` annotations for [instance-based authorization](/CAP/InstanceBasedAuthorization):

::: code-group
```cds [SalesOrderService.cds]
annotate SalesOrder with @ams.attributes: {
    Region: (region),
    Budget: (total)
};
```

```dcl [schema.dcl]
SCHEMA {
  Region : String,
  Budget : Number
}
```
:::

### Base Policy Upload
Unless disabled, a [policy deployer application](/Authorization/DeployDCL#ams-policies-deployer-app) is generated in:

- [**Node.js**] `<cds.build.target>/policies` which defaults to `gen/policies`
- [**Java**] `srv/src/gen/policies`

During `cds add ams`, deployment descriptors like `mta.yaml` or `helm` charts are automatically configured to deploy policies to AMS with the policy deployer application in the default location.

### Validation
It validates `@ams.attributes` annotations for syntactic correctness and type coherence. This gives early feedback about the correctness of the annotations during development:

- validates that `@ams.attributes` annotations map AMS attributes syntactically correct to cds elements via cds expressions.
- if a generated `schema.dcl` is used, validates that the inferred type of each AMS attribute is coherent across all `@ams.attributes` mappings in which it's mapped to a cds element.
- if a manually written/adjusted `schema.dcl` is used, validates that all AMS attributes mapped using `@ams.attributes` annotations exist and have a type that fits each cds element to which they are mapped.

## Configuration
The cds plugin for AMS is configured inside the `requires.auth.ams` property of the [cds env](https://cap.cloud.sap/docs/node.js/cds-env#project-settings).\
It supports the following properties with the following [`default`]:

- **generateDcl** *true/false* [`true`]: unless set to `false`, generates `basePolicies.dcl` and `schema.dcl` from the cds model (see [Base Policy Generation](#base-policy-generation))
- **dclRoot** *string* [`ams/dcl` / `srv/src/gen/ams` (Java)]: the root DCL folder (containing the `schema.dcl`) which is used for generating DCL, compiling DCL to DCN during development, uploading DCL etc.
- **dclGenerationPackage** *string* [`cap`]: name of the DCL package to which basePolicies.dcl is generated
- **dcnRoot** *string* [`gen/dcn` / `srv/src/gen/ams/dcn` (Java)]:  folder for DCL to DCN compilation results during development (see [Testing](/Authorization/Testing#compiling-dcl-to-dcn))
- **generatePoliciesDeployer** *"auto"/false* [`"auto"`]: unless set to `false`, generates a policy deployer application during `cds build` (see [Base Policy Upload](#base-policy-upload))
- **policyDeployerRoot** *string* [`gen/policies` / `srv/src/gen/policies` (Java)]: folder of the AMS policy deployer application created during `cds build` (see [Base Policy Upload](#base-policy-upload))

### Node.js specific configuration
- **autoDeployDcl** *true/false* [`false`]:  if enabled, uploads the base policies to the AMS server on application start and after DCL changes (see [Hybrid Testing](https://cap.cloud.sap/docs/advanced/hybrid-testing)).

:::tip
All *requires.auth.ams* properties also work in lowercase (for example `generatedcl`), and lowercase has priority over the camel case version (for example `generateDcl`) of properties. This means that all [cds env sources](https://cap.cloud.sap/docs/node.js/cds-env#sources-for-cds-env) including the case-insensitive ones are supported, such as setting properties using environment variables (`CDS_REQUIRES_AUTH_AMS_GENERATEDCL`), which are mapped to lowercase versions of the property. 
:::