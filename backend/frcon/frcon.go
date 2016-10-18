package frcon

import (
	"bytes"
	"crypto/md5"
	"encoding/binary"
	"encoding/hex"
	"errors"
	"io"
	"net"
	"strings"
	"sync"
	"time"

	// log "github.com/Sirupsen/logrus"
	"github.com/djherbis/buffer"
)

// TODO: implement proper logging, with levels
// TODO: add timeout / sleep customization
// TODO: make auto-login and auto-events-enabling opt-outable
// TODO: document the API
// TODO: move as much code as possible into pure functions, so it can be tested
// TODO: implement a mock game server
// TODO: test for memory leaks and race conditions, including after closing

const (
	retrySleep      = 5 * time.Second
	dialTimeout     = 5 * time.Second
	readTimeout     = 5 * time.Second
	writeTimeout    = 5 * time.Second
	requestTimeout  = 10 * time.Second
	callbackTimeout = time.Second
	emitTimeout     = 30 * time.Second
	maxSize         = 0x4000     // 16384
	isClientBit     = 0x80000000 // 0 = server, 1 = client
	isResponseBit   = 0x40000000 // 0 = request, 1 = response
	seqNumMask      = 0x3FFFFFFF // this is also the max, thanks to lucky alignment
)

var (
	ErrBadSequence    = errors.New("bad sequence number")
	ErrDone           = errors.New("done")
	ErrNoWords        = errors.New("I have no words for this")
	ErrPayloadTooLong = errors.New("payload too long")
	ErrRequestTimeout = errors.New("request timeout")

	errTimeout = errors.New("timeout")
)

type timeoutConn struct {
	netconn net.Conn
}

func (conn timeoutConn) Read(b []byte) (int, error) {
	conn.netconn.SetReadDeadline(time.Now().Add(readTimeout))
	n, err := conn.netconn.Read(b)
	if neterr, ok := err.(net.Error); ok && neterr.Timeout() {
		return n, errTimeout
	}

	return n, err
}

func (conn timeoutConn) Write(b []byte) (int, error) {
	conn.netconn.SetWriteDeadline(time.Now().Add(writeTimeout))
	n, err := conn.netconn.Write(b)
	if neterr, ok := err.(net.Error); ok && neterr.Timeout() {
		return n, errTimeout
	}

	return n, err
}

func (conn timeoutConn) Close() error {
	return conn.netconn.Close()
}

type EventType int

const (
	EWords        EventType = iota
	EConnected    EventType = iota
	ELoggedIn     EventType = iota
	EDisconnected EventType = iota
	EBadPassword  EventType = iota
	EBadRead      EventType = iota
)

type Event struct {
	Type      EventType
	Timestamp time.Time
	Words     []string
}

type Session struct {
	host       string
	password   string
	conn       chan timeoutConn
	writeMutex sync.Mutex
	badconn    chan timeoutConn
	loggedin   chan struct{}
	relogin    chan struct{}
	done       chan struct{}
	seqNumIter int
	seqMutex   sync.Mutex
	callbacks  map[int]chan []string
	cbMutex    sync.Mutex
	events     chan Event
}

func Dial(host, password string) (*Session, chan Event) {
	s := &Session{
		host:      host,
		password:  password,
		conn:      make(chan timeoutConn),
		badconn:   make(chan timeoutConn),
		loggedin:  make(chan struct{}),
		relogin:   make(chan struct{}),
		done:      make(chan struct{}),
		callbacks: make(map[int]chan []string),
		events:    make(chan Event),
	}
	go s.maintainConnection()
	go s.stayLoggedIn()
	go s.keepReading()
	return s, s.events
}

func (s *Session) emit(t EventType, w []string) {
	// In case you're unsure, the Event struct initialization below is evaluated
	// before select is blocked, not after the channel is ready to consume it, so
	// the timestamp doesn't get offset or anything.
	//
	//  https://play.golang.org/p/mdJ42XdAKz
	select {
	case s.events <- Event{t, time.Now(), w}:
	case <-time.After(emitTimeout):
	case <-s.done:
	}
}

