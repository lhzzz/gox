package structure

import "reflect"

// CommonResultSlicePageLimit 通用的结果分页的方法
// offset: 截取切片的起始位置
// maxResult: 一次最多截取的结果数量
// inSlice: 要截取的切片，类型slice
// outSlice: 保存输出结果的切片，类型slice pointer
// return: 返回分页后下一页数据的起始位置，-1代表没有下一页
func CommonResultSlicePageLimit(offset int, maxResult int, inSlice interface{}, outSlice interface{}) int {
	inSliceV := reflect.ValueOf(inSlice)
	inSliceV.Type().Name()
	outSliceV := reflect.ValueOf(outSlice).Elem()
	if inSliceV.Len() < offset {
		offset = -1
	} else if inSliceV.Len() <= offset+maxResult {
		outSliceV.Set(inSliceV.Slice(offset, inSliceV.Len()))
		offset = -1
	} else {
		outSliceV.Set(inSliceV.Slice(offset, offset+maxResult))
		offset += maxResult
	}
	return offset
}

// MapKeyConversionSlice 将map的所有key放入slice中
// inMap: 类型map
// outSlice: 类型slice pointer
func MapKeyConversionSlice(inMap interface{}, outSlice interface{}) {
	outSliceV := reflect.ValueOf(outSlice).Elem()
	inMapV := reflect.ValueOf(inMap)
	outSliceV.Set(reflect.Append(outSliceV, inMapV.MapKeys()...))
}

// MapValueConversionSlice 将map的所有value放入slice中
// inMap: 类型map
// outSlice: 类型slice pointer
func MapValueConversionSlice(inMap interface{}, outSlice interface{}) {
	inMapV := reflect.ValueOf(inMap)
	outSliceV := reflect.ValueOf(outSlice).Elem()
	tempSlice := reflect.MakeSlice(outSliceV.Type(), 0, 0)
	for it := inMapV.MapRange(); it.Next(); {
		tempSlice = reflect.Append(tempSlice, it.Value())
	}
	outSliceV.Set(tempSlice)
}
