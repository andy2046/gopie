package fileutil_test

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
	"time"

	. "github.com/andy2046/gopie/pkg/fileutil"
)

func TestIsDirWriteable(t *testing.T) {
	tmpdir, err := ioutil.TempDir("", "foo")
	if err != nil {
		t.Fatalf("ioutil.TempDir error %v", err)
	}
	defer os.RemoveAll(tmpdir)
	if err = IsDirWriteable(tmpdir); err != nil {
		t.Fatalf("IsDirWriteable error %v", err)
	}
}

func TestCreateDirAll(t *testing.T) {
	tmpdir, err := ioutil.TempDir(os.TempDir(), "foo")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpdir)

	tmpdir2 := filepath.Join(tmpdir, "bar")
	if err = CreateDirAll(tmpdir2); err != nil {
		t.Fatal(err)
	}

	if err = ioutil.WriteFile(filepath.Join(tmpdir2, "test.txt"),
		[]byte("test"), PrivateFileMode); err != nil {
		t.Fatal(err)
	}

	if err = CreateDirAll(tmpdir2); err == nil || !strings.Contains(err.Error(), "to be empty, got") {
		t.Fatalf("CreateDirAll error %v", err)
	}
}

func TestExist(t *testing.T) {
	fdir := filepath.Join(os.TempDir(), fmt.Sprint(time.Now().UnixNano()))
	os.RemoveAll(fdir)
	if err := os.Mkdir(fdir, 0666); err != nil {
		t.Skip(err)
	}
	defer os.RemoveAll(fdir)
	if !Exist(fdir) {
		t.Fatal("Exist expected to be true")
	}

	f, err := ioutil.TempFile(os.TempDir(), "fileutil")
	if err != nil {
		t.Skip(err)
	}
	f.Close()

	if !Exist(f.Name()) {
		t.Error("Exist expected to be true")
	}

	os.Remove(f.Name())
	if Exist(f.Name()) {
		t.Error("Exist expected to be false")
	}
}

func TestReadDir(t *testing.T) {
	tmpdir, err := ioutil.TempDir("", "")
	defer os.RemoveAll(tmpdir)
	if err != nil {
		t.Fatalf("ioutil.TempDir error %v", err)
	}

	files := []string{"def", "abc", "xyz", "ghi"}
	for _, f := range files {
		fh, err := os.Create(filepath.Join(tmpdir, f))
		if err != nil {
			t.Skip(err)
		}
		if err = fh.Close(); err != nil {
			t.Skip(err)
		}
	}
	fs, err := ReadDir(tmpdir)
	if err != nil {
		t.Fatalf("ReadDir error %v", err)
	}
	wfs := []string{"abc", "def", "ghi", "xyz"}
	if !reflect.DeepEqual(fs, wfs) {
		t.Fatalf("ReadDir error got %v, want %v", fs, wfs)
	}

	files = []string{"def.wal", "abc.wal", "xyz.wal", "ghi.wal"}
	for _, f := range files {
		fh, err := os.Create(filepath.Join(tmpdir, f))
		if err != nil {
			t.Skip(err)
		}
		if err = fh.Close(); err != nil {
			t.Skip(err)
		}
	}
	fs, err = ReadDir(tmpdir, WithExt(".wal"))
	if err != nil {
		t.Fatalf("ReadDir error %v", err)
	}
	wfs = []string{"abc.wal", "def.wal", "ghi.wal", "xyz.wal"}
	if !reflect.DeepEqual(fs, wfs) {
		t.Fatalf("ReadDir error got %v, want %v", fs, wfs)
	}
}
