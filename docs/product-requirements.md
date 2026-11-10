# Product Requirements — TCG

**Version:** 1.0
**Status:** Definido

Este documento define os requisitos funcionais e comportamentais do produto TCG.

Ele representa **o que o sistema deve oferecer aos seus usuários**, independentemente da tecnologia utilizada para implementar essas funcionalidades.

As regras específicas de funcionamento do duelo devem ser consultadas no documento:

```text
docs/game-rules.md
```

Este documento não deve substituir o Game Rules nem definir regras de combate.

---

# 1. Visão do Produto

O produto é um TCG competitivo 1 contra 1 com dois modos principais:

```text
Story Mode
Multiplayer Online
```

O jogador deve poder:

* criar uma conta;
* receber um Deck inicial;
* possuir uma coleção de cartas;
* criar e gerenciar múltiplos Decks;
* jogar o Story Mode;
* enfrentar adversários controlados por IA;
* jogar partidas online contra outros jogadores;
* receber cartas como recompensa;
* evoluir sua coleção e seus Decks ao longo da progressão.

O produto deve utilizar o mesmo núcleo de regras de duelo para Story Mode e Multiplayer.

---

# 2. Cadastro e Conta

## PR-ACCOUNT-001 — Cadastro

O sistema deve permitir que um novo usuário crie uma conta.

O cadastro deve possuir:

* nome;
* email;
* senha.

O sistema deve validar os dados necessários para criação da conta.

---

## PR-ACCOUNT-002 — Identidade do jogador

Após o cadastro, o usuário deve possuir uma identidade de jogador dentro do sistema.

O jogador deve possuir seus próprios:

* Decks;
* coleção de cartas;
* progresso;
* histórico de partidas;
* demais dados relacionados à sua conta.

---

## PR-ACCOUNT-003 — Login

O usuário deve poder autenticar-se utilizando sua conta.

Após autenticar-se, o sistema deve permitir acesso aos recursos correspondentes ao jogador.

---

## PR-ACCOUNT-004 — Persistência da conta

Os dados relevantes da conta devem permanecer disponíveis após o usuário sair do sistema e retornar posteriormente.

---

# 3. Deck Inicial

## PR-STARTER-001 — Recebimento do Deck inicial

Ao concluir o cadastro, o jogador deve receber automaticamente um Deck inicial.

O Deck inicial será escolhido aleatoriamente entre Decks iniciais previamente definidos pelo jogo.

O sistema não deve simplesmente gerar 30 cartas aleatórias.

---

## PR-STARTER-002 — Deck inicial jogável

Todos os Decks disponíveis como Deck inicial devem ser válidos de acordo com as regras de construção de Deck definidas pelo jogo.

O jogador deve conseguir iniciar sua experiência de jogo utilizando o Deck recebido.

---

# 4. Coleção de Cartas

## PR-COLLECTION-001 — Coleção

Cada jogador deve possuir uma coleção própria de cartas.

A coleção deve permitir identificar quais cartas o jogador possui e a quantidade de cópias de cada uma.

---

## PR-COLLECTION-002 — Quantidade de cópias

O sistema deve controlar a quantidade de cópias de cada carta pertencente ao jogador.

As regras de utilização dessas cópias em Decks devem respeitar as regras definidas em `game-rules.md`.

---

## PR-COLLECTION-003 — Visualização da coleção

O jogador deve poder visualizar sua coleção.

A interface deve permitir identificar, no mínimo:

* cartas disponíveis;
* quantidade de cópias;
* informações necessárias para construção de Decks.

---

## PR-COLLECTION-004 — Recebimento de cartas

O jogador pode receber novas cartas através de diferentes mecanismos do produto.

No MVP, uma das fontes de aquisição será:

```text
Recompensas do Story Mode
```

Outras formas de aquisição poderão ser adicionadas posteriormente.

---

# 5. Deck Building

## PR-DECK-001 — Múltiplos Decks

O jogador pode possuir múltiplos Decks.

Não existe uma regra de produto limitando o jogador a apenas um Deck.

---

## PR-DECK-002 — Criar Deck

O jogador deve poder criar novos Decks utilizando cartas disponíveis em sua coleção.

---

## PR-DECK-003 — Editar Deck

O jogador deve poder modificar um Deck existente.

A edição deve permitir adicionar e remover cartas conforme as cartas disponíveis na coleção.

