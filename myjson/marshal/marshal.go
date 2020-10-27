package marshal

// JsonMarshal 用于将结构数据解析成为字节流，支持的类型包括（嵌套）结构体，map，结构体数组。结构内的数据必须是go的基本类型
func JsonMarshal(v interface{}) ([]byte, error) {
	//存储v的编组之后的bytes.Buffer对象（见下 具体数据结构）
	e := &EncodeState{}
	//调用编组对象的marshal方法，encOpts 对象：编码过程的配置
	err := e.marshal(v)
	if err != nil {
		return nil, err
	}
	return e.Bytes(), nil
}
