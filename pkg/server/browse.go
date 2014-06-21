package server

import (
	"html/template"
	"net/http"
	"os"
	"path"
	"strings"

	"github.com/facette/facette/pkg/catalog"
	"github.com/facette/facette/pkg/library"
	"github.com/facette/facette/pkg/logger"
)

func (server *Server) serveBrowse(writer http.ResponseWriter, request *http.Request) {
	var err error

	if request.Method != "GET" && request.Method != "HEAD" {
		server.serveResponse(writer, nil, http.StatusMethodNotAllowed)
		return
	}

	// Redirect to default location
	if request.URL.Path == "/" {
		http.Redirect(writer, request, server.Config.URLPrefix+urlBrowsePath, 301)
		return
	}

	tmpl := template.New("layout.html").Funcs(template.FuncMap{
		"asset": server.templateAsset,
		"eq":    templateEqual,
		"ne":    templateNotEqual,
		"dump":  templateDumpMap,
		"hl":    templateHighlight,
	})

	setHTTPCacheHeaders(writer)

	if strings.HasPrefix(request.URL.Path, urlBrowsePath+"collections/") {
		err = server.serveBrowseCollection(writer, request, tmpl)
	} else if request.URL.Path == urlBrowsePath+"search" {
		err = server.serveBrowseSearch(writer, request, tmpl)
	} else if request.URL.Path == urlBrowsePath {
		err = server.serveBrowseIndex(writer, request, tmpl)
	} else {
		err = os.ErrNotExist
	}

	if os.IsNotExist(err) {
		server.serveError(writer, http.StatusNotFound)
	} else if err != nil {
		logger.Log(logger.LevelError, "server", "%s", err)
		server.serveError(writer, http.StatusInternalServerError)
	}
}

func (server *Server) serveBrowseIndex(writer http.ResponseWriter, request *http.Request,
	tmpl *template.Template) error {

	var data struct {
		URLPrefix string
	}

	// Set template data
	data.URLPrefix = server.Config.URLPrefix

	// Execute template
	tmpl, err := tmpl.ParseFiles(
		path.Join(server.Config.BaseDir, "html", "layout.html"),
		path.Join(server.Config.BaseDir, "html", "common", "element.html"),
		path.Join(server.Config.BaseDir, "html", "browse", "layout.html"),
		path.Join(server.Config.BaseDir, "html", "browse", "index.html"),
	)
	if err != nil {
		return err
	}

	return tmpl.Execute(writer, data)
}

func (server *Server) serveBrowseCollection(writer http.ResponseWriter, request *http.Request,
	tmpl *template.Template) error {

	var (
		err error
	)

	type collectionData struct {
		*library.Collection
		Parent string
	}

	var data struct {
		URLPrefix  string
		Collection *collectionData
		Request    *http.Request
	}

	data.URLPrefix = server.Config.URLPrefix

	data.Collection = &collectionData{Collection: &library.Collection{}}
	data.Collection.ID = strings.TrimPrefix(request.URL.Path, urlBrowsePath+"collections/")

	item, err := server.Library.GetItem(data.Collection.ID, library.LibraryItemCollection)
	if err != nil {
		return err
	}

	data.Collection.Collection = item.(*library.Collection)

	if request.FormValue("q") != "" {
		data.Collection.Collection = server.Library.FilterCollection(data.Collection.Collection, request.FormValue("q"))
	}

	if data.Collection.Collection.Parent != nil {
		data.Collection.Parent = data.Collection.Collection.Parent.ID
	} else {
		data.Collection.Parent = "null"
	}

	// Execute template
	tmpl, err = tmpl.ParseFiles(
		path.Join(server.Config.BaseDir, "html", "layout.html"),
		path.Join(server.Config.BaseDir, "html", "common", "element.html"),
		path.Join(server.Config.BaseDir, "html", "common", "graph.html"),
		path.Join(server.Config.BaseDir, "html", "browse", "layout.html"),
		path.Join(server.Config.BaseDir, "html", "browse", "collection.html"),
	)
	if err != nil {
		return err
	}

	data.Request = request

	return tmpl.Execute(writer, data)
}

func (server *Server) serveBrowseSearch(writer http.ResponseWriter, request *http.Request,
	tmpl *template.Template) error {

	var data struct {
		URLPrefix   string
		Count       int
		Request     *http.Request
		Sources     []*catalog.Source
		Collections []*library.Collection
	}

	// Set template data
	data.URLPrefix = server.Config.URLPrefix
	data.Request = request

	// Perform search filtering
	if request.FormValue("q") != "" {
		chunks := make([]string, 0)

		for _, chunk := range strings.Split(strings.ToLower(request.FormValue("q")), " ") {
			chunks = append(chunks, strings.Trim(chunk, " \t"))
		}

		for _, origin := range server.Catalog.Origins {
			for _, source := range origin.Sources {
				for _, chunk := range chunks {
					if strings.Index(strings.ToLower(source.Name), chunk) == -1 {
						goto nextOrigin
					}
				}

				data.Sources = append(data.Sources, source)
			nextOrigin:
			}
		}

		for _, collection := range server.Library.Collections {
			for _, chunk := range chunks {
				if strings.Index(strings.ToLower(collection.Name), chunk) == -1 {
					goto nextCollection
				}
			}

			data.Collections = append(data.Collections, collection)
		nextCollection:
		}
	}

	data.Count = len(data.Sources) + len(data.Collections)

	// Execute template
	tmpl, err := tmpl.ParseFiles(
		path.Join(server.Config.BaseDir, "html", "layout.html"),
		path.Join(server.Config.BaseDir, "html", "common", "element.html"),
		path.Join(server.Config.BaseDir, "html", "browse", "layout.html"),
		path.Join(server.Config.BaseDir, "html", "browse", "search.html"),
	)
	if err != nil {
		return err
	}

	return tmpl.Execute(writer, data)
}
