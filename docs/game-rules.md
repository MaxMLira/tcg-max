# Game Rules — TCG

**Version:** 1.0
**Status:** Definido

Este documento define as regras de negócio e funcionamento do jogo de cartas.

Ele representa **o que o jogo deve fazer**, independentemente da tecnologia utilizada para implementá-lo.

> **Importante:** este documento é a fonte de verdade das regras do jogo.

A implementação técnica não deve criar ou alterar regras que não estejam definidas aqui.

> **Regras gerais x Story Mode MVP:** este documento descreve as regras gerais do jogo completo. Os blocos marcados como **Story Mode MVP** descrevem o comportamento do MVP e, onde diferem da regra geral, **prevalecem no MVP**. O escopo do MVP está em `docs/mvp-story-mode.md`.

---

# 1. Objetivo

O jogo é um **TCG competitivo 1 contra 1**.

Cada jogador utiliza um Deck de cartas para enfrentar o adversário em um duelo.

O objetivo principal é reduzir os **Life Points (LP)** do adversário para **0 ou menos**.

---

# 2. Vitória

Um jogador vence imediatamente quando os Life Points do adversário chegam a:

```text
LP <= 0
```

Exemplo:

```text
Jogador A: 500 LP

Jogador B causa 700 de dano

Jogador A: -200 LP

Resultado:
Jogador A perde o duelo.
```

O duelo termina imediatamente quando uma condição de vitória ou derrota é atingida.

Nenhuma ação posterior deve ser executada.

---

# 3. Life Points

Cada jogador começa o duelo com uma quantidade de Life Points determinada pela progressão do jogo.

## 3.1 Life Points iniciais

Nos primeiros níveis:

```text
4000 LP
```

A partir do **nível 5**, os jogadores passam a iniciar os duelos com:

```text
8000 LP
```

Progressão:

```text
Níveis 1–4
    ↓
4000 LP

Níveis 5–10
    ↓
8000 LP
```

---

# 4. Deck

Cada jogador possui um Deck com exatamente:

```text
30 cartas
```

O Deck pode conter qualquer combinação dos tipos de cartas disponíveis:

* Monstros;
* Magias;
* Armadilhas.

Não existem cartas obrigatórias para construção do Deck.

## 4.1 Quantidade máxima de cópias

O jogador pode possuir no máximo:

```text
2 cópias da mesma carta
```

Exemplo válido:

```text
Carta A × 2
Carta B × 2
Carta C × 1
...
Total = 30
```

Exemplo inválido:

```text
Carta A × 3
```

---

# 5. Mão

## 5.1 Mão inicial

**Regra geral (fora do MVP):** cada jogador começa o duelo com:

```text
3 cartas
```

O jogador que começa o duelo realiza sua primeira compra normalmente.

Assim:

```text
Mão inicial
3 cartas

↓ Draw

4 cartas
```

> **Story Mode MVP:** cada jogador começa o duelo com exatamente **4 cartas aleatórias do próprio Deck**. Não existe carta inicial adicional nem vantagem de compra para o primeiro jogador.

## 5.2 Limite da mão

A mão possui limite máximo de:

```text
4 cartas
```

Um jogador nunca possui mais de 4 cartas na mão.

> **Story Mode MVP:** a mão é reposta automaticamente até 4 cartas (seção 7).

---

# 6. Primeiro jogador

O primeiro jogador é definido aleatoriamente.

Não existe escolha manual de quem começa.

O jogador escolhido para começar:

* compra cartas normalmente;
* pode realizar suas ações;
* pode realizar Summons;
* pode colocar Magias e Armadilhas;
* pode ativar cartas conforme suas regras;
* **não pode atacar no primeiro turno**.

> **Story Mode MVP:** o primeiro jogador não recebe cartas adicionais. A única diferença é que ele não pode atacar no primeiro turno.

---

# 7. Compra

No início de cada turno ocorre a fase de compra.

O jogador compra cartas até completar sua mão:

```text
Mão < 4
    ↓
Compra
    ↓
Mão = 4
```

Exemplos:

```text
Mão com 3 cartas
↓
Compra 1
↓
4 cartas
```

```text
Mão com 1 carta
↓
Compra
↓
Compra
↓
Compra
↓
4 cartas
```

Se o Deck não possuir cartas suficientes para completar a compra necessária, o jogador perde imediatamente ao tentar comprar uma carta inexistente.

> **Story Mode MVP:** não existe ação manual de compra. Sempre que uma ação resolvida deixa o jogador com menos de 4 cartas na mão, a mão é reposta automaticamente até 4 cartas. Se o Deck não puder fornecer as cartas necessárias, o jogador perde imediatamente.

---

# 8. Deck Exaurido

Um jogador perde imediatamente quando precisa comprar uma carta e seu Deck não possui mais cartas.

Não existe empate causado por Deck vazio.

Exemplo:

```text
Deck: 0 cartas
Mão: 3 cartas

Início do turno
↓
Precisa comprar
↓
Não existem cartas
↓
Derrota imediata
```

---

# 9. Campo

Cada jogador possui duas zonas principais.

## 9.1 Zona de Monstros

Cada jogador pode possuir no máximo:

