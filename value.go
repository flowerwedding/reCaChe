/**
 * @Title  value
 * @description  内存控制，只考虑value的内存，不考虑key的内存
 * @Author  沈来
 * @Update  2020/8/17 14:40
 **/
package reCaChe

//其他类型的长度，要求自己实现该方法
type Value interface {
	Len() int
}