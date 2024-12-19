package gui

import (
	"idz_ais/internal/pc"
	"idz_ais/internal/service"
	"image"
	"image/color"
	"strings"

	"gioui.org/app"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
)

func Loop(w *app.Window) error {
	th := material.NewTheme()
	var ops op.Ops
	var startIp widget.Editor
	var endIp widget.Editor
	var outputText string
	var scanButton widget.Clickable
	var aliveHosts []string
	var selectedHost widget.Editor
	var userNameHost widget.Editor
	var passwordHost widget.Editor
	var userInfoButton widget.Clickable
	var processInfoButton widget.Clickable
	var shutdownButton widget.Clickable
	var rebootButton widget.Clickable
	var selectedPC *pc.PC
	var outPutCommandText string
	var scroll widget.List
	scroll.Axis = layout.Vertical

	for {
		switch e := w.Event().(type) {
		case app.DestroyEvent:
			return e.Err
		case app.FrameEvent:
			gtx := app.NewContext(&ops, e)

			// Обработка нажатия кнопки "Просканировать сеть"
			if scanButton.Clicked(gtx) {
				startIpText := startIp.Text()
				endIpText := endIp.Text()
				aliveHosts = service.ScanNetwork(startIpText, endIpText)
				outputText = strings.Join(aliveHosts, "\n")
			}

			// Обработка нажатия кнопок для выбранного хоста
			if textSelectedHost := selectedHost.Text(); textSelectedHost != "" {
				if userInfoButton.Clicked(gtx) {
					selectedPC = pc.NewPC(textSelectedHost, userNameHost.Text(), passwordHost.Text())
					outPutCommandText = selectedPC.GetUserInfo()
				}
				if processInfoButton.Clicked(gtx) {
					selectedPC = pc.NewPC(textSelectedHost, userNameHost.Text(), passwordHost.Text())
					outPutCommandText = selectedPC.GetProcessInfo()
				}
				if shutdownButton.Clicked(gtx) {
					selectedPC = pc.NewPC(textSelectedHost, userNameHost.Text(), passwordHost.Text())
					outPutCommandText = selectedPC.ShutdownHost()
				}
				if rebootButton.Clicked(gtx) {
					selectedPC = pc.NewPC(textSelectedHost, userNameHost.Text(), passwordHost.Text())
					outPutCommandText = selectedPC.RebootHost()
				}
			}

			// Разметка интерфейса
			layout.Flex{Axis: layout.Vertical}.Layout(gtx,
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					return material.Editor(th, &startIp, "Начальный адрес подсети").Layout(gtx)
				}),
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					return material.Editor(th, &endIp, "Конечный адрес подсети").Layout(gtx)
				}),
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					return material.Button(th, &scanButton, "Просканировать сеть").Layout(gtx)
				}),
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					return material.Label(th, unit.Sp(16), outputText).Layout(gtx)
				}),
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					rect := image.Rect(0, 0, gtx.Constraints.Max.X, 2)
					paint.FillShape(gtx.Ops, color.NRGBA{R: 128, G: 128, B: 128, A: 255}, clip.Rect(rect).Op())
					return layout.Dimensions{Size: rect.Max}
				}),
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
						layout.Rigid(func(gtx layout.Context) layout.Dimensions {
							return material.Editor(th, &selectedHost, "Ip адрес ПК").Layout(gtx)
						}),
						layout.Rigid(func(gtx layout.Context) layout.Dimensions {
							return material.Editor(th, &userNameHost, "Имя пользователя ПК").Layout(gtx)
						}),
						layout.Rigid(func(gtx layout.Context) layout.Dimensions {
							return material.Editor(th, &passwordHost, "Пароль").Layout(gtx)
						}),
						layout.Rigid(func(gtx layout.Context) layout.Dimensions {
							return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
								layout.Rigid(func(gtx layout.Context) layout.Dimensions {
									return material.Button(th, &userInfoButton, "Пользователь").Layout(gtx)
								}),
								layout.Rigid(func(gtx layout.Context) layout.Dimensions {
									return material.Button(th, &processInfoButton, "Процессы").Layout(gtx)
								}),
								layout.Rigid(func(gtx layout.Context) layout.Dimensions {
									return material.Button(th, &shutdownButton, "Выключить").Layout(gtx)
								}),
								layout.Rigid(func(gtx layout.Context) layout.Dimensions {
									return material.Button(th, &rebootButton, "Перезагрузить").Layout(gtx)
								}),
							)
						}),
						layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
							return material.List(th, &scroll).Layout(gtx, 1, func(gtx layout.Context, index int) layout.Dimensions {
								return material.Label(th, unit.Sp(16), outPutCommandText).Layout(gtx)
							})
						}),
					)
				}),
			)
			e.Frame(gtx.Ops)
		}
	}
}
