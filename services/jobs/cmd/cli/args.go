package main

import (
	"flag"
	"strings"
)

// splitCommandLineArgs separates global flags from flags belonging to the selected task.
func splitCommandLineArgs(flags *flag.FlagSet, args []string) (globalArgs []string, taskArgs []string) {
	taskSelected := false

	for i := 0; i < len(args); i++ {
		arg := args[i]

		// Everything after an explicit "--" separator belongs to the task
		if arg == "--" {
			return globalArgs, append(taskArgs, args[i+1:]...)
		}

		name, inlineValue, hasInlineValue, isFlag := parseFlagArgument(arg)
		if !isFlag {
			// After a task has been selected, a positional argument marks the
			// beginning of its arguments, otherwise it stays with the globals
			if taskSelected {
				return globalArgs, append(taskArgs, args[i:]...)
			}
			return append(globalArgs, args[i:]...), taskArgs
		}

		globalFlag := flags.Lookup(name)
		if globalFlag == nil {
			// Unknown flags that we have not registered yet belong to the task
			return globalArgs, append(taskArgs, args[i:]...)
		}

		// This is a recognized global flag -> keep it for the global parser
		globalArgs = append(globalArgs, arg)

		if name == "name" && hasInlineValue && inlineValue != "" {
			// Task is only selected if -name has a value
			taskSelected = true
		}

		// Inline values and boolean flags are complete in the current argument
		if hasInlineValue || isBooleanFlag(globalFlag) {
			continue
		}
		if i+1 >= len(args) {
			continue
		}

		// Other global flags consume the following argument as their value
		i++
		globalArgs = append(globalArgs, args[i])
		if name == "name" && args[i] != "" {
			taskSelected = true
		}
	}

	return globalArgs, taskArgs
}

func parseFlagArgument(arg string) (name string, value string, hasValue bool, ok bool) {
	// Ignore positional arguments and the special flag separators
	if len(arg) < 2 || arg[0] != '-' || arg == "-" || arg == "--" {
		return "", "", false, false
	}

	// Strip either the single-dash or double-dash prefix
	prefixLength := 1
	if arg[1] == '-' {
		prefixLength = 2
	}
	nameAndValue := arg[prefixLength:]

	// An equals sign means the flag contains an inline value
	if before, after, isOk := strings.Cut(nameAndValue, "="); isOk {
		return before, after, true, true
	}
	return nameAndValue, "", false, true
}

func isBooleanFlag(f *flag.Flag) bool {
	booleanFlag, ok := f.Value.(interface {
		IsBoolFlag() bool
	})
	return ok && booleanFlag.IsBoolFlag()
}
