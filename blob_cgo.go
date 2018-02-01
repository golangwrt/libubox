package libubox

/*
#cgo LDFLAGS: -lblobmsg_json -lubox
#include <libubox/blobmsg_json.h>

int blobmsg_print(struct blob_buf *buf, const char *name, const char *str)
{
	return blobmsg_printf(buf, name, "%s", str);
}

*/
import "C"
import (
	"encoding/json"
	"fmt"
	"os"
	"reflect"
	"strings"
	"unsafe"
)

// BlobBuf encapsulates struct blob_buf
type BlobBuf struct {
	ptr *C.struct_blob_buf
	head *BlobAttr
}

// BlobAttr encapsulates struct blob_attr
type BlobAttr struct {
	ptr *C.struct_blob_attr
}

// New create BlobBuf using C malloc for underlying struct blob_buf
func NewBlobBuf() *BlobBuf {
	p := C.calloc(1, C.sizeof_struct_blob_buf)
	C.memset(p, 0, C.sizeof_struct_blob_buf)
	return &BlobBuf{
		ptr: (*C.struct_blob_buf)(p),
	}
}

// Init initialize the underlying field 'struct blob_buf' with blob_buf_init
// int blob_buf_init(struct blob_buf *buf, int id)
func (buf *BlobBuf) Init(id int) int {
	return int(C.blob_buf_init(buf.ptr, C.int(id)))
}

// Free deallocate underlying struct blob_buf field with blob_buf_free
func (buf *BlobBuf) Free() {
	C.blob_buf_free(buf.ptr)
}

// AddU32 encapsulate blobmsg_add_u32
// int blobmsg_add_u32(struct blob_buf *buf, const char *name, uint32_t val)
func (buf *BlobBuf) AddU32(name string, val uint32) int {
	cname := C.CString(name)
	defer C.free(unsafe.Pointer(cname))

	return int(C.blobmsg_add_u32(buf.ptr, cname, C.uint32_t(val)))
}

// AddU16 encapsulate blobmsg_add_u16
// int blobmsg_add_u16(struct blob_buf *buf, const char *name, uint16_t val)
func (buf *BlobBuf) AddU16(name string, val uint16) int {
	cname := C.CString(name)
	defer C.free(unsafe.Pointer(cname))

	return int(C.blobmsg_add_u16(buf.ptr, cname, C.uint16_t(val)))
}

// AddU8 encapsulate blobmsg_add_u8
// blobmsg_add_u8(struct blob_buf *buf, const char *name, uint8_t val)
func (buf *BlobBuf) AddU8(name string, val uint8) int {
	cname := C.CString(name)
	defer C.free(unsafe.Pointer(cname))

	return int(C.blobmsg_add_u8(buf.ptr, cname, C.uint8_t(val)))
}

// AddBool invoke AddU8 with value 0 or 1
func (buf *BlobBuf) AddBool(name string, val bool) int {
	var tmp uint8
	if val {
		tmp = 1
	}
	return buf.AddU8(name, tmp)
}

// AddU64 encapsulate blobmsg_add_u64
// int blobmsg_add_u64(struct blob_buf *buf, const char *name, uint64_t val)
func (buf *BlobBuf) AddU64(name string, val uint64) int {
	cname := C.CString(name)
	defer C.free(unsafe.Pointer(cname))

	return int(C.blobmsg_add_u64(buf.ptr, cname, C.uint64_t(val)))
}

// AddDouble encapsulate blobmsg_add_double
// int blobmsg_add_double(struct blob_buf *buf, const char *name, double val)
func (buf *BlobBuf) AddDouble(name string, val float64) int {
	cname := C.CString(name)
	defer C.free(unsafe.Pointer(cname))

	return int(C.blobmsg_add_double(buf.ptr, cname, C.double(val)))
}

// Head return the front BlobAttr (head filed) of BlobBuf
func (buf *BlobBuf) Head() *BlobAttr {
	if buf.head == nil {
		buf.head = &BlobAttr{
			ptr: buf.ptr.head,
		}
	}
	return buf.head
}

// OpenNested encapsulate blobmsg_open_nested
// void *blobmsg_open_nested(struct blob_buf *buf, const char *name, bool array)
func (buf *BlobBuf) OpenNested(name string, isArray bool) unsafe.Pointer {
	cname := C.CString(name)
	defer C.free(unsafe.Pointer(cname))

	return C.blobmsg_open_nested(buf.ptr, cname, C.bool(isArray))
}

