// Package p2p implements the vite P2P network

package p2p

import (
	"errors"
	"fmt"
	"github.com/vitelabs/go-vite/crypto/ed25519"
	"github.com/vitelabs/go-vite/log15"
	"github.com/vitelabs/go-vite/monitor"
	"github.com/vitelabs/go-vite/p2p/block"
	"github.com/vitelabs/go-vite/p2p/discovery"
	"github.com/vitelabs/go-vite/p2p/nat"
	"net"
	"strconv"
	"sync"
	"sync/atomic"
	"time"
)

var errSvrStarted = errors.New("server has started")

type Discovery interface {
	Start() error
	Stop()
	SubNodes(ch chan<- *discovery.Node)
	UnSubNodes(ch chan<- *discovery.Node)
	Mark(id discovery.NodeID, lifetime int64)
	Block(id discovery.NodeID, ip net.IP)
	Need(n uint)
}

type Config struct {
	Name            string
	NetID           NetworkID          // which network server runs on
	MaxPeers        uint               // max peers can be connected
	MaxPendingPeers uint               // max peers waiting for connect
	MaxInboundRatio uint               // max inbound peers: MaxPeers / MaxInboundRatio
	Port            uint               // TCP and UDP listen port
	DataDir         string             // the directory for storing node table, default is "~/viteisbest/p2p"
	PrivateKey      ed25519.PrivateKey // use for encrypt message, the corresponding public key use for NodeID
	Protocols       []*Protocol        // protocols server supported
	BootNodes       []string           // nodes as discovery seed
	StaticNodes     []string           // nodes to connect
}

type Server struct {
	*Config
	addr        *net.TCPAddr
	StaticNodes []*discovery.Node

	running   int32          // atomic
	wg        sync.WaitGroup // Wait for all jobs done
	term      chan struct{}
	pending   chan struct{} // how many connection can wait for handshake
	addPeer   chan *conn
	delPeer   chan *Peer
	discv     Discovery
	handshake *Handshake
	peers     *PeerSet
	blockList *block.Set
	self      *discovery.Node
	ln        net.Listener
	nodeChan  chan *discovery.Node // sub discovery nodes
	log       log15.Logger

	dialer *net.Dialer
}

func New(cfg *Config) (svr *Server, err error) {
	cfg = EnsureConfig(cfg)

	addr := "0.0.0.0:" + strconv.FormatUint(uint64(cfg.Port), 10)
	// udp discover
	udpAddr, err := net.ResolveUDPAddr("udp", addr)
	if err != nil {
		return
	}

	// tcp listener
	tcpAddr, err := net.ResolveTCPAddr("tcp", addr)
	if err != nil {
		return
	}

	ID, err := discovery.Priv2NodeID(cfg.PrivateKey)
	if err != nil {
		return
	}

	node := &discovery.Node{
		ID:  ID,
		IP:  udpAddr.IP,
		UDP: uint16(udpAddr.Port),
		TCP: uint16(tcpAddr.Port),
	}

	svr = &Server{
		Config:      cfg,
		addr:        tcpAddr,
		StaticNodes: parseNodes(cfg.StaticNodes),
		peers:       NewPeerSet(),
		pending:     make(chan struct{}, cfg.MaxPendingPeers),
		addPeer:     make(chan *conn, 1),
		delPeer:     make(chan *Peer, 1),
		blockList:   block.New(100),
		self:        node,
		nodeChan:    make(chan *discovery.Node, 10),
		log:         log15.New("module", "p2p/server"),
		dialer:      &net.Dialer{Timeout: 3 * time.Second},
	}

	svr.discv = discovery.New(&discovery.Config{
		Priv:      cfg.PrivateKey,
		DBPath:    cfg.DataDir,
		BootNodes: parseNodes(cfg.BootNodes),
		Addr:      udpAddr,
		Self:      node,
	})

	return
}

