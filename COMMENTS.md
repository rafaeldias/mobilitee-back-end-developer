## Arquitetura
Optei por uma [arquitetura limpa](https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html), optando pela injeção de dependência em quase todas as funções, para facilitar o teste unitário, deixando o código modular. Evitando também deixar a camada de lógica de negócio associada a uma tecnologia específica (banco de dados, por exemplo), ficando essa camada agnóstica a tudo.

## Bibliotecas 
- [go-sqlmock](https://github.com/DATA-DOG/go-sqlmock)
- [GORM](https://github.com/jinzhu/gorm)
- [httprouter](https://github.com/julienschmidt/httprouter)

## O que melhoraria

- Código de Estado das respostas HTTP;
- Validação de Cabeçalhos HTTP enviados (`Content-Type`, por exemplo);
- Endpoints (No lugar de `/api/devices` usaria `/api/users/:id/devices`, me parece mais semântico agora);
- Refatoração de testes para evitar códigos repetidos;

## Requisitos Obrigatórios

Não houve falta da entrega dos requisitos obrigatórios. O requisito desejável de integração contínua não foi entregue por falta de tempo.
