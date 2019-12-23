# TeamCity Go Test Json

Whilst newer versions of [JetBrains TeamCity](https://www.jetbrains.com/teamcity) support exporting tests via `go test -json` - older versions do not.

This application assists with this by intercepting and parsing the output and outputting it using [TeamCity's service messages](https://confluence.jetbrains.com/display/TCD9/Build+Script+Interaction+with+TeamCity)

## Example Usage

```
teamcity-go-test-json -scope ./package -parallelism=3 -count=1
```

## Supported Arguments

The following arguments are supported on the command-line:

* `count` - The number of times which each test should be run. Defaults to `1`.

* `prefix` - The prefix for tests which should be run e.g. `TestAcc`.

* `scope` - The directory/scope at which tests should be run. Note: the suffix `/...` will be appended if not specified. e.g. `./package`)

* `parallelism` - The number of tests which should be run in parallel.

* `timeout` - The maximum test duration in hours. Defaults to `1`.

You can also set the Environment Variable `DEBUG` to any value to see the commands/output which is being parsed.

## Licence

Apache 2.0