# cloud-identity-authorizations-golang-library

The official Go client library is https://github.com/SAP/cloud-identity-authorizations-golang-library.

## Documentation
Unfortunately, documentation for the Go client library is not yet available. Once it becomes available, it will be released here.

## Configuration

### Memory Usage
The memory usage of AMS in Go is very similar to the memory usage in Java.\
The formula to calculate the memory usage is: 
````
memory_usage_in_kb = 0.2 * number_tenants + 0.25 * number_user + 0.1 * number_assignments
````
