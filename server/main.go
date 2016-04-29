package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha256"
	"io"
	"log"
	"math/rand"
	"net"
	"os"
	"time"

	"github.com/codegangsta/cli"
	"github.com/hashicorp/yamux"
	"github.com/spance/suft/protocol"
)

type secureConn struct {
	encoder cipher.Stream
	decoder cipher.Stream
	conn    net.Conn
}

func newSecureConn(key string, conn net.Conn, iv []byte) *secureConn {
	sc := new(secureConn)
	sc.conn = conn
	commkey := sha256.Sum256([]byte(key))

	// encoder
	block, err := aes.NewCipher(commkey[:])
	if err != nil {
		log.Println(err)
		return nil
	}
	sc.encoder = cipher.NewCFBEncrypter(block, iv)

	// decoder
	block, err = aes.NewCipher(commkey[:])
	if err != nil {
		log.Println(err)
		return nil
	}
	sc.decoder = cipher.NewCFBDecrypter(block, iv)
	return sc
}

func (sc *secureConn) Read(p []byte) (n int, err error) {
	n, err = sc.conn.Read(p)
	if err == nil {
		sc.decoder.XORKeyStream(p[:n], p[:n])
	}
	return
}

func (sc *secureConn) Write(p []byte) (n int, err error) {
	sc.encoder.XORKeyStream(p, p)
	return sc.conn.Write(p)
}

func (sc *secureConn) Close() (err error) {
	return sc.conn.Close()
}

// handle multiplex-ed connection
func handleMux(conn *suft.Conn, key, target string, tuncrypt bool) {
	// read iv
	iv := make([]byte, aes.BlockSize)
	if _, err := io.ReadFull(conn, iv); err != nil {
		log.Println(err)
		conn.Close()
		return
	}

	// stream multiplex
	var mux *yamux.Session
	if tuncrypt {
		scon := newSecureConn(key, conn, iv)
		m, err := yamux.Server(scon, nil)
		if err != nil {
			log.Println(err)
			return
		}
		mux = m
	} else {
		m, err := yamux.Server(conn, nil)
		if err != nil {
			log.Println(err)
			return
		}
		mux = m
	}
	defer mux.Close()

	for {
		p1, err := mux.Accept()
		if err != nil {
			log.Println(err)
			return
		}
		p2, err := net.Dial("tcp", target)
		if err != nil {
			log.Println(err)
			return
		}
		go handleClient(p1, p2)
	}
}

func handleClient(p1, p2 net.Conn) {
	log.Println("stream opened")
	defer log.Println("stream closed")
	defer p1.Close()
	defer p2.Close()

	// start tunnel
	p1die := make(chan struct{})
	go func() {
		io.Copy(p1, p2)
		close(p1die)
	}()

	p2die := make(chan struct{})
	go func() {
		io.Copy(p2, p1)
		close(p2die)
	}()

	// wait for tunnel termination
	select {
	case <-p1die:
	case <-p2die:
	}
}

func checkErr(e error) {
	if e != nil {
		log.Fatalln(e)
	}
}

func main() {
	rand.Seed(int64(time.Now().Nanosecond()))
	myApp := cli.NewApp()
	myApp.Name = "sufttun"
	myApp.Usage = "suft server"
	myApp.Version = "1.0"
	myApp.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "listen, l",
			Value: ":29900",
			Usage: "suft server listen addr:",
		},
		cli.StringFlag{
			Name:  "target, t",
			Value: "127.0.0.1:12948",
			Usage: "target server addr",
		},
		cli.StringFlag{
			Name:   "key",
			Value:  "it's a secrect",
			Usage:  "key for communcation, must be the same as sufttun client",
			EnvVar: "KCPTUN_KEY",
		},
		cli.IntFlag{
			Name:  "bandwidth, b",
			Value: 10,
			Usage: "your bandwidth",
		},
		cli.BoolFlag{
			Name:  "tuncrypt",
			Usage: "enable tunnel encryption, adds extra secrecy for data transfer",
		},
	}
	myApp.Action = func(c *cli.Context) {
		// var timeWaiting int64
		var raddr string
		var p suft.Params
		p.LocalAddr = c.String("listen")
		raddr = c.String("target")
		p.IsServ = true
		p.FastRetransmit = true
		p.FlatTraffic = true
		p.Bandwidth = int64(c.Int("bandwidth"))
		p.EnablePprof = false
		p.Stacktrace = false
		p.Debug = 0
		// timeWaiting = 0

		if raddr == "" {
			log.Println("missing raddr")
			return
		}

		if p.LocalAddr == "" {
			log.Println("missing localaddr")
			return
		}

		e, err := suft.NewEndpoint(&p)
		checkErr(err)
		defer e.Close()


		log.Println("listening on ", p.LocalAddr)
		log.Println("tunnel encryption:", c.Bool("tuncrypt"))
		var conn *suft.Conn
		for {
			conn = e.Listen()
			log.Println("connected from", conn.RemoteAddr())
			go handleMux(conn, c.String("key"), c.String("target"), c.Bool("tuncrypt"))
		}
	}
	myApp.Run(os.Args)
}
