package main

import (
	"crypto/aes"
	"crypto/cipher"
	crand "crypto/rand"
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

func checkError(err error) {
	if err != nil {
		log.Println(err)
		os.Exit(-1)
	}
}

func main() {
	rand.Seed(int64(time.Now().Nanosecond()))
	myApp := cli.NewApp()
	myApp.Name = "sufttun"
	myApp.Usage = "sufttun client"
	myApp.Version = "1.0"
	myApp.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "localaddr, l",
			Value: ":12948",
			Usage: "local listen addr:",
		},
		cli.StringFlag{
			Name:  "remoteaddr, r",
			Value: "vps:29900",
			Usage: "suft server addr",
		},
		cli.StringFlag{
			Name:   "key",
			Value:  "it's a secrect",
			Usage:  "key for communcation, must be the same as kcptun server",
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
		addr, err := net.ResolveTCPAddr("tcp", c.String("localaddr"))
		checkError(err)
		listener, err := net.ListenTCP("tcp", addr)
		checkError(err)
		log.Println("listening on:", listener.Addr())

	START_SUFT:
		var raddr string
		var p suft.Params
		p.LocalAddr = c.String("localaddr")
		raddr = c.String("remoteaddr")
		p.IsServ = false
		p.FastRetransmit = true
		p.FlatTraffic = true
		p.Bandwidth = int64(c.Int("bandwidth"))
		p.EnablePprof = false
		p.Stacktrace = false
		p.Debug = 0


		e, err := suft.NewEndpoint(&p)
		checkError(err)
		defer e.Close()

		suftconn, err := e.Dial(raddr)
		checkError(err)
		log.Println("connected to", suftconn.RemoteAddr())
		// generate & send iv
		iv := make([]byte, aes.BlockSize)
		io.ReadFull(crand.Reader, iv)
		_, err = suftconn.Write(iv)
		checkError(err)

		// stream multiplex
		var mux *yamux.Session
		if c.Bool("tuncrypt") {
			scon := newSecureConn(c.String("key"), suftconn, iv)
			session, err := yamux.Client(scon, nil)
			checkError(err)
			mux = session
		} else {
			session, err := yamux.Client(suftconn, nil)
			checkError(err)
			mux = session
		}
		log.Println("tunnel encryption:", c.Bool("tuncrypt"))

		for {
			p1, err := listener.AcceptTCP()
			if err != nil {
				log.Println(err)
				continue
			}
			p2, err := mux.Open()
			if err != nil { // yamux failure
				log.Println(err)
				suftconn.Close()
				p1.Close()
				goto START_SUFT
			}
			go handleClient(p1, p2)
		}
	}
	myApp.Run(os.Args)
}
