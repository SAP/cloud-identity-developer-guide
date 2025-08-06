# Value Help

The administration console of SAP Cloud Identity Services provides a value help feature when an administrator creates restrictions of base policies. For example, using the following base policy:

```dcl
POLICY ReadProducts {
    GRANT read ON products WHERE category IS NOT RESTRICTED;
}
```

An administrator can use the administration console to create an authorization policy that is a restriction of the base policy for a specific category, such as "electronics":

```dcl
POLICY ReadElectronics {
    USE ReadProducts RESTRICT category = 'electronics';
}
```

The value help feature allows the administrator to select the "electronics" category (or any other category that exists in the application) from a list of available categories. This way, the administrator doesn't have to guess valid values for the `category` attribute, but can select valid values from a list.

## Value Help Requests

To retrieve the list of available values for a specific attribute, the Authorization Management Service (**AMS**) sends a value help request to the application. In the response, AMS receives a list of valid values for the requested attribute from the application, which it can then present to the administrator in the administration console.

## Authorizing Value Help Requests

The value help endpoints in your application MUST be protected because they return business data. To allow the application to authorize value help requests, the AMS server calls the application with an [App-To-App](/Authorization/TechnicalCommunication#app-to-app) principal propagation token based on the administrator who requests value help in the administration console.

### API Permission Group

The API permission group consumed by the AMS server can be freely chosen in the service configuration of the AMS instance.
Make sure to setup a policy for this API permission group that grants the necessary privileges for the value help endpoints as described in the [App-To-App](/Authorization/TechnicalCommunication#app-to-app) documentation.
Note that the administrator using the administration console must also have the necessary privileges to access the value help endpoints in your application.

::: tip
Depending on your application, it may not be necessary to create a dedicated role or *action*/*resource* for the value help endpoints.
For example, if there is already a policy that grants read access to categories, you can reuse this policy in the *internal* policy to which the value help API permission group is mapped.
:::

### Certificate Validation
The value help request from the AMS server uses TLS with a certificate that must be used to validate ownership of the token.

::: tip
The official BTP security libraries for authentication provide the necessary proof-of-ownership validation under the name `proof token validation`.
:::

::: warning Platform-specific certificate handling
As the ingress of the cloud platform terminates TLS, the certificate of the caller needs to be forwarded to your application, by default in the `x-forwarded-client-cert` header.

If your application is deployed on Cloud Foundry, the `.cert` domain must be used for the value help callback URL.
Cloud Foundry accepts client certificates only on this domain. In this case, it automatically fills the `x-forwarded-client-cert` header that is used during the validation.

If your application is deployed on Kyma, `Istio` must be configured to forward the TLS certificate to your application using the `x-forwarded-client-cert` header.
:::