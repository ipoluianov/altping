package forms

import (
	"strings"
	"testing"
)

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

// The keys in the hints are the same in every language: they name the keys to press
func TestHotkeysKept(t *testing.T) {
	key := func(s string) string {
		if i := strings.LastIndex(s, " ("); i >= 0 {
			return s[i:]
		}
		return ""
	}
	hints := func(s *Strings) []string {
		return []string{s.ToolConfigs, s.ToolAddHost, s.ToolEditHost, s.ToolRemoveHosts, s.ToolDetails,
			s.ToolTray, s.ToolShareOpen, s.MenuEdit, s.MenuRemove}
	}
	want := hints(&en)
	for _, lang := range catalog.Languages() {
		for i, hint := range hints(catalog.Get(lang)) {
			if key(hint) != key(want[i]) {
				t.Errorf("%s: %q, want the key %q", lang, hint, key(want[i]))
			}
		}
	}
}
