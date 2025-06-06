package dbHelper

import (
	"bytes"
	"crypto/md5"
	"encoding/gob"
	"encoding/hex"
	"fmt"
	"github.com/google/uuid"
	"io/ioutil"
	"math/rand"
	"os"
	"reflect"
	"sync"
	"time"
)

func FileExists(name string) bool {
	if _, err := os.Stat(name); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}

// PathMkdir 目录不存在则创建
func PathMkdir(path string) {
	_, err := os.Stat(path)
	if err != nil && os.IsNotExist(err) {
		_ = os.MkdirAll(path, 0777)
	}
}

func DeepCopy[T any](dst, src T) error {
	return deepCopy(dst, src)
}

func deepCopy[T any](dst, src T) error {
	var buf bytes.Buffer
	if err := gob.NewEncoder(&buf).Encode(src); err != nil {
		return err
	}
	return gob.NewDecoder(bytes.NewBuffer(buf.Bytes())).Decode(dst)
}

// IdWorker 雪花Id
type IdWorker struct {
	startTime             int64
	workerIdBits          uint
	datacenterIdBits      uint
	maxWorkerId           int64
	maxDatacenterId       int64
	sequenceBits          uint
	workerIdLeftShift     uint
	datacenterIdLeftShift uint
	timestampLeftShift    uint
	sequenceMask          int64
	workerId              int64
	datacenterId          int64
	sequence              int64
	lastTimestamp         int64
	signMask              int64
	idLock                *sync.Mutex
}

func (idw *IdWorker) InitIdWorker(workerId, datacenterId int64) error {
	var baseValue int64 = -1
	idw.startTime = 1463834116272
	idw.workerIdBits = 5
	idw.datacenterIdBits = 5
	idw.maxWorkerId = baseValue ^ (baseValue << idw.workerIdBits)
	idw.maxDatacenterId = baseValue ^ (baseValue << idw.datacenterIdBits)
	idw.sequenceBits = 12
	idw.workerIdLeftShift = idw.sequenceBits
	idw.datacenterIdLeftShift = idw.workerIdBits + idw.workerIdLeftShift
	idw.timestampLeftShift = idw.datacenterIdBits + idw.datacenterIdLeftShift
	idw.sequenceMask = baseValue ^ (baseValue << idw.sequenceBits)
	idw.sequence = 0
	idw.lastTimestamp = -1
	idw.signMask = ^baseValue + 1
	idw.idLock = &sync.Mutex{}
	if idw.workerId < 0 || idw.workerId > idw.maxWorkerId {
		return fmt.Errorf("workerId[%v] is less than 0 or greater than maxWorkerId[%v]",
			workerId, datacenterId)
	}
	if idw.datacenterId < 0 || idw.datacenterId > idw.maxDatacenterId {
		return fmt.Errorf("datacenterId[%d] is less than 0 or greater than maxDatacenterId[%d]",
			workerId, datacenterId)
	}
	idw.workerId = workerId
	idw.datacenterId = datacenterId
	return nil
}

// NextId 返回一个唯一的 INT64 ID
func (idw *IdWorker) NextId() (int64, error) {
	idw.idLock.Lock()
	timestamp := time.Now().UnixNano()
	if timestamp < idw.lastTimestamp {
		return -1, fmt.Errorf(fmt.Sprintf("Clock moved backwards.  Refusing to generate id for %d milliseconds",
			idw.lastTimestamp-timestamp))
	}
	if timestamp == idw.lastTimestamp {
		idw.sequence = (idw.sequence + 1) & idw.sequenceMask
		if idw.sequence == 0 {
			timestamp = idw.tilNextMillis()
			idw.sequence = 0
		}
	} else {
		idw.sequence = 0
	}
	idw.lastTimestamp = timestamp
	idw.idLock.Unlock()
	id := ((timestamp - idw.startTime) << idw.timestampLeftShift) |
		(idw.datacenterId << idw.datacenterIdLeftShift) |
		(idw.workerId << idw.workerIdLeftShift) |
		idw.sequence
	if id < 0 {
		id = -id
	}
	return id, nil
}

