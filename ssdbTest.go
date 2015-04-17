// ssdb_test
package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"strconv"
)

type Config struct {
	ssdbUrl string
	count   int
	cmd     string
	prefix  string
}

var (
	cfg = Config{
		ssdbUrl: "10.39.80.182:8888",
		count:   1000,
		cmd:     "set",
		prefix:  "",
	}
)

func configure() {
	flag.StringVar(&cfg.ssdbUrl, "s", cfg.ssdbUrl, "ssdb ip:port")
	flag.IntVar(&cfg.count, "c", cfg.count, "iterations")
	flag.StringVar(&cfg.cmd, "cmd", cfg.cmd, "[set|qpush]")
	flag.StringVar(&cfg.prefix, "p", cfg.prefix, "prefix")

	flag.Parse()

	if cfg.cmd != "set" && cfg.cmd != "qpush" {
		log.Panicln("Cmd err: %+v\n")
	}

	fmt.Printf("Config: %+v\n", cfg)
}

func ssdbConnect() *net.TCPConn {
	ssdbAddr, _ := net.ResolveTCPAddr("tcp", cfg.ssdbUrl)
	ssdb, err := net.DialTCP("tcp", nil, ssdbAddr)
	if err != nil {
		log.Panicf("Dial err: %+v\n", err.Error())
	}
	return ssdb
}

func writeCmd(conn *net.TCPConn, cmd string) {
	conn.Write([]byte(cmd))
	buf := make([]byte, 128)
	n, err := conn.Read(buf)
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Printf("Rec: %q\n", string(buf[:n]))
}

func main() {

	fmt.Println(net.ParseIP(cfg.ssdbUrl))

	configure()

	conn := ssdbConnect()
	defer conn.Close()

	for i := 0; i < cfg.count; i++ {
		var cmd string
		if cfg.cmd == "qpush" {
			cmd = fmt.Sprintf("%d\n%s\n%d\n%d\n%d\n%d\n\n", len(cfg.cmd), cfg.cmd, len(strconv.Itoa(0)), 0, len(strconv.Itoa(i)), i)
		} else {
			cmd = fmt.Sprintf("%d\n%s\n%d\n%d\n%d\n%d\n\n", len(cfg.cmd), cfg.cmd, len(strconv.Itoa(i)), i, len(strconv.Itoa(i)), i)
		}
		fmt.Printf("Sent: %q\n", cmd)
		go writeCmd(conn, cmd)
	}
}
