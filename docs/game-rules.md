# Game Rules — TCG

**Version:** 0.4
**Status:** Em definição

Este documento define as regras de negócio e funcionamento do jogo de cartas.

Ele representa **o que o jogo deve fazer**, independentemente da tecnologia utilizada para implementá-lo.

> **Importante:** este documento é a fonte de verdade das regras do jogo.
> A implementação técnica não deve criar ou alterar regras que não estejam definidas aqui.

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
```

Jogador A perde o duelo.

---

# 3. Life Points

Cada jogador começa o duelo com uma quantidade de Life Points determinada pelo nível/progressão do jogo.

### Regra inicial

```text
Life Points iniciais: 4000
```

A progressão poderá aumentar esse valor posteriormente.

Exemplo planejado:

```text
4000 → ... → 8000 LP
```

A definição de em quais níveis o valor aumenta ainda está pendente.

---

# 4. Deck

Cada jogador possui um Deck com:

```text
30 cartas
```

O Deck pode conter diferentes tipos de cartas, incluindo:

* Monstros
* Magias
* Armadilhas

A quantidade máxima de cópias da mesma carta dentro do Deck ainda não foi definida.

### Regras pendentes

* Quantidade máxima de cópias por carta.
* Regras para construção do Deck.
* Cartas obrigatórias/opcionais.
* Regras de Deck inválido.

---

# 5. Mão inicial

No início do duelo, o jogador recebe:

```text
3 cartas
```

No primeiro turno, o jogador compra 1 carta.

Assim:

```text
Mão inicial
3 cartas

↓ compra 1

4 cartas
```

A definição de quem começa o duelo ainda precisa ser especificada.

---

# 6. Compra e descarte

A partir do segundo turno, o jogador:

1. Compra 1 carta.
2. Realiza suas ações.
3. Descarta 1 carta.

Exemplo:

```text
Início do turno:
4 cartas

Compra:
5 cartas

Final do turno:
descarta 1

Resultado:
4 cartas
```

A regra exata de descarte ainda precisa ser definida.

### Questões pendentes

* O jogador pode escolher qualquer carta?
* Existem cartas que não podem ser descartadas?
* O descarte acontece obrigatoriamente no final do turno?
* O que acontece se o jogador não possuir cartas suficientes?

---

# 7. Campo

Cada jogador possui dois espaços principais no campo.

## 7.1 Zona de Monstros

Cada jogador pode possuir no máximo:

```text
4 Monstros
```

---

## 7.2 Zona de Magias/Armadilhas

Cada jogador pode possuir no máximo:

```text
4 cartas de Magia/Armadilha
```

Portanto, o campo máximo de cada jogador é:

```text
4 Monstros
+
4 Magias/Armadilhas
```

---

# 8. Turno

Cada turno possui fases.

## Primeiro turno

No primeiro turno:

1. Compra 1 carta.
2. O jogador pode realizar suas ações normalmente.
3. O jogador **não pode atacar**.

---

## Turnos seguintes

Nos turnos seguintes:

1. Compra 1 carta.
2. Descarta 1 carta.
3. Executa a **Phase 1**.
4. Executa a **Phase 2**.

Fluxo:

```text
Draw
 ↓
Discard
 ↓
Phase 1
 ↓
Phase 2
 ↓
