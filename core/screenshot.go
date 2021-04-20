// Command screenshot is a chromedp example demonstrating how to take a
// screenshot of a specific element and of the entire browser viewport.
package core

import (
	"context"
	"io/ioutil"
	"log"
	"math"

	"github.com/chromedp/cdproto/emulation"
	"github.com/chromedp/cdproto/page"
	"github.com/chromedp/chromedp"
)

func Screenshot() {
	// create context
	ctx, cancel := chromedp.NewContext(
		context.Background(),
		// chromedp.WithDebugf(log.Printf),
	)
	defer cancel()

	// capture screenshot of an element
	var buf []byte
	if err := chromedp.Run(ctx, elementScreenshot(`https://pkg.go.dev/`, `img.Homepage-logo`, &buf)); err != nil {
		log.Fatal(err)
	}
	if err := ioutil.WriteFile("elementScreenshot.png", buf, 0o644); err != nil {
		log.Fatal(err)
	}

	// capture entire browser viewport, returning png with quality=90
	if err := chromedp.Run(ctx, fullScreenshot(`https://www.monitorapp.com/`, 90, &buf)); err != nil {
		log.Fatal(err)
	}
	if err := ioutil.WriteFile("fullScreenshot.png", buf, 0o644); err != nil {
		log.Fatal(err)
	}

	log.Printf("wrote elementScreenshot.png and fullScreenshot.png")
}

// elementScreenshot takes a screenshot of a specific element.
func elementScreenshot(urlstr, sel string, res *[]byte) chromedp.Tasks {
	return chromedp.Tasks{
		chromedp.Navigate(urlstr),
		chromedp.Screenshot(sel, res, chromedp.NodeVisible),
	}
}

// fullScreenshot takes a screenshot of the entire browser viewport.
//
// Liberally copied from puppeteer's source.
//
// Note: this will override the viewport emulation settings.
func fullScreenshot(urlstr string, quality int64, res *[]byte) chromedp.Tasks {
	return chromedp.Tasks{
		chromedp.Navigate(urlstr),
		chromedp.ActionFunc(func(ctx context.Context) error {
			// get layout metrics
			// 컨텐츠 만큼의 크기만 구해옴 *dom.Rect
			_, _, cssContentSize, err := page.GetLayoutMetrics().Do(ctx)
			if err != nil {
				return err
			}
			// cssContentSize의 최소 값을 구해서 width, height에 넣는 과정
			width, height := int64(math.Ceil(cssContentSize.Width)), int64(math.Ceil(cssContentSize.Height))

			// force viewport emulation
			// 뷰포트 강제로 재정의하는 부분
			// 만약 컨텐츠만큼의 크기가 아니라, 픽셀을 지정해서 스크린샷 찍고 싶으면, SetDeviceMetricsOverride의 width, height에 픽셀 값 넣으면 됨
			err = emulation.SetDeviceMetricsOverride(width, height, 1, false).
				WithScreenOrientation(&emulation.ScreenOrientation{
					Type:  emulation.OrientationTypePortraitPrimary,
					Angle: 0,
				}).
				Do(ctx)
			if err != nil {
				return err
			}

			// capture screenshot
			// 만약 컨텐츠만큼의 크기가 아니라, 픽셀을 지정해서 스크린샷 찍고 싶으면, SetDeviceMetricsOverride의 width, height에 픽셀 값 넣으면 됨
			*res, err = page.CaptureScreenshot().
				WithQuality(quality).
				WithClip(&page.Viewport{
					X:      cssContentSize.X,
					Y:      cssContentSize.Y,
					Width:  cssContentSize.Width,
					Height: cssContentSize.Height,
					Scale:  1,
				}).Do(ctx)
			if err != nil {
				return err
			}
			return nil
		}),
	}
}
