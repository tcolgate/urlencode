package main

import (
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"net/url"
	"os"
)

func main() {
	var err error
	defer func() {
		if err != nil {
			fmt.Fprintf(os.Stderr, "%v\n", err)
		}
		os.Exit(1)
	}()

	ustr := flag.String("u", "", "a base url")
	k := flag.String("k", "", "a query key, default to \"data\" if -u is specified")
	flag.Parse()

	var bs []byte
	switch len(flag.Args()) {
	case 0:
		bs, err = ioutil.ReadAll(os.Stdin)
		if err != nil {
			return
		}
	case 1:
		bs = []byte(flag.Args()[0])
	default:
		err = errors.New("too many args, give one string to encode, or pass text on stdin")
	}

	out := url.QueryEscape(string(bs))

	var u *url.URL
	var vs url.Values
	if *ustr != "" {
		u, err = url.Parse(*ustr)
		if err != nil {
			return
		}
		vs, err = url.ParseQuery(u.RawQuery)
		if *k != "" {
			vs.Add(*k, string(bs))
		} else {
			vs.Add("data", string(bs))
		}
		u.RawQuery = vs.Encode()
		out = u.String()
	}

	if *ustr == "" && *k != "" {
		vs := url.Values{*k: []string{string(bs)}}
		out = vs.Encode()
	}

	fmt.Println(out)
}
