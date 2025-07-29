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
	'KEY_USAGE_C_R_L_SIGN' : 'KEY_USAGE_CRL_SIGN'
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


def format_strings_code(name):
	keyword = format_keyword(name)
	runcode = format_tabline(1,'x509temp.%s = []string{}'%(name))
	runcode += format_tabline(1,'valarr, ok = mapv[%s].([]interface{})'%(keyword))
	runcode += format_tabline(1,'if ok {')
	runcode += format_tabline(2,'logutil.Debug("[%%s] parse", %s)'%(keyword))
	runcode += format_tabline(2,'arrs, err = trans_inter_to_string(valarr,%s)'%(keyword))
	runcode += format_tabline(2, 'if err != nil {')
	runcode += format_tabline(3,'return')
	runcode += format_tabline(2,'}')
	runcode += format_tabline(2,'x509temp.%s = arrs'%(name))
	runcode += format_tabline(1,'}')
	kdefine = '%s = "%s"'%(keyword,name.lower())
	return kdefine,runcode

def format_int_code(name):
	keyword = format_keyword(name)
	runcode = format_tabline(1,'x509temp.%s = 0'%(name))
	runcode += format_tabline(1,'intval, ok = mapv[%s]'%(keyword))
	runcode += format_tabline(1,'if ok {')
	runcode += format_tabline(2,'logutil.Debug("[%%s] parse", %s)'%(keyword))
	runcode += format_tabline(2,'x509temp.%s, err = get_int_value(intval,%s)'%(name,keyword))
	runcode += format_tabline(2,'if err != nil {')
	runcode += format_tabline(3,'return')
	runcode += format_tabline(2,'}')
	runcode += format_tabline(1,'}')
	kdefine = '%s = "%s"'%(keyword,name.lower())
	return kdefine,runcode

def format_keyusage_code(name):
	keyword = format_keyword(name)
	runcode = format_tabline(1,'x509temp.%s = 0'%(name))
	runcode += format_tabline(1,'valarr, ok = mapv[%s].([]interface{})'%(keyword))
	runcode += format_tabline(1,'if ok {')
	runcode += format_tabline(2,'logutil.Debug("[%%s] parse", %s)'%(keyword))
	runcode += format_tabline(2,'arrs, err = trans_inter_to_string(valarr,%s)'%(keyword))
	runcode += format_tabline(2, 'if err != nil {')
	runcode += format_tabline(3,'return')
	runcode += format_tabline(2,'}')
	runcode += format_tabline(2,'x509temp.%s, err = get_keyusage_value(arrs,%s)'%(name,keyword))
	runcode += format_tabline(2,'if err != nil {')
	runcode += format_tabline(3,'return')
	runcode += format_tabline(2,'}')
	runcode += format_tabline(1,'}')
	kdefine = '%s = "%s"'%(keyword,name.lower())
	return kdefine,runcode

def format_bytes_code(name):
	keyword = format_keyword(name)
	runcode = format_tabline(1,'x509temp.%s = []byte{}'%(name))
	runcode += format_tabline(1,'valarr, ok = mapv[%s].([]interface{})'%(keyword))
	runcode += format_tabline(1,'if ok {')
	runcode += format_tabline(2,'logutil.Debug("[%%s] parse", %s)'%(keyword))
	runcode += format_tabline(2,'x509temp.%s, err = get_bytes_value(valarr,%s)'%(name,keyword))
	runcode += format_tabline(2,'if err != nil {')
	runcode += format_tabline(3,'return')
	runcode += format_tabline(2,'}')
	runcode += format_tabline(1,'}')
	kdefine = '%s = "%s"'%(keyword,name.lower())
	return kdefine,runcode

def format_bool_code(name):
	keyword = format_keyword(name)
	runcode = format_tabline(1,'x509temp.%s = false'%(name))
	runcode += format_tabline(1,'valb, ok = mapv[%s].(bool)'%(keyword))
	runcode += format_tabline(1,'if ok {')
	runcode += format_tabline(2,'logutil.Debug("[%%s] parse",%s)'%(keyword))
	runcode += format_tabline(2,'x509temp.%s = valb'%(name))
	runcode += format_tabline(1,'}')
	kdefine = '%s = "%s"'%(keyword,name.lower())
	return kdefine,runcode

