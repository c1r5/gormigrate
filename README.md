# Gormigrate

`gormigrate` é uma biblioteca simples e eficiente para gerenciar migrações de banco de dados utilizando o [GORM](https://gorm.io/). Ela permite que você defina versões de migração, descrições e funções para subir (`Up`) ou reverter (`Down`) as alterações de forma controlada.

## Instalação

```bash
go get github.com/c1r5/gormigrate
```

## Funcionalidades

- **Controle por Versão:** Cada migração possui um número de versão sequencial.
- **Transações:** Migrações são executadas dentro de transações (se suportado pelo banco).
- **Helpers Úteis:** Funções para verificar existência de tabelas, colunas, índices e constraints.
- **Registro Global:** Facilidade para registrar migrações em diferentes partes do projeto.

## Exemplo de Uso

### 1. Definindo Migrações

Você pode registrar suas migrações em qualquer lugar do seu código, geralmente em um arquivo dedicado ou na inicialização do pacote.

```go
package migrations

import (
	"github.com/c1r5/gormigrate"
	"gorm.io/gorm"
)

func init() {
	gormigrate.Register(gormigrate.Migration{
		Version:     1,
		Description: "Criação da tabela de usuários",
		Up: func(db *gorm.DB) error {
			return db.Exec("CREATE TABLE users (id INT PRIMARY KEY, name VARCHAR(100))").Error
		},
		Down: func(db *gorm.DB) error {
			return db.Exec("DROP TABLE users").Error
		},
	})

	gormigrate.Register(gormigrate.Migration{
		Version:     2,
		Description: "Adiciona coluna de email",
		Up: func(db *gorm.DB) error {
			return db.Exec("ALTER TABLE users ADD COLUMN email VARCHAR(100)").Error
		},
		Down: func(db *gorm.DB) error {
			return db.Exec("ALTER TABLE users DROP COLUMN email").Error
		},
	})
}
```

### 2. Executando as Migrações

No seu `main.go`, você pode chamar a função `Run` para aplicar as migrações até uma versão específica.

```go
package main

import (
	"log"
	"github.com/c1r5/gormigrate"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	dsn := "user:pass@tcp(127.0.0.1:3306)/dbname?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}

	// Aplica todas as migrações até a versão 2
	targetVersion := 2
	if err := gormigrate.Run(db, targetVersion); err != nil {
		log.Fatalf("Erro ao executar migrações: %v", err)
	}
}
```

### 3. Revertendo Migrações

Se precisar voltar para uma versão anterior:

```go
// Reverte da versão 2 para a versão 1
if err := gormigrate.Rollback(db, 2, 1); err != nil {
    log.Fatalf("Erro ao reverter migrações: %v", err)
}
```

## Helpers

A biblioteca fornece helpers para facilitar verificações manuais dentro das migrações:

```go
if gormigrate.TableExists(db, "my_schema", "users") {
    // ...
}

if gormigrate.ColumnExists(db, "my_schema", "users", "email") {
    // ...
}
```

## Licença

MIT
