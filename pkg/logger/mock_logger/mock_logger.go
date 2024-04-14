package mocklogger

import "fmt"

type MockLogger struct {
}

func New() *MockLogger {
	return &MockLogger{}
}

func (m *MockLogger) Infoln(args ...interface{}) {
	fmt.Println(args...)
}

func (m *MockLogger) Errorln(args ...interface{}) {
	fmt.Println(args...)
}

func (m *MockLogger) Fatal(args ...interface{}) {
	fmt.Println(args...)
}

func (m *MockLogger) Errorf(format string, v ...interface{}) {
	fmt.Printf(format, v...)
}
func (m *MockLogger) Warnf(format string, v ...interface{}) {
	fmt.Printf(format, v...)
}
func (m *MockLogger) Debugf(format string, v ...interface{}) {
	fmt.Printf(format, v...)
}

func (m *MockLogger) Infow(msg string, keysAndValues ...interface{}) {
	k := keysAndValues[:]
	k = append(k, msg)
	fmt.Println(k...)
}
