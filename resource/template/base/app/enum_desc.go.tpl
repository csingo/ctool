
var ##ENUM##_desc = map[int32]string{##ENUMVALUE##}

func (x ##ENUM##) Desc() string {
	return ##ENUM##_desc[int32(x)]
}
