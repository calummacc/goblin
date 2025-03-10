package enums

// ShutdownSignal represents system signals that shut down a process
type ShutdownSignal string

const (
	SIGHUP  ShutdownSignal = "SIGHUP"
	SIGINT  ShutdownSignal = "SIGINT"
	SIGQUIT ShutdownSignal = "SIGQUIT"
	SIGILL  ShutdownSignal = "SIGILL"
	SIGTRAP ShutdownSignal = "SIGTRAP"
	SIGABRT ShutdownSignal = "SIGABRT"
	SIGBUS  ShutdownSignal = "SIGBUS"
	SIGFPE  ShutdownSignal = "SIGFPE"
	SIGSEGV ShutdownSignal = "SIGSEGV"
	SIGUSR2 ShutdownSignal = "SIGUSR2"
	SIGTERM ShutdownSignal = "SIGTERM"
)

// signalStrings maps ShutdownSignal values to their string representations
var signalStrings = map[ShutdownSignal]bool{
	SIGHUP:  true,
	SIGINT:  true,
	SIGQUIT: true,
	SIGILL:  true,
	SIGTRAP: true,
	SIGABRT: true,
	SIGBUS:  true,
	SIGFPE:  true,
	SIGSEGV: true,
	SIGUSR2: true,
	SIGTERM: true,
}

// IsValid checks if a shutdown signal is valid
func (s ShutdownSignal) IsValid() bool {
	_, exists := signalStrings[s]
	return exists
}

// String returns the string representation of the ShutdownSignal
func (s ShutdownSignal) String() string {
	return string(s)
}
