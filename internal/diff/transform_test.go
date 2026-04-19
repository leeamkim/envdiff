package diff

import (
	"testing"
)

func TestTransform_NoOpts(t *testing.T) {
	env := map[string]string{"A": "hello", "B": "world"}
	res := Transform(env, nil)
	if res.Out["A"] != "hello" || res.Out["B"] != "world" {
		t.Fatal("expected unchanged map")
	}
	if len(res.Changed) != 0 {
		t.Fatalf("expected no changes, got %v", res.Changed)
	}
}

func TestTransform_AllKeys(t *testing.T) {
	env := map[string]string{"A": "  hello  ", "B": "  world  "}
	res := Transform(env, []TransformOption{{Fn: TransformTrimSpace()}})
	if res.Out["A"] != "hello" || res.Out["B"] != "world" {
		t.Fatalf("unexpected values: %v", res.Out)
	}
	if len(res.Changed) != 2 {
		t.Fatalf("expected 2 changed, got %v", res.Changed)
	}
}

func TestTransform_SpecificKeys(t *testing.T) {
	env := map[string]string{"A": "hello", "B": "world"}
	res := Transform(env, []TransformOption{
		{Keys: []string{"A"}, Fn: TransformUppercase()},
	})
	if res.Out["A"] != "HELLO" {
		t.Fatalf("expected HELLO, got %s", res.Out["A"])
	}
	if res.Out["B"] != "world" {
		t.Fatal("B should be unchanged")
	}
	if len(res.Changed) != 1 || res.Changed[0] != "A" {
		t.Fatalf("expected [A], got %v", res.Changed)
	}
}

func TestTransform_DoesNotMutateOriginal(t *testing.T) {
	env := map[string]string{"KEY": "value"}
	Transform(env, []TransformOption{{Fn: TransformUppercase()}})
	if env["KEY"] != "value" {
		t.Fatal("original map was mutated")
	}
}

func TestTransform_NoChangeWhenAlreadyTransformed(t *testing.T) {
	env := map[string]string{"KEY": "VALUE"}
	res := Transform(env, []TransformOption{{Fn: TransformUppercase()}})
	if len(res.Changed) != 0 {
		t.Fatalf("expected no changes, got %v", res.Changed)
	}
}

func TestTransform_MissingKeyIgnored(t *testing.T) {
	env := map[string]string{"A": "hello"}
	res := Transform(env, []TransformOption{
		{Keys: []string{"MISSING"}, Fn: TransformUppercase()},
	})
	if len(res.Changed) != 0 {
		t.Fatalf("expected no changes, got %v", res.Changed)
	}
}
