#! /usr/bin/env python

import extargsparse
import logging
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


def read_file(infile=None):
    fin = sys.stdin
    if infile is not None:
        fin = open(infile,'r+b')
    rets = ''
    for l in fin:
        s = l
        if 'b' in fin.mode:
            if sys.version[0] == '3':
                s = l.decode('utf-8')
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

def format_tab_line(tab,s):
    rets = ''
    for i in range(tab):
        rets += '    '
    rets += s
    rets += '\n'
    return rets


def gencfg_handler(args,parser):
    set_logging(args)
    s = read_file(args.input)
    sarr = re.split('\n',s)
    outs = ''
    structs = ''

    strexpr = re.compile('^\\s*type\\s+([^\\s]+)\\s+struct\\s+\\{$')
    strname = None
    mname = None
    mtype = None
    mjsonexpr = re.compile('.*`json\\s*:.*`')
    for l in sarr:
        l = l.rstrip('\r\n')
        m = strexpr.findall(l)
        if m is not None and len(m) > 0:
            strname = m[0]
            structs += '%s\n'%(l)
        else:
            cl = l.strip(' \t')
            sarr = re.split('\\s+',cl)
            if len(sarr) >= 2:
                mname = sarr[0]
                mtype = sarr[1]
                outs += format_tab_line(0,'// format %s type %s set'%(mname,mtype))
                outs += format_tab_line(0,'func (retp *%s) Set%s(val %s) *%s {'%(strname,mname,mtype,strname))
                outs += format_tab_line(1,'retp.%s = val'%(mname))
                outs += format_tab_line(1,'return retp')
                outs += format_tab_line(0,'}')
                outs += format_tab_line(0,' ')
                outs += format_tab_line(0,'// get %s type %s value'%(mname,mtype))
                outs += format_tab_line(0, 'func (retp *%s) Get%s() %s {'%(strname,mname,mtype) )
                outs += format_tab_line(1,'return retp.%s'%(mname))
                outs += format_tab_line(0,'}')
                outs += format_tab_line(0,' ')
                # now to make sure the
                nl = None
                if sarr[0][0] >= 'A' and sarr[0][0] <= 'Z':
                    if not mjsonexpr.match(l):
                        lname = mname.lower()
                        nl = format_tab_line(1,'%s %s `json:"%s"`'%(mname,mtype,lname))
                        nl = nl.rstrip('\r\n')
                if nl is not None:
                    structs += '%s\n'%(nl)
                else:
                    structs += '%s\n'%(l)

            else:
                structs += '%s\n'%(l)
    nouts = ''
    nouts += structs
    nouts += format_tab_line(0,'')
    nouts += outs
    write_file(nouts,args.output)
    sys.exit(0)
    return



def main():
    commandline='''
    {
        "input|i" : null,
        "output|o" : null,
        "gencfg<gencfg_handler>##from input file go code to output go code##" : {
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