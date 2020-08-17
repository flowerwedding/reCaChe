/**
 * @Title  value
 * @description  内存控制，只考虑value的内存，不考虑key的内存
 * @Author  沈来
 * @Update  2020/8/17 14:40
 **/
package reCaChe

type Value interface {
	Len() int
}