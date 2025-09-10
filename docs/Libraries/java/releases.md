# Release Notes for AMS Client Library Java 

## 3.7.0

- Maintenance release with updated dependencies and fixes for the Maven Central release process.

## 3.6.0

- The property `cds.security.mock.enabled` is now used to enable the mock users in the
  `cap-ams-support` module.
- A new property `ams.properties.bundleGatewayUpdater.maxFailedUpdates` is introduced to configure the maximum
  number of failed updates of the bundle gateway before it logs an error message. The default value is `0`.
