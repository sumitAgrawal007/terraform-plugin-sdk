package structure

import (
	"reflect"

	"github.com/sumitAgrawal007/terraform-plugin-sdk/v2/helper/schema"
)

func SuppressJsonDiff(k, old, new string, d *schema.ResourceData) bool {
	oldMap, err := ExpandJsonFromString(old)
	if err != nil {
		return false
	}

	newMap, err := ExpandJsonFromString(new)
	if err != nil {
		return false
	}

	return reflect.DeepEqual(oldMap, newMap)
}
