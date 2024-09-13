package server

import (
	"context"
	"embed"
	"fmt"
	"io"
	"io/fs"
	"log/slog"
	"mime/multipart"
	"net/http"
	"os"
	"text/template"

	"github.com/google/uuid"
	"github.com/oc-docker/libra/internal/store"
	"github.com/pkg/errors"
	sloghttp "github.com/samber/slog-http"
)

//go:embed assets
var assetsFS embed.FS

//go:embed templates
var templatesFS embed.FS

type ServerOptions struct {
	MaxUploadSize int64
	Store         store.Store
	BaseURL       string
}

type Server struct {
	handler http.Handler
	opts    *ServerOptions
}

// ServeHTTP implements http.Handler.
func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.handler.ServeHTTP(w, r)
}

func (s *Server) renderHome(w http.ResponseWriter, r *http.Request) {
	s.renderPage(w, r, "base", "home", struct {
		MaxUploadSize int64
	}{
		MaxUploadSize: s.opts.MaxUploadSize / 1024 / 1024,
	})
}

func (s *Server) renderDownload(w http.ResponseWriter, r *http.Request) {
	identifier := r.PathValue("identifier")
	s.renderPage(w, r, "base", "download", struct {
		RemoteIdentifier string
		BaseURL          string
	}{
		RemoteIdentifier: identifier,
		BaseURL:          s.opts.BaseURL,
	})
}

func (s *Server) handleDownload(w http.ResponseWriter, r *http.Request) {
	part := r.PathValue("part")

	if part != "data" && part != "info" {
		slog.Warn("unexpected required part", slog.Any("part", part))
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	identifier := r.PathValue("identifier")
	if _, err := uuid.Parse(identifier); err != nil {
		slog.Warn("could not parse identifier", slog.Any("error", errors.WithStack(err)), slog.Any("identifier", identifier))
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	reader, err := s.opts.Store.Reader(r.Context(), identifier+"."+part)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			return
		}

		slog.Error("could not open reader", slog.Any("error", errors.WithStack(err)), slog.Any("identifier", identifier), slog.Any("part", part))
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	if _, err := io.Copy(w, reader); err != nil {
		slog.Error("could not copy data", slog.Any("error", errors.WithStack(err)), slog.Any("identifier", identifier), slog.Any("part", part))
		return
	}

	if part != "data" {
		return
	}

	if err := s.opts.Store.Remove(context.Background(), identifier+".data"); err != nil {
		slog.Error("could not remove blob", slog.Any("error", errors.WithStack(err)), slog.Any("identifier", identifier), slog.Any("part", "data"))
	}

	if err := s.opts.Store.Remove(context.Background(), identifier+".info"); err != nil {
		slog.Error("could not remove blob", slog.Any("error", errors.WithStack(err)), slog.Any("identifier", identifier), slog.Any("part", "info"))
	}
}

func (s *Server) handleUpload(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, s.opts.MaxUploadSize)

	if err := r.ParseMultipartForm(s.opts.MaxUploadSize); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	files := r.MultipartForm.File["file"]

	if len(files) != 2 {
		http.Error(w, fmt.Sprintf("Unexpected number of uploaded files '%d'", len(files)), http.StatusBadRequest)
		return
	}

	formLocalIdentifier := r.MultipartForm.Value["local-identifier"]

	if len(formLocalIdentifier) == 0 {
		http.Error(w, "Missing required informations", http.StatusBadRequest)
		return
	}

	localIdentifier := formLocalIdentifier[0]
	remoteIdentifier := uuid.New().String()

	writeFile := func(fileHeader *multipart.FileHeader, ext string) error {
		file, err := fileHeader.Open()
		if err != nil {
			return errors.WithStack(err)
		}

		defer func() {
			if err := file.Close(); err != nil {
				slog.Error("could not close file", slog.Any("error", errors.WithStack(err)))
			}
		}()

		filename := remoteIdentifier + ext

		slog.Debug("writing file", slog.Any("filename", filename))

		writer, err := s.opts.Store.Writer(r.Context(), filename)
		if err != nil {
			return errors.WithStack(err)
		}

		defer func() {
			if err := writer.Close(); err != nil {
				slog.Error("could not close writer", slog.Any("error", errors.WithStack(err)))
			}
		}()

		if _, err := io.Copy(writer, file); err != nil {
			return errors.WithStack(err)
		}

		slog.Debug("file written", slog.Any("filename", filename))

		return nil
	}

	fileHeader := files[0]
	if fileHeader.Size > s.opts.MaxUploadSize {
		http.Error(w, fmt.Sprintf("The uploaded file is too big: %s. Please use a file less than %dMB in size", fileHeader.Filename, s.opts.MaxUploadSize/1024/1024), http.StatusBadRequest)
		return
	}

	if err := writeFile(fileHeader, ".data"); err != nil {
		slog.Error("could not write file", slog.Any("error", errors.WithStack(err)))
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	fileHeader = files[1]
	if fileHeader.Size > s.opts.MaxUploadSize {
		http.Error(w, fmt.Sprintf("The uploaded file is too big: %s. Please use a file less than %dMB in size", fileHeader.Filename, s.opts.MaxUploadSize/1024/1024), http.StatusBadRequest)
		return
	}

	if err := writeFile(fileHeader, ".info"); err != nil {
		slog.Error("could not write file", slog.Any("error", errors.WithStack(err)))
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/upload?r=%s&l=%s", remoteIdentifier, localIdentifier), http.StatusSeeOther)
}

func (s *Server) renderUpload(w http.ResponseWriter, r *http.Request) {
	remoteIdentifier := r.URL.Query().Get("r")
	localIdentifier := r.URL.Query().Get("l")
	s.renderPage(w, r, "base", "upload", struct {
		RemoteIdentifier string
		LocalIdentifier  string
		BaseURL          string
	}{RemoteIdentifier: remoteIdentifier, LocalIdentifier: localIdentifier, BaseURL: s.opts.BaseURL})
}

func (s *Server) renderPage(w http.ResponseWriter, r *http.Request, layout string, page string, data any) {
	tmpl, err := template.ParseFS(templatesFS, "templates/layouts/*.gohtml", "templates/partials/*.gohtml", "templates/pages/"+page+".gohtml")
	if err != nil {
		slog.Error("could not parse templates", slog.Any("error", err))
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	if err := tmpl.ExecuteTemplate(w, layout, data); err != nil {
		slog.Error("could not render template", slog.Any("error", err), slog.Any("data", data), slog.Any("page", page))
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
}

func (s *Server) init() {
	mux := http.NewServeMux()

	assetsDir, err := fs.Sub(assetsFS, "assets")
	if err != nil {
		slog.Error("could not create assets filesystem", slog.Any("error", err))
		os.Exit(1)
	}

	assetsHandler := http.FileServer(http.FS(assetsDir))

	mux.HandleFunc("/{$}", s.renderHome)
	mux.HandleFunc("POST /upload", s.handleUpload)
	mux.HandleFunc("GET /upload", s.renderUpload)
	mux.HandleFunc("GET /dl/{identifier}", s.renderDownload)
	mux.HandleFunc("POST /dl/{identifier}/{part}", s.handleDownload)
	mux.Handle("/", assetsHandler)

	s.handler = sloghttp.New(slog.Default())(mux)
}

func New(opts *ServerOptions) *Server {
	server := &Server{
		opts: opts,
	}
	server.init()
	return server
}

var _ http.Handler = &Server{}
