package main

import (
	"bytes"
	"crypto"
	"crypto/x509"
	"crypto/x509/pkix"
	"dbgutil"
	"errors"
	"fileop"
	"fmt"
	"jsonext"
	"logutil"
	"net"
	"net/url"
	"reflect"
	"strings"
	"time"
)

// VerifyOptions contains parameters for Certificate.Verify.
type VerifyOptionsF struct {
	// DNSName, if set, is checked against the leaf certificate with
	// Certificate.VerifyHostname or the platform verifier.
	DNSName string

	// Intermediates is an optional pool of certificates that are not trust
	// anchors, but can be used to form a chain from the leaf certificate to a
	// root certificate.
	Intermediates *CertPool
	// Roots is the set of trusted root certificates the leaf certificate needs
	// to chain up to. If nil, the system roots or the platform verifier are used.
	Roots *CertPool

	// CurrentTime is used to check the validity of all certificates in the
	// chain. If zero, the current time is used.
	CurrentTime time.Time

	// KeyUsages specifies which Extended Key Usage values are acceptable. A
	// chain is accepted if it allows any of the listed values. An empty list
	// means ExtKeyUsageServerAuth. To accept any key usage, include ExtKeyUsageAny.
	KeyUsages []x509.ExtKeyUsage

	// MaxConstraintComparisions is the maximum number of comparisons to
	// perform when checking a given certificate's name constraints. If
	// zero, a sensible default is used. This limit prevents pathological
	// certificates from consuming excessive amounts of CPU time when
	// validating. It does not apply to the platform verifier.
	MaxConstraintComparisions int
}

func load_certs_value(cert *CertPool, arrs []string, note string) (err error) {
	var i int
	var s string
	var pemb []byte
	var ok bool
	for i, s = range arrs {
		pemb, err = fileop.ReadFileBytes(s)
		if err != nil {
			return
		}
		ok = cert.AppendCertsFromPEM(pemb)
		if !ok {
			err = dbgutil.FormatError("[%s].[%d] [%s] load cert error", note, i, s)
			return
		}
	}
	err = nil
	return
}

func get_x509_verify_options(cfgfile string) (retv VerifyOptionsF, err error) {
	var s string
	var vmap map[string]interface{}
	var vals string
	var valarr []interface{}
	var arrs []string
	var ok bool
	var intval interface{}
	s, err = fileop.ReadFile(cfgfile)
	if err != nil {
		return
	}
	vmap, err = jsonext.GetJsonMap(s)
	if err != nil {
		return
	}

	retv = VerifyOptionsF{}
	retv.Intermediates = NewCertPool()
	retv.Roots = NewCertPool()

	vals, ok = vmap[KEYWORD_DNSNAME].(string)
	if ok {
		logutil.Debug("[%s] parsed", KEYWORD_DNSNAME)
		retv.DNSName = vals
	}

	valarr, ok = vmap[KEYWORD_ROOTS].([]interface{})
	if ok {
		logutil.Debug("[%s] parsed", KEYWORD_ROOTS)
		arrs, err = trans_inter_to_string(valarr, KEYWORD_ROOTS)
		if err != nil {
			return
		}
		err = load_certs_value(retv.Roots, arrs, KEYWORD_ROOTS)
		if err != nil {
			return
		}
	}

	valarr, ok = vmap[KEYWORD_INTERMEDIATES].([]interface{})
	if ok {
		logutil.Debug("[%s] parsed", KEYWORD_INTERMEDIATES)
		arrs, err = trans_inter_to_string(valarr, KEYWORD_INTERMEDIATES)
		if err != nil {
			return
		}
		err = load_certs_value(retv.Intermediates, arrs, KEYWORD_INTERMEDIATES)
		if err != nil {
			return
		}
	}

	s, ok = vmap[KEYWORD_CURRENTTIME].(string)
	if ok {
		logutil.Debug("[%s] parsed", KEYWORD_CURRENTTIME)
		retv.CurrentTime, err = get_time_value(s, KEYWORD_CURRENTTIME)
		if err != nil {
			return
		}
	}

	valarr, ok = vmap[KEYWORD_EXT_KEY_USAGE].([]interface{})
	if ok {
		logutil.Debug("[%s] parsed", KEYWORD_EXT_KEY_USAGE)
		arrs, err = trans_inter_to_string(valarr, KEYWORD_EXT_KEY_USAGE)
		if err != nil {
			return
		}
		retv.KeyUsages, err = get_key_ext_usage(arrs, KEYWORD_EXT_KEY_USAGE)
		if err != nil {
			return
		}
	}

	intval, ok = vmap[KEYWORD_MAX_CONSTRAINTS_COMPARISIONS].(interface{})
	if ok {
		logutil.Debug("[%s] parsed", KEYWORD_MAX_CONSTRAINTS_COMPARISIONS)
		retv.MaxConstraintComparisions, err = get_int_value(intval, KEYWORD_MAX_CONSTRAINTS_COMPARISIONS)
		if err != nil {
			return
		}
	}

	err = nil
	return
}

