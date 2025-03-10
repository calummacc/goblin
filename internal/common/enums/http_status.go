package enums

// HttpStatus defines standard HTTP status codes
type HttpStatus int

const (
	// 1xx: Informational responses
	CONTINUE            HttpStatus = 100
	SWITCHING_PROTOCOLS HttpStatus = 101
	PROCESSING          HttpStatus = 102
	EARLY_HINTS         HttpStatus = 103

	// 2xx: Successful responses
	OK                            HttpStatus = 200
	CREATED                       HttpStatus = 201
	ACCEPTED                      HttpStatus = 202
	NON_AUTHORITATIVE_INFORMATION HttpStatus = 203
	NO_CONTENT                    HttpStatus = 204
	RESET_CONTENT                 HttpStatus = 205
	PARTIAL_CONTENT               HttpStatus = 206
	MULTI_STATUS                  HttpStatus = 207
	ALREADY_REPORTED              HttpStatus = 208
	CONTENT_DIFFERENT             HttpStatus = 210

	// 3xx: Redirection messages
	AMBIGUOUS          HttpStatus = 300
	MOVED_PERMANENTLY  HttpStatus = 301
	FOUND              HttpStatus = 302
	SEE_OTHER          HttpStatus = 303
	NOT_MODIFIED       HttpStatus = 304
	TEMPORARY_REDIRECT HttpStatus = 307
	PERMANENT_REDIRECT HttpStatus = 308

	// 4xx: Client error responses
	BAD_REQUEST                     HttpStatus = 400
	UNAUTHORIZED                    HttpStatus = 401
	PAYMENT_REQUIRED                HttpStatus = 402
	FORBIDDEN                       HttpStatus = 403
	NOT_FOUND                       HttpStatus = 404
	METHOD_NOT_ALLOWED              HttpStatus = 405
	NOT_ACCEPTABLE                  HttpStatus = 406
	PROXY_AUTHENTICATION_REQUIRED   HttpStatus = 407
	REQUEST_TIMEOUT                 HttpStatus = 408
	CONFLICT                        HttpStatus = 409
	GONE                            HttpStatus = 410
	LENGTH_REQUIRED                 HttpStatus = 411
	PRECONDITION_FAILED             HttpStatus = 412
	PAYLOAD_TOO_LARGE               HttpStatus = 413
	URI_TOO_LONG                    HttpStatus = 414
	UNSUPPORTED_MEDIA_TYPE          HttpStatus = 415
	REQUESTED_RANGE_NOT_SATISFIABLE HttpStatus = 416
	EXPECTATION_FAILED              HttpStatus = 417
	I_AM_A_TEAPOT                   HttpStatus = 418
	MISDIRECTED                     HttpStatus = 421
	UNPROCESSABLE_ENTITY            HttpStatus = 422
	LOCKED                          HttpStatus = 423
	FAILED_DEPENDENCY               HttpStatus = 424
	PRECONDITION_REQUIRED           HttpStatus = 428
	TOO_MANY_REQUESTS               HttpStatus = 429
	UNRECOVERABLE_ERROR             HttpStatus = 456

	// 5xx: Server error responses
	INTERNAL_SERVER_ERROR      HttpStatus = 500
	NOT_IMPLEMENTED            HttpStatus = 501
	BAD_GATEWAY                HttpStatus = 502
	SERVICE_UNAVAILABLE        HttpStatus = 503
	GATEWAY_TIMEOUT            HttpStatus = 504
	HTTP_VERSION_NOT_SUPPORTED HttpStatus = 505
	INSUFFICIENT_STORAGE       HttpStatus = 507
	LOOP_DETECTED              HttpStatus = 508
)

// statusStrings maps HTTP status codes to their respective descriptions
var statusStrings = map[HttpStatus]string{
	CONTINUE:                        "Continue",
	SWITCHING_PROTOCOLS:             "Switching Protocols",
	PROCESSING:                      "Processing",
	EARLY_HINTS:                     "Early Hints",
	OK:                              "OK",
	CREATED:                         "Created",
	ACCEPTED:                        "Accepted",
	NON_AUTHORITATIVE_INFORMATION:   "Non-Authoritative Information",
	NO_CONTENT:                      "No Content",
	RESET_CONTENT:                   "Reset Content",
	PARTIAL_CONTENT:                 "Partial Content",
	MULTI_STATUS:                    "Multi-Status",
	ALREADY_REPORTED:                "Already Reported",
	CONTENT_DIFFERENT:               "Content Different",
	AMBIGUOUS:                       "Ambiguous",
	MOVED_PERMANENTLY:               "Moved Permanently",
	FOUND:                           "Found",
	SEE_OTHER:                       "See Other",
	NOT_MODIFIED:                    "Not Modified",
	TEMPORARY_REDIRECT:              "Temporary Redirect",
	PERMANENT_REDIRECT:              "Permanent Redirect",
	BAD_REQUEST:                     "Bad Request",
	UNAUTHORIZED:                    "Unauthorized",
	PAYMENT_REQUIRED:                "Payment Required",
	FORBIDDEN:                       "Forbidden",
	NOT_FOUND:                       "Not Found",
	METHOD_NOT_ALLOWED:              "Method Not Allowed",
	NOT_ACCEPTABLE:                  "Not Acceptable",
	PROXY_AUTHENTICATION_REQUIRED:   "Proxy Authentication Required",
	REQUEST_TIMEOUT:                 "Request Timeout",
	CONFLICT:                        "Conflict",
	GONE:                            "Gone",
	LENGTH_REQUIRED:                 "Length Required",
	PRECONDITION_FAILED:             "Precondition Failed",
	PAYLOAD_TOO_LARGE:               "Payload Too Large",
	URI_TOO_LONG:                    "URI Too Long",
	UNSUPPORTED_MEDIA_TYPE:          "Unsupported Media Type",
	REQUESTED_RANGE_NOT_SATISFIABLE: "Requested Range Not Satisfiable",
	EXPECTATION_FAILED:              "Expectation Failed",
	I_AM_A_TEAPOT:                   "I'm a Teapot",
	MISDIRECTED:                     "Misdirected",
	UNPROCESSABLE_ENTITY:            "Unprocessable Entity",
	LOCKED:                          "Locked",
	FAILED_DEPENDENCY:               "Failed Dependency",
	PRECONDITION_REQUIRED:           "Precondition Required",
	TOO_MANY_REQUESTS:               "Too Many Requests",
	UNRECOVERABLE_ERROR:             "Unrecoverable Error",
	INTERNAL_SERVER_ERROR:           "Internal Server Error",
	NOT_IMPLEMENTED:                 "Not Implemented",
	BAD_GATEWAY:                     "Bad Gateway",
	SERVICE_UNAVAILABLE:             "Service Unavailable",
	GATEWAY_TIMEOUT:                 "Gateway Timeout",
	HTTP_VERSION_NOT_SUPPORTED:      "HTTP Version Not Supported",
	INSUFFICIENT_STORAGE:            "Insufficient Storage",
	LOOP_DETECTED:                   "Loop Detected",
}

// String returns the string representation of an HTTP status code
func (s HttpStatus) String() string {
	if str := statusStrings[s]; str != "" {
		return str
	}
	return "Unknown Status"
}

// IsValid checks whether an HTTP status code is valid
func (s HttpStatus) IsValid() bool {
	_, exists := statusStrings[s]
	return exists
}