def format_ipnet_code(name):
	keyword = format_keyword(name)
	runcode = format_tabline(1,'x509temp.%s = []*net.IPNet{}'%(name))
	runcode += format_tabline(1,'valarr, ok = mapv[%s].([]interface{})'%(keyword))
	runcode += format_tabline(1,'if ok {')
	runcode += format_tabline(2,'logutil.Debug("[%%s] parse",%s)'%(keyword))
	runcode += format_tabline(2,'x509temp.%s , err = get_netip_value(valarr, %s)'%(name,keyword))
	runcode += format_tabline(2,'if err != nil {')
	runcode += format_tabline(3,'return')
	runcode += format_tabline(2,'}')
	runcode += format_tabline(1,'}')
	kdefine = '%s = "%s"'%(keyword,name.lower())
	return kdefine,runcode

def format_ip_code(name):
	keyword = format_keyword(name)
	runcode = format_tabline(1,'x509temp.%s = []net.IP{}'%(name))
	runcode += format_tabline(1,'valarr, ok = mapv[%s].([]interface{})'%(keyword))
	runcode += format_tabline(1,'if ok {')
	runcode += format_tabline(2,'logutil.Debug("[%%s] parse",%s)'%(keyword))
	runcode += format_tabline(2,'x509temp.%s , err = get_ip_value(valarr, %s)'%(name,keyword))
	runcode += format_tabline(2,'if err != nil {')
	runcode += format_tabline(3,'return')
	runcode += format_tabline(2,'}')
	runcode += format_tabline(1,'}')
	kdefine = '%s = "%s"'%(keyword,name.lower())
	return kdefine,runcode

def format_time_code(name):
	keyword = format_keyword(name)
	runcode = format_tabline(1,'x509temp.%s = time.Now()'%(name))
	runcode += format_tabline(1,'vals, ok = mapv[%s].(string)'%(keyword))
	runcode += format_tabline(1,'if ok {')
	runcode += format_tabline(2,'logutil.Debug("[%%s] parse",%s)'%(keyword))
	runcode += format_tabline(2,'x509temp.%s , err = get_time_value(vals, %s)'%(name,keyword))
	runcode += format_tabline(2,'if err != nil {')
	runcode += format_tabline(3,'return')
	runcode += format_tabline(2,'}')
	runcode += format_tabline(1,'}')
	kdefine = '%s = "%s"'%(keyword,name.lower())
	return kdefine,runcode

def format_timeafter_code(name):
	keyword = format_keyword(name)
	runcode = format_tabline(1,'x509temp.%s = time.Now().AddDate(20,0,0)'%(name))
	runcode += format_tabline(1,'vals, ok = mapv[%s].(string)'%(keyword))
	runcode += format_tabline(1,'if ok {')
	runcode += format_tabline(2,'logutil.Debug("[%%s] parse",%s)'%(keyword))
	runcode += format_tabline(2,'x509temp.%s , err = get_time_value(vals, %s)'%(name,keyword))
	runcode += format_tabline(2,'if err != nil {')
	runcode += format_tabline(3,'return')
	runcode += format_tabline(2,'}')
	runcode += format_tabline(1,'}')
	kdefine = '%s = "%s"'%(keyword,name.lower())
	return kdefine,runcode


def format_bigint_code(name):
	keyword = format_keyword(name)
	runcode = format_tabline(1,'x509temp.%s = big.NewInt(0)'%(name))
	runcode += format_tabline(1,'vals, ok = mapv[%s].(string)'%(keyword))
	runcode += format_tabline(1,'if ok {')
	runcode += format_tabline(2,'x509temp.%s , err = get_bigint_value(vals,%s)'%(name,keyword))
	runcode += format_tabline(2,'if err != nil {')
	runcode += format_tabline(3,'return')
	runcode += format_tabline(2,'}')
	runcode += format_tabline(1,'}')
	kdefine = '%s = "%s"'%(keyword,name.lower())
	return kdefine,runcode

def format_pkixname_code(name):
	keyword = format_keyword(name)
	runcode = format_tabline(1,'x509temp.%s = pkix.Name{}'%(name))
	runcode += format_tabline(1,'valmap, ok = mapv[%s].(map[string]interface{})'%(keyword))
	runcode += format_tabline(1,'if ok {')
	runcode += format_tabline(2,'logutil.Debug("[%%s] parse",%s)'%(keyword))
	runcode += format_tabline(2,'x509temp.%s , err = get_pkixname_value(valmap,%s)'%(name,keyword))
	runcode += format_tabline(2,'if err != nil {')
	runcode += format_tabline(3,'return')
	runcode += format_tabline(2,'}')
	runcode += format_tabline(1,'}')
	kdefine = '%s = "%s"'%(keyword,name.lower())
	return kdefine,runcode

