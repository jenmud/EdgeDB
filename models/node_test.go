package models_test

import (
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/jenmud/edgedb/models"
)

// CampareNodes is a go-cmp composer used to compare two nodes for equality in tests.
func CompareNodes(a, b *models.Node) bool {
	if a == nil || b == nil {
		return a == b
	}

	return a.ID == b.ID &&
		a.Label == b.Label &&
		cmp.Equal(a.Properties, b.Properties, cmpopts.EquateEmpty())
}
