package main

import (
	"fmt"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/gocolly/colly"
)

const amd_6700xt = "6700xt"
const nv_3060 = "3060"
const nv_3060ti = "3060ti"

type gphcrd struct {
	vendor string
	price  float64
	name   string
	desc   string
	url    string
	date   string
}

var graphicCards = []gphcrd{}

func main() {
	ddtechScrap()
	pcelScrap()
	coreGaiming()
	cbpScrap()
	amd6700 := gphcrd{
		price: 99999,
	}
	nv3060 := gphcrd{
		price: 99999,
	}
	nv3060ti := gphcrd{
		price: 99999,
	}

	sort.Slice(graphicCards, func(i, j int) bool {
		return graphicCards[i].name < graphicCards[j].name
	})
	for _, card := range graphicCards {
		fmt.Printf("%s \t%.2f %s\n", card.name, card.price, card.vendor)
		switch {
		case card.name == amd_6700xt && card.price < amd6700.price:
			amd6700 = card
		case card.name == nv_3060 && card.price < nv3060.price:
			nv3060 = card
		case card.name == nv_3060ti && card.price < nv3060ti.price:
			nv3060ti = card

		}
	}

	printCard(amd6700)
	printCard(nv3060ti)
	printCard(nv3060)

}

func printCard(card gphcrd) {
	println("\n-----------------\n")
	fmt.Printf("%s %.2f \n%s \nurl: %s \n", card.name, card.price, card.desc, card.url)
}

func cbpScrap() {
	urls := []string{
		"https://www.cyberpuerta.mx/Computo-Hardware/Componentes/Tarjetas-de-Video/Filtro/Procesador-grafico/GeForce-RTX-3060/Procesador-grafico/Radeon-RX-6700-XT/Procesador-grafico/NVIDIA-GeForce-RTX-3070/Estatus/En-existencia/",
	}
	c := colly.NewCollector()

	// Find and visit all links
	c.OnHTML("div.emproduct_right", func(e *colly.HTMLElement) {
		reg, _ := regexp.Compile("[^0-9]+")
		strPrice := reg.ReplaceAllString(e.ChildText(".price"), "")
		strDel := reg.ReplaceAllString(e.ChildText(".deliveryvalue"), "")
		price, _ := strconv.ParseFloat(strPrice, 32)
		delivery, _ := strconv.ParseFloat(strDel, 32)
		price = price / 100
		delivery = delivery / 100
		desc := e.ChildText(".emproduct_right > a")
		url := e.ChildAttr(".emproduct_right > a", "href")

		println(price)
		storeGphcrd(desc, price+delivery, url)
	})

	c.Visit(urls[0])

}

func coreGaiming() {

	urls := []string{
		"https://coregaming.com.mx/tienda/componentes/tarjetas-de-video",
	}
	c := colly.NewCollector(
		colly.AllowedDomains("https://coregaming.com.mx", "coregaming.com.mx"),
	)

	// Find and visit all links
	c.OnHTML("div.card-item", func(e *colly.HTMLElement) {
		reg, _ := regexp.Compile("[^0-9]+")
		processedString := reg.ReplaceAllString(e.ChildText(".current"), "")
		price, _ := strconv.ParseFloat(processedString, 32)
		price = price / 100
		desc := e.ChildText("div h3")
		url := e.ChildAttr("a", "href")

		storeGphcrd(desc, price, url)
	})

	c.Visit(urls[0])
}

func pcelScrap() {

	urls := []string{
		"https://pcel.com/tarjetas-de-video?sucursal=0&sort=p.price&order=DESC&limit=100",
	}
	c := colly.NewCollector()

	c.Limit(&colly.LimitRule{
		// Set a delay between requests to these domains
		Delay: 2 * time.Second,
		// Add an additional random delay
		RandomDelay: 1 * time.Second,
	})

	// Find and visit all links
	c.OnHTML("div.product-list > table > tbody > tr", func(e *colly.HTMLElement) {
		reg, _ := regexp.Compile("[^0-9]+")
		v1, _ := strconv.ParseFloat(reg.ReplaceAllString(e.ChildText(".price-new"), ""), 32)
		v2, _ := strconv.ParseFloat(reg.ReplaceAllString(e.ChildText(".price"), ""), 32)
		var price float64
		if v1 != 0 {
			price = v1
		} else if v2 != 0 {
			price = v2
		} else {
			return
		}
		price = price / 100
		desc := e.ChildText(".productClick")
		url := e.ChildAttr(".productClick", "href")
		storeGphcrd(desc, price, url)
	})

	for _, url := range urls {
		c.Visit(url)
	}
}

func ddtechScrap() {

	urls := []string{
		"https://ddtech.mx/productos/componentes/tarjetas-de-video?radeon-rx-6000[]=rx-6700xt&stock=con-existencia&orden=primero-existencia&precio=1:99999",
		"https://ddtech.mx/productos/componentes/tarjetas-de-video?geforce-rtx-serie-30[]=rtx-3060-ti&stock=con-existencia&orden=primero-existencia&precio=1:99999",
		"https://ddtech.mx/productos/componentes/tarjetas-de-video?geforce-rtx-serie-30[]=rtx-3060&stock=con-existencia&orden=primero-existencia&precio=1:99999",
	}
	c := colly.NewCollector(
		colly.AllowedDomains("https://ddtech.mx", "ddtech.mx"),
	)

	// Find and visit all links
	c.OnHTML("div.product", func(e *colly.HTMLElement) {
		reg, _ := regexp.Compile("[^0-9]+")
		processedString := reg.ReplaceAllString(e.ChildText(".price"), "")
		price, _ := strconv.ParseFloat(processedString, 32)
		price = price / 100
		desc := e.ChildText("div h3 a")
		url := e.ChildAttr("div h3 a", "href")

		storeGphcrd(desc, price, url)
	})

	for _, url := range urls {
		c.Visit(url)
	}
}

func storeGphcrd(desc string, price float64, url string) {

	dt := time.Now()
	desc = strings.ToLower(desc)
	is6700 := strings.Contains(desc, "6700")
	is3060 := strings.Contains(desc, "3060")
	is3060ti := strings.Contains(desc, "ti")
	name := ""
	if is6700 {
		name = amd_6700xt
	} else if is3060 && is3060ti {
		name = nv_3060ti
	} else if is3060 {
		name = nv_3060
	} else {
		return
	}

	re := regexp.MustCompile(`^(?:https?:\/\/)?(?:[^@\/\n]+@)?(?:www\.)?([^:\/\n]+)`)
	vendor := re.FindString(url)
	tmp := gphcrd{
		vendor: vendor,
		price:  price,
		name:   name,
		desc:   desc,
		url:    url,
		date:   dt.Format("01-02-2006 15:04:05"),
	}
	graphicCards = append(graphicCards, tmp)
}
