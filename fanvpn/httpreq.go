package main

import (
	"bytes"
	"crypto/tls"
	"dbgutil"
	"io/ioutil"
	"logutil"
	"net/http"
	"strings"
	"time"
)

func GetHttp(url string, mills int, bcheck bool) (respstr string, err error) {
	var nbytes []byte
	var req *http.Request
	var resp *http.Response
	var client *http.Client
	var secure bool = false

	if strings.HasPrefix(url, "https:") {
		secure = true
	}

	nbytes = []byte{}
	req, err = http.NewRequest("GET", url, bytes.NewBuffer(nbytes))
	if err != nil {
		err = dbgutil.FormatError("can not format request error[%s]", err.Error())
		return
	}
	if mills == 0 {
		if secure {
			if bcheck {
				client = &http.Client{
					Transport: &http.Transport{
						TLSClientConfig: &tls.Config{InsecureSkipVerify: false},
					},
				}
			} else {
				client = &http.Client{
					Transport: &http.Transport{
						TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
					},
				}
			}
		} else {
			client = &http.Client{}
		}

	} else {
		if secure {
			if bcheck {
				client = &http.Client{
					Timeout: time.Duration(mills) * time.Millisecond,
					Transport: &http.Transport{
						TLSClientConfig: &tls.Config{InsecureSkipVerify: false},
					},
				}
			} else {
				client = &http.Client{
					Timeout: time.Duration(mills) * time.Millisecond,
					Transport: &http.Transport{
						TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
					},
				}
			}
		} else {
			client = &http.Client{
				Timeout: time.Duration(mills) * time.Millisecond,
			}
		}
	}

	resp, err = client.Do(req)
	if err != nil {
		err = dbgutil.FormatError("get [%s]  error[%s]", url, err.Error())
		return
	}
	defer resp.Body.Close()
	nbytes, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		err = dbgutil.FormatError("read [%s] response error[%s]", url, err.Error())
		return
	}
	respstr = string(nbytes)
	logutil.Debug("req-----------\n%s\n------------\nresp++++++++++++++++\n%s\n++++++++++++++++\n", url, respstr)
	err = nil
	return

}
