package plugin

import (
	"sort"
	"strings"
)

// PlatformsFromVariants builds grouped platforms/packages from legacy variants.
// Plugins should call this when composing Version payloads.
func PlatformsFromVariants(version, releaseDate, officialURL string, variants []Variant) []PlatformRelease {
	if len(variants) <= 1 {
		if len(variants) == 0 {
			return nil
		}
		v := Version{Version: version, ReleaseDate: releaseDate, OfficialURL: officialURL}
		return buildPlatformsFromVariants(v, variants)
	}

	grouped := make(map[string][]Variant)
	platformOrder := make([]string, 0, len(variants))
	for _, variant := range variants {
		platform := strings.TrimSpace(variant.Platform)
		if platform == "" {
			platform = "Unknown"
		}
		if _, ok := grouped[platform]; !ok {
			platformOrder = append(platformOrder, platform)
		}
		grouped[platform] = append(grouped[platform], variant)
	}

	sort.Strings(platformOrder)
	out := make([]Variant, 0, len(variants))
	for _, platform := range platformOrder {
		group := grouped[platform]
		sort.SliceStable(group, func(i, j int) bool {
			return strings.ToLower(strings.TrimSpace(group[i].Architecture)) < strings.ToLower(strings.TrimSpace(group[j].Architecture))
		})
		out = append(out, group...)
	}
	v := Version{Version: version, ReleaseDate: releaseDate, OfficialURL: officialURL}
	return buildPlatformsFromVariants(v, out)
}

func buildPlatformsFromVariants(version Version, variants []Variant) []PlatformRelease {
	if len(variants) == 0 {
		return nil
	}
	grouped := make(map[string][]PlatformPackage)
	platformOrder := make([]string, 0, len(variants))
	for _, variant := range variants {
		platform := strings.TrimSpace(variant.Platform)
		if platform == "" {
			platform = "Unknown"
		}
		if _, ok := grouped[platform]; !ok {
			platformOrder = append(platformOrder, platform)
		}
		grouped[platform] = append(grouped[platform], PlatformPackage{
			Architecture: variant.Architecture,
			Links:        variant.Links,
		})
	}

	sort.Strings(platformOrder)
	out := make([]PlatformRelease, 0, len(platformOrder))
	for _, platform := range platformOrder {
		packages := grouped[platform]
		sort.SliceStable(packages, func(i, j int) bool {
			return strings.ToLower(strings.TrimSpace(packages[i].Architecture)) < strings.ToLower(strings.TrimSpace(packages[j].Architecture))
		})
		out = append(out, PlatformRelease{
			Platform:    platform,
			Version:     strings.TrimSpace(version.Version),
			ReleaseDate: strings.TrimSpace(version.ReleaseDate),
			OfficialURL: strings.TrimSpace(version.OfficialURL),
			Packages:    packages,
		})
	}
	return out
}
