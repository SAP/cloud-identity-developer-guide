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
Some example data.json sizes can be found in this table: 
| Tenants | User  | Assignments | Measured Difference to empty data.json | KB per Tenant (T)/User (U)/Assignment (A) |
|---------|-------|-------------|----------------------------------------|-------------------------------------------|
| 10      | 0     | 0           | 4                                      | 0.4 (T)                                   |
| 1000    | 0     | 0           | 216                                    | 0.216 (T)                                 |
| 10000   | 0     | 0           | 2517                                   | 0.2517 (T)                                |
| 10      | 100   | 0           | 28                                     | 0.26 (U)                                  |
| 1       | 10    | 20          | 5                                      | 0.3 (A)                                   |
| 1       | 100   | 200         | 45                                     | 0.085 (A)                                 |
| 10      | 1000  | 2000        | 421                                    | 0.069 (A)                                 |
| 100     | 10000 | 20000       | 4104                                   | 0.1138 (A)                                |
| 1000    | 10000 | 200000      | 24116                                  | 0.1077 (A)                                |

The increase in memory usage per tenant, user and policy assignment in Java is approximately linear. 