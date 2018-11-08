package main

import (
	"flag"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
)

type PackageDep struct {
	m_gopaths       []string
	m_curfilepaths  []string
	m_lastfilepaths []string
	m_curfullname   string
	m_lastfullnames []string
	m_searched      []string
	m_needsearch    []string
	m_notarchs      []string
	m_notos         []string
}

var allos []string
var allarchs []string
var defpkgs []string

func init() {
	allos = []string{}
	allarchs = []string{}

	allos = append(allos, "linux")
	allos = append(allos, "darwin")
	allos = append(allos, "windows")
	allos = append(allos, "freebsd")
	allos = append(allos, "netbsd")
	allos = append(allos, "plan9")
	allos = append(allos, "openbsd")
	allos = append(allos, "nacl")
	allos = append(allos, "unix")
	allos = append(allos, "test")

	allarchs = append(allarchs, "386")
	allarchs = append(allarchs, "amd64")
	allarchs = append(allarchs, "amd64p32")
	allarchs = append(allarchs, "arm")
	allarchs = append(allarchs, "armx")
	allarchs = append(allarchs, "arm_gen")
	allarchs = append(allarchs, "arm64")
	allarchs = append(allarchs, "arm64p32")
	allarchs = append(allarchs, "ppc64")
	allarchs = append(allarchs, "ppc")
	allarchs = append(allarchs, "ppc64le")

	defpkgs = append(defpkgs, "runtime")
	defpkgs = append(defpkgs, "os")
	defpkgs = append(defpkgs, "fmt")
	defpkgs = append(defpkgs, "C")
	defpkgs = append(defpkgs, "a")
}

func NewPackageDep() *PackageDep {
	p := &PackageDep{}
	p.m_gopaths = []string{}
	p.m_curfilepaths = []string{}
	p.m_lastfilepaths = []string{}
	p.m_needsearch = []string{}
	p.m_searched = []string{}
	p.m_notarchs = []string{}
	p.m_notos = []string{}
	p.m_curfullname = ""
	p.m_lastfullnames = []string{}
	return p
}

func (pkg *PackageDep) push_inner_current_dir(fname string) error {
	if len(pkg.m_curfilepaths) > 0 {
		for _, f := range pkg.m_curfilepaths {
			pkg.m_lastfilepaths = append(pkg.m_lastfilepaths, f)
		}
	}

	ff, err2 := filepath.Abs(fname)
	pkg.m_curfilepaths = []string{}
	if err2 == nil {
		pkg.m_curfilepaths = append(pkg.m_curfilepaths, filepath.Dir(ff))
	}
	return nil
}

func (pkg *PackageDep) pop_inner_current_dir() string {
	var currentdir string
	currentdir = ""
	if len(pkg.m_curfilepaths) > 0 {
		currentdir = pkg.m_curfilepaths[0]
	}
	pkg.m_curfilepaths = []string{}
	if len(pkg.m_lastfilepaths) > 0 {
		pkg.m_curfilepaths = append(pkg.m_curfilepaths, pkg.m_lastfilepaths[(len(pkg.m_lastfilepaths)-1)])
		pkg.m_lastfilepaths = pkg.m_lastfilepaths[:(len(pkg.m_lastfilepaths) - 1)]
	}
	return currentdir
}

func (pkg *PackageDep) push_curfile(fullname string) bool {
	ff, err2 := filepath.Abs(fullname)
	if err2 != nil {
		return false
	}

	if len(pkg.m_curfullname) > 0 {
		pkg.m_lastfullnames = append(pkg.m_lastfullnames, pkg.m_curfullname)
	}
	pkg.m_curfullname = ff
	return true
}

func (pkg *PackageDep) pop_curfile() string {
	retfull := pkg.m_curfullname
	if len(pkg.m_lastfullnames) > 0 {
		pkg.m_curfullname = pkg.m_lastfullnames[(len(pkg.m_lastfullnames) - 1)]
		pkg.m_lastfullnames = pkg.m_lastfullnames[:(len(pkg.m_lastfullnames) - 1)]
	}
	return retfull
}

