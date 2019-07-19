package files

import (
	"fmt"
	"os"
	"strings"
	"sync"
)

var once sync.Once

// FileData : 文件数据对象
type FileData struct {
	data []byte
	sz   int64
}

// GetBytes : 获取数据
func (fd *FileData) GetBytes() []byte {
	return fd.data
}

// GetSize : 获取当前数据的字节数
func (fd *FileData) GetSize() int64 {
	return fd.sz
}

// IsNil : 是否是空数据
func (fd *FileData) IsNil() bool {
	if fd.sz == 0 || fd.data == nil {
		return true
	}
	return false
}

func (fd *FileData) close() {
	fd.data = nil
	fd.sz = 0
}

type virtaulFile struct {
	v          map[string]FileData
	cachedFile int
	cachedMem  int64
	vpath      string
	lock       sync.RWMutex
}

var (
	defaultV = virtaulFile{v: make(map[string]FileData), cachedFile: 0, cachedMem: 0}
)

// WithRootPath : 设置虚拟文件系统根目录
func WithRootPath(rootPath string) {
	defaultV.vpath = rootPath
}

// IsFileExist : 判断在根目录下此文件是否存在
func IsFileExist(filename string) bool {
	fullPath := GetFullPathForFilename(filename)
	defaultV.lock.RLock()
	_, ok := defaultV.v[fullPath]
	if ok {
		defaultV.lock.RUnlock()
		return true
	}
	defaultV.lock.RUnlock()

	_, err := os.Stat(fullPath)
	if err != nil {
		return false
	}
	return true
}

// GetFullPathForFilename : 获取文件的全路径
func GetFullPathForFilename(filename string) string {
	if strings.HasPrefix(filename, "./") {
		filename = strings.Replace(filename, "./", "", 1)
	} else if strings.HasPrefix(filename, "/") {
		filename = strings.Replace(filename, "/", "", 1)
	}

	return defaultV.vpath + filename
}

// GetCachedInfo : 获取虚拟文件系统缓存文件信息
func GetCachedInfo() (int, int64) {
	return defaultV.cachedFile, defaultV.cachedMem
}

func GetVirtualPath() string {
	return defaultV.vpath
}

//GetDataFromFile : 获得指定文件的文件数据快
func GetDataFromFile(fullPath string) FileData {
	defaultV.lock.RLock()
	v, ok := defaultV.v[fullPath]
	if ok {
		defaultV.lock.RUnlock()
		return v
	}
	defaultV.lock.RUnlock()

	defaultV.lock.Lock()
	defer defaultV.lock.Unlock()

	f, err := os.Open(fullPath)
	if err != nil {
		return FileData{}
	}
	defer f.Close()

	flen, _ := f.Seek(0, os.SEEK_END)
	f.Seek(0, os.SEEK_SET)

	fm := FileData{data: make([]byte, flen), sz: flen}

	n, err := f.Read(fm.data)
	if err != nil {
		panic(fmt.Sprintf("Read File Fail:%s", err.Error()))
	} else if int64(n) < fm.sz {
		panic(fmt.Sprintf("Read File Fail[%d:%d]", n, fm.sz))
	}

	defaultV.cachedFile++
	defaultV.cachedMem += flen

	defaultV.v[fullPath] = fm
	return defaultV.v[fullPath]
}

// Close : 关闭虚拟文件系统
func Close() {
	defaultV.lock.Lock()
	defer defaultV.lock.Unlock()

	for k, v := range defaultV.v {
		v.close()
		delete(defaultV.v, k)
	}
}