```text
4 Monstros
```

## 9.2 Zona de Magias/Armadilhas

Cada jogador pode possuir no máximo:

```text
4 cartas
```

Esses espaços podem conter Magias ou Armadilhas.

Portanto:

```text
Campo de cada jogador

4 Monstros
+
4 Magias/Armadilhas
```

---

# 10. Tipos de cartas

Existem três tipos principais de cartas:

```text
MONSTER
MAGIC
TRAP
```

## 10.1 Monster

Monstros são utilizados principalmente para:

* atacar;
* defender;
* realizar tributos;
* possuir efeitos próprios.

## 10.1.1 Subtipos de Monstro

Todo Monstro possui um **subtipo**. Magic e Trap não possuem subtipo de Monstro.

O subtipo é um conceito extensível: novos subtipos podem ser adicionados no futuro sem redesenhar o modelo de domínio.

> **Story Mode MVP:** os subtipos disponíveis são `Dragon`, `Warrior`, `Mage` e `Aquatic`. O subtipo é usado para determinar quais Monstros uma Magic Equipment aceita (seção 35.2).

## 10.2 Magic

Magias são cartas ativadas pelo jogador de acordo com suas condições.

Podem possuir diferentes efeitos, como:

* `ADD_TO_HAND`;
* `DESTROY`;
* `RECOVER_LP`;
* `BUFF_ATK`;
* `DEBUFF_ATK`;
* `FIELD_TYPE`;
* entre outros.

> **Story Mode MVP:** as únicas Magics são as quatro Magic Equipment (seção 35.2). Burn e Power Up não fazem parte do MVP.

## 10.3 Trap

Armadilhas são cartas associadas a eventos ou condições específicas.

Elas podem responder a ações do jogo, como:

* ataques;
* ativação de Magias;
* outros eventos definidos pela carta.

Traps podem responder a outras Magias ou Traps.

> **Story Mode MVP:** Trap continua sendo um tipo de carta válido no domínio, mas jogabilidade, ativação e efeitos de Trap estão fora do MVP.

---

# 11. Turno

Cada turno possui as seguintes etapas:

```text
Draw
 ↓
Discard
 ↓
Phase 1
 ↓
Phase 2
 ↓
End Turn
```

No primeiro turno, o jogador não pode atacar.

> **Story Mode MVP:** o turno possui as etapas `Phase 1 → Phase 2 → End Turn`. Não há etapa manual de Draw nem de Discard: a reposição da mão é automática (seção 7) e o descarte faz parte da ação de seleção da seção 13.3.

---

# 12. Draw

No início do turno:

```text
Se mão < 4
    ↓
Comprar até chegar a 4
```

Se a mão já possuir 4 cartas:

```text
Nenhuma compra é realizada.
```

> **Story Mode MVP:** não existe fase de Draw. A mão é reposta automaticamente até 4 cartas após cada ação resolvida que deixe a mão com menos de 4 cartas (seção 7). Se o Deck não puder completar a mão, o jogador perde imediatamente (seção 8).

---

# 13. Discard

> **Story Mode MVP:** não existe etapa de Discard nem ação de descarte isolada. Descarte e colocação formam uma única ação (seção 13.3). As seções 13.1 e 13.2 são a regra geral.

O Discard está relacionado à preparação de uma colocação de carta.

O jogador pode descartar cartas antes de realizar uma colocação.

O descarte **não é obrigatório em todo turno**.

Se o jogador já possui Monstros no campo e deseja apenas:

* atacar;
* alterar posições;
* ativar uma Magic;
* ou realizar outras ações permitidas;

não precisa realizar um descarte.

## 13.1 Quantidade de cartas descartadas

O jogador pode descartar:

```text
1 ou mais cartas
```

Inclusive pode descartar todas as cartas que possui na mão.

Exemplo:

```text
Mão:
A
B
C
D

Pode descartar:

A
B
C
D
```

Mesmo descartando quatro cartas, o jogador continua limitado a **uma única colocação de carta na Phase 1**.

## 13.2 Destino do descarte

Toda carta descartada vai para o:

```text
Graveyard
```

## 13.3 Seleção, descarte e colocação (Story Mode MVP)

No Story Mode MVP, descarte e colocação formam **uma única ação atômica**, realizada na Phase 1 e sujeita ao limite de uma colocação por Phase 1 (seção 14).

* O jogador seleciona de **1 a 4** cartas da mão. A ordem de seleção importa.
* A **última** carta selecionada é a carta colocada no campo.
* Todas as outras cartas selecionadas vão para o Graveyard.
* A última carta precisa satisfazer todas as regras de colocação do seu tipo.
* Se a última carta não puder ser colocada legalmente, a ação inteira é rejeitada e nenhuma carta selecionada é descartada.
* Quando a ação tem sucesso, a mão é reposta automaticamente até 4 cartas (seção 7). Se o Deck não puder completar a mão, o jogador perde imediatamente.
* Não existe ação de descarte isolada.

Exemplos:

```text
Seleção: A
A → campo (nenhuma carta descartada)
```

```text
Mão: A, B, C, D
Seleção: A, B, C, D

A → Graveyard
B → Graveyard
C → Graveyard
D → campo
```

