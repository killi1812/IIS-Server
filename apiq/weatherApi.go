package apiq

import (
	"encoding/xml"
	"fmt"
	"net/http"
	"strings"

	"github.com/killi1812/libxml2/parser"
	"github.com/killi1812/libxml2/types"
	"github.com/killi1812/libxml2/xpath"
	"go.uber.org/zap"
)

type WeaterService struct{}

func (*WeaterService) GetWeatherForCity(query string) ([]City, error) {
	res, err := http.Get("https://vrijeme.hr/hrvatska_n.xml")
	if err != nil {
		return nil, err
	}
	// TODO: ne parsa zadni node
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
	// defer root.Free()

	return FindCity(root, query)
}

func FindCity(node types.Node, query string) ([]City, error) {
	xpathQ := fmt.Sprintf("//Grad[contains(GradIme,'%s')]", query)
	rez := xpath.NodeList(node.Find(xpathQ))
	city := make([]City, 0, len(rez))

	for _, cnode := range rez {
		var loc City
		dec := strings.NewReader(cnode.String())
		err := xml.NewDecoder(dec).Decode(&loc)
		if err != nil {
			return nil, err
		}
		zap.S().Infof("adding location: %+v", loc)
		city = append(city, loc)
	}

	return city, nil
}
