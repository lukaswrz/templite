// Package templite provides a simple caching mechanism for rendering text templates.
package templite

import (
	"errors"
	"io"
	"io/fs"
	"text/template"

	"github.com/alphadose/haxmap"
)

// ErrUncachedTemplate is returned when a requested template isn't found in the cache.
var ErrUncachedTemplate = errors.New("tr: template not found in cache")

// Renderer manages cached templates for rendering.
// FS is the file system used to load templates, and c is the template cache.
type Renderer struct {
	FS fs.FS
	c  *haxmap.Map[string, *template.Template]
}

// Cache parses templates from patterns and stores them under the given name.
// Returns an error if parsing fails.
func (r Renderer) Cache(name string, patterns ...string) error {
	t, err := template.New(name).ParseFS(r.FS, patterns...)
	if err != nil {
		return err
	}
	r.c.Set(name, t)
	return nil
}

// Render writes the output of a cached template to w using the provided data.
// Returns ErrUncachedTemplate if the template is not found, or any execution error.
func (r Renderer) Render(w io.Writer, data any, name string) error {
	t, cached := r.c.Get(name)
	if !cached {
		return ErrUncachedTemplate
	}
	return t.Execute(w, data)
}
