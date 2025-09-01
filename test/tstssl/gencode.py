#! /usr/bin/env python

import extargsparse
import logging
import os
import sys
import re
import importlib
import struct


def read_file(infile=None):
    fin = sys.stdin
    if infile is not None:
        fin = open(infile,'rb')
    rets = ''
    if 'b' in fin.mode:
        rdata = b''
        while True:
            try:
                l = fin.read(64 * 1024)
                if l is None or len(l) == 0:
                    break
                rdata += l
            except:
                break
        if sys.version[0] == '3':
            rets = rdata.decode('utf-8')
        else:
            rets = rdata
    else:        
        for l in fin:
            s = l
            rets += s

    if fin != sys.stdin:
        fin.close()
    fin = None
    return rets

def read_file_bytes(infile=None):
    fin = sys.stdin
    if infile is not None:
        fin = open(infile,'rb')
    retb = b''
    while True:
        if fin != sys.stdin:
            curb = fin.read(1024 * 1024)
        else:
            curb = fin.buffer.read()
        if curb is None or len(curb) == 0:
            break
        retb += curb
    if fin != sys.stdin:
        fin.close()
    fin = None
    return retb


def write_file(s,outfile=None):
    fout = sys.stdout
    if outfile is not None:
        fout = open(outfile, 'w+b')
    outs = s
    if 'b' in fout.mode:
        outs = s.encode('utf-8')
    fout.write(outs)
    if fout != sys.stdout:
        fout.close()
    fout = None
    return 

def write_file_bytes(sarr,outfile=None):
    fout = sys.stdout
    if outfile is not None:
        fout = open(outfile, 'wb')
    if 'b' not in fout.mode:
        fout.buffer.write(sarr)
    else:        
        fout.write(sarr)
    if fout != sys.stdout:
        fout.close()
    fout = None
    return 


def set_logging(args):
    loglvl= logging.ERROR
    if args.verbose >= 3:
        loglvl = logging.DEBUG
    elif args.verbose >= 2:
        loglvl = logging.INFO
    curlog = logging.getLogger(args.lognames)
    #sys.stderr.write('curlog [%s][%s]\n'%(args.logname,curlog))
    curlog.setLevel(loglvl)
    if len(curlog.handlers) > 0 :
        curlog.handlers = []
    formatter = logging.Formatter('%(asctime)s:%(filename)s:%(funcName)s:%(lineno)d<%(levelname)s>\t%(message)s')
    if not args.lognostderr:
        logstderr = logging.StreamHandler()
        logstderr.setLevel(loglvl)
        logstderr.setFormatter(formatter)
        curlog.addHandler(logstderr)

    for f in args.logfiles:
        flog = logging.FileHandler(f,mode='w',delay=False)
        flog.setLevel(loglvl)
        flog.setFormatter(formatter)
        curlog.addHandler(flog)
    for f in args.logappends:       
        if args.logrotate:
            flog = logging.handlers.RotatingFileHandler(f,mode='a',maxBytes=args.logmaxbytes,backupCount=args.logbackupcnt,delay=0)
        else:
            sys.stdout.write('appends [%s] file\n'%(f))
            flog = logging.FileHandler(f,mode='a',delay=0)
        flog.setLevel(loglvl)
        flog.setFormatter(formatter)
        curlog.addHandler(flog)
    return


def load_log_commandline(parser):
    logcommand = '''
    {
        "verbose|v" : "+",
        "logname" : "root",
        "logfiles" : [],
        "logappends" : [],
        "logrotate" : true,
        "logmaxbytes" : 10000000,
        "logbackupcnt" : 2,
        "lognostderr" : false
    }
    '''
    parser.load_command_line_string(logcommand)
    return parser

def write_file(s,outfile=None):
    fout = sys.stdout
    if outfile is not None:
        fout = open(outfile, 'w+b')
    outs = s
    if 'b' in fout.mode:
        outs = s.encode('utf-8')
    fout.write(outs)
    if fout != sys.stdout:
        fout.close()
    fout = None
    return 


