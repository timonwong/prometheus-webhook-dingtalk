// Code generated for package deftmpl by go-bindata DO NOT EDIT. (@generated)
// sources:
// template/default.tmpl
package deftmpl

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func bindataRead(data []byte, name string) ([]byte, error) {
	gz, err := gzip.NewReader(bytes.NewBuffer(data))
	if err != nil {
		return nil, fmt.Errorf("Read %q: %v", name, err)
	}

	var buf bytes.Buffer
	_, err = io.Copy(&buf, gz)
	clErr := gz.Close()

	if err != nil {
		return nil, fmt.Errorf("Read %q: %v", name, err)
	}
	if clErr != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

type asset struct {
	bytes []byte
	info  os.FileInfo
}

type bindataFileInfo struct {
	name    string
	size    int64
	mode    os.FileMode
	modTime time.Time
}

// Name return file name
func (fi bindataFileInfo) Name() string {
	return fi.name
}

// Size return file size
func (fi bindataFileInfo) Size() int64 {
	return fi.size
}

// Mode return file mode
func (fi bindataFileInfo) Mode() os.FileMode {
	return fi.mode
}

// Mode return file modify time
func (fi bindataFileInfo) ModTime() time.Time {
	return fi.modTime
}

// IsDir return file whether a directory
func (fi bindataFileInfo) IsDir() bool {
	return fi.mode&os.ModeDir != 0
}

// Sys return file is sys mode
func (fi bindataFileInfo) Sys() interface{} {
	return nil
}

var _templateDefaultTmpl = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\xb4\x53\xc1\x6e\xdb\x30\x0c\xbd\xe7\x2b\x08\xfb\x12\x0b\x98\x7a\x2f\xb0\x0d\xc5\xb0\xf5\x12\x0c\x43\x82\xec\xb2\x06\x81\x1a\x33\xae\x5a\x89\xca\x64\xba\x0d\xe0\xe8\xdf\x07\xd9\x86\x63\x2d\xd9\x6d\xbd\x49\x14\xf9\x1e\x1f\x1f\xd5\xb6\x50\xe2\x5e\x13\x42\xb6\xdd\xd6\xcd\xe3\x33\xee\x38\x83\x10\x7e\xb5\x2d\xc8\x15\x2b\x6e\x6a\x38\x01\xbb\xf5\xe1\x80\x1e\x42\x68\x5b\xd0\x7b\xc0\xdf\xe3\x63\xb6\xd7\x5e\x53\x15\x6b\x6e\x63\xcd\x9d\x41\xcf\xb5\xfc\xd6\x45\xe1\x04\x06\xa9\x2f\x43\x2a\x21\x84\x0d\xc4\xa4\x7b\xef\x9a\xc3\x42\x3d\xa2\xa9\xe5\xca\x79\xc6\xf2\x87\xd2\xbe\x96\x3f\x95\x69\x30\x12\x3e\x3b\x4d\x90\x41\x44\x85\x9e\xb2\x62\x98\x47\x2c\xf9\xc5\x59\xeb\xa8\x2f\x2e\x86\xd8\x04\xaf\x80\x10\xe6\x6d\x0b\x6f\x9a\x9f\xd2\x64\xb9\x44\xeb\x5e\x31\x65\xff\xae\x2c\xd6\x7d\x83\x57\xd9\xc7\xc6\x8b\xf1\x34\x1e\x66\xc9\xf0\x54\x14\x6e\x15\xa9\x0a\xfd\x7a\xb9\x18\x8a\xe5\xd7\x23\xa3\x27\x65\xd6\xcb\x05\x84\x70\x93\xdf\x74\x79\xf5\x67\x8f\x3b\xd4\xaf\xe8\x3f\xc6\xa4\xe5\x70\x49\xd0\x53\x78\xc6\x23\xf7\x1c\x5b\xa3\x6b\x1e\xe0\xbd\xa2\x0a\x41\xc6\x74\x21\x7a\x49\x42\xcc\xce\x0f\x97\x33\x86\x10\x3e\xc1\x87\xce\x85\xa8\x3d\xda\x06\xa3\x78\x38\x81\x55\xfe\xa5\x74\x6f\x04\x27\x78\x62\x6b\x06\x99\x43\x4b\x42\xdc\x11\x39\x56\xac\x1d\xa5\x44\x93\xf8\x7f\x64\x5b\xb9\xc6\xef\xf0\x56\x08\xe8\xf6\xf1\x1e\x09\xbd\x62\xe7\xfb\x61\x6e\xe6\x57\x82\xc5\x6c\x76\xc5\xa9\xe9\x2c\x4b\x4d\x95\x34\x9a\x5e\x24\x6b\x36\x38\x4c\x92\xd1\x1e\x8c\xe2\xf4\x1f\xc8\x7f\xd9\x7d\xc6\xd8\x39\x62\xa4\xce\x8f\x3c\xcf\x73\x78\x78\xaf\x9f\xf3\xb0\x01\x21\x22\xb8\xa6\x12\x8f\xc9\x16\x43\xd6\x2d\x06\x29\xdb\xa9\xe9\xe6\x32\xd5\x73\xb1\x9a\x51\x57\xd1\xfb\x37\xcd\xbb\xd8\xb1\xbf\x7a\x4a\xdc\xf9\x13\x00\x00\xff\xff\xaa\x4e\x32\x4b\x3a\x04\x00\x00")

func templateDefaultTmplBytes() ([]byte, error) {
	return bindataRead(
		_templateDefaultTmpl,
		"template/default.tmpl",
	)
}

func templateDefaultTmpl() (*asset, error) {
	bytes, err := templateDefaultTmplBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "template/default.tmpl", size: 1082, mode: os.FileMode(420), modTime: time.Unix(1, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

// Asset loads and returns the asset for the given name.
// It returns an error if the asset could not be found or
// could not be loaded.
func Asset(name string) ([]byte, error) {
	cannonicalName := strings.Replace(name, "\\", "/", -1)
	if f, ok := _bindata[cannonicalName]; ok {
		a, err := f()
		if err != nil {
			return nil, fmt.Errorf("Asset %s can't read by error: %v", name, err)
		}
		return a.bytes, nil
	}
	return nil, fmt.Errorf("Asset %s not found", name)
}

// MustAsset is like Asset but panics when Asset would return an error.
// It simplifies safe initialization of global variables.
func MustAsset(name string) []byte {
	a, err := Asset(name)
	if err != nil {
		panic("asset: Asset(" + name + "): " + err.Error())
	}

	return a
}

// AssetInfo loads and returns the asset info for the given name.
// It returns an error if the asset could not be found or
// could not be loaded.
func AssetInfo(name string) (os.FileInfo, error) {
	cannonicalName := strings.Replace(name, "\\", "/", -1)
	if f, ok := _bindata[cannonicalName]; ok {
		a, err := f()
		if err != nil {
			return nil, fmt.Errorf("AssetInfo %s can't read by error: %v", name, err)
		}
		return a.info, nil
	}
	return nil, fmt.Errorf("AssetInfo %s not found", name)
}

// AssetNames returns the names of the assets.
func AssetNames() []string {
	names := make([]string, 0, len(_bindata))
	for name := range _bindata {
		names = append(names, name)
	}
	return names
}

// _bindata is a table, holding each asset generator, mapped to its name.
var _bindata = map[string]func() (*asset, error){
	"template/default.tmpl": templateDefaultTmpl,
}

// AssetDir returns the file names below a certain
// directory embedded in the file by go-bindata.
// For example if you run go-bindata on data/... and data contains the
// following hierarchy:
//     data/
//       foo.txt
//       img/
//         a.png
//         b.png
// then AssetDir("data") would return []string{"foo.txt", "img"}
// AssetDir("data/img") would return []string{"a.png", "b.png"}
// AssetDir("foo.txt") and AssetDir("notexist") would return an error
// AssetDir("") will return []string{"data"}.
func AssetDir(name string) ([]string, error) {
	node := _bintree
	if len(name) != 0 {
		cannonicalName := strings.Replace(name, "\\", "/", -1)
		pathList := strings.Split(cannonicalName, "/")
		for _, p := range pathList {
			node = node.Children[p]
			if node == nil {
				return nil, fmt.Errorf("Asset %s not found", name)
			}
		}
	}
	if node.Func != nil {
		return nil, fmt.Errorf("Asset %s not found", name)
	}
	rv := make([]string, 0, len(node.Children))
	for childName := range node.Children {
		rv = append(rv, childName)
	}
	return rv, nil
}

type bintree struct {
	Func     func() (*asset, error)
	Children map[string]*bintree
}

var _bintree = &bintree{nil, map[string]*bintree{
	"template": &bintree{nil, map[string]*bintree{
		"default.tmpl": &bintree{templateDefaultTmpl, map[string]*bintree{}},
	}},
}}

// RestoreAsset restores an asset under the given directory
func RestoreAsset(dir, name string) error {
	data, err := Asset(name)
	if err != nil {
		return err
	}
	info, err := AssetInfo(name)
	if err != nil {
		return err
	}
	err = os.MkdirAll(_filePath(dir, filepath.Dir(name)), os.FileMode(0755))
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(_filePath(dir, name), data, info.Mode())
	if err != nil {
		return err
	}
	err = os.Chtimes(_filePath(dir, name), info.ModTime(), info.ModTime())
	if err != nil {
		return err
	}
	return nil
}

// RestoreAssets restores an asset under the given directory recursively
func RestoreAssets(dir, name string) error {
	children, err := AssetDir(name)
	// File
	if err != nil {
		return RestoreAsset(dir, name)
	}
	// Dir
	for _, child := range children {
		err = RestoreAssets(dir, filepath.Join(name, child))
		if err != nil {
			return err
		}
	}
	return nil
}

func _filePath(dir, name string) string {
	cannonicalName := strings.Replace(name, "\\", "/", -1)
	return filepath.Join(append([]string{dir}, strings.Split(cannonicalName, "/")...)...)
}
