package main

import (
	"context"
	"fmt"
	"log"
	"path"

	"github.com/skip2/go-qrcode"
	"github.com/spf13/pflag"
)

const version = "v1.1.0"

func main() {
	if err := run(context.Background()); err != nil {
		log.Fatal(err)
	}
}

func run(ctx context.Context) error {
	fmt.Printf("Welcome To ZAPS. This version is %s\n", version)

	slug := pflag.StringP("slug", "s", "", "Custom slug for your shorten link")
	qr := pflag.BoolP("qr", "q", false, "Generate a QR Code for the shorten link")
	saveQR := pflag.BoolP("save", "S", false, "Save the QR Code to a file")
	pflag.Parse()

	originalLink := pflag.Arg(0)

	shURL, err := sendLinkToCreateLink(originalLink, *slug)
	if err != nil {
		return err
	}
	fmt.Printf("This is your shorten link: %s\n", shURL)

	if *qr || *saveQR {
		qrc, err := qrcode.New(shURL, qrcode.Medium)
		if err != nil {
			return err
		}

		if *qr {
			fmt.Println("This is your QR Code:")
			fmt.Println(qrc.ToSmallString(true))
		}

		if *saveQR {
			fileName := path.Base(shURL) + ".png"
			if err := qrc.WriteFile(256, fileName); err != nil {
				return err
			}
			fmt.Printf("The QR Code has been saved as %s\n", fileName)
		}
	}
	return nil
}
