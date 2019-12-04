// Code generated by go-bindata.
// sources:
// template/default.tmpl
// DO NOT EDIT!

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

func (fi bindataFileInfo) Name() string {
	return fi.name
}
func (fi bindataFileInfo) Size() int64 {
	return fi.size
}
func (fi bindataFileInfo) Mode() os.FileMode {
	return fi.mode
}
func (fi bindataFileInfo) ModTime() time.Time {
	return fi.modTime
}
func (fi bindataFileInfo) IsDir() bool {
	return false
}
func (fi bindataFileInfo) Sys() interface{} {
	return nil
}

var _templateDefaultTmpl = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\xb4\x53\xc1\x6e\xdb\x30\x0c\xbd\xfb\x2b\x08\xfb\x52\x0b\x98\x7a\x0f\xb0\x0d\xc5\xb0\xf5\x12\x0c\x43\xb2\xec\xb2\x06\x81\x1a\x33\xae\x5a\x89\xca\x64\xba\x0d\xe0\x68\xdf\x3e\x48\x36\xdc\xb8\xc9\x6e\xdb\x4d\xa2\xc8\xf7\xf8\xf8\xa8\xae\x83\x0a\x77\x9a\x10\xf2\xcd\xa6\x69\xef\x1f\x71\xcb\x39\x84\xf0\xb3\xeb\x40\x2e\x59\x71\xdb\xc0\x11\xd8\xad\xf6\x7b\xf4\x10\x42\xd7\x81\xde\x01\xfe\x1a\x1f\xf3\x9d\xf6\x9a\xea\x58\x33\x8b\x35\x37\x06\x3d\x37\xf2\x4b\x8a\xc2\x11\x0c\x52\x5f\x86\x54\x41\x08\x6b\x88\x49\xb7\xde\xb5\xfb\xb9\xba\x47\xd3\xc8\xa5\xf3\x8c\xd5\x37\xa5\x7d\x23\x7f\x28\xd3\x62\x24\x7c\x74\x9a\x20\x87\x88\x0a\x3d\x65\xcd\x70\x15\xb1\xe4\x27\x67\xad\xa3\xbe\xb8\x1c\x62\x27\x78\x25\x84\x70\xd5\x75\xf0\xa2\xf9\x61\x9a\x2c\x17\x68\xdd\x33\x4e\xd9\xbf\x2a\x8b\x4d\xdf\xe0\x45\xf6\xb1\xf1\x72\x3c\x8d\x87\x6c\x32\x3c\x15\x85\x5b\x45\xaa\x46\xbf\x5a\xcc\x87\x62\xf9\xf9\xc0\xe8\x49\x99\xd5\x62\x0e\x21\x5c\x17\xd7\x29\xaf\xf9\xe8\x71\x8b\xfa\x19\xfd\xfb\x98\xb4\x18\x2e\x13\xf4\x29\x3c\xe3\x81\x7b\x8e\x8d\xd1\x0d\x0f\xf0\x5e\x51\x8d\x20\x63\xba\x10\xbd\x24\x21\xb2\xd7\x87\xf3\x19\x43\x08\x1f\xe0\x5d\x72\x21\x6a\x8f\xb6\xc1\x28\x1e\x8e\x60\x95\x7f\xaa\xdc\x0b\xc1\x11\x1e\xd8\x9a\x41\xe6\xd0\x92\x10\x37\x44\x8e\x15\x6b\x47\x53\xa2\x93\xf8\x3f\x64\x5b\xba\xd6\x6f\x71\x26\x04\xa4\x7d\xbc\x45\x42\xaf\xd8\xf9\x7e\x98\xeb\xab\x0b\xc1\x32\xcb\x2e\x38\x75\x3a\xcb\x4a\x53\x2d\x8d\xa6\x27\xc9\x9a\x0d\x0e\x93\x64\xb4\x7b\xa3\x78\xfa\x0f\xe4\xdf\xec\x7e\xc5\xd8\x3a\x62\xa4\xe4\x47\x51\x14\x05\xdc\xfd\xaf\x9f\x73\xb7\x06\x21\x22\xb8\xa6\x0a\x0f\x93\x2d\x86\x3c\x2d\x06\x29\x9b\xd4\xa4\xb9\x9c\xea\x39\x5b\xcd\xa8\xab\x14\x22\x4b\x94\xac\x2d\xce\xe0\x37\x8c\x3d\x7c\xd7\xc9\xa9\xec\x0d\xca\xd9\x06\xbe\xe9\x78\xe2\xdd\x9f\x00\x00\x00\xff\xff\x44\xdc\xe3\xea\x58\x04\x00\x00")

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

	info := bindataFileInfo{name: "template/default.tmpl", size: 1112, mode: os.FileMode(420), modTime: time.Unix(1, 0)}
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
