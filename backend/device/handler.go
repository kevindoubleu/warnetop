package device

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"
)

const urlParamKeyID = "id"

type DeviceHandler struct {
	service *DeviceService
}

func NewDeviceHandler(s *DeviceService) *DeviceHandler {
	return &DeviceHandler{
		service: s,
	}
}

func (h *DeviceHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.deviceGet(w, r)
	case http.MethodPost:
		h.devicePost(w, r)
	case http.MethodPut:
		h.devicePut(w, r)
	case http.MethodDelete:
		h.deviceDelete(w, r)
	}
}

func (h *DeviceHandler) deviceGet(w http.ResponseWriter, r *http.Request) {
	deviceID, err := strconv.Atoi(r.FormValue(urlParamKeyID))
	if err != nil {
		http.Error(w, "invalid param: id", http.StatusBadRequest)
		return
	}

	retrieved, err := h.service.getDevice(deviceID)
	if err == ErrNotFound {
		http.Error(w, ErrNotFound.Error(), http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(retrieved)
}

func (h *DeviceHandler) devicePost(w http.ResponseWriter, r *http.Request) {
	newDevice, err := parseDeviceJSON(r, false)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	persistedDevice, err := h.service.createNewDevice(*newDevice)
	if err == ErrCreateFail {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(persistedDevice)
}

func (h *DeviceHandler) devicePut(w http.ResponseWriter, r *http.Request) {
	newDevice, err := parseDeviceJSON(r, true)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	updatedDevice, err := h.service.updateDevice(*newDevice)
	if err == ErrNotFound {
		http.Error(w, ErrNotFound.Error(), http.StatusNotFound)
		return
	}
	if err == ErrUpdateFail {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(updatedDevice)
}

func (h *DeviceHandler) deviceDelete(w http.ResponseWriter, r *http.Request) {
	deviceID, err := strconv.Atoi(r.FormValue(urlParamKeyID))
	if err != nil {
		http.Error(w, "invalid param: id", http.StatusBadRequest)
		return
	}

	deleted, err := h.service.deleteDevice(deviceID)
	if err == ErrNotFound {
		http.Error(w, ErrNotFound.Error(), http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(deleted)
}

// Parses request body for a Device in JSON format
func parseDeviceJSON(r *http.Request, withId bool) (*Device, error) {
	parsedDevice := &Device{}
	err := json.NewDecoder(r.Body).Decode(parsedDevice)
	if err != nil {
		log.Printf("Error in decoding json: %s", err)
		return nil, errors.New("invalid json")
	}

	if withId && parsedDevice.ID == 0 {
		return nil, errors.New("invalid param: id")
	}
	if parsedDevice.Rate == 0 {
		return nil, errors.New("invalid param: rate")
	}
	if parsedDevice.Model == "" {
		return nil, errors.New("invalid param: model")
	}

	return parsedDevice, nil
}
