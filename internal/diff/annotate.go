package diff

// Annotation holds a note attached to a specific key.
type Annotation struct {
	Key  string
	Note string
}

// AnnotatedResult pairs a Result with per-key annotations.
type AnnotatedResult struct {
	Result      Result
	Annotations map[string]string
}

// Annotate builds an AnnotatedResult by applying a set of annotators to each
// entry in the result. An annotator receives a DiffEntry and returns a note
// string (empty string means no annotation).
func Annotate(r Result, annotators []func(DiffEntry) string) AnnotatedResult {
	notes := make(map[string]string)
	entries := Flatten(r)
	for _, e := range entries {
		for _, fn := range annotators {
			if note := fn(e); note != "" {
				notes[e.Key] = note
				break
			}
		}
	}
	return AnnotatedResult{Result: r, Annotations: notes}
}

// AnnotatorMissingInB flags keys absent from the right-hand file.
func AnnotatorMissingInB(e DiffEntry) string {
	if e.Status == "missing_in_b" {
		return "key not present in right file"
	}
	return ""
}

// AnnotatorMissingInA flags keys absent from the left-hand file.
func AnnotatorMissingInA(e DiffEntry) string {
	if e.Status == "missing_in_a" {
		return "key not present in left file"
	}
	return ""
}

// AnnotatorMismatched flags keys whose values differ.
func AnnotatorMismatched(e DiffEntry) string {
	if e.Status == "mismatched" {
		return "value mismatch between files"
	}
	return ""
}
