package services

import (
	"errors"

	"recursiveDine/internal/repositories"
)

type TableService struct {
	tableRepo *repositories.TableRepository
}

func NewTableService(tableRepo *repositories.TableRepository) *TableService {
	return &TableService{
		tableRepo: tableRepo,
	}
}

func (s *TableService) GetTableByQRCode(qrCode string) (*repositories.Table, error) {
	table, err := s.tableRepo.GetByQRCode(qrCode)
	if err != nil {
		return nil, errors.New("table not found")
	}

	if !table.IsAvailable {
		return nil, errors.New("table is not available")
	}

	return table, nil
}

func (s *TableService) GetTableByID(id uint) (*repositories.Table, error) {
	return s.tableRepo.GetByID(id)
}

func (s *TableService) GetTableByNumber(number int) (*repositories.Table, error) {
	return s.tableRepo.GetByNumber(number)
}

func (s *TableService) GetAllAvailableTables() ([]repositories.Table, error) {
	return s.tableRepo.GetAllAvailable()
}

func (s *TableService) CreateTable(table *repositories.Table) error {
	// Check if table number already exists
	if exists, err := s.tableRepo.IsNumberExists(table.Number); err != nil {
		return errors.New("failed to check table number")
	} else if exists {
		return errors.New("table number already exists")
	}

	// Check if QR code already exists
	if exists, err := s.tableRepo.IsQRCodeExists(table.QRCode); err != nil {
		return errors.New("failed to check QR code")
	} else if exists {
		return errors.New("QR code already exists")
	}

	return s.tableRepo.Create(table)
}

func (s *TableService) UpdateTable(table *repositories.Table) error {
	// Check if table exists
	existing, err := s.tableRepo.GetByID(table.ID)
	if err != nil {
		return errors.New("table not found")
	}

	// Check if table number changed and if new number already exists
	if existing.Number != table.Number {
		if exists, err := s.tableRepo.IsNumberExists(table.Number); err != nil {
			return errors.New("failed to check table number")
		} else if exists {
			return errors.New("table number already exists")
		}
	}

	// Check if QR code changed and if new QR code already exists
	if existing.QRCode != table.QRCode {
		if exists, err := s.tableRepo.IsQRCodeExists(table.QRCode); err != nil {
			return errors.New("failed to check QR code")
		} else if exists {
			return errors.New("QR code already exists")
		}
	}

	return s.tableRepo.Update(table)
}

func (s *TableService) DeleteTable(id uint) error {
	return s.tableRepo.Delete(id)
}

func (s *TableService) SetTableAvailability(id uint, available bool) error {
	table, err := s.tableRepo.GetByID(id)
	if err != nil {
		return errors.New("table not found")
	}

	table.IsAvailable = available
	return s.tableRepo.Update(table)
}
