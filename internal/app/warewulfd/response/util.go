package response

import (
	"github.com/hpcng/warewulf/internal/pkg/assets"
	"github.com/hpcng/warewulf/internal/pkg/errors"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
)

func getSanity(req *http.Request) (assets.NodeInfo, error) {
	url := strings.Split(req.URL.Path, "/")

	hwaddr := strings.ReplaceAll(url[2], "-", ":")
	node, err := assets.FindByHwaddr(hwaddr)
	if err != nil {
		return node, errors.New("Could not find node by HW address")
	}

	if node.Fqdn == "" {
		log.Printf("UNKNOWN: %15s: %s\n", hwaddr, req.URL.Path)
		return node, errors.New("Unknown node HW address: " + hwaddr)
	} else {
		log.Printf("REQ:   %15s: %s\n", node.Fqdn, req.URL.Path)
	}

	return node, nil
}

func sendFile(w http.ResponseWriter, filename string, sendto string) error {

	fd, err := os.Open(filename)
	if err != nil {
		return err
	}

	FileHeader := make([]byte, 512)
	fd.Read(FileHeader)
	FileContentType := http.DetectContentType(FileHeader)
	FileStat, _ := fd.Stat()
	FileSize := strconv.FormatInt(FileStat.Size(), 10)

	w.Header().Set("Content-Disposition", "attachment; filename=kernel")
	w.Header().Set("Content-Type", FileContentType)
	w.Header().Set("Content-Length", FileSize)

	fd.Seek(0, 0)
	io.Copy(w, fd)

	fd.Close()
	return nil
}