package main

import (
	"bytes"
	"flag"
	"fmt"
	"html/template"
	"io"
	"os"
	"os/exec"
	"runtime"
	"time"

	"github.com/microcosm-cc/bluemonday"
	"github.com/russross/blackfriday/v2"
)

const defaultTemplate = `<!DOCTYPE html>
<html>
<head>
	<meta http-equiv="content-type" content="text/html; charset=utf-8">
	<title>{{.Title}}</title>
</head>
<body>
{{.Body}}
</body>
</html>
`

type content struct {
	Title string
	Body  template.HTML
}

func main() {
	filename := flag.String("file", "", "File containing the Markdown")
	skipPreview := flag.Bool("s", false, "Skip auto-preview")
	templateFname := flag.String("t", "", "HTML template filename")

	flag.Parse() //don't forget to parse

	if *filename == "" {
		flag.Usage()
		os.Exit(1)
	}

	if err := run(*filename, *templateFname, os.Stdout, *skipPreview); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run(filename string, tFname string, output io.Writer, skipPreview bool) error {
	markdown, err := os.ReadFile(filename)
	if err != nil {
		return err
	}

	htmlContent, err := convertToHtml(markdown, tFname)
	if err != nil {
		return err
	}

	htmlFilename, err := getTempFile()
	if err != nil {
		return err
	}

	fmt.Fprintln(output, htmlFilename)

	if err := saveFile(htmlContent, htmlFilename); err != nil {
		return err
	}

	if skipPreview {
		return nil
	}

	defer os.Remove(htmlFilename)

	return preview(htmlFilename)
}

func convertToHtml(markdown []byte, tFname string) ([]byte, error) {
	output := blackfriday.Run(markdown)
	sanitizedHtml := bluemonday.UGCPolicy().SanitizeBytes(output)

	t := template.New("mdp")
	var err error

	if tFname != "" {
		t, err = template.ParseFiles(tFname)
		if err != nil {
			return nil, err
		}
	} else {
		t, err = t.Parse(defaultTemplate)
		if err != nil {
			return nil, err
		}
	}

	var buffer bytes.Buffer

	data := content{Title: "Markdown Preview Tool", Body: template.HTML(sanitizedHtml)}
	if err := t.Execute(&buffer, data); err != nil {
		return nil, err
	}

	return buffer.Bytes(), nil
}

func getTempFile() (string, error) {
	tempfile, err := os.CreateTemp("", "mdp*.html")
	if err != nil {
		return "", err
	}

	if err := tempfile.Close(); err != nil {
		return "", err
	}

	return tempfile.Name(), nil
}

func saveFile(content []byte, filename string) error {
	return os.WriteFile(filename, content, 0644)
}

func preview(htmlFile string) error {
	cName, cParams, err := getPreviewCommand()
	if err != nil {
		return err
	}

	cParams = append(cParams, htmlFile)
	cPath, err := exec.LookPath(cName)
	if err != nil {
		return err
	}

	err = exec.Command(cPath, cParams...).Run()
	time.Sleep(2 * time.Second)
	return err
}

func getPreviewCommand() (string, []string, error) {

	switch runtime.GOOS {
	case "linux":
		return "xdg-open", []string{}, nil
	case "windows":
		params := []string{"/C", "start"}
		return "cmd.exe", params, nil
	case "darwin":
		return "open", []string{}, nil
	default:
		return "", nil, fmt.Errorf("Not supported runtime %s", runtime.GOOS)
	}
}
