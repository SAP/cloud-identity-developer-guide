# Getting Started

This document provides the basic information required to use Authorization Management Service (**AMS**) for authorization checks in your application. It provides an overview of the available library modules, their features, and how to integrate them into your projects.


## Supported Languages and Frameworks

The client libraries of AMS consist of different modules for the following programming languages and frameworks:

- **Java** (Maven):
    - [jakarta-ams](/Libraries/java/jakarta/jakarta-ams.md)
    - [spring-ams](/Libraries/java/spring/spring-ams.md)
    - [cap-ams](/Libraries/java/cap/cap-ams.md)
- **JavaScript** (Node.js):
    - [@sap/ams](/Libraries/nodejs/sap_ams/sap_ams.md)
    - [@sap/ams-dev](https://www.npmjs.com/package/@sap/ams-dev)
- **Go**:
    - [cloud-identity-authorizations-golang-library](/Libraries/go/go-ams)

The next section lists the required module dependencies for different application setups, depending on the programming language and framework you are using.

## Dependency Setup

::: tip
In CAP applications, the [`cds add ams`](https://cap.cloud.sap/docs/tools/cds-cli#cds-add) command can be executed with the *latest version* of [`@sap/cds-dk`](https://cap.cloud.sap/docs/tools/cds-cli#cli). It automatically adds the correct dependencies.
:::

The following tables give an overview of the required AMS module dependencies for different application setups.

::: warning
The recommended modules and versions have changed over time (see [Historical Setups](#historical-setups))

**Please begin new projects with the currently recommended modules**.
:::

### Java

##### Runtime Dependencies

| Project Type                | jakarta-ams | spring-ams |spring-boot-starter-ams|   cap-ams  | spring-boot-starter-cap-ams |
|-----------------------------|:-----------:|:----------:|:---------------------:|:----------:|:---------------------------:|
| Jakarta EE                  |     ✓       |     -      |            -          |     -      |              -              |
| Spring Boot                 |     *       |     *      |            ✓          |     -      |              -              |
| Jakarta EE (CAP)            |     *       |     -      |            -          |     ✓      |              -              |
| Spring Boot (CAP)           |     *       |     -      |            -          |     *      |              ✓              |

\* transitive dependency

##### Test-Scoped Dependencies

| Project Type                | spring-boot-starter-ams-test | spring-boot-starter-cap-ams-test|
|-----------------------------|:----------------------------:|:-------------------------------:|
| Spring Boot                 |                ✓             |                 -               |     
| Spring Boot (CAP)           |                -             |                 ✓               |

##### Tooling Dependencies

::: tip
In CAP Java projects, the (optional) Node.js module `@sap/ams` *should* be added in the `package.json` as a *devDependency* with version `^3` to provide dev-time features as [cds build plugin](/CAP/cds-Plugin).
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

If you operate productive applications with a dependency setup different from the recommended one, you can usually continue using the modules you already have installed for some time. However, we recommend migrating to the new modules and major versions eventually in discussion with us.

### Java Library Version 3
For major version 3 of the Java libraries, the following dependency setup was recommended:

##### Runtime Dependencies

| Project Type                | jakarta-ams | spring-ams | cap-ams-support  |
|-----------------------------|:-----------:|:----------:|:----------------:|
| Jakarta EE                  |     ✓       |     -      |         -        |
| Spring Boot                 |     *       |     ✓      |         -        |
| Jakarta EE (CAP)            |     *       |     -      |         ✓        |
| Spring Boot (CAP)           |     *       |     -      |         ✓        |

\* transitive dependency

##### Test-Scoped Dependencies

| Project Type                |       jakarta-ams-test       | spring-ams-test-starter |
|-----------------------------|:----------------------------:|:-----------------------:|
| Jakarta EE                  |                ✓             |             -           |
| Spring Boot                 |                -             |             ✓           |
| Jakarta EE (CAP)            |                -             |             -           |
| Spring Boot (CAP)           |                -             |             -           |

##### Tooling Dependencies

::: tip
In CAP Java projects, the (optional) Node.js module `@sap/ams` *should* be added in the `package.json` as a *devDependency* with version `^3` to provide dev-time features as [cds build plugin](/CAP/cds-Plugin).
:::