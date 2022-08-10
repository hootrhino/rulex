package sensor_server

/*
*
* 默认设备，未来可能有很多别的设备
*
 */

type Sensor struct {
	session Session
	sn      string
	Authed  bool
}

func (s *Sensor) Sn() string {
	return s.sn
}
func (s *Sensor) Session() Session {

	return s.session
}
func (s *Sensor) Ping() []byte {
	return []byte{}

}
func (s *Sensor) OnRegister(sn string) error {
	s.sn = sn
	return nil
}
func (s *Sensor) OnLine() {

}
func (s *Sensor) OffLine() {

}
func (s *Sensor) OnError(error) {

}
func (s *Sensor) OnData([]byte) {

}

func NewSensor(session Session) ISensor {
	return &Sensor{session: session}
}
