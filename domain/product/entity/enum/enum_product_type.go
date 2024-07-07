package enum

type ProductType int

const (
	product_all = iota
)

var ProductTypeMap = map[ProductType]string{
	product_all: "全部商品",
}
