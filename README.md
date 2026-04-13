# chainctl

CLI in Go per interagire con nodi EVM, Quorum e Besu.

`chainctl` permette di:

- interrogare un nodo RPC
- leggere balance, transazioni e receipt
- inviare ETH
- chiamare o inviare transazioni a smart contract via ABI
- gestire operazioni sui validator tramite RPC

## Requisiti

- Go `1.24+` per build ed esecuzione via sorgente
- accesso a un endpoint JSON-RPC EVM compatibile

## Installazione

Compilazione locale:

```bash
go build -o chainctl .
```

Esecuzione diretta senza build esplicita:

```bash
go run . --help
```

## Quick Start

1. Crea un file `config.yaml` partendo da [config.example.yaml](/home/alessandro/personalProject/chainctl/config.example.yaml:1)
2. Se devi firmare transazioni, crea anche un `.env` partendo da [.env.example](/home/alessandro/personalProject/chainctl/.env.example:1)
3. Verifica la connessione al nodo:

```bash
./chainctl status
```

4. Scopri i parametri disponibili:

```bash
./chainctl --help
./chainctl send-eth --help
./chainctl contract call --help
```

## Come Scoprire Gli Input Da CLI

Per capire cosa puoi passare da riga di comando usa sempre `--help`.

Help globale:

```bash
./chainctl --help
```

Mostra:

- elenco dei comandi
- flag globali validi per tutti i comandi
- descrizione generale del tool

Help di un comando specifico:

```bash
./chainctl <comando> --help
```

Esempi:

```bash
./chainctl balance --help
./chainctl send-eth --help
./chainctl contract --help
./chainctl contract call --help
./chainctl contract send --help
./chainctl validator --help
./chainctl validator propose-add --help
```

Regola pratica:

- usa `./chainctl --help` per vedere i parametri globali
- usa `./chainctl <comando> --help` per vedere gli input obbligatori e opzionali del comando

## Configurazione

`chainctl` risolve i parametri con questo ordine di precedenza:

1. flag da linea di comando
2. variabili ambiente esportate o lette da `.env`
3. `config.yaml`
4. default interni

Questo significa che puoi:

- usare solo `config.yaml`
- usare solo variabili d'ambiente o `.env`
- usare solo flag CLI
- combinare le fonti e sovrascrivere da CLI quello che ti serve

### Ricerca Del File `.env`

Il file `.env` viene cercato in:

- directory corrente
- directory del file passato con `--config`
- `$HOME/.chainctl/.env`

### Esempio `config.yaml`

Riferimento: [config.example.yaml](/home/alessandro/personalProject/chainctl/config.example.yaml:1)

```yaml
rpc_url: "http://127.0.0.1:8545"
chain_id: 1337
from_address: "0x0000000000000000000000000000000000000000"
private_key: ""
private_key_env: "CHAINCTL_PRIVATE_KEY"
timeout_seconds: 10
poll_interval_seconds: 3
validator:
  list_method: "qbft_getValidatorsByBlockNumber"
  propose_method: "clique_propose"
```

Significato dei campi:

- `rpc_url`: endpoint JSON-RPC del nodo
- `chain_id`: chain ID usato per firmare le transazioni
- `from_address`: address atteso della chiave privata
- `private_key`: chiave privata diretta in esadecimale
- `private_key_env`: nome della variabile ambiente da cui leggere la chiave
- `timeout_seconds`: timeout per le chiamate RPC
- `poll_interval_seconds`: intervallo di polling per operazioni ripetute come `wait-tx`
- `validator.list_method`: metodo RPC usato per leggere i validator
- `validator.propose_method`: metodo RPC usato per proporre add/remove validator

### Esempio `.env`

Riferimento: [.env.example](/home/alessandro/personalProject/chainctl/.env.example:1)

```dotenv
CHAINCTL_RPC_URL=http://127.0.0.1:8545
CHAINCTL_CHAIN_ID=1337
CHAINCTL_FROM_ADDRESS=0x0000000000000000000000000000000000000000
CHAINCTL_PRIVATE_KEY=your_private_key_here
CHAINCTL_TIMEOUT_SECONDS=10
CHAINCTL_POLL_INTERVAL_SECONDS=3
CHAINCTL_VALIDATOR_LIST_METHOD=qbft_getValidatorsByBlockNumber
CHAINCTL_VALIDATOR_PROPOSE_METHOD=clique_propose
```

Nota sulla chiave privata:

- puoi usare `CHAINCTL_PRIVATE_KEY` in `.env`
- oppure puoi mettere la chiave in un'altra variabile e indicarne il nome con `private_key_env` o `--private-key-env`
- oppure puoi passare la chiave direttamente con `--private-key`

## Flag Globali

Questi flag sono disponibili su tutti i comandi:

```text
--config
--rpc-url
--chain-id
--from-address
--private-key
--private-key-env
--timeout-seconds
--poll-interval-seconds
--validator-list-method
--validator-propose-method
```

In pratica:

- `--config` sceglie un file YAML specifico
- `--rpc-url` sovrascrive l'endpoint RPC
- `--chain-id` sovrascrive il chain ID
- `--from-address` imposta l'address mittente atteso
- `--private-key` passa la chiave direttamente da CLI
- `--private-key-env` indica il nome della variabile ambiente contenente la chiave
- `--timeout-seconds` cambia il timeout RPC
- `--poll-interval-seconds` cambia il polling delle operazioni periodiche
- `--validator-list-method` e `--validator-propose-method` permettono di adattarsi a RPC custom

