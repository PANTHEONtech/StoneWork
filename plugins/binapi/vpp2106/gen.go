package binapi

//go:generate binapi-generator --output-dir=. --input-dir=api abx
//go:generate binapi-generator --output-dir=. --input-dir=api bfd
//go:generate binapi-generator --output-dir=. --input-dir=api nat64
//go:generate binapi-generator --output-dir=. --input-dir=api vpe
//go:generate binapi-generator --output-dir=. --input-dir=api isisx

const Version = "21.06"
