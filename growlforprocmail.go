package main

import (
	"bufio"
	"bytes"
	"code.google.com/p/go.text/encoding/simplifiedchinese"
	"code.google.com/p/go.text/transform"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"regexp"
	"strconv"
	"strings"
)

type qDecoder struct {
	r       io.Reader
	scratch [2]byte
}

func GBKDecode(s []byte) ([]byte, error) {
	I := bytes.NewReader(s)
	O := transform.NewReader(I, simplifiedchinese.GBK.NewDecoder())
	d, e := ioutil.ReadAll(O)
	if e != nil {
		return nil, e
	}
	return d, nil
}

func GBKEncode(s []byte) ([]byte, error) {
	I := bytes.NewReader(s)
	O := transform.NewReader(I, simplifiedchinese.GBK.NewEncoder())
	d, e := ioutil.ReadAll(O)
	if e != nil {
		return nil, e
	}
	return d, nil
}

func decodeRFC2047Word(s string) (string, error) {
	fields := strings.Split(s, "?")
	if len(fields) != 5 || fields[0] != "=" || fields[4] != "=" {
		return "", errors.New("mail: address not RFC 2047 encoded")
	}
	charset, enc := strings.ToLower(fields[1]), strings.ToLower(fields[2])
	// if charset != "iso-8859-1" && charset != "utf-8" {
	// return "", fmt.Errorf("mail: charset not supported: %q", charset)
	// }
	in := bytes.NewBufferString(fields[3])
	var r io.Reader
	switch enc {
	case "b":
		r = base64.NewDecoder(base64.StdEncoding, in)
	case "q":
		r = qDecoder{r: in}
	default:
		return "", fmt.Errorf("mail: RFC 2047 encoding not supported: %q", enc)
	}
	dec, err := ioutil.ReadAll(r)
	if err != nil {
		return "", err
	}
	switch charset {
	case "iso-8859-1":
		b := new(bytes.Buffer)
		for _, c := range dec {
			b.WriteRune(rune(c))
		}
		return b.String(), nil
	case "utf-8":
		return string(dec), nil
	case "gbk":
		str, err := GBKDecode(dec)
		return string(str), err
	case "gb2312":
		str, err := GBKDecode(dec)
		return string(str), err
	default:
		return string(dec), nil
	}
	panic("unreachable")
}

func (qd qDecoder) Read(p []byte) (n int, err error) {
	// This method writes at most one byte into p.
	if len(p) == 0 {
		return 0, nil
	}
	if _, err := qd.r.Read(qd.scratch[:1]); err != nil {
		return 0, err
	}
	switch c := qd.scratch[0]; {
	case c == '=':
		if _, err := io.ReadFull(qd.r, qd.scratch[:2]); err != nil {
			return 0, err
		}
		x, err := strconv.ParseInt(string(qd.scratch[:2]), 16, 64)
		if err != nil {
			return 0, fmt.Errorf("mail: invalid RFC 2047 encoding: %q", qd.scratch[:2])
		}
		p[0] = byte(x)
	case c == '_':
		p[0] = ' '
	default:
		p[0] = c
	}
	return 1, nil
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			// You may check here if err == io.EOF
			break
		}
		length := len(line)
		if length >= 9 {
			reg := regexp.MustCompile(`From[: ]+(.*) +(<.*>)`)
			if reg.MatchString(line) {
				matched := reg.FindAllStringSubmatch(line, -1)
				str, err := decodeRFC2047Word(matched[0][1])
				if err != nil {
					fmt.Println(matched[0][1] + matched[0][2])
				} else {
					fmt.Println(str + matched[0][2])
				}
			}
			reg = regexp.MustCompile(`Subject[: ]+(.*)`)
			if reg.MatchString(line) {
				matched := reg.FindAllStringSubmatch(line, 1)[0][1]
				str, err := decodeRFC2047Word(matched)
				if err != nil {
					fmt.Println(matched)
				} else {
					fmt.Println(str)
				}
			}
			reg = regexp.MustCompile(`filename[= ](.*)`)
			if reg.MatchString(line) {
				matched := reg.FindAllStringSubmatch(line, 1)[0][1]
				str, err := decodeRFC2047Word(matched)
				if err != nil {
					fmt.Println(matched)
				} else {
					fmt.Println(str)
				}
			}

		}
	}
}