func hasNameConstraints_Certificate(c *x509.Certificate) bool {
	return oidInExtensions(oidExtensionNameConstraints, c.Extensions)
}

func hasSANExtension_Certificate(c *x509.Certificate) bool {
	return oidInExtensions(oidExtensionSubjectAltName, c.Extensions)
}

func getSANExtension_Certificate(c *x509.Certificate) []byte {
	for _, e := range c.Extensions {
		if e.Id.Equal(oidExtensionSubjectAltName) {
			return e.Value
		}
	}
	return nil
}

// checkNameConstraints checks that c permits a child certificate to claim the
// given name, of type nameType. The argument parsedName contains the parsed
// form of name, suitable for passing to the match function. The total number
// of comparisons is tracked in the given count and should not exceed the given
// limit.
func checkNameConstraints_Certificate(c *x509.Certificate, count *int,
	maxConstraintComparisons int,
	nameType string,
	name string,
	parsedName any,
	match func(parsedName, constraint any) (match bool, err error),
	permitted, excluded any) error {

	excludedValue := reflect.ValueOf(excluded)

	*count += excludedValue.Len()
	if *count > maxConstraintComparisons {
		return x509.CertificateInvalidError{c, x509.TooManyConstraints, ""}
	}

	for i := 0; i < excludedValue.Len(); i++ {
		constraint := excludedValue.Index(i).Interface()
		match, err := match(parsedName, constraint)
		if err != nil {
			return x509.CertificateInvalidError{c, x509.CANotAuthorizedForThisName, err.Error()}
		}

		if match {
			return x509.CertificateInvalidError{c, x509.CANotAuthorizedForThisName, fmt.Sprintf("%s %q is excluded by constraint %q", nameType, name, constraint)}
		}
	}

	permittedValue := reflect.ValueOf(permitted)

	*count += permittedValue.Len()
	if *count > maxConstraintComparisons {
		return x509.CertificateInvalidError{c, x509.TooManyConstraints, ""}
	}

	ok := true
	for i := 0; i < permittedValue.Len(); i++ {
		constraint := permittedValue.Index(i).Interface()

		var err error
		if ok, err = match(parsedName, constraint); err != nil {
			return x509.CertificateInvalidError{c, x509.CANotAuthorizedForThisName, err.Error()}
		}

		if ok {
			break
		}
	}

	if !ok {
		return x509.CertificateInvalidError{c, x509.CANotAuthorizedForThisName, fmt.Sprintf("%s %q is not permitted by any constraint", nameType, name)}
	}

	return nil
}

