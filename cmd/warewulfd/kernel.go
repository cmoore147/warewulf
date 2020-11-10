package main

import (
	"fmt"
	"github.com/hpcng/warewulf/internal/pkg/config"
	"log"
	"net/http"
)

func kernelSend(w http.ResponseWriter, req *http.Request) {

	node, err := getSanity(req)
	if err != nil {
		w.WriteHeader(404)
		log.Println(err)
		return
	}

	if node.KernelVersion != "" {
		fileName := fmt.Sprintf("%s/provision/kernel/vmlinuz-%s", config.LocalStateDir, node.KernelVersion)

		err := sendFile(w, fileName, node.Fqdn)
		if err != nil {
			log.Printf("ERROR: %s\n", err)
		} else {
			log.Printf("SEND:  %15s: %s\n", node.Fqdn, fileName)
		}

	} else {
		w.WriteHeader(503)
		log.Printf("ERROR: No 'kernel version' set for node %s\n", node.Fqdn)
	}

	return
}
