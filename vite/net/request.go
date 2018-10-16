package net

import (
	"errors"
	"fmt"
	"github.com/vitelabs/go-vite/ledger"
	"github.com/vitelabs/go-vite/p2p"
	"github.com/vitelabs/go-vite/vite/net/message"
	"sort"
	"strconv"
	"strings"
	"time"
)

type reqState byte

const (
	reqWaiting reqState = iota
	reqPending
	reqRespond
	reqDone
	reqError
	reqCancel
)

var reqStatus = [...]string{
	reqWaiting: "waiting",
	reqPending: "pending",
	reqRespond: "respond",
	reqDone:    "done",
	reqError:   "error",
	reqCancel:  "canceled",
}

func (s reqState) String() string {
	return reqStatus[s]
}

type context interface {
	Add(r Request)
	Retry(id uint64, err error)
	FC() *fileClient
	Get(id uint64) (Request, bool)
}

type Request interface {
	Handle(ctx context, msg *p2p.Msg, peer *Peer)
	ID() uint64
	SetID(id uint64)
	Run(ctx context)
	Done()
	Catch(err error)
	Expired() bool
	State() reqState
	Req() Request
	Band() (from, to uint64)
	SetBand(from, to uint64)
	SetPeer(peer *Peer)
	Peer() *Peer
}

type reqRec interface {
	receiveSnapshotBlock(block *ledger.SnapshotBlock)
	receiveAccountBlock(block *ledger.AccountBlock)
}
type errCallback = func(id uint64, err error)

var errMissingPeer = errors.New("request missing peer")
var errUnExpectedRes = errors.New("unexpected response")

const minBlocks uint64 = 3600  // minimal snapshot blocks per subLedger request
const maxBlocks uint64 = 10800 // maximal snapshot blocks per subLedger request
const maxRetry = 3

type subLedgerPiece struct {
	from, to uint64
	peer     *Peer
}

// split large subledger request to many small pieces
func splitSubLedger(from, to uint64, peers Peers) (cs []*subLedgerPiece) {
	// sort peers from low to high
	sort.Sort(peers)

	total := to - from + 1
	if total < minBlocks {
		cs = append(cs, &subLedgerPiece{
			from: from,
			to:   to,
			peer: peers[len(peers)-1], // choose the tallest peer
		})
		return
	}

	var pCount, pTo uint64 // piece length
	for _, peer := range peers {
		if peer.height > from+minBlocks {
			pTo = peer.height - minBlocks
			pCount = pTo - from + 1

			// peer not high enough
			if pCount < minBlocks {
				continue
			}

			// piece too large
			if pCount > maxBlocks {
				pCount = maxBlocks
			}

			// piece to
			pTo = from + pCount - 1

			// piece to exceed target height
			if pTo > to {
				pTo = to
				pCount = to - from + 1
			}
			// reset piece is too small, then collapse to one piece
			if to < pTo+minBlocks {
				pCount += to - pTo
			}

			cs = append(cs, &subLedgerPiece{
				from: from,
				to:   pTo,
				peer: peer,
			})

			from = pTo + 1
			if from > to {
				break
			}
		}
	}

	// reset piece, alloc to best peer
	if from < to {
		// if peer is not too much, has rest chunk, then collapse the rest chunk with last chunk
		lastChunk := cs[len(cs)-1]
		if lastChunk.peer == peers[len(peers)-1] {
			lastChunk.to = to
		} else {
			cs = append(cs, &subLedgerPiece{
				from: from,
				to:   to,
				peer: peers[len(peers)-1],
			})
		}
	}

	return
}

func splitChunk(from, to uint64) (chunks [][2]uint64) {
	// chunks may be only one block, then from == to
	if from > to || to == 0 {
		return
	}

	total := (to-from)/minBlocks + 1
	chunks = make([][2]uint64, total)

	var cTo uint64
	var i int
	for from <= to {
		if to > from+minBlocks {
			cTo = from + minBlocks - 1
		} else {
			cTo = to
		}

		chunks[i] = [2]uint64{from, cTo}

		from = cTo + 1
		i++
	}

	return chunks
}

