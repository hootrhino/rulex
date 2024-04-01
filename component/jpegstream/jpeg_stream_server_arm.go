//go:build arm
// +build arm

package jpegstream

import (
	"fmt"
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

/*
*
* Manage API
*
 */

func (s *JpegStreamServer) RegisterJpegStreamSource(liveId string) error {

	return fmt.Errorf("stream already exists")
}

func (s *JpegStreamServer) GetJpegStreamSource(liveId string) (*JpegStream, error) {

	return nil, nil

}

func (s *JpegStreamServer) Exists(liveId string) bool {
	return true
}
func (s *JpegStreamServer) DeleteJpegStreamSource(liveId string) {

}

func (s *JpegStreamServer) JpegStreamSourceList() []JpegStream {
	return nil
}
func (s *JpegStreamServer) JpegStreamFlush() {
}
