# Authorizations in the Token

With AMS you define fine-grained, attribute-based permissions (**ABAC**)
as [authorization policies](/Authorization/AuthorizationBundle#authorization-policies). These policies are normally
distributed to your applications over a back-channel in the form of
an [Authorization Bundle](/Authorization/AuthorizationBundle) and evaluated locally by
the [AMS client libraries](/Authorization/AuthorizationChecks).

In some scenarios it is helpful to *additionally* expose some authorizations directly in the
JWT token, for example to drive **visibility conditions in a UI** or to let a gateway-like
component (such as SAP Build Work Zone) know the permissions of a user across downstream
applications.

This page explains what can be put into the token, the important restrictions that apply, and
the concrete use cases and configuration steps.

## What Ends Up in the Token

::: info The token is size-limited
A JWT is transported in HTTP headers and must stay within a size limit, so the token never
contains a full policy definition, only the *functional authorization* of a policy. See
[Token Size and the Fallback to Introspection](#token-size-and-the-fallback-to-introspection)
for details.
:::

The *functional authorization* is the functional part of an AMS policy:

- for `GRANT <action> ON <resource>` it is the `action` and `resource`
- for the CAP-specific [`ASSIGN ROLE <role>`](/CAP/Basics#assign-role-keyword) it is the `role`

**Conditions (`WHERE` statements) are never written into the token.**

In the token, functional authorizations are represented as a simple list. Action and resource
are concatenated with an `@` character; CAP roles are listed by their name:

```json
{
  "functionalAuthz": ["read@salesOrders", "write@salesOrders"]
}
```

```json
{
  "functionalAuthz": ["SalesManager"]
}
```

::: info The claim name is up to you
`functionalAuthz` is only an example. The actual claim name is the name you give to the
[assertion attribute](#step-2-add-the-assertion-attribute) that carries the values.
:::

::: warning Only the access token carries the functional authorizations
The claim is added to the **access token** only; the **ID token never contains it**. If
functional authorizations are configured, an access token is always issued alongside the ID token,
so make sure your UI reads the claim from the access token.
:::

## Functional Authorization vs. Full Policy

| | Authorization Bundle | Functional Authorization in Token |
|---|---|---|
| Content | Full policies incl. conditions | Only `action@resource` / `role` |
| Conditions (`WHERE`) | ✅ evaluated | ❌ not included |
| Grants from `DEFAULT` policies | ✅ included | ❌ not included |
| Distribution | Back-channel, polled by the client library | Front-channel, inside the JWT |
| Evaluated by | AMS client library in the backend | Anyone who can read the token (e.g. a UI) |
| Typical use | Authoritative authorization checks | UI visibility, pre-checks, shared/gateway scenarios |

::: danger Only use it for full checks if there are no conditions
The functional authorization in the token may be used for a **full authorization check only if
the AMS instance does not use any conditions** (no `WHERE` statements in any policy).

As soon as conditions are involved, the claim may only be used as a **visibility condition** or
as a **pre-check** (for example in a middleware). The authoritative check, including
condition evaluation, must still be performed in the backend with the AMS client library based
on the [Authorization Bundle](/Authorization/AuthorizationBundle). See also
[Conditional Policies](/Authorization/AuthorizationChecks#conditional-policies) and
[Querying Potential Privileges](/Authorization/AuthorizationChecks#querying-potential-privileges).
:::

## Enabling Functional Authorizations

Getting functional authorizations into the token takes two steps:

1. Mark which authorizations are relevant for the token.
2. Add an assertion attribute to the IAS application configuration so that the values are
   written into the issued tokens.

### Step 1: Mark the Relevant Authorizations

Nothing is written into the token by default. You mark the actions that are relevant with the
`@jwtRelevant` annotation in your `schema.dcl`, as described in
[Use Case 2](#use-case-2-ui-visibility-for-your-own-application). 

Instances that run **without** an authorization bundle are the exception. There the token is the
only carrier of the authorizations, so every action is marked automatically, see
[Use Case 1](#use-case-1-simple-permission-model-no-abac).

### Step 2: Add the Assertion Attribute

The values are only written into the token if the IAS application declares a corresponding
**assertion attribute**. An assertion attribute maps a **user attribute** (the source of the
values) to an **assertion attribute name** (the claim name in the token):

| User attribute | Contains | Use case |
|---|---|---|
| `applicationFunctionalAuthz` | Functional authorizations **of the user** (token subject) in **this** application | [Simple model](#use-case-1-simple-permission-model-no-abac), [UI visibility](#use-case-2-ui-visibility-for-your-own-application) |
| `dependencyFunctionalAuthz` | Functional authorizations **of the user** in **dependent** applications, keyed by their global tenant ID | [Shared authorization](#use-case-3-shared-authorization-work-zone-joule-gateways) |

For example, to expose the user's functional authorizations for your own application under a
`functionalAuthz` claim, add an assertion attribute with the user attribute
`applicationFunctionalAuthz` and the assertion attribute name `functionalAuthz`:

```json5
{
  "assertion-attributes": {
    "functionalAuthz": "applicationFunctionalAuthz" // whereas functionalAuthz is the application-chosen claim name as it appears in the token
  }
}
```

The example above uses the Identity Broker service instance parameter syntax. See the
[Identity Broker reference](https://help.sap.com/docs/cloud-identity-services/cloud-identity-services/reference-information-for-identity-service-of-sap-btp)
for the full configuration.

## Token Size and the Fallback to Introspection

Even though only the functional authorization is written into the token, a user with many
permissions could still exceed the token size limit. To prevent HTTP requests from being
rejected because of an oversized header, the following safeguards apply:

- A functional-authorization claim is only written into the token if it has **fewer than 100
  entries**. For claims that list dependent applications, the same limit applies both to the number
  of applications and to the list of each individual application.
- If the limit is reached, the claim is **not** added to the token. Instead, the token contains a
  `sap_introspect` claim listing the names of the omitted claims.
- The client can then fetch the full values online from the **token introspection endpoint**
  (`/oauth2/introspect`) or the **userinfo endpoint** (`/oauth2/userinfo`). Both return the omitted
  values for the access token you pass in.

`sap_introspect` is an array of the omitted claim names, so a client can tell from the token alone
whether it has to fetch anything online, and for which claims:

```json
{
  "sap_introspect": ["functionalAuthz"]
}
```

`sap_introspect` is a protected claim name, so it cannot be used for an assertion attribute of your
own.

::: tip Keep the token small
To reduce the risk that claims have to be fetched online, expose only the functional
authorizations you actually need in the token. For the UI-visibility use case this is done by
[annotating relevant actions](#use-case-2-ui-visibility-for-your-own-application) with
`@jwtRelevant`.
:::

## Use Cases

### Use Case 1: Simple Permission Model (No ABAC)

If your application has a simple permission model and does **not** use any conditions
(no `WHERE` statements in any policy), you can rely entirely on the token and do not need the
authorization bundle or the client library at all.

Disable the authorization bundle for the AMS instance, in the same `authorization` block that
carries `enabled`
(see [Provisioning of AMS instances](/Authorization/GettingStarted#provisioning-of-ams-instances)):

```yml [mta.yaml]
resources:
  - name: my-app-ias
    type: org.cloudfoundry.managed-service
    parameters:
      service: identity
      service-plan: application
      config:
        authorization:
          enabled: true
          disable-authz-bundle: true
```

::: info Configuration style depends on your provisioning system
The example above shows the service configuration of the SAP Cloud Identity Services service
broker. The concrete syntax can differ depending on how your application is provisioned: if you
configure the SAP Cloud Identity Services application directly (for example through the
Application API), the same setting is named `disableAuthzBundle`
and is nested inside an `amsConfiguration` object. Please check the documentation of the
provisioning system you use.
:::

With this configuration:

- **All** `action`/`resource` tuples of all policies are automatically marked as relevant for
  the token, so no annotation is required. Disabling the bundle is what switches this on: without
  a bundle the token is the only carrier of the authorizations. Instances that keep their bundle
  expose only explicitly annotated actions
  (see [use case 2](#use-case-2-ui-visibility-for-your-own-application)).
- **No authorization bundle** is produced for the AMS instance, i.e. the AMS client library
  cannot be used for this instance.
- The IAS application configures the `applicationFunctionalAuthz` assertion attribute.

Authorization checks are then performed by a simple **string comparison** of the token claim
against the expected `action@resource` tuple:

```text
Token claim:            functionalAuthz: ["read@salesOrders", "write@salesOrders"]
Expected authorization: "read@salesOrders"
Result:                 authorized ✅
```

::: danger Reminder
This full-check-via-token approach is only valid because there are no conditions. Do not use it
if any policy of the instance contains a `WHERE` condition.
:::

::: warning Handle a missing claim explicitly
Because the token is the only source of authorizations in this use case, you must handle the case
that the claim is [omitted due to its size](#token-size-and-the-fallback-to-introspection). If the
claim name is listed in `sap_introspect`, fetch the values online before deciding. Never treat a
missing claim as "no permissions" (which locks the user out) or as "all permissions" (which grants
too much).
:::

::: details Complete access token example
```json
{
  "iss": "https://mytenant.accounts.ondemand.com",
  "sub": "8a2c4f10-6b3d-4e21-9f7a-1c2d3e4f5a6b",
  "aud": "my-app-client-id",
  "azp": "my-app-client-id",
  "iat": 1735894000,
  "exp": 1735897600,
  "jti": "b1d4e8f2-3a5c-4d6e-8f90-1a2b3c4d5e6f",
  "scim_id": "8a2c4f10-6b3d-4e21-9f7a-1c2d3e4f5a6b",
  "email": "alice.miller@example.com",
  "given_name": "Alice",
  "family_name": "Miller",
  "app_tid": "a1b2c3d4-5e6f-7a8b-9c0d-1e2f3a4b5c6d",
  "functionalAuthz": ["read@salesOrders", "write@salesOrders"]
}
```

The `functionalAuthz` claim is populated from the `applicationFunctionalAuthz` source.
:::

### Use Case 2: UI Visibility for Your Own Application

Applications that *do* have complex, conditional authorization but want to control the
**visibility of UI features** based on permissions can use the functional authorizations from the
token to show or hide features. For any API call, the backend still performs the authoritative
authorization check with the AMS client library based on the
[Authorization Bundle](/Authorization/AuthorizationBundle).

In this case you keep the bundle, so do **not** set `disable-authz-bundle`. Because typically not
all authorizations of an application are relevant for the UI, and to keep the token small, you must
**explicitly annotate** which actions of which resources are relevant. Only annotated actions are
exposed through the `applicationFunctionalAuthz` assertion attribute.

The annotation is done with the `@jwtRelevant` annotation on actions inside a `RESOURCE`
declaration in your `schema.dcl`:

```dcl [schema.dcl]
SCHEMA {
    // your attributes
}

RESOURCE salesOrders { // [!code focus:7]
    ACTIONS {
        @jwtRelevant
        read,
        write
    }
}
```

In this example only `read` is annotated, so only `read@salesOrders` becomes visible in the
token (independent of which policy grants it). Note that `write` is still listed in the
`ACTIONS` block even though it is not annotated, which is required by the rule below.

::: warning Minimum tooling versions required
`RESOURCE` declarations are a recent DCL language feature. Older tooling reports them as invalid
syntax, which breaks builds that compile DCL:

- `@sap/ams-dev` (DCL compilation before local tests): version 3.0.2 or higher
- `dcl-compiler-plugin` (Maven builds): version 1.5 or higher
- DCL VSCode extension: TBD

:::

The `ACTIONS` list also has a compact single-line notation, which supports the annotation as well:

```dcl
RESOURCE salesOrders {
    ACTIONS @jwtRelevant read, write; // jwtRelevant only applies to read, needs to be repeated on every action individually
}
```

::: warning Declaring a resource requires listing all of its actions 
Once you declare a resource with an `ACTIONS` block, that block must list **every** action that any policy uses on the resource, not only the ones you mark
`@jwtRelevant`. This completeness is validated when DCL is compiled/deployed: an action that is used in a policy but
missing from the resource declaration fails this validation. The `@jwtRelevant` annotation then selects which of those
declared actions are exposed in the token.
:::

```dcl
POLICY "Manage Sales Orders" {
    GRANT read ON salesOrders WHERE region IS NOT RESTRICTED;
    GRANT write ON salesOrders WHERE region IS NOT RESTRICTED;
}

POLICY "Read Sales Orders" {
    GRANT read ON salesOrders WHERE region IS NOT RESTRICTED;
}
```

In this example with the `RESOURCE` definition above, even if the `Manage Sales Orders` policy is assigned, only
`read@salesOrders` appears in the token. The `write` on `salesOrders` is still granted to the user, but is only
evaluable via the AMS client library in the backend.

**Inheritance through admin policies.** Customer administrators can derive restricted admin
policies from your base policies. The `@jwtRelevant` marking is inherited through `USE`:

```dcl
POLICY "Read Sales Orders Germany" {
    USE ams."Read Sales Orders" RESTRICT region = 'DE';
}
```

A user assigned to this admin policy also has `read@salesOrders` in the token, inherited from the
base policy. As always, the `region = 'DE'` condition must still be evaluated in the backend with
the AMS client library.

::: details Complete access token example
```json
{
  "iss": "https://mytenant.accounts.ondemand.com",
  "sub": "8a2c4f10-6b3d-4e21-9f7a-1c2d3e4f5a6b",
  "aud": "my-app-client-id",
  "azp": "my-app-client-id",
  "iat": 1735894000,
  "exp": 1735897600,
  "jti": "b1d4e8f2-3a5c-4d6e-8f90-1a2b3c4d5e6f",
  "scim_id": "8a2c4f10-6b3d-4e21-9f7a-1c2d3e4f5a6b",
  "email": "alice.miller@example.com",
  "given_name": "Alice",
  "family_name": "Miller",
  "app_tid": "a1b2c3d4-5e6f-7a8b-9c0d-1e2f3a4b5c6d",
  "functionalAuthz": ["read@salesOrders"]
}
```

The `functionalAuthz` claim is populated from the `applicationFunctionalAuthz` source. Only
`read@salesOrders` appears because only `read` is annotated with `@jwtRelevant`.
:::

#### CAP applications: roles via `$SCOPES`

In CAP applications the authorization model is role-based: roles are modeled as actions on the
special [`$SCOPES`](/CAP/Basics#assign-role-keyword) resource. To expose selected roles in the
token, declare `$SCOPES` and annotate the relevant roles:

```dcl
RESOURCE $SCOPES {
    ACTIONS {
        @jwtRelevant
        SalesManager,
        SalesViewer
    }
}
```

```dcl
POLICY SalesManager {
    ASSIGN ROLE SalesManager;   // equivalent to GRANT SalesManager ON $SCOPES
}
```

For `$SCOPES`, the token contains only the role name (**without** the `@resource` suffix):

```json
{
  "functionalAuthz": ["SalesManager"]
}
```

The same completeness rule applies: every role assigned via `ASSIGN ROLE` (i.e. granted on
`$SCOPES`) in any policy must be listed in the `ACTIONS` block.

::: tip Bundle-based alternative
If your backend already loads the authorization bundle, it can compute the same visibility
information locally via
[`getPotentialPrivileges`](/Authorization/AuthorizationChecks#getpotentialprivileges) instead of
relying on the token. Use the token claim when the decision must be made where the bundle is not
available (for example in a pure frontend).
:::

### Use Case 3: Shared Authorization (Work Zone, Joule, Gateways)

This model applies wherever one application (the **consuming application**) needs to know the
authorizations a user has in *another* application (the **consumed application** that owns those
authorizations). The classic example is SAP Build Work Zone (consuming), which maps the
authorizations of your application (consumed) to Launchpad content (Common Data Model / CDM
content).

The configuration is split across the two applications:

**On the consumed application** (owns the authorizations, e.g. your Sales Orders app):

- Annotate the authorizations you want to share with `@jwtRelevant` in `schema.dcl`, exactly as in
  [Use Case 2](#use-case-2-ui-visibility-for-your-own-application).
- No assertion attribute and no dependency configuration are needed on this side.

**On the consuming application** (its token carries the authorizations, e.g. Work Zone):

- Declares the consumed application as a `sharedAuthorizationDependencies` entry in its
  configuration. See
  [Configure Shared Authorization Dependencies](https://help.sap.com/docs/cloud-identity-services/cloud-identity-services/configure-shared-authorization-dependencies?version=Cloud&ai=true)
  for how to set this up in the Admin Console.
- Needs to add the `dependencyFunctionalAuthz` assertion attribute (in the following example under the
  claim name `dependencyFuncAuthz`). This is a one-time setup in the consuming application, not to be done per integrated application

Note the direction: your application never puts these authorizations into its *own* token for this
use case; it only makes them available through `@jwtRelevant` so the consuming application can pick
them up. Every authorization the consuming application needs (for example to map it to CDM content)
must therefore be annotated in your `schema.dcl`, otherwise it never reaches that application.

The consuming application's token then contains the functional authorizations of the user (token
subject) for each configured dependent application, keyed by the global tenant ID of that
application:

```json
{
  "dependencyFuncAuthz": {
    "<salesOrdersAppGlobalTenantId>": ["read@salesOrders"],
    "<purchaseOrdersAppGlobalTenantId>": ["read@purchaseOrders"]
  }
}
```

The user must have the corresponding policies assigned in the consumed application for its
authorizations to appear.

::: details Complete access token example
```json
{
  "iss": "https://mytenant.accounts.ondemand.com",
  "sub": "8a2c4f10-6b3d-4e21-9f7a-1c2d3e4f5a6b",
  "aud": "workzone-client-id",
  "azp": "workzone-client-id",
  "iat": 1735894000,
  "exp": 1735897600,
  "jti": "b1d4e8f2-3a5c-4d6e-8f90-1a2b3c4d5e6f",
  "scim_id": "8a2c4f10-6b3d-4e21-9f7a-1c2d3e4f5a6b",
  "email": "alice.miller@example.com",
  "given_name": "Alice",
  "family_name": "Miller",
  "app_tid": "a1b2c3d4-5e6f-7a8b-9c0d-1e2f3a4b5c6d",
  "dependencyFuncAuthz": {
    "<salesOrdersAppGlobalTenantId>": ["read@salesOrders"],
    "<purchaseOrdersAppGlobalTenantId>": ["read@purchaseOrders"]
  }
}
```

The `dependencyFuncAuthz` claim is populated from the `dependencyFunctionalAuthz` source.
:::

See [App-to-App](/Authorization/App2App) for the general concepts of
app-to-app dependencies and token exchange.

## Constraints and Pitfalls

::: warning Actions and resources must not contain `@`
Because the `@` character concatenates the action and resource into a single functional
authorization string (`action@resource`), actions and resources themselves **must not** contain the
`@` character. Otherwise the string becomes ambiguous.

This is **not** validated when your DCL is compiled or deployed: an `@` in an action or resource
name does not cause an error, it only makes the resulting functional authorization ambiguous at
runtime. Avoiding it is therefore your responsibility.
:::

- **No conditions in the token.** Conditions (`WHERE`) are never included; always evaluate them
  in the backend with the client library.
- **Deduplication.** The same functional authorization can be granted by multiple policies; the
  token contains a deduplicated list of distinct values.

## Related Pages

- [Authorization Bundle](/Authorization/AuthorizationBundle): the authoritative, condition-aware
  authorization data.
- [Authorization Checks](/Authorization/AuthorizationChecks): performing checks and
  [querying potential privileges](/Authorization/AuthorizationChecks#querying-potential-privileges).
- [Getting Started](/Authorization/GettingStarted#provisioning-of-ams-instances): provisioning
  and instance configuration.
- [CAP Integration](/CAP/Basics#assign-role-keyword): `ASSIGN ROLE` and the `$SCOPES` resource.
- [App-to-App](/Authorization/App2App): app-to-app dependencies.
