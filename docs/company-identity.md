# Company Identity

`Organization` is the platform's company record. Its database `id` is the internal company ID and its `registration_number` is the canonical external company identity.

## Registration Rules

- A company name is globally unique, case-insensitively.
- A registration number is globally unique after normalization.
- Registration APIs require a company name. A registration number is optional and
	is normalized and uniqueness-checked when provided.
- Allowed input characters are letters, digits, spaces, periods, underscores, slashes, and hyphens.
- Separators are removed and the identifier is uppercased before storage. For example, `za-2024 / 123456` becomes `ZA2024123456`.
- The canonical identifier must contain 4-64 characters after normalization and include a digit.

## Domain Verification

New organizations must provide a contact email address and an HTTPS website. Before
creation, the API rejects public email providers, verifies that the corporate email
domain has MX records, verifies that the website domain resolves to a public IP,
and requires a valid TLS connection with an HTTPS `200` response. The website host
must match the email domain or its `www` subdomain.

Verification is performed by `atonixcorp.services.domain_verification` for both
`/api/organizations/` and the v1 organization endpoint. Each attempt creates a
platform audit event without storing the raw contact email address. DNS and website
failures are returned as field-level errors so the creation interface can explain
what needs correction.

## Identity Scope

Users and roles remain company-scoped through `TeamMember.organization`. The LDAP-compatible directory projects the company registration number into the root DN and user attributes, so Founder and member identities inherit the same company root. Governance YAML recovery verifies the registration number before any records are restored.

## Verification API

`POST /api/organizations/verify_registration_number/` validates syntax and reports whether a supplied normalized identifier is available. For an omitted registration number, it returns a valid `not_provided` result and organization creation proceeds without an identity-verification request. It intentionally reports `external_registry_verified: false` until a country-specific registry provider is configured. A future CIPC, SEC, or other provider must validate the already-normalized identifier and never replace the local uniqueness check.

## Legacy Records

Organizations without a registration number remain readable and can be created while registration information is unavailable. When an external registry provider is enabled in the future, it must support a staged rollout rather than invalidating existing organizations without a number.