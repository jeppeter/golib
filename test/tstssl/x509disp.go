package main

import (
	"crypto/x509"
	"crypto/x509/pkix"
	"fmt"
	"net"
	"net/url"
	"reflect"
)

func format_tabs(tabs int) {
	var i int
	for i = 0; i < tabs; i += 1 {
		fmt.Printf("    ")
	}
	return
}

func display_strings(arrs []string, note string, tabs int) {
	format_tabs(tabs)
	fmt.Printf("%s", note)
	var i int
	for i = 0; i < len(arrs); i += 1 {
		if (i % 4) == 0 {
			fmt.Printf("\n")
			format_tabs(tabs)
			fmt.Printf("0x%04x", i)
		}
		fmt.Printf(" %s", arrs[i])
	}
	fmt.Printf("\n")
}

func display_buffer(buf []byte, note string, tabs int) {
	var i, lasti int
	fmt.Printf("%s buffer [%d:0x%x]", note, len(buf), len(buf))
	lasti = 0
	for i = 0; i < len(buf); i += 1 {
		if (i % 16) == 0 {
			if i > 0 {
				fmt.Printf("    ")
				for lasti < i {
					if buf[lasti] >= byte(' ') && buf[lasti] <= byte('~') {
						fmt.Printf("%c", buf[lasti])
					} else {
						fmt.Printf(".")
					}

					lasti += 1
				}
				fmt.Printf("\n")
			}
			fmt.Printf("0x%08x ", i)
		}
		fmt.Printf(" 0x%02x", buf[i])
	}

	if lasti != i {
		for (i % 16) != 0 {
			fmt.Printf("     ")
			i += 1
		}
		fmt.Printf("    ")
		for lasti < len(buf) {
			if buf[lasti] >= byte(' ') && buf[lasti] <= byte('~') {
				fmt.Printf("%c", buf[lasti])
			} else {
				fmt.Printf(".")
			}

			lasti += 1
		}
		fmt.Printf("\n")
	}
	return
}

func display_names(ns []string, note string, tabs int) {
	if len(ns) > 0 {
		format_tabs(tabs)
		fmt.Printf("%s\n", note)
		format_tabs(tabs + 1)
		for i, c := range ns {
			if i > 0 {
				fmt.Printf(" ")
			}
			fmt.Printf("%s", c)
		}
		fmt.Printf("\n")
	} else {
		format_tabs(tabs)
		fmt.Printf("no %s\n", note)
	}
}

func display_attr(ns []pkix.AttributeTypeAndValue, note string, tabs int) {
	var objid []int
	var types string
	if len(ns) > 0 {
		format_tabs(tabs)
		fmt.Printf("%s\n", note)
		format_tabs(tabs + 1)
		for i, c := range ns {
			if i > 0 {
				fmt.Printf(" ")
			}
			objid = c.Type
			switch c.Value.(type) {
			case string:
				types = "string"
			case float64:
				types = "float64"
			case []interface{}:
				types = "array"
			case map[string]interface{}:
				types = "map"
			case bool:
				types = "bool"
			default:
				types = fmt.Sprintf("%s", reflect.TypeOf(c.Value))
			}
			fmt.Printf("type [%s] %v value [%s] type[%s]", c.Type, objid, c.Value, types)
		}
		fmt.Printf("\n")
	} else {
		format_tabs(tabs)
		fmt.Printf("no attr %s\n", note)
	}
}

func display_keyusage(u *x509.KeyUsage, note string, tabs int) {
	format_tabs(tabs)
	fmt.Printf("%s %v\n", note, *u)
}

func display_extkeyusage(u []x509.ExtKeyUsage, note string, tabs int) {
	if len(u) > 0 {
		format_tabs(tabs)
		fmt.Printf("%s\n", note)
		format_tabs(tabs + 1)
		for i, c := range u {
			if i > 0 {
				fmt.Printf(" ")
			}
			fmt.Printf("%v", c)
		}
		fmt.Printf("\n")
	} else {
		format_tabs(tabs)
		fmt.Printf("no %s ExtKeyUsage\n", note)
	}
}

func debug_pkix_name(n *pkix.Name, note string, tabs int) {
	format_tabs(tabs)
	fmt.Printf("pkix.Name %s:\n", note)
	display_names(n.Country, "Country", tabs+1)
	display_names(n.Organization, "Organization", tabs+1)
	display_names(n.OrganizationalUnit, "OrganizationalUnit", tabs+1)
	display_names(n.Locality, "Locality", tabs+1)
	display_names(n.Province, "Province", tabs+1)
	display_names(n.StreetAddress, "StreetAddress", tabs+1)
	display_names(n.PostalCode, "PostalCode", tabs+1)
	format_tabs(tabs + 1)
	fmt.Printf("SerialNumber [%s]\n", n.SerialNumber)
	format_tabs(tabs + 1)
	fmt.Printf("CommonName [%s]\n", n.CommonName)
	display_attr(n.Names, "Names", tabs+1)
	display_attr(n.ExtraNames, "ExtraNames", tabs+1)
	return
}

func display_url(u *url.URL, note string, tabs int) {
	format_tabs(tabs)
	fmt.Printf("%s %s\n", note, u.String())
}

func display_urls(us []*url.URL, note string, tabs int) {
	for i := 0; i < len(us); i += 1 {
		display_url(us[i], fmt.Sprintf("%s[%d]", note, i), tabs+1)
	}
}

