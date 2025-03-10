package enums

// RequestMethod defines HTTP request methods
type RequestMethod int

const (
	GET RequestMethod = iota
	POST
	PUT
	DELETE
	PATCH
	ALL
	OPTIONS
	HEAD
	SEARCH
	PROPFIND
	PROPPATCH
	MKCOL
	COPY
	MOVE
	LOCK
	UNLOCK
)

// methodStrings maps RequestMethod values to their string representations
var methodStrings = map[RequestMethod]string{
	GET:       "GET",
	POST:      "POST",
	PUT:       "PUT",
	DELETE:    "DELETE",
	PATCH:     "PATCH",
	ALL:       "ALL",
	OPTIONS:   "OPTIONS",
	HEAD:      "HEAD",
	SEARCH:    "SEARCH",
	PROPFIND:  "PROPFIND",
	PROPPATCH: "PROPPATCH",
	MKCOL:     "MKCOL",
	COPY:      "COPY",
	MOVE:      "MOVE",
	LOCK:      "LOCK",
	UNLOCK:    "UNLOCK",
}

// String returns the string representation of the RequestMethod
func (m RequestMethod) String() string {
	if str, exists := methodStrings[m]; exists {
		return str
	}
	return "UNKNOWN"
}

// IsValid checks if the request method is valid
func (m RequestMethod) IsValid() bool {
	_, exists := methodStrings[m]
	return exists
}