func (pkg *PackageDep) get_imports_inner(fname string) (imports []string, err error) {
	var fileast *ast.File
	var fset *token.FileSet
	var curpkg string
	imports = []string{}
	err = nil

	fset = token.NewFileSet()
	fileast, err = parser.ParseFile(fset, fname, nil, parser.ImportsOnly)
	if err != nil {
		return
	}

	for _, imps := range fileast.Imports {
		curpkg, err = strconv.Unquote(imps.Path.Value)
		if err == nil {
			imports = append(imports, curpkg)
		} else {
			imports = append(imports, imps.Path.Value)
		}
	}
	err = nil
	//Debug("%s imports %s", fname, imports)
	return
}

func (pkg *PackageDep) find_in_searched(importdir string) int {
	for i, f := range pkg.m_searched {
		if f == importdir {
			return i
		}
	}
	return -1
}

func (pkg *PackageDep) is_dir_local(importdir string) bool {
	if strings.HasPrefix(importdir, "./") ||
		strings.HasPrefix(importdir, "../") ||
		importdir == "." || importdir == ".." {
		return true
	}
	return false
}

func (pkg *PackageDep) insert_in_not_searched(importdir string) int {
	for i, f := range pkg.m_needsearch {
		if f == importdir {
			return i
		}
	}

	if pkg.is_dir_local(importdir) {
		bname := filepath.Dir(pkg.m_curfullname)
		ff, err2 := filepath.Abs(path.Join(bname, importdir))
		if err2 != nil {
			Error("[%s][%s] not set", bname, importdir)
			return 0
		}
		Debug("[%s]add [%s]", pkg.m_curfullname, ff)
		pkg.m_needsearch = append(pkg.m_needsearch, ff)
		return len(pkg.m_needsearch)
	}

	Debug("[%s] insert [%s]", pkg.m_curfullname, importdir)
	pkg.m_needsearch = append(pkg.m_needsearch, importdir)
	return len(pkg.m_needsearch)
}

func (pkg *PackageDep) set_imports_inner(imports []string) {
	for _, imps := range imports {
		if pkg.find_in_searched(imps) >= 0 {
			continue
		}

		pkg.insert_in_not_searched(imps)
	}
}

func (pkg *PackageDep) handle_one_file(fname string) error {
	var err error
	var imports []string
	imports, err = pkg.get_imports_inner(fname)
	if err != nil {
		return err
	}

	pkg.set_imports_inner(imports)
	return nil
}

func (pkg *PackageDep) search_dir(path string) error {
	var err error
	var files []os.FileInfo
	var finfo os.FileInfo
	var archname string
	var osname string
	var tmpfiltername string
	var tmposname string
	var filterd int
	var fullname string
	var added bool
	var curadded bool
	added = false
	curadded = false
	files, err = ioutil.ReadDir(path)
	if err != nil {
		return err
	}

	defer func() {
		if added {
			pkg.pop_inner_current_dir()
			added = false
		}
		if curadded {
			pkg.pop_curfile()
			curadded = false
		}
	}()

	for _, finfo = range files {
		filterd = 0
		//Debug("finfo %s", finfo.Name())
		if !finfo.IsDir() && strings.HasSuffix(finfo.Name(), ".go") {
			for _, archname = range pkg.m_notarchs {
				if filterd > 0 {
					break
				}

				for _, osname = range pkg.m_notos {
					tmpfiltername = fmt.Sprintf("_%s_%s.go", osname, archname)
					tmposname = fmt.Sprintf("_%s.go", osname)
					if strings.HasSuffix(finfo.Name(), tmpfiltername) || strings.HasSuffix(finfo.Name(), tmposname) {
						filterd = 1
						break
					}
				}
			}

			if filterd == 0 {
				fullname = fmt.Sprintf("%s%c%s", path, filepath.Separator, finfo.Name())
				//Debug("fullname [%s]", fullname)
				err2 := pkg.push_inner_current_dir(fullname)
				if err2 == nil {
					added = true
				}
				curadded = pkg.push_curfile(fullname)
				err = pkg.handle_one_file(fullname)
				if err != nil {
					return err
				}
			}
		} else if finfo.IsDir() {
			fullname = fmt.Sprintf("%s%c%s", path, filepath.Separator, finfo.Name())
			err2 := pkg.push_inner_current_dir(fullname)
			if err2 == nil {
				added = true
			}
			curadded = pkg.push_curfile(fullname)
			err = pkg.search_dir(fullname)
			if err != nil {
				return err
			}
		}
		if added {
			pkg.pop_inner_current_dir()
			added = false
		}

		if curadded {
			pkg.pop_curfile()
			curadded = false
		}
	}

	return nil
}

