package main

import (
	"encoding/json"
  "encoding/xml"
	"fmt"
	"os"
	"strings"
  "net/http"
	"github.com/antchfx/xmlquery"
  "html/template"
	// "text/template"
  "io/ioutil"
	"github.com/gorilla/mux"
	"crypto/rand"
	"golang.org/x/crypto/nacl/secretbox"
	"io"
	"log"
	"flag"
  "time"
)


/*шифрование */
const keySize = 32
const nonceSize = 24

// “test_encrypt” ключ по умолчанию
var userKey = flag.String("k", "test_encrypt", "encryption key")
var pad = []byte("«encrypt on Go - android app dev»")



//задаём родительскую структуру json
type DataS struct {
    XMLName    xml.Name `xml:"rates" json:"-"`
    CoinList []Coin `xml:"item" json:"coinList"`
}
//задаем структуру дочернюю
type Coin struct {
    XMLName   xml.Name `xml:"item" json:"-"`
    CoinName string   `xml:"from" json:"CoinName"`
    To string   `xml:"to" json:"to"`
    In string `xml:"in" json:"in"`
    Out string `xml:"out" json:"out"`
    Amount string `xml:"amount" json:"amount"`
    Minamount string `xml:"minamount" json:"minamount"`
    Maxamount string `xml:"maxamount" json:"in"`
    City string `xml:"city" json:"city"`
}
//задаём структуру json
type Contents struct {
	Coins []DynamicValues
}

type DynamicValues struct {
	Coin map[string]interface{}
}

// функция на добавление к списку содержимого массива
func (contents *Contents) AddItemToDynamicList(dynamicValue DynamicValues) []DynamicValues {
	contents.Coins = append(contents.Coins, dynamicValue)
	return contents.Coins
}

  //динамическое обогащение объекта значениями xml

func (dynamicValue *DynamicValues) enrichDynamicValue(nodeToTransform *xmlquery.Node) map[string]interface{} {

	newValue := map[string]interface{}{}
	// цикл прохода по узлу
	for cryptoC := nodeToTransform.FirstChild; cryptoC != nil; cryptoC = cryptoC.NextSibling {
		// проверка на пустоту строк
		nodeIsBlank := strings.TrimSpace(cryptoC.Data) == ""
		valueIsBlank := strings.TrimSpace(cryptoC.InnerText()) == ""
		// если не пустые, присваиваем их значение newValue
		if nodeIsBlank && valueIsBlank {
			continue
		}
		// присваивание имени ключа к тегу
		newValue[cryptoC.Data] = cryptoC.InnerText()
		// присваиваем значение Value новому Value
		dynamicValue.Coin = newValue
	}
	return dynamicValue.Coin
}

var iPage = template.Must(template.ParseFiles("index.html"))
func indexHandler(w http.ResponseWriter, r *http.Request) {
    iPage.Execute(w, nil)
}

func coinsHandler(w http.ResponseWriter, r *http.Request) {
		tpl := template.Must(template.ParseFiles("crypto.json"))
		w.Header().Set("Content-Type", "application/json")
		tpl.Execute(w, nil)
	}

func coinsHandler3(w http.ResponseWriter, r *http.Request) {
				tpl3 := template.Must(template.ParseFiles("decrypto.json"))
				w.Header().Set("Content-Type", "application/json")
				tpl3.Execute(w, nil)
			}

func getCourses() {
				for{
				resp, errR := http.Get("https://test.cryptohonest.ru/request-exportxml.xml")
			    if errR != nil{
			         fmt.Println(errR)
			         }
			  //кэширование в файл
			    xmlFILE, errX := os.Create("crypto.xml")
			    if errX != nil {
			    		fmt.Println("Error to create file: ", errX)
			    	}
			    defer xmlFILE.Close()
			    resp.Write(xmlFILE)
					time.Sleep(10*time.Second)
				}
			}

