package enums

// RouteParamTypes defines parameter types for route handlers
type RouteParamTypes int

const (
	REQUEST RouteParamTypes = iota
	RESPONSE
	NEXT
	BODY
	QUERY
	PARAM
	HEADERS
	SESSION
	FILE
	FILES
	HOST
	IP
	RAW_BODY
)

// paramTypeStrings maps RouteParamTypes values to their string representations
var paramTypeStrings = map[RouteParamTypes]string{
	REQUEST:  "REQUEST",
	RESPONSE: "RESPONSE",
	NEXT:     "NEXT",
	BODY:     "BODY",
	QUERY:    "QUERY",
	PARAM:    "PARAM",
	HEADERS:  "HEADERS",
	SESSION:  "SESSION",
	FILE:     "FILE",
	FILES:    "FILES",
	HOST:     "HOST",
	IP:       "IP",
	RAW_BODY: "RAW_BODY",
}

// String returns the string representation of the RouteParamTypes
func (p RouteParamTypes) String() string {
	if str, exists := paramTypeStrings[p]; exists {
		return str
	}
	return "UNKNOWN"
}

// IsValid checks if the route parameter type is valid
func (p RouteParamTypes) IsValid() bool {
	_, exists := paramTypeStrings[p]
	return exists
}