func matchDomainConstraint_func(domain, constraint string) (bool, error) {
	// The meaning of zero length constraints is not specified, but this
	// code follows NSS and accepts them as matching everything.
	if len(constraint) == 0 {
		return true, nil
	}

	domainLabels, ok := domainToReverseLabels(domain)
	if !ok {
		return false, fmt.Errorf("x509: internal error: cannot parse domain %q", domain)
	}

	// RFC 5280 says that a leading period in a domain name means that at
	// least one label must be prepended, but only for URI and email
	// constraints, not DNS constraints. The code also supports that
	// behaviour for DNS constraints.

	mustHaveSubdomains := false
	if constraint[0] == '.' {
		mustHaveSubdomains = true
		constraint = constraint[1:]
	}

	constraintLabels, ok := domainToReverseLabels(constraint)
	if !ok {
		return false, fmt.Errorf("x509: internal error: cannot parse domain %q", constraint)
	}

	if len(domainLabels) < len(constraintLabels) ||
		(mustHaveSubdomains && len(domainLabels) == len(constraintLabels)) {
		return false, nil
	}

	for i, constraintLabel := range constraintLabels {
		if !strings.EqualFold(constraintLabel, domainLabels[i]) {
			return false, nil
		}
	}

	return true, nil
}

// rfc2821Mailbox represents a “mailbox” (which is an email address to most
// people) by breaking it into the “local” (i.e. before the '@') and “domain”
// parts.
type rfc2821MailboxF struct {
	local, domain string
}

func matchEmailConstraint_func(mailbox rfc2821MailboxF, constraint string) (bool, error) {
	// If the constraint contains an @, then it specifies an exact mailbox
	// name.
	if strings.Contains(constraint, "@") {
		constraintMailbox, ok := parseRFC2821Mailbox(constraint)
		if !ok {
			return false, fmt.Errorf("x509: internal error: cannot parse constraint %q", constraint)
		}
		return mailbox.local == constraintMailbox.local && strings.EqualFold(mailbox.domain, constraintMailbox.domain), nil
	}

	// Otherwise the constraint is like a DNS constraint of the domain part
	// of the mailbox.
	return matchDomainConstraint_func(mailbox.domain, constraint)
}

func matchURIConstraint_func(uri *url.URL, constraint string) (bool, error) {
	// From RFC 5280, Section 4.2.1.10:
	// “a uniformResourceIdentifier that does not include an authority
	// component with a host name specified as a fully qualified domain
	// name (e.g., if the URI either does not include an authority
	// component or includes an authority component in which the host name
	// is specified as an IP address), then the application MUST reject the
	// certificate.”

	host := uri.Host
	if len(host) == 0 {
		return false, fmt.Errorf("URI with empty host (%q) cannot be matched against constraints", uri.String())
	}

	if strings.Contains(host, ":") && !strings.HasSuffix(host, "]") {
		var err error
		host, _, err = net.SplitHostPort(uri.Host)
		if err != nil {
			return false, err
		}
	}

	if strings.HasPrefix(host, "[") && strings.HasSuffix(host, "]") ||
		net.ParseIP(host) != nil {
		return false, fmt.Errorf("URI with IP (%q) cannot be matched against constraints", uri.String())
	}

	return matchDomainConstraint_func(host, constraint)
}

func matchIPConstraint_func(ip net.IP, constraint *net.IPNet) (bool, error) {
	if len(ip) != len(constraint.IP) {
		return false, nil
	}

	for i := range ip {
		if mask := constraint.Mask[i]; ip[i]&mask != constraint.IP[i]&mask {
			return false, nil
		}
	}

	return true, nil
}

// boringAllowCert reports whether c is allowed to be used
// in a certificate chain by the current fipstls enforcement setting.
// It is called for each leaf, intermediate, and root certificate.
func boringAllowCert_func(c *x509.Certificate) bool {
	return true
}