// NestEnd encapsulate blob_nest_end
// void blob_nest_end(struct blob_buf *buf, void *cookie)
func (buf *BlobBuf) NestEnd(nest unsafe.Pointer) {
	C.blob_nest_end(buf.ptr, nest)
}

// AddString encapsulate blobmsg_add_string
// int blobmsg_add_string(struct blob_buf *buf, const char *name, const char *string)
func (buf *BlobBuf) AddString(name string, val string) int {
	cname := C.CString(name)
	defer C.free(unsafe.Pointer(cname))

	cval := C.CString(val)
	defer C.free(unsafe.Pointer(cval))

	return int(C.blobmsg_add_string(buf.ptr, cname, cval))
}

func (buf *BlobBuf) AddObject(name string, o interface{}) {
	var err error
	var p unsafe.Pointer
	var isNamedField = false

	if o == nil {
		return
	}
	if len(name) > 0 {
		isNamedField = true
	}
	v := reflect.ValueOf(o)
	for v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	switch v.Kind() {
	case reflect.Map:
		if isNamedField {
			p = buf.OpenNested(name, false)
		}
		for _, mapKey := range v.MapKeys() {
			buf.AddObject(fmt.Sprintf("%v", mapKey.Interface()), v.MapIndex(mapKey).Interface())
		}
		if isNamedField {
			buf.NestEnd(p)
		}
	case reflect.Struct:
		if isNamedField {
			p = buf.OpenNested(name, false)
		}
		buf.AddStruct(v.Interface())
		if isNamedField {
			buf.NestEnd(p)
		}
	case reflect.Slice:
		p = buf.OpenNested(name, true)
		for i := 0; i < v.Len(); i++ {
			buf.AddObject("", v.Index(i).Interface())
		}
		buf.NestEnd(p)
	case reflect.Bool:
		buf.AddBool(name, v.Bool())
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		buf.AddU64(name, uint64(v.Int()))
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		buf.AddU64(name, v.Uint())
	case reflect.Float32, reflect.Float64:
		buf.AddDouble(name, v.Float())
	case reflect.String:
		buf.AddString(name, v.String())
	default:
		err = fmt.Errorf("AddObject: skip unsupported object with type: %s\n", v.Kind().String())
		goto out
	}
out:
	if err != nil {
		buf.AddString("error", err.Error())
	}
}

func (buf *BlobBuf) AddStruct(s interface{}) {
	var err error
	if s == nil {
		return
	}

	v := reflect.ValueOf(s)
	for v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	if v.Kind() != reflect.Struct {
		err = fmt.Errorf("AddObject: param <%v> is not a struct, it's %s\n", s, v.Kind().String())
		goto out
	}
out:
	if err != nil {
		buf.AddString("error", err.Error())
		return
	}
	var f reflect.StructField
	var fv reflect.Value

	for i := 0; i < v.NumField(); i++ {
		var key string
		f = v.Type().Field(i)
		if len(f.Tag) > 0 {
			fields := strings.SplitN(f.Tag.Get("json"), ",", 2)
			if len(fields[0]) > 0 {
				key = fields[0]
				if strings.HasPrefix(key, "-") {
					continue
				}
			} else {
				key = f.Name
			}
		}
		if len(key) == 0 {
			key = f.Name
		}
		fv = v.Field(i)
		for fv.Kind() == reflect.Ptr {
			fv = fv.Elem()
		}
		buf.AddObject(key, fv.Interface())
	}
}

// AddJsonFrom marshal object o using json encoding, and then invoke AddJsonFromString
// refer to AddJsonFromString for limitations
func (buf *BlobBuf) AddJsonFrom(o interface{}) {
	out, err := json.Marshal(o)
	if err != nil {
		buf.AddString("error", err.Error())
	} else {
		buf.AddJsonFromString(string(out))
	}
}

