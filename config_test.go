package main

import "testing"

func TestShouldInclude(t *testing.T) {
	tests := []struct {
		name           string
		itemTags       []string
		includeTags    []string
		requireAnyTags bool
		want           bool
	}{
		{
			name:           "Include all when no includeTags and no tag requirement",
			itemTags:       []string{"dev", "prod"},
			includeTags:    []string{},
			requireAnyTags: false,
			want:           true,
		},
		{
			name:           "Exclude untagged item when requireAnyTags is true",
			itemTags:       []string{},
			includeTags:    []string{},
			requireAnyTags: true,
			want:           false,
		},
		{
			name:           "Include untagged item when requireAnyTags is false",
			itemTags:       []string{},
			includeTags:    []string{},
			requireAnyTags: false,
			want:           true,
		},
		{
			name:           "Include item with matching tag",
			itemTags:       []string{"dev", "prod"},
			includeTags:    []string{"dev"},
			requireAnyTags: false,
			want:           true,
		},
		{
			name:           "Exclude item with non-matching tags",
			itemTags:       []string{"infra"},
			includeTags:    []string{"dev", "prod"},
			requireAnyTags: false,
			want:           false,
		},
		{
			name:           "Case-insensitive tag match",
			itemTags:       []string{"DeV"},
			includeTags:    []string{"dev"},
			requireAnyTags: false,
			want:           true,
		},
		{
			name:           "Matching tag with requireAnyTags = true",
			itemTags:       []string{"tools"},
			includeTags:    []string{"tools"},
			requireAnyTags: true,
			want:           true,
		},
		{
			name:           "Untagged item with includeTags should be excluded",
			itemTags:       []string{},
			includeTags:    []string{"dev"},
			requireAnyTags: false,
			want:           false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := shouldInclude(tc.itemTags, tc.includeTags, tc.requireAnyTags)
			if got != tc.want {
				t.Errorf("shouldInclude(%v, %v, %v) = %v; want %v",
					tc.itemTags, tc.includeTags, tc.requireAnyTags, got, tc.want)
			}
		})
	}
}
