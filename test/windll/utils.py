#! /usr/bin/env python

import extargsparse
import sys
import logging


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

def parse_int(v):
    c = v
    base = 10
    if c.startswith('0x') or c.startswith('0X') :
        base = 16
        c = c[2:]
    elif c.startswith('x') or c.startswith('X'):
        base = 16
        c = c[1:]
    return int(c,base)

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


def gendecl_handler(args,parser):
	set_logging(args)
	num = 10
	prefix = 'print'
	if len(args.subnargs) > 0:
		num = parse_int(args.subnargs[0])
	if len(args.subnargs) > 1:
		prefix = args.subnargs[1]

	i = 0
	outs = ''
	while i < num:
		outs += 'WINLIB_API int %s_%d('%(prefix,i)
		j = 0
		while j < i:
			if j > 0:
				outs += ', '
			outs += 'int a%d'%(j)
			j += 1
		outs += ');\n'
		i += 1

	write_file(outs,args.output)
	sys.exit(0)
	return

def genfunc_handler(args,parser):
	set_logging(args)
	num = 10
	prefix = 'print'
	if len(args.subnargs) > 0:
		num = parse_int(args.subnargs[0])
	if len(args.subnargs) > 1:
		prefix = args.subnargs[1]

	i = 0
	outs = ''
	while i < num:
		if i > 0:
			outs += '\n\n'
		outs += 'int %s_%d('%(prefix,i)
		j = 0
		while j < i:
			if j > 0:
				outs += ', '
			outs += 'int a%d'%(j)
			j += 1
		outs += ')\n'
		outs += format_tab_line(0,'{')
		outs += format_tab_line(1,'printf("call %s_%d\\n");'%(prefix,i))
		outs += format_tab_line(0,' ')
		j = 0
		while j < i:
			outs += format_tab_line(1,'printf("a%d = %%d\\n",a%d);'%(j,j))
			j += 1
		outs += format_tab_line(1,'return 0;')
		outs += format_tab_line(0,'}')
		i += 1

	write_file(outs,args.output)
	sys.exit(0)
	return


def genchardecl_handler(args,parser):
    set_logging(args)
    num = 10
    prefix = 'charfunc'
    if len(args.subnargs) > 0:
        num = parse_int(args.subnargs[0])
    if len(args.subnargs) > 1:
        prefix = args.subnargs[1]

    i = 0
    outs = ''
    while i < num:
        outs += 'WINLIB_API char* %s_%d('%(prefix,i)
        j = 0
        while j < i:
            if j > 0:
                outs += ', '
            outs += 'int a%d'%(j)
            j += 1
        outs += ');\n'
        i += 1

    write_file(outs,args.output)
    sys.exit(0)
    return


def gencharfunc_handler(args,parser):
    set_logging(args)
    num = 10
    prefix = 'charfunc'
    if len(args.subnargs) > 0:
        num = parse_int(args.subnargs[0])
    if len(args.subnargs) > 1:
        prefix = args.subnargs[1]

    i = 0
    outs = ''
    while i < num:
        if i > 0:
            outs += '\n\n'
        outs += 'char* %s_%d('%(prefix,i)
        j = 0
        while j < i:
            if j > 0:
                outs += ', '
            outs += 'int a%d'%(j)
            j += 1
        outs += ')\n'
        outs += format_tab_line(0,'{')
        outs += format_tab_line(1,'int totalv;')
        outs += format_tab_line(1,'char* pret=NULL;')
        outs += format_tab_line(1,'printf("call %s_%d\\n");'%(prefix,i))
        outs += format_tab_line(0,' ')
        j = 0
        while j < i:
            outs += format_tab_line(1,'printf("a%d = %%d\\n",a%d);'%(j,j))
            j += 1
        curs = ''
        j = 0
        while j < i:
            if j > 0:
                curs += '+'
            curs += 'a%d'%(j)
            j += 1
        if i > 0:
            outs += format_tab_line(1,'totalv = %s;'%(curs))
        else:
            outs += format_tab_line(1,'totalv = 30;')
        outs += format_tab_line(1,' ')
        outs += format_tab_line(1,'pret = (char*) malloc(totalv);')
        outs += format_tab_line(1,'printf("pret %p\\n",pret);')
        outs += format_tab_line(1,'return pret;')
        outs += format_tab_line(0,'}')
        i += 1

    write_file(outs,args.output)
    sys.exit(0)
    return


def genstrdecl_handler(args,parser):
    set_logging(args)
    num = 10
    prefix = 'strfunc'
    if len(args.subnargs) > 0:
        num = parse_int(args.subnargs[0])
    if len(args.subnargs) > 1:
        prefix = args.subnargs[1]

    i = 0
    outs = ''
    while i < num:
        outs += 'WINLIB_API char* %s_%d('%(prefix,i)
        j = 0
        while j < i:
            if j > 0:
                outs += ', '
            outs += 'char* a%d'%(j)
            j += 1
        outs += ');\n'
        i += 1

    write_file(outs,args.output)
    sys.exit(0)
    return