---

## PR-DECK-004 — Validar Deck

Antes de um Deck ser utilizado em uma partida, o sistema deve verificar se ele é válido.

A validação deve respeitar as regras definidas no `game-rules.md`.

---

## PR-DECK-005 — Selecionar Deck

Quando possuir múltiplos Decks, o jogador deve poder escolher qual Deck utilizar em uma partida.

---

## PR-DECK-006 — Deck durante o Story Mode

Após receber novas cartas como recompensa, o jogador pode:

* manter o Deck atual;
* editar o Deck atual;
* criar um novo Deck;
* utilizar outro Deck já existente.

O jogador pode realizar essas alterações antes de continuar para o próximo duelo.

---

# 6. Story Mode

## PR-STORY-001 — Acesso ao Story Mode

O jogador deve poder acessar o Story Mode.

---

## PR-STORY-002 — Progressão

O Story Mode deve possuir a progressão definida pelo Game Rules.

Atualmente:

```text
10 níveis
4 duelos por nível
40 duelos
```

---

## PR-STORY-003 — Adversários

Cada duelo do Story Mode deve possuir um adversário controlado pela IA.

---

## PR-STORY-004 — Dificuldade progressiva

A dificuldade dos adversários deve aumentar conforme o jogador avança na campanha.

Os Decks dos adversários podem se tornar progressivamente mais fortes.

---

## PR-STORY-005 — Progressão por vitória

O jogador deve vencer todos os duelos necessários de um nível para avançar para o próximo nível.

A lógica exata de vitória e derrota pertence ao `game-rules.md`.

---

## PR-STORY-006 — Recompensas

Ao vencer um duelo do Story Mode, o jogador deve receber:

```text
3 cartas
```

Essas cartas devem ser adicionadas à coleção do jogador.

---

## PR-STORY-007 — Persistência da campanha

O progresso do jogador no Story Mode deve ser salvo.

O jogador não deve perder sua progressão ao sair do sistema.

---

## PR-STORY-008 — Retomar campanha

O jogador deve poder retornar posteriormente ao Story Mode e continuar sua progressão.

---

# 7. Sistema de Duelo

## PR-DUEL-001 — Criar duelo

O sistema deve ser capaz de iniciar uma partida válida entre dois participantes.

Os participantes podem ser:

```text
Jogador + IA
```

ou:

```text
Jogador + Jogador
```

---

## PR-DUEL-002 — Validação antes do duelo

Antes de iniciar uma partida, o sistema deve validar as condições necessárias para que os participantes possam jogar.

Entre elas está a validade do Deck utilizado.

---

## PR-DUEL-003 — Duel Engine compartilhado

Story Mode e Multiplayer devem utilizar o mesmo núcleo de regras de duelo.

O Duel Engine não deve possuir regras diferentes para:

```text
Story Mode
```

e:

```text
Multiplayer
```

A diferença deve estar apenas no controlador do adversário.

---

## PR-DUEL-004 — Estado do duelo

O sistema deve manter o estado necessário para representar uma partida em andamento.

Esse estado deve ser suficiente para determinar as ações permitidas aos participantes e o resultado da partida.

---

## PR-DUEL-005 — Resultado

Ao finalizar uma partida, o sistema deve registrar seu resultado.

O resultado deve identificar pelo menos:

* vencedor;
* perdedor;
* modo da partida;
* participantes;
* informações necessárias para histórico e progressão.

---

## PR-DUEL-006 — Histórico

O jogador deve poder consultar seu histórico de partidas.

O detalhamento exato das informações exibidas poderá evoluir posteriormente.

---

# 8. Multiplayer Online

## PR-MULTI-001 — Partidas online

O jogador deve poder disputar partidas contra outros jogadores.

---

## PR-MULTI-002 — Matchmaking

O sistema deve possuir um mecanismo de matchmaking para encontrar adversários.

O algoritmo de matchmaking poderá evoluir conforme o produto amadurecer.

---

## PR-MULTI-003 — Lobby

O produto deve possuir uma área de lobby relacionada às partidas online.

O lobby deve permitir ao jogador visualizar e acessar as funcionalidades disponíveis para iniciar ou encontrar partidas.

A definição detalhada do funcionamento do lobby será estabelecida durante o desenvolvimento.

---