// isValid performs validity checks on c given that it is a candidate to append
// to the chain in currentChain.
func isValid_Certificate(c *x509.Certificate, certType int, currentChain []*x509.Certificate, opts *VerifyOptionsF) error {
	if len(c.UnhandledCriticalExtensions) > 0 {
		logutil.Debug(" ")
		return x509.UnhandledCriticalExtension{}
	}

	if len(currentChain) > 0 {
		child := currentChain[len(currentChain)-1]
		if !bytes.Equal(child.RawIssuer, c.RawSubject) {
			logutil.Debug(" ")
			return x509.CertificateInvalidError{c, x509.NameMismatch, ""}
		}
	}

	now := opts.CurrentTime
	if now.IsZero() {
		now = time.Now()
	}
	if now.Before(c.NotBefore) {
		logutil.Debug(" ")
		return x509.CertificateInvalidError{
			Cert:   c,
			Reason: x509.Expired,
			Detail: fmt.Sprintf("current time %s is before %s", now.Format(time.RFC3339), c.NotBefore.Format(time.RFC3339)),
		}
	} else if now.After(c.NotAfter) {
		logutil.Debug(" ")
		return x509.CertificateInvalidError{
			Cert:   c,
			Reason: x509.Expired,
			Detail: fmt.Sprintf("current time %s is after %s", now.Format(time.RFC3339), c.NotAfter.Format(time.RFC3339)),
		}
	}

	maxConstraintComparisons := opts.MaxConstraintComparisions
	if maxConstraintComparisons == 0 {
		maxConstraintComparisons = 250000
	}
	comparisonCount := 0

	if certType == intermediateCertificate || certType == rootCertificate {
		if len(currentChain) == 0 {
			logutil.Debug(" ")
			return errors.New("x509: internal error: empty chain when appending CA cert")
		}
	}

	if (certType == intermediateCertificate || certType == rootCertificate) &&
		hasNameConstraints_Certificate(c) {
		toCheck := []*x509.Certificate{}
		for _, c := range currentChain {
			if hasSANExtension_Certificate(c) {
				toCheck = append(toCheck, c)
			}
		}
		for _, sanCert := range toCheck {
			err := forEachSAN(getSANExtension_Certificate(sanCert), func(tag int, data []byte) error {
				switch tag {
				case nameTypeEmail:
					name := string(data)
					mailbox, ok := parseRFC2821Mailbox(name)
					if !ok {
						logutil.Debug(" ")
						return fmt.Errorf("x509: cannot parse rfc822Name %q", mailbox)
					}

					if err := checkNameConstraints_Certificate(c, &comparisonCount, maxConstraintComparisons, "email address", name, mailbox,
						func(parsedName, constraint any) (bool, error) {
							return matchEmailConstraint_func(parsedName.(rfc2821MailboxF), constraint.(string))
						}, c.PermittedEmailAddresses, c.ExcludedEmailAddresses); err != nil {
						logutil.Debug(" ")
						return err
					}

				case nameTypeDNS:
					name := string(data)
					if _, ok := domainToReverseLabels(name); !ok {
						logutil.Debug(" ")
						return fmt.Errorf("x509: cannot parse dnsName %q", name)
					}

					if err := checkNameConstraints_Certificate(c, &comparisonCount, maxConstraintComparisons, "DNS name", name, name,
						func(parsedName, constraint any) (bool, error) {
							return matchDomainConstraint_func(parsedName.(string), constraint.(string))
						}, c.PermittedDNSDomains, c.ExcludedDNSDomains); err != nil {
						logutil.Debug(" ")
						return err
					}

				case nameTypeURI:
					name := string(data)
					uri, err := url.Parse(name)
					if err != nil {
						logutil.Debug(" ")
						return fmt.Errorf("x509: internal error: URI SAN %q failed to parse", name)
					}

					if err := checkNameConstraints_Certificate(c, &comparisonCount, maxConstraintComparisons, "URI", name, uri,
						func(parsedName, constraint any) (bool, error) {
							return matchURIConstraint_func(parsedName.(*url.URL), constraint.(string))
						}, c.PermittedURIDomains, c.ExcludedURIDomains); err != nil {
						logutil.Debug(" ")
						return err
					}

				case nameTypeIP:
					ip := net.IP(data)
					if l := len(ip); l != net.IPv4len && l != net.IPv6len {
						return fmt.Errorf("x509: internal error: IP SAN %x failed to parse", data)
					}

					if err := checkNameConstraints_Certificate(c, &comparisonCount, maxConstraintComparisons, "IP address", ip.String(), ip,
						func(parsedName, constraint any) (bool, error) {
							return matchIPConstraint_func(parsedName.(net.IP), constraint.(*net.IPNet))
						}, c.PermittedIPRanges, c.ExcludedIPRanges); err != nil {
						logutil.Debug(" ")
						return err
					}

				default:
					// Unknown SAN types are ignored.
				}

				return nil
			})

			if err != nil {
				return err
			}
		}
	}

	// KeyUsage status flags are ignored. From Engineering Security, Peter
	// Gutmann: A European government CA marked its signing certificates as
	// being valid for encryption only, but no-one noticed. Another
	// European CA marked its signature keys as not being valid for
	// signatures. A different CA marked its own trusted root certificate
	// as being invalid for certificate signing. Another national CA
	// distributed a certificate to be used to encrypt data for the
	// country’s tax authority that was marked as only being usable for
	// digital signatures but not for encryption. Yet another CA reversed
	// the order of the bit flags in the keyUsage due to confusion over
	// encoding endianness, essentially setting a random keyUsage in
	// certificates that it issued. Another CA created a self-invalidating
	// certificate by adding a certificate policy statement stipulating
	// that the certificate had to be used strictly as specified in the
	// keyUsage, and a keyUsage containing a flag indicating that the RSA
	// encryption key could only be used for Diffie-Hellman key agreement.

	if certType == intermediateCertificate && (!c.BasicConstraintsValid || !c.IsCA) {
		logutil.Debug(" ")
		return x509.CertificateInvalidError{c, x509.NotAuthorizedToSign, ""}
	}

	if c.BasicConstraintsValid && c.MaxPathLen >= 0 {
		numIntermediates := len(currentChain) - 1
		if numIntermediates > c.MaxPathLen {
			logutil.Debug(" ")
			return x509.CertificateInvalidError{c, x509.TooManyIntermediates, ""}
		}
	}

	if !boringAllowCert_func(c) {
		// IncompatibleUsage is not quite right here,
		// but it's also the "no chains found" error
		// and is close enough.
		logutil.Debug(" ")
		return x509.CertificateInvalidError{c, x509.IncompatibleUsage, ""}
	}

	return nil
}

