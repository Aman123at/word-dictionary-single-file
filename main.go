package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

const (
	fileName = "dictionary.txt"
)

type Dictionary struct {
	word       string
	definition string
}

var wordOffsets = map[string]int64{}

func createDictionary() error {
	file, err := os.Create(fileName)
	if err != nil {
		return err
	}
	defer file.Close()
	writer := bufio.NewWriter(file)

	// Sample dictionary entries
	entries := []Dictionary{
		{"apple", "A round fruit with red or green skin and white flesh."},
		{"book", "A written or printed work consisting of pages glued or sewn together along one side."},
		{"computer", "An electronic device for storing and processing data."},
		{"dog", "A domesticated carnivorous mammal that typically has a long snout, an acute sense of smell, and a barking, howling, or whining voice."},
		{"elephant", "A very large plant-eating mammal with a prehensile trunk, long curved ivory tusks, and large ears."},
		{"flower", "The seed-bearing part of a plant, consisting of reproductive organs surrounded by petals and sepals."},
		{"guitar", "A stringed musical instrument with a fretted fingerboard, typically incurved sides, and six or twelve strings."},
		{"house", "A building for human habitation, especially one that is lived in by a family or small group of people."},
		{"internet", "A global computer network providing a variety of information and communication facilities."},
		{"jacket", "An outer garment extending either to the waist or the hips, typically having sleeves and fastening down the front."},
		{"kitchen", "A room or area where food is prepared and cooked."},
		{"lamp", "A device for giving light, either one consisting of an electric bulb together with its holder and shade or cover, or one burning gas or oil."},
		{"mountain", "A large natural elevation of the earth's surface rising abruptly from the surrounding level."},
		{"notebook", "A small book with blank or ruled pages for writing notes in."},
		{"ocean", "A very large expanse of sea, in particular each of the main areas into which the sea is divided geographically."},
		{"pencil", "An instrument for writing or drawing, consisting of a thin stick of graphite or a similar substance enclosed in a long thin piece of wood or fixed in a metal or plastic case."},
		{"queen", "The female ruler of an independent state, especially one who inherits the position by right of birth."},
		{"river", "A large natural stream of water flowing in a channel to the sea, a lake, or another such stream."},
		{"sun", "The star around which the earth orbits."},
		{"telephone", "A system for transmitting voices over a distance using wire or radio, by converting acoustic vibrations to electrical signals."},
		{"umbrella", "A device consisting of a circular canopy of cloth on a folding metal frame supported by a central rod, used as protection against rain or sometimes sun."},
		{"violin", "A stringed musical instrument of treble pitch, played with a horsehair bow."},
		{"window", "An opening in the wall or roof of a building or vehicle that is fitted with glass or other transparent material in a frame to admit light or air and allow people to see out."},
		{"xylophone", "A musical instrument played by striking a series of wooden bars of graduated length with small wooden or plastic mallets."},
		{"yoga", "A Hindu spiritual and ascetic discipline, a part of which, including breath control, simple meditation, and the adoption of specific bodily postures, is widely practiced for health and relaxation."},
		{"zebra", "An African wild horse with black-and-white stripes and an erect mane."},
		{"airplane", "A powered flying vehicle with fixed wings and a weight greater than that of the air it displaces."},
		{"banana", "A long curved fruit which grows in clusters and has soft pulpy flesh and yellow skin when ripe."},
		{"camera", "A device for recording visual images in the form of photographs, film, or video signals."},
		{"dolphin", "A small gregarious toothed whale that typically has a beaklike snout and a curved fin on the back."},
		{"eagle", "A large bird of prey with a massive hooked bill and long broad wings."},
		{"forest", "A large area covered chiefly with trees and undergrowth."},
		{"galaxy", "A system of millions or billions of stars, together with gas and dust, held together by gravitational attraction."},
		{"helicopter", "A type of aircraft which derives both lift and propulsion from one or more sets of horizontally revolving overhead rotors."},
		{"island", "A piece of land surrounded by water."},
		{"jungle", "An area of land overgrown with dense forest and tangled vegetation, typically in the tropics."},
		{"kangaroo", "A large plant-eating marsupial with a long powerful tail and strongly developed hindlimbs that enable it to travel by leaping."},
		{"lighthouse", "A tower or other structure containing a beacon light to warn or guide ships at sea."},
		{"microscope", "An optical instrument used for viewing very small objects, such as mineral samples or animal or plant cells."},
		{"necklace", "An ornamental chain or string of beads, jewels, or links worn around the neck."},
		{"octopus", "A cephalopod mollusc with eight sucker-bearing arms, a soft body, strong beaklike jaws, and no internal shell."},
		{"pyramid", "A monumental structure with a square or triangular base and sloping sides that meet in a point at the top."},
		{"quilt", "A warm bed covering made of padding enclosed between layers of fabric and kept in place by lines of stitching."},
		{"rainbow", "An arch of colors visible in the sky, caused by the refraction and dispersion of the sun's light by rain or other water droplets in the atmosphere."},
		{"satellite", "An artificial body placed in orbit around the earth or moon or another planet in order to collect information or for communication."},
		{"tornado", "A mobile, destructive vortex of violently rotating winds having the appearance of a funnel-shaped cloud and advancing beneath a large storm system."},
		{"unicorn", "A mythical animal typically represented as a horse with a single straight horn projecting from its forehead."},
		{"volcano", "A mountain or hill, typically conical, having a crater or vent through which lava, rock fragments, hot vapor, and gas are or have been erupted from the earth's crust."},
		{"waterfall", "A cascade of water falling from a height, formed when a river or stream flows over a precipice or steep incline."},
	}

	// calculate total header offset
	totalHeaderOffset := int(0)
	for _, entry := range entries {
		totalHeaderOffset += int(len(fmt.Sprintf("%s=>\n", entry.word)))
	}

	// estimate total characters of every offset value combined , take 1000 bits for safer side, bounding header withing range
	// taking 64 bit for offset
	totalHeaderOffset += len(entries)*64 + 1000

	// intialize definition offset from total header offset, so we can write word:definition after headers
	defOffsets := totalHeaderOffset
	for _, entry := range entries {
		fmt.Fprintf(writer, "%s=>%d\n", entry.word, defOffsets)
		defOffsets += len(fmt.Sprintf("%s:%s\n", entry.word, entry.definition))
	}
	// flush the header content
	writer.Flush()

	// seek to the offset where header ends and body content starts
	file.Seek(int64(totalHeaderOffset), os.SEEK_SET)

	// initiate new writer for body content
	writer = bufio.NewWriter(file)
	for _, entry := range entries {
		fmt.Fprintf(writer, "%s:%s\n", entry.word, entry.definition)
	}

	// flush the body content into file
	return writer.Flush()
}

