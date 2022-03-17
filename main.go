package main

import (
	"fmt"
	"log"
	"os"

	"github.com/alicekaerast/ioffice/lib"
	"github.com/araddon/dateparse"
	"github.com/spf13/viper"
)

func usage() {
	fmt.Printf("Please use one of the following commands:\n\nlist\ncreate <yyyy-mm-dd> [roomID]\ncheckin <reservation ID>\ndelete <reservation ID>")
}

func main() {
	if len(os.Args) < 2 {
		usage()
	} else {
		viper.SetConfigName("ioffice")
		viper.SetConfigType("yaml")
		viper.AddConfigPath("$HOME")
		viper.AddConfigPath("$HOME/.config")
		viper.AddConfigPath(".")
		err := viper.ReadInConfig()
		if err != nil { // Handle errors reading the config file
			panic(fmt.Errorf("Fatal error config file: %w \n", err))
		}
		viper.SetEnvPrefix("ioffice")
		viper.AutomaticEnv()

		username := viper.GetString("username")
		password := viper.GetString("password")
		roomID := viper.GetInt("roomID")
		hostname := viper.GetString("hostname")
		session := viper.GetString("session")

		ioffice := lib.NewIOffice(hostname, username, password, session)

		me := ioffice.GetMe()
		if !ioffice.WasOkay() {
			log.Println("Stopping now as auth failed.  Are you on SSO?  See README.md on how to authenticate.")
			return
		}


		switch os.Args[1] {
		case "list":
			ioffice.ListReservations()
		case "create":
			if len(os.Args) == 2 {
				usage()
			}
			if len(os.Args) == 3 {
				ioffice.CreateReservation(me, roomID, dateparse.MustParse(os.Args[2]))
			}
			if len(os.Args) == 4 {
				room := ioffice.GetRoom(os.Args[2])
				ioffice.CreateReservation(me, room.ID, dateparse.MustParse(os.Args[2]))
			}
		case "checkin":
			reservationID := os.Args[2]
			ioffice.CheckIn(reservationID)
		case "cancel":
			reservationID := os.Args[2]
			ioffice.CancelReservation(reservationID)
		default:
			usage()
		}
	}
}
