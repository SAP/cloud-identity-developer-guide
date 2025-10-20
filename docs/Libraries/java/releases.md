# Release Notes for AMS Client Library Java 

## 3.8.0

- This release removes the dependencies from `com.sap.cloud.security.ams.dcl` artifacts. All required classes,
interfaces, etc., are now part of the `jakarta-ams` module using the same packages. So, everything should continue
to work without any changes. Please remove any direct dependencies on `com.sap.cloud.security.ams.dcl` artifacts.

## 3.7.0

- Maintenance release with updated dependencies and fixes for the Maven Central release process.

## 3.6.0

- The property `cds.security.mock.enabled` is now used to enable the mock users in the
  `cap-ams-support` module.
- A new property `ams.properties.bundleGatewayUpdater.maxFailedUpdates` is introduced to configure the maximum
  number of failed updates of the bundle gateway before it logs an error message. The default value is `0`.
