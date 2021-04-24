package handlers

import "net/http"

// GRCHandler contains fields shared among the different handlers
type GRCHandler struct {
	AssetsFS http.FileSystem
}
