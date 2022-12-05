package binapi

//go:generate binapi-generator --output-dir=. abx
//go:generate binapi-generator --output-dir=. bfd
//go:generate binapi-generator --output-dir=. memclnt
//go:generate binapi-generator --output-dir=. nat64
//go:generate binapi-generator --output-dir=. vpe
//go:generate binapi-generator --output-dir=. isisx

const Version = "22.10"
