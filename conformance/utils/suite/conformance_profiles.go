/*
Copyright 2023 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package suite

import (
	"fmt"

	"k8s.io/apimachinery/pkg/util/sets"
)

// -----------------------------------------------------------------------------
// Conformance Profiles - Public Types
// -----------------------------------------------------------------------------

// ConformanceProfile is a group of features that have a related purpose, e.g.
// to cover a specific feature present in Network Policy V2 API.
//
// For more details see the relevant GEP: https://network-policy-api.sigs.k8s.io/npeps/npep-137/
type ConformanceProfile struct {
	Name             ConformanceProfileName
	CoreFeatures     sets.Set[SupportedFeature]
	ExtendedFeatures sets.Set[SupportedFeature]
}

type ConformanceProfileName string

const (
	// ANPConformanceProfileName indicates the name of the conformance profile
	// which covers AdminNetworkPolicy core API
	ANPConformanceProfileName ConformanceProfileName = "AdminNetworkPolicy"

	// BANPConformanceProfileName indicates the name of the conformance profile
	// which covers BaselineAdminNetworkPolicy core API
	BANPConformanceProfileName ConformanceProfileName = "BaselineAdminNetworkPolicy"
)

// -----------------------------------------------------------------------------
// Conformance Profiles - Public Vars
// -----------------------------------------------------------------------------

var (
	// ANPConformanceProfile is a ConformanceProfile that covers testing ANP API
	ANPConformanceProfile = ConformanceProfile{
		Name: ANPConformanceProfileName,
		CoreFeatures: sets.New(
			SupportAdminNetworkPolicy,
		),
		ExtendedFeatures: sets.New(
			SupportAdminNetworkPolicyNamedPorts,
		),
	}

	// BANPConformanceProfile is a ConformanceProfile that covers testing BANP API
	BANPConformanceProfile = ConformanceProfile{
		Name: BANPConformanceProfileName,
		CoreFeatures: sets.New(
			SupportBaselineAdminNetworkPolicy,
		),
		ExtendedFeatures: sets.New(
			SupportBaselineAdminNetworkPolicyNamedPorts,
		),
	}
)

// -----------------------------------------------------------------------------
// Conformance Profiles - Private Profile Mapping Helpers
// -----------------------------------------------------------------------------

// conformanceProfileMap maps short human-readable names to their respective
// ConformanceProfiles.
var conformanceProfileMap = map[ConformanceProfileName]ConformanceProfile{
	ANPConformanceProfileName:  ANPConformanceProfile,
	BANPConformanceProfileName: BANPConformanceProfile,
}

// getConformanceProfileForName retrieves a known ConformanceProfile by it's simple
// human readable ConformanceProfileName.
func getConformanceProfileForName(name ConformanceProfileName) (ConformanceProfile, error) {
	profile, ok := conformanceProfileMap[name]
	if !ok {
		return profile, fmt.Errorf("%s is not a valid conformance profile", name)
	}

	return profile, nil
}

// getConformanceProfilesForTest retrieves the ConformanceProfiles a test belongs to.
func getConformanceProfilesForTest(test ConformanceTest, conformanceProfiles sets.Set[ConformanceProfileName]) sets.Set[*ConformanceProfile] {
	matchingConformanceProfiles := sets.New[*ConformanceProfile]()
	for _, conformanceProfileName := range conformanceProfiles.UnsortedList() {
		cp := conformanceProfileMap[conformanceProfileName]
		hasAllFeatures := true
		for _, feature := range test.Features {
			if !cp.CoreFeatures.Has(feature) && !cp.ExtendedFeatures.Has(feature) {
				hasAllFeatures = false
				break
			}
		}
		if hasAllFeatures {
			matchingConformanceProfiles.Insert(&cp)
		}
	}

	return matchingConformanceProfiles
}
