package binapi

//go:generate binapi-generator --output-dir=. --input-dir=api abx
//go:generate binapi-generator --output-dir=. --input-dir=api bfd
//go:generate binapi-generator --output-dir=. --input-dir=api nat
//go:generate binapi-generator --output-dir=. --input-dir=api vpe

const Version = "20.09"
