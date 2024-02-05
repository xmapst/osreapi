package selfupdate

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"os"
	"testing"

	"github.com/xmapst/osreapi/internal/binarydist"
)

var (
	oldFile         = []byte{0xDE, 0xAD, 0xBE, 0xEF}
	newFile         = []byte{0x01, 0x02, 0x03, 0x04, 0x05, 0x06}
	newFileChecksum = sha256.Sum256(newFile)
)

func cleanup(path string) {
	_ = os.Remove(path)
	_ = os.Remove(fmt.Sprintf(".%s.new", path))
}

// we write with a separate name for each test so that we can run them in parallel
func writeOldFile(path string, t *testing.T) {
	if err := os.WriteFile(path, oldFile, 0777); err != nil {
		t.Fatalf("Failed to write file for testing preparation: %v", err)
	}
}

func validateUpdate(path string, err error, t *testing.T) {
	if err != nil {
		t.Fatalf("Failed to update: %v", err)
	}

	buf, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("Failed to read file post-update: %v", err)
	}

	if !bytes.Equal(buf, newFile) {
		t.Fatalf("File was not updated! Bytes read: %v, Bytes expected: %v", buf, newFile)
	}
}

func TestApplySimple(t *testing.T) {
	fName := "TestApplySimple"
	defer cleanup(fName)
	writeOldFile(fName, t)

	err := Apply(bytes.NewReader(newFile), Options{
		TargetPath: fName,
	})
	validateUpdate(fName, err, t)
}

func TestApplyOldSavePath(t *testing.T) {
	fName := "TestApplyOldSavePath"
	defer cleanup(fName)
	writeOldFile(fName, t)

	oldfName := "OldSavePath"

	err := Apply(bytes.NewReader(newFile), Options{
		TargetPath:  fName,
		OldSavePath: oldfName,
	})
	validateUpdate(fName, err, t)

	if _, err := os.Stat(oldfName); os.IsNotExist(err) {
		t.Fatalf("Failed to find the old file: %v", err)
	}

	cleanup(oldfName)
}

func TestVerifyChecksum(t *testing.T) {
	fName := "TestVerifyChecksum"
	defer cleanup(fName)
	writeOldFile(fName, t)

	err := Apply(bytes.NewReader(newFile), Options{
		TargetPath: fName,
		Checksum:   newFileChecksum[:],
	})
	validateUpdate(fName, err, t)
}

func TestVerifyChecksumNegative(t *testing.T) {
	fName := "TestVerifyChecksumNegative"
	defer cleanup(fName)
	writeOldFile(fName, t)

	badChecksum := []byte{0x0A, 0x0B, 0x0C, 0xFF}
	err := Apply(bytes.NewReader(newFile), Options{
		TargetPath: fName,
		Checksum:   badChecksum,
	})
	if err == nil {
		t.Fatalf("Failed to detect bad checksum!")
	}
}

func TestApplyPatch(t *testing.T) {
	fName := "TestApplyPatch"
	defer cleanup(fName)
	writeOldFile(fName, t)

	patch := new(bytes.Buffer)
	err := binarydist.Diff(bytes.NewReader(oldFile), bytes.NewReader(newFile), patch)
	if err != nil {
		t.Fatalf("Failed to create patch: %v", err)
	}

	err = Apply(patch, Options{
		TargetPath: fName,
		Patcher:    NewBSDiffPatcher(),
	})
	validateUpdate(fName, err, t)
}

func TestCorruptPatch(t *testing.T) {
	fName := "TestCorruptPatch"
	defer cleanup(fName)
	writeOldFile(fName, t)

	badPatch := []byte{0x44, 0x38, 0x86, 0x3c, 0x4f, 0x8d, 0x26, 0x54, 0xb, 0x11, 0xce, 0xfe, 0xc1, 0xc0, 0xf8, 0x31, 0x38, 0xa0, 0x12, 0x1a, 0xa2, 0x57, 0x2a, 0xe1, 0x3a, 0x48, 0x62, 0x40, 0x2b, 0x81, 0x12, 0xb1, 0x21, 0xa5, 0x16, 0xed, 0x73, 0xd6, 0x54, 0x84, 0x29, 0xa6, 0xd6, 0xb2, 0x1b, 0xfb, 0xe6, 0xbe, 0x7b, 0x70}
	err := Apply(bytes.NewReader(badPatch), Options{
		TargetPath: fName,
		Patcher:    NewBSDiffPatcher(),
	})
	if err == nil {
		t.Fatalf("Failed to detect corrupt patch!")
	}
}

func TestVerifyChecksumPatchNegative(t *testing.T) {
	fName := "TestVerifyChecksumPatchNegative"
	defer cleanup(fName)
	writeOldFile(fName, t)

	patch := new(bytes.Buffer)
	anotherFile := []byte{0x77, 0x11, 0x22, 0x33, 0x44, 0x55, 0x66}
	err := binarydist.Diff(bytes.NewReader(oldFile), bytes.NewReader(anotherFile), patch)
	if err != nil {
		t.Fatalf("Failed to create patch: %v", err)
	}

	err = Apply(patch, Options{
		TargetPath: fName,
		Checksum:   newFileChecksum[:],
		Patcher:    NewBSDiffPatcher(),
	})
	if err == nil {
		t.Fatalf("Failed to detect patch to wrong file!")
	}
}
