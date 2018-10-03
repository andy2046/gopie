package deadline

import (
	"errors"
	"testing"
	"time"
)

var (
	errText = "xxx"
)

func TestDeadline(t *testing.T) {
	d := New(10 * time.Millisecond)

	if err := d.Go(fiveMillisecondFunc); err != nil {
		t.Error(err)
	}

	if err := d.Go(twentyMillisecondFunc); err != ErrTimeout {
		t.Error(err)
	}

	if err := d.Go(errorFunc); err.Error() != errText {
		t.Error(err)
	}
}

func fiveMillisecondFunc(<-chan struct{}) error {
	time.Sleep(5 * time.Millisecond)
	return nil
}

func twentyMillisecondFunc(<-chan struct{}) error {
	time.Sleep(20 * time.Millisecond)
	return nil
}

func errorFunc(<-chan struct{}) error {
	return errors.New(errText)
}
