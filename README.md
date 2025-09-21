# Uni-Demo API

API de demonstração para consulta da agenda de professores universitarios. Fornecendo endpoints para consultar os horários de professores e as salas ocupadas por disciplinas.

Utilizamos Golang com o framework Chi para lidar com as requisições. O banco de dados utilizado é o PostgreSQL.

## Endpoints

* GET /professor-hours: Retorna os horários de professores.
* GET /room-schedules: Retorna as salas ocupadas por disciplinas.

## Como usar

Para usar a API, basta executar o comando `docker compose up` na pasta raiz do projeto. Isso irá inicializar o servidor na porta 8080 e o container do banco de dados localmente.

## Documentação
A documentação Swagger da API pode ser acessada no caminho `/swagger/index.html`
