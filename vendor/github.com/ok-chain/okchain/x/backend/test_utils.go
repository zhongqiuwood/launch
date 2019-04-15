package backend

type TestLogger struct {
}

func (l *TestLogger) Debug(msg string, keyvals ...interface{}) {}
func (l *TestLogger) Info(msg string, keyvals ...interface{})  {}
func (l *TestLogger) Error(msg string, keyvals ...interface{}) {}
func (l *TestLogger) With(keyvals ...interface{})              {}
