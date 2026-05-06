package unit_test

import (
	"regexp"
	"strings"
	"testing"
)

// Expression régulière reproduisant la logique de génération de slug.
var nonAlphanumeric = regexp.MustCompile(`[^a-z0-9]+`)

func generateSlug(name string) string {
	s := nonAlphanumeric.ReplaceAllString(strings.ToLower(name), "-")
	return strings.Trim(s, "-")
}

// Vérifie que le slug est entièrement en minuscules.
func TestSlug_Lowercase(t *testing.T) {
	slug := generateSlug("ACME Corp")
	if slug != strings.ToLower(slug) {
		t.Fatalf("le slug doit être en minuscules, obtenu : %q", slug)
	}
}

// Vérifie que les espaces et caractères spéciaux sont remplacés par des tirets.
func TestSlug_SpecialChars(t *testing.T) {
	slug := generateSlug("Acme & Co. Solutions!")
	if strings.ContainsAny(slug, " &.!") {
		t.Fatalf("le slug ne doit pas contenir de caractères spéciaux, obtenu : %q", slug)
	}
}

// Vérifie que le slug ne commence ni ne finit par un tiret.
func TestSlug_NoLeadingOrTrailingDash(t *testing.T) {
	slug := generateSlug("  --Acme--  ")
	if strings.HasPrefix(slug, "-") || strings.HasSuffix(slug, "-") {
		t.Fatalf("le slug ne doit pas commencer ou finir par un tiret, obtenu : %q", slug)
	}
}

// Vérifie qu'un nom simple produit le slug attendu.
func TestSlug_SimpleCase(t *testing.T) {
	cases := map[string]string{
		"Acme Corp":    "acme-corp",
		"Hello World":  "hello-world",
		"Skema CMS":    "skema-cms",
		"Test123":      "test123",
	}
	for input, expected := range cases {
		got := generateSlug(input)
		if got != expected {
			t.Errorf("slug(%q) = %q, attendu %q", input, got, expected)
		}
	}
}