func (svr *Server) Start() error {
	if !atomic.CompareAndSwapInt32(&svr.running, 0, 1) {
		return errSvrStarted
	}

	svr.term = make(chan struct{})

	// setHandshake in method Start, because svr.Protocols may be modified
	svr.setHandshake()

	listener, err := net.ListenTCP("tcp", svr.addr)
	if err != nil {
		return err
	}
	svr.log.Info(fmt.Sprintf("tcp listen at %s", svr.addr))
	svr.ln = listener

	svr.term = make(chan struct{})

	svr.log.Info("p2p server start")

	// mapping udp
	svr.wg.Add(1)
	go func() {
		nat.Map(svr.term, "udp", int(svr.self.UDP), int(svr.self.UDP), "vite p2p udp", 0, svr.updateNode)
		svr.wg.Done()
	}()

	// mapping tcp
	svr.wg.Add(1)
	go func() {
		nat.Map(svr.term, "tcp", int(svr.self.TCP), int(svr.self.TCP), "vite p2p tcp", 0, svr.updateNode)
		svr.wg.Done()
	}()

	// subscribe nodes
	svr.discv.SubNodes(svr.nodeChan)

	err = svr.discv.Start()
	if err != nil {
		svr.ln.Close()
		return err
	}

	svr.wg.Add(1)
	go svr.dialLoop()

	// tcp listener
	svr.wg.Add(1)
	go svr.listenLoop()

	// peer manager
	svr.wg.Add(1)
	go svr.loop()

	svr.log.Info("p2p server started")
	return nil
}

func (svr *Server) Stop() {
	if svr.term == nil {
		return
	}

	select {
	case <-svr.term:
	default:
		svr.log.Warn("p2p server stop")

		close(svr.term)
		svr.discv.Stop()
		svr.discv.UnSubNodes(svr.nodeChan)
		svr.wg.Wait()

		svr.log.Warn("p2p server stopped")
	}
}

func (svr *Server) updateNode(addr *nat.Addr) {
	if addr.Proto == "tcp" {
		svr.self.TCP = uint16(addr.Port)
	} else {
		svr.self.UDP = uint16(addr.Port)
	}
}

func (svr *Server) setHandshake() {
	cmdsets := make([]*CmdSet, len(svr.Protocols))
	for i, p := range svr.Protocols {
		cmdsets[i] = p.CmdSet()
	}

	svr.handshake = &Handshake{
		Version: Version,
		Name:    svr.Name,
		NetID:   svr.NetID,
		ID:      svr.self.ID,
		CmdSets: cmdsets,
	}
}

func (svr *Server) dialLoop() {
	defer svr.wg.Done()

	// connect to static node first
	for _, node := range svr.StaticNodes {
		svr.dial(node.ID, node.TCPAddr(), static)
	}

	for {
		select {
		case <-svr.term:
			return
		case node := <-svr.nodeChan:
			svr.dial(node.ID, node.TCPAddr(), outbound)
		}
	}
}

// when peer is disconnected, maybe we want to reconnect it.
// we can get ID and addr only from peer, but not Node
// so dial(id, addr, flag) not dial(Node, flag)
func (svr *Server) dial(id discovery.NodeID, addr *net.TCPAddr, flag connFlag) {
	// has been blocked
	if svr.blockList.Has(id[:]) {
		return
	}

	if err := svr.checkConn(id, flag); err != nil {
		return
	}

	svr.pending <- struct{}{}
	if conn, err := svr.dialer.Dial("tcp", addr.String()); err == nil {
		go svr.setupConn(conn, flag)
	} else {
		<-svr.pending
		svr.log.Error(fmt.Sprintf("dial node %s@%s failed: %v", id, addr, err))
		svr.blockList.Add(id[:])
	}
}

func (svr *Server) listenLoop() {
	defer svr.wg.Done()
	defer svr.ln.Close()

	var conn net.Conn
	var err error

	for {
		select {
		case svr.pending <- struct{}{}:
			for {
				if conn, err = svr.ln.Accept(); err == nil {
					break
				}
			}

			go svr.setupConn(conn, inbound)
		case <-svr.term:
			return
		}
	}
}

func (svr *Server) setupConn(c net.Conn, flag connFlag) {
	ts := &conn{
		AsyncMsgConn: NewAsyncMsgConn(c, nil),
		flags:        flag,
	}

	// handshake data, add remoteIP and remotePort
	handshake := *svr.handshake
	tcpAddr := c.RemoteAddr().(*net.TCPAddr)
	handshake.RemoteIP = tcpAddr.IP
	handshake.RemotePort = uint16(tcpAddr.Port)
	data, err := handshake.Serialize()
	if err != nil {
		ts.Close(nil)
		return
	}
	sig := ed25519.Sign(svr.PrivateKey, data)
	data = append(sig, data...)

	their, err := ts.Handshake(data)

	if err != nil {
		ts.Close(err)
		svr.log.Error(fmt.Sprintf("handshake with %s error: %v", c.RemoteAddr(), err))
	} else if their.NetID != svr.NetID {
		err = fmt.Errorf("different NetID: our %s, their %s", svr.NetID, their.NetID)
		ts.Close(err)
		svr.log.Error(fmt.Sprintf("handshake with %s error: %v", c.RemoteAddr(), err))
	} else {
		ts.name = their.Name
		ts.cmdSets = their.CmdSets

		// use to discribe the connection
		ts.remoteID = their.ID
		ts.remoteIP = handshake.RemoteIP
		ts.remotePort = handshake.RemotePort

		ts.localID = svr.self.ID
		ts.localIP = their.RemoteIP
		ts.localPort = their.RemotePort

		svr.addPeer <- ts
	}

	<-svr.pending
}

