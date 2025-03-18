package validation

import (
	"fmt"
	"mime/multipart"
	"net/url"
	"regexp"
	"strings"
	"time"
)

// PasswordValidator validates password strength
type PasswordValidator struct {
	minLength        int
	requireUppercase bool
	requireLowercase bool
	requireNumbers   bool
	requireSpecial   bool
}

// NewPasswordValidator creates a new password validator
func NewPasswordValidator(minLength int, requireUppercase, requireLowercase, requireNumbers, requireSpecial bool) *PasswordValidator {
	return &PasswordValidator{
		minLength:        minLength,
		requireUppercase: requireUppercase,
		requireLowercase: requireLowercase,
		requireNumbers:   requireNumbers,
		requireSpecial:   requireSpecial,
	}
}

// Validate implements the Validator interface
func (v *PasswordValidator) Validate(value interface{}, metadata ValidationMetadata) error {
	password, ok := value.(string)
	if !ok {
		return fmt.Errorf("password must be a string")
	}

	if len(password) < v.minLength {
		return fmt.Errorf("password must be at least %d characters long", v.minLength)
	}

	if v.requireUppercase && !regexp.MustCompile(`[A-Z]`).MatchString(password) {
		return fmt.Errorf("password must contain at least one uppercase letter")
	}

	if v.requireLowercase && !regexp.MustCompile(`[a-z]`).MatchString(password) {
		return fmt.Errorf("password must contain at least one lowercase letter")
	}

	if v.requireNumbers && !regexp.MustCompile(`[0-9]`).MatchString(password) {
		return fmt.Errorf("password must contain at least one number")
	}

	if v.requireSpecial && !regexp.MustCompile(`[!@#$%^&*]`).MatchString(password) {
		return fmt.Errorf("password must contain at least one special character")
	}

	return nil
}

// DateValidator validates date formats and ranges
type DateValidator struct {
	format string
	min    time.Time
	max    time.Time
}

// NewDateValidator creates a new date validator
func NewDateValidator(format string, min, max time.Time) *DateValidator {
	return &DateValidator{
		format: format,
		min:    min,
		max:    max,
	}
}

// Validate implements the Validator interface
func (v *DateValidator) Validate(value interface{}, metadata ValidationMetadata) error {
	dateStr, ok := value.(string)
	if !ok {
		return fmt.Errorf("date must be a string")
	}

	date, err := time.Parse(v.format, dateStr)
	if err != nil {
		return fmt.Errorf("invalid date format. Expected format: %s", v.format)
	}

	if !v.min.IsZero() && date.Before(v.min) {
		return fmt.Errorf("date must be after %s", v.min.Format(v.format))
	}

	if !v.max.IsZero() && date.After(v.max) {
		return fmt.Errorf("date must be before %s", v.max.Format(v.format))
	}

	return nil
}

// PhoneNumberValidator validates phone number formats
type PhoneNumberValidator struct {
	countryCode string
	format      string
}

// NewPhoneNumberValidator creates a new phone number validator
func NewPhoneNumberValidator(countryCode, format string) *PhoneNumberValidator {
	return &PhoneNumberValidator{
		countryCode: countryCode,
		format:      format,
	}
}

// Validate implements the Validator interface
func (v *PhoneNumberValidator) Validate(value interface{}, metadata ValidationMetadata) error {
	phone, ok := value.(string)
	if !ok {
		return fmt.Errorf("phone number must be a string")
	}

	// Remove any non-digit characters
	phone = regexp.MustCompile(`[^\d]`).ReplaceAllString(phone, "")

	// Check country code if specified
	if v.countryCode != "" {
		if !strings.HasPrefix(phone, v.countryCode) {
			return fmt.Errorf("phone number must start with country code %s", v.countryCode)
		}
	}

	// Check format if specified
	if v.format != "" {
		matched, err := regexp.MatchString(v.format, phone)
		if err != nil {
			return fmt.Errorf("invalid phone number format: %v", err)
		}
		if !matched {
			return fmt.Errorf("phone number does not match required format")
		}
	}

	return nil
}

// URLValidator validates URL formats
type URLValidator struct {
	requireHTTPS   bool
	allowedDomains []string
}

// NewURLValidator creates a new URL validator
func NewURLValidator(requireHTTPS bool, allowedDomains []string) *URLValidator {
	return &URLValidator{
		requireHTTPS:   requireHTTPS,
		allowedDomains: allowedDomains,
	}
}

// Validate implements the Validator interface
func (v *URLValidator) Validate(value interface{}, metadata ValidationMetadata) error {
	urlStr, ok := value.(string)
	if !ok {
		return fmt.Errorf("URL must be a string")
	}

	url, err := url.Parse(urlStr)
	if err != nil {
		return fmt.Errorf("invalid URL format: %v", err)
	}

	if v.requireHTTPS && url.Scheme != "https" {
		return fmt.Errorf("URL must use HTTPS")
	}

	if len(v.allowedDomains) > 0 {
		allowed := false
		for _, domain := range v.allowedDomains {
			if url.Hostname() == domain || strings.HasSuffix(url.Hostname(), "."+domain) {
				allowed = true
				break
			}
		}
		if !allowed {
			return fmt.Errorf("URL domain not allowed")
		}
	}

	return nil
}

// FileValidator validates file uploads
type FileValidator struct {
	maxSize      int64
	allowedTypes []string
	requireImage bool
}

// NewFileValidator creates a new file validator
func NewFileValidator(maxSize int64, allowedTypes []string, requireImage bool) *FileValidator {
	return &FileValidator{
		maxSize:      maxSize,
		allowedTypes: allowedTypes,
		requireImage: requireImage,
	}
}

// Validate implements the Validator interface
func (v *FileValidator) Validate(value interface{}, metadata ValidationMetadata) error {
	file, ok := value.(*multipart.FileHeader)
	if !ok {
		return fmt.Errorf("value must be a file upload")
	}

	if file.Size > v.maxSize {
		return fmt.Errorf("file size exceeds maximum allowed size of %d bytes", v.maxSize)
	}

	if len(v.allowedTypes) > 0 {
		allowed := false
		for _, t := range v.allowedTypes {
			if strings.HasPrefix(file.Header.Get("Content-Type"), t) {
				allowed = true
				break
			}
		}
		if !allowed {
			return fmt.Errorf("file type not allowed")
		}
	}

	if v.requireImage {
		contentType := file.Header.Get("Content-Type")
		if !strings.HasPrefix(contentType, "image/") {
			return fmt.Errorf("file must be an image")
		}
	}

	return nil
}
