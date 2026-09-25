package domain

// GenerateDMARCRecord returns a recommended DMARC TXT value for a domain.
func (s *Service) GenerateDMARCRecord(domain, ruaEmail string) string {
	// Basic DMARC policy: quarantine suspicious, send reports
	// _dmarc.<domain> TXT "v=DMARC1; p=quarantine; rua=mailto:reports@domain"
	return "v=DMARC1; p=quarantine; rua=mailto:" + ruaEmail
}
