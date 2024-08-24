package main

import (
  "encoding/binary"
  "fmt"
  "log"
  "net"
  "sync"
  "time"
)

const (
  msgTypeData = 0x01
  msgTypeAck = 0x02
  msgTypeNack = 0x03
  msgTypeClose = 0x04

  header = 8
  maxMsgSize = 1024

  retransmitTimeout = 500 * time.Millisecond
  maxRetransmit = 3
)

type udpClient struct {
  conn *net.UDPConn
  addr *net.UDPAddr
  seqNum uint32
  ackNum uint32
  mutex sync.Mutex
  closing bool
  closed bool
}

func newUDPClient(addr string) (*udpClient, error){
  udpAddr, err != net.ResolveUDPAddr("udp", addr)
  if err != nil {
    return nil, err

  } 
  conn, err := net.DialUDP("udp", nil, updAddr)
  if err != nil {
    return nil, err
  }

  return &udpClient{
    conn : conn,
    addr: udpAddr,
    seqNum : 1,
    ackNum : 1, 
    closing : false,
    closed : false,
  }, nil 

}

func (c *udpClient) sendMessage(msgTypes byte, data []byte) error {
  header := make([]byte, headerSize)
  binary.BigEndian.PutUint32(header[:4], c.seqNum)
  header[4] = msgType
  binary.BigEndian.PutUint32(header[5:], c.ackNum)

  for i := 0; i < len(data); i += maxMsgSize - headerSize {
    fragment := data[i:min]
  }
}