// @request for subLedger, will get FileList and Chunk
type subLedgerRequest struct {
	id         uint64 // id & child_id
	from, to   uint64
	peer       *Peer
	state      reqState
	expiration time.Time
	catch      errCallback
	rec        reqRec
}

func (s *subLedgerRequest) SetID(id uint64) {
	s.id = id
}

func (s *subLedgerRequest) Band() (from, to uint64) {
	return s.from, s.to
}

func (s *subLedgerRequest) SetBand(from, to uint64) {
	s.from = from
	s.to = to
}

func (s *subLedgerRequest) SetPeer(peer *Peer) {
	s.peer = peer
}

func (s *subLedgerRequest) Peer() *Peer {
	return s.peer
}

func (s *subLedgerRequest) State() reqState {
	return s.state
}

func (s *subLedgerRequest) Handle(ctx context, pkt *p2p.Msg, peer *Peer) {
	defer staticDuration("handle-filelist", time.Now())

	if cmd(pkt.Cmd) == FileListCode {
		s.state = reqRespond

		msg := new(message.FileList)
		err := msg.Deserialize(pkt.Payload)
		if err != nil {
			ctx.Retry(s.id, err)
			return
		}

		if len(msg.Files) != 0 {
			// sort as StartHeight
			sort.Sort(files(msg.Files))

			from := msg.Files[0].StartHeight
			to := msg.Files[len(msg.Files)-1].EndHeight
			if from < s.from {
				from = s.from
			}
			if to > s.to {
				to = s.to
			}

			var blockNumber uint64
			for _, file := range msg.Files {
				blockNumber += file.BlockNumbers
			}

			// request files
			ctx.FC().request(&fileRequest{
				from:    from,
				to:      to,
				files:   msg.Files,
				nonce:   msg.Nonce,
				peer:    peer,
				current: from,
				rec:     s.rec,
				ctx:     ctx,
				catch:   s.catch,
			})
		}

		// request chunks
		for _, chunk := range msg.Chunks {
			// maybe chunk is too large
			cs := splitChunk(chunk[0], chunk[1])
			for _, c := range cs {
				ctx.Add(&chunkRequest{
					from: c[0],
					to:   c[1],
					peer: peer,
					rec:  s.rec,
				})
			}
		}

		s.Done()
	} else {
		ctx.Retry(s.id, errUnExpectedRes)
	}
}

func (s *subLedgerRequest) ID() uint64 {
	return s.id
}

func (s *subLedgerRequest) Run(ctx context) {
	s.expiration = time.Now().Add(subledgerTimeout)

	err := s.peer.Send(GetSubLedgerCode, s.id, &message.GetSubLedger{
		From:    ledger.HashHeight{Height: s.from},
		Count:   s.to - s.from + 1,
		Forward: true,
	})

	if err != nil {
		s.peer.log.Error(fmt.Sprintf("send GetSubLedgerMsg<%d-%d> to %s error: %v", s.from, s.to, s.peer, err))
		ctx.Retry(s.id, err)
	} else {
		s.state = reqPending
		s.peer.log.Info(fmt.Sprintf("send GetSubLedgerMsg<%d-%d> to %s done", s.from, s.to, s.peer))
	}
}

func (s *subLedgerRequest) Done() {
	s.state = reqDone
}

func (s *subLedgerRequest) Catch(err error) {
	s.state = reqError
	s.catch(s.id, err)
}

func (s *subLedgerRequest) Expired() bool {
	return time.Now().After(s.expiration)
}

func (s *subLedgerRequest) Req() Request {
	return s
}

