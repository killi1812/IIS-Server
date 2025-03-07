package apiq

import (
	"encoding/xml"
	"fmt"
	"net/http"
	"strings"

	"github.com/killi1812/libxml2/parser"
	"github.com/killi1812/libxml2/types"
	"go.uber.org/zap"
)

func GetWeatherForCity(query string) (*City, error) {
	res, err := http.Get("https://vrijeme.hr/hrvatska_n.xml")
	if err != nil {
		return nil, err
	}

	p := parser.New()
	doc, err := p.ParseReader(res.Body)
	defer res.Body.Close()
	if err != nil {
		return nil, err
	}
	defer doc.Free()

	root, err := doc.DocumentElement()
	if err != nil {
		zap.S().Debugf("Failed to fetch document element: %s", err)
		return nil, err
	}

	/*
		ctx.RegisterNS("atom", "http://www.w3.org/2005/Atom")
		title := xpath.String(ctx.Find("/atom:feed/atom:title/text()"))
	*/

	return FindCity(root, query)
}

func FindCity(node types.Node, query string) (*City, error) {
	///Hrvatska/Grad[1]/GradIme
	xpathQ := fmt.Sprintf("/Hrvatska/Grad[GradIme='%s']", query)
	rez, err := node.Find(xpathQ)
	if err != nil {
		return nil, err
	}
	defer rez.Free()

	if !rez.Bool() {
		zap.S().Infof("Failed to find node q: %v, rez: %v", xpathQ, rez)
		return nil, ErrCityNotFound
	}

	fmt.Printf("rez: %v\n", rez.Bool())

	cityNode := rez.NodeIter().Node()

	dec := strings.NewReader(cityNode.String())
	var grad City
	err = xml.NewDecoder(dec).Decode(&grad)
	if err != nil {
		return nil, err
	}

	return &grad, nil
}
