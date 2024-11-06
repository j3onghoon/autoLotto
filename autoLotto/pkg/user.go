package pkg

import (
	"autoLotto/ent"
	"autoLotto/ent/user"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"golang.org/x/net/publicsuffix"
	"golang.org/x/text/encoding/korean"
	"golang.org/x/text/transform"
	"io"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strconv"
	"strings"

	_ "github.com/lib/pq"
)

type Lotto645Ticket struct {
	ArrGameChoiceNum string `json:"arrGameChoiceNum,omitempty"`
	GenType          string `json:"genType"`
	Alphabet         string `json:"alpabet"`
}

type LotteryUser struct {
	Client          *http.Client
	User            *ent.User
	JSESSIONID      string
	Round           string
	Deposit         string
	Lotto645Tickets []Lotto645Ticket
}

func encodeToEUCKR(input string) (string, error) {
	var euckrBuilder strings.Builder
	writer := transform.NewWriter(&euckrBuilder, korean.EUCKR.NewEncoder())

	_, err := writer.Write([]byte(input))
	if err != nil {
		return "", err
	}

	err = writer.Close()
	if err != nil {
		return "", err
	}

	return euckrBuilder.String(), nil
}

func GetUserByID(ctx context.Context, userID string) (*ent.User, error) {
	user, err := DbClient.User.Get(ctx, userID)
	if err != nil {
		fmt.Println("failed querying user: %v\n", err)
		return nil, err
	}
	return user, nil
}

func createUser(ctx context.Context, id string, password string) (*ent.User, error) {
	user, err := DbClient.User.Create().
		SetID(id).
		SetWeeklyPurchase(false).
		SetPassword(password).
		Save(ctx)
	if err != nil {
		fmt.Println("failed creating user: %v\n", err)
	}
	log.Println("User created: %v\n", user)
	return user, nil
}

func UpdateUser(ctx context.Context, dbUser *ent.User) error {
	// todo
	// 아래 2가지 방법 중 하나로
	// 1. 로또 구매를 하는 동안에는 update를 못하게 막기
	// 2. db 스냅샷을 만들고 그 스냅샷을 이용하게 만들기
	newDbUser, err := DbClient.User.UpdateOneID(dbUser.ID).
		SetWeeklyPurchase(dbUser.WeeklyPurchase).
		SetPassword(dbUser.Password).
		SetTicketCount(dbUser.TicketCount).
		Save(ctx)
	if err != nil {
		fmt.Println("failed updating user: %v\n", err)
	}
	log.Println("User updated: %v\n", newDbUser)
	return nil
}

func DeleteUser(ctx context.Context, dbUser *ent.User) error {
	err := DbClient.User.DeleteOneID(dbUser.ID).Exec(ctx)
	if err != nil {
		fmt.Println("failed deleting user: %v\n", err)
	}
	log.Println("User deleted: %v\n", dbUser.ID)
	return err
}

func userExists(ctx context.Context, userID string) (bool, error) {
	exists, err := DbClient.User.Query().Where(user.IDEQ(userID)).Exist(ctx)
	if err != nil {
		return false, err
	}
	return exists, nil

}

func setDefaultConfig(lotteryUser *LotteryUser) error {
	mainUrl := "https://dhlottery.co.kr/common.do?method=main"
	req, err := http.NewRequest("GET", mainUrl, nil)
	setHeader(req, *lotteryUser)
	resp, err := lotteryUser.Client.Do(req)
	defer resp.Body.Close()
	reader := transform.NewReader(resp.Body, korean.EUCKR.NewDecoder())
	decodedBytes, err := io.ReadAll(reader)
	htmlstring := string(decodedBytes)

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(htmlstring))
	if err != nil {
		log.Fatal(err)
	}
	var deposit string
	lottoNumber := doc.Find("strong#lottoDrwNo").Text()
	doc.Find("li.money").Each(func(i int, s *goquery.Selection) {
		spanText := s.Find("span").Text()
		if spanText == "예치금" {
			s.Find("a[href='/myPage.do?method=depositListView']").Each(func(j int, a *goquery.Selection) {
				deposit = a.Find("strong").Text()
			})
		}
	})
	deposit = strings.ReplaceAll(deposit, "원", "")
	deposit = strings.ReplaceAll(deposit, ",", "")

	if deposit == "" {
		return errors.New("로그인에 실패했습니다.")
	}
	round, err := strconv.Atoi(lottoNumber)
	lotteryUser.Round = strconv.Itoa(round + 1)
	lotteryUser.Deposit = deposit

	return err
}

func BuyLottery(lotteryUser LotteryUser) error {
	getSocketUrl := "https://ol.dhlottery.co.kr/olotto/game/egovUserReadySocket.json"
	req, err := http.NewRequest("POST", getSocketUrl, nil)
	resp, err := lotteryUser.Client.Do(req)

	defer resp.Body.Close()
	var result map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		fmt.Println("Error parsing JSON:", err)
		return err
	}
	buyLotteryURL := "https://ol.dhlottery.co.kr/olotto/game/execBuy.do"

	data := url.Values{}
	data.Set("round", lotteryUser.Round)
	data.Set("direct", result["ready_ip"].(string))
	data.Set("nBuyAmount", strconv.Itoa(1000*len(lotteryUser.Lotto645Tickets)))
	data.Set("param", transformLottoTicket(&lotteryUser))
	//data.Set("param", `[{"arrGameChoiceNum":"23","genType":"2","alpabet":"A"}]`)
	data.Set("gameCnt", "1")
	euckrEncodedData, err := encodeToEUCKR(data.Encode())
	req, err = http.NewRequest("POST", buyLotteryURL, strings.NewReader(euckrEncodedData))
	setHeader(req, lotteryUser)
	resp, err = lotteryUser.Client.Do(req)
	defer resp.Body.Close()

	//reader := transform.NewReader(resp.Body, korean.EUCKR.NewDecoder())
	//var bodyBuilder strings.Builder
	//_, err = io.Copy(&bodyBuilder, reader)
	//
	//fmt.Println(bodyBuilder.String())
	//err = json.NewDecoder(resp.Body).Decode(&result)
	//if err != nil {
	//	fmt.Println("Error parsing JSON:", err)
	//	return err
	//}
	return err
}

