package mqtt

import (
	"bytes"
	"errors"
	"io"
)

// OoS only support QoS 0
const (
	QosAtMostOnce = TagQosLevel(iota)
	QosAtLeastOnce
	QosExactlyOnce
	QosInvalid
)

// Max Payload size
const (
	MaxPayloadSize = (1 << (4 * 7)) - 1
)

type TagQosLevel uint8

func (qos TagQosLevel) IsValid() bool {
	return qos < QosInvalid && qos >= QosAtMostOnce
}

func (qos TagQosLevel) HasId() bool {
	return qos == QosAtLeastOnce || qos == QosExactlyOnce
}

// Message Type
const (
	MsgConnect = TagMessageType(iota + 1)
	MsgConnAck
	MsgPublish
	MsgPubAck
	MsgPubRec
	MsgPubRel
	MsgPubComp
	MsgSubscribe
	MsgSubAck
	MsgUnsubscribe
	MsgUnsubAck
	MsgPingReq
	MsgPingResp
	MsgDisconnect
	MsgInvalid
)

//  retcode
const (
	RetCodeAccepted = TagRetCode(iota)
	RetCodeUnacceptableProtocolVersion
	RetCodeIdentifierRejected
	RetCodeServerUnavailable
	RetCodeBadUsernameOrPassword
	RetCodeNotAuthorized
	RetCodeInvalid
)

type TagRetCode uint8

func (rc TagRetCode) IsValid() bool {
	return rc >= RetCodeAccepted && rc < RetCodeInvalid
}

type TagMessageType uint8

func (msg TagMessageType) IsValid() bool {
	return msg >= MsgConnect && msg < MsgInvalid
}

// message interface
type Message interface {
	Encode(w io.Writer) error
	Decode(r io.Reader, hdr Header, packetRemaining int32) error
}

// message fix header
type Header struct {
	DupFlag  bool
	QosLevel TagQosLevel
	Retain   bool
}

func (hdr *Header) Encode(w io.Writer, msgType TagMessageType, remainingLength int32) error {
	buf := new(bytes.Buffer)
	err := hdr.EncodeInto(buf, msgType, remainingLength)
	if err != nil {
		return err
	}

	_, err = w.Write(buf.Bytes())

	return err
}

func (hdr *Header) EncodeInto(buf *bytes.Buffer, msgType TagMessageType, remainingLength int32) error {
	if !hdr.QosLevel.IsValid() {
		return errors.New("Invalid Qos level")
	}

	if !msgType.IsValid() {
		return errors.New("Invalid MsgType")
	}

	val := byte(msgType) << 4
	val |= (boolToByte(hdr.DupFlag) << 3)
	val |= byte(hdr.QosLevel) << 1
	val |= boolToByte(hdr.Retain)
	buf.WriteByte(val)
	encodeLength(remainingLength, buf)

	return nil
}

func (hdr *Header) Decode(r io.Reader) (msgType TagMessageType, remainingLength int32, err error) {
	var buf [1]byte

	if _, err := io.ReadFull(r, buf[:]); err != nil {
		return 0, 0, err
	}

	byte1 := buf[0]
	msgType = TagMessageType(byte1 & 0xf0 >> 4)

	*hdr = Header{
		DupFlag:  byte1&0x08 > 0,
		QosLevel: TagQosLevel(byte1 & 0x06 >> 1),
		Retain:   byte1&0x01 > 0,
	}

	remainingLength, err = decodeLength(r)

	return msgType, remainingLength, err
}

func writeMessage(w io.Writer, msgType TagMessageType, hdr *Header, payloadBuf *bytes.Buffer, extraLength int32) error {
	totalPayloadLength := int64(len(payloadBuf.Bytes())) + int64(extraLength)
	if totalPayloadLength > MaxPayloadSize {
		return errors.New("message too long")
	}

	buf := new(bytes.Buffer)
	err := hdr.EncodeInto(buf, msgType, int32(totalPayloadLength))
	if err != nil {
		return err
	}

	buf.Write(payloadBuf.Bytes())
	_, err = w.Write(buf.Bytes())

	return err
}

// Connect represents an MQTT CONNECT message.
type Connect struct {
	Header
	ProtocolName               string
	ProtocolVersion            uint8
	WillRetain                 bool
	WillFlag                   bool
	CleanSession               bool
	WillQos                    TagQosLevel
	KeepAliveTimer             uint16
	ClientId                   string
	WillTopic, WillMessage     string
	UsernameFlag, PasswordFlag bool
	Username, Password         string
}

