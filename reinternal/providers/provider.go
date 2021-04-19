package providers

import (
	"github.com/sumitAgrawal007/terraform-plugin-sdk/v2/reinternal/configs/configschema"
)

// Schema pairs a provider or resource schema with that schema's version.
// This is used to be able to upgrade the schema in UpgradeResourceState.
type Schema struct {
	Version int64
	Block   *configschema.Block
}