func main() {

  valuesToInput := []DynamicValues{}
	contentWithValues := Contents{valuesToInput}


	go getCourses()

  resp, errR := http.Get("https://test.cryptohonest.ru/request-exportxml.xml")
    if errR != nil{
         fmt.Println(errR)
         }
  //кэширование в файл
    xmlFILE, errX := os.Create("crypto.xml")
    if errX != nil {
    		fmt.Println("Error to create file: ", errX)
    	}
    defer xmlFILE.Close()
    resp.Write(xmlFILE)


	// открытие сохранённого выше файла
	xmlFile, err := os.Open("crypto.xml")
	if err != nil {
		fmt.Println("Error to open file: ")
	}

  //парсинг файла
	doc, erro := xmlquery.Parse(xmlFile)
	//обработка ошибки
	if erro != nil {
		fmt.Println("Error to parcing xml")
	}

	// ищем тег rates
	nodecryptoCs := xmlquery.FindOne(doc, "//rates")

	// ищем тэг валюты, который будет возвращать дерево узлов с данными по валюте
	for node := nodecryptoCs.FirstChild; node != nil; node = node.NextSibling {
		// ищем Node валюты
		cryptoCNode := xmlquery.FindOne(node, "//item")
		if cryptoCNode == nil {
			continue
		}
		dynamicValue := DynamicValues{}
		dynamicValue.enrichDynamicValue(cryptoCNode)
		contentWithValues.AddItemToDynamicList(dynamicValue)
	}

	//конвертация массива в json v1
	encjson, _ := json.Marshal(contentWithValues)
	data := []byte(encjson)
  file, err0 := os.Create("crypto.json")
    if err0 != nil{
        fmt.Println("Unable to create file:", err0)
        os.Exit(1)
    }
		defer file.Close()
    file.Write(data)
    fmt.Println("Done.")
	// вывод на консоль результата
//	fmt.Println(string(encjson))
	fmt.Println("Successfully parsed.xml")
	//закрытие
	defer xmlFile.Close()


//конвертация в json v2
//открытие файла xml
    xmlFile2, err2 := os.Open("crypto.xml")
    if err2 != nil {
        fmt.Println(err2)
      }
      fmt.Println("Successfully Opened crypto.xml")
      defer xmlFile2.Close()

//чтение в byteValue сождержимого xmlFile2
    byteValue, _ := ioutil.ReadAll(xmlFile2)
    var dataS DataS
		//конвертация
    xml.Unmarshal(byteValue, &dataS)
    jsonData, _ := json.Marshal(dataS)

    dataR := []byte(jsonData)
    xmlFILEx, errXs := os.Create("cryptoFormat.txt")
    if errXs != nil {
    		fmt.Println("Erorr to create file: ", errXs)
    	}
    defer xmlFILEx.Close()
    xmlFILEx.Write(dataR)

//encod
	var message = []byte(data)

	flag.Parse()

	key := []byte(*userKey)

	key = append(key, pad...)

	naclKey := new([keySize]byte)
	copy(naclKey[:], key[:keySize])

	nonce := new([nonceSize]byte)

	_ , errz := io.ReadFull(rand.Reader, nonce[:])
	if errz != nil {
		log.Fatalln("Could not read from random:", errz)
	}


	out := make([]byte, nonceSize)
	copy(out, nonce[:])
	out = secretbox.Seal(out, message, nonce, naclKey)

	err = ioutil.WriteFile("encoding.json", out, 0777)
	if err != nil {
		log.Fatalln("Error while writing encrypto.json: ", err)
	}

	fmt.Printf("The encoding.json is: '%s'\n", out)


	// dataRi := []byte(out)
	// xmlFILExi, errXsi := os.Create("cryptoCourses.txt")
	// if errXsi != nil {
	// 		fmt.Println("Erorr to create file: ", errXs)
	// 	}
	// 		xmlFILExi.Write(dataRi)
	// defer xmlFILExi.Close()
	//
	//
	// fmt.Printf("Message encrypted succesfully. Total size is %d bytes,"+
	// 	" of which %d bytes is the message, "+
	// 	"%d bytes is the nonce and %d bytes is the overhead.\n",
	// 	len(out), len(message), nonceSize, secretbox.Overhead)
	// fmt.Printf("The encryption key is: '%s'\n", naclKey[:])
	// fmt.Printf("The nonce is: '%v'\n", nonce[:])



//deco

	in, err := ioutil.ReadFile("encoding.json")
	if err != nil {
		log.Fatalln(err)
	}

	copy(nonce[:], in[:nonceSize])

	message, ok := secretbox.Open(nil, in[nonceSize:], nonce, naclKey)
	if ok {
		err = ioutil.WriteFile("decrypto.json", message, 0777)
		if err != nil {
			log.Fatalln("Error while writing decrypto.json: ", err)
		}

		fmt.Println("Message decrypted successfully.")
		fmt.Printf("The encryption key is: '%s'\n", naclKey[:])
		fmt.Printf("The nonce is: '%v'\n", nonce[:])
		fmt.Printf("The message is: '%s'\n", message)
	} else {
		log.Fatalln("Could not decrypt the message.")
	}

		//создание локального сервера
		r := mux.NewRouter()
		r.HandleFunc("/courses", coinsHandler).Methods("GET")
		r.HandleFunc("/encoding", coinsHandler2).Methods("GET")
		r.HandleFunc("/decoding", coinsHandler3).Methods("GET")
		r.HandleFunc("/", indexHandler)
    http.Handle("/", r)
		http.ListenAndServe(":8181", nil)

}

func coinsHandler2(w http.ResponseWriter, r *http.Request) {
		tpl2 := template.Must(template.ParseFiles("encoding.json"))
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(tpl2)
	}