Fim do turno
```

---

# 9. Phase 1

A Phase 1 é responsável principalmente pela colocação de cartas no campo.

Durante cada Phase 1 o jogador possui:

```text
1 colocação
```

O jogador pode escolher **uma** das opções:

### Opção A — Colocar Monstro

Invocar ou colocar um Monstro no campo.

### Opção B — Colocar Magia/Armadilha

Colocar uma carta de Magia ou Armadilha no campo.

---

## Regra importante

O jogador **não pode realizar duas colocações** na mesma Phase 1.

Exemplo inválido:

```text
Invocar Monstro
+
Colocar Magia
```

na mesma Phase 1.

Também não é permitido:

```text
Invocar Monstro A
+
Invocar Monstro B
```

na mesma Phase 1.

---

# 10. Invocação de Monstros

Cada Monstro possui um atributo chamado:

```text
Weight
```

O Weight determina quantos tributos são necessários para invocá-lo.

---

# 11. Weight

Os monstros possuem Weight de:

```text
1+
```

A quantidade de tributos é determinada pela seguinte tabela:

| Weight | Tributos necessários |
| ------ | -------------------: |
| 1–5    |                    0 |
| 6–8    |                    1 |
| 9–12   |                    2 |
| 13+    |                    3 |

O número de tributos é determinado pela faixa do Weight.

O Weight dos monstros utilizados como tributo **não precisa ser igual ao Weight do monstro sendo invocado**.

---

# 12. Tributos

Quando um Monstro exige tributos, os tributos devem ser monstros que já estejam no campo.

Um monstro que ainda esteja na mão não pode ser utilizado como tributo.

Exemplo:

```text
Monstro alvo:
Weight 10

Tributos necessários:
2
```

O jogador pode utilizar:

```text
Weight 3
+
Weight 8
=
Weight 11
```

Isso é válido.

O valor combinado dos tributos pode ser maior que o Weight do monstro invocado.

---

## 12.1 Preparação de tributos

Como a Phase 1 permite apenas uma colocação, o jogador não pode realizar uma sequência de invocações para preparar tributos na mesma Phase 1.

Exemplo:

O jogador possui:

```text
Mão:
Monstro A — Weight 3
Monstro B — Weight 8
Monstro C — Weight 10
```

Campo vazio.

O jogador não pode:

```text
Invocar A
Invocar B
Invocar C usando A+B
```

na mesma Phase 1.

Os monstros utilizados como tributo precisam estar previamente no campo.

---

# 13. Posição do Monstro

Todo Monstro possui uma posição.

As posições disponíveis são:

```text
Attack
Defense
```

Ao colocar um Monstro no campo, o jogador escolhe sua posição.

---

# 14. Atributos do Monstro

Um Monstro pode possuir:

```text
Weight
Attack
Defense
Position
Effects
```

Exemplo conceitual:

```text
Monster:
    Weight: 7
    Attack: 2500
    Defense: 1800
    Position: Attack
    Effects:
        Piercing Damage
```

---

# 15. Mudança de posição

Um Monstro pode alterar sua posição durante o jogo.

### Regra padrão

Um Monstro só pode mudar de posição se:

```text
não tiver atacado naquele turno
```

Portanto, normalmente:

```text
Attack
↓
ataca
↓
não pode mudar para Defense
```

---

# 16. Exceção: Double Attack

Um Monstro com o efeito:

```text
Double Attack
```

possui uma exceção à regra de mudança de posição.

Esse Monstro pode:

```text
Atacar
↓
Atacar novamente
↓
Mudar para Defense
```

mesmo tendo atacado naquele turno.

Essa exceção existe porque o efeito do Monstro modifica a regra normal de mudança de posição.

---

# 17. Phase 2

Durante a Phase 2 o jogador pode realizar ações relacionadas ao combate.

Entre elas:

* Atacar com Monstros.
* Ativar Magias que atendam às suas condições.
* Ativar Armadilhas quando suas condições forem atendidas.
* Executar efeitos permitidos pelas regras das cartas.

---

# 18. Ataques

Normalmente, cada Monstro pode atacar:

```text
1 vez por turno
```

Um Monstro com o efeito:

```text
Double Attack
```

pode atacar:

```text
2 vezes por turno
```

---

# 19. Ataque direto

Um jogador pode realizar um ataque direto quando o adversário:

```text
não possui Monstros no campo
```

Nesse caso, o dano causado é igual ao Attack do Monstro atacante.

Exemplo:

```text
Monster Attack: 2500

