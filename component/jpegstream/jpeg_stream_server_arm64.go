//go:build arm64
// +build arm64

package jpegstream

import (
	"github.com/hootrhino/rulex/component"
	"github.com/hootrhino/rulex/typex"
)

type JpegStreamServer struct {
}

func InitJpegStreamServer(rulex typex.RuleX) {

}

func (s *JpegStreamServer) Init(cfg map[string]any) error {

	return nil
}
func (s *JpegStreamServer) Start(r typex.RuleX) error {

	return nil
}
func (s *JpegStreamServer) Stop() error {
	return nil
}
func (s *JpegStreamServer) PluginMetaInfo() component.XComponentMetaInfo {
	return component.XComponentMetaInfo{}
}
