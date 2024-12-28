package handler

import (
	"WaterSportsRental/internal/configs/dumpConfig"
	"WaterSportsRental/internal/entity"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"time"
)

func (h *Handler) createDump(w http.ResponseWriter, r *http.Request) {

	dumpCfg := dumpConfig.MustLoadDumpConfig()
	currentTime := time.Now()
	timeString := currentTime.Format("2006-01-02_15-04-05")
	filePath := filepath.Join(dumpCfg.Dir, fmt.Sprintf("%s_%s", dumpCfg.Prefix, timeString))

	cmd := exec.Command("docker", "exec", dumpCfg.ContainerName, "pg_dump", "-U", dumpCfg.Username, "-F", "c", dumpCfg.DbName)
	outputFile, err := os.Create(filePath)
	if err != nil {
		log.Fatal(err)
	}
	defer outputFile.Close()

	cmd.Stdout = outputFile

	if err := cmd.Run(); err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	err = h.services.Dump.InsertDump(r.Context(), filePath)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (h *Handler) restoreDump(w http.ResponseWriter, r *http.Request) {

	var fileName entity.Dump

	if err := json.NewDecoder(r.Body).Decode(&fileName); err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	dumpCfg := dumpConfig.MustLoadDumpConfig()

	currentTime := time.Now()
	timeString := currentTime.Format("2006-01-02_15-04-05")
	filePath := fmt.Sprintf("%s_%s", dumpCfg.RestorePrefix, timeString)

	copyCmd := exec.Command("docker", "cp", fileName.Filename, fmt.Sprintf("%s:%s", dumpCfg.ContainerName, filePath))
	if err := copyCmd.Run(); err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusBadRequest)
		return
	}
	restoreCmd := exec.Command("docker", "exec", dumpCfg.ContainerName, "pg_restore", "-U", dumpCfg.Username, "--clean", "-d", dumpCfg.DbName, filePath)
	if err := restoreCmd.Run(); err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (h *Handler) getAllDumps(w http.ResponseWriter, r *http.Request) {
	dumps, err := h.services.Dump.GetAllDumps(r.Context())
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Add("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(dumps)
}
