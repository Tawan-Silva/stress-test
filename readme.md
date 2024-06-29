## Stress Test Go
### Como rodar o programa usando Docker
Docker
Este projeto também inclui um Dockerfile, então você pode construir e rodar o programa em um container Docker. 

Aqui estão os comandos para fazer isso:
```bash
docker build -t stress-test-go .
docker run stress-test-go --url=<URL do serviço a ser testado> --requests=<número total de requests> --concurrency=<número de chamadas simultâneas>
```

Substitua `<URL do serviço a ser testado>` pelo endereço do serviço que você deseja testar, `<número total de requests>` pelo número total de requests que você deseja fazer e `<número de chamadas simultâneas>` pelo número de chamadas simultâneas que você deseja fazer.