def split_capital_name(name):
	outarr = []
	curb = b''
	inb = None
	if sys.version[0] == '3':
		inb = name.encode('utf-8')
	else:
		inb = byte(name)
	idx = 0
	aintval = 0x41
	zintval = 0x5a
	while idx < len(inb):
		if inb[idx] >= aintval and inb[idx] <= zintval:
			if len(curb) > 0:
				if sys.version[0] == '3':
					curs = curb.decode('utf-8')
				else:
					curs = string(curb)
				outarr.append(curs)
			curb = b''
		curb += struct.pack('B',inb[idx])
		idx += 1
	if len(curb) > 0:
		if sys.version[0] == '3':
			curs = curb.decode('utf-8')
		else:
			curs = string(curb)
		outarr.append(curs)
	return outarr

errkeys = {
	'KEYWORD_C_R_L_DISTRIBUTION_POINTS' : 'KEYWORD_CRL_DISTRIBUTION_POINTS',
	'KEYWORD_D_N_S_NAMES' : 'KEYWORD_DNS_NAMES',
	'KEYWORD_EXCLUDED_D_N_S_DOMAINS' : 'KEYWORD_EXCLUDED_DNS_DOMAINS',
	'KEYWORD_EXCLUDED_I_P_RANGES' : 'KEYWORD_EXCLUDED_IP_RANGES',
	'KEYWORD_I_P_ADDRESSES' : 'KEYWORD_IP_ADDRESSES',
	'KEYWORD_IS_C_A' : 'KEYWORD_IS_CA',
	'KEYWORD_ISSUING_CERTIFICATE_U_R_L' : 'KEYWORD_ISSUING_CERTIFICATE_URL',
	'KEYWORD_PERMITTED_D_N_S_DOMAINS' : 'KEYWORD_PERMITTED_DNS_DOMAINS',
	'KEYWORD_PERMITTED_D_N_S_DOMAINS_CRITICAL' : 'KEYWORD_PERMITTED_DNS_DOMAINS_CRITICAL',
	'KEYWORD_PERMITTED_I_P_RANGES' : 'KEYWORD_PERMITTED_IP_RANGES',
	'KEYWORD_PERMITTED_U_R_I_DOMAINS' : 'KEYWORD_PERMITTED_URI_DOMAINS',
	'KEYWORD_U_R_IS' : 'KEYWORD_URIS',
	'KEY_USAGE_C_R_L_SIGN' : 'KEY_USAGE_CRL_SIGN',
	'KEYWORD_EXCLUDED_U_R_I_DOMAINS' : 'KEYWORD_EXCLUDED_URI_DOMAINS'
}


def format_keyword(name):
	outarr = split_capital_name(name)
	keyword = 'KEYWORD'
	for s in outarr:
		keyword += '_%s'%(s.upper())

	if keyword in errkeys.keys():
		keyword = errkeys[keyword]
	return keyword

def format_uppername(name):
	outarr = split_capital_name(name)
	keyword = ''
	for s in outarr:
		if len(keyword) > 0:
			keyword += '_'
		keyword += '%s'%(s.upper())

	if keyword in errkeys.keys():
		keyword = errkeys[keyword]
	return keyword



def format_tabline(tab,l):
	s = ''
	for i in range(tab):
		s += '    '
	s += '%s\n'%(l)
	return s


def format_strings_code(name,varname='x509temp'):
	keyword = format_keyword(name)
	runcode = format_tabline(1,'%s.%s = []string{}'%(varname,name))
	runcode += format_tabline(1,'valarr, ok = mapv[%s].([]interface{})'%(keyword))
	runcode += format_tabline(1,'if ok {')
	runcode += format_tabline(2,'logutil.Debug("[%%s] parse", %s)'%(keyword))
	runcode += format_tabline(2,'arrs, err = trans_inter_to_string(valarr,%s)'%(keyword))
	runcode += format_tabline(2, 'if err != nil {')
	runcode += format_tabline(3,'return')
	runcode += format_tabline(2,'}')
	runcode += format_tabline(2,'%s.%s = arrs'%(varname,name))
	runcode += format_tabline(1,'}')
	kdefine = '%s = "%s"'%(keyword,name.lower())
	return kdefine,runcode

def format_int_code(name,varname='x509temp'):
	keyword = format_keyword(name)
	runcode = format_tabline(1,'%s.%s = 0'%(varname,name))
	runcode += format_tabline(1,'intval, ok = mapv[%s]'%(keyword))
	runcode += format_tabline(1,'if ok {')
	runcode += format_tabline(2,'logutil.Debug("[%%s] parse", %s)'%(keyword))
	runcode += format_tabline(2,'%s.%s, err = get_int_value(intval,%s)'%(varname,name,keyword))
	runcode += format_tabline(2,'if err != nil {')
	runcode += format_tabline(3,'return')
	runcode += format_tabline(2,'}')
	runcode += format_tabline(1,'}')
	kdefine = '%s = "%s"'%(keyword,name.lower())
	return kdefine,runcode

