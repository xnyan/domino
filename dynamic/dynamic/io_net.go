package dynamic

import (
	"bufio"
	"encoding/gob"
	"errors"
	"fmt"
	"net"
	"sync"
)

const (
	Msg_Type_PaxosProposeReq uint8 = iota
	Msg_Type_PaxosProposeReply
	Msg_Type_FpProposeReq
	Msg_Type_FpProposeReply
	Msg_Type_ReplicaMsg
	//Msg_Type_EmptyReply
	Msg_Type_TestReq
	Msg_Type_TestReply
)

type NetIo interface {
	SendMsg(msgType uint8, msg interface{}) error
	RecvMsg() (msgType uint8, msg interface{}, err error)
	SendByte(b uint8) error
	RecvByte() (uint8, error)
	Close() error
}

type netIoImpl struct {
	// Network Connection
	conn net.Conn

	// Reader
	reader *bufio.Reader
	dec    *gob.Decoder

	// Writer
	writer *bufio.Writer
	enc    *gob.Encoder
}

func NewNetIo(conn net.Conn) NetIo {
	n := &netIoImpl{
		conn: conn,
	}
	n.reader = bufio.NewReader(n.conn)
	n.dec = gob.NewDecoder(n.reader)
	n.writer = bufio.NewWriter(n.conn)
	n.enc = gob.NewEncoder(n.writer)
	return n
}

func (n *netIoImpl) SendByte(b uint8) error {
	err := n.writer.WriteByte(b)
	if err != nil {
		return err
	}

	err = n.writer.Flush()
	if err != nil {
		return err
	}

	return nil
}

func (n *netIoImpl) RecvByte() (uint8, error) {
	b, err := n.reader.ReadByte()
	return b, err
}

func (n *netIoImpl) SendMsg(msgType uint8, msg interface{}) error {
	err := n.writer.WriteByte(msgType)
	if err != nil {
		return err
	}

	err = n.enc.Encode(msg)
	if err != nil {
		return err
	}

	err = n.writer.Flush()
	if err != nil {
		return err
	}

	return nil
}

func (n *netIoImpl) RecvMsg() (msgType uint8, msg interface{}, err error) {
	if msgType, err = n.reader.ReadByte(); err != nil {
		return
	}

	switch uint8(msgType) {
	case Msg_Type_PaxosProposeReq:
		var m PaxosProposeReq
		if err = n.dec.Decode(&m); err == nil {
			msg = &m
		}
	case Msg_Type_PaxosProposeReply:
		var m PaxosProposeReply
		if err = n.dec.Decode(&m); err == nil {
			msg = &m
		}
	case Msg_Type_FpProposeReq:
		var m FpProposeReq
		if err = n.dec.Decode(&m); err == nil {
			msg = &m
		}
	case Msg_Type_FpProposeReply:
		var m FpProposeReply
		if err = n.dec.Decode(&m); err == nil {
			msg = &m
		}
	case Msg_Type_ReplicaMsg:
		var m ReplicaMsg
		if err = n.dec.Decode(&m); err == nil {
			msg = &m
		}
	case Msg_Type_TestReq:
		var m TestReq
		if err = n.dec.Decode(&m); err == nil {
			msg = &m
		}
	case Msg_Type_TestReply:
		var m TestReply
		if err = n.dec.Decode(&m); err == nil {
			msg = &m
		}
	default:
		err = errors.New(fmt.Sprintf("Unkonwn message type %d", msgType))
	}

	return
}

func (n *netIoImpl) Close() error {
	return n.conn.Close()
}

type SyncNetIo struct {
	netIo NetIo
	lock  sync.Mutex
}

func NewSyncNetIo(netIo NetIo) *SyncNetIo {
	n := &SyncNetIo{
		netIo: netIo,
	}
	return n
}

func (n *SyncNetIo) SendMsg(msgType uint8, msg interface{}) error {
	n.lock.Lock()
	defer n.lock.Unlock()

	return n.netIo.SendMsg(msgType, msg)
}
