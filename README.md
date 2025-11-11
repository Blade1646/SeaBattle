# Sea Battle: Windows console game for up to 2 players

This console application, written in **Go**, reimagines the classic Sea Battle under the name Admiral. It allows players to custom-place their ships on the field, attack enemy vessels, and move their own ships to safety during the fight — adding a new layer of tactics and control to the familiar gameplay.

## Features

- Play in one of two modes — compete with another player on the same device or challenge a computer opponent
- Face an adaptive AI that studies the battlefield, tracks its previous shots, and refines its targeting as the match progresses
- Define your fleet size (from 4 to 10 ships) before each game — the arena is populated automatically, after which you can freely reposition and rotate your ships to prepare your formation
- Command your fleet turn by turn — maneuver and fire strategically, with each ship’s condition naturally affecting its available actions
- Follow every exchange with a clear visual tracker that marks hits, misses, and sunk ships on both grids, helping you plan ahead without losing track of the battle
- Pause the action at any time and continue later using a save system that preserves the full game state, including ship layouts and turn order

## How to Play

### Installation

1. Download the repository and place it in a separate folder on your system.

2. Make sure you have the latest version of Go installed. Initialize the module by running `go mod init <module_name>` inside the project folder.

3. If not installed automatically, fetch the required dependency, e.g., `go get github.com/gdamore/tcell/v2`.

4. To start the game, run the following command in Windows Terminal (recommended for proper character rendering): `wt -w 0 nt go run .`.

### Gameplay

1. **Setup the match**
    
    - Choose the game mode (Player vs Player or Player vs AI).
    - Select whether to start a new game or load a saved one.
    - If starting a new game, choose the number of ships (from 4 to 10).
        
2. **Arrange your ships**
    
    - Use `W`, `A`, `S`, `D` keys to move and `Q`, `E` to rotate your ship.
    - For the second player, use the numeric keypad: `8`, `5`, `4`, `6` to move and `7`, `9` to rotate.
    - When all ships are positioned, press `Space` to confirm.
        
3. **Play the battle**
    
    - On your turn, select a ship to act with.
    - A healthy ship can perform one movement (optional) and one shot per turn.
    - To fire, select a target cell. If you hit, you continue shooting; if you miss, the turn passes to the other player.
    - A damaged ship can either move or shoot — not both.
    - A sunk ship can no longer act and remains on the field as wreckage.
        
4.  **Winning the game**
    
    - The match continues until one side loses all ships.
    - The player who sinks the entire opposing fleet wins the battle.
        
5. **Exiting and saving**
    
    - After the game ends, press `Space` to return to the main menu.
    - To exit mid-game, press `Esc` and confirm with `Enter`.
    - The game state (map layout, ship status, turn order) is saved automatically and can be resumed later.

## Preview

The full video demonstration is available on [YouTube](https://youtu.be/iK2a3EZqNuY).