func (s *Session) maintainConnection() {
	var conn timeoutConn

	connect := func() {
		isRetry := false
		for {
			select {
			case <-s.done:
				return
			default:
			}

			if isRetry {
				select {
				case <-time.After(retrySleep):
				case <-s.done:
					return
				}
			} else {
				isRetry = true
			}

			netconn, err := net.DialTimeout("tcp", s.host, dialTimeout)
			if err != nil {
				continue
			}

			go s.emit(EConnected, nil)
			conn = timeoutConn{netconn}
			break
		}
	}

	connect()
	for {
		select {
		case s.conn <- conn:
		case badconn := <-s.badconn:
			if badconn == conn {
				go s.emit(EDisconnected, nil)
				conn.Close()
				// s.cbMutex.Lock()
				// for _, c := range s.callbacks {
				// 	close(c)
				// } // TODO: check sanity
				// s.cbMutex.Unlock()
				connect()

				// Could check if there's a login procedure in progress, but let's keep
				// this simple, for a safe measure.
				go func() { s.relogin <- struct{}{} }()
			}
		case <-s.done:
			conn.Close()
			return
		}
	}
}

func (s *Session) stayLoggedIn() {
	login := func() {
		for {
			select {
			case <-s.done:
				return
			default:
			}

			// TODO: consider sleeping on retry, like in connect()

			resp, err := s.request([]string{"login.hashed"}, false)
			if err != nil {
				continue
			} else if len(resp) < 2 || resp[0] != "OK" {
				continue
			}

			challenge, err := hex.DecodeString(resp[1])
			if err != nil {
				// TODO: report connection?
				continue
			}
			hash := md5.Sum(append(challenge, []byte(s.password)...))
			resp, err = s.request([]string{
				"login.hashed", strings.ToUpper(hex.EncodeToString(hash[:])),
			}, false)
			if err != nil {
				continue
			} else if len(resp) != 1 || resp[0] != "OK" {
				go s.emit(EBadPassword, nil)
				continue
			}

			resp, err = s.request([]string{"admin.eventsEnabled", "true"}, false)
			if err != nil {
				continue
			} else if len(resp) != 1 || resp[0] != "OK" {
				continue
			}

			go s.emit(ELoggedIn, nil)
			break
		}
	}

	login()
	for {
		select {
		case s.loggedin <- struct{}{}:
		case <-s.relogin:
			login()
		case <-s.done:
			return
		}
	}
}

const headerSize = 12

type header struct {
	Sequence uint32
	Size     int32
	NumWords int32
}

func (s *Session) getCallback(seqNum int) (c chan []string, ok bool) {
	s.cbMutex.Lock()
	defer s.cbMutex.Unlock()
	c, ok = s.callbacks[seqNum]
	return
}

func (s *Session) keepReading() {
	var conn timeoutConn
	var b [maxSize]byte

	buf := buffer.NewRing(buffer.New(maxSize * 2))
	queued := int32(0)
	var h *header

	reset := func(_ string) {
		go s.emit(EBadRead, nil)

		select {
		case s.badconn <- conn:
		case <-s.done:
			return
		}

		buf.Reset()
		queued = int32(0)
		h = nil
	}

OUTER:
	for {
		select {
		case conn = <-s.conn:
		case <-s.done:
			return
		}

		n, err := conn.Read(b[:]) // we read only here, no mutex needed
		// TODO: if n is zero, then the connection is closed
		if err == errTimeout {
			continue
		} else if err == io.EOF {
			reset("EOF")
			continue
		} else if err != nil {
			reset("bad read")
			continue
		}

		queued += int32(n)
		_, err = buf.Write(b[:n])
		if err != nil {
			panic(err)
		}

		for queued >= headerSize {
			if h == nil {
				h = new(header)
				binary.Read(buf, binary.LittleEndian, h)
				if h.Size > maxSize {
					reset("h.Size > maxSize")
					continue OUTER
				}

				queued -= headerSize
			}

			if queued < h.Size-headerSize {
				continue OUTER
			}

			targetQueued := queued + headerSize - h.Size
			var wordSize int32
			words := make([]string, 0, h.NumWords)
			for i := int32(0); i < h.NumWords; i++ {
				binary.Read(buf, binary.LittleEndian, &wordSize)
				if wordSize+5 > queued {
					reset("wordSize+5 > queued")
					continue OUTER
				}
				buf.Read(b[:wordSize+1]) // +1 is there to accommodate NULL
				words = append(words, string(b[:wordSize]))
				queued -= wordSize + 5
			}

			if queued != targetQueued {
				reset("queued != targetQueued")
				continue OUTER
			}

			seqNum := int(h.Sequence & seqNumMask)
			if h.Sequence&isResponseBit == 0 {
				go s.emit(EWords, words)
				// send's errors are irrelevant in this case
				s.send(seqNum, []string{"OK"}, false, false)
			} else {
				if c, ok := s.getCallback(seqNum); ok {
					go func() {
						select {
						case c <- words: // words slice is goroutine-safe
						case <-time.After(callbackTimeout):
							// TODO: print an error or something
						case <-s.done:
						}
					}()
				}
			}

			h = nil
		}
	}
}