func lookupWord(word string) (string, error) {
	file, err := os.Open(fileName)
	if err != nil {
		return "", err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var offset int64

	// scan headers and extract offset for particular word
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.SplitN(line, "=>", 2)
		if len(parts) != 2 {
			break
		}
		if parts[0] == word {
			fmt.Sscanf(parts[1], "%d", &offset)
			break
		}
	}

	// edge case, no word found in dictionary
	if offset == 0 {
		return "", fmt.Errorf("word not found")
	}

	// seek to found offset to get definition
	_, err = file.Seek(offset, os.SEEK_SET)
	if err != nil {
		return "", err
	}

	scanner = bufio.NewScanner(file)
	// extract definition from word, extraction logic here
	if scanner.Scan() {
		line := scanner.Text()
		parts := strings.SplitN(line, ":", 2)
		if len(parts) == 2 && parts[0] == word {
			return parts[1], nil
		}
	}

	return "", fmt.Errorf("error reading definition")
}

// this works only when dictionary.txt file already present with all contents
func init() {
	file, err := os.Open("dictionary.txt")
	if err != nil {
		log.Println("Unable to load dictionary headers.")
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.SplitN(line, "=>", 2)
		if len(parts) != 2 {
			break
		} else {
			offsetInt, convErr := strconv.ParseInt(parts[1], 10, 64)
			if convErr != nil {
				log.Println("Unable to conver offset to integer : ", convErr.Error())
			}

			wordOffsets[parts[0]] = offsetInt
		}
	}
	log.Println("Headers loaded in memory successfully.")
}

func lookupWordFast(word string) (string, error) {
	file, err := os.Open("dictionary.txt")
	if err != nil {
		log.Println("Unable to load dictionary headers.")
	}
	defer file.Close()
	offset := 0

	value, exists := wordOffsets[word]
	if exists {
		offset = int(value)
	} else {
		return "", fmt.Errorf("word not found")
	}

	file.Seek(int64(offset), os.SEEK_SET)
	scanner := bufio.NewScanner(file)
	if scanner.Scan() {
		line := scanner.Text()
		parts := strings.SplitN(line, ":", 2)
		return parts[1], nil
	}

	return "", fmt.Errorf("error reading definition")
}

func main() {
	log.Println("Word dictionary using single file ")

	err := createDictionary()

	if err != nil {
		fmt.Println("Error creating dictionary:", err)
		return
	}

	startTime := time.Now()
	definition, err := lookupWord("waterfall")
	if err != nil {
		fmt.Printf("Error looking up '%s': %v\n", "waterfall", err)
	} else {
		fmt.Printf("Time taken : %f \nDefinition of '%s': %s\n", time.Since(startTime).Seconds(), "waterfall", definition)
	}

}
