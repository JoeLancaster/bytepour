package tracker

import (
	"math"
)

// URLBuffer is a working space
// for the url builder
// its size is sufficiently large for any
// sane url.
type URLBuffer [2048]byte

// AnnounceEvent is a three state
// type representing the state
// the downloader is in.
type AnnounceEvent byte

const (
	// Must be sent when we quit a download
	EventStopped = AnnounceEvent(0)
	// Must be sent when we first contact the tracker
	EventStarted = AnnounceEvent(1)
	// Must be sent when the download has completed
	// but not if the download is already completed when we contact
	// the tracker
	EventCompleted = AnnounceEvent(2)
)

// String implements the Stringer interface for AnnounceEvent
// String must inline to avoid heap allocs.
func (e AnnounceEvent) String() string {
	switch e {
	case EventStopped:
		return "stopped"
	case EventStarted:
		return "started"
	case EventCompleted:
		return "completed"
	default:
		return "stopped"
	}
}

// AnnounceRequest is a collection of query param
// items that will be constructed into a URL
// to request at the tracker's announce endpoint.
type AnnounceRequest struct {
	InfoHash   [20]byte
	PeerId     [20]byte
	Port       uint64
	Uploaded   uint64
	Downloaded uint64
	Left       uint64
	NumWant    uint64
	Event      AnnounceEvent
}

