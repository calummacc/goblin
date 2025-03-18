package validation

import (
	"fmt"
	"reflect"
	"regexp"
	"strings"
)

// ValidationError represents a validation error
type ValidationError struct {
	Field   string
	Message string
	Value   interface{}
	Rule    string
}

// Error returns the error message
func (e *ValidationError) Error() string {
	return fmt.Sprintf("validation failed for field '%s': %s (value: %v, rule: %s)",
		e.Field, e.Message, e.Value, e.Rule)
}

// Validator interface defines the contract for custom validators
type Validator interface {
	Validate(value interface{}, metadata ValidationMetadata) error
}

// ValidationMetadata contains validation rules and messages
type ValidationMetadata struct {
	Field      string
	Rules      map[string]interface{}
	Messages   map[string]string
	StopOnFail bool
}

// ValidationPipe handles validation of input data
type ValidationPipe struct {
	validators map[string]Validator
	stopOnFail bool
}

// NewValidationPipe creates a new ValidationPipe
func NewValidationPipe(stopOnFail bool) *ValidationPipe {
	return &ValidationPipe{
		validators: make(map[string]Validator),
		stopOnFail: stopOnFail,
	}
}

// RegisterValidator registers a custom validator
func (p *ValidationPipe) RegisterValidator(name string, validator Validator) {
	p.validators[name] = validator
}

// Transform validates and transforms input data
func (p *ValidationPipe) Transform(value interface{}, metadata ValidationMetadata) (interface{}, error) {
	errors := make([]*ValidationError, 0)
	valueType := reflect.TypeOf(value)

	// Handle struct validation
	if valueType.Kind() == reflect.Struct {
		valueValue := reflect.ValueOf(value)
		for i := 0; i < valueType.NumField(); i++ {
			field := valueType.Field(i)
			fieldValue := valueValue.Field(i)

			// Get validation tags
			if tag, ok := field.Tag.Lookup("validate"); ok {
				fieldMetadata := ValidationMetadata{
					Field:      field.Name,
					Rules:      make(map[string]interface{}),
					Messages:   make(map[string]string),
					StopOnFail: p.stopOnFail,
				}

				// Parse validation rules
				rules := strings.Split(tag, "|")
				for _, rule := range rules {
					parts := strings.Split(rule, ":")
					ruleName := parts[0]
					var ruleValue interface{}

					if len(parts) > 1 {
						ruleValue = parts[1]
					}

					fieldMetadata.Rules[ruleName] = ruleValue

					// Get custom message if specified
					if msgTag, ok := field.Tag.Lookup("msg"); ok {
						fieldMetadata.Messages[ruleName] = msgTag
					}
				}

				// Validate field
				if err := p.validateField(fieldValue.Interface(), fieldMetadata); err != nil {
					if validationErr, ok := err.(*ValidationError); ok {
						errors = append(errors, validationErr)
						if p.stopOnFail {
							return nil, validationErr
						}
					} else {
						return nil, err
					}
				}
			}
		}
	}

	// If there are any errors, return them
	if len(errors) > 0 {
		return nil, fmt.Errorf("validation failed: %v", errors)
	}

	return value, nil
}

// validateField validates a single field
func (p *ValidationPipe) validateField(value interface{}, metadata ValidationMetadata) error {
	for ruleName, ruleValue := range metadata.Rules {
		// Check if we have a custom validator
		if validator, exists := p.validators[ruleName]; exists {
			if err := validator.Validate(value, metadata); err != nil {
				return &ValidationError{
					Field:   metadata.Field,
					Message: err.Error(),
					Value:   value,
					Rule:    ruleName,
				}
			}
			continue
		}

		// Handle built-in validators
		switch ruleName {
		case "required":
			if value == nil || (reflect.TypeOf(value).Kind() == reflect.String && value == "") {
				return &ValidationError{
					Field:   metadata.Field,
					Message: getMessage(metadata.Messages, ruleName, "field is required"),
					Value:   value,
					Rule:    ruleName,
				}
			}

		case "email":
			if str, ok := value.(string); ok {
				if !isValidEmail(str) {
					return &ValidationError{
						Field:   metadata.Field,
						Message: getMessage(metadata.Messages, ruleName, "invalid email format"),
						Value:   value,
						Rule:    ruleName,
					}
				}
			}

		case "min":
			if min, ok := ruleValue.(string); ok {
				if err := validateMin(value, min); err != nil {
					return &ValidationError{
						Field:   metadata.Field,
						Message: getMessage(metadata.Messages, ruleName, err.Error()),
						Value:   value,
						Rule:    ruleName,
					}
				}
			}

		case "max":
			if max, ok := ruleValue.(string); ok {
				if err := validateMax(value, max); err != nil {
					return &ValidationError{
						Field:   metadata.Field,
						Message: getMessage(metadata.Messages, ruleName, err.Error()),
						Value:   value,
						Rule:    ruleName,
					}
				}
			}

		case "regex":
			if pattern, ok := ruleValue.(string); ok {
				if err := validateRegex(value, pattern); err != nil {
					return &ValidationError{
						Field:   metadata.Field,
						Message: getMessage(metadata.Messages, ruleName, err.Error()),
						Value:   value,
						Rule:    ruleName,
					}
				}
			}
		}
	}

	return nil
}

// Helper functions for built-in validators
func isValidEmail(email string) bool {
	// Basic email validation regex
	emailRegex := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	return validateRegex(email, emailRegex) == nil
}

func validateMin(value interface{}, min string) error {
	switch v := value.(type) {
	case string:
		if len(v) < len(min) {
			return fmt.Errorf("string length must be at least %s", min)
		}
	case int:
		minInt := 0
		fmt.Sscanf(min, "%d", &minInt)
		if v < minInt {
			return fmt.Errorf("value must be at least %d", minInt)
		}
	case float64:
		minFloat := 0.0
		fmt.Sscanf(min, "%f", &minFloat)
		if v < minFloat {
			return fmt.Errorf("value must be at least %f", minFloat)
		}
	}
	return nil
}

func validateMax(value interface{}, max string) error {
	switch v := value.(type) {
	case string:
		if len(v) > len(max) {
			return fmt.Errorf("string length must not exceed %s", max)
		}
	case int:
		maxInt := 0
		fmt.Sscanf(max, "%d", &maxInt)
		if v > maxInt {
			return fmt.Errorf("value must not exceed %d", maxInt)
		}
	case float64:
		maxFloat := 0.0
		fmt.Sscanf(max, "%f", &maxFloat)
		if v > maxFloat {
			return fmt.Errorf("value must not exceed %f", maxFloat)
		}
	}
	return nil
}

func validateRegex(value interface{}, pattern string) error {
	if str, ok := value.(string); ok {
		matched := regexp.MustCompile(pattern).MatchString(str)
		if !matched {
			return fmt.Errorf("value does not match pattern %s", pattern)
		}
	}
	return nil
}

func getMessage(messages map[string]string, rule string, defaultMsg string) string {
	if msg, ok := messages[rule]; ok {
		return msg
	}
	return defaultMsg
}