## PR-MULTI-004 — Estado em tempo real

Durante uma partida online, os participantes devem receber as atualizações necessárias para acompanhar o estado atual do duelo.

---

## PR-MULTI-005 — Ações do jogador

As ações realizadas durante uma partida devem ser refletidas para o adversário.

Exemplos incluem:

* Summons;
* ataques;
* ativação de cartas;
* mudanças de posição;
* resolução de efeitos;
* mudanças de Life Points.

As regras dessas ações são definidas pelo `game-rules.md`.

---

## PR-MULTI-006 — Surrender

O jogador deve poder desistir de uma partida online.

O comportamento da desistência é definido pelo `game-rules.md`.

---

## PR-MULTI-007 — Abandono

O sistema deve detectar o abandono de uma partida quando aplicável.

O comportamento decorrente do abandono deve respeitar o `game-rules.md`.

---

## PR-MULTI-008 — Desconexão

O sistema deve detectar desconexões relevantes durante uma partida.

O comportamento da desconexão deve respeitar o `game-rules.md`.

---

## PR-MULTI-009 — Timeout

O sistema deve controlar o tempo de turno durante partidas.

O comportamento do timeout é definido pelo `game-rules.md`.

---

# 9. Interface do Produto

O produto deve possuir interfaces para as principais áreas funcionais.

## PR-UI-001 — Cadastro

Tela para criação de conta.

---

## PR-UI-002 — Login

Tela para autenticação.

---

## PR-UI-003 — Home / Lobby

Área principal após autenticação.

Deve fornecer acesso às principais funcionalidades do produto.

---

## PR-UI-004 — Coleção

Tela para visualização das cartas do jogador.

---

## PR-UI-005 — Deck Builder

Tela para:

* visualizar Decks;
* criar Decks;
* editar Decks;
* selecionar Deck.

---

## PR-UI-006 — Story Mode

Interface para:

* visualizar progressão;
* acessar duelos;
* visualizar recompensas;
* continuar a campanha.

---

## PR-UI-007 — Duel

Interface para acompanhar e interagir com uma partida.

---

## PR-UI-008 — Resultado

Após uma partida, o jogador deve visualizar o resultado do duelo.

Quando aplicável, também deve visualizar suas recompensas.

---

## PR-UI-009 — Multiplayer

Interface para acesso às funcionalidades de matchmaking e lobby.

---

# 10. Segurança e Integridade do Jogo

## PR-SECURITY-001 — Autoridade do servidor

O cliente não deve ser considerado a fonte de verdade para informações críticas do jogo.

O sistema deve validar no servidor as ações relevantes.

---

## PR-SECURITY-002 — Integridade da coleção

O jogador não deve conseguir adicionar ou remover arbitrariamente cartas de sua coleção através do cliente.

Alterações devem ocorrer através de operações válidas do sistema.

---

## PR-SECURITY-003 — Integridade dos Decks

O jogador não deve conseguir utilizar um Deck inválido em uma partida.

---

## PR-SECURITY-004 — Integridade do duelo

O cliente não deve conseguir determinar arbitrariamente:

* vencedor;
* perdedor;
* dano;
* Life Points;
* cartas recebidas;
* estado do campo;
* resultado de efeitos.

Essas informações devem ser determinadas pelo sistema responsável pelo duelo.

---

## PR-SECURITY-005 — Isolamento de jogadores

Um jogador não deve conseguir executar ações em nome de outro jogador.

---

# 11. IA dos Adversários

## PR-AI-001 — Adversários controlados por IA

O Story Mode deve possuir adversários controlados por IA.

---

## PR-AI-002 — Ações válidas

A IA deve realizar somente ações permitidas pelas regras do jogo.

---

## PR-AI-003 — Mesmo conjunto de regras

A IA deve jogar utilizando o mesmo Duel Engine utilizado pelos jogadores humanos.

A IA não deve possuir regras especiais de funcionamento do duelo.

---

# 12. Persistência do Produto

Os dados necessários para a continuidade da experiência do jogador devem ser persistidos.

Entre eles:

* conta;
* coleção;
* Decks;
* progresso do Story Mode;
* histórico de partidas;
* demais dados necessários para continuidade da experiência.

A tecnologia utilizada para persistência será definida no documento de Architecture.

---

# 13. Funcionalidades Fora do MVP