func display_ips(ips []net.IP, note string, tabs int) {
	if len(ips) > 0 {
		format_tabs(tabs)
		fmt.Printf("%s\n", note)
		format_tabs(tabs + 1)
		for i, c := range ips {
			if i > 0 {
				fmt.Printf(" ")
			}
			fmt.Printf("[%v]", c)
		}
		fmt.Printf("\n")
	} else {
		format_tabs(tabs)
		fmt.Printf("no %s\n", note)
	}
}

func display_ipnets(ips []*net.IPNet, note string, tabs int) {
	if len(ips) > 0 {
		format_tabs(tabs)
		fmt.Printf("%s\n", note)
		format_tabs(tabs + 1)
		for i, c := range ips {
			if i > 0 {
				fmt.Printf(" ")
			}
			fmt.Printf("[IP %v][Mask %v]", c.IP, c.Mask)
		}
		fmt.Printf("\n")
	} else {
		format_tabs(tabs)
		fmt.Printf("no %s\n", note)
	}
}

func display_x509_cert(cert *x509.Certificate) {
	var pkname *pkix.Name
	fmt.Printf("version %d\n", cert.Version)
	fmt.Printf("SerialNumber %v\n", cert.SerialNumber)
	fmt.Printf("NotBefore %v NotAfter %v\n", cert.NotBefore, cert.NotAfter)
	fmt.Printf("PublicKeyAlgorithm %v\n", cert.PublicKeyAlgorithm)
	fmt.Printf("PublicKey         %v\n", cert.PublicKey)
	fmt.Printf("IsCA %v\n", cert.IsCA)
	pkname = &cert.Issuer
	debug_pkix_name(pkname, "Issuer", 0)
	pkname = &cert.Subject
	debug_pkix_name(pkname, "Subject", 0)
	display_keyusage(&cert.KeyUsage, "KeyUsage", 0)
	display_extkeyusage(cert.ExtKeyUsage, "ExtKeyUsage", 0)
	fmt.Printf("AuthorityKeyId %v\n", cert.AuthorityKeyId)
	fmt.Printf("BasicConstraintsValid %v\n", cert.BasicConstraintsValid)
	fmt.Printf("CRLDistributionPoints %v\n", cert.CRLDistributionPoints)
	fmt.Printf("DNSNames %v\n", cert.DNSNames)
	fmt.Printf("EmailAddresses %v\n", cert.EmailAddresses)
	fmt.Printf("ExcludedDNSDomains %v\n", cert.ExcludedDNSDomains)
	fmt.Printf("ExcludedEmailAddresses %v\n", cert.ExcludedEmailAddresses)
	fmt.Printf("ExcludedIPRanges %v\n", cert.ExcludedIPRanges)
	fmt.Printf("ExcludedURIDomains %v\n", cert.ExcludedURIDomains)
	display_ips(cert.IPAddresses, "IPAddresses", 0)
	fmt.Printf("IssuingCertificateURL %v\n", cert.IssuingCertificateURL)
	fmt.Printf("MaxPathLen %d\n", cert.MaxPathLen)
	fmt.Printf("MaxPathLenZero %v\n", cert.MaxPathLenZero)
	fmt.Printf("OCSPServer %v\n", cert.OCSPServer)
	fmt.Printf("PermittedDNSDomains %v\n", cert.PermittedDNSDomains)
	fmt.Printf("PermittedDNSDomainsCritical  %v\n", cert.PermittedDNSDomainsCritical)
	fmt.Printf("PermittedEmailAddresses %v\n", cert.PermittedEmailAddresses)
	display_ipnets(cert.PermittedIPRanges, "PermittedIPRanges", 0)
	fmt.Printf("PermittedURIDomains %v\n", cert.PermittedURIDomains)
	fmt.Printf("PolicyIdentifiers %v\n", cert.PolicyIdentifiers)
	fmt.Printf("Policies %v\n", cert.Policies)
	fmt.Printf("SignatureAlgorithm %v\n", cert.SignatureAlgorithm)
	fmt.Printf("SubjectKeyId %v\n", cert.SubjectKeyId)
	fmt.Printf("URIs %v\n", cert.URIs)
	fmt.Printf("UnknownExtKeyUsage %v\n", cert.UnknownExtKeyUsage)
	return
}

func display_pkix_attrvalue(val *pkix.AttributeTypeAndValue, tabs int, note string) {
	format_tabs(tabs)
	fmt.Printf("%s type %s value %v", note, val.Type.String(), val.Value)
}

func display_attri_set(set *pkix.AttributeTypeAndValueSET, note string, tabs int) {
	format_tabs(tabs)
	fmt.Printf("%s ObjectIdentifier %s", note, set.Type.String())
	var idx int
	var ns string

	for idx = 0; idx < len(set.Value); idx += 1 {
		ns = fmt.Sprintf("%s[%d]", note, idx)
		display_attr(set.Value[idx], ns, tabs+1)
	}
	return
}

func display_attri_array(set []pkix.AttributeTypeAndValueSET, note string, tabs int) {
	var ns string
	var cval *pkix.AttributeTypeAndValueSET
	var i int
	for i = 0; i < len(set); i += 1 {
		cval = &set[i]
		ns = fmt.Sprintf("%s[%d]", note, i)
		display_attri_set(cval, ns, tabs+1)
	}
}

func display_extension(curset *pkix.Extension, note string, tabs int) {
	display_buffer(curset.Value, fmt.Sprintf("%s.Id %s .Critical %v", note, curset.Id.String(), curset.Critical), tabs)
}

func display_extensions(ext []pkix.Extension, note string, tabs int) {
	var ns string
	var i int
	var curv *pkix.Extension
	for i = 0; i < len(ext); i += 1 {
		curv = &ext[i]
		ns = fmt.Sprintf("%s[%d]", note, i)
		display_extension(curv, ns, tabs+1)
	}
}
