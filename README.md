# Projeto de Microsserviços com gRPC

Este projeto consiste na implementação de três microsserviços (Order, Payment e Shipping) em Go, utilizando a arquitetura hexagonal e gRPC para comunicação entre eles. Os microsserviços trabalham juntos para simular as etapas de uma compra em um sistema de comércio eletrônico.

## Microsserviços
- **Order (Pedido)**: Responsável por receber requisições de compras e orquestrar a cobrança e o envio. Ele verifica a existência dos produtos no estoque, salva o pedido no banco de dados, comunica-se com o serviço de pagamento e, se bem-sucedido, com o serviço de envio. Escuta na porta 3000.
- **Payment (Pagamento)**: Responsável por aprovar ou negar o pagamento. Rejeita pagamentos superiores a 1000. Escuta na porta 3001.
- **Shipping (Envio)**: Responsável por calcular o prazo de entrega. O prazo mínimo é 1 dia, sendo adicionado mais 1 dia a cada 5 unidades compradas. Escuta na porta 3002.

## Pré-requisitos
- Docker
- Docker Compose
- [grpcurl](https://github.com/fullstorydev/grpcurl) (Para realizar requisições gRPC aos serviços).

## Como Executar

Para iniciar todo o ambiente (Banco de dados MySQL, Order, Payment e Shipping), basta usar o Docker Compose a partir da raiz do projeto:

```bash
docker-compose up --build
```
Os serviços estarão rodando localmente nas seguintes portas:
- Order: `localhost:3000`
- Payment: `localhost:3001`
- Shipping: `localhost:3002`

O script `init.sql` cria automaticamente os bancos de dados (`order` e `payment`) necessários. Produtos pré-cadastrados no banco de dados para testes: `P1`, `P2`, `P3`.

## Como Testar

Abaixo estão os comandos para validar as regras de negócio implementadas nas várias partes do projeto, utilizando a ferramenta `grpcurl`.

**1. Pedido Válido (Comunicação com Payment e Shipping)**
Este comando realiza um pedido de 2 unidades do produto `P1`. O valor é inferior a 1000 e a quantidade não passa de 50. O pedido será aceito, pago e o prazo de entrega será calculado.
```bash
grpcurl -plaintext -d "{\"costumer_id\": 1, \"order_items\": [{\"product_code\": \"P1\", \"unit_price\": 10.5, \"quantity\": 2}]}" localhost:3000 Order/Create
```

**2. Retornar Erro (Preço superior ao limite de 1000)**
O serviço de `Payment` foi implementado para rejeitar compras que superam o limite de 1000. O erro retornado ao cliente provém da recusa do pagamento.
```bash
grpcurl -plaintext -d "{\"costumer_id\": 1, \"order_items\": [{\"product_code\": \"P1\", \"unit_price\": 1500, \"quantity\": 1}]}" localhost:3000 Order/Create
```

**3. Retornar Erro (Quantidade de itens acima de 50)**
O serviço de `Order` bloqueia a requisição antes mesmo de salvar ou acionar a cobrança se a soma das quantidades dos produtos exceder 50.
```bash
grpcurl -plaintext -d "{\"costumer_id\": 1, \"order_items\": [{\"product_code\": \"P1\", \"unit_price\": 10, \"quantity\": 60}]}" localhost:3000 Order/Create
```

**4. Retornar Erro (Produto Inexistente)**
O serviço de `Order` verifica no banco de dados se os itens solicitados existem (neste caso, `P1`, `P2` e `P3` existem). Se você passar um produto inexistente (ex: `P99`), a chamada retornará um erro:
```bash
grpcurl -plaintext -d "{\"costumer_id\": 1, \"order_items\": [{\"product_code\": \"P99\", \"unit_price\": 10, \"quantity\": 2}]}" localhost:3000 Order/Create
```
