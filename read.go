package gma

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
)

func ReadGMA(r *bufio.Reader) ([]Entry, error) {
	// Header
	{
		var requiredHeader [4]byte
		if _, err := io.ReadFull(r, requiredHeader[:]); err != nil {
			return nil, err
		}
		if !bytes.Equal(requiredHeader[:], HEADER) {
			return nil, fmt.Errorf("invalid header: %s", requiredHeader)
		}
	}
	// Version
	{
		var version int8
		if err := binary.Read(r, binary.LittleEndian, &version); err != nil {
			return nil, err
		}
		if version != VERSION {
			return nil, fmt.Errorf("invalid version: %d", version)
		}
	}
	// SteamID64
	if _, err := r.Discard(8); err != nil {
		return nil, err
	}
	// Timestamp
	if _, err := r.Discard(8); err != nil {
		return nil, err
	}
	// Required content
	if _, err := r.Discard(1); err != nil {
		return nil, err
	}
	// Addon name
	if _, err := r.ReadString('\x00'); err != nil {
		return nil, err
	}
	// Addon description
	if _, err := r.ReadString('\x00'); err != nil {
		return nil, err
	}
	// Addon author
	if _, err := r.ReadString('\x00'); err != nil {
		return nil, err
	}
	// Addon version
	if _, err := r.Discard(4); err != nil {
		return nil, err
	}

	fileCount := 0
	entries := make([]Entry, 0, 10)
	for {
		// loop till we hit null
		{
			var idx uint32
			if err := binary.Read(r, binary.LittleEndian, &idx); err != nil {
				return nil, err
			}
			if idx == 0 {
				break
			}
		}
		// File name
		name, err := r.ReadString('\x00')
		if err != nil {
			return nil, err
		}
		name = name[:len(name)-1] // remove null byte
		var size int64
		if err := binary.Read(r, binary.LittleEndian, &size); err != nil {
			return nil, err
		}
		// File CRC32
		if _, err := r.Discard(4); err != nil {
			return nil, err
		}
		entries = append(entries, Entry{
			name: name,
			size: uint64(size),
		})
		fileCount++
	}
	for i := 0; i < fileCount; i++ {
		var content = make([]byte, entries[i].size)
		if _, err := io.ReadFull(r, content); err != nil {
			return nil, err
		}
		entries[i].content = content
	}
	if _, err := r.Discard(4); err != nil {
		return nil, err
	}
	return entries, nil
}
