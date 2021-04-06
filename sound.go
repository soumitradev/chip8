package main

func playSound() {
	if soundTimer > 0 {
		playFreq(440)
	} else {
		stopFreq()
	}
}

func playFreq(freq int) {

}

func stopFreq() {

}
