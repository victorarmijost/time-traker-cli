package bairestt

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"time"

	"github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/cdproto/target"
	"github.com/chromedp/chromedp"
)

//Emulate web browser login to get access token
func (t *Bairestt) emulate_login(ctx context.Context, password string) (string, error) {
	if t.email == "" {
		return "", fmt.Errorf("missing email")
	}

	//Remove unwanted errors printing
	tErr := os.Stderr
	tOut := os.Stdout
	os.Stderr = nil
	os.Stdout = nil

	defer func() {
		os.Stderr = tErr
		os.Stdout = tOut
	}()

	ctx, cancel := chromedp.NewContext(ctx)
	defer cancel()

	// Grab the first spawned tab that isn't blank. Used to catch the login popup
	ch := chromedp.WaitNewTarget(ctx, func(info *target.Info) bool {
		return info.URL != ""
	})

	//Starts the main page
	var nodes []*cdp.Node
	if err := chromedp.Run(ctx,
		chromedp.Navigate("https://employees.bairesdev.com/login"),
		chromedp.WaitReady(`Sign in with Google`, chromedp.BySearch),
		chromedp.Nodes(`Sign in with Google`, &nodes, chromedp.BySearch),
	); err != nil {
		return "", fmt.Errorf("error loading initial page, %w", err)
	}

	//Search the login with google node
	var login_node *cdp.Node
	for _, node := range nodes {
		if node.NodeValue == "Sign in with Google" {
			login_node = node
		}
	}

	//Click the login
	if err := chromedp.Run(ctx,
		chromedp.MouseClickNode(login_node),
	); err != nil {
		return "", fmt.Errorf("error clicking login button, %w", err)
	}

	//Waits for the popup and attach it
	newCtx, cancelNew := chromedp.NewContext(ctx, chromedp.WithTargetID(<-ch))
	defer cancelNew()

	//Handle the login form
	var buf []byte
	if err := chromedp.Run(newCtx,
		chromedp.WaitReady(`Correo electrónico o teléfono`, chromedp.BySearch),
		chromedp.SendKeys(`Correo electrónico o teléfono`, t.email, chromedp.BySearch),
		chromedp.Click(`Siguiente`, chromedp.BySearch, chromedp.NodeVisible),
		chromedp.WaitReady(`Enter your password`, chromedp.BySearch),
		chromedp.SendKeys(`Passwd`, password, chromedp.BySearch),
		chromedp.Submit(`submit`, chromedp.ByID, chromedp.NodeVisible),
		chromedp.WaitNotPresent(`submit`, chromedp.ByID),
		chromedp.WaitReady(`submit_approve_access`, chromedp.ByID),
		chromedp.Sleep(2*time.Second),
		chromedp.FullScreenshot(&buf, 100),
		chromedp.Click(`submit_approve_access`, chromedp.ByID),
	); err != nil {
		return "", fmt.Errorf("error handling login popup, %w", err)
	}

	if t.debug {
		if err := ioutil.WriteFile("login-form.png", buf, 0o644); err != nil {
			return "", fmt.Errorf("error saving screenshoot, %w", err)
		}
	}

	var buf2 []byte
	var res string
	if err := chromedp.Run(ctx,
		chromedp.WaitReady(`Track your hours`, chromedp.BySearch),
		chromedp.Evaluate(`localStorage.bairesdev_token`, &res),
		chromedp.CaptureScreenshot(&buf2),
	); err != nil {
		return "", fmt.Errorf("error getting token, %w", err)
	}

	if t.debug {
		if err := ioutil.WriteFile("main.png", buf2, 0o644); err != nil {
			return "", fmt.Errorf("error saving screenshoot, %w", err)
		}
	}

	return res, nil
}
