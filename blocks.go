package geoip

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"sort"
	"strconv"
)

// Block represents block IPs with their associated location id
type Block struct {
	// StartIPNum is the starting IP number for a block.
	StartIPNum int
	// EndIPNum is the ending IP number for a block.
	EndIPNum int
	// LocationID is the locations id for lookup in LocationMap.
	LocationID int
}

// NewBlock initializes a blocks ip mapping to location id.
func NewBlock(record []string) (*Block, error) {
	var err error
	b := &Block{}
	b.StartIPNum, err = strconv.Atoi(record[0])
	if err != nil {
		return nil, err
	}
	b.EndIPNum, err = strconv.Atoi(record[1])
	if err != nil {
		return nil, err
	}
	b.LocationID, err = strconv.Atoi(record[2])
	if err != nil {
		return nil, err
	}
	return b, nil
}

// BlocksMap holds IP to location id information. It is a
// list of range boundaries
type BlocksMap struct {
	ranges []int
	blocks []*Block
}

// NewBlocksMap will load a blocks csv from file and return a
// mapping of ip to location id.
func NewBlocksMap(f io.Reader) (*BlocksMap, error) {
	bm := &BlocksMap{
		ranges: []int{-1, 16777215},
		blocks: []*Block{nil, nil},
	}
	r := csv.NewReader(f)
	r.FieldsPerRecord = 3
	fmt.Println("Blocks:")
	i := -1
READ_LOOP:
	for {
		i++
		record, err := r.Read()
		if err == io.EOF {
			break READ_LOOP
		}
		if err != nil {
			if perr, ok := err.(*csv.ParseError); ok && perr.Err == csv.ErrFieldCount {
				if i != 0 {
					log.Println(err)
				}
				continue READ_LOOP
			}
			return nil, err
		}
		if len(record) == 3 {
			b, err := NewBlock(record)
			if err != nil {
				log.Println(err)
				continue READ_LOOP
			}
			bm.ranges = append(bm.ranges, b.EndIPNum)
			bm.blocks = append(bm.blocks, b)
			fmt.Printf("\r %v", i)
		}
	}
	fmt.Printf("\ndone\n\n")
	return bm, nil
}

// Lookup will lookup a location id by ip address.
func (l *BlocksMap) Lookup(ip string) (int, bool) {
	ipInt, err := ipToInt(ip)
	if err != nil {
		return 0, false
	}
	idx := sort.SearchInts(l.ranges, ipInt)
	if idx >= len(l.blocks) {
		// returned out of range
		return -1, false
	}
	b := l.blocks[idx]
	if b == nil {
		return -1, false
	}
	return b.LocationID, true
}
