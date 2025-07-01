[![REUSE status](https://api.reuse.software/badge/github.com/SAP/cloud-identity-authorizations-libraries)](https://api.reuse.software/info/github.com/SAP/cloud-identity-authorizations-libraries)

# AMS Client Libraries documentation

## About this project
This repository contains documentation for the usage of the client libraries of the Authorization Management Service (**AMS**) which is part of [SAP Cloud Identity Services](https://help.sap.com/docs/cloud-identity-services?locale=en-US) (**SCI**). The libraries provide APIs for policy-based authorization checks in applications of different programming languages and frameworks.

## Documentation
The documentation for the client libraries is available in the [docs](/docs) directory.

### Concepts
- [**Getting Started**](/docs/GettingStarted.md): The basics to find the right client library modules for your application and jump-start development.
- [**Authorization Checks**](/docs/AuthorizationChecks.md): How to perform authorization checks in your application.
- [**Testing**](/docs/Testing.md): How to test authorization checks in your application.
- [**Technical communication**](/docs/TechnicalCommunication.md): How to authorize requests from other applications.
- [**Deploying DCL Policies**](/docs/DeployDCL.md): How to deploy DCL policies to the AMS server.
- [**ValueHelp**](/docs/ValueHelp.md): How to provide value help for attribute restrictions in the SCI cockpit.
- [**Support**](/docs/Support.md): Information on how to report bugs, receive consultation, provide feedback, and raise feature requests.

### Modules

- **Java (Maven)**:
    - [jakarta-ams](/docs/java/jakarta-ams/jakarta-ams.md)
    - [spring-ams](/docs/java/spring-ams/spring-ams.md)
    - [cap-ams-support](/docs/java/cap-ams-support/cap-ams-support.md) (replaces [~~`cap-support`~~](/docs/java/cap-support/cap-support.md))
- **Javascript (Node.js)**:
    - [@sap/ams](/docs/nodejs/sap_ams/sap_ams.md)
    - [@sap/ams-dev](/docs/nodejs/sap_ams-dev/sap_ams-dev.md)
- **Go**:
    - [cloud-identity-authorizations-golang-library](/docs/go/go-ams/go-ams.md)

## Support, Feedback, Contributing

This project is open to feature requests/suggestions, bug reports etc. via [GitHub issues](https://github.com/SAP/cloud-identity-authorizations-libraries/issues). Contribution and feedback are encouraged and always welcome. For more information about how to contribute, the project structure, as well as additional contribution information, see our [Contribution Guidelines](CONTRIBUTING.md).

## Security / Disclosure
If you find any bug that may be a security problem, please follow our instructions at [in our security policy](https://github.com/SAP/cloud-identity-authorizations-libraries/security/policy) on how to report it. Please do not create GitHub issues for security-related doubts or problems.

## Code of Conduct

We as members, contributors, and leaders pledge to make participation in our community a harassment-free experience for everyone. By participating in this project, you agree to abide by its [Code of Conduct](https://github.com/SAP/.github/blob/main/CODE_OF_CONDUCT.md) at all times.

## Licensing

Copyright 2025 SAP SE or an SAP affiliate company and cloud-identity-authorizations-libraries contributors. Please see our [LICENSE](LICENSE) for copyright and license information. Detailed information including third-party components and their licensing/copyright information is available [via the REUSE tool](https://api.reuse.software/info/github.com/SAP/cloud-identity-authorizations-libraries).