Premissas técnicas (a regra mais simples consistente com as decisões):

* Para um Monstro que exige Tributos, os Tributos, a posição e a face fazem parte da mesma ação e são validados como parte da colocação. Se forem inválidos, a ação inteira é rejeitada.
* Cartas Trap não aparecem nos Decks do MVP e sua colocação é rejeitada.
* A mesma carta não pode aparecer duas vezes na seleção.

---

# 14. Phase 1

A Phase 1 é responsável principalmente pela preparação do campo.

Durante cada Phase 1, o jogador possui no máximo:

```text
1 colocação
```

A colocação pode ser:

* Summon de um Monstro;
* colocação de uma Magic;
* colocação de uma Trap.

## 14.1 Limite de colocação

Não é permitido:

```text
Summon Monster A
+
Summon Monster B
```

na mesma Phase 1.

Também não é permitido:

```text
Summon Monster
+
Set Magic
```

na mesma Phase 1.

Apenas uma colocação é permitida.

## 14.2 Alteração de posição

Alterar a posição de um Monstro que já está no campo **não conta como uma colocação**.

Portanto, o jogador pode:

```text
Alterar posição
+
Realizar um Summon
```

na mesma Phase 1, desde que todas as outras condições sejam atendidas.

---

# 15. Summon de Monstros

Todo Monstro colocado no campo é considerado um **Summon**.

Não existe uma ação de `Set Monster` separada de `Summon`.

O jogador pode realizar um Summon:

```text
Face-up Attack
Face-up Defense
Face-down Defense
```

A escolha depende das regras de tributo e posição.

---

# 16. Weight

Cada Monstro possui um valor chamado:

```text
Weight
```

O Weight determina quantos tributos são necessários para realizar seu Summon.

| Weight | Tributos |
| ------ | -------: |
| 1–5    |        0 |
| 6–8    |        1 |
| 9–12   |        2 |
| 13+    |        3 |

O Weight dos Monstros utilizados como tributo não precisa corresponder ao Weight do Monstro sendo invocado.

Exemplo:

```text
Monstro alvo:
Weight 10

Tributos necessários:
2
```

Podem ser utilizados:

```text
Weight 3
+
Weight 8
```

Total dos tributos:

```text
11
```

O Summon continua válido.

---

# 17. Tributos

Quando um Monstro exige tributos:

* os tributos devem ser Monstros que já estejam no campo;
* Monstros na mão não podem ser utilizados;
* Monstros face-down podem ser utilizados;
* os Monstros tributados vão para o Graveyard.

## 17.1 Preparação de tributos

Como existe apenas uma colocação por Phase 1, não é possível realizar uma sequência de Summons para criar tributos na mesma Phase 1.

Exemplo:

```text
Campo vazio

Mão:
Monster A — Weight 3
Monster B — Weight 8
Monster C — Weight 10
```

Não é permitido:

```text
Summon A
↓
Summon B
↓
Tributar A+B
↓
Summon C
```

na mesma Phase 1.

Os tributos precisam estar previamente no campo.

---

# 18. Face-up e Face-down

Um Monstro pode estar:

```text
Face-up
```

ou:

```text
Face-down
```

## 18.1 Face-up

O adversário pode visualizar as informações do Monstro.

## 18.2 Face-down

Quando um Monstro está face-down:

* todas as suas informações ficam ocultas para o adversário;
* o adversário não conhece ATK;
* o adversário não conhece DEF;
* o adversário não conhece Weight;
* o adversário não conhece seus efeitos;
* o adversário não conhece suas demais informações.

O próprio jogador continua podendo visualizar as informações da sua carta.

Um Monstro face-down obrigatoriamente está em:

```text
Defense
```

---

# 19. Summon Face-down

Um Monstro de Weight 1–5 pode ser Summoned sem tributos e colocado:

```text
Face-down Defense
```

Exemplo:

```text
Monster:
Weight: 2
ATK: 0
DEF: 2000

Summon
↓
Face-down Defense
```

Esse Monstro foi **Summoned normalmente**, apenas foi colocado face-down.

Monstros Summoned através de tributos devem obrigatoriamente entrar:

```text
Face-up
```

---

# 20. Revelação de Monstro Face-down

Um Monstro face-down pode ser revelado pelo próprio controlador durante sua Phase 1.

Ao revelar o Monstro, o jogador pode definir sua posição de acordo com as regras de posição.

Também existe uma revelação automática durante o combate.

Quando um Monstro face-down é atacado:

```text
Face-down Defense
       ↓
Ataque declarado
       ↓
Monstro é revelado
       ↓
Combate é resolvido
```

Se o Monstro sobreviver:

```text
Turno do adversário termina
       ↓
Próximo turno do controlador
       ↓
Phase 1
       ↓
Pode alterar sua posição
```

O Monstro não pode ser utilizado para atacar durante o turno do adversário apenas porque foi revelado.

---

# 21. Posição dos Monstros

As posições disponíveis são:

```text
Attack
Defense
```

Um Monstro em Defense não pode atacar.

## 21.1 Mudança de posição

O jogador pode alterar a posição de um Monstro durante sua Phase 1.

