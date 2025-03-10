package enums

// VersioningType defines the different types of API versioning methods
type VersioningType int

const (
	URI VersioningType = iota
	HEADER
	MEDIA_TYPE
	CUSTOM
)

// versioningTypeStrings maps VersioningType values to their string representations
var versioningTypeStrings = map[VersioningType]string{
	URI:        "URI",
	HEADER:     "HEADER",
	MEDIA_TYPE: "MEDIA_TYPE",
	CUSTOM:     "CUSTOM",
}

// String returns the string representation of the VersioningType
func (v VersioningType) String() string {
	if str, exists := versioningTypeStrings[v]; exists {
		return str
	}
	return "UNKNOWN"
}

// IsValid checks if a VersioningType is valid
func (v VersioningType) IsValid() bool {
	_, exists := versioningTypeStrings[v]
	return exists
}
