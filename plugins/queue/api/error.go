// SiGG-Satellite-Network-SII  //

package api

import "errors"

var (
	ErrEmpty  = errors.New("cannot read data when the queue is empty")
	ErrFull   = errors.New("cannot write data when the queue is full")
	ErrClosed = errors.New("cannot enqueue or dequeue when the queue is closed")
)
