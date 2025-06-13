package apiserver

import (
	"backend/internal/app/model"
	"backend/internal/app/store"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

type APIServer struct {
	config *Config
	logger *logrus.Logger
	router *mux.Router
	store  *store.Store
}

func New(config *Config) *APIServer {
	return &APIServer{
		config: config,
		logger: logrus.New(),
		router: mux.NewRouter(),
	}
}

func (s *APIServer) Start() error {
	if err := s.configureLogger(); err != nil {
		return err
	}

	if err := s.configureStore(); err != nil {
		return err
	}

	s.configureRouter()

	s.logger.Info("START API-SERVER")
	return http.ListenAndServe(s.config.BinAddr, s.router)
}

func (s *APIServer) configureLogger() error {
	level, err := logrus.ParseLevel(s.config.LogLevel)
	if err != nil {
		return err
	}
	fmt.Println(level)
	s.logger.SetLevel(level)
	return nil
}

func (s *APIServer) configureStore() error {
	minioClient, err := store.NewMinioClient(
		s.config.MinIO.Endpoint,
		s.config.MinIO.AccessKey,
		s.config.MinIO.SecretKey,
		s.config.MinIO.Bucket,
		s.config.MinIO.UseSSL,
	)
	if err != nil {
		return fmt.Errorf("failed to init MinIO: %v", err)
	}

	st := store.New(s.config.Store, minioClient)
	if err := st.Open(); err != nil {
		return err
	}
	s.store = st

	return nil
}

func (s *APIServer) configureRouter() {
	s.router.Use(s.loggingMiddleware)

	s.router.HandleFunc("/hello", s.handleHello())
	s.router.HandleFunc("/tracks", s.handleGetAllTracks()).Methods("GET")
	s.router.HandleFunc("/tracks", s.handleUploadTrack()).Methods("POST")
	s.router.HandleFunc("/tracks/{id}/audio", s.handleTrackAudio).Methods("GET")
}

func (s *APIServer) loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		s.logger.Infof("%s %s %s", r.RemoteAddr, r.Method, r.RequestURI)
		next.ServeHTTP(w, r)
	})
}

func (s *APIServer) handleHello() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, "Hello")
	}
}

func (s *APIServer) handleGetAllTracks() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tracks, err := s.store.Track().GetAll()
		if err != nil {
			s.logger.Errorf("Failed to get tracks: %v", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(tracks)
	}
}

func (s *APIServer) handleUploadTrack() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseMultipartForm(10 << 20); err != nil {
			s.logger.Errorf("File too large: %v", err)
			http.Error(w, "File too large", http.StatusBadRequest)
			return
		}

		file, _, err := r.FormFile("file")
		if err != nil {
			s.logger.Errorf("Invalid file: %v", err)
			http.Error(w, "Invalid file", http.StatusBadRequest)
			return
		}
		defer file.Close()

		audioData, err := io.ReadAll(file)
		if err != nil {
			s.logger.Errorf("Failed to read file: %v", err)
			http.Error(w, "Failed to process file", http.StatusInternalServerError)
			return
		}

		track := &model.Track{
			Metadata: model.Metadata{
				Title:  r.FormValue("title"),
				Artist: r.FormValue("artist"),
			},
		}

		createdTrack, err := s.store.Track().Create(track, audioData)
		if err != nil {
			s.logger.Errorf("Failed to save track: %v", err)
			http.Error(w, "Failed to save track", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(createdTrack)
	}
}

// в handlers.go или apiserver/server.go

func (s *APIServer) handleTrackAudio(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]
	if idStr == "" {
		http.Error(w, "missing track ID", http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "invalid track ID", http.StatusBadRequest)
		return
	}

	track, err := s.store.Track().FindByID(id)
	if err != nil {
		http.Error(w, "track not found", http.StatusNotFound)
		return
	}

	reader, size, err := s.store.Track().GetAudioStream(track.ObjectKey)
	if err != nil {
		http.Error(w, "could not get audio stream", http.StatusInternalServerError)
		return
	}
	defer reader.Close()

	w.Header().Set("Content-Type", "audio/mpeg")
	w.Header().Set("Accept-Ranges", "bytes")
	w.Header().Set("Content-Disposition", fmt.Sprintf("inline; filename=\"%s - %s.mp3\"",
		track.Metadata.Artist, track.Metadata.Title))

	rangeHeader := r.Header.Get("Range")
	if rangeHeader == "" {
		// Без Range — отдаем весь файл
		w.Header().Set("Content-Length", fmt.Sprintf("%d", size))
		w.WriteHeader(http.StatusOK)
		io.Copy(w, reader)
		return
	}

	// Обработка Range: bytes=START-END
	var start, end int64
	_, err = fmt.Sscanf(rangeHeader, "bytes=%d-%d", &start, &end)
	if err != nil || end == 0 {
		end = size - 1
	}

	// Проверка границ
	if start > end || end >= size {
		w.Header().Set("Content-Range", fmt.Sprintf("bytes */%d", size))
		http.Error(w, "invalid range", http.StatusRequestedRangeNotSatisfiable)
		return
	}

	// Смещаемся
	_, err = reader.Seek(start, io.SeekStart)
	if err != nil {
		http.Error(w, "seek failed", http.StatusInternalServerError)
		return
	}

	contentLength := end - start + 1
	w.Header().Set("Content-Length", fmt.Sprintf("%d", contentLength))
	w.Header().Set("Content-Range", fmt.Sprintf("bytes %d-%d/%d", start, end, size))
	w.WriteHeader(http.StatusPartialContent)

	// Отдаем диапазон
	io.CopyN(w, reader, contentLength)
}