func IsAbsPath(fname string) bool {
	ff, err2 := filepath.Abs(fname)
	if err2 != nil {
		return false
	}

	if ff == fname {
		return true
	}

	return false
}

func (pkg *PackageDep) handle_one_package(pkgname string) error {
	var err error
	var pathname, curdir string
	var sarr []string
	var fullpath, findpath string
	var finded int
	var dinfo os.FileInfo
	var sep string

	sarr = strings.Split(pkgname, "/")
	sep = fmt.Sprintf("%c", os.PathSeparator)
	pathname = strings.Join(sarr, sep)

	finded = 0
	findpath = ""
	for _, curdir = range pkg.m_gopaths {
		fullpath = fmt.Sprintf("%s%c%s", curdir, filepath.Separator, pathname)
		dinfo, err = os.Stat(fullpath)
		if err == nil && dinfo.IsDir() {
			finded = 1
			findpath = fullpath
			Debug("findpath [%s] for [%s]", findpath, pathname)
			break
		}
	}

	if finded == 0 {
		if IsAbsPath(pkgname) {
			finded = 1
			findpath = pkgname
			Debug("findpath [%s]", findpath)
		}
		if finded == 0 {
			err = fmt.Errorf("[%s]%s not find in gopaths", pkg.m_curfullname, pkgname)
			return err
		}
	}

	err = pkg.search_dir(findpath)
	if err != nil {
		return err
	}
	return nil
}

func (pkg *PackageDep) SearchDir(path string) error {
	var err error
	var curpkg string
	var pkgnames []string
	var errmsg string
	var added bool
	added = false

	err = pkg.search_dir(path)
	if err != nil {
		return err
	}

	defer func() {
		if added {
			pkg.pop_curfile()
			added = false
		}
	}()

	for len(pkg.m_needsearch) > 0 {
		curpkg = pkg.m_needsearch[0]
		added = pkg.push_curfile(curpkg)
		//Debug("curpkg [%s]", curpkg)
		err = pkg.handle_one_package(curpkg)
		if err != nil {
			return err
		}

		if added {
			pkg.pop_curfile()
			added = false
		}

		if curpkg != pkg.m_needsearch[0] {
			errmsg = fmt.Sprintf("%s != needsearch[0] %s", curpkg, pkg.m_needsearch[0])
			panic(errmsg)
		}

		/*we put into the search */
		pkgnames = pkg.m_needsearch[1:]
		pkg.m_needsearch = pkgnames
		pkg.m_searched = append(pkg.m_searched, curpkg)
	}
	return nil
}

func (pkg *PackageDep) SetGoPath(paths []string) int {
	var finded int
	var cnt int
	cnt = 0
	for _, curpath := range paths {
		finded = 0
		for _, spath := range pkg.m_gopaths {
			if spath == curpath {
				finded = 1
				break
			}
		}

		if finded == 0 {
			pkg.m_gopaths = append(pkg.m_gopaths, curpath)
			cnt++
		}
	}

	Info("gopaths %s", pkg.m_gopaths)
	return cnt
}

func (pkg *PackageDep) SetSearch(paths []string) int {
	var finded int
	var cnt int
	cnt = 0
	for _, curpath := range paths {
		finded = 0
		for _, spath := range pkg.m_searched {
			if spath == curpath {
				finded = 1
				break
			}
		}

		if finded == 0 {
			pkg.m_searched = append(pkg.m_searched, curpath)
			cnt++
		}
	}
	Info("searched %s", pkg.m_searched)
	return cnt
}