// tilNextMillis
func (idw *IdWorker) tilNextMillis() int64 {
	timestamp := time.Now().UnixNano()
	if timestamp <= idw.lastTimestamp {
		timestamp = time.Now().UnixNano() / int64(time.Millisecond)
	}
	return timestamp
}

func ID64() (int64, error) {
	currWorker := &IdWorker{}
	err := currWorker.InitIdWorker(1000, 2)
	if err != nil {
		return 0, err
	}
	return currWorker.NextId()
}

func ID() int64 {
	id, _ := ID64()
	return id
}

func IDStr() string {
	currWorker := &IdWorker{}
	err := currWorker.InitIdWorker(1000, 2)
	if err != nil {
		return ""
	}
	id, err := currWorker.NextId()
	if err != nil {
		return ""
	}
	return AnyToString(id)
}

func IDMd5() string {
	return Get16MD5Encode(IDStr())
}

// MD5 MD5
func MD5(str string) string {
	h := md5.New()
	h.Write([]byte(str))
	return hex.EncodeToString(h.Sum(nil))
}

// Get16MD5Encode 返回一个16位md5加密后的字符串
func Get16MD5Encode(data string) string {
	return GetMD5Encode(data)[8:24]
}

// GetMD5Encode 获取Md5编码
func GetMD5Encode(data string) string {
	h := md5.New()
	h.Write([]byte(data))
	return hex.EncodeToString(h.Sum(nil))
}

// PanicToError panic -> error
func PanicToError(fn func()) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("Panic error: %v", r)
		}
	}()
	fn()
	return
}

// P2E panic -> error
func P2E() {
	defer func() {
		if r := recover(); r != nil {
			Error("Panic error: ", r)
		}
	}()
}

func GetUUID() string {
	return uuid.New().String()
}

const (
	TimeTemplate       = "2006-01-02 15:04:05"
	TimeTemplateNotSec = "2006-01-02 15:04"
)

func TimeStr2Unix(timeStr string) int64 {
	t, err := time.ParseInLocation(TimeTemplate, timeStr, time.Local)
	if err != nil {
		ErrorF("时间字符串 %s 转时间戳错误: %v", timeStr, err.Error())
		return 0
	}
	return t.UnixMilli()
}

func Timestamp2Time(timestamp int64) time.Time {
	return time.Unix(timestamp, 0)
}

func NowTimestampStr() string {
	return AnyToString(time.Now().Unix())
}

func TimeStr2UnixByLayOut(layout, timeStr string) int64 {
	if len(timeStr) < 1 {
		return 0
	}
	t, err := time.ParseInLocation(layout, timeStr, time.Local)
	if err != nil {
		//logger.ErrorTimes(3, "时间字符串 %s 转时间戳错误: %v", timeStr, err.Error())
		return 0
	}
	return t.UnixMilli()
}

func Timestamp2Week(timestamp int64) string {
	tm := time.UnixMilli(timestamp)
	switch tm.Weekday() {
	case time.Sunday:
		return "周天"
	case time.Monday:
		return "周一"
	case time.Tuesday:
		return "周二"
	case time.Wednesday:
		return "周三"
	case time.Thursday:
		return "周四"
	case time.Friday:
		return "周五"
	case time.Saturday:
		return "周六"
	}
	return ""
}

func GetDate() string {
	now := time.Now()
	return now.Format("2006-01-02")
}

func TimestampSubDay(timestamp int64) int64 {
	timestampTime := time.Unix(timestamp, 0)
	now := time.Now()
	day := int64(now.Sub(timestampTime).Hours() / 24)
	return day
}

func TimeHM(t time.Time) string {
	return t.Format("15:04")
}

func TimeYMDCN(t time.Time) string {
	return t.Format("2006年01月02日")
}

func TimeToYMDHMS(t time.Time) string {
	return t.Format(TimeTemplate)
}

func TimeToYMD(t time.Time) string {
	return t.Format("2006-01-02")
}

// FileMd5sum 文件 Md5
func FileMd5sum(fileName string) string {
	fin, err := os.OpenFile(fileName, os.O_RDONLY, 0644)
	if err != nil {
		Error(fileName, err)
		return ""
	}
	defer func() {
		_ = fin.Close()
	}()
	buf, bufErr := ioutil.ReadFile(fileName)
	if bufErr != nil {
		Error(fileName, bufErr)
		return ""
	}
	m := md5.Sum(buf)
	return hex.EncodeToString(m[:16])
}

