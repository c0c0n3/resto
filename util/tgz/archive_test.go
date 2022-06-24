package tgz

import (
	"io/fs"
	"os"
	"path"
	"testing"

	"github.com/c0c0n3/resto/util/bytez"
	"github.com/c0c0n3/resto/util/file"
)

func TestWriteFileArchiveErrOnNilSink(t *testing.T) {
	withTempDir(t, func(tempDirPath string) {
		sourceDir := findTestDataDir()
		if err := WriteFileArchive(sourceDir, nil); err == nil {
			t.Errorf("want: error; got: nil")
		}
	})
}

func TestWriteFileArchiveVisitorErr(t *testing.T) {
	withTempDir(t, func(tempDirPath string) {
		sink := bytez.NewBuffer()
		sourceDir, _ := file.NewAbsPathParser().Resolve(tempDirPath)
		os.Chmod(tempDirPath, 0200) // visitor can't scan it

		err := WriteFileArchive(sourceDir, sink)
		if _, ok := err.(*file.VisitError); !ok {
			t.Errorf("want: visit error; got: %v", err)
		}
	})
}

func TestMakeTarballOpenFileErr(t *testing.T) {
	withTempDir(t, func(tempDirPath string) {
		sourceDir := findTestDataDir()
		tarball, _ := file.NewAbsPathParser().Resolve(
			path.Join(tempDirPath, "test.tgz"))
		os.Chmod(tempDirPath, 0400) // can't write tarball to it

		err := MakeTarball(sourceDir, tarball)
		if _, ok := err.(*fs.PathError); !ok {
			t.Errorf("want: visit error; got: %v", err)
		}
	})
}
