package main

import (
	"encoding/json"
	"fmt"
	"github.com/codegangsta/cli"
	//"github.com/jeppeter/jsonext"
	"io/ioutil"
	"os"
	"reflect"
	"regexp"
	"strings"
)

const (
	ATTR_SPLIT       = "split"
	ATTR_SPLIT_EQUAL = "split="
)

type keyAttr map[string]string

func NewKeyAttr() *keyAttr {
	p := &keyAttr{}
	return p
}

func (k *keyAttr) GetAttr(name string) (val string, err error) {
	var ok bool
	val, ok = k[name]
	if !ok {
		err = fmt.Errorf("no [%s] key", name)
		return
	}
	err = nil
	return
}

func (k *keyAttr) SetAttr(name string, value string) {
	k[name] = value
}

func (k *)

func parseAttr(attr string) (kattr map[string]string, err error) {
	var lattr string
	var splitchar string = ";"
	var splitstrings []string
	var splitexpr *regexp.Regexp
	var equalexpr *regexp.Regexp
	var vk []string
	var curs string

	kattr = nil
	err = nil
	lattr = strings.ToLower(attr)
	if strings.HasPrefix(lattr, ATTR_SPLIT_EQUAL) {
		splitchar = lattr[len(ATTR_SPLIT_EQUAL):(len(ATTR_SPLIT_EQUAL) + 1)]
		switch splitchar {
		case "\\":
			splitchar = "\\\\"
		case ".":
			splitchar = "\\."
		case "/":
			splitchar = "/"
		case ":":
			splitchar = ":"
		case "+":
			splitchar = "\\+"
		case "@":
			splitchar = "@"
		default:
			return nil, fmt.Errorf("unknown splitchar [%s]", splitchar)
		}
	}
	splitexpr, err = regexp.Compile(splitchar)
	if err != nil {
		return
	}
	equalexpr, err = regexp.Compile("=")
	if err != nil {
		return
	}

	kattr = make(map[string]string)
	splitstrings = splitexpr.Split(lattr, -1)
	for _, curs = range splitstrings {
		if strings.HasPrefix(curs, ATTR_SPLIT_EQUAL) || curs == "" {
			continue
		}
		vk = equalexpr.Split(curs, 2)
		if len(vk) < 2 {
			continue
		}
		kattr[vk[0]] = vk[1]
	}

	err = nil
	return
}

func setAttr(attr interface{}) (kattr map[string]string, err error) {
	var k string
	var v interface{}
	var vmap map[string]interface{}
	var vstr string
	kattr = make(map[string]string)

	switch attr.(type) {
	case map[string]interface{}:
		vmap = attr.(map[string]interface{})
		for k, v = range vmap {
			if strings.ToLower(k) == ATTR_SPLIT || k == "" {
				continue
			}
			switch v.(type) {
			case string:
				vstr = v.(string)
			default:
				vstr = fmt.Sprintf("%v", v)
			}
			kattr[k] = vstr
		}
	default:
		return kattr, fmt.Errorf("not valid type [%s]", reflect.TypeOf(attr))
	}

	err = nil
	return
}

func debug_kattr(pre string, kattr map[string]string) {
	fmt.Fprintf(os.Stdout, "parse [%s]\n", pre)
	fmt.Fprintf(os.Stdout, "---------------------------\n")
	for k, v := range kattr {
		fmt.Fprintf(os.Stdout, "[%s]=[%s]\n", k, v)
	}
	fmt.Fprintf(os.Stdout, "+++++++++++++++++++++++++++\n")
}

func makeStringCommand() cli.Command {
	cmd := cli.Command{}
	cmd.Name = "string"
	cmd.ShortName = "st"
	cmd.Usage = "strings..."
	cmd.Action = func(c *cli.Context) {
		if len(c.Args()) < 1 {
			fmt.Fprintf(os.Stderr, "string %s\n", cmd.Usage)
			os.Exit(5)
		}
		for _, cs := range c.Args() {
			kattr, err := parseAttr(cs)
			if err != nil {
				fmt.Fprintf(os.Stderr, "parse [%s] err[%s]\n", cs, err.Error())
				os.Exit(5)
			}
			debug_kattr(cs, kattr)
		}
	}
	return cmd
}

func makeJsonCommand() cli.Command {
	var v interface{}
	cmd := cli.Command{}
	cmd.Name = "json"
	cmd.ShortName = "js"
	cmd.Usage = "files..."
	cmd.Action = func(c *cli.Context) {
		if len(c.Args()) < 1 {
			fmt.Fprintf(os.Stderr, "json %s\n", cmd.Usage)
			os.Exit(5)
		}
		for _, curf := range c.Args() {

			data, err := ioutil.ReadFile(curf)
			if err != nil {
				fmt.Fprintf(os.Stderr, "read [%s] error [%s]\n", curf, err.Error())
				os.Exit(5)
			}

			err = json.Unmarshal(data, &v)
			if err != nil {
				fmt.Fprintf(os.Stderr, "can not parse [%s] [%s] error[%s]\n", curf, string(data), err.Error())
				os.Exit(5)
			}

			kattr, err := setAttr(v)
			if err != nil {
				fmt.Fprintf(os.Stderr, "parse error [%s]\n", err.Error())
				os.Exit(5)
			}

			debug_kattr(curf, kattr)
		}
	}
	return cmd
}

func main() {
	app := cli.NewApp()
	app.Version = "1.0.2"
	app.Commands = append(app.Commands, makeStringCommand())
	app.Commands = append(app.Commands, makeJsonCommand())
	app.Run(os.Args)

}
