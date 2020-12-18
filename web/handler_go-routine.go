package web

import (
	"encoding/json"
	"net/http"
)

func goRoutineHandler(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	var data []int
	for i := 0; i < 100000; i++ {
		if i%1000 == 0 {
			data = append(data, i)
		}
	}

	bytes, err := json.Marshal(data)
	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}

	w.Write(bytes)
}