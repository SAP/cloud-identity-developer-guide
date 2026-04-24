# DCL Compiler Maven Plugin

Maven plugin for compiling DCL (Data Control Language) files into DCN (Data Control Notation) format for local
application [testing](/Authorization/Testing).

## Installation

Add the plugin to your `pom.xml` (see also the [Testing guide](/Authorization/Testing#compiling-dcl-to-dcn) for full
integration examples):

```xml
<build>
    <plugins>
        <plugin>
            <groupId>com.sap.cloud.security.ams.dcl</groupId>
            <artifactId>dcl-compiler-plugin</artifactId>
            <version>${sap.cloud.security.ams.dcl-compiler.version}</version>
            <executions>
                <execution>
                    <goals>
                        <goal>compile</goal>
                    </goals>
                </execution>
            </executions>
        </plugin>
    </plugins>
</build>
```

The latest version can be found on [Maven Central](https://mvnrepository.com/artifact/com.sap.cloud.security.ams.dcl/dcl-compiler-plugin).

## Goals

### `dcl:compile`

Compiles DCL files to DCN format. Bound to the `generate-test-resources` phase by default.

### `dcl:validate`

Validates DCL files without generating output. Bound to the `validate` phase by default. Useful for CI/CD pipelines.

## Configuration

| Parameter | Default | Description |
|-----------|---------|-------------|
| `sourceDirectory` | `.../src/main/resources/ams/dcl` | Directory containing DCL source files |
| `outputDirectory` | `.../generated-test-resources/ams/dcn` | Output directory for compiled DCN files (compile only) |
| `skip` | `false` | Skip plugin execution |
| `verbose` | `false` | Enable verbose output |
| `failOn` | `error` | Failure threshold: `error`, `warning`, or `deprecation` |
| `readDcn` | `false` | Allow reading `.dcn` files as input |
| `timeout` | `60000` | CLI timeout in milliseconds |
| `additionalArguments` | - | Additional CLI arguments |

Parameters can be set in `<configuration>` or overridden via command line with `-Ddcl.<parameter>`, e.g. `mvn compile -Ddcl.verbose=true`.

## Platform-Specific Binaries

The plugin bundles platform-specific AMS CLI binaries for:
- macOS (Intel and Apple Silicon)
- Linux (x86_64 and ARM64)
- Windows (x86_64)

The relevant binaries are automatically extracted and
cached in `~/.ams-cli/binaries/{version}/` on first use.

## Troubleshooting

**Timeout errors** — Increase the timeout for large projects: `<timeout>120000</timeout>`

**Permission denied (Unix)** — The plugin sets executable permissions automatically. If this fails:
```bash
chmod +x ~/.ams-cli/binaries/{version}/{os}-{arch}/ams
```
