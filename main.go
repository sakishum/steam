/*!
 * Note!  This is a "test case", it's used for ease of development
 * This will turn into a library.  */
package main

import (
	"log"
	"os"
	//	"time"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	twoFactorCode, err := GenerateTwoFactorCode(os.Getenv("steamSharedSecret"))
	if err != nil {
		log.Fatal(err)
	}

	community := Community{}
	if err := community.login(os.Getenv("steamAccount"), os.Getenv("steamPassword"), twoFactorCode); err != nil {
		log.Fatal(err)
	}
	log.Print("Login successful")

	key, err := community.getWebAPIKey()
	if err != nil {
		log.Fatal(err)
	}
	log.Print("Key: ", key)

	sid := SteamID(76561198078821986)
	inven, err := community.GetInventory(&sid, 730, 2, false)
	if err != nil {
		log.Fatal(err)
	}

	for _, item := range inven {
		log.Printf("Item: %s = %d\n", item.MarketHashName, item.AssetID)
	}

	marketPrices, err := community.GetMarketItemPriceHistory(730, "P90 | Asiimov (Factory New)")
	if err != nil {
		log.Fatal(err)
	}

	for k, v := range marketPrices {
		log.Printf("%d: %s -> %.2f (%s)\n", k, v.Date, v.Price, v.Count)
	}

	/*
		sent, _, err := community.GetTradeOffers(TradeFilterSentOffers|TradeFilterRecvOffers, time.Now())
		if err != nil {
			log.Fatal(err)
		}

		var receiptID uint64
		for k := range sent {
			offer := sent[k]
			var sid SteamID
			sid.Parse(offer.Partner, AccountInstanceDesktop, AccountTypeIndividual, UniversePublic)

			if receiptID == 0 && len(offer.ReceiveItems) != 0 && offer.State == TradeStateAccepted {
				receiptID = offer.ReceiptID
			}

			log.Printf("Offer id: %d, Receipt ID: %d", offer.ID, offer.ReceiptID)
			log.Printf("Offer partner SteamID 64: %d", uint64(sid))
		}

		items, err := community.GetTradeReceivedItems(receiptID)
		if err != nil {
			log.Fatal(err)
		}

		for _, item := range items {
			log.Printf("New asset id: %d", item.AssetID)
		}
	*/

	key, err = GenerateConfirmationCode(os.Getenv("steamIdentitySecret"), "conf")
	if err != nil {
		log.Fatal(err)
	}

	confirmations, err := community.GetConfirmations(key)
	if err != nil {
		log.Fatal(err)
	}

	for i := range confirmations {
		c := confirmations[i]
		log.Printf("Confirmation ID: %d, Key: %d\n", c.ID, c.Key)
		log.Printf("-> Title %s\n", c.Title)
		log.Printf("-> Receiving %s\n", c.Receiving)
		log.Printf("-> Since %s\n", c.Since)

		key, err = GenerateConfirmationCode(os.Getenv("steamIdentitySecret"), "details")
		if err != nil {
			log.Fatal(err)
		}

		tid, err := community.GetConfirmationOfferID(key, c.ID)
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("-> OfferID %d\n", tid)

		key, err = GenerateConfirmationCode(os.Getenv("steamIdentitySecret"), "allow")
		err = community.AnswerConfirmation(c, key, "allow")
		if err != nil {
			log.Fatal(err)
		}

		log.Printf("Accepted %d\n", c.ID)
	}

	log.Println("Bye!")
}
