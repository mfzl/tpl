package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

func env(key string) string {
	return os.Getenv(key)
}

type templatePair struct {
	source      string
	destination string
}

func (p *templatePair) FromString(value string) error {
	pair := strings.Split(value, ":")

	if len(pair) < 2 {
		return fmt.Errorf("invalid pair: '%s'", value)
	}

	srcInfo, err := os.Stat(pair[0])
	if err != nil {
		if errIsFileNotFound(err) {
			return fmt.Errorf("'%s' is not found", pair[0])
		}

		return fmt.Errorf("stat: %w", err)
	}

	if srcInfo.IsDir() {
		return fmt.Errorf("source template is a directory")
	}

	p.source = pair[0]
	p.destination = pair[1]

	return nil
}

func main() {

	pairs := []*templatePair{}

	for _, arg := range os.Args[1:] {
		tplPair := &templatePair{}

		err := tplPair.FromString(arg)
		exitOnError(err)

		pairs = append(pairs, tplPair)
	}

	if len(pairs) < 1 {
		log.Fatalf("minimum of 1 template and a destination is required")
	}

	for _, p := range pairs {
		err := compileTemplate(p)
		if errIsFileNotFound(err) {
			fmt.Printf("'%s' is not found\n", p.source)
			os.Exit(1)
		}
		exitOnError(err)
	}
}

func compileTemplate(p *templatePair) error {
	sourceAbs, err := filepath.Abs(p.source)
	if err != nil {
		return fmt.Errorf("get absolute path: %w", err)
	}

	tpl, err := template.New(filepath.Base(p.source)).Funcs(template.FuncMap{"env": env}).ParseFiles(sourceAbs)
	if err != nil {
		return fmt.Errorf("init template: %w", err)
	}

	resultFile, err := os.Create(p.destination)
	if err != nil {
		return fmt.Errorf("create destination file: %w", err)
	}
	defer resultFile.Close()

	err = tpl.Execute(resultFile, nil)
	if err != nil {
		return fmt.Errorf("compiling template: %w", err)
	}

	return nil
}

func errIsFileNotFound(err error) bool {
	return errors.Is(err, os.ErrNotExist)
}

func exitOnError(err error) {
	if err != nil {
		log.Printf("[ERRO] %+v\n", err)
		os.Exit(1)
	}
}
