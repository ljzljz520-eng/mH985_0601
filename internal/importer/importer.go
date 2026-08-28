package importer

import (
	"fmt"
	"strings"

	"trainingdesk/internal/catalog"
	"trainingdesk/internal/model"
)

type Service struct {
	catalog *catalog.Catalog
}

func New(c *catalog.Catalog) *Service {
	return &Service{catalog: c}
}

func (s *Service) Import(rows []model.ImportRow, startSeq int64) model.ImportReport {
	report := model.ImportReport{Errors: make([]string, 0), IDs: make([]string, 0)}
	for i, row := range rows {
		if strings.TrimSpace(row.ID) == "" {
			report.Rejected++
			report.Errors = append(report.Errors, fmt.Sprintf("row %d: id is required", i+1))
			continue
		}
		seq := startSeq + int64(i)
		if _, err := s.catalog.Register(row, seq); err != nil {
			report.Rejected++
			report.Errors = append(report.Errors, fmt.Sprintf("row %d: %v", i+1, err))
			continue
		}
		report.Accepted++
		report.IDs = append(report.IDs, row.ID)
	}
	return report
}

func (s *Service) ValidateRows(rows []model.ImportRow) []string {
	errs := make([]string, 0)
	seen := make(map[string]bool)
	for i, row := range rows {
		if seen[row.ID] {
			errs = append(errs, fmt.Sprintf("row %d: duplicate id", i+1))
		}
		seen[row.ID] = true
		if row.StoreID == "" {
			errs = append(errs, fmt.Sprintf("row %d: store id is required", i+1))
		}
		if row.SortKey < 0 {
			errs = append(errs, fmt.Sprintf("row %d: sort key is negative", i+1))
		}
	}
	return errs
}
