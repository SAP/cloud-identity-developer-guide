# Deploying DCL

Deploying DCL to your AMS instance has a lot of similarities to applying [DDL](https://en.wikipedia.org/wiki/Data_definition_language) to a database schema:
- Its availability is a crucial precondition for the application to start up.
- Changes in business logic may depend on changes in DCL.
- Multiple microservices may depend on central DCL files.
- Changes need to be backward compatible to not break existing applications.
- You need a reliable mechanism to get migrations right the first time

To help with this process, this document recommends best practices and showcases the official DCL deployment solution that is available for different platforms and scenarios.

::: tip CAP
In CAP projects, `@sap/ams` registers a custom build task that automatically creates an `ams-policies-deployer` application during `cds build`.
Additionally, the deployment artefacts below are automatically configured during `cds add ams` to deploy it.

However, if your application consists of multiple microservices, it makes sense to [disable](/Libraries/nodejs/sap_ams/sap_ams#configuration) the automatic creation of the deployer application and manually deploy from a [central repository](#microservice-applications).
:::

## AMS Policies Deployer App
The official solution for deploying DCL files to an AMS instance is to do it with an *AMS Policies Deployer App*.
It is a minimalistic Node.js application that uses the credentials for your *identity* service instance and the [deploy-dcl](#deploy-dcl-script) script from the `@sap/ams` Node.js client library to upload the DCL files to the AMS server.

The deployer app is a minimalistic deployment artefact that is not part of your backend service.
This separation of concerns allows you to deploy the DCL files independently from your backend service, typically as a pre-step, and to manage the DCL files in a dedicated repository in microservice architectures.

For example, in a monorepo, assuming `bookshop` is your backend service, a typical project structure could look like this:

```text title="Procurement"
├─ bookshop
└─ ams-policies-deployer
    ├─ package.json
    └─ dcl
        ├─ sales
        │  └─ contracts.dcl
        │  └─ orders.dcl
        └─ schema.dcl
```

Given this structure, you would typically replace the placeholders in the snippets below with:

<code v-pre>`{{appName}}`</code>: **bookshop**\
<code v-pre>`{{identityServiceInstanceName}}`</code>: **bookshop-identity**\
<code v-pre>`{{dclDeployerAppFolder}}`</code>: **ams-policies-deployer**

### package.json
The `package.json` from the [ams-dcl-content-deployer](https://www.npmjs.com/package/@sap/ams?activeTab=code) folder of `@sap/ams` is ready to be used for the AMS deployer in your project, independent of the languages used in other parts of your project, e.g. Java, Javascript or Go.
It defines a minimalistic Node.js module that simply executes the [deploy-dcl](/Libraries/nodejs/sap_ams/sap_ams#deploy-dcl-script) script from the `@sap/ams` Node.js client library.  

::: tip
You do not need a local Node.js installation for development or deployment as the script does not have to be executed locally.
:::

### DCL Files
While the deployer needs to contain your DCL files, you should not commit additional copies of the DCL files to version control for the deployer if they are already contained in your microservices. In the [Microservice Applications](#microservice-applications) section we give recommendations for avoiding that.

In monolithic applications, we recommend to place the DCL files in your application folder and copy them to the `ams-policies-deployer` folder during a build step before deployment.

## Platforms
The following section describes how to deploy an AMS policy deployer application to different cloud platforms.
Those platforms provide different mechanisms for deploying static content that are discussed in the next sections.

### Cloud Foundry
In Cloud Foundry, [Tasks](https://docs.cloudfoundry.org/devguide/using-tasks.html) are the recommended mechanism for running *one-off* tasks such as deploying DCL files.

For [MTA](#mta) based deployments, you can use the provided snippet to automatically execute a task at the beginning of the deployment.
Unfortunately, this is not possible when deploying via a [CF Manifest](#cf-manifest). As manual task registration is documented poorly, we show an alternative that uses an *app* deployment instead.

#### MTA
```yaml title=mta.yaml
modules:
  - name: {{appName}}-ams-policies-deployer
      type: javascript.nodejs
      path: {{dclDeployerAppFolder}}
      parameters:
        buildpack: nodejs_buildpack
        no-route: true
        no-start: true
        tasks:
          - name: deploy-dcl
            command: npm start
            memory: 512M
      requires:
        - name: {{identityServiceInstanceName}}
          parameters:
            config:
              credential-type: X509_GENERATED
              app-identifier: policy-deployer
```

We recommend to use a [deployed-after hook](https://help.sap.com/docs/SAP_HANA_PLATFORM/4505d0bdaf4948449b7f7379d24d0f0d/4050fee4c469498ebc31b10f2ae15ff2.html#parameters) to delay the deployment of the backend server until the base policies have been successfully uploaded:

```yaml title=mta.yaml
modules:
  ### Backend Server
  - name: {{appName}}
    deployed-after: 
      - {{appName}}-ams-policies-deployer
```

#### CF Manifest
Unlike MTAs, CF Manifests do not support the scripted registration of CF tasks and their manual registration is poorly documented. Therefore, the proposed manifest snippet defines another *application* instead of a task. As CF considers applications whose process exits as crashed (even with status code 0), the deployer application needs to idle to prevent constant restarts. It can be manually stopped after the DCL deployment to free resources. It will report about success or failure in its logs.

```yaml title="manifest.yml"
applications:
  - name: {{appName}}-ams-policies-deployer
    path: {{dclDeployerAppFolder}}
    no-route: true
    health-check-type: none
    memory: 256M
    instances: 1
    buildpack: nodejs_buildpack
    command: (npm start && echo "This application may now be stopped to free resources." || echo "AMS Policy Deployment unsuccessful.") && sleep infinity
    services:
      - name: {{identityServiceInstanceName}}
        parameters:
          credential-type: X509_GENERATED
          app-identifier: policy-deployer
```

### Kubernetes
A fitting resource on Kubernetes to do *one-off* tasks such as deploying static content, is a [Job](https://kubernetes.io/docs/Authorization/workloads/controllers/job/). The snippets below help you define such a job but it needs a container image to run.

#### Kyma
The following sections describe how to run a job on SAP BTP Kyma based on a policy deployer image. 

::: info Outside Kyma
Plain (i.e. non-Kyma) Kubernetes clusters should be able to use the strategies below but additional steps might be necessary, e.g.
- installing the BTP service operator
- making sure the service bindings format implements the SAP Kubernetes Service Bindings Spec
:::

#### Build Deployer Image
To build the deployer image, you can use this provided `Dockerfile` from the the [ams-dcl-content-deployer](https://www.npmjs.com/package/@sap/ams?activeTab=code) folder of `@sap/ams` which is based on a distro-less [Chainguard](https://edu.chainguard.dev/chainguard/chainguard-images/getting-started/node/) base image for Node with 0 known vulnerabilities at the time of writing.

Adjust the `Dockerfile` as follows to fit it to your project structure:
1. Set the correct source path to the `package.json` of the policy deployer that gets copied in line `COPY ./package.json /app/package.json`
1. Add an additional `COPY` step in the next line to copy your base policies into `/app/dcl`.
1. (Optional) Adjust the version of `@sap/ams` in the `package.json` of the policy deployer to a specific version.

Alternatively, you might also be able to write a custom `Dockerfile` that builds on top of our pre-built image and inserts your DCL policy files similar to step 2.

When a new version of `@sap/ams` is released that contains changes for the base policy upload, you need to build a new deployer image for that version to benefit from those changes.
Usually, this should not happen frequently, as the deploy script is stable by now.

::: warning Cache Warning
**Explicitly specify an exact `@sap/ams` version in `package.json` on subsequent builds to prevent docker from re-using the cached (previous) version of `@sap/ams` when building the image.**
:::

##### Job for self-built image
To run an AMS Deployer Image built and deployed by yourself, you can use the following Job descriptor:

```yml title="amsDclDeployer.yml"
apiVersion: batch/v1
kind: Job
metadata:
  name: {{appName}}-ams-policies-deployer
spec:
  completions: 1
  parallelism: 1
  ttlSecondsAfterFinished: 300 # 5 minutes
  template:
    spec:
      imagePullSecrets:
        - name: {{imagePullSecret}}
      containers:
      - image: {{URL to your self-built ams dcl deployer image}}
        name: ams-policies-deployer
        env:
        - name: SERVICE_BINDING_ROOT
          value: /bindings
        securityContext:
          allowPrivilegeEscalation: false
          capabilities:
            drop:
            - ALL
          privileged: false
          runAsNonRoot: true
          readOnlyRootFilesystem: false
        volumeMounts:        
        - mountPath: /bindings/identity/
          name: identity-binding
          readOnly: true
      restartPolicy: OnFailure
      volumes:      
      - name: identity-binding
        secret:
          secretName: {{certBasedIdentityBinding}}
```


## Microservice Applications

::: tip
The information in this section is our recommendation to the best of our knowledge and experience. It is not a requirement to follow it.

If you'd like to challenge our reasoning, we would love to get feedback via the Support channels to improve the experience for all customers.
:::

Applications whose architecture consists of different microservices that use a common authorization model defined via schema and base policies, should make sure to deploy DCL centrally and not from individual microservices.

::: warning
Subsequent DCL deployments will override each other on the AMS server because AMS does **not** support merging when deploying DCL.  
:::

### Motivation for central authorization model
Given the nature of deploying base policies being similar to database schema migrations, there are good reasons against a distributed authorization model.

Typically, the microservices share at least **some** parts of their authorization model, e.g. for modeling business central roles like `auditor`. There are two ways to deal with this requirement.

If each microservice should have its own, independent model, the common aspects of the model would lead to redundancy in the models. However, redundancy must be avoided at all costs - not just for development, but also because it will be visible to customer administrators. Imagine the same privilege being defined in multiple policies of individual microservices. That privilege would need to be assigned in *all* microservice DCL packages that are pushed with a redundant policy independently. This results in a terrible user experience for the customer.

::: tip
In addition, the microservice architecture should be abstracted away from customer administrators. A grouping of policies by microservices may not always feel like a coherent application.
:::

Alternatively, one could decide to share certain aspects of the model, by not splitting microservices up into completely different DCL packages.\
But that raises the question: *why not share everything centrally?*\
Is it really worth the trouble of having to solve both problems simultaneously: distributing shared base policies and schema attributes across microservices **and** merging individual base policies and schemata? We think not.

### Accessing central DCL files
If you maintain DCL files (and the policy deployer) in a dedicated, central repository, it is best practice to place your DCL files a `dcl` folder at root level and make the files available to your microservice repositories by including the central DCL repository as a git submodule for local tests and validation.

There are techniques to only use the relevant files from the central repository:

- [git sparse checkout](https://git-scm.com/docs/git-sparse-checkout) in cone mode allows you to only check out the `dcl` folder of the submodule
- `symbolic links` allow you to link specifically the `dcl` folder of the submodule into specific paths of your microservice repository
- The AMS CAP plugins can be configured to look for DCL files in a different root folder, e.g. in your submodule folder

::: info Alternative Solutions
If you do not want to use git submodules, there are alternative solutions available on different levels of the tech stack to distribute common files across microservices:

- [Git Monorepo](https://www.atlassian.com/git/tutorials/monorepos)
- [NPM Workspaces](https://docs.npmjs.com/cli/v8/using-npm/workspaces)
- [Maven Modules](https://maven.apache.org/guides/mini/guide-multiple-modules.html)
:::

## deploy-dcl script

The `deploy-dcl` script from `@sap/ams` pushes a DCL bundle (including schema.dcl, DCL root package and all subpackages) to the Identity service instance from the environment (see `deploy-dcl --help`):

```
Usage: npx --package=@sap/ams deploy-dcl -d [DCL_ROOT_DIR] -c [CREDENTIALS_FILE] -n [DEPLOYER_APP_NAME]

Options:
      --help         Show help                                         [boolean]
      --version      Show version number                               [boolean]
  -d, --dcl          [optional] path to the directory that contains the DCL root
                     package. If a path is provided via environment variable
                     AMS_DCL_ROOT, it overrides this option.
                                                       [string] [default: "dcl"]
  -c, --credentials  [optional] path to a JSON file containing the credentials
                     object of an Identity service binding. If omitted, will try
                     to find and use an Identity service binding from the
                     process environment.                               [string]
  -n, --name         [optional] a descriptive name of this deployer application
                     to trace back the currently deployed DCL bundle on the AMS
                     server to its source when DCL is deployed from more than
                     one source. If a name is provided directly via environment
                     variable AMS_APP_NAME or indirectly as application_name via
                     VCAP_APPLICATION on Cloud Foundry or the pod name on K8s,
                     it overrides this option.
                                   [string] [default: "@sap/ams:deploy-dcl"]

Examples:
  deploy-dcl                                Pushes the DCL content in ./dcl
                                            (including schema.dcl, DCL root
                                            package and all subpackages) to the
                                            identity service instance from the
                                            environment.
  deploy-dcl -d src/dcl -c config/ias.json  Pushes the DCL content from
  -n bookshop-dcl-deployer                  ./src/dcl using the SAP Cloud Identity Services
                                            credentials in ./config/ias.json.
                                            The deployer app name for this
                                            upload will be set to
                                            "bookshop-dcl-deployer" to be able
                                            to trace back the upload source to
                                            this deployer.
```