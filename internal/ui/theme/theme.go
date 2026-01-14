package theme

import (
	"image/color"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
)

type CustomTheme struct{}

var (
	// 主色调 - 粉色系
	PrimaryPink   = color.NRGBA{R: 255, G: 105, B: 180, A: 255} // #FF69B4 热粉色
	SecondaryPink = color.NRGBA{R: 255, G: 20, B: 147, A: 255}  // #FF1493 深粉色
	LightPink     = color.NRGBA{R: 255, G: 182, B: 193, A: 255} // #FFB6C1 浅粉色
	VeryLightPink = color.NRGBA{R: 255, G: 228, B: 225, A: 255} // #FFE4E1 雾粉色

	// 辅助色彩
	White     = color.NRGBA{R: 255, G: 255, B: 255, A: 255}
	SoftGray  = color.NRGBA{R: 248, G: 248, B: 248, A: 255}
	WarmGray  = color.NRGBA{R: 245, G: 245, B: 245, A: 255}
	TextDark  = color.NRGBA{R: 74, G: 74, B: 74, A: 255}
	TextLight = color.NRGBA{R: 136, G: 136, B: 136, A: 255}

	// 装饰色彩
	RoseGold = color.NRGBA{R: 183, G: 110, B: 121, A: 255}
	Lavender = color.NRGBA{R: 230, G: 230, B: 250, A: 255}
	Peach    = color.NRGBA{R: 255, G: 218, B: 185, A: 255}
)

func (m CustomTheme) Color(name fyne.ThemeColorName, variant fyne.ThemeVariant) color.Color {
	switch name {
	case theme.ColorNameBackground:
		return VeryLightPink
	case theme.ColorNameButton:
		return PrimaryPink
	case theme.ColorNameDisabled:
		return LightPink
	case theme.ColorNameDisabledButton:
		return WarmGray
	case theme.ColorNameError:
		return color.NRGBA{R: 255, G: 107, B: 107, A: 255}
	case theme.ColorNameFocus:
		return SecondaryPink
	case theme.ColorNameForeground:
		return TextDark
	case theme.ColorNameHover:
		return LightPink
	case theme.ColorNameInputBackground:
		return White
	case theme.ColorNameInputBorder:
		return LightPink
	case theme.ColorNamePlaceHolder:
		return TextLight
	case theme.ColorNamePrimary:
		return PrimaryPink
	case theme.ColorNameScrollBar:
		return LightPink
	case theme.ColorNameSelection:
		return Lavender
	case theme.ColorNameShadow:
		return color.NRGBA{R: 0, G: 0, B: 0, A: 25}
	default:
		return theme.DefaultTheme().Color(name, variant)
	}
}

func (m CustomTheme) Font(style fyne.TextStyle) fyne.Resource {
	return theme.DefaultTheme().Font(style)
}

func (m CustomTheme) Icon(name fyne.ThemeIconName) fyne.Resource {
	return theme.DefaultTheme().Icon(name)
}

func (m CustomTheme) Size(name fyne.ThemeSizeName) float32 {
	switch name {
	case theme.SizeNameCaptionText:
		return 12
	case theme.SizeNameInlineIcon:
		return 20
	case theme.SizeNamePadding:
		return 8
	case theme.SizeNameScrollBar:
		return 8
	case theme.SizeNameText:
		return 14
	case theme.SizeNameInputBorder:
		return 2
	default:
		return theme.DefaultTheme().Size(name)
	}
}

// 动画配置
type AnimationConfig struct {
	Duration time.Duration
	Curve    fyne.AnimationCurve
}

var (
	// 动画配置
	FadeInAnimation = AnimationConfig{
		Duration: 300 * time.Millisecond,
		Curve:    fyne.AnimationEaseInOut,
	}

	SlideAnimation = AnimationConfig{
		Duration: 400 * time.Millisecond,
		Curve:    fyne.AnimationEaseOut,
	}

	HoverAnimation = AnimationConfig{
		Duration: 200 * time.Millisecond,
		Curve:    fyne.AnimationEaseIn,
	}
)
