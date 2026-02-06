package qqwry

import (
	"encoding/binary"
	"errors"
	"io"
	"net"
	"os"
	"strings"

	"golang.org/x/text/encoding/simplifiedchinese"
)

const (
	indexLen      = 7
	redirectMode1 = 0x01
	redirectMode2 = 0x02
)

type QQWry struct {
	Data     []byte
	FilePath string
	fd       *os.File
	Offset   int64
}

func NewQQWry(filePath string) (*QQWry, error) {
	if len(filePath) == 0 {
		return nil, errors.New("file path is empty")
	}

	_, err := os.Stat(filePath)
	if err != nil && os.IsNotExist(err) {
		return nil, errors.New("file does not exist")
	}
	wry := &QQWry{}

	wry.fd, err = os.OpenFile(filePath, os.O_RDONLY, 0400)
	if err != nil {
		return nil, err
	}
	defer wry.fd.Close()

	wry.Data, err = io.ReadAll(wry.fd)
	if err != nil {
		return nil, err
	}

	return wry, nil
}

func (q *QQWry) readData(num int, offset ...int64) (rs []byte) {
	if len(offset) > 0 {
		q.Offset = offset[0]
	}
	nums := int64(num)
	end := q.Offset + nums
	dataNum := int64(len(q.Data))
	if q.Offset > dataNum {
		return nil
	}

	if end > dataNum {
		end = dataNum
	}

	rs = q.Data[q.Offset:end]
	q.Offset = end
	return
}

func (q *QQWry) Find(ip string) (string, string, string) {
	if strings.Count(ip, ".") != 3 {
		return "", "", ""
	}
	offset := q.searchIndex(binary.BigEndian.Uint32(net.ParseIP(ip).To4()))
	if offset <= 0 {
		return "", "", ""
	}

	var country []byte
	var area []byte

	mode := q.readMode(offset + 4)
	if mode == redirectMode1 {
		countryOffset := q.readUInt24()
		mode = q.readMode(countryOffset)
		if mode == redirectMode2 {
			c := q.readUInt24()
			country = q.readString(c)
			countryOffset += 4
		} else {
			country = q.readString(countryOffset)
			countryOffset += uint32(len(country) + 1)
		}
		area = q.readArea(countryOffset)
	} else if mode == redirectMode2 {
		countryOffset := q.readUInt24()
		country = q.readString(countryOffset)
		area = q.readArea(offset + 8)
	} else {
		country = q.readString(offset + 4)
		area = q.readArea(offset + uint32(5+len(country)))
	}

	enc := simplifiedchinese.GBK.NewDecoder()
	cou, _ := enc.String(string(country))
	are, _ := enc.String(string(area))

	return ip, cou, are
}

func (q *QQWry) readMode(offset uint32) byte {
	mode := q.readData(1, int64(offset))
	return mode[0]
}

func (q *QQWry) readArea(offset uint32) []byte {
	mode := q.readMode(offset)
	if mode == redirectMode1 || mode == redirectMode2 {
		areaOffset := q.readUInt24()
		if areaOffset == 0 {
			return []byte("")
		}
		return q.readString(areaOffset)
	}
	return q.readString(offset)
}

func (q *QQWry) readString(offset uint32) []byte {
	q.Offset = int64(offset)
	data := make([]byte, 0, 30)
	buf := make([]byte, 1)
	for {
		buf = q.readData(1)
		if buf[0] == 0 {
			break
		}
		data = append(data, buf[0])
	}
	return data
}

func (q *QQWry) searchIndex(ip uint32) uint32 {
	header := q.readData(8, 0)

	start := binary.LittleEndian.Uint32(header[:4])
	end := binary.LittleEndian.Uint32(header[4:])

	buf := make([]byte, indexLen)
	mid := uint32(0)
	_ip := uint32(0)

	for {
		mid = q.getMiddleOffset(start, end)
		buf = q.readData(indexLen, int64(mid))
		_ip = binary.LittleEndian.Uint32(buf[:4])

		if end-start == indexLen {
			offset := byteToUInt32(buf[4:])
			buf = q.readData(indexLen)
			if ip < binary.LittleEndian.Uint32(buf[:4]) {
				return offset
			}
			return 0
		}

		if _ip > ip {
			end = mid
		} else if _ip < ip { // 找到的比较小，向后移
			start = mid
		} else if _ip == ip {
			return byteToUInt32(buf[4:])
		}
	}
}

func (q *QQWry) readUInt24() uint32 {
	buf := q.readData(3)
	return byteToUInt32(buf)
}

func (q *QQWry) getMiddleOffset(start uint32, end uint32) uint32 {
	records := ((end - start) / indexLen) >> 1
	return start + records*indexLen
}

func byteToUInt32(data []byte) uint32 {
	i := uint32(data[0]) & 0xff
	i |= (uint32(data[1]) << 8) & 0xff00
	i |= (uint32(data[2]) << 16) & 0xff0000
	return i
}
