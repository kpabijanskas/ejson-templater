package main

import (
	"encoding/json"
	"flag"
	"log"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/Masterminds/sprig/v3"
	"github.com/Shopify/ejson"
	homedir "github.com/mitchellh/go-homedir"
)

func main() {
	os.Exit(run())
}

func run() int {
	keyDir := flag.String("keydir", "~/.config/ejson/", "/path/to/ejson/keydir")
	secretsFile := flag.String("secretsFile", "secrets.ejson", "/path/to/ejson/secrets")
	srcDir := flag.String("srcDir", ".", "/path/to/source/template/directory")
	dstDir := flag.String("dstDir", ".", "/path/to/destination/directory")
	pattern := flag.String("pattern", "*.gotmpl", "files to match")
	suffixTrim := flag.String("suffixTrim", ".gotmpl", "suffix to trim from filenames")
	flag.Parse()

	var err error
	*secretsFile, err = homedir.Expand(*secretsFile)
	if err != nil {
		log.Printf("unable to expand secretsFile path: %v", err)
		return 1
	}
	*keyDir, err = homedir.Expand(*keyDir)
	if err != nil {
		log.Printf("unable to expand keyDir path: %v", err)
		return 1
	}
	data, err := ejson.DecryptFile(*secretsFile, *keyDir, "")
	if err != nil {
		log.Printf("unable to decrypt secrets file: %v", err)
		return 1
	}

	vars := make(map[string]any)
	if err = json.Unmarshal(data, &vars); err != nil {
		log.Printf("unable to parse decrypted secrets file as json: %v", err)
		return 1
	}

	*srcDir, err = homedir.Expand(*srcDir)
	if err != nil {
		log.Printf("unable to expand srcDir path: %v", err)
		return 1
	}

	templates, err := template.ParseGlob(filepath.Join(*srcDir, *pattern))
	if err != nil {
		log.Printf("unable to parse templates: %v", err)
		return 1
	}

	templates.Funcs(sprig.FuncMap())
	templates.Option("missingkey=error")

	for _, t := range templates.Templates() {
		filename := strings.TrimSuffix(t.Name(), *suffixTrim)
		var f *os.File
		f, err = os.Create(filepath.Join(*dstDir, filename))
		if err != nil {
			log.Printf("error creating file %s: %v", f.Name(), err)
			return 1
		}
		defer f.Close()

		if err = t.Execute(f, vars); err != nil {
			log.Printf("error executing template %s: %v", t.Name(), err)
			return 1
		}
	}
	return 0
}
