package main

import (
	"bytes"
	"fmt"
	"go/build"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
)

const (
	goSourceCodeExt = ".go"

	goTemplateExt    = ".go.template"
	goTemplateExtLen = len(goTemplateExt)

	separatedGoTemplateExt    = ".sep.go.template"
	separatedGoTemplateExtLen = len(separatedGoTemplateExt)
)

var separatedSuffixes = []string{
	"state.go.template",
	"message.go.template",
}

var replacers = [][2]string{
	{
		"\"github.com/filecoin-project/lotus/chain/types\"",
		"types \"github.com/filecoin-project/venus/venus-shared/chain\"",
	},
	{
		"github.com/filecoin-project/lotus/chain/actors",
		"github.com/filecoin-project/venus/venus-shared/actors",
	},
	{
		"\"github.com/filecoin-project/lotus/node/modules/dtypes\"",
		"",
	},
	{
		"dtypes.NetworkName",
		"string",
	},
}

func findActorsPkgDir() (string, error) {
	pkg, err := build.Import("github.com/filecoin-project/lotus/chain/actors", ".", build.FindOnly)
	if err != nil {
		return "", fmt.Errorf("find local build path for louts: %w", err)
	}

	return pkg.Dir, nil
}

func fetch(src, dst string, paths []string) error {
	err := os.MkdirAll(dst, 0755)
	if err != nil {
		return fmt.Errorf("mkdir-all for %s: %w", dst, err)
	}

	for _, rel := range paths {
		if err := fetchOne(src, dst, rel); err != nil {
			return fmt.Errorf("fetch template for %s: %w", rel, err)
		}

		log.Printf("\t%s done", rel)
	}

	return nil
}

func fetchOne(srcDir, dstDir string, rel string) error {
	dstRel := rel
	for _, suffix := range separatedSuffixes {
		if strings.HasSuffix(rel, suffix) {
			dstRel = strings.ReplaceAll(rel, goTemplateExt, separatedGoTemplateExt)
			break
		}
	}

	fsrc, err := os.Open(filepath.Join(srcDir, rel))
	if err != nil {
		return fmt.Errorf("open source file: %w", err)
	}

	defer fsrc.Close() // nolint: errcheck

	dstPath := filepath.Join(dstDir, dstRel)
	err = os.MkdirAll(filepath.Dir(dstPath), 0755)
	if err != nil {
		return fmt.Errorf("mkdir for %s: %w", dstPath, err)
	}

	fdst, err := os.OpenFile(dstPath, os.O_CREATE|os.O_TRUNC|os.O_RDWR, 0644)
	if err != nil {
		return fmt.Errorf("open dst file: %w", err)
	}

	defer fdst.Close() // nolint: errcheck

	var buf bytes.Buffer

	if _, err := buf.WriteString(fmt.Sprintf("// FETCHED FROM LOTUS: %s\n\n", rel)); err != nil {
		return fmt.Errorf("write file header: %w", err)
	}

	_, err = io.Copy(&buf, fsrc)
	if err != nil {
		return fmt.Errorf("copy to buffer: %w", err)
	}

	data := buf.Bytes()
	for _, replacer := range replacers {
		data = bytes.ReplaceAll(data, []byte(replacer[0]), []byte(replacer[1]))
	}

	_, err = io.Copy(fdst, bytes.NewReader(data))
	if err != nil {
		return fmt.Errorf("copy to dst file: %w", err)
	}

	err = fdst.Sync()
	if err != nil {
		return fmt.Errorf("dst file sync: %w", err)
	}

	return nil
}