func (msg *Connect) Encode(w io.Writer) (err error) {
	if msg.WillQos > QosInvalid {
		return errors.New("invalid Qos")
	}

	buf := new(bytes.Buffer)

	flags := boolToByte(msg.UsernameFlag) << 7
	flags |= boolToByte(msg.PasswordFlag) << 6
	flags |= boolToByte(msg.WillRetain) << 5
	flags |= byte(msg.WillQos) << 3
	flags |= boolToByte(msg.WillFlag) << 2
	flags |= boolToByte(msg.CleanSession) << 1

	setString(msg.ProtocolName, buf)
	setUint8(msg.ProtocolVersion, buf)
	buf.WriteByte(flags)
	setUint16(msg.KeepAliveTimer, buf)
	setString(msg.ClientId, buf)
	if msg.WillFlag {
		setString(msg.WillTopic, buf)
		setString(msg.WillMessage, buf)
	}
	if msg.UsernameFlag {
		setString(msg.Username, buf)
	}
	if msg.PasswordFlag {
		setString(msg.Password, buf)
	}

	return writeMessage(w, MsgConnect, &msg.Header, buf, 0)
}

func (msg *Connect) Decode(r io.Reader, hdr Header, packetRemaining int32) (err error) {
	protocolName, err := getString(r, &packetRemaining)
	if err != nil {
		return err
	}
	protocolVersion, err := getUint8(r, &packetRemaining)
	if err != nil {
		return err
	}
	flags, err := getUint8(r, &packetRemaining)
	if err != nil {
		return err
	}
	keepAliveTimer, err := getUint16(r, &packetRemaining)
	if err != nil {
		return err
	}
	clientId, err := getString(r, &packetRemaining)
	if err != nil {
		return err
	}

	*msg = Connect{
		Header:          hdr,
		ProtocolName:    protocolName,
		ProtocolVersion: protocolVersion,
		UsernameFlag:    flags&0x80 > 0,
		PasswordFlag:    flags&0x40 > 0,
		WillRetain:      flags&0x20 > 0,
		WillQos:         TagQosLevel(flags & 0x18 >> 3),
		WillFlag:        flags&0x04 > 0,
		CleanSession:    flags&0x02 > 0,
		KeepAliveTimer:  keepAliveTimer,
		ClientId:        clientId,
	}

	if msg.WillFlag {
		msg.WillTopic, err = getString(r, &packetRemaining)
		if err != nil {
			return err
		}
		msg.WillMessage, err = getString(r, &packetRemaining)
		if err != nil {
			return err
		}
	}
	if msg.UsernameFlag {
		msg.Username, err = getString(r, &packetRemaining)
		if err != nil {
			return err
		}
	}
	if msg.PasswordFlag {
		msg.Password, err = getString(r, &packetRemaining)
		if err != nil {
			return err
		}
	}

	if packetRemaining != 0 {
		return errors.New("message too long")
	}

	return nil
}

// ConnAck represents an MQTT CONNACK message.
type ConnAck struct {
	Header
	ReturnCode TagRetCode
}

func (msg *ConnAck) Encode(w io.Writer) (err error) {
	buf := new(bytes.Buffer)

	buf.WriteByte(byte(0)) // Reserved byte.
	setUint8(uint8(msg.ReturnCode), buf)

	return writeMessage(w, MsgConnAck, &msg.Header, buf, 0)
}

func (msg *ConnAck) Decode(r io.Reader, hdr Header, packetRemaining int32) (err error) {
	msg.Header = hdr

	_, err = getUint8(r, &packetRemaining) // Skip reserved byte.
	if err != nil {
		return err
	}

	code, err := getUint8(r, &packetRemaining)
	if err != nil {
		return err
	}
	msg.ReturnCode = TagRetCode(code)
	if !msg.ReturnCode.IsValid() {
		return errors.New("invliad retcode")
	}

	if packetRemaining != 0 {
		return errors.New("message too long")
	}

	return nil
}

// Publish represents an MQTT PUBLISH message.
type Publish struct {
	Header
	TopicName string
	MessageId uint16
	Payload   Payload
}

func (msg *Publish) Encode(w io.Writer) (err error) {
	buf := new(bytes.Buffer)

	setString(msg.TopicName, buf)
	if msg.Header.QosLevel.HasId() {
		setUint16(msg.MessageId, buf)
	}

	if err = msg.Payload.WritePayload(buf); err != nil {
		return err
	}

	if err = writeMessage(w, MsgPublish, &msg.Header, buf, int32(0)); err != nil {
		return err
	}

	return nil
}

func (msg *Publish) Decode(r io.Reader, hdr Header, packetRemaining int32) (err error) {
	msg.Header = hdr

	msg.TopicName, err = getString(r, &packetRemaining)
	if err != nil {
		return err
	}
	if msg.Header.QosLevel.HasId() {
		msg.MessageId, err = getUint16(r, &packetRemaining)
		if err != nil {
			return err
		}
	}

	payloadReader := &io.LimitedReader{r, int64(packetRemaining)}
	msg.Payload = make(BytesPayload, int(packetRemaining))

	return msg.Payload.ReadPayload(payloadReader, int(packetRemaining))
}

