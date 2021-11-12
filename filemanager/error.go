package filemanager

// Error is the packages error type
type Error struct {
	err     error
	logData map[string]interface{}
}

// Error satisfies the standard library Go error interface
func (e Error) Error() string {
	if e.err == nil {
		return "nil"
	}
	return e.err.Error()
}

// Unwrap implements the standard library Go unwrapper interface
func (e Error) Unwrap() error {
	return e.err
}

// LogData satisfies the dataLogger interface which is used to recover
// log data from an error
func (e Error) LogData() map[string]interface{} {
	return e.logData
}