def format_keyusage_code(name,varname='x509temp'):
	keyword = format_keyword(name)
	runcode = format_tabline(1,'%s.%s = 0'%(varname,name))
	runcode += format_tabline(1,'valarr, ok = mapv[%s].([]interface{})'%(keyword))
	runcode += format_tabline(1,'if ok {')
	runcode += format_tabline(2,'logutil.Debug("[%%s] parse", %s)'%(keyword))
	runcode += format_tabline(2,'arrs, err = trans_inter_to_string(valarr,%s)'%(keyword))
	runcode += format_tabline(2, 'if err != nil {')
	runcode += format_tabline(3,'return')
	runcode += format_tabline(2,'}')
	runcode += format_tabline(2,'%s.%s, err = get_keyusage_value(arrs,%s)'%(varname,name,keyword))
	runcode += format_tabline(2,'if err != nil {')
	runcode += format_tabline(3,'return')
	runcode += format_tabline(2,'}')
	runcode += format_tabline(1,'}')
	kdefine = '%s = "%s"'%(keyword,name.lower())
	return kdefine,runcode

def format_bytes_code(name,varname='x509temp'):
	keyword = format_keyword(name)
	runcode = format_tabline(1,'%s.%s = []byte{}'%(varname,name))
	runcode += format_tabline(1,'valarr, ok = mapv[%s].([]interface{})'%(keyword))
	runcode += format_tabline(1,'if ok {')
	runcode += format_tabline(2,'logutil.Debug("[%%s] parse", %s)'%(keyword))
	runcode += format_tabline(2,'%s.%s, err = get_bytes_value(valarr,%s)'%(varname,name,keyword))
	runcode += format_tabline(2,'if err != nil {')
	runcode += format_tabline(3,'return')
	runcode += format_tabline(2,'}')
	runcode += format_tabline(1,'}')
	kdefine = '%s = "%s"'%(keyword,name.lower())
	return kdefine,runcode

def format_bool_code(name,varname='x509temp'):
	keyword = format_keyword(name)
	runcode = format_tabline(1,'%s.%s = false'%(varname,name))
	runcode += format_tabline(1,'valb, ok = mapv[%s].(bool)'%(keyword))
	runcode += format_tabline(1,'if ok {')
	runcode += format_tabline(2,'logutil.Debug("[%%s] parse",%s)'%(keyword))
	runcode += format_tabline(2,'%s.%s = valb'%(varname,name))
	runcode += format_tabline(1,'}')
	kdefine = '%s = "%s"'%(keyword,name.lower())
	return kdefine,runcode

def format_ipnet_code(name,varname='x509temp'):
	keyword = format_keyword(name)
	runcode = format_tabline(1,'%s.%s = []*net.IPNet{}'%(varname,name))
	runcode += format_tabline(1,'valarr, ok = mapv[%s].([]interface{})'%(keyword))
	runcode += format_tabline(1,'if ok {')
	runcode += format_tabline(2,'logutil.Debug("[%%s] parse",%s)'%(keyword))
	runcode += format_tabline(2,'%s.%s , err = get_netip_value(valarr, %s)'%(varname,name,keyword))
	runcode += format_tabline(2,'if err != nil {')
	runcode += format_tabline(3,'return')
	runcode += format_tabline(2,'}')
	runcode += format_tabline(1,'}')
	kdefine = '%s = "%s"'%(keyword,name.lower())
	return kdefine,runcode

def format_ip_code(name,varname='x509temp'):
	keyword = format_keyword(name)
	runcode = format_tabline(1,'%s.%s = []net.IP{}'%(varname,name))
	runcode += format_tabline(1,'valarr, ok = mapv[%s].([]interface{})'%(keyword))
	runcode += format_tabline(1,'if ok {')
	runcode += format_tabline(2,'logutil.Debug("[%%s] parse",%s)'%(keyword))
	runcode += format_tabline(2,'%s.%s , err = get_ip_value(valarr, %s)'%(varname,name,keyword))
	runcode += format_tabline(2,'if err != nil {')
	runcode += format_tabline(3,'return')
	runcode += format_tabline(2,'}')
	runcode += format_tabline(1,'}')
	kdefine = '%s = "%s"'%(keyword,name.lower())
	return kdefine,runcode

