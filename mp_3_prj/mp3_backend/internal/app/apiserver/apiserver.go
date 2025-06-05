package apiserver

import (
	"backend/internal/app/model"
	"backend/internal/app/store"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

//APIServer

type APIServer struct {
	config *Config
	logger *logrus.Logger
	router *mux.Router
	store  *store.Store
	minio  *store.MinioClient
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

	s.configureRouter()

	if err := s.configureStore(); err != nil {
		return err
	}
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

func (s *APIServer) configureRouter() {
	s.router.HandleFunc("/hello", s.handleHello())
	s.router.HandleFunc("/tracks", s.handleGetAllTracks()).Methods("GET")
	s.router.HandleFunc("/tracks", s.handleUploadTrack()).Methods("POST")
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
	s.minio = minioClient

	st := store.New(s.config.Store)
	if err := st.Open(); err != nil {
		return err
	}

	s.store = st
	return nil
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
		if err := json.NewEncoder(w).Encode(tracks); err != nil {
			s.logger.Errorf("Failed to encode response: %v", err)
		}
	}
}

func (s *APIServer) handleUploadTrack() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Ограничиваем размер файла (например, 10 МБ)
		if err := r.ParseMultipartForm(10 << 20); err != nil {
			s.logger.Errorf("File too large: %v", err)
			http.Error(w, "File size exceeds 10MB", http.StatusBadRequest)
			return
		}

		// Получаем файл из запроса
		file, header, err := r.FormFile("file")
		if err != nil {
			s.logger.Errorf("Invalid file: %v", err)
			http.Error(w, "Invalid file", http.StatusBadRequest)
			return
		}
		defer file.Close()
		fmt.Println(header)

		// Читаем файл в []byte
		audioData, err := io.ReadAll(file)
		if err != nil {
			s.logger.Errorf("Failed to read file: %v", err)
			http.Error(w, "Failed to process file", http.StatusInternalServerError)
			return
		}

		// Получаем метаданные из формы
		title := r.FormValue("title")
		artist := r.FormValue("artist")

		// Создаем трек
		track := &model.Track{
			Metadata: model.Metadata{
				Title:    title,
				Artist:   artist,
				Duration: 0, // Можно вычислить длительность через ffmpeg
			},
		}

		// Сохраняем в MinIO и БД
		createdTrack, err := s.store.Track().Create(track, audioData)
		if err != nil {
			s.logger.Errorf("Failed to save track: %v", err)
			http.Error(w, "Failed to save track", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		if err := json.NewEncoder(w).Encode(createdTrack); err != nil {
			s.logger.Errorf("Failed to encode response: %v", err)
		}
	}
}
