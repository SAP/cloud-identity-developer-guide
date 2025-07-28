# Getting Started

This document provides the basic information required to use AMS for authorization checks in your application. It provides an overview of the available library modules, their features, and how to integrate them into your projects.



## Supported Languages and Frameworks

The AMS client libraries consist of different modules for the following programming languages and frameworks:

- **Java** (Maven):
    - [jakarta-ams](/java/jakarta-ams/jakarta-ams.md)
    - [spring-ams](/java/spring-ams/spring-ams.md)
    - [cap-ams-support](/java/cap-ams-support/cap-ams-support.md)
- **JavaScript** (Node.js):
    - [@sap/ams](/nodejs/sap_ams/sap_ams.md)
    - [@sap/ams-dev](https://www.npmjs.com/package/@sap/ams-dev)
- **Go**:
    - [cloud-identity-authorizations-golang-library](/go/go-ams.md)

The next section lists the required module dependencies for different application setups, depending on the programming language and framework you are using.

## Dependency Setup

::: tip
In CAP applications, the [`cds add ams`](https://cap.cloud.sap/docs/tools/cds-cli#cds-add) command can be executed with the *latest version* of [`@sap/cds-dk`](https://cap.cloud.sap/docs/tools/cds-cli#cli) to add the correct dependencies automatically.
:::

The following tables give an overview of the required AMS module dependencies for different application setups.

::: warning
The recommended modules and versions have changed over time (see [Historical Setups](#historical-setups))

**Please begin new projects with the currently recommended modules**.
:::

### Java

| Project Type                | jakarta-ams | spring-ams | jakarta-ams-test | cap-ams-support | @sap/ams    |
|-----------------------------|:-----------:|:----------:|:----------------:|:---------------:|:-----------:|
| Jakarta EE                  |     ✓       |     -      |        [✓]      |        -        |      -
| Spring Boot                 |     -       |     ✓      |        [✓]      |        -        |      -
| CAP                         |     -       |     -      |         -       |        ✓        |     (✓)\*  

[✓] = *test-scoped **(!)** maven dependency*\
(✓) = *devDependency in package.json*

::: tip *
The (optional) Node.js module `@sap/ams` *can* be added in the `package.json` as a *devDependency* with version `^3` to provide dev-time features as [cds build plugin](/CAP/cds-Plugin).

:::

### Node.js

| Project Type        | @sap/ams | @sap/ams-dev   | Java JDK |
|---------------------|:--------:|:--------------:|:---------------:|
| Plain Node.js       |   ✓ ^3   |      (✓)* ^2   |    (✓)* 17+
| express (Node.js)   |   ✓ ^3   |      (✓)* ^2   |    (✓)* 17+
| CAP (Node.js)       |   ✓ ^3   |      (✓)* ^2   |    (✓)* 17+

(✓) = *devDependency*

::: tip *
only required to compile DCL files before running local tests. We are currently finishing a compiler in Javascript that will make these dependencies obsolete.
:::

### Go

| Project Type | cloud-identity-authorizations-golang-library |
|--------------|:--------------------------------------------:|
| Go           |                    ✓                         |

## Samples
For practical examples of how to set up and use the AMS client libraries, refer to the [Samples](/Samples) section. It contains sample applications demonstrating the necessary setup for authorization with AMS in various programming languages and frameworks.

## Historical Setups

If you operate productive applications with a dependency setup different from the recommended one, you can usually continue using the modules you already have installed for some time, but we recommend migrating to the new modules and major versions eventually in discussion with us.

### JDK < 17
For Java versions < 17, the modules `java-ams` and `java-ams-test` are a drop-in replacement for `jakarta-ams` and `jakarta-ams-test`.