// Package diff provides utilities for comparing, transforming, and analyzing
// .env files across environments.
//
// # Transform
//
// The Transform function applies one or more TransformOptions to an env map,
// returning a new map with the changes applied along with a list of keys that
// were modified. The original map is never mutated.
//
// Built-in transform functions:
//   - TransformTrimSpace: removes leading/trailing whitespace from values.
//   - TransformUppercase: converts values to UPPER CASE.
//   - TransformLowercase: converts values to lower case.
//
// Example:
//
//	res := diff.Transform(env, []diff.TransformOption{
//		{Keys: []string{"SECRET"}, Fn: diff.TransformTrimSpace()},
//	})
//	for _, k := range res.Changed {
//		fmt.Printf("%s => %s\n", k, res.Out[k])
//	}
package diff
