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

const monitoramentos = 3
const delayMonitoramento = 20

func main() {

	sites := []string{}

	exibeIntroducao()

	for {
		exibeMenu()
		comando := leComando()
		controleFluxo(comando, &sites)
	}

}

func exibeIntroducao() {
	nome := "Sidney" // Por padrão o valor inicial da variavel string é ""
	versao := 0.1    // Por padrão o valor inicial da variavel float é 0.0
	fmt.Println("Olá", nome)
	fmt.Println("Versão do programa:", versao)
}

func exibeMenu() {
	fmt.Println(`
		[1]-INICIAR MONITORAMENTO
		[2]-EXIBIR OS LOGS
		[3]-INSERIR SITE
		[4]-REMOVER SITE
		[0]-SAIR DO PROGRAMA
	`)

	fmt.Print("OPÇÃO: ")
}

func leComando() int {
	var comando int

	//fmt.Scanf("%d", &comando) | & indica que o valor sera apontado para o endereço da variavel comando
	fmt.Scan(&comando) //Scan dispensa o uso do tipificador %d, %f, %s, pois o tipo ja esta declarado na variavel

	//fmt.Println("O tipo da variavel nome é:", reflect.TypeOf(nome)) | reflect.TypeOf() indica o tipo da variavel

	//fmt.Println("O endereço da variavel comando é:", &comando)
	//fmt.Println("O comando escolhido foi", comando)

	return comando
}

func controleFluxo(comando int, sites *[]string) {

	// if comando == 1 {
	// 	fmt.Println("Monitorando...")
	// } else if comando == 2 {
	// 	fmt.Println("Exibindo Logs...")
	// } else if comando == 0 {
	// 	fmt.Println("Saindo do programa")
	// } else {
	// 	fmt.Println("Não conheço este comando")
	// }

	switch comando {
	case 1:
		testaSite(*sites)
	case 2:
		exibirLogs()
	case 3:
		inserirSite()
	case 4:
		removerSite()
	case 0:
		fmt.Println("Saindo...")
		os.Exit(0) //Sair do programa indicando ao sistema o sucesso na execusão
	default:
		fmt.Println("Não conheço este comando")
		os.Exit(-1) //Indica ao sistema que ocorreu algum erro na execusão do codigo
	}

}

func inserirSite() {

	fmt.Println("\n------ ADICIONANDO SITES ------")

	fmt.Println("Digite o endereço do site")
	var site string
	fmt.Scan(&site)

	site = "http://" + site

	gravarSite(site)
}

func testaSite(sites []string) {

	fmt.Println("\n------ MONITORANDO ------")

	sites = leSitesDoArquivo()

	for i := 0; i < monitoramentos; i++ {

		for _, site := range sites {
			resp, err := http.Get(site)

			if err != nil {
				fmt.Println("Ocorreu um erro: ", err.Error())
			} else {

				if resp != nil {
					if resp.StatusCode == 200 {
						fmt.Println("Site:", site, "Status:", resp.StatusCode)
						gravarLog(site, true)
						time.Sleep(delayMonitoramento * time.Second)
					} else {
						fmt.Println("Verificar acesso ao site:", site, "Status:", resp.StatusCode)
						gravarLog(site, false)
					}
				} else {
					fmt.Println("Site invalido", site)
				}

			}
		}

	}
}

func exibirLogs() {

	fmt.Println("\n------ EXIBINDO LOGS ------")

	arquivo, err := ioutil.ReadFile("file/log.txt")

	if err != nil {
		fmt.Println("Ocorreu um erro", err.Error())
	}

	fmt.Println(string(arquivo))

}

func removerSite() {

	fmt.Println("\n------ EXIBINDO LOGS ------")

	sites := leSitesDoArquivo()

	fmt.Println("Digite o site de deseja remover: ")
	var remover string
	fmt.Scan(&remover)

	for i, site := range sites {
		if strings.Contains(site, remover) {
			fmt.Println("Site encontrado, vamos remover [s/n]")
			fmt.Println(site)
			var confirma string
			fmt.Scan(&confirma)
			if confirma == "s" || confirma == "S" {
				sites = append(sites[:i], sites[i+1:]...)
				fmt.Println("Site removido")
			}
		}
	}
}

func leSitesDoArquivo() []string {

	var sites = []string{}
	arquivo, err := os.Open("file/sites.txt") // Somente abrir
	// arquivo, err := ioutil.ReadFile("../sites.txt") // Somente imprimir

	if err != nil {
		fmt.Println("Ocorreu um erro: ", err.Error())
	} else {

		reader := bufio.NewReader(arquivo)

		for {
			linha, err := reader.ReadString('\n')
			linha = strings.TrimSpace(linha) // remove espaçoes e quebras de linha
			sites = append(sites, linha)

			if err == io.EOF {
				break
			}
		}
	}

	arquivo.Close()

	return sites
}

func leLogDoArquivo() []string {

	var logs = []string{}
	arquivo, err := os.Open("file/log.txt")
	//ioutil.ReadFile("")

	if err != nil {
		fmt.Println("Ocorreu um erro: ", err.Error())
	} else {
		reader := bufio.NewReader(arquivo)

		for {
			linha, err := reader.ReadString('\n')
			linha = strings.TrimSpace(linha)
			logs = append(logs, linha)

			if err == io.EOF {
				break
			}
		}

	}
	arquivo.Close()
	return logs
}

func gravarSite(site string) {
	arquivo, err := os.OpenFile("file/sites.txt", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)

	if err != nil {
		fmt.Println("Ocorreu um erro", err.Error())
	} else {
		_, err := arquivo.WriteString("\n" + site)
		if err != nil {
			fmt.Println("Ocorreu um erro", err.Error())
		}
	}
	arquivo.Close()

}

func gravarLog(site string, status bool) {
	arquivo, err := os.OpenFile("file/log.txt", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)

	if err != nil {
		fmt.Println("Ocorreu um erro", err.Error())
	} else {
		_, err := arquivo.WriteString(time.Now().Format("02/01/2006 15:04:05") + " - " + site + " - status: " + strconv.FormatBool(status) + "\n")
		if err != nil {
			fmt.Println("Ocorreu um erro", err.Error())
		}
	}
	arquivo.Close()
}
