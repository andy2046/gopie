package log

import "testing"

func TestNewLogger(t *testing.T) {
	logger := NewLogger(func(c *Config) error {
		c.Level = DEBUG
		c.Prefix = "test:"
		return nil
	})

	logger.Debug("Debug")
	logger.Debugf("format:%s\n", "Debug")

	logger.SetLevel(INFO)

	logger.Info("Info")
	logger.Infof("format:%s\n", "Info")

	logger.Warn("Warn")
	logger.Warnf("format:%s\n", "Warn")

	logger.Error("Error")
	logger.Errorf("format:%s\n", "Error")
}
