# Instance-Based Authorization

Policies that assign roles can be extended with attribute filters for instance-based authorization. This allows administrators of different tenants to define fine-grained policies with individual restrictions at runtime to customize the authorization model for their tenant.

## Motivation

Let's imagine a scenario where we want to empower tenant administrators to be able to give users the `SalesRepresentative` role but only for a specific `Region` and `ProductCategory`.
For example, the following policy would grant the `SalesRepresentative` role only in the `EU` region and only for products in the `Electronics` category:

```SQL
POLICY SalesRepresentativeEUElectronics {
    ASSIGN ROLE SalesRepresentative 
    WHERE Region = 'EU' AND ProductCategory = 'Electronics'; -- [!code ++]
}
```

However, administrators cannot simply write arbitrary DCL policies as free-text because it is important that the application supports the conditions that are used in the policies. In this example, the application must understand how and where to apply the conditions for `Region` and `ProductCategory`.

In the next sections, we will explain the necessary steps in the application to achieve this example.

## Schema Definition

As a basis, attributes used in policies must be defined in a `schema.dcl` file. This file must be located in the [DCL root folder](/CAP/Basics#dcl-root-folder) of the CAP application.

In this file, the data type of the attributes is defined, and it can also contain additional metadata annotations, such as those for [Value Help](/concepts/ValueHelp).

```sql
SCHEMA {
  @valueHelp { ... }  // details omitted for brevity
  Region : String,
  
  @valueHelp { ... } // details omitted for brevity
  ProductCategory : String,
}
```

::: tip Schema Generation
The AMS Node.js module [`@sap/ams`](https://www.npmjs.com/package/@sap/ams) can be used to [generate](/CAP/cds-Plugin#base-policy-generation) the `schema.dcl` from [`@ams.attributes`](#annotating-the-cds-model) annotations in a cds model.
:::


## Restricted Role Assignments

To provide guard rails for runtime policies, it is possible to extend role assignments in the base policies with a `WHERE` condition. In that condition, the `IS NOT RESTRICTED` and `IS RESTRICTED` keywords can be used to list the attributes that can (must) be restricted by the tenant administrator at runtime.

##### IS RESTRICTED
When assigning a base policy with an attribute that `IS RESTRICTED`, the attribute condition evaluates to `false` which typically means access by assigning the base policy is not possible at all. To gain access, administrators **must** derive a runtime policy from this base policy that restricts the attribute to a specific value.

##### IS NOT RESTRICTED
When assigning a base policy with an attribute that `IS NOT RESTRICTED`, the attribute condition evaluates to `true` which means the base policy can be assigned for unfiltered access. However, administrators **can** derive a runtime policy from this base policy to restrict the attribute to a specific value.

In our example, we want to allow (but not force) the tenant administrator to restrict the `Region` and `ProductCategory` attributes, so we extend the `SalesRepresentative` base policy as follows:

```SQL
POLICY SalesRepresentative {
    ASSIGN ROLE SalesRepresentative
    WHERE Region IS NOT RESTRICTED AND ProductCategory IS NOT RESTRICTED; -- [!code ++]
}
```

::: warning
In CAP, `WHERE` conditions behind role assignments must only be used to determine **what** users with the role are allowed to see. It must not be misunderstood as a means to control **if** the role is assigned or not based on attribute conditions.
:::

### Combining Attribute Conditions
Usually, when there are multiple attributes in the `WHERE` condition of a rule, developers want to combine the filter conditions of all attributes using `AND`. However, it is also possible to use `OR` to combine attribute conditions.

::: tip Partial Restrictions
In either case, it is important to consider the effects of the `IS RESTRICTED` and `IS NOT RESTRICTED` keywords. For example, administrators can decide to leave some attributes as `IS RESTRICTED` (`IS NOT RESTRICTED`) in a derived policy by restricting only a subset of attributes to specific conditions. The evaluation to `false` (`true`) of the `IS RESTRICTED` (`IS NOT RESTRICTED`) attributes will short-circuit the evaluation of the remaining attributes in the same `AND` (`OR`) clause.

For this reason, when using `IS NOT RESTRICTED`, we discourage the use of `OR` in the `WHERE` clause of role assignments because it can lead to unintended full access.
:::

## Runtime Policy Creation

Once the previous steps are in place, the tenant administrator can use the `SCI admin cockpit` to create a runtime policy from the base policy that restricts the `Region` and `ProductCategory` attributes to specific values.

::: tip
For local tests, such a derived policy can be written in a DCL file inside the [`local`](/concepts/Testing#test-policies) DCL package.
:::

```SQL
POLICY SalesRepresentativeEUElectronics {
    USE cap.SalesRepresentative
    RESTRICT Region = 'EU', ProductCategory = 'Electronics';
}
```

This derived policy is equivalent to the policy defined in the [Motivation](#motivation) section. It can be assigned to users in the same tenant like any other policy.

## Annotating the CDS Model

Finally, via `@ams.attributes` annotations, the AMS attributes are mapped to elements (or association paths) in the cds model via compile-safe cds expressions. Whenever requests access the annotated resources, the result is filtered based on the attribute conditions computed by AMS.

```js
annotate Product with @ams.attributes: { // [!code ++:8]
    ProductCategory: (category),
};

annotate SalesOrder with @ams.attributes: {
    Region: (region),
    ProductCategory: (product.category),
};

annotate Product with @restrict: [
    {
      grant: ['READ'],
      to: [ 'SalesManager', 'SalesRepresentative' ]
    }
];

annotate SalesOrder with @restrict: [
    {
      grant: [ 'CREATE', 'READ', 'UPDATE', 'DELETE' ],
      to: 'SalesManager',
    },
    {
      grant: [ 'READ' ],
      to: 'SalesRepresentative',
    },
];
```

::: tip
`ams.attributes` annotations are supported on *aspects*, *entities* and *actions/functions bound to a single entity* as those are the [cds resources that support *where* conditions](https://cap.cloud.sap/docs/guides/security/authorization#supported-combinations-with-cds-resources).
:::

## Effect of attribute filters

When a user is assigned the `SalesRepresentativeEUElectronics` policy, the AMS CAP modules will dynamically adjust the cds `where` condition to inject the AMS attribute conditions.

For example, when accessing the `SalesOrder` entity with this policy, the AMS module will add a `where` condition to the privilege for the `SalesRepresentative` role that looks like this:

```js
@restrict: [
    {
      grant: [ 'READ' ],
      to: 'SalesRepresentative',
      where: 'region = "EU" AND product.category = "Electronics"', // [!code ++]
    },
]
```

If a user has more than one cds role that grants access to a resource, the AMS module will combine the attribute conditions of all roles using `OR`. For example, if the user also has the `SalesManager` role assigned with unfiltered access, the resulting `where` condition on `Product` would effectively look like this (before being simplified by the AMS module):

```js
@restrict: [
    {
      grant: ['READ'],
      to: [ 'SalesManager', 'SalesRepresentative' ],
      where: 'true OR (region = "EU" AND product.category = "Electronics")', // [!code ++]
    }
]
```

::: tip
When there is already a static *where* condition on a cds privilege, the AMS module will combine the static *where* condition with the attribute conditions from AMS by using `AND`.
:::