def format_time_code(name,varname='x509temp'):
	keyword = format_keyword(name)
	runcode = format_tabline(1,'%s.%s = time.Now()'%(varname,name))
	runcode += format_tabline(1,'vals, ok = mapv[%s].(string)'%(keyword))
	runcode += format_tabline(1,'if ok {')
	runcode += format_tabline(2,'logutil.Debug("[%%s] parse",%s)'%(keyword))
	runcode += format_tabline(2,'%s.%s , err = get_time_value(vals, %s)'%(varname,name,keyword))
	runcode += format_tabline(2,'if err != nil {')
	runcode += format_tabline(3,'return')
	runcode += format_tabline(2,'}')
	runcode += format_tabline(1,'}')
	kdefine = '%s = "%s"'%(keyword,name.lower())
	return kdefine,runcode

def format_timeafter_code(name,varname='x509temp'):
	keyword = format_keyword(name)
	runcode = format_tabline(1,'%s.%s = time.Now().AddDate(20,0,0)'%(varname,name))
	runcode += format_tabline(1,'vals, ok = mapv[%s].(string)'%(keyword))
	runcode += format_tabline(1,'if ok {')
	runcode += format_tabline(2,'logutil.Debug("[%%s] parse",%s)'%(keyword))
	runcode += format_tabline(2,'%s.%s , err = get_time_value(vals, %s)'%(varname,name,keyword))
	runcode += format_tabline(2,'if err != nil {')
	runcode += format_tabline(3,'return')
	runcode += format_tabline(2,'}')
	runcode += format_tabline(1,'}')
	kdefine = '%s = "%s"'%(keyword,name.lower())
	return kdefine,runcode


def format_bigint_code(name,varname='x509temp'):
	keyword = format_keyword(name)
	runcode = format_tabline(1,'%s.%s = big.NewInt(0)'%(varname,name))
	runcode += format_tabline(1,'vals, ok = mapv[%s].(string)'%(keyword))
	runcode += format_tabline(1,'if ok {')
	runcode += format_tabline(2,'%s.%s , err = get_bigint_value(vals,%s)'%(varname,name,keyword))
	runcode += format_tabline(2,'if err != nil {')
	runcode += format_tabline(3,'return')
	runcode += format_tabline(2,'}')
	runcode += format_tabline(1,'}')
	kdefine = '%s = "%s"'%(keyword,name.lower())
	return kdefine,runcode

def format_pkixname_code(name,varname='x509temp'):
	keyword = format_keyword(name)
	runcode = format_tabline(1,'%s.%s = pkix.Name{}'%(varname,name))
	runcode += format_tabline(1,'valmap, ok = mapv[%s].(map[string]interface{})'%(keyword))
	runcode += format_tabline(1,'if ok {')
	runcode += format_tabline(2,'logutil.Debug("[%%s] parse",%s)'%(keyword))
	runcode += format_tabline(2,'%s.%s , err = get_pkixname_value(valmap,%s)'%(varname,name,keyword))
	runcode += format_tabline(2,'if err != nil {')
	runcode += format_tabline(3,'return')
	runcode += format_tabline(2,'}')
	runcode += format_tabline(1,'}')
	kdefine = '%s = "%s"'%(keyword,name.lower())
	return kdefine,runcode

def format_objoids_code(name,varname='x509temp'):
	keyword = format_keyword(name)
	runcode = format_tabline(1,'%s.%s = []asn1.ObjectIdentifier{}'%(varname,name))
	runcode += format_tabline(1,'valarr, ok = mapv[%s].([]interface{})'%(keyword))
	runcode += format_tabline(1,'if ok {')
	runcode += format_tabline(2,'logutil.Debug("[%%s] parse",%s)'%(keyword))
	runcode += format_tabline(2,'%s.%s , err = get_objoids_value(valarr,%s)'%(varname,name,keyword))
	runcode += format_tabline(2,'if err != nil {')
	runcode += format_tabline(3,'return')
	runcode += format_tabline(2,'}')
	runcode += format_tabline(1,'}')
	kdefine = '%s = "%s"'%(keyword,name.lower())
	return kdefine,runcode

