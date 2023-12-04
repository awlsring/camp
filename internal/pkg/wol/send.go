package wol

import (
	"context"
	"fmt"

	"github.com/sabhiram/go-wol/wol"
)

func (w *wakeOnLanClient) SendSignal(ctx context.Context, mac string) error {
	packet, err := wol.New(mac)
	if err != nil {
		return err
	}

	bytes, err := packet.Marshal()
	if err != nil {
		return err
	}

	n, err := w.conn.Write(bytes)
	if err != nil {
		return err
	}
	if n != 102 {
		err = fmt.Errorf("sending magic packet returned %d bytes, howerver expected 102 bytes as response).", n)
	}

	return nil
}
