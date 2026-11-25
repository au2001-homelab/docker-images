package main

import (
	"flag"
	"io/fs"
	"log"
	"net/http"
	"time"
)

// parentLinkFS adds a ".." entry to every directory except the root, so that
// http.FileServer's directory listings include a link to the parent directory.
type parentLinkFS struct {
	http.FileSystem
}

func (fsys parentLinkFS) Open(name string) (http.File, error) {
	f, err := fsys.FileSystem.Open(name)
	if err != nil || name == "/" {
		return f, err
	}
	return parentLinkFile{f}, nil
}

// parentLinkFile embeds the http.File interface (not *os.File), so it doesn't
// implement fs.ReadDirFile and http.FileServer lists it through Readdir.
type parentLinkFile struct {
	http.File
}

func (f parentLinkFile) Readdir(count int) ([]fs.FileInfo, error) {
	entries, err := f.File.Readdir(count)
	if err != nil || count > 0 {
		return entries, err
	}
	return append([]fs.FileInfo{parentDirInfo{}}, entries...), nil
}

type parentDirInfo struct{}

func (parentDirInfo) Name() string       { return ".." }
func (parentDirInfo) Size() int64        { return 0 }
func (parentDirInfo) Mode() fs.FileMode  { return fs.ModeDir | 0o555 }
func (parentDirInfo) ModTime() time.Time { return time.Time{} }
func (parentDirInfo) IsDir() bool        { return true }
func (parentDirInfo) Sys() any           { return nil }

func main() {
	addr := flag.String("addr", ":8000", "address to listen on")
	dir := flag.String("dir", "/shares", "directory to serve")
	flag.Parse()

	srv := &http.Server{
		Addr:              *addr,
		Handler:           http.FileServer(parentLinkFS{http.Dir(*dir)}),
		ReadHeaderTimeout: 10 * time.Second,
	}

	log.Printf("serving %s on %s", *dir, *addr)
	log.Fatal(srv.ListenAndServe())
}
