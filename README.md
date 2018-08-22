### Consulta de CEP rápida e simples

Fastcep unifica consultas de `CEP` de diversos serviços em apenas 1 local, oque te permite economizar tempo
para realizar suas consultas.

### Como?
Fastcep utiliza concorrência diponibilizada em [`Golang`](http://golang.org/) para realizar
diversas requisições simultaneas em diversos serviços de `CEP` sendo capaz de selecionar apenas a resposta mais rapida dentro os serviços utilizados.

### Por que?
Existem diversos serviços de consulta de `CEP` espalhados e cada um possui um corpo de resposta diferente,
Fastcep consegue consultar estes serviços e retornar uma resposta uniforme independente de qual foi a fonte de consulta, assim você não precisa
se preocupar com qual foi o serviço que respondeu e sim com os dados em sí.

### Exemplo

Abaixo existe um comparativo de uma requisição única feita em diversos serviços, note o tempo de resposta de cada um.

```sh
# viacep - tempo de resposta: 754ms
curl https://viacep.com.br/ws/07400885/json/

# postmon - tempo de resposta: 118ms
curl http://api.postmon.com.br/v1/cep/07400885

# republica virtual - tempo de resposta: 66ms
curl http://republicavirtual.com.br/web_cep.php?cep=07400885&formato=json
```

Fastcep é capaz de realizar as requisições simultaneas em todos os serviços e devolver
a resposta mais rapida, desta maneira você não precisa ficar refem de apenas um serviço e consequentemente economizara tempo nas suas requisições. :tada:

#### Requisição com Fastcep: (link provisório)

```sh
# fastcep - tempo de resposta: em média 5ms a mais do menor tempo
curl https://fastcep-kzvqeaxdzi.now.sh/v1/07400885
```

#### Quer adicionar outro serviço para consulta?
- substitua o valor cep por `%s` no endpoint
- Adicione o nome e endpoint na [lista de endpoint](https://github.com/rafa-acioly/fastcep/blob/master/main.go#L21)
