# Go编码规范

## 说明
本规范尽量参考并遵从Go语言官方建议、Go语言标准库以及社区编码规范。

## 动词解释：
- **不允许**：表示该行为或风格在本规范中是严厉禁止的
- **不建议**：表示该行为或风格在本规范中是不建议出现的，除非特殊原因，尽量不要出现
- **建议**：表示本规范建议实施该行为或风格，但不是强制的
- **必须**：表示该行为或风格在本规范中是强制要求的
- **尽量**：表示除非特殊原因，尽可能参考该行为或风格

## 开发环境
+ 建议采用mac os或者ubuntu进行代码编译及调试；不建议采用windows开发环境，对于windows pc的情况，建议安装虚拟机
+ 不对开发ide进行强制要求，建议使用sublime编辑器+GoSublime插件

## 目录结构
+ 所有业务相关模块都必须放在src/pandocloud.com目录内，并按照程序模块的功能进行命名
+ 对于可执行程序包，**必须**含有main.go文件；对于库包，必须含有库名对应的.go文件，其中包含包功能说明注释以及主要对外到处接口
+ 程序目录中**不允许**出现不相关的文件（比如废弃的代码）以及不相关的文件夹（比如XXX.bak备份文件）

## 代码缩进
+ 所有代码**必须**经过**go fmt**进行格式化(推荐使用sublime编辑器+GoSublime插件实现自动格式化)
+ 不对代码行宽度进行字符限制，但是对于过长的行建议多行排版并从第二行开始采用tab进行缩进，如：

``` golang
r.Post("/users/verification",
    binding.Json(actions.UserVerifyArgs{}), actions.VendorAuth, actions.ProductAuth,
    actions.SendVerifyCode)
```

## 注释
+ 编码阶段**必须**同步写好变量、函数和包注释，注释**必须**采用英文。
+ 所有注释全部采用双斜杠，注释内容必须是完整的句子，需要以注释的内容作为开头，句号作为结尾，第一个字母和双斜杠之间空一格。
+ 每个程序包必须包含一个包注释，一般在包目录同名的go源码文件或main.go下，如bytes包在bytes.go中的包注释：

``` golang
// Package bytes implements functions for the manipulation of byte slices.
// It is analogous to the facilities of the strings package.
package bytes

```

+ 程序中每一个大写的（将被导出）的名称（变量或者函数），都必须辅以一个文档注释，如：

``` golang
// IndexAny interprets s as a sequence of UTF-8-encoded Unicode code points.
// It returns the byte index of the first occurrence in s of any of the Unicode
// code points in chars.  It returns -1 if chars is empty or if there is no code
// point in common.
func IndexAny(s []byte, chars string) int {
  if len(chars) > 0 {
    var r rune
    var width int
    for i := 0; i < len(s); i += width {
      r = rune(s[i])
      if r < utf8.RuneSelf {
        width = 1
      } else {
        r, width = utf8.DecodeRune(s[i:])
      }
      for _, ch := range chars {
        if r == ch {
          return i
        }
      }
    }
  }
  return -1
}
```

+ 程序中小写的（不会被导出）的名称（变量或者函数），建议辅以一个说明注释，如：

``` golang
// explode splits s into a slice of UTF-8 sequences, one per Unicode code point (still slices of bytes),
// up to a maximum of n byte slices. Invalid UTF-8 sequences are chopped into individual bytes.
func explode(s []byte, n int) [][]byte {
  if n <= 0 {
    n = len(s)
  }
  a := make([][]byte, n)
  var size int
  na := 0
  for len(s) > 0 {
    if na+1 >= n {
      a[na] = s
      na++
      break
    }
    _, size = utf8.DecodeRune(s)
    a[na] = s[0:size]
    s = s[size:]
    na++
  }
  return a[0:na]
}
```

## 命名
+ 建议使用短命名，Go认为文档注释比长名字更容易解释变量意义。
+ 由于Go语言通过变量首字母大小写区分是否导出变量，所有首字母大小写完全由变量是否需要导出来决定。（包括结构体中的变量名与结构体名都采用相同的首字母大小写标准）
+ 包名统一使用小写字母，可以包含数字，库包**不允许**使用下划线，可执行程序包可以使用下划线
+ 接口名**尽量**以er结尾，如Reader, Writer, ReadWriter，**不允许**出现下划线
+ 全局变量名、函数名必须采用驼峰命名法，如setContent(), SetContent(), MaxLength, nBytes，**不允许**出现下划线
+ 局部变量名必须采用全小写命名法，如maxlevel, nbytes, **不允许**出现下划线
+ 函数参数及返回值统一采用全小写，**不允许**出现下划线
+ 对于名称中含有缩写词的变量或者函数，缩写词必须采用全大写形式，如ServerHTTP()
+ 对于接收者命名，**不允许**采用"me", "this" or "self"这种名称，一般采用一两个能代表接收者的名称如"c"或者"cli"表示"Client"，并且保持所有成员函数接收者名称统一，如

``` golang
// sendPackageToDevice retreive package from sending chan and try to send it to device.
func (dev *Device) sendPackageToDevice() {
...
}
```

## 错误处理
+ 普通情况下**不允许**使用panic进行错误处理，必须使用error或多返回值进行错误抛出
+ 错误字符串*必须*采用全小写，比如使用fmt.Errorf("something bad")而不是fmt.Errorf("Something bad")
+ 对于返回值带error的函数，*必须*进程错误处理，如

``` golang
// Calling Close does not close the wrapped io.Reader originally passed to NewReader.
func (z *reader) Close() error {
  if z.err != nil {
    return z.err
  }
  z.err = z.decompressor.Close()
  return z.err
}
```
+ 实现返回值有error的函数时，遇到出错尽早返回

+ **必须**将主线逻辑尽可能放在最外层，当进行错误处理时不允许将正常流程放在else中，如：

``` golang
if err != nil {
    // error handling
    return // or continue, etc.
}
// normal code

```

``` golang
x, err := f()
if err != nil {
    // error handling
    return
}
// use x

```

## 参数传递及返回值
+ 对于int这种小数据，**必须**要传递值。
+ 对于比较大的struct，**必须**传指针。
+ 对于本身就是引用类型的参数，尽量直接传值，如map，string，chan，interface接口等
+ 对于结构体成员函数，定义的接收者(receiver)的类型建议同上一条
+ 对于函数返回值类型都不相同的情况，建议采用匿名返回值；对于返回值类型中有两个或多个相同类型的情况，建议采用命名返回值：

``` golang
func (n *Node) Parent1() *Node

func (n *Node) Parent2() (*Node, error)

// Location returns f's latitude and longitude.
// Negative values mean south and west, respectively.
func (f *Foo) Location() (lat, long float64, err error)

```

## 单元测试
+ 每个可复用的工具函数和工具模块**必须**有单元测试程序(xxx_test.go文件)
+ 单元测试case**必须**覆盖边界值（0值和超大值), 正常值以及异常值
+ 测试中出错的错误输出**必须**打印出输入，错误的输出以及期待的输出，如：

``` golang
if got != tt.want {
    t.Errorf("Foo(%q) = %d; want %d", tt.in, got, tt.want) // or Fatalf, if test can't test anything more past this point
}
```

+ 测试用例比较多的情况下，建议从数组中读取测试输入并依次测试，避免不必要的重复代码拷贝