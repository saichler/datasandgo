package main

import (
	"github.com/nats-io/go-nats"
	"fmt"
	"strconv"
	"time"
	"log"
)

func main(){
	count := 0
	//ncAside, err := nats.Connect("tls://10.157.157.6:4443", nats.RootCAs("/home/ubuntu/ca.pem"))
	ncAside, err := nats.Connect("nats://10.157.157.6:4222")
	if err!=nil {
		log.Fatal("Failed to connect to NATS\n", err)
	}
	//ncZside, _ := nats.Connect("tls://10.157.157.6:4443", nats.RootCAs("/home/ubuntu/ca.pem"))
	ncZside, _ := nats.Connect("nats://10.157.157.6:4222")

	ncAside.Publish("natstest", []byte("Hello World"))

	start := time.Now().UnixNano()
	end := time.Now().UnixNano()

	go ncZside.Subscribe("natstest", func(m *nats.Msg) {
		count++
		end = time.Now().UnixNano()
	})

	for i:=0;i<100000;i++ {
		time.Sleep(2000)
		ncAside.Publish("natstest", []byte("Hello World, my name is Sharon and I am testing NATS TLS Performance"))
	}

	for count<100000 {
		time.Sleep(time.Millisecond*500)
		fmt.Println("c="+strconv.Itoa(count))
	}

	took := end-start

	fmt.Println("Count="+strconv.Itoa(count)+" time="+strconv.Itoa(int(took)))

	//With TLS
	tlsAvrg := (float64(14422031333) +
		float64(14345262499) +
		float64(14541710038) +
		float64(15111623073) +
		float64(14819799038) +
		float64(15442325170) +
		float64(14992261313) +
		float64(14923149471) +
		float64(14961652224) +
		float64(14425457414))/10

	//14422031333
	//14345262499
	//14541710038
	//15111623073
	//14819799038
	//15442325170
	//14992261313
	//14923149471
	//14961652224
	//14425457414


	//15442325170 -
	//13970837371 =
	// 1471487799

	noTlsAvrg := (float64(14324694133) +
		float64(14794422933) +
		float64(13970837371) +
		float64(14957771203) +
		float64(14439994881) +
		float64(14582934783) +
		float64(14552152269) +
		float64(14089358012) +
		float64(14038611229) +
		float64(14280897982))/10

	//Wintout TLS
	//14324694133
	//14794422933
	//13970837371
	//14957771203
	//14439994881
	//14582934783
	//14552152269
	//14089358012
	//14038611229
	//14280897982

	p := noTlsAvrg/tlsAvrg*100
	diff := int64(tlsAvrg)-int64(noTlsAvrg)
	fmt.Println("Diff in nanoseconds:\n",diff,"\n", p)
}

