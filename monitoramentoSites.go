package main

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

const monitoramento = 3 
const delay = 5 

func main(){
	introducao()	
	
	for {
		exibeMenu()
		option := retornaOpcao()

		switch option{
		case 1: 
			monitorar()
		case 2: 
			fmt.Println("Exibindo logs")
			imprimirLogs()
		case 3: 
			fmt.Println("Saindo")
			os.Exit(0)
		default: fmt.Println("Inválido")
			os.Exit(-1)
		}
	}
}

func introducao(){
	version := 1.6
	fmt.Println("App de monitoramento de sites versão", version)
}

func exibeMenu(){
	fmt.Println("1 - Iniciar monitoramento")
	fmt.Println("2 - Exibir logs")
	fmt.Println("3 - Sair")
}

func retornaOpcao() int {
	var option int
	fmt.Scan(&option)
	fmt.Println("Opcao escolhida: ", option)

	return option
}

func monitorar(){
	sites := lerArquivosDeSite("sites.txt")		

	for i := 0; i < monitoramento; i++ {
		for _, site := range(sites){
			if siteOn(site)	{
				fmt.Println("Site", site, "online")		
				escreverLog(site, true)
			} else {
				fmt.Println("Site", site, "offline")		
				escreverLog(site, false)
			}
		}
		time.Sleep(delay * time.Second)
		fmt.Println(" ")
	}
}

func siteOn(site string) bool {
	response, err := http.Get(site)
	if err != nil {
		return false
	}
	return response.StatusCode == 200
}

func lerArquivosDeSite(path string) []string {
	var sites []string

	if _, err := os.Stat(path); err == nil {		
		arquivo, err_file := os.Open(path)
		
		if err_file == nil {			
			leitor := bufio.NewReader(arquivo)
			for {
				linha, err_linha := leitor.ReadString('\n')
				sites = append(sites, strings.TrimSpace(linha))	
				if err_linha == io.EOF {
					break
				}				
			}
		} else {
			print(err_file)
		}
		arquivo.Close()
	}
	return sites
}

func escreverLog(site string, status bool){
	if arquivo, err := os.OpenFile("log.txt", os.O_CREATE|os.O_RDWR|os.O_APPEND, 0666); err == nil{
		arquivo.WriteString( time.Now().Format("02/01/2006 15:04:05") + " - " + site + " - online: " + strconv.FormatBool(status) + "\n")
		arquivo.Close()
	} else {
		arquivo.Close()
		fmt.Println("Falha na criacao do arquivo de log")
	}	
}

func imprimirLogs(){
	if arquivo, err := ioutil.ReadFile("log.txt"); err == nil{
		fmt.Println(string(arquivo))
	} else {
		fmt.Println("Falha ao imprimiri os logs")	
	}
}