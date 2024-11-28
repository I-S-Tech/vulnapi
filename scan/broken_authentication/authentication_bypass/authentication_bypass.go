package authenticationbypass

import (
	"github.com/cerberauth/vulnapi/internal/auth"
	"github.com/cerberauth/vulnapi/internal/operation"
	"github.com/cerberauth/vulnapi/internal/scan"
	"github.com/cerberauth/vulnapi/report"
)

const (
	AcceptsUnauthenticatedOperationScanID   = "generic.accept_unauthenticated_operation"
	AcceptsUnauthenticatedOperationScanName = "Accept Unauthenticated Operation"
)

var issue = report.Issue{
	ID:   "broken_authentication.authentication_bypass",
	Name: "Authentication is expected but can be bypassed",

	Classifications: &report.Classifications{
		OWASP: report.OWASP_2023_BrokenAuthentication,
	},

	CVSS: report.CVSS{
		Version: 4.0,
		Vector:  "CVSS:4.0/AV:N/AC:L/AT:N/PR:N/UI:N/VC:H/VI:H/VA:N/SC:N/SI:N/SA:N",
		Score:   9.3,
	},
}

func ScanHandler(op *operation.Operation, securityScheme *auth.SecurityScheme) (*report.ScanReport, error) {
	vulnReport := report.NewIssueReport(issue).WithOperation(op).WithSecurityScheme(securityScheme)

	r := report.NewScanReport(AcceptsUnauthenticatedOperationScanID, AcceptsUnauthenticatedOperationScanName, op)
	if securityScheme.GetType() == auth.None {
		return r.AddIssueReport(vulnReport.Skip()).End(), nil
	}

	noAuthSecurityScheme := auth.MustNewNoAuthSecurityScheme()
	vsa, err := scan.ScanURL(op, noAuthSecurityScheme)
	if err != nil {
		return r, err
	}
	vulnReport.WithBooleanStatus(scan.IsUnauthorizedStatusCodeOrSimilar(vsa.Response))
	r.AddIssueReport(vulnReport).AddScanAttempt(vsa).End()

	return r, nil
}