Exemplo:

```text
Defense
↓
Phase 1
↓
Attack
↓
Phase 2
↓
Pode atacar
```

Um Monstro que atacou naquele turno não pode normalmente mudar sua posição naquele mesmo turno.

A alteração ocorre no turno seguinte.

---

# 22. Double Attack

`Double Attack` é um efeito passivo de Monstro.

Um Monstro com Double Attack pode realizar:

```text
2 ataques por turno
```

Os ataques podem ser realizados contra:

* Monstros diferentes;
* o mesmo Monstro, caso ele sobreviva ao primeiro ataque;
* diretamente contra o adversário quando não houver Monstros no campo adversário.

O Double Attack **não permite** mudar para Defense depois de atacar.

Exemplo:

```text
Phase 2

Attack
↓
Attack

Monster continua em Attack
↓
Próximo turno
↓
Pode mudar de posição
```

---

# 23. Phase 2

Durante a Phase 2 o jogador pode:

* atacar com Monstros elegíveis;
* ativar Magias que atendam às condições;
* executar outras ações permitidas;
* resolver efeitos automáticos;
* participar de Chains quando aplicável.

O jogador pode atacar com todos os Monstros elegíveis.

Exemplo:

```text
Monster A → 1 ataque
Monster B → 1 ataque
Monster C → Double Attack → 2 ataques

Total:
4 ataques
```

---

# 24. Ataques

Normalmente, um Monstro pode atacar:

```text
1 vez por turno
```

Monstros com Double Attack podem atacar duas vezes.

Um Monstro em Defense não pode atacar.

---

# 25. Ataque direto

Um ataque direto é permitido quando o adversário não possui Monstros no campo.

A existência de uma Trap no campo adversário, por si só, **não impede um ataque direto**.

Uma Trap pode interferir no ataque apenas se sua própria condição de ativação for satisfeita e seu efeito modificar o combate.

O dano de um ataque direto é igual ao ATK do Monstro.

Exemplo:

```text
Monster:
2500 ATK

Ataque direto:
2500 LP de dano
```

---

# 26. Combate — Attack vs Attack

Quando um Monstro em Attack ataca outro Monstro em Attack, seus valores de ATK são comparados.

## 26.1 Atacante possui ATK maior

Exemplo:

```text
Atacante: 2500 ATK
Defensor: 1800 ATK
```

Resultado:

```text
Defensor destruído.

2500 - 1800 = 700

Jogador defensor:
-700 LP
```

O Monstro defensor vai para o Graveyard.

---

## 26.2 ATK igual

Exemplo:

```text
Atacante: 2000 ATK
Defensor: 2000 ATK
```

Resultado:

```text
Nenhum Monstro é destruído.

Dano:
0
```

---

## 26.3 Atacante possui ATK menor

Exemplo:

```text
Atacante: 1500 ATK
Defensor: 2000 ATK
```

Resultado:

```text
Nenhum Monstro é destruído.

2000 - 1500 = 500

Controlador do atacante:
-500 LP
```

O Monstro atacante permanece no campo.

---

# 27. Combate — Attack vs Defense

Quando um Monstro em Attack ataca um Monstro em Defense:

```text
ATK do atacante
vs
DEF do defensor
```

---

## 27.1 ATK maior que DEF

Exemplo:

```text
Atacante:
3000 ATK

Defensor:
1800 DEF
```

Resultado:

```text
Defensor destruído.

Dano ao jogador:
0
```

A exceção é quando o atacante possui `Piercing Damage`.

---

## 27.2 ATK menor que DEF

Exemplo:

```text
Atacante:
1500 ATK

Defensor:
2000 DEF
```

Resultado:

```text
Nenhum Monstro é destruído.

2000 - 1500 = 500

Controlador do atacante:
-500 LP
```

---

## 27.3 ATK igual à DEF

Exemplo:

```text
Atacante:
2000 ATK

Defensor:
2000 DEF
```

Resultado:

```text
Nenhum Monstro é destruído.

Dano:
0
```

---

# 28. Piercing Damage

`Piercing Damage` é um efeito de Monstro.

Quando um Monstro com Piercing possui ATK maior que a DEF de um Monstro defensor:

```text
Defensor é destruído
```

e a diferença entre ATK e DEF é causada como dano ao jogador defensor.

Exemplo:

```text
Atacante:
3000 ATK

Defensor:
1800 DEF

3000 - 1800 = 1200
```

Resultado:

```text
Defensor → Graveyard

Jogador defensor:
-1200 LP
```

Se:

```text
ATK < DEF
```

a regra normal permanece:

```text
Atacante permanece no campo
Jogador atacante recebe a diferença
```

Se:

```text
ATK = DEF
```

não existe destruição nem dano.

---

# 29. Múltiplos ataques

Um jogador pode atacar com todos os Monstros elegíveis durante sua Phase 2.

Exemplo:

```text
Monster A → ataca Monster X
Monster B → ataca Monster Y
Monster C → ataca diretamente
```

A situação do campo é avaliada após cada ataque.

Exemplo:

```text
Campo adversário:

Monster X
Monster Y

↓

Monster A destrói X

↓

Monster B destrói Y

↓

Não existem mais Monstros

↓

Monster C pode realizar ataque direto
```

