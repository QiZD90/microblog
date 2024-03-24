package logger

type logStub struct{}

func (log logStub) Errorf(format string, args ...interface{}) {}
func (log logStub) Infof(format string, args ...interface{})  {}
func (log logStub) WithFields(fields map[string]any) Logger   { return log }
