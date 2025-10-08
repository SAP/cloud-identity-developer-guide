# Go-ams

## Client Library
For details of the Go Authorization Client Library please refer to the repository [here](https://github.com/SAP/cloud-identity-authorizations-golang-library).

## Configuration

### Memory Usage
The memory usage of AMS in Go is very similar to the memory usage in Java. \n
The formula to calculate the memory usage is: 
````
memory_usage_in_kb = 0.2 * number_tenants + 0.25 * number_user + 0.1 * number_assignments
````