Se um Monstro com Double Attack possuir dois ataques disponíveis, ele pode utilizar o segundo ataque após a resolução do primeiro, desde que continue elegível.

---

# 30. Fim imediato do duelo

Sempre que uma ação causar:

```text
LP <= 0
```

o duelo termina imediatamente.

Nenhum ataque ou efeito posterior é executado.

---

# 31. Graveyard

O Graveyard é uma zona permanente do duelo.

Cartas enviadas ao Graveyard incluem:

* Monstros destruídos;
* Monstros utilizados como tributo;
* Magias utilizadas;
* Traps utilizadas;
* cartas descartadas;
* Field Magics substituídas.

O Graveyard é:

```text
Público
```

Ambos os jogadores podem visualizar seu conteúdo.

## 31.1 Ordem

A ordem das cartas no Graveyard corresponde à ordem em que chegaram.

Entretanto, quando múltiplas cartas são destruídas simultaneamente, a ordem entre essas cartas **não possui significado de gameplay**.

Exemplo:

```text
Monster A destruído
Monster B destruído

A → Graveyard
B → Graveyard
```

Não existe diferença de resultado caso a ordem seja invertida.

---

# 32. Posição no Graveyard

Quando uma carta vai para o Graveyard, sua posição anterior deixa de existir.

Exemplo:

```text
Defense
↓
Destroyed
↓
Graveyard
```

Se posteriormente for revivida, a posição anterior não é restaurada.

O jogador escolhe a nova posição conforme as regras de Revival.

---

# 33. Revival

Monstros podem ser revividos por efeitos de Magic.

Uma Magic pode possuir diferentes condições de Revival.

Exemplos:

```text
Revive qualquer Monstro
```

ou:

```text
Revive um Monstro específico
```

ou:

```text
Revive um Monstro de determinado Weight
```

O Revival só pode trazer para o campo uma carta pertencente ao próprio jogador.

Não é possível reviver um Monstro pertencente ao adversário.

---

# 34. Regras do Monstro Revivido

Quando um Monstro é revivido:

* o jogador escolhe sua posição;
* ele não pode atacar naquele mesmo turno;
* ele não pode ser utilizado como tributo naquele mesmo turno.

Se for destruído novamente:

```text
Monstro
↓
Graveyard
```

Não existe limite de quantidade de vezes que uma carta pode ser revivida.

---

# 35. Magias

Toda Magic precisa ser colocada no campo antes de sua ativação.

A Magic é colocada em um dos espaços disponíveis da Zona de Magias/Armadilhas.

## 35.1 Ativação

Magics são ativadas pelo jogador quando suas condições forem atendidas.

Uma Magic pode possuir condições específicas de ativação.

Algumas Magics podem ser ativadas durante o turno do próprio jogador.

Determinadas Magics podem possuir condições que permitam sua ativação durante o turno do adversário.

Essa possibilidade depende da própria carta.

---

## 35.2 Magic Equipment

Magic Equipment é um tipo de Magic que concede um bônus a um Monstro compatível do campo.

A compatibilidade é definida pelo **subtipo** do Monstro (seção 10.1.1). Cada Equipment aceita somente Monstros do subtipo correspondente.

> **Story Mode MVP:**
>
> * A Magic Equipment precisa primeiro ser colocada em um espaço livre da Zona de Magias/Armadilhas, pela ação da seção 13.3.
> * A ativação só é permitida na **Phase 2**.
> * Não é possível ativar uma Magic diretamente da mão.
> * Colocar um Monstro e ativar uma Equipment não podem ser combinados em uma única ação.
> * Se os quatro espaços da Zona de Magias/Armadilhas estiverem ocupados, não é possível colocar outra Magic. Substituir uma carta existente não é permitido.
> * O MVP possui somente as quatro Magic Equipment do catálogo abaixo. Burn e Power Up não fazem parte do MVP.

### 35.2.1 Efeito

Cada Magic Equipment do MVP concede **+800 ATK** ao Monstro equipado, desde que ele seja de seu subtipo compatível.

Catálogo do MVP:

| Carta | Efeito |
| --- | --- |
| Dragon Equipment | +800 ATK a Monstros Dragon |
| Warrior Equipment | +800 ATK a Monstros Warrior |
| Mage Equipment | +800 ATK a Monstros Mage |
| Aquatic Equipment | +800 ATK a Monstros Aquatic |

O bônus segue a regra de duração da seção 38: permanece enquanto a carta permanecer equipada.

### 35.2.2 Relação com o Monstro equipado

A Magic Equipment mantém uma relação com o Monstro ao qual está equipada.

Um Monstro pode possuir várias Magic Equipment compatíveis equipadas ao mesmo tempo, e os bônus se acumulam.

### 35.2.3 Destino

Quando um Monstro equipado deixa o campo, todas as Magic Equipment equipadas a ele vão para o Graveyard.

Exemplo:

```text
Monstro Dragon + Dragon Equipment + Dragon Equipment
↓
Monstro deixa o campo
↓
Monstro → Graveyard
Dragon Equipment → Graveyard
Dragon Equipment → Graveyard
```

