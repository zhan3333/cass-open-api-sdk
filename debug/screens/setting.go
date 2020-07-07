package screens

import (
	"fmt"
	"fyne.io/fyne"
	"fyne.io/fyne/dialog"
	"fyne.io/fyne/widget"
	"go-skysharing-openapi/debug/conf"
)

func SettingScreen(w fyne.Window) fyne.CanvasObject {
	return widget.NewTabContainer(
		widget.NewTabItem("Form", func() fyne.CanvasObject {
			uri := widget.NewEntry()
			uri.SetPlaceHolder("api uri")
			uri.SetText(conf.C.Uri)

			appId := widget.NewEntry()
			appId.SetPlaceHolder("APP ID")
			appId.SetText(conf.C.AppId)

			sysPubKey := widget.NewMultiLineEntry()
			sysPubKey.SetPlaceHolder("System Public Key")
			sysPubKey.SetText(conf.C.SystemPublicKey)
			sysPubKey.BaseWidget.Resize(sysPubKey.MinSize())
			sysPubKey.MinSize()

			userPubKey := widget.NewMultiLineEntry()
			userPubKey.SetPlaceHolder("User Public Key")
			userPubKey.SetText(conf.C.UserPublicKey)

			userPriKey := widget.NewMultiLineEntry()
			userPriKey.SetPlaceHolder("User Private Key")
			userPriKey.SetText(conf.C.UserPrivateKey)

			form := &widget.Form{
				OnCancel: func() {
					fmt.Println("Cancelled")
				},
				OnSubmit: func() {
					err := conf.Set(uri.Text, appId.Text, sysPubKey.Text, userPubKey.Text, userPriKey.Text)
					if err != nil {
						dialog.ShowError(err, w)
					} else {
						dialog.ShowInformation("ok", "Save success", w)
					}
					fmt.Printf("Form submitted: %+v", err)
				},
			}
			form.Append("URI", uri)
			form.Append("APP ID", appId)
			form.Append("System Public Key", sysPubKey)
			form.Append("User Public Key", userPubKey)
			form.Append("User Private Key", userPriKey)
			return form
		}()),
	)
}
