### Consulta de CEP rápida e simples

Fastcep unifica consultas de `CEP` de diversos serviços em apenas 1 local, oque te permite economizar tempo
para realizar suas consultas.

### Como?
Fastcep utiliza concorrência diponibilizada em [`Golang`](http://golang.org/) para realizar
diversas requisições simultaneas em diversos serviços de `CEP` sendo capaz de selecionar apenas a resposta mais rapida dentro os serviços utilizados.

### Por que?
Existem diversos serviços de consulta de `CEP` espalhados e cada um possui um corpo de resposta diferente,
Fastcep consegue consultar estes serviços e retornar uma resposta uniforme independente de qual foi a fonte de consulta, assim você não precisa se preocupar com qual foi o serviço que respondeu e sim com os dados em sí, além de ter outros serviços de "backup" caso algum deles falhe, isso significa que você sempre terá uma resposta válida o que torna suas consultas "a prova de falhas"

### Executando

1. Faça o download do [código fonte](https://github.com/rafa-acioly/fastcep/archive/master.zip)
2. Extraia os arquivos
3. Compile o binário executando o comando em seu terminal `go build`
4. Rode o projeto com `./fastcep`

> Fastcep tentará utilizar a porta definida na variável de ambiente `PORT`, caso não encontre a porta `3000` será utilizada.

### Exemplo

Abaixo um comparativo de uma requisição única feita em diversos serviços, note o tempo de resposta de cada um.

```sh
# viacep - tempo de resposta: 754ms
curl https://viacep.com.br/ws/07400885/json/

# postmon - tempo de resposta: 118ms
curl http://api.postmon.com.br/v1/cep/07400885

# republica virtual - tempo de resposta: 66ms
curl http://republicavirtual.com.br/web_cep.php?cep=07400885&formato=json
```


```sh
# fastcep - utilizando como exemplo as requisições acima
# a resposta seria de 66ms + 5ms(tempo médio de processamento)
curl http://fastcep.appspot.com/v1/07400885
```

Fastcep é capaz de realizar as requisições simultaneas em todos os serviços e devolver
a resposta mais rapida, desta maneira você não precisa ficar refem de apenas um serviço e consequentemente economizara tempo nas suas requisições,
se algum dos serviços ficar indisponivel ele sera automaticamente descartado até que normalize e outros serviços continuaram respondendo normalmente assim você não corre
o risco de ficar sem os dados de `CEP` em nenhum momento. :tada:


#### Quer adicionar outro serviço para consulta?
- substitua o valor cep por `%s` no endpoint
- Adicione o nome e endpoint na [lista de endpoint](https://github.com/rafa-acioly/fastcep/blob/master/main.go#L21)


![](https://preview.ibb.co/fnkAxU/gopher.png)
