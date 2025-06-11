# Getting Started

This document helps you get started with the correct client library setup in your application. It provides an overview of the available libraries, their features, and how to integrate them into your projects.



## Table of Contents

- [Supported Languages](#supported-languages)
- [Project Dependencies](#project-dependencies)
    - [Java](#java)
    - [Node.js](#nodejs)
    - [Go](#go)
    - [Historic Project Setups](#historic-project-setups)
- [Deploying DCL Policies](#deploying-dcl-policies)



## Supported Languages
The AMS client libraries are available for the following programming languages and frameworks:

- **Java**:
    - [jakarta-ams](/docs/libraries/java/jakarta-ams/jakarta-ams.md)
    - [spring-ams](/docs/libraries/java/spring-ams/spring-ams.md)
    - [cap-ams-support](/docs/libraries/java/cap-ams-support/cap-ams-support.md)
    - ~~`cap-support`~~ (superseded by `cap-ams-support`)
- **JavaScript**:
    - [@sap/ams](/docs/libraries/nodejs/ams/ams.md)
    - [@sap/ams-dev](/docs/libraries/nodejs/ams-dev/ams-dev.md)
- **Go**:
    - [cloud-identity-authorizations-golang-library](/docs/libraries/go/go-ams/go-ams.md)



## Project Dependencies
The following tables give an overview of the required AMS module dependencies in your project setup.

In CAP applications, the [`cds add ams`](https://cap.cloud.sap/docs/tools/cds-cli#cds-add) command can be executed with the *latest version* of [`@sap/cds-dk`](https://cap.cloud.sap/docs/tools/cds-cli#cli) to add the correct dependencies automatically.

The recommended modules and versions have changed over time (see [History](#history)), so please begin new projects with the currently recommended modules. If you have existing projects, you can usually continue using the modules you already have installed for some time, but we recommend migrating to the new modules eventually in discussion with us.

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
| Plain Node.js       |   ✔️ ^3  |      (✔️)* ^2    |    (✔️)** 17+
| express (Node.js)   |   ✔️ ^3  |      (✔️)* ^2    |    (✔️)** 17+
| CAP (Node.js)       |   ✔️ ^3  |      (✔️)* ^2    |    (✔️)** 17+

\* / \*\*  required to compile DCL files before running local tests. We are currently finishing a compiler in Javascript that will make these dependencies obsolete.

### Go

| Project Type | cloud-identity-authorizations-golang-library |
|--------------|:-------------------------------------------:|
| Go           |                    ✔️                        |

### Historic Project Setups



## Deploying DCL Policies
Learn how to deploy DCL policies together with your application in the dedicated [DCL Policy Deployment Guide](/docs/DeployDCL.md).