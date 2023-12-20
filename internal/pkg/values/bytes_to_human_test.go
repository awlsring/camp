package values

import "testing"

func Test_BytesToHumanReadable_MB(t *testing.T) {
	var b uint64 = 1024 * 1024
	human := BytesToHumanReadable(b)
	if human != "1 MB" {
		t.Errorf("Expected 1.0 MB, got %s", human)
	}
}

func Test_BytesToHumanReadable_MB_Float(t *testing.T) {
	var b uint64 = 1024 * 1024 * 3.5
	human := BytesToHumanReadable(b)
	if human != "3.5 MB" {
		t.Errorf("Expected 3.5 MB, got %s", human)
	}
}

func Test_BytesToHumanReadable_GB(t *testing.T) {
	var b uint64 = 1024 * 1024 * 1024
	human := BytesToHumanReadable(b)
	if human != "1 GB" {
		t.Errorf("Expected 1 GB, got %s", human)
	}
}

func Test_BytesToHumanReadable_TB(t *testing.T) {
	var b uint64 = 1024 * 1024 * 1024 * 1024
	human := BytesToHumanReadable(b)
	if human != "1 TB" {
		t.Errorf("Expected 1.0 TB, got %s", human)
	}
}
