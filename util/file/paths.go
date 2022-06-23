package file

import (
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// AbsPath represents an absolute path to a filesystem node.
type AbsPath struct{ data string }

// Get the string representation of this AbsPath. That's the canonical
// string path Go's standard lib deals with---e.g. '/some/where'.
func (d AbsPath) Value() string {
	return d.data
}

func IsStringPath(value interface{}) error {
	s, _ := value.(string)
	_, err := NewAbsPathParser().Resolve(s)
	return err
}

// AbsPathResolver encapsulates the process of turning a, possibly relative,
// string path into an AbsPath.
type AbsPathResolver interface {
	// Figure out what's the absolute path of the given string path.
	Resolve(path string) (AbsPath, error)
}

type pathParser struct {
	makeAbs func(path string) (string, error)
}

// Make an AbsPathParser.
func NewAbsPathParser() AbsPathResolver {
	return &pathParser{
		makeAbs: filepath.Abs,
	}
}

func (p *pathParser) Resolve(path string) (AbsPath, error) {
	path = strings.TrimSpace(path) // (*)
	if len(path) == 0 {
		return AbsPath{}, pathResolveWhitespaceErr()
	}
	if p, err := p.makeAbs(path); err != nil {
		return AbsPath{}, pathResolveAbsErr(err)
	} else {
		return AbsPath{data: p}, nil
	}

	// (*) filepath.Abs doesn't trim space, e.g.
	// filepath.Abs('/a/b ') == '/a/b '.
}

// Append the given relative path to this absolute path.
func (d AbsPath) Join(relativePath string) AbsPath {
	rest := strings.TrimSpace(relativePath) // (1)
	return AbsPath{
		data: filepath.Join(d.Value(), rest), // (2)
	}

	// (1) Join doesn't trim space, e.g. Join("/a", "/b ") == "/a/b "
	// (2) In principle this is wrong since we don't know if relativePath
	// is a valid path according to the FS we're running on. (Join doesn't
	// check that.) So we could potentially return an inconsistent AbsPath.
	// Go's standard lib is quite weak in the handling of abstract paths,
	// i.e. independent of OS, so this is the best we can do. See e.g.
	// - https://stackoverflow.com/questions/35231846
}

// Is this a path to a directory?
// If there's an error accessing the file system return that error.
// Otherwise return nil if the path points to a directory or a NotADirectory
// error if it doesn't.
func (d AbsPath) IsDir() error {
	if f, err := os.Stat(d.Value()); err != nil {
		return err
	} else {
		if !f.IsDir() {
			return notADirErr(d)
		}
	}
	return nil
}

// ListPaths collects, recursively, the paths of all the directories and
// files inside dirPath. Each collected path is relative to dirPath, so
// for example, if dirPath = "b" and f is a file at "b/d/f", then "d/f"
// gets returned. ListPaths sorts the returned paths in alphabetical
// order.
func ListPaths(dirPath string) ([]string, []error) {
	visitedPaths := []string{}
	errs := []error{}

	targetDir, err := NewAbsPathParser().Resolve(dirPath)
	if err != nil {
		errs = append(errs, err)
		return visitedPaths, errs
	}

	scanner := NewTreeScanner(targetDir)
	es := scanner.Visit(func(node TreeNode) error {
		if node.RelPath != "" {
			visitedPaths = append(visitedPaths, node.RelPath)
		}
		return nil
	})

	errs = append(errs, es...)
	sort.Strings(visitedPaths)

	return visitedPaths, errs
}

// ListSubDirectoryNames returns the names of any directory found just
// below dirPath. ListSubDirectoryNames sorts the returned names in
// alphabetical order. Also, ListSubDirectoryNames will return an empty
// list if an error happens.
func ListSubDirectoryNames(dirPath string) ([]string, error) {
	dirs := []string{}
	entries, err := os.ReadDir(dirPath)
	if err != nil {
		return dirs, err
	}

	for _, e := range entries {
		if e.IsDir() {
			dirs = append(dirs, e.Name())
		}
	}

	sort.Strings(dirs)

	return dirs, nil
}