Esta regra não introduz nenhum comportamento de ativação de Trap.

### 35.2.4 Ativação e vínculo

* A Equipment deve estar colocada em um espaço da Zona de Magias/Armadilhas antes de ser ativada.
* A ativação ocorre na Phase 2.
* O jogador precisa selecionar um Monstro compatível que já esteja no campo dele e esteja **face-up**.
* Se não existir um Monstro compatível face-up, a ativação é rejeitada.
* Após a ativação, a Equipment permanece no campo, vinculada ao Monstro escolhido, e continua ocupando seu espaço da Zona de Magias/Armadilhas.
* Várias Equipment podem estar vinculadas ao mesmo Monstro.
* Se o Monstro equipado sair do campo, todas as Equipment vinculadas a ele vão para o Graveyard (seção 35.2.3).

---

# 36. Tipos de efeitos de Magic

Os efeitos são classificados por comportamento, e não necessariamente por uma carta específica.

Exemplos:

```text
ADD_TO_HAND
DESTROY
RECOVER_LP
BUFF_ATK
DEBUFF_ATK
FIELD_TYPE
```

Podem existir várias cartas utilizando o mesmo tipo de efeito.

Exemplo:

```text
Magic A:
DESTROY
Condição:
Destrói Monstro com ATK <= 2000

Magic B:
DESTROY
Condição:
Destrói qualquer Monstro
```

Ambas possuem:

```text
EffectType = DESTROY
```

mas possuem regras diferentes.

---

# 37. Magic de campo — Field

Uma Magic do tipo `FIELD` permanece ativa no campo.

Ela não vai imediatamente para o Graveyard após sua ativação.

Enquanto estiver ativa, seu efeito permanece válido conforme as condições da carta.

## 37.1 Substituição

Só pode existir uma Field Magic ativa por vez.

Quando uma nova Field Magic é ativada:

```text
Field Magic antiga
        ↓
Graveyard

Field Magic nova
        ↓
Campo
```

A nova Field substitui a anterior.

Exemplo:

```text
FIELD A ativa
↓
FIELD B ativa
↓
FIELD A → Graveyard
↓
FIELD B permanece ativa
```

---

# 38. Duração dos efeitos

Os efeitos que modificam atributos normalmente permanecem ativos durante todo o duelo.

Exemplo:

```text
+500 ATK
```

O bônus permanece enquanto a regra da carta não determinar sua remoção.

Entretanto, uma carta pode especificar uma duração diferente.

Exemplo:

```text
+1000 ATK
até o final do turno
```

Nesse caso, o modificador é removido ao final do turno.

A duração do efeito é definida pela própria carta quando aplicável.

---

# 39. Efeitos de Monstros

Efeitos de Monstros são **automáticos**.

O jogador não escolhe se o efeito será ativado quando sua condição for satisfeita.

Exemplo:

```text
"Quando este Monstro for Summoned,
ganhe +300 ATK."
```

Ao realizar o Summon:

```text
Summon
↓
Condição satisfeita
↓
Efeito executado automaticamente
```

---

# 40. Repetição de efeitos de Monstros

Um efeito automático pode ocorrer múltiplas vezes se sua condição ocorrer múltiplas vezes.

Exemplo:

```text
Efeito:
"Ao atacar, ganhe +100 ATK."
```

Se o Monstro atacar duas vezes:

```text
Primeiro ataque
↓
+100 ATK

Segundo ataque
↓
+100 ATK
```

Resultado:

```text
+200 ATK
```

O efeito é avaliado novamente a cada ocorrência da condição.

---

# 41. Efeitos opcionais

Efeitos de Monstros não possuem escolha manual de ativação.

Quando a condição do efeito é satisfeita:

```text
Condição satisfeita
↓
Efeito automático
```

Traps possuem comportamento diferente conforme definido na seção de Armadilhas.

---

# 42. Armadilhas — Traps

Traps precisam ser colocadas no campo antes de poderem atuar.

As Traps possuem condições/eventos que determinam quando podem ser utilizadas.

Exemplo:

```text
Condição:
Oponente declara um ataque.

↓
Condição satisfeita

↓
Trap pode atuar.
```

Traps podem atuar durante o turno do próprio jogador ou do adversário, conforme sua condição.

---

# 43. Ativação automática de Traps

Traps são diferentes das Magics.

Quando o evento da Trap acontece:

```text
Evento
↓
Condição da Trap satisfeita
↓
Trap é ativada
```

O jogador não precisa iniciar manualmente a ativação da Trap.

Quando existir uma escolha associada à resolução da Trap, o jogador possui uma janela de:

```text
5 segundos
```

para tomar a decisão necessária.

---

# 44. Trap respondendo a Magic ou Trap

Traps podem responder a outras cartas.

Exemplo:

```text
Magic A
"Ganhe +500 ATK"

↓

Trap B
"Cancele o efeito"

↓

Chain
```

A Trap pode responder à Magic porque sua condição permite.

---

# 45. Chain

O jogo possui um sistema de Chain para representar respostas entre cartas.

Uma Chain é formada quando uma carta responde à ativação de outra.

Exemplo:

