/*
programma che legge dal file binario la file table e la riordina per valore
*/
package main

import (
	"os"
	"fmt"
	"encoding/binary"
	"sort"
)

// QUESTE INFORMAZIONI SERVONO PER ESEGUIRE IL SORTING
// dichiaro un alias di tipo
type offsets []uint32
// metodo che mi permette di avere la dimensione di una struttura offsets 
func (i offsets) Len() int {
	return len(i)
}
// metodo che mi permette di scambiare due elementi all'interno della struttura
func (p offsets) Swap(i, j int) {
	p[i], p[j] = p[j], p[i]
} 
// metodo usato per comparare due elementi nella struttura (SE VOLESSI IN ORDINE DECRESCENTE MODIFICO QUA)
func (p offsets) Less(i, j int) bool {
	return p[i] < p[j]
}

const DEFAULT_VALUE = -1

func NextOffset(o offsets, pos int ) int {
	for i:= pos+1; i < o.Len(); i++{
		if o[i] != o[pos] {
			return i
		}
	}
	// non ho trovato alcun offset diverso: retituisco valore di controllo (default)
	return DEFAULT_VALUE
}

func main(){
	// leggo tutto il file e lo metto in una slice di byte
	data, _ := os.ReadFile(os.Args[1])
	// leggo il file count (sarebbero i primi 4 byte)
	file_count := binary.LittleEndian.Uint32(data[0:4])
	fmt.Printf("letto: %#04x, value: %d\n", file_count, file_count )
	file_table := data[4:file_count*4]
	// sorto la file_table
	// creo una slice di uint32 per semplificare sorting
	var file_table_intera []uint32 = make([]uint32, file_count);
	// aggiungo alla mia slice tutti gli offsets
	for i:=0; i < int(file_count); i++{
		// leggo un offset e lo metto all'interno della slice RISPETTANDO LA POSIZIONE
		file_table_intera[i] = binary.LittleEndian.Uint32(file_table[i*4:(i+1)*4])
	} 
	// sorto file_table_intera
	sort.Sort(offsets(file_table_intera));
	// ora ho la file table in ordine crescente all'interno di file_table_intera
	// for i  := 0; i < int(file_count); i++{
	// 	// leggo il valore dell'offset del file in pos i
	// 	valore_attuale := file_table_intera[i]
	// 	// stampo le informazioni riguardanti l'offset in posizione i
	// 	// fmt.Printf("%d\t:\t%#04x\t%d\n",i, valore_attuale,valore_attuale)
	// }
	// leggo file ed estraggo
	// eseguo l'operazione per ogni file
	for i := 0; i < int(file_count); i++ {
		// leggo l'offset di partenza
		file_start := int(file_table_intera[i]);
		// leggo l'offset di arrivo
		file_end := NextOffset(file_table_intera, i)

		// se è stata restituita default value vuol dire che devo leggere tutto il file, quindi metto file_end alla dimensione di data
		if file_end == DEFAULT_VALUE{
			// leggi tutto il file
			file_end = len(data)
		}else{
			file_end = int(file_table_intera[file_end])
		}
		// quanto byte devo leggere?
		lenght := file_end - file_start
		// fmt.Println("(DEBUG) lenght:", lenght, "file_start:", file_start, "file_end:", file_end)
		dati_da_file := data[file_start:file_start+lenght]
		// apro il file in scrittura
		f, err := os.Create("files_out" + string(os.PathSeparator) + fmt.Sprintf("%08d", i) + ".bin")
		if err != nil{
			fmt.Println("ERROR opening output file:", err)
		}else{
			f.Write(dati_da_file)
		}
		f.Close()
	}
}