// alreadyInChain checks whether a candidate certificate is present in a chain.
// Rather than doing a direct byte for byte equivalency check, we check if the
// subject, public key, and SAN, if present, are equal. This prevents loops that
// are created by mutual cross-signatures, or other cross-signature bridge
// oddities.
func alreadyInChain_func(candidate *x509.Certificate, chain []*x509.Certificate) bool {
	type pubKeyEqual interface {
		Equal(crypto.PublicKey) bool
	}

	var candidateSAN *pkix.Extension
	for _, ext := range candidate.Extensions {
		if ext.Id.Equal(oidExtensionSubjectAltName) {
			candidateSAN = &ext
			break
		}
	}

	for _, cert := range chain {
		if !bytes.Equal(candidate.RawSubject, cert.RawSubject) {
			continue
		}
		if !candidate.PublicKey.(pubKeyEqual).Equal(cert.PublicKey) {
			continue
		}
		var certSAN *pkix.Extension
		for _, ext := range cert.Extensions {
			if ext.Id.Equal(oidExtensionSubjectAltName) {
				certSAN = &ext
				break
			}
		}
		if candidateSAN == nil && certSAN == nil {
			return true
		} else if candidateSAN == nil || certSAN == nil {
			return false
		}
		if bytes.Equal(candidateSAN.Value, certSAN.Value) {
			return true
		}
	}
	return false
}

func appendToFreshChain_func(chain []*x509.Certificate, cert *x509.Certificate) []*x509.Certificate {
	n := make([]*x509.Certificate, len(chain)+1)
	copy(n, chain)
	n[len(chain)] = cert
	return n
}

// UnknownAuthorityError results when the certificate issuer is unknown
type UnknownAuthorityError_s struct {
	Cert *x509.Certificate
	// hintErr contains an error that may be helpful in determining why an
	// authority wasn't found.
	HintErr error
	// hintCert contains a possible authority certificate that was rejected
	// because of the error in hintErr.
	HintCert *x509.Certificate
}

