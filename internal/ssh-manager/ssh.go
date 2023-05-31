package ssh_manager

import (
	"os"
	"path/filepath"
	"strings"
)

type Manager struct {
	path       string
	files      []string
	privateKey string
	publicKey  string
}

func scanDir(path, ext string) []string {
	var files []string
	dir, err := os.ReadDir(path)

	if err != nil {
		panic(err)
	}

	for _, file := range dir {
		if file.IsDir() {
			files = append(files, scanDir(path+"/"+file.Name(), ext)...)
			continue
		}
		f := strings.Split(file.Name(), ".")
		if f[len(f)-1] == ext {
			files = append(files, path+"/"+file.Name())
		}
	}

	return files
}

func NewManager(filePath string) *Manager {
	absFile, err := filepath.Abs(filePath)
	if err != nil {
		panic(err)
	}
	manager := &Manager{
		path:  absFile,
		files: []string{},
	}
	files := scanDir(filePath, "ssh")

	for _, file := range files {
		abs, err := filepath.Abs(file)
		if err != nil {
			panic(err)
		}
		manager.files = append(manager.files, abs)
	}

	return manager
}

func (m *Manager) GetFiles() []string {
	return m.files
}