// GetAllFile 获取目录下的所有文件
func GetAllFile(pathname string) ([]string, error) {
	s := make([]string, 0)
	rd, err := ioutil.ReadDir(pathname)
	if err != nil {
		Error("read dir fail:", err)
		return s, err
	}
	for _, fi := range rd {
		if !fi.IsDir() {
			fullName := pathname + "/" + fi.Name()
			s = append(s, fullName)
		}
	}
	return s, nil
}

// RandomIntCaptcha 生成 captchaLen 位随机数，理论上会重复
func RandomIntCaptcha(captchaLen int) string {
	var arr string
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < captchaLen; i++ {
		arr = arr + fmt.Sprintf("%d", r.Intn(10))
	}
	return arr
}

// DeepEqual 深度比较任意类型的两个变量的是否相等,类型一样值一样反回true
// 如果元素都是nil，且类型相同，则它们是相等的; 如果它们是不同的类型，它们是不相等的
func DeepEqual(a, b interface{}) bool {
	ra := reflect.Indirect(reflect.ValueOf(a))
	rb := reflect.Indirect(reflect.ValueOf(b))
	if raValid, rbValid := ra.IsValid(), rb.IsValid(); !raValid && !rbValid {
		return reflect.TypeOf(a) == reflect.TypeOf(b)
	} else if raValid != rbValid {
		return false
	}
	return reflect.DeepEqual(ra.Interface(), rb.Interface())
}

var randObj = rand.New(rand.NewSource(time.Now().UnixNano()))

func SliceContains[V comparable](a []V, v V) bool {
	l := len(a)
	if l == 0 {
		return false
	}
	for i := 0; i < l; i++ {
		if a[i] == v {
			return true
		}
	}
	return false
}

func SliceDeduplicate[V comparable](a []V) []V {
	l := len(a)
	if l < 2 {
		return a
	}
	seen := make(map[V]struct{})
	j := 0
	for i := 0; i < l; i++ {
		if _, ok := seen[a[i]]; ok {
			continue
		}
		seen[a[i]] = struct{}{}
		a[j] = a[i]
		j++
	}
	return a[:j]
}

func SliceDel[V comparable](a []V, i int) []V {
	l := len(a)
	if l == 0 {
		return a
	}
	if i < 0 || i > l-1 {
		return a
	}
	return append(a[:i], a[i+1:]...)
}

type number interface {
	int | int8 | int16 | int32 | int64 |
		uint | uint8 | uint16 | uint32 | uint64 |
		float32 | float64
}

func SliceMax[V number](a []V) V {
	l := len(a)
	if l == 0 {
		var none V
		return none
	}
	maxV := a[0]
	for k := 1; k < l; k++ {
		if a[k] > maxV {
			maxV = a[k]
		}
	}
	return maxV
}

func SliceMin[V number](a []V) V {
	l := len(a)
	if l == 0 {
		return 0
	}
	minV := a[0]
	for k := 1; k < l; k++ {
		if a[k] < minV {
			minV = a[k]
		}
	}
	return minV
}

func SlicePop[V comparable](a []V) (V, []V) {
	if len(a) == 0 {
		var none V
		return none, a
	}
	return a[len(a)-1], a[:len(a)-1]
}

func SliceReverse[V comparable](a []V) []V {
	l := len(a)
	if l == 0 {
		return a
	}
	for s, e := 0, len(a)-1; s < e; {
		a[s], a[e] = a[e], a[s]
		s++
		e--
	}
	return a
}

func SliceShuffle[V comparable](a []V) []V {
	l := len(a)
	if l <= 1 {
		return a
	}
	randObj.Shuffle(l, func(i, j int) {
		a[i], a[j] = a[j], a[i]
	})
	return a
}

func SliceCopy[V comparable](a []V) []V {
	return append(a[:0:0], a...)
}

func SliceRand[V comparable](a []V) V {
	rand.New(rand.NewSource(time.Now().UnixNano()))
	randomIndex := rand.Intn(len(a))
	return a[randomIndex]
}
