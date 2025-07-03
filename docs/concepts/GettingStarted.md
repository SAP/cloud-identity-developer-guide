# Getting Started

This document provides the basic information required to use AMS for authorization checks in your application. It provides an overview of the available library modules, their features, and how to integrate them into your projects.



## Supported Languages and Frameworks

The AMS client libraries consist of different modules for the following programming languages and frameworks:

- **Java** (Maven):
    - [jakarta-ams](/java/jakarta-ams/jakarta-ams.md)
    - [spring-ams](/java/spring-ams/spring-ams.md)
    - [cap-ams-support](/java/cap-ams-support/cap-ams-support.md) (replaces [~~`cap-support`~~](/java/cap-support/cap-support.md))
- **JavaScript** (Node.js):
    - [@sap/ams](/nodejs/sap_ams/sap_ams.md)
    - [@sap/ams-dev](/nodejs/sap_ams-dev/sap_ams-dev.md)
- **Go**:
    - [cloud-identity-authorizations-golang-library](/go/go-ams.md)

The next section lists the required module dependencies for different application setups, depending on the programming language and framework you are using.

## Dependency Setup

::: tip
In CAP applications, the [`cds add ams`](https://cap.cloud.sap/docs/tools/cds-cli#cds-add) command can be executed with the *latest version* of [`@sap/cds-dk`](https://cap.cloud.sap/docs/tools/cds-cli#cli) to add the correct dependencies automatically.
:::

The following tables give an overview of the required AMS module dependencies for different application setups.

::: warning
The recommended modules and versions have changed over time (see [Historic Setups](#historic-setups))

**Please begin new projects with the currently recommended modules**.
:::

**Legend**: ✔️ runtime dependency (✔️) development dependency

### Java

| Project Type                | jakarta-ams | spring-ams | cap-ams-support | @sap/ams    |
|-----------------------------|:-----------:|:----------:|:---------------:|:-----------:|
| Jakarta EE                  |     ✔️      |     -      |        -        |      -
| Spring Boot                 |     -       |     ✔️     |        -        |      -
| CAP (Spring Boot)           |     ✔️\*    |     -\*    |       ✔️        |     (✔️)\*\*

\* Yes, `jakarta-ams` is required but `spring-ams` should not be installed even when the CAP application uses Spring.\
\*\* The (optional) Node.js module `@sap/ams` *can* be added in the `package.json` as a *devDependency* with version `^3` to provide the following dev-time features as [cds build plugin](https://cap.cloud.sap/docs/guides/deployment/custom-builds#custom-build-plugins):
- Generation of DCL during `cds build`
- Generation of DCL policy deployer application during `cds build`
- Validation of `@ams.attributes` annotations against `schema.dcl` during `cds build`

### Node.js

| Project Type        | @sap/ams | @sap/ams-dev   | Java JDK |
|---------------------|:--------:|:--------------:|:----------:|
| Plain Node.js       |   ✔️ ^3  |      (✔️)* ^2    |    (✔️)* 17+
| express (Node.js)   |   ✔️ ^3  |      (✔️)* ^2    |    (✔️)* 17+
| CAP (Node.js)       |   ✔️ ^3  |      (✔️)* ^2    |    (✔️)* 17+

\* only required to compile DCL files before running local tests. We are currently finishing a compiler in Javascript that will make these dependencies obsolete.

### Go

| Project Type | cloud-identity-authorizations-golang-library |
|--------------|:-------------------------------------------:|
| Go           |                    ✔️                        |

## Samples
For practical examples of how to set up and use the AMS client libraries, refer to the [Samples](/Samples) section. It contains sample applications demonstrating the necessary setup for authorization with AMS in various programming languages and frameworks.

## Historical Setups

If you operate productive applications with a dependency setup different from the recommended one, you can usually continue using the modules you already have installed for some time, but we recommend migrating to the new modules eventually in discussion with us.