 package routes

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	chimw "github.com/go-chi/chi/v5/middleware"
	"github.com/manassp21/drivex/backend/internal/auth"
	"github.com/manassp21/drivex/backend/database"
	authmw "github.com/manassp21/drivex/backend/internal/middleware"
	"github.com/manassp21/drivex/backend/internal/models"
	"github.com/google/uuid"
	"io"
	"log"
	"strings"

	"github.com/manassp21/drivex/backend/internal/storage"
	"github.com/manassp21/drivex/backend/internal/middleware"

)

var fileStorage storage.Storage

type RegisterRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type CreateFolderRequest struct {
	Name     string  `json:"name"`
	ParentID *string `json:"parent_id,omitempty"`
}

type LoginResponse struct {
	User  *models.User `json:"user"`
	Token string       `json:"token"`
}

func NewRouter(store storage.Storage) *chi.Mux {
	fileStorage = store
	r := chi.NewRouter()

	r.Use(middleware.RateLimit)

	r.Use(chimw.Logger)
	r.Use(chimw.Recoverer)

	r.Get("/health", healthHandler)

	r.Route("/api/users", func(r chi.Router) {
		r.Post("/register", RegisterUser)
		r.Post("/login", loginUser)
		r.Get("/{id}", getUser)
	})

	r.Group(func(r chi.Router) {
		r.Use(authmw.AuthMiddleware)

		r.Route("/api/folders", func(r chi.Router) {
			r.Get("/", listFolders)
			r.Post("/", createFolder)
			r.Get("/{id}", getFolder)
			r.Delete("/{id}", deleteFolder)
		})

		r.Route("/api/files", func(r chi.Router) {
			r.Get("/", listFiles)
			r.Post("/", uploadFile)
			r.Get("/{id}", getFile)
			r.Get("/{id}/download", downloadFile)
			r.Delete("/{id}", deleteFile)
		})
	})

	return r
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("ok"))
}

func RegisterUser(w http.ResponseWriter, r *http.Request) {
	var req RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	if len(req.Password) < 8{
		http.Error(w, "password must contain atleast  characters", http.StatusBadRequest)
		return
	}

	if !strings.Contains(req.Email, "@") {
    	http.Error(w, "invalid email", http.StatusBadRequest)
    	return
	}

	ctx := r.Context()
	pwHash, err := auth.HashPassword(req.Password)
	if err != nil {
		http.Error(w, "error hashing password", http.StatusInternalServerError)
		return
	}


user, err := database.CreateUser(ctx, req.Email, pwHash)
if err != nil {
    log.Println("CreateUser error:", err)
    http.Error(w, "internal server error", http.StatusInternalServerError)
    return
}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(user)
}

func loginUser(w http.ResponseWriter, r *http.Request) {
	var req RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	ctx := r.Context()
	user, err := database.GetUserByEmail(ctx, req.Email)
	if err != nil {
		http.Error(w, "user not found", http.StatusBadRequest)
		return
	}

	if err := auth.ComparePassword(req.Password, user.PasswordHash); err != nil {
		http.Error(w, "password doesnt match", http.StatusUnauthorized)
		return
	}

	token, err := auth.GenerateToken(user.ID)
	if err != nil {
		http.Error(w, "failed to generate token", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(LoginResponse{User: user, Token: token})
}

func getUser(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	ctx := r.Context()
	user, err := database.GetUserByID(ctx, id)
	if err != nil {
		http.Error(w, "user not found", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(user)
}

func createFolder(w http.ResponseWriter, r *http.Request) {
	ownerID, ok := authmw.GetUserID(r)
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	var req CreateFolderRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}
	if req.Name == "" {
		http.Error(w, "folder name is required", http.StatusBadRequest)
		return
	}

	folder, err := database.CreateFolder(r.Context(), req.Name, ownerID, req.ParentID)
	if err != nil {
		http.Error(w, "failed to create folder", http.StatusInternalServerError)
		return
	}


	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(folder)
}

func listFolders(w http.ResponseWriter, r *http.Request) {
	ownerID, ok := authmw.GetUserID(r)
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	var parentID *string
	if pid := r.URL.Query().Get("parent_id"); pid != "" {
		parentID = &pid
	}

	folders, err := database.ListFolders(r.Context(), ownerID, parentID)
	if err != nil {
		http.Error(w, "failed to list folders", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(folders)
}

func getFolder(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	folder, err := database.GetFolder(r.Context(), id)
	if err != nil {
		http.Error(w, "folder not found", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(folder)
}

func deleteFolder(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if err := database.DeleteFolder(r.Context(), id); err != nil {
		http.Error(w, "folder not found", http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func listFiles(w http.ResponseWriter, r *http.Request)    {
	ownerID, ok := authmw.GetUserID(r)
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	var folderID *string
	if fid := r.URL.Query().Get("folder_id"); fid != "" {
		folderID = &fid
	}

	files, err := database.ListFiles(r.Context(), ownerID, folderID)
	if err != nil {
		http.Error(w, "failed to list files", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(files)
}

func uploadFile(w http.ResponseWriter, r *http.Request)   {
	ownerID, ok := authmw.GetUserID(r)
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	r.Body = http.MaxBytesReader(w, r.Body, 25<<20) // 25 MB
	if err := r.ParseMultipartForm(25 << 20); err != nil {
		http.Error(w, "file too large", http.StatusRequestEntityTooLarge)
	    return
	}

	if err := r.ParseMultipartForm(32 << 20); err != nil {
		http.Error(w, "could not parse form", http.StatusBadRequest)
		return
	}

	file, header, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "file field is required", http.StatusBadRequest)
		return
	}
	defer file.Close()

	var folderID *string
	if fid := r.FormValue("folder_id"); fid != "" {
		folderID = &fid
	}

	storageKey := uuid.NewString() + "_" + header.Filename

	mimeType := header.Header.Get("Content-Type")
	if mimeType == "" {
		mimeType = "application/octet-stream"
	}

	if err := fileStorage.Upload(r.Context(), storageKey, file, header.Size, mimeType); err != nil {
		log.Println("s3 upload error:", err)
		http.Error(w, "failed to save file", http.StatusInternalServerError)
		return
	}

	fileRecord, err := database.CreateFile(r.Context(), header.Filename, ownerID, folderID, header.Size, mimeType, storageKey)

	if err != nil {
		http.Error(w, "failed to save file record", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(fileRecord)
}

func getFile(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	file, err := database.GetFile(r.Context(), id)
	if err != nil {
		http.Error(w, "file not found", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(file)
}

func downloadFile(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	file, err := database.GetFile(r.Context(), id)
	if err != nil {
		http.Error(w, "file not found", http.StatusNotFound)
		return
	}

		body, err := fileStorage.Download(r.Context(), file.StorageKey)
	if err != nil {
		log.Println("s3 download error:", err)
		http.Error(w, "file missing from storage", http.StatusInternalServerError)
		return
	}
	defer body.Close()

	w.Header().Set("Content-Disposition", "attachment; filename=\""+file.Name+"\"")
	w.Header().Set("Content-Type", file.MimeType)
	io.Copy(w, body)

}

func deleteFile(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	file, err := database.GetFile(r.Context(), id)
	if err != nil {
		http.Error(w, "file not found", http.StatusNotFound)
		return
	}

	if err := database.DeleteFile(r.Context(), id); err != nil {
		http.Error(w, "failed to delete file record", http.StatusInternalServerError)
		return
	}

	if err := fileStorage.Delete(r.Context(), file.StorageKey); err != nil {
		log.Println("s3 delete error:", err)
	}

	w.WriteHeader(http.StatusNoContent)
}