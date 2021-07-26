package portgate

import (
	"errors"
	"html/template"
	"io"
	"io/fs"
	"path/filepath"
	"strings"
)

// Templates houses all the Portgate page templates and can render them with correct styling.
type Templates struct {
	templateMap map[string]*template.Template
}

// LoadAllTemplates loads all Portgate templates from ./assets/templates/ and bundles them
// into a single Templates instance.
func LoadAllTemplates() (Templates, error) {
	templateMap := make(map[string]*template.Template)

	// We walk through every file in the templates folder.
	err := filepath.Walk("./assets/templates", func(path string, info fs.FileInfo, err error) error {
		// We only care about files which are actual templates and are not of special use (like "_base")
		if !strings.HasPrefix(info.Name(), "_") && strings.HasSuffix(path, ".template.html") {
			// We bundle the templates together with the base template so that we can render them together later.
			t, err := template.ParseFiles("./assets/templates/_base.template.html", path)
			if err == nil {
				// We keep the parsed template in the templateMap by it's filename.
				templateMap[info.Name()] = t
			}
		}

		return err
	})

	if err != nil {
		return Templates{}, err
	}

	return Templates{
		templateMap: templateMap,
	}, nil
}

// ExecuteTemplate executes a single template with the given name and data and writes it to the given
// io.Writer, which is usually fasthttp.RequestCtx.
func (templates *Templates) ExecuteTemplate(w io.Writer, name string, data interface{}) error {
	t := templates.templateMap[name]
	if t == nil {
		return errors.New("Unknown template name: " + name)
	}

	return t.ExecuteTemplate(w, "_base.template.html", data)
}
