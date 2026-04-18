// Package archive provides functionality to export and import profile archives
// as compressed tar bundles for backup and portability.
package archive

import (
	"archive/tar"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// DefaultOptions returns sensible defaults.
func DefaultOptions() Options {
	return Options{Compress: true}
}

// Options controls archive behaviour.
type Options struct {
	Compress bool
}

// Entry represents a single profile inside an archive.
type Entry struct {
	Name string            `json:"name"`
	Vars map[string]string `json:"vars"`
}

// Pack writes profiles from store dir into w as a gzipped tar archive.
func Pack(storeDir string, profiles []string, w io.Writer) error {
	gw := gzip.NewWriter(w)
	defer gw.Close()
	tw := tar.NewWriter(gw)
	defer tw.Close()

	for _, name := range profiles {
		path := filepath.Join(storeDir, name+".json")
		data, err := os.ReadFile(path)
		if err != nil {
			return fmt.Errorf("archive: read profile %q: %w", name, err)
		}
		hdr := &tar.Header{
			Name: name + ".json",
			Mode: 0600,
			Size: int64(len(data)),
		}
		if err := tw.WriteHeader(hdr); err != nil {
			return err
		}
		if _, err := tw.Write(data); err != nil {
			return err
		}
	}
	return nil
}

// Unpack reads a gzipped tar archive from r and writes profile JSON files
// into storeDir. Returns the list of profile names extracted.
func Unpack(r io.Reader, storeDir string, overwrite bool) ([]string, error) {
	if err := os.MkdirAll(storeDir, 0755); err != nil {
		return nil, err
	}
	gr, err := gzip.NewReader(r)
	if err != nil {
		return nil, fmt.Errorf("archive: gzip open: %w", err)
	}
	defer gr.Close()
	tr := tar.NewReader(gr)

	var names []string
	for {
		hdr, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("archive: tar read: %w", err)
		}
		if !strings.HasSuffix(hdr.Name, ".json") {
			continue
		}
		dest := filepath.Join(storeDir, filepath.Base(hdr.Name))
		if !overwrite {
			if _, err := os.Stat(dest); err == nil {
				return nil, fmt.Errorf("archive: %q already exists", hdr.Name)
			}
		}
		data, err := io.ReadAll(tr)
		if err != nil {
			return nil, err
		}
		// validate JSON
		var v map[string]string
		if err := json.Unmarshal(data, &v); err != nil {
			return nil, fmt.Errorf("archive: invalid JSON in %q: %w", hdr.Name, err)
		}
		if err := os.WriteFile(dest, data, 0600); err != nil {
			return nil, err
		}
		names = append(names, strings.TrimSuffix(hdr.Name, ".json"))
	}
	return names, nil
}
