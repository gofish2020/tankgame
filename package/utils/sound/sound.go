package sound

import (
	"bytes"
	_ "embed"
	"io"
	"log"

	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/audio/mp3"
	"github.com/hajimehoshi/ebiten/v2/audio/wav"
)

var (
	currentplayer *audio.Player = nil

	bgmPlayer *audio.Player = nil

	//go:embed boom.wav
	boom []byte

	//go:embed bgm.mp3
	bgm []byte

	mSound map[string][]byte
)

func init() {
	mSound = make(map[string][]byte)
}

func LoadSound() {

	boomStream, err := wav.DecodeWithSampleRate(44100, bytes.NewReader(boom))
	if err != nil {
		log.Fatal(err)
	}
	// 这里的data 应该是pcm原始数据
	data, err := io.ReadAll(boomStream)
	if err != nil {
		log.Fatal(err)
	}

	mSound["boom"] = data

	bgmStream, err := mp3.DecodeWithSampleRate(44100, bytes.NewReader(bgm))
	if err != nil {
		log.Fatal(err)
	}
	// 这里的data 应该是pcm原始数据
	data, err = io.ReadAll(bgmStream)
	if err != nil {
		log.Fatal(err)
	}

	mSound["bgm"] = data
}

func PlaySound(s string) {

	pcm, ok := mSound[s]
	if ok {
		if currentplayer != nil && currentplayer.IsPlaying() {
			currentplayer.Close()
		}
		currentplayer = audio.CurrentContext().NewPlayerFromBytes(pcm)
		currentplayer.SetVolume(.5)
		currentplayer.Play()
	}
}

func PlayBGM() {

	if bgmPlayer != nil && bgmPlayer.IsPlaying() {
		return
	}

	if bgmPlayer != nil {
		bgmPlayer.Close()
	}
	bgmPlayer = audio.CurrentContext().NewPlayerFromBytes(mSound["bgm"])
	bgmPlayer.SetVolume(.2)
	bgmPlayer.Play()
}
