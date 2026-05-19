// Package nork is Nodes with ORdered Kids - a way to create
// a hierarchical tree of nodes, where the nodes are ordered
// and keep their order. It does not need or use Go generics.
//
// Node order is of course important for markup in general 
// and XML mixedcontent in particular, while unimportant 
// when XML is used purely for data records.
//
// interface [Norker] is implemented not for type [Nork] but 
// rather for `*Nork´. This is so that nodes are shared-writable. 
//
package nork
