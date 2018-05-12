package main

import (
	"fmt"
	"math/rand"
	"time"
	"sync"
)

// Define Cryptographer Actor
type Crypto struct {
	coin    bool	// Result of coin flip
	paid    bool	// Is the cryptographer payer?
	xorval  bool	// Outcome of comparision
	identity	string	// Identity of cryptographer
}

//Define NSA Actor
type NSA struct {
	paid bool	// Is the NSA paying?
	identity	string
}

//Waitgroup to sync processes
var wg[4] sync.WaitGroup

//Simulates flipping a coin
func (c *Crypto) Flip(channel chan bool, channel_0 chan bool) {
	random := rand.Int()
	if random%2 == 0 { 
		c.coin = true
	} else {
		c.coin = false
	}
	fmt.Println(c.identity," coin value: ",c.coin)
	channel <- c.coin
	channel_0 <- c.coin
	close(channel)
	wg[0].Done()
}

// Compare function
func (c *Crypto) compare(left_coin bool, channel_out chan bool, channel_0 chan bool) {
	if c.paid {
		c.xorval = reverse_xor(c.coin, left_coin)
	} else {
		c.xorval = xor(c.coin, left_coin)
	}
	if c.xorval {
		fmt.Println(c.identity," declares outcome: Different")
	} else {
		fmt.Println(c.identity," declares outcome: Same")
	}
	channel_out <- c.xorval
	channel_0 <- c.xorval
	wg[1].Done()
}

func restaurant_owner(a bool, b bool, c bool) {
	var same, diff int
	var coin[3] bool
	coin[0]=a
	coin[1]=b
	coin[2]=c
	same = 0
	diff = 0
	for i :=0; i<3; i++ {
		if coin[i] {
			diff = diff + 1
		} else {
			same = same + 1
		}
	} 	
	fmt.Println("Result Count -> ||Same:",same,"||  ||Different:",diff,"||\n")	
	if (diff % 2 != 0){ 
		fmt.Println("*** Odd count of \"Different\" Uttered. So.. ***")
		fmt.Println("=> Restaurant Owner declares a cryptographer paid !!")
	} else {
		fmt.Println("*** Even count of \"Different\" Uttered. So.. ***")
		fmt.Println("=> Restaurant Owner declares the NSA Paid !!")
	}
	wg[2].Done()
}

func crypt0(aCoin bool, bCoin bool, cCoin bool, axor bool, bxor bool, cxor bool ) {
	fmt.Println("-----------------------------------")	
	fmt.Println("     ==> Cryptographer A <==       ")
	fmt.Println("-----------------------------------")
	fmt.Println("Computed Result: A xor C =>",xor(aCoin, cCoin))
	fmt.Println("Declared Result by A     =>",axor)

	fmt.Println("\n-----------------------------------")	
	fmt.Println("        Cryptographer B        ")
	fmt.Println("-----------------------------------")
	fmt.Println("Computed Result: B xor A =>",xor(bCoin, aCoin))
	fmt.Println("Declared Result by B     =>",bxor)

	fmt.Println("\n-----------------------------------")	
	fmt.Println("        Cryptographer C        ")
	fmt.Println("-----------------------------------")
	fmt.Println("Computed Result: C xor B =>",xor(cCoin, bCoin))
	fmt.Println("Declared Result by C     =>",cxor)
	fmt.Println("-----------------------------------\n") 
	if reverse_xor(aCoin, cCoin) == axor {
		fmt.Println("*** Cryptographer A said opposite --> Cryptographer A Paid! ***")
	} else if reverse_xor(bCoin, aCoin) == bxor {
		fmt.Println("*** Cryptographer B said opposite --> Cryptographer B Paid!")
	} else if reverse_xor(cCoin, bCoin) == cxor {
		fmt.Println("*** Cryptographer C said opposite --> Cryptographer C Paid!")
	} else {
		fmt.Println("*** Nobody said opposite so The NSA paid! ***\n\n")
	}
	wg[3].Done()
}

// XOR operation
func xor(a bool, b bool) bool {
	return a != b
}

// Reverse XOR operation
func reverse_xor(a bool, b bool) bool {
	return !(a != b)
}

func main() {

	// Channels for communication between cryptographers
	channel_A := make(chan bool, 1) // Channel for Crytographer A's coin
	channel_B := make(chan bool, 1) // Channel for Crytographer B's coin
	channel_C := make(chan bool, 1) // Channel for Crytographer C's coin
	channel_O := make(chan bool, 3) // Channel for Crytographers outcome
	channel_0A := make(chan bool, 2) // Channel shared by Crytographers 0,A
	channel_0B := make(chan bool, 2) // Channel shared by Crytographers 0,B
	channel_0C := make(chan bool, 2) // Channel shared by Crytographers 0,A
	//get random seed as current time
	rand.Seed(time.Now().UTC().UnixNano())
	
	// Set defaults
	a := Crypto{paid: false, identity: "Cryptographer A"}
	b := Crypto{paid: false, identity: "Cryptographer B"}
	c := Crypto{paid: false, identity: "Cryptographer C"}
	n := NSA{paid: false, identity: "NSA"}

	//determine payer randomly
	payer := rand.Int() % 4
	fmt.Println("-------------------------------------")
	fmt.Println("  ==> Randomly Selecting Payer <== 	   ")
	fmt.Println("-------------------------------------")
	if payer == 0 {
		n.paid = true
		fmt.Println("  Selected Payer:", n.identity)
	}
	if payer == 1 {
		a.paid = true
		fmt.Println("  Selected Payer:", a.identity)
	}
	if payer == 2 {
		b.paid = true
		fmt.Println("  Selected Payer:", b.identity)
	}
	if payer == 3 {
		c.paid = true
		fmt.Println("  Selected Payer:", c.identity)
	}
	fmt.Println("-------------------------------------\n")

	// All crytographers flip coins
	fmt.Println("-----------------------------------------------")
	fmt.Println(" ==> Flip results visible to Crytographer0 <==  ")
	fmt.Println("-----------------------------------------------\n")
	wg[0].Add(3)
	go a.Flip(channel_A, channel_0A)
	go b.Flip(channel_B, channel_0B)
	go c.Flip(channel_C, channel_0C)
	wg[0].Wait()
	fmt.Println("---------------------------------------\n")
	// All crytographers compare coins and broadcast
	fmt.Println("------------------------------------------")
	fmt.Println(" ==> Cryptographers declare outcomes <==   ")	
	fmt.Println("------------------------------------------")
	fmt.Println("** Visible to: Observer, Restaurant Owner, Cryptographer 0,A,B,C **\n")	
	wg[1].Add(3)
	go a.compare(<-channel_C, channel_O, channel_0A)
	go b.compare(<-channel_A, channel_O, channel_0B)
	go c.compare(<-channel_B, channel_O, channel_0C)
	wg[1].Wait()
	fmt.Println("---------------------------------------\n")	
	fmt.Println("------------------------------------------")
	fmt.Println(" ==> Restaurant Owner computes payer <==    ")
	fmt.Println("------------------------------------------")
	//Restaurant Owner
	wg[2].Add(1)
	go restaurant_owner(<-channel_O,<-channel_O,<-channel_O)
	wg[2].Wait()
	fmt.Println("---------------------------------------\n")
	fmt.Println("---------------------------------------")
	fmt.Println(" ==> Cryptographer0 verifies payer <== ")
	fmt.Println("---------------------------------------")	
	//Cryptographer0
	wg[3].Add(1)
	go crypt0(<-channel_0A , <-channel_0B, <-channel_0C, <-channel_0A , <-channel_0B, <-channel_0C)
	wg[3].Wait()
}
