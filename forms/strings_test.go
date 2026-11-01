package forms

import "testing"

// Every language has all the texts: a new field in Strings without a
// translation fails here instead of silently showing English
func TestAllTranslated(t *testing.T) {
	for lang, fields := range catalog.Missing() {
		t.Errorf("%s: not translated: %v", lang, fields)
	}
}

// Every translation is offered in the settings
func TestLanguagesOffered(t *testing.T) {
	offered := make(map[string]bool)
	for _, l := range languages {
		offered[l.tag] = true
	}
	for _, lang := range catalog.Languages() {
		if !offered[lang] {
			t.Errorf("%s: not in the languages list", lang)
		}
	}
}

// Every language has the key at the end of the hints the dialog titles are made of
func TestWithoutHotkey(t *testing.T) {
	for _, lang := range append(catalog.Languages(), "en") {
		s := catalog.Get(lang)
		for _, hint := range []string{s.ToolAddHost, s.ToolEditHost} {
			title := withoutHotkey(hint)
			if title == hint || title == "" {
				t.Errorf("%s: %q has no key to cut", lang, hint)
			}
		}
	}
}
