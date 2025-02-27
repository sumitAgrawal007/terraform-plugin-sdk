package validation

import (
	"fmt"
	"reflect"

	"github.com/hashicorp/go-cty/cty"
	"github.com/sumitAgrawal007/terraform-plugin-sdk/v2/diag"
	"github.com/sumitAgrawal007/terraform-plugin-sdk/v2/helper/schema"
)

// NoZeroValues is a SchemaValidateFunc which tests if the provided value is
// not a zero value. It's useful in situations where you want to catch
// explicit zero values on things like required fields during validation.
func NoZeroValues(i interface{}, k string) (s []string, es []error) {
	if reflect.ValueOf(i).Interface() == reflect.Zero(reflect.TypeOf(i)).Interface() {
		switch reflect.TypeOf(i).Kind() {
		case reflect.String:
			es = append(es, fmt.Errorf("%s must not be empty, got %v", k, i))
		case reflect.Int, reflect.Float64:
			es = append(es, fmt.Errorf("%s must not be zero, got %v", k, i))
		default:
			// this validator should only ever be applied to TypeString, TypeInt and TypeFloat
			panic(fmt.Errorf("can't use NoZeroValues with %T attribute %s", i, k))
		}
	}
	return
}

// All returns a SchemaValidateFunc which tests if the provided value
// passes all provided SchemaValidateFunc
func All(validators ...schema.SchemaValidateFunc) schema.SchemaValidateFunc {
	return func(i interface{}, k string) ([]string, []error) {
		var allErrors []error
		var allWarnings []string
		for _, validator := range validators {
			validatorWarnings, validatorErrors := validator(i, k)
			allWarnings = append(allWarnings, validatorWarnings...)
			allErrors = append(allErrors, validatorErrors...)
		}
		return allWarnings, allErrors
	}
}

// Any returns a SchemaValidateFunc which tests if the provided value
// passes any of the provided SchemaValidateFunc
func Any(validators ...schema.SchemaValidateFunc) schema.SchemaValidateFunc {
	return func(i interface{}, k string) ([]string, []error) {
		var allErrors []error
		var allWarnings []string
		for _, validator := range validators {
			validatorWarnings, validatorErrors := validator(i, k)
			if len(validatorWarnings) == 0 && len(validatorErrors) == 0 {
				return []string{}, []error{}
			}
			allWarnings = append(allWarnings, validatorWarnings...)
			allErrors = append(allErrors, validatorErrors...)
		}
		return allWarnings, allErrors
	}
}

// ToDiagFunc is a wrapper for legacy schema.SchemaValidateFunc
// converting it to schema.SchemaValidateDiagFunc
func ToDiagFunc(validator schema.SchemaValidateFunc) schema.SchemaValidateDiagFunc {
	return func(i interface{}, p cty.Path) diag.Diagnostics {
		var diags diag.Diagnostics

		attr := p[len(p)-1].(cty.GetAttrStep)
		ws, es := validator(i, attr.Name)

		for _, w := range ws {
			diags = append(diags, diag.Diagnostic{
				Severity:      diag.Warning,
				Summary:       w,
				AttributePath: p,
			})
		}
		for _, e := range es {
			diags = append(diags, diag.Diagnostic{
				Severity:      diag.Error,
				Summary:       e.Error(),
				AttributePath: p,
			})
		}
		return diags
	}
}
