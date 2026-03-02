# ams-test

`ams-test` is a test support module that provides utilities for testing AMS authorization policies directly in Java unit
tests. It allows you to test your DCL policies with different input combinations and verify expected authorization
decisions without the overhead of integration tests.

## Installation

The module is available as a Maven dependency:

```xml

<dependency>
    <groupId>com.sap.cloud.security.ams</groupId>
    <artifactId>ams-test</artifactId>
    <scope>test</scope>
</dependency>
```

## When to Use

Use `ams-test` when you want to:

- **Unit test individual policies**: Test DCL policies in isolation with various input combinations
- **Verify authorization logic**: Ensure that your policies grant or deny access as expected
- **Test edge cases**: Verify behavior with partial inputs, null values, or boundary conditions
- **Debug policy issues**: Quickly identify problems in your authorization logic without running the full application

::: tip
This module is particularly useful for testing complex policies with multiple conditions before deploying to production.
It complements integration tests that verify the full application behavior.
:::

## Prerequisites

Before running tests, your DCL files must be compiled to DCN format.
See [Testing - Compiling DCL to DCN](/Authorization/Testing#compiling-dcl-to-dcn) for setup instructions.

## Usage

### Basic Setup with JUnit 5

Register the `AmsTestExtension` in your test class to load DCN files and perform authorization checks:

```java
import static com.sap.cloud.security.ams.Assertions.*;

import com.sap.cloud.security.ams.AmsTestExtension;
import com.sap.cloud.security.ams.api.Authorizations;
import org.junit.jupiter.api.Test;
import org.junit.jupiter.api.extension.RegisterExtension;

import java.nio.file.Path;

public class PolicyTest {

    @RegisterExtension
    static AmsTestExtension amsTest =
            AmsTestExtension.fromLocalDcn(Path.of("target/generated-test-resources/ams/dcn"));

    @Test
    void testReadAccess_granted() {
        Authorizations auth = amsTest.getAuthorizations("mypackage.ReadPolicy");
        assertGranted(() -> auth.checkPrivilege("read", "orders"));
    }

    @Test
    void testDeleteAccess_denied() {
        Authorizations auth = amsTest.getAuthorizations("mypackage.ReadPolicy");
        assertDenied(() -> auth.checkPrivilege("delete", "orders"));
    }
}
```

### Factory Methods

The `AmsTestExtension` provides several factory methods:

```java
// Use default DCN path (target/generated-test-resources/ams/dcn)
AmsTestExtension.fromLocalDcn();

// Use custom DCN path
AmsTestExtension.fromLocalDcn(Path.of("path/to/dcn"));
```

### DCL Packages

Use `forPackage()` to create a scoped accessor that prepends the package name to policy names. Define the package
accessor at the class level to reuse it across test cases:

```java

@RegisterExtension
static AmsTestExtension amsTest =
        AmsTestExtension.fromLocalDcn(Path.of("target/generated-test-resources/ams/dcn"));

// Create a package-scoped accessor for reuse across test cases
static TestAuthorizationsProvider ordersPackage = amsTest.forPackage("shopping.orders");

@Test
void testReadPolicy() {
    // Instead of "shopping.orders.ReadPolicy"
    Authorizations auth = ordersPackage.getAuthorizations("ReadPolicy");
    assertGranted(() -> auth.checkPrivilege("read", "orders"));
}

@Test
void testWritePolicy() {
    Authorizations auth = ordersPackage.getAuthorizations("WritePolicy");
    assertGranted(() -> auth.checkPrivilege("write", "orders"));
}
```

### Testing with Input Attributes

For policies with conditions based on input attributes, provide the attribute values in the authorization check:

```java
import com.sap.cloud.security.ams.api.expression.AttributeName;

import java.util.Map;

@Test
void testConditionalPolicy_withMatchingInput() {
    Authorizations auth = amsTest.getAuthorizations("mypackage.ConditionalPolicy");

    Map<AttributeName, Object> input = Map.of(
            AttributeName.of("amount"), 100,
            AttributeName.of("currency"), "EUR"
    );

    assertGranted(() -> auth.checkPrivilege("approve", "orders", input));
}

@Test
void testConditionalPolicy_withNonMatchingInput() {
    Authorizations auth = amsTest.getAuthorizations("mypackage.ConditionalPolicy");

    Map<AttributeName, Object> input = Map.of(
            AttributeName.of("amount"), 10000,
            AttributeName.of("currency"), "EUR"
    );

    assertDenied(() -> auth.checkPrivilege("approve", "orders", input));
}
```

### Handling Conditional Results

When a policy cannot be fully evaluated due to missing input attributes, the result is `CONDITIONAL`. This indicates
that the decision depends on further conditions over attributes for which no value was provided. For example:

```java

@Test
void testPartialInput_isConditional() {
    Authorizations auth = amsTest.getAuthorizations("mypackage.TwoVarPolicy");

    // Policy: GRANT action ON * WHERE a AND b
    // Only providing one variable results in CONDITIONAL
    Map<AttributeName, Object> partialInput = Map.of(
            AttributeName.of("a"), true
            // 'b' is not provided
    );

    assertConditional(() -> auth.checkPrivilege("action", "resource", partialInput));
}
```

### Testing with Empty Unknowns

To treat missing attributes as `null` rather than unknown (resulting in `DENIED` or `GRANTED` instead of `CONDITIONAL`),
pass an empty set of unknowns:

```java
import java.util.Set;

@Test
void testMissingInput_withoutUnknowns_isDenied() {
    Authorizations auth = amsTest.getAuthorizations("mypackage.RequiredInputPolicy");

    // With empty unknowns, missing attributes are treated as null
    assertDenied(() -> auth.checkPrivilege("action", "resource", Map.of(), Set.of()));
}
```

## Assertions API

The `Assertions` class provides static methods for verifying authorization decisions:

| Method                                               | Description                                           |
|------------------------------------------------------|-------------------------------------------------------|
| `assertGranted(Supplier<Decision>)`                  | Asserts that the decision is `GRANTED`                |
| `assertDenied(Supplier<Decision>)`                   | Asserts that the decision is `DENIED`                 |
| `assertConditional(Supplier<Decision>)`              | Asserts that the decision is `CONDITIONAL`            |
| `assertDecision(DecisionResult, Supplier<Decision>)` | Asserts that the decision matches the expected result |

### Example with assertDecision

```java
import com.sap.cloud.security.ams.api.DecisionResult;

@ParameterizedTest
@CsvSource({
        "true,  true,  GRANTED",
        "true,  false, DENIED",
        "false, true,  DENIED",
        "false, false, DENIED"
})
void testAndPolicy(boolean var1, boolean var2, DecisionResult expected) {
    Authorizations auth = amsTest.getAuthorizations("mypackage.AndPolicy");

    Map<AttributeName, Object> input = Map.of(
            AttributeName.of("var1"), var1,
            AttributeName.of("var2"), var2
    );

    assertDecision(expected, () -> auth.checkPrivilege("action", "resource", input));
}
```