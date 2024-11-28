package nullsignature

import (
	"github.com/cerberauth/vulnapi/internal/auth"
	"github.com/cerberauth/vulnapi/internal/operation"
	"github.com/cerberauth/vulnapi/internal/scan"
	"github.com/cerberauth/vulnapi/jwt"
	"github.com/cerberauth/vulnapi/report"
)

const (
	NullSignatureScanID   = "jwt.null_signature"
	NullSignatureScanName = "JWT Null Signature"
)

var issue = report.Issue{
	ID:   "broken_authentication.null_signature",
	Name: "JWT Token has a null signature",
	URL:  "https://vulnapi.cerberauth.com/docs/vulnerabilities/broken-authentication/jwt-null-signature?utm_source=vulnapi",

	Classifications: &report.Classifications{
		OWASP: report.OWASP_2023_BrokenAuthentication,
		CWE:   report.CWE_345_Insufficient_Verification_Authenticity,
	},

	CVSS: report.CVSS{
		Version: 4.0,
		Vector:  "CVSS:4.0/AV:N/AC:L/AT:N/PR:N/UI:N/VC:H/VI:H/VA:N/SC:N/SI:N/SA:N",
		Score:   9.3,
	},
}

func ShouldBeScanned(securityScheme *auth.SecurityScheme) bool {
	return securityScheme != nil && securityScheme.GetType() != auth.None && (securityScheme.GetTokenFormat() == nil || *securityScheme.GetTokenFormat() == auth.JWTTokenFormat)
}

func ScanHandler(op *operation.Operation, securityScheme *auth.SecurityScheme) (*report.ScanReport, error) {
	vulnReport := report.NewIssueReport(issue).WithOperation(op).WithSecurityScheme(securityScheme)
	r := report.NewScanReport(NullSignatureScanID, NullSignatureScanName, op)

	if !ShouldBeScanned(securityScheme) {
		r.AddIssueReport(vulnReport.Skip()).End()
		return r, nil
	}

	var token string
	if securityScheme.HasValidValue() {
		token = securityScheme.GetToken()
	} else {
		token = jwt.FakeJWT
	}

	valueWriter, err := jwt.NewJWTWriter(token)
	if err != nil {
		return r, err
	}

	newToken, err := valueWriter.WithoutSignature()
	if err != nil {
		return r, err
	}
	if err = securityScheme.SetAttackValue(newToken); err != nil {
		return r, err
	}
	vsa, err := scan.ScanURL(op, securityScheme)
	if err != nil {
		return r, err
	}
	r.AddScanAttempt(vsa).End()
	vulnReport.WithBooleanStatus(scan.IsUnauthorizedStatusCodeOrSimilar(vsa.Response))
	r.AddIssueReport(vulnReport)

	return r, nil
}
