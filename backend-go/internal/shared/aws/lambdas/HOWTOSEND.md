
# Deploy de AWS Lambda em Go (Runtime Custom / bootstrap)

## 1. Compilar o binário

```bash
GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o bootstrap main.go
```

### Parâmetros

- `GOOS=linux`
  - Define o sistema operacional de destino como Linux.
  - AWS Lambda executa funções em ambiente Linux.

- `GOARCH=amd64`
  - Define a arquitetura do processador.
  - `amd64` = x86_64 (arquitetura padrão suportada pela Lambda).

- `CGO_ENABLED=0`
  - Desabilita CGO.
  - Gera um binário estático sem dependências nativas do sistema operacional.
  - Evita erros de compatibilidade no ambiente da Lambda.

- `go build`
  - Compila o código Go.

- `-o bootstrap`
  - Define o nome do arquivo gerado.
  - Em runtimes customizados da Lambda, o executável deve se chamar `bootstrap`.

- `main.go`
  - Arquivo de entrada da aplicação.

---

## 2. Dar permissão de execução

```bash
chmod +x bootstrap
```

### Parâmetros

- `chmod`
  - Altera permissões de arquivos.

- `+x`
  - Adiciona permissão de execução.

- `bootstrap`
  - Binário que será executado pela AWS Lambda.

---

## 3. Criar o ZIP

```bash
zip function.zip bootstrap
```

### Parâmetros

- `zip`
  - Utilitário para criar arquivos compactados.

- `function.zip`
  - Nome do arquivo ZIP gerado.

- `bootstrap`
  - Arquivo incluído no ZIP.

### Importante

O arquivo `bootstrap` deve estar na raiz do ZIP:

✅ Correto:

```text
function.zip
└── bootstrap
```

❌ Incorreto:

```text
function.zip
└── dist/
    └── bootstrap
```

---

## 4. Atualizar a função Lambda

```bash
aws lambda update-function-code \
  --function-name SUA_FUNCAO \
  --zip-file fileb://function.zip
```

### Parâmetros

- `aws`
  - AWS CLI.

- `lambda`
  - Serviço AWS Lambda.

- `update-function-code`
  - Atualiza o código da função existente.

- `--function-name`
  - Nome da função Lambda que receberá o deploy.

- `SUA_FUNCAO`
  - Nome real da função na AWS.

- `--zip-file`
  - Informa o pacote ZIP que será enviado.

- `fileb://function.zip`
  - Envia o arquivo em formato binário.
  - O prefixo `fileb://` é obrigatório para upload correto do ZIP.

---

## Fluxo completo

```bash
GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o bootstrap main.go

chmod +x bootstrap

zip function.zip bootstrap

aws lambda update-function-code \
  --function-name SUA_FUNCAO \
  --zip-file fileb://function.zip
```

