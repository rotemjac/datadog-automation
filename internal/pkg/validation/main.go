package validation

import (
	"github.com/go-playground/validator/v10"
	"strings"
)

// use a single instance of Validate, it caches struct info
var validate *validator.Validate

func ValidateStruct(structToValidate any) error {
	validate = validator.New()

	// returns nil or ValidationErrors ( []FieldError )
	err := validate.Struct(structToValidate)
	return err
	//if err != nil {

	// This check is only needed when your code could produce an invalid value for validation such as interface with nil
	//if _, ok := err.(*validator.InvalidValidationError); ok {
	//	fmt.Println(err)
	//}
	//for _, err := range err.(validator.ValidationErrors) {
	//	fmt.Println(err.Namespace())
	//	fmt.Println(err.Field())
	//	fmt.Println(err.StructNamespace())
	//	fmt.Println(err.StructField())
	//	fmt.Println(err.Tag())
	//	fmt.Println(err.ActualTag())
	//	fmt.Println(err.Kind())
	//	fmt.Println(err.Type())
	//	fmt.Println(err.Value())
	//	fmt.Println(err.Param())
	//	fmt.Println()
	//}
	//return false
}

func ValidateVar(varToValidate any, tag string) error {
	validate = validator.New()
	err := validate.Var(varToValidate, tag)
	return err
}

func RegisterStructValidationMapRules(rules map[string]string, types ...interface{}) {
	validate = validator.New()
	validate.RegisterStructValidationMapRules(rules, types)
}

func GetValidationRuleForInt(isRequired bool, gt int, lt int) {

}

func GetValidationRuleForString(isRequired bool, eqToOneOf string) string {
	//`validate:"required,oneof=dev staging prod"`
	var sb strings.Builder

	if isRequired {
		sb.WriteString("required")
	}
	//Check if string pointer in not null
	if &eqToOneOf != nil {
		// If buffer is not empty - means that we need to append with "," in order to separate from previous value
		if sb.Len() > 0 {
			sb.WriteString(",")
		}
		sb.WriteString("oneof=")
		sb.WriteString(eqToOneOf)
	}

	return sb.String()
}