```text
Magic A
    ↓
Trap B responde
    ↓
Trap C responde
```

Quando não houver novas ativações:

```text
Chain termina de receber novas ativações
```

A resolução ocorre na ordem inversa da ativação.

```text
C resolve primeiro
↓
B resolve
↓
A resolve
```

Esse comportamento é:

```text
Last In, First Out
(LIFO)
```

---

# 46. Resolução da Chain

Durante a resolução da Chain:

**nenhuma nova ativação pode ocorrer.**

Exemplo:

```text
Magic A
↓
Trap B
↓
Trap C

Fim das ativações

↓

C resolve
↓

B resolve
↓

A resolve
```

Uma resolução não pode automaticamente iniciar outra ativação.

---

# 47. Efeitos não geram novas ativações

Um efeito resolvido não pode automaticamente criar uma nova ativação de outro efeito.

Exemplo:

```text
Magic A resolve
↓
modifica ATK
↓
não inicia automaticamente outra Magic/Trap
```

Uma nova ativação só ocorre através de uma nova oportunidade válida de ativação.

---

# 48. Prioridade das regras

As regras gerais do jogo representam o comportamento padrão.

Uma carta pode criar uma exceção explícita.

Exemplo:

```text
Regra normal:

Monstro que atacou
não pode mudar de posição naquele turno.
```

Uma carta pode possuir um efeito que explicitamente modifique essa regra.

Portanto:

> **Efeitos de cartas podem modificar regras gerais quando essa alteração estiver explicitamente definida pela carta.**

---

# 49. Stop Defense e Set Defense

Existem Magics capazes de alterar a posição dos Monstros.

## 49.1 Stop Defense

Uma Magic do tipo `Stop Defense` pode alterar a posição de Monstros que estejam em Defense conforme as condições da carta.

## 49.2 Set Defense

Uma Magic do tipo `Set Defense` possui comportamento oposto, podendo alterar Monstros em Attack para Defense conforme as condições da carta.

A Magic determina quais Monstros são afetados.

---

# 50. Efeitos de alteração de posição

Efeitos de Magic podem modificar a posição dos Monstros.

Esses efeitos podem ocorrer independentemente da regra normal de mudança de posição.

Exemplo:

```text
Monster:
Attack

Magic:
Set Defense

↓

Monster:
Defense
```

---

# 51. Surrender

No modo Online, o jogador pode escolher:

```text
Surrender
```

Ao realizar Surrender:

```text
Jogador que desistiu → derrota
Adversário → vitória
```

---

# 52. Abandono da partida

Se o jogador fechar o navegador ou abandonar a partida durante um duelo Online:

```text
Abandono
↓
Derrota imediata
```

O adversário vence.

---

# 53. Desconexão

Uma desconexão durante uma partida Online resulta em:

```text
Desconexão
↓
Derrota imediata
```

O adversário vence.

---

# 54. Timeout

Cada jogador possui:

```text
40 segundos por turno
```

Quando o tempo chega a zero:

```text
Turno termina imediatamente
↓
Vez do adversário
```

O jogador **não perde o duelo** simplesmente por atingir o timeout.

---

# 55. Progressão

O Story Mode possui:

```text
10 níveis
```

Cada nível possui:

```text
4 duelos
```

Total:

```text
10 × 4 = 40 duelos
```

---

# 56. Progressão dentro do nível

Para avançar para o próximo nível, o jogador precisa vencer todos os quatro duelos do nível atual.

Exemplo:

```text
Level 1

Duel 1 → Vitória
Duel 2 → Vitória
Duel 3 → Vitória
Duel 4 → Vitória

↓
Level 2
```

Se perder qualquer duelo:

```text
Progressão não avança.
```

---

# 57. Dificuldade

A dificuldade aumenta dentro de cada nível.

Os quatro duelos possuem dificuldade progressiva.

Os adversários possuem Decks progressivamente mais fortes.

Conforme o jogador avança:

* aparecem cartas mais fortes;
* aparecem novas estratégias;
* os adversários possuem combinações mais complexas;
* novas cartas podem ser desbloqueadas.

---

# 58. Recompensas

Ao vencer um duelo do Story Mode, o jogador recebe:

```text
3 cartas
```

Cartas também poderão ser obtidas posteriormente através de um sistema de Card Shop.

O objetivo é permitir que o jogador fortaleça e modifique seu Deck ao longo da progressão.

---

# 59. Balanceamento

O jogo deve incentivar diferentes estratégias.

Não deve existir uma única carta ou combinação obrigatória que determine automaticamente o resultado do duelo.

O balanceamento deve considerar:

* Attack;
* Defense;
* Weight;
* quantidade de tributos;
* efeitos;
* condições;
* duração dos efeitos;
* quantidade máxima de cópias;
* frequência de compra;
* interação entre cartas;
* progressão do jogador;
* construção do Deck.

O objetivo é permitir que diferentes combinações de cartas sejam competitivas.

---

# 60. Princípio fundamental de Game Design

O jogo deve priorizar:

```text
Estratégia
+
Construção de Deck
+
Combinação de cartas
+
Tomada de decisão
```

em vez de depender exclusivamente de:

```text
"Quem possui a carta mais forte vence."
```

