package places

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func GetSaloonSlots(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.WriteHeader(http.StatusNotImplemented)
	_, _ = w.Write([]byte("GetSaloonSlots not implemented yet"))
}

func PostSaloonSlot(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.WriteHeader(http.StatusNotImplemented)
	_, _ = w.Write([]byte("PostSaloonSlot not implemented yet"))
}

func PutSaloonSlot(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.WriteHeader(http.StatusNotImplemented)
	_, _ = w.Write([]byte("PutSaloonSlot not implemented yet"))
}

func DeleteSaloonSlot(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.WriteHeader(http.StatusNotImplemented)
	_, _ = w.Write([]byte("DeleteSaloonSlot not implemented yet"))
}

func BookSaloonSlot(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.WriteHeader(http.StatusNotImplemented)
	_, _ = w.Write([]byte("BookSaloonSlot not implemented yet"))
}