def format_objoids_code(name):
	keyword = format_keyword(name)
	runcode = format_tabline(1,'x509temp.%s = []asn1.ObjectIdentifier{}'%(name))
	runcode += format_tabline(1,'valarr, ok = mapv[%s].([]interface{})'%(keyword))
	runcode += format_tabline(1,'if ok {')
	runcode += format_tabline(2,'logutil.Debug("[%%s] parse",%s)'%(keyword))
	runcode += format_tabline(2,'x509temp.%s , err = get_objoids_value(valarr,%s)'%(name,keyword))
	runcode += format_tabline(2,'if err != nil {')
	runcode += format_tabline(3,'return')
	runcode += format_tabline(2,'}')
	runcode += format_tabline(1,'}')
	kdefine = '%s = "%s"'%(keyword,name.lower())
	return kdefine,runcode

def format_oids_code(name):
	keyword = format_keyword(name)
	runcode = format_tabline(1,'x509temp.%s = []x509.OID{}'%(name))
	runcode += format_tabline(1,'valarr, ok = mapv[%s].([]interface{})'%(keyword))
	runcode += format_tabline(1,'if ok {')
	runcode += format_tabline(2,'logutil.Debug("[%%s] parse",%s)'%(keyword))
	runcode += format_tabline(2,'arrs, err = trans_inter_to_string(valarr,%s)'%(keyword))
	runcode += format_tabline(2, 'if err != nil {')
	runcode += format_tabline(3,'return')
	runcode += format_tabline(2,'}')
	runcode += format_tabline(2,'x509temp.%s , err = get_oids_value(arrs,%s)'%(name,keyword))
	runcode += format_tabline(2,'if err != nil {')
	runcode += format_tabline(3,'return')
	runcode += format_tabline(2,'}')
	runcode += format_tabline(1,'}')
	kdefine = '%s = "%s"'%(keyword,name.lower())
	return kdefine,runcode

def format_urls_code(name):
	keyword = format_keyword(name)
	runcode = format_tabline(1,'x509temp.%s = []*url.URL{}'%(name))
	runcode += format_tabline(1,'valarr, ok = mapv[%s].([]interface{})'%(keyword))
	runcode += format_tabline(1,'if ok {')
	runcode += format_tabline(2,'logutil.Debug("[%%s] parse",%s)'%(keyword))
	runcode += format_tabline(2,'arrs, err = trans_inter_to_string(valarr,%s)'%(keyword))
	runcode += format_tabline(2, 'if err != nil {')
	runcode += format_tabline(3,'return')
	runcode += format_tabline(2,'}')
	runcode += format_tabline(2,'x509temp.%s , err = get_urls_value(arrs,%s)'%(name,keyword))
	runcode += format_tabline(2,'if err != nil {')
	runcode += format_tabline(3,'return')
	runcode += format_tabline(2,'}')
	runcode += format_tabline(1,'}')
	kdefine = '%s = "%s"'%(keyword,name.lower())
	return kdefine,runcode

def format_extkeyusage_code(name):
	keyword = format_keyword(name)
	runcode = format_tabline(1,'x509temp.%s = []x509.ExtKeyUsage{}'%(name))
	runcode += format_tabline(1,'valarr, ok = mapv[%s].([]interface{})'%(keyword))
	runcode += format_tabline(1,'if ok {')
	runcode += format_tabline(2,'logutil.Debug("[%%s] parse",%s)'%(keyword))
	runcode += format_tabline(2,'arrs, err = trans_inter_to_string(valarr,%s)'%(keyword))
	runcode += format_tabline(2, 'if err != nil {')
	runcode += format_tabline(3,'return')
	runcode += format_tabline(2,'}')
	runcode += format_tabline(2,'x509temp.%s , err = get_key_ext_usage(arrs,%s)'%(name,keyword))
	runcode += format_tabline(2,'if err != nil {')
	runcode += format_tabline(3,'return')
	runcode += format_tabline(2,'}')
	runcode += format_tabline(1,'}')
	kdefine = '%s = "%s"'%(keyword,name.lower())
	return kdefine,runcode

def format_algorithm_code(name):
	keyword = format_keyword(name)
	runcode = format_tabline(1,'x509temp.%s = x509.SHA256WithRSA'%(name))
	runcode += format_tabline(1,'vals, ok = mapv[%s].(string)'%(keyword))
	runcode += format_tabline(1,'if ok {')
	runcode += format_tabline(2,'logutil.Debug("[%%s] parse",%s)'%(keyword))
	runcode += format_tabline(2,'x509temp.%s , err = get_algorithm_value(vals,%s)'%(name,keyword))
	runcode += format_tabline(2,'if err != nil {')
	runcode += format_tabline(3,'return')
	runcode += format_tabline(2,'}')
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