def genstrfunc_handler(args,parser):
    set_logging(args)
    num = 10
    prefix = 'strfunc'
    if len(args.subnargs) > 0:
        num = parse_int(args.subnargs[0])
    if len(args.subnargs) > 1:
        prefix = args.subnargs[1]

    i = 0
    outs = ''
    while i < num:
        if i > 0:
            outs += '\n\n'
        outs += 'char* %s_%d('%(prefix,i)
        j = 0
        while j < i:
            if j > 0:
                outs += ', '
            outs += 'char* a%d'%(j)
            j += 1
        outs += ')\n'
        outs += format_tab_line(0,'{')
        outs += format_tab_line(1,'int totalv;')
        outs += format_tab_line(1,'char* pret=NULL;')
        outs += format_tab_line(1,'printf("call %s_%d\\n");'%(prefix,i))
        outs += format_tab_line(0,' ')
        j = 0
        while j < i:
            outs += format_tab_line(1,'printf("a%d = [%%s] %%p\\n",a%d,a%d);'%(j,j,j))
            j += 1
        if i > 0:
            outs += format_tab_line(1,'totalv = %d;'%(i))
        else:
            outs += format_tab_line(1,'totalv = 30;')
        outs += format_tab_line(1,' ')
        outs += format_tab_line(1,'pret = (char*) malloc(totalv);')
        outs += format_tab_line(1,'printf("pret %p\\n",pret);')
        outs += format_tab_line(1,'return pret;')
        outs += format_tab_line(0,'}')
        i += 1

    write_file(outs,args.output)
    sys.exit(0)
    return


def gencallbackdecl_handler(args,parser):
    set_logging(args)
    num = 10
    prefix = 'callbackfunc'
    if len(args.subnargs) > 0:
        num = parse_int(args.subnargs[0])
    if len(args.subnargs) > 1:
        prefix = args.subnargs[1]
    idx = 0
    outs = ''
    ins = ''
    while idx < num:
        # to generate the declare
        if len(outs) > 0:
            outs += format_tab_line(0,' ')
        outs += format_tab_line(0,'// to generate callback with %d'%(idx))
        ins = ''
        jdx = 0
        while jdx < idx:
            if len(ins) > 0:
                ins += ','
            ins += 'char* a%d'%(jdx)
            jdx += 1
        outs += format_tab_line(0,'typedef int (%s_%d_func_t)(%s);'%(prefix,idx,ins))
        idx += 1

    outs += format_tab_line(0,' ')
    outs += format_tab_line(0,' ')
    idx = 0
    while idx < num:

        ins = ''
        jdx = 0
        while jdx < idx:
            ins += ','
            ins += 'char* a%d'%(jdx)
            jdx += 1
        outs += format_tab_line(0,'WINLIB_API int %s_%d(%s_%d_func_t pfunc%s);'%(prefix,idx,prefix,idx,ins))

        idx += 1

    write_file(outs,args.output)
    sys.exit(0)
    return


def gencallbackfunc_handler(args,parser):
    set_logging(args)
    logging.info('call gencallbackfunc_handler')
    num = 10
    prefix = 'callbackfunc'
    if len(args.subnargs) > 0:
        num = parse_int(args.subnargs[0])
    if len(args.subnargs) > 1:
        prefix = args.subnargs[1]
    outs = ''
    idx = 0
    while idx < num:
        outs += format_tab_line(0,' ')
        outs += format_tab_line(0,'// to call func %d params'%(idx))
        ins = ''
        jdx = 0
        while jdx < idx:
            ins += ','
            ins += 'char* a%d'%(jdx)
            jdx += 1
        outs += format_tab_line(0,'int %s_%d(%s_%d_func_t pfunc%s)'%(prefix,idx,prefix,idx,ins))
        outs += format_tab_line(0,'{')
        ins = ''
        jdx = 0
        while jdx < idx:
            if len(ins) > 0:
                ins += ','
            ins += 'a%d'%(jdx)
            jdx += 1
        outs += format_tab_line(1,'int retv;')

        outs += format_tab_line(1,'retv = pfunc(%s);'%(ins))
        jdx = 0
        while jdx < idx:
            outs += format_tab_line(1,'printf("C:a%d=[%%s]\\n",a%d);'%(jdx,jdx))
            jdx += 1
        outs += format_tab_line(1,'printf("C:retv=%d\\n",retv);')
        outs += format_tab_line(1,'return retv;')
        outs += format_tab_line(0,'}')
        idx += 1
    write_file(outs,args.output)
    sys.exit(0)
    return


def main():
    commandline='''
    {
        "input|i" : null,
        "output|o" : null,
        "gendecl<gendecl_handler>##[num] [prefix] to generate function with default prefix print num default 10##" : {
        	"$" : "*"
        },
        "genfunc<genfunc_handler>##[num] [prefix] to generate function with default prefix print num default 10##" : {
        	"$" : "*"
        },
        "genchardecl<genchardecl_handler>##[num] [prefix] to generate function with default prefix print num default 10##" : {
            "$" : "*"
        },
        "gencharfunc<gencharfunc_handler>##[num] [prefix] to generate char* functions##" : {
            "$" : "*"
        },
        "genstrdecl<genstrdecl_handler>##[num] [prefix] to generate function with default prefix print num default 10##" : {
            "$" : "*"
        },
        "genstrfunc<genstrfunc_handler>##[num] [prefix] to generate function with default prefix print num default 10##" : {
            "$" : "*"
        },
        "gencallbackdecl<gencallbackdecl_handler>##[num] [prefix] to generate function with call back functions##" : {
            "$" : "*"
        },
        "gencallbackfunc<gencallbackfunc_handler>##[num] [prefix] to generate function with call back functions##" : {
            "$" : "*"
        }
    }
    '''
    parser = extargsparse.ExtArgsParse()
    load_log_commandline(parser)
    parser.load_command_line_string(commandline)
    parser.parse_command_line(None,parser)
    raise Exception('can not reach here')
    return

if __name__ == '__main__':
    main()