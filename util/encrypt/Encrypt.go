package encrypt

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"fmt"
	"os"
)

func Test() {
	// 生成密钥对，保存到文件
	GenerateRSAKey(2048)
	message := []byte("rsa 加密解密测试！")
	// 加密
	cipherText := RSA_Encrypt(message, "public.pem")
	cipherText_base64 := base64.StdEncoding.EncodeToString(cipherText) // 将 []byte 类型的密文转换为 base64 字符串
	fmt.Println("加密后为(base64)：", cipherText_base64)
	fmt.Println("加密后为([]byte)：", cipherText)
	// 解密
	cipherText, _ = base64.StdEncoding.DecodeString(cipherText_base64) // 若转换为输入为 base64 字符串，则需先解码为 []byte
	plainText := RSA_Decrypt(cipherText, "private.pem")
	fmt.Println("解密后为([]byte)：", plainText)
	fmt.Println("解密后为(string)：", string(plainText))
}

// 生成 RSA 私钥和公钥，保存到文件中
func GenerateRSAKey(bits int) {
	// GenerateKey 函数使用随机数据生成器 random 生成一对具有指定字位数的 RSA 密钥
	// Reader 是一个全局、共享的密码用强随机数生成器
	privateKey, err := rsa.GenerateKey(rand.Reader, bits)
	if err != nil {
		panic(err)
	}
	// 保存私钥
	// 通过 x509 标准将得到的 ras 私钥序列化为 ASN.1 的 DER 编码字符串
	X509PrivateKey := x509.MarshalPKCS1PrivateKey(privateKey)
	// 使用 pem 格式对 x509 输出的内容进行编码
	// 创建文件保存私钥
	privateFile, err := os.Create("private.pem")
	if err != nil {
		panic(err)
	}
	defer privateFile.Close()
	// 构建一个 pem.Block 结构体对象
	privateBlock := pem.Block{Type: "RSA Private Key", Bytes: X509PrivateKey}
	// 将数据保存到文件
	pem.Encode(privateFile, &privateBlock)
	// 保存公钥
	// 获取公钥的数据
	publicKey := privateKey.PublicKey
	// X509 对公钥编码
	X509PublicKey, err := x509.MarshalPKIXPublicKey(&publicKey)
	if err != nil {
		panic(err)
	}
	// pem 格式编码
	// 创建用于保存公钥的文件
	publicFile, err := os.Create("public.pem")
	if err != nil {
		panic(err)
	}
	defer publicFile.Close()
	// 创建一个 pem.Block 结构体对象
	publicBlock := pem.Block{Type: "RSA Public Key", Bytes: X509PublicKey}
	// 保存到文件
	pem.Encode(publicFile, &publicBlock)
}

// RSA 加密
func RSA_Encrypt(plainText []byte, path string) []byte {
	// 打开文件
	file, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	// 读取文件的内容
	info, _ := file.Stat()
	buf := make([]byte, info.Size())
	file.Read(buf)
	// pem 解码
	block, _ := pem.Decode(buf)
	// x509 解码
	publicKeyInterface, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		panic(err)
	}
	// 类型断言
	publicKey := publicKeyInterface.(*rsa.PublicKey)
	// 对明文进行加密
	cipherText, err := rsa.EncryptPKCS1v15(rand.Reader, publicKey, plainText)
	if err != nil {
		panic(err)
	}
	// 返回 []byte 密文
	return cipherText
}

// RSA 解密
func RSA_Decrypt(cipherText []byte, path string) []byte {
	// 打开文件
	file, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	// 获取文件内容
	info, _ := file.Stat()
	buf := make([]byte, info.Size())
	file.Read(buf)
	// pem 解码
	block, _ := pem.Decode(buf)
	// X509 解码
	privateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		panic(err)
	}
	// 对密文进行解密
	plainText, _ := rsa.DecryptPKCS1v15(rand.Reader, privateKey, cipherText)
	// 返回明文
	return plainText
}
