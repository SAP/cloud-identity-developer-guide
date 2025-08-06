# Support

## Provided Services
The support channels mentioned below are primarily intended for **bugs** of the provided libraries and the service itself.
However, you can also use the support channels to raise feature requests and give feedback.

::: warning Out of scope
Please understand that, due to the high volume of requests, we **can't** offer individual consulting via the support channels unless it's part of a customer's service plan.

For the same reason, we can't handle tickets in the spirit of *"Here is my code, why does it not work?"*.
:::


## Support Channels

::: danger Important
Before opening support tickets, please try to solve the problem by following the steps in the [Troubleshooting](./Troubleshooting.md) section.

Make sure to provide the [required information](#required-information) when creating a support ticket.
:::

As registered SAP customers, report your issue by creating a ticket on the [SAP Support Portal](https://support.sap.com/en/index.html) for one of the following components:

| Name                                        | Service Now Component | Description                                                                                              |
| --------------------------------------------|-----------------------|----------------------------------------------------------------------------------------------------------|
| Identity Authentication Service (**IAS**)   | BC-IAM-IDS            | SCI tenant, User Management, Token Flows, Technical Communication etc.                                   |
| Authorization Management Service (**AMS**)  | BC-IAM-CAS            | AMS Server, Authorization Management UI, DCL compiler, buildpack, bundle gateway                         |
| Security Client Libraries                   | BC-CP-CF-SEC-LIB      | Java, Node.js and Go client libraries for SAP Cloud Identity Services, Authorization Management, XSUAA   |
| BTP `identity` Service Broker               | BC-IAM-IB             | Service creation of SAP Cloud Identity Services instances on BTP                                         | 

::: warning Related Components
The following components are not supported by the SCI team but related. Please report issues to these components directly if they are not related to SCI:
:::

| Name                                  | Service Now Component | Description                                       |
|---------------------------------------|-----------------------|---------------------------------------------------|
| Approuter                             | BC-XS-APR             | @sap/approuter, managed application router etc.   |
| **XSUAA**                             | BC-CP-CF-SEC-IAM      | Extended Services User Account and Authentication |   


See also [Getting Support](https://help.sap.com/docs/btp/sap-business-technology-platform/btp-getting-support) in the SAP BTP documentation.

## Required Information

Make sure that you include the following information in the support ticket to give us a chance to understand and resolve your issue:

- Installed **modules** and **versions** of the client libraries for SAP Cloud Identity Services *or* full dependency tree (for example,, `mvn dependency:tree` for Java, `npm ls` for Node.js, or `go list -m all` for Go)
- For exceptions: Stack trace that includes the executed **code locations of the SCI client libraries** that lead to the exception
- For unexpected authentication results (401 or Error): the relevant **debug log output** of the token validation in the libraries
- For unexpected authorization results (403 or too many privileges): the relevant **debug log output** of the privilege check in the libraries for Authorization Management that you don't understand
- Potential steps you have already tried to fix the problem
- Reason why you believe it's a bug