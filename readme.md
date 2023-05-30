# gptWhatsapp


# Entendendo alguns conceitos bases

## Tokens 

No contexto do ChatGPT, um token é uma palavra ou conjunto de palavras que representa uma pergunta ou uma declaração feita pelo usuário.

Quando um usuário envia uma mensagem para o ChatGPT, o modelo divide essa mensagem em tokens. Cada token representa uma palavra ou conjunto de palavras que o modelo tentará interpretar e responder de acordo. Os tokens podem incluir palavras-chave, frases, símbolos e outras informações relevantes para a compreensão do significado da mensagem.

Por exemplo, se o usuário enviar a mensagem "Qual é a capital do Brasil?", o ChatGPT dividirá essa mensagem em três tokens: "Qual", "é a capital", "do Brasil". Cada um desses tokens representa uma parte da pergunta e é usado pelo modelo para gerar uma resposta apropriada.

## Models

Os modelos do ChatGPT são divididos em diferentes tamanhos e capacidades, cada um com um número diferente de parâmetros e habilidades para lidar com diferentes tarefas de linguagem natural. Atualmente, os principais modelos do ChatGPT são:

GPT-3: é o modelo mais avançado do ChatGPT, com 175 bilhões de parâmetros. Ele é capaz de realizar uma ampla gama de tarefas de linguagem natural, como tradução de idiomas, geração de texto, respostas a perguntas, entre outras. O GPT-3 é treinado em um conjunto de dados muito grande e é capaz de gerar textos extremamente coerentes e convincentes.

GPT-2: é o modelo anterior ao GPT-3, com 1,5 bilhão de parâmetros. Ele ainda é capaz de realizar várias tarefas de linguagem natural, como geração de texto, respostas a perguntas, entre outras. O GPT-2 é treinado em um conjunto de dados menor que o GPT-3, mas ainda é muito poderoso e eficaz.


## Como Funciona a API do chatGPT?

O sistema da API do ChatGPT é composto por várias camadas que trabalham juntas para processar as solicitações do usuário e fornecer uma resposta em linguagem natural. Aqui está uma explicação detalhada de como o sistema da API do ChatGPT funciona:

1. Recebimento da solicitação: O processo começa quando o usuário envia uma solicitação para o modelo do ChatGPT através da API. A API recebe a solicitação HTTP POST contendo a chave de API e o texto de entrada.

2. Autenticação: A API verifica a chave de API fornecida pelo usuário para garantir que a solicitação seja legítima e que o usuário tenha permissão para acessar o modelo do ChatGPT.

3. Pré-processamento: O texto de entrada é pré-processado para remover caracteres especiais, palavras sem sentido e outras informações desnecessárias que possam afetar o desempenho do modelo do ChatGPT.

4. Envio para o modelo: O texto de entrada pré-processado é enviado para o modelo do ChatGPT, que processa a entrada e gera uma resposta em linguagem natural.

5. Pós-processamento: A resposta gerada pelo modelo é pós-processada para remover informações desnecessárias, como caracteres especiais, espaços extras e outras informações que possam afetar a legibilidade da resposta.

6. Envio da resposta: A resposta processada é enviada de volta para o usuário como uma resposta HTTP.

Todo o processo ocorre em questão de segundos e é projetado para ser altamente escalável e capaz de lidar com um grande número de solicitações simultâneas. O sistema da API do ChatGPT é executado em servidores poderosos, que são capazes de processar grandes volumes de dados e fornecer respostas em tempo real.

# Definindo o projeto

## Objetivo

Desenvolver duas interfaces, a primeira sera parecida com o chatgpt, e a outra sera 
via whatsapp.

## Dinamica

Criar um microsservico de chat: Next.js (back-end) se comunicando com o Next.js (front-end)


### Como vai funcionar a interacao com sistema?

<img src="./images/flow-gptwhatsapp.png">

**Web**

O usuario esta no front no nextjs com componentes react -->

escrever um texto -->

Chama o back-end --> 

Back-end vai chamar o microsservico chat via gRPC (protocolo em microsservico) -->

O microsservico vai chamar o openIa -->

O chat microsservico vai falar com as api de openIa

**Whatsapp**

Twillio vai fornecedor o numero do whatsapp -->

Vai mandar mensagem no Whatsapp --> 

O Twillio vai receber a mensagem que foi enviada para o whatsapp -->

O twillio vai chamar atraves de request Http o microsservico do chat

**OBS: O microsservico de chat tera dois servidores tera dois servidores um gRPC e outro HTTP**

**OBS: A chamada via Twillio vai ser no formato Webhook**


## Regras de negocio

**Exlucuir mensagens para liberar espaco para outros tokens**: Temos a limitacao de tokens isso gera a necessidade de criar uma condicao para excluir alguns palavras/tokens nao tao importantes para o contexto de um dialogo de um usuario, para uma mensagem nova poder entrar.

**Contagem de tokens:**

**Armazenamento dos tokens**


## Contexto dos tokens

- Saber o modelo para identificar a quantidade de tokens posso armazenar
- A ideia aqui e fazer a contegem de tokens para saber quantos tokens posso acumular
- Quanto mais tokens o modelo suporta mais mensagens e melhor a resposta 


## Interface Web

### Tecnologias

* Nextjs
* React
* gRPC
* Node
* Docker
- Linguagem Go
- MySQL

## Interface Whatsapp

### Tecnologias

* Webhook
* Twillio
* Whatsapp
* Node
* Http




## Arquitetura

<img src="./images/CleanArchitecture.jpeg">

### Separando responsabilidade nas camadas

**Entidades:** No meio da aplicacao, as regras de negocios, acumlacao de tokens, contagem de tokens, escolher qual mensagem vai ser removida ou nao, chamado de Enterprise Busines Rules.

**Use cases:** A conversacao dos usuarios, a intencao de algo, acontece no usecases, ele ira utilizar as entidades e orquestrar o meio do caminho. Vai ser responsavel por chamar um repositorio para falar com um banco de dados para guardar os dados no banco, ou quando receber chamada no servidor http/gRPC, ele vai chamar o usecase para fazer a mensagem. `Application Business Rules`

**Controllers:** Basicamente serao responsaveis por receber as requisições HTTP, processá-las e enviar as respostas adequadas. Em Node.js, podemos implementar um Controller da seguinte forma:



 ## Pontos importantes

- Coracao da aplicacao deve ter suas regras de negocios consolidadas
- Coracao da aplicacao nao sabe que a API da OpenAI existe
- Armazenar todas conversacoes em banco de dados 
- Usuario podera informar seu user_id com referencia para ter acesso as conversas 
de um determinado usuario
- Servidor Web gRPC para realizar conversas
- Precisaremos gerar token no site da OpenAI para termos acesso API
A autenticacao de nosso microsservico tambem sera realizada via um token fixo em
um arquivo de configuracao



## Configurando o Go

### 1. Instalar o go

acessar o site: `https://go.dev/doc/install`

### 2. Instalar extensao e instalar tools
* Apos instalar a extensao instale o tools: ctrl/command + shift + p
* Digite Go: Install/Update

### 3. Iniciando modulo

terminal command: `go mod init github.com/gpt_chat/chat_service`

### 4. Estrutura da aplicacao

**internal:** tudo relacionado a aplicacao, nao tem relacao a ambiente externos nem pacotes

**domain:** E o dominio da aplicacao onde tera as regras de negocios

**entity:** E o centro da aplicacao onde ficam as entidades da aplicacao

* chat.go: E a entidade que representa o chat
* message.go: Representa a identidade das mensagens
* model.go: Representa o modelo do GPT, gpt 3.5, gpt turbo, etc.