// @request file
type fileRequest struct {
	from, to uint64
	files    []*ledger.CompressedFileMeta
	nonce    uint64
	peer     *Peer
	current  uint64 // the tallest snapshotBlock have received, as the breakpoint resume
	rec      reqRec
	ctx      context     // for create new GetSubLedgerMsg, retry
	catch    errCallback // for retry GetSubLedger, catch error
}

// file1/0-3600 file2/3601-7200
func (r *fileRequest) String() string {
	files := make([]string, len(r.files))
	for i, file := range r.files {
		files[i] = file.Filename + "/" +
			strconv.FormatUint(file.StartHeight, 10) +
			"-" +
			strconv.FormatUint(file.EndHeight, 10)
	}

	return strings.Join(files, " ")
}

func (r *fileRequest) Catch(err error) {
	from := r.current
	if from == 0 {
		from = r.from
	}

	req := &subLedgerRequest{
		from:  from,
		to:    r.to,
		catch: r.catch,
		rec:   r.rec,
	}

	r.ctx.Add(req)
}

func (r *fileRequest) Addr() string {
	return r.peer.FileAddress().String()
}

// @request for chunk
type chunkRequest struct {
	id         uint64
	from, to   uint64
	peer       *Peer
	state      reqState
	expiration time.Time
	catch      errCallback
	rec        reqRec
}

func (c *chunkRequest) SetID(id uint64) {
	c.id = id
}

func (c *chunkRequest) Band() (from, to uint64) {
	return c.from, c.to
}

func (c *chunkRequest) SetBand(from, to uint64) {
	c.from = from
	c.to = to
}

func (c *chunkRequest) SetPeer(peer *Peer) {
	c.peer = peer
}

func (c *chunkRequest) Peer() *Peer {
	return c.peer
}

func (c *chunkRequest) Req() Request {
	return c
}

func (c *chunkRequest) State() reqState {
	return c.state
}

func (c *chunkRequest) Handle(ctx context, pkt *p2p.Msg, peer *Peer) {
	defer staticDuration("handle-chunk", time.Now())

	if cmd(pkt.Cmd) == SubLedgerCode {
		c.state = reqRespond

		msg := new(message.SubLedger)
		err := msg.Deserialize(pkt.Payload)
		if err != nil {
			fmt.Println("chunkRequest handle error: ", err)
			ctx.Retry(c.id, err)
			return
		}

		for _, block := range msg.SBlocks {
			c.rec.receiveSnapshotBlock(block)
		}
		for _, block := range msg.ABlocks {
			c.rec.receiveAccountBlock(block)
		}

		c.Done()
	} else {
		ctx.Retry(c.id, errUnExpectedRes)
	}
}

func (c *chunkRequest) ID() uint64 {
	return c.id
}

func (c *chunkRequest) Run(ctx context) {
	c.state = reqWaiting

	err := c.peer.Send(GetChunkCode, c.id, &message.GetChunk{
		Start: c.from,
		End:   c.to,
	})

	if err != nil {
		c.peer.log.Error(fmt.Sprintf("send GetChunkMsg<%d/%d> to %s error: %v", c.from, c.to, c.peer, err))
		ctx.Retry(c.id, err)
	} else {
		c.state = reqPending
		c.peer.log.Info(fmt.Sprintf("send GetChunkMsg<%d/%d> to %s done", c.from, c.to, c.peer))
	}
}

func (c *chunkRequest) Done() {
	c.state = reqDone
}

func (c *chunkRequest) Catch(err error) {
	c.state = reqError
	c.catch(c.id, err)
}

func (c *chunkRequest) Expired() bool {
	return time.Now().After(c.expiration)
}

// helper
type files []*ledger.CompressedFileMeta

func (a files) Len() int {
	return len(a)
}

func (a files) Less(i, j int) bool {
	return a[i].StartHeight < a[j].StartHeight
}

func (a files) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}

// helper
func u64ToDuration(n uint64) time.Duration {
	return time.Duration(int64(n/100) * int64(time.Second))
}
