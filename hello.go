package main // indica que é o pacote principal e que o programa deve começar por ele

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"net/http" // pacote para intereções http
	"os"       // pacote para intereção com o sistema operacional
	"strconv"  // pacote para conversão para string
	"strings"
	"time"
)

// declaração de constante
const monitoramentos = 5
const delayTime = 5

// declaração da função principal do programa
// não recebe nenhum parâmetro
// não retorna nada, quando retorna o programa termina
func main() {
	exibeIntroducao()

	// A instrução for sem parâmetros roda indefinidamente
	for {
		exibiMenu()
		comando := leComando()

		// ------------------------ Sobre 'if' 'else'
		//Sempre tem que ser uma expressão que retorne um boolean: 'true' ou 'false'
		/*
			if comando == 1 {
				fmt.Println("Monitorando...")

			} else if comando == 2 {
				fmt.Println("Exibindo Logs...")

			} else if comando == 0 {
				fmt.Println("Saindo do programa...")

			} else {
				fmt.Println("Comando não reconhecido")
			}

		*/

		// -------------------- Sobre switch
		// Não precisa do break
		switch comando {
		case 1:
			iniciarMonitoramento()
		case 2:
			imprimeLogs()
		case 0:
			fmt.Println("Saindo do programa...")
			os.Exit(0)
		default:
			fmt.Println("Comando não reconhecido")
			os.Exit(-1)
		}
	}
}

func exibeIntroducao() {
	nome := "Douglas"
	idade := 24
	versao := 1.1
	fmt.Println("Olá sr.", nome, "sua idade é", idade)
	fmt.Println("Este programa está na versão", versao)
}

func exibiMenu() {
	fmt.Println("1- Iniciar Monitoramento")
	fmt.Println("2- Exibir Logs")
	fmt.Println("0- Sair do Programa")
	fmt.Println("")
}

func leComando() int {
	var comandoLido int

	// fmt.Scanf é preciso explicitar o tipo de dado:
	// Primeiro parâmetro define o tipo de dado e o segundo aponta para a variável (com &)
	// fmt.Scanf("%d", &comando) //Scanf lê um valor e salva na variável selecionada

	// É possivel acessar o endereço em memória da variável com:
	//fmt.Println("O endereço da variável comando é:", &comando)

	//fmt.Scan não é preciso explicitar o valor pois ja foi explicitado na declaração
	fmt.Scan(&comandoLido)
	fmt.Println("")

	return comandoLido
}

func iniciarMonitoramento() {
	fmt.Println("Monitorando...")
	fmt.Println("")

	sites := leSitesDoArquivo()

	// podemos percorrer o slice com o "for" tradicional:
	// for i := 0; i < len(sites); i++  // utilizo o "i" para cada item do slice

	// Ou podemos utilizar um método melhor com o "range":
	// "range" retorna o índice atual e o item referente àquela posição
	for i := 0; i < monitoramentos; i++ {
		for _, site := range sites {
			testaSite(site)
		}

		// Forçamos um delay na execução do proximo laço da repetição
		time.Sleep(delayTime * time.Second)
	}

	fmt.Println(" ")
}

func testaSite(site string) {
	// Faz uma requisição get para uma url
	// http.Get()

	// A função .Get retornar mais de um valor, sendo assim podemos capturar os dois valores assim:
	// resp, err := http.Get(site)

	// Ou ignorar um dos valores com underline ao invés do nome da variável, assim:
	resp, err := http.Get(site)

	if err != nil {
		fmt.Println("Ocorreu um erro durante a requisição:", err)
	}

	if resp.StatusCode == 200 {
		fmt.Println("Site:", site, "foi carregado com sucesso!")
		registraLog(site, true)
	} else {
		fmt.Println("Site:", site, "esta com problemas. Status Code", resp.StatusCode)
		registraLog(site, false)
	}
}

func leSitesDoArquivo() []string {
	var sites []string

	// Com os.Open() irá retornar o apontamento para memória desse arquivo
	arquivo, err := os.Open("sites.txt")

	// Com ioutil.ReadFile() vai retornar um array de bytes
	// arquivo, err := ioutil.ReadFile("sites.txt")
	// depois pode ser lido assim:
	// string(arquivo)

	if err != nil {
		fmt.Println("Ocorreu um erro durante a abertura do arquivo:", err)
	}

	// Cria uma nova instancia de leitor do bufio
	leitor := bufio.NewReader(arquivo)

	for {
		// Define a leitura até encontrar o primeiro byte de \n
		// Byte específico é definido com aspas simples ''
		linha, error := leitor.ReadString('\n')
		linha = strings.TrimSpace(linha)

		sites = append(sites, linha)

		if error == io.EOF {
			break
		}
	}

	// É uma boa prática fechar o arquivo no final da função em que ele é aberto
	arquivo.Close()

	return sites
}