func (e UnknownAuthorityError_s) Error() string {
	s := "x509: certificate signed by unknown authority"
	if e.HintErr != nil {
		certName := e.HintCert.Subject.CommonName
		if len(certName) == 0 {
			if len(e.HintCert.Subject.Organization) > 0 {
				certName = e.HintCert.Subject.Organization[0]
			} else {
				certName = "serial:" + e.HintCert.SerialNumber.String()
			}
		}
		s += fmt.Sprintf(" (possibly because of %q while trying to verify candidate authority certificate %q)", e.HintErr, certName)
	}
	return s
}

// maxChainSignatureChecks is the maximum number of CheckSignatureFrom calls
// that an invocation of buildChains will (transitively) make. Most chains are
// less than 15 certificates long, so this leaves space for multiple chains and
// for failed checks due to different intermediates having the same Subject.
const maxChainSignatureChecks_var = 100

func buildChains_Certificate(c *x509.Certificate, currentChain []*x509.Certificate, sigChecks *int, opts *VerifyOptionsF) (chains [][]*x509.Certificate, err error) {
	var (
		hintErr  error
		hintCert *x509.Certificate
	)

	logutil.Debug(" ")

	considerCandidate := func(certType int, candidate potentialParent) {
		logutil.Debug(" ")
		if candidate.cert.PublicKey == nil || alreadyInChain_func(candidate.cert, currentChain) {
			return
		}

		if sigChecks == nil {
			sigChecks = new(int)
		}
		*sigChecks++
		logutil.Debug(" ")
		if *sigChecks > maxChainSignatureChecks_var {
			err = errors.New("x509: signature check attempts limit reached while verifying certificate chain")
			return
		}

		logutil.Debug(" ")
		if err := c.CheckSignatureFrom(candidate.cert); err != nil {
			logutil.Debug(" ")
			if hintErr == nil {
				hintErr = err
				hintCert = candidate.cert
			}
			logutil.Debug(" ")
			return
		}

		logutil.Debug(" ")
		err = isValid_Certificate(candidate.cert, certType, currentChain, opts)
		if err != nil {
			logutil.Debug(" ")
			if hintErr == nil {
				hintErr = err
				hintCert = candidate.cert
			}
			return
		}

		logutil.Debug(" ")

		if candidate.constraint != nil {
			if err := candidate.constraint(currentChain); err != nil {
				logutil.Debug(" ")
				if hintErr == nil {
					hintErr = err
					hintCert = candidate.cert
				}
				return
			}
		}

		switch certType {
		case rootCertificate:
			chains = append(chains, appendToFreshChain_func(currentChain, candidate.cert))
		case intermediateCertificate:
			var childChains [][]*x509.Certificate
			childChains, err = buildChains_Certificate(candidate.cert, appendToFreshChain_func(currentChain, candidate.cert), sigChecks, opts)
			chains = append(chains, childChains...)
		}
	}

	logutil.Debug(" ")
	for _, root := range opts.Roots.findPotentialParents(c) {
		considerCandidate(rootCertificate, root)
	}

	logutil.Debug(" ")
	for _, intermediate := range opts.Intermediates.findPotentialParents(c) {
		considerCandidate(intermediateCertificate, intermediate)
	}

	if len(chains) > 0 {
		err = nil
	}
	if len(chains) == 0 && err == nil {
		logutil.Debug(" ")
		err = UnknownAuthorityError_s{c, hintErr, hintCert}
	}

	return
}

func checkChainForKeyUsage_func(chain []*x509.Certificate, keyUsages []x509.ExtKeyUsage) bool {
	usages := make([]x509.ExtKeyUsage, len(keyUsages))
	copy(usages, keyUsages)

	if len(chain) == 0 {
		return false
	}

	usagesRemaining := len(usages)

	// We walk down the list and cross out any usages that aren't supported
	// by each certificate. If we cross out all the usages, then the chain
	// is unacceptable.

NextCert:
	for i := len(chain) - 1; i >= 0; i-- {
		cert := chain[i]
		if len(cert.ExtKeyUsage) == 0 && len(cert.UnknownExtKeyUsage) == 0 {
			// The certificate doesn't have any extended key usage specified.
			continue
		}

		for _, usage := range cert.ExtKeyUsage {
			if usage == x509.ExtKeyUsageAny {
				// The certificate is explicitly good for any usage.
				continue NextCert
			}
		}

		const invalidUsage x509.ExtKeyUsage = -1

	NextRequestedUsage:
		for i, requestedUsage := range usages {
			if requestedUsage == invalidUsage {
				continue
			}

			for _, usage := range cert.ExtKeyUsage {
				if requestedUsage == usage {
					continue NextRequestedUsage
				}
			}

			usages[i] = invalidUsage
			usagesRemaining--
			if usagesRemaining == 0 {
				return false
			}
		}
	}

	return true
}