// Build adds required query params to buf, for the announce request
// to the tracker.
func Build(buf *URLBuffer, announce []byte, req *AnnounceRequest) []byte {
	const (
		param = byte('?')
		is    = byte('=')
		and   = byte('&')

		// query param keys
		info_hash  = "info_hash"
		peer_id    = "peer_id"
		port       = "port"
		uploaded   = "uploaded"
		downloaded = "downloaded"
		left       = "left"
		event      = "event"
		numwant    = "numwant"
	)

	var n int

	// length of the announce URL info hash, peer id, and all keys are known a priori
	// we also calculate the length of the integer fields after they're converted to strings
	// hence we know the length of the entire URL before we construct it
	// so we bounds check once and fail early.
	//
	// infohash and peerid lengths are 60 as this is the worst case
	// url encode.
	lengthUpperBound := len(announce) + 1 /*len(param)*/ +
		len(info_hash) + 1 /*len(is)*/ + 60 /*len(InfoHash)*/ + 1 /*len(and)*/ +
		len(peer_id) + 1 /*len(is)*/ + 60 /*len(PeerId)*/ + 1 /*len(and)*/ +
		len(downloaded) + 1 /*len(is)*/ + uintLen(req.Downloaded) + 1 /*len(and)*/ +
		len(uploaded) + 1 /*len(is)*/ + uintLen(req.Uploaded) + 1 /*len(and)*/ +
		len(port) + 1 /*len(is)*/ + uintLen(req.Port) + 1 /*len(and)*/ +
		len(left) + 1 /*len(is)*/ + uintLen(req.Left) + 1 /*len(and)*/ +
		len(numwant) + 1 /*len(is)*/ + uintLen(req.NumWant) + 1 /*len(and)*/ +
		len(event) + 1 /*len(is)*/ + len("completed")

	if lengthUpperBound >= len(buf) {
		return nil
	}

	// a.com:9000
	n += copy(buf[n:], announce)

	// a.com:9000?
	buf[n] = param
	n += 1

	// a.com:9000?info_hash
	n += copy(buf[n:], info_hash)

	// a.com:9000?info_hash=
	buf[n] = is
	n += 1

	// a.com:9000?info_hash=a93ef199cd398209802
	n += Escape20(buf[n:], &req.InfoHash)

	// a.com:9000?info_hash=a93ef199cd398209802&
	buf[n] = and
	n += 1

	// a.com:9000?info_hash=a93ef199cd398209802&peer_id
	// n.b. as of writing the ID we're using doesn't
	// actually need escaping.
	n += copy(buf[n:], peer_id)

	// a.com:9000?info_hash=a93ef199cd398209802&peer_id=
	buf[n] = is
	n += 1

	// a.com:9000?info_hash=a93ef199cd398209802&peer_id=BitTorrent1234567890
	n += Escape20(buf[n:], &req.PeerId)

	// a.com:9000?info_hash=a93ef199cd398209802&peer_id=BitTorrent1234567890&
	buf[n] = and
	n += 1

	// a.com:9000?info_hash=a93ef199cd398209802&peer_id=BitTorrent1234567890&downloaded
	n += copy(buf[n:], downloaded)

	// a.com:9000?info_hash=a93ef199cd398209802&peer_id=BitTorrent1234567890&downloaded=
	buf[n] = is
	n += 1

	// a.com:9000?info_hash=a93ef199cd398209802&peer_id=BitTorrent1234567890&downloaded=93
	n += putUint(buf[n:], req.Downloaded)

	// a.com:9000?info_hash=a93ef199cd398209802&peer_id=BitTorrent1234567890&downloaded=93&
	buf[n] = and
	n += 1

	// a.com:9000?info_hash=a93ef199cd398209802&peer_id=BitTorrent1234567890&downloaded=93&uploaded
	n += copy(buf[n:], uploaded)

	// a.com:9000?info_hash=a93ef199cd398209802&peer_id=BitTorrent1234567890&downloaded=93&uploaded=
	buf[n] = is
	n += 1

	// a.com:9000?info_hash=a93ef199cd398209802&peer_id=BitTorrent1234567890&downloaded=93&uploaded=23
	n += putUint(buf[n:], req.Uploaded)

	// a.com:9000?info_hash=a93ef199cd398209802&peer_id=BitTorrent1234567890&downloaded=93&uploaded=23&
	buf[n] = and
	n += 1

	// a.com:9000?info_hash=a93ef199cd398209802&peer_id=BitTorrent1234567890&downloaded=93&uploaded=23&numwant
	n += copy(buf[n:], numwant)

	// a.com:9000?info_hash=a93ef199cd398209802&peer_id=BitTorrent1234567890&downloaded=93&uploaded=23&numwant=
	buf[n] = is
	n += 1

	// a.com:9000?info_hash=a93ef199cd398209802&peer_id=BitTorrent1234567890&downloaded=93&uploaded=23&numwant=50
	n += putUint(buf[n:], req.NumWant)

	// a.com:9000?info_hash=a93ef199cd398209802&peer_id=BitTorrent1234567890&downloaded=93&uploaded=23&numwant=50&
	buf[n] = and
	n += 1

	// a.com:9000?info_hash=a93ef199cd398209802&peer_id=BitTorrent1234567890&downloaded=93&uploaded=23&numwant=50&port
	n += copy(buf[n:], port)

	// a.com:9000?info_hash=a93ef199cd398209802&peer_id=BitTorrent1234567890&downloaded=93&uploaded=23&numwant=50&port=
	buf[n] = is
	n += 1

	// a.com:9000?info_hash=a93ef199cd398209802&peer_id=BitTorrent1234567890&downloaded=93&\
	// uploaded=23&numwant=50&port=6007
	n += putUint(buf[n:], req.Port)

	// a.com:9000?info_hash=a93ef199cd398209802&peer_id=BitTorrent1234567890&downloaded=93&\
	// uploaded=23&numwant=50&port=6007&
	buf[n] = and
	n += 1

	// a.com:9000?info_hash=a93ef199cd398209802&peer_id=BitTorrent1234567890&downloaded=93&\
	// uploaded=23&numwant=50&port=6007&left
	n += copy(buf[n:], left)

	// a.com:9000?info_hash=a93ef199cd398209802&peer_id=BitTorrent1234567890&downloaded=93&\
	// uploaded=23&numwant=50&port=6007&left=
	buf[n] = is
	n += 1

	// a.com:9000?info_hash=a93ef199cd398209802&peer_id=BitTorrent1234567890&downloaded=93&\
	// uploaded=23&numwant=50&port=6007&left=23904
	n += putUint(buf[n:], req.Left)

	// a.com:9000?info_hash=a93ef199cd398209802&peer_id=BitTorrent1234567890&downloaded=93&\
	// uploaded=23&numwant=50&port=6007&left=23904&
	buf[n] = and
	n += 1

	// a.com:9000?info_hash=a93ef199cd398209802&peer_id=BitTorrent1234567890&downloaded=93&\
	// uploaded=23&numwant=50&port=6007&left=23904&event
	n += copy(buf[n:], []byte(event))

	// a.com:9000?info_hash=a93ef199cd398209802&peer_id=BitTorrent1234567890&downloaded=93&\
	// uploaded=23&numwant=50&port=6007&left=23904&event=
	buf[n] = is
	n += 1

	// a.com:9000?info_hash=a93ef199cd398209802&peer_id=BitTorrent1234567890&downloaded=93&\
	// uploaded=23&numwant=50&port=6007&left=23904&event=completed
	n += copy(buf[n:], []byte(req.Event.String()))

	return buf[:n]
}

// uintLen gives the length of 'x' if it was
// converted to a base-10 string.
func uintLen(x uint64) int {
	const lim = uint64(math.MaxUint64 / 10)

	n := 1

	var powerOfTen uint64 = 10

	for x >= powerOfTen {
		n++

		if powerOfTen > lim {
			break
		}

		powerOfTen *= 10
	}

	return n
}

// putUint copies the string representation of x into
// the slice at dst.
func putUint(dst []byte, x uint64) int {
	const smallsString = "00010203040506070809" +
		"10111213141516171819" +
		"20212223242526272829" +
		"30313233343536373839" +
		"40414243444546474849" +
		"50515253545556575859" +
		"60616263646566676869" +
		"70717273747576777879" +
		"80818283848586878889" +
		"90919293949596979899"

	var a [65]byte
	i := len(a)

	us := x
	for us >= 100 {
		is := us % 100 * 2
		us /= 100
		i -= 2
		a[i+1] = smallsString[is+1]
		a[i+0] = smallsString[is+0]
	}

	is := us * 2
	i--
	a[i] = smallsString[is+1]
	if us >= 10 {
		i--
		a[i] = smallsString[is]
	}

	return copy(dst, a[i:])
}