// AddJsonFromString encapsulate blobmsg_add_json_from_string
// bool blobmsg_add_json_from_string(struct blob_buf *b, const char *str)
//
// pay attention for following limitations:
//
// 1. blobmsg_add_json_from_string can not present integer larger than
// int32 max, it's limitation in blobmsg implementation of libubox, which use int32
// for json type int. Use AddObject instead to avoid this limitation.
//
// 2. the underlying C implementation only supports string standards for json object
//
func (buf *BlobBuf) AddJsonFromString(str string) error {
	cstr := C.CString(str)
	defer C.free(unsafe.Pointer(cstr))
	_, err := C.blobmsg_add_json_from_string(buf.ptr, cstr)
	return err
}

func (buf *BlobBuf) Printf(name string, format string, args ...interface{}) error {
	cname := C.CString(name)
	defer C.free(unsafe.Pointer(cname))
	str := fmt.Sprintf(format, args...)
	cstr := C.CString(str)
	defer C.free(unsafe.Pointer(cstr))
	_, err := C.blobmsg_print(buf.ptr, cname, cstr)
	return err
}

// PadLen return the padded length of an attribute (incuding the header)
// unsigned int blob_pad_len(const struct blob_attr *attr)
func (attr *BlobAttr) PadLen() uint {
	return uint(C.blob_pad_len(attr.ptr))
}

// RawLen return the complete length of an attribute (including the header)
// unsigned int blob_raw_len(const struct blob_attr *attr)
func (attr *BlobAttr) RawLen() uint {
	return uint(C.blob_raw_len(attr.ptr))
}

// Length return the lenght of the attribute's payload
// unsigned int blob_len(const struct blob_attr *attr)
func (attr *BlobAttr) Length() uint {
	return uint(C.blob_len(attr.ptr))
}

// Data return the data pointer for an attribute
// void* blob_data(const struct blob_attr *attr)
func (attr *BlobAttr) Data() unsafe.Pointer {
	return C.blob_data(attr.ptr)
}

// ID return the id of of an attribute
// unsigned int blob_id(const struct blob_attr *attr)
func (attr *BlobAttr) ID() uint {
	return uint(C.blob_id(attr.ptr))
}

// IsExtended test whether BLOB_ATTR_EXTENDED (0x80000000) is set
// bool blob_is_extended(const struct blob_attr *attr)
func (attr *BlobAttr) IsExtended() bool {
	return bool(C.blob_is_extended(attr.ptr))
}

// GetString method return the string value of the blob attribute
func (attr *BlobAttr) GetString() string {
	return C.GoString((*C.char)(attr.Data()))
}

func (attr *BlobAttr) GetU8() uint8 {
	return uint8(C.blob_get_u8(attr.ptr))
}

// FormatJSON encapsulate blobmsg_format_json
// char *blobmsg_format_json(struct blob_attr *attr, bool list)
func (attr *BlobAttr) FormatJSON(list bool) string {
	cstr, err := C.blobmsg_format_json(attr.ptr, C.bool(list))
	if cstr != nil {
		defer C.free(unsafe.Pointer(cstr))
		return C.GoString(cstr)
	} else {
		fmt.Fprintf(os.Stderr, "blobmsg_format_json return NULL, %s\n", err)
		return ""
	}
}

// FormatJSONValue encapsulate blobmsg_format_json_value
// char *blobmsg_format_json_value(struct blob_attr *attr)
func (attr *BlobAttr) FormatJSONValue() string {
	cstr, err := C.blobmsg_format_json_value(attr.ptr)
	if cstr != nil {
		defer C.free(unsafe.Pointer(cstr))
		return C.GoString(cstr)
	} else {
		fmt.Fprintf(os.Stderr, "blobmsg_format_json return NULL, %s\n", err)
		return ""
	}
}

// Unmarshal obj shall be a pointer of required object type
func (attr *BlobAttr) Unmarshal(obj interface{}) error {
	data := attr.FormatJSON(true)
	err := json.Unmarshal([]byte(data), obj)
	if err != nil {
		return fmt.Errorf("BlobAttr json unmarshal to %T failed , error: %s", obj, err)
	}
	return nil
}

// Pointer return the underlying *C.struct_blob_attr
func (attr *BlobAttr) Pointer() unsafe.Pointer {
	return unsafe.Pointer(attr.ptr)
}

// NewBlobAttrFrom create *BlobAttr from ptr (must be *C.struct_blob_attr)
func NewBlobAttrFromPointer(ptr unsafe.Pointer) *BlobAttr {
	return &BlobAttr{
		ptr: (*C.struct_blob_attr)(ptr),
	}
}

