package car

import (
	"bytes"
	"io"
	"os"
	"testing"

	"github.com/ipfs/go-cid"
	cidlink "github.com/ipld/go-ipld-prime/linking/cid"
	"github.com/storacha/go-ucanto/core/ipld"
)

type fixture struct {
	path   string
	root   ipld.Link
	blocks []ipld.Link
}

var fixtures = []fixture{
	{
		path: "testdata/lost-dog.jpg.car",
		root: cidlink.Link{Cid: cid.MustParse("bafybeif4owy5gno5lwnixqm52rwqfodklf76hsetxdhffuxnplvijskzqq")},
		blocks: []ipld.Link{
			cidlink.Link{Cid: cid.MustParse("bafkreifau35r7vi37tvbvfy3hdwvgb4tlflqf7zcdzeujqcjk3rsphiwte")},
			cidlink.Link{Cid: cid.MustParse("bafkreicj3ozpzd46nx26hflpoi6hgm5linwo65cvphd5ol3ke3vk5nb7aa")},
			cidlink.Link{Cid: cid.MustParse("bafybeihkqv2ukwgpgzkwsuz7whmvneztvxglkljbs3zosewgku2cfluvba")},
			cidlink.Link{Cid: cid.MustParse("bafybeif4owy5gno5lwnixqm52rwqfodklf76hsetxdhffuxnplvijskzqq")},
		},
	},
}

func TestDecodeCAR(t *testing.T) {
	file, err := os.Open(fixtures[0].path)
	if err != nil {
		t.Fatal(err)
	}

	roots, blocks, err := Decode(file)
	if err != nil {
		t.Fatal(err)
	}
	if len(roots) != 1 {
		t.Fatalf("unexpected number of roots: %d, expected: 1", len(roots))
	}
	if roots[0].String() != fixtures[0].root.String() {
		t.Fatalf("unexpected root: %s, expected: %s", roots[0], fixtures[0].root)
	}

	var blks []CarBlock
	for b, err := range blocks {
		if err != nil {
			t.Fatalf("reading blocks: %s", err)
		}
		cb, ok := b.(CarBlock)
		if !ok {
			t.Fatalf("should have returned a car block")
		}
		blks = append(blks, cb)
	}

	if len(blks) != len(fixtures[0].blocks) {
		t.Fatalf("incorrect number of blocks: %d, expected: %d", len(blks), len(fixtures[0].blocks))
	}
	for i, b := range fixtures[0].blocks {
		if b.String() != blks[i].Link().String() {
			t.Fatalf("unexpected block: %s, expected: %s", b, blks[i].Link())
		}
		// verify offset and length can be used to directly read the block in the CAR file
		file.Seek(int64(blks[i].Offset()), io.SeekStart)
		data := make([]byte, blks[i].Length())
		_, err := file.Read(data)
		if err != nil {
			t.Fatalf("error reading block from raw file")
		}
		hashed, err := blks[i].Link().(cidlink.Link).Cid.Prefix().Sum(data)
		if err != nil {
			t.Fatalf("error hashing block from raw file")
		}

		if hashed.String() != blks[i].Link().String() {
			t.Fatalf("raw read from offset block: %s, expected: %s", hashed, blks[i].Link())
		}

	}
}

func TestEncodeCAR(t *testing.T) {
	file, err := os.Open(fixtures[0].path)
	if err != nil {
		t.Fatal(err)
	}

	fbytes, err := io.ReadAll(file)
	if err != nil {
		t.Fatal(err)
	}

	roots, blocks, err := Decode(bytes.NewReader(fbytes))
	if err != nil {
		t.Fatal(err)
	}

	rd := Encode(roots, blocks)

	dbytes, err := io.ReadAll(rd)
	if err != nil {
		t.Fatal(err)
	}

	if !bytes.Equal(fbytes, dbytes) {
		t.Fatal("failed to round trip")
	}
}