func transformLottoTicket(lotteryUser *LotteryUser) string {
	order := 'A'
	for i := range lotteryUser.Lotto645Tickets {
		lotteryUser.Lotto645Tickets[i].Alphabet = string(order)
		order += 1
	}
	jsonData, err := json.Marshal(lotteryUser.Lotto645Tickets)
	if err != nil {
		fmt.Println("Error parsing JSON:", err)
	}
	return string(jsonData)
}

func MakeLotto645Tickets(lotteryUser *LotteryUser, mode string, numbers string, count int) {
	if mode == "0" {
		for i := 0; i < count; i++ {
			lotteryUser.Lotto645Tickets = append(lotteryUser.Lotto645Tickets, Lotto645Ticket{ArrGameChoiceNum: "", GenType: mode})
		}
	} else if mode == "1" {
		for i := 0; i < count; i++ {
			lotteryUser.Lotto645Tickets = append(lotteryUser.Lotto645Tickets, Lotto645Ticket{ArrGameChoiceNum: numbers, GenType: mode})
		}
	} else if mode == "2" {
		for i := 0; i < count; i++ {
			lotteryUser.Lotto645Tickets = append(lotteryUser.Lotto645Tickets, Lotto645Ticket{ArrGameChoiceNum: numbers, GenType: mode})
		}
	}
}

func setHeader(req *http.Request, lotteryUser LotteryUser) {
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.77 Safari/537.36")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Cache-Control", "max-age=0")
	req.Header.Set("sec-ch-ua", `" Not;A Brand";v="99", "Google Chrome";v="91", "Chromium";v="91"`)
	req.Header.Set("sec-ch-ua-mobile", "?0")
	req.Header.Set("Upgrade-Insecure-Requests", "1")
	req.Header.Set("Origin", "https://dhlottery.co.kr")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9")
	req.Header.Set("Referer", "https://dhlottery.co.kr")
	req.Header.Set("Sec-Fetch-Site", "same-site")
	req.Header.Set("Sec-Fetch-Mode", "navigate")
	req.Header.Set("Sec-Fetch-User", "?1")
	req.Header.Set("Sec-Fetch-Dest", "document")
	req.Header.Set("Accept-Language", "ko,en-US;q=0.9,en;q=0.8,ko-KR;q=0.7")
	req.Header.Set("X-Requested-With", "XMLHttpRequest")
	req.Header.Set("Cookie", "JSESSIONID="+lotteryUser.JSESSIONID)
}

func login(lotteryUser LotteryUser, id string, password string) error {
	data := url.Values{}
	data.Set("userId", id)
	data.Set("password", password)
	data.Set("checkSave", "on")
	data.Set("returnUrl", "https://www.dhlottery.co.kr/common.do?method=main")
	euckrEncodedData, err := encodeToEUCKR(data.Encode())

	loginURL := "https://www.dhlottery.co.kr/userSsl.do?method=login"
	req, err := http.NewRequest("POST", loginURL, strings.NewReader(euckrEncodedData))
	setHeader(req, lotteryUser)
	resp, err := lotteryUser.Client.Do(req)
	defer resp.Body.Close()

	//reader := transform.NewReader(resp.Body, korean.EUCKR.NewDecoder())
	//var bodyBuilder strings.Builder
	//_, err = io.Copy(&bodyBuilder, reader)
	//fmt.Println("resp.body: ", bodyBuilder.String())
	//fmt.Println("resp.status: ", resp.Status)
	return err
}

func setDefaultSession(lotteryUser *LotteryUser) error {
	defaultSessionUrl := "https://dhlottery.co.kr/gameResult.do?method=byWin&wiselog=H_C_1_1"
	req, err := http.NewRequest("GET", defaultSessionUrl, nil)

	client := lotteryUser.Client
	resp, err := client.Do(req)
	defer resp.Body.Close()

	for _, cookie := range resp.Cookies() {
		if cookie.Name == "JSESSIONID" {
			lotteryUser.JSESSIONID = cookie.Value
		}
	}
	return err
}

func setUserInfo(lotteryUser *LotteryUser, id string, password string) error {
	err := login(*lotteryUser, id, password)
	err = setDefaultConfig(lotteryUser)
	if err != nil {
		return err
	}

	exists, err := userExists(context.Background(), id)
	if err != nil {
		return err
	}

	if exists {
		dbUser, err := GetUserByID(context.Background(), id)
		if err != nil {
			return err
		}
		lotteryUser.User = dbUser
		lotteryUser.User.Password = password
	} else {
		dbUser, err := createUser(context.Background(), id, password)
		if err != nil {
			return err
		}
		lotteryUser.User = dbUser
	}

	UpdateUser(context.Background(), lotteryUser.User)
	return nil
}

func SetLottoUserInfo(id string, password string) (*LotteryUser, error) {
	lotteryUser := LotteryUser{}
	jar, _ := cookiejar.New(&cookiejar.Options{
		PublicSuffixList: publicsuffix.List,
	})
	lotteryUser.Client = &http.Client{
		Jar: jar,
	}
	err := setDefaultSession(&lotteryUser)
	err = setUserInfo(&lotteryUser, id, password)
	return &lotteryUser, err
}