func (svr *Server) checkConn(id discovery.NodeID, flag connFlag) error {
	if id == svr.self.ID {
		return DiscSelf
	}

	if svr.peers.Has(id) {
		return DiscAlreadyConnected
	}

	// static can be connected even if peers too many
	if flag == static {
		return nil
	}

	if uint(svr.peers.Size()) >= svr.MaxPeers {
		return DiscTooManyPeers
	}

	if flag.is(inbound) && uint(svr.peers.inbound) >= svr.maxInboundPeers() {
		return DiscTooManyInboundPeers
	}

	return nil
}

func (svr *Server) loop() {
	defer svr.wg.Done()

	var peersCount uint
loop:
	for {
		select {
		case <-svr.term:
			break loop
		case c := <-svr.addPeer:
			err := svr.checkConn(c.remoteID, c.flags)

			if err == nil {
				if p, err := NewPeer(c, svr.Protocols); err == nil {
					svr.peers.Add(p)
					peersCount = svr.peers.Size()
					svr.log.Info(fmt.Sprintf("create new peer %s, total: %d", p, peersCount))

					monitor.LogDuration("p2p/peer", "count", int64(peersCount))
					monitor.LogEvent("p2p/peer", "create")

					svr.wg.Add(1)
					go svr.runPeer(p)
				} else {
					svr.log.Error(fmt.Sprintf("create new peer error: %v", err))
				}
			} else {
				c.Close(err)
				svr.log.Error(fmt.Sprintf("can`t create new peer: %v", err))
			}

		case p := <-svr.delPeer:
			svr.peers.Del(p)
			peersCount = svr.peers.Size()
			svr.log.Info(fmt.Sprintf("delete peer %s, total: %d", p, peersCount))

			monitor.LogDuration("p2p/peer", "count", int64(peersCount))
			monitor.LogEvent("p2p/peer", "delete")

			if p.ts.is(static) {
				svr.dial(p.ID(), p.RemoteAddr(), static)
			}

			//if peersCount < svr.MaxPeers {
			//	svr.discv.Need(svr.MaxPeers - peersCount)
			//}
		}
	}

	svr.peers.Traverse(func(id discovery.NodeID, p *Peer) {
		p.Disconnect(DiscQuitting)
	})
}

func (svr *Server) runPeer(p *Peer) {
	defer svr.wg.Done()

	err := p.run()
	if err != nil {
		svr.log.Error("run peer error", "error", err)
	}
	svr.delPeer <- p
}

func (svr *Server) Peers() []*PeerInfo {
	return svr.peers.Info()
}

func (svr *Server) PeersCount() uint {
	return svr.peers.Size()
}

func (svr *Server) NodeInfo() *NodeInfo {
	protocols := make([]string, len(svr.Protocols))
	for i, protocol := range svr.Protocols {
		protocols[i] = protocol.String()
	}

	return &NodeInfo{
		ID:    svr.self.ID.String(),
		Name:  svr.Name,
		Url:   svr.self.String(),
		NetID: svr.NetID,
		Address: &address{
			IP:  svr.self.IP,
			TCP: svr.self.TCP,
			UDP: svr.self.UDP,
		},
		Protocols: protocols,
	}
}

func (svr *Server) URL() string {
	return svr.self.String()
}

func (svr *Server) Available() bool {
	return svr.PeersCount() > 0
}

func (svr *Server) maxOutboundPeers() uint {
	return svr.MaxPeers - svr.maxInboundPeers()
}

func (svr *Server) maxInboundPeers() uint {
	return svr.MaxPeers / svr.MaxInboundRatio
}

// @section NodeInfo
type NodeInfo struct {
	ID        string    `json:"remoteID"`
	Name      string    `json:"name"`
	Url       string    `json:"url"`
	NetID     NetworkID `json:"netId"`
	Address   *address  `json:"address"`
	Protocols []string  `json:"protocols"`
}

type address struct {
	IP  net.IP `json:"ip"`
	TCP uint16 `json:"tcp"`
	UDP uint16 `json:"udp"`
}
