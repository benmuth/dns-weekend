package dns

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"math/rand"
	"strings"
)

// class DNSHeader:
//     id: int
//     flags: int
//     num_questions: int = 0
//     num_answers: int = 0
//     num_authorities: int = 0
//     num_additionals: int = 0

type DNSHeader struct {
	id              uint16
	flags           uint16
	num_questions   uint16
	num_answers     uint16
	num_authorities uint16
	num_additionals uint16
}

type DNSQuestion struct {
	name  []byte
	type_ uint16
	class uint16
}

func headerToBytes(header DNSHeader) ([]byte, error) {
	buf := new(bytes.Buffer)
	err := binary.Write(buf, binary.BigEndian, header)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func questionToBytes(question DNSQuestion) ([]byte, error) {
	buf := new(bytes.Buffer)
	err := binary.Write(buf, binary.BigEndian, [2]uint16{question.type_, question.class})
	if err != nil {
		return nil, err
	}
	return append(question.name, buf.Bytes()...), nil
}

func encodeDNSName(domain string) ([]byte, error) {
	buf := new(bytes.Buffer)
	for _, part := range strings.Split(domain, ".") {
		_ = buf.WriteByte(byte(len(part)))
		if _, err := buf.WriteString(part); err != nil {
			return nil, err
		}
	}
	_ = buf.WriteByte(0)
	return buf.Bytes(), nil
}

const TYPE_A = 1
const CLASS_IN = 1

func buildQuery(domainName string, recordType uint16) ([]byte, error) {
	name, err := encodeDNSName(domainName)
	if err != nil {
		panic(err)
	}
	id := uint16(rand.Intn(65535))
	recursionDesired := uint16(1 << 8)
	header := DNSHeader{id: id, num_questions: 1, flags: recursionDesired}
	question := DNSQuestion{name: name, type_: recordType, class: CLASS_IN}
	headerBytes, err := headerToBytes(header)
	if err != nil {
		return nil, err
	}
	questionBytes, err := questionToBytes(question)
	if err != nil {
		return nil, err
	}
	return append(headerBytes, questionBytes...), nil
}

func main() {
	fmt.Println("Hello world!")
	header := DNSHeader{
		1, 2, 3, 4, 5, 6,
	}
	headerToBytes(header)
}
