package main

import (
	"encoding/xml"
	"io"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/kalginnick/go-lambda-talk/pkg/client"
)

func main() {
	h := handler{
		region: os.Getenv("AWS_REGION"),
		bucket: os.Getenv("BUCKET"),
		file:   os.Getenv("FILENAME"),
	}

	lambda.Start(h.handle)
}

type handler struct {
	bucket string
	file   string
	region string
}

func (h *handler) handle(event events.S3Event) error {
	for _, record := range event.Records {
		err := func() error {
			file, err := client.ReadS3(record.AWSRegion, record.S3.Bucket.Name, record.S3.Object.Key)
			if err != nil {
				return err
			}
			defer file.Close()

			return client.WriteS3(h.region, h.bucket, h.file, transform(file))
		}()
		if err != nil {
			return err
		}
	}

	return nil
}

/*
<entry>
	<SKU>4G680</SKU>
	<Description>Tenda 4G LTE 300Mbps WiFi Router | 4G680</Description>
	<CPT>1</CPT>
	<JHB>2</JHB>
	<DBN>3</DBN>
	<TotalStock>6</TotalStock>
	<DealerPrice>995</DealerPrice>
	<RetailPrice>1252.1700</RetailPrice>
	<Manufacturer>Tenda</Manufacturer>
	<ImageURL>https://scoop.co.za/download/marketing/images/4G680.jpg</ImageURL>
</entry>
*/
type entry struct {
	XMLName           xml.Name `xml:"entry"`
	SKU               string   `xml:"SKU"`
	Description       string   `xml:"Description"`
	CapeTownStock     int      `xml:"CPT"`
	JohannesburgStock int      `xml:"JHB"`
	DurbanStock       int      `xml:"DBN"`
	TotalStock        int      `xml:"TotalStock"`
	DealerPrice       float64  `xml:"DealerPrice"`
	RetailPrice       float64  `xml:"RetailPrice"`
	Manufacturer      string   `xml:"Manufacturer"`
	ImageURL          string   `xml:"ImageURL"`
}

/*
<offer ref="4G680">
	<name></name>
	<description>Tenda 4G LTE 300Mbps WiFi Router | 4G680</description>
	<price>1252.1700</price>
	<stock>2</stock>
	<picture>https://scoop.co.za/download/marketing/images/4G680.jpg</picture>
	<manufacturer>Tenda</manufacturer>
</offer>
*/
type offer struct {
	XMLName      xml.Name `xml:"offer"`
	Ref          string   `xml:"ref,attr"`
	Name         string   `xml:"name"`
	Description  string   `xml:"description"`
	Price        float64  `xml:"price"`
	Stock        int      `xml:"stock"`
	Picture      string   `xml:"picture"`
	Manufacturer string   `xml:"manufacturer"`
}

func (e entry) ToOffer() offer {
	return offer{
		XMLName:      xml.Name{},
		Ref:          e.SKU,
		Description:  e.Description,
		Price:        e.RetailPrice,
		Stock:        e.JohannesburgStock,
		Picture:      e.ImageURL,
		Manufacturer: e.Manufacturer,
	}
}

func transform(r io.Reader) io.Reader {
	pr, pw := io.Pipe()
	from := xml.NewDecoder(r)
	to := xml.NewEncoder(pw)
	go func() {
		defer pw.Close()

		to.EncodeToken(xml.ProcInst{Target: "xml", Inst: []byte(`version="1.0" encoding="utf-8"`)})
		to.EncodeToken(xml.CharData("\n"))
		to.EncodeToken(xml.StartElement{Name: xml.Name{Local: "offers"}})
		to.Flush()

		for {
			token, err := from.Token()
			if err != nil {
				break
			}
			switch element := token.(type) {
			case xml.StartElement:
				if element.Name.Local == "entry" {
					item := entry{}
					err = from.DecodeElement(&item, &element)
					if err != nil {
						pw.CloseWithError(err)
					}
					to.Encode(item.ToOffer())
				}
			}
		}

		to.EncodeToken(xml.EndElement{Name: xml.Name{Local: "offers"}})
		to.Flush()
	}()
	return pr
}
