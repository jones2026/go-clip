package main

import (
	"log"
	"os"
	"time"

	"github.com/playwright-community/playwright-go"
)

func main() {
	clipCoupons()
}

func clipCoupons() {
	err := playwright.Install()
	if err != nil {
		log.Fatalf("could not install playwright: %v", err)
	}
	pw, err := playwright.Run()
	if err != nil {
		log.Fatalf("could not launch playwright: %v", err)
	}
	browser, err := pw.WebKit.Launch(playwright.BrowserTypeLaunchOptions{
		Headless: playwright.Bool(false),
	})
	if err != nil {
		log.Fatalf("could not launch Chromium: %v", err)
	}

	page, err := browser.NewPage()
	if err != nil {
		log.Fatalf("could not create page: %v", err)
	}

	if _, err = page.Goto("https://www.hy-vee.com/account/login.aspx?retUrl=/deals/coupons.aspx"); err != nil {
		log.Fatalf("could not goto: %v", err)
	}

	username := os.Getenv("usernamez")
	password := os.Getenv("passwordz")

	page.Type("#username", username)
	page.Type("#password", password)

	page.Click("#__next > main > section > section > form > button > span")
	if _, err = page.Goto("https://www.hy-vee.com/deals/coupons.aspx", playwright.PageGotoOptions{
		WaitUntil: playwright.WaitUntilStateNetworkidle,
	}); err != nil {
		log.Fatalf("could not goto: %v", err)
	}

	coupons, err := page.QuerySelectorAll("#__next > div > div.sc-gsTEea.cgMhGp > main > div.coupon-grid-container > div.coupons-grid > div.coupons-column.coupon-grid > div > div > div")
	if err != nil {
		log.Fatalf("could not get entries: %v", err)
	}
	log.Printf("Found %d coupons to clip\n", len(coupons))
	clippedCount := 0
	for _, coupon := range coupons {
		clipButton, err := coupon.QuerySelector("button > span")
		if err != nil {
			log.Fatalf("could not get loadButton element: %v", err)
		}
		err = clipButton.Click()
		if err == nil {
			clippedCount++
		}
		time.Sleep(100 * time.Millisecond)
	}

	if err = browser.Close(); err != nil {
		log.Fatalf("could not close browser: %v", err)
	}
	if err = pw.Stop(); err != nil {
		log.Fatalf("could not stop Playwright: %v", err)
	}
	log.Printf("Succesfully clipped %d coupons\n", len(coupons))
}