Cartas poderosas devem possuir condições, custos, limitações ou oportunidades de interação que permitam ao adversário responder.

---

# 61. Story Mode e Multiplayer

O jogo possui dois modos principais:

```text
Story Mode
Online Multiplayer
```

Ambos utilizam o mesmo núcleo de regras de duelo.

## 61.1 Story Mode

No Story Mode:

```text
Duel Engine
    ↓
Jogador
    +
AI Opponent
```

A IA controla o adversário.

## 61.2 Multiplayer

No Multiplayer:

```text
Duel Engine
    ↓
Jogador A
    +
Jogador B
```

O Duel Engine não deve depender do tipo de oponente.

A diferença entre:

```text
AI
```

e:

```text
Human
```

pertence à camada responsável pelo controle do jogador, não às regras do duelo.

---

# 62. Sistema de Stars e apostas

O jogo poderá possuir uma moeda interna chamada:

```text
Stars
```

Exemplo:

```text
Jogador:
40 Stars
```

## 62.1 Aposta mínima

A aposta mínima é:

```text
10 Stars
```

Para iniciar um duelo com aposta:

```text
Jogador A aposta 10
+
Jogador B aceita 10
=
Pot de 20 Stars
```

Ao final:

```text
Vencedor → recebe 20 Stars
Perdedor → perde 10 Stars
```

## 62.2 Recusa

Se o jogador desafiado recusar a aposta:

```text
Duelo com aposta não começa.
```

O jogador pode realizar uma nova tentativa sem aposta.

## 62.3 Transferência

Stars não podem ser transferidas diretamente entre jogadores.

A transferência de Stars entre jogadores ocorre somente através do sistema de aposta.

## 62.4 Economia com dinheiro real

A possibilidade futura de comprar Stars com dinheiro real pertence a uma camada de economia separada.

Exemplo conceitual:

```text
20 Stars = R$5
```

Essa funcionalidade não faz parte do núcleo do Duel Engine.

---

# 63. Regras de informação

As informações das cartas possuem diferentes níveis de visibilidade.

## 63.1 Cartas na mão

A mão é privada.

O adversário não pode visualizar as cartas que o jogador possui na mão.

## 63.2 Cartas face-up

Cartas face-up possuem suas informações visíveis ao adversário.

## 63.3 Cartas face-down

Todas as informações são ocultadas do adversário.

Somente o controlador conhece os detalhes da carta.

## 63.4 Graveyard

O Graveyard é público.

Ambos os jogadores podem visualizar as cartas presentes nele.

---

# 64. Estado de uma carta durante o duelo

Uma carta possui uma definição permanente e um estado específico durante um duelo.

O estado pode incluir informações como:

```text
Owner
Location
Face-up / Face-down
Position
Current ATK
Current DEF
Attacks performed
```

O estado da carta pode mudar durante o duelo.

Exemplo:

```text
Deck
 ↓
Hand
 ↓
Field
 ↓
Graveyard
```

A mesma instância de carta permanece sendo a mesma carta durante sua movimentação pelo duelo.

---

# 65. Regras não implícitas

A implementação não deve criar regras por conta própria.

Quando uma carta ou funcionalidade exigir uma regra que não esteja definida neste documento:

1. a regra ausente deve ser identificada;
2. nenhum comportamento arbitrário deve ser criado;
3. a decisão deve ser tomada no nível de Game Design;
4. este documento deve ser atualizado;
5. somente depois a implementação definitiva deve ser realizada.

---

# 66. Fonte de verdade

Este documento representa a **fonte de verdade das regras do jogo**.

Qualquer implementação deve respeitar as regras descritas aqui.

A implementação técnica pode utilizar:

```text
Go
React
TypeScript
PostgreSQL
Docker
WebSocket
REST
```

ou qualquer outra tecnologia.

Entretanto, essas tecnologias **não definem as regras do jogo**.

As regras pertencem exclusivamente ao domínio definido neste documento.

---

# 67. Estado atual

**Game Rules v1.0**

O núcleo do jogo está definido:

* objetivo e vitória;
* Life Points;
* Deck;
* limite de cópias;
* mão;
* compra;
* Deck vazio;
* campo;
* tipos de cartas;
* Turno;
* Discard;
* Phase 1;
* Summon;
* Weight;
* Tributos;
* Face-up/Face-down;
* Attack/Defense;
* mudança de posição;
* Double Attack;
* ataques;
* ataque direto;
* combate;
* Piercing;
* Graveyard;
* Revival;
* Magics;
* Field Magic;
* Traps;
* efeitos automáticos;
* duração de efeitos;
* Chain;
* resolução LIFO;
* timeout;
* surrender;
* abandono;
* progressão;
* Story Mode;
* Multiplayer;
* Stars e apostas como sistema futuro.

Este documento deve ser tratado como a referência de Game Design para a implementação do domínio.

O próximo passo é transformar essas regras em um **modelo de domínio implementável**, mantendo a separação entre:

```text
Game Rules
     ↓
Domain Model
     ↓
Application Layer
     ↓
Infrastructure
     ↓
API
```

A implementação não deve antecipar regras que ainda não tenham sido definidas pelo Game Design.
