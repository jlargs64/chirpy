package handlers

import (
	"sync/atomic"

	"github.com/jlargs64/chirpy/internal/database"
)

type APIConfig struct {
	FileserverHits atomic.Int32
	DBQueries      *database.Queries
}