func (s *Session) request(words []string, waitLogin bool) ([]string, error) {
	s.seqMutex.Lock()
	seqNum := s.seqNumIter
	s.seqNumIter++
	if s.seqNumIter > seqNumMask {
		s.seqNumIter = 0
	}
	s.seqMutex.Unlock()

	c := make(chan []string)
	s.cbMutex.Lock()
	s.callbacks[seqNum] = c
	s.cbMutex.Unlock()
	defer func() {
		s.cbMutex.Lock()
		delete(s.callbacks, seqNum)
		s.cbMutex.Unlock()
	}()

	// TODO: if the connection keeps failing, or the password is incorrect / the
	// game server is misbehaving, then send can block for ages, hence, we
	// probably need a timeout
	err := s.send(seqNum, words, true, waitLogin)
	if err != nil {
		return nil, err
	}
	select {
	case resp := <-c:
		return resp, nil
	case <-time.After(requestTimeout):
		return nil, ErrRequestTimeout
	case <-s.done:
		return nil, ErrDone
	}
}

func (s *Session) send(
	seqNum int, words []string, isRequest bool, waitLogin bool,
) error {
	if seqNum < 0 || seqNum > seqNumMask {
		return ErrBadSequence
	} else if len(words) == 0 {
		return ErrNoWords
	} // TODO: not sure if matters

	h := header{
		Sequence: uint32(seqNum),
		Size:     headerSize,
		NumWords: int32(len(words)),
	}
	for _, word := range words {
		h.Size += 5 + int32(len(word)) // 5 = len(word) + NULL terminator
	}
	if h.Size > maxSize {
		return ErrPayloadTooLong
	}

	buf := bytes.NewBuffer(make([]byte, 0, h.Size))
	if isRequest {
		h.Sequence |= isClientBit
	} else {
		h.Sequence |= isResponseBit
	}

	binary.Write(buf, binary.LittleEndian, &h)
	for _, word := range words {
		binary.Write(buf, binary.LittleEndian, int32(len(word)))
		buf.WriteString(word)
		buf.WriteByte(0)
	}

	var conn timeoutConn
	for {
		// Login waiting comes before conn because conn will resolve instantly if
		// everything is fine, while login may block for a while.
		if waitLogin {
			select {
			case <-s.loggedin:
			case <-s.done:
				return ErrDone
			}
		}

		select {
		case conn = <-s.conn:
		case <-s.done:
			return ErrDone
		}

		s.writeMutex.Lock()
		_, err := conn.Write(buf.Bytes())
		s.writeMutex.Unlock()
		if err != nil {
			select {
			case s.badconn <- conn:
			case <-s.done:
				return ErrDone
			}
			continue
		}

		break
	}

	return nil
}

func (s *Session) Request(words []string) ([]string, error) {
	return s.request(words, true)
}

// don't wait for the login procedure to complete
func (s *Session) RequestPublic(words []string) ([]string, error) {
	return s.request(words, false)
}

func (s *Session) Host() string {
	return s.host
}

func (s *Session) Close() {
	close(s.done)
}
