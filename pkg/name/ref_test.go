// Copyright 2018 Google LLC All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package name

import (
	"testing"
)

var (
	testDefaultRegistry = "registry.upbound.io"
	testDefaultTag      = "stable"
	inputDefaultNames   = []string{
		"crossplane/provider-gcp",
		"crossplane/provider-gcp:v0.14.0",
		"ubuntu",
		"gcr.io/crossplane/provider-gcp:latest",
	}
	outputDefaultNames = []string{
		"registry.upbound.io/crossplane/provider-gcp:stable",
		"registry.upbound.io/crossplane/provider-gcp:v0.14.0",
		"registry.upbound.io/ubuntu:stable",
		"gcr.io/crossplane/provider-gcp:latest",
	}
)

func TestParseReferenceDefaulting(t *testing.T) {
	for i, name := range inputDefaultNames {
		ref, err := ParseReference(name, WithDefaultRegistry(testDefaultRegistry), WithDefaultTag(testDefaultTag))
		if err != nil {
			t.Errorf("ParseReference(%q); %v", name, err)
		}
		if ref.Name() != outputDefaultNames[i] {
			t.Errorf("ParseReference(%q); got %v, want %v", name, ref.String(), outputDefaultNames[i])
		}
	}
}

func TestParseReference(t *testing.T) {
	for _, name := range goodWeakValidationDigestNames {
		ref, err := ParseReference(name, WeakValidation)
		if err != nil {
			t.Errorf("ParseReference(%q); %v", name, err)
		}
		dig, err := NewDigest(name, WeakValidation)
		if err != nil {
			t.Errorf("NewDigest(%q); %v", name, err)
		}
		if ref != dig {
			t.Errorf("ParseReference(%q) != NewDigest(%q); got %v, want %v", name, name, ref, dig)
		}
	}

	for _, name := range goodStrictValidationDigestNames {
		ref, err := ParseReference(name, StrictValidation)
		if err != nil {
			t.Errorf("ParseReference(%q); %v", name, err)
		}
		dig, err := NewDigest(name, StrictValidation)
		if err != nil {
			t.Errorf("NewDigest(%q); %v", name, err)
		}
		if ref != dig {
			t.Errorf("ParseReference(%q) != NewDigest(%q); got %v, want %v", name, name, ref, dig)
		}
	}

	for _, name := range badDigestNames {
		if _, err := ParseReference(name, WeakValidation); err == nil {
			t.Errorf("ParseReference(%q); expected error, got none", name)
		}
	}

	for _, name := range goodWeakValidationTagNames {
		ref, err := ParseReference(name, WeakValidation)
		if err != nil {
			t.Errorf("ParseReference(%q); %v", name, err)
		}
		tag, err := NewTag(name, WeakValidation)
		if err != nil {
			t.Errorf("NewTag(%q); %v", name, err)
		}
		if ref != tag {
			t.Errorf("ParseReference(%q) != NewTag(%q); got %v, want %v", name, name, ref, tag)
		}
	}

	for _, name := range goodStrictValidationTagNames {
		ref, err := ParseReference(name, StrictValidation)
		if err != nil {
			t.Errorf("ParseReference(%q); %v", name, err)
		}
		tag, err := NewTag(name, StrictValidation)
		if err != nil {
			t.Errorf("NewTag(%q); %v", name, err)
		}
		if ref != tag {
			t.Errorf("ParseReference(%q) != NewTag(%q); got %v, want %v", name, name, ref, tag)
		}
	}

	for _, name := range badTagNames {
		if _, err := ParseReference(name, WeakValidation); err == nil {
			t.Errorf("ParseReference(%q); expected error, got none", name)
		}
	}
}

func TestMustParseReference(t *testing.T) {
	goodWeakValidationNames := append(goodWeakValidationTagNames, goodWeakValidationDigestNames...)
	for _, name := range goodWeakValidationNames {
		func() {
			defer func() {
				if err := recover(); err != nil {
					t.Errorf("MustParseReference(%q, WeakValidation); panic: %v", name, err)
				}
			}()
			MustParseReference(name, WeakValidation)
		}()
	}

	goodStrictValidationNames := append(goodStrictValidationTagNames, goodStrictValidationDigestNames...)
	for _, name := range goodStrictValidationNames {
		func() {
			defer func() {
				if err := recover(); err != nil {
					t.Errorf("MustParseReference(%q, StrictValidation); panic: %v", name, err)
				}
			}()
			MustParseReference(name, StrictValidation)
		}()
	}

	badNames := append(badTagNames, badDigestNames...)
	for _, name := range badNames {
		func() {
			defer func() { recover() }()
			ref := MustParseReference(name, WeakValidation)
			t.Errorf("MustParseReference(%q, WeakValidation) should panic, got: %#v", name, ref)
		}()
	}
}
