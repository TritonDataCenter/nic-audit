# NIC Compliance Auditing Tool

This tool will allow you to scan a set of Triton accounts for the presence of 
non-compliant network configurations. Typically, you would configure the tool
to search for instances where a public network and a secure private network are
attached to the same instance in order to enforce network security compliance
rules. Additionally, automatic removal of networks upon detection of a 
non-compliant configuration is possible.

## Configuration

The `nic-audit` tool supports a single parameter `-c` or `--config` which
specifies the path to the configuration file for the utility. The configuration
file is in the [json5](https://github.com/json5/json5) format. An example 
configuration file can be found [here](example/nic-audit.json5).

## Runtime

When run informational messages are written to STDERR and audit or compliance
messages are written to STDOUT. If configured, the tool can send emails 
containing aggregated alerts per Triton account. After a successful execution,
the utility will exit.