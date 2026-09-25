package domain

// GenerateSPFRecord returns a recommended SPF TXT value for a domain.
func (s *Service) GenerateSPFRecord(domain string, mxHost string) string {
	// Simple SPF: allow our MX + sending IPs
	// You can extend this later with include: etc.
	return "v=spf1 mx a:" + mxHost + " ~all"
}