Dano:
2500 LP
```

---

# 20. Combate — Attack vs Attack

Quando um Monstro em posição de Attack ataca outro Monstro em posição de Attack, os valores de Attack são comparados.

---

## 20.1 Atacante possui Attack maior

Exemplo:

```text
Atacante: 2500 ATK
Defensor: 1800 ATK
```

Resultado:

```text
Defensor destruído.

Diferença:
2500 - 1800 = 700
```

O jogador defensor recebe:

```text
700 de dano
```

---

## 20.2 Attack igual

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

## 20.3 Atacante possui Attack menor

Exemplo:

```text
Atacante: 1500 ATK
Defensor: 2000 ATK
```

Resultado:

```text
Nenhum Monstro é destruído.

Diferença:
2000 - 1500 = 500
```

O jogador do Monstro atacante recebe:

```text
500 de dano
```

O Monstro atacante permanece no campo.

---

# 21. Combate — Attack vs Defense

Quando um Monstro em posição de Attack ataca um Monstro em posição de Defense, o Attack do atacante é comparado com a Defense do defensor.

---

## 21.1 Attack maior que Defense

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
```

Normalmente:

```text
Dano ao jogador:
0
```

A exceção é quando o atacante possui o efeito:

```text
Piercing Damage
```

---

## 21.2 Attack menor que Defense

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
```

Diferença:

```text
2000 - 1500 = 500
```

O jogador do Monstro atacante recebe:

```text
500 de dano
```

---

## 21.3 Attack igual à Defense

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

# 22. Piercing Damage

Um Monstro que possui o efeito:

```text
Piercing Damage
```

pode causar dano ao jogador mesmo quando ataca um Monstro em posição de Defense.

Exemplo:

```text
Atacante:
3000 ATK
Piercing Damage

Defensor:
1800 DEF
```

Resultado:

```text
Defensor destruído.

Diferença:
3000 - 1800 = 1200
```

O jogador defensor recebe:

```text
1200 LP de dano
```

---

# 23. Monstro que perde o combate

Quando um Monstro ataca e possui valores inferiores ao Monstro defensor, ele não é destruído automaticamente.

Exemplo:

```text
Atacante:
1500 ATK

Defensor:
2000 ATK
```

Resultado:

```text
Atacante permanece no campo.
```

O jogador do atacante recebe o dano correspondente à diferença.

O Monstro também:

```text
não pode mudar para Defense naquele mesmo turno
```

pela regra normal de mudança de posição.

---

# 24. Magias

As cartas de Magia podem possuir diferentes tipos de comportamento.

Uma Magia pode ser colocada no campo e posteriormente ativada pelo jogador.

Exemplos de efeitos possíveis:

```text
+500 ATK
Comprar 2 cartas
Destruir um Monstro
Recuperar Life Points
```

Esses exemplos representam possibilidades de efeitos e ainda não constituem a lista definitiva do MVP.

---

# 25. Armadilhas

As Armadilhas são cartas que podem ser colocadas no campo e ativadas quando uma determinada condição ocorre.

Exemplo:

```text
Condição:
Oponente declara um ataque.