def format_oids_code(name,varname='x509temp'):
	keyword = format_keyword(name)
	runcode = format_tabline(1,'%s.%s = []x509.OID{}'%(varname,name))
	runcode += format_tabline(1,'valarr, ok = mapv[%s].([]interface{})'%(keyword))
	runcode += format_tabline(1,'if ok {')
	runcode += format_tabline(2,'logutil.Debug("[%%s] parse",%s)'%(keyword))
	runcode += format_tabline(2,'arrs, err = trans_inter_to_string(valarr,%s)'%(keyword))
	runcode += format_tabline(2, 'if err != nil {')
	runcode += format_tabline(3,'return')
	runcode += format_tabline(2,'}')
	runcode += format_tabline(2,'%s.%s , err = get_oids_value(arrs,%s)'%(varname,name,keyword))
	runcode += format_tabline(2,'if err != nil {')
	runcode += format_tabline(3,'return')
	runcode += format_tabline(2,'}')
	runcode += format_tabline(1,'}')
	kdefine = '%s = "%s"'%(keyword,name.lower())
	return kdefine,runcode

def format_urls_code(name,varname='x509temp'):
	keyword = format_keyword(name)
	runcode = format_tabline(1,'%s.%s = []*url.URL{}'%(varname,name))
	runcode += format_tabline(1,'valarr, ok = mapv[%s].([]interface{})'%(keyword))
	runcode += format_tabline(1,'if ok {')
	runcode += format_tabline(2,'logutil.Debug("[%%s] parse",%s)'%(keyword))
	runcode += format_tabline(2,'arrs, err = trans_inter_to_string(valarr,%s)'%(keyword))
	runcode += format_tabline(2, 'if err != nil {')
	runcode += format_tabline(3,'return')
	runcode += format_tabline(2,'}')
	runcode += format_tabline(2,'%s.%s , err = get_urls_value(arrs,%s)'%(varname,name,keyword))
	runcode += format_tabline(2,'if err != nil {')
	runcode += format_tabline(3,'return')
	runcode += format_tabline(2,'}')
	runcode += format_tabline(1,'}')
	kdefine = '%s = "%s"'%(keyword,name.lower())
	return kdefine,runcode

def format_extkeyusage_code(name,varname='x509temp'):
	keyword = format_keyword(name)
	runcode = format_tabline(1,'%s.%s = []x509.ExtKeyUsage{}'%(varname,name))
	runcode += format_tabline(1,'valarr, ok = mapv[%s].([]interface{})'%(keyword))
	runcode += format_tabline(1,'if ok {')
	runcode += format_tabline(2,'logutil.Debug("[%%s] parse",%s)'%(keyword))
	runcode += format_tabline(2,'arrs, err = trans_inter_to_string(valarr,%s)'%(keyword))
	runcode += format_tabline(2, 'if err != nil {')
	runcode += format_tabline(3,'return')
	runcode += format_tabline(2,'}')
	runcode += format_tabline(2,'%s.%s , err = get_key_ext_usage(arrs,%s)'%(varname,name,keyword))
	runcode += format_tabline(2,'if err != nil {')
	runcode += format_tabline(3,'return')
	runcode += format_tabline(2,'}')
	runcode += format_tabline(1,'}')
	kdefine = '%s = "%s"'%(keyword,name.lower())
	return kdefine,runcode

def format_algorithm_code(name,varname='x509temp'):
	keyword = format_keyword(name)
	runcode = format_tabline(1,'%s.%s = x509.SHA256WithRSA'%(varname,name))
	runcode += format_tabline(1,'vals, ok = mapv[%s].(string)'%(keyword))
	runcode += format_tabline(1,'if ok {')
	runcode += format_tabline(2,'logutil.Debug("[%%s] parse",%s)'%(keyword))
	runcode += format_tabline(2,'%s.%s , err = get_algorithm_value(vals,%s)'%(varname,name,keyword))
	runcode += format_tabline(2,'if err != nil {')
	runcode += format_tabline(3,'return')
	runcode += format_tabline(2,'}')
	runcode += format_tabline(1,'}')
	kdefine = '%s = "%s"'%(keyword,name.lower())
	return kdefine,runcode


