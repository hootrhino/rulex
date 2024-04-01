//go:build arm
// +build arm

package jpegstream

import (
	"github.com/hootrhino/rulex/component"
	"github.com/hootrhino/rulex/typex"
)
var __DefaultJpegStreamServer *JpegStreamServer

type JpegStreamServer struct {
}

func InitJpegStreamServer(rulex typex.RuleX) {
	__DefaultJpegStreamServer = &JpegStreamServer{}
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
