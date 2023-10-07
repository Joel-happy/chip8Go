# CHIP8Go
Pour le compte d'Ymmersion, nous avons réalisé un projet d'émulateur CHIP-8 en Golang.

## Description
L'émulateur CHIP8Go est une implémentation de l'émulateur CHIP-8, une machine virtuelle de jeu rétro.
CHIP-8 a été initialement développé dans les années 1970 et est connu pour sa simplicité et sa facilité de programmation, en en faisant un choix populaire pour les développeurs de jeux amateurs et les amateurs de rétro-gaming.

## Comment utiliser notre émulateur?
Il vous suffit d'entrer la commande suivante dans le terminal : "./main.exe -rom (rom)/(nom de la rom que vous souhaitez charger en mémoire)"
Par exemple, la commande "./main.exe -rom .\rom\1-chip8-logo.ch8" lancera la ROM "1-chip8-logo.ch8" contenue dans le répertoire "rom".


## Caractéristiques de CHIP8Go 
  
-Émulateur CHIP-8 en Golang.
-Basé sur la machine virtuelle de jeu rétro CHIP-8.
-Supporte les ROMs CHIP-8.
-Chargement et exécution de ROMs depuis la ligne de commande.
-Prise en charge des jeux rétro et de la programmation amateur.
-Interface utilisateur graphique basée sur la bibliothèque Ebiten pour afficher les graphismes.
-Gestion des entrées clavier pour émuler les entrées de la CHIP-8.
-Émulation des caractères de la CHIP-8 à l'aide d'un ensemble de polices intégré.
-Gestion des timers de délai et de son de la CHIP-8.
-Capacité à réinitialiser l'état de la machine CHIP-8.
-Exécution des instructions de la CHIP-8 et interprétation des opcodes.
-Affichage de l'écran de jeu CHIP-8 avec possibilité de dessiner des pixels.
-Émulation de la fonction sonore de la CHIP-8 avec lecture audio MP3.


## Répartition des tâches
**Zineb :**

- Fontset
- Cpu
- Setupkeys
- README


**Joël :**

- Main
- Cpu
- README
