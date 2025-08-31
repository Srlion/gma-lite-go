package gma

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"io"
	"strings"
	"time"
)

type Builder struct {
	name        string
	steamID64   uint64
	author      string
	description string
	entries     []Entry
}

func NewBuilder(name string, steamid64 int64, size int) *Builder {
	return &Builder{
		name:      name,
		steamID64: uint64(steamid64),
		author:    "unknown",
		entries:   make([]Entry, 0, size),
	}
}

func (builder *Builder) SetDescription(desc string) {
	builder.description = desc
}

func (builder *Builder) SetAuthor(author string) {
	builder.author = author
}

func (builder *Builder) FileFromBytes(name string, bytes []byte) {
	builder.entries = append(builder.entries, Entry{
		name:    name,
		content: bytes,
	})
}

func (builder *Builder) FileFromString(name, content string) {
	builder.FileFromBytes(name, []byte(content))
}

func (builder *Builder) WriteGMATo(w io.Writer) error {
	// Wrap the writer with a buffered writer
	bw := bufio.NewWriter(w)

	if _, err := bw.Write(HEADER); err != nil {
		return err
	}

	if err := binary.Write(bw, binary.LittleEndian, VERSION); err != nil {
		return err
	}

	if err := binary.Write(bw, binary.LittleEndian, builder.steamID64); err != nil {
		return err
	}

	unixTime := uint64(time.Now().Unix())
	if err := binary.Write(bw, binary.LittleEndian, unixTime); err != nil {
		return err
	}

	// required content (unused)
	if err := binary.Write(bw, binary.LittleEndian, uint8(0)); err != nil {
		return err
	}

	if err := writeCString(bw, builder.name); err != nil {
		return err
	}

	if err := writeCString(bw, builder.description); err != nil {
		return err
	}

	if err := writeCString(bw, builder.author); err != nil {
		return err
	}

	// version (unused)
	if err := binary.Write(bw, binary.LittleEndian, int32(1)); err != nil {
		return err
	}

	// Write metadata for each file entry.
	for i, e := range builder.entries {
		// Write the file index (1-based).
		if err := binary.Write(bw, binary.LittleEndian, uint32(i+1)); err != nil {
			return err
		}
		// Write the entry name as a C string.
		if err := writeCString(bw, e.name); err != nil {
			return err
		}
		// Write the file size.
		if err := binary.Write(bw, binary.LittleEndian, int64(len(e.content))); err != nil {
			return err
		}
		// Write the CRC (unused by the game, so we write zero).
		if err := binary.Write(bw, binary.LittleEndian, uint32(0)); err != nil {
			return err
		}
	}

	// Write a zero to indicate the end of the metadata.
	if err := binary.Write(bw, binary.LittleEndian, uint32(0)); err != nil {
		return err
	}

	// Write the content of each file entry.
	for _, e := range builder.entries {
		if _, err := bw.Write(e.content); err != nil {
			return err
		}
	}

	// Write a zero to indicate the end of the file.
	if err := binary.Write(bw, binary.LittleEndian, uint32(0)); err != nil {
		return err
	}

	// Flush the buffered writer
	if err := bw.Flush(); err != nil {
		return err
	}

	return nil
}

func writeCString(w io.Writer, s string) error {
	if strings.ContainsRune(s, 0) {
		return fmt.Errorf("string contains null byte")
	}

	if _, err := io.WriteString(w, s); err != nil {
		return err
	}

	return binary.Write(w, binary.LittleEndian, uint8(0))
}