Efeito:
Ativar a Armadilha.
```

As condições e efeitos específicos ainda precisam ser definidos.

---

# 26. Efeitos

Cartas podem possuir efeitos capazes de modificar as regras normais do jogo.

Exemplos:

```text
Piercing Damage
Double Attack
```

Os efeitos podem:

* alterar regras de combate;
* alterar quantidade de ataques;
* alterar dano;
* alterar posição;
* modificar Attack/Defense;
* permitir ações adicionais;
* modificar compra de cartas;
* modificar Life Points;
* destruir cartas;
* entre outros.

---

# 27. Prioridade das regras

As regras gerais do jogo são o comportamento padrão.

Um efeito de carta pode criar uma exceção explícita a uma regra.

Exemplo:

Regra normal:

```text
Monstro que atacou não pode mudar para Defense.
```

Efeito:

```text
Double Attack
```

Exceção:

```text
Monstro com Double Attack pode atacar e depois mudar para Defense.
```

Portanto, efeitos devem ser tratados como modificadores controlados das regras do jogo.

---

# 28. Progressão

O jogo terá um sistema de progressão baseado em níveis.

Planejamento inicial:

```text
~10 níveis
~40 duelos
```

Conforme o jogador progride:

* novos Monstros podem ser desbloqueados;
* cartas podem possuir valores maiores;
* novas estratégias podem aparecer;
* Life Points podem aumentar.

A progressão de Life Points poderá seguir aproximadamente:

```text
4000 → 8000
```

A definição exata ainda está pendente.

---

# 29. Balanceamento

O jogo deve ser construído para incentivar diferentes estratégias.

Não deve existir uma única carta ou combinação obrigatória que determine automaticamente o resultado do duelo.

O objetivo é permitir que diferentes combinações de cartas sejam competitivas.

O balanceamento deverá considerar, entre outros fatores:

* Attack;
* Defense;
* Weight;
* quantidade de tributos;
* efeitos;
* custo/oportunidade;
* quantidade de cópias;
* frequência de compra;
* interação entre cartas;
* progressão do jogador.

---

# 30. Regras ainda não definidas

As seguintes regras ainda precisam ser especificadas antes de considerar o Game Rules como fechado:

### Deck

* Quantidade máxima de cópias da mesma carta.
* Regras de construção do Deck.

### Compra

* O que acontece quando o Deck acaba.
* O jogador perde ao tentar comprar sem cartas?
* Existe alguma condição de empate?

### Mão

* Limite máximo de cartas na mão.
* Regras exatas do descarte.

### Turno

* Quem começa.
* Duração máxima do turno.
* Regras de timeout.

### Combate

* Regras adicionais de ataque.
* Interações entre múltiplos efeitos.
* O que acontece quando um Monstro é removido durante o combate.

### Posição

* Quantas vezes um Monstro pode mudar de posição.
* Se existe custo para mudar de posição.
* Se efeitos podem forçar mudança de posição.

### Magias

* Tipos definitivos de Magia.
* Lista de efeitos do MVP.
* Quando cada Magia pode ser ativada.

### Armadilhas

* Condições de ativação.
* Quando podem ser ativadas.
* Se podem responder a outras cartas.

### Efeitos

* Sistema de prioridade.
* Ordem de resolução.
* Interação entre efeitos simultâneos.
* Interação entre efeitos que modificam outros efeitos.

### Progressão

* Quantidade exata de níveis.
* Quantidade exata de duelos.
* Recompensas.
* Desbloqueio de cartas.
* Momentos em que os Life Points aumentam.

### Duelo

* Surrender.
* Empate.
* Timeout.
* Reconexão.
* Abandono da partida.

---

# 31. Princípio fundamental do Game Design

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

As cartas mais poderosas devem possuir custos, condições ou limitações que permitam ao adversário interagir com elas.

---

# 32. Fonte de verdade

Este documento representa a **fonte de verdade das regras do jogo**.

Qualquer implementação deve respeitar as regras descritas aqui.

A implementação técnica:

```text
Go
React
PostgreSQL
Docker
etc.
```

não deve definir regras de negócio por conta própria.

Quando uma regra necessária para uma implementação não estiver definida neste documento, a implementação deve:

1. identificar a regra ausente;
2. evitar assumir um comportamento arbitrário;
3. solicitar uma decisão de Game Design;
4. atualizar este documento após a decisão.

---

# 33. Estado atual

**Game Rules v0.4**

O núcleo do combate, sistema de Weight/Tributos, posições, campo, turnos e conceito de efeitos já está definido.

As regras pendentes devem ser resolvidas antes de transformar o domínio completo em código definitivo.

O próximo passo recomendado é modelar o domínio do jogo a partir destas regras.
