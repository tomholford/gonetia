package main

import (
	"strings"
	"testing"
)

func withFixtures(t *testing.T) {
	t.Helper()

	orig := [4]map[string]bool{
		singleEnglishWords,
		doubleEnglishWords,
		singleApproxWords,
		doubleApproxWords,
	}
	t.Cleanup(func() {
		singleEnglishWords = orig[0]
		doubleEnglishWords = orig[1]
		singleApproxWords = orig[2]
		doubleApproxWords = orig[3]
	})

	singleEnglishWords = map[string]bool{"zod": true, "mar": true}
	doubleEnglishWords = map[string]bool{"sampel": true}
	singleApproxWords = map[string]bool{"wic": true}
	doubleApproxWords = map[string]bool{"dozzod": true}
}

func TestValidate(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		in      string
		wantErr string
	}{
		{name: "star", in: "~marzod"},
		{name: "planet", in: "~sampel-palnet", wantErr: "must be a star"},
		{name: "galaxy", in: "~zod", wantErr: "must be a star"},
		{name: "garbage", in: "not-a-patp", wantErr: "invalid patp"},
		{name: "empty", in: "", wantErr: "invalid patp"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			err := validate(tt.in)
			if tt.wantErr == "" {
				if err != nil {
					t.Fatalf("validate(%q) = %v, want nil", tt.in, err)
				}
				return
			}
			if err == nil || err.Error() != tt.wantErr {
				t.Fatalf("validate(%q) = %v, want %q", tt.in, err, tt.wantErr)
			}
		})
	}
}

func TestDoubles(t *testing.T) {
	t.Parallel()

	if !doubles("~datnut-datnut") {
		t.Fatal("expected doubles to match identical segments")
	}
	if doubles("~datnut-marzod") {
		t.Fatal("expected doubles to reject different segments")
	}
}

func TestAlliteration(t *testing.T) {
	t.Parallel()

	if !alliteration("~bacbel-baldut") {
		t.Fatal("expected alliteration to match the same first letter")
	}
	if alliteration("~bacbel-datnut") {
		t.Fatal("expected alliteration to reject different first letters")
	}
}

func TestEnglishFilters(t *testing.T) {
	withFixtures(t)

	if !onlyEnglish("~zod-mar") {
		t.Fatal("expected onlyEnglish when both segments are english")
	}
	if onlyEnglish("~zod-wic") {
		t.Fatal("expected onlyEnglish to reject an approx segment")
	}
	if !anyEnglish("~zod-wic") {
		t.Fatal("expected anyEnglish when one segment is english")
	}
	if anyEnglish("~wic-dozzod") {
		t.Fatal("expected anyEnglish to reject approx-only names")
	}
}

func TestApproxFilters(t *testing.T) {
	withFixtures(t)

	if !anyApprox("~wic-zod") {
		t.Fatal("expected anyApprox when one segment is approx")
	}
	if anyApprox("~zod-mar") {
		t.Fatal("expected anyApprox to reject english-only names")
	}
	if !onlyApprox("~wic-zod") {
		t.Fatal("expected onlyApprox when both segments are in any wordlist")
	}
	if onlyApprox("~wic-xxxx") {
		t.Fatal("expected onlyApprox to reject an unknown segment")
	}
}

func TestFilterPlanetsIncludesFirst(t *testing.T) {
	withFixtures(t)

	planets := []string{"~zod-zod", "~mar-mar", "~wic-xxxx"}
	got := filterPlanets(planets, Doubles)
	want := []string{"~zod-zod", "~mar-mar"}
	if strings.Join(got, ",") != strings.Join(want, ",") {
		t.Fatalf("filterPlanets doubles = %v, want %v", got, want)
	}
}

func TestFilterPlanetsStrategies(t *testing.T) {
	withFixtures(t)

	planets := []string{"~zod-mar", "~wic-zod", "~datnut-datnut"}

	got := filterPlanets(planets, OnlyEnglish)
	if strings.Join(got, ",") != "~zod-mar" {
		t.Fatalf("OnlyEnglish = %v", got)
	}

	got = filterPlanets(planets, AnyApprox)
	if strings.Join(got, ",") != "~wic-zod" {
		t.Fatalf("AnyApprox = %v", got)
	}

	got = filterPlanets(planets, All)
	if strings.Join(got, ",") != strings.Join(planets, ",") {
		t.Fatalf("All = %v, want %v", got, planets)
	}
}

func TestVersionRequested(t *testing.T) {
	t.Parallel()

	if versionRequested([]string{"gonetia"}) {
		t.Fatal("bare invocation is not a version request")
	}
	if !versionRequested([]string{"gonetia", "--version"}) {
		t.Fatal("expected --version")
	}
	if !versionRequested([]string{"gonetia", "-version"}) {
		t.Fatal("expected -version")
	}
	if versionRequested([]string{"gonetia", "~marzod"}) {
		t.Fatal("a star argument is not a version request")
	}
}

func TestVersionDefault(t *testing.T) {
	t.Parallel()

	if Version != "dev" {
		t.Fatalf("Version = %q, want dev", Version)
	}
}

func TestLoadWordlists(t *testing.T) {
	if err := loadWordlists(); err != nil {
		t.Fatal(err)
	}
	if len(singleEnglishWords) == 0 || len(doubleEnglishWords) == 0 ||
		len(singleApproxWords) == 0 || len(doubleApproxWords) == 0 {
		t.Fatal("expected embedded wordlists to be non-empty")
	}
}
