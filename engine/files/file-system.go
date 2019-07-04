package files

import (
  "sync"
  "fmt"
  "strings"
  "os"
)

var once sync.Once

type FileData struct {
  _data []byte
  _sz int64
}

type FileSystem struct {
  _vfile map[string]*FileData
  _rpath string
  _rwmtx sync.RWMutex
}

var instance *FileSystem

func GetInstance() *FileSystem {
  once.Do(func() {
    instance = new(FileSystem)
  })

  return instance
}

func (D *FileData)GetData()*byte {
  return &D._data[0]
}

func (D *FileData)GetBytes() int64 {
  return D._sz
}

func (D *FileData)IsNull() bool {
  if D._data == nil {
    return true
  }
  return false
}


func (D *FileData)SetData(f *os.File) {
  n, err := f.Read(D._data)
  if err != nil {
    panic(fmt.Sprintf("Read File Fail:%s",  err.Error()))
  } else if int64(n) < D._sz {
      panic(fmt.Sprintf("Read File Fail[%d:%d]", n, D._sz))
  }
}

func (D *FileData) Clear() {
  D._data = nil
  D._sz = 0
}

func (F *FileSystem)SetRootPath(rootPath string) {
  if (!strings.HasSuffix(rootPath, "/")) {
   rootPath = rootPath + "/"
  }

  F._rpath = rootPath
}

func (F *FileSystem)GetRootPath() string {
  return F._rpath
}

func (F *FileSystem)Close() {
  F._rwmtx.Lock()
  defer F._rwmtx.Unlock()
  for key, val := range F._vfile {
      val.Clear()
      F._vfile[key] = nil
  }
}

func (F *FileSystem)IsFileExist(filename string) bool {
  fullPath := F.GetFullPathForFilename(filename)
  F._rwmtx.RLock()
  _, ok := F._vfile[fullPath]
  if ok {
    F._rwmtx.RUnlock()
    return true
  }
  F._rwmtx.RUnlock()

  _, err := os.Stat(fullPath)
  if err != nil {
    return false
  }
  return true
}

func (F *FileSystem)GetFullPathForFilename(filename string) string {
  if (strings.HasPrefix(filename, "./")) {
    filename = strings.Replace(filename, "./", "", 1)
  } else if (strings.HasPrefix(filename, "/")) {
    filename = strings.Replace(filename, "/", "", 1)
  }

  return F._rpath + filename
}

func (F *FileSystem)GetDataFromFile(fullPath string) *FileData {
  F._rwmtx.RLock()
  v, ok := F._vfile[fullPath]
  if ok {
    F._rwmtx.RUnlock()
    return v
  }
  F._rwmtx.RUnlock()

  F._rwmtx.Lock()
  defer F._rwmtx.Unlock()

  f, err := os.Open(fullPath)
  if err != nil {
    return nil
  }
  defer f.Close()

  flen,_ := f.Seek(0, os.SEEK_END)
  f.Seek(0, os.SEEK_SET)

  fi := &FileData{make([]byte, flen), flen}

  fi.SetData(f)
  F._vfile[fullPath] = fi

  return fi
}
