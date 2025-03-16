package device

import (
	"encoding/json"
	"errors"
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
	newDevice := &Device{}
	err := json.NewDecoder(r.Body).Decode(newDevice)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if newDevice.Rate == 0 {
		http.Error(w, errors.New("invalid param: rate").Error(), http.StatusBadRequest)
		return
	}
	if newDevice.Model == "" {
		http.Error(w, errors.New("invalid param: model").Error(), http.StatusBadRequest)
		return
	}

	persistedDevice, err := h.service.createNewDevice(*newDevice)
	if err == ErrCreateFail {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(persistedDevice)
}
