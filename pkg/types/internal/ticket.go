package internal

import (
	"bytes"
	"fmt"
	"math/big"

	"github.com/minio/blake2b-simd"
)

// A Ticket is a marker of a tick of the blockchain's clock.  It is the source
// of randomness for proofs of storage and leader election.  It is generated
// by the miner of a newBlock using a VRF.
type Ticket struct {
	// A proof output by running a VRF on the VRFProof of the parent ticket
	VRFProof VRFPi
}

// String returns the string representation of the VRFProof of the ticket
func (t Ticket) String() string {
	return fmt.Sprintf("%x", t.VRFProof)
}

func (t *Ticket) Compare(o *Ticket) int {
	tDigest := t.VRFProof.Digest()
	oDigest := o.VRFProof.Digest()
	return bytes.Compare(tDigest[:], oDigest[:])
}

func (t *Ticket) Less(o *Ticket) bool {
	return t.Compare(o) < 0
}

func (t *Ticket) Quality() float64 {
	ticketHash := blake2b.Sum256(t.VRFProof)
	ticketNum := BigFromBytes(ticketHash[:]).Int
	ticketDenu := big.NewInt(1)
	ticketDenu.Lsh(ticketDenu, 256)
	tv, _ := new(big.Rat).SetFrac(ticketNum, ticketDenu).Float64()
	tq := 1 - tv
	return tq
}
