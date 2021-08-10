package app

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/DABronskikh/bgo-3_14.1/pkg/appErr"
	"github.com/DABronskikh/bgo-3_14.1/pkg/dto"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"log"
	"net/http"
	"net/url"
)

type Server struct {
	mux  *http.ServeMux
	ctx  context.Context
	conn *pgxpool.Conn
}

func NewServer(mux *http.ServeMux, ctx context.Context, conn *pgxpool.Conn) *Server {
	return &Server{mux: mux, ctx: ctx, conn: conn}
}

func (s *Server) Init() {
	s.mux.HandleFunc("/getCards", s.getCards)
	s.mux.HandleFunc("/getTransactions", s.getTransactions)
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.mux.ServeHTTP(w, r)
}

func (s *Server) getCards(w http.ResponseWriter, r *http.Request) {
	userId, ok := url.Parse(r.URL.Query().Get("userId"))
	userIdstr := fmt.Sprintf("%v", userId)
	if ok != nil || userIdstr == "" {
		dtos := dto.ErrDTO{Err: appErr.ErrUserIdNotFound.Error()}
		prepareResponseErr(w, r, dtos)
		return
	}

	cardDB := []*dto.CardDTO{}
	rows, err := s.conn.Query(s.ctx, `
		SELECT id, number, balance, issuer, holder, owner_id, status, created
		FROM cards 
		WHERE owner_id = $1
		LIMIT 50
	`, userIdstr)

	defer rows.Close()

	for rows.Next() {
		cardEl := &dto.CardDTO{}
		err = rows.Scan(&cardEl.Id, &cardEl.Number, &cardEl.Balance, &cardEl.Issuer, &cardEl.Holder, &cardEl.OwnerId, &cardEl.Status, &cardEl.Created)
		if err != nil {
			dtos := dto.ErrDTO{Err: appErr.ErrDB.Error()}
			prepareResponseErr(w, r, dtos)
			return
		}
		cardDB = append(cardDB, cardEl)
	}

	if err != nil {
		if err != pgx.ErrNoRows {
			dtos := dto.ErrDTO{Err: appErr.ErrDB.Error()}
			prepareResponseErr(w, r, dtos)
			return
		}
	}

	prepareResponseCard(w, r, cardDB)
}

func (s *Server) getTransactions(w http.ResponseWriter, r *http.Request) {
	cardId, ok := url.Parse(r.URL.Query().Get("cardId"))
	cardIdstr := fmt.Sprintf("%v", cardId)
	if ok != nil || cardIdstr == "" {
		dtos := dto.ErrDTO{Err: appErr.ErrCardIdNotFound.Error()}
		prepareResponseErr(w, r, dtos)
		return
	}

	transactionDB := []*dto.TransactionDTO{}
	rows, err := s.conn.Query(s.ctx, `
		SELECT id, card_id, amount, created, status, mcc_id, description, supplier_icon_id
		FROM transactions 
		WHERE card_id = $1
		LIMIT 50
	`, cardIdstr)

	defer rows.Close()

	for rows.Next() {
		trEl := &dto.TransactionDTO{}
		err = rows.Scan(&trEl.Id, &trEl.CardId, &trEl.Amount, &trEl.Created, &trEl.Status, &trEl.MccId, &trEl.Description, &trEl.SupplierIconId)
		if err != nil {
			dtos := dto.ErrDTO{Err: appErr.ErrDB.Error()}
			prepareResponseErr(w, r, dtos)
			return
		}
		transactionDB = append(transactionDB, trEl)
	}

	if err != nil {
		if err != pgx.ErrNoRows {
			dtos := dto.ErrDTO{Err: appErr.ErrDB.Error()}
			prepareResponseErr(w, r, dtos)
			return
		}
	}

	prepareResponseTransaction(w, r, transactionDB)
}

func prepareResponseCard(w http.ResponseWriter, r *http.Request, dtos []*dto.CardDTO) {
	respBody, err := json.Marshal(dtos)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	_, err = w.Write(respBody)
	if err != nil {
		log.Println(err)
	}
}

func prepareResponseTransaction(w http.ResponseWriter, r *http.Request, dtos []*dto.TransactionDTO) {
	respBody, err := json.Marshal(dtos)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	_, err = w.Write(respBody)
	if err != nil {
		log.Println(err)
	}
}

func prepareResponseErr(w http.ResponseWriter, r *http.Request, dtos dto.ErrDTO) {
	respBody, err := json.Marshal(dtos)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	_, err = w.Write(respBody)
	if err != nil {
		log.Println(err)
	}
}
