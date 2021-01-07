package gonfig

// Provider is used to provide values
// It can implement either Unmarshaler or Filler interface or both
// Name method is used for error messages
type Provider interface {
	Name() string
}

// Unmarshaler can be implemented by providers to receive struct pointer and unmarshal values into it
type Unmarshaler interface {
	UnmarshalStruct(i interface{}) (err error)
}

// Filler can be implemented by providers to receive struct fields and set their value
type Filler interface {
	Fill(in *Input) (err error)
}