As funcionalidades abaixo não fazem parte da primeira versão do produto.

## 13.1 Card Shop

O Card Shop será desenvolvido em uma segunda etapa.

Ele poderá permitir:

* visualizar cartas disponíveis;
* comprar cartas;
* adicionar cartas à coleção.

---

## 13.2 Stars

O sistema de Stars e apostas ficará fora do MVP.

Funcionalidades planejadas para uma etapa futura incluem:

* carteira de Stars;
* saldo;
* apostas;
* aceite/recusa de apostas;
* distribuição de recompensas;
* histórico de transações.

---

## 13.3 Compra de Stars com dinheiro real

A compra de Stars utilizando dinheiro real ficará fora do MVP.

Qualquer implementação futura desse recurso deverá ser tratada como uma área própria do produto.

---

## 13.4 Ranking

O sistema de ranking competitivo ficará fora do MVP.

Futuramente poderá incluir:

* ranking;
* classificação dos jogadores;
* temporadas;
* divisões;
* recompensas competitivas.

---

## 13.5 Outras funcionalidades futuras

Também podem ser adicionados posteriormente:

* torneios;
* sistema social;
* chat;
* achievements;
* espectador;
* replay de partidas;
* outras funcionalidades competitivas.

Esses itens não fazem parte dos requisitos atuais.

---

# 14. Escopo do MVP

O MVP deve possuir **dois modos de jogo funcionais**:

```text
Story Mode
+
Multiplayer Online
```

Ambos devem utilizar o mesmo Duel Engine e as mesmas Game Rules.

## 14.1 Story Mode

O jogador deve poder:

```text
Cadastro
   ↓
Receber Deck inicial
   ↓
Visualizar coleção
   ↓
Criar/editar Decks
   ↓
Iniciar Story Mode
   ↓
Jogar contra IA
   ↓
Receber recompensas
   ↓
Modificar/criar Deck
   ↓
Continuar progressão
```

## 14.2 Multiplayer

O jogador deve poder:

```text
Login
   ↓
Lobby
   ↓
Matchmaking
   ↓
Encontrar adversário
   ↓
Selecionar Deck
   ↓
Iniciar duelo
   ↓
Jogar contra outro jogador
   ↓
Finalizar partida
   ↓
Registrar resultado
```

O Multiplayer faz parte do **MVP**.

O MVP não precisa possuir, inicialmente, sistemas competitivos avançados como ranking, temporadas ou divisões.

## 14.3 Ordem de implementação

Embora Story Mode e Multiplayer façam parte do MVP, a implementação pode ocorrer incrementalmente.

A ordem recomendada é:

```text
Duel Engine
    ↓
Story Mode
    ↓
Multiplayer
    ↓
Refinamentos do Multiplayer
```

Essa ordem é uma decisão de desenvolvimento e não altera o escopo do produto.

O objetivo é validar o Duel Engine primeiro em um ambiente controlado com IA antes de utilizá-lo em partidas entre jogadores reais.

---

# 15. Funcionalidades Fora do MVP

As seguintes funcionalidades ficam fora do MVP:

```text
Card Shop
Stars
Apostas
Compra de Stars com dinheiro real
Ranking competitivo
Temporadas
Divisões competitivas
Torneios
Sistema social
Chat
Achievements
Espectador
Replay
```

Essas funcionalidades podem ser adicionadas em versões posteriores.

---

# 16. Prioridade de implementação do MVP

O MVP pode ser implementado nas seguintes etapas:

### Etapa 1 — Núcleo do jogo

```text
Game Rules
    ↓
Domain Model
    ↓
Duel Engine
    ↓
Testes do domínio
```

### Etapa 2 — Story Mode

```text
Duel Engine
    +
AI Opponent
    +
Progression
    +
Rewards
```

### Etapa 3 — Conta e coleção

```text
Cadastro
    +
Login
    +
Collection
    +
Deck Building
```

### Etapa 4 — Multiplayer

```text
Lobby
    +
Matchmaking
    +
Human vs Human
    +
Sincronização da partida
```

### Etapa 5 — Refinamento

```text
Observabilidade
Segurança
Performance
UX
```

A divisão acima representa a **ordem de implementação**, não diferentes versões do produto.

O MVP somente estará completo quando os dois modos — **Story Mode e Multiplayer Online** — estiverem funcionais.
