package main

import (
	"encoding/gob"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"sort"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (s *Service) MakeNotice(msgText string, userID int64, msgID int) {
	msg := tgbotapi.NewMessage(userID, msgText, "MarkdownV2", true)
	rep, err := s.bot.Send(msg)
	if err != nil {
		s.ReportToAdmin(err.Error())
	}

	if msgID != 0 {
		s.DeleteOldMsg(userID, msgID)
	}

	s.UpdateUserOldMsg(userID, rep.MessageID)
	s.messageCleaner(userID, msgID)
}

func (s *Service) BackupCache(filename string) {
	var file *os.File
	var err error
	if file, err = os.Create(fmt.Sprintf("%s.gob", filename)); err != nil {
		s.ReportToAdmin(err.Error())
	}
	defer file.Close()
	encoder := gob.NewEncoder(file)

	if err = encoder.Encode(s.Cache); nil != err {
		s.ReportToAdmin(err.Error())
	}
}

func (s *Service) LoadBackup(filename string) {
	var file *os.File
	var err error
	if file, err = os.Open(fmt.Sprintf("%s.gob", filename)); err != nil {
		s.ReportToAdmin(err.Error())
	}
	defer file.Close()
	decoder := gob.NewDecoder(file)
	if err = decoder.Decode(s.Cache); err != nil {
		s.ReportToAdmin(err.Error())
	}
}

func (s *Service) ReportToAdmin(msgText string) {
	if _, err := s.SendMessage(s.SuperAdmin, msgText, "", false); err != nil {
		fmt.Printf("errors %s", err.Error())
	}

}

func (s *Service) DownloadFile(filepath string, url string) (err error) {

	// Create the file
	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Check server response
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("bad status: %s", resp.Status)
	}

	// Writer the body to file
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}

	return nil
}

func (s *Service) Exists(name string) bool {
	_, err := os.Stat(name)
	if err == nil {
		return true
	}
	if errors.Is(err, os.ErrNotExist) {
		return false
	}
	return false
}

func (s *Service) AskCache() {
	var FileID string
	s.ReportToAdmin("please send to me a last updated Cache")
	for update := range s.updates {
		if update.Message != nil {
			if s.SuperAdmin == update.Message.From.ID {
				if update.Message.Text == "!" || update.Message.Command() == "start" {
					s.messageCleaner(update.Message.Chat.ID, update.Message.MessageID)
					return
				}
				if update.Message.Document != nil {
					FileID = update.Message.Document.FileID
					break
				}
			}
			s.messageCleaner(update.Message.Chat.ID, update.Message.MessageID)
		}
	}
	if fileURL, err := s.bot.GetFileDirectURL(FileID); err == nil {
		if err = s.DownloadFile("cache.gob", fileURL); err == nil {
			s.LoadBackup("cache")
			s.ReportToAdmin("successfully Updated Cache")
			return
		}
		s.ReportToAdmin(err.Error())
	}
}

func PrettyPrint(data ...interface{}) {
	fmt.Println("[")
	for i, d := range data {

		var p []byte
		//    var err := error
		p, err := json.MarshalIndent(d, "", "\t")
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Printf("%s \n", p)
		if i+1 != len(data) {
			fmt.Println(",")
		}
	}
	fmt.Println("]")
}

func (s *Service) Doctor(sec int64) {
	for {
		time.Sleep(time.Duration(sec) * time.Second)
		s.BackupCache("cache")
		fileBackUp := tgbotapi.NewDocumentChannel("-1001777959481", tgbotapi.FilePath("cache.gob"))
		if _, err := s.bot.Send(fileBackUp); err != nil {
			s.ReportToAdmin(err.Error())
		}
	}
}

func (s *Service) Dispatcher(sec int64) {
	for {
		time.Sleep(time.Duration(sec) * time.Second)
		s.Request = make(map[int64]int)
	}
}

func (s *Service) Timer() {
	var conn int8
	for {
		time.Sleep(time.Second)
		conn++
		println(conn, " Second")
		if conn == 15 {
			println("Clear Request")
			return
		}
	}
}

func (s *Service) Leaderboard(sec int64) {
	for {
		time.Sleep(time.Duration(sec) * time.Second)
	}
}

func (s *Service) RefreshLeaderboard() {
	s.mu.Lock()
	defer s.mu.Unlock()

	totalWord := float32(len(s.WordList))
	for _, u := range s.Users {
		recordGarde := float32(len(u.Datasets)) / totalWord
		confirmedGarde := float32(u.Confirmed) / totalWord
		invitionGarde := float32(len(u.Invited)) / float32(s.HighestInvention)
		u.Avarage = (recordGarde + confirmedGarde + invitionGarde) * 100 / 3
	}

	rank := make(RankList, len(s.Users))
	i := 0

	for k, u := range s.Users {
		rank[i] = Pair{k, u.Avarage}
		i++
	}

	sort.Sort(rank)
	s.RankList = rank

	for newRank, ranks := range s.RankList {
		s.Users[ranks.Key].Rank = newRank + 1
	}
}

// func timeTrack(start time.Time, name string) {
// 	elapsed := time.Since(start)
// 	log.Printf("%s took %s", name, elapsed)
// }
