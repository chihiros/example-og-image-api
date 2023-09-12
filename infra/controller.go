package infra

import (
	"image"
	"image/color"
	"net/http"
	"strings"

	"github.com/fogleman/gg"
	"github.com/nfnt/resize"
	"golang.org/x/image/font/opentype"
)

type Controller interface {
	GenOgImage(w http.ResponseWriter, r *http.Request)
}

type controller struct{}

func NewController() Controller {
	return &controller{}
}

func (c *controller) GenOgImage(w http.ResponseWriter, r *http.Request) {
	title := r.URL.Query().Get("title")
	userName := r.URL.Query().Get("user")

	// 1200x630の画像を生成
	dc := gg.NewContext(1200, 630)

	// 左上から右下に向かってグラデーション
	grad := gg.NewLinearGradient(0, 0, 1200, 630)
	grad.AddColorStop(0, color.RGBA{255, 197, 193, 255})
	grad.AddColorStop(0.25, color.RGBA{244, 222, 244, 255})
	grad.AddColorStop(0.6943, color.RGBA{255, 249, 195, 255})
	grad.AddColorStop(1, color.RGBA{206, 249, 255, 255})
	dc.SetFillStyle(grad)

	// 画像全体の矩形を描画してグラデーションを適用
	dc.DrawRectangle(0, 0, 1200, 630)
	dc.Fill()

	// 図形のサイズと位置を計算
	rectWidth := 1200 - 2*43
	rectHeight := 630 - 2*41
	rectX := 43
	rectY := 41

	// 背景色を設定
	dc.SetColor(color.RGBA{255, 255, 255, 255})
	dc.DrawRoundedRectangle(float64(rectX), float64(rectY), float64(rectWidth), float64(rectHeight), 16)
	dc.Fill()

	// フォントを読み込む
	fontFace, err := opentype.Parse(fontTitle)
	if err != nil {
		http.Error(w, "Failed to parse font", http.StatusInternalServerError)
		return
	}

	// フォントのWeightはWeightMediumにする
	face, err := opentype.NewFace(fontFace, &opentype.FaceOptions{
		Size: 64,
		DPI:  72,
	})

	if err != nil {
		http.Error(w, "Failed to create font face", http.StatusInternalServerError)
		return
	}
	dc.SetFontFace(face)

	// 文字を挿入
	dc.SetRGB(0, 0, 0) // 文字色を黒に設定

	maxWidth := 910.0
	formatTitle := ""

	tmp := 0.0
	for _, word := range title {
		fw, _ := dc.MeasureString(string(word))
		if tmp+fw > maxWidth {
			formatTitle += "\n"
			tmp = 0.0
		}

		formatTitle += string(word)
		tmp += fw
	}

	x := 145.0
	y := 175.0
	for _, line := range strings.Split(formatTitle, "\n") {
		dc.DrawString(line, x, y)
		y += 64
	}

	// サイトのロゴを挿入
	logoImg, _, err := image.Decode(strings.NewReader(string(logo)))
	if err != nil {
		http.Error(w, "Failed to decode logo image", http.StatusInternalServerError)
		return
	}

	// ロゴのサイズを変更
	resizedLogoImg := resize.Resize(0, 150, logoImg, resize.Lanczos3)

	// ロゴを挿入
	dc.DrawImage(resizedLogoImg, 800, 430)

	// フォントを読み込む
	fontFace, err = opentype.Parse(fontUserName)
	if err != nil {
		http.Error(w, "Failed to parse font", http.StatusInternalServerError)
		return
	}

	face, _ = opentype.NewFace(fontFace, &opentype.FaceOptions{
		Size: 48,
		DPI:  72,
	})
	dc.SetFontFace(face)

	// 投稿者の名前を挿入
	dc.SetRGB(0, 0, 0) // 文字色を黒に設定

	x = 160.0
	y = 520.0
	for _, line := range strings.Split(userName, "\n") {
		dc.DrawString(line, x, y)
		y += 48
	}

	// 画像をレスポンスとして返す
	w.Header().Set("Content-Type", "image/png")
	dc.EncodePNG(w)
}
