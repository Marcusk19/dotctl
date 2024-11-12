package test

import (
	"os"
	"testing"

	"github.com/Marcusk19/dotctl/tools"
	"github.com/spf13/afero"
)

func init() {
	tools.SetTestFs()
}

func TestCopyFile(t *testing.T) {
	fs := afero.NewMemMapFs()

	fs.MkdirAll("test/src", 0755)
	fs.MkdirAll("test/dest", 0755)
	err := afero.WriteFile(fs, "test/src/a.txt", []byte("file a"), 0644)
	if err != nil {
		t.Errorf("problem creating source file: %s", err.Error())
	}

	err = tools.CopyFile(fs, "test/src/a.txt", "test/dest/a.txt")
	if err != nil {
		t.Error(err.Error())
	}
	_, err = fs.Stat("test/dest/a.txt")
	if os.IsNotExist(err) {
		t.Errorf("expected destination file does not exist")
	}

	result, err := afero.ReadFile(fs, "test/dest/a.txt")
	if err != nil {
		t.Error(err.Error())
	}

	if string(result) != "file a" {
		t.Errorf("expected 'file a' got '%s'", string(result))
	}

}

func TestCopyDir(t *testing.T) {
	fs := afero.NewMemMapFs()

	fs.MkdirAll("test/src/dirA", 0755)
	fs.MkdirAll("test/dest/", 0755)
	fs.Mkdir("test/src/dirA/dirB", 0755)

	err := afero.WriteFile(fs, "test/src/dirA/a.txt", []byte("file a"), 0644)
	if err != nil {
		t.Error(err.Error())
	}
	err = afero.WriteFile(fs, "test/src/dirA/dirB/b.txt", []byte("file b"), 0644)
	if err != nil {
		t.Error(err.Error())
	}

	err = tools.CopyDir(fs, "test/src", "test/dest")
	if err != nil {
		t.Error(err.Error())
	}

	result, err := afero.ReadFile(fs, "test/dest/dirA/a.txt")
	if err != nil {
		t.Error(err.Error())
	}

	if string(result) != "file a" {
		t.Errorf("expected 'file a' got '%s'", string(result))
	}

	result, err = afero.ReadFile(fs, "test/dest/dirA/dirB/b.txt")
	if err != nil {
		t.Error(err.Error())
	}

	if string(result) != "file b" {
		t.Errorf("expected 'file b' got '%s'", string(result))
	}

}