func (pkg *PackageDep) SetFilterOS(oss []string) int {
	var finded int
	var cnt int
	cnt = 0
	for _, curos := range oss {
		finded = 0
		for _, sos := range pkg.m_notos {
			if sos == curos {
				finded = 1
				break
			}
		}

		if finded == 0 {
			pkg.m_notos = append(pkg.m_notos, curos)
			cnt++
		}
	}
	Info("filteros %s", pkg.m_notos)
	return cnt
}

func (pkg *PackageDep) SetFilterArch(archs []string) int {
	var finded int
	var cnt int
	cnt = 0
	for _, curarch := range archs {
		finded = 0
		for _, sarch := range pkg.m_notarchs {
			if sarch == curarch {
				finded = 1
				break
			}
		}

		if finded == 0 {
			pkg.m_notarchs = append(pkg.m_notarchs, curarch)
			cnt++
		}
	}
	Info("filterarch %s", pkg.m_notarchs)
	return cnt
}

type ArrayVar struct {
	m_arr []string
}

func NewArrayVar() *ArrayVar {
	p := &ArrayVar{}
	p.m_arr = []string{}
	return p
}

func (arr *ArrayVar) String() string {
	var s string
	s = "["

	for i, curv := range arr.m_arr {
		if i != 0 {
			s += ","
		}
		s += fmt.Sprintf("%s", curv)
	}

	s += "]"
	return s
}

func (arr *ArrayVar) Set(val string) error {
	for _, curv := range arr.m_arr {
		if curv == val {
			return nil
		}
	}

	arr.m_arr = append(arr.m_arr, val)
	return nil
}

func main() {
	var filters []string
	var finded int
	suparchs := NewArrayVar()
	supos := NewArrayVar()
	filterpkg := NewArrayVar()
	gopaths := NewArrayVar()
	pkg := NewPackageDep()

	paths := os.Getenv("GOPATH")
	if len(paths) > 0 {
		sep := fmt.Sprintf("%c", os.PathListSeparator)
		patharr := strings.Split(paths, sep)
		for _, curdir := range patharr {
			gopaths.Set(curdir)
			gopaths.Set(fmt.Sprintf("%s%cvendor", curdir, os.PathSeparator))
			gopaths.Set(fmt.Sprintf("%s%csrc", curdir, os.PathSeparator))
			gopaths.Set(fmt.Sprintf("%s%csrc%cvendor", curdir, os.PathSeparator, os.PathSeparator))
		}
	}
	suparchs.Set(runtime.GOARCH)
	supos.Set(runtime.GOOS)

	for _, curpkg := range defpkgs {
		filterpkg.Set(curpkg)
	}

	flag.Var(suparchs, "arch", "add support archs")
	flag.Var(suparchs, "A", "add support archs")
	flag.Var(supos, "O", "add support os")
	flag.Var(supos, "os", "add support os")
	flag.Var(filterpkg, "P", "add filter package")
	flag.Var(filterpkg, "pkg", "add filter package")
	flag.Var(gopaths, "p", "add gopaths to search")
	flag.Var(gopaths, "path", "add gopaths to search")

	flag.Parse()

	pkg.SetGoPath(gopaths.m_arr)
	pkg.SetSearch(filterpkg.m_arr)
	filters = []string{}

	for _, curos := range allos {
		finded = 0
		for _, notos := range supos.m_arr {
			if curos == notos {
				finded = 1
				break
			}
		}

		if finded == 0 {
			filters = append(filters, curos)
		}
	}

	pkg.SetFilterOS(filters)

	filters = []string{}
	for _, curarch := range allarchs {
		finded = 0
		for _, notarch := range suparchs.m_arr {
			if curarch == notarch {
				finded = 1
				break
			}
		}

		if finded == 0 {
			filters = append(filters, curarch)
		}
	}

	pkg.SetFilterArch(filters)

	for _, curdir := range flag.Args() {
		err := pkg.SearchDir(curdir)
		if err != nil {
			Error("%s", err.Error())
			os.Exit(5)
		}
	}
	os.Exit(0)
}
