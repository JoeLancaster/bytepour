package testdata

import (
	_ "embed"

	"github.com/joelancaster/bytepour/pkg/metainfo"
)

//go:embed debian_pieces.dat
var pieces []byte

//go:embed debian_info_dict.dat
var infodict []byte

var WantDebianMetaInfo = metainfo.MetaInfoPreCompute{
	Announce: []byte("http://bttracker.debian.org:6969/announce"),
	Comment:  []byte(`"Debian CD from cdimage.debian.org"`),
	InfoDict: infodict,
	Info: metainfo.Info{
		Length:      659554304,
		Name:        []byte("debian-12.5.0-amd64-netinst.iso"),
		PieceLength: 262144,
		Pieces:      pieces,
	},
}