func registraLog(site string, status bool) {
	// os.OpenFile() da mais opções para tratar os arquivos, como o "os.O_RDWR" que permite a leitura e escrita
	// e o "os.O_CREATE" que permite a criação do arquivo caso não exista
	// A flag "os.O_APPEND" permite a adição de novas linhas, ao invés de escrever a partir do começo do arquivo
	arquivo, err := os.OpenFile("log.txt", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)

	if err != nil {
		fmt.Println("D")
	}

	// Escreve a string passada como parâmetro
	// Para mais infos sobre os formatos de tempo consultar a doc
	arquivo.WriteString(time.Now().Format("02/01/2006 15:04:05") + " - " + site + " - online: " + strconv.FormatBool(status) + "\n")

	arquivo.Close()
}

func imprimeLogs() {
	fmt.Println("Exibindo Logs...")
	fmt.Println("")
	// Não é necessário fechar o arquivo com ioutil.ReadFile()
	arquivo, err := ioutil.ReadFile("log.txt")

	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(string(arquivo))
}

// -------------------------- Sobre Arrays e Slices

// func trataArraysESlices() {
// ---------------- Array
// arrays em GO tem que ter um tamanho fixo, isso gera um problema
// var sites [4]string // declaração do array e seu tamanho
// sites[0] = "https://www.alura.com.br"
// sites[1] = "https://www.alura.com.br/dsadassda"
// sites[2] = "https://www.caelum.com.br"

// ---------------- Slice
// Não precisa declarar um tamanho, e utiliza o array por baixo dos panos
// nomes := []string{"Luiz", "Gustavo", "Pedrosa"} // [Luiz Gustavo Pedrosa]

// podemos utilizar a função "len()" para verificar o tamanho do array (não a capacidade)
// len(nomes) // 3

// podemos utilizar a função "cap()" para verificar capacidade do array
// cap(nomes) // 3

// No slice podemos adicionar novos itens com a função "append()"
// nomes = append(nomes, "Melo")

// Quando a capacidade de um slice é "estourada" a própria linguagem GO dobra a capacidade
// }

// -------------------------- Sobre funções -------------
/*
------ vai retornar dois valores, um int e uma string
	func retornaDoisValore() (string, int) {
		nome := "Gustavo"
		idade := 24

		return nome, idade
	}
*/

// --------------------------- Sobre Variáveis -----------------

/*
func variaveis() {
	// Variaveis não utilizadas geram erros
	// Quando a variável é declarada sem um valor fica "vazia": Para string seria string vazia, para int 0...

	var nome string = "Gustavo" // Declaração de variável sem inferência de tipo
	var idade = 32              // Declaração de variável com inferência de tipo
	var versao float32 = 1.1    // No caso do float sem inferência ele manda pro maior (float64)

	fmt.Println("Olá, sr.", nome) // Toda função que vem de um pacote externo tem a primeira letra maiuscula
	fmt.Println("Sua idade:", idade)
	fmt.Println("Este programa está na versão", versao)

	// Descobrindo o tipo das variáveis
	fmt.Println("O tipo da variável nome é:", reflect.TypeOf(nome))
	fmt.Println("O tipo da variável idade é:", reflect.TypeOf(idade))
	fmt.Println("O tipo da variável versao é:", reflect.TypeOf(versao))

	// Declaração de variável sem a palavra "var" é com o operador de declaração curta ":="
	nome1 := "Luiz"
	idade1 := 64
	versao1 := 2.5

	fmt.Println("Olá, sr.", nome1)
	fmt.Println("Sua idade:", idade1)
	fmt.Println("Este programa está na versão", versao1)

	fmt.Println("O tipo da variável nome1 é:", reflect.TypeOf(nome1))
	fmt.Println("O tipo da variável idade1 é:", reflect.TypeOf(idade1))
	fmt.Println("O tipo da variável versao1 é:", reflect.TypeOf(versao1))
}

*/