def format_attribute_set_code(name,varname='x509temp'):
	keyword = format_keyword(name)
	runcode = format_tabline(1,'%s.%s = []pkix.AttributeTypeAndValueSET{}'%(varname,name))
	runcode += format_tabline(1,'valarr, ok = mapv[%s].([]interface{})'%(keyword))
	runcode += format_tabline(1,'if ok {')
	runcode += format_tabline(1,'logutil.Debug("[%%s] parse",%s)'%(keyword))
	runcode += format_tabline(1,'%s.%s, err = get_attribute_set_value(valarr,%s)'%(varname,name,keyword))
	runcode += format_tabline(1,'if err != nil {')
	runcode += format_tabline(2,'return')
	runcode += format_tabline(1,'}')
	runcode += format_tabline(1,'}')
	kdefine = '%s = "%s"'%(keyword,name.lower())
	return kdefine,runcode


def format_pkix_extensions_code(name,varname='x509temp'):
	keyword = format_keyword(name)
	runcode = format_tabline(1,'%s.%s = []pkix.Extension{}'%(varname,name))
	runcode += format_tabline(1,'valarr, ok = mapv[%s].([]interface{})'%(keyword))
	runcode += format_tabline(1,'if ok {')
	runcode += format_tabline(1,'logutil.Debug("[%%s] parse",%s)'%(keyword))
	runcode += format_tabline(1,'%s.%s, err = get_pkix_extensions_value(valarr,%s)'%(varname,name,keyword))
	runcode += format_tabline(1,'if err != nil {')
	runcode += format_tabline(2,'return')
	runcode += format_tabline(1,'}')
	runcode += format_tabline(1,'}')
	kdefine = '%s = "%s"'%(keyword,name.lower())
	return kdefine,runcode


KEYWORDS = ['ExtKeyUsage',
'AuthorityKeyId',
'BasicConstraintsValid',
'CRLDistributionPoints',
'DNSNames',
'EmailAddresses',
'ExcludedDNSDomains',
'ExcludedEmailAddresses',
'ExcludedIPRanges',
'ExcludedURIDomains',
'IPAddresses',
'IsCA',
'IssuingCertificateURL',
'KeyUsage',
'MaxPathLen',
'MaxPathLenZero',
'NotAfter',
'NotBefore',
'OCSPServer',
'PermittedDNSDomains',
'PermittedDNSDomainsCritical',
'PermittedEmailAddresses',
'PermittedIPRanges',
'PermittedURIDomains',
'PolicyIdentifiers',
'Policies',
'SerialNumber',
'SignatureAlgorithm',
'Subject',
'SubjectKeyId',
'URIs',
'UnknownExtKeyUsage']

KEYMAPS = {
'ExtKeyUsage': 'extkeyusage',
'AuthorityKeyId': 'bytes',
'BasicConstraintsValid':'bool' ,
'CRLDistributionPoints': 'strings',
'DNSNames': 'strings',
'EmailAddresses': 'strings',
'ExcludedDNSDomains': 'strings',
'ExcludedEmailAddresses': 'strings',
'ExcludedIPRanges': 'ipnet',
'ExcludedURIDomains' : 'strings',
'IPAddresses': 'ip',
'IsCA': 'bool',
'IssuingCertificateURL': 'strings',
'KeyUsage': 'keyusage',
'MaxPathLen': 'int',
'MaxPathLenZero': 'bool',
'NotAfter': 'timeafter',
'NotBefore': 'time',
'OCSPServer': 'strings',
'PermittedDNSDomains': 'strings',
'PermittedDNSDomainsCritical': 'bool',
'PermittedEmailAddresses': 'strings',
'PermittedIPRanges': 'ipnet',
'PermittedURIDomains': 'strings',
'PolicyIdentifiers': 'objoids',
'Policies': 'oids',
'SerialNumber': 'bigint',
'SignatureAlgorithm': 'algorithm',
'Subject': 'pkixname',
'SubjectKeyId': 'bytes',
'URIs': 'urls',
'UnknownExtKeyUsage': 'objoids' 
}

def gencode_handler(args,parser):
	set_logging(args)
	defines = []
	outcodes = ''
	logging.info('ccc')

	for k in KEYWORDS:
		types = KEYMAPS[k]
		logging.info('[%s]=[%s]'%(k,types))
		funcname = 'format_%s_code'%(types)
		m = importlib.import_module(__name__)
		val = getattr(m,funcname)
		kd,code = val(k)
		defines.append(kd)
		outcodes += format_tabline(1,'')
		outcodes += code
	write_file(outcodes,args.output)
	defs = 'const (\n'
	for d in defines:
		defs += format_tabline(1,'%s'%(d))
	defs += ')\n'
	sys.stdout.write('%s'%(defs))
	sys.exit(0)
	return

