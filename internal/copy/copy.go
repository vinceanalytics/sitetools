package copy

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/vinceanalytics/sitetools/data"
)

// Copies all assets to dst
func Copy(dst string) error {
	os.RemoveAll(dst)
	err := fs.WalkDir(data.Assets, ".", func(path string, d fs.DirEntry, err error) error {
		dstPath := filepath.Join(dst, path)
		if d.IsDir() {
			return os.MkdirAll(dstPath, 0755)
		}
		data, err := data.Assets.ReadFile(path)
		if err != nil {
			return err
		}
		fmt.Println("=>", dstPath)
		return os.WriteFile(dstPath, data, 0600)
	})
	if err != nil {
		return err
	}
	cname := filepath.Join(dst, "CNAME")
	fmt.Println("=>", cname)
	install := filepath.Join(dst, "install.sh")
	fmt.Println("=>", install)
	return errors.Join(
		os.WriteFile(cname, data.Cname, 0600),
		os.WriteFile(install, data.Install, 0600),
	)
}
