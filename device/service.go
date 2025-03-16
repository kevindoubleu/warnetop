package device

import "log"

type DeviceService struct {
	repo DeviceRepository
}

func NewDeviceService(r DeviceRepository) *DeviceService {
	return &DeviceService{
		repo: r,
	}
}

func (s *DeviceService) createNewDevice(d Device) (*Device, error) {
	log.Printf("Creating new Device %s", d)
	return s.repo.Create(d)
}

func (s *DeviceService) getDevice(id int) (*Device, error) {
	retrieved, err := s.repo.Read(id)
	if err == ErrNotFound {
		return nil, ErrNotFound
	}

	return retrieved, nil
}
