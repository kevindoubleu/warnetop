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
	return s.repo.Read(id)
}

func (s *DeviceService) updateDevice(d Device) (*Device, error) {
	log.Printf("Updating Device %s", d)
	return s.repo.Update(d)
}

func (s *DeviceService) deleteDevice(id int) (*Device, error) {
	log.Printf("Deleting Device with ID: %d", id)
	return s.repo.Delete(id)
}