// PubAck represents an MQTT PUBACK message.
type PubAck struct {
	Header
	MessageId uint16
}

func (msg *PubAck) Encode(w io.Writer) error {
	return encodeAckCommon(w, &msg.Header, msg.MessageId, MsgPubAck)
}

func (msg *PubAck) Decode(r io.Reader, hdr Header, packetRemaining int32) (err error) {
	msg.Header = hdr
	return decodeAckCommon(r, packetRemaining, &msg.MessageId)
}

// PubRec represents an MQTT PUBREC message.
type PubRec struct {
	Header
	MessageId uint16
}

func (msg *PubRec) Encode(w io.Writer) error {
	return encodeAckCommon(w, &msg.Header, msg.MessageId, MsgPubRec)
}

func (msg *PubRec) Decode(r io.Reader, hdr Header, packetRemaining int32) (err error) {
	msg.Header = hdr
	return decodeAckCommon(r, packetRemaining, &msg.MessageId)
}

// PubRel represents an MQTT PUBREL message.
type PubRel struct {
	Header
	MessageId uint16
}

func (msg *PubRel) Encode(w io.Writer) error {
	return encodeAckCommon(w, &msg.Header, msg.MessageId, MsgPubRel)
}

func (msg *PubRel) Decode(r io.Reader, hdr Header, packetRemaining int32) (err error) {
	msg.Header = hdr
	return decodeAckCommon(r, packetRemaining, &msg.MessageId)
}

// PubComp represents an MQTT PUBCOMP message.
type PubComp struct {
	Header
	MessageId uint16
}

func (msg *PubComp) Encode(w io.Writer) error {
	return encodeAckCommon(w, &msg.Header, msg.MessageId, MsgPubComp)
}

func (msg *PubComp) Decode(r io.Reader, hdr Header, packetRemaining int32) (err error) {
	msg.Header = hdr
	return decodeAckCommon(r, packetRemaining, &msg.MessageId)
}

// Subscribe represents an MQTT SUBSCRIBE message.
type Subscribe struct {
	Header
	MessageId uint16
	Topics    []TopicQos
}

type TopicQos struct {
	Topic string
	Qos   TagQosLevel
}

func (msg *Subscribe) Encode(w io.Writer) (err error) {
	buf := new(bytes.Buffer)
	if msg.Header.QosLevel.HasId() {
		setUint16(msg.MessageId, buf)
	}
	for _, topicSub := range msg.Topics {
		setString(topicSub.Topic, buf)
		setUint8(uint8(topicSub.Qos), buf)
	}

	return writeMessage(w, MsgSubscribe, &msg.Header, buf, 0)
}

func (msg *Subscribe) Decode(r io.Reader, hdr Header, packetRemaining int32) (err error) {
	msg.Header = hdr

	if msg.Header.QosLevel.HasId() {
		msg.MessageId, err = getUint16(r, &packetRemaining)
		if err != nil {
			return err
		}
	}
	var topics []TopicQos
	for packetRemaining > 0 {
		topic, err := getString(r, &packetRemaining)
		if err != nil {
			return err
		}
		qos, err := getUint8(r, &packetRemaining)
		if err != nil {
			return err
		}
		topics = append(topics, TopicQos{
			Topic: topic,
			Qos:   TagQosLevel(qos),
		})
	}
	msg.Topics = topics

	return nil
}

// SubAck represents an MQTT SUBACK message.
type SubAck struct {
	Header
	MessageId uint16
	TopicsQos []TagQosLevel
}

func (msg *SubAck) Encode(w io.Writer) (err error) {
	buf := new(bytes.Buffer)
	setUint16(msg.MessageId, buf)
	for i := 0; i < len(msg.TopicsQos); i += 1 {
		setUint8(uint8(msg.TopicsQos[i]), buf)
	}

	return writeMessage(w, MsgSubAck, &msg.Header, buf, 0)
}

func (msg *SubAck) Decode(r io.Reader, hdr Header, packetRemaining int32) (err error) {
	msg.Header = hdr

	msg.MessageId, err = getUint16(r, &packetRemaining)
	if err != nil {
		return err
	}
	topicsQos := make([]TagQosLevel, 0)
	for packetRemaining > 0 {
		qos, err := getUint8(r, &packetRemaining)
		if err != nil {
			return err
		}
		grantedQos := TagQosLevel(qos & 0x03)
		topicsQos = append(topicsQos, grantedQos)
	}
	msg.TopicsQos = topicsQos

	return nil
}

// Unsubscribe represents an MQTT UNSUBSCRIBE message.
type Unsubscribe struct {
	Header
	MessageId uint16
	Topics    []string
}

