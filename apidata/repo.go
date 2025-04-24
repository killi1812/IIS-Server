package apidata

import (
	"encoding/xml"
	"errors"
	"fmt"
	"iis_server/apiq"
	"os"
	"strings"

	"github.com/killi1812/libxml2/parser"
	"github.com/killi1812/libxml2/types"
	"github.com/killi1812/libxml2/xpath"
	"go.uber.org/zap"
)

var ErrNotFound = errors.New("file not found")

// TODO: format xml style
type userInfoRepo struct {
	users []apiq.UserInfo
}

func Search(username string) ([]apiq.UserInfo, error) {
	// TODO: see about mode
	file, err := os.OpenFile("userInfoRepo.xml", os.O_RDONLY, os.ModeCharDevice)
	if err != nil {
		fmt.Println("Error reading file:", err)
		return nil, ErrNotFound
	}

	p := parser.New()
	doc, err := p.ParseReader(file)
	if err != nil {
		return nil, err
	}
	defer doc.Free()

	root, err := doc.DocumentElement()
	if err != nil {
		zap.S().Debugf("Failed to fetch document element: %s", err)
		return nil, err
	}
	defer root.Free()

	rez, err := find(root, "")
	if err != nil {
		zap.S().Debugf("Failed to find node : %s", err)
		return nil, err
	}

	return rez, nil
}

func find(node types.Node, query string) ([]apiq.UserInfo, error) {
	xpathQ := fmt.Sprintf("//UserInfo[contains(Username,'%s')]", query)
	rez := xpath.NodeList(node.Find(xpathQ))
	userInfos := make([]apiq.UserInfo, 0, len(rez))

	for _, cnode := range rez {
		var loc apiq.UserInfo
		dec := strings.NewReader(cnode.String())
		err := xml.NewDecoder(dec).Decode(&loc)
		if err != nil {
			return nil, err
		}
		zap.S().Infof("adding location: %+v", loc)
		userInfos = append(userInfos, loc)
	}

	return userInfos, nil
}
