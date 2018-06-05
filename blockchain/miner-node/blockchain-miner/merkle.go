package main

import (
	"crypto/sha256"
	"encoding/hex"
)

var (
	lastindex     = []string{}
	Hashresult    string
	templastindex string
)

func Hashing(text string) string { //Gelen parametrelerin Hashini alan fonksiyon
	hasher := sha256.New()
	hasher.Write([]byte(text))
	return hex.EncodeToString(hasher.Sum(nil))
}
func merkleroot(Hashtext []string) string { //Merkle kök değeri hesaplayan fonksiyon
	//Fonksiyon 2 tane slice tanımladık ve bu slice ları sırayla sıfırlayarak dönüşümlü olarak kullandık
	Datatemp := []string{}
	for len(Hashtext) != 1 { //Gelen Hastext slice'sını 1 oluncuya kadar dönen fonksiyon (While kullanımı)
		Datatemp = Datatemp[:0] //Datatemp slice ının içini boşaltıyoruz
		lengthashtex := len(Hashtext)
		if lengthashtex%2 != 0 { //Gelen Hashtext boyutunu kontrol ediyoruz eğer tek sayı ise son indexi lastindex slicenın içine atıyoruz
			lastindex = append(lastindex, Hashtext[lengthashtex-1])
			Hashtext = Hashtext[:len(Hashtext)-1]
		}
		for i := 0; i < len(Hashtext); i += 2 { //Çift haneli olan Hashtext imizi 2 şer li topalayarak hashleyip Datatemp slice ımıza ekliyoruz
			temp := Hashtext[i] + Hashtext[i+1]
			temp = Hashing(temp)
			Datatemp = append(Datatemp, temp)

		}
		if len(Datatemp) == 1 { //Eğer Datatemp sonucu 1 ise elimizdeki hash sonucumuz oluyor döngüden çıkıyoruz
			Hashresult = Datatemp[0]
			Hashresult = Hashing(Hashresult)
			return Hashresult
		}
		if len(Datatemp)%2 != 0 { //Eğer Datatemp slice ı tek sayı ise son indexi alıp last indexe atıyoruz
			lastindex = append(lastindex, Datatemp[len(Datatemp)-1])
			Datatemp = Datatemp[:len(Datatemp)-1]

		}
		Hashtext = Hashtext[:0]                 //Hashtext slice ımızı sıfırlıyoruz
		for i := 0; i < len(Datatemp); i += 2 { //Datatemp slice ı içindeki özetlerimizi 2 li toplayıp hashliyoruz
			temp := Datatemp[i] + Datatemp[i+1]
			temp = Hashing(temp)
			Hashtext = append(Hashtext, temp) //Çıkan sonucu Hashtext e yazıyoruz
		}
		if len(Hashtext) == 1 { //Eğer Hashtextimiz 1 ise elimizdeki sonuc hashımız oluyor
			Hashresult = Hashtext[0]
			Hashresult = Hashing(Hashresult)
			return Hashresult
		}
	}
	if len(lastindex) != 0 { //Lastindeximize boyutu tek sayı olan Slice lardan veri geldimi kontrol ediyoruz
		for item := range lastindex { //Eğer veri geldiyse bunları topluyoruz
			templastindex = templastindex + lastindex[item]
		}
	}
	Hashresult = Hashing(templastindex + Hashresult) //Döngüden çıktıktan sonra Hashresult ve templeindex değerlerimizi toplayarak sonuç hashımızı üretiyoruz
	return Hashresult
}
