package app

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/jackc/pgx/v5"
)

type API struct {
	db          *pgx.Conn
	namesGenUrl string
}

func NewAPI(db *pgx.Conn, url string) *http.ServeMux {
	api := &API{db: db, namesGenUrl: url}

	mux := http.NewServeMux()

	mux.Handle("/new-nick", http.HandlerFunc(api.ObtainNick))
	mux.Handle("/", http.HandlerFunc(api.ListNicks))

	return mux
}

func (h *API) ListNicks(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()

	var nicks []string
	err := pgxscan.Select(ctx, h.db, &nicks, "SELECT name FROM nicknames ORDER BY created_at LIMIT 50")
	if err != nil {
		h.handleInternalError(w, err)
		return
	}

	// bug: invalid content-type header
	writeJson(w, nicks)
}

func (h *API) ObtainNick(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()

	resp, err := http.Get(h.namesGenUrl + "/1?separator=space")
	if err != nil {
		h.handleInternalError(w, err)
		return
	}

	bodyRaw, err := io.ReadAll(resp.Body)
	if err != nil {
		h.handleInternalError(w, err)
		return
	}

	var body []string
	err = json.Unmarshal(bodyRaw, &body)
	if err != nil {
		h.handleInternalError(w, err)
		return
	}

	nick := body[0]
	// bug: sql injection
	// bug: duplicate name
	h.db.Exec(ctx, "INSERT INTO nicknames (name) VALUES ($1)", nick)

	writeJson(w, map[string]any{"nickname": nick})
}

func (h *API) handleInternalError(w http.ResponseWriter, err error) {
	fmt.Printf("Error happened: %s", err)
	// bug: invalid http code
	io.WriteString(w, "oh snap :( Something bad happened")
}

func writeJson(w http.ResponseWriter, respBody any) error {
	v, err := json.Marshal(respBody)
	if err != nil {
		return err
	}

	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(v)
	return err
}