def grconst_handler(args,parser):
	set_logging(args)
	s = read_file(args.input)
	sarr = re.split('\n', s)
	ms = re.compile('^\\s*([a-zA-Z0-9]+)')
	shiftval = 0
	outs = ''
	for l in sarr:
		l = l.rstrip('\r')
		m = ms.findall(l)
		if m is not None and len(m) > 0 :
			name = m[0]
			keyword = format_uppername(name)
			val = 1 << shiftval
			outs += format_tabline(0,'pub const %s :u32 = %d;'%(keyword,val))
			shiftval += 1
	write_file(outs,args.output)
	sys.exit(0)
	return

def genenum_handler(args,parser):
	set_logging(args)
	enumname = args.subnargs[0]
	s = read_file(args.input)
	mexpr = re.compile('^\\s*([a-zA-Z0-9]+)')
	sarr = re.split('\n',s)
	members = []
	for l in sarr:
		l = l.rstrip('\r')
		m = mexpr.findall(l)
		if m is not None and len(m) > 0:
			members.append(m[0])


	outs = format_tabline(0,'pub enum %s {'%(enumname))
	for m in members:
		outs += format_tabline(1,'%s,'%(m))
	outs += format_tabline(0,'}')

	outs += format_tabline(0,'')
	outs += format_tabline(0,'impl PartialEq for %s {'%(enumname))
	outs += format_tabline(1,'fn ne(&self,other :&Self) -> bool {')
	outs += format_tabline(2,'return !self.eq(other);')
	outs += format_tabline(1,'}')

	outs += format_tabline(1,'')

	outs += format_tabline(1,'fn eq(&self,other :&Self) -> bool{')
	outs += format_tabline(2,'let mut retval :bool = false;');
	outs += format_tabline(2,'match self {')
	for m in members:
		outs += format_tabline(3,'%s::%s => {'%(enumname,m))
		outs += format_tabline(4,'match other {')
		outs += format_tabline(5,'%s::%s => {retval = true;},'%(enumname,m))
		outs += format_tabline(5,'_ => {},')
		outs += format_tabline(4,'}')
		outs += format_tabline(3,'},')
	outs += format_tabline(2,'}')
	outs += format_tabline(2,'return retval;')
	outs += format_tabline(1,'}')

	outs += format_tabline(0,'}')

	write_file(outs,args.output)
	sys.exit(0)
	return

def genenumserde_handler(args,parser):
	set_logging(args)
	enumname = args.subnargs[0]
	s = read_file(args.input)
	mexpr = re.compile('^\\s*([a-zA-Z0-9]+)')
	sarr = re.split('\n',s)
	members = []
	for l in sarr:
		l = l.rstrip('\r')
		m = mexpr.findall(l)
		if m is not None and len(m) > 0:
			members.append(m[0])


	outs = format_tabline(0,'impl serde::ser::Serialize for %s{'%(enumname))
	outs += format_tabline(1,'fn serialize<S>(&self,serializer: S) -> Result<S::Ok, S::Error> where S: serde::ser::Serializer {')
	outs += format_tabline(2,'match *self{')
	for m in members:
		outs += format_tabline(3,'%s::%s => {'%(enumname,m))
		outs += format_tabline(4,'return serializer.serialize_str("%s");'%(m.lower()))
		outs += format_tabline(3,'},')
	outs += format_tabline(2,'}')
	outs += format_tabline(1,'}')
	outs += format_tabline(0,'}')

	outs += format_tabline(0,'')
	outs += format_tabline(0,'impl<\'de: \'a, \'a> serde::de::Deserialize<\'de> for %s {'%(enumname))
	outs += format_tabline(1,'fn deserialize<D>(deserializer :D) -> Result<Self, D::Error>')
	outs += format_tabline(2,'where D: serde::de::Deserializer<\'de> {')
	outs += format_tabline(3,'let vs :StringVisitor = StringVisitor("".to_string());')
	outs += format_tabline(3,'let ores = deserializer.deserialize_str(vs);')
	outs += format_tabline(3,'let mut retv :Self = %s::%s;'%(enumname,members[0]))
	outs += format_tabline(3,'if ores.is_ok() {')
	outs += format_tabline(4,'let cmps = ores.unwrap();')
	idx = 0
	for m in members:
		if idx == 0:
			outs += format_tabline(4,'if cmps == "%s" {'%(m.lower()))
		else:
			outs += format_tabline(4,'} else if cmps == "%s" {'%(m.lower()))
		outs += format_tabline(5,'retv = %s::%s'%(enumname,m))
		idx += 1
	if idx > 0:
		outs += format_tabline(4,'}')
	outs += format_tabline(3,'}')
	outs += format_tabline(2,'return Ok(retv);')
	outs += format_tabline(1,'}')
	outs += format_tabline(0,'}')
	write_file(outs,args.output)
	sys.exit(0)
	return


