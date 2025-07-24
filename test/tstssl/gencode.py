#! /usr/bin/env python

import extargsparse
import logging
import os
import sys
import re


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

ARRS_DEFINE = []

def split_capital_name(name):
	outarr = []
	curb = b''
	inb = None
	if sys.version[0] == '3':
		inb = name.encode('utf-8')
	else:
		inb = byte(name)
	idx = 0
	while idx < len(b):
		if inb[idx] >= b'A' and inb[idx] <= b'Z':
			if len(curb) > 0:
				if sys.version[0] == '3':
					curs = curb.decode('utf-8')
				else:
					curs = string(curb)
				outarr.append(curs)
			curb = b''
		curb += inb[idx]
		idx += 1
	if len(curb) > 0:
		if sys.version[0] == '3':
			curs = curb.decode('utf-8')
		else:
			curs = string(curb)
		outarr.append(curs)
	return outarr

def format_keyword(name):
	outarr = split_capital_name(name)
	keyword = 'KEYWORD'
	for s in outarr:
		keyword += '_%s'%(s.upper())
	return keyword


def format_tabline(tab,l):
	s = ''
	for i in range(tab):
		s += '    '
	s += '%s'%(l)
	return s


def format_arrs_code(name):
	keyword = format_keyword(name)
	runcode = format_tabline(1,'x509temp.%s = []string{}'%(name))
	runcode += format_tabline(1,'arrs, err = mapv[%s].([]string)'%(keyword))
	runcode += format_tabline(1,'if ok {')
	runcode += format_tabline(2,'logutil.Debug("[%%s] parse", %s)'%(keyword))
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
	runcode += format_tabline(2,'x509temp.%s, err = get_int_val(intval,%s)'%(name,keyword))
	runcode += format_tabline(2,'if err != nil {')
	runcode += format_tabline(3,'return')
	runcode += format_tabline(2,'}')
	runcode += format_tabline(1,'}')
	kdefine = '%s = "%s"'%(keyword,name.lower())
	return kdefine,runcode


def format_bool_code(name):
	keyword = format_keyword(name)
	runcode = format_tabline(1,'x509temp.%s = false'%(name))
	runcode += format_tabline(1,'bval, ok = mapv[%s].(bool)'%(keyword))
	runcode += format_tabline(1,'if ok {')
	runcode += format_tabline(2,'logutil.Debug("[%%s] parse",%s)'%(keyword))
	runcode += format_tabline(2,'x509temp.%s = bval'%(name))
	runcode += format_tabline(1,'}')
	kdefine = '%s = "%s"'%(keyword,name.lower())
	return kdefine,runcode

def format_netip_code(name):
	keyword = format_keyword(name)
	runcode = format_tabline(1,'x509temp.%s = []*net.NetIP{}'%(name))
	runcode += format_tabline(1,'intval, ok = mapv[%s].([]interface{})'%(keyword))
	runcode += format_tabline(1,'if ok {')
	runcode += format_tabline(2,'logutil.Debug("[%%s] parse",%s)'%(keyword))
	runcode += format_tabline(2,'x509temp.%s , err = get_netip_value(intval, %s)'%(name,keyword))
	runcode += format_tabline(2,'if err != nil {')
	runcode += format_tabline(3,'return')
	runcode += format_tabline(2,'}')
	runcode += format_tabline(1,'}')
	kdefine = '%s = "%s"'%(keyword,name.lower())
	return kdefine,runcode

def format_ip_code(name):
	keyword = format_keyword(name)
	runcode = format_tabline(1,'x509temp.%s = []net.IP{}'%(name))
	runcode += format_tabline(1,'intval, ok = mapv[%s].([]interface{})'%(keyword))
	runcode += format_tabline(1,'if ok {')
	runcode += format_tabline(2,'logutil.Debug("[%%s] parse",%s)'%(keyword))
	runcode += format_tabline(2,'x509temp.%s , err = get_ip_value(intval, %s)'%(name,keyword))
	runcode += format_tabline(2,'if err != nil {')
	runcode += format_tabline(3,'return')
	runcode += format_tabline(2,'}')
	runcode += format_tabline(1,'}')
	kdefine = '%s = "%s"'%(keyword,name.lower())
	return kdefine,runcode

def format_time_code(name,after = False):
	keyword = format_keyword(name)
	if after :
		runcode = format_tabline(1,'x509temp.%s = time.Now().AddDate(20,0,0)'%(name))
	else:
		runcode = format_tabline(1,'x509temp.%s = time.Now()'%(name))
	runcode += format_tabline(1,'sval, ok = mapv[%s].(string)'%(keyword))
	runcode += format_tabline(1,'if ok {')
	runcode += format_tabline(2,'logutil.Debug("[%%s] parse",%s)'%(keyword))
	runcode += format_tabline(2,'x509temp.%s , err = get_time_value(sval, %s)'%(name,keyword))
	runcode += format_tabline(2,'if err != nil {')
	runcode += format_tabline(3,'return')
	runcode += format_tabline(2,'}')
	runcode += format_tabline(1,'}')
	kdefine = '%s = "%s"'%(keyword,name.lower())
	return kdefine,runcode



def gencode_handler(args,parser):
	set_logging(args)
	consts = []
	sys.exit(0)
	return

def main():
    commandline='''
    {
        "input|i" : null,
        "output|o" : null,
        "gencode<gencode_handler>##to format code to output##" : {
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