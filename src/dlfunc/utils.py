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


def genlxfunc_handler(args,parser):
    set_logging(args)
    num = 10
    basename  = 'call'
    bridgname = 'bridge'
    if len(args.subnargs) > 0:
        num = parse_int(args.subnargs[0])
    if len(args.subnargs) > 1:
        basename = args.subnargs[1]
    if len(args.subnargs) > 2:
        bridgname = args.subnargs[2]

    i = 0
    outs = ''
    while i < num:
        if i == 0:
            outs += format_tab_line(1,'if num == %d {'%(i))
        else:
            outs += format_tab_line(1,'} else if num == %d {'%(i))

        outs += format_tab_line(2,'func_ptr := C.%s_%d_func_t(fptr.funcptr)'%(basename,i))
        curs = 'retval = C.%s_call_%d(func_ptr'%(bridgname,i)
        j = 0
        while j < i:
            curs += ' ,'
            curs += 'input_var[%d]'%(j)
            j += 1
        curs += ')'
        outs += format_tab_line(2,'%s'%(curs))
        i += 1

    outs += format_tab_line(1,'} else {')
    outs += format_tab_line(2,'err = fmt.Errorf("%%d >= %d",num)')
    outs += format_tab_line(2,'return')
    outs += format_tab_line(1,'}')

    write_file(outs,args.output)
    sys.exit(0)
    return

def genccode_handler(args,parser):
    set_logging(args)
    num = 10
    basename  = 'call'
    bridgname = 'bridge'
    if len(args.subnargs) > 0:
        num = parse_int(args.subnargs[0])
    if len(args.subnargs) > 1:
        basename = args.subnargs[1]
    if len(args.subnargs) > 2:
        bridgname = args.subnargs[2]

    outs = ''

    #// typedef int (*call_0_func_t) ();
    #//
    #// int
    #// bridge_call_0(call_0_func_t f)
    #// {
    #//      return f();
    #// }
    #//
    i = 0
    while i < num:
        curs = ''
        j = 0
        while j < i:
            if j > 0:
                curs += ' ,'
            curs += 'unsigned long'
            j += 1
        outs += format_tab_line(0,'// typedef int (*%s_%d_func_t)(%s);'%(basename,i,curs))
        outs += format_tab_line(0,'// ')
        outs += format_tab_line(0,'// int')
        curs = ''
        j = 0
        while j < i:
            curs += ' ,'
            curs += 'unsigned long a%d'%(j)
            j += 1
        outs += format_tab_line(0,'// %s_call_%d(%s_%d_func_t f%s)'%(bridgname,i,basename,i,curs))
        outs += format_tab_line(0,'// {')
        j = 0
        curs = ''
        while j < i:
            if j > 0:
                curs += ','
            curs += 'a%d'%(j)
            j += 1
        outs += format_tab_line(0,'//     return f(%s);'%(curs))
        outs += format_tab_line(0,'// }')
        outs += format_tab_line(0,'// ')
        i += 1
    write_file(outs,args.output)
    sys.exit(0)
    return



def main():
    commandline='''
    {
        "input|i" : null,
        "output|o" : null,
        "genlxfunc<genlxfunc_handler>##[num] [prefix] to generate function with default prefix print num default 10##" : {
            "$" : "*"
        },
        "genccode<genccode_handler>##num [prefix] to generate c code before import C##" : {
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