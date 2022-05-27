package bairestt

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"time"
	"varmijo/time-tracker/utils"

	"github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/cdproto/target"
	"github.com/chromedp/chromedp"
	"github.com/sirupsen/logrus"
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

	logrus.Debug("start login process")

	// Grab the first spawned tab that isn't blank. Used to catch the login popup
	ch := chromedp.WaitNewTarget(ctx, func(info *target.Info) bool {
		return info.URL != ""
	})

	logrus.Debug("opening main page")

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

	logrus.Debug("hitting sign with google button")

	//Click the login
	if err := chromedp.Run(ctx,
		chromedp.MouseClickNode(login_node),
	); err != nil {
		return "", fmt.Errorf("error clicking login button, %w", err)
	}

	logrus.Debug("waiting for popup login form")

	var tid target.ID
	select {
	case tid = <-ch:
	case <-time.After(10 * time.Second):
		logrus.Fatal("login popup timeout")
	}

	//Waits for the popup and attach it
	newCtx, cancelNew := chromedp.NewContext(ctx, chromedp.WithTargetID(tid))
	defer cancelNew()

	logrus.Debug("sending email")

	//Handle the login form
	var buf []byte
	if err := chromedp.Run(newCtx,
		chromedp.WaitReady(`Email`, chromedp.ByID),
		chromedp.FullScreenshot(&buf, 100),
	); err != nil {
		return "", fmt.Errorf("error handling login popup email, %w", err)
	}

	if logrus.GetLevel() == logrus.DebugLevel {
		if err := ioutil.WriteFile(utils.GeAppPath("login-form-email.png"), buf, 0o644); err != nil {
			return "", fmt.Errorf("error saving screenshoot, %w", err)
		}
	}

	if err := chromedp.Run(newCtx,
		chromedp.SendKeys(`Email`, t.email, chromedp.ByID),
		chromedp.Click(`next`, chromedp.ByID, chromedp.NodeVisible),
	); err != nil {
		return "", fmt.Errorf("error handling login popup email, %w", err)
	}

	logrus.Debug("sending password")

	//Handle the login form
	if err := chromedp.Run(newCtx,
		chromedp.WaitReady(`password`, chromedp.ByID),
		chromedp.FullScreenshot(&buf, 100),
	); err != nil {
		return "", fmt.Errorf("error handling login popup password, %w", err)
	}

	if logrus.GetLevel() == logrus.DebugLevel {
		if err := ioutil.WriteFile(utils.GeAppPath("login-form-password.png"), buf, 0o644); err != nil {
			return "", fmt.Errorf("error saving screenshoot, %w", err)
		}
	}

	if err := chromedp.Run(newCtx,
		chromedp.SendKeys(`password`, password, chromedp.ByID),
		chromedp.Submit(`submit`, chromedp.ByID, chromedp.NodeVisible),
		chromedp.WaitNotPresent(`submit`, chromedp.ByID),
		chromedp.WaitReady(`submit_approve_access`, chromedp.ByID),
		chromedp.Sleep(2*time.Second),
		chromedp.Click(`submit_approve_access`, chromedp.ByID),
	); err != nil {
		return "", fmt.Errorf("error handling login popup password, %w", err)
	}

	logrus.Debug("getting token after login")

	var buf2 []byte
	var res string
	if err := chromedp.Run(ctx,
		chromedp.WaitReady(`Track your hours`, chromedp.BySearch),
		chromedp.Evaluate(`localStorage.bairesdev_token`, &res),
		chromedp.CaptureScreenshot(&buf2),
	); err != nil {
		return "", fmt.Errorf("error getting token, %w", err)
	}

	if logrus.GetLevel() == logrus.DebugLevel {
		if err := ioutil.WriteFile(utils.GeAppPath("main.png"), buf2, 0o644); err != nil {
			return "", fmt.Errorf("error saving screenshoot, %w", err)
		}
	}

	logrus.Debug("login completed!")

	return res, nil
}
