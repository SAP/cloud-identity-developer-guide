# CAP Integration
This document describes the basics of the integration concept for the use of AMS in a [CAP](https://cap.cloud.sap/docs/) application.

::: info
The core concepts described in the other pages of the documentation apply also to CAP applications unless indicated otherwise.
:::

## Role Policies
In CAP applications, the authorization model is [role-based](https://cap.cloud.sap/docs/guides/security/authorization#role-based-authorization).
The cds authorization annotations - such as `@requires` and `@restrict` - already define the actions which users with specific roles may perform on specific resources (*entities*, *services*, *aspects*...).
It makes no sense to duplicate this information in *action*/*resource* based AMS policies. Instead, AMS is used to assign the roles from the cds authorization model to users.

For example, the following policy would grant the `SalesRepresentative` role when assigned to a user:

```dcl
POLICY SalesRepresentative {
    ASSIGN ROLE SalesRepresentative;
}
```

::: tip Policy Generation
The AMS Node.js module [`@sap/ams`](https://www.npmjs.com/package/@sap/ams) can be used to [generate](/CAP/cds-Plugin#base-policy-generation) base policies from a cds model.
:::

### Role Policy Guidelines
There is no technical requirement for the policy name to match the role name. For example, policies can also be used to define higher-level business roles that assign lower-level technical roles from the cds model:
```dcl
POLICY SalesRepresentative {
    ASSIGN ROLE SalesViewer;
    ASSIGN ROLE SalesCampaignViewer;
}
```

If a policy assigns exactly one role, it is a good practice to use the same name of the role as policy name.

## ASSIGN ROLE keyword
The `ASSIGN ROLE` keyword is CAP-specific syntactic sugar of DCL to abstract away from the underlying *action*/*resource* model of AMS and to allow a more intuitive way of defining policies that assign cds roles to users.

The policy above is equivalent to the following syntax:

```dcl
POLICY SalesRepresentative {
    GRANT SalesRepresentative ON $SCOPES;
}
```

As you can see, cds roles in AMS are modelled as actions on a special resource called `$SCOPES` which is typically the only AMS resource in a CAP application.

## DCL Root Folder
The `*.dcl` files in which authorization policies are defined are located in the **DCL root folder** of a CAP application.
In this folder, there is the `schema.dcl` file on root level alongside one or multiple DCL packages (subfolders) that contain the policies.

The DCL root folder is expected by the AMS modules in the following places by default, depending on the project type:
- **Java**: `srv/src/main/resources/ams`
- **Node.js**: `ams/dcl`
