package main

import (
	"fmt"
	"log"
	"os"

	"github.com/mxschmitt/playwright-go"
)

func main() {
	clipCoupons()
}

func clipCoupons() {
	pw, err := playwright.Run()
	if err != nil {
		log.Fatalf("could not launch playwright: %v", err)
	}
	browser, err := pw.Chromium.Launch(playwright.BrowserTypeLaunchOptions{
		Headless: playwright.Bool(true),
	})
	if err != nil {
		log.Fatalf("could not launch Chromium: %v", err)
	}

	device := pw.Devices["iPhone 11 Pro"]
	context, err := browser.NewContext(playwright.BrowserNewContextOptions{
		Locale:    playwright.String("en-US"),
		UserAgent: playwright.String(device.UserAgent),
	})
	if err != nil {
		log.Fatalf("could not create context: %v", err)
	}
	page, err := context.NewPage()
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
		WaitUntil: playwright.String("networkidle"),
	}); err != nil {
		log.Fatalf("could not goto: %v", err)
	}

	coupons, err := page.QuerySelectorAll("[id$='_liOffer']")
	if err != nil {
		log.Fatalf("could not get entries: %v", err)
	}
	for i, coupon := range coupons {
		titleElement, err := coupon.QuerySelector("div.product-desc.text-left > span")
		if err != nil {
			log.Fatalf("could not get title element: %v", err)
		}
		title, err := titleElement.TextContent()
		if err != nil {
			log.Fatalf("could not get text content: %v", err)
		}
		clipButton, err := coupon.QuerySelector("[id$='_divLoad'] > a")
		if err != nil {
			log.Fatalf("could not get loadButton element: %v", err)
		}
		clipButton.Click()
		fmt.Printf("%d: %s\n", i+1, title)
	}

	// #spuAvailableCoupons_rptOffers_ctl01_divLoad > span

	if err = browser.Close(); err != nil {
		log.Fatalf("could not close browser: %v", err)
	}
	if err = pw.Stop(); err != nil {
		log.Fatalf("could not stop Playwright: %v", err)
	}
}
