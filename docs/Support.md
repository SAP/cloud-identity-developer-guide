# Support

## Provided Services
The support channels mentioned below are primarily intended for bugs of the provided libraries and the service itself.
But you can also use the support channels to raise feature requests and give feedback.

::: warning Out of scope
Please understand that, due to the high volume of requests, individual consulting can **not** be offered via the support channels unless it is part of a customer's service plan.

For the same reason, we cannot handle tickets whose essence comes down to

"Here is my code, why does it not work?".
:::


## Support Channels

::: danger Important
Before opening support tickets, please try to solve the problem by following the steps in the [Troubleshooting](./Troubleshooting.md) section.

Make sure to provide the [required information](#required-information) when creating a support ticket.
:::

As registered SAP customers, report your issue by creating a ticket on the [SAP Support Portal](https://support.sap.com/en/index.html) for this component:

**BC-CP-CF-SEC-LIB** (Problem in AMS client libraries)
**BC-IAM-CAS** (Problem in AMS service)

See also [Getting Support](https://help.sap.com/docs/btp/sap-business-technology-platform/btp-getting-support) in the SAP BTP documentation.

## Required Information

Make sure to include the following information in the support ticket to give us a chance to understand and resolve your issue:

- Installed **modules** and **versions** of the AMS libraries *or* full dependency tree (e.g., `mvn dependency:tree` for Java, `npm ls` for Node.js, or `go list -m all` for Go)
- For exceptions: Stack trace that includes the executed **code locations of the AMS libraries** that lead to the exception
- For unexpected authorization results (403 or too many privileges): the relevant **debug log output** of the privilege check in the AMS libraries that you do not understand
- Potential steps you have already tried to fix the problem
- Reason why you believe it is a bug