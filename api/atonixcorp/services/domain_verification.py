"""Bounded verification for organization contact and website domains."""
from __future__ import annotations

import ipaddress
import socket
import ssl
from urllib.error import HTTPError, URLError
from urllib.parse import urlparse
from urllib.request import HTTPSHandler, HTTPRedirectHandler, Request, build_opener

import dns.exception
import dns.resolver
from django.core.validators import validate_email
from django.core.exceptions import ValidationError as DjangoValidationError


PUBLIC_EMAIL_DOMAINS = frozenset({
    'aol.com', 'gmail.com', 'gmx.com', 'hotmail.com', 'icloud.com', 'live.com',
    'mail.com', 'msn.com', 'outlook.com', 'proton.me', 'protonmail.com',
    'yahoo.com', 'yandex.com', 'zoho.com',
})
DNS_TIMEOUT_SECONDS = 3
HTTPS_TIMEOUT_SECONDS = 5


def _result(status, reason, **details):
    return {'status': status, 'reason': reason, **details}


def _normalize_domain(value):
    domain = str(value or '').strip().rstrip('.').lower()
    if not domain or len(domain) > 253:
        return None
    try:
        normalized = domain.encode('idna').decode('ascii')
    except UnicodeError:
        return None
    labels = normalized.split('.')
    if len(labels) < 2 or any(not label or len(label) > 63 or not label.replace('-', '').isalnum() or label.startswith('-') or label.endswith('-') for label in labels):
        return None
    return normalized


def _website_host(website):
    raw_website = str(website or '').strip()
    parsed = urlparse(raw_website)
    if parsed.scheme.lower() != 'https' or not parsed.hostname or parsed.username or parsed.password:
        return None
    try:
        port = parsed.port
    except ValueError:
        return None
    if port not in (None, 443):
        return None
    return _normalize_domain(parsed.hostname)


def verify_email_domain(email):
    """Verify a corporate email address and its MX records."""
    normalized_email = str(email or '').strip().lower()
    try:
        validate_email(normalized_email)
    except DjangoValidationError:
        return _result('fail', 'Email must use a valid corporate address.')

    domain = _normalize_domain(normalized_email.rsplit('@', 1)[-1])
    if not domain:
        return _result('fail', 'Email must use a valid corporate domain.')
    if domain in PUBLIC_EMAIL_DOMAINS:
        return _result('fail', 'Email must use a corporate domain, not a public email provider.', domain=domain)

    try:
        mx_records = dns.resolver.resolve(domain, 'MX', lifetime=DNS_TIMEOUT_SECONDS)
        if not list(mx_records):
            return _result('fail', 'Invalid MX record.', domain=domain)
    except (dns.resolver.NXDOMAIN, dns.resolver.NoAnswer, dns.resolver.NoNameservers, dns.exception.Timeout):
        return _result('fail', 'Invalid MX record.', domain=domain)

    return _result('success', 'Corporate email domain verified.', domain=domain)


def verify_website_domain(website):
    """Verify an externally reachable HTTPS website with a valid certificate."""
    domain = _website_host(website)
    if not domain:
        return _result('fail', 'Website must use a valid HTTPS URL on port 443.')

    try:
        addresses = {
            item[4][0]
            for item in socket.getaddrinfo(domain, 443, type=socket.SOCK_STREAM)
        }
    except socket.gaierror:
        return _result('fail', 'Website domain does not resolve to an IP address.', domain=domain)

    if not addresses or any(not ipaddress.ip_address(address).is_global for address in addresses):
        return _result('fail', 'Website domain must resolve to a public IP address.', domain=domain)

    class NoRedirect(HTTPRedirectHandler):
        def redirect_request(self, request, fp, code, message, headers, newurl):
            return None

    request = Request(
        f'https://{domain}/',
        headers={'User-Agent': 'AtonixCorp-Domain-Verification/1.0'},
        method='GET',
    )
    opener = build_opener(NoRedirect(), HTTPSHandler(context=ssl.create_default_context()))
    try:
        with opener.open(request, timeout=HTTPS_TIMEOUT_SECONDS) as response:
            if response.status != 200:
                return _result('fail', 'Website must respond with HTTPS status 200.', domain=domain)
    except HTTPError as error:
        return _result('fail', f'Website must respond with HTTPS status 200 (received {error.code}).', domain=domain)
    except (URLError, ssl.SSLError, TimeoutError, OSError):
        return _result('fail', 'Website HTTPS connection or TLS certificate could not be verified.', domain=domain)

    return _result('success', 'Website domain verified.', domain=domain)


def match_domain_email_website(email_domain, website_domain):
    """Accept an exact domain match or the conventional www subdomain."""
    email_domain = _normalize_domain(email_domain)
    website_domain = _normalize_domain(website_domain)
    if not email_domain or not website_domain:
        return _result('fail', 'Website domain does not match email domain.')
    if website_domain == email_domain or website_domain == f'www.{email_domain}':
        return _result('success', 'Email and website domains match.', domain=email_domain)
    return _result('fail', 'Website domain does not match email domain.', email_domain=email_domain, website_domain=website_domain)


def verify_organization_domains(email, website):
    """Return structured outcomes for every organization-domain verification step."""
    email_result = verify_email_domain(email)
    website_result = verify_website_domain(website)
    match_result = match_domain_email_website(
        email_result.get('domain'),
        website_result.get('domain'),
    )
    failed = next(
        (result for result in (email_result, website_result, match_result) if result['status'] != 'success'),
        None,
    )
    return {
        'status': 'fail' if failed else 'success',
        'reason': failed['reason'] if failed else 'Organization domains verified.',
        'email': email_result,
        'website': website_result,
        'match': match_result,
    }


# Compatibility aliases for the service contract used in architecture documentation.
verifyEmailDomain = verify_email_domain
verifyWebsiteDomain = verify_website_domain
matchDomainEmailWebsite = match_domain_email_website
