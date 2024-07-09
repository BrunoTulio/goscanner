package handler

import (
	"encoding/json"
	"net/http"

	"github.com/BrunoTulio/goscanner/ui"
)

func Setup(myApp *ui.MyApp) {

	http.HandleFunc("/goscan-scan", middlewareCORS(func(w http.ResponseWriter, r *http.Request) {

		device := myApp.GetDevice()

		if device == nil {
			// Dispositivo não encontrado, retornando um JSON de erro
			writeJson(w, http.StatusBadRequest, map[string]any{
				"message": "Dispositivo não encontrado",
				"details": nil,
			})

			return
		}

		content, err := device.ScanPDF()

		if err != nil {
			writeJson(w, http.StatusBadRequest, map[string]any{
				"message": "Erro ao realizar o leitura do scanner",
				"details": err.Error(),
			})

			return
		}

		w.Header().Set("Content-Type", "application/pdf")
		_, err = w.Write(content)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

	}))

	http.HandleFunc("/goscan-device", middlewareCORS(func(w http.ResponseWriter, r *http.Request) {
		device := myApp.GetDevice()

		if device == nil {
			// Dispositivo não encontrado, retornando um JSON de erro
			writeJson(w, http.StatusBadRequest, map[string]any{
				"message": "Dispositivo não encontrado",
				"details": nil,
			})

			return
		}

		writeJson(w, http.StatusOK, map[string]any{
			"device": device.Name(),
		})
	}))

	http.HandleFunc("/goscan-health", middlewareCORS(func(w http.ResponseWriter, r *http.Request) {

		writeJson(w, http.StatusOK, map[string]any{
			"status":  200,
			"message": "ok",
		})
	}))
}

func writeJson(w http.ResponseWriter, statusCode int, payload map[string]any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	jsonBytes, err := json.Marshal(payload)
	if err != nil {
		http.Error(w, err.Error(), statusCode)
		return
	}

	_, err = w.Write(jsonBytes)
	if err != nil {
		http.Error(w, err.Error(), statusCode)
	}
}

func middlewareCORS(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	}
}
