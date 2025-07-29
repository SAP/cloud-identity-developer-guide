# Value Help

The IAS Administration Console provides a value help feature when creating restrictions of base policies. For example, given the following base policy:

```dcl
POLICY ReadProducts {
    GRANT read ON products WHERE category IS NOT RESTRICTED;
}
```

An administrator may use the IAS Administration Console to create a policy that is a restriction of the base policy for a specific category, such as "electronics":

```dcl
POLICY ReadElectronics {
    USE ReadProducts RESTRICT category = 'electronics';
}
```

The value help feature allows the administrator to select the "electronics" category (or any other category that exists in the application) from a list of available categories. This way, the administrator does not have to guess valid values for the `category` attribute, but can select valid values from a list.

## Value Help Requests

To retrieve the list of available values for a specific attribute, the AMS server sends a value help request to the application. In the response, the AMS server receives a list of valid values for the requested attribute from the application, which it can then present to the administrator in the IAS Administration Console.

<!-- TODO: insert link to AMS Value Help documentation on help.sap.com once available -->
This documentation describes how to annotate the DCL schema to enable value help for attributes as well as the required response format which applications need to implement in the value help request handlers.

## Authorizing Value Help Requests

The value help endpoints in your application MUST be protected because they return business data. To allow the application to authorize value help requests, the AMS server calls the application with an [App-To-App](/Authorization/TechnicalCommunication#app-to-app) principal propagation token based on the administrator that is requesting value help in the IAS Administration Console.

### API Permission Group

The API permission group consumed by the AMS server can be freely chosen in the service configuration of the AMS instance.
Make sure to setup a policy for this API permission group that grants the necessary privileges for the value help endpoints as described in the [App-To-App](/Authorization/TechnicalCommunication#app-to-app) documentation.
Note that the administrator using the IAS Administration Console must also have the necessary privileges to access the value help endpoints in your application.

::: tip
Depending on your application, it may not be necessary to create a dedicated role or *action*/*resource* for the value help endpoints.
For example, if there is already a policy that grants read access to categories, you can re-use this policy in the *internal* policy to which the value help API permission group is mapped.
:::

### Certificate Validation
The value help request from the AMS server uses TLS with a certificate that must be used to validate ownership of the token.

::: tip
The official BTP security libraries for authentication provide the necessary proof-of-ownership validation under the name `proof token validation`.
:::

::: warning Platform-specific certificate handling
As the ingress of the cloud platform terminates TLS, the certificate of the caller needs to be forwarded to your application, by default in the `x-forwarded-client-cert` header.

If your application is deployed on Cloud Foundry, the `.cert` domain must be used for the value help callback URL.
Cloud Foundry accepts client certificates only on this domain. In this case, it automatically fill the `x-forwarded-client-cert` header that is used during the validation.

If your application is deployed on Kyma, `Istio` needs to be configured to forward the TLS certificate via the `x-forwarded-client-cert` header to your application.
:::