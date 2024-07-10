package sound

import (
	"bytes"
	"embed"
	_ "embed"
	"io"
	"io/fs"
	"log"
	"path/filepath"
	"strings"

	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/audio/mp3"
	"github.com/hajimehoshi/ebiten/v2/audio/wav"
)

var (
	currentplayer *audio.Player = nil

	bgmPlayer *audio.Player = nil

	mSound map[string][]byte

	//go:embed *
	f embed.FS
)

func init() {
	mSound = make(map[string][]byte)
	loadSound()
}

// 加载音频文件
func loadSound() {

	a, _ := fs.ReadDir(f, ".")

	for _, v := range a {
		// 读取文件内容
		data, _ := f.ReadFile(v.Name())

		// 去掉文件名后缀（只剩下文件名）
		name := strings.TrimSuffix(v.Name(), filepath.Ext(v.Name()))
		// 文件后缀
		ext := filepath.Ext(v.Name())

		switch ext {
		case ".mp3":

			mp3Stream, err := mp3.DecodeWithSampleRate(44100, bytes.NewReader(data))
			if err != nil {
				log.Fatal(err)
			}
			// 这里的data 应该是pcm原始数据
			data, err := io.ReadAll(mp3Stream)
			if err != nil {
				log.Fatal(err)
			}
			mSound[name] = data
		case ".wav":
			stream, err := wav.DecodeWithSampleRate(44100, bytes.NewReader(data))
			if err != nil {
				log.Fatal(err)
			}
			// 这里的data 应该是pcm原始数据
			data, err := io.ReadAll(stream)
			if err != nil {
				log.Fatal(err)
			}

			mSound[name] = data
		}
	}
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
