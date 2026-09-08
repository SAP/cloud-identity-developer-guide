# System-to-System Communication

System-to-System communication describes a programmatic data exchange between two applications or services, without a user sitting in front of a browser driving the request.

The Authorization Management Service (**AMS**) supports authorization of system-to-system communication for both *technical users* (systems acting on their own behalf) and *principal propagation* (systems acting on behalf of a user whose request is forwarded).

The main difference between the communication patterns is how trust between the parties is established and how they authenticate. Understanding the authentication side is a prerequisite for the authorization side described here. See [Consume APIs from Other Applications](https://help.sap.com/docs/cloud-identity-services/cloud-identity-services/consume-apis-from-other-applications) for the authentication details.

## Authorization Approaches

When you develop an application that needs to authorize requests coming from another system, the first question is **what** should be authorized. The answer depends entirely on the use case, and the two approaches are:

- **Technical Communication** (also referred to as *technical access*): the calling application itself is the principal. The caller is the entity that is authenticated and authorized, and no user is involved in the decision. Different authorization levels are typically mapped to different usage scenarios of the callee.
- **Principal Propagation**: the calling application acts on behalf of a user. The application is authenticated, but the user is the principal, so the effective authorizations depend on both the caller's and the user's access levels. This is the right approach when the data in the receiving application is owned by the user, or when actions should be performed on the user's behalf.

An application can support both approaches at the same time, and it decides per request which one applies based on the incoming token.

## Internal Policies

Regardless of which approach you choose, AMS lets you define special `INTERNAL` policies for authorizing the calling application. Internal policies have the following characteristics:

- They're **not assignable to users** and don't appear in the administration console, because they aren't relevant for tenant administrators.
- The same DCL schema applies as for all other policies. We recommend to define them in a dedicated package, typically named `internal`.
- The mapping from an internal policy to the caller's usage scenario is considered static and is therefore implemented in the application itself, rather than configured by an administrator.
- Every application is in charge of enforcing its own internal policies, as described on the pattern pages below.

For the `INTERNAL` keyword itself, see the [Data Control Language (DCL) reference](https://help.sap.com/docs/cloud-identity-services/cloud-identity-services/data-control-language-dcl).

## Effective Authorizations

In technical communication, the policies derived from the caller are the only input to the authorization decision, so they determine the caller's privileges directly.

In principal propagation, the policies assigned to the user form the basis of the decision, and the policies derived from the caller act as an **upper limit** on what may be granted during that request. The effective privileges are the intersection of both. This lets an application restrict what an external caller can do on behalf of a user, while still respecting the user's own assignments.

The `SciAuthorizationsProvider` implements this combination out of the box. See [Authorization Checks](/Authorization/AuthorizationChecks#sciauthorizationsprovider) for the exact rules, including what happens when only one of the two layers is present.

## Communication Patterns

### App-to-App

Communication between two applications that are registered in the same SAP Cloud Identity Services tenant, based on JWT tokens. The caller's authorizations are derived from the **API permission groups** it consumes, which are carried in the `ias_apis` claim of the token.

This pattern doesn't require the participating applications to follow the SAP BTP tenancy model, so it also works between BTP and non-BTP applications.

See [App-to-App](/Authorization/App2App) for the full documentation.
