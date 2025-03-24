SparkCI is a Go CLI tool using Cobra library that enhances GitLab CI workflows with various utilities.

For SparkCI code, follow these Go coding patterns:

Main commands are organized in their own package under `cmd/` (e.g., `cmd/gitlab/`) with the main command defined in a file with the same name (e.g., `gitlab.go`).

Main command files export a PascalCase command variable (e.g., `GitlabCommand`) that other files can import.

Subcommands are defined in separate files within the same package, with camelCase variables (e.g., `printEnv`). The filename should match the command's functionality (e.g., `printEnv.go`).

The main command's `init()` function adds all subcommands using `MainCommand.AddCommand(subcommand)`.

Use `RunE` instead of `Run` when the command can return errors, with `SilenceErrors: true` for custom error handling.

Command flags are defined in each command's `init()` function.

Core functionality is implemented in packages under `pkg/`. Import these packages in command files.

For logging, use the project's logger utilities: `utils.Info()`, `utils.Debug()`, `utils.Warn()`, `utils.Error()`, and `utils.Fatal()`.

Keep each command in its own file, with the command name matching the filename.

Main commands should provide a default `Run` function that calls `cmd.Help()` when no subcommand is specified.

To add a new top-level command, import and add it to the root command's `init()` function in `cmd/root.go`.

To add a new subcommand, create a new file in the appropriate package and add the subcommand to its parent in the `init()` function.

For commands processing CLI arguments, use `DisableFlagParsing: true` when custom argument handling is needed.

Provide examples of command usage in the `Long` description for complex commands.
