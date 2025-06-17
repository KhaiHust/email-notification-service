package properties

type Sasl struct {
	Mechanism string `default:"PLAIN"`
	Username  string `default:""`
	Password  string `default:""`
}