## Comandi Disponibili

### `status`

Mostra lo stato base del nodo:

- `rpc_url`
- `chain_id`
- `block_number`

Uso:

```bash
./chainctl status
```

### `balance [address]`

Legge il balance di un address.

Input:

- argomento obbligatorio: `address`

Uso:

```bash
./chainctl balance 0x1111111111111111111111111111111111111111
```

### `tx [hash]`

Mostra i dati base di una transazione.

Input:

- argomento obbligatorio: `hash`

Uso:

```bash
./chainctl tx 0xaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa
```

### `receipt [hash]`

Mostra la receipt di una transazione.

Input:

- argomento obbligatorio: `hash`

Uso:

```bash
./chainctl receipt 0xaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa
```

### `wait-tx [hash]`

Attende finché una transazione non viene minata.

Input:

- argomento obbligatorio: `hash`
- flag opzionale: `--timeout`

Uso:

```bash
./chainctl wait-tx 0xaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa --timeout 180
```

### `send-eth`

Invia una transazione ETH semplice.

Input:

- flag obbligatorio: `--to`
- flag obbligatorio: `--value`
- flag opzionale: `--gas-limit`
- richiede anche configurazione RPC e una chiave privata valida

Uso con `.env`:

```bash
./chainctl send-eth \
  --rpc-url http://127.0.0.1:8545 \
  --from-address 0x0000000000000000000000000000000000000000 \
  --to 0x1111111111111111111111111111111111111111 \
  --value 0.1
```

Uso tutto da CLI:

```bash
./chainctl send-eth \
  --rpc-url http://127.0.0.1:8545 \
  --chain-id 1337 \
  --from-address 0x0000000000000000000000000000000000000000 \
  --private-key 0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef \
  --to 0x1111111111111111111111111111111111111111 \
  --value 0.1
```

### `contract call`

Chiama una funzione read-only di un contratto usando un ABI JSON.

Input:

- flag obbligatorio: `--abi`
- flag obbligatorio: `--address`
- flag obbligatorio: `--method`
- flag opzionale: `--args`

Uso:

```bash
./chainctl contract call \
  --abi ./MyContract.abi.json \
  --address 0x2222222222222222222222222222222222222222 \
  --method balanceOf \
  --args 0x1111111111111111111111111111111111111111
```

Nota: `--args` usa il formato Cobra `--args a,b,c`.

### `contract send`

Invia una transazione a una funzione write di un contratto.

Input:

- flag obbligatorio: `--abi`
- flag obbligatorio: `--address`
- flag obbligatorio: `--method`
- flag opzionale: `--args`
- flag opzionale: `--value`
- flag opzionale: `--gas-limit`
- richiede anche configurazione RPC e una chiave privata valida

Uso:

```bash
./chainctl contract send \
  --abi ./MyContract.abi.json \
  --address 0x2222222222222222222222222222222222222222 \
  --method transfer \
  --args 0x3333333333333333333333333333333333333333,1000 \
  --value 0
```

### `validator list`

Lista i validator correnti usando il metodo RPC configurato.

Uso:

```bash
./chainctl validator list
```

### `validator propose-add [address]`

Propone l'aggiunta di un validator.

Input:

- argomento obbligatorio: `address`

Uso:

```bash
./chainctl validator propose-add 0x4444444444444444444444444444444444444444
```

### `validator propose-remove [address]`

Propone la rimozione di un validator.

Input:

- argomento obbligatorio: `address`

Uso:

```bash
./chainctl validator propose-remove 0x4444444444444444444444444444444444444444
```

## Esempi Di Override

Usare un file di config ma sovrascrivere l'endpoint RPC da riga di comando:

```bash
./chainctl status --config ./config.yaml --rpc-url http://127.0.0.1:20000
```

Usare solo variabili d'ambiente:

```bash
export CHAINCTL_RPC_URL=http://127.0.0.1:8545
export CHAINCTL_CHAIN_ID=1337
./chainctl status
```

Usare un nome custom per la variabile contenente la private key:

```bash
export MY_NODE_PK=0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef
./chainctl send-eth \
  --private-key-env MY_NODE_PK \
  --rpc-url http://127.0.0.1:8545 \
  --from-address 0x0000000000000000000000000000000000000000 \
  --to 0x1111111111111111111111111111111111111111 \
  --value 0.1
```

## Note Operative

- l'output dei comandi e in formato JSON
- se `from_address` e impostato, la chiave privata deve corrispondere a quell'address
- `send-eth` e `contract send` stimano automaticamente il gas se `--gas-limit` non viene passato
- `wait-tx` usa `poll_interval_seconds` della configurazione per il polling

## Troubleshooting

Errore su chiave privata vuota:

- verifica `CHAINCTL_PRIVATE_KEY`
- oppure passa `--private-key`
- oppure indica il nome corretto con `--private-key-env`

Errore di mismatch su `from_address`:

- la chiave privata caricata non corrisponde all'address configurato
- correggi `from_address` oppure usa la chiave giusta

Errore RPC:

- verifica `rpc_url`
- controlla che il nodo sia raggiungibile
- se serve, aumenta `--timeout-seconds`

## File Utili Nel Repository

- [README.md](/home/alessandro/personalProject/chainctl/README.md:1)
- [config.example.yaml](/home/alessandro/personalProject/chainctl/config.example.yaml:1)
- [.env.example](/home/alessandro/personalProject/chainctl/.env.example:1)
- [cmd/root.go](/home/alessandro/personalProject/chainctl/cmd/root.go:27)
