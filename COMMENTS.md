## Arquitetura
Optei por uma [arquitetura limpa](https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html), optando pela inje��o de depend�ncia em quase todas as fun��es, para facilitar o teste unit�rio, deixando o c�digo modular. Evitando tamb�m deixar a camada de l�gica de neg�cio associada a uma tecnologia espec�fica (banco de dados, por exemplo), ficando essa camada agn�stica a tudo.

## Bibliotecas 
- [go-sqlmock](https://github.com/DATA-DOG/go-sqlmock)
- [GORM](https://github.com/jinzhu/gorm)
- [httprouter](https://github.com/julienschmidt/httprouter)

## O que melhoraria

- C�digo de Estado das respostas HTTP;
- Valida��o de Cabe�alhos HTTP enviados (`Content-Type`, por exemplo);
- Endpoints (No lugar de `/api/devices` usaria `/api/users/:id/devices`, me parece mais sem�ntico agora);
- Refatora��o de testes para evitar c�digos repetidos;

## Requisitos Obrigat�rios

N�o houve falta da entrega dos requisitos obrigat�rios. O requisito desej�vel de integra��o cont�nua n�o foi entregue por falta de tempo.