REQ_KEYS=[
'Version',
'Subject',
'SignatureAlgorithm',
'Attributes',
'Extensions',
'ExtraExtensions',
'DNSNames',
'EmailAddresses',
'IPAddresses',
'URIs'
]

REQ_MAPS = {
	'Version' : 'int',
	'Subject' : 'pkixname',
	'Attributes': 'attribute_set',
	'SignatureAlgorithm': 'algorithm',
	'Extensions' : 'pkix_extensions',
	'ExtraExtensions' : 'pkix_extensions',
	'DNSNames' : 'strings',
	'EmailAddresses' : 'strings',
	'IPAddresses' : 'ip',
	'URIs' : 'urls'
}

def genreq_handler(args,parser):
	set_logging(args)
	varname = 'req'
	defines = []
	outcodes = ''

	outcodes = format_tabline(0,'func parse_x509_req_json(jsonfile string) (%s *x509.CertificateRequest,err error) {'%(varname))
	outcodes += format_tabline(1,'var s string')
	outcodes += format_tabline(1,'var mapv map[string]interface{}')
	outcodes += format_tabline(1,'var intval interface{}')
	outcodes += format_tabline(1,'var ok bool')
	outcodes += format_tabline(1,'var valmap map[string]interface{}')
	outcodes += format_tabline(1,'var valarr []interface{}')
	outcodes += format_tabline(1,'var arrs []string')
	outcodes += format_tabline(1,'var vals string')
	outcodes += format_tabline(1,'s , err= fileop.ReadFile(jsonfile)')
	outcodes += format_tabline(1,'if err != nil {')
	outcodes += format_tabline(2,'return')
	outcodes += format_tabline(1,'}')

	outcodes += format_tabline(1,'')
	outcodes += format_tabline(1,'mapv, err = jsonext.GetJsonMap(s)')
	outcodes += format_tabline(1,'if err != nil {')
	outcodes += format_tabline(2,'return')
	outcodes += format_tabline(1,'}')

	outcodes += format_tabline(1,'')
	outcodes += format_tabline(1,'%s = &x509.CertificateRequest{}'%(varname))

	for k in REQ_KEYS:
		types = REQ_MAPS[k]
		logging.info('[%s]=[%s]'%(k,types))
		funcname = 'format_%s_code'%(types)
		m = importlib.import_module(__name__)
		val = getattr(m,funcname)
		kd,code = val(k,varname)
		defines.append(kd)
		outcodes += format_tabline(1,'')
		outcodes += code

	outcodes += format_tabline(1,'')
	outcodes += format_tabline(1,'err = nil')
	outcodes += format_tabline(1,'return')
	outcodes += format_tabline(0,'}')
	write_file(outcodes,args.output)
	defs = 'const (\n'
	for d in defines:
		defs += format_tabline(1,'%s'%(d))
	defs += ')\n'
	sys.stdout.write('%s'%(defs))
	sys.exit(0)
	return

def main():
    commandline='''
    {
        "input|i" : null,
        "output|o" : null,
        "gencode<gencode_handler>##to format code to output##" : {
        	"$" : 0
        },
        "grconst<grconst_handler>##to format rust const values##" : {
        	"$" : 0
        },
        "genenum<genenum_handler>##enumname to generate enum with Partial##" : {
        	"$" : 1
        },
        "genreq<genreq_handler>##to format code for x509req##" : {
        	"$" : 0
        },
        "genenumserde<genenumserde_handler>##enumname to generate enum with serde trait##" : {
        	"$" : 1
        }
    }
    '''
    parser = extargsparse.ExtArgsParse()
    parser.load_command_line_string(commandline)
    load_log_commandline(parser)
    parser.parse_command_line(None,parser)
    raise Exception('can not reach here')
    return

if __name__ == '__main__':
    main()    