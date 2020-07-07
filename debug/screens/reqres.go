package screens

import (
	"encoding/json"
	"fmt"
	"fyne.io/fyne"
	"fyne.io/fyne/layout"
	"fyne.io/fyne/theme"
	"fyne.io/fyne/widget"
	"go-skysharing-openapi/debug/conf"
	"go-skysharing-openapi/pkg/cass"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

var (
	reqInput = widget.NewMultiLineEntry()
	resInput = widget.NewMultiLineEntry()
	method   cass.Method
)

func ReqResScreen() fyne.CanvasObject {
	tool := toolbox()
	return fyne.NewContainerWithLayout(
		layout.NewBorderLayout(tool, nil, nil, nil),
		tool,
		fyne.NewContainerWithLayout(
			layout.NewGridLayoutWithRows(3),
			widget.NewSelect(cass.M.GetOptions(), func(name string) {
				method = cass.M.GetOption(name)
			}),
			fyne.NewContainerWithLayout(
				layout.NewGridLayoutWithColumns(2),
				req(),
				res(),
			),
			submit(),
		),
	)
}

func toolbox() *widget.Toolbar {
	return widget.NewToolbar(
		widget.NewToolbarAction(theme.DeleteIcon(), func() {
			reqInput.SetText("")
			resInput.SetText("")
		}),
	)
}

func submit() fyne.CanvasObject {
	b := widget.NewButton("submit", func() {
		var err error
		fmt.Printf("method: %+v \n", method)
		fmt.Printf("request: %+v \n", reqInput.Text)
		fmt.Printf("response: %+v \n", resInput.Text)
		err = func() error {
			req := map[string]interface{}{}
			err = json.Unmarshal([]byte(reqInput.Text), &req)
			if err != nil {
				return err
			}
			r, err := cass.NewRequest(
				conf.C.UserPrivateKey,
				conf.C.AppId,
				"JSON",
				"UTF-8",
				"1.0",
				"RSA2",
			)
			if err != nil {
				return err
			}
			r.BizParam = map[string]interface{}{}
			r.Params.Method = method.Method
			client := http.Client{
				Transport:     nil,
				CheckRedirect: nil,
				Jar:           nil,
				Timeout:       10 * time.Second,
			}
			err = r.BuildParams()
			if err != nil {
				return err
			}
			buildQuery, err := r.Params.BuildQuery()
			if err != nil {
				return err
			}
			post, err := client.Post(conf.C.Uri, "application/html; charset=utf-8", strings.NewReader(buildQuery))
			if err != nil {
				return err
			}
			bites, err := ioutil.ReadAll(post.Body)
			if err != nil {
				return err
			}
			fmt.Printf("response: %+v", string(bites))
			resInput.SetText(string(bites))
			return nil
		}()
		if err != nil {
			fmt.Printf("request err: %+v \n", err)
		}
	})
	return b
}

func req() fyne.CanvasObject {
	reqInput.PlaceHolder = "request"
	return widget.NewVBox(reqInput)
}

func res() fyne.CanvasObject {
	resInput.PlaceHolder = "response"
	return widget.NewVBox(resInput)
}
