package util

import (
	"os"
	"testing"
	"time"
)

func TestGetTodayInKST(t *testing.T) {
	// Arrange
	expected := time.Now().In(kstLocation)
	// Act
	actual := GetTodayInKST()
	// Assert
	if expected.Year() != actual.Year() || expected.Month() != actual.Month() || expected.Day() != actual.Day() {
		t.Errorf("expected: %v, actual: %v", expected, actual)
	}
}

func TestParseDate_ValidDate(t *testing.T) {
	testCases := []struct {
		name     string
		dateStr  string
		expected time.Time
	}{
		{"ValidDate: 2021-09-01", "2021-09-01", time.Date(2021, 9, 1, 0, 0, 0, 0, kstLocation)},
		{"ValidDate: 2024-09-11", "2024-09-11", time.Date(2024, 9, 11, 0, 0, 0, 0, kstLocation)},
		{"ValidDate: 2024-11-22", "2024-11-22", time.Date(2024, 11, 22, 0, 0, 0, 0, kstLocation)},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Act
			actual, err := ParseDate(tc.dateStr)
			// Assert
			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}
			if tc.expected.Year() != actual.Year() || tc.expected.Month() != actual.Month() || tc.expected.Day() != actual.Day() {
				t.Errorf("expected: %v, actual: %v", tc.expected, actual)
			}
		})
	}
}

func TestParseDate_InvalidDate(t *testing.T) {
	testCases := []struct {
		name    string
		dateStr string
	}{
		{"InvalidDate: 2021-09-31", "2021-09-31"},
		{"InvalidDate: 2024-11-31", "2024-11-31"},
		{"InvalidDate: 2024-02-30", "2024-02-30"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Act
			_, err := ParseDate(tc.dateStr)
			// Assert
			if err == nil {
				t.Errorf("expected error, but nil")
			}
		})
	}
}

func TestLoadFile_ValidFile(t *testing.T) {
	// Arrange
	tempFile, err := os.CreateTemp("", "testFile-*.txt")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	defer func(name string) {
		err := os.Remove(name)
		if err != nil {
			t.Fatalf("failed to remove temp file: %v", err)
		}
	}(tempFile.Name())

	content := []byte("Hello, World!")
	if _, err := tempFile.Write(content); err != nil {
		t.Fatalf("failed to write to temp file: %v", err)
	}

	if err := tempFile.Close(); err != nil {
		t.Fatalf("failed to close temp file: %v", err)
	}

	// Act
	readContent, err := LoadFile(tempFile.Name())
	if err != nil {
		t.Fatalf("LoadFile failed: %v", err)
	}

	// Assert
	if string(readContent) != string(content) {
		t.Errorf("expected %s but got %s", string(content), string(readContent))
	}
}

func TestLoadFile_InvalidFile(t *testing.T) {
	// Act
	_, err := LoadFile("invalid-file-path")

	// Assert
	if err == nil {
		t.Errorf("expected error, but nil")
	}
}

func TestPickRandomNumber(t *testing.T) {
	// Arrange
	upperBound := 10
	// Act
	actual := PickRandomNumber(upperBound)
	// Assert
	if actual < 0 || actual >= upperBound {
		t.Errorf("expected 0 <= actual < %d, but got %d", upperBound, actual)
	}
}