func (msg *Unsubscribe) Encode(w io.Writer) (err error) {
	buf := new(bytes.Buffer)
	if msg.Header.QosLevel.HasId() {
		setUint16(msg.MessageId, buf)
	}
	for _, topic := range msg.Topics {
		setString(topic, buf)
	}

	return writeMessage(w, MsgUnsubscribe, &msg.Header, buf, 0)
}

func (msg *Unsubscribe) Decode(r io.Reader, hdr Header, packetRemaining int32) (err error) {
	msg.Header = hdr

	if msg.Header.QosLevel.HasId() {
		msg.MessageId, err = getUint16(r, &packetRemaining)
		if err != nil {
			return err
		}
	}
	topics := make([]string, 0)
	for packetRemaining > 0 {
		topic, err := getString(r, &packetRemaining)
		if err != nil {
			return err
		}
		topics = append(topics, topic)
	}
	msg.Topics = topics

	return nil
}

// UnsubAck represents an MQTT UNSUBACK message.
type UnsubAck struct {
	Header
	MessageId uint16
}

func (msg *UnsubAck) Encode(w io.Writer) error {
	return encodeAckCommon(w, &msg.Header, msg.MessageId, MsgUnsubAck)
}

func (msg *UnsubAck) Decode(r io.Reader, hdr Header, packetRemaining int32) (err error) {
	msg.Header = hdr
	return decodeAckCommon(r, packetRemaining, &msg.MessageId)
}

// PingReq represents an MQTT PINGREQ message.
type PingReq struct {
	Header
}

func (msg *PingReq) Encode(w io.Writer) error {
	return msg.Header.Encode(w, MsgPingReq, 0)
}

func (msg *PingReq) Decode(r io.Reader, hdr Header, packetRemaining int32) error {
	if packetRemaining != 0 {
		return errors.New("msg too long")
	}
	return nil
}

// PingResp represents an MQTT PINGRESP message.
type PingResp struct {
	Header
}

func (msg *PingResp) Encode(w io.Writer) error {
	return msg.Header.Encode(w, MsgPingResp, 0)
}

func (msg *PingResp) Decode(r io.Reader, hdr Header, packetRemaining int32) error {
	if packetRemaining != 0 {
		return errors.New("msg too long")
	}
	return nil
}

// Disconnect represents an MQTT DISCONNECT message.
type Disconnect struct {
	Header
}

func (msg *Disconnect) Encode(w io.Writer) error {
	return msg.Header.Encode(w, MsgDisconnect, 0)
}

func (msg *Disconnect) Decode(r io.Reader, hdr Header, packetRemaining int32) error {
	if packetRemaining != 0 {
		return errors.New("msg too long")
	}
	return nil
}

func encodeAckCommon(w io.Writer, hdr *Header, messageId uint16, msgType TagMessageType) error {
	buf := new(bytes.Buffer)
	setUint16(messageId, buf)
	return writeMessage(w, msgType, hdr, buf, 0)
}

func decodeAckCommon(r io.Reader, packetRemaining int32, messageId *uint16) (err error) {
	*messageId, err = getUint16(r, &packetRemaining)
	if err != nil {
		return err
	}

	if packetRemaining != 0 {
		return errors.New("msg too long")
	}

	return nil
}

// DecodeOneMessage decodes one message from r. config provides specifics on
// how to decode messages, nil indicates that the DefaultDecoderConfig should
// be used.
func DecodeOneMessage(r io.Reader) (msg Message, err error) {
	var hdr Header
	var msgType TagMessageType
	var packetRemaining int32
	msgType, packetRemaining, err = hdr.Decode(r)
	if err != nil {
		return nil, err
	}

	msg, err = NewMessage(msgType)
	if err != nil {
		return nil, err
	}

	return msg, msg.Decode(r, hdr, packetRemaining)
}

func NewMessage(msgType TagMessageType) (msg Message, err error) {
	switch msgType {
	case MsgConnect:
		msg = new(Connect)
	case MsgConnAck:
		msg = new(ConnAck)
	case MsgPublish:
		msg = new(Publish)
	case MsgPubAck:
		msg = new(PubAck)
	case MsgPubRec:
		msg = new(PubRec)
	case MsgPubRel:
		msg = new(PubRel)
	case MsgPubComp:
		msg = new(PubComp)
	case MsgSubscribe:
		msg = new(Subscribe)
	case MsgUnsubAck:
		msg = new(UnsubAck)
	case MsgSubAck:
		msg = new(SubAck)
	case MsgUnsubscribe:
		msg = new(Unsubscribe)
	case MsgPingReq:
		msg = new(PingReq)
	case MsgPingResp:
		msg = new(PingResp)
	case MsgDisconnect:
		msg = new(Disconnect)
	default:
		return nil, errors.New("msgType error")
	}

	return msg, nil
}
