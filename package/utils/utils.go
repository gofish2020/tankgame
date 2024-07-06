package utils

import (
	"image"
	"image/color"
	"image/png"
	"os"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/nfnt/resize"
)

var (
	whiteImage    = ebiten.NewImage(3, 3)
	whiteSubImage = whiteImage.SubImage(image.Rect(1, 1, 2, 2)).(*ebiten.Image)
)

func init() {
	b := whiteImage.Bounds()
	pix := make([]byte, 4*b.Dx()*b.Dy())
	for i := range pix {
		pix[i] = 0xff
	}
	// This is hacky, but WritePixels is better than Fill in term of automatic texture packing.
	whiteImage.WritePixels(pix)
}

func Resize(path string, w, h uint) *image.Image {
	file, err := os.Open(path)
	if err != nil {
		return nil
	}
	defer file.Close()
	img, err := png.Decode(file)
	if err != nil {
		return nil
	}

	m := resize.Resize(w, h, img, resize.Lanczos2)
	out, _ := os.Create("tank.png")
	defer out.Close()
	png.Encode(out, m)
	return &m
}

func DrawSector(screen *ebiten.Image, x, y float32, lineWidth float32, radius float32, startAngle, endAngle float32, clr color.Color, isFill bool) {
	var path vector.Path

	//theta2 := math.Pi * float64(count) / 180 / 3
	path.MoveTo(x, y)
	path.Arc(x, y, radius, startAngle, endAngle, vector.Clockwise)
	path.Close()

	var vs []ebiten.Vertex
	var is []uint16
	if !isFill {
		op := &vector.StrokeOptions{}
		op.Width = lineWidth
		op.LineJoin = vector.LineJoinRound
		vs, is = path.AppendVerticesAndIndicesForStroke(nil, nil, op)
	} else {
		vs, is = path.AppendVerticesAndIndicesForFilling(nil, nil)
	}

	r, g, b, a := clr.RGBA()
	for i := range vs {
		vs[i].SrcX = 1
		vs[i].SrcY = 1
		vs[i].ColorR = float32(r) / 0xffff
		vs[i].ColorG = float32(g) / 0xffff
		vs[i].ColorB = float32(b) / 0xffff
		vs[i].ColorA = float32(a) / 0xffff
	}

	op := &ebiten.DrawTrianglesOptions{}
	op.AntiAlias = true
	op.FillRule = ebiten.FillAll
	screen.DrawTriangles(vs, is, whiteSubImage, op)
}