// Verify attempts to verify c by building one or more chains from c to a
// certificate in opts.Roots, using certificates in opts.Intermediates if
// needed. If successful, it returns one or more chains where the first
// element of the chain is c and the last element is from opts.Roots.
//
// If opts.Roots is nil, the platform verifier might be used, and
// verification details might differ from what is described below. If system
// roots are unavailable the returned error will be of type SystemRootsError.
//
// Name constraints in the intermediates will be applied to all names claimed
// in the chain, not just opts.DNSName. Thus it is invalid for a leaf to claim
// example.com if an intermediate doesn't permit it, even if example.com is not
// the name being validated. Note that DirectoryName constraints are not
// supported.
//
// Name constraint validation follows the rules from RFC 5280, with the
// addition that DNS name constraints may use the leading period format
// defined for emails and URIs. When a constraint has a leading period
// it indicates that at least one additional label must be prepended to
// the constrained name to be considered valid.
//
// Extended Key Usage values are enforced nested down a chain, so an intermediate
// or root that enumerates EKUs prevents a leaf from asserting an EKU not in that
// list. (While this is not specified, it is common practice in order to limit
// the types of certificates a CA can issue.)
//
// Certificates that use SHA1WithRSA and ECDSAWithSHA1 signatures are not supported,
// and will not be used to build chains.
//
// Certificates other than c in the returned chains should not be modified.
//
// WARNING: this function doesn't do any revocation checking.
func Verify_Certificate(c *x509.Certificate, opts VerifyOptionsF) (chains [][]*x509.Certificate, err error) {
	// Platform-specific verification needs the ASN.1 contents so
	// this makes the behavior consistent across platforms.
	if len(c.Raw) == 0 {
		return nil, errNotParsed
	}
	for i := 0; i < opts.Intermediates.len(); i++ {
		c, _, err := opts.Intermediates.cert(i)
		if err != nil {
			return nil, fmt.Errorf("crypto/x509: error fetching intermediate: %w", err)
		}
		if len(c.Raw) == 0 {
			return nil, errNotParsed
		}
	}

	if opts.Roots == nil {
		return nil, errNoRootError
	}

	logutil.Debug(" ")

	err = isValid_Certificate(c, leafCertificate, nil, &opts)
	if err != nil {
		return
	}

	logutil.Debug(" ")
	if len(opts.DNSName) > 0 {
		err = c.VerifyHostname(opts.DNSName)
		if err != nil {
			return
		}
	}

	logutil.Debug(" ")
	var candidateChains [][]*x509.Certificate
	if opts.Roots.contains(c) {
		logutil.Debug(" ")
		candidateChains = [][]*x509.Certificate{{c}}
	} else {
		candidateChains, err = buildChains_Certificate(c, []*x509.Certificate{c}, nil, &opts)
		if err != nil {
			return nil, err
		}
	}

	logutil.Debug(" ")
	if len(opts.KeyUsages) == 0 {
		opts.KeyUsages = []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth}
	}

	logutil.Debug(" ")
	for _, eku := range opts.KeyUsages {
		if eku == x509.ExtKeyUsageAny {
			// If any key usage is acceptable, no need to check the chain for
			// key usages.
			return candidateChains, nil
		}
	}

	chains = make([][]*x509.Certificate, 0, len(candidateChains))
	for _, candidate := range candidateChains {
		if checkChainForKeyUsage_func(candidate, opts.KeyUsages) {
			chains = append(chains, candidate)
		}
	}

	if len(chains) == 0 {
		return nil, x509.CertificateInvalidError{c, x509.IncompatibleUsage, ""}
	}

	return chains, nil
}
