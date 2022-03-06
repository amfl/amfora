package display

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/makeworld-the-better-one/amfora/config"
)

// This file contains functions to construct a local mirror of everything that
// amfora loads. Since you're loading a file, you might as well save it, too!

// mirror will save a page in a subdirectory of the download directory
// determined by the url of the page. Existing files will be clobbered.
func mirror(u string, content string) (string, error) {
	// fmt.Fprintf(os.Stderr, "Saving URL: %s\n", u)
	// fmt.Fprintf(os.Stderr, "Page contents: %s\n", content)

	parsed, _ := url.Parse(u)

	// Convert URL to filesystem-safe path
	savePath := filepath.Join(config.DownloadsDir, parsed.Hostname(), parsed.Path)
	if strings.HasSuffix(u, "/") {
		// TODO don't assume extension?
		savePath = filepath.Join(savePath, "index.gmi")
	}

	// fmt.Fprintf(os.Stderr, "Writing to: %s\n", savePath)

	// Assert that path is a subdir of config.DownloadsDir to avoid security
	// issues
	rel, err := (filepath.Rel(config.DownloadsDir, savePath))
	if err != nil || strings.HasPrefix(rel, "..") {
		errText := fmt.Sprintf("Invalid rel path: %s\n", rel)
		// fmt.Fprintf(os.Stderr, errText)
		return "", errors.New(errText)
	}

	// Create any subdirectories as needed
	os.MkdirAll(filepath.Dir(savePath), os.ModePerm)

	// Write file, clobbering anything that is there already
	err = ioutil.WriteFile(savePath, []byte(content), 0644)
	if err != nil {
		// Just in case
		os.Remove(savePath)
		return "", err
	}

	return savePath